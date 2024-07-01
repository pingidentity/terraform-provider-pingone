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

type SDKInterfaceFunc func() (any, *http.Response, error)
type CustomError func(model.P1Error) diag.Diagnostics

var (
	DefaultCustomError = func(error model.P1Error) diag.Diagnostics { return nil }

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

	CustomErrorInvalidValue = func(error model.P1Error) diag.Diagnostics {
		var diags diag.Diagnostics

		// Value not allowed
		if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
			if target, ok := details[0].GetTargetOk(); ok && details[0].GetCode() == "INVALID_VALUE" && *target == "name" {
				diags = diag.FromErr(fmt.Errorf(details[0].GetMessage()))

				return diags
			}
		}

		return nil
	}
)

func ParseResponse(ctx context.Context, f SDKInterfaceFunc, sdkMethod string, customError CustomError, customRetryConditions Retryable) (interface{}, diag.Diagnostics) {
	defaultTimeout := 10
	return ParseResponseWithCustomTimeout(ctx, f, sdkMethod, customError, customRetryConditions, time.Duration(defaultTimeout)*time.Minute)
}

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

	return
}
