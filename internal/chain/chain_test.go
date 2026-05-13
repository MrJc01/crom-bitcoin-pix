package chain

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		network string
		wantURL string
	}{
		{"mainnet", "mainnet", EsploraMainnet},
		{"testnet", "testnet", EsploraTestnet},
		{"default", "", EsploraMainnet},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.network)
			if c.baseURL != tt.wantURL {
				t.Errorf("NewClient(%s) baseURL = %s, want %s", tt.network, c.baseURL, tt.wantURL)
			}
			if c.httpClient == nil {
				t.Fatal("httpClient is nil")
			}
		})
	}
}

func TestNewTxBuilder(t *testing.T) {
	tests := []struct {
		name    string
		network string
	}{
		{"mainnet", "mainnet"},
		{"testnet", "testnet"},
		{"regtest", "regtest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := NewTxBuilder(tt.network, 5)
			if tb.network == nil {
				t.Fatal("network params is nil")
			}
			if tb.feeRate != 5 {
				t.Errorf("feeRate = %d, want 5", tb.feeRate)
			}
		})
	}
}

func TestEstimateFee(t *testing.T) {
	tb := NewTxBuilder("mainnet", 10)

	// 1 input, 2 outputs: ~(11 + 68 + 62) = 141 vBytes
	fee := tb.EstimateFee(1, 2)
	if fee <= 0 {
		t.Errorf("EstimateFee(1, 2) = %d, want > 0", fee)
	}

	expected := int64((11 + 68 + 62) * 10)
	if fee != expected {
		t.Errorf("EstimateFee(1, 2) = %d, want %d", fee, expected)
	}

	// Mais inputs = mais fee
	fee2 := tb.EstimateFee(3, 2)
	if fee2 <= fee {
		t.Errorf("3 inputs (%d) deveria custar mais que 1 input (%d)", fee2, fee)
	}
}

func TestClientNetwork(t *testing.T) {
	c := NewClient("testnet")
	if c.Network() != "testnet" {
		t.Errorf("Network() = %s, want testnet", c.Network())
	}
}
