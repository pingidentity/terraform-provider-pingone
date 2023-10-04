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
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccGroupNesting_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group_nesting.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var groupNestingID, groupID, environmentID string

	var ctx = context.Background()
	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		t.Fatalf("Failed to get API client: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.TestAccCheckGroupNestingDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccGroupNestingConfig_Full(resourceName, name),
				Check:  sso.TestAccGetGroupNestingIDs(resourceFullName, &environmentID, &groupID, &groupNestingID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					sso.GroupNesting_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, groupID, groupNestingID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the group
			{
				Config: testAccGroupNestingConfig_Full(resourceName, name),
				Check:  sso.TestAccGetGroupNestingIDs(resourceFullName, &environmentID, &groupID, &groupNestingID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					sso.Group_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, groupID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccGroupNestingConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  sso.TestAccGetGroupNestingIDs(resourceFullName, &environmentID, &groupID, &groupNestingID),
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

func TestAccGroupNesting_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group_nesting.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.TestAccCheckGroupNestingDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupNestingConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "group_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "nested_group_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "type", "DIRECT"),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["group_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGroupNesting_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group_nesting.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.TestAccCheckGroupNestingDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccGroupNestingConfig_Full(resourceName, name),
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

func testAccGroupNestingConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_group" "%[3]s" {
	environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_group" "%[3]s-nesting" {
	environment_id = pingone_environment.%[2]s.id

  name = "%[4]s-nesting"
}

resource "pingone_group_nesting" "%[3]s" {
	environment_id = pingone_environment.%[2]s.id
  group_id        = pingone_group.%[3]s.id
  nested_group_id = pingone_group.%[3]s-nesting.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccGroupNestingConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_group" "%[2]s-nesting" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-nesting"
}

resource "pingone_group_nesting" "%[2]s" {
  environment_id  = data.pingone_environment.general_test.id
  group_id        = pingone_group.%[2]s.id
  nested_group_id = pingone_group.%[2]s-nesting.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
