package lightning

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	tests := []struct {
		name    string
		network string
	}{
		{"mainnet", "mainnet"},
		{"testnet", "testnet"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig(tt.network)
			if cfg.Host == "" {
				t.Fatal("Host is empty")
			}
			if cfg.Network != tt.network {
				t.Errorf("Network = %s, want %s", cfg.Network, tt.network)
			}
			if cfg.MacaroonPath == "" {
				t.Fatal("MacaroonPath is empty")
			}
		})
	}
}

func TestStubClient(t *testing.T) {
	stub := NewStubClient("testnet")

	if stub.IsConnected() {
		t.Error("StubClient should not be connected")
	}

	// Todos os métodos devem retornar ErrNotConnected
	_, err := stub.GetInfo()
	if err != ErrNotConnected {
		t.Errorf("GetInfo() err = %v, want ErrNotConnected", err)
	}

	_, err = stub.CreateInvoice(1000, "test")
	if err != ErrNotConnected {
		t.Errorf("CreateInvoice() err = %v, want ErrNotConnected", err)
	}

	_, err = stub.PayInvoice("lnbc...")
	if err != ErrNotConnected {
		t.Errorf("PayInvoice() err = %v, want ErrNotConnected", err)
	}

	_, err = stub.ListChannels()
	if err != ErrNotConnected {
		t.Errorf("ListChannels() err = %v, want ErrNotConnected", err)
	}

	_, err = stub.LookupInvoice("hash")
	if err != ErrNotConnected {
		t.Errorf("LookupInvoice() err = %v, want ErrNotConnected", err)
	}

	_, err = stub.DecodeInvoice("lnbc...")
	if err != ErrNotConnected {
		t.Errorf("DecodeInvoice() err = %v, want ErrNotConnected", err)
	}

	_, err = stub.OpenChannel("pubkey", 100000)
	if err != ErrNotConnected {
		t.Errorf("OpenChannel() err = %v, want ErrNotConnected", err)
	}

	_, err = stub.CloseChannel(123)
	if err != ErrNotConnected {
		t.Errorf("CloseChannel() err = %v, want ErrNotConnected", err)
	}

	// Connect retorna erro informativo (não ErrNotConnected)
	if err := stub.Connect(); err == nil {
		t.Error("Connect() should fail for stub")
	}

	// Close é no-op
	if err := stub.Close(); err != nil {
		t.Errorf("Close() err = %v, want nil", err)
	}
}

func TestLNDClientNew(t *testing.T) {
	cfg := DefaultConfig("testnet")
	client := NewLNDClient(cfg)

	if client == nil {
		t.Fatal("NewLNDClient returned nil")
	}
	if client.IsConnected() {
		t.Error("should not be connected before Connect()")
	}
}
