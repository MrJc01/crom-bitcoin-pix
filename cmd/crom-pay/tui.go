package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/MrJc01/crom-bitcoin-pix/internal/chain"
	"github.com/MrJc01/crom-bitcoin-pix/internal/nostr"
	"github.com/MrJc01/crom-bitcoin-pix/internal/ui"
	"github.com/MrJc01/crom-bitcoin-pix/internal/wallet"
)

func init() {
	rootCmd.AddCommand(tuiCmd)
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "🖥️  Interface visual interativa no terminal",
	Long:  "Abre a interface TUI (Terminal User Interface) com dashboard, saldo, QR Code e navegação por teclado.",
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

		// Obter dados da wallet
		walletInfo, err := w.GetInfo()
		if err != nil {
			return err
		}

		// Consultar saldo real via API
		client := chain.NewClient(network)
		confirmed, unconfirmed, _ := client.GetBalance(walletInfo.Address)
		blockHeight, _ := client.GetBlockHeight()

		// Derivar Nostr
		var npub string
		seed, err := w.GetSeed(password)
		if err == nil {
			keys, err := nostr.DeriveFromSeed(seed, 0, 0)
			wallet.ZeroBytes(seed)
			if err == nil {
				npub = keys.Npub()
			}
		}

		// Gerar QR
		qrContent := ui.BitcoinURI(walletInfo.Address, 0, "")
		qr, _ := ui.GenerateQRInverted(qrContent)

		// Montar dados
		data := ui.WalletData{
			Address:     walletInfo.Address,
			Balance:     confirmed,
			Unconfirmed: unconfirmed,
			Network:     walletInfo.Network,
			BlockHeight: blockHeight,
			NostrNpub:   npub,
			Synced:      true,
			QRCode:      qr,
		}

		return ui.RunTUI(data)
	},
}
