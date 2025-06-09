// Copyright © 2025 Ping Identity Corporation

package framework

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pingidentity/pingone-go-client/pingone"
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
