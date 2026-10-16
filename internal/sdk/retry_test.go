// Copyright © 2026 Ping Identity Corporation

package sdk_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/pingcli"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

func notAuthorizedServer(t *testing.T, requests *atomic.Int32, bearerToken string) *management.APIClient {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"id":"1","code":"ACCESS_FAILED","message":"The actor attempting to perform the request is not authorized."}`))
	}))
	t.Cleanup(server.Close)

	cfg := management.NewConfiguration()
	cfg.Servers = management.ServerConfigurations{{URL: server.URL}}
	if bearerToken != "" {
		// The generated client adds this header to every request, standing in
		// for the oauth2 transport's Authorization header in production.
		cfg.DefaultHeader["Authorization"] = "Bearer " + bearerToken
	}

	return management.NewAPIClient(cfg)
}

func readEnvironment(client *management.APIClient) func() (any, *http.Response, error) {
	return func() (any, *http.Response, error) {
		return client.EnvironmentsApi.ReadOneEnvironment(context.Background(), "00000000-0000-0000-0000-000000000000").Execute()
	}
}

func TestRetryWrapper_RetriesNotAuthorizedForClientCredentials(t *testing.T) {
	// A short session window proves it does not apply to worker credentials.
	previous := pingcli.SessionAuthorizationRetryWindow
	pingcli.SessionAuthorizationRetryWindow = 300 * time.Millisecond

	t.Cleanup(func() { pingcli.SessionAuthorizationRetryWindow = previous })

	var requests atomic.Int32
	client := notAuthorizedServer(t, &requests, "worker-token")

	start := time.Now()

	_, _, err := sdk.RetryWrapper(context.Background(), 3*time.Second, readEnvironment(client), sdk.DefaultCreateReadRetryable)
	if err == nil {
		t.Fatal("RetryWrapper returned nil error, want the not-authorized error")
	}

	if got := requests.Load(); got < 2 {
		t.Errorf("requests = %d, want the not-authorized response to be retried", got)
	}

	if elapsed := time.Since(start); elapsed < 2*time.Second {
		t.Errorf("RetryWrapper stopped after %s, want worker credentials retried for the full timeout", elapsed)
	}
}

func useSessionToken(t *testing.T, token string, window time.Duration) {
	t.Helper()

	// The token is registered exactly as a loaded Ping CLI session registers it.
	pingcli.RegisterSessionToken(token)

	previous := pingcli.SessionAuthorizationRetryWindow
	pingcli.SessionAuthorizationRetryWindow = window

	t.Cleanup(func() {
		pingcli.SessionAuthorizationRetryWindow = previous
		pingcli.ResetSessionTokens()
	})
}

func TestRetryWrapper_BoundsNotAuthorizedRetryForPingCLISession(t *testing.T) {
	useSessionToken(t, "pingcli-session-token", time.Second)

	var requests atomic.Int32
	client := notAuthorizedServer(t, &requests, "pingcli-session-token")

	start := time.Now()

	_, _, err := sdk.RetryWrapper(context.Background(), 30*time.Second, readEnvironment(client), sdk.DefaultCreateReadRetryable)
	if err == nil {
		t.Fatal("RetryWrapper returned nil error, want the not-authorized error")
	}

	if got := requests.Load(); got < 2 {
		t.Errorf("requests = %d, want the error to be retried within the window", got)
	}

	if elapsed := time.Since(start); elapsed > 10*time.Second {
		t.Errorf("RetryWrapper took %s, want it to stop shortly after the %s window rather than the full timeout", elapsed, time.Second)
	}
}

func TestRetryWrapper_DoesNotRetryWhereCallSiteDoesNotForPingCLISession(t *testing.T) {
	useSessionToken(t, "pingcli-session-token", time.Minute)

	var requests atomic.Int32
	client := notAuthorizedServer(t, &requests, "pingcli-session-token")

	_, _, err := sdk.RetryWrapper(context.Background(), 30*time.Second, readEnvironment(client), sdk.DefaultRetryable)
	if err == nil {
		t.Fatal("RetryWrapper returned nil error, want the not-authorized error")
	}

	if got := requests.Load(); got != 1 {
		t.Errorf("requests = %d, want exactly 1: the window must not add retries a call site never had", got)
	}
}
