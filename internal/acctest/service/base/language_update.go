// Copyright © 2025 Ping Identity Corporation

package base

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func LanguageUpdate_CheckDestroy(s *terraform.State) error {
	return nil
}

func LanguageUpdate_GetIDs(resourceName string, environmentID, languageID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if languageID != nil {
			*languageID = rs.Primary.Attributes["language_id"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}
