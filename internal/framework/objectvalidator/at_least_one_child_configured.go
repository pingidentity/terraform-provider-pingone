// Copyright © 2025 Ping Identity Corporation

package objectvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Object = AtLeastOneChildConfiguredValidator{}

// AtLeastOneChildConfiguredValidator validates that if an object is configured,
// at least one of the specified nested attributes must be configured (not null).
// This prevents empty object blocks in Terraform configurations.
type AtLeastOneChildConfiguredValidator struct {
	AttributeNames []string
}

func (v AtLeastOneChildConfiguredValidator) Description(_ context.Context) string {
	return fmt.Sprintf("At least one of these attributes must be configured: %s", strings.Join(v.AttributeNames, ", "))
}

func (v AtLeastOneChildConfiguredValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v AtLeastOneChildConfiguredValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attributes := req.ConfigValue.Attributes()

	// Check if at least one of the specified attributes is configured (not null)
	for _, attrName := range v.AttributeNames {
		if attrValue, exists := attributes[attrName]; exists && !attrValue.IsNull() {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Missing Required Configuration",
		fmt.Sprintf("At least one of these attributes must be configured: %s", strings.Join(v.AttributeNames, ", ")),
	)
}

func AtLeastOneChildConfigured(attributeNames ...string) validator.Object {
	return AtLeastOneChildConfiguredValidator{
		AttributeNames: attributeNames,
	}
}
