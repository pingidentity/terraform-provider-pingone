// Copyright © 2025 Ping Identity Corporation

package sweep_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	_ "github.com/pingidentity/terraform-provider-pingone/internal/service/base"
	_ "github.com/pingidentity/terraform-provider-pingone/internal/service/davinci"
	_ "github.com/pingidentity/terraform-provider-pingone/internal/service/sso"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
