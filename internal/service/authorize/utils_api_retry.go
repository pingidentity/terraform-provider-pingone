package authorize

import (
	"context"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

var (
	retryAuthorizeEditorDelete = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

		defaultRetryable := sdk.DefaultRetryable(ctx, r, p1error)
		if defaultRetryable {
			return defaultRetryable
		}

		if p1error != nil {

			// Permissions may not have propagated by this point
			m, err := regexp.MatchString("^REQUEST_FAILED", p1error.GetCode())
			if err == nil && m {
				tflog.Warn(ctx, "Delete request failed.. evaluating details..")

				if details, ok := p1error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
					if message, ok := details[0].GetMessageOk(); ok {
						m, err := regexp.MatchString("entity is referenced", *message)
						if err == nil && m {
							tflog.Warn(ctx, "Delete request failed as entity is referenced.  Retrying...")
							return true
						}
					}
				}
			}
			if err != nil {
				tflog.Warn(ctx, "Cannot match error string for retry")
				return false
			}

		}

		return false
	}

	retryAuthorizeEditorCreateUpdate = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

		defaultRetryable := sdk.DefaultCreateReadRetryable(ctx, r, p1error)
		if defaultRetryable {
			return defaultRetryable
		}

		// if p1error != nil {

		// Permissions may not have propagated by this point
		// m, err := regexp.MatchString("^INVALID_DATA", p1error.GetCode())
		// if err == nil && m {
		// 	tflog.Warn(ctx, "Create/update request failed.. evaluating details..")

		// 	if details, ok := p1error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
		// 		if message, ok := details[0].GetMessageOk(); ok {
		// 			m, err := regexp.MatchString(fmt.Sprintf("No definition with id %s exists", verify.P1ResourceIDRegexp.String()), *message)
		// 			if err == nil && m {
		// 				tflog.Warn(ctx, "Create/update request failed as entity with UUID not found.  Retrying...")
		// 				return true
		// 			}
		// 		}
		// 	}
		// }
		// if err != nil {
		// 	tflog.Warn(ctx, "Cannot match error string for retry")
		// 	return false
		// }

		// }

		return false
	}
)
