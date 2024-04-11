package sso_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccApplicationFlowPolicyAssignment_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_flow_policy_assignment.%s", resourceName)

	// environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	// licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var flowPolicyAssignmentID, applicationID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationFlowPolicyAssignment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccApplicationFlowPolicyAssignmentConfig_Single(resourceName, name),
				Check:  sso.ApplicationFlowPolicyAssignment_GetIDs(resourceFullName, &environmentID, &applicationID, &flowPolicyAssignmentID),
			},
			{
				PreConfig: func() {
					sso.ApplicationFlowPolicyAssignment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, applicationID, flowPolicyAssignmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the application
			{
				Config: testAccApplicationFlowPolicyAssignmentConfig_Single(resourceName, name),
				Check:  sso.ApplicationFlowPolicyAssignment_GetIDs(resourceFullName, &environmentID, &applicationID, &flowPolicyAssignmentID),
			},
			{
				PreConfig: func() {
					sso.Application_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, applicationID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			// {
			// 	Config: testAccApplicationFlowPolicyAssignmentConfig_Single(environmentName, licenseID, resourceName, name),
			// 	Check:  sso.TestAccGetApplicationFlowPolicyAssignmentIDs(resourceFullName, &environmentID, &applicationID, &flowPolicyAssignmentID),
			// },
			// {
			// 	PreConfig: func() {
			// 		base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
			// 	},
			// 	RefreshState:       true,
			// 	ExpectNonEmptyPlan: true,
			// },
		},
	})
}

func TestAccApplicationFlowPolicyAssignment_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_flow_policy_assignment.%s", resourceName)

	name := resourceName

	singleStep := resource.TestStep{
		Config: testAccApplicationFlowPolicyAssignmentConfig_Single(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "flow_policy_id", verify.P1DVResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "priority", "1"),
		),
	}

	singleStepChange := resource.TestStep{
		Config: testAccApplicationFlowPolicyAssignmentConfig_Change(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "flow_policy_id", verify.P1DVResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "priority", "1"),
		),
	}

	multipleStep := resource.TestStep{
		Config: testAccApplicationFlowPolicyAssignmentConfig_Multiple(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "flow_policy_id", verify.P1DVResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "priority", "2"),
			resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "application_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "flow_policy_id", verify.P1DVResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "priority", "1"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationFlowPolicyAssignment_CheckDestroy,
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
			singleStepChange,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplicationFlowPolicyAssignment_SystemApplication(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_flow_policy_assignment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationFlowPolicyAssignment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationFlowPolicyAssignmentConfig_SystemApplication(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "flow_policy_id", verify.P1DVResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "priority", "1"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplicationFlowPolicyAssignment_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_flow_policy_assignment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationFlowPolicyAssignment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccApplicationFlowPolicyAssignmentConfig_Single(resourceName, name),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
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

  oidc_options = {
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

func testAccApplicationFlowPolicyAssignmentConfig_Change(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
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

  oidc_options = {
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

  priority = 2
}

resource "pingone_application_flow_policy_assignment" "%[2]s-2" {
  environment_id = data.pingone_environment.davinci_test.id
  application_id = pingone_application.%[2]s.id

  flow_policy_id = data.pingone_flow_policies.%[2]s.ids[1]

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
