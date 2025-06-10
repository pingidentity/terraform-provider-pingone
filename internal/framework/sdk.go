// Copyright © 2025 Ping Identity Corporation

package framework

import (
	"context"
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

type pingOneError struct {
	message string
	code    *string
}

type CustomError func(*http.Response, *pingOneError) diag.Diagnostics

var (
	DefaultCustomError = func(_ *http.Response, _ *pingOneError) diag.Diagnostics { return nil }

	CustomErrorResourceNotFoundWarning = func(r *http.Response, p1Error *pingOneError) diag.Diagnostics {
		var diags diag.Diagnostics

		// Deleted outside of TF
		if p1Error != nil && p1Error.code != nil && *p1Error.code == "NOT_FOUND" {
			diags.AddWarning("Requested resource not found", fmt.Sprintf("The requested resource configuration cannot be found in the PingOne service.  If the requested resource is managed in Terraform's state, it may have been removed outside of Terraform.\nAPI error: %s", p1Error.message))

			return diags
		}

		if r != nil && r.StatusCode == 404 {
			diags.AddWarning("Requested resource not found", "The requested resource configuration cannot be found in the PingOne service.  If the requested resource is managed in Terraform's state, it may have been removed outside of Terraform.")

			return diags
		}

		return nil
	}
)

func CheckEnvironmentExistsOnPermissionsError(ctx context.Context, apiClient *pingone.APIClient, environmentID string, fO any, fR *http.Response, fErr error) (any, *http.Response, error) {
	if fR != nil && (fR.StatusCode == http.StatusUnauthorized || fR.StatusCode == http.StatusForbidden || fR.StatusCode == http.StatusBadRequest) {
		environmentIdUuid, err := uuid.Parse(environmentID)
		if err != nil {
			return fO, nil, fmt.Errorf("unable to parse environment id '%s' as uuid: %v", environmentID, err)
		}

		_, fER, fEErr := apiClient.EnvironmentApi.GetEnvironmentById(ctx, environmentIdUuid).Execute()

		if fER.StatusCode == http.StatusNotFound {
			tflog.Warn(ctx, "API responded with 400, 401 or 403, and the provider determined the environment doesn't exist.  Overriding resource response.")
			return fO, fER, fEErr
		}
	}

	return fO, fR, fErr
}

func ParseResponse(ctx context.Context, f SDKInterfaceFunc, requestID string, customError CustomError, customRetryConditions Retryable, targetObject any) diag.Diagnostics {
	defaultTimeout := 10
	return ParseResponseWithCustomTimeout(ctx, f, requestID, customError, customRetryConditions, targetObject, time.Duration(defaultTimeout)*time.Minute)
}

func ParseResponseWithCustomTimeout(ctx context.Context, f SDKInterfaceFunc, requestID string, customError CustomError, customRetryConditions Retryable, targetObject any, timeout time.Duration) diag.Diagnostics {
	var diags diag.Diagnostics

	if customError == nil {
		customError = DefaultCustomError
	}

	resp, r, err := RetryWrapper(
		ctx,
		timeout,
		f,
		customRetryConditions,
	)

	if err != nil || r.StatusCode >= 300 {

		var genericError *pingone.ServiceError

		switch t := err.(type) {
		case pingone.AccessFailedError:
			code := string(t.Code)
			diags = customError(r, &pingOneError{
				message: t.Message,
				code:    &code,
			})
			if diags != nil {
				return diags
			}

			genericError = &pingone.ServiceError{
				Details:              t.Details,
				Id:                   t.Id,
				Message:              t.Message,
				Code:                 &code,
				AdditionalProperties: t.AdditionalProperties,
			}
		case pingone.InvalidDataError:
			code := string(t.Code)
			diags = customError(r, &pingOneError{
				message: t.Message,
				code:    &code,
			})
			if diags != nil {
				return diags
			}

			genericError = &pingone.ServiceError{
				Details:              t.Details,
				Id:                   t.Id,
				Message:              t.Message,
				Code:                 &code,
				AdditionalProperties: t.AdditionalProperties,
			}
		case pingone.InvalidRequestError:
			code := string(t.Code)
			diags = customError(r, &pingOneError{
				message: t.Message,
				code:    &code,
			})
			if diags != nil {
				return diags
			}

			genericError = &pingone.ServiceError{
				Details:              t.Details,
				Id:                   t.Id,
				Message:              t.Message,
				Code:                 &code,
				AdditionalProperties: t.AdditionalProperties,
			}
		case pingone.NotFoundError:
			code := string(t.Code)
			diags = customError(r, &pingOneError{
				message: t.Message,
				code:    &code,
			})
			if diags != nil {
				return diags
			}

			genericError = &pingone.ServiceError{
				Details:              t.Details,
				Id:                   t.Id,
				Message:              t.Message,
				Code:                 &code,
				AdditionalProperties: t.AdditionalProperties,
			}
		case pingone.RequestFailedError:
			code := string(t.Code)
			diags = customError(r, &pingOneError{
				message: t.Message,
				code:    &code,
			})
			if diags != nil {
				return diags
			}

			genericError = &pingone.ServiceError{
				Details:              t.Details,
				Id:                   t.Id,
				Message:              t.Message,
				Code:                 &code,
				AdditionalProperties: t.AdditionalProperties,
			}
		case pingone.RequestLimitedError:
			code := string(t.Code)
			diags = customError(r, &pingOneError{
				message: t.Message,
				code:    &code,
			})
			if diags != nil {
				return diags
			}

			genericError = &pingone.ServiceError{
				Details:              t.Details,
				Id:                   t.Id,
				Message:              t.Message,
				Code:                 &code,
				AdditionalProperties: t.AdditionalProperties,
			}
		case pingone.ServiceError:
			diags = customError(r, &pingOneError{
				message: t.Message,
				code:    t.Code,
			})
			if diags != nil {
				return diags
			}

			genericError = err.(*pingone.ServiceError)
		case pingone.UnexpectedServiceError:
			code := string(t.Code)
			diags = customError(r, &pingOneError{
				message: t.Message,
				code:    &code,
			})
			if diags != nil {
				return diags
			}

			genericError = &pingone.ServiceError{
				Details:              t.Details,
				Id:                   t.Id,
				Message:              t.Message,
				Code:                 &code,
				AdditionalProperties: t.AdditionalProperties,
			}
		case pingone.APIError:
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

		if genericError != nil {
			summaryText, detailText := FormatPingOneError(requestID, *genericError)

			diags.AddError(summaryText, detailText)

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

func FormatPingOneError(sdkMethod string, v pingone.ServiceError) (summaryText, detailText string) {
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

				if v, ok := innerError.GetAllowedPatternOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Allowed Pattern:\t%s\n", *v)
				}

				//TODO this field currently is just an array of interfaces
				/*if v, ok := innerError.GetAllowedValuesOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Allowed Values:\t[%s]\n", strings.Join(v, ", "))
				}*/

				if v, ok := innerError.GetMaximumValueOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Max Value:\t%f\n", *v)
				}

				if v, ok := innerError.GetQuotaLimitOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Quota Limit:\t%f\n", *v)
				}

				if v, ok := innerError.GetQuotaResetTimeOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Quota Reset Time:\t%s\n", v.Format(time.RFC3339))
				}

				if v, ok := innerError.GetRangeMaximumValueOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Range Max Value:\t%f\n", *v)
				}

				if v, ok := innerError.GetRangeMinimumValueOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Range Min Value:\t%f\n", *v)
				}

				if v, ok := innerError.GetRetryAfterOk(); ok {
					innerDetailsStr += fmt.Sprintf("      Referenced Values:\t%s\n", *v)
				}

				detailsStr += fmt.Sprintf("  %s Data:\n%s", nextLineMarker, innerDetailsStr)
			}

			detailsStrList = append(detailsStrList, detailsStr)
		}

		detailText += fmt.Sprintf("\nDetails:\n%s", strings.Join(detailsStrList, "\n"))
	}

	return summaryText, detailText
}
