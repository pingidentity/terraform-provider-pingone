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

func dataConditionComparandObjectLeftSchemaAttributes() (attributes map[string]schema.Attribute) {

	allowedValues := []authorize.EnumAuthorizeEditorDataConditionsComparandDTOType{
		"ATTRIBUTE",
	}

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the authorization condition comparand type.",
	).AllowedValuesEnum(allowedValues)

	idDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the ID of the authorization attribute in the trust framework to use as the condition comparand.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONSCOMPARANDDTOTYPE_ATTRIBUTE))).AppendMarkdownString("Must be a valid PingOne resource ID.")

	attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description:         typeDescription.Description,
			MarkdownDescription: typeDescription.MarkdownDescription,
			Required:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(allowedValues)...),
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
	}

	return attributes
}

func dataConditionComparandObjectRightSchemaAttributes() (attributes map[string]schema.Attribute) {

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

type editorDataConditionComparandLeftResourceModel struct {
	Type types.String                 `tfsdk:"type"`
	Id   pingonetypes.ResourceIDValue `tfsdk:"id"`
}

type editorDataConditionComparandRightResourceModel struct {
	Type  types.String                 `tfsdk:"type"`
	Id    pingonetypes.ResourceIDValue `tfsdk:"id"`
	Value types.String                 `tfsdk:"value"`
}

var editorDataConditionComparandLeftTFObjectTypes = map[string]attr.Type{
	"type": types.StringType,
	"id":   pingonetypes.ResourceIDType{},
}

var editorDataConditionComparandRightTFObjectTypes = map[string]attr.Type{
	"type":  types.StringType,
	"id":    pingonetypes.ResourceIDType{},
	"value": types.StringType,
}

func expandEditorDataConditionLeftComparand(ctx context.Context, conditionComparand basetypes.ObjectValue) (conditionComparandObject *authorize.AuthorizeEditorDataConditionsComparandLeftDTO, diags diag.Diagnostics) {
	var plan *editorDataConditionComparandLeftResourceModel
	diags.Append(conditionComparand.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return
	}

	conditionComparandObject, d := plan.expandLeft()
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	return
}

func expandEditorDataConditionRightComparand(ctx context.Context, conditionComparand basetypes.ObjectValue) (conditionComparandObject *authorize.AuthorizeEditorDataConditionsComparandRightDTO, diags diag.Diagnostics) {
	var plan *editorDataConditionComparandRightResourceModel
	diags.Append(conditionComparand.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return
	}

	conditionComparandObject, d := plan.expandRight()
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	return
}

func (p *editorDataConditionComparandLeftResourceModel) expandLeft() (*authorize.AuthorizeEditorDataConditionsComparandLeftDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.AuthorizeEditorDataConditionsComparandLeftDTO{}

	switch authorize.EnumAuthorizeEditorDataConditionsComparandDTOType(p.Type.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATACONDITIONSCOMPARANDDTOTYPE_ATTRIBUTE:
		data.AuthorizeEditorDataConditionsComparandsAttributeComparandDTO = p.expandAttributeConditionComparand()
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

func (p *editorDataConditionComparandRightResourceModel) expandRight() (*authorize.AuthorizeEditorDataConditionsComparandRightDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.AuthorizeEditorDataConditionsComparandRightDTO{}

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

func (p *editorDataConditionComparandLeftResourceModel) expandAttributeConditionComparand() *authorize.AuthorizeEditorDataConditionsComparandsAttributeComparandDTO {
	return expandAttributeConditionComparand(p.Id.ValueString(), p.Type.ValueString())
}

func (p *editorDataConditionComparandRightResourceModel) expandAttributeConditionComparand() *authorize.AuthorizeEditorDataConditionsComparandsAttributeComparandDTO {
	return expandAttributeConditionComparand(p.Id.ValueString(), p.Type.ValueString())
}

func expandAttributeConditionComparand(comparandId, comparandType string) *authorize.AuthorizeEditorDataConditionsComparandsAttributeComparandDTO {

	data := authorize.NewAuthorizeEditorDataConditionsComparandsAttributeComparandDTO(
		authorize.EnumAuthorizeEditorDataConditionsComparandDTOType(comparandType),
		comparandId,
	)

	return data
}

func (p *editorDataConditionComparandRightResourceModel) expandConstantConditionComparand() *authorize.AuthorizeEditorDataConditionsComparandsConstantComparandDTO {

	data := authorize.NewAuthorizeEditorDataConditionsComparandsConstantComparandDTO(
		authorize.EnumAuthorizeEditorDataConditionsComparandDTOType(p.Type.ValueString()),
		p.Value.ValueString(),
	)

	return data
}

func editorDataConditionComparandLeftOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataConditionsComparandLeftDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorize.AuthorizeEditorDataConditionsComparandLeftDTO{}) {
		return types.ObjectNull(editorDataConditionComparandLeftTFObjectTypes), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case *authorize.AuthorizeEditorDataConditionsComparandsAttributeComparandDTO:

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["id"] = framework.PingOneResourceIDOkToTF(t.GetIdOk())

	default:
		tflog.Error(ctx, "Invalid left condition comparand type", map[string]interface{}{
			"condition comparand type": t,
		})
		diags.AddError(
			"Invalid left condition comparand type",
			"The condition comparand type is not supported.  Please raise an issue with the provider maintainers.",
		)
		return types.ObjectNull(editorDataConditionComparandLeftTFObjectTypes), diags
	}

	attributeMap = editorDataConditionComparandLeftConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(editorDataConditionComparandLeftTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorDataConditionComparandRightOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataConditionsComparandRightDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorize.AuthorizeEditorDataConditionsComparandRightDTO{}) {
		return types.ObjectNull(editorDataConditionComparandRightTFObjectTypes), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case *authorize.AuthorizeEditorDataConditionsComparandsAttributeComparandDTO:

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["id"] = framework.PingOneResourceIDOkToTF(t.GetIdOk())

	case *authorize.AuthorizeEditorDataConditionsComparandsConstantComparandDTO:

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["value"] = framework.StringOkToTF(t.GetValueOk())

	default:
		tflog.Error(ctx, "Invalid right condition comparand type", map[string]interface{}{
			"condition comparand type": t,
		})
		diags.AddError(
			"Invalid right condition comparand type",
			"The condition comparand type is not supported.  Please raise an issue with the provider maintainers.",
		)
		return types.ObjectNull(editorDataConditionComparandRightTFObjectTypes), diags
	}

	attributeMap = editorDataConditionComparandRightConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(editorDataConditionComparandRightTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorDataConditionComparandLeftConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
	return editorDataConditionComparandConvertEmptyValuesToTFNulls(attributeMap, map[string]attr.Value{
		"type": types.StringNull(),
		"id":   pingonetypes.NewResourceIDNull(),
	})
}

func editorDataConditionComparandRightConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
	return editorDataConditionComparandConvertEmptyValuesToTFNulls(attributeMap, map[string]attr.Value{
		"type":  types.StringNull(),
		"id":    pingonetypes.NewResourceIDNull(),
		"value": types.StringNull(),
	})
}

func editorDataConditionComparandConvertEmptyValuesToTFNulls(attributeMap, nullMap map[string]attr.Value) map[string]attr.Value {
	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}
