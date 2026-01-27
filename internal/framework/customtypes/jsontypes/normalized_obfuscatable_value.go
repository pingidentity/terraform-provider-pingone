// Copyright © 2025 Ping Identity Corporation

package jsontypes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces
var _ basetypes.StringValuable = NormalizedObfuscatableValue{}
var _ basetypes.StringValuableWithSemanticEquals = NormalizedObfuscatableValue{}
var _ xattr.ValidateableAttribute = NormalizedObfuscatableValue{}

// Alias useful in schemas and structs
type NormalizedObfuscatable = NormalizedObfuscatableValue

// NormalizedObfuscatableValue is a custom value for handling JSON strings that might contain
// obfuscated fields (represented as a string with non-zero length made up of only asterisks).
// It allows for semantic equality checking that ignores these obfuscated fields.
type NormalizedObfuscatableValue struct {
	basetypes.StringValue
}

func (v NormalizedObfuscatableValue) Equal(o attr.Value) bool {
	other, ok := o.(NormalizedObfuscatableValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v NormalizedObfuscatableValue) Type(ctx context.Context) attr.Type {
	return NormalizedObfuscatableType{}
}

func (v NormalizedObfuscatableValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	// The framework should always pass the correct value type, but always check
	newValue, ok := newValuable.(NormalizedObfuscatableValue)

	if !ok {
		diags.AddError(
			"JSON String Type Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	return semanticCompareJsonIgnoreObfuscated(v.ValueString(), newValue.ValueString()), diags
}

// Semantically compare the json blobs, ignoring fields that are strings of asterisks, which may have been obfuscated
func semanticCompareJsonIgnoreObfuscated(json1, json2 string) bool {
	var marshalled1, marshalled2 map[string]interface{}

	// Unmarshal both JSON blobs into generic maps. Return false if parsing fails.
	if err := json.Unmarshal([]byte(json1), &marshalled1); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(json2), &marshalled2); err != nil {
		return false
	}

	ignoreRedactedFields := cmp.FilterPath(func(p cmp.Path) bool {
		// Get the values being compared at the current path
		vx, vy := p.Last().Values()

		if !vx.IsValid() || !vy.IsValid() {
			return false
		}

		vxStr, okX := vx.Interface().(string)
		vyStr, okY := vy.Interface().(string)

		// If the values are strings, and one or both are strings of asterisks, we ignore this path
		if okX && okY {
			return isAllAsterisks(vxStr) || isAllAsterisks(vyStr)
		}
		return false
	}, cmp.Ignore())

	return cmp.Equal(marshalled1, marshalled2, ignoreRedactedFields)
}

func isAllAsterisks(s string) bool {
	// Empty string is not considered to be all asterisks
	if len(s) == 0 {
		return false
	}
	for _, char := range s {
		if char != '*' {
			return false
		}
	}
	return true
}

func NormalizedObfuscatableNull() NormalizedObfuscatableValue {
	return NormalizedObfuscatableValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NormalizedObfuscatableUnknown() NormalizedObfuscatableValue {
	return NormalizedObfuscatableValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NormalizedObfuscatableStringValue(value string) NormalizedObfuscatableValue {
	return NormalizedObfuscatableValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NormalizedObfuscatableStringPointerValue(value *string) NormalizedObfuscatableValue {
	return NormalizedObfuscatableValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}

func (v NormalizedObfuscatableValue) ValidateAttribute(ctx context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsNull() || v.IsUnknown() {
		return
	}

	// Validate that the string is valid JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(v.ValueString()), &jsonData); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"JSON String Type Validation Error",
			fmt.Sprintf("The value could not be unmarshalled as JSON: %s", err.Error()),
		)
	}
}
