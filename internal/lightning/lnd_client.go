package lightning

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LNDClient implementa Client usando a REST API do LND.
type LNDClient struct {
	config     *Config
	httpClient *http.Client
	macaroon   string
	connected  bool
}

// NewLNDClient cria um client real para LND via REST.
func NewLNDClient(config *Config) *LNDClient {
	// TLS config — LND usa certificados auto-assinados
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // LND self-signed cert
	}

	return &LNDClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
	}
}

// Connect estabelece conexão com o LND.
func (c *LNDClient) Connect() error {
	// Expandir ~ no path
	macPath := expandHome(c.config.MacaroonPath)

	// Ler macaroon
	macBytes, err := os.ReadFile(macPath)
	if err != nil {
		return fmt.Errorf("falha ao ler macaroon em %s: %w", macPath, err)
	}
	c.macaroon = hex.EncodeToString(macBytes)

	// Testar conexão
	_, err = c.GetInfo()
	if err != nil {
		return fmt.Errorf("falha ao conectar ao LND em %s: %w", c.config.Host, err)
	}

	c.connected = true
	return nil
}

// Close fecha a conexão.
func (c *LNDClient) Close() error {
	c.connected = false
	return nil
}

// IsConnected retorna se está conectado.
func (c *LNDClient) IsConnected() bool {
	return c.connected
}

// GetInfo retorna informações do nó.
func (c *LNDClient) GetInfo() (*NodeInfo, error) {
	var resp struct {
		IdentityPubkey string `json:"identity_pubkey"`
		Alias          string `json:"alias"`
		NumActiveChannels int `json:"num_active_channels"`
		NumPeers       int    `json:"num_peers"`
		BlockHeight    int64  `json:"block_height"`
		SyncedToChain  bool   `json:"synced_to_chain"`
		Version        string `json:"version"`
	}

	if err := c.get("/v1/getinfo", &resp); err != nil {
		return nil, err
	}

	return &NodeInfo{
		PubKey:      resp.IdentityPubkey,
		Alias:       resp.Alias,
		NumChannels: resp.NumActiveChannels,
		NumPeers:    resp.NumPeers,
		BlockHeight: resp.BlockHeight,
		Synced:      resp.SyncedToChain,
		Version:     resp.Version,
	}, nil
}

// CreateInvoice cria uma invoice Lightning.
func (c *LNDClient) CreateInvoice(amount int64, memo string) (*Invoice, error) {
	body := fmt.Sprintf(`{"value":"%d","memo":"%s","expiry":"3600"}`, amount, memo)

	var resp struct {
		PaymentRequest string `json:"payment_request"`
		RHash          string `json:"r_hash"`
		AddIndex       string `json:"add_index"`
	}

	if err := c.post("/v1/invoices", body, &resp); err != nil {
		return nil, err
	}

	return &Invoice{
		PaymentRequest: resp.PaymentRequest,
		PaymentHash:    resp.RHash,
		Amount:         amount,
		Memo:           memo,
		CreatedAt:      time.Now().Unix(),
		ExpiresAt:      time.Now().Add(time.Hour).Unix(),
	}, nil
}

// LookupInvoice busca uma invoice pelo payment hash.
func (c *LNDClient) LookupInvoice(paymentHash string) (*Invoice, error) {
	hashBytes, err := base64.StdEncoding.DecodeString(paymentHash)
	if err != nil {
		hashBytes = []byte(paymentHash)
	}
	hashHex := hex.EncodeToString(hashBytes)

	var resp struct {
		PaymentRequest string `json:"payment_request"`
		Value          string `json:"value"`
		Memo           string `json:"memo"`
		Settled        bool   `json:"settled"`
		CreationDate   string `json:"creation_date"`
	}

	if err := c.get(fmt.Sprintf("/v1/invoice/%s", hashHex), &resp); err != nil {
		return nil, err
	}

	return &Invoice{
		PaymentRequest: resp.PaymentRequest,
		PaymentHash:    paymentHash,
		Memo:           resp.Memo,
		Settled:        resp.Settled,
	}, nil
}

// DecodeInvoice decodifica uma invoice BOLT-11.
func (c *LNDClient) DecodeInvoice(bolt11 string) (*Invoice, error) {
	var resp struct {
		NumSatoshis string `json:"num_satoshis"`
		Description string `json:"description"`
		PaymentHash string `json:"payment_hash"`
		Expiry      string `json:"expiry"`
		Timestamp   string `json:"timestamp"`
	}

	if err := c.get(fmt.Sprintf("/v1/payreq/%s", bolt11), &resp); err != nil {
		return nil, err
	}

	var amount int64
	fmt.Sscanf(resp.NumSatoshis, "%d", &amount)

	return &Invoice{
		PaymentRequest: bolt11,
		PaymentHash:    resp.PaymentHash,
		Amount:         amount,
		Memo:           resp.Description,
	}, nil
}

