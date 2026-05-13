package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/MrJc01/crom-bitcoin-pix/internal/nostr"
	"github.com/MrJc01/crom-bitcoin-pix/internal/wallet"
)

func init() {
	nostrCmd := &cobra.Command{
		Use:   "nostr",
		Short: "🌐 Identidade Nostr descentralizada",
	}

	nostrCmd.AddCommand(nostrIdentityCmd)
	nostrCmd.AddCommand(nostrPublishCmd)
	nostrCmd.AddCommand(nostrRelaysCmd)
	nostrCmd.AddCommand(nostrVerifyCmd)

	rootCmd.AddCommand(nostrCmd)
}

// ─── nostr identity ──────────────────────────────────────────────────────────

var nostrIdentityCmd = &cobra.Command{
	Use:   "identity",
	Short: "Exibir identidade Nostr derivada do seed Bitcoin",
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

		// Derivar chaves Nostr do mesmo seed
		seed, err := w.GetSeed(password)
		if err != nil {
			return fmt.Errorf("❌ falha ao obter seed: %w", err)
		}

		keys, err := nostr.DeriveFromSeed(seed, 0, 0)
		wallet.ZeroBytes(seed)
		if err != nil {
			return fmt.Errorf("❌ falha ao derivar chaves Nostr: %w", err)
		}

		fmt.Println()
		info.Println("🌐 Identidade Nostr")
		fmt.Println()
		success.Printf("  npub: %s\n", keys.Npub())
		dim.Printf("  hex:  %s\n", keys.PublicKey)
		dim.Printf("  fingerprint: %s\n", keys.Fingerprint())
		fmt.Println()
		warn.Println("  ⚠️  Sua nsec (chave privada) é derivada do mesmo seed Bitcoin.")
		dim.Println("  💡 Quem tem suas 24 palavras tem acesso ao seu Nostr também.")
		fmt.Println()

		return nil
	},
}

// ─── nostr publish ───────────────────────────────────────────────────────────

var nostrPublishCmd = &cobra.Command{
	Use:   "publish <mensagem>",
	Short: "Publicar nota em relays Nostr",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		message := args[0]

		password, err := readPassword("Senha da carteira: ")
		if err != nil {
			return err
		}

		w, err := wallet.Open(dataDir, password)
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}
		defer w.Close()

		seed, err := w.GetSeed(password)
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}

		keys, err := nostr.DeriveFromSeed(seed, 0, 0)
		wallet.ZeroBytes(seed)
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}

		// Criar evento
		event, err := keys.CreateNoteEvent(message)
		if err != nil {
			return fmt.Errorf("❌ %w", err)
		}

		// Publicar em relays
		info.Println("\n🔄 Conectando a relays...")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		pool := nostr.NewRelayPool(nostr.DefaultRelays())
		if err := pool.Connect(ctx); err != nil {
			return fmt.Errorf("❌ %w", err)
		}
		defer pool.Close()

		if err := pool.Publish(ctx, *event); err != nil {
			return fmt.Errorf("❌ falha ao publicar: %w", err)
		}

		fmt.Println()
		success.Println("✅ Nota publicada com sucesso!")
		dim.Printf("  ID: %s\n", event.ID[:16]+"...")
		dim.Printf("  npub: %s\n", keys.Npub())
		fmt.Println()

		return nil
	},
}

// ─── nostr relays ────────────────────────────────────────────────────────────

var nostrRelaysCmd = &cobra.Command{
	Use:   "relays",
	Short: "Listar relays Nostr configurados",
	RunE: func(cmd *cobra.Command, args []string) error {
		relays := nostr.DefaultRelays()

		fmt.Println()
		info.Println("🌐 Relays Nostr")
		fmt.Println()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		pool := nostr.NewRelayPool(relays)
		pool.Connect(ctx)
		defer pool.Close()

		for _, url := range relays {
			// Tentar conectar para verificar status
			connected := false
			for _, r := range pool.URLs {
				if r == url {
					connected = true
					break
				}
			}

			if connected {
				success.Printf("  ✅ %s\n", url)
			} else {
				danger.Printf("  ❌ %s\n", url)
			}
		}

		fmt.Println()
		return nil
	},
}

// ─── nostr verify ────────────────────────────────────────────────────────────

var nostrVerifyCmd = &cobra.Command{
	Use:   "verify <nip05@domain>",
	Short: "Verificar identidade NIP-05",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nip05Addr := args[0]

		fmt.Println()
		info.Printf("🔍 Verificando %s...\n", nip05Addr)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Verificação sem pubkey específico (apenas testar se existe)
		verified, err := nostr.VerifyNIP05(ctx, nip05Addr, "")
		if err != nil {
			warn.Printf("  ⚠️  Erro: %v\n", err)
		} else if verified {
			success.Printf("  ✅ %s verificado!\n", nip05Addr)
		} else {
			dim.Printf("  ℹ️  Endereço encontrado (pubkey diferente do esperado)\n")
		}

		fmt.Println()
		return nil
	},
}
