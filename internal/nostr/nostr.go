package nostr

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

// Constantes de derivação NIP-06
// Path: m/44'/1237'/0'/0/0
const (
	NostrCoinType = 1237
)

// Keys contém o par de chaves Nostr derivado do seed Bitcoin.
type Keys struct {
	PrivateKey string // hex
	PublicKey  string // hex
}

// DeriveFromSeed deriva chaves Nostr do mesmo seed BIP-39 usando NIP-06.
// Path: m/44'/1237'/account'/0/index
func DeriveFromSeed(seed []byte, account, index uint32) (*Keys, error) {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar master key: %w", err)
	}
	defer masterKey.Zero()

	// m/44' (purpose BIP-44)
	purpose, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return nil, fmt.Errorf("derivação purpose: %w", err)
	}

	// m/44'/1237' (coin type Nostr — NIP-06)
	coin, err := purpose.Derive(hdkeychain.HardenedKeyStart + NostrCoinType)
	if err != nil {
		return nil, fmt.Errorf("derivação coin: %w", err)
	}

	// m/44'/1237'/account'
	acct, err := coin.Derive(hdkeychain.HardenedKeyStart + account)
	if err != nil {
		return nil, fmt.Errorf("derivação account: %w", err)
	}

	// m/44'/1237'/account'/0
	chain, err := acct.Derive(0)
	if err != nil {
		return nil, fmt.Errorf("derivação chain: %w", err)
	}

	// m/44'/1237'/account'/0/index
	key, err := chain.Derive(index)
	if err != nil {
		return nil, fmt.Errorf("derivação index: %w", err)
	}

	// Extrair chave privada
	privKey, err := key.ECPrivKey()
	if err != nil {
		return nil, fmt.Errorf("falha ao extrair privkey: %w", err)
	}

	privBytes := privKey.Serialize()
	privHex := hex.EncodeToString(privBytes)

	// Derivar chave pública
	pubKey, err := nostr.GetPublicKey(privHex)
	if err != nil {
		return nil, fmt.Errorf("falha ao derivar pubkey: %w", err)
	}

	return &Keys{
		PrivateKey: privHex,
		PublicKey:  pubKey,
	}, nil
}

// Npub retorna a chave pública em formato bech32 (npub1...).
func (k *Keys) Npub() string {
	npub, _ := nip19.EncodePublicKey(k.PublicKey)
	return npub
}

// Nsec retorna a chave privada em formato bech32 (nsec1...).
func (k *Keys) Nsec() string {
	nsec, _ := nip19.EncodePrivateKey(k.PrivateKey)
	return nsec
}

// Identity contém metadados do perfil Nostr.
type Identity struct {
	Name    string `json:"name,omitempty"`
	About   string `json:"about,omitempty"`
	Picture string `json:"picture,omitempty"`
	NIP05   string `json:"nip05,omitempty"`
	LUD16   string `json:"lud16,omitempty"` // Lightning Address
}

// SignEvent assina um evento Nostr com a chave privada.
func (k *Keys) SignEvent(event *nostr.Event) error {
	event.PubKey = k.PublicKey
	event.CreatedAt = nostr.Timestamp(time.Now().Unix())
	event.ID = event.GetID()
	return event.Sign(k.PrivateKey)
}

// CreateMetadataEvent cria um evento kind 0 (metadata/perfil).
func (k *Keys) CreateMetadataEvent(identity Identity) (*nostr.Event, error) {
	content := fmt.Sprintf(
		`{"name":"%s","about":"%s","picture":"%s","nip05":"%s","lud16":"%s"}`,
		identity.Name, identity.About, identity.Picture,
		identity.NIP05, identity.LUD16,
	)

	event := &nostr.Event{
		Kind:    nostr.KindProfileMetadata,
		Content: content,
		Tags:    nostr.Tags{},
	}

	if err := k.SignEvent(event); err != nil {
		return nil, fmt.Errorf("falha ao assinar metadata: %w", err)
	}

	return event, nil
}

