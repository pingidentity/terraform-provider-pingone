package sso_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckApplicationFlowPolicyAssignmentDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_application_flow_policy_assignment" {
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

		body, r, err := apiClient.ApplicationFlowPolicyAssignmentsApi.ReadOneFlowPolicyAssignment(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Application Flow Policy assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccApplicationFlowPolicyAssignment_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_flow_policy_assignment.%s", resourceName)

	name := resourceName

	singleStep := resource.TestStep{
		Config: testAccApplicationFlowPolicyAssignmentConfig_Single(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "flow_policy_id", verify.P1DVResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "priority", "1"),
		),
	}

	multipleStep := resource.TestStep{
		Config: testAccApplicationFlowPolicyAssignmentConfig_Multiple(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "flow_policy_id", verify.P1DVResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "priority", "2"),
			resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "application_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "flow_policy_id", verify.P1DVResourceIDRegexp),
			resource.TestCheckResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "priority", "1"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckApplicationFlowPolicyAssignmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Single from new
			singleStep,
			{
				Config:  testAccApplicationFlowPolicyAssignmentConfig_Single(resourceName, name),
				Destroy: true,
			},
			// Multiple from new
			multipleStep,
			{
				Config:  testAccApplicationFlowPolicyAssignmentConfig_Multiple(resourceName, name),
				Destroy: true,
			},
			// Changes
			singleStep,
			multipleStep,
			singleStep,
		},
	})
}

func TestAccApplicationFlowPolicyAssignment_SystemApplication(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_flow_policy_assignment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckApplicationFlowPolicyAssignmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationFlowPolicyAssignmentConfig_SystemApplication(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "flow_policy_id", verify.P1DVResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "priority", "1"),
				),
			},
		},
	})
}

func testAccApplicationFlowPolicyAssignmentConfig_Single(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id
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

data "pingone_flow_policies" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id

  scim_filter = "(trigger.type eq \"AUTHENTICATION\")"
}

resource "pingone_application_flow_policy_assignment" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id
  application_id = pingone_application.%[2]s.id

  flow_policy_id = data.pingone_flow_policies.%[2]s.ids[0]

  priority = 1
}`, acctest.DaVinciFlowPolicySandboxEnvironment(), resourceName, name)
}

func testAccApplicationFlowPolicyAssignmentConfig_Multiple(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id
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

data "pingone_flow_policies" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id

  scim_filter = "(trigger.type eq \"AUTHENTICATION\")"
}

resource "pingone_application_flow_policy_assignment" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id
  application_id = pingone_application.%[2]s.id

  flow_policy_id = data.pingone_flow_policies.%[2]s.ids[1]

  priority = 2
}

resource "pingone_application_flow_policy_assignment" "%[2]s-2" {
  environment_id = data.pingone_environment.davinci_test.id
  application_id = pingone_application.%[2]s.id

  flow_policy_id = data.pingone_flow_policies.%[2]s.ids[0]

  priority = 1
}`, acctest.DaVinciFlowPolicySandboxEnvironment(), resourceName, name)
}

func testAccApplicationFlowPolicyAssignmentConfig_SystemApplication(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id
  type           = "PING_ONE_PORTAL"
  enabled        = true
}

data "pingone_flow_policies" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id

  scim_filter = "(trigger.type eq \"AUTHENTICATION\")"
}

resource "pingone_application_flow_policy_assignment" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id
  application_id = pingone_system_application.%[2]s.id

  flow_policy_id = data.pingone_flow_policies.%[2]s.ids[0]

  priority = 1
}`, acctest.DaVinciFlowPolicySandboxEnvironment(), resourceName, name)
}
