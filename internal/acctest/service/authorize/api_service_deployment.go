package authorize

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func APIServiceDeployment_CheckDestroy(s *terraform.State) error {
	return nil
}

func APIServiceDeployment_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*environmentID = rs.Primary.Attributes["environment_id"]
		*resourceID = rs.Primary.Attributes["api_service_id"]

		return nil
	}
}
