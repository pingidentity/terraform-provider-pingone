package stringvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = shouldNotContainValidator{}

// shouldNotContainValidator validates that the value does not match one of the values.
type shouldNotContainValidator struct {
	values []types.String
}

func (v shouldNotContainValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v shouldNotContainValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must not contain: %s", v.values)
}

func (v shouldNotContainValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	for _, otherValue := range v.values {
		if !strings.Contains(value.ValueString(), otherValue.ValueString()) {
			continue
		}

		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value.String(),
		))

		break
	}
}

// ShouldNotContain checks that the String held in the attribute
// is none of the given `values`.
func ShouldNotContain(values ...string) validator.String {
	frameworkValues := make([]types.String, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.StringValue(value))
	}

	return shouldNotContainValidator{
		values: frameworkValues,
	}
}
