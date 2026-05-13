package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Flags globais
	dataDir string
	network string
)

// Banner ASCII do Crom-Pay
const banner = `
   ██████╗██████╗  ██████╗ ███╗   ███╗      ██████╗  █████╗ ██╗   ██╗
  ██╔════╝██╔══██╗██╔═══██╗████╗ ████║      ██╔══██╗██╔══██╗╚██╗ ██╔╝
  ██║     ██████╔╝██║   ██║██╔████╔██║█████╗██████╔╝███████║ ╚████╔╝ 
  ██║     ██╔══██╗██║   ██║██║╚██╔╝██║╚════╝██╔═══╝ ██╔══██║  ╚██╔╝  
  ╚██████╗██║  ██║╚██████╔╝██║ ╚═╝ ██║      ██║     ██║  ██║   ██║   
   ╚═════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═╝      ╚═╝     ╚═╝  ╚═╝   ╚═╝   
`

var rootCmd = &cobra.Command{
	Use:   "crom-pay",
	Short: "⚡ Crom-Pay — Bitcoin como Pix, sem intermediários",
	Long: fmt.Sprintf(`%s
  Pagamentos Bitcoin instantâneos, descentralizados e self-custodial.
  
  Zero Backend · Zero Custódia · Soberania Digital
  
  Use 'crom-pay wallet create' para começar.`,
		color.New(color.FgHiYellow, color.Bold).Sprint(banner)),
	Version: Version,
}

func init() {
	// Diretório de dados padrão: ~/.crom-pay
	homeDir, _ := os.UserHomeDir()
	defaultDataDir := filepath.Join(homeDir, ".crom-pay")

	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", defaultDataDir,
		"Diretório de dados da carteira")
	rootCmd.PersistentFlags().StringVar(&network, "network", "testnet",
		"Rede Bitcoin: mainnet, testnet ou regtest")
}
