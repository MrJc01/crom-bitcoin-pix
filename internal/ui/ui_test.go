package ui

import (
	"strings"
	"testing"
)

func TestGenerateQR(t *testing.T) {
	qr, err := GenerateQR("bitcoin:bc1qtest")
	if err != nil {
		t.Fatalf("GenerateQR: %v", err)
	}
	if qr == "" {
		t.Fatal("QR is empty")
	}
	if !strings.Contains(qr, "\n") {
		t.Error("QR should have newlines")
	}
}

func TestGenerateQRInverted(t *testing.T) {
	qr, err := GenerateQRInverted("bitcoin:bc1qtest")
	if err != nil {
		t.Fatalf("GenerateQRInverted: %v", err)
	}
	if qr == "" {
		t.Fatal("QR inverted is empty")
	}
}

func TestBitcoinURI(t *testing.T) {
	tests := []struct {
		name    string
		address string
		amount  int64
		label   string
		want    string
	}{
		{"basic", "bc1qtest", 0, "", "bitcoin:bc1qtest"},
		{"with amount", "bc1qtest", 100000, "", "bitcoin:bc1qtest?amount=0.00100000"},
		{"with label", "bc1qtest", 0, "Crom", "bitcoin:bc1qtest?label=Crom"},
		{"full", "bc1qtest", 50000, "Pay", "bitcoin:bc1qtest?amount=0.00050000&label=Pay"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BitcoinURI(tt.address, tt.amount, tt.label)
			if got != tt.want {
				t.Errorf("BitcoinURI() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestFormatSats(t *testing.T) {
	tests := []struct {
		sats int64
		want string
	}{
		{0, "0 sats"},
		{500, "500 sats"},
		{50000, "50000 sats"},
		{1500000, "1500000 sats (15.00 mBTC)"},
		{100000000, "1.00000000 BTC"},
		{250000000, "2.50000000 BTC"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := FormatSats(tt.sats)
			if got != tt.want {
				t.Errorf("FormatSats(%d) = %s, want %s", tt.sats, got, tt.want)
			}
		})
	}
}

func TestFormatAddress(t *testing.T) {
	short := "bc1qtest"
	if FormatAddress(short) != short {
		t.Errorf("short address should not be truncated")
	}

	long := "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq"
	formatted := FormatAddress(long)
	if !strings.Contains(formatted, "...") {
		t.Error("long address should be truncated with ...")
	}
	if len(formatted) > 20 {
		t.Error("truncated address too long")
	}
}

func TestNewModel(t *testing.T) {
	data := WalletData{
		Address: "bc1qtest",
		Balance: 50000,
		Network: "mainnet",
	}

	m := NewModel(data)
	if m.screen != ScreenDashboard {
		t.Error("initial screen should be Dashboard")
	}
	if m.walletData.Balance != 50000 {
		t.Error("balance not set")
	}
}
