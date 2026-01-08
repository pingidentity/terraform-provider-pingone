// Copyright © 2025 Ping Identity Corporation

package pingcli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/pingidentity/pingone-go-client/config"
	svcOAuth2 "github.com/pingidentity/pingone-go-client/oauth2"
	"golang.org/x/oauth2"
)

// Config represents the structure for reading PingCLI configuration
type Config struct {
	koanf      *koanf.Koanf
	configPath string
}

// ProfileConfig holds the authentication configuration from a PingCLI profile
type ProfileConfig struct {
	// Grant type: "authorization_code", "device_code", "client_credentials"
	GrantType string

	// Common fields
	RegionCode    string
	EnvironmentID string

	// Client credentials fields
	ClientID     string
	ClientSecret string

	// Authorization code fields
	RedirectURIPath string
	RedirectURIPort string

	// OAuth scopes (optional)
	Scopes []string
}

// NewConfig creates a new PingCLI config reader
func NewConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, fmt.Errorf("config path cannot be empty")
	}

	// Expand ~ to home directory if present
	if strings.HasPrefix(configPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, configPath[2:])
	}

	// Check if file exists
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file does not exist: %s", configPath)
		}
		return nil, fmt.Errorf("error accessing config file %s: %w", configPath, err)
	}

	k := koanf.New(".")
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
	}

	return &Config{
		koanf:      k,
		configPath: configPath,
	}, nil
}

// GetActiveProfile returns the name of the active profile from the config
func (c *Config) GetActiveProfile() (string, error) {
	// Try both activeProfile (current PingCLI format) and activeprofile (legacy)
	if c.koanf.Exists("activeProfile") {
		activeProfile := c.koanf.String("activeProfile")
		if activeProfile == "" {
			return "", fmt.Errorf("active profile is empty in config file")
		}
		return activeProfile, nil
	}
	if c.koanf.Exists("activeprofile") {
		activeProfile := c.koanf.String("activeprofile")
		if activeProfile == "" {
			return "", fmt.Errorf("active profile is empty in config file")
		}
		return activeProfile, nil
	}
	return "", fmt.Errorf("no active profile found in config file")
}

