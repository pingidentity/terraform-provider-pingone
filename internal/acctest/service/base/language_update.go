package base

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCheckLanguageUpdateDestroy(s *terraform.State) error {
	return nil
}

func TestAccGetLanguageUpdateIDs(resourceName string, environmentID, languageID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*languageID = rs.Primary.Attributes["language_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}
