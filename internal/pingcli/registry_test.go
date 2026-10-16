// Copyright © 2026 Ping Identity Corporation

package pingcli_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/pingidentity/terraform-provider-pingone/internal/pingcli"
	"github.com/zalando/go-keyring"
)

func TestRegisterSessionToken(t *testing.T) {
	pingcli.ResetSessionTokens()
	t.Cleanup(pingcli.ResetSessionTokens)

	if pingcli.IsSessionToken("token-a") {
		t.Error("unregistered token reported as a session token")
	}

	pingcli.RegisterSessionToken("token-a")

	if !pingcli.IsSessionToken("token-a") {
		t.Error("registered token not reported as a session token")
	}

	if pingcli.IsSessionToken("token-b") {
		t.Error("unregistered token reported as a session token")
	}

	if pingcli.IsSessionToken("") {
		t.Error("empty token reported as a session token")
	}
}

func TestResetSessionTokens(t *testing.T) {
	pingcli.ResetSessionTokens()
	t.Cleanup(pingcli.ResetSessionTokens)

	pingcli.RegisterSessionToken("token-a")
	pingcli.ResetSessionTokens()

	if pingcli.IsSessionToken("token-a") {
		t.Error("token still registered after ResetSessionTokens")
	}
}

func TestRegisterSessionToken_ConcurrentUse(t *testing.T) {
	pingcli.ResetSessionTokens()
	t.Cleanup(pingcli.ResetSessionTokens)

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			pingcli.RegisterSessionToken("concurrent-token")
			if !pingcli.IsSessionToken("concurrent-token") {
				t.Error("registered token not reported as a session token")
			}
		}()
	}
	wg.Wait()
}

func TestLoad_RegistersLoadedSessionToken(t *testing.T) {
	keyring.MockInit()

	account := sessionAccount(t, pingcli.GrantTypeClientCredentials, "default")
	seedKeychainSession(t, account, "registry-token", time.Now().Add(time.Hour))

	pingcli.ResetSessionTokens()
	t.Cleanup(pingcli.ResetSessionTokens)

	if _, err := pingcli.Load(context.Background(), pingcli.Plan{ConfigProfile: "default"}); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if !pingcli.IsSessionToken("registry-token") {
		t.Error("loaded session token was not registered")
	}
}
