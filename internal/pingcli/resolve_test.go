// Copyright © 2026 Ping Identity Corporation

package pingcli_test

import (
	"strings"
	"testing"

	"github.com/pingidentity/terraform-provider-pingone/internal/pingcli"
)

func TestResolve(t *testing.T) {
	t.Parallel()

	const (
		envRegion = "NA"
		region    = "EU"
	)

	tests := []struct {
		name         string
		inputs       pingcli.Inputs
		env          map[string]string
		expectPlan   pingcli.Plan
		expectDetail []string
	}{
		{
			name:       "no pingcli configuration yields empty plan",
			inputs:     pingcli.Inputs{RegionCode: region, HasHCLCredential: true},
			expectPlan: pingcli.Plan{},
		},
		{
			name:       "no pingcli configuration and no credentials yields empty plan",
			inputs:     pingcli.Inputs{RegionCode: region},
			expectPlan: pingcli.Plan{},
		},
		{
			name:       "explicit config_profile enables",
			inputs:     pingcli.Inputs{ConfigProfile: "dev", RegionCode: region},
			expectPlan: pingcli.Plan{Enabled: true, ConfigProfile: "dev"},
		},
		{
			name:       "explicit values are trimmed",
			inputs:     pingcli.Inputs{ConfigProfile: "  dev  ", RegionCode: region},
			expectPlan: pingcli.Plan{Enabled: true, ConfigProfile: "dev"},
		},
		{
			name:   "env activates implicit mode",
			inputs: pingcli.Inputs{RegionCode: region},
			env:    map[string]string{"PINGCLI_PROFILE": "dev"},
			expectPlan: pingcli.Plan{
				Enabled:       true,
				ConfigProfile: "dev",
			},
		},
		{
			name:       "env activation suppressed by HCL credential",
			inputs:     pingcli.Inputs{RegionCode: region, HasHCLCredential: true},
			env:        map[string]string{"PINGCLI_PROFILE": "dev"},
			expectPlan: pingcli.Plan{},
		},
		{
			name:       "env activation suppressed by env credential",
			inputs:     pingcli.Inputs{RegionCode: region},
			env:        map[string]string{"PINGCLI_PROFILE": "dev", "PINGONE_CLIENT_ID": "abc"},
			expectPlan: pingcli.Plan{},
		},
		{
			name: "explicit with HCL credential conflicts",
			inputs: pingcli.Inputs{
				ConfigProfile:    "dev",
				RegionCode:       region,
				HasHCLCredential: true,
			},
			expectPlan: pingcli.Plan{Enabled: true, ConfigProfile: "dev"},
			expectDetail: []string{
				"Conflicting provider authentication configuration",
			},
		},
		{
			name: "explicit config_profile wins over env credentials",
			inputs: pingcli.Inputs{
				ConfigProfile: "dev",
				RegionCode:    region,
			},
			env:        map[string]string{"PINGONE_CLIENT_ID": "abc", "PINGONE_ENVIRONMENT_ID": "env", "PINGONE_API_ACCESS_TOKEN": "tok"},
			expectPlan: pingcli.Plan{Enabled: true, ConfigProfile: "dev"},
		},
		{
			name: "explicit without region code is an error",
			inputs: pingcli.Inputs{
				ConfigProfile: "dev",
				RegionCode:    "",
			},
			env: map[string]string{"PINGONE_REGION_CODE": ""},
			expectPlan: pingcli.Plan{
				Enabled:       true,
				ConfigProfile: "dev",
			},
			expectDetail: []string{"Missing region code"},
		},
		{
			name: "explicit with env region code is fine",
			inputs: pingcli.Inputs{
				ConfigProfile: "dev",
				RegionCode:    envRegion,
			},
			expectPlan: pingcli.Plan{
				Enabled:       true,
				ConfigProfile: "dev",
			},
		},
		{
			name: "whitespace-only env values are treated as unset",
			inputs: pingcli.Inputs{
				RegionCode: region,
			},
			env:        map[string]string{"PINGCLI_PROFILE": "   "},
			expectPlan: pingcli.Plan{},
		},
		{
			name:       "explicit config_grant_type with env profile enables",
			inputs:     pingcli.Inputs{ConfigGrantType: "device_code", RegionCode: region},
			env:        map[string]string{"PINGCLI_PROFILE": "dev"},
			expectPlan: pingcli.Plan{Enabled: true, ConfigProfile: "dev", ConfigGrantType: "device_code"},
		},
		{
			name:       "explicit config_grant_type without any profile is an error",
			inputs:     pingcli.Inputs{ConfigGrantType: "device_code", RegionCode: region},
			expectPlan: pingcli.Plan{Enabled: true, ConfigGrantType: "device_code"},
			expectDetail: []string{
				"Missing Ping CLI profile",
			},
		},
		{
			name:       "explicit profile and grant type are kept",
			inputs:     pingcli.Inputs{ConfigProfile: "dev", ConfigGrantType: "authorization_code", RegionCode: region},
			expectPlan: pingcli.Plan{Enabled: true, ConfigProfile: "dev", ConfigGrantType: "authorization_code"},
		},
		{
			name:       "explicit grant type is trimmed",
			inputs:     pingcli.Inputs{ConfigProfile: "dev", ConfigGrantType: " client_credentials ", RegionCode: region},
			expectPlan: pingcli.Plan{Enabled: true, ConfigProfile: "dev", ConfigGrantType: "client_credentials"},
		},
		{
			name:       "env grant type applies to an explicit profile",
			inputs:     pingcli.Inputs{ConfigProfile: "dev", RegionCode: region},
			env:        map[string]string{"PINGCLI_PINGONE_AUTHENTICATION_TYPE": "device_code"},
			expectPlan: pingcli.Plan{Enabled: true, ConfigProfile: "dev", ConfigGrantType: "device_code"},
		},
		{
			name:       "explicit grant type wins over env grant type",
			inputs:     pingcli.Inputs{ConfigProfile: "dev", ConfigGrantType: "client_credentials", RegionCode: region},
			env:        map[string]string{"PINGCLI_PINGONE_AUTHENTICATION_TYPE": "device_code"},
			expectPlan: pingcli.Plan{Enabled: true, ConfigProfile: "dev", ConfigGrantType: "client_credentials"},
		},
		{
			name:       "env profile and grant type activate implicit mode",
			inputs:     pingcli.Inputs{RegionCode: region},
			env:        map[string]string{"PINGCLI_PROFILE": "dev", "PINGCLI_PINGONE_AUTHENTICATION_TYPE": "authorization_code"},
			expectPlan: pingcli.Plan{Enabled: true, ConfigProfile: "dev", ConfigGrantType: "authorization_code"},
		},
		{
			name:       "env grant type alone does not enable",
			inputs:     pingcli.Inputs{RegionCode: region},
			env:        map[string]string{"PINGCLI_PINGONE_AUTHENTICATION_TYPE": "client_credentials"},
			expectPlan: pingcli.Plan{},
		},
		{
			name:       "invalid grant type is an error",
			inputs:     pingcli.Inputs{ConfigProfile: "dev", ConfigGrantType: "password", RegionCode: region},
			expectPlan: pingcli.Plan{Enabled: true, ConfigProfile: "dev", ConfigGrantType: "password"},
			expectDetail: []string{
				"Invalid Ping CLI grant type",
			},
		},
		{
			name:       "grant type with HCL credential conflicts",
			inputs:     pingcli.Inputs{ConfigGrantType: "device_code", RegionCode: region, HasHCLCredential: true},
			env:        map[string]string{"PINGCLI_PROFILE": "dev"},
			expectPlan: pingcli.Plan{Enabled: true, ConfigProfile: "dev", ConfigGrantType: "device_code"},
			expectDetail: []string{
				"Conflicting provider authentication configuration",
			},
		},
		{
			name: "PINGCLI_CONFIG alone does not enable",
			inputs: pingcli.Inputs{
				RegionCode: region,
			},
			env:        map[string]string{"PINGCLI_CONFIG": "/tmp/config.yaml"},
			expectPlan: pingcli.Plan{},
		},
		{
			name: "explicit config_profile wins over environment profile",
			inputs: pingcli.Inputs{
				ConfigProfile: "dev",
				RegionCode:    region,
			},
			env:        map[string]string{"PINGCLI_PROFILE": "prod"},
			expectPlan: pingcli.Plan{Enabled: true, ConfigProfile: "dev"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			getenv := func(name string) string {
				if tt.env == nil {
					return ""
				}
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
				if !contains(combined, want) {
					t.Errorf("problem %d = %q, want it to contain %q", i, combined, want)
				}
			}
		})
	}
}

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}
