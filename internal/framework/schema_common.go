package framework

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

// Common models
type DataFilterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.List   `tfsdk:"values"`
}

// Common schema attributes
func Attr_ID() schema.StringAttribute {
	return Attr_IDCustomType(pingonetypes.ResourceIDType{})
}

func Attr_IDCustomType(customType basetypes.StringTypable) schema.StringAttribute {
	return schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},

		CustomType: customType,
	}
}

func Attr_LinkID(description SchemaAttributeDescription) schema.StringAttribute {
	return Attr_LinkIDWithValidators(description, []validator.String{})
}

func Attr_LinkIDWithValidators(description SchemaAttributeDescription, validators []validator.String) schema.StringAttribute {

	if description.MarkdownDescription == "" {
		description.MarkdownDescription = description.Description
	}

	description = description.AppendMarkdownString("Must be a valid PingOne resource ID.").RequiresReplace()

	return schema.StringAttribute{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,
		Required:            true,

		CustomType: pingonetypes.ResourceIDType{},

		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
		Validators: validators,
	}
}

func Attr_SCIMFilter(description SchemaAttributeDescription, acceptableAttributes []string, mutuallyExclusiveAttributes []string) schema.StringAttribute {
	filterMinLength := 1

	description = description.Clean(true)

	description.MarkdownDescription = fmt.Sprintf("%s.  The SCIM filter can use the following attributes: `%s`.", description.MarkdownDescription, strings.Join(acceptableAttributes, "`, `"))
	description.Description = fmt.Sprintf("%s.  The SCIM filter can use the following attributes: \"%s\".", description.Description, strings.Join(acceptableAttributes, "\", \""))

	description = description.ExactlyOneOf(mutuallyExclusiveAttributes)

	validators := make([]validator.String, 0)
	validators = append(validators, stringvalidator.LengthAtLeast(filterMinLength))

	paths := make([]path.Expression, 0)
	for _, v := range mutuallyExclusiveAttributes {
		paths = append(paths, path.MatchRelative().AtParent().AtName(v))
	}
	validators = append(validators, stringvalidator.ExactlyOneOf(paths...))

	return schema.StringAttribute{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,
		Optional:            true,
		Validators:          validators,
	}
}

func Attr_DataFilter(description SchemaAttributeDescription, acceptableAttributes []string, mutuallyExclusiveAttributes []string) schema.ListNestedAttribute {
	attrMinLength := 1

	description = description.Clean(true)

	description.MarkdownDescription = fmt.Sprintf("%s.  Allowed attributes to filter: `%s`", description.MarkdownDescription, strings.Join(acceptableAttributes, "`, `"))
	description.Description = fmt.Sprintf("%s.  Allowed attributes to filter: \"%s\"", description.Description, strings.Join(acceptableAttributes, "\", \""))

	description = description.ExactlyOneOf(mutuallyExclusiveAttributes)

	childNameAttrDescriptionFmt := fmt.Sprintf("The attribute name to filter on.  Must be one of the following values: `%s`.", strings.Join(acceptableAttributes, "`, `"))
	childNameDescription := SchemaAttributeDescription{
		MarkdownDescription: childNameAttrDescriptionFmt,
		Description:         strings.Replace(childNameAttrDescriptionFmt, "`", "\"", -1),
	}

	childValueAttrDescriptionFmt := "The possible values (case sensitive) of the attribute defined in the `name` parameter to filter."
	childValueDescription := SchemaAttributeDescription{
		MarkdownDescription: childValueAttrDescriptionFmt,
		Description:         strings.Replace(childValueAttrDescriptionFmt, "`", "\"", -1),
	}

	// The parent attribute validators
	validators := make([]validator.List, 0)

	paths := make([]path.Expression, 0)
	for _, v := range mutuallyExclusiveAttributes {
		paths = append(paths, path.MatchRelative().AtParent().AtName(v))
	}
	validators = append(validators, listvalidator.ExactlyOneOf(paths...))

	return schema.ListNestedAttribute{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,
		Optional:            true,

		NestedObject: schema.NestedAttributeObject{

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
		Validators: validators,
	}
}

func Attr_DataSourceReturnIDs(description SchemaAttributeDescription) schema.ListAttribute {
	return Attr_DataSourceReturnIDsByElement(description, pingonetypes.ResourceIDType{})
}

func Attr_DataSourceReturnIDsByElement(description SchemaAttributeDescription, elementType attr.Type) schema.ListAttribute {
	if description.MarkdownDescription == "" {
		description.MarkdownDescription = description.Description
	}

	return schema.ListAttribute{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,
		Computed:            true,
		ElementType:         elementType,
	}
}

func Attr_DataSourceReturnIDsSet(description SchemaAttributeDescription) schema.SetAttribute {
	return Attr_DataSourceReturnIDsByElementSet(description, pingonetypes.ResourceIDType{})
}

func Attr_DataSourceReturnIDsByElementSet(description SchemaAttributeDescription, elementType attr.Type) schema.SetAttribute {
	if description.MarkdownDescription == "" {
		description.MarkdownDescription = description.Description
	}

	return schema.SetAttribute{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,
		Computed:            true,
		ElementType:         elementType,
	}
}
