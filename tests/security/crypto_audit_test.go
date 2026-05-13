//go:build security

package security

import (
	"testing"
	"time"

	"github.com/MrJc01/crom-bitcoin-pix/internal/wallet"
)

// ─── Auditoria da Cifra AES-256-GCM + Argon2id ──────────────────────────────
//
// Verifica propriedades criptográficas essenciais da implementação.

func TestCrypto_CiphertextDiffers(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando que cifras com mesma entrada produzem outputs diferentes")

	plaintext := []byte("semente secreta de 64 bytes simulando uma seed BIP-39 real!!!!")
	password := "senha-de-teste-123"

	enc1, _ := wallet.EncryptData(plaintext, password)
	enc2, _ := wallet.EncryptData(plaintext, password)

	if string(enc1) == string(enc2) {
		t.Fatal("🔴 CRÍTICO: Duas cifras com mesma entrada são IDÊNTICAS! Salt/nonce não está aleatório.")
	}

	t.Log("✅ Salt/nonce aleatórios garantem outputs diferentes")
}

func TestCrypto_TamperDetection(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando detecção de adulteração (tamper detection)")

	plaintext := []byte("dados protegidos por AES-GCM autenticado")
	password := "minha-senha"

	encrypted, _ := wallet.EncryptData(plaintext, password)

	// Adulterar 1 byte no meio do ciphertext
	tampered := make([]byte, len(encrypted))
	copy(tampered, encrypted)
	tampered[len(tampered)/2] ^= 0xFF // flip bits

	_, err := wallet.DecryptData(tampered, password)
	if err == nil {
		t.Fatal("🔴 CRÍTICO: Dados adulterados foram aceitos pela decifra! GCM auth tag não está funcionando.")
	}

	t.Log("✅ AES-GCM detectou adulteração corretamente")
}

func TestCrypto_TruncationDetection(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando detecção de truncamento")

	plaintext := []byte("dados que não podem ser cortados")
	password := "senha-teste"

	encrypted, _ := wallet.EncryptData(plaintext, password)

	// Truncar o ciphertext
	truncated := encrypted[:len(encrypted)-10]

	_, err := wallet.DecryptData(truncated, password)
	if err == nil {
		t.Fatal("🔴 Dados truncados foram aceitos!")
	}

	t.Log("✅ Truncamento detectado corretamente")
}

func TestCrypto_WrongPasswordTiming(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando resistência a timing attack na senha")

	plaintext := []byte("dados sensíveis")
	encrypted, _ := wallet.EncryptData(plaintext, "senha-correta-longa-123")

	// Medir tempo com senha errada curta
	start1 := time.Now()
	wallet.DecryptData(encrypted, "a")
	dur1 := time.Since(start1)

	// Medir tempo com senha errada longa
	start2 := time.Now()
	wallet.DecryptData(encrypted, "senha-completamente-errada-e-muito-longa-demais-para-tudo")
	dur2 := time.Since(start2)

	// Argon2id deve dominar o tempo — senhas diferentes devem ter tempo similar
	ratio := float64(dur1) / float64(dur2)
	t.Logf("📊 Tempo senha curta: %v | Tempo senha longa: %v | Ratio: %.2f", dur1, dur2, ratio)

	// Ratio entre 0.5 e 2.0 indica que Argon2id domina (bom)
	if ratio < 0.3 || ratio > 3.0 {
		t.Logf("⚠️ Ratio de timing fora do ideal: %.2f (pode indicar leak de informação)", ratio)
	} else {
		t.Log("✅ Timing consistente — Argon2id domina o tempo de processamento")
	}
}

func TestCrypto_Argon2idSlowness(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando que Argon2id é suficientemente lento (anti-bruteforce)")

	plaintext := []byte("dados para teste de velocidade")
	password := "senha-teste-123"

	start := time.Now()
	wallet.EncryptData(plaintext, password)
	duration := time.Since(start)

	t.Logf("📊 Tempo de cifra (inclui Argon2id): %v", duration)

	// Argon2id com 64MB RAM e 3 iterações deve levar pelo menos 100ms
	if duration < 100*time.Millisecond {
		t.Errorf("🔴 Cifra muito rápida (%v) — Argon2id pode estar mal configurado!", duration)
	} else {
		t.Logf("✅ Tempo adequado para resistir a brute-force (>100ms por tentativa)")
	}

	// Calcular tentativas por segundo
	attemptsPerSec := float64(time.Second) / float64(duration)
	t.Logf("📊 Brute-force rate: ~%.0f tentativas/segundo", attemptsPerSec)

	if attemptsPerSec > 100 {
		t.Errorf("🔴 Rate muito alto! Atacante pode tentar %.0f senhas/segundo", attemptsPerSec)
	}
}

func TestCrypto_EmptyPasswordRejected(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando rejeição de senha vazia")

	_, err := wallet.EncryptData([]byte("dados"), "")
	if err == nil {
		t.Fatal("🔴 Senha vazia foi aceita pela cifra!")
	}

	t.Log("✅ Senha vazia rejeitada corretamente")
}

func TestCrypto_LargePayload(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando cifra com payload grande")

	// Simular seed + metadata grande (1MB)
	payload := make([]byte, 1024*1024)
	for i := range payload {
		payload[i] = byte(i % 256)
	}

	encrypted, err := wallet.EncryptData(payload, "senha-forte-123")
	if err != nil {
		t.Fatalf("Falha ao cifrar payload grande: %v", err)
	}

	decrypted, err := wallet.DecryptData(encrypted, "senha-forte-123")
	if err != nil {
		t.Fatalf("Falha ao decifrar payload grande: %v", err)
	}

	if string(decrypted) != string(payload) {
		t.Fatal("🔴 Payload grande corrompido após cifra/decifra!")
	}

	t.Logf("✅ Payload de 1MB cifrado/decifrado corretamente (%d bytes cifrados)", len(encrypted))
}
