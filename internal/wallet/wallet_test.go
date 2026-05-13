package wallet

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ─── Testes de Seed (BIP-39) ──────────────────────────────────────────────────

func TestGenerateMnemonic(t *testing.T) {
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		t.Fatalf("GenerateMnemonic() erro: %v", err)
	}

	words := strings.Split(mnemonic, " ")
	if len(words) != 24 {
		t.Errorf("esperava 24 palavras, recebeu %d", len(words))
	}

	if !ValidateMnemonic(mnemonic) {
		t.Error("mnemônico gerado falhou na validação")
	}
}

func TestGenerateMnemonic_Unique(t *testing.T) {
	m1, _ := GenerateMnemonic()
	m2, _ := GenerateMnemonic()

	if m1 == m2 {
		t.Error("dois mnemônicos gerados são idênticos — falha de entropia!")
	}
}

func TestValidateMnemonic_Valid(t *testing.T) {
	// Mnemônico de teste conhecido (BIP-39 test vector)
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	if !ValidateMnemonic(mnemonic) {
		t.Error("mnemônico válido rejeitado pela validação")
	}
}

func TestValidateMnemonic_Invalid(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
	}{
		{"vazio", ""},
		{"palavras erradas", "foo bar baz qux"},
		{"checksum errado", "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ValidateMnemonic(tt.mnemonic) {
				t.Errorf("mnemônico inválido '%s' passou na validação", tt.name)
			}
		})
	}
}

func TestMnemonicToSeed(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	seed, err := MnemonicToSeed(mnemonic, "")
	if err != nil {
		t.Fatalf("MnemonicToSeed() erro: %v", err)
	}

	if len(seed) != 64 { // 512 bits = 64 bytes
		t.Errorf("esperava seed de 64 bytes, recebeu %d", len(seed))
	}
}

func TestMnemonicToSeed_InvalidMnemonic(t *testing.T) {
	_, err := MnemonicToSeed("invalid mnemonic words", "")
	if err == nil {
		t.Error("esperava erro para mnemônico inválido")
	}
}

func TestMnemonicToSeed_Deterministic(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	seed1, _ := MnemonicToSeed(mnemonic, "")
	seed2, _ := MnemonicToSeed(mnemonic, "")

	if string(seed1) != string(seed2) {
		t.Error("mesma semente deve gerar seed idêntico")
	}
}

// ─── Testes de Keychain (BIP-32/84) ──────────────────────────────────────────

func TestNewKeychain(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed, _ := MnemonicToSeed(mnemonic, "")

	kc, err := NewKeychain(seed, "testnet")
	if err != nil {
		t.Fatalf("NewKeychain() erro: %v", err)
	}

	if kc.Network() != "testnet" {
		t.Errorf("esperava testnet, recebeu %s", kc.Network())
	}
}

func TestNewKeychain_InvalidSeed(t *testing.T) {
	_, err := NewKeychain([]byte("short"), "testnet")
	if err == nil {
		t.Error("esperava erro para seed curta")
	}
}

func TestNewKeychain_InvalidNetwork(t *testing.T) {
	seed := make([]byte, 64)
	_, err := NewKeychain(seed, "foonet")
	if err == nil {
		t.Error("esperava erro para rede inválida")
	}
}

func TestDeriveAddress_Testnet(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed, _ := MnemonicToSeed(mnemonic, "")

	kc, _ := NewKeychain(seed, "testnet")
	addr, err := kc.DeriveAddress(0)
	if err != nil {
		t.Fatalf("DeriveAddress() erro: %v", err)
	}

	if !strings.HasPrefix(addr, "tb1q") {
		t.Errorf("endereço testnet deve começar com 'tb1q', recebeu '%s'", addr)
	}
}

func TestDeriveAddress_Mainnet(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed, _ := MnemonicToSeed(mnemonic, "")

	kc, _ := NewKeychain(seed, "mainnet")
	addr, err := kc.DeriveAddress(0)
	if err != nil {
		t.Fatalf("DeriveAddress() erro: %v", err)
	}

	if !strings.HasPrefix(addr, "bc1q") {
		t.Errorf("endereço mainnet deve começar com 'bc1q', recebeu '%s'", addr)
	}
}

func TestDeriveAddress_Deterministic(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed, _ := MnemonicToSeed(mnemonic, "")

	kc1, _ := NewKeychain(seed, "testnet")
	kc2, _ := NewKeychain(seed, "testnet")

	addr1, _ := kc1.DeriveAddress(0)
	addr2, _ := kc2.DeriveAddress(0)

	if addr1 != addr2 {
		t.Error("mesmo seed deve gerar mesmo endereço")
	}
}

