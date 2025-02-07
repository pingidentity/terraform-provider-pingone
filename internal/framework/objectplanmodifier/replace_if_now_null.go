// Copyright © 2025 Ping Identity Corporation

package objectplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func RequiresReplaceIfNowNull() objectplanmodifier.RequiresReplaceIfFunc {
	return func(ctx context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifier.RequiresReplaceIfFuncResponse) {
		// If the configuration is unknown, this cannot be sure what to do yet.
		if req.ConfigValue.IsUnknown() {
			resp.RequiresReplace = false
			return
		}

		// If the state is not null and the config value is null, replace the resource.
		if !req.StateValue.IsNull() && req.ConfigValue.IsNull() {
			resp.RequiresReplace = true
			return
		}

		resp.RequiresReplace = false
	}
}
