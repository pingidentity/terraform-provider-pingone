// Influenced from github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier/*
package setplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func UnmodifiableDataLossProtectionIfPreviouslyNull() UnmodifiableDataLossProtectionIfFunc {
	return func(ctx context.Context, req planmodifier.SetRequest, resp *UnmodifiableDataLossProtectionIfFuncResponse) {
		// If the configuration is unknown, this cannot be sure what to do yet.
		if req.ConfigValue.IsUnknown() {
			resp.Error = false
			return
		}

		// If the state is null and the config value is not null, error
		if req.StateValue.IsNull() && !req.ConfigValue.IsNull() {
			resp.Error = true
			return
		}

		resp.Error = false
	}
}
