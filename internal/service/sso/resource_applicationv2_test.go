// Copyright © 2025 Ping Identity Corporation

package sso_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func TestAccApplicationV2_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_applicationv2.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var applicationID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccApplicationV2Config_OIDC_MinimalWeb(resourceName, name),
				Check:  sso.Application_GetIDs(resourceFullName, &environmentID, &applicationID),
			},
			{
				PreConfig: func() {
					sso.Application_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, applicationID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccApplicationV2Config_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  sso.Application_GetIDs(resourceFullName, &environmentID, &applicationID),
			},
			{
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccApplicationV2Config_OIDC_MinimalWeb(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_applicationv2" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  protocol = "OPENID_CONNECT"
    type                       = "WEB_APP"
    grant_types                = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
    redirect_uris              = ["https://www.pingidentity.com"]
  
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationV2Config_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_applicationv2" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  enabled        = true

  protocol = "OPENID_CONNECT"
    type                       = "WEB_APP"
    grant_types                = ["AUTHORIZATION_CODE", "REFRESH_TOKEN"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
    redirect_uris              = ["https://www.pingidentity.com"]
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
