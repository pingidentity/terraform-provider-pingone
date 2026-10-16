// Copyright © 2026 Ping Identity Corporation

package pingcli_test

import (
	"os"
	"slices"
	"testing"

	"github.com/pingidentity/terraform-provider-pingone/internal/pingcli"
)

func TestMaskCredentialEnv(t *testing.T) {
	// t.Setenv restores every variable at test end.
	for _, name := range []string{
		"PINGONE_CLIENT_ID",
		"PINGONE_AUTH_GRANT_TYPE",
		"PINGONE_STORAGE_TYPE",
	} {
		t.Setenv(name, "value")
	}

	pingcli.MaskCredentialEnv()

	for _, name := range []string{
		"PINGONE_CLIENT_ID",
		"PINGONE_CLIENT_SECRET",
		"PINGONE_ENVIRONMENT_ID",
		"PINGONE_API_ACCESS_TOKEN",
		"PINGONE_AUTH_GRANT_TYPE",
		"PINGONE_STORAGE_TYPE",
		"PINGONE_STORAGE_NAME",
		"PINGONE_STORAGE_OPTIONS",
		"PINGONE_STORAGE_OPTIONAL_SUFFIX",
	} {
		if os.Getenv(name) != "" {
			t.Errorf("%s = %q after MaskCredentialEnv, want unset", name, os.Getenv(name))
		}
	}
}

func TestConflictEnvVarsStayCredentialVars(t *testing.T) {
	// Resolve iterates ConflictEnvVars for conflict detection; it must contain
	// only the four credential variables. The grant type and storage variables
	// are masked but must not make a Ping CLI configuration a hard conflict.
	want := []string{
		"PINGONE_CLIENT_ID",
		"PINGONE_CLIENT_SECRET",
		"PINGONE_ENVIRONMENT_ID",
		"PINGONE_API_ACCESS_TOKEN",
	}
	for _, name := range want {
		if !slices.Contains(pingcli.ConflictEnvVars, name) {
			t.Errorf("ConflictEnvVars is missing %s", name)
		}
	}
	if len(pingcli.ConflictEnvVars) != len(want) {
		t.Errorf("ConflictEnvVars = %v, want exactly %v", pingcli.ConflictEnvVars, want)
	}
}
