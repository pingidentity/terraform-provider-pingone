// Copyright © 2025 Ping Identity Corporation

package acctest

import (
	"testing"

	"github.com/pingidentity/terraform-provider-pingone/internal/provider/sdkv2"
)

// TestProvider validates the internal structure and configuration of the PingOne Terraform provider.
// This test ensures that the provider instance created for testing is properly configured
// and passes internal validation checks required for correct operation.
// The t parameter is the testing instance used for reporting test failures.
func TestProvider(t *testing.T) {
	if err := sdkv2.New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
