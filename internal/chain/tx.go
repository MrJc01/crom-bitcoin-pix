package chain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// TxBuilder constrói, assina e serializa transações Bitcoin.
type TxBuilder struct {
	network   *chaincfg.Params
	feeRate   int64 // sat/vByte
	client    *Client
}

// TxResult contém o resultado de uma transação construída.
type TxResult struct {
	TxID      string
	RawTxHex  string
	Fee       int64
	Size      int
	VSize     int
}

// NewTxBuilder cria um builder de transações.
func NewTxBuilder(network string, feeRate int64) *TxBuilder {
	params := &chaincfg.MainNetParams
	if network == "testnet" {
		params = &chaincfg.TestNet3Params
	} else if network == "regtest" {
		params = &chaincfg.RegressionNetParams
	}

	return &TxBuilder{
		network: params,
		feeRate: feeRate,
		client:  NewClient(network),
	}
}

// BuildTransaction constrói uma transação P2WPKH.
//   - fromAddress: endereço de origem (bc1q.../tb1q...)
//   - toAddress: endereço de destino
//   - amount: valor em satoshis
//   - privKey: chave privada para assinar
//
// Retorna a transação serializada pronta para broadcast.
func (tb *TxBuilder) BuildTransaction(fromAddress, toAddress string, amount int64, privKey *btcec.PrivateKey) (*TxResult, error) {
	// 1. Buscar UTXOs disponíveis
	utxos, err := tb.client.GetUTXOs(fromAddress)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar UTXOs: %w", err)
	}

	if len(utxos) == 0 {
		return nil, fmt.Errorf("sem UTXOs disponíveis para %s", fromAddress)
	}

	// 2. Selecionar UTXOs (coin selection — simples: maior primeiro)
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Value > utxos[j].Value
	})

	var selectedUTXOs []UTXO
	var totalInput int64
	estimatedSize := int64(11 + 31*2) // overhead + 2 outputs P2WPKH

	for _, utxo := range utxos {
		selectedUTXOs = append(selectedUTXOs, utxo)
		totalInput += utxo.Value
		estimatedSize += 68 // tamanho estimado de input P2WPKH

		estimatedFee := estimatedSize * tb.feeRate
		if totalInput >= amount+estimatedFee {
			break
		}
	}

	fee := estimatedSize * tb.feeRate
	if totalInput < amount+fee {
		return nil, fmt.Errorf("saldo insuficiente: disponível %d sats, necessário %d + %d fee", totalInput, amount, fee)
	}

	// 3. Construir transação
	tx := wire.NewMsgTx(wire.TxVersion)

	// Inputs
	for _, utxo := range selectedUTXOs {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, fmt.Errorf("hash inválido: %w", err)
		}
		outPoint := wire.NewOutPoint(hash, uint32(utxo.Vout))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		txIn.Sequence = 0xfffffffd // sinaliza RBF
		tx.AddTxIn(txIn)
	}

	// Output: pagamento
	destAddr, err := btcutil.DecodeAddress(toAddress, tb.network)
	if err != nil {
		return nil, fmt.Errorf("endereço de destino inválido: %w", err)
	}
	destScript, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		return nil, fmt.Errorf("script de destino: %w", err)
	}
	tx.AddTxOut(wire.NewTxOut(amount, destScript))

	// Output: troco (se necessário)
	change := totalInput - amount - fee
	if change > 546 { // dust limit
		changeAddr, err := btcutil.DecodeAddress(fromAddress, tb.network)
		if err != nil {
			return nil, fmt.Errorf("endereço de troco inválido: %w", err)
		}
		changeScript, err := txscript.PayToAddrScript(changeAddr)
		if err != nil {
			return nil, fmt.Errorf("script de troco: %w", err)
		}
		tx.AddTxOut(wire.NewTxOut(change, changeScript))
	} else {
		// Troco menor que dust vai pra mineradores como fee extra
		fee += change
	}

	// 4. Assinar inputs (P2WPKH witness)
	pubKeyHash := btcutil.Hash160(privKey.PubKey().SerializeCompressed())
	p2wpkhScript, err := txscript.NewScriptBuilder().
		AddOp(txscript.OP_0).
		AddData(pubKeyHash).
		Script()
	if err != nil {
		return nil, fmt.Errorf("script P2WPKH: %w", err)
	}

	for i, utxo := range selectedUTXOs {
		sigHashes := txscript.NewTxSigHashes(tx, txscript.NewCannedPrevOutputFetcher(
			p2wpkhScript, utxo.Value,
		))

		witness, err := txscript.WitnessSignature(
			tx, sigHashes, i, utxo.Value, p2wpkhScript,
			txscript.SigHashAll, privKey, true,
		)
		if err != nil {
			return nil, fmt.Errorf("falha ao assinar input %d: %w", i, err)
		}
		tx.TxIn[i].Witness = witness
	}

	// 5. Serializar
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return nil, fmt.Errorf("falha ao serializar TX: %w", err)
	}

	rawHex := hex.EncodeToString(buf.Bytes())

	return &TxResult{
		TxID:     tx.TxHash().String(),
		RawTxHex: rawHex,
		Fee:      fee,
		Size:     buf.Len(),
	}, nil
}

// Broadcast envia uma transação para a rede.
func (tb *TxBuilder) Broadcast(rawTxHex string) (string, error) {
	return tb.client.BroadcastTxRaw(rawTxHex)
}

// EstimateFee retorna a taxa estimada para uma transação.
func (tb *TxBuilder) EstimateFee(numInputs, numOutputs int) int64 {
	// P2WPKH: ~68 vBytes por input, ~31 vBytes por output, ~11 overhead
	vsize := int64(11 + numInputs*68 + numOutputs*31)
	return vsize * tb.feeRate
}