func TestDeriveAddress_DifferentIndexes(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed, _ := MnemonicToSeed(mnemonic, "")

	kc, _ := NewKeychain(seed, "testnet")
	addr0, _ := kc.DeriveAddress(0)
	addr1, _ := kc.DeriveAddress(1)

	if addr0 == addr1 {
		t.Error("índices diferentes devem gerar endereços diferentes")
	}
}

// ─── Testes de Crypto (AES-256-GCM) ──────────────────────────────────────────

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	plaintext := []byte("dados secretos da carteira Bitcoin")
	password := "minha-senha-forte-123!"

	encrypted, err := EncryptData(plaintext, password)
	if err != nil {
		t.Fatalf("EncryptData() erro: %v", err)
	}

	decrypted, err := DecryptData(encrypted, password)
	if err != nil {
		t.Fatalf("DecryptData() erro: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Error("dados decifrados não coincidem com os originais")
	}
}

func TestDecrypt_WrongPassword(t *testing.T) {
	plaintext := []byte("segredo")
	encrypted, _ := EncryptData(plaintext, "senha-correta")

	_, err := DecryptData(encrypted, "senha-errada")
	if err == nil {
		t.Error("esperava erro com senha errada")
	}
}

func TestEncrypt_EmptyPassword(t *testing.T) {
	_, err := EncryptData([]byte("dados"), "")
	if err == nil {
		t.Error("esperava erro com senha vazia")
	}
}

func TestEncrypt_DifferentOutputs(t *testing.T) {
	plaintext := []byte("mesmos dados")
	password := "mesma-senha"

	enc1, _ := EncryptData(plaintext, password)
	enc2, _ := EncryptData(plaintext, password)

	if string(enc1) == string(enc2) {
		t.Error("cifras devem ser diferentes por causa do salt/nonce aleatórios")
	}
}

// ─── Testes de Wallet (integração) ───────────────────────────────────────────

func TestWallet_CreateAndOpen(t *testing.T) {
	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "test-wallet")

	// Criar carteira
	mnemonic, info, err := Create(dataDir, "test-password-123", "testnet")
	if err != nil {
		t.Fatalf("Create() erro: %v", err)
	}

	if mnemonic == "" {
		t.Error("mnemônico não deve ser vazio")
	}

	if info.Address == "" {
		t.Error("endereço não deve ser vazio")
	}

	if !strings.HasPrefix(info.Address, "tb1q") {
		t.Errorf("endereço testnet deve começar com tb1q, recebeu: %s", info.Address)
	}

	// Abrir carteira
	w, err := Open(dataDir, "test-password-123")
	if err != nil {
		t.Fatalf("Open() erro: %v", err)
	}
	defer w.Close()

	walletInfo, err := w.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo() erro: %v", err)
	}

	// Endereço deve ser o mesmo
	if walletInfo.Address != info.Address {
		t.Errorf("endereço diferente após reabrir: '%s' vs '%s'", walletInfo.Address, info.Address)
	}
}

func TestWallet_CreateDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "test-wallet")

	// Criar primeira vez
	_, _, err := Create(dataDir, "pass123456", "testnet")
	if err != nil {
		t.Fatalf("Create() erro: %v", err)
	}

	// Tentar criar de novo
	_, _, err = Create(dataDir, "pass123456", "testnet")
	if err == nil {
		t.Error("esperava erro ao criar carteira duplicada")
	}
}

func TestWallet_OpenWrongPassword(t *testing.T) {
	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "test-wallet")

	Create(dataDir, "senha-correta", "testnet")

	_, err := Open(dataDir, "senha-errada")
	if err == nil {
		t.Error("esperava erro com senha errada")
	}
}

func TestWallet_OpenNonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "nao-existe")

	_, err := Open(dataDir, "qualquer")
	if err == nil {
		t.Error("esperava erro ao abrir carteira inexistente")
	}
}

func TestWallet_RestoreAndVerify(t *testing.T) {
	tmpDir := t.TempDir()

	// Criar carteira original
	dataDir1 := filepath.Join(tmpDir, "wallet-original")
	mnemonic, originalInfo, err := Create(dataDir1, "pass123456", "testnet")
	if err != nil {
		t.Fatalf("Create() erro: %v", err)
	}

	// Restaurar com mesma semente em diretório diferente
	dataDir2 := filepath.Join(tmpDir, "wallet-restored")
	restoredInfo, err := Restore(dataDir2, mnemonic, "outra-senha1", "testnet")
	if err != nil {
		t.Fatalf("Restore() erro: %v", err)
	}

	// Endereços devem ser idênticos (mesma semente = mesmo endereço)
	if originalInfo.Address != restoredInfo.Address {
		t.Errorf("endereço restaurado diferente: '%s' vs '%s'",
			originalInfo.Address, restoredInfo.Address)
	}
}

// ─── Limpeza ──────────────────────────────────────────────────────────────────

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
