package sso_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckGroupDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_group" {
			continue
		}

		_, r, err := apiClient.ManagementAPIsGroupsApi.ReadOneGroup(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("PingOne Population Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGroup_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group.%s", resourceName)

	name := resourceName
	description := "Test description"

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	userFilter := `email ew "@test.com"`
	externalID := "external_1234"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckGroupDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig_Full(resourceName, name, description, licenseID, region, userFilter, externalID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					//resource.TestCheckResourceAttr(resourceFullName, "population_id", populationID),
					resource.TestCheckResourceAttr(resourceFullName, "user_filter", userFilter),
					resource.TestCheckResourceAttr(resourceFullName, "external_id", externalID),
				),
			},
		},
	})
}

func TestAccGroup_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group.%s", resourceName)

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckGroupDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig_Minimal(resourceName, name, licenseID, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func testAccGroupConfig_Full(resourceName, name, description, licenseID, region, userFilter, externalID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[2]s"
			type = "SANDBOX"
			license_id = "%[4]s"
			region = "%[5]s"
			default_population {}
			service {}
		}

		resource "pingone_group" "%[1]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			name = "%[2]s"
			description = "%[3]s"
			population_id = "${pingone_environment.%[1]s.default_population_id}"
			user_filter = %[6]q
			external_id = "%[7]s"
		}`, resourceName, name, description, licenseID, region, userFilter, externalID)
}

func testAccGroupConfig_Minimal(resourceName, name, licenseID, region string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[2]s"
			type = "SANDBOX"
			license_id = "%[3]s"
			region = "%[4]s"
			default_population {}
			service {}
		}

		resource "pingone_group" "%[1]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			name = "%[2]s"
		}`, resourceName, name, licenseID, region)
}
