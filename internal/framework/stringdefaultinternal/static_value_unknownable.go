package stringdefault

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// StaticString returns a static string value default handler.
//
// Use StaticString if a static default value for a string should be set.
func StaticStringUnknownable(defaultVal basetypes.StringValue) defaults.String {
	return staticStringUnknownableDefault{
		defaultVal: defaultVal,
	}
}

// staticStringDefault is static value default handler that
// sets a value on a string attribute.
type staticStringUnknownableDefault struct {
	defaultVal basetypes.StringValue
}

// Description returns a human-readable description of the default value handler.
func (d staticStringUnknownableDefault) Description(_ context.Context) string {
	return fmt.Sprintf("value defaults to %s", d.defaultVal)
}

// MarkdownDescription returns a markdown description of the default value handler.
func (d staticStringUnknownableDefault) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value defaults to `%s`", d.defaultVal)
}

// DefaultString implements the static default value logic.
func (d staticStringUnknownableDefault) DefaultString(_ context.Context, req defaults.StringRequest, resp *defaults.StringResponse) {
	resp.PlanValue = d.defaultVal
}
