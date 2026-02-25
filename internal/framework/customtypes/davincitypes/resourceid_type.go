// Copyright © 2026 Ping Identity Corporation

package davincitypes

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

type ResourceIDType struct {
	basetypes.StringType
}

func (t ResourceIDType) Equal(o attr.Type) bool {
	other, ok := o.(ResourceIDType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t ResourceIDType) String() string {
	return "pingonetypes.ResourceIDType"
}

func (t ResourceIDType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	// ResourceIDValue defined in the value type section
	value := ResourceIDValue{
		StringValue: in,
	}

	return value, nil
}

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

func (t ResourceIDType) ValueType(ctx context.Context) attr.Value {
	// ResourceIDValue defined in the value type section
	return ResourceIDValue{}
}

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

	if !verify.P1DVResourceIDRegexpFullString.MatchString(valueString) {
		diags.AddAttributeError(
			path,
			"PingOne Resource ID Type Validation Error",
			fmt.Sprintf("The PingOne resource ID is malformed. Must match regex %q", verify.P1DVResourceIDRegexpFullString),
		)
	}

	return diags
}
