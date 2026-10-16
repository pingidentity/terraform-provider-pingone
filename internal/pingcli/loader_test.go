// Copyright © 2026 Ping Identity Corporation

package pingcli_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pingidentity/terraform-provider-pingone/internal/pingcli"
)

func TestAuthenticate_ConcurrentCallsShareOneRead(t *testing.T) {
	var reads atomic.Int32

	restore := pingcli.SetLoadForTest(func(_ context.Context, plan pingcli.Plan) (pingcli.SessionResult, error) {
		reads.Add(1)

		// A slow read widens the window in which a second caller could also miss
		// the memo.
		time.Sleep(50 * time.Millisecond)

		return pingcli.SessionResult{AccessToken: "token", Expiry: time.Now().Add(time.Hour), Profile: plan.ConfigProfile}, nil
	})
	t.Cleanup(restore)

	plan := pingcli.Plan{Enabled: true, ConfigProfile: "default"}

	var wg sync.WaitGroup

	for range 20 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			result, err := pingcli.Authenticate(context.Background(), plan)
			if err != nil || result.AccessToken != "token" {
				t.Errorf("Authenticate = (%q, %v), want the shared session", result.AccessToken, err)
			}
		}()
	}

	wg.Wait()

	if got := reads.Load(); got != 1 {
		t.Errorf("keychain reads = %d, want exactly 1 for concurrent callers", got)
	}
}

func TestAuthenticate_DoesNotCacheFailures(t *testing.T) {
	var reads atomic.Int32

	restore := pingcli.SetLoadForTest(func(context.Context, pingcli.Plan) (pingcli.SessionResult, error) {
		reads.Add(1)

		return pingcli.SessionResult{}, context.DeadlineExceeded
	})
	t.Cleanup(restore)

	plan := pingcli.Plan{Enabled: true, ConfigProfile: "default"}

	for range 2 {
		if _, err := pingcli.Authenticate(context.Background(), plan); err == nil {
			t.Fatal("Authenticate returned nil error, want the load failure")
		}
	}

	if got := reads.Load(); got != 2 {
		t.Errorf("keychain reads = %d, want a failed load to be retried", got)
	}
}
