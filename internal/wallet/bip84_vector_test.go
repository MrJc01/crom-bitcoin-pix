package wallet

import (
	"fmt"
	"testing"
)

// TestBIP84_OfficialVector verifica o endereço contra o vetor oficial da spec BIP-84.
// Referência: https://github.com/bitcoin/bips/blob/master/bip-0084.mediawiki
func TestBIP84_OfficialVector(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	seed, err := MnemonicToSeed(mnemonic, "")
	if err != nil {
		t.Fatalf("MnemonicToSeed: %v", err)
	}

	// Mainnet
	kc, err := NewKeychain(seed, "mainnet")
	if err != nil {
		t.Fatalf("NewKeychain mainnet: %v", err)
	}

	addr0, _ := kc.DeriveAddress(0)
	addr1, _ := kc.DeriveAddress(1)

	// Vetor BIP-84 spec oficial:
	// Account 0, first receiving = bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu
	expected0 := "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu"

	t.Logf("Addr[0]: %s", addr0)
	t.Logf("Addr[1]: %s", addr1)

	if addr0 != expected0 {
		t.Fatalf("❌ BIP-84 MISMATCH!\n  Esperado: %s\n  Obtido:   %s", expected0, addr0)
	}
	t.Log("✅ Addr[0] confere com BIP-84 spec oficial")

	// Testnet
	kcT, _ := NewKeychain(seed, "testnet")
	addrT, _ := kcT.DeriveAddress(0)
	t.Logf("Testnet[0]: %s", addrT)

	expectedT := "tb1q6rz28mcfaxtmd6v789l9rrlrusdprr9pqcpvkl"
	if addrT != expectedT {
		t.Fatalf("❌ Testnet MISMATCH!\n  Esperado: %s\n  Obtido:   %s", expectedT, addrT)
	}
	t.Log("✅ Testnet[0] confere com derivação BIP-84")

	fmt.Println() // silent linter
}
