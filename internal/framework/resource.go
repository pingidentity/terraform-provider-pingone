// Copyright © 2025 Ping Identity Corporation

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
	"github.com/pingidentity/pingone-go-client/pingone"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/davincitypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

type ResourceType struct {
	Client *pingone.APIClient
}

func PingOneResourceIDToTF(v string) pingonetypes.ResourceIDValue {
	if v == "" {
		return pingonetypes.NewResourceIDNull()
	} else {
		return pingonetypes.NewResourceIDValue(v)
	}
}

func PingOneResourceIDOkToTF(v *string, ok bool) pingonetypes.ResourceIDValue {
	if !ok || v == nil {
		return pingonetypes.NewResourceIDNull()
	} else {
		return pingonetypes.NewResourceIDValue(*v)
	}
}

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

func DaVinciResourceIDToTF(v string) davincitypes.ResourceIDValue {
	if v == "" {
		return davincitypes.NewResourceIDNull()
	} else {
		return davincitypes.NewResourceIDValue(v)
	}
}

func DaVinciResourceIDOkToTF(v *string, ok bool) davincitypes.ResourceIDValue {
	if !ok || v == nil {
		return davincitypes.NewResourceIDNull()
	} else {
		return davincitypes.NewResourceIDValue(*v)
	}
}

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

func JSONNormalizedOkToTF(v map[string]interface{}, ok bool) (jsontypes.Normalized, diag.Diagnostics) {
	if !ok || v == nil {
		return jsontypes.NewNormalizedNull(), nil
	} else {
		return JSONNormalizedToTF(v)
	}
}

func StringToTF(v string) basetypes.StringValue {
	if v == "" {
		return types.StringNull()
	} else {
		return types.StringValue(v)
	}
}

func StringOkToTF(v *string, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(*v)
	}
}

func Int32ToTF(i int32) basetypes.Int32Value {
	return types.Int32Value(i)
}

func Int64ToTF(i int64) basetypes.Int64Value {
	return types.Int64Value(i)
}

func EnumToTF(v interface{}) basetypes.StringValue {
	if v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(utils.EnumToString(v))
	}
}

func EnumOkToTF(v interface{}, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(utils.EnumToString(v))
	}
}

func Int32OkToTF(i *int32, ok bool) basetypes.Int32Value {
	if !ok || i == nil {
		return types.Int32Null()
	} else {
		return types.Int32Value(*i)
	}
}

func Int64OkToTF(i *int64, ok bool) basetypes.Int64Value {
	if !ok || i == nil {
		return types.Int64Null()
	} else {
		return types.Int64Value(*i)
	}
}

func Float32OkToTF(i *float32, ok bool) basetypes.Float32Value {
	if !ok || i == nil {
		return types.Float32Null()
	} else {
		return types.Float32Value(*i)
	}
}

func Float64OkToTF(i *float64, ok bool) basetypes.Float64Value {
	if !ok || i == nil {
		return types.Float64Null()
	} else {
		return types.Float64Value(*i)
	}
}

func BoolOkToTF(b *bool, ok bool) basetypes.BoolValue {
	if !ok || b == nil {
		return types.BoolNull()
	} else {
		return types.BoolValue(*b)
	}
}

func TimeOkToTF(v *time.Time, ok bool) timetypes.RFC3339 {
	if !ok || v == nil {
		return timetypes.NewRFC3339Null()
	} else {
		return timetypes.NewRFC3339TimeValue(*v)
	}
}

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

type ImportComponent struct {
	Label     string
	Regexp    *regexp.Regexp
	PrimaryID bool
}

// Parse Import ID format
func ParseImportID(id string, components ...ImportComponent) (map[string]string, error) {

	keys := make([]string, len(components))
	regexpList := make([]string, len(components))

	i := 0
	for _, v := range components {
		if v.Regexp == nil {
			return nil, fmt.Errorf("cannot parse import ID as component %d has no Regexp", i)
		}
		keys[i] = v.Label
		regexpList[i] = v.Regexp.String()
		i++
	}

	compiledRegexpString := fmt.Sprintf("^%s$", strings.Join(regexpList, `\/`))

	m, err := regexp.MatchString(compiledRegexpString, id)
	if err != nil {
		return nil, fmt.Errorf("cannot verify import ID regex: %s", err)
	}

	if !m {
		return nil, fmt.Errorf("invalid import ID specified (\"%s\").  The ID should be in the format \"%s\" and must match regex: %s", id, strings.Join(keys, "/"), compiledRegexpString)
	}

	attributeValues := strings.SplitN(id, "/", len(components))

	if len(attributeValues) != len(components) {
		return nil, fmt.Errorf("invalid import ID specified (\"%s\").  The ID should be in the format \"%s\"", id, strings.Join(keys, "/"))
	}

	attributes := make(map[string]string)

	i = 0
	for _, v := range components {
		attributes[v.Label] = attributeValues[i]
		i++
	}

	return attributes, nil
}
