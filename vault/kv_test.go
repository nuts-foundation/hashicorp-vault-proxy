package vault

import (
	"encoding/base64"
	"errors"
	"net/url"
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

var encodedKid = url.PathEscape(kid)
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
		assert.Equal(t, []string{encodedKid}, result)
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

	// func TestVaultKVStorage_CheckHealth(t *testing.T) {
	// 	t.Run("ok", func(t *testing.T) {
	// 		vaultStorage := KVStorage{config: DefaultConfig(), client: mockVaultClient{store: map[string]map[string]interface{}{"auth/token/lookup-self": {"key": []byte("foo")}}}}
	// 		result := vaultStorage.CheckHealth()
	//
	// 		assert.Equal(t, core.HealthStatusUp, result[StorageType].Status)
	// 		assert.Empty(t, result[StorageType].Details)
	// 	})
	//
	// 	t.Run("error - lookup token endpoint returns empty response", func(t *testing.T) {
	// 		vaultStorage := KVStorage{config: DefaultConfig(), client: mockVaultClient{store: map[string]map[string]interface{}{}}}
	// 		result := vaultStorage.CheckHealth()
	//
	// 		assert.Equal(t, core.HealthStatusDown, result[StorageType].Status)
	// 		assert.Equal(t, "could not read token information on auth/token/lookup-self", result[StorageType].Details)
	// 	})
	// }
	//
	// func TestVaultKVStorage_ListPrivateKeys(t *testing.T) {
	// 	s := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	// 		writer.Write([]byte("{\"request_id\":\"d728876e-ea1e-8a58-f297-dcd4cd0a41bb\",\"lease_id\":\"\",\"renewable\":false,\"lease_duration\":0,\"data\":{\"keys\":[\"did:nuts:8AB7Jf8KYgNHC52sfyTTK2f2yGnDoSHkgzDgeqvrUBLo#45KSfeG71ZMh9NjGzSWFfcMsmu5587J93prf8Io1wf4\",\"did:nuts:8AB7Jf8KYgNHC52sfyTTK2f2yGnDoSHkgzDgeqvrUBLo#6Cc91cQQze7txdcEor_zkM4YSwX0kH1wsiMyeV9nedA\",\"did:nuts:8AB7Jf8KYgNHC52sfyTTK2f2yGnDoSHkgzDgeqvrUBLo#MaNou-G07aPD7oheretmI2C_VElG1XaHiqh89SlfkWQ\",\"did:nuts:8AB7Jf8KYgNHC52sfyTTK2f2yGnDoSHkgzDgeqvrUBLo#alt3OIpy21VxDlWao0jRumIyXi3qHBPG-ir5q8zdv8w\",\"did:nuts:8AB7Jf8KYgNHC52sfyTTK2f2yGnDoSHkgzDgeqvrUBLo#wumme98rwUOQVle-sT_MP3pRg_oqblvlanv3zYR2scc\",\"did:nuts:8AB7Jf8KYgNHC52sfyTTK2f2yGnDoSHkgzDgeqvrUBLo#yBLHNjVq_WM3qzsRQ_zi2yOcedjY9FfVfByp3HgEbR8\",\"did:nuts:8AB7Jf8KYgNHC52sfyTTK2f2yGnDoSHkgzDgeqvrUBLo#yREqK5id7I6SP1Iq7teThin2o53w17tb9sgEXZBIcDo\"]},\"wrap_info\":null,\"warnings\":null,\"auth\":null}"))
	// 	}))
	// 	defer s.Close()
	// 	storage, _ := NewVaultKVStorage(Config{Address: s.URL})
	// 	keys := storage.ListPrivateKeys()
	// 	assert.Len(t, keys, 7)
	// 	// Assert first and last entry, rest should be OK then
	// 	assert.Equal(t, "did:nuts:8AB7Jf8KYgNHC52sfyTTK2f2yGnDoSHkgzDgeqvrUBLo#45KSfeG71ZMh9NjGzSWFfcMsmu5587J93prf8Io1wf4", keys[0])
	// 	assert.Equal(t, "did:nuts:8AB7Jf8KYgNHC52sfyTTK2f2yGnDoSHkgzDgeqvrUBLo#yREqK5id7I6SP1Iq7teThin2o53w17tb9sgEXZBIcDo", keys[6])
	// }
	//
	// func Test_PrivateKeyPath(t *testing.T) {
	// 	t.Run("it removes dot-dot-slash paths from the kid", func(t *testing.T) {
	// 		assert.Equal(t, "kv/nuts-private-keys/did:nuts:123#abc", privateKeyPath("kv", "did:nuts:123#abc"))
	// 		assert.Equal(t, "kv/nuts-private-keys/did:nuts:123#abc", privateKeyPath("kv", "../did:nuts:123#abc"))
	// 		assert.Equal(t, "kv/nuts-private-keys/did:nuts:123#abc", privateKeyPath("kv", "/../did:nuts:123#abc"))
	// 	})
	// }
	//
	// func TestVaultKVStorage_configure(t *testing.T) {
	// 	t.Run("ok - configure a new vault store", func(t *testing.T) {
	// 		_, err := configureVaultClient(Config{
	// 			Token:   "tokenString",
	// 			Address: "http://localhost:123",
	// 		})
	// 		assert.NoError(t, err)
	// 	})
	//
	// 	t.Run("error - invalid address", func(t *testing.T) {
	// 		_, err := configureVaultClient(Config{
	// 			Token:   "tokenString",
	// 			Address: "%zzzzz",
	// 		})
	// 		assert.Error(t, err)
	// 		assert.EqualError(t, err, "vault address invalid: failed to set address: parse \"%zzzzz\": invalid URL escape \"%zz\"")
	// 	})
	// 	t.Run("VAULT_TOKEN not overriden by empty config", func(t *testing.T) {
	// 		t.Setenv("VAULT_TOKEN", "123")
	//
	// 		client, err := configureVaultClient(Config{
	// 			Address: "http://localhost:123",
	// 		})
	//
	// 		require.NoError(t, err)
	// 		assert.Equal(t, "123", client.Token())
	// 	})
	// }
	//
	// func TestNewVaultKVStorage(t *testing.T) {
	// 	t.Run("ok - data", func(t *testing.T) {
	// 		s := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	// 			writer.Write([]byte("{\"data\": {\"keys\":[]}}"))
	// 		}))
	// 		defer s.Close()
	// 		storage, err := NewVaultKVStorage(Config{Address: s.URL})
	// 		assert.NoError(t, err)
	// 		assert.NotNil(t, storage)
	// 	})
	//
	// 	t.Run("error - vault StatusUnauthorized", func(t *testing.T) {
	// 		s := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	// 			writer.WriteHeader(http.StatusUnauthorized)
	// 		}))
	// 		defer s.Close()
	// 		storage, err := NewVaultKVStorage(Config{Address: s.URL})
	// 		assert.Error(t, err)
	// 		assert.True(t, strings.HasPrefix(err.Error(), "unable to connect to Vault: unable to retrieve token status: Error making API request"))
	// 		assert.Nil(t, storage)
	// 	})
	//
	// 	t.Run("error - wrong URL", func(t *testing.T) {
	// 		storage, err := NewVaultKVStorage(Config{Address: "http://non-existing"})
	// 		require.Error(t, err)
	// 		assert.Regexp(t, `no such host|Temporary failure in name resolution`, err.Error())
	// 		assert.Nil(t, storage)
	// 	})
	// }
	//
	// func TestVaultKVStorage_checkConnection(t *testing.T) {
	// 	t.Run("ok", func(t *testing.T) {
	// 		vaultStorage := KVStorage{config: DefaultConfig(), client: mockVaultClient{store: map[string]map[string]interface{}{"auth/token/lookup-self": {"key": []byte("foo")}}}}
	// 		err := vaultStorage.checkConnection()
	// 		assert.NoError(t, err)
	// 	})
	//
	// 	t.Run("error - lookup token endpoint empty", func(t *testing.T) {
	// 		vaultStorage := KVStorage{config: DefaultConfig(), client: mockVaultClient{store: map[string]map[string]interface{}{}}}
	// 		err := vaultStorage.checkConnection()
	// 		assert.EqualError(t, err, "could not read token information on auth/token/lookup-self")
	// 	})
	//
	// 	t.Run("error - vault error while reading", func(t *testing.T) {
	// 		var vaultError = errors.New("vault error")
	// 		vaultStorage := KVStorage{client: mockVaultClient{err: vaultError}}
	// 		err := vaultStorage.checkConnection()
	// 		assert.EqualError(t, err, "unable to connect to Vault: unable to retrieve token status: vault error")
	// 	})
}
