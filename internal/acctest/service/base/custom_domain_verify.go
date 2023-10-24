package base

import (
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func CustomDomainVerify_CheckDestroy(s *terraform.State) error {
	return nil
}
