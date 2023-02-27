package boolplanmodifierinternal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func BoolDefaultValue(v types.Bool, description, markdownDescription string) planmodifier.Bool {
	return &boolDefaultValuePlanModifier{
		v,
		description,
		markdownDescription,
	}
}

type boolDefaultValuePlanModifier struct {
	defaultValue        types.Bool
	description         string
	markdownDescription string
}

var _ planmodifier.Bool = (*boolDefaultValuePlanModifier)(nil)

func (apm *boolDefaultValuePlanModifier) Description(ctx context.Context) string {
	return apm.description
}

func (apm *boolDefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return apm.markdownDescription
}

func (apm *boolDefaultValuePlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, res *planmodifier.BoolResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	res.PlanValue = apm.defaultValue
}
