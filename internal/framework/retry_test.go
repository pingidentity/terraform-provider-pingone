// Copyright © 2026 Ping Identity Corporation

package framework_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/pingidentity/pingone-go-client/pingone"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/pingcli"
)

type apiError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (e apiError) Error() string { return e.Message }

// notAuthorizedCall returns an SDK call that fails with a "not authorized"
// error, authenticating every request with the given bearer token.
func notAuthorizedCall(calls *int, bearerToken string) framework.SDKInterfaceFunc {
	return func() (any, *http.Response, error) {
		*calls++

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.pingone.test", nil)
		req.Header.Set("Authorization", "Bearer "+bearerToken)

		return nil, &http.Response{StatusCode: http.StatusForbidden, Request: req}, apiError{
			ID:      "1",
			Message: "The request could not be completed. You do not have access to this resource.",
		}
	}
}

func TestRetryWrapper_RetriesNotAuthorizedForClientCredentials(t *testing.T) {
	// A short session window proves it does not apply to worker credentials.
	previous := pingcli.SessionAuthorizationRetryWindow
	pingcli.SessionAuthorizationRetryWindow = 300 * time.Millisecond

	t.Cleanup(func() { pingcli.SessionAuthorizationRetryWindow = previous })

	var calls int

	start := time.Now()

	_, _, err := framework.RetryWrapper(context.Background(), 3*time.Second, notAuthorizedCall(&calls, "worker-token"), "test", framework.DefaultCreateReadRetryable)
	if err == nil {
		t.Fatal("RetryWrapper returned nil error, want the not-authorized error")
	}

	if calls < 2 {
		t.Errorf("calls = %d, want the not-authorized response to be retried", calls)
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

	var calls int

	start := time.Now()

	_, _, err := framework.RetryWrapper(context.Background(), 30*time.Second, notAuthorizedCall(&calls, "pingcli-session-token"), "test", framework.DefaultCreateReadRetryable)
	if err == nil {
		t.Fatal("RetryWrapper returned nil error, want the not-authorized error")
	}

	if calls < 2 {
		t.Errorf("calls = %d, want the error to be retried within the window", calls)
	}

	if elapsed := time.Since(start); elapsed > 10*time.Second {
		t.Errorf("RetryWrapper took %s, want it to stop shortly after the %s window rather than the full timeout", elapsed, time.Second)
	}
}

func TestRetryWrapper_DoesNotRetryWhereCallSiteDoesNotForPingCLISession(t *testing.T) {
	useSessionToken(t, "pingcli-session-token", time.Minute)

	var calls int

	_, _, err := framework.RetryWrapper(context.Background(), 30*time.Second, notAuthorizedCall(&calls, "pingcli-session-token"), "test", framework.DefaultRetryable)
	if err == nil {
		t.Fatal("RetryWrapper returned nil error, want the not-authorized error")
	}

	if calls != 1 {
		t.Errorf("calls = %d, want exactly 1: the window must not add retries a call site never had", calls)
	}
}

func TestRetryWrapper_OtherRetryableErrorsStillRetryWhenSessionEnabled(t *testing.T) {
	pingcli.RegisterSessionToken("pingcli-session-token")
	t.Cleanup(pingcli.ResetSessionTokens)

	var calls int

	call := func() (any, *http.Response, error) {
		calls++

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.pingone.test", nil)
		req.Header.Set("Authorization", "Bearer pingcli-session-token")

		return nil, &http.Response{StatusCode: http.StatusServiceUnavailable, Request: req}, apiError{ID: "2", Message: "Service unavailable"}
	}

	always := func(context.Context, *http.Response, *pingone.GeneralError) bool { return true }

	_, _, err := framework.RetryWrapper(context.Background(), 3*time.Second, call, "test", always)
	if err == nil {
		t.Fatal("RetryWrapper returned nil error, want the final error")
	}

	if calls < 2 {
		t.Errorf("calls = %d, want other retryable errors to still be retried", calls)
	}
}

func TestRetryWrapper_ResponseWithoutRequestRetries(t *testing.T) {
	// A response without a request cannot be attributed to a credential, so
	// the not-authorized retry must apply.
	var calls int

	call := func() (any, *http.Response, error) {
		calls++
		return nil, &http.Response{StatusCode: http.StatusForbidden}, apiError{
			ID:      "1",
			Message: "The request could not be completed. You do not have access to this resource.",
		}
	}

	pingcli.RegisterSessionToken("pingcli-session-token")
	t.Cleanup(pingcli.ResetSessionTokens)

	_, _, err := framework.RetryWrapper(context.Background(), 3*time.Second, call, "test", framework.DefaultCreateReadRetryable)
	if err == nil {
		t.Fatal("RetryWrapper returned nil error, want the not-authorized error")
	}

	if calls < 2 {
		t.Errorf("calls = %d, want the not-authorized response to be retried", calls)
	}
}
