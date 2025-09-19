// Copyright © 2025 Ping Identity Corporation

// Package sso provides utility functions for managing application secrets in PingOne SSO service configurations.
// This file contains functions for handling application secret data transformations and retry logic.
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

// ApplicationSecretPreviousTFObjectTypes defines the Terraform Framework attribute types for previous application secret objects.
// This variable maps attribute names to their corresponding Terraform types for previous application secret schema definition.
var (
	ApplicationSecretPreviousTFObjectTypes = map[string]attr.Type{
		"secret":     types.StringType,
		"expires_at": timetypes.RFC3339Type{},
		"last_used":  timetypes.RFC3339Type{},
	}
)

// applicationSecretPreviousOkToTF converts a PingOne API application secret previous object to a Terraform Framework object value.
// It returns a types.Object containing the previous secret information and any diagnostics encountered during conversion.
// The apiObject parameter contains the previous application secret data from the PingOne API response.
// The ok parameter indicates whether the previous secret data was successfully retrieved from the API.
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

// applicationOIDCSecretDataSourceRetryConditions determines whether a failed application OIDC secret API call should be retried.
// It returns true if the error indicates authorization issues that may be resolved by retrying the request.
// The ctx parameter provides context for logging retry decisions.
// The r parameter contains the HTTP response from the failed API call.
// The p1error parameter contains the parsed PingOne API error details.
func applicationOIDCSecretDataSourceRetryConditions(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

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
