package wallet

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"golang.org/x/crypto/ripemd160"
)

// Constantes de derivação BIP-84 (Native SegWit / Bech32)
//
// Path: m/84'/0'/0'/0/0
//   84'  = purpose (BIP-84: Native SegWit)
//   0'   = coin type (0 = Bitcoin mainnet, 1 = testnet)
//   0'   = account
//   0    = change (0 = external/receiving, 1 = internal/change)
//   0    = address index
const (
	PurposeBIP84  = 84
	CoinTypeBTC   = 0
	CoinTypeTest  = 1
	DefaultAcct   = 0
	ExternalChain = 0
)

var (
	ErrSeedTooShort  = errors.New("seed deve ter pelo menos 16 bytes")
	ErrDeriveFailed  = errors.New("falha na derivação de chaves")
	ErrInvalidNetwork = errors.New("rede inválida: use 'mainnet' ou 'testnet'")
)

// Keychain gerencia a derivação hierárquica de chaves (HD Wallet).
type Keychain struct {
	masterKey *hdkeychain.ExtendedKey
	network   *chaincfg.Params
}

// NewKeychain cria um Keychain a partir de uma seed BIP-39 (512 bits).
//
// network: "mainnet" ou "testnet"
func NewKeychain(seed []byte, network string) (*Keychain, error) {
	if len(seed) < 16 {
		return nil, ErrSeedTooShort
	}

	params, err := getNetworkParams(network)
	if err != nil {
		return nil, err
	}

	masterKey, err := hdkeychain.NewMaster(seed, params)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDeriveFailed, err)
	}

	return &Keychain{
		masterKey: masterKey,
		network:   params,
	}, nil
}

// DeriveAddress gera o endereço SegWit nativo Bech32 (bc1q... / tb1q...) no índice especificado.
//
// Segue o path BIP-84: m/84'/coin'/0'/0/index
// Usa P2WPKH (Pay-to-Witness-Public-Key-Hash) conforme BIP-84.
func (kc *Keychain) DeriveAddress(index uint32) (string, error) {
	key, err := kc.deriveKey(index)
	if err != nil {
		return "", err
	}

	// Extrair chave pública comprimida
	pubKey, err := key.ECPubKey()
	if err != nil {
		return "", fmt.Errorf("falha ao extrair pubkey: %w", err)
	}

	// Hash160 = RIPEMD160(SHA256(pubkey)) — padrão Bitcoin
	pubKeyBytes := pubKey.SerializeCompressed()
	sha := sha256.Sum256(pubKeyBytes)
	rip := ripemd160.New()
	rip.Write(sha[:])
	pubKeyHash := rip.Sum(nil)

	// Criar endereço P2WPKH (SegWit nativo / Bech32)
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, kc.network)
	if err != nil {
		return "", fmt.Errorf("falha ao gerar endereço SegWit: %w", err)
	}

	return addr.EncodeAddress(), nil
}

// DerivePublicKey retorna a chave pública comprimida (33 bytes) no índice especificado.
func (kc *Keychain) DerivePublicKey(index uint32) ([]byte, error) {
	key, err := kc.deriveKey(index)
	if err != nil {
		return nil, err
	}

	pubKey, err := key.ECPubKey()
	if err != nil {
		return nil, fmt.Errorf("falha ao extrair pubkey: %w", err)
	}

	return pubKey.SerializeCompressed(), nil
}

// Network retorna o nome da rede configurada.
func (kc *Keychain) Network() string {
	if kc.network.Net == chaincfg.MainNetParams.Net {
		return "mainnet"
	}
	return "testnet"
}

// Zero limpa a chave mestra da memória.
// DEVE ser chamado quando o Keychain não for mais necessário.
func (kc *Keychain) Zero() {
	if kc.masterKey != nil {
		kc.masterKey.Zero()
		kc.masterKey = nil
	}
}

// deriveKey faz a derivação completa BIP-84: m/84'/coin'/0'/0/index
func (kc *Keychain) deriveKey(index uint32) (*hdkeychain.ExtendedKey, error) {
	coinType := uint32(CoinTypeBTC)
	if kc.network.Net != chaincfg.MainNetParams.Net {
		coinType = CoinTypeTest
	}

	// m/84' (purpose)
	purpose, err := kc.masterKey.Derive(hdkeychain.HardenedKeyStart + uint32(PurposeBIP84))
	if err != nil {
		return nil, fmt.Errorf("%w: purpose: %v", ErrDeriveFailed, err)
	}

	// m/84'/coin' (coin type)
	coin, err := purpose.Derive(hdkeychain.HardenedKeyStart + coinType)
	if err != nil {
		return nil, fmt.Errorf("%w: coin: %v", ErrDeriveFailed, err)
	}

	// m/84'/coin'/0' (account)
	account, err := coin.Derive(hdkeychain.HardenedKeyStart + uint32(DefaultAcct))
	if err != nil {
		return nil, fmt.Errorf("%w: account: %v", ErrDeriveFailed, err)
	}

	// m/84'/coin'/0'/0 (external chain)
	chain, err := account.Derive(uint32(ExternalChain))
	if err != nil {
		return nil, fmt.Errorf("%w: chain: %v", ErrDeriveFailed, err)
	}

	// m/84'/coin'/0'/0/index (address index)
	addrKey, err := chain.Derive(index)
	if err != nil {
		return nil, fmt.Errorf("%w: index: %v", ErrDeriveFailed, err)
	}

	return addrKey, nil
}

// getNetworkParams retorna os parâmetros de rede do Bitcoin.
func getNetworkParams(network string) (*chaincfg.Params, error) {
	switch network {
	case "mainnet":
		return &chaincfg.MainNetParams, nil
	case "testnet":
		return &chaincfg.TestNet3Params, nil
	case "regtest":
		return &chaincfg.RegressionNetParams, nil
	default:
		return nil, fmt.Errorf("%w: '%s'", ErrInvalidNetwork, network)
	}
}
