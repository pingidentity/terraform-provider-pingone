package framework

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	pingone "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

type ResourceType struct {
	Client *pingone.Client
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

func Int32ToTF(i int32) basetypes.Int64Value {
	return types.Int64Value(int64(i))
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

func Int32OkToTF(i *int32, ok bool) basetypes.Int64Value {
	if !ok || i == nil {
		return types.Int64Null()
	} else {
		return types.Int64Value(int64(*i))
	}
}

func Float32OkToTF(i *float32, ok bool) basetypes.Float64Value {
	if !ok || i == nil {
		return types.Float64Null()
	} else {
		return types.Float64Value(float64(*i))
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
