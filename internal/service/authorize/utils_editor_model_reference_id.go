package authorize

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

func referenceIdObjectSchemaAttributes() (attributes map[string]schema.Attribute) {
	attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A string that specifies the %s resource's parent ID.  Must be a valid PingOne resource ID.", descriptionName)).Description,
			Required:    true,

			CustomType: pingonetypes.ResourceIDType{},
		},
	}

	return
}

type editorReferenceDataResourceModel struct {
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
}

var (
	editorReferenceObjectTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}
)

func expandEditorReferenceData(ctx context.Context, referenceData basetypes.ObjectValue) (referenceDataObject *authorize.AuthorizeEditorDataReferenceObjectDTO, diags diag.Diagnostics) {
	var plan *editorReferenceDataResourceModel
	diags.Append(referenceData.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return
	}

	referenceDataObject = plan.expand()

	return
}

func (p *editorReferenceDataResourceModel) expand() *authorize.AuthorizeEditorDataReferenceObjectDTO {
	data := authorize.NewAuthorizeEditorDataReferenceObjectDTO(p.Id.ValueString())
	return data
}

func editorDataReferenceObjectOkToTF(apiObject *authorize.AuthorizeEditorDataReferenceObjectDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorReferenceObjectTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorReferenceObjectTFObjectTypes, map[string]attr.Value{
		"id": framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
