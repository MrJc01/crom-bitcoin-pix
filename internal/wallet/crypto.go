package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"runtime"

	"golang.org/x/crypto/argon2"
)

// Parâmetros Argon2id para derivação de chave de cifra.
// Escolhidos para resistir a ataques de brute-force em hardware moderno.
const (
	argonTime    = 3           // iterações
	argonMemory  = 64 * 1024  // 64 MB de RAM
	argonThreads = 4          // paralelismo
	argonKeyLen  = 32         // 256 bits (AES-256)
	saltLen      = 16         // 128 bits de salt
	nonceLen     = 12         // 96 bits para AES-GCM
)

var (
	ErrEncryptFailed  = errors.New("falha ao cifrar dados")
	ErrDecryptFailed  = errors.New("falha ao decifrar dados")
	ErrPasswordEmpty  = errors.New("senha não pode ser vazia")
)

// EncryptData cifra dados usando AES-256-GCM com chave derivada via Argon2id.
//
// O formato do output é: salt (16 bytes) || nonce (12 bytes) || ciphertext
//
// SEGURANÇA:
// - Argon2id resiste a ataques de GPU e side-channel
// - AES-GCM provê confidencialidade + integridade (AEAD)
// - Salt e nonce são gerados com crypto/rand (CSPRNG)
func EncryptData(plaintext []byte, password string) ([]byte, error) {
	if password == "" {
		return nil, ErrPasswordEmpty
	}

	// Gerar salt aleatório
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("%w: salt: %v", ErrEncryptFailed, err)
	}

	// Derivar chave AES-256 com Argon2id
	key := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	// Criar cifra AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w: aes: %v", ErrEncryptFailed, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%w: gcm: %v", ErrEncryptFailed, err)
	}

	// Gerar nonce aleatório
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("%w: nonce: %v", ErrEncryptFailed, err)
	}

	// Cifrar: salt || nonce || ciphertext (com tag de autenticação)
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	result := make([]byte, 0, saltLen+len(nonce)+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	// Limpar chave da memória
	zeroBytes(key)

	return result, nil
}

// DecryptData decifra dados cifrados com EncryptData.
//
// Espera o formato: salt (16 bytes) || nonce (12 bytes) || ciphertext
func DecryptData(encrypted []byte, password string) ([]byte, error) {
	if password == "" {
		return nil, ErrPasswordEmpty
	}

	minLen := saltLen + nonceLen + 1
	if len(encrypted) < minLen {
		return nil, fmt.Errorf("%w: dados muito curtos (min %d bytes)", ErrDecryptFailed, minLen)
	}

	// Extrair salt e nonce
	salt := encrypted[:saltLen]
	nonce := encrypted[saltLen : saltLen+nonceLen]
	ciphertext := encrypted[saltLen+nonceLen:]

	// Derivar mesma chave com Argon2id
	key := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	// Criar cifra AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		zeroBytes(key)
		return nil, fmt.Errorf("%w: aes: %v", ErrDecryptFailed, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		zeroBytes(key)
		return nil, fmt.Errorf("%w: gcm: %v", ErrDecryptFailed, err)
	}

	// Decifrar e verificar integridade
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		zeroBytes(key)
		return nil, fmt.Errorf("%w: senha incorreta ou dados corrompidos", ErrDecryptFailed)
	}

	// Limpar chave da memória
	zeroBytes(key)

	return plaintext, nil
}

// zeroBytes limpa um slice de bytes da memória (zeroing).
// Previne que dados sensíveis permaneçam em RAM após uso.
// runtime.KeepAlive impede que o compilador otimize este loop.
func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
	runtime.KeepAlive(&b)
}
