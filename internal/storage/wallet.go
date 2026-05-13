package storage

// SaveEncryptedSeed persiste a semente cifrada no BadgerDB.
func (s *Store) SaveEncryptedSeed(ciphertext []byte) error {
	return s.Put(prefixSeed, []byte("master"), ciphertext)
}

// LoadEncryptedSeed carrega a semente cifrada do BadgerDB.
func (s *Store) LoadEncryptedSeed() ([]byte, error) {
	return s.Get(prefixSeed, []byte("master"))
}

// HasWallet verifica se já existe uma carteira criada.
func (s *Store) HasWallet() bool {
	return s.Has(prefixSeed, []byte("master"))
}

// SaveConfig salva uma configuração no formato key=value.
func (s *Store) SaveConfig(key, value string) error {
	return s.Put(prefixConfig, []byte(key), []byte(value))
}

// LoadConfig carrega uma configuração pelo nome.
func (s *Store) LoadConfig(key string) (string, error) {
	val, err := s.Get(prefixConfig, []byte(key))
	if err != nil {
		return "", err
	}
	return string(val), nil
}
