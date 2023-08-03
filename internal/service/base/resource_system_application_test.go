package base_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckSystemApplicationDestroy(s *terraform.State) error {
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

func TestAccSystemApplication_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_system_application.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSystemApplicationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSystemApplicationConfig_NewEnv(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccSystemApplication_Portal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_system_application.%s", resourceName)

	applicationType := string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL)

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "type", applicationType),
		resource.TestCheckResourceAttr(resourceFullName, "name", "PingOne Application Portal"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "2"),
		resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.1", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ALL_GROUPS"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "type", applicationType),
		resource.TestCheckResourceAttr(resourceFullName, "name", "PingOne Application Portal"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "access_control_role_type"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSystemApplicationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccSystemApplicationConfig_Full(resourceName, applicationType, false),
				Check:  fullCheck,
			},
			{
				Config:  testAccSystemApplicationConfig_Full(resourceName, applicationType, false),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccSystemApplicationConfig_Minimal(resourceName, applicationType),
				Check:  minimalCheck,
			},
			{
				Config:  testAccSystemApplicationConfig_Minimal(resourceName, applicationType),
				Destroy: true,
			},
			// Change
			{
				Config: testAccSystemApplicationConfig_Full(resourceName, applicationType, false),
				Check:  fullCheck,
			},
			{
				Config: testAccSystemApplicationConfig_Minimal(resourceName, applicationType),
				Check:  minimalCheck,
			},
			{
				Config: testAccSystemApplicationConfig_Full(resourceName, applicationType, false),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccSystemApplication_SelfService(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_system_application.%s", resourceName)

	applicationType := string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE)

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "type", applicationType),
		resource.TestCheckResourceAttr(resourceFullName, "name", "PingOne Self-Service - MyAccount"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "2"),
		resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.1", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ALL_GROUPS"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "type", applicationType),
		resource.TestCheckResourceAttr(resourceFullName, "name", "PingOne Self-Service - MyAccount"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "access_control_role_type"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSystemApplicationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccSystemApplicationConfig_Full(resourceName, applicationType, true),
				Check:  fullCheck,
			},
			{
				Config:  testAccSystemApplicationConfig_Full(resourceName, applicationType, true),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccSystemApplicationConfig_Minimal(resourceName, applicationType),
				Check:  minimalCheck,
			},
			{
				Config:  testAccSystemApplicationConfig_Minimal(resourceName, applicationType),
				Destroy: true,
			},
			// Change
			{
				Config: testAccSystemApplicationConfig_Full(resourceName, applicationType, true),
				Check:  fullCheck,
			},
			{
				Config: testAccSystemApplicationConfig_Minimal(resourceName, applicationType),
				Check:  minimalCheck,
			},
			{
				Config: testAccSystemApplicationConfig_Full(resourceName, applicationType, true),
				Check:  fullCheck,
			},
			{
				Config:  testAccSystemApplicationConfig_Full(resourceName, applicationType, true),
				Destroy: true,
			},
		},
	})
}

func testAccSystemApplicationConfig_NewEnv(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_system_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  type    = "PING_ONE_PORTAL"
  enabled = false
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccSystemApplicationConfig_Full(resourceName, applicationType string, enabled bool) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_group" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s-1"
}

resource "pingone_group" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s-2"
}

resource "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type    = "%[3]s"
  enabled = %[4]t

  access_control_role_type = "ADMIN_USERS_ONLY"
  access_control_group_options = {
    groups = [
      pingone_group.%[2]s-2.id,
      pingone_group.%[2]s-1.id,
    ]

    type = "ALL_GROUPS"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, applicationType, enabled)
}

func testAccSystemApplicationConfig_Minimal(resourceName, applicationType string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type    = "%[3]s"
  enabled = true
}`, acctest.GenericSandboxEnvironment(), resourceName, applicationType)
}
