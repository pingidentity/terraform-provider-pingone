// Copyright © 2026 Ping Identity Corporation

package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider"
)

// TestMuxProviderSchema guards against schema drift between the three provider
// flavours: the mux server fails with "Invalid Provider Server Combination" when
// schemas or descriptions differ.
func TestMuxProviderSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	server, err := provider.ProviderServerFactoryV6(ctx, "test")
	if err != nil {
		t.Fatalf("failed to build mux server: %v", err)
	}

	schemaResp, err := server().GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema failed: %v", err)
	}

	for _, diagnostic := range schemaResp.Diagnostics {
		if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
			t.Fatalf("mux schema has error diagnostic: %s: %s", diagnostic.Summary, diagnostic.Detail)
		}
	}

	found := map[string]bool{
		"config_profile":    false,
		"config_grant_type": false,
	}

	for _, attribute := range schemaResp.Provider.Block.Attributes {
		if attribute.Name == "config_path" {
			t.Errorf("provider schema must not contain the %q attribute", attribute.Name)
		}

		if _, wanted := found[attribute.Name]; wanted {
			found[attribute.Name] = true
		}
	}

	for name, ok := range found {
		if !ok {
			t.Errorf("provider schema is missing the %q attribute", name)
		}
	}
}
