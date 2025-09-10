// Copyright © 2025 Ping Identity Corporation

// Package framework provides utilities for Terraform Plugin Framework implementation in the PingOne provider.
// This package contains data type conversion functions, resource ID utilities, import parsing, and other
// common functionality that bridges the gap between PingOne SDK responses and Terraform state management.
package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	pingone "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/davincitypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// ResourceType represents the configuration structure for a framework-based resource.
// It contains the client connection information required for PingOne API interactions
// and is used to pass provider-level configuration to individual resources.
type ResourceType struct {
	// Client is the configured PingOne API client instance used for making requests to the PingOne platform
	Client *pingone.Client
}

// PingOneResourceIDToTF converts a PingOne resource ID string to a Terraform resource ID value.
// It returns a null resource ID value if the input string is empty, otherwise returns a new resource ID value.
// The v parameter must be a string representing the PingOne resource identifier.
func PingOneResourceIDToTF(v string) pingonetypes.ResourceIDValue {
	if v == "" {
		return pingonetypes.NewResourceIDNull()
	} else {
		return pingonetypes.NewResourceIDValue(v)
	}
}

// PingOneResourceIDOkToTF converts a PingOne resource ID pointer and ok boolean to a Terraform resource ID value.
// It returns a null resource ID value if ok is false or the pointer is nil, otherwise returns a new resource ID value.
// The v parameter must be a pointer to a string representing the PingOne resource identifier.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func PingOneResourceIDOkToTF(v *string, ok bool) pingonetypes.ResourceIDValue {
	if !ok || v == nil {
		return pingonetypes.NewResourceIDNull()
	} else {
		return pingonetypes.NewResourceIDValue(*v)
	}
}

// PingOneResourceIDSetOkToTF converts a slice of PingOne resource ID strings to a Terraform set value.
// It returns a null set if ok is false or the slice is nil, otherwise returns a set containing resource ID values.
// The v parameter must be a slice of strings representing PingOne resource identifiers.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func PingOneResourceIDSetOkToTF(v []string, ok bool) basetypes.SetValue {
	if !ok || v == nil {
		return types.SetNull(pingonetypes.ResourceIDType{})
	} else {
		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, PingOneResourceIDToTF(item))
		}

		return types.SetValueMust(pingonetypes.ResourceIDType{}, list)
	}
}

// PingOneResourceIDSetToTF converts a slice of PingOne resource ID strings to a Terraform set value.
// It returns a null set if the slice is nil, otherwise returns a set containing resource ID values.
// The v parameter must be a slice of strings representing PingOne resource identifiers.
func PingOneResourceIDSetToTF(v []string) basetypes.SetValue {
	if v == nil {
		return types.SetNull(pingonetypes.ResourceIDType{})
	} else {
		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, PingOneResourceIDToTF(item))
		}

		return types.SetValueMust(pingonetypes.ResourceIDType{}, list)
	}
}

// PingOneResourceIDListOkToTF converts a slice of PingOne resource ID strings to a Terraform list value.
// It returns a null list if ok is false or the slice is nil, otherwise returns a list containing resource ID values.
// The v parameter must be a slice of strings representing PingOne resource identifiers.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func PingOneResourceIDListOkToTF(v []string, ok bool) basetypes.ListValue {
	if !ok || v == nil {
		return types.ListNull(pingonetypes.ResourceIDType{})
	} else {
		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, PingOneResourceIDToTF(item))
		}

		return types.ListValueMust(pingonetypes.ResourceIDType{}, list)
	}
}

// PingOneResourceIDListToTF converts a slice of PingOne resource ID strings to a Terraform list value.
// It returns a null list if the slice is nil, otherwise returns a list containing resource ID values.
// The v parameter must be a slice of strings representing PingOne resource identifiers.
func PingOneResourceIDListToTF(v []string) basetypes.ListValue {
	if v == nil {
		return types.ListNull(pingonetypes.ResourceIDType{})
	} else {
		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, PingOneResourceIDToTF(item))
		}

		return types.ListValueMust(pingonetypes.ResourceIDType{}, list)
	}
}

