// Copyright © 2025 Ping Identity Corporation

// Package framework provides utilities for Terraform Plugin Framework implementation in the PingOne provider.
// This package contains common schema attributes, data models, and schema builders for consistent
// resource and data source definitions across the provider.
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

// DataFilterModel represents the structure for data source filtering configuration.
// It defines the filter criteria used to narrow down results in data source queries.
type DataFilterModel struct {
	// Name specifies the attribute name to filter on
	Name types.String `tfsdk:"name"`
	// Values contains the list of acceptable values for the named attribute
	Values types.List `tfsdk:"values"`
}

// Attr_ID creates a standard ID attribute schema for resources using the default PingOne resource ID type.
// It returns a computed string attribute that uses state for unknown values and cannot be modified by users.
// This attribute is automatically populated by the provider with the resource's unique identifier.
func Attr_ID() schema.StringAttribute {
	return Attr_IDCustomType(pingonetypes.ResourceIDType{})
}

// Attr_IDCustomType creates a standard ID attribute schema using a custom string type.
// It returns a computed string attribute with the specified custom type for specialized ID handling.
// The customType parameter must implement basetypes.StringTypable for custom ID value processing.
// This function is useful for resources that need specialized ID types beyond the standard PingOne resource ID.
func Attr_IDCustomType(customType basetypes.StringTypable) schema.StringAttribute {
	return schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},

		CustomType: customType,
	}
}

// Attr_LinkID creates a required link ID attribute schema with basic validation.
// It returns a required string attribute that references another PingOne resource and triggers replacement when changed.
// The description parameter provides documentation for the attribute's purpose and usage.
// This attribute uses the PingOne resource ID custom type and requires replacement when modified.
func Attr_LinkID(description SchemaAttributeDescription) schema.StringAttribute {
	return Attr_LinkIDWithValidators(description, []validator.String{})
}

// Attr_LinkIDWithValidators creates a required link ID attribute schema with custom validators.
// It returns a required string attribute with additional validation rules beyond the basic link ID validation.
// The description parameter provides documentation for the attribute's purpose and usage.
// The validators parameter allows adding custom validation logic specific to the resource's requirements.
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

// Attr_SCIMFilter creates a SCIM filter attribute schema for data source filtering.
// It returns an optional string attribute that accepts SCIM filter expressions for querying resources.
// The description parameter provides documentation for the attribute's purpose and usage.
// The acceptableAttributes parameter lists the resource attributes that can be used in the filter.
// The mutuallyExclusiveAttributes parameter lists other filtering attributes that cannot be used simultaneously.
// This attribute includes validation for minimum length and exactly-one-of constraints with other filter attributes.
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

// Attr_DataFilter creates a data filter attribute schema for structured data source filtering.
// It returns an optional list of nested objects that define filter criteria for data source queries.
// The description parameter provides documentation for the attribute's purpose and usage.
// The acceptableAttributes parameter lists the resource attributes that can be used in the filter.
// The mutuallyExclusiveAttributes parameter lists other filtering attributes that cannot be used simultaneously.
// Each filter object contains a name (attribute to filter on) and values (acceptable values for that attribute).
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

// Attr_DataSourceReturnIDs creates a computed list attribute for returning resource IDs from data sources.
// It returns a computed list attribute containing PingOne resource IDs that match the data source query.
// The description parameter provides documentation for the attribute's purpose and usage.
// This attribute is automatically populated by the provider with matching resource identifiers.
func Attr_DataSourceReturnIDs(description SchemaAttributeDescription) schema.ListAttribute {
	return Attr_DataSourceReturnIDsByElement(description, pingonetypes.ResourceIDType{})
}

// Attr_DataSourceReturnIDsByElement creates a computed list attribute with custom element type.
// It returns a computed list attribute containing values of the specified element type.
// The description parameter provides documentation for the attribute's purpose and usage.
// The elementType parameter specifies the type of elements contained in the returned list.
// This function is useful for data sources that return lists of non-standard ID types.
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

// Attr_DataSourceReturnIDsSet creates a computed set attribute for returning unique resource IDs from data sources.
// It returns a computed set attribute containing PingOne resource IDs that match the data source query.
// The description parameter provides documentation for the attribute's purpose and usage.
// This attribute is automatically populated by the provider with unique matching resource identifiers.
func Attr_DataSourceReturnIDsSet(description SchemaAttributeDescription) schema.SetAttribute {
	return Attr_DataSourceReturnIDsByElementSet(description, pingonetypes.ResourceIDType{})
}

// Attr_DataSourceReturnIDsByElementSet creates a computed set attribute with custom element type.
// It returns a computed set attribute containing unique values of the specified element type.
// The description parameter provides documentation for the attribute's purpose and usage.
// The elementType parameter specifies the type of elements contained in the returned set.
// This function is useful for data sources that return sets of non-standard ID types.
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
