package authorize

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

func parentObjectSchema(descriptionName string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("An object that specifies configuration settings for the %s resource's parent.", descriptionName)).Description,
		Optional:    true,

		Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the parent resource.")),
	}
}

type editorParentResourceModel editorReferenceDataResourceModel

func expandEditorParent(ctx context.Context, parent basetypes.ObjectValue) (parentObject *authorize.AuthorizeEditorDataReferenceObjectDTO, diags diag.Diagnostics) {
	var plan *editorParentResourceModel
	diags.Append(parent.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return
	}

	parentObject = plan.expand()

	return
}

func (p *editorParentResourceModel) expand() *authorize.AuthorizeEditorDataReferenceObjectDTO {
	referenceDataResourceModel := editorReferenceDataResourceModel(*p)
	return referenceDataResourceModel.expand()
}

func editorParentOkToTF(apiObject *authorize.AuthorizeEditorDataReferenceObjectDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	return editorDataReferenceObjectOkToTF(apiObject, ok)
}
