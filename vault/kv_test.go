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
	"encoding/base64"
	"errors"
	"strings"
	"testing"

	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

type mockVaultClient struct {
	// when set, the methods return this error
	err   error
	store map[string]map[string]interface{}
}

func (m mockVaultClient) Read(path string) (*vault.Secret, error) {
	if m.err != nil {
		return nil, m.err
	}
	data, ok := m.store[path]
	if !ok {
		return nil, nil
	}
	return &vault.Secret{
		Data: data,
	}, nil
}

func (m mockVaultClient) ReadWithData(path string, _ map[string][]string) (*vault.Secret, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &vault.Secret{
		Data: m.store[path],
	}, nil
}

func (m mockVaultClient) Write(path string, data map[string]interface{}) (*vault.Secret, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.store[path] = data
	return &vault.Secret{
		Data: data,
	}, nil
}

func (m mockVaultClient) List(path string) (*vault.Secret, error) {
	if m.err != nil {
		return nil, m.err
	}
	var keys []interface{}
	for path := range m.store {
		parts := strings.Split(path, "/")
		keys = append(keys, parts[len(parts)-1])
	}
	return &vault.Secret{
		Data: map[string]interface{}{
			"keys": keys,
		},
	}, nil
}

func (m mockVaultClient) Delete(path string) (*vault.Secret, error) {
	if m.err != nil {
		return nil, m.err
	}
	delete(m.store, path)
	return &vault.Secret{}, nil
}

var secret = []byte("secret-value")
var encodedSecret = []byte(base64.StdEncoding.EncodeToString(secret))

const prefix = "kv"
const kid = "did:nuts:123#abc"

var vaultError = errors.New("vault error")

func TestVaultKVStorage(t *testing.T) {

	t.Run("ok - store and retrieve a secret", func(t *testing.T) {
		vaultStorage := KVStorage{client: mockVaultClient{store: map[string]map[string]interface{}{}}}
		result, err := vaultStorage.GetSecret(kid)
		assert.EqualError(t, err, ErrNotFound.Error(), "secret should not be found")
		assert.Nil(t, result, "result should be nil")

		assert.NoError(t, vaultStorage.StoreSecret(kid, secret), "storing secret should work")

		result, err = vaultStorage.GetSecret(kid)
		assert.NoError(t, err)
		assert.Equal(t, secret, result, "result should equal the secret")
	})

	t.Run("error - while writing", func(t *testing.T) {
		v := KVStorage{client: mockVaultClient{err: vaultError}}
		err := v.StoreSecret(kid, secret)
		assert.Error(t, err, "saving should fail")
		assert.ErrorIs(t, err, vaultError)
	})

	t.Run("error - while reading", func(t *testing.T) {
		v := KVStorage{client: mockVaultClient{err: vaultError}}
		_, err := v.GetSecret(kid)
		assert.Error(t, err, "saving should fail")
		assert.ErrorIs(t, err, vaultError)
	})

	t.Run("error - key not found (empty response)", func(t *testing.T) {
		v := KVStorage{client: mockVaultClient{store: map[string]map[string]interface{}{}}}
		_, err := v.GetSecret(kid)
		assert.Error(t, err, "expected error on unknown kid")
		assert.EqualError(t, err, ErrNotFound.Error())
	})

	t.Run("error - key not found (key not present in data field)", func(t *testing.T) {
		store := map[string]map[string]interface{}{
			storagePath(prefix, kid): {"other-key": "other-value"},
		}
		v := KVStorage{pathPrefix: prefix, client: mockVaultClient{store: store}}
		_, err := v.GetSecret(kid)
		assert.Error(t, err, "expected error on unknown kid")
		assert.EqualError(t, err, ErrNotFound.Error())
	})
	t.Run("error - encoding issues", func(t *testing.T) {
		path := storagePath("kv", kid)
		store := map[string]map[string]interface{}{
			path: {"key": []byte("foo")},
		}
		v := KVStorage{pathPrefix: "kv", client: mockVaultClient{store: store}}

		t.Run("GetPrivateKey", func(t *testing.T) {
			_, err := v.GetSecret(kid)
			assert.Error(t, err, "expected type conversion error on byte array")
			assert.EqualError(t, err, "unable to convert key result to string")
		})
	})

	t.Run("error - key already exists", func(t *testing.T) {
		v := KVStorage{pathPrefix: prefix, client: mockVaultClient{store: map[string]map[string]interface{}{storagePath(prefix, kid): {"key": string(encodedSecret)}}}}
		assert.EqualError(t, v.StoreSecret(kid, secret), ErrKeyAlreadyExists.Error())
	})
}

func TestVaultKVStorage_ListKeys(t *testing.T) {

	t.Run("ok - list keys", func(t *testing.T) {
		v := KVStorage{pathPrefix: prefix, client: mockVaultClient{store: map[string]map[string]interface{}{storagePath(prefix, kid): {"key": string(encodedSecret)}}}}
		result, err := v.ListKeys()
		assert.NoError(t, err)
		assert.Equal(t, []string{kid}, result)
	})

	t.Run("error - while listing", func(t *testing.T) {
		v := KVStorage{pathPrefix: prefix, client: mockVaultClient{err: vaultError}}
		_, err := v.ListKeys()
		assert.Error(t, err, "listing should fail")
		assert.ErrorIs(t, err, vaultError)
	})
}

func TestVaultKVStorage_DeleteKey(t *testing.T) {
	t.Run("ok - delete key", func(t *testing.T) {
		v := KVStorage{pathPrefix: prefix, client: mockVaultClient{store: map[string]map[string]interface{}{storagePath(prefix, kid): {"key": string(encodedSecret)}}}}
		assert.NoError(t, v.DeleteSecret(kid))
		result, err := v.GetSecret(kid)
		assert.EqualError(t, err, ErrNotFound.Error(), "secret should not be found")
		assert.Nil(t, result, "result should be nil")
	})

	t.Run("error - while deleting", func(t *testing.T) {
		v := KVStorage{pathPrefix: prefix, client: mockVaultClient{err: vaultError}}
		err := v.DeleteSecret(kid)
		assert.Error(t, err, "deleting should fail")
		assert.ErrorIs(t, err, vaultError)
	})

	t.Run("error - key not found", func(t *testing.T) {
		v := KVStorage{pathPrefix: prefix, client: mockVaultClient{store: map[string]map[string]interface{}{}}}
		assert.EqualError(t, v.DeleteSecret(kid), ErrNotFound.Error())
	})
}
