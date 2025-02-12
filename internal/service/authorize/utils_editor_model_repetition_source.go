package authorize

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

func repetitionSourceObjectSchema(descriptionName string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("An object that specifies configuration settings for the %s resource's repetition source.", descriptionName)).Description,
		Optional:    true,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A string that specifies the %s resource's repetition source ID.  Must be a valid PingOne resource ID.", descriptionName)).Description,
				Required:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},
		},
	}
}

func dataSourceRepetitionSourceObjectSchema(descriptionName string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("An object that specifies configuration settings for the %s resource's repetition source.", descriptionName)).Description,
		Computed:    true,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A string that specifies the %s resource's repetition source ID.  Must be a valid PingOne resource ID.", descriptionName)).Description,
				Computed:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},
		},
	}
}

type editorRepetitionSourceResourceModel editorReferenceDataResourceModel

func expandEditorRepetitionSource(ctx context.Context, repetitionSource basetypes.ObjectValue) (repetitionSourceObject *authorize.AuthorizeEditorDataReferenceObjectDTO, diags diag.Diagnostics) {
	var plan *editorRepetitionSourceResourceModel
	diags.Append(repetitionSource.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return
	}

	repetitionSourceObject = plan.expand()

	return
}

func (p *editorRepetitionSourceResourceModel) expand() *authorize.AuthorizeEditorDataReferenceObjectDTO {
	referenceDataResourceModel := editorReferenceDataResourceModel(*p)
	return referenceDataResourceModel.expand()
}

func editorRepetitionSourceOkToTF(apiObject *authorize.AuthorizeEditorDataReferenceObjectDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	return editorDataReferenceObjectOkToTF(apiObject, ok)
}
