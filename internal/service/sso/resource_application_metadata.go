// Copyright © 2026 Ping Identity Corporation

package sso

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// The application metadata sub-resource endpoint is not present in either PingOne
// Go SDK (management@v0.70.0 or pingone-go-client), so requests are constructed
// directly from the management client's configuration, reusing its server URL,
// default (authorization) headers and HTTP client.  API errors are formatted
// consistently with legacysdk.ParseResponse.

type applicationMetadataResponse struct {
	ApplicationId string                 `json:"applicationId,omitempty"`
	EnvironmentId string                 `json:"environmentId,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

type applicationMetadataPayload struct {
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

const applicationMetadataTimeoutMinutes = 10

// readApplicationMetadata issues a GET against the application metadata
// sub-resource.  A 404 response means no metadata is set on the application,
// which is not an error; a nil map is returned in that case.
func readApplicationMetadata(ctx context.Context, managementClient *management.APIClient, environmentID, applicationID string) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	endpoint, err := applicationMetadataEndpoint(ctx, managementClient, environmentID, applicationID)
	if err != nil {
		diags.AddError("Cannot determine application metadata endpoint", err.Error())
		return nil, diags
	}

	resp, body, err := applicationMetadataExecute(ctx, managementClient, http.MethodGet, endpoint, nil)
	if err != nil {
		// No metadata is set on the application; a null value is the correct state
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, diags
		}

		diags.Append(applicationMetadataError("ReadApplicationMetadata", resp, err)...)
		return nil, diags
	}

	var response applicationMetadataResponse
	if err := json.Unmarshal(body, &response); err != nil {
		diags.AddError(
			"Cannot parse application metadata response",
			fmt.Sprintf("The PingOne API returned an unexpected response for `ReadApplicationMetadata`: %s", err.Error()),
		)
		return nil, diags
	}

	return response.Metadata, diags
}

// writeApplicationMetadata issues a PUT against the application metadata
// sub-resource.  A non-nil `metadata` value sets (or fully replaces) the stored
// metadata; a nil value sends an empty body, which removes existing metadata
// (documented API behavior for a PUT without the `metadata` key).
func writeApplicationMetadata(ctx context.Context, managementClient *management.APIClient, environmentID, applicationID string, metadata map[string]interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	endpoint, err := applicationMetadataEndpoint(ctx, managementClient, environmentID, applicationID)
	if err != nil {
		diags.AddError("Cannot determine application metadata endpoint", err.Error())
		return diags
	}

	var payload any
	if metadata != nil {
		payload = applicationMetadataPayload{Metadata: metadata}
	} else {
		payload = &applicationMetadataPayload{}
	}

	resp, _, err := applicationMetadataExecute(ctx, managementClient, http.MethodPut, endpoint, payload)
	if err != nil {
		diags.Append(applicationMetadataError("UpdateApplicationMetadata", resp, err)...)
		return diags
	}

	return diags
}

// applicationMetadataEndpoint resolves the full URL for the application
// metadata sub-resource from the management client configuration.
func applicationMetadataEndpoint(ctx context.Context, managementClient *management.APIClient, environmentID, applicationID string) (string, error) {
	cfg := managementClient.GetConfig()

	// OperationServers is empty in this SDK version, so the endpoint key is
	// unused and resolution falls back to the default server configuration,
	// honoring region and api hostname overrides
	basePath, err := cfg.ServerURLWithContext(ctx, "ApplicationsApiService.ReadOneApplication")
	if err != nil {
		return "", fmt.Errorf("cannot resolve the PingOne API base URL: %w", err)
	}

	return fmt.Sprintf(
		"%s/environments/%s/applications/%s/metadata",
		basePath,
		url.PathEscape(environmentID),
		url.PathEscape(applicationID),
	), nil
}

// applicationMetadataExecute performs the metadata API request via
// sdk.RetryWrapper, retrying transient statuses.  It returns the response with
// a re-readable body and the raw response body.
func applicationMetadataExecute(ctx context.Context, managementClient *management.APIClient, method, endpoint string, payload any) (*http.Response, []byte, error) {
	var (
		resp *http.Response
		body []byte
	)

	_, _, err := sdk.RetryWrapper(ctx, applicationMetadataTimeoutMinutes*time.Minute,
		func() (any, *http.Response, error) {
			r, b, rErr := applicationMetadataRequest(ctx, managementClient, method, endpoint, payload)
			resp, body = r, b

			if rErr != nil {
				return nil, r, rErr
			}

			return b, r, nil
		},
		applicationMetadataRetryable,
	)
	if err != nil {
		return resp, body, err
	}

	return resp, body, nil
}

// applicationMetadataRequest executes a single HTTP request against the
// application metadata endpoint and returns the response with a re-readable
// body.  A status >= 300 is returned as an error alongside the response.
func applicationMetadataRequest(ctx context.Context, managementClient *management.APIClient, method, endpoint string, payload any) (*http.Response, []byte, error) {
	cfg := managementClient.GetConfig()

	var bodyReader io.Reader
	if payload != nil {
		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot marshal application metadata request: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bodyReader)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", cfg.UserAgent)
	for key, value := range cfg.DefaultHeader {
		req.Header.Add(key, value)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	r, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}
	r.Body = io.NopCloser(bytes.NewReader(body))

	if r.StatusCode >= 300 {
		return r, body, fmt.Errorf("received unexpected response from the PingOne API: %s", r.Status)
	}

	return r, body, nil
}

// applicationMetadataError formats diagnostics for a failed metadata API call
// consistently with legacysdk.ParseResponse.
func applicationMetadataError(requestID string, resp *http.Response, err error) diag.Diagnostics {
	var diags diag.Diagnostics

	if resp == nil {
		diags.AddError(
			fmt.Sprintf("Error when calling `%s`: %v", requestID, err),
			"A generic error has occurred.",
		)
		return diags
	}

	// The response body is re-readable; attempt to parse a PingOne error object
	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr == nil && len(body) > 0 {
		var p1Error management.P1Error
		if err := json.Unmarshal(body, &p1Error); err == nil && p1Error.GetId() != "" {
			p1ErrorModel, err := model.RemarshalErrorObj(p1Error)
			if err == nil && p1ErrorModel != nil {
				summaryText, detailText := sdk.FormatPingOneError(requestID, *p1ErrorModel)
				diags.AddError(summaryText, detailText)
				return diags
			}
		}
	}

	diags.AddError(
		fmt.Sprintf("Error when calling `%s`: %v", requestID, err),
		fmt.Sprintf("A generic error has occurred.\nError details: %s", utils.ResponseErrorDetails(resp)),
	)

	return diags
}

// applicationMetadataRetryable retries transient HTTP statuses, mirroring the
// SDK's internal testForRetryable list.
var applicationMetadataRetryable = func(ctx context.Context, resp *http.Response, p1Error *model.P1Error) bool {
	if resp == nil {
		return false
	}

	switch resp.StatusCode {
	case http.StatusRequestTimeout,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}

	return false
}
