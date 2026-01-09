// Copyright © 2025 Ping Identity Corporation

package pingcli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
activeProfile: test-profile
test-profile:
  service:
    pingOne:
      authentication:
        type: client_credentials
        environmentID: env-123
        clientCredentials:
          clientID: client-123
          clientSecret: secret-123
      regionCode: NA
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Test successful load
	cfg, err := NewConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, configPath, cfg.configPath)

	// Test non-existent file
	_, err = NewConfig("non-existent.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config file does not exist")

	// Test empty path
	_, err = NewConfig("")
	assert.Error(t, err)
	assert.Equal(t, "config path cannot be empty", err.Error())
}

func TestGetActiveProfile(t *testing.T) {
	tests := []struct {
		name          string
		configContent string
		expected      string
		expectError   bool
	}{
		{
			name: "activeProfile standard",
			configContent: `
activeProfile: standard
`,
			expected: "standard",
		},
		{
			name: "activeprofile legacy",
			configContent: `
activeprofile: legacy
`,
			expected: "legacy",
		},
		{
			name: "no active profile",
			configContent: `
other: value
`,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			err := os.WriteFile(configPath, []byte(tc.configContent), 0644)
			require.NoError(t, err)

			cfg, err := NewConfig(configPath)
			require.NoError(t, err)

			profile, err := cfg.GetActiveProfile()
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, profile)
			}
		})
	}
}

func TestGetProfileConfig(t *testing.T) {
	tests := []struct {
		name          string
		configContent string
		profileName   string
		expected      *ProfileConfig
		expectError   bool
		errorContains string
	}{
		{
			name: "client credentials success",
			configContent: `
test-profile:
  service:
    pingOne:
      regionCode: NA
      authentication:
        type: client_credentials
        environmentID: env-id
        clientCredentials:
          clientID: client-id
          clientSecret: client-secret
`,
			profileName: "test-profile",
			expected: &ProfileConfig{
				GrantType:     "client_credentials",
				RegionCode:    "NA",
				EnvironmentID: "env-id",
				ClientID:      "client-id",
				ClientSecret:  "client-secret",
				Scopes:        []string{"openid"},
			},
		},
		{
			name: "authorization code success",
			configContent: `
auth-code:
  service:
    pingOne:
      regionCode: EU
      authentication:
        type: authorization_code
        environmentID: env-id
        authorizationCode:
          clientID: client-id
          redirectURIPath: /callback
          redirectURIPort: 8080
`,
			profileName: "auth-code",
			expected: &ProfileConfig{
				GrantType:       "authorization_code",
				RegionCode:      "EU",
				EnvironmentID:   "env-id",
				ClientID:        "client-id",
				RedirectURIPath: "/callback",
				RedirectURIPort: "8080",
				Scopes:          []string{"openid"},
			},
		},
		{
			name: "missing profile",
			configContent: `
other: value
`,
			profileName:   "missing",
			expectError:   true,
			errorContains: "not found",
		},
		{
			name: "missing auth type",
			configContent: `
invalid:
  service:
    pingOne:
      regionCode: NA
`,
			profileName:   "invalid",
			expectError:   true,
			errorContains: "authentication type not configured",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			err := os.WriteFile(configPath, []byte(tc.configContent), 0644)
			require.NoError(t, err)

			cfg, err := NewConfig(configPath)
			require.NoError(t, err)

			pConfig, err := cfg.GetProfileConfig(tc.profileName)
			if tc.expectError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, pConfig)
			}
		})
	}
}

func TestLoadStoredToken_FileStorage(t *testing.T) {
	// Setup user home for file storage test
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome) // unix
	t.Setenv("USERPROFILE", tmpHome) // windows

	credsDir := filepath.Join(tmpHome, ".pingcli", "credentials")
	err := os.MkdirAll(credsDir, 0755)
	require.NoError(t, err)

	profileName := "test-profile-file"
	tokenPattern := fmt.Sprintf("token-hash_pingone_client_credentials_%s.json", profileName)
	tokenPath := filepath.Join(credsDir, tokenPattern)

	// Create valid token file
	validTokenContent := `{"access_token":"valid-token","token_type":"Bearer","expiry":"` + time.Now().Add(1 * time.Hour).Format(time.RFC3339) + `"}`
	err = os.WriteFile(tokenPath, []byte(validTokenContent), 0644)
	require.NoError(t, err)

	profileConfig := &ProfileConfig{
		GrantType: "client_credentials",
	}

	// Test successful load
	// We pass a background context as LoadStoredToken now requires context
	token, err := LoadStoredToken(context.Background(), profileConfig, profileName)
	require.NoError(t, err)
	require.NotNil(t, token)
	assert.Equal(t, "valid-token", token.AccessToken)

	// Test expired token
	expiredTokenContent := `{"access_token":"expired-token","token_type":"Bearer","expiry":"` + time.Now().Add(-1 * time.Hour).Format(time.RFC3339) + `"}`
	err = os.WriteFile(tokenPath, []byte(expiredTokenContent), 0644)
	require.NoError(t, err)

	token, err = LoadStoredToken(context.Background(), profileConfig, profileName)
	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Error(), "no valid stored token found")
}

func TestIsFileStorageEnabled(t *testing.T) {
	tests := []struct {
		name          string
		configContent string
		expected      bool
	}{
		{
			name: "file storage explicitly enabled",
			configContent: `
test-profile:
  login:
    storage:
      type: file_system
`,
			expected: true,
		},
		{
			name: "file storage case insensitive",
			configContent: `
test-profile:
  login:
    storage:
      type: FILE_SYSTEM
`,
			expected: true,
		},
		{
			name: "storage type other",
			configContent: `
test-profile:
  login:
    storage:
      type: keychain
`,
			expected: false,
		},
		{
			name: "storage type missing",
			configContent: `
test-profile:
  other: value
`,
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			err := os.WriteFile(configPath, []byte(tc.configContent), 0644)
			require.NoError(t, err)

			cfg, err := NewConfig(configPath)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, cfg.IsFileStorageEnabled("test-profile"))
		})
	}
}
