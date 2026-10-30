// Copyright © 2026 Ping Identity Corporation

package base

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pingidentity/pingone-go-sdk-legacy/management"
	"github.com/pingidentity/pingone-go-sdk-legacy/pingone"
	"github.com/pingidentity/pingone-go-sdk-legacy/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func TestLanguageUpdate_NotFound(t *testing.T) {
	t.Parallel()

	const environmentID = "11111111-1111-1111-1111-111111111111"
	const languageID = "22222222-2222-2222-2222-222222222222"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/environments/"+environmentID+"/languages/"+languageID {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"id":"33333333-3333-3333-3333-333333333333","code":"NOT_FOUND","message":"Language not found."}`))
	}))
	defer server.Close()

	config := management.NewConfiguration()
	config.Servers = management.ServerConfigurations{{URL: server.URL}}
	config.HTTPClient = server.Client()
	p1Client := &client.Client{API: &pingone.Client{ManagementAPIClient: management.NewAPIClient(config)}}
	for _, operation := range []string{"create", "update"} {
		t.Run(operation, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, ResourceLanguageUpdate().Schema, map[string]interface{}{
				"environment_id": environmentID,
				"language_id":    languageID,
				"enabled":        true,
				"default":        true,
			})

			var diags diag.Diagnostics
			if operation == "create" {
				// The SDK marks new resources before invoking CreateContext.
				d.MarkNewResource()
				diags = resourceLanguageUpdateCreate(context.Background(), d, p1Client)
				if len(diags) != 2 || !diags.HasError() || diags[1].Summary != "Cannot create language update" {
					t.Fatalf("expected an actionable missing-language error, got %v", diags)
				}
				wantDetail := "The language " + languageID + " was not found in environment " + environmentID + ". Ensure `language_id` references an existing language in that environment before creating this resource."
				if diags[1].Detail != wantDetail {
					t.Fatalf("expected detail %q, got %q", wantDetail, diags[1].Detail)
				}
			} else {
				d.SetId(languageID)
				diags = resourceLanguageUpdateUpdate(context.Background(), d, p1Client)
				if diags.HasError() || len(diags) != 1 || diags[0].Severity != diag.Warning {
					t.Fatalf("expected a not-found warning, got %v", diags)
				}
			}
			if d.Id() != "" {
				t.Fatalf("expected no state for the missing language, got ID %q", d.Id())
			}
		})
	}
}

func TestLanguageUpdate_AlreadyDefault(t *testing.T) {
	t.Parallel()

	for _, locale := range []string{"en", "fr"} {
		t.Run(locale, func(t *testing.T) {
			t.Parallel()

			language := management.NewLanguage(true, true, locale)
			// A nil API client proves that this no-op never issues an update.
			diags := updateLanguageEnabledDefaultSequence(context.Background(), nil, "environment", "language", *language, true, true, false)
			if len(diags) != 0 {
				t.Fatalf("expected no diagnostics for an unchanged default language, got %v", diags)
			}
		})
	}
}

func TestLanguageUpdate_InvalidDefault(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		enabled      bool
		defaultValue bool
		wantSummary  string
	}{
		{"disable default", false, true, "Invalid combination of `enabled` and `default`, the language must be enabled to be made default. Got enabled=false, default=true."},
		{"unset default English", true, false, "The default language cannot be un-defaulted"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			language := management.NewLanguage(true, true, "en")
			diags := updateLanguageEnabledDefaultSequence(context.Background(), nil, "environment", "language", *language, tt.enabled, tt.defaultValue, false)
			if len(diags) != 1 || !diags.HasError() || diags[0].Summary != tt.wantSummary {
				t.Fatalf("expected error %q, got %v", tt.wantSummary, diags)
			}
		})
	}
}

func TestRetryLanguageDefault(t *testing.T) {
	t.Parallel()

	p1Error := model.NewP1Error()
	p1Error.SetMessage("The language must be enabled before it is set as the default.")
	if !retryLanguageDefault(context.Background(), nil, p1Error) {
		t.Fatal("expected enablement propagation to be retryable")
	}
	p1Error.SetMessage("You cannot reset the default language.")
	if retryLanguageDefault(context.Background(), nil, p1Error) {
		t.Fatal("unexpected retry for an unrelated language error")
	}
	if retryLanguageDefault(context.Background(), nil, nil) {
		t.Fatal("unexpected retry without an API error")
	}
}
