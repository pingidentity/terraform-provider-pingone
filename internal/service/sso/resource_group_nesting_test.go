package sso_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckGroupNestingDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_group_nesting" {
			continue
		}

		body, r, err := apiClient.GroupsApi.ReadOneGroupNesting(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["group_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Group Nesting Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGroupNesting_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group_nesting.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupNestingDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupNestingConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "group_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "nested_group_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "type", "DIRECT"),
				),
			},
		},
	})
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
