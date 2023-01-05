package sso_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/pingidentity/terraform-provider-pingone/internal/service/sso"
)

func TestSignOnPolicyActionStateUpgradeV0(t *testing.T) {

	tests := []struct {
		Objective     string
		TestState     map[string]interface{}
		ExpectedState map[string]interface{}
	}{
		{
			Objective:     "no existing state",
			TestState:     nil,
			ExpectedState: nil,
		},

		{
			Objective: "no conditions",
			TestState: map[string]interface{}{
				"enforce_lockout_for_identity_providers": false,
				"environment_id":                         "environmentID1234",
				"id":                                     "ID1234",
				"login": []map[string]interface{}{
					{
						"recovery_enabled": true,
					},
				},
				"priority":                             1,
				"registration_confirm_user_attributes": false,
				"registration_external_href":           "",
				"registration_local_population_id":     "",
				"sign_on_policy_id":                    "policyID1234",
			},
			ExpectedState: map[string]interface{}{
				"enforce_lockout_for_identity_providers": false,
				"environment_id":                         "environmentID1234",
				"id":                                     "ID1234",
				"login": []map[string]interface{}{
					{
						"recovery_enabled": true,
					},
				},
				"priority":                             1,
				"registration_confirm_user_attributes": false,
				"registration_external_href":           "",
				"registration_local_population_id":     "",
				"sign_on_policy_id":                    "policyID1234",
			},
		},

		{
			Objective: "has conditions but no user_attribute_equals maps",
			TestState: map[string]interface{}{
				"conditions": []map[string]interface{}{
					{
						"anonymous_network_detected":      false,
						"geovelocity_anomaly_detected":    false,
						"ip_reputation_high_risk":         false,
						"last_sign_on_older_than_seconds": 86400,
					},
				},
				"enforce_lockout_for_identity_providers": false,
				"environment_id":                         "environmentID1234",
				"id":                                     "ID1234",
				"login": []map[string]interface{}{
					{
						"recovery_enabled": true,
					},
				},
				"priority":                             1,
				"registration_confirm_user_attributes": false,
				"registration_external_href":           "",
				"registration_local_population_id":     "",
				"sign_on_policy_id":                    "policyID1234",
			},
			ExpectedState: map[string]interface{}{
				"conditions": []map[string]interface{}{
					{
						"anonymous_network_detected":      false,
						"geovelocity_anomaly_detected":    false,
						"ip_reputation_high_risk":         false,
						"last_sign_on_older_than_seconds": 86400,
					},
				},
				"enforce_lockout_for_identity_providers": false,
				"environment_id":                         "environmentID1234",
				"id":                                     "ID1234",
				"login": []map[string]interface{}{
					{
						"recovery_enabled": true,
					},
				},
				"priority":                             1,
				"registration_confirm_user_attributes": false,
				"registration_external_href":           "",
				"registration_local_population_id":     "",
				"sign_on_policy_id":                    "policyID1234",
			},
		},

		{
			Objective: "has conditions with single user_attribute_equals map",
			TestState: map[string]interface{}{
				"conditions": []map[string]interface{}{
					{
						"anonymous_network_detected":      false,
						"geovelocity_anomaly_detected":    false,
						"ip_reputation_high_risk":         false,
						"last_sign_on_older_than_seconds": 86400,
						"user_attribute_equals": []map[string]interface{}{
							{
								"attribute_reference": "${user.lifecycle.status}",
								"value":               "VERIFICATION_REQUIRED",
							},
						},
					},
				},
				"enforce_lockout_for_identity_providers": false,
				"environment_id":                         "environmentID1234",
				"id":                                     "ID1234",
				"login": []map[string]interface{}{
					{
						"recovery_enabled": true,
					},
				},
				"priority":                             1,
				"registration_confirm_user_attributes": false,
				"registration_external_href":           "",
				"registration_local_population_id":     "",
				"sign_on_policy_id":                    "policyID1234",
			},
			ExpectedState: map[string]interface{}{
				"conditions": []map[string]interface{}{
					{
						"anonymous_network_detected":      false,
						"geovelocity_anomaly_detected":    false,
						"ip_reputation_high_risk":         false,
						"last_sign_on_older_than_seconds": 86400,
						"user_attribute_equals": []map[string]interface{}{
							{
								"attribute_reference": "${user.lifecycle.status}",
								"value_string":        "VERIFICATION_REQUIRED",
							},
						},
					},
				},
				"enforce_lockout_for_identity_providers": false,
				"environment_id":                         "environmentID1234",
				"id":                                     "ID1234",
				"login": []map[string]interface{}{
					{
						"recovery_enabled": true,
					},
				},
				"priority":                             1,
				"registration_confirm_user_attributes": false,
				"registration_external_href":           "",
				"registration_local_population_id":     "",
				"sign_on_policy_id":                    "policyID1234",
			},
		},

		{
			Objective: "has conditions with multiple user_attribute_equals maps",
			TestState: map[string]interface{}{
				"conditions": []map[string]interface{}{
					{
						"anonymous_network_detected":      false,
						"geovelocity_anomaly_detected":    false,
						"ip_reputation_high_risk":         false,
						"last_sign_on_older_than_seconds": 86400,
						"user_attribute_equals": []map[string]interface{}{
							{
								"attribute_reference": "${user.lifecycle.status}",
								"value":               "VERIFICATION_REQUIRED",
							},
							{
								"attribute_reference": "${user.name.family}",
								"value":               "Wayne",
							},
						},
					},
				},
				"enforce_lockout_for_identity_providers": false,
				"environment_id":                         "environmentID1234",
				"id":                                     "ID1234",
				"login": []map[string]interface{}{
					{
						"recovery_enabled": true,
					},
				},
				"priority":                             1,
				"registration_confirm_user_attributes": false,
				"registration_external_href":           "",
				"registration_local_population_id":     "",
				"sign_on_policy_id":                    "policyID1234",
			},
			ExpectedState: map[string]interface{}{
				"conditions": []map[string]interface{}{
					{
						"anonymous_network_detected":      false,
						"geovelocity_anomaly_detected":    false,
						"ip_reputation_high_risk":         false,
						"last_sign_on_older_than_seconds": 86400,
						"user_attribute_equals": []map[string]interface{}{
							{
								"attribute_reference": "${user.lifecycle.status}",
								"value_string":        "VERIFICATION_REQUIRED",
							},
							{
								"attribute_reference": "${user.name.family}",
								"value_string":        "Wayne",
							},
						},
					},
				},
				"enforce_lockout_for_identity_providers": false,
				"environment_id":                         "environmentID1234",
				"id":                                     "ID1234",
				"login": []map[string]interface{}{
					{
						"recovery_enabled": true,
					},
				},
				"priority":                             1,
				"registration_confirm_user_attributes": false,
				"registration_external_href":           "",
				"registration_local_population_id":     "",
				"sign_on_policy_id":                    "policyID1234",
			},
		},
	}

	for _, test := range tests {

		t.Run(test.Objective, func(t *testing.T) {
			returnedState, err := sso.ResourceSignOnPolicyActionStateUpgradeV0(context.Background(), test.TestState, nil)

			if err != nil {
				t.Fatalf("State cannot be migrated: %s", err)
			}

			if !reflect.DeepEqual(test.ExpectedState, returnedState) {
				t.Fatalf("\nExpected: \t%#v\ngot:\t\t%#v", test.ExpectedState, returnedState)
			}
		})
	}
}