// GetProfileConfig extracts authentication configuration for a given profile
func (c *Config) GetProfileConfig(profileName string) (*ProfileConfig, error) {
	if profileName == "" {
		return nil, fmt.Errorf("profile name cannot be empty")
	}

	// Check if profile exists
	if !c.koanf.Exists(profileName) {
		return nil, fmt.Errorf("profile '%s' not found in config file", profileName)
	}

	profile := &ProfileConfig{}

	// Get authentication type
	authTypeKey := fmt.Sprintf("%s.service.pingOne.authentication.type", profileName)
	if c.koanf.Exists(authTypeKey) {
		profile.GrantType = c.koanf.String(authTypeKey)
	} else {
		return nil, fmt.Errorf("authentication type not configured for profile '%s'", profileName)
	}

	// Get region code
	regionCodeKey := fmt.Sprintf("%s.service.pingOne.regionCode", profileName)
	if c.koanf.Exists(regionCodeKey) {
		profile.RegionCode = c.koanf.String(regionCodeKey)
	} else {
		return nil, fmt.Errorf("region code not configured for profile '%s'", profileName)
	}

	// Get environment ID - this is now centralized
	envIDKey := fmt.Sprintf("%s.service.pingOne.authentication.environmentID", profileName)
	if c.koanf.Exists(envIDKey) {
		profile.EnvironmentID = c.koanf.String(envIDKey)
	}

	// Hardcode scopes value for grant types. This may change later if we decide to allow users to send more scopes.
	profile.Scopes = []string{string(config.OIDCScopeOpenID)}

	// Get auth-type specific configuration
	switch profile.GrantType {
	case "client_credentials":
		clientIDKey := fmt.Sprintf("%s.service.pingOne.authentication.clientCredentials.clientID", profileName)
		clientSecretKey := fmt.Sprintf("%s.service.pingOne.authentication.clientCredentials.clientSecret", profileName)

		if c.koanf.Exists(clientIDKey) {
			profile.ClientID = c.koanf.String(clientIDKey)
		} else {
			return nil, fmt.Errorf("client ID not configured for client_credentials auth in profile '%s'", profileName)
		}

		if c.koanf.Exists(clientSecretKey) {
			profile.ClientSecret = c.koanf.String(clientSecretKey)
		} else {
			return nil, fmt.Errorf("client secret not configured for client_credentials auth in profile '%s'", profileName)
		}

		if !c.koanf.Exists(envIDKey) || profile.EnvironmentID == "" {
			return nil, fmt.Errorf("environment ID not configured for client_credentials auth in profile '%s'", profileName)
		}

	case "authorization_code":
		clientIDKey := fmt.Sprintf("%s.service.pingOne.authentication.authorizationCode.clientID", profileName)
		redirectURIPathKey := fmt.Sprintf("%s.service.pingOne.authentication.authorizationCode.redirectURIPath", profileName)
		redirectURIPortKey := fmt.Sprintf("%s.service.pingOne.authentication.authorizationCode.redirectURIPort", profileName)

		if c.koanf.Exists(clientIDKey) {
			profile.ClientID = c.koanf.String(clientIDKey)
		} else {
			return nil, fmt.Errorf("client ID not configured for authorization_code auth in profile '%s'", profileName)
		}

		if !c.koanf.Exists(envIDKey) || profile.EnvironmentID == "" {
			return nil, fmt.Errorf("environment ID not configured for authorization_code auth in profile '%s'", profileName)
		}

		if c.koanf.Exists(redirectURIPathKey) {
			profile.RedirectURIPath = c.koanf.String(redirectURIPathKey)
		}

		if c.koanf.Exists(redirectURIPortKey) {
			profile.RedirectURIPort = c.koanf.String(redirectURIPortKey)
		}

	case "device_code":
		clientIDKey := fmt.Sprintf("%s.service.pingOne.authentication.deviceCode.clientID", profileName)

		if c.koanf.Exists(clientIDKey) {
			profile.ClientID = c.koanf.String(clientIDKey)
		} else {
			return nil, fmt.Errorf("client ID not configured for device_code auth in profile '%s'", profileName)
		}

		if !c.koanf.Exists(envIDKey) || profile.EnvironmentID == "" {
			return nil, fmt.Errorf("environment ID not configured for device_code auth in profile '%s'", profileName)
		}

	case "worker":
		clientIDKey := fmt.Sprintf("%s.service.pingOne.authentication.worker.clientID", profileName)
		clientSecretKey := fmt.Sprintf("%s.service.pingOne.authentication.worker.clientSecret", profileName)

		if c.koanf.Exists(clientIDKey) {
			profile.ClientID = c.koanf.String(clientIDKey)
		} else {
			return nil, fmt.Errorf("client ID not configured for worker auth in profile '%s'", profileName)
		}

		workerAppEnvironmentID := fmt.Sprintf("%s.service.pingOne.authentication.worker.environmentID", profileName)
		if c.koanf.Exists(envIDKey) {
			profile.EnvironmentID = c.koanf.String(envIDKey)
		} else if c.koanf.Exists(workerAppEnvironmentID) {
			profile.EnvironmentID = c.koanf.String(workerAppEnvironmentID)
		} else {
			return nil, fmt.Errorf("environment ID not configured for worker auth in profile '%s'", profileName)
		}

		if c.koanf.Exists(clientSecretKey) {
			profile.ClientSecret = c.koanf.String(clientSecretKey)
		} else {
			return nil, fmt.Errorf("client secret not configured for worker auth in profile '%s'", profileName)
		}

	default:
		return nil, fmt.Errorf("unsupported authentication type '%s' for profile '%s'", profile.GrantType, profileName)
	}

	return profile, nil
}

func (c *Config) IsFileStorageEnabled(profileName string) bool {
	if profileName == "" {
		return false
	}

	storageTypeKey := fmt.Sprintf("%s.login.storage.type", profileName)
	if c.koanf.Exists(storageTypeKey) {
		v := strings.TrimSpace(strings.ToLower(c.koanf.String(storageTypeKey)))
		return v == "file_system"
	}

	// Default to secure local (keychain) if not set
	return false
}

// LoadProfileConfig is a convenience function that loads config and extracts profile configuration
// If profileName is empty, it will use the active profile
func LoadProfileConfig(configPath string, profileName string) (*ProfileConfig, error) {
	config, err := NewConfig(configPath)
	if err != nil {
		return nil, err
	}

	// If no profile name provided, get the active profile
	if profileName == "" {
		profileName, err = config.GetActiveProfile()
		if err != nil {
			return nil, fmt.Errorf("failed to get active profile: %w", err)
		}
	}

	return config.GetProfileConfig(profileName)
}

// tokenFileData represents the structure of the credentials file used by PingCLI
type tokenFileData struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

