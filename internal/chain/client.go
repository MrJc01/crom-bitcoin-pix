package chain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIs suportadas para consulta on-chain.
const (
	EsploraMainnet = "https://mempool.space/api"
	EsploraTestnet = "https://mempool.space/testnet/api"
)

// Client consulta a blockchain Bitcoin via API Esplora (mempool.space).
// Permite obter saldo, UTXOs e broadcast de transações sem full node.
type Client struct {
	baseURL    string
	network    string
	httpClient *http.Client
}

// UTXO representa uma saída de transação não gasta.
type UTXO struct {
	TxID   string `json:"txid"`
	Vout   int    `json:"vout"`
	Value  int64  `json:"value"`
	Status struct {
		Confirmed   bool  `json:"confirmed"`
		BlockHeight int64 `json:"block_height"`
	} `json:"status"`
}

// TxInfo contém informações sobre uma transação.
type TxInfo struct {
	TxID     string `json:"txid"`
	Fee      int64  `json:"fee"`
	Confirmed bool
	BlockHeight int64
}

// FeeEstimate contém estimativas de taxa.
type FeeEstimate struct {
	FastestFee  int64 `json:"fastestFee"`
	HalfHourFee int64 `json:"halfHourFee"`
	HourFee     int64 `json:"hourFee"`
	MinimumFee  int64 `json:"minimumFee"`
}

// NewClient cria um client para a rede especificada.
func NewClient(network string) *Client {
	baseURL := EsploraMainnet
	if network == "testnet" {
		baseURL = EsploraTestnet
	}

	return &Client{
		baseURL: baseURL,
		network: network,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetBalance retorna o saldo total (confirmado + não confirmado) em satoshis.
func (c *Client) GetBalance(address string) (confirmed, unconfirmed int64, err error) {
	type balanceResp struct {
		Address    string `json:"address"`
		ChainStats struct {
			FundedSum int64 `json:"funded_txo_sum"`
			SpentSum  int64 `json:"spent_txo_sum"`
		} `json:"chain_stats"`
		MempoolStats struct {
			FundedSum int64 `json:"funded_txo_sum"`
			SpentSum  int64 `json:"spent_txo_sum"`
		} `json:"mempool_stats"`
	}

	var resp balanceResp
	if err := c.get(fmt.Sprintf("/address/%s", address), &resp); err != nil {
		return 0, 0, fmt.Errorf("falha ao consultar saldo: %w", err)
	}

	confirmed = resp.ChainStats.FundedSum - resp.ChainStats.SpentSum
	unconfirmed = resp.MempoolStats.FundedSum - resp.MempoolStats.SpentSum

	return confirmed, unconfirmed, nil
}

// GetUTXOs retorna as UTXOs de um endereço.
func (c *Client) GetUTXOs(address string) ([]UTXO, error) {
	var utxos []UTXO
	if err := c.get(fmt.Sprintf("/address/%s/utxo", address), &utxos); err != nil {
		return nil, fmt.Errorf("falha ao consultar UTXOs: %w", err)
	}
	return utxos, nil
}

// GetFeeEstimates retorna estimativas de taxa em sat/vByte.
func (c *Client) GetFeeEstimates() (*FeeEstimate, error) {
	var fees FeeEstimate
	if err := c.get("/v1/fees/recommended", &fees); err != nil {
		return nil, fmt.Errorf("falha ao consultar fees: %w", err)
	}
	return &fees, nil
}

// BroadcastTx envia uma transação raw para a rede.
func (c *Client) BroadcastTx(rawTxHex string) (string, error) {
	url := fmt.Sprintf("%s/tx", c.baseURL)
	resp, err := c.httpClient.Post(url, "text/plain", io.NopCloser(
		io.LimitReader(
			io.NopCloser(nil), 0,
		),
	))
	if err != nil {
		return "", fmt.Errorf("falha ao broadcast: %w", err)
	}
	// Na prática, enviamos o hex como body
	_ = resp
	return "", fmt.Errorf("broadcast via API: use BroadcastTxRaw")
}

// BroadcastTxRaw envia uma transação serializada.
func (c *Client) BroadcastTxRaw(rawTxHex string) (string, error) {
	url := fmt.Sprintf("%s/tx", c.baseURL)

	resp, err := c.httpClient.Post(url, "text/plain",
		io.NopCloser(io.LimitReader(
			mustReader(rawTxHex), int64(len(rawTxHex)),
		)),
	)
	if err != nil {
		return "", fmt.Errorf("falha ao broadcast: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("broadcast rejeitado: %s", string(body))
	}

	return string(body), nil // retorna TXID
}

// GetBlockHeight retorna a altura atual da blockchain.
func (c *Client) GetBlockHeight() (int64, error) {
	url := fmt.Sprintf("%s/blocks/tip/height", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var height int64
	if err := json.Unmarshal(body, &height); err != nil {
		return 0, fmt.Errorf("falha ao parsear altura: %w", err)
	}
	return height, nil
}

// Network retorna a rede configurada.
func (c *Client) Network() string {
	return c.network
}

// get faz um GET request e decodifica o JSON response.
func (c *Client) get(path string, target interface{}) error {
	url := c.baseURL + path
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API retornou %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// mustReader cria um io.Reader de uma string.
func mustReader(s string) io.Reader {
	return io.NopCloser(io.LimitReader(
		readerFromString(s), int64(len(s)),
	))
}

type stringReader struct {
	s string
	i int
}

func (r *stringReader) Read(p []byte) (n int, err error) {
	if r.i >= len(r.s) {
		return 0, io.EOF
	}
	n = copy(p, r.s[r.i:])
	r.i += n
	return n, nil
}

func readerFromString(s string) io.Reader {
	return &stringReader{s: s}
}
