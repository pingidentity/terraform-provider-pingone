// Copyright © 2026 Ping Identity Corporation

package acctest

import "testing"

func TestPingCLIAuthSkipReason(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		env      map[string]string
		wantSkip bool
	}{
		{name: "gate unset skips", env: map[string]string{}, wantSkip: true},
		{name: "gate false skips", env: map[string]string{"TESTACC_PINGCLI_AUTH": "false"}, wantSkip: true},
		{name: "gate set locally runs", env: map[string]string{"TESTACC_PINGCLI_AUTH": "true"}, wantSkip: false},
		{name: "CI skips even with the gate set", env: map[string]string{"TESTACC_PINGCLI_AUTH": "true", "CI": "true"}, wantSkip: true},
		{name: "any non-empty CI value skips", env: map[string]string{"TESTACC_PINGCLI_AUTH": "true", "CI": "1"}, wantSkip: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reason := pingCLIAuthSkipReason(func(name string) string { return tt.env[name] })

			if skipped := reason != ""; skipped != tt.wantSkip {
				t.Errorf("skip reason = %q, want skip = %v", reason, tt.wantSkip)
			}
		})
	}
}
