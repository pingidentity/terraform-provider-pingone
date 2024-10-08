package authorize

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

func dataConditionComparandObjectSchemaAttributes() (attributes map[string]schema.Attribute) {

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the authorization condition comparand type.",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataConditionsComparandDTOTypeEnumValues)

	idDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the ID of the authorization attribute in the trust framework to use as the condition comparand.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONSCOMPARANDDTOTYPE_ATTRIBUTE))).AppendMarkdownString("Must be a valid PingOne resource ID.")

	valueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a constant text value to use as the condition comparand.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONSCOMPARANDDTOTYPE_CONSTANT)))

	attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description:         typeDescription.Description,
			MarkdownDescription: typeDescription.MarkdownDescription,
			Required:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataConditionsComparandDTOTypeEnumValues)...),
			},
		},

		// type == "ATTRIBUTE"
		"id": schema.StringAttribute{
			Description:         idDescription.Description,
			MarkdownDescription: idDescription.MarkdownDescription,
			Optional:            true,

			CustomType: pingonetypes.ResourceIDType{},

			Validators: []validator.String{
				stringvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONSCOMPARANDDTOTYPE_ATTRIBUTE)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},
		},

		// type == "CONSTANT"
		"value": schema.StringAttribute{
			Description:         valueDescription.Description,
			MarkdownDescription: valueDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.String{
				stringvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONSCOMPARANDDTOTYPE_CONSTANT)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},
		},
	}

	return attributes
}

type editorDataConditionComparandResourceModel struct {
	Type  types.String `tfsdk:"type"`
	Id    types.String `tfsdk:"id"`
	Value types.String `tfsdk:"value"`
}

var editorDataConditionComparandTFObjectTypes = map[string]attr.Type{
	"type":  types.StringType,
	"id":    types.StringType,
	"value": types.StringType,
}

func expandEditorDataConditionComparand(ctx context.Context, conditionComparand basetypes.ObjectValue) (conditionComparandObject *authorize.AuthorizeEditorDataConditionsComparandDTO, diags diag.Diagnostics) {
	var plan *editorDataConditionComparandResourceModel
	diags.Append(conditionComparand.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return
	}

	conditionComparandObject, d := plan.expand()
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	return
}

func (p *editorDataConditionComparandResourceModel) expand() (*authorize.AuthorizeEditorDataConditionsComparandDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.AuthorizeEditorDataConditionsComparandDTO{}

	switch authorize.EnumAuthorizeEditorDataConditionsComparandDTOType(p.Type.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATACONDITIONSCOMPARANDDTOTYPE_ATTRIBUTE:
		data.AuthorizeEditorDataConditionsComparandsAttributeComparandDTO = p.expandAttributeConditionComparand()
	case authorize.ENUMAUTHORIZEEDITORDATACONDITIONSCOMPARANDDTOTYPE_CONSTANT:
		data.AuthorizeEditorDataConditionsComparandsConstantComparandDTO = p.expandConstantConditionComparand()
	default:
		diags.AddError(
			"Invalid condition comparand type",
			fmt.Sprintf("The condition comparand type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &data, diags
}

func (p *editorDataConditionComparandResourceModel) expandAttributeConditionComparand() *authorize.AuthorizeEditorDataConditionsComparandsAttributeComparandDTO {

	data := authorize.NewAuthorizeEditorDataConditionsComparandsAttributeComparandDTO(
		authorize.EnumAuthorizeEditorDataConditionsComparandDTOType(p.Type.ValueString()),
		p.Id.ValueString(),
	)

	return data
}

func (p *editorDataConditionComparandResourceModel) expandConstantConditionComparand() *authorize.AuthorizeEditorDataConditionsComparandsConstantComparandDTO {

	data := authorize.NewAuthorizeEditorDataConditionsComparandsConstantComparandDTO(
		authorize.EnumAuthorizeEditorDataConditionsComparandDTOType(p.Type.ValueString()),
		p.Value.ValueString(),
	)

	return data
}

func editorDataConditionComparandsOkToSetTF(ctx context.Context, apiObject []authorize.AuthorizeEditorDataConditionsComparandDTO, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: editorDataConditionComparandTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := editorDataConditionComparandOkToTF(ctx, &v, true)
		diags.Append(d...)
		if diags.HasError() {
			return types.SetNull(tfObjType), diags
		}

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func editorDataConditionComparandOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataConditionsComparandDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorize.AuthorizeEditorDataConditionsComparandDTO{}) {
		return types.ObjectNull(editorDataConditionComparandTFObjectTypes), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case authorize.AuthorizeEditorDataConditionsComparandsAttributeComparandDTO:

		attributeMap = map[string]attr.Value{
			"type": framework.EnumOkToTF(t.GetTypeOk()),
			"id":   framework.StringOkToTF(t.GetIdOk()),
		}

	case authorize.AuthorizeEditorDataConditionsComparandsConstantComparandDTO:

		attributeMap = map[string]attr.Value{
			"type":  framework.EnumOkToTF(t.GetTypeOk()),
			"value": framework.StringOkToTF(t.GetValueOk()),
		}

	default:
		tflog.Error(ctx, "Invalid condition comparand type", map[string]interface{}{
			"condition comparand type": t,
		})
		diags.AddError(
			"Invalid condition comparand type",
			"The condition comparand type is not supported.  Please raise an issue with the provider maintainers.",
		)
	}

	attributeMap = editorDataConditionComparandConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(editorDataConditionComparandTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorDataConditionComparandConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
	nullMap := map[string]attr.Value{
		"type":  types.StringNull(),
		"id":    types.StringNull(),
		"value": types.StringNull(),
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}
