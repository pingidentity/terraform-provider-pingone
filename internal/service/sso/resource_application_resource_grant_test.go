package sso_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckApplicationResourceGrantDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_application_resource_grant" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.ApplicationResourceGrantsApi.ReadOneApplicationGrant(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID).Execute()

		if err != nil {

			if r == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if r.StatusCode == 404 {
				continue
			}

			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Application Role Assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccApplicationResourceGrant_OpenIDResource(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_resource_grant.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckApplicationResourceGrantDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationResourceGrantConfig_OpenIDResource(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexp),
				),
			},
			// Test error catch on update
			{
				Config:      testAccApplicationResourceGrantConfig_OpenIDResource_InvalidOpenIDScope(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid scope`),
			},
			{
				Config:  testAccApplicationResourceGrantConfig_OpenIDResource(resourceName, name),
				Destroy: true,
			},
			// Test error catch on from new
			{
				Config:      testAccApplicationResourceGrantConfig_OpenIDResource_InvalidOpenIDScope(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid scope`),
			},
		},
	})
}

func TestAccApplicationResourceGrant_CustomResource(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_resource_grant.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckApplicationResourceGrantDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationResourceGrantConfig_CustomResource(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.2", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccApplicationResourceGrant_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_resource_grant.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckApplicationResourceGrantDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationResourceGrantConfig_OpenIDResource(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexp),
				),
			},
			{
				Config: testAccApplicationResourceGrantConfig_CustomResource(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.2", verify.P1ResourceIDRegexp),
				),
			},
			{
				Config: testAccApplicationResourceGrantConfig_OpenIDResource(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func testAccApplicationResourceGrantConfig_OpenIDResource(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "openid"
}

data "pingone_resource_scope" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id

  name = "email"
}

data "pingone_resource_scope" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id

  name = "profile"
}

resource "pingone_application_resource_grant" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  resource_id = data.pingone_resource.%[2]s.id
  scopes = [
    data.pingone_resource_scope.%[2]s-1.id,
    data.pingone_resource_scope.%[2]s-2.id,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationResourceGrantConfig_OpenIDResource_InvalidOpenIDScope(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "openid"
}

data "pingone_resource_scope" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id

  name = "email"
}

data "pingone_resource_scope" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id

  name = "profile"
}

data "pingone_resource_scope" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id

  name = "openid"
}

resource "pingone_application_resource_grant" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  resource_id = data.pingone_resource.%[2]s.id
  scopes = [
    data.pingone_resource_scope.%[2]s-1.id,
    data.pingone_resource_scope.%[2]s-2.id,
    data.pingone_resource_scope.%[2]s-3.id,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationResourceGrantConfig_CustomResource(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_scope" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s-1"
}

resource "pingone_resource_scope" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s-2"
}

resource "pingone_resource_scope" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s-3"
}

resource "pingone_application_resource_grant" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  resource_id = pingone_resource.%[2]s.id
  scopes = [
    pingone_resource_scope.%[2]s-1.id,
    pingone_resource_scope.%[2]s-2.id,
    pingone_resource_scope.%[2]s-3.id
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
