package config

import (
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type Provider string

const (
	EnvProvider   Provider = "env"
	FileProvider  Provider = "file"
	VaultProvider Provider = "vault"
)

// KeyProviderConfig defines the source of cryptographic keys
type KeyProviderConfig struct {
	Provider    Provider        `json:"provider"`
	VaultConfig *VaultKeyConfig `json:"vault_config,omitempty"`
	FileConfig  *FileKeyConfig  `json:"file_config,omitempty"`
	EnvConfig   *EnvKeyConfig   `json:"env_config,omitempty"`
}

type VaultKeyConfig struct {
	Address   string `json:"address"`
	Path      string `json:"path"`
	TokenPath string `json:"token_path,omitempty"`
	RoleID    string `json:"role_id,omitempty"`
	SecretID  string `json:"secret_id,omitempty"`
	KeyName   string `json:"key_name"`
}

type FileKeyConfig struct {
	KeyPath        string `json:"key_path"`
	PrivateKeyPath string `json:"private_key_path"`
	PublicKeyPath  string `json:"public_key_path"`
	Permissions    string `json:"permissions,omitempty"`
	AutoCreate     bool   `json:"auto_create,omitempty"`
}

type EnvKeyConfig struct {
	PrivateKeyVar string `json:"private_key_var"`
	PublicKeyVar  string `json:"public_key_var"`
}

// Config provides complete configuration for the encryption system
type Config struct {
	// General settings
	IsServer          bool `json:"is_server"`
	RequireEncryption bool `json:"require_encryption"`
	DisableEncryption bool `json:"disable_encryption,omitempty"`

	// Timeout settings
	KeyExchangeTimeout time.Duration `json:"key_exchange_timeout"`
	SessionTimeout     time.Duration `json:"session_timeout"`

	// Key rotation
	EnableKeyRotation   bool          `json:"enable_key_rotation"`
	KeyRotationInterval time.Duration `json:"key_rotation_interval,omitempty"`

	// EncryptionProtocol settings
	EncryptionProtocol types.EncryptionProtocol `json:"encryption_protocol"`
	// KeyExchangeProtocolOptions  []keyexchange.ProtocolFactoryOption
	EncryptionFallbackProtocols []types.EncryptionProtocol `json:"encryption_fallback_protocols,omitempty"`

	// Key management
	KeyProvider KeyProviderConfig `json:"key_provider"`

	// Security settings
	ReplayProtection  bool          `json:"replay_protection"`
	NonceReplayWindow time.Duration `json:"nonce_replay_window,omitempty"`
}

// DefaultConfig provides sensible defaults for the encryption system
func DefaultConfig() Config {
	return Config{
		RequireEncryption:   true,
		DisableEncryption:   false,
		KeyExchangeTimeout:  30 * time.Second,
		SessionTimeout:      24 * time.Hour,
		EnableKeyRotation:   true,
		KeyRotationInterval: 1 * time.Hour,
		EncryptionProtocol:  types.ProtocolV2,
		// KeyExchangeProtocolOptions:  []keyexchange.ProtocolFactoryOption{keyexchange.WithKeySignature()},
		EncryptionFallbackProtocols: []types.EncryptionProtocol{types.ProtocolV1},
		KeyProvider: KeyProviderConfig{
			Provider: EnvProvider,
			EnvConfig: &EnvKeyConfig{
				PrivateKeyVar: "SERVER_ENCRYPT_PRIV_KEY",
				PublicKeyVar:  "SERVER_ENCRYPT_PUB_KEY",
			},
		},
		ReplayProtection:  true,
		NonceReplayWindow: 5 * time.Minute,
	}
}

// ValidateConfig checks if a configuration is valid and complete
func ValidateConfig(config Config) error {
	if config.KeyExchangeTimeout < 5*time.Second {
		return encryptionerr.ErrInvalidConfig
	}

	if config.RequireEncryption && config.DisableEncryption {
		return encryptionerr.ErrInvalidConfig
	}

	if config.EnableKeyRotation && config.KeyRotationInterval < time.Minute {
		return encryptionerr.ErrInvalidConfig
	}

	// Validate key provider configuration
	switch config.KeyProvider.Provider {
	case VaultProvider:
		if config.KeyProvider.VaultConfig == nil {
			return encryptionerr.ErrInvalidProvider
		}

	case FileProvider:
		if config.KeyProvider.FileConfig == nil {
			return encryptionerr.ErrInvalidProvider
		}

	case EnvProvider:
		if config.KeyProvider.EnvConfig == nil {
			return encryptionerr.ErrInvalidProvider
		}

	default:
		return encryptionerr.ErrInvalidProvider
	}

	return nil
}
