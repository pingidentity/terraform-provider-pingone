// Copyright © 2025 Ping Identity Corporation

package pingonetypes

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces
var _ basetypes.StringValuable = ResourceIDValue{}
var _ basetypes.StringValuableWithSemanticEquals = ResourceIDValue{}

type ResourceIDValue struct {
	basetypes.StringValue
}

func (v ResourceIDValue) Equal(o attr.Value) bool {
	other, ok := o.(ResourceIDValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v ResourceIDValue) Type(ctx context.Context) attr.Type {
	// ResourceIDType defined in the schema type section
	return ResourceIDType{}
}

func (v ResourceIDValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	// The framework should always pass the correct value type, but always check
	newValue, ok := newValuable.(ResourceIDValue)

	if !ok {
		diags.AddError(
			"PingOne Resource ID Type Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	// Check whether the flows are equal, ignoring environment metadata and designer UI cues.  Just the flow configuration
	return cmp.Equal(v.StringValue.ValueString(), newValue.ValueString()), diags
}

func NewResourceIDNull() ResourceIDValue {
	return ResourceIDValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewResourceIDUnknown() ResourceIDValue {
	return ResourceIDValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewResourceIDValue(value string) ResourceIDValue {
	return ResourceIDValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewResourceIDPointerValue(value *string) ResourceIDValue {
	return ResourceIDValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}

func ResourceIDNull() ResourceIDValue {
	return NewResourceIDNull()
}

func ResourceIDUnknown() ResourceIDValue {
	return NewResourceIDUnknown()
}
