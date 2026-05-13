//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MrJc01/crom-bitcoin-pix/internal/wallet"
)

// ─── Testes de Integração: Fluxo Completo ────────────────────────────────────

func TestFlow_CreateOpenClose(t *testing.T) {
	t.Log("🔄 INTEGRAÇÃO: Create → Open → GetInfo → Close")

	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "int-wallet")
	password := "integration-test-pw"

	// Create
	mnemonic, info, err := wallet.Create(dataDir, password, "testnet")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	t.Logf("  Mnemônico: %s...", mnemonic[:30])
	t.Logf("  Endereço: %s", info.Address)

	// Open
	w, err := wallet.Open(dataDir, password)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	// GetInfo
	info2, err := w.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo: %v", err)
	}

	if info.Address != info2.Address {
		t.Errorf("Endereço diferente: create=%s open=%s", info.Address, info2.Address)
	}

	// Close
	if err := w.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}

	t.Log("✅ Fluxo create→open→info→close OK")
}

func TestFlow_CreateRestoreMatch(t *testing.T) {
	t.Log("🔄 INTEGRAÇÃO: Create em Dir-A → Restore em Dir-B → Endereços devem coincidir")

	tmpDir := t.TempDir()

	// Create original
	dirA := filepath.Join(tmpDir, "wallet-a")
	mnemonic, infoA, err := wallet.Create(dirA, "pw-original1", "testnet")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	// Restore em diretório diferente
	dirB := filepath.Join(tmpDir, "wallet-b")
	infoB, err := wallet.Restore(dirB, mnemonic, "pw-diferente", "testnet")
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}

	if infoA.Address != infoB.Address {
		t.Fatalf("🔴 Endereços diferentes! A=%s B=%s", infoA.Address, infoB.Address)
	}

	t.Logf("✅ Endereço idêntico: %s", infoA.Address)
}

func TestFlow_MultipleAddresses(t *testing.T) {
	t.Log("🔄 INTEGRAÇÃO: Derivar múltiplos endereços do mesmo wallet")

	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "multi-addr")

	wallet.Create(dataDir, "test-pass-123", "testnet")
	w, _ := wallet.Open(dataDir, "test-pass-123")
	defer w.Close()

	addresses := make(map[string]bool)
	for i := uint32(0); i < 10; i++ {
		addr, err := w.GetAddress(i)
		if err != nil {
			t.Fatalf("GetAddress(%d): %v", i, err)
		}

		if !strings.HasPrefix(addr, "tb1q") {
			t.Errorf("Endereço %d não é SegWit testnet: %s", i, addr)
		}

		if addresses[addr] {
			t.Fatalf("🔴 Endereço duplicado no índice %d: %s", i, addr)
		}
		addresses[addr] = true
		t.Logf("  [%d] %s", i, addr)
	}

	t.Logf("✅ 10 endereços únicos derivados")
}

func TestFlow_NetworkIsolation(t *testing.T) {
	t.Log("🔄 INTEGRAÇÃO: Mesma seed em redes diferentes → endereços diferentes")

	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	tmpDir := t.TempDir()

	// Mainnet
	dirMain := filepath.Join(tmpDir, "mainnet")
	infoMain, _ := wallet.Restore(dirMain, mnemonic, "pw-main1234", "mainnet")

	// Testnet
	dirTest := filepath.Join(tmpDir, "testnet")
	infoTest, _ := wallet.Restore(dirTest, mnemonic, "pw-test1234", "testnet")

	if infoMain.Address == infoTest.Address {
		t.Fatal("🔴 Endereços mainnet e testnet são IGUAIS! Derivação está errada.")
	}

	if !strings.HasPrefix(infoMain.Address, "bc1q") {
		t.Errorf("Mainnet deveria começar com bc1q: %s", infoMain.Address)
	}
	if !strings.HasPrefix(infoTest.Address, "tb1q") {
		t.Errorf("Testnet deveria começar com tb1q: %s", infoTest.Address)
	}

	t.Logf("✅ Mainnet: %s | Testnet: %s", infoMain.Address, infoTest.Address)
}

func TestFlow_DuplicateCreateFails(t *testing.T) {
	t.Log("🔄 INTEGRAÇÃO: Segundo create no mesmo dir deve falhar")

	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "dup-test")

	_, _, err := wallet.Create(dataDir, "first-pass1", "testnet")
	if err != nil {
		t.Fatalf("Primeiro create falhou: %v", err)
	}

	_, _, err = wallet.Create(dataDir, "second-pas1", "testnet")
	if err == nil {
		t.Fatal("🔴 Segundo create deveria falhar!")
	}

	t.Logf("✅ Segundo create rejeitado: %v", err)
}

func TestFlow_DataPersistence(t *testing.T) {
	t.Log("🔄 INTEGRAÇÃO: Dados persistem entre Open/Close")

	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "persist-test")

	// Criar
	_, info1, _ := wallet.Create(dataDir, "persist-pass", "testnet")

	// Abrir e fechar 5 vezes
	for i := 0; i < 5; i++ {
		w, err := wallet.Open(dataDir, "persist-pass")
		if err != nil {
			t.Fatalf("Open iteração %d: %v", i, err)
		}
		infoN, _ := w.GetInfo()
		if infoN.Address != info1.Address {
			t.Fatalf("Endereço mudou na iteração %d!", i)
		}
		w.Close()
	}

	t.Log("✅ Dados persistentes após 5 ciclos open/close")
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
