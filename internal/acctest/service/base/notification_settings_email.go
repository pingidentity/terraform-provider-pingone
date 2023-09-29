package base

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCheckNotificationSettingsEmailDestroy(s *terraform.State) error {
	return nil
}

func TestAccGetNotificationSettingsEmailIDs(resourceName string, environmentID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}
