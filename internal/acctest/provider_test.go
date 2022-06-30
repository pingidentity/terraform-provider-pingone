package acctest

import (
	"testing"

	"github.com/pingidentity/terraform-provider-pingone/internal/provider"
)

func TestProvider(t *testing.T) {
	if err := provider.New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
