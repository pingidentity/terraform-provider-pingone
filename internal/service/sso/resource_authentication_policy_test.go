package sso_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	pingone "github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckAuthenticationPolicyDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_authentication_policy" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if rEnv.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		body, r, err := apiClient.SignOnPoliciesSignOnPoliciesApi.ReadOneSignOnPolicy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Authentication Policy %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccAuthenticationPolicy_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authentication_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName
	description := "Test description"

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAuthenticationPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPolicyConfig_Full(environmentName, licenseID, resourceName, name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
				),
			},
		},
	})
}

func TestAccAuthenticationPolicy_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authentication_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAuthenticationPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPolicyConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
				),
			},
		},
	})
}

func testAccAuthenticationPolicyConfig_Full(environmentName, licenseID, resourceName, name, description string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"
			description = "%[5]s"

			policy_action {
				action_type = "LOGIN"
			}
		}`, environmentName, licenseID, resourceName, name, description)
}

func testAccAuthenticationPolicyConfig_Minimal(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "LOGIN"
			}
		}`, environmentName, licenseID, resourceName, name)
}
