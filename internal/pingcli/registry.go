// Copyright © 2026 Ping Identity Corporation

package pingcli

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// SessionAuthorizationRetryWindow is how long a "not authorized" error is retried
// for a Ping CLI session token. Worker roles can take a while to propagate, so
// they retry for the full timeout; a user session only needs a short window for
// resources created during the run.
var SessionAuthorizationRetryWindow = 30 * time.Second

// sessionTokens holds hashes of the loaded Ping CLI session tokens, so the retry
// wrappers can tell session requests from worker credentials. Hashed, so no
// plaintext token is retained.
var sessionTokens = struct {
	sync.Mutex
	hashes map[string]struct{}
}{hashes: map[string]struct{}{}}

// RegisterSessionToken records the access token of a loaded Ping CLI session,
// so requests authenticated with it are recognised by the retry wrappers.
func RegisterSessionToken(accessToken string) {
	sessionTokens.Lock()
	defer sessionTokens.Unlock()

	sessionTokens.hashes[accessTokenHash(accessToken)] = struct{}{}
}

// IsSessionToken reports whether accessToken is the access token of a Ping CLI
// session loaded in this process.
func IsSessionToken(accessToken string) bool {
	if accessToken == "" {
		return false
	}

	sessionTokens.Lock()
	defer sessionTokens.Unlock()

	_, ok := sessionTokens.hashes[accessTokenHash(accessToken)]
	return ok
}

// ResetSessionTokens forgets every registered session token. Exposed for tests
// that seed state between scenarios.
func ResetSessionTokens() {
	sessionTokens.Lock()
	defer sessionTokens.Unlock()

	sessionTokens.hashes = map[string]struct{}{}
}

func accessTokenHash(accessToken string) string {
	sum := sha256.Sum256([]byte(accessToken))
	return hex.EncodeToString(sum[:])
}
