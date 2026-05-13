package contacts

import (
	"os"
	"testing"

	"github.com/MrJc01/crom-bitcoin-pix/internal/storage"
)

func setupTestStore(t *testing.T) (*storage.Store, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "contacts-test-*")
	if err != nil {
		t.Fatal(err)
	}
	store, err := storage.NewStore(dir)
	if err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}
	return store, func() {
		store.Close()
		os.RemoveAll(dir)
	}
}

func TestAddAndGet(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	mgr := NewManager(store)

	contact := Contact{
		Name:           "Alice",
		BitcoinAddress: "bc1qtest123",
		NIP05:          "alice@crom.run",
	}

	if err := mgr.Add(contact); err != nil {
		t.Fatalf("Add: %v", err)
	}

	got, err := mgr.Get("Alice")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if got.Name != "Alice" {
		t.Errorf("Name = %s, want Alice", got.Name)
	}
	if got.BitcoinAddress != "bc1qtest123" {
		t.Errorf("BitcoinAddress = %s, want bc1qtest123", got.BitcoinAddress)
	}
	if got.NIP05 != "alice@crom.run" {
		t.Errorf("NIP05 = %s, want alice@crom.run", got.NIP05)
	}
	if got.CreatedAt == 0 {
		t.Error("CreatedAt should be set")
	}
}

func TestAddEmptyName(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	mgr := NewManager(store)
	err := mgr.Add(Contact{})
	if err == nil {
		t.Error("Add with empty name should fail")
	}
}

func TestRemove(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	mgr := NewManager(store)
	mgr.Add(Contact{Name: "Bob", BitcoinAddress: "bc1qbob"})

	if err := mgr.Remove("Bob"); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	_, err := mgr.Get("Bob")
	if err == nil {
		t.Error("Get after Remove should fail")
	}
}

func TestList(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	mgr := NewManager(store)
	mgr.Add(Contact{Name: "Alice"})
	mgr.Add(Contact{Name: "Bob"})
	mgr.Add(Contact{Name: "Charlie"})

	list, err := mgr.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("List() len = %d, want 3", len(list))
	}
}

func TestListEmpty(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	mgr := NewManager(store)
	list, err := mgr.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("empty List() len = %d, want 0", len(list))
	}
}

func TestResolve(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	mgr := NewManager(store)
	mgr.Add(Contact{Name: "Alice", BitcoinAddress: "bc1qalice"})

	// Resolver por nome de contato
	c, err := mgr.Resolve("Alice")
	if err != nil {
		t.Fatalf("Resolve Alice: %v", err)
	}
	if c.BitcoinAddress != "bc1qalice" {
		t.Error("should resolve saved contact")
	}

	// Resolver endereço Bitcoin direto
	c, err = mgr.Resolve("bc1qdirect")
	if err != nil {
		t.Fatalf("Resolve btc: %v", err)
	}
	if c.BitcoinAddress != "bc1qdirect" {
		t.Error("should resolve BTC address")
	}

	// Resolver NIP-05
	c, err = mgr.Resolve("user@crom.run")
	if err != nil {
		t.Fatalf("Resolve nip05: %v", err)
	}
	if c.NIP05 != "user@crom.run" {
		t.Error("should resolve NIP-05")
	}

	// Destino inválido
	_, err = mgr.Resolve("invalid")
	if err == nil {
		t.Error("Resolve invalid should fail")
	}
}

func TestUpdate(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	mgr := NewManager(store)
	mgr.Add(Contact{Name: "Alice", BitcoinAddress: "old"})
	mgr.Add(Contact{Name: "Alice", BitcoinAddress: "new"})

	got, _ := mgr.Get("Alice")
	if got.BitcoinAddress != "new" {
		t.Errorf("update failed: got %s, want new", got.BitcoinAddress)
	}
}
