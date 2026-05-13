package nostr

import (
	"testing"

	"github.com/tyler-smith/go-bip39"
)

func TestDeriveFromSeed(t *testing.T) {
	// Gerar seed determinístico para teste
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed := bip39.NewSeed(mnemonic, "")

	keys, err := DeriveFromSeed(seed, 0, 0)
	if err != nil {
		t.Fatalf("DeriveFromSeed: %v", err)
	}

	// Verificar que as chaves foram geradas
	if keys.PrivateKey == "" {
		t.Fatal("PrivateKey is empty")
	}
	if keys.PublicKey == "" {
		t.Fatal("PublicKey is empty")
	}
	if len(keys.PrivateKey) != 64 {
		t.Errorf("PrivateKey hex len = %d, want 64", len(keys.PrivateKey))
	}
	if len(keys.PublicKey) != 64 {
		t.Errorf("PublicKey hex len = %d, want 64", len(keys.PublicKey))
	}
}

func TestDeriveConsistency(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed := bip39.NewSeed(mnemonic, "")

	keys1, err := DeriveFromSeed(seed, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	keys2, err := DeriveFromSeed(seed, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	// Mesmo seed + path = mesmas chaves
	if keys1.PublicKey != keys2.PublicKey {
		t.Error("mesma derivação produziu chaves diferentes")
	}
}

func TestDifferentAccounts(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed := bip39.NewSeed(mnemonic, "")

	keys0, _ := DeriveFromSeed(seed, 0, 0)
	keys1, _ := DeriveFromSeed(seed, 1, 0)

	if keys0.PublicKey == keys1.PublicKey {
		t.Error("accounts diferentes produziram mesma chave")
	}
}

func TestNpubNsec(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed := bip39.NewSeed(mnemonic, "")

	keys, _ := DeriveFromSeed(seed, 0, 0)

	npub := keys.Npub()
	if len(npub) < 10 || npub[:4] != "npub" {
		t.Errorf("Npub() = %s, want npub1...", npub)
	}

	nsec := keys.Nsec()
	if len(nsec) < 10 || nsec[:4] != "nsec" {
		t.Errorf("Nsec() = %s, want nsec1...", nsec)
	}
}

func TestFingerprint(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed := bip39.NewSeed(mnemonic, "")

	keys, _ := DeriveFromSeed(seed, 0, 0)

	fp := keys.Fingerprint()
	if len(fp) != 8 { // 4 bytes = 8 hex chars
		t.Errorf("Fingerprint() len = %d, want 8", len(fp))
	}
}

func TestCreateNoteEvent(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed := bip39.NewSeed(mnemonic, "")

	keys, _ := DeriveFromSeed(seed, 0, 0)

	event, err := keys.CreateNoteEvent("hello from crom-pay")
	if err != nil {
		t.Fatalf("CreateNoteEvent: %v", err)
	}

	if event.Content != "hello from crom-pay" {
		t.Errorf("Content = %s, want 'hello from crom-pay'", event.Content)
	}
	if event.Kind != 1 {
		t.Errorf("Kind = %d, want 1", event.Kind)
	}
	if event.PubKey != keys.PublicKey {
		t.Error("PubKey doesn't match")
	}
	if event.ID == "" {
		t.Error("Event ID is empty")
	}
	if event.Sig == "" {
		t.Error("Event Sig is empty")
	}
}

func TestDefaultRelays(t *testing.T) {
	relays := DefaultRelays()
	if len(relays) == 0 {
		t.Fatal("DefaultRelays() returned empty")
	}
	for _, r := range relays {
		if r[:6] != "wss://" {
			t.Errorf("relay %s doesn't start with wss://", r)
		}
	}
}

func TestNewRelayPool(t *testing.T) {
	// Com URLs custom
	pool := NewRelayPool([]string{"wss://test.relay"})
	if len(pool.URLs) != 1 {
		t.Errorf("URLs len = %d, want 1", len(pool.URLs))
	}

	// Sem URLs = usa defaults
	pool2 := NewRelayPool(nil)
	if len(pool2.URLs) == 0 {
		t.Error("empty URLs should use defaults")
	}
}