// TFTypePingOneResourceIDSliceToStringSlice converts a slice of PingOne resource ID values to a slice of strings.
// It returns a slice of string values and any diagnostics encountered during conversion.
// The v parameter must be a slice of PingOne resource ID values from the Terraform state.
// The path parameter specifies the attribute path for error reporting in case of conversion failures.
// This function validates that all values are neither unknown nor null before conversion.
func TFTypePingOneResourceIDSliceToStringSlice(v []pingonetypes.ResourceIDValue, path path.Path) ([]string, diag.Diagnostics) {
	var sliceOut []string
	var diags diag.Diagnostics

	for _, vElement := range v {
		if vElement.IsUnknown() {
			diags.AddAttributeError(
				path,
				"Unexpected unknown resource ID slice value",
				"Cannot convert a resource ID slice value to string as the slice value is unknown.  Please report this to the provider maintainers.",
			)
			continue
		}
		if vElement.IsNull() {
			diags.AddAttributeError(
				path,
				"Unexpected null slice value",
				"Cannot convert a resource ID slice value to string as the slice value is null.  Please report this to the provider maintainers.",
			)
			continue
		}
		sliceOut = append(sliceOut, vElement.ValueString())
	}

	return sliceOut, diags
}

// DaVinciResourceIDToTF converts a DaVinci resource ID string to a Terraform resource ID value.
// It returns a null resource ID value if the input string is empty, otherwise returns a new resource ID value.
// The v parameter must be a string representing the DaVinci resource identifier.
func DaVinciResourceIDToTF(v string) davincitypes.ResourceIDValue {
	if v == "" {
		return davincitypes.NewResourceIDNull()
	} else {
		return davincitypes.NewResourceIDValue(v)
	}
}

// DaVinciResourceIDOkToTF converts a DaVinci resource ID pointer and ok boolean to a Terraform resource ID value.
// It returns a null resource ID value if ok is false or the pointer is nil, otherwise returns a new resource ID value.
// The v parameter must be a pointer to a string representing the DaVinci resource identifier.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func DaVinciResourceIDOkToTF(v *string, ok bool) davincitypes.ResourceIDValue {
	if !ok || v == nil {
		return davincitypes.NewResourceIDNull()
	} else {
		return davincitypes.NewResourceIDValue(*v)
	}
}

// DaVinciResourceIDSetOkToTF converts a slice of DaVinci resource ID strings to a Terraform set value.
// It returns a null set if ok is false or the slice is nil, otherwise returns a set containing resource ID values.
// The v parameter must be a slice of strings representing DaVinci resource identifiers.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func DaVinciResourceIDSetOkToTF(v []string, ok bool) basetypes.SetValue {
	if !ok || v == nil {
		return types.SetNull(davincitypes.ResourceIDType{})
	} else {
		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, DaVinciResourceIDToTF(item))
		}

		return types.SetValueMust(davincitypes.ResourceIDType{}, list)
	}
}

// DaVinciResourceIDSetToTF converts a slice of DaVinci resource ID strings to a Terraform set value.
// It returns a null set if the slice is nil, otherwise returns a set containing resource ID values.
// The v parameter must be a slice of strings representing DaVinci resource identifiers.
func DaVinciResourceIDSetToTF(v []string) basetypes.SetValue {
	if v == nil {
		return types.SetNull(davincitypes.ResourceIDType{})
	} else {
		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, DaVinciResourceIDToTF(item))
		}

		return types.SetValueMust(davincitypes.ResourceIDType{}, list)
	}
}

