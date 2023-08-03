package sso_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckResourceScopeOpenIDDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	re, err := regexp.Compile(`^(address|email|openid|phone|profile)$`)
	if err != nil {
		return fmt.Errorf("Cannot compile regex check for predefined scopes.")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_resource_scope_openid" {
			continue
		}

		if m := re.MatchString(rs.Primary.Attributes["name"]); m {
			return nil
		} else {

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

			body, r, err := apiClient.ResourcesApi.ReadOneResource(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

			return fmt.Errorf("PingOne Resource scope Instance %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccGetMFAPolicyIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccMFAPolicy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	var resourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
				Check:  testAccGetMFAPolicyIDs(resourceFullName, &environmentID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.MFAAPIClient

					if environmentID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Resource ID: %s", environmentID, resourceID)
					}

					_, err = apiClient.DeviceAuthenticationPolicyApi.DeleteDeviceAuthenticationPolicy(ctx, environmentID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete MFA Policy: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceScopeOpenID_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_openid.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopeOpenIDDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeOpenIDConfig_Full(resourceName, name, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My resource scope"),
					resource.TestCheckResourceAttr(resourceFullName, "mapped_claims.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "mapped_claims.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "mapped_claims.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "mapped_claims.2", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccResourceScopeOpenID_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_openid.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopeOpenIDDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeOpenIDConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "mapped_claims.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceScopeOpenID_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_openid.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopeOpenIDDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeOpenIDConfig_Full(resourceName, name, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My resource scope"),
					resource.TestCheckResourceAttr(resourceFullName, "mapped_claims.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "mapped_claims.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "mapped_claims.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "mapped_claims.2", verify.P1ResourceIDRegexp),
				),
			},
			{
				Config: testAccResourceScopeOpenIDConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					//resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "mapped_claims.#", "0"),
				),
			},
			{
				Config: testAccResourceScopeOpenIDConfig_Full(resourceName, name, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My resource scope"),
					resource.TestCheckResourceAttr(resourceFullName, "mapped_claims.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "mapped_claims.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "mapped_claims.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "mapped_claims.2", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccResourceScopeOpenID_OverridePredefined(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_openid.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopeOpenIDDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeOpenIDConfig_OverridePredefined(environmentName, licenseID, resourceName, "email"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "mapped_claims.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "mapped_claims.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "mapped_claims.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "mapped_claims.2", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccResourceScopeOpenID_InvalidParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopeOpenIDDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceScopeOpenIDConfig_Full(resourceName, name, "email"),
				ExpectError: regexp.MustCompile("The scope `email` is an existing platform scope.  The description cannot be changed."),
			},
		},
	})
}

func testAccResourceScopeOpenIDConfig_Full(resourceName, attributeName, scopeName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource_attribute" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  resource_name  = "openid"

  name  = "%[3]s-1"
  value = "$${user.name.given}"
}

resource "pingone_resource_attribute" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  resource_name  = "openid"

  name  = "%[3]s-2"
  value = "$${user.name.family}"
}

resource "pingone_resource_attribute" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  resource_name  = "openid"

  name  = "%[3]s-3"
  value = "$${user.email}"
}

resource "pingone_resource_scope_openid" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name        = "%[4]s"
  description = "My resource scope"

  mapped_claims = [
    pingone_resource_attribute.%[2]s-2.id,
    pingone_resource_attribute.%[2]s-3.id,
    pingone_resource_attribute.%[2]s-1.id
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, attributeName, scopeName)
}

func testAccResourceScopeOpenIDConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource_scope_openid" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceScopeOpenIDConfig_OverridePredefined(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource_attribute" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id
  resource_name  = "openid"

  name  = "%[4]s-1"
  value = "$${user.name.given}"
}

resource "pingone_resource_attribute" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id
  resource_name  = "openid"

  name  = "%[4]s-2"
  value = "$${user.name.family}"
}

resource "pingone_resource_attribute" "%[3]s-3" {
  environment_id = pingone_environment.%[2]s.id
  resource_name  = "openid"

  name  = "%[4]s-3"
  value = "$${user.email}"
}

resource "pingone_resource_scope_openid" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"

  mapped_claims = [
    pingone_resource_attribute.%[3]s-2.id,
    pingone_resource_attribute.%[3]s-3.id,
    pingone_resource_attribute.%[3]s-1.id
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
