package sweep_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	_ "github.com/pingidentity/terraform-provider-pingone/internal/service/base"
	_ "github.com/pingidentity/terraform-provider-pingone/internal/service/sso"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
