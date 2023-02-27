package framework

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type SchemaDescription struct {
	Description         string
	MarkdownDescription string
}

// Common models
type DataFilterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.List   `tfsdk:"values"`
}

// Common schema attributes
func Attr_ID() schema.StringAttribute {
	return schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func Attr_EnvironmentID(description SchemaDescription) schema.StringAttribute {
	if description.MarkdownDescription == "" {
		description.MarkdownDescription = description.Description
	}

	return schema.StringAttribute{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
		Validators: []validator.String{
			verify.P1ResourceIDValidator(),
		},
	}
}

func Attr_ApplicationID(description SchemaDescription) schema.StringAttribute {
	if description.MarkdownDescription == "" {
		description.MarkdownDescription = description.Description
	}

	return schema.StringAttribute{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
		Validators: []validator.String{
			verify.P1ResourceIDValidator(),
		},
	}
}

func Attr_SCIMFilter(description SchemaDescription, acceptableAttributes []string, mutuallyExclusiveAttributes []string) schema.StringAttribute {
	filterMinLength := 1

	description.Clean(true)

	description.MarkdownDescription = fmt.Sprintf("%s.  The SCIM filter can use the following attributes: `%s`.", description.MarkdownDescription, strings.Join(acceptableAttributes, "`, `"))
	description.Description = fmt.Sprintf("%s.  The SCIM filter can use the following attributes: \"%s\".", description.Description, strings.Join(acceptableAttributes, "\", \""))

	stringValidators := make([]validator.String, 0)
	stringValidators = append(stringValidators, stringvalidator.LengthAtLeast(filterMinLength))
	for _, v := range mutuallyExclusiveAttributes {
		stringValidators = append(stringValidators, stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName(v)))
	}

	return schema.StringAttribute{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,
		Optional:            true,
		Validators:          stringValidators,
	}
}

func Attr_DataFilter(description SchemaDescription, acceptableAttributes []string, mutuallyExclusiveAttributes []string) schema.ListNestedBlock {
	attrMinLength := 1

	description.Clean(true)

	description.MarkdownDescription = fmt.Sprintf("%s.  Allowed attributes to filter: `%s`", description.MarkdownDescription, strings.Join(acceptableAttributes, "`, `"))
	description.Description = fmt.Sprintf("%s.  Allowed attributes to filter: \"%s\"", description.Description, strings.Join(acceptableAttributes, "\", \""))

	childNameAttrDescriptionFmt := fmt.Sprintf("The attribute name to filter on.  Must be one of the following values: `%s`.", strings.Join(acceptableAttributes, "`, `"))
	childNameDescription := SchemaDescription{
		MarkdownDescription: childNameAttrDescriptionFmt,
		Description:         strings.Replace(childNameAttrDescriptionFmt, "`", "\"", -1),
	}

	childValueAttrDescriptionFmt := "The possible values (case sensitive) of the attribute defined in the `name` parameter to filter."
	childValueDescription := SchemaDescription{
		MarkdownDescription: childValueAttrDescriptionFmt,
		Description:         strings.Replace(childValueAttrDescriptionFmt, "`", "\"", -1),
	}

	// The parent attribute validators
	listValidators := make([]validator.List, 0)
	for _, v := range mutuallyExclusiveAttributes {
		listValidators = append(listValidators, listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName(v)))
	}

	return schema.ListNestedBlock{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,

		NestedObject: schema.NestedBlockObject{

			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description:         childNameDescription.Description,
					MarkdownDescription: childNameDescription.MarkdownDescription,
					Required:            true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(attrMinLength),
						stringvalidator.OneOf(acceptableAttributes...),
					},
				},
				"values": schema.ListAttribute{
					ElementType:         types.StringType,
					Description:         childValueDescription.Description,
					MarkdownDescription: childValueDescription.MarkdownDescription,
					Required:            true,
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
						listvalidator.UniqueValues(),
						listvalidator.ValueStringsAre(
							stringvalidator.LengthAtLeast(attrMinLength),
						),
					},
				},
			},
		},
		Validators: listValidators,
	}
}

func Attr_DataSourceReturnIDs(description SchemaDescription) schema.ListAttribute {
	if description.MarkdownDescription == "" {
		description.MarkdownDescription = description.Description
	}

	return schema.ListAttribute{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

// Helpers
func (description *SchemaDescription) Clean(removeTrailingStop bool) {

	// Trim trailing fullstop
	if removeTrailingStop {
		trailingDot := regexp.MustCompile(`(\.\s*)$`)
		description.Description = trailingDot.ReplaceAllString(description.Description, "")
		description.MarkdownDescription = trailingDot.ReplaceAllString(description.MarkdownDescription, "")
	}

	description.Description = strings.TrimSpace(description.Description)
	description.MarkdownDescription = strings.TrimSpace(description.MarkdownDescription)

	if description.MarkdownDescription == "" && description.Description != "" {
		// Prefil the blank markdown description with the description value
		description.MarkdownDescription = description.Description
	}

	if description.MarkdownDescription != "" && description.Description == "" {
		// Prefil the blank description with the markdown description value, ignoring MD formatting
		description.Description = description.MarkdownDescription
	}

}
