package vault

import (
	"encoding/base64"
	"fmt"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
	"net/url"
	"path/filepath"
)

type KVStorage struct {
	config Config
	client logicaler
}

// NewKVStore creates a new Vault backend using the kv version 1 secret engine: https://www.vaultproject.io/docs/secrets/kv
// It currently only supports token authentication which should be provided by the token param.
// If config.Address is empty, the VAULT_ADDR environment should be set.
// If config.Token is empty, the VAULT_TOKEN environment should be is set.
func NewKVStore(config Config) (Storage, error) {
	client, err := configureVaultClient(config)
	if err != nil {
		return nil, err
	}

	vaultStorage := KVStorage{client: client.Logical(), config: config}
	if err = vaultStorage.checkConnection(); err != nil {
		return nil, err
	}
	return vaultStorage, nil
}

func configureVaultClient(cfg Config) (*vaultapi.Client, error) {
	vaultConfig := vaultapi.DefaultConfig()
	vaultConfig.Timeout = cfg.Timeout
	client, err := vaultapi.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}
	// The Vault client will automatically use the env var VAULT_TOKEN
	// the client.SetToken overrides this value, so only set when not empty
	if cfg.Token != "" {
		client.SetToken(cfg.Token)
	}
	if cfg.Address != "" {
		if err = client.SetAddress(cfg.Address); err != nil {
			return nil, fmt.Errorf("vault address invalid: %w", err)
		}
	}
	logrus.Infof("Proxying to Vault at %s", client.Address)
	return client, nil
}

func (v KVStorage) checkConnection() error {
	// Perform a token introspection to test the connection. This should be allowed by the default vault token policy.
	logrus.Debug("Verifying Vault connection...")
	secret, err := v.client.Read("auth/token/lookup-self")
	if err != nil {
		return fmt.Errorf("unable to connect to Vault: unable to retrieve token status: %w", err)
	}
	if secret == nil || len(secret.Data) == 0 {
		return fmt.Errorf("could not read token information on auth/token/lookup-self")
	}
	logrus.Info("Connected to Vault.")
	return nil
}

func (v KVStorage) GetSecret(key string) ([]byte, error) {
	path := storagePath(v.config.PathPrefix, key)
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
		return nil, errKeyNotFound
	}
	rawValue, ok := result.Data[key]
	if !ok {
		return nil, errKeyNotFound
	}
	value, ok := rawValue.(string)
	if !ok {
		return nil, fmt.Errorf("unable to convert key result to string")
	}
	return base64.StdEncoding.DecodeString(value)
}
func (v KVStorage) storeValue(path, key string, value []byte) error {
	_, err := v.client.Write(path, map[string]interface{}{key: value})
	if err != nil {
		return fmt.Errorf("unable to write secret to vault: %w", err)
	}
	return nil
}

// ListSecretPaths returns a list of all keys in the vault storage
func (v KVStorage) ListKeys() ([]string, error) {
	path := privateKeyListPath(v.config.PathPrefix)
	response, err := v.client.ReadWithData(path, map[string][]string{"list": {"true"}})
	if err != nil {
		logrus.WithError(err).Error("Could not list private keys in Vault")
		return nil, err
	}
	if response == nil {
		logrus.Warnf("Vault returned nothing while fetching private keys, maybe the path prefix ('%s') is incorrect or the engine doesn't exist?", v.config.PathPrefix)
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
	path := storagePath(v.config.PathPrefix, key)

	_, err := v.getValue(path, keyName)
	if err == errKeyNotFound {
		return v.storeValue(path, keyName, value)
	}
	if err != nil {
		return err
	}

	return errKeyAlreadyExists
}
