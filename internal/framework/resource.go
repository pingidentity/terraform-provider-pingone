package framework

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	pingone "github.com/pingidentity/terraform-provider-pingone/internal/client"
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

func Int32ToTF(i int32) basetypes.Int64Value {
	return types.Int64Value(int64(i))
}

func StringOkToTF(v *string, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(*v)
	}
}

func Int32OkToTF(i *int32, ok bool) basetypes.Int64Value {
	if !ok || i == nil {
		return types.Int64Null()
	} else {
		return types.Int64Value(int64(*i))
	}
}

func BoolOkToTF(b *bool, ok bool) basetypes.BoolValue {
	if !ok || b == nil {
		return types.BoolNull()
	} else {
		return types.BoolValue(*b)
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
