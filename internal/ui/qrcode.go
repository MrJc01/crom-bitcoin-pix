package ui

import (
	"fmt"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// GenerateQR gera um QR Code como string ASCII para exibição no terminal.
// Usa blocos Unicode para melhor resolução.
func GenerateQR(content string) (string, error) {
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("falha ao gerar QR: %w", err)
	}

	bitmap := qr.Bitmap()
	size := len(bitmap)

	var sb strings.Builder

	// Usar blocos Unicode para compactar 2 linhas em 1
	// ▀ (upper half), ▄ (lower half), █ (full), ' ' (empty)
	for y := 0; y < size; y += 2 {
		for x := 0; x < size; x++ {
			upper := bitmap[y][x]
			lower := false
			if y+1 < size {
				lower = bitmap[y+1][x]
			}

			switch {
			case upper && lower:
				sb.WriteString("█")
			case upper && !lower:
				sb.WriteString("▀")
			case !upper && lower:
				sb.WriteString("▄")
			default:
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// GenerateQRInverted gera QR com cores invertidas (melhor para terminais escuros).
func GenerateQRInverted(content string) (string, error) {
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("falha ao gerar QR: %w", err)
	}

	bitmap := qr.Bitmap()
	size := len(bitmap)

	var sb strings.Builder

	for y := 0; y < size; y += 2 {
		for x := 0; x < size; x++ {
			upper := !bitmap[y][x] // invertido
			lower := true
			if y+1 < size {
				lower = !bitmap[y+1][x] // invertido
			}

			switch {
			case upper && lower:
				sb.WriteString("█")
			case upper && !lower:
				sb.WriteString("▀")
			case !upper && lower:
				sb.WriteString("▄")
			default:
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// BitcoinURI gera uma URI BIP-21 para pagamento Bitcoin.
// Formato: bitcoin:<address>?amount=<btc>&label=<label>
func BitcoinURI(address string, amountSats int64, label string) string {
	uri := fmt.Sprintf("bitcoin:%s", address)

	params := []string{}
	if amountSats > 0 {
		btc := float64(amountSats) / 100_000_000.0
		params = append(params, fmt.Sprintf("amount=%.8f", btc))
	}
	if label != "" {
		params = append(params, fmt.Sprintf("label=%s", label))
	}

	if len(params) > 0 {
		uri += "?" + strings.Join(params, "&")
	}

	return uri
}

// FormatSats formata satoshis para exibição legível.
func FormatSats(sats int64) string {
	if sats >= 100_000_000 {
		btc := float64(sats) / 100_000_000.0
		return fmt.Sprintf("%.8f BTC", btc)
	}
	if sats >= 1_000_000 {
		return fmt.Sprintf("%d sats (%.2f mBTC)", sats, float64(sats)/100_000.0)
	}
	if sats >= 1_000 {
		return fmt.Sprintf("%d sats", sats)
	}
	return fmt.Sprintf("%d sats", sats)
}

// FormatAddress trunca um endereço Bitcoin para display.
func FormatAddress(address string) string {
	if len(address) <= 16 {
		return address
	}
	return address[:8] + "..." + address[len(address)-8:]
}
