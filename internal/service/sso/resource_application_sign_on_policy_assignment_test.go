package sso_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccApplicationSignOnPolicyAssignment_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_sign_on_policy_assignment.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var signOnPolicyID, signOnPolicyAssignmentID, applicationID, environmentID string

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
		CheckDestroy:             sso.ApplicationSignOnPolicyAssignment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccApplicationSignOnPolicyAssignmentConfig_Single(resourceName, name),
				Check:  sso.ApplicationSignOnPolicyAssignment_GetIDs(resourceFullName, &environmentID, &applicationID, &signOnPolicyID, &signOnPolicyAssignmentID),
			},
			{
				PreConfig: func() {
					sso.ApplicationSignOnPolicyAssignment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, applicationID, signOnPolicyAssignmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the sign on policy
			{
				Config: testAccApplicationSignOnPolicyAssignmentConfig_Single(resourceName, name),
				Check:  sso.ApplicationSignOnPolicyAssignment_GetIDs(resourceFullName, &environmentID, &applicationID, &signOnPolicyID, &signOnPolicyAssignmentID),
			},
			{
				PreConfig: func() {
					sso.SignOnPolicy_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, signOnPolicyID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the application
			{
				Config: testAccApplicationSignOnPolicyAssignmentConfig_Single(resourceName, name),
				Check:  sso.ApplicationSignOnPolicyAssignment_GetIDs(resourceFullName, &environmentID, &applicationID, &signOnPolicyID, &signOnPolicyAssignmentID),
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
				Config: testAccApplicationSignOnPolicyAssignmentConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  sso.ApplicationSignOnPolicyAssignment_GetIDs(resourceFullName, &environmentID, &applicationID, &signOnPolicyID, &signOnPolicyAssignmentID),
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

func TestAccApplicationSignOnPolicyAssignment_Single(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_sign_on_policy_assignment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationSignOnPolicyAssignment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationSignOnPolicyAssignmentConfig_Single(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
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

func TestAccApplicationSignOnPolicyAssignment_Multiple(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_sign_on_policy_assignment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationSignOnPolicyAssignment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationSignOnPolicyAssignmentConfig_Multiple(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "priority", "2"),
					resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "priority", "1"),
				),
			},
		},
	})
}

func TestAccApplicationSignOnPolicyAssignment_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_sign_on_policy_assignment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationSignOnPolicyAssignment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationSignOnPolicyAssignmentConfig_Single(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "priority", "1"),
				),
			},
			{
				Config: testAccApplicationSignOnPolicyAssignmentConfig_Multiple(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "priority", "2"),
					resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "priority", "1"),
				),
			},
			{
				Config: testAccApplicationSignOnPolicyAssignmentConfig_Single(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "priority", "1"),
				),
			},
		},
	})
}

func TestAccApplicationSignOnPolicyAssignment_SystemApplication(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_sign_on_policy_assignment.%s", resourceName)

	name := resourceName

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationSignOnPolicyAssignment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationSignOnPolicyAssignmentConfig_SystemApplication(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
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

func TestAccApplicationSignOnPolicyAssignment_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_sign_on_policy_assignment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationSignOnPolicyAssignment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccApplicationSignOnPolicyAssignmentConfig_Single(resourceName, name),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/application_id/sign_on_policy_assignment_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/application_id/sign_on_policy_assignment_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/application_id/sign_on_policy_assignment_id" and must match regex: .*`),
			},
		},
	})
}

func testAccApplicationSignOnPolicyAssignmentConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
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

resource "pingone_sign_on_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_application_sign_on_policy_assignment" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = pingone_application.%[3]s.id

  sign_on_policy_id = pingone_sign_on_policy.%[3]s.id
  priority          = 1
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationSignOnPolicyAssignmentConfig_Single(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
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

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_application_sign_on_policy_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id
  priority          = 1
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationSignOnPolicyAssignmentConfig_Multiple(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
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

resource "pingone_sign_on_policy" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s_1"
}

resource "pingone_sign_on_policy" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s_2"
}

resource "pingone_application_sign_on_policy_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  sign_on_policy_id = pingone_sign_on_policy.%[2]s-1.id
  priority          = 2
}

resource "pingone_application_sign_on_policy_assignment" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  sign_on_policy_id = pingone_sign_on_policy.%[2]s-2.id
  priority          = 1
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationSignOnPolicyAssignmentConfig_SystemApplication(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_system_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  type           = "PING_ONE_PORTAL"
  enabled        = true
}

resource "pingone_sign_on_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_application_sign_on_policy_assignment" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = pingone_system_application.%[3]s.id

  sign_on_policy_id = pingone_sign_on_policy.%[3]s.id
  priority          = 1
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
