// Copyright © 2025 Ping Identity Corporation

package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

type ImageResourceModel struct {
	Id   pingonetypes.ResourceIDValue `tfsdk:"id"`
	Href types.String                 `tfsdk:"href"`
}

var (
	ImageTFObjectTypes = map[string]attr.Type{
		"id":   pingonetypes.ResourceIDType{},
		"href": types.StringType,
	}
)

func ImageOkToTF(logo interface{}, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || logo == nil {
		return types.ObjectNull(ImageTFObjectTypes), diags
	}

	b, e := json.Marshal(logo)
	if e != nil {
		diags.AddError(
			"Invalid data object",
			fmt.Sprintf("Cannot remap the data object to JSON: %s.  Please report this to the provider maintainers.", e),
		)
		return types.ObjectNull(ImageTFObjectTypes), diags
	}

	var s map[string]string
	e = json.Unmarshal(b, &s)
	if e != nil {
		diags.AddError(
			"Invalid data object",
			fmt.Sprintf("Cannot remap the data object to map: %s.  Please report this to the provider maintainers.", e),
		)
		return types.ObjectNull(ImageTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{}

	if s["href"] != "" {
		attributesMap["href"] = framework.StringToTF(s["href"])
	} else {
		attributesMap["href"] = types.StringNull()
	}

	if s["id"] != "" {
		attributesMap["id"] = framework.PingOneResourceIDToTF(s["id"])
	} else {
		attributesMap["id"] = pingonetypes.NewResourceIDNull()
	}

	returnVar, d := types.ObjectValue(ImageTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func ImageListToObjectSchemaUpgrade(ctx context.Context, planAttribute types.List) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := ImageTFObjectTypes

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []ImageResourceModel
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := priorStateData[0]

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}
