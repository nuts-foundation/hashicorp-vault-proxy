package vault

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"path/filepath"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

const keyName = "key"

type KVStorage struct {
	client     vaultClient
	pathPrefix string
}

// vaultClient is an interface which has been implemented by the mockVaultClient and real vault.Logical to allow testing vault without the server.
type vaultClient interface {
	Read(path string) (*vaultapi.Secret, error)
	Write(path string, data map[string]interface{}) (*vaultapi.Secret, error)
	List(path string) (*vaultapi.Secret, error)
	ReadWithData(path string, data map[string][]string) (*vaultapi.Secret, error)
	Delete(path string) (*vaultapi.Secret, error)
}

// NewKVStore creates a new Vault backend using the kv version 1 secret engine: https://www.vaultproject.io/docs/secrets/kv
// It currently only supports token authentication which should be provided by the token param.
// If config.Address is empty, the VAULT_ADDR environment should be set.
// If config.Token is empty, the VAULT_TOKEN environment should be is set.
func NewKVStore(pathPrefix string) (Storage, error) {
	client, err := configureVaultClient()
	if err != nil {
		return nil, err
	}

	vaultStorage := KVStorage{client: client.Logical(), pathPrefix: pathPrefix}
	if err = vaultStorage.Ping(); err != nil {
		return nil, err
	}
	return vaultStorage, nil
}

func configureVaultClient() (*vaultapi.Client, error) {
	vaultConfig := vaultapi.DefaultConfig()
	client, err := vaultapi.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}
	logrus.Infof("Proxying to Vault at %s", client.Address())
	return client, nil
}

func (v KVStorage) Ping() error {
	// Perform a token introspection to test the connection. This should be allowed by the default vault token policy.
	logrus.Debug("Verifying Vault connection...")
	secret, err := v.client.Read("auth/token/lookup-self")
	if err != nil {
		return fmt.Errorf("unable to connect to Vault: unable to retrieve token status: %w", err)
	}
	if secret == nil || len(secret.Data) == 0 {
		return fmt.Errorf("could not read token information on auth/token/lookup-self")
	}
	logrus.Debug("Vault connection verified")
	return nil
}

func (v KVStorage) GetSecret(key string) ([]byte, error) {
	path := storagePath(v.pathPrefix, key)
	value, err := v.getValue(path, keyName)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// getValue extracts a field with name as provided by the key param from the Vault response.
func (v KVStorage) getValue(path, key string) ([]byte, error) {
	result, err := v.client.Read(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read key from vault: %w", err)
	}
	if result == nil || result.Data == nil {
		return nil, ErrNotFound
	}
	rawValue, ok := result.Data[key]
	if !ok {
		return nil, ErrNotFound
	}
	value, ok := rawValue.(string)

	if !ok {
		return nil, fmt.Errorf("unable to convert key result to string")
	}
	return base64.StdEncoding.DecodeString(string(value))
}

func (v KVStorage) storeValue(path, key string, value []byte) error {
	encodedValue := base64.StdEncoding.EncodeToString(value)
	_, err := v.client.Write(path, map[string]interface{}{key: encodedValue})
	if err != nil {
		return fmt.Errorf("unable to write secret to vault: %w", err)
	}
	return nil
}

func (v KVStorage) DeleteSecret(key string) error {
	_, err := v.GetSecret(key)
	if err != nil {
		return err
	}
	path := storagePath(v.pathPrefix, key)
	_, err = v.client.Delete(path)
	if err != nil {
		return fmt.Errorf("unable to delete secret from vault: %w", err)
	}
	return nil
}

// ListKeys returns a list of all keys in the vault storage for the given path.
func (v KVStorage) ListKeys() ([]string, error) {
	path := privateKeyListPath(v.pathPrefix)
	response, err := v.client.List(path)
	if err != nil {
		logrus.WithError(err).Error("Could not list private keys in Vault")
		return nil, err
	}
	if response == nil {
		logrus.Warnf("Vault returned nothing while fetching private keys, maybe the path prefix ('%s') is incorrect or the engine doesn't exist?", v.pathPrefix)
		return nil, fmt.Errorf("vault returned nothing while fetching private keys")
	}
	keys, _ := response.Data["keys"].([]interface{})
	var result []string
	for _, key := range keys {
		keyStr, ok := key.(string)
		if ok {
			result = append(result, keyStr)
		}
	}
	return result, nil
}

// storagePath cleans the key by removing optional slashes and dots and constructs the key path
// This prevents “dot-dot-slash” aka “directory traversal” attacks.
func storagePath(prefix, key string) string {
	// Clean the key by encoding optional slashes and dots and prepend the prefix
	cleanKey := url.PathEscape(key)
	return filepath.Clean(fmt.Sprintf("%s/%s", prefix, filepath.Base(cleanKey)))
}

func privateKeyListPath(prefix string) string {
	path := fmt.Sprintf("%s", prefix)
	return filepath.Clean(path)
}

func (v KVStorage) StoreSecret(key string, value []byte) error {
	path := storagePath(v.pathPrefix, key)

	_, err := v.getValue(path, keyName)
	if err == ErrNotFound {
		return v.storeValue(path, keyName, value)
	}
	if err != nil {
		return err
	}

	return ErrKeyAlreadyExists
}
