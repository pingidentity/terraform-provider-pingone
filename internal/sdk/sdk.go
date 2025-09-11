// Copyright © 2025 Ping Identity Corporation

// Package sdk provides SDK wrapper functions and error handling utilities for the PingOne Terraform provider.
// This package brokers the interaction between PingOne SDK functions and the Terraform provider using the v5 protocol/SDKv2 SDK.
// It includes functions for processing API responses, handling retries, and formatting error messages.
package sdk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
)

// SDKInterfaceFunc represents a function signature for PingOne SDK API calls.
// This function type is used as a wrapper for all SDK method invocations to enable
// consistent error handling and retry logic across the provider.
type SDKInterfaceFunc func() (any, *http.Response, error)

// CustomError represents a function that processes PingOne API errors and returns custom diagnostics.
// This allows resources to override default error handling with context-specific error messages
// or warnings based on the API error details.
type CustomError func(model.P1Error) diag.Diagnostics

var (
	// DefaultCustomError is the default error handler that returns no custom diagnostics.
	// This allows the standard error formatting to be used when no custom error handling is required.
	DefaultCustomError = func(error model.P1Error) diag.Diagnostics { return nil }

	// CustomErrorResourceNotFoundWarning provides custom error handling for resource not found scenarios.
	// It converts NOT_FOUND errors to warnings instead of errors, useful for cases where resources
	// may have been deleted outside of Terraform and should be gracefully handled.
	CustomErrorResourceNotFoundWarning = func(error model.P1Error) diag.Diagnostics {
		var diags diag.Diagnostics

		// Deleted outside of TF
		if error.GetCode() == "NOT_FOUND" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  error.GetMessage(),
			})

			return diags
		}

		return nil
	}

	// CustomErrorInvalidValue provides custom error handling for invalid value scenarios.
	// It extracts specific error details for INVALID_VALUE errors targeting the "name" field
	// and formats them as user-friendly error messages.
	CustomErrorInvalidValue = func(error model.P1Error) diag.Diagnostics {
		var diags diag.Diagnostics

		// Value not allowed
		if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
			if target, ok := details[0].GetTargetOk(); ok && details[0].GetCode() == "INVALID_VALUE" && *target == "name" {
				diags = diag.FromErr(fmt.Errorf("%s", details[0].GetMessage()))

				return diags
			}
		}

		return nil
	}
)

// ParseResponse processes the result of a PingOne SDK API call with standard error handling and retry logic.
// It returns the API response data and any diagnostics encountered during processing.
// The f parameter is the SDK function to execute, which should return the API response data, HTTP response, and any error.
// The sdkMethod parameter is used for logging and error identification purposes.
// The customError parameter allows overriding default error handling with resource-specific error processing.
// The customRetryConditions parameter defines when API calls should be retried based on response conditions.
// This function uses a default timeout of 10 minutes for retry operations.
func ParseResponse(ctx context.Context, f SDKInterfaceFunc, sdkMethod string, customError CustomError, customRetryConditions Retryable) (interface{}, diag.Diagnostics) {
	defaultTimeout := 10
	return ParseResponseWithCustomTimeout(ctx, f, sdkMethod, customError, customRetryConditions, time.Duration(defaultTimeout)*time.Minute)
}

// ParseResponseWithCustomTimeout processes the result of a PingOne SDK API call with custom timeout settings.
// It returns the API response data and any diagnostics encountered during processing.
// The f parameter is the SDK function to execute, which should return the API response data, HTTP response, and any error.
// The sdkMethod parameter is used for logging and error identification purposes.
// The customError parameter allows overriding default error handling with resource-specific error processing.
// The customRetryConditions parameter defines when API calls should be retried based on response conditions.
// The timeout parameter specifies the maximum duration to wait for successful API completion including retries.
func ParseResponseWithCustomTimeout(ctx context.Context, f SDKInterfaceFunc, sdkMethod string, customError CustomError, customRetryConditions Retryable, timeout time.Duration) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if customError == nil {
		customError = DefaultCustomError
	}

	if customRetryConditions == nil {
		customRetryConditions = DefaultRetryable
	}

	resp, r, err := RetryWrapper(
		ctx,
		timeout,
		f,
		customRetryConditions,
	)

	if err != nil || r.StatusCode >= 300 {

		switch t := err.(type) {
		case *model.GenericOpenAPIError:

			if v, ok := t.Model().(model.P1Error); ok && v.GetId() != "" {

				diags = customError(v)
				if diags != nil {
					return nil, diags
				}

				summaryText, detailText := FormatPingOneError(sdkMethod, v)

				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  summaryText,
					Detail:   detailText,
				})

				return nil, diags
			}

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `%s`: %v", sdkMethod, t.Error()),
			})

			tflog.Error(ctx, fmt.Sprintf("Error when calling `%s`: %v\n\nFull response body: %+v", sdkMethod, t.Error(), r.Body))

			return nil, diags

		case *url.Error:
			tflog.Warn(ctx, fmt.Sprintf("Detected HTTP error %s", t.Err.Error()))

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `%s`: %v", sdkMethod, t.Err.Error()),
			})

			return nil, diags

		default:
			tflog.Warn(ctx, fmt.Sprintf("Detected unknown error (SDK) %+v", t))

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `%s`: %v", sdkMethod, t.Error()),
				Detail:   fmt.Sprintf("A generic error has occurred.\nError details: %+v", t),
			})

			return nil, diags
		}

	}

	return resp, diags

}

// FormatPingOneError creates formatted error summary and detail messages from a PingOne API error.
// It returns both a summary text suitable for brief error display and detailed text with comprehensive error information.
// The sdkMethod parameter identifies the API method that generated the error for context.
// The v parameter contains the PingOne error details including error code, message, and nested error information.
// The formatted output includes error ID, code, message, and detailed breakdowns of any nested error constraints or validation rules.
func FormatPingOneError(sdkMethod string, v model.P1Error) (summaryText, detailText string) {
	summaryText = fmt.Sprintf("Error when calling `%s`: %v", sdkMethod, v.GetMessage())
	detailText = fmt.Sprintf("PingOne Error Details:\nID:\t\t%s\nCode:\t\t%s\nMessage:\t%s", v.GetId(), v.GetCode(), v.GetMessage())

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
				innerDetailsStr := ""

				if v, ok := innerError.GetRangeMinimumValueOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Range Min Value:\t%d\n", *v)
				}

				if v, ok := innerError.GetRangeMaximumValueOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Range Max Value:\t%d\n", *v)
				}

				if v, ok := innerError.GetAllowedPatternOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Allowed Pattern:\t%s\n", *v)
				}

				if v, ok := innerError.GetAllowedValuesOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Allowed Values:\t[%s]\n", strings.Join(v, ", "))
				}

				if v, ok := innerError.GetMaximumValueOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Max Value:\t%d\n", *v)
				}

				if v, ok := innerError.GetReferencedValuesOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Referenced Values:\t[%s]\n", strings.Join(v, ", "))
				}

				detailsStr += fmt.Sprintf("  %s Data:\n%s", nextLineMarker, innerDetailsStr)
			}

			detailsStrList = append(detailsStrList, detailsStr)
		}

		detailText += fmt.Sprintf("\nDetails:\n%s", strings.Join(detailsStrList, "\n"))
	}

	return summaryText, detailText
}
