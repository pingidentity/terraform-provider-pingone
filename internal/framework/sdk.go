package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

type CustomError func(model.P1Error) diag.Diagnostics

var (
	DefaultCustomError = func(error model.P1Error) diag.Diagnostics { return nil }

	CustomErrorResourceNotFoundWarning = func(error model.P1Error) diag.Diagnostics {
		var diags diag.Diagnostics

		// Deleted outside of TF
		if error.GetCode() == "NOT_FOUND" {
			diags.AddWarning("Resource not found", fmt.Sprintf("The requested resource object cannot be found.  Error returned: %s.", error.GetMessage()))

			return diags
		}

		return nil
	}

	CustomErrorInvalidValue = func(error model.P1Error) diag.Diagnostics {
		var diags diag.Diagnostics

		// Value not allowed
		if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
			if target, ok := details[0].GetTargetOk(); ok && details[0].GetCode() == "INVALID_VALUE" && *target == "name" {
				diags.AddError("Invalid Value", details[0].GetMessage())

				return diags
			}
		}

		return nil
	}
)

func ParseResponse(ctx context.Context, f sdk.SDKInterfaceFunc, requestID string, customError CustomError, customRetryConditions sdk.Retryable) (interface{}, diag.Diagnostics) {
	defaultTimeout := 10
	return ParseResponseWithCustomTimeout(ctx, f, requestID, customError, customRetryConditions, time.Duration(defaultTimeout)*time.Minute)
}

func ParseResponseWithCustomTimeout(ctx context.Context, f sdk.SDKInterfaceFunc, requestID string, customError CustomError, customRetryConditions sdk.Retryable, timeout time.Duration) (interface{}, diag.Diagnostics) {
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

			if v, ok := t.Model().(model.P1Error); ok && v.GetId() != "" {

				summaryText := fmt.Sprintf("Error when calling `%s`: %v", requestID, v.GetMessage())
				detailText := fmt.Sprintf("PingOne Error Details:\nID: %s\nCode: %s\nMessage: %s", v.GetId(), v.GetCode(), v.GetMessage())

				diags = customError(v)
				if diags != nil {
					return nil, diags
				}

				if details, ok := v.GetDetailsOk(); ok {
					detailsBytes, err := json.Marshal(details)
					if err != nil {
						diags.AddWarning("Cannot parse details object", "There is an internal problem with the provider.  Please raise an issue with the provider's maintainers.")
					}

					detailText = fmt.Sprintf("%s\nDetails object: %+v", detailText, string(detailsBytes[:]))
				}

				diags.AddError(summaryText, detailText)

				return nil, diags
			}

			diags.AddError(fmt.Sprintf("Error when calling `%s`: %v", requestID, t.Error()), "")

			tflog.Error(ctx, fmt.Sprintf("Error when calling `%s`: %v\n\nFull response body: %+v", requestID, t.Error(), r.Body))

			return nil, diags

		case *url.Error:
			tflog.Warn(ctx, fmt.Sprintf("Detected HTTP error %s", t.Err.Error()))

			diags.AddError(fmt.Sprintf("Error when calling `%s`: %v", requestID, t.Error()), "")

			return nil, diags

		default:
			tflog.Warn(ctx, fmt.Sprintf("Detected unknown error (SDK) %+v", t))

			diags.AddError(fmt.Sprintf("Error when calling `%s`: %v", requestID, t.Error()), fmt.Sprintf("A generic error has occurred.\nError details: %+v", t))

			return nil, diags
		}

	}

	return resp, diags

}
