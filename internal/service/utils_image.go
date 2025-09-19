// Copyright © 2025 Ping Identity Corporation

// Package service provides utility functions and common types for handling image resources
// across different PingOne services in the Terraform provider.
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

// ImageResourceModel represents the structure for image resource data in Terraform configurations.
// It contains the resource ID and HREF URL for images associated with PingOne resources.
type ImageResourceModel struct {
	// Id is the unique identifier for the image resource in PingOne
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
	// Href is the URL where the image can be accessed
	Href types.String `tfsdk:"href"`
}

// ImageTFObjectTypes defines the Terraform Framework attribute types for image objects.
// This variable maps attribute names to their corresponding Terraform types for proper schema definition.
var (
	ImageTFObjectTypes = map[string]attr.Type{
		"id":   pingonetypes.ResourceIDType{},
		"href": types.StringType,
	}
)

// ImageOkToTF converts a PingOne API image object to a Terraform Framework object value.
// It returns a types.Object containing the image ID and HREF, and any diagnostics encountered during conversion.
// The logo parameter should be an interface{} containing image data from the PingOne API response.
// The ok parameter indicates whether the image data was successfully retrieved from the API.
// This function handles JSON marshaling/unmarshaling to extract the ID and HREF from the image data.
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

// ImageListToObjectSchemaUpgrade performs schema upgrade for image attributes from list to object format.
// It returns a types.Object representing the upgraded image data and any diagnostics encountered during the upgrade.
// The ctx parameter provides the context for the upgrade operation.
// The planAttribute parameter contains the prior state data in list format that needs to be upgraded to object format.
// This function is typically used during provider version upgrades when schema changes from list to object structure.
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
