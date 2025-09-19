// Copyright © 2025 Ping Identity Corporation

// Package stringplanmodifier provides custom string plan modifiers for the Terraform Plugin Framework.
// This package contains plan modifiers that handle string attribute lifecycle management,
// including data loss protection and immutability enforcement for the PingOne provider.
//
// Influenced from github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier/*
package stringplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// UnmodifiableDataLossProtection creates a plan modifier that prevents changes to protect against data loss.
// It returns a plan modifier that prevents any modifications to a string attribute once it has been set,
// requiring manual resource replacement to change the value. This is used for critical attributes where
// modification could result in data loss and requires explicit user action.
//
// The plan modifier will generate an error for any attempted changes, directing users to use
// Terraform's replace command option to manually replace the resource when changes are needed.
func UnmodifiableDataLossProtection() planmodifier.String {
	return UnmodifiableDataLossProtectionIf(
		func(_ context.Context, _ planmodifier.StringRequest, resp *UnmodifiableDataLossProtectionIfFuncResponse) {
			resp.Error = true
		},
		"This field is immutable and cannot be changed once defined.  To protect against accidental data loss, this resource must be replaced manually (for example, by using Terraform's plan `-replace` command option https://developer.hashicorp.com/terraform/cli/commands/plan#replace-address).  Any data that is stored against this resource must be manually exported before the resource is removed and re-imported once the resource has been replaced.",
		"This field is immutable and cannot be changed once defined.  To protect against accidental data loss, this resource must be replaced manually (for example, by using Terraform's [plan `-replace` command option](https://developer.hashicorp.com/terraform/cli/commands/plan#replace-address)).  Any data that is stored against this resource must be manually exported before the resource is removed and re-imported once the resource has been replaced.",
	)
}
