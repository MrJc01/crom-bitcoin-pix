package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Estilos do TUI
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F7931A")).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#F7931A")).
		Padding(0, 2)

	menuStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF88")).
		Background(lipgloss.Color("#333333")).
		Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00BFFF"))

	balanceStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F7931A")).
		Padding(0, 1)

	dimStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666"))

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF4444")).
		Bold(true)

	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF88")).
		Bold(true)

	boxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#F7931A")).
		Padding(1, 2).
		Width(60)
)

// Screen define as telas do TUI.
type Screen int

const (
	ScreenDashboard Screen = iota
	ScreenReceive
	ScreenSend
	ScreenNostr
	ScreenSettings
)

// WalletData contém os dados exibidos no TUI.
type WalletData struct {
	Address       string
	Balance       int64
	Unconfirmed   int64
	Network       string
	BlockHeight   int64
	NostrNpub     string
	NostrNIP05    string
	LightningPub  string
	Synced        bool
	QRCode        string
}

// Model é o model do bubbletea TUI.
type Model struct {
	screen      Screen
	menuIndex   int
	walletData  WalletData
	spinner     spinner.Model
	loading     bool
	message     string
	messageType string // "info", "error", "success"
	width       int
	height      int
	quitting    bool
}

// menuItems define os itens do menu principal.
var menuItems = []string{
	"💰 Dashboard",
	"📥 Receber",
	"📤 Enviar",
	"🌐 Nostr",
	"⚙️  Configurações",
	"🚪 Sair",
}

// NewModel cria um novo model TUI.
func NewModel(data WalletData) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#F7931A"))

	return Model{
		screen:     ScreenDashboard,
		walletData: data,
		spinner:    s,
		width:      80,
		height:     24,
	}
}

// Init implementa tea.Model.
func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update implementa tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.menuIndex > 0 {
				m.menuIndex--
			}

		case "down", "j":
			if m.menuIndex < len(menuItems)-1 {
				m.menuIndex++
			}

		case "enter":
			switch m.menuIndex {
			case 0:
				m.screen = ScreenDashboard
			case 1:
				m.screen = ScreenReceive
			case 2:
				m.screen = ScreenSend
			case 3:
				m.screen = ScreenNostr
			case 4:
				m.screen = ScreenSettings
			case 5:
				m.quitting = true
				return m, tea.Quit
			}

		case "esc":
			m.screen = ScreenDashboard
			m.message = ""

		case "1":
			m.screen = ScreenDashboard
		case "2":
			m.screen = ScreenReceive
		case "3":
			m.screen = ScreenSend
		case "4":
			m.screen = ScreenNostr
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View implementa tea.Model.
func (m Model) View() string {
	if m.quitting {
		return "\n  ⚡ Até logo! — Crom Bitcoin Pix\n\n"
	}

	var sb strings.Builder

	// Header
	header := titleStyle.Render("⚡ Crom Bitcoin Pix")
	sb.WriteString("\n" + header + "\n\n")

	// Layout: Menu à esquerda + Conteúdo à direita
	menu := m.renderMenu()
	content := m.renderContent()

	// Join horizontal
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, menu, "  ", content))

	// Footer
	sb.WriteString("\n\n")
	sb.WriteString(dimStyle.Render("  [↑↓] navegar  [enter] selecionar  [esc] voltar  [q] sair"))
	sb.WriteString("\n")

	// Mensagem
	if m.message != "" {
		var msgStyle lipgloss.Style
		switch m.messageType {
		case "error":
			msgStyle = errorStyle
		case "success":
			msgStyle = successStyle
		default:
			msgStyle = infoStyle
		}
		sb.WriteString("\n  " + msgStyle.Render(m.message) + "\n")
	}

	return sb.String()
}

