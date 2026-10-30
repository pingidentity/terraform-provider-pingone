// Copyright © 2026 Ping Identity Corporation

package base

import (
	"testing"

	"github.com/pingidentity/pingone-go-sdk-legacy/pingone/model"
)

func TestNotificationTemplateCustomWriteError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		code        string
		message     string
		detailCode  string
		target      string
		wantSummary string
	}{
		{
			name:        "read-only content",
			code:        "INVALID_DATA",
			message:     "The resource is read only and can not be modified.",
			wantSummary: "The default template content is read-only and cannot be modified.",
		},
		{
			name:    "unrelated invalid data",
			code:    "INVALID_DATA",
			message: "Invalid notification content.",
		},
		{
			name:    "different error code",
			code:    "INVALID_REQUEST",
			message: "The resource is read only and can not be modified.",
		},
		{
			name:        "invalid delivery method",
			code:        "INVALID_DATA",
			detailCode:  "INVALID_VALUE",
			target:      "deliveryMethod",
			wantSummary: "The configured delivery method does not apply to the selected template.",
		},
		{
			name:        "invalid locale",
			code:        "INVALID_DATA",
			detailCode:  "INVALID_VALUE",
			target:      "language",
			wantSummary: "The locale is not valid for the environment.",
		},
		{
			name:        "missing value",
			code:        "INVALID_DATA",
			detailCode:  "REQUIRED_VALUE",
			wantSummary: "Content body is missing a required value.",
		},
		{
			name:        "duplicate content",
			code:        "INVALID_DATA",
			detailCode:  "UNIQUENESS_VIOLATION",
			wantSummary: "Customized content for the template, locale and variant combination already exists.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p1Error := model.NewP1Error()
			p1Error.SetCode(tt.code)
			p1Error.SetMessage(tt.message)
			if tt.detailCode != "" {
				detail := model.NewP1ErrorDetailsInner()
				detail.SetCode(tt.detailCode)
				detail.SetTarget(tt.target)
				detail.SetMessage("Notification content validation failed.")
				p1Error.SetDetails([]model.P1ErrorDetailsInner{*detail})
			}

			diags := notificationTemplateCustomWriteError(nil, p1Error)
			if tt.wantSummary == "" {
				if len(diags) != 0 {
					t.Fatalf("expected generic error handling, got %v", diags)
				}
				return
			}
			if len(diags) != 1 || !diags.HasError() || diags[0].Summary() != tt.wantSummary {
				t.Fatalf("expected error %q, got %v", tt.wantSummary, diags)
			}
		})
	}

	if diags := notificationTemplateCustomWriteError(nil, nil); len(diags) != 0 {
		t.Fatalf("expected no custom diagnostic without an API error, got %v", diags)
	}
}
