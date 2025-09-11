// Copyright © 2025 Ping Identity Corporation

// Package pingonetypes provides custom Terraform types specific to PingOne resources.
// This package contains type definitions, validation logic, and value handling for
// PingOne-specific data types used throughout the provider.
package pingonetypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Ensure the implementation satisfies the expected interfaces
var _ basetypes.StringTypable = ResourceIDType{}
var _ xattr.TypeWithValidate = ResourceIDType{}

// ResourceIDType represents a custom Terraform type for PingOne resource identifiers.
// It extends the base string type with PingOne-specific validation to ensure
// resource IDs match the expected format and structure required by the PingOne API.
type ResourceIDType struct {
	basetypes.StringType
}

// Equal determines if this ResourceIDType is equal to another attr.Type.
// It returns true if the other type is also a ResourceIDType with an equivalent base StringType.
func (t ResourceIDType) Equal(o attr.Type) bool {
	other, ok := o.(ResourceIDType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// String returns the string representation of the ResourceIDType.
// It returns a human-readable type identifier for debugging and error messages.
func (t ResourceIDType) String() string {
	return "pingonetypes.ResourceIDType"
}

// ValueFromString creates a ResourceIDValue from a base string value.
// It returns a ResourceIDValue that wraps the provided string with PingOne resource ID semantics.
// The ctx parameter provides the context for the conversion operation.
// The in parameter is the base string value to be converted to a ResourceIDValue.
func (t ResourceIDType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	// ResourceIDValue defined in the value type section
	value := ResourceIDValue{
		StringValue: in,
	}

	return value, nil
}

// ValueFromTerraform creates a ResourceIDValue from a Terraform value.
// It returns an attr.Value containing the converted ResourceIDValue with proper type validation.
// The ctx parameter provides the context for the conversion operation.
// The in parameter is the tftypes.Value from the Terraform state or configuration.
func (t ResourceIDType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

// ValueType returns the attr.Value type that this ResourceIDType represents.
// It returns a zero-value ResourceIDValue as a type representative.
// The ctx parameter provides the context for the type operation.
func (t ResourceIDType) ValueType(ctx context.Context) attr.Value {
	// ResourceIDValue defined in the value type section
	return ResourceIDValue{}
}

// Validate performs PingOne resource ID format validation on the provided value.
// It returns diagnostics containing any validation errors encountered during the check.
// The validation ensures the string matches the expected PingOne resource ID format using regex.
// Null and unknown values are considered valid and skip validation.
//
// The ctx parameter provides the context for the validation operation.
// The in parameter is the tftypes.Value to be validated.
// The path parameter specifies the attribute path for error reporting.
func (t ResourceIDType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Is(tftypes.String) {
		err := fmt.Errorf("expected String value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"PingOne Resource ID Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var valueString string

	if err := in.As(&valueString); err != nil {
		diags.AddAttributeError(
			path,
			"PingOne Resource ID Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)

		return diags
	}

	if !verify.P1ResourceIDRegexpFullString.MatchString(valueString) {
		diags.AddAttributeError(
			path,
			"PingOne Resource ID Type Validation Error",
			fmt.Sprintf("The PingOne resource ID is malformed. Must match regex %q", verify.P1ResourceIDRegexpFullString),
		)
	}

	return diags
}
