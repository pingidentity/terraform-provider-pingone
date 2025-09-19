// Copyright © 2025 Ping Identity Corporation

// Package framework provides utilities for Terraform Plugin Framework implementation in the PingOne provider.
// This package contains SDK response parsing, error handling, and retry logic for PingOne API interactions.
package framework

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// CustomError defines a function type for handling custom error responses from the PingOne API.
// It takes an HTTP response and a P1Error model and returns diagnostics for error handling.
// Custom error functions allow specific resources to implement specialized error handling logic
// beyond the default error processing, such as treating certain errors as warnings or providing
// additional context for specific error conditions.
type CustomError func(*http.Response, *model.P1Error) diag.Diagnostics

var (
	// DefaultCustomError is the default custom error handler that performs no special processing.
	// It always returns nil diagnostics, allowing the standard error handling logic to proceed.
	DefaultCustomError = func(_ *http.Response, _ *model.P1Error) diag.Diagnostics { return nil }

	// CustomErrorResourceNotFoundWarning handles resource not found errors by converting them to warnings.
	// This is useful for read operations where a missing resource should not cause an error but should
	// warn the user that the resource may have been deleted outside of Terraform.
	// It checks for HTTP 404 status codes or P1Error codes of "NOT_FOUND".
	CustomErrorResourceNotFoundWarning = func(r *http.Response, p1Error *model.P1Error) diag.Diagnostics {
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

	// CustomErrorInvalidValue handles invalid value errors by converting them to specific error diagnostics.
	// This custom error handler looks for P1Error responses with "INVALID_VALUE" codes targeting the "name" field
	// and converts them into user-friendly error messages. This is commonly used for validation errors
	// where the API rejects a provided name value due to format or constraint violations.
	CustomErrorInvalidValue = func(_ *http.Response, p1Error *model.P1Error) diag.Diagnostics {
		var diags diag.Diagnostics

		// Value not allowed
		if p1Error != nil {
			if details, ok := p1Error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
				if target, ok := details[0].GetTargetOk(); ok && details[0].GetCode() == "INVALID_VALUE" && *target == "name" {
					diags.AddError("Invalid Value", details[0].GetMessage())

					return diags
				}
			}
		}

		return nil
	}
)

// CheckEnvironmentExistsOnPermissionsError verifies environment existence when permission errors occur.
// It returns the original response objects, potentially modified to reflect environment not found status.
// This function is used to distinguish between permission errors on existing environments versus
// permission errors caused by non-existent environments. When a 400, 401, or 403 error occurs,
// it makes an additional API call to check if the environment exists, and if not, overrides the
// original error response to indicate the environment was not found.
//
// The ctx parameter provides the context for the API call.
// The managementClient parameter is the PingOne management API client.
// The environmentID parameter is the ID of the environment to check.
// The fO, fR, and fErr parameters are the original response objects from the failed API call.
func CheckEnvironmentExistsOnPermissionsError(ctx context.Context, managementClient *management.APIClient, environmentID string, fO any, fR *http.Response, fErr error) (any, *http.Response, error) {
	if fR != nil && (fR.StatusCode == http.StatusUnauthorized || fR.StatusCode == http.StatusForbidden || fR.StatusCode == http.StatusBadRequest) {
		_, fER, fEErr := managementClient.EnvironmentsApi.ReadOneEnvironment(ctx, environmentID).Execute()

		if fER.StatusCode == http.StatusNotFound {
			tflog.Warn(ctx, "API responded with 400, 401 or 403, and the provider determined the environment doesn't exist.  Overriding resource response.")
			return fO, fER, fEErr
		}
	}

	return fO, fR, fErr
}

// ParseResponse processes SDK interface function responses with error handling and retry logic.
// It returns diagnostics for any errors encountered during the API call or response processing.
// This function wraps SDK API calls with standardized error handling, retry logic, and response parsing.
// It uses a default timeout of 10 minutes for the operation.
//
// The ctx parameter provides the context for the API call.
// The f parameter is the SDK interface function to execute.
// The requestID parameter is used for error reporting and logging.
// The customError parameter defines custom error handling logic (can be nil for default handling).
// The customRetryConditions parameter defines custom retry logic (can be nil for default retry conditions).
// The targetObject parameter receives the parsed response object (can be nil if response is not needed).
func ParseResponse(ctx context.Context, f sdk.SDKInterfaceFunc, requestID string, customError CustomError, customRetryConditions sdk.Retryable, targetObject any) diag.Diagnostics {
	defaultTimeout := 10
	return ParseResponseWithCustomTimeout(ctx, f, requestID, customError, customRetryConditions, targetObject, time.Duration(defaultTimeout)*time.Minute)
}

// ParseResponseWithCustomTimeout processes SDK interface function responses with custom timeout.
// It returns diagnostics for any errors encountered during the API call or response processing.
// This function provides the same functionality as ParseResponse but allows specifying a custom timeout
// duration for operations that may require longer processing times than the default 10 minutes.
//
// The ctx parameter provides the context for the API call.
// The f parameter is the SDK interface function to execute.
// The requestID parameter is used for error reporting and logging.
// The customError parameter defines custom error handling logic (can be nil for default handling).
// The customRetryConditions parameter defines custom retry logic (can be nil for default retry conditions).
// The targetObject parameter receives the parsed response object (can be nil if response is not needed).
// The timeout parameter specifies the maximum duration to wait for the operation to complete.
func ParseResponseWithCustomTimeout(ctx context.Context, f sdk.SDKInterfaceFunc, requestID string, customError CustomError, customRetryConditions sdk.Retryable, targetObject any, timeout time.Duration) diag.Diagnostics {
	var diags diag.Diagnostics

	if customError == nil {
		customError = DefaultCustomError
	}

	if customRetryConditions == nil {
		customRetryConditions = sdk.DefaultRetryable
	}

	resp, r, err := sdk.RetryWrapper(
		ctx,
		timeout,
		f,
		customRetryConditions,
	)

	if err != nil || r.StatusCode >= 300 {

		switch t := err.(type) {
		case *model.GenericOpenAPIError:

			model, ok := t.Model().(model.P1Error)

			diags = customError(r, &model)
			if diags != nil {
				return diags
			}

			if ok && model.GetId() != "" {
				summaryText, detailText := sdk.FormatPingOneError(requestID, model)

				diags.AddError(summaryText, detailText)

				return diags
			}

			diags.AddError(fmt.Sprintf("Error when calling `%s`: %v", requestID, t.Error()), "")

			tflog.Error(ctx, fmt.Sprintf("Error when calling `%s`: %v\n\nResponse code: %d\nResponse content-type: %s\nFull response body: %+v", requestID, t.Error(), r.StatusCode, r.Header.Get("Content-Type"), r.Body))

			return diags

		case *url.Error:
			tflog.Warn(ctx, fmt.Sprintf("Detected HTTP error %s\n\nResponse code: %d\nResponse content-type: %s", t.Err.Error(), r.StatusCode, r.Header.Get("Content-Type")))

			diags.AddError(fmt.Sprintf("Error when calling `%s`: %v", requestID, t.Error()), "")

			return diags

		default:
			tflog.Warn(ctx, fmt.Sprintf("Detected unknown error (SDK) %+v", t))

			diags.AddError(fmt.Sprintf("Error when calling `%s`: %v", requestID, t.Error()), fmt.Sprintf("A generic error has occurred.\nError details: %+v", t))

			return diags
		}

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
