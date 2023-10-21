package base

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func NotificationSettingsEmail_CheckDestroy(s *terraform.State) error {
	return nil
}

func NotificationSettingsEmail_GetIDs(resourceName string, environmentID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}
