package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewStore_CreateAndClose(t *testing.T) {
	dir := t.TempDir()
	store, err := NewStore(filepath.Join(dir, "test"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestStore_PutGet(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(filepath.Join(dir, "test"))
	defer store.Close()

	prefix := []byte("test:")
	key := []byte("mykey")
	value := []byte("myvalue")

	if err := store.Put(prefix, key, value); err != nil {
		t.Fatalf("Put: %v", err)
	}

	got, err := store.Get(prefix, key)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if string(got) != string(value) {
		t.Errorf("Get retornou %q, esperado %q", got, value)
	}
}

func TestStore_GetNotFound(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(filepath.Join(dir, "test"))
	defer store.Close()

	_, err := store.Get([]byte("x:"), []byte("naoexiste"))
	if err != ErrKeyNotFound {
		t.Errorf("Esperado ErrKeyNotFound, obtido: %v", err)
	}
}

func TestStore_Has(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(filepath.Join(dir, "test"))
	defer store.Close()

	prefix := []byte("has:")
	key := []byte("k1")

	if store.Has(prefix, key) {
		t.Error("Has retornou true para chave inexistente")
	}

	store.Put(prefix, key, []byte("v1"))

	if !store.Has(prefix, key) {
		t.Error("Has retornou false para chave existente")
	}
}

func TestStore_Delete(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(filepath.Join(dir, "test"))
	defer store.Close()

	prefix := []byte("del:")
	key := []byte("k1")
	store.Put(prefix, key, []byte("val"))

	if err := store.Delete(prefix, key); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if store.Has(prefix, key) {
		t.Error("Chave ainda existe após Delete")
	}
}

func TestStore_Overwrite(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(filepath.Join(dir, "test"))
	defer store.Close()

	prefix := []byte("ow:")
	key := []byte("k1")
	store.Put(prefix, key, []byte("v1"))
	store.Put(prefix, key, []byte("v2"))

	got, _ := store.Get(prefix, key)
	if string(got) != "v2" {
		t.Errorf("Overwrite falhou: obtido %q, esperado %q", got, "v2")
	}
}

func TestStore_PrefixIsolation(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(filepath.Join(dir, "test"))
	defer store.Close()

	store.Put([]byte("a:"), []byte("key"), []byte("value-a"))
	store.Put([]byte("b:"), []byte("key"), []byte("value-b"))

	gotA, _ := store.Get([]byte("a:"), []byte("key"))
	gotB, _ := store.Get([]byte("b:"), []byte("key"))

	if string(gotA) != "value-a" {
		t.Errorf("Prefix A corrompido: %q", gotA)
	}
	if string(gotB) != "value-b" {
		t.Errorf("Prefix B corrompido: %q", gotB)
	}
}

func TestStore_WalletSeedRoundtrip(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(filepath.Join(dir, "test"))
	defer store.Close()

	seed := []byte("fake-encrypted-seed-data-64-bytes-for-testing-purposes-here!!!")

	if store.HasWallet() {
		t.Fatal("HasWallet true antes de salvar")
	}

	if err := store.SaveEncryptedSeed(seed); err != nil {
		t.Fatalf("SaveEncryptedSeed: %v", err)
	}

	if !store.HasWallet() {
		t.Fatal("HasWallet false após salvar")
	}

	loaded, err := store.LoadEncryptedSeed()
	if err != nil {
		t.Fatalf("LoadEncryptedSeed: %v", err)
	}

	if string(loaded) != string(seed) {
		t.Errorf("Seed corrompida: obtido %q", loaded)
	}
}

func TestStore_ConfigRoundtrip(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(filepath.Join(dir, "test"))
	defer store.Close()

	store.SaveConfig("network", "testnet")
	store.SaveConfig("version", "1.0.0")

	net, _ := store.LoadConfig("network")
	ver, _ := store.LoadConfig("version")

	if net != "testnet" {
		t.Errorf("network: %q", net)
	}
	if ver != "1.0.0" {
		t.Errorf("version: %q", ver)
	}

	_, err := store.LoadConfig("inexistente")
	if err == nil {
		t.Error("LoadConfig deveria falhar para chave inexistente")
	}
}

func TestStore_ClosedStoreErrors(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewStore(filepath.Join(dir, "test"))
	store.Close()
	store.db = nil // forçar estado fechado

	if err := store.Put([]byte("x:"), []byte("k"), []byte("v")); err != ErrStoreNotOpen {
		t.Errorf("Put em store fechado: %v", err)
	}
	if _, err := store.Get([]byte("x:"), []byte("k")); err != ErrStoreNotOpen {
		t.Errorf("Get em store fechado: %v", err)
	}
	if err := store.Delete([]byte("x:"), []byte("k")); err != ErrStoreNotOpen {
		t.Errorf("Delete em store fechado: %v", err)
	}
}

func TestStore_Persistence(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "persist")

	// Escrever
	s1, _ := NewStore(dbPath)
	s1.Put([]byte("p:"), []byte("k"), []byte("persistent"))
	s1.Close()

	// Reabrir e ler
	s2, _ := NewStore(dbPath)
	defer s2.Close()

	got, err := s2.Get([]byte("p:"), []byte("k"))
	if err != nil {
		t.Fatalf("Get após reopen: %v", err)
	}
	if string(got) != "persistent" {
		t.Errorf("Dado não persistiu: %q", got)
	}
}

func TestStore_DirectoryPermissions(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "perms", "wallet.db")

	NewStore(filepath.Join(dir, "perms"))

	info, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}

	mode := info.Mode().Perm()
	if mode&0077 != 0 {
		t.Errorf("Permissões muito abertas: %o (esperado 700)", mode)
	}
}
