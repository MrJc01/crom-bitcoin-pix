//go:build stress

package stress

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/MrJc01/crom-bitcoin-pix/internal/wallet"
)

// ─── Testes de Stress: Carga e Limites ───────────────────────────────────────

func TestStress_ConcurrentWalletCreation(t *testing.T) {
	t.Log("💪 STRESS: Criando 20 carteiras simultaneamente")

	tmpDir := t.TempDir()
	var wg sync.WaitGroup
	errors := make(chan error, 20)

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			dir := filepath.Join(tmpDir, fmt.Sprintf("wallet-%d", idx))
			_, _, err := wallet.Create(dir, fmt.Sprintf("password-%d!!", idx), "testnet")
			if err != nil {
				errors <- fmt.Errorf("wallet %d: %v", idx, err)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	errCount := 0
	for err := range errors {
		t.Errorf("🔴 %v", err)
		errCount++
	}

	if errCount == 0 {
		t.Log("✅ 20 carteiras criadas simultaneamente sem erros")
	}
}

func TestStress_RapidOpenClose(t *testing.T) {
	t.Log("💪 STRESS: 50 ciclos open/close rápidos")

	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "rapid-test")
	password := "rapid-test-pw12"

	wallet.Create(dataDir, password, "testnet")

	for i := 0; i < 50; i++ {
		w, err := wallet.Open(dataDir, password)
		if err != nil {
			t.Fatalf("Open falhou na iteração %d: %v", i, err)
		}
		_, err = w.GetInfo()
		if err != nil {
			t.Fatalf("GetInfo falhou na iteração %d: %v", i, err)
		}
		w.Close()
	}

	t.Log("✅ 50 ciclos open/close sem erro")
}

func TestStress_ManyAddressDerivations(t *testing.T) {
	t.Log("💪 STRESS: Derivando 1000 endereços do mesmo wallet")

	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "many-addr")
	password := "many-addr-test1"

	wallet.Create(dataDir, password, "testnet")
	w, _ := wallet.Open(dataDir, password)
	defer w.Close()

	addresses := make(map[string]bool)
	for i := uint32(0); i < 1000; i++ {
		addr, err := w.GetAddress(i)
		if err != nil {
			t.Fatalf("GetAddress(%d): %v", i, err)
		}
		if addresses[addr] {
			t.Fatalf("🔴 Colisão de endereço no índice %d!", i)
		}
		addresses[addr] = true
	}

	t.Logf("✅ 1000 endereços únicos derivados (total: %d)", len(addresses))
}

func TestStress_EncryptDecryptBatch(t *testing.T) {
	t.Log("💪 STRESS: 10 ciclos encrypt/decrypt (Argon2id ~3s/ciclo)")

	data := []byte("dados de teste para stress de cifra AES-256-GCM com Argon2id")
	password := "stress-password1"

	for i := 0; i < 10; i++ {
		enc, err := wallet.EncryptData(data, password)
		if err != nil {
			t.Fatalf("Encrypt falhou na iteração %d: %v", i, err)
		}
		dec, err := wallet.DecryptData(enc, password)
		if err != nil {
			t.Fatalf("Decrypt falhou na iteração %d: %v", i, err)
		}
		if string(dec) != string(data) {
			t.Fatalf("Dados corrompidos na iteração %d!", i)
		}
	}

	t.Log("✅ 10 ciclos encrypt/decrypt sem corrupção")
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
