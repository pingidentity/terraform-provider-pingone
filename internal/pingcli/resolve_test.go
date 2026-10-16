// Copyright © 2026 Ping Identity Corporation

package pingcli_test

import (
	"strings"
	"testing"

	"github.com/pingidentity/terraform-provider-pingone/internal/pingcli"
)

func TestResolve(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		inputs       pingcli.Inputs
		env          map[string]string
		expectPlan   pingcli.Plan
		expectDetail []string
	}{
		{
			name:       "boolean unset yields empty plan, even with profile attributes set",
			inputs:     pingcli.Inputs{PingCLIProfile: "dev", PingCLIGrantType: "device_code", RegionCode: "EU"},
			expectPlan: pingcli.Plan{},
		},
		{
			name:       "boolean enables with defaults",
			inputs:     pingcli.Inputs{UsePingCLIAuth: true, RegionCode: "EU"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "default"},
		},
		{
			name:       "profile and grant type are used",
			inputs:     pingcli.Inputs{UsePingCLIAuth: true, PingCLIProfile: "dev", PingCLIGrantType: "device_code", RegionCode: "EU"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "dev", PingCLIGrantType: "device_code"},
		},
		{
			name:       "environment variables supply profile and grant type",
			inputs:     pingcli.Inputs{UsePingCLIAuth: true, RegionCode: "EU"},
			env:        map[string]string{"PINGCLI_PROFILE": "dev", "PINGCLI_PINGONE_AUTHENTICATION_TYPE": "authorization_code"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "dev", PingCLIGrantType: "authorization_code"},
		},
		{
			name:       "HCL values win over environment variables",
			inputs:     pingcli.Inputs{UsePingCLIAuth: true, PingCLIProfile: "hcl", PingCLIGrantType: "device_code", RegionCode: "EU"},
			env:        map[string]string{"PINGCLI_PROFILE": "env", "PINGCLI_PINGONE_AUTHENTICATION_TYPE": "authorization_code"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "hcl", PingCLIGrantType: "device_code"},
		},
		{
			name:       "values are trimmed",
			inputs:     pingcli.Inputs{UsePingCLIAuth: true, PingCLIProfile: "  dev  ", RegionCode: "EU"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "dev"},
		},
		{
			name:         "HCL credentials conflict",
			inputs:       pingcli.Inputs{UsePingCLIAuth: true, RegionCode: "EU", HasHCLCredential: true},
			expectPlan:   pingcli.Plan{Enabled: true, PingCLIProfile: "default"},
			expectDetail: []string{"Conflicting provider authentication configuration"},
		},
		{
			name:       "exported credential environment variables do not conflict",
			inputs:     pingcli.Inputs{UsePingCLIAuth: true, RegionCode: "EU"},
			env:        map[string]string{"PINGONE_CLIENT_ID": "abc", "PINGONE_ENVIRONMENT_ID": "env", "PINGONE_API_ACCESS_TOKEN": "tok"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "default"},
		},
		{
			name:         "missing region code is an error",
			inputs:       pingcli.Inputs{UsePingCLIAuth: true},
			expectPlan:   pingcli.Plan{Enabled: true, PingCLIProfile: "default"},
			expectDetail: []string{"Missing region code"},
		},
		{
			name:       "region code from the environment is fine",
			inputs:     pingcli.Inputs{UsePingCLIAuth: true, RegionCode: "NA"},
			env:        map[string]string{"PINGONE_REGION_CODE": ""},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "default"},
		},
		{
			name:         "invalid grant type is an error",
			inputs:       pingcli.Inputs{UsePingCLIAuth: true, PingCLIGrantType: "password", RegionCode: "EU"},
			expectPlan:   pingcli.Plan{Enabled: true, PingCLIProfile: "default", PingCLIGrantType: "password"},
			expectDetail: []string{"Invalid Ping CLI grant type"},
		},
		{
			name:       "invalid grant type from the environment is an error",
			inputs:     pingcli.Inputs{UsePingCLIAuth: true, RegionCode: "EU"},
			env:        map[string]string{"PINGCLI_PINGONE_AUTHENTICATION_TYPE": "password"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "default", PingCLIGrantType: "password"},
			expectDetail: []string{
				"Invalid Ping CLI grant type",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			getenv := func(name string) string {
				return tt.env[name]
			}

			plan, problems := pingcli.Resolve(tt.inputs, getenv)

			if plan != tt.expectPlan {
				t.Errorf("plan = %+v, want %+v", plan, tt.expectPlan)
			}

			if len(problems) != len(tt.expectDetail) {
				t.Fatalf("got %d problems, want %d: %+v", len(problems), len(tt.expectDetail), problems)
			}

			for i, want := range tt.expectDetail {
				combined := problems[i].Summary + "\n" + problems[i].Detail
				if !strings.Contains(combined, want) {
					t.Errorf("problem %d = %q, want it to contain %q", i, combined, want)
				}
			}
		})
	}
}
