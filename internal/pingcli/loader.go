// Copyright © 2026 Ping Identity Corporation

package pingcli

import (
	"context"
	"sync"
	"time"
)

// loadSession is Load; it is a variable only so tests can count keychain reads.
var loadSession = Load

type memoKey struct {
	configProfile   string
	configGrantType string
}

type memoEntry struct {
	result    SessionResult
	expiresAt time.Time
}

var (
	memoMu sync.Mutex
	memo   = map[memoKey]memoEntry{}
)

// Authenticate returns the Ping CLI session for the plan, memoised until shortly
// before expiry so the muxed provider flavours share one keychain read.
func Authenticate(ctx context.Context, plan Plan) (SessionResult, error) {
	if !plan.Enabled {
		return SessionResult{}, nil
	}

	key := memoKey{configProfile: plan.ConfigProfile, configGrantType: plan.ConfigGrantType}

	// Hold the lock across the load so concurrent Configure calls share one
	// keychain read (and one OS prompt).
	memoMu.Lock()
	defer memoMu.Unlock()

	if entry, ok := memo[key]; ok && time.Now().Before(entry.expiresAt) {
		return entry.result, nil
	}

	result, err := loadSession(ctx, plan)
	if err != nil {
		// Do not cache failures: a user may log in and re-run in the same
		// process (for example in tests or long-lived tooling).
		delete(memo, key)
		return SessionResult{}, err
	}

	expiresAt := result.Expiry.Add(-expirySkew)
	if result.Expiry.IsZero() {
		// Unknown expiry: memoise briefly so one run shares one read, without
		// pinning a possibly short-lived token for long.
		expiresAt = time.Now().Add(time.Minute)
	}

	memo[key] = memoEntry{result: result, expiresAt: expiresAt}

	return result, nil
}

// ResetMemo clears any memoised sessions. Exposed for tests that seed state
// between scenarios.
func ResetMemo() {
	memoMu.Lock()
	defer memoMu.Unlock()
	memo = map[memoKey]memoEntry{}
}
