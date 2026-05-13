package wallet

import (
	"errors"
	"fmt"

	"github.com/MrJc01/crom-bitcoin-pix/internal/storage"
)

const MinPasswordLen = 8

var (
	ErrWalletExists   = errors.New("carteira já existe neste diretório")
	ErrWalletNotFound = errors.New("carteira não encontrada — use 'wallet create' primeiro")
	ErrWrongPassword  = errors.New("senha incorreta ou dados corrompidos")
	ErrPasswordWeak   = errors.New("senha deve ter pelo menos 8 caracteres")
)

// Wallet é a struct central que agrega seed, keychain e storage.
type Wallet struct {
	keychain *Keychain
	store    *storage.Store
	network  string
}

// WalletInfo contém informações públicas da carteira para display.
type WalletInfo struct {
	Address string
	Network string
	Balance int64 // em satoshis (placeholder por enquanto)
}

// Create gera uma nova carteira e retorna o mnemônico para backup.
//
// ATENÇÃO: O mnemônico retornado deve ser exibido ao usuário UMA VEZ
// e depois descartado da memória. Ele é salvo cifrado no disco.
func Create(dataDir, password, network string) (mnemonic string, info *WalletInfo, err error) {
	if len(password) < MinPasswordLen {
		return "", nil, ErrPasswordWeak
	}

	store, err := storage.NewStore(dataDir)
	if err != nil {
		return "", nil, fmt.Errorf("falha ao abrir storage: %w", err)
	}

	// Verificar se já existe uma carteira
	if store.HasWallet() {
		store.Close()
		return "", nil, ErrWalletExists
	}

	// Gerar mnemônico BIP-39 (24 palavras)
	mnemonic, err = GenerateMnemonic()
	if err != nil {
		store.Close()
		return "", nil, err
	}

	// Converter para seed
	seed, err := MnemonicToSeed(mnemonic, "")
	if err != nil {
		store.Close()
		return "", nil, err
	}

	// FIX P0 #1: Derivar endereço ANTES de zerar a seed (evita recriação)
	kc, err := NewKeychain(seed, network)
	if err != nil {
		store.Close()
		zeroBytes(seed)
		return "", nil, err
	}

	addr, err := kc.DeriveAddress(0)
	if err != nil {
		store.Close()
		zeroBytes(seed)
		kc.Zero()
		return "", nil, err
	}

	// Cifrar seed e salvar
	encrypted, err := EncryptData(seed, password)
	zeroBytes(seed) // seed zerada imediatamente após cifra
	kc.Zero()       // chave mestra zerada
	if err != nil {
		store.Close()
		return "", nil, err
	}

	if err := store.SaveEncryptedSeed(encrypted); err != nil {
		store.Close()
		return "", nil, fmt.Errorf("falha ao salvar seed: %w", err)
	}

	// Salvar configuração de rede
	if err := store.SaveConfig("network", network); err != nil {
		store.Close()
		return "", nil, err
	}

	store.Close()

	return mnemonic, &WalletInfo{
		Address: addr,
		Network: network,
		Balance: 0,
	}, nil
}

// Open carrega uma carteira existente do disco.
func Open(dataDir, password string) (*Wallet, error) {
	store, err := storage.NewStore(dataDir)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir storage: %w", err)
	}

	if !store.HasWallet() {
		store.Close()
		return nil, ErrWalletNotFound
	}

	// Carregar e decifrar seed
	encrypted, err := store.LoadEncryptedSeed()
	if err != nil {
		store.Close()
		return nil, fmt.Errorf("falha ao carregar seed: %w", err)
	}

	seed, err := DecryptData(encrypted, password)
	if err != nil {
		store.Close()
		return nil, ErrWrongPassword
	}

	// Carregar rede
	network, err := store.LoadConfig("network")
	if err != nil {
		store.Close()
		zeroBytes(seed)
		return nil, fmt.Errorf("rede não configurada: %w", err)
	}

	// Criar keychain
	kc, err := NewKeychain(seed, network)
	zeroBytes(seed)
	if err != nil {
		store.Close()
		return nil, err
	}

	return &Wallet{
		keychain: kc,
		store:    store,
		network:  network,
	}, nil
}

// Restore restaura uma carteira a partir de um mnemônico BIP-39 existente.
func Restore(dataDir, mnemonic, password, network string) (*WalletInfo, error) {
	if len(password) < MinPasswordLen {
		return nil, ErrPasswordWeak
	}

	if !ValidateMnemonic(mnemonic) {
		return nil, ErrInvalidMnemonic
	}

	store, err := storage.NewStore(dataDir)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir storage: %w", err)
	}

	if store.HasWallet() {
		store.Close()
		return nil, ErrWalletExists
	}

	// Converter para seed
	seed, err := MnemonicToSeed(mnemonic, "")
	if err != nil {
		store.Close()
		return nil, err
	}

	// Cifrar e salvar
	encrypted, err := EncryptData(seed, password)
	if err != nil {
		store.Close()
		zeroBytes(seed)
		return nil, err
	}

	if err := store.SaveEncryptedSeed(encrypted); err != nil {
		store.Close()
		zeroBytes(seed)
		return nil, err
	}

	if err := store.SaveConfig("network", network); err != nil {
		store.Close()
		zeroBytes(seed)
		return nil, err
	}

	// Derivar endereço
	kc, err := NewKeychain(seed, network)
	zeroBytes(seed)
	if err != nil {
		store.Close()
		return nil, err
	}

	addr, err := kc.DeriveAddress(0)
	if err != nil {
		store.Close()
		return nil, err
	}

	store.Close()

	return &WalletInfo{
		Address: addr,
		Network: network,
		Balance: 0,
	}, nil
}

// GetAddress retorna o endereço no índice especificado.
func (w *Wallet) GetAddress(index uint32) (string, error) {
	return w.keychain.DeriveAddress(index)
}

// GetInfo retorna informações públicas da carteira.
func (w *Wallet) GetInfo() (*WalletInfo, error) {
	addr, err := w.keychain.DeriveAddress(0)
	if err != nil {
		return nil, err
	}

	return &WalletInfo{
		Address: addr,
		Network: w.network,
		Balance: 0, // placeholder — Neutrino vem no Milestone 02
	}, nil
}

// Close fecha a carteira e libera recursos sensíveis da memória.
func (w *Wallet) Close() error {
	// FIX P0 #2: Zerar masterKey antes de fechar
	if w.keychain != nil {
		w.keychain.Zero()
	}
	if w.store != nil {
		return w.store.Close()
	}
	return nil
}
