package stringplanmodifierinternal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringDefaultValue(v types.String, description, markdownDescription string) planmodifier.String {
	return &stringDefaultValuePlanModifier{
		v,
		description,
		markdownDescription,
	}
}

type stringDefaultValuePlanModifier struct {
	defaultValue        types.String
	description         string
	markdownDescription string
}

var _ planmodifier.String = (*stringDefaultValuePlanModifier)(nil)

func (apm *stringDefaultValuePlanModifier) Description(ctx context.Context) string {
	return apm.description
}

func (apm *stringDefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return apm.markdownDescription
}

func (apm *stringDefaultValuePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, res *planmodifier.StringResponse) {
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
