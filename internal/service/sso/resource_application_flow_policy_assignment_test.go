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

func TestAccCheckApplicationFlowPolicyAssignmentDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

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

func TestAccGetApplicationFlowPolicyAssignmentIDs(resourceName string, environmentID, applicationID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*applicationID = rs.Primary.Attributes["application_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccApplicationFlowPolicyAssignment_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_flow_policy_assignment.%s", resourceName)

	name := resourceName

	var resourceID, applicationID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.TestAccCheckApplicationFlowPolicyAssignmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccApplicationFlowPolicyAssignmentConfig_Single(resourceName, name),
				Check:  sso.TestAccGetApplicationFlowPolicyAssignmentIDs(resourceFullName, &environmentID, &applicationID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient

					if environmentID == "" || applicationID == "" || resourceID == "" {
						t.Fatalf("One of environment ID, application ID or resource ID cannot be determined. Environment ID: %s, Application ID: %s, Resource ID: %s", environmentID, applicationID, resourceID)
					}

					_, err = apiClient.ApplicationFlowPolicyAssignmentsApi.DeleteFlowPolicyAssignment(ctx, environmentID, applicationID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete Application flow policy assignment: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the application
			{
				Config: testAccApplicationFlowPolicyAssignmentConfig_Single(resourceName, name),
				Check:  sso.TestAccGetApplicationFlowPolicyAssignmentIDs(resourceFullName, &environmentID, &applicationID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient

					if environmentID == "" || applicationID == "" || resourceID == "" {
						t.Fatalf("One of environment ID, application ID or resource ID cannot be determined. Environment ID: %s, Application ID: %s, Resource ID: %s", environmentID, applicationID, resourceID)
					}

					_, err = apiClient.ApplicationsApi.DeleteApplication(ctx, environmentID, applicationID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete Application: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
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
		CheckDestroy:             sso.TestAccCheckApplicationFlowPolicyAssignmentDestroy,
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
		CheckDestroy:             sso.TestAccCheckApplicationFlowPolicyAssignmentDestroy,
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
		CheckDestroy:             sso.TestAccCheckApplicationFlowPolicyAssignmentDestroy,
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

func testAccApplicationFlowPolicyAssignmentConfig_Change(resourceName, name string) string {
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
