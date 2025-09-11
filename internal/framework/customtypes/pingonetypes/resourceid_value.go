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

// ResourceIDValue represents a PingOne resource identifier value with semantic comparison capabilities.
// It extends the base string value with PingOne-specific semantics and comparison logic
// for resource identifiers used throughout the provider.
type ResourceIDValue struct {
	basetypes.StringValue
}

// Equal determines if this ResourceIDValue is equal to another attr.Value.
// It returns true if the other value is also a ResourceIDValue with an equivalent base StringValue.
func (v ResourceIDValue) Equal(o attr.Value) bool {
	other, ok := o.(ResourceIDValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// Type returns the ResourceIDType that this value implements.
// It returns the corresponding type definition for this value in the type system.
// The ctx parameter provides the context for the type operation.
func (v ResourceIDValue) Type(ctx context.Context) attr.Type {
	// ResourceIDType defined in the schema type section
	return ResourceIDType{}
}

// StringSemanticEquals performs semantic equality comparison between ResourceIDValue instances.
// It returns true if the two values are semantically equivalent and any diagnostics encountered.
// This method performs string-based comparison of the resource identifier values.
//
// The ctx parameter provides the context for the comparison operation.
// The newValuable parameter is the other StringValuable to compare against.
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

// NewResourceIDNull creates a new ResourceIDValue with a null state.
// It returns a ResourceIDValue that represents the absence of a resource ID value.
func NewResourceIDNull() ResourceIDValue {
	return ResourceIDValue{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewResourceIDUnknown creates a new ResourceIDValue with an unknown state.
// It returns a ResourceIDValue that represents a resource ID value that is not yet known.
func NewResourceIDUnknown() ResourceIDValue {
	return ResourceIDValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewResourceIDValue creates a new ResourceIDValue with the specified string value.
// It returns a ResourceIDValue containing the provided resource identifier.
// The value parameter must be a string representing the PingOne resource identifier.
func NewResourceIDValue(value string) ResourceIDValue {
	return ResourceIDValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewResourceIDPointerValue creates a new ResourceIDValue from a string pointer.
// It returns a ResourceIDValue containing the value from the pointer, or null if the pointer is nil.
// The value parameter must be a pointer to a string representing the PingOne resource identifier.
func NewResourceIDPointerValue(value *string) ResourceIDValue {
	return ResourceIDValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}

// ResourceIDNull creates a new null ResourceIDValue.
// It returns a ResourceIDValue that represents the absence of a resource ID value.
// This is an alias for NewResourceIDNull() for consistency with framework patterns.
func ResourceIDNull() ResourceIDValue {
	return NewResourceIDNull()
}

// ResourceIDUnknown creates a new unknown ResourceIDValue.
// It returns a ResourceIDValue that represents a resource ID value that is not yet known.
// This is an alias for NewResourceIDUnknown() for consistency with framework patterns.
func ResourceIDUnknown() ResourceIDValue {
	return NewResourceIDUnknown()
}
