// Copyright © 2026 Ping Identity Corporation

package pingcli

import "context"

// SetLoadForTest replaces the session loader used by Authenticate and clears
// the memo. The returned function restores the production loader.
func SetLoadForTest(f func(ctx context.Context, plan Plan) (SessionResult, error)) (restore func()) {
	ResetMemo()

	previous := loadSession
	loadSession = f

	return func() {
		loadSession = previous

		ResetMemo()
	}
}
