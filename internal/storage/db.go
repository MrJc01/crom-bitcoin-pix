package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v4"
)

var (
	ErrStoreNotOpen = errors.New("store não está aberto")
	ErrKeyNotFound  = errors.New("chave não encontrada")
)

// Prefixos de bucket para organizar dados no BadgerDB.
var (
	prefixSeed    = []byte("seed:")
	prefixConfig  = []byte("config:")
	prefixTx      = []byte("tx:")
	prefixContact = []byte("contact:")
)

// Store é o wrapper do BadgerDB para persistência local.
type Store struct {
	db   *badger.DB
	path string
}

// NewStore abre ou cria um BadgerDB no diretório especificado.
//
// O diretório é criado automaticamente se não existir.
// O DB usa compressão Snappy e garbage collection automático.
func NewStore(dataDir string) (*Store, error) {
	dbPath := filepath.Join(dataDir, "wallet.db")

	// Criar diretório se não existir
	if err := os.MkdirAll(dbPath, 0700); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório: %w", err)
	}

	opts := badger.DefaultOptions(dbPath).
		WithLoggingLevel(badger.ERROR). // Silenciar logs do BadgerDB
		WithCompression(0).            // Sem compressão para simplicidade
		WithValueLogFileSize(16 << 20) // 16MB value log files

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir BadgerDB: %w", err)
	}

	return &Store{
		db:   db,
		path: dbPath,
	}, nil
}

// Close fecha o BadgerDB de forma segura.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Put grava um par chave-valor com o prefixo especificado.
func (s *Store) Put(prefix, key, value []byte) error {
	if s.db == nil {
		return ErrStoreNotOpen
	}

	// FIX: Alocar novo slice para evitar mutação do prefix compartilhado
	fullKey := make([]byte, 0, len(prefix)+len(key))
	fullKey = append(fullKey, prefix...)
	fullKey = append(fullKey, key...)
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(fullKey, value)
	})
}

// Get lê um valor pela chave com o prefixo especificado.
func (s *Store) Get(prefix, key []byte) ([]byte, error) {
	if s.db == nil {
		return nil, ErrStoreNotOpen
	}

	fullKey := make([]byte, 0, len(prefix)+len(key))
	fullKey = append(fullKey, prefix...)
	fullKey = append(fullKey, key...)
	var value []byte

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(fullKey)
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return ErrKeyNotFound
			}
			return err
		}

		value, err = item.ValueCopy(nil)
		return err
	})

	return value, err
}

// Delete remove um par chave-valor.
func (s *Store) Delete(prefix, key []byte) error {
	if s.db == nil {
		return ErrStoreNotOpen
	}

	fullKey := make([]byte, 0, len(prefix)+len(key))
	fullKey = append(fullKey, prefix...)
	fullKey = append(fullKey, key...)
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(fullKey)
	})
}

// Has verifica se uma chave existe no DB.
func (s *Store) Has(prefix, key []byte) bool {
	_, err := s.Get(prefix, key)
	return err == nil
}

// List retorna todos os valores com o prefixo dado.
func (s *Store) List(prefix []byte) ([][]byte, error) {
	if s.db == nil {
		return nil, ErrStoreNotOpen
	}

	var results [][]byte
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			results = append(results, val)
		}
		return nil
	})

	return results, err
}
