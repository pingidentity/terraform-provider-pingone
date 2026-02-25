// Copyright © 2026 Ping Identity Corporation

package objectplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func RequiresReplaceIfExistenceChanges() planmodifier.Object {
	return requiresReplaceIfExistenceChangesModifier{}
}

type requiresReplaceIfExistenceChangesModifier struct {
}

// Description returns a human-readable description of the plan modifier.
func (m requiresReplaceIfExistenceChangesModifier) Description(_ context.Context) string {
	return "If the object is set and becomes unset, or is unset and becomes set, Terraform will destroy and recreate the resource."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m requiresReplaceIfExistenceChangesModifier) MarkdownDescription(_ context.Context) string {
	return "If the object is set and becomes unset, or is unset and becomes set, Terraform will destroy and recreate the resource."
}

// PlanModifyString implements the plan modification logic.
func (m requiresReplaceIfExistenceChangesModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// Creation plan
	if req.State.Raw.IsNull() {
		return
	}

	// Destruction plan
	if req.Plan.Raw.IsNull() {
		return
	}

	// If the state is not null and the plan value is null or unknown, replace the resource.
	if !req.StateValue.IsNull() && (req.PlanValue.IsNull() || req.PlanValue.IsUnknown()) {
		resp.RequiresReplace = true
		return
	}

	// If the state is null and the plan value is not null, replace the resource.
	if req.StateValue.IsNull() && !req.PlanValue.IsNull() {
		resp.RequiresReplace = true
		return
	}

	resp.RequiresReplace = false
}
