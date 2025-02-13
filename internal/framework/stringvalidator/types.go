package stringvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

type CustomStringValidatorModel struct {
	Validators  []validator.String
	Description framework.SchemaAttributeDescription
}
