// Copyright © 2025 Ping Identity Corporation

package setplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
)

func RequiresReplaceIfPreviouslyNull() setplanmodifier.RequiresReplaceIfFunc {
	return func(ctx context.Context, req planmodifier.SetRequest, resp *setplanmodifier.RequiresReplaceIfFuncResponse) {
		// If the configuration is unknown, this cannot be sure what to do yet.
		if req.ConfigValue.IsUnknown() {
			resp.RequiresReplace = false
			return
		}

		// If the state is null and the config value is not null, replace the resource.
		if req.StateValue.IsNull() && !req.ConfigValue.IsNull() {
			resp.RequiresReplace = true
			return
		}

		resp.RequiresReplace = false
	}
}
