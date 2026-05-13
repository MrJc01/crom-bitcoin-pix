package wallet

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/tyler-smith/go-bip39"
)

const (
	// EntropyBits é o tamanho da entropia para gerar 24 palavras (256 bits)
	EntropyBits = 256
	// EntropyBytes é EntropyBits convertido para bytes
	EntropyBytes = EntropyBits / 8
)

var (
	ErrInvalidMnemonic = errors.New("mnemônico inválido: checksum ou palavras incorretas")
	ErrEntropyFailed   = errors.New("falha ao gerar entropia criptográfica")
)

// GenerateMnemonic gera um mnemônico BIP-39 de 24 palavras usando entropia
// criptograficamente segura do sistema operacional (crypto/rand).
//
// SEGURANÇA: Usa CSPRNG do OS — NUNCA math/rand.
func GenerateMnemonic() (string, error) {
	entropy, err := generateEntropy(EntropyBytes)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrEntropyFailed, err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("falha ao criar mnemônico: %w", err)
	}

	return mnemonic, nil
}

// ValidateMnemonic verifica se um mnemônico BIP-39 é válido
// (palavras corretas + checksum válido).
func ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

// MnemonicToSeed converte um mnemônico BIP-39 em seed de 512 bits
// usando PBKDF2-HMAC-SHA512 com 2048 iterações.
//
// O passphrase é opcional (string vazia = sem passphrase adicional).
// A seed resultante é a raiz para toda derivação de chaves (BIP-32).
func MnemonicToSeed(mnemonic, passphrase string) ([]byte, error) {
	if !ValidateMnemonic(mnemonic) {
		return nil, ErrInvalidMnemonic
	}

	seed := bip39.NewSeed(mnemonic, passphrase)
	return seed, nil
}

// generateEntropy gera bytes criptograficamente seguros via crypto/rand.
func generateEntropy(numBytes int) ([]byte, error) {
	entropy := make([]byte, numBytes)
	_, err := rand.Read(entropy)
	if err != nil {
		return nil, err
	}
	return entropy, nil
}
