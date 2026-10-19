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
			name:       "environment variable enables when the attribute is unset",
			inputs:     pingcli.Inputs{RegionCode: "EU"},
			env:        map[string]string{"PINGONE_TF_USE_PINGCLI_AUTH": "true"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "default"},
		},
		{
			name:       "environment variable false leaves it off",
			inputs:     pingcli.Inputs{RegionCode: "EU"},
			env:        map[string]string{"PINGONE_TF_USE_PINGCLI_AUTH": "false"},
			expectPlan: pingcli.Plan{},
		},
		{
			name:       "explicit false beats the environment variable",
			inputs:     pingcli.Inputs{UsePingCLIAuth: boolPtr(false), RegionCode: "EU"},
			env:        map[string]string{"PINGONE_TF_USE_PINGCLI_AUTH": "true"},
			expectPlan: pingcli.Plan{},
		},
		{
			name:       "explicit true beats the environment variable",
			inputs:     pingcli.Inputs{UsePingCLIAuth: boolPtr(true), RegionCode: "EU"},
			env:        map[string]string{"PINGONE_TF_USE_PINGCLI_AUTH": "false"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "default"},
		},
		{
			name:         "invalid environment variable is an error",
			inputs:       pingcli.Inputs{RegionCode: "EU"},
			env:          map[string]string{"PINGONE_TF_USE_PINGCLI_AUTH": "maybe"},
			expectPlan:   pingcli.Plan{},
			expectDetail: []string{"PINGONE_TF_USE_PINGCLI_AUTH environment variable must be a boolean"},
		},
		{
			name:       "boolean enables with defaults",
			inputs:     pingcli.Inputs{UsePingCLIAuth: boolPtr(true), RegionCode: "EU"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "default"},
		},
		{
			name:       "profile and grant type are used",
			inputs:     pingcli.Inputs{UsePingCLIAuth: boolPtr(true), PingCLIProfile: "dev", PingCLIGrantType: "device_code", RegionCode: "EU"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "dev", PingCLIGrantType: "device_code"},
		},
		{
			name:       "environment variables supply profile and grant type",
			inputs:     pingcli.Inputs{UsePingCLIAuth: boolPtr(true), RegionCode: "EU"},
			env:        map[string]string{"PINGCLI_PROFILE": "dev", "PINGCLI_PINGONE_AUTHENTICATION_TYPE": "authorization_code"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "dev", PingCLIGrantType: "authorization_code"},
		},
		{
			name:       "HCL values win over environment variables",
			inputs:     pingcli.Inputs{UsePingCLIAuth: boolPtr(true), PingCLIProfile: "hcl", PingCLIGrantType: "device_code", RegionCode: "EU"},
			env:        map[string]string{"PINGCLI_PROFILE": "env", "PINGCLI_PINGONE_AUTHENTICATION_TYPE": "authorization_code"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "hcl", PingCLIGrantType: "device_code"},
		},
		{
			name:       "values are trimmed",
			inputs:     pingcli.Inputs{UsePingCLIAuth: boolPtr(true), PingCLIProfile: "  dev  ", RegionCode: "EU"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "dev"},
		},
		{
			name:         "HCL credentials conflict",
			inputs:       pingcli.Inputs{UsePingCLIAuth: boolPtr(true), RegionCode: "EU", HasHCLCredential: true},
			expectPlan:   pingcli.Plan{Enabled: true, PingCLIProfile: "default"},
			expectDetail: []string{"Conflicting provider authentication configuration"},
		},
		{
			name:         "exported credential environment variables conflict",
			inputs:       pingcli.Inputs{UsePingCLIAuth: boolPtr(true), RegionCode: "EU"},
			env:          map[string]string{"PINGONE_CLIENT_ID": "abc", "PINGONE_ENVIRONMENT_ID": "env", "PINGONE_API_ACCESS_TOKEN": "tok"},
			expectPlan:   pingcli.Plan{Enabled: true, PingCLIProfile: "default"},
			expectDetail: []string{"PINGONE_CLIENT_ID, PINGONE_ENVIRONMENT_ID, PINGONE_API_ACCESS_TOKEN"},
		},
		{
			name:       "exported grant type and storage variables are not credentials",
			inputs:     pingcli.Inputs{UsePingCLIAuth: boolPtr(true), RegionCode: "EU"},
			env:        map[string]string{"PINGONE_AUTH_GRANT_TYPE": "client_credentials", "PINGONE_STORAGE_TYPE": "NONE"},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "default"},
		},
		{
			name:         "missing region code is an error",
			inputs:       pingcli.Inputs{UsePingCLIAuth: boolPtr(true)},
			expectPlan:   pingcli.Plan{Enabled: true, PingCLIProfile: "default"},
			expectDetail: []string{"Missing region code"},
		},
		{
			name:       "region code from the environment is fine",
			inputs:     pingcli.Inputs{UsePingCLIAuth: boolPtr(true), RegionCode: "NA"},
			env:        map[string]string{"PINGONE_REGION_CODE": ""},
			expectPlan: pingcli.Plan{Enabled: true, PingCLIProfile: "default"},
		},
		{
			name:         "invalid grant type is an error",
			inputs:       pingcli.Inputs{UsePingCLIAuth: boolPtr(true), PingCLIGrantType: "password", RegionCode: "EU"},
			expectPlan:   pingcli.Plan{Enabled: true, PingCLIProfile: "default", PingCLIGrantType: "password"},
			expectDetail: []string{"Invalid Ping CLI grant type"},
		},
		{
			name:       "invalid grant type from the environment is an error",
			inputs:     pingcli.Inputs{UsePingCLIAuth: boolPtr(true), RegionCode: "EU"},
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

func boolPtr(b bool) *bool { return &b }
