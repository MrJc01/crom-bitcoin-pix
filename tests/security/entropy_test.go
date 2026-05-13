//go:build security

package security

import (
	"bytes"
	"crypto/rand"
	"math"
	"testing"

	"github.com/MrJc01/crom-bitcoin-pix/internal/wallet"
)

// ─── Teste de Qualidade de Entropia ──────────────────────────────────────────
//
// Verifica se a entropia gerada pelo sistema é criptograficamente adequada.
// Um CSPRNG ruim pode gerar carteiras previsíveis = roubo de fundos.

func TestEntropy_RandomnessQuality(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando qualidade da entropia criptográfica")

	const samples = 1000
	const size = 32 // 256 bits

	allSamples := make([][]byte, samples)
	for i := 0; i < samples; i++ {
		buf := make([]byte, size)
		_, err := rand.Read(buf)
		if err != nil {
			t.Fatalf("crypto/rand.Read falhou: %v", err)
		}
		allSamples[i] = buf
	}

	// Verificar unicidade — nenhum par deve ser igual
	for i := 0; i < samples; i++ {
		for j := i + 1; j < samples; j++ {
			if bytes.Equal(allSamples[i], allSamples[j]) {
				t.Fatalf("🔴 CRÍTICO: Duas amostras de entropia IDÊNTICAS (i=%d, j=%d)!", i, j)
			}
		}
	}
	t.Log("✅ 1000 amostras de entropia são todas únicas")
}

func TestEntropy_ShannonEntropy(t *testing.T) {
	t.Log("🔍 PENTEST: Calculando Shannon Entropy do CSPRNG")

	buf := make([]byte, 10000)
	rand.Read(buf)

	// Calcular frequência de cada byte
	freq := make(map[byte]float64)
	for _, b := range buf {
		freq[b]++
	}

	// Shannon entropy: H = -Σ p(x) * log2(p(x))
	var entropy float64
	n := float64(len(buf))
	for _, count := range freq {
		p := count / n
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}

	// Entropia máxima para bytes = 8.0 bits
	// Um bom CSPRNG deve estar acima de 7.9
	t.Logf("📊 Shannon Entropy: %.4f bits (máximo: 8.0)", entropy)

	if entropy < 7.9 {
		t.Errorf("🔴 Entropia muito baixa: %.4f (mínimo aceitável: 7.9)", entropy)
	} else {
		t.Log("✅ Entropia dentro dos padrões criptográficos")
	}
}

func TestEntropy_MnemonicUniqueness(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando unicidade de mnemônicos gerados")

	const count = 100
	mnemonics := make(map[string]bool)

	for i := 0; i < count; i++ {
		m, err := wallet.GenerateMnemonic()
		if err != nil {
			t.Fatalf("GenerateMnemonic() falhou na iteração %d: %v", i, err)
		}
		if mnemonics[m] {
			t.Fatalf("🔴 CRÍTICO: Mnemônico duplicado na iteração %d!", i)
		}
		mnemonics[m] = true
	}

	t.Logf("✅ %d mnemônicos gerados, todos únicos", count)
}

func TestEntropy_NoBiasInBytes(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando bias nos bytes gerados")

	buf := make([]byte, 100000)
	rand.Read(buf)

	// Contar frequência de cada valor de byte
	counts := make([]int, 256)
	for _, b := range buf {
		counts[b]++
	}

	// Frequência esperada: 100000/256 ≈ 390.6
	expected := float64(len(buf)) / 256.0
	maxDeviation := 0.0

	for i, c := range counts {
		deviation := math.Abs(float64(c)-expected) / expected * 100
		if deviation > maxDeviation {
			maxDeviation = deviation
		}
		// Desvio > 20% indica bias sério
		if deviation > 20.0 {
			t.Errorf("🔴 Byte 0x%02x tem desvio de %.1f%% (count=%d, esperado=%.0f)", i, deviation, c, expected)
		}
	}

	t.Logf("✅ Desvio máximo: %.2f%% (limite: 20%%)", maxDeviation)
}
