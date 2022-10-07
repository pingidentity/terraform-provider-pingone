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
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
)

type SDKInterfaceFunc func() (interface{}, *http.Response, error)
type CustomError func(interface{}) diag.Diagnostics

var (
	DefaultCustomError = func(error interface{}) diag.Diagnostics { return nil }

	CustomErrorResourceNotFoundWarning = func(error interface{}) diag.Diagnostics {
		var diags diag.Diagnostics

		errorObj, err := model.RemarshalErrorObj(error)
		if err != nil {
			return diag.FromErr(err)
		}

		// Deleted outside of TF
		if errorObj.GetCode() == "NOT_FOUND" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  errorObj.GetMessage(),
			})

			return diags
		}

		return nil
	}

	CustomErrorInvalidValue = func(error interface{}) diag.Diagnostics {
		var diags diag.Diagnostics

		errorObj, err := model.RemarshalErrorObj(error)
		if err != nil {
			return diag.FromErr(err)
		}

		if details, ok := errorObj.GetDetailsOk(); ok && details != nil && len(details) > 0 {
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

		var error *model.GenericOpenAPIError

		switch t := err.(type) {
		case *management.GenericOpenAPIError:
			error, err = model.RemarshalGenericOpenAPIErrorObj(t)
			if err != nil {
				return nil, diag.FromErr(err)
			}

		case *mfa.GenericOpenAPIError:
			error, err = model.RemarshalGenericOpenAPIErrorObj(t)
			if err != nil {
				return nil, diag.FromErr(err)
			}

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
				Summary:  fmt.Sprintf("Error when calling `%s`: %v", sdkMethod, t.Error()),
				Detail:   fmt.Sprintf("A generic error has occurred.\nError details: %+v", t),
			})

			return nil, diags
		}

		if error.Model() != nil {
			errorModel := error.Model().(model.P1Error)

			summaryText := fmt.Sprintf("Error when calling `%s`: %v", sdkMethod, errorModel.GetMessage())
			detailText := fmt.Sprintf("PingOne Error Details:\nID: %s\nCode: %s\nMessage: %s", errorModel.GetId(), errorModel.GetCode(), errorModel.GetMessage())

			diags = customError(errorModel)
			if diags != nil {
				return nil, diags
			}

			if details, ok := errorModel.GetDetailsOk(); ok {
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

	}

	return resp, diags

}
