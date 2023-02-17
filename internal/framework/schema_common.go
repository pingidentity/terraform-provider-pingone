package framework

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

type SchemaDescription struct {
	Description         string
	MarkdownDescription string
}

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
