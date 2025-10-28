// Copyright © 2025 Ping Identity Corporation

package jsontypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisfies the expected interfaces
var _ basetypes.StringTypable = NormalizedObfuscatableType{}

// NormalizedObfuscatableType is a custom type for handling JSON strings that might contain
// obfuscated fields (represented as a string with non-zero length made up of only asterisks).
// It allows for semantic equality checking that ignores these obfuscated fields.
type NormalizedObfuscatableType struct {
	basetypes.StringType
}

func (t NormalizedObfuscatableType) Equal(o attr.Type) bool {
	other, ok := o.(NormalizedObfuscatableType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t NormalizedObfuscatableType) String() string {
	return "jsontypes.NormalizedObfuscatableType"
}

func (t NormalizedObfuscatableType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	// NormalizedObfuscatableValue defined in the value type file
	value := NormalizedObfuscatableValue{
		StringValue: in,
	}

	return value, nil
}

func (t NormalizedObfuscatableType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t NormalizedObfuscatableType) ValueType(ctx context.Context) attr.Value {
	// NormalizedObfuscatableValue defined in the value type file
	return NormalizedObfuscatableValue{}
}
