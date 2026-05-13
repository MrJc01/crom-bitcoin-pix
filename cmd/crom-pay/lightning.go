package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/MrJc01/crom-bitcoin-pix/internal/lightning"
)

func init() {
	lnCmd := &cobra.Command{
		Use:   "lightning",
		Short: "⚡ Gerenciar nó Lightning Network",
	}

	lnCmd.AddCommand(lnInfoCmd)
	lnCmd.AddCommand(lnInvoiceCmd)
	lnCmd.AddCommand(lnChannelsCmd)
	lnCmd.AddCommand(lnSetupCmd)

	rootCmd.AddCommand(lnCmd)
}

// ─── lightning info ──────────────────────────────────────────────────────────

var lnInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Informações do nó Lightning",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := lightning.NewStubClient(network)
		nodeInfo, err := client.GetInfo()
		if err != nil {
			fmt.Println()
			warn.Println("⚡ Lightning Network")
			fmt.Println()
			dim.Println("  Status: Não configurado")
			dim.Println()
			dim.Println("  Para conectar a um nó LND:")
			dim.Println("    crom-pay lightning setup --host localhost:10009")
			dim.Println()
			dim.Println("  Para instalar LND:")
			dim.Println("    https://docs.lightning.engineering/lightning-network-tools/lnd/run-lnd")
			fmt.Println()
			return nil
		}

		fmt.Println()
		success.Println("⚡ Lightning Network")
		fmt.Println()
		info.Printf("  PubKey:   %s\n", nodeInfo.PubKey)
		info.Printf("  Alias:    %s\n", nodeInfo.Alias)
		dim.Printf("  Canais:   %d\n", nodeInfo.NumChannels)
		dim.Printf("  Peers:    %d\n", nodeInfo.NumPeers)
		dim.Printf("  Bloco:    %d\n", nodeInfo.BlockHeight)
		dim.Printf("  Versão:   %s\n", nodeInfo.Version)
		fmt.Println()

		return nil
	},
}

// ─── lightning invoice ───────────────────────────────────────────────────────

var lnInvoiceCmd = &cobra.Command{
	Use:   "invoice <amount> [memo]",
	Short: "Gerar invoice Lightning para recebimento",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var amount int64
		fmt.Sscanf(args[0], "%d", &amount)
		if amount <= 0 {
			return fmt.Errorf("❌ valor deve ser positivo (em satoshis)")
		}

		memo := "Pagamento Crom-Pay"
		if len(args) > 1 {
			memo = args[1]
		}

		client := lightning.NewStubClient(network)
		inv, err := client.CreateInvoice(amount, memo)
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}

		fmt.Println()
		success.Printf("⚡ Invoice: %s\n", inv.PaymentRequest)
		fmt.Println()

		return nil
	},
}

// ─── lightning channels ──────────────────────────────────────────────────────

var lnChannelsCmd = &cobra.Command{
	Use:   "channels",
	Short: "Listar canais Lightning",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := lightning.NewStubClient(network)
		channels, err := client.ListChannels()
		if err != nil {
			fmt.Println()
			dim.Println("  Nenhum canal aberto. Use 'lightning setup' primeiro.")
			fmt.Println()
			return nil
		}

		fmt.Println()
		info.Println("⚡ Canais Lightning")
		fmt.Println()
		for _, ch := range channels {
			status := "✅"
			if !ch.Active {
				status = "❌"
			}
			fmt.Printf("  %s Canal %d | Local: %d sats | Remoto: %d sats\n",
				status, ch.ChanID, ch.LocalBalance, ch.RemoteBalance)
		}
		fmt.Println()

		return nil
	},
}

// ─── lightning setup ─────────────────────────────────────────────────────────

var lnSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configurar conexão com nó LND",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println()
		info.Println("⚡ Setup Lightning Network")
		fmt.Println()
		dim.Println("  O Crom-Pay suporta dois modos Lightning:")
		fmt.Println()
		bold.Println("  1. Nó externo (recomendado para começar)")
		dim.Println("     Conecta a um LND já rodando via gRPC.")
		dim.Println("     crom-pay lightning setup --host localhost:10009 \\")
		dim.Println("       --tls-cert ~/.lnd/tls.cert \\")
		dim.Println("       --macaroon ~/.lnd/data/chain/bitcoin/mainnet/admin.macaroon")
		fmt.Println()
		bold.Println("  2. Nó embutido (all-in-one)")
		dim.Println("     LND integrado ao binário com Neutrino.")
		dim.Println("     crom-pay lightning start --embedded")
		dim.Println("     ⚠️  Requer ~200MB RAM e 10-30 min para sync inicial.")
		fmt.Println()
		dim.Println("  Instalar LND: https://docs.lightning.engineering/")
		fmt.Println()

		return nil
	},
}
