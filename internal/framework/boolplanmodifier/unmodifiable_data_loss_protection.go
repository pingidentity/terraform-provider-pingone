// Copyright © 2026 Ping Identity Corporation

// Influenced from github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier/*

package boolplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func UnmodifiableDataLossProtection() planmodifier.Bool {
	return UnmodifiableDataLossProtectionIf(
		func(_ context.Context, _ planmodifier.BoolRequest, resp *UnmodifiableDataLossProtectionIfFuncResponse) {
			resp.Error = true
		},
		"This field is immutable and cannot be changed once defined.  To protect against accidental data loss, this resource must be replaced manually (for example, by using Terraform's plan `-replace` command option https://developer.hashicorp.com/terraform/cli/commands/plan#replace-address).  Any data that is stored against this resource must be manually exported before the resource is removed and re-imported once the resource has been replaced.",
		"This field is immutable and cannot be changed once defined.  To protect against accidental data loss, this resource must be replaced manually (for example, by using Terraform's [plan `-replace` command option](https://developer.hashicorp.com/terraform/cli/commands/plan#replace-address)).  Any data that is stored against this resource must be manually exported before the resource is removed and re-imported once the resource has been replaced.",
	)
}
