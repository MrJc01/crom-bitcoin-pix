package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/MrJc01/crom-bitcoin-pix/internal/wallet"
)

// Cores reutilizáveis
var (
	bold    = color.New(color.Bold)
	success = color.New(color.FgHiGreen, color.Bold)
	warn    = color.New(color.FgHiYellow)
	danger  = color.New(color.FgHiRed, color.Bold)
	info    = color.New(color.FgHiCyan)
	dim     = color.New(color.FgWhite)
)

// stdinScanner é compartilhado entre todas as leituras de stdin em modo pipe,
// evitando que um bufio.Reader consuma todo o buffer de uma vez.
var stdinScanner = bufio.NewScanner(os.Stdin)

func init() {
	walletCmd := &cobra.Command{
		Use:   "wallet",
		Short: "💰 Gerenciar carteira Bitcoin",
	}

	walletCmd.AddCommand(walletCreateCmd)
	walletCmd.AddCommand(walletBalanceCmd)
	walletCmd.AddCommand(walletRestoreCmd)
	walletCmd.AddCommand(walletAddressCmd)

	rootCmd.AddCommand(walletCmd)
}

// ─── wallet create ───────────────────────────────────────────────────────────

var walletCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Criar nova carteira Bitcoin",
	Long:  "Gera uma nova carteira HD com semente BIP-39 de 24 palavras.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println()
		info.Println("🔐 Criando carteira soberana...")
		fmt.Println()

		// Pedir senha
		password, err := readPassword("Digite uma senha para proteger sua carteira: ")
		if err != nil {
			return err
		}

		confirm, err := readPassword("Confirme a senha: ")
		if err != nil {
			return err
		}

		if password != confirm {
			return fmt.Errorf("❌ senhas não coincidem")
		}

		// Criar carteira (validação de senha mínima é feita no domínio)
		mnemonic, walletInfo, err := wallet.Create(dataDir, password, network)
		if err != nil {
			return fmt.Errorf("❌ falha ao criar carteira: %w", err)
		}

		// Exibir semente
		fmt.Println()
		danger.Println("╔════════════════════════════════════════════════════════════╗")
		danger.Println("║  ⚠️  ANOTE ESTAS PALAVRAS EM PAPEL — NUNCA COMPARTILHE!   ║")
		danger.Println("╠════════════════════════════════════════════════════════════╣")

		words := strings.Split(mnemonic, " ")
		for i := 0; i < len(words); i += 4 {
			end := i + 4
			if end > len(words) {
				end = len(words)
			}
			line := "║  "
			for j := i; j < end; j++ {
				line += fmt.Sprintf("%2d. %-12s", j+1, words[j])
			}
			for len(line) < 61 {
				line += " "
			}
			line += "║"
			warn.Println(line)
		}

		danger.Println("╚════════════════════════════════════════════════════════════╝")

		// Exibir informações
		fmt.Println()
		success.Printf("⚡ Endereço Bitcoin: %s\n", walletInfo.Address)
		info.Printf("🌐 Rede: %s\n", walletInfo.Network)
		dim.Printf("📁 Dados: %s\n", dataDir)

		fmt.Println()
		success.Println("✅ Carteira criada com sucesso!")
		dim.Println("💡 Envie Bitcoin para o endereço acima para começar.")
		fmt.Println()

		return nil
	},
}

// ─── wallet balance ──────────────────────────────────────────────────────────

var walletBalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Consultar saldo da carteira",
	RunE: func(cmd *cobra.Command, args []string) error {
		password, err := readPassword("Senha da carteira: ")
		if err != nil {
			return err
		}

		w, err := wallet.Open(dataDir, password)
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}
		defer w.Close()

		walletInfo, err := w.GetInfo()
		if err != nil {
			return err
		}

		fmt.Println()
		bold.Println("╔═══════════════════════════════════════╗")
		bold.Println("║        💰 SALDO CROM-PAY              ║")
		bold.Println("╠═══════════════════════════════════════╣")
		fmt.Println("║                                       ║")
		info.Printf("║  Endereço:  %s\n", walletInfo.Address)
		warn.Printf("║  Saldo:     %d sats\n", walletInfo.Balance)
		dim.Printf("║  Rede:      %s\n", walletInfo.Network)
		fmt.Println("║                                       ║")
		bold.Println("╚═══════════════════════════════════════╝")
		fmt.Println()

		dim.Println("💡 Saldo real disponível após integração Neutrino (Milestone 02)")
		fmt.Println()

		return nil
	},
}

// ─── wallet address ──────────────────────────────────────────────────────────

var walletAddressCmd = &cobra.Command{
	Use:   "address",
	Short: "Exibir endereço Bitcoin de recebimento",
	RunE: func(cmd *cobra.Command, args []string) error {
		password, err := readPassword("Senha da carteira: ")
		if err != nil {
			return err
		}

		w, err := wallet.Open(dataDir, password)
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}
		defer w.Close()

		addr, err := w.GetAddress(0)
		if err != nil {
			return err
		}

		fmt.Println()
		success.Printf("⚡ Endereço Bitcoin: %s\n", addr)
		fmt.Println()

		return nil
	},
}

// ─── wallet restore ──────────────────────────────────────────────────────────

var walletRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restaurar carteira a partir de semente BIP-39",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println()
		info.Println("🔑 Restauração de carteira")
		fmt.Println()

		// Ler mnemônico
		warn.Print("Digite sua semente (24 palavras separadas por espaço):\n> ")
		mnemonic, err := readLine()
		if err != nil {
			return fmt.Errorf("falha ao ler semente: %w", err)
		}

		if !wallet.ValidateMnemonic(mnemonic) {
			return fmt.Errorf("❌ semente inválida — verifique as palavras e tente novamente")
		}

		// Pedir senha
		password, err := readPassword("Digite uma senha para proteger a carteira: ")
		if err != nil {
			return err
		}

		confirm, err := readPassword("Confirme a senha: ")
		if err != nil {
			return err
		}

		if password != confirm {
			return fmt.Errorf("❌ senhas não coincidem")
		}

		// Restaurar
		info.Println()
		info.Println("🔄 Restaurando carteira...")

		walletInfo, err := wallet.Restore(dataDir, mnemonic, password, network)
		if err != nil {
			return fmt.Errorf("❌ falha ao restaurar: %w", err)
		}

		fmt.Println()
		success.Printf("⚡ Endereço Bitcoin: %s\n", walletInfo.Address)
		info.Printf("🌐 Rede: %s\n", walletInfo.Network)

		fmt.Println()
		success.Println("✅ Carteira restaurada com sucesso!")
		fmt.Println()

		return nil
	},
}

// ─── Utilitários ──────────────────────────────────────────────────────────────

// readPassword lê senha do terminal sem exibir os caracteres digitados.
// Em modo pipe (CI/testes), usa o scanner compartilhado para evitar EOF.
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	// Tentar ler sem eco (funciona em terminais reais)
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		password, err := term.ReadPassword(fd)
		fmt.Println() // newline após input oculto
		if err != nil {
			return "", fmt.Errorf("falha ao ler senha: %w", err)
		}
		return string(password), nil
	}

	// Fallback para ambientes não-terminais (pipes, CI)
	// Usa scanner compartilhado para não consumir todo o stdin
	return readLine()
}

// readLine lê uma linha do stdin usando o scanner compartilhado.
func readLine() (string, error) {
	if stdinScanner.Scan() {
		return strings.TrimSpace(stdinScanner.Text()), nil
	}
	if err := stdinScanner.Err(); err != nil {
		return "", fmt.Errorf("falha ao ler entrada: %w", err)
	}
	return "", fmt.Errorf("entrada vazia (EOF)")
}