// DaVinciResourceIDListOkToTF converts a slice of DaVinci resource ID strings to a Terraform list value.
// It returns a null list if ok is false or the slice is nil, otherwise returns a list containing resource ID values.
// The v parameter must be a slice of strings representing DaVinci resource identifiers.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func DaVinciResourceIDListOkToTF(v []string, ok bool) basetypes.ListValue {
	if !ok || v == nil {
		return types.ListNull(pingonetypes.ResourceIDType{})
	} else {
		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, PingOneResourceIDToTF(item))
		}

		return types.ListValueMust(pingonetypes.ResourceIDType{}, list)
	}
}

// DaVinciResourceIDListToTF converts a slice of DaVinci resource ID strings to a Terraform list value.
// It returns a null list if the slice is nil, otherwise returns a list containing resource ID values.
// The v parameter must be a slice of strings representing DaVinci resource identifiers.
func DaVinciResourceIDListToTF(v []string) basetypes.ListValue {
	if v == nil {
		return types.ListNull(pingonetypes.ResourceIDType{})
	} else {
		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, PingOneResourceIDToTF(item))
		}

		return types.ListValueMust(pingonetypes.ResourceIDType{}, list)
	}
}

// JSONNormalizedToTF converts a map of interface values to a Terraform normalized JSON type.
// It returns a normalized JSON value and any diagnostics encountered during marshaling.
// The v parameter must be a map containing the JSON data to be normalized and stored in Terraform state.
// This function marshals the map to JSON and creates a normalized JSON type that preserves formatting.
func JSONNormalizedToTF(v map[string]interface{}) (jsontypes.Normalized, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v == nil {
		return jsontypes.NewNormalizedNull(), diags
	} else {
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			diags.Append(diag.NewErrorDiagnostic("Normalized JSON Type Conversion Error", err.Error()))
			return jsontypes.NewNormalizedNull(), diags
		}
		return jsontypes.NewNormalizedValue(string(jsonBytes[:])), diags
	}
}

// JSONNormalizedOkToTF converts a map of interface values to a Terraform normalized JSON type with ok boolean check.
// It returns a null normalized JSON value if ok is false or the map is nil, otherwise returns the converted value.
// The v parameter must be a map containing the JSON data to be normalized and stored in Terraform state.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func JSONNormalizedOkToTF(v map[string]interface{}, ok bool) (jsontypes.Normalized, diag.Diagnostics) {
	if !ok || v == nil {
		return jsontypes.NewNormalizedNull(), nil
	} else {
		return JSONNormalizedToTF(v)
	}
}

// StringToTF converts a string value to a Terraform string type.
// It returns a null string value if the input string is empty, otherwise returns a new string value.
// The v parameter must be a string value to be stored in Terraform state.
// StringToTF converts a string to a Terraform string type.
// It returns a null string value if the input string is empty, otherwise returns a new string value.
// The v parameter must be a string value from the API response.
func StringToTF(v string) basetypes.StringValue {
	if v == "" {
		return types.StringNull()
	} else {
		return types.StringValue(v)
	}
}

// StringOkToTF converts a string pointer and ok boolean to a Terraform string type.
// It returns a null string value if ok is false or the pointer is nil, otherwise returns a new string value.
// The v parameter must be a pointer to a string value from the API response.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func StringOkToTF(v *string, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(*v)
	}
}

// Int32ToTF converts an int32 value to a Terraform int32 type.
// It returns a new int32 value for storage in Terraform state.
// The i parameter must be an int32 value from the API response.
func Int32ToTF(i int32) basetypes.Int32Value {
	return types.Int32Value(i)
}

// Int64ToTF converts an int64 value to a Terraform int64 type.
// It returns a new int64 value for storage in Terraform state.
// The i parameter must be an int64 value from the API response.
func Int64ToTF(i int64) basetypes.Int64Value {
	return types.Int64Value(i)
}

// EnumToTF converts an enum interface value to a Terraform string type.
// It returns a null string value if the enum is nil, otherwise returns the string representation of the enum.
// The v parameter must be an enum value from the PingOne SDK that implements string conversion.
func EnumToTF(v interface{}) basetypes.StringValue {
	if v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(utils.EnumToString(v))
	}
}

