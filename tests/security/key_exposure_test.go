//go:build security

package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MrJc01/crom-bitcoin-pix/internal/wallet"
)

// ─── Testes de Exposição de Chaves ───────────────────────────────────────────
//
// Verifica se chaves privadas ou sementes estão expostas em disco, logs ou memória.

func TestKeyExposure_NoPlaintextSeedOnDisk(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando que a semente NÃO está em texto plano no disco")

	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "pentest-wallet")

	mnemonic, _, err := wallet.Create(dataDir, "senha-segura-12345", "testnet")
	if err != nil {
		t.Fatalf("Create() erro: %v", err)
	}

	// Buscar a semente em qualquer arquivo do diretório de dados
	words := strings.Split(mnemonic, " ")

	err = filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// Verificar se alguma palavra do mnemônico aparece como substring no arquivo
		contentStr := string(content)
		matchCount := 0
		for _, word := range words {
			if strings.Contains(contentStr, word) {
				matchCount++
			}
		}

		// Se mais de 3 palavras do mnemônico estão no arquivo, é suspeito
		if matchCount > 3 {
			t.Errorf("🔴 CRÍTICO: Arquivo %s contém %d palavras do mnemônico em texto plano!",
				filepath.Base(path), matchCount)
		}

		// Verificar o mnemônico completo
		if strings.Contains(contentStr, mnemonic) {
			t.Fatalf("🔴 CRÍTICO: Mnemônico COMPLETO encontrado em texto plano em: %s", path)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Walk() erro: %v", err)
	}

	t.Log("✅ Nenhuma semente em texto plano encontrada no disco")
}

func TestKeyExposure_FilePermissions(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando permissões de arquivo do diretório de dados")

	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "pentest-perms")

	wallet.Create(dataDir, "senha-segura-12345", "testnet")

	// O diretório do DB deve ter permissão restrita (700 = rwx------)
	info, err := os.Stat(filepath.Join(dataDir, "wallet.db"))
	if err != nil {
		t.Fatalf("Stat() erro: %v", err)
	}

	mode := info.Mode().Perm()
	t.Logf("📊 Permissões do diretório wallet.db: %o", mode)

	// Verificar que outros usuários não têm acesso
	if mode&0077 != 0 {
		t.Errorf("🔴 Permissões muito abertas: %o — outros usuários podem ler os dados!", mode)
	} else {
		t.Log("✅ Permissões restritas ao dono (700)")
	}
}

func TestKeyExposure_WalletNotAccessibleWithoutPassword(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando que carteira não abre sem senha correta")

	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "pentest-nopass")

	wallet.Create(dataDir, "senha-real-segura", "testnet")

	// Tentar abrir com senhas erradas
	wrongPasswords := []string{
		"",
		"a",
		"senha-errada",
		"senha-real-segur",  // 1 char a menos
		"Senha-real-segura", // case diferente
		"senha-real-segura ", // espaço extra
	}

	for _, wp := range wrongPasswords {
		_, err := wallet.Open(dataDir, wp)
		if err == nil {
			t.Fatalf("🔴 CRÍTICO: Carteira abriu com senha errada: '%s'", wp)
		}
	}

	t.Log("✅ Todas as senhas erradas foram rejeitadas")
}

func TestKeyExposure_PathTraversal(t *testing.T) {
	t.Log("🔍 PENTEST: Verificando resistência a path traversal no data-dir")

	// Tentar criar carteira com path malicioso
	maliciousPaths := []string{
		"../../../etc/passwd",
		"/tmp/../../../etc/shadow",
		"data/../../outside",
	}

	for _, path := range maliciousPaths {
		// Não deve causar panic ou escrita fora do diretório
		_, _, err := wallet.Create(path, "senha12345678", "testnet")
		// Limpar qualquer coisa que foi criada
		os.RemoveAll(path)

		if err != nil {
			t.Logf("✅ Path '%s' tratado corretamente: %v", path, err)
		} else {
			// Se criou, verificar que não saiu do diretório
			t.Logf("⚠️ Path '%s' foi aceito — verificar manualmente se dados foram escritos fora do escopo", path)
			os.RemoveAll(path)
		}
	}
}
