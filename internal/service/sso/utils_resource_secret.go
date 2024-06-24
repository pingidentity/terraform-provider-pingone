package sso

import (
	"context"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

var (
	ResourceSecretPreviousTFObjectTypes = map[string]attr.Type{
		"secret":     types.StringType,
		"expires_at": timetypes.RFC3339Type{},
		"last_used":  timetypes.RFC3339Type{},
	}
)

func resourceSecretPreviousOkToTF(apiObject *management.ResourceSecretPrevious, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(ResourceSecretPreviousTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"secret":     framework.StringOkToTF(apiObject.GetSecretOk()),
		"expires_at": framework.TimeOkToTF(apiObject.GetExpiresAtOk()),
		"last_used":  framework.TimeOkToTF(apiObject.GetLastUsedOk()),
	}

	returnVar, d := types.ObjectValue(ResourceSecretPreviousTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func resourceOIDCSecretDataSourceRetryConditions(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

	if p1error != nil {

		m, err := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage())
		if err == nil && m {
			tflog.Warn(ctx, "Insufficient PingOne privileges detected")
			return true
		} else if err != nil {
			tflog.Warn(ctx, "Cannot match error string for retry")
			return false
		}

	}

	return false
}