// EnumOkToTF converts an enum interface value and ok boolean to a Terraform string type.
// It returns a null string value if ok is false or the enum is nil, otherwise returns the string representation.
// The v parameter must be an enum value from the PingOne SDK that implements string conversion.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func EnumOkToTF(v interface{}, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(utils.EnumToString(v))
	}
}

// Int32OkToTF converts an int32 pointer and ok boolean to a Terraform int32 type.
// It returns a null int32 value if ok is false or the pointer is nil, otherwise returns a new int32 value.
// The i parameter must be a pointer to an int32 value from the API response.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func Int32OkToTF(i *int32, ok bool) basetypes.Int32Value {
	if !ok || i == nil {
		return types.Int32Null()
	} else {
		return types.Int32Value(*i)
	}
}

// Int64OkToTF converts an int64 pointer and ok boolean to a Terraform int64 type.
// It returns a null int64 value if ok is false or the pointer is nil, otherwise returns a new int64 value.
// The i parameter must be a pointer to an int64 value from the API response.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func Int64OkToTF(i *int64, ok bool) basetypes.Int64Value {
	if !ok || i == nil {
		return types.Int64Null()
	} else {
		return types.Int64Value(*i)
	}
}

// Float32OkToTF converts a float32 pointer and ok boolean to a Terraform float32 type.
// It returns a null float32 value if ok is false or the pointer is nil, otherwise returns a new float32 value.
// The i parameter must be a pointer to a float32 value from the API response.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func Float32OkToTF(i *float32, ok bool) basetypes.Float32Value {
	if !ok || i == nil {
		return types.Float32Null()
	} else {
		return types.Float32Value(*i)
	}
}

// Float64OkToTF converts a float64 pointer and ok boolean to a Terraform float64 type.
// It returns a null float64 value if ok is false or the pointer is nil, otherwise returns a new float64 value.
// The i parameter must be a pointer to a float64 value from the API response.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func Float64OkToTF(i *float64, ok bool) basetypes.Float64Value {
	if !ok || i == nil {
		return types.Float64Null()
	} else {
		return types.Float64Value(*i)
	}
}

// BoolOkToTF converts a boolean pointer and ok boolean to a Terraform boolean type.
// It returns a null boolean value if ok is false or the pointer is nil, otherwise returns a new boolean value.
// The b parameter must be a pointer to a boolean value from the API response.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func BoolOkToTF(b *bool, ok bool) basetypes.BoolValue {
	if !ok || b == nil {
		return types.BoolNull()
	} else {
		return types.BoolValue(*b)
	}
}

// TimeOkToTF converts a time.Time pointer and ok boolean to a Terraform RFC3339 time type.
// It returns a null RFC3339 time value if ok is false or the pointer is nil, otherwise returns a new time value.
// The v parameter must be a pointer to a time.Time value from the API response.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func TimeOkToTF(v *time.Time, ok bool) timetypes.RFC3339 {
	if !ok || v == nil {
		return timetypes.NewRFC3339Null()
	} else {
		return timetypes.NewRFC3339TimeValue(*v)
	}
}

// StringSetOkToTF converts a slice of strings and ok boolean to a Terraform set type.
// It returns a null set if ok is false or the slice is nil, otherwise returns a set containing string values.
// The v parameter must be a slice of strings from the API response.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func StringSetOkToTF(v []string, ok bool) basetypes.SetValue {
	if !ok || v == nil {
		return types.SetNull(types.StringType)
	} else {
		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, StringToTF(item))
		}

		return types.SetValueMust(types.StringType, list)
	}
}

// StringSetToTF converts a slice of strings to a Terraform set type.
// It returns a null set if the slice is nil, otherwise returns a set containing string values.
// The v parameter must be a slice of strings from the API response.
func StringSetToTF(v []string) basetypes.SetValue {
	if v == nil {
		return types.SetNull(types.StringType)
	} else {
		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, StringToTF(item))
		}

		return types.SetValueMust(types.StringType, list)
	}
}

