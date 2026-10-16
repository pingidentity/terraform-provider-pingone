// Copyright © 2026 Ping Identity Corporation

package sdkv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ConfigureForTest exposes the provider configure function to tests in this
// package's _test external test package. It lives in an export_test.go so the
// identifier is excluded from production builds.
func ConfigureForTest(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return configure("test")(ctx, d)
}
