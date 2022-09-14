package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

type SDKInterfaceFunc func() (interface{}, *http.Response, error)
type CustomError func(management.P1Error) diag.Diagnostics

var (
	DefaultCustomError = func(error management.P1Error) diag.Diagnostics { return nil }

	CustomErrorResourceNotFoundWarning = func(error management.P1Error) diag.Diagnostics {
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

	CustomErrorInvalidValue = func(error management.P1Error) diag.Diagnostics {
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

func ParseResponse(ctx context.Context, f SDKInterfaceFunc, sdkMethod string, customError CustomError, retryable Retryable) (interface{}, diag.Diagnostics) {
	defaultTimeout := 10
	return ParseResponseWithCustomTimeout(ctx, f, sdkMethod, customError, retryable, time.Duration(defaultTimeout)*time.Minute)
}

func ParseResponseWithCustomTimeout(ctx context.Context, f SDKInterfaceFunc, sdkMethod string, customError CustomError, retryable Retryable, timeout time.Duration) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if customError == nil {
		customError = DefaultCustomError
	}

	if retryable == nil {
		retryable = DefaultRetryable
	}

	resp, r, err := RetryWrapper(
		ctx,
		timeout,
		f,
		retryable,
	)

	if err != nil || r.StatusCode >= 300 {

		switch t := err.(type) {
		case *management.GenericOpenAPIError:
			error := t

			if error.Model() != nil {
				model := error.Model().(management.P1Error)

				summaryText := fmt.Sprintf("Error when calling `%s`: %v", sdkMethod, model.GetMessage())
				detailText := fmt.Sprintf("PingOne Error Details:\nID: %s\nCode: %s\nMessage: %s", model.GetId(), model.GetCode(), model.GetMessage())

				diags = customError(model)
				if diags != nil {
					return nil, diags
				}

				if details, ok := model.GetDetailsOk(); ok {
					detailsBytes, err := json.Marshal(details)
					if err != nil {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Warning,
							Summary:  "Cannot parse details object",
						})
					}

					detailText = fmt.Sprintf("%s\nDetails object: %+v", detailText, string(detailsBytes[:]))
				}

				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  summaryText,
					Detail:   detailText,
				})

				return nil, diags
			}

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `%s`: %v", sdkMethod, error.Error()),
				Detail:   fmt.Sprintf("Full response body: %+v", r.Body),
			})

			return nil, diags

		case *url.Error:
			tflog.Warn(ctx, fmt.Sprintf("Detected HTTP error %s", t.Err.Error()))

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `%s`: %v", sdkMethod, t.Err.Error()),
			})

			return nil, diags

		default:
			tflog.Warn(ctx, fmt.Sprintf("Detected unknown error %+v", t))

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `%s`", sdkMethod),
				Detail:   "The error is not properly handled and further information is not available.  Please raise an issue with the provider maintainers.",
			})

			return nil, diags
		}

	}

	return resp, diags

}
