package sso

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

var (
	ApplicationSecretPreviousTFObjectTypes = map[string]attr.Type{
		"secret":     types.StringType,
		"expires_at": types.StringType,
		"last_used":  types.StringType,
	}
)

func applicationSecretPreviousOkToTF(apiObject *management.ApplicationSecretPrevious, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(ApplicationSecretPreviousTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"secret":     framework.StringOkToTF(apiObject.GetSecretOk()),
		"expires_at": framework.TimeOkToTF(apiObject.GetExpiresAtOk()),
		"last_used":  framework.TimeOkToTF(apiObject.GetLastUsedOk()),
	}

	returnVar, d := types.ObjectValue(ApplicationSecretPreviousTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}
