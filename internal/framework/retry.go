// Copyright © 2025 Ping Identity Corporation

package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/pingidentity/pingone-go-client/pingone"
)

type Retryable func(context.Context, *http.Response, *pingone.GeneralError) bool

var (
	DefaultRetryable = func(ctx context.Context, r *http.Response, p1error *pingone.GeneralError) bool { return false }

	// Similar but not identical to DefaultCreateReadRetryable in sdk/retry.go
	DefaultCreateReadRetryable = func(ctx context.Context, r *http.Response, p1error *pingone.GeneralError) bool {
		if p1error != nil {
			// Permissions may not have propagated by this point
			m, err := regexp.MatchString("^The request could not be completed. You do not have access to this resource.", p1error.GetMessage())
			if err == nil && m {
				tflog.Warn(ctx, "Insufficient PingOne privileges detected")
				return true
			}
			if err != nil {
				return false
			}
		}

		return false
	}
)

func RetryWrapper(ctx context.Context, timeout time.Duration, f SDKInterfaceFunc, requestID string, isRetryable Retryable) (interface{}, *http.Response, error) {
	var resp interface{}
	var r *http.Response

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		var err error

		// SDK handles most typical retry logic already
		resp, r, err = f()

		if err != nil || r.StatusCode >= 300 {

			var errorModel *pingone.GeneralError

			switch t := err.(type) {
			case pingone.APIError:
				tflog.Error(ctx, fmt.Sprintf("Error when calling `%s`: %v\n\nResponse code: %d\nResponse content-type: %s\nFull response body: %+v", requestID, t.Error(), r.StatusCode, r.Header.Get("Content-Type"), r.Body))
			case *url.Error:
				tflog.Warn(ctx, fmt.Sprintf("Detected HTTP error %s\n\nResponse code: %d\nResponse content-type: %s", t.Err.Error(), r.StatusCode, r.Header.Get("Content-Type")))
			default:
				// Attempt to marshal the error into pingone.GeneralError
				errorUnmarshaled := false
				errBytes, jsonErr := json.Marshal(t)
				if jsonErr == nil {
					var targetError pingone.GeneralError
					jsonErr = json.Unmarshal(errBytes, &targetError)
					if jsonErr == nil && isValidGeneralError(targetError) {
						errorModel = &targetError
						errorUnmarshaled = true
					}
				}
				if !errorUnmarshaled {
					tflog.Warn(ctx, fmt.Sprintf("Detected unknown error (SDK) %+v", t))
				}
			}

			if errorModel != nil && isRetryable != nil && isRetryable(ctx, r, errorModel) {
				tflog.Warn(ctx, "Retrying ... ")
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)

		}
		return nil
	})

	if err != nil {
		return nil, r, err
	}

	return resp, r, nil
}
