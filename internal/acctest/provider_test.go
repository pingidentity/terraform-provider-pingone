// Copyright © 2026 Ping Identity Corporation

package acctest

import (
	"testing"

	"github.com/pingidentity/terraform-provider-pingone/internal/provider/sdkv2"
)

func TestProvider(t *testing.T) {
	if err := sdkv2.New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
