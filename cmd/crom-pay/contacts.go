package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/MrJc01/crom-bitcoin-pix/internal/contacts"
	"github.com/MrJc01/crom-bitcoin-pix/internal/storage"
)

func init() {
	contactsCmd := &cobra.Command{
		Use:   "contacts",
		Short: "📇 Gerenciar contatos de pagamento",
	}

	contactsCmd.AddCommand(contactAddCmd)
	contactsCmd.AddCommand(contactListCmd)
	contactsCmd.AddCommand(contactRemoveCmd)
	contactsCmd.AddCommand(contactShowCmd)

	rootCmd.AddCommand(contactsCmd)
}

// ─── contacts add ────────────────────────────────────────────────────────────

var contactAddCmd = &cobra.Command{
	Use:   "add <nome> [--btc addr] [--ln addr] [--npub npub] [--nip05 user@domain]",
	Short: "Adicionar contato",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		btcAddr, _ := cmd.Flags().GetString("btc")
		lnAddr, _ := cmd.Flags().GetString("ln")
		npub, _ := cmd.Flags().GetString("npub")
		nip05, _ := cmd.Flags().GetString("nip05")

		store, err := storage.NewStore(dataDir)
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}
		defer store.Close()

		mgr := contacts.NewManager(store)
		contact := contacts.Contact{
			Name:           name,
			BitcoinAddress: btcAddr,
			LightningAddr:  lnAddr,
			NostrNpub:      npub,
			NIP05:          nip05,
		}

		if err := mgr.Add(contact); err != nil {
			return fmt.Errorf("❌ %w", err)
		}

		fmt.Println()
		success.Printf("✅ Contato '%s' salvo\n", name)
		if btcAddr != "" {
			dim.Printf("  BTC:   %s\n", btcAddr)
		}
		if lnAddr != "" {
			dim.Printf("  LN:    %s\n", lnAddr)
		}
		if npub != "" {
			dim.Printf("  Nostr: %s\n", npub)
		}
		if nip05 != "" {
			dim.Printf("  NIP05: %s\n", nip05)
		}
		fmt.Println()

		return nil
	},
}

// ─── contacts list ───────────────────────────────────────────────────────────

var contactListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar todos os contatos",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewStore(dataDir)
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}
		defer store.Close()

		mgr := contacts.NewManager(store)
		list, err := mgr.List()
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}

		fmt.Println()
		info.Println("📇 Contatos")
		fmt.Println()

		if len(list) == 0 {
			dim.Println("  Nenhum contato salvo.")
			dim.Println("  Use 'crom-pay contacts add <nome> --btc <endereço>' para adicionar.")
		}

		for _, c := range list {
			bold.Printf("  %s\n", c.Name)
			if c.BitcoinAddress != "" {
				dim.Printf("    BTC:   %s\n", c.BitcoinAddress)
			}
			if c.LightningAddr != "" {
				dim.Printf("    LN:    %s\n", c.LightningAddr)
			}
			if c.NostrNpub != "" {
				dim.Printf("    Nostr: %s\n", c.NostrNpub)
			}
			if c.NIP05 != "" {
				dim.Printf("    NIP05: %s\n", c.NIP05)
			}
			fmt.Println()
		}

		return nil
	},
}

// ─── contacts remove ─────────────────────────────────────────────────────────

var contactRemoveCmd = &cobra.Command{
	Use:   "remove <nome>",
	Short: "Remover contato",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		store, err := storage.NewStore(dataDir)
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}
		defer store.Close()

		mgr := contacts.NewManager(store)
		if err := mgr.Remove(name); err != nil {
			return fmt.Errorf("❌ %w", err)
		}

		success.Printf("\n✅ Contato '%s' removido\n\n", name)
		return nil
	},
}

// ─── contacts show ───────────────────────────────────────────────────────────

var contactShowCmd = &cobra.Command{
	Use:   "show <nome>",
	Short: "Exibir detalhes de um contato",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		store, err := storage.NewStore(dataDir)
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}
		defer store.Close()

		mgr := contacts.NewManager(store)
		c, err := mgr.Get(name)
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}

		fmt.Println()
		bold.Printf("📇 %s\n", c.Name)
		if c.BitcoinAddress != "" {
			info.Printf("  BTC:   %s\n", c.BitcoinAddress)
		}
		if c.LightningAddr != "" {
			info.Printf("  LN:    %s\n", c.LightningAddr)
		}
		if c.NostrNpub != "" {
			info.Printf("  Nostr: %s\n", c.NostrNpub)
		}
		if c.NIP05 != "" {
			info.Printf("  NIP05: %s\n", c.NIP05)
		}
		fmt.Println()

		return nil
	},
}

func init() {
	contactAddCmd.Flags().String("btc", "", "Endereço Bitcoin")
	contactAddCmd.Flags().String("ln", "", "Lightning address")
	contactAddCmd.Flags().String("npub", "", "Nostr public key")
	contactAddCmd.Flags().String("nip05", "", "NIP-05 (user@domain)")
}
