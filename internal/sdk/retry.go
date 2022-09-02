package sdk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

type Retryable func(context.Context, *http.Response, *management.P1Error) bool

var (
	DefaultRetryable = func(ctx context.Context, r *http.Response, p1error *management.P1Error) bool {

		// Gateway errors
		if r.StatusCode >= 502 && r.StatusCode <= 504 {
			tflog.Warn(ctx, "Gateway error detected, available for retry")
			return true
		}

		return false
	}

	DefaultCreateReadRetryable = func(ctx context.Context, r *http.Response, p1error *management.P1Error) bool {

		if p1error != nil {
			var err error

			// Permissions may not have propagated by this point
			if m, err := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
				tflog.Warn(ctx, "Insufficient PingOne privileges detected")
				return true
			}
			if err != nil {
				tflog.Warn(ctx, "Cannot match error string for retry")
				return false
			}

		}

		return false
	}

	RoleAssignmentRetryable = func(ctx context.Context, r *http.Response, p1error *management.P1Error) bool {

		if p1error != nil {
			var err error

			// Permissions may not have propagated by this point (1)
			if m, err := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
				tflog.Warn(ctx, "Insufficient PingOne privileges detected")
				return true
			}
			if err != nil {
				tflog.Warn(ctx, "Cannot match error string for retry")
				return false
			}

			// Permissions may not have propagated by this point (2)
			if details, ok := p1error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
				if m, err := regexp.MatchString("^Must have role at the same or broader scope", details[0].GetMessage()); err == nil && m {
					tflog.Warn(ctx, "Insufficient PingOne privileges detected")
					return true
				}
				if err != nil {
					tflog.Warn(ctx, "Cannot match error string for retry")
					return false
				}
			}

		}

		return false
	}
)

func RetryWrapper(ctx context.Context, timeout time.Duration, f SDKInterfaceFunc, isRetryable Retryable) (interface{}, *http.Response, error) {

	var resp interface{}
	var r *http.Response

	err := resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		var err error

		resp, r, err = f()

		if err != nil || r.StatusCode >= 300 {

			var model management.P1Error

			switch t := err.(type) {
			case *management.GenericOpenAPIError:
				error := t

				if error.Model() != nil {
					model = error.Model().(management.P1Error)
				}

			case *url.Error:
				tflog.Warn(ctx, fmt.Sprintf("Detected URL error %s", t.Err.Error()))

			default:
				tflog.Warn(ctx, fmt.Sprintf("Detected unknown error %+v", t))
			}

			if (model.Id != nil || r != nil) && (isRetryable(ctx, r, &model) || DefaultRetryable(ctx, r, &model)) {
				tflog.Warn(ctx, "Retrying ... ")
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)

		}
		return nil
	})

	if err != nil {
		return nil, r, err
	}

	return resp, r, nil
}
