package sdk

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

type Retryable func(*http.Response, error) (bool, error)

var (
	DefaultRetryable = func(r *http.Response, err error) (bool, error) { return false, err }
)

func RetryWrapper(ctx context.Context, timeout time.Duration, f SDKInterfaceFunc, isRetryable Retryable) (interface{}, *http.Response, error) {

	var resp interface{}
	var r *http.Response

	err := resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		var err error
		var retry bool

		resp, r, err = f()

		if err != nil || r.StatusCode >= 300 {
			error := err.(*management.GenericOpenAPIError)

			retry, err = isRetryable(r, error)

			if retry {
				return resource.RetryableError(err)
			}

			if err != nil {
				return resource.NonRetryableError(err)
			}

		}
		return nil
	})

	if err != nil {
		return nil, r, err
	}

	return resp, r, nil
}
