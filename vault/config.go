package vault

import "time"

const defaultPathPrefix = "kv"
const keyName = "key"

// Config VaultConfig contains the config options to configure the VaultKVStorage backend
type Config struct {
	// Token to authenticate to the Vault cluster.
	Token string `koanf:"token"`
	// Address of the Vault cluster
	Address string `koanf:"address"`
	// PathPrefix can be used to overwrite the default 'kv' path.
	PathPrefix string `koanf:"pathprefix"`
	// Timeout specifies the Vault client timeout.
	Timeout time.Duration
}

// DefaultVaultConfig returns a VaultConfig with the PathPrefix containing the default value.
func DefaultVaultConfig() Config {
	return Config{
		PathPrefix: defaultPathPrefix,
		Timeout:    5 * time.Second,
	}
}
