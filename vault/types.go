/*
 * Copyright (C) 2023 Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
