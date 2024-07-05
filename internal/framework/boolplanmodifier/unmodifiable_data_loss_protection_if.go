// Influenced from github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier/*

package boolplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func UnmodifiableDataLossProtectionIf(f UnmodifiableDataLossProtectionIfFunc, description, markdownDescription string) planmodifier.Bool {
	return requiresReplaceIfModifier{
		ifFunc:              f,
		description:         description,
		markdownDescription: markdownDescription,
	}
}

type requiresReplaceIfModifier struct {
	ifFunc              UnmodifiableDataLossProtectionIfFunc
	description         string
	markdownDescription string
}

func (m requiresReplaceIfModifier) Description(_ context.Context) string {
	return m.description
}

func (m requiresReplaceIfModifier) MarkdownDescription(_ context.Context) string {
	return m.markdownDescription
}

func (m requiresReplaceIfModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	// Do not check on resource creation.
	if req.State.Raw.IsNull() {
		return
	}

	// Do not check on resource destroy.
	if req.Plan.Raw.IsNull() {
		return
	}

	// Do not error if the plan and state values are equal.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	ifFuncResp := &UnmodifiableDataLossProtectionIfFuncResponse{}

	m.ifFunc(ctx, req, ifFuncResp)

	resp.Diagnostics.Append(ifFuncResp.Diagnostics...)

	if ifFuncResp.Error {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Data Loss Protection",
			m.Description(ctx),
		)
	}
}
