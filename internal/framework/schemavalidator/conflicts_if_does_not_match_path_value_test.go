// Copyright © 2026 Ping Identity Corporation

package schemavalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/schemavalidator"
)

func TestConflictsIfDoesNotMatchPathValueValidator(t *testing.T) {
	t.Parallel()

	// A minimal schema with a controlling "action" attribute and a guarded
	// "guarded_attr" attribute. The validator permits guarded_attr only when
	// action == "MATCH".
	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"action":       schema.StringAttribute{Optional: true},
			"guarded_attr": schema.StringAttribute{Optional: true},
		},
	}

	objType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"action":       tftypes.String,
			"guarded_attr": tftypes.String,
		},
	}

	newConfig := func(action, guarded *string) tfsdk.Config {
		actionVal := tftypes.NewValue(tftypes.String, nil)
		if action != nil {
			actionVal = tftypes.NewValue(tftypes.String, *action)
		}
		guardedVal := tftypes.NewValue(tftypes.String, nil)
		if guarded != nil {
			guardedVal = tftypes.NewValue(tftypes.String, *guarded)
		}

		return tfsdk.Config{
			Schema: testSchema,
			Raw: tftypes.NewValue(objType, map[string]tftypes.Value{
				"action":       actionVal,
				"guarded_attr": guardedVal,
			}),
		}
	}

	strPtr := func(s string) *string { return &s }

	testCases := map[string]struct {
		configValue basetypes.StringValue
		action      *string
		guarded     *string
		expectError bool
	}{
		"guarded-null-no-error": {
			configValue: types.StringNull(),
			action:      strPtr("OTHER"),
			guarded:     nil,
			expectError: false,
		},
		"guarded-unknown-no-error": {
			configValue: types.StringUnknown(),
			action:      strPtr("OTHER"),
			guarded:     nil,
			expectError: false,
		},
		"action-matches-no-error": {
			configValue: types.StringValue("some-id"),
			action:      strPtr("MATCH"),
			guarded:     strPtr("some-id"),
			expectError: false,
		},
		"action-does-not-match-error": {
			configValue: types.StringValue("some-id"),
			action:      strPtr("OTHER"),
			guarded:     strPtr("some-id"),
			expectError: true,
		},
		"action-null-no-error": {
			configValue: types.StringValue("some-id"),
			action:      nil,
			guarded:     strPtr("some-id"),
			expectError: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.StringRequest{
				Config:         newConfig(testCase.action, testCase.guarded),
				ConfigValue:    testCase.configValue,
				Path:           path.Root("guarded_attr"),
				PathExpression: path.MatchRoot("guarded_attr"),
			}
			response := validator.StringResponse{}

			schemavalidator.ConflictsIfDoesNotMatchPathValueValidator{
				TargetValue: types.StringValue("MATCH"),
				Expressions: path.Expressions{path.MatchRoot("action")},
			}.ValidateString(context.Background(), request, &response)

			if response.Diagnostics.HasError() && !testCase.expectError {
				t.Fatalf("unexpected error: %s", response.Diagnostics)
			}
			if !response.Diagnostics.HasError() && testCase.expectError {
				t.Fatal("expected error, got none")
			}
		})
	}
}
