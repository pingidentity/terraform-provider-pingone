package framework

import (
	"context"
	"encoding/json"
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

type CustomError func(model.P1Error) diag.Diagnostics

var (
	DefaultCustomError = func(error model.P1Error) diag.Diagnostics { return nil }

	CustomErrorResourceNotFoundWarning = func(error model.P1Error) diag.Diagnostics {
		var diags diag.Diagnostics

		// Deleted outside of TF
		if error.GetCode() == "NOT_FOUND" {
			diags.AddWarning("Requested resource not found", fmt.Sprintf("The requested resource configuration cannot be found in the PingOne service.  If the requested resource is managed in Terraform's state, it may have been removed outside of Terraform.\nAPI error: %s", error.GetMessage()))

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

func CheckEnvironmentExistsOnPermissionsError(ctx context.Context, managementClient *management.APIClient, environmentID string, fO any, fR *http.Response, fErr error) (any, *http.Response, error) {
	if fR.StatusCode == http.StatusUnauthorized || fR.StatusCode == http.StatusForbidden || fR.StatusCode == http.StatusBadRequest {
		_, fER, fEErr := managementClient.EnvironmentsApi.ReadOneEnvironment(ctx, environmentID).Execute()

		if fER.StatusCode == http.StatusNotFound {
			tflog.Warn(ctx, "API responded with 400, 401 or 403, and the provider determined the environment doesn't exist.  Overriding resource response.")
			return fO, fER, fEErr
		}
	}

	return fO, fR, fErr
}

func ParseResponse(ctx context.Context, f sdk.SDKInterfaceFunc, requestID string, customError CustomError, customRetryConditions sdk.Retryable, targetObject any) diag.Diagnostics {
	defaultTimeout := 10
	return ParseResponseWithCustomTimeout(ctx, f, requestID, customError, customRetryConditions, targetObject, time.Duration(defaultTimeout)*time.Minute)
}

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

			if v, ok := t.Model().(model.P1Error); ok && v.GetId() != "" {

				summaryText := fmt.Sprintf("Error when calling `%s`: %v", requestID, v.GetMessage())
				detailText := fmt.Sprintf("PingOne Error Details:\nID: %s\nCode: %s\nMessage: %s", v.GetId(), v.GetCode(), v.GetMessage())

				diags = customError(v)
				if diags != nil {
					return diags
				}

				if details, ok := v.GetDetailsOk(); ok {
					detailsBytes, err := json.Marshal(details)
					if err != nil {
						diags.AddWarning("Cannot parse details object", "There is an internal problem with the provider.  Please raise an issue with the provider's maintainers.")
					}

					detailText = fmt.Sprintf("%s\nDetails object: %+v", detailText, string(detailsBytes[:]))
				}

				diags.AddError(summaryText, detailText)

				return diags
			}

			diags.AddError(fmt.Sprintf("Error when calling `%s`: %v", requestID, t.Error()), "")

			tflog.Error(ctx, fmt.Sprintf("Error when calling `%s`: %v\n\nFull response body: %+v", requestID, t.Error(), r.Body))

			return diags

		case *url.Error:
			tflog.Warn(ctx, fmt.Sprintf("Detected HTTP error %s", t.Err.Error()))

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
