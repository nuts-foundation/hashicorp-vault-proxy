package vault

import (
	"errors"
	vaultapi "github.com/hashicorp/vault/api"
)

// ErrNotFound indicates that the specified crypto storage entry couldn't be found.
var ErrNotFound = errors.New("entry not found")
var errKeyNotFound = errors.New("key not found")
var errKeyAlreadyExists = errors.New("key already exists")

// Storage interface containing functions for storing and retrieving keys.
type Storage interface {
	// GetSecret from the storage backend and return its value.
	GetSecret(key string) ([]byte, error)
	// StoreSecret stores the key under the key in the storage backend.
	StoreSecret(key string, value []byte) error
	// ListKeys returns a list of all keys in the storage backend.
	ListKeys() ([]string, error)
}

// logicaler is an interface which has been implemented by the mockVaultClient and real vault.Logical to allow testing vault without the server.
type logicaler interface {
	Read(path string) (*vaultapi.Secret, error)
	Write(path string, data map[string]interface{}) (*vaultapi.Secret, error)
	ReadWithData(path string, data map[string][]string) (*vaultapi.Secret, error)
}
