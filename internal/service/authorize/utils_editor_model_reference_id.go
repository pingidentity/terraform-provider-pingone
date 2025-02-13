package authorize

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

func referenceIdObjectSchemaAttributes(description framework.SchemaAttributeDescription) (attributes map[string]schema.Attribute) {

	description = description.AppendMarkdownString("Must be a valid PingOne resource ID.")

	attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         description.Description,
			MarkdownDescription: description.MarkdownDescription,
			Required:            true,

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

func editorDataReferenceObjectOkToSetTF(apiObject []authorize.AuthorizeEditorDataReferenceObjectDTO, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := editorDataReferenceObjectOkToTF(&v, true)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
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
