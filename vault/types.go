package vault

import (
	"errors"
)

// ErrNotFound indicates that the specified crypto storage entry couldn't be found.
var ErrNotFound = errors.New("entry not found")
var errKeyNotFound = errors.New("key not found")
var errKeyAlreadyExists = errors.New("key already exists")

// Storage interface containing functions for storing and retrieving keys.
type Storage interface {
	// GetSecret from the storage backend and return its value.
	GetSecret(key string) ([]byte, error)
	// StoreSecret stores the secret under the key in the storage backend.
	StoreSecret(key string, value []byte) error
	// DeleteSecret the key under the given key in the storage backend.
	DeleteSecret(key string) error
	// ListKeys returns a list of all keys in the storage backend.
	ListKeys() ([]string, error)
}
