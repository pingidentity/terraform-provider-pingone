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

func TestAccEditorProcessor_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_editor_processor.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var fido2PolicyID, environmentID string

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
		CheckDestroy:             authorize.EditorProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccEditorProcessorConfig_Minimal(resourceName, name),
				Check:  authorize.EditorProcessor_GetIDs(resourceFullName, &environmentID, &fido2PolicyID),
			},
			{
				PreConfig: func() {
					authorize.EditorProcessor_RemovalDrift_PreConfig(ctx, p1Client.API.AuthorizeAPIClient, t, environmentID, fido2PolicyID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccEditorProcessorConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.EditorProcessor_GetIDs(resourceFullName, &environmentID, &fido2PolicyID),
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

func TestAccEditorProcessor_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_editor_processor.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.EditorProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEditorProcessorConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccEditorProcessor_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_editor_processor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test processor"),
		resource.TestCheckResourceAttr(resourceFullName, "full_name", name),
		resource.TestMatchResourceAttr(resourceFullName, "parent.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test child processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.expression", "$$.data.item2"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.value_type.type", "STRING"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckNoResourceAttr(resourceFullName, "description"),
		resource.TestCheckNoResourceAttr(resourceFullName, "full_name"),
		resource.TestCheckNoResourceAttr(resourceFullName, "parent"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.expression", "$$.data.item"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.value_type.type", "STRING"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.EditorProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccEditorProcessorConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccEditorProcessorConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccEditorProcessorConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccEditorProcessorConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccEditorProcessorConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccEditorProcessorConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccEditorProcessorConfig_Full(resourceName, name),
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

func TestAccEditorProcessor_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_editor_processor.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.EditorProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccEditorProcessorConfig_Minimal(resourceName, name),
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

func testAccEditorProcessorConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_editor_processor" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[3]s"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccEditorProcessorConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_editor_processor" "%[2]s-parent" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test parent processor"
    type = "JSON_PATH"

    expression = "$.data.item.parent"
    value_type = {
      type = "STRING"
    }
  }
}

resource "pingone_authorize_editor_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test processor"
  full_name      = "%[3]s"

  parent = {
    id = pingone_authorize_editor_processor.%[2]s-parent.id
  }

  processor = {
    name = "%[3]s Test child processor"
    type = "JSON_PATH"

    expression = "$.data.item2"
    value_type = {
      type = "STRING"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccEditorProcessorConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_editor_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "JSON_PATH"

    expression = "$.data.item"
    value_type = {
      type = "STRING"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
