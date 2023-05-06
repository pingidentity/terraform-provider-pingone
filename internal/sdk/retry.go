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
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
)

type Retryable func(context.Context, *http.Response, *model.P1Error) bool

var (
	DefaultRetryable = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool { return false }

	DefaultCreateReadRetryable = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

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

	RoleAssignmentRetryable = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

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

		// error could be management, mfa, authorize, credentials
		resp, r, err = f()

		if err != nil || r.StatusCode >= 300 {

			var errorModel *model.P1Error
			var err1 error

			switch t := err.(type) {
			case *authorize.GenericOpenAPIError:

				if t.Model() != nil {
					errorModel, err1 = model.RemarshalErrorObj(t.Model().(authorize.P1Error))
					if err1 != nil {
						tflog.Error(ctx, fmt.Sprintf("Cannot remarshal type %s", err1))
						return resource.NonRetryableError(err)
					}
				}

				err, err1 = model.RemarshalGenericOpenAPIErrorObj(t)
				if err1 != nil {
					tflog.Error(ctx, fmt.Sprintf("Cannot remarshal type %s", err1))
					return resource.NonRetryableError(err)
				}

			case *credentials.GenericOpenAPIError:

				if t.Model() != nil {
					errorModel, err1 = model.RemarshalErrorObj(t.Model().(credentials.P1Error))
					if err1 != nil {
						tflog.Error(ctx, fmt.Sprintf("Cannot remarshal type %s", err1))
						return resource.NonRetryableError(err)
					}
				}

				err, err1 = model.RemarshalGenericOpenAPIErrorObj(t)
				if err1 != nil {
					tflog.Error(ctx, fmt.Sprintf("Cannot remarshal type %s", err1))
					return resource.NonRetryableError(err)
				}

			case *management.GenericOpenAPIError:

				if t.Model() != nil {
					errorModel, err1 = model.RemarshalErrorObj(t.Model().(management.P1Error))
					if err1 != nil {
						tflog.Error(ctx, fmt.Sprintf("Cannot remarshal type %s", err1))
						return resource.NonRetryableError(err)
					}
				}

				err, err1 = model.RemarshalGenericOpenAPIErrorObj(t)
				if err1 != nil {
					tflog.Error(ctx, fmt.Sprintf("Cannot remarshal type %s", err1))
					return resource.NonRetryableError(err)
				}

			case *mfa.GenericOpenAPIError:

				if t.Model() != nil {
					errorModel, err1 = model.RemarshalErrorObj(t.Model().(mfa.P1Error))
					if err1 != nil {
						tflog.Error(ctx, fmt.Sprintf("Cannot remarshal type %s", err1))
						return resource.NonRetryableError(err)
					}
				}

				err, err1 = model.RemarshalGenericOpenAPIErrorObj(t)
				if err1 != nil {
					tflog.Error(ctx, fmt.Sprintf("Cannot remarshal type %s", err1))
					return resource.NonRetryableError(err)
				}

			case *url.Error:
				tflog.Warn(ctx, fmt.Sprintf("Detected HTTP error %s", t.Err.Error()))

			default:
				tflog.Warn(ctx, fmt.Sprintf("Detected unknown error (retry) %+v", t))
			}

			if ((errorModel != nil && errorModel.Id != nil) || r != nil) && (isRetryable(ctx, r, errorModel) || DefaultRetryable(ctx, r, errorModel)) {
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
