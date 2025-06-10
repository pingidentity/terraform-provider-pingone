// Copyright © 2025 Ping Identity Corporation

package framework

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type Retryable func(context.Context, *http.Response, error) bool

var (
	DefaultRetryable = func(ctx context.Context, r *http.Response, p1error error) bool { return false }

	DefaultCreateReadRetryable = func(ctx context.Context, r *http.Response, p1error error) bool {

		if p1error != nil {
			var err error

			// Permissions may not have propagated by this point
			m, err := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.Error())
			if err == nil && m {
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
)

func RetryWrapper(ctx context.Context, timeout time.Duration, f SDKInterfaceFunc, isRetryable Retryable) (interface{}, *http.Response, error) {

	var resp interface{}
	var r *http.Response

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		var err error

		resp, r, err = f()

		if err != nil || r.StatusCode >= 300 {
			//TODO is it necessary to check for an id here?
			if isRetryable != nil && isRetryable(ctx, r, err) {
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
