package vault

import (
	"errors"
)

// ErrNotFound indicates that the specified crypto storage entry couldn't be found.
var ErrNotFound = errors.New("key not found")
var ErrKeyAlreadyExists = errors.New("key already exists")

// Storage interface containing functions for storing and retrieving keys.
type Storage interface {
	// Ping checks if the server is available and the credentials are correct.
	Ping() error
	// GetSecret from the storage backend and return its value.
	GetSecret(key string) ([]byte, error)
	// StoreSecret stores the secret under the key in the storage backend.
	StoreSecret(key string, value []byte) error
	// DeleteSecret the key under the given key in the storage backend.
	DeleteSecret(key string) error
	// ListKeys returns a list of all keys in the storage backend.
	ListKeys() ([]string, error)
}
