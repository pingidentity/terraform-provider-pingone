// Copyright © 2025 Ping Identity Corporation

package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pingidentity/pingone-go-client/pingone"
)

type SDKInterfaceFunc func() (any, *http.Response, error)

type CustomError func(*http.Response, *pingone.GeneralError) diag.Diagnostics

var (
	DefaultCustomError = func(_ *http.Response, _ *pingone.GeneralError) diag.Diagnostics { return nil }

	CustomErrorResourceNotFoundWarning = func(r *http.Response, p1Error *pingone.GeneralError) diag.Diagnostics {
		var diags diag.Diagnostics

		// Deleted outside of TF
		if p1Error != nil && p1Error.GetCode() == "NOT_FOUND" {
			diags.AddWarning("Requested resource not found", fmt.Sprintf("The requested resource configuration cannot be found in the PingOne service.  If the requested resource is managed in Terraform's state, it may have been removed outside of Terraform.\nAPI error: %s", p1Error.GetMessage()))

			return diags
		}

		if r != nil && r.StatusCode == 404 {
			diags.AddWarning("Requested resource not found", "The requested resource configuration cannot be found in the PingOne service.  If the requested resource is managed in Terraform's state, it may have been removed outside of Terraform.")

			return diags
		}

		return nil
	}

	regionMappingSuffixMap = map[string]string{
		"na": "com",
		"eu": "eu",
		"ap": "asia",
		"au": "com.au",
		"ca": "ca",
		"sg": "sg",
	}
)

func RegionSuffixFromCode(regionCode string) (string, bool) {
	suffix, ok := regionMappingSuffixMap[regionCode]
	return suffix, ok
}

func UserAgent(suffix, version string) string {
	var agentBuilder strings.Builder
	agentBuilder.WriteString("terraform-provider-pingone/")
	agentBuilder.WriteString(version)
	if suffix != "" {
		agentBuilder.WriteString(" ")
		agentBuilder.WriteString(suffix)
	}
	return agentBuilder.String()
}

func CheckEnvironmentExistsOnPermissionsError(ctx context.Context, apiClient *pingone.APIClient, environmentID string, fO any, fR *http.Response, fErr error) (any, *http.Response, error) {
	if fR != nil && (fR.StatusCode == http.StatusUnauthorized || fR.StatusCode == http.StatusForbidden || fR.StatusCode == http.StatusBadRequest) {
		environmentIdUuid, err := uuid.Parse(environmentID)
		if err != nil {
			return fO, nil, fmt.Errorf("unable to parse environment id '%s' as uuid: %v", environmentID, err)
		}

		//TODO remove placeholder expand once pingone-go-client is updated to remove this requirement
		_, fER, fEErr := apiClient.EnvironmentsApi.GetEnvironmentById(ctx, environmentIdUuid).Expand("placeholder").Execute()

		if fER.StatusCode == http.StatusNotFound {
			tflog.Warn(ctx, "API responded with 400, 401 or 403, and the provider determined the environment doesn't exist.  Overriding resource response.")
			return fO, fER, fEErr
		}
	}

	return fO, fR, fErr
}

func ParseResponse(ctx context.Context, f SDKInterfaceFunc, requestID string, customError CustomError, targetObject any) diag.Diagnostics {
	defaultTimeout := 10
	return ParseResponseWithCustomTimeout(ctx, f, requestID, customError, targetObject, time.Duration(defaultTimeout)*time.Minute)
}