// StringListOkToTF converts a slice of strings and ok boolean to a Terraform list type.
// It returns a null list if ok is false or the slice is nil, otherwise returns a list containing string values.
// The v parameter must be a slice of strings from the API response.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func StringListOkToTF(v []string, ok bool) basetypes.ListValue {
	if !ok || v == nil {
		return types.ListNull(types.StringType)
	} else {
		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, StringToTF(item))
		}

		return types.ListValueMust(types.StringType, list)
	}
}

// StringListToTF converts a slice of strings to a Terraform list type.
// It returns a null list if the slice is nil, otherwise returns a list containing string values.
// The v parameter must be a slice of strings from the API response.
func StringListToTF(v []string) basetypes.ListValue {
	if v == nil {
		return types.ListNull(types.StringType)
	} else {
		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, StringToTF(item))
		}

		return types.ListValueMust(types.StringType, list)
	}
}

// StringMapOkToTF converts a map of strings pointer and ok boolean to a Terraform map type.
// It returns a null map if ok is false or the pointer is nil, otherwise returns a map containing string values.
// The v parameter must be a pointer to a map of strings from the API response.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func StringMapOkToTF(v *map[string]string, ok bool) basetypes.MapValue {
	if !ok || v == nil {
		return types.MapNull(types.StringType)
	} else {
		list := make(map[string]attr.Value, 0)
		for key, item := range *v {
			list[key] = StringToTF(item)
		}

		return types.MapValueMust(types.StringType, list)
	}
}

// EnumSetOkToTF converts an enum slice interface and ok boolean to a Terraform set type.
// It returns a null set if ok is false or the enum is nil, otherwise returns a set containing string values.
// The v parameter must be an enum slice from the PingOne SDK that implements string slice conversion.
// The ok parameter indicates whether the value was successfully retrieved from the source API response.
func EnumSetOkToTF(v interface{}, ok bool) basetypes.SetValue {
	if !ok || v == nil {
		return types.SetNull(types.StringType)
	} else {
		list := make([]attr.Value, 0)
		for _, item := range utils.EnumSliceToStringSlice(v) {
			list = append(list, StringToTF(item))
		}

		return types.SetValueMust(types.StringType, list)
	}
}

// StringSliceToTF converts a slice of strings to a Terraform list type with error handling.
// It returns a list value and any diagnostics encountered during conversion.
// The v parameter must be a slice of strings to be converted to a Terraform list.
// This function differs from StringListToTF by providing diagnostic information for conversion errors.
func StringSliceToTF(v []string) (basetypes.ListValue, diag.Diagnostics) {
	if v == nil {
		return types.ListNull(types.StringType), nil
	} else {

		list := make([]attr.Value, 0)
		for _, item := range v {
			list = append(list, StringToTF(item))
		}

		return types.ListValue(types.StringType, list)
	}
}

// StringSliceToTFSet converts a slice of strings to a Terraform set type with error handling.
// It returns a set value and any diagnostics encountered during conversion.
// The v parameter must be a slice of strings to be converted to a Terraform set.
// This function differs from StringSetToTF by providing diagnostic information for conversion errors.
func StringSliceToTFSet(v []string) (basetypes.SetValue, diag.Diagnostics) {
	if v == nil {
		return types.SetNull(types.StringType), nil
	} else {

		set := make([]attr.Value, 0)
		for _, item := range v {
			set = append(set, StringToTF(item))
		}

		return types.SetValue(types.StringType, set)
	}
}

// TFListToStringSlice converts a Terraform list to a slice of string pointers.
// It returns a slice of string pointers extracted from the Terraform list.
// The ctx parameter provides the context for the conversion operation.
// The v parameter must be a Terraform list containing string values.
// This function returns nil if the list is null, unknown, or if conversion errors occur.
func TFListToStringSlice(ctx context.Context, v types.List) []*string {
	var sliceOut []*string

	if v.IsNull() || v.IsUnknown() {
		return nil
	} else {

		if v.ElementsAs(ctx, &sliceOut, false).HasError() {
			return nil
		}
	}

	return sliceOut
}

