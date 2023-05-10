package credentials

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type isbase64EncodedValidator struct {
	//StringValue string
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v isbase64EncodedValidator) Description(ctx context.Context) string {
	return "string value must be a properly base64 encoded image"
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v isbase64EncodedValidator) MarkdownDescription(ctx context.Context) string {
	return "string value must be a properly base64 encoded image"
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v isbase64EncodedValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	s := req.ConfigValue.ValueString()

	// determine if the input value has an image encoding prefix
	re := regexp.MustCompile(`^data:image\/(\w+);base64,`)
	matches := re.Split(s, 2)

	if len(matches) == 2 {
		// parse out an image prefix if present
		s = matches[1]
	} else {
		// no match - evaluate input as-is
		s = matches[0]
	}

	// check encoding
	_, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid String",
			fmt.Sprintf("string must be a properly base64 encoded image: %s", req.Path),
		)

		return
	}
}
