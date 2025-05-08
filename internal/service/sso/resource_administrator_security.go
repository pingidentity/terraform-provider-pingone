// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.ResourceWithModifyPlan = &administratorSecurityResource{}
)

func (r *administratorSecurityResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *administratorSecurityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if plan == nil {
		return
	}

	// Validate the identity_provider and authentication_method fields
	if plan.IdentityProvider.IsUnknown() || plan.AuthenticationMethod.IsUnknown() {
		return
	}

	if plan.AuthenticationMethod.ValueString() == "PINGONE" && !plan.IdentityProvider.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("identity_provider"),
			"Invalid configuration",
			"The `identity_provider` field must not be set if `authentication_method` is set to `PINGONE`."+
				" To configure an identity provider, set `authentication_method` to `EXTERNAL` or `HYBRID`.",
		)
	}

	if (plan.AuthenticationMethod.ValueString() == "EXTERNAL" || plan.AuthenticationMethod.ValueString() == "HYBRID") &&
		plan.IdentityProvider.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("identity_provider"),
			"Invalid configuration",
			"The `identity_provider` field must be set if `authentication_method` is set to `EXTERNAL` or `HYBRID`.",
		)
	}
}