func ParseResponseWithCustomTimeout(ctx context.Context, f SDKInterfaceFunc, requestID string, customError CustomError, targetObject any, timeout time.Duration) diag.Diagnostics {
	var diags diag.Diagnostics

	if customError == nil {
		customError = DefaultCustomError
	}

	// Note - retry logic is handled by the client SDK
	resp, r, err := f()

	if err != nil || r.StatusCode >= 300 {
		switch t := err.(type) {
		case pingone.APIError:
			diags.AddError(fmt.Sprintf("Error when calling `%s`: %v", requestID, t.Error()), "")
			tflog.Error(ctx, fmt.Sprintf("Error when calling `%s`: %v\n\nResponse code: %d\nResponse content-type: %s\nFull response body: %+v", requestID, t.Error(), r.StatusCode, r.Header.Get("Content-Type"), r.Body))
		case *url.Error:
			tflog.Warn(ctx, fmt.Sprintf("Detected HTTP error %s\n\nResponse code: %d\nResponse content-type: %s", t.Err.Error(), r.StatusCode, r.Header.Get("Content-Type")))
			diags.AddError(fmt.Sprintf("Error when calling `%s`: %v", requestID, t.Error()), "")
		default:
			// Attempt to marshal the error into pingone.ServiceError
			errorUnmarshaled := false
			errBytes, jsonErr := json.Marshal(t)
			if jsonErr == nil {
				var targetError pingone.GeneralError
				jsonErr = json.Unmarshal(errBytes, &targetError)
				if jsonErr == nil {
					// Apply custom error handler
					diags = customError(r, &targetError)
					// If no custom error handling was applied, format the error for output
					if len(diags) == 0 {
						summaryText, detailText := FormatPingOneError(requestID, targetError)
						diags.AddError(summaryText, detailText)
					}
					errorUnmarshaled = true
				}
			}
			if !errorUnmarshaled {
				tflog.Warn(ctx, fmt.Sprintf("Detected unknown error (SDK) %+v", t))
				diags.AddError(fmt.Sprintf("Error when calling `%s`: %v", requestID, t.Error()), fmt.Sprintf("A generic error has occurred.\nError details: %+v", t))
			}
		}
		return diags
	}

	if targetObject != nil {
		v := reflect.ValueOf(targetObject)
		if v.Kind() != reflect.Ptr {
			diags.AddError(
				"Invalid target object",
				"Target object must be a pointer.  This is always a problem with the provider, please raise an issue with the provider maintainers.",
			)
			return diags
		}
		if !v.Elem().IsValid() {
			diags.AddError(
				"Invalid target object",
				"Target object is not valid.  This is always a problem with the provider, please raise an issue with the provider maintainers.",
			)
			return diags
		}

		if resp != nil {
			v.Elem().Set(reflect.ValueOf(resp))
		}
	}

	return diags

}

func FormatPingOneError(sdkMethod string, v pingone.GeneralError) (summaryText, detailText string) {
	summaryText = fmt.Sprintf("Error when calling `%s`: %v", sdkMethod, v.GetMessage())
	var detailTextBuilder strings.Builder
	detailTextBuilder.WriteString("PingOne Error Details:\n")
	if v.Id != nil {
		detailTextBuilder.WriteString(fmt.Sprintf("ID:\t\t%s\n", v.GetId()))
	}
	if v.Code != nil {
		detailTextBuilder.WriteString(fmt.Sprintf("Code:\t\t%s\n", v.GetCode()))
	}
	detailTextBuilder.WriteString(fmt.Sprintf("Message:\t%s\n", v.GetMessage()))

	if details, ok := v.GetDetailsOk(); ok {

		detailsStrList := make([]string, 0, len(details))
		for _, detail := range details {
			detailsStr := ""
			nextLineMarker := "-"

			if code, ok := detail.GetCodeOk(); ok {
				detailsStr += fmt.Sprintf("  %s Code:\t%s\n", nextLineMarker, *code)
				nextLineMarker = " "
			}

			if message, ok := detail.GetMessageOk(); ok {
				detailsStr += fmt.Sprintf("  %s Message:\t%s\n", nextLineMarker, *message)
				nextLineMarker = " "
			}

			if target, ok := detail.GetTargetOk(); ok {
				detailsStr += fmt.Sprintf("  %s Target:\t%s\n", nextLineMarker, *target)
				nextLineMarker = " "
			}

			if innerError, ok := detail.GetInnerErrorOk(); ok {
				// Attempt to convert the interface map into json bytes
				innerErrorBytes, err := json.Marshal(innerError)
				if err != nil {
					detailsStr += fmt.Sprintf("  %s Data:\n%s", nextLineMarker, string(innerErrorBytes))
				}
			}

			detailsStrList = append(detailsStrList, detailsStr)
		}

		detailTextBuilder.WriteString(fmt.Sprintf("\nDetails:\n%s", strings.Join(detailsStrList, "\n")))
	}

	detailText = detailTextBuilder.String()
	return summaryText, detailText
}
