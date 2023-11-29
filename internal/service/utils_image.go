package service

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

type ImageResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Href types.String `tfsdk:"href"`
}

var (
	ImageTFObjectTypes = map[string]attr.Type{
		"id":   types.StringType,
		"href": types.StringType,
	}
)

func ImageOkToTF(logo interface{}, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: ImageTFObjectTypes}

	if !ok || logo == nil {
		return types.ListNull(tfObjType), diags
	}

	b, e := json.Marshal(logo)
	if e != nil {
		diags.AddError(
			"Invalid data object",
			fmt.Sprintf("Cannot remap the data object to JSON: %s.  Please report this to the provider maintainers.", e),
		)
		return types.ListNull(tfObjType), diags
	}

	var s map[string]string
	e = json.Unmarshal(b, &s)
	if e != nil {
		diags.AddError(
			"Invalid data object",
			fmt.Sprintf("Cannot remap the data object to map: %s.  Please report this to the provider maintainers.", e),
		)
		return types.ListNull(tfObjType), diags
	}

	attributesMap := map[string]attr.Value{}

	if s["href"] != "" {
		attributesMap["href"] = framework.StringToTF(s["href"])
	} else {
		attributesMap["href"] = types.StringNull()
	}

	if s["id"] != "" {
		attributesMap["id"] = framework.StringToTF(s["id"])
	} else {
		attributesMap["id"] = types.StringNull()
	}

	flattenedObj, d := types.ObjectValue(ImageTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}