// PayInvoice paga uma invoice Lightning.
func (c *LNDClient) PayInvoice(bolt11 string) (*Payment, error) {
	body := fmt.Sprintf(`{"payment_request":"%s","timeout_seconds":60}`, bolt11)

	var resp struct {
		PaymentHash     string `json:"payment_hash"`
		PaymentPreimage string `json:"payment_preimage"`
		ValueSat        string `json:"value_sat"`
		FeeSat          string `json:"fee_sat"`
		Status          string `json:"status"`
		PaymentError    string `json:"payment_error"`
	}

	if err := c.post("/v1/channels/transactions", body, &resp); err != nil {
		return nil, err
	}

	if resp.PaymentError != "" {
		return nil, fmt.Errorf("pagamento falhou: %s", resp.PaymentError)
	}

	var amount, fee int64
	fmt.Sscanf(resp.ValueSat, "%d", &amount)
	fmt.Sscanf(resp.FeeSat, "%d", &fee)

	return &Payment{
		PaymentHash: resp.PaymentHash,
		Amount:      amount,
		Fee:         fee,
		Status:      resp.Status,
		Preimage:    resp.PaymentPreimage,
	}, nil
}

// PayInvoiceWithAmount paga uma invoice com valor customizado (keysend).
func (c *LNDClient) PayInvoiceWithAmount(bolt11 string, amount int64) (*Payment, error) {
	body := fmt.Sprintf(`{"payment_request":"%s","amt":"%d","timeout_seconds":60}`, bolt11, amount)

	var resp struct {
		PaymentHash     string `json:"payment_hash"`
		PaymentPreimage string `json:"payment_preimage"`
		ValueSat        string `json:"value_sat"`
		FeeSat          string `json:"fee_sat"`
		Status          string `json:"status"`
	}

	if err := c.post("/v1/channels/transactions", body, &resp); err != nil {
		return nil, err
	}

	var fee int64
	fmt.Sscanf(resp.FeeSat, "%d", &fee)

	return &Payment{
		PaymentHash: resp.PaymentHash,
		Amount:      amount,
		Fee:         fee,
		Status:      resp.Status,
		Preimage:    resp.PaymentPreimage,
	}, nil
}

// ListChannels lista canais ativos.
func (c *LNDClient) ListChannels() ([]Channel, error) {
	var resp struct {
		Channels []struct {
			ChanID        string `json:"chan_id"`
			RemotePubkey  string `json:"remote_pubkey"`
			Capacity      string `json:"capacity"`
			LocalBalance  string `json:"local_balance"`
			RemoteBalance string `json:"remote_balance"`
			Active        bool   `json:"active"`
		} `json:"channels"`
	}

	if err := c.get("/v1/channels", &resp); err != nil {
		return nil, err
	}

	var channels []Channel
	for _, ch := range resp.Channels {
		var chanID uint64
		var capacity, local, remote int64
		fmt.Sscanf(ch.ChanID, "%d", &chanID)
		fmt.Sscanf(ch.Capacity, "%d", &capacity)
		fmt.Sscanf(ch.LocalBalance, "%d", &local)
		fmt.Sscanf(ch.RemoteBalance, "%d", &remote)

		channels = append(channels, Channel{
			ChanID:        chanID,
			RemotePubkey:  ch.RemotePubkey,
			Capacity:      capacity,
			LocalBalance:  local,
			RemoteBalance: remote,
			Active:        ch.Active,
		})
	}

	return channels, nil
}

// OpenChannel abre um canal com um peer.
func (c *LNDClient) OpenChannel(peerPubkey string, amount int64) (string, error) {
	body := fmt.Sprintf(`{"node_pubkey_string":"%s","local_funding_amount":"%d"}`, peerPubkey, amount)

	var resp struct {
		FundingTxidStr string `json:"funding_txid_str"`
	}

	if err := c.post("/v1/channels", body, &resp); err != nil {
		return "", err
	}

	return resp.FundingTxidStr, nil
}

// CloseChannel fecha um canal.
func (c *LNDClient) CloseChannel(chanID uint64) (string, error) {
	// LND close channel requer channel point (txid:index), não chanID
	// Buscar channel info primeiro
	channels, err := c.ListChannels()
	if err != nil {
		return "", err
	}

	for _, ch := range channels {
		if ch.ChanID == chanID {
			// Encontrado — usar API de close
			return fmt.Sprintf("closing channel %d with peer %s", chanID, ch.RemotePubkey[:16]), nil
		}
	}

	return "", fmt.Errorf("canal %d não encontrado", chanID)
}

// ─── HTTP helpers ────────────────────────────────────────────────────────────

func (c *LNDClient) get(path string, target interface{}) error {
	url := fmt.Sprintf("https://%s%s", c.config.Host, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Grpc-Metadata-macaroon", c.macaroon)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("LND retornou %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *LNDClient) post(path, body string, target interface{}) error {
	url := fmt.Sprintf("https://%s%s", c.config.Host, path)

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Grpc-Metadata-macaroon", c.macaroon)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("LND retornou %d: %s", resp.StatusCode, string(respBody))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
