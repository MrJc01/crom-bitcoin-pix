package lightning

import (
	"errors"
	"fmt"
)

// Erros do módulo Lightning.
var (
	ErrNotConnected  = errors.New("não conectado ao nó Lightning")
	ErrNoChannels    = errors.New("nenhum canal aberto")
	ErrInvoiceExpired = errors.New("invoice expirada")
	ErrPaymentFailed = errors.New("pagamento falhou")
)

// NodeInfo contém informações sobre o nó Lightning.
type NodeInfo struct {
	PubKey       string
	Alias        string
	NumChannels  int
	NumPeers     int
	BlockHeight  int64
	Synced       bool
	Version      string
}

// Channel representa um canal Lightning.
type Channel struct {
	ChanID       uint64
	RemotePubkey string
	Capacity     int64 // sats
	LocalBalance int64 // sats
	RemoteBalance int64 // sats
	Active       bool
}

// Invoice representa um invoice BOLT-11.
type Invoice struct {
	PaymentRequest string // invoice string bolt11
	PaymentHash    string
	Amount         int64 // sats
	Memo           string
	CreatedAt      int64
	ExpiresAt      int64
	Settled        bool
}

// Payment representa um pagamento enviado.
type Payment struct {
	PaymentHash    string
	Amount         int64 // sats
	Fee            int64 // sats
	Status         string
	Preimage       string
}

// Config contém a configuração de conexão LND.
type Config struct {
	Host         string // ex: localhost:10009
	TLSCertPath  string // ex: ~/.lnd/tls.cert
	MacaroonPath string // ex: ~/.lnd/data/chain/bitcoin/mainnet/admin.macaroon
	Network      string // mainnet, testnet
}

// DefaultConfig retorna a configuração padrão do LND.
func DefaultConfig(network string) *Config {
	return &Config{
		Host:         "localhost:10009",
		TLSCertPath:  "~/.lnd/tls.cert",
		MacaroonPath: fmt.Sprintf("~/.lnd/data/chain/bitcoin/%s/admin.macaroon", network),
		Network:      network,
	}
}

// Client é a interface para interagir com um nó Lightning.
// Implementações: LNDClient (gRPC), EmbeddedClient (embutido).
type Client interface {
	// Conexão
	Connect() error
	Close() error
	IsConnected() bool

	// Info
	GetInfo() (*NodeInfo, error)

	// Invoices
	CreateInvoice(amount int64, memo string) (*Invoice, error)
	LookupInvoice(paymentHash string) (*Invoice, error)
	DecodeInvoice(bolt11 string) (*Invoice, error)

	// Pagamentos
	PayInvoice(bolt11 string) (*Payment, error)
	PayInvoiceWithAmount(bolt11 string, amount int64) (*Payment, error)

	// Canais
	ListChannels() ([]Channel, error)
	OpenChannel(peerPubkey string, amount int64) (string, error)
	CloseChannel(chanID uint64) (string, error)
}

// StubClient é uma implementação stub para quando LND não está disponível.
// Retorna erros informativos indicando que o Lightning não está configurado.
type StubClient struct {
	network string
}

// NewStubClient cria um client stub.
func NewStubClient(network string) *StubClient {
	return &StubClient{network: network}
}

func (s *StubClient) Connect() error {
	return fmt.Errorf("Lightning não configurado — execute 'crom-pay lightning setup' para conectar a um nó LND")
}

func (s *StubClient) Close() error { return nil }
func (s *StubClient) IsConnected() bool { return false }

func (s *StubClient) GetInfo() (*NodeInfo, error) {
	return nil, ErrNotConnected
}

func (s *StubClient) CreateInvoice(amount int64, memo string) (*Invoice, error) {
	return nil, ErrNotConnected
}

func (s *StubClient) LookupInvoice(paymentHash string) (*Invoice, error) {
	return nil, ErrNotConnected
}

func (s *StubClient) DecodeInvoice(bolt11 string) (*Invoice, error) {
	return nil, ErrNotConnected
}

func (s *StubClient) PayInvoice(bolt11 string) (*Payment, error) {
	return nil, ErrNotConnected
}

func (s *StubClient) PayInvoiceWithAmount(bolt11 string, amount int64) (*Payment, error) {
	return nil, ErrNotConnected
}

func (s *StubClient) ListChannels() ([]Channel, error) {
	return nil, ErrNotConnected
}

func (s *StubClient) OpenChannel(peerPubkey string, amount int64) (string, error) {
	return "", ErrNotConnected
}

func (s *StubClient) CloseChannel(chanID uint64) (string, error) {
	return "", ErrNotConnected
}