// renderMenu renderiza o menu lateral.
func (m Model) renderMenu() string {
	var items []string
	for i, item := range menuItems {
		if i == m.menuIndex {
			items = append(items, selectedStyle.Render("▶ "+item))
		} else {
			items = append(items, menuStyle.Render("  "+item))
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// renderContent renderiza o conteúdo da tela atual.
func (m Model) renderContent() string {
	switch m.screen {
	case ScreenDashboard:
		return m.renderDashboard()
	case ScreenReceive:
		return m.renderReceive()
	case ScreenSend:
		return m.renderSend()
	case ScreenNostr:
		return m.renderNostr()
	case ScreenSettings:
		return m.renderSettings()
	}
	return ""
}

// renderDashboard renderiza a tela principal.
func (m Model) renderDashboard() string {
	d := m.walletData

	balance := FormatSats(d.Balance)
	if d.Unconfirmed > 0 {
		balance += fmt.Sprintf(" (+%s pendente)", FormatSats(d.Unconfirmed))
	}

	content := fmt.Sprintf(
		"%s\n\n"+
			"%s %s\n\n"+
			"%s %s\n"+
			"%s %s\n"+
			"%s %s",
		balanceStyle.Render("💰 "+balance),
		infoStyle.Render("⚡ Endereço:"), d.Address,
		dimStyle.Render("🌐 Rede:"), d.Network,
		dimStyle.Render("🔗 Bloco:"), fmt.Sprintf("%d", d.BlockHeight),
		dimStyle.Render("⚡ Lightning:"), m.lightningStatus(),
	)

	return boxStyle.Render(content)
}

// renderReceive renderiza a tela de recebimento com QR Code.
func (m Model) renderReceive() string {
	d := m.walletData

	content := fmt.Sprintf(
		"%s\n\n"+
			"%s\n\n"+
			"%s",
		infoStyle.Render("📥 Receber Bitcoin"),
		successStyle.Render("⚡ "+d.Address),
		dimStyle.Render("Escaneie o QR abaixo com qualquer wallet:"),
	)

	if d.QRCode != "" {
		content += "\n\n" + d.QRCode
	}

	return boxStyle.Render(content)
}

// renderSend renderiza a tela de envio.
func (m Model) renderSend() string {
	content := fmt.Sprintf(
		"%s\n\n"+
			"%s\n\n"+
			"%s\n"+
			"%s\n"+
			"%s",
		infoStyle.Render("📤 Enviar Bitcoin"),
		dimStyle.Render("Destinos suportados:"),
		dimStyle.Render("  • Endereço Bitcoin (bc1q...)"),
		dimStyle.Render("  • Invoice Lightning (lnbc...)"),
		dimStyle.Render("  • Nostr NIP-05 (user@crom.run)"),
	)

	return boxStyle.Render(content)
}

// renderNostr renderiza a tela Nostr.
func (m Model) renderNostr() string {
	d := m.walletData

	npub := d.NostrNpub
	if npub == "" {
		npub = dimStyle.Render("(não configurado)")
	}
	nip05 := d.NostrNIP05
	if nip05 == "" {
		nip05 = dimStyle.Render("(não configurado)")
	}

	content := fmt.Sprintf(
		"%s\n\n"+
			"%s %s\n"+
			"%s %s\n\n"+
			"%s",
		infoStyle.Render("🌐 Identidade Nostr"),
		dimStyle.Render("npub:"), npub,
		dimStyle.Render("NIP-05:"), nip05,
		dimStyle.Render("Sua identidade descentralizada no ecossistema Crom"),
	)

	return boxStyle.Render(content)
}

// renderSettings renderiza a tela de configurações.
func (m Model) renderSettings() string {
	d := m.walletData

	content := fmt.Sprintf(
		"%s\n\n"+
			"%s %s\n"+
			"%s %s\n"+
			"%s %v",
		infoStyle.Render("⚙️  Configurações"),
		dimStyle.Render("Rede:"), d.Network,
		dimStyle.Render("Lightning:"), m.lightningStatus(),
		dimStyle.Render("Sincronizado:"), d.Synced,
	)

	return boxStyle.Render(content)
}

func (m Model) lightningStatus() string {
	if m.walletData.LightningPub != "" {
		return successStyle.Render("✅ Conectado")
	}
	return dimStyle.Render("⬜ Não configurado")
}

// RunTUI inicia o TUI interativo.
func RunTUI(data WalletData) error {
	p := tea.NewProgram(NewModel(data), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
