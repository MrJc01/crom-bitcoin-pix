package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MrJc01/crom-bitcoin-pix/internal/chain"
	"github.com/MrJc01/crom-bitcoin-pix/internal/ui"
	"github.com/MrJc01/crom-bitcoin-pix/internal/wallet"
)

func init() {
	rootCmd.AddCommand(receiveCmd)
	rootCmd.AddCommand(payCmd)
}

// ─── receive ─────────────────────────────────────────────────────────────────

var receiveCmd = &cobra.Command{
	Use:   "receive [amount]",
	Short: "📥 Receber Bitcoin — gera QR Code para pagamento",
	Long:  "Gera um endereço Bitcoin com QR Code para recebimento. Opcionalmente, especifique o valor em satoshis.",
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

		// Montar URI BIP-21
		var amountSats int64
		if len(args) > 0 {
			fmt.Sscanf(args[0], "%d", &amountSats)
		}
		uri := ui.BitcoinURI(addr, amountSats, "Crom-Pay")

		// Gerar QR
		qr, err := ui.GenerateQRInverted(uri)
		if err != nil {
			return fmt.Errorf("falha ao gerar QR: %w", err)
		}

		fmt.Println()
		info.Println("📥 Receber Bitcoin")
		fmt.Println()
		success.Printf("⚡ Endereço: %s\n", addr)
		if amountSats > 0 {
			warn.Printf("💰 Valor: %s\n", ui.FormatSats(amountSats))
		}
		info.Printf("🔗 URI: %s\n", uri)
		fmt.Println()
		dim.Println("Escaneie o QR Code abaixo:")
		fmt.Println()
		fmt.Println(qr)
		dim.Println("💡 Compatível com qualquer wallet BIP-21 (BlueWallet, Sparrow, etc.)")
		fmt.Println()

		return nil
	},
}

// ─── pay ─────────────────────────────────────────────────────────────────────

var payCmd = &cobra.Command{
	Use:   "pay <destino> <amount>",
	Short: "📤 Enviar Bitcoin — paga endereço, invoice ou NIP-05",
	Long: `Envia Bitcoin para um destino. Destinos suportados:
  • Endereço Bitcoin (bc1q..., tb1q...)
  • Invoice Lightning (lnbc..., lntb...)
  • Nostr NIP-05 (user@crom.run)`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		destination := args[0]

		password, err := readPassword("Senha da carteira: ")
		if err != nil {
			return err
		}

		w, err := wallet.Open(dataDir, password)
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}
		defer w.Close()

		fmt.Println()

		// Detectar tipo de destino
		switch {
		case isLightningInvoice(destination):
			// Lightning invoice
			info.Println("⚡ Pagamento Lightning detectado")
			danger.Println("❌ Lightning requer nó LND configurado — use 'crom-pay lightning setup'")

		case isBitcoinAddress(destination):
			// Endereço Bitcoin on-chain
			if len(args) < 2 {
				return fmt.Errorf("❌ especifique o valor: pay <endereço> <sats>")
			}

			var amount int64
			fmt.Sscanf(args[1], "%d", &amount)
			if amount <= 0 {
				return fmt.Errorf("❌ valor deve ser positivo")
			}

			info.Println("🔗 Pagamento on-chain")
			fmt.Println()
			bold.Printf("  Destino: %s\n", destination)
			warn.Printf("  Valor:   %s\n", ui.FormatSats(amount))

			// Consultar fee estimada
			client := chain.NewClient(network)
			fees, err := client.GetFeeEstimates()
			if err == nil {
				dim.Printf("  Taxa:    ~%d sat/vB (média)\n", fees.HalfHourFee)
			}

			fmt.Println()
			danger.Println("⚠️  Transação on-chain será implementada no Milestone 02B")
			dim.Println("💡 Use 'receive' para gerar um QR de recebimento")

		case isNIP05(destination):
			// NIP-05 address
			info.Printf("🌐 Pagamento Nostr para %s\n", destination)
			danger.Println("❌ Pagamento via NIP-05 requer Lightning — use 'crom-pay lightning setup'")

		default:
			return fmt.Errorf("❌ destino inválido: %s\n  Formatos: bc1q..., lnbc..., user@domain", destination)
		}

		fmt.Println()
		return nil
	},
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func isLightningInvoice(s string) bool {
	return len(s) > 4 && (s[:4] == "lnbc" || s[:4] == "lntb" || s[:4] == "lnbs")
}

func isBitcoinAddress(s string) bool {
	return len(s) > 3 && (s[:3] == "bc1" || s[:3] == "tb1" || s[:1] == "1" || s[:1] == "3")
}

func isNIP05(s string) bool {
	for _, c := range s {
		if c == '@' {
			return true
		}
	}
	return false
}

// cores extras (referência ao arquivo wallet.go)
var boldPay = color.New(color.Bold)

func init() {
	_ = boldPay // evitar unused
}