// CreateNoteEvent cria um evento kind 1 (nota/post).
func (k *Keys) CreateNoteEvent(content string) (*nostr.Event, error) {
	event := &nostr.Event{
		Kind:    nostr.KindTextNote,
		Content: content,
		Tags:    nostr.Tags{},
	}

	if err := k.SignEvent(event); err != nil {
		return nil, fmt.Errorf("falha ao assinar nota: %w", err)
	}

	return event, nil
}

// CreateZapRequest cria um evento kind 9734 (zap request).
func (k *Keys) CreateZapRequest(recipientPubkey string, amountMsats int64, relays []string) (*nostr.Event, error) {
	tags := nostr.Tags{
		{"p", recipientPubkey},
		{"amount", fmt.Sprintf("%d", amountMsats)},
		{"relays"},
	}

	// Adicionar relays à tag
	for _, r := range relays {
		tags[2] = append(tags[2], r)
	}

	event := &nostr.Event{
		Kind:    9734,
		Content: "",
		Tags:    tags,
	}

	if err := k.SignEvent(event); err != nil {
		return nil, fmt.Errorf("falha ao assinar zap request: %w", err)
	}

	return event, nil
}

// RelayPool gerencia conexões com múltiplos relays Nostr.
type RelayPool struct {
	URLs    []string
	relays  []*nostr.Relay
}

// DefaultRelays retorna relays públicos confiáveis.
func DefaultRelays() []string {
	return []string{
		"wss://relay.damus.io",
		"wss://relay.nostr.band",
		"wss://nos.lol",
		"wss://relay.snort.social",
		"wss://nostr.wine",
	}
}

// NewRelayPool cria um pool de relays.
func NewRelayPool(urls []string) *RelayPool {
	if len(urls) == 0 {
		urls = DefaultRelays()
	}
	return &RelayPool{URLs: urls}
}

// Connect conecta a todos os relays do pool.
func (rp *RelayPool) Connect(ctx context.Context) error {
	for _, url := range rp.URLs {
		relay, err := nostr.RelayConnect(ctx, url)
		if err != nil {
			continue // skip relays offline
		}
		rp.relays = append(rp.relays, relay)
	}

	if len(rp.relays) == 0 {
		return fmt.Errorf("nenhum relay disponível")
	}

	return nil
}

// Publish publica um evento em todos os relays conectados.
func (rp *RelayPool) Publish(ctx context.Context, event nostr.Event) error {
	var lastErr error
	published := 0

	for _, relay := range rp.relays {
		if err := relay.Publish(ctx, event); err != nil {
			lastErr = err
			continue
		}
		published++
	}

	if published == 0 {
		return fmt.Errorf("falha ao publicar em todos os relays: %w", lastErr)
	}

	return nil
}

// Close fecha todas as conexões.
func (rp *RelayPool) Close() {
	for _, relay := range rp.relays {
		relay.Close()
	}
}

// VerifyNIP05 verifica um endereço NIP-05 (user@domain).
func VerifyNIP05(ctx context.Context, nip05Address, expectedPubkey string) (bool, error) {
	// Parsear name@domain
	var name, domain string
	for i, c := range nip05Address {
		if c == '@' {
			name = nip05Address[:i]
			domain = nip05Address[i+1:]
			break
		}
	}
	if domain == "" {
		return false, fmt.Errorf("formato NIP-05 inválido: %s", nip05Address)
	}

	// Buscar .well-known/nostr.json
	url := fmt.Sprintf("https://%s/.well-known/nostr.json?name=%s", domain, name)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false, fmt.Errorf("falha ao verificar NIP-05: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Names map[string]string `json:"names"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("falha ao parsear NIP-05: %w", err)
	}

	pubkey, ok := result.Names[name]
	if !ok {
		return false, nil
	}

	return pubkey == expectedPubkey, nil
}

// Fingerprint retorna um hash curto da chave pública para display.
func (k *Keys) Fingerprint() string {
	hash := sha256.Sum256([]byte(k.PublicKey))
	return hex.EncodeToString(hash[:4])
}
