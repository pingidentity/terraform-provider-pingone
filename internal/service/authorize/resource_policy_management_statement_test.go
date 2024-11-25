package authorize_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccPolicyManagementStatement_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_policy_management_statement.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var statementID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.PolicyManagementStatement_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccPolicyManagementStatementConfig_Minimal(resourceName, name),
				Check:  authorize.PolicyManagementStatement_GetIDs(resourceFullName, &environmentID, &statementID),
			},
			{
				PreConfig: func() {
					authorize.PolicyManagementStatement_RemovalDrift_PreConfig(ctx, p1Client.API.AuthorizeAPIClient, t, environmentID, statementID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccPolicyManagementStatementConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.PolicyManagementStatement_GetIDs(resourceFullName, &environmentID, &statementID),
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

func TestAccPolicyManagementStatement_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_policy_management_statement.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.PolicyManagementStatement_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyManagementStatementConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccPolicyManagementStatement_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_policy_management_statement.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test statement full"),
		resource.TestCheckResourceAttr(resourceFullName, "code", "my statement"),
		resource.TestCheckResourceAttr(resourceFullName, "applies_to", "PERMIT"),
		resource.TestCheckResourceAttr(resourceFullName, "applies_if", "FINAL_DECISION_MATCHES"),
		resource.TestCheckResourceAttr(resourceFullName, "payload", "{\"foo\":\"bar\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "obligatory", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "attributes.#", "3"),
		resource.TestMatchResourceAttr(resourceFullName, "attributes.0.id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "attributes.1.id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "attributes.2.id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test statement"),
		resource.TestCheckResourceAttr(resourceFullName, "code", "my statement 1"),
		resource.TestCheckResourceAttr(resourceFullName, "applies_to", "DENY"),
		resource.TestCheckResourceAttr(resourceFullName, "applies_if", "PATH_MATCHES"),
		resource.TestCheckResourceAttr(resourceFullName, "payload", "{\"foo\":\"bar\",\"foo2\":\"bar2\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "obligatory", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "attributes.#", "2"),
		resource.TestMatchResourceAttr(resourceFullName, "attributes.0.id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "attributes.1.id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.PolicyManagementStatement_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccPolicyManagementStatementConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccPolicyManagementStatementConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccPolicyManagementStatementConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccPolicyManagementStatementConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccPolicyManagementStatementConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccPolicyManagementStatementConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccPolicyManagementStatementConfig_Full(resourceName, name),
				Check:  fullCheck,
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPolicyManagementStatement_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_policy_management_statement.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.PolicyManagementStatement_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccPolicyManagementStatementConfig_Minimal(resourceName, name),
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
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccPolicyManagementStatementConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-1" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[3]s-1"
  description    = "Test attribute"

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-2" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[3]s-2"
  description    = "Test attribute"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "test1"
    }
  ]

  value_type = {
    type = "JSON"
  }
}

resource "pingone_authorize_policy_management_statement" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[3]s"
  description    = "Test statement"

  code = "my statement 1"

  applies_to = "DENY"
  applies_if = "PATH_MATCHES"

  payload = jsonencode({
    "foo" : "bar",
    "foo2" : "bar2"
  })

  attributes = [
    {
      id = pingone_authorize_trust_framework_attribute.%[2]s-2.id
    },
    {
      id = pingone_authorize_trust_framework_attribute.%[2]s-1.id
    },
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPolicyManagementStatementConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"
  description    = "Test attribute"

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
  description    = "Test attribute"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "test1"
    }
  ]

  value_type = {
    type = "JSON"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-3"
  description    = "Test attribute"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "test"
    }
  ]

  value_type = {
    type = "JSON"
  }
}

resource "pingone_authorize_policy_management_statement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test statement full"

  code = "my statement"

  applies_to = "PERMIT"
  applies_if = "FINAL_DECISION_MATCHES"

  payload = jsonencode({
    "foo" : "bar"
  })

  obligatory = true

  attributes = [
    {
      id = pingone_authorize_trust_framework_attribute.%[2]s-1.id
    },
    {
      id = pingone_authorize_trust_framework_attribute.%[2]s-2.id
    },
    {
      id = pingone_authorize_trust_framework_attribute.%[2]s-3.id
    },
  ]
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccPolicyManagementStatementConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"
  description    = "Test attribute"

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
  description    = "Test attribute"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "test1"
    }
  ]

  value_type = {
    type = "JSON"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-3"
  description    = "Test attribute"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "test"
    }
  ]

  value_type = {
    type = "JSON"
  }
}

resource "pingone_authorize_policy_management_statement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test statement"

  code = "my statement 1"

  applies_to = "DENY"
  applies_if = "PATH_MATCHES"

  payload = jsonencode({
    "foo" : "bar",
    "foo2" : "bar2"
  })

  attributes = [
    {
      id = pingone_authorize_trust_framework_attribute.%[2]s-2.id
    },
    {
      id = pingone_authorize_trust_framework_attribute.%[2]s-1.id
    },
  ]
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}