// LoadStoredToken attempts to load a stored token from PingCLI's file storage
// It scans the credentials directory for any valid token files matching the profile name
// This allows using tokens from any auth flow (authorization_code, device_code, client_credentials)
func LoadStoredToken(ctx context.Context, profileConfig *ProfileConfig, profileName string) (*oauth2.Token, error) {
	if profileConfig == nil {
		return nil, fmt.Errorf("profile config cannot be nil")
	}

	// 1) Try Keychain directly via SDK oauth2 keychain storage using the same token key derivation
	tflog.Debug(ctx, "[pingcli] Attempting to load token from keychain", map[string]any{
		"grantType":     profileConfig.GrantType,
		"environmentID": profileConfig.EnvironmentID,
		"region":        profileConfig.RegionCode,
		"profile":       profileName,
	})

	if profileConfig.GrantType != "" {
		// Build the optional suffix used by pingcli
		if profileName == "" {
			profileName = "default"
		}
		suffix := fmt.Sprintf("_pingone_%s_%s", strings.TrimSpace(profileConfig.GrantType), profileName)

		// Derive the exact keychain account name using the SDK helper
		envID := strings.TrimSpace(profileConfig.EnvironmentID)
		clientID := strings.TrimSpace(profileConfig.ClientID)
		grant := strings.TrimSpace(profileConfig.GrantType)
		if envID != "" && clientID != "" && grant != "" {
			account := svcOAuth2.GenerateKeychainAccountNameWithSuffix(envID, clientID, grant, suffix)
			ks, kerr := svcOAuth2.NewKeychainStorage("pingcli", account)
			if kerr == nil {
				if token, loadErr := ks.LoadToken(); loadErr == nil && token != nil && token.AccessToken != "" {
					tflog.Debug(ctx, "[pingcli] Keychain token loaded directly", map[string]any{
						"account":    account,
						"expires":    token.Expiry.Format(time.RFC3339),
						"hasRefresh": token.RefreshToken != "",
					})
					return token, nil
				}
			} else {
				tflog.Debug(ctx, "[pingcli] Keychain storage init failed", map[string]any{
					"account": account,
					"error":   kerr,
				})
			}
		} else {
			tflog.Debug(ctx, "[pingcli] Insufficient inputs for keychain account", map[string]any{
				"environmentID": envID,
				"clientID":      clientID,
				"grant":         grant,
			})
		}
	}

	// 2) File-based fallback to maintain compatibility when fileStorage is enabled
	// Get credentials directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	credentialsDir := filepath.Join(homeDir, ".pingcli", "credentials")
	tflog.Debug(ctx, "[pingcli] Checking file storage for tokens", map[string]any{
		"directory": credentialsDir,
	})

	if profileName == "" {
		profileName = "default"
	}

	// Scan directory for any token files matching this profile and provider
	// Token filename format: token-<hash>_<provider>_<grantType>_<profile>.json
	pattern := fmt.Sprintf("token-*_pingone_*_%s.json", profileName)
	glob := filepath.Join(credentialsDir, pattern)
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("failed to scan credentials directory: %w", err)
	}
	tflog.Debug(ctx, "[pingcli] Scanned credentials directory", map[string]any{
		"globPattern": glob,
		"matchCount":  len(matches),
	})

	// Try each matching file
	for _, filePath := range matches {
		data, err := os.ReadFile(filePath)
		if err != nil {
			tflog.Debug(ctx, "[pingcli] Failed to read token file", map[string]any{
				"path":  filePath,
				"error": err,
			})
			continue // Try next file
		}

		var tokenData tokenFileData
		if err := json.Unmarshal(data, &tokenData); err != nil {
			tflog.Debug(ctx, "[pingcli] Failed to unmarshal token file", map[string]any{
				"path":  filePath,
				"error": err,
			})
			continue // Try next file
		}

		token := &oauth2.Token{
			AccessToken:  tokenData.AccessToken,
			TokenType:    tokenData.TokenType,
			RefreshToken: tokenData.RefreshToken,
			Expiry:       tokenData.Expiry,
		}

		// Check if token is still valid
		if token.Valid() {
			tflog.Debug(ctx, "[pingcli] Valid file token found", map[string]any{
				"type":       token.TokenType,
				"expires":    token.Expiry.Format(time.RFC3339),
				"hasRefresh": token.RefreshToken != "",
			})
			return token, nil
		}
		duration := time.Until(token.Expiry)
		tflog.Debug(ctx, "[pingcli] Token invalid or expired", map[string]any{
			"expires":   token.Expiry.Format(time.RFC3339),
			"expiresIn": duration.String(),
		})
	}

	return nil, fmt.Errorf("no valid stored token found for profile '%s'; checked Keychain and %s. Ensure storage alignment and run `pingcli login` to obtain a token", profileName, credentialsDir)
}