// TFSetToStringSlice converts a Terraform set to a slice of string pointers.
// It returns a slice of string pointers extracted from the Terraform set.
// The ctx parameter provides the context for the conversion operation.
// The v parameter must be a Terraform set containing string values.
// This function returns nil if the set is null, unknown, or if conversion errors occur.
func TFSetToStringSlice(ctx context.Context, v types.Set) []*string {
	var sliceOut []*string

	if v.IsNull() || v.IsUnknown() {
		return nil
	} else {

		if v.ElementsAs(ctx, &sliceOut, false).HasError() {
			return nil
		}
	}

	return sliceOut
}

// TFTypeStringSliceToStringSlice converts a slice of Terraform string types to a slice of strings.
// It returns a slice of string values and any diagnostics encountered during conversion.
// The v parameter must be a slice of Terraform string types from the state.
// The path parameter specifies the attribute path for error reporting in case of conversion failures.
// This function validates that all values are neither unknown nor null before conversion.
func TFTypeStringSliceToStringSlice(v []types.String, path path.Path) ([]string, diag.Diagnostics) {
	var sliceOut []string
	var diags diag.Diagnostics

	for _, vElement := range v {
		if vElement.IsUnknown() {
			diags.AddAttributeError(
				path,
				"Unexpected unknown slice value",
				"Cannot convert an slice value to string as the slice value is unknown.  Please report this to the provider maintainers.",
			)
			continue
		}
		if vElement.IsNull() {
			diags.AddAttributeError(
				path,
				"Unexpected null slice value",
				"Cannot convert an slice value to string as the slice value is null.  Please report this to the provider maintainers.",
			)
			continue
		}
		sliceOut = append(sliceOut, vElement.ValueString())
	}

	return sliceOut, diags
}

// ImportComponent represents a single component of a resource import ID format.
// It defines the label, validation regex, and whether this component represents the primary resource ID.
type ImportComponent struct {
	// Label is the human-readable name for this component, used in error messages and documentation
	Label string
	// Regexp is the regular expression that validates the format of this component
	Regexp *regexp.Regexp
	// PrimaryID indicates whether this component represents the primary resource identifier
	PrimaryID bool
}

// ParseImportID parses a resource import ID string according to the specified component format.
// It returns a map of component labels to their extracted values, or an error if parsing fails.
// The id parameter must be the import ID string provided by the user during resource import.
// The components parameter defines the expected format and validation rules for each ID component.
// Components are expected to be separated by forward slashes in the import ID string.
func ParseImportID(id string, components ...ImportComponent) (map[string]string, error) {

	keys := make([]string, len(components))
	regexpList := make([]string, len(components))

	i := 0
	for _, v := range components {
		keys[i] = v.Label
		regexpList[i] = v.Regexp.String()
		i++
	}

	compiledRegexpString := fmt.Sprintf("^%s$", strings.Join(regexpList, `\/`))

	m, err := regexp.MatchString(compiledRegexpString, id)
	if err != nil {
		return nil, fmt.Errorf("Cannot verify import ID regex: %s", err)
	}

	if !m {
		return nil, fmt.Errorf("Invalid import ID specified (\"%s\").  The ID should be in the format \"%s\" and must match regex: %s", id, strings.Join(keys, "/"), compiledRegexpString)
	}

	attributeValues := strings.SplitN(id, "/", len(components))

	if len(attributeValues) != len(components) {
		return nil, fmt.Errorf("Invalid import ID specified (\"%s\").  The ID should be in the format \"%s\".", id, strings.Join(keys, "/"))
	}

	attributes := make(map[string]string)

	i = 0
	for _, v := range components {
		attributes[v.Label] = attributeValues[i]
		i++
	}

	return attributes, nil
}
