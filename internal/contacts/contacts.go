package contacts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MrJc01/crom-bitcoin-pix/internal/storage"
)

var contactPrefix = []byte("contact:")

// Contact representa um contato salvo.
type Contact struct {
	Name           string `json:"name"`
	BitcoinAddress string `json:"bitcoin_address,omitempty"`
	LightningAddr  string `json:"lightning_address,omitempty"`
	NostrNpub      string `json:"nostr_npub,omitempty"`
	NIP05          string `json:"nip05,omitempty"`
	Notes          string `json:"notes,omitempty"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
}

// Manager gerencia contatos no BadgerDB.
type Manager struct {
	store *storage.Store
}

// NewManager cria um gerenciador de contatos.
func NewManager(store *storage.Store) *Manager {
	return &Manager{store: store}
}

// Add adiciona ou atualiza um contato.
func (m *Manager) Add(contact Contact) error {
	if contact.Name == "" {
		return fmt.Errorf("nome do contato é obrigatório")
	}

	now := time.Now().Unix()
	if contact.CreatedAt == 0 {
		contact.CreatedAt = now
	}
	contact.UpdatedAt = now

	data, err := json.Marshal(contact)
	if err != nil {
		return fmt.Errorf("falha ao serializar contato: %w", err)
	}

	return m.store.Put(contactPrefix, []byte(contact.Name), data)
}

// Get retorna um contato pelo nome.
func (m *Manager) Get(name string) (*Contact, error) {
	data, err := m.store.Get(contactPrefix, []byte(name))
	if err != nil {
		return nil, fmt.Errorf("contato '%s' não encontrado", name)
	}

	var contact Contact
	if err := json.Unmarshal(data, &contact); err != nil {
		return nil, fmt.Errorf("falha ao parsear contato: %w", err)
	}

	return &contact, nil
}

// Remove remove um contato pelo nome.
func (m *Manager) Remove(name string) error {
	return m.store.Delete(contactPrefix, []byte(name))
}

// List retorna todos os contatos.
func (m *Manager) List() ([]Contact, error) {
	entries, err := m.store.List(contactPrefix)
	if err != nil {
		return nil, err
	}

	var contacts []Contact
	for _, data := range entries {
		var contact Contact
		if err := json.Unmarshal(data, &contact); err != nil {
			continue
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// Resolve tenta resolver um destino de pagamento.
func (m *Manager) Resolve(destination string) (*Contact, error) {
	contact, err := m.Get(destination)
	if err == nil {
		return contact, nil
	}

	if isBitcoinAddr(destination) {
		return &Contact{Name: destination, BitcoinAddress: destination}, nil
	}
	if isLightningAddr(destination) {
		return &Contact{Name: destination, LightningAddr: destination}, nil
	}
	if isNIP05Addr(destination) {
		return &Contact{Name: destination, NIP05: destination}, nil
	}

	return nil, fmt.Errorf("destino não reconhecido: %s", destination)
}

func isBitcoinAddr(s string) bool {
	return len(s) > 3 && (s[:3] == "bc1" || s[:3] == "tb1" || s[:1] == "1" || s[:1] == "3")
}

func isLightningAddr(s string) bool {
	return len(s) > 4 && (s[:4] == "lnbc" || s[:4] == "lntb")
}

func isNIP05Addr(s string) bool {
	for _, c := range s {
		if c == '@' {
			return true
		}
	}
	return false
}
