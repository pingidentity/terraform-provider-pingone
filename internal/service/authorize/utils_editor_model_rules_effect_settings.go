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
	objectvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

func dataRulesEffectSettingsObjectSchemaAttributes() (attributes map[string]schema.Attribute) {

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of the policy combination effect settings.",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataRulesEffectSettingsDTOTypeEnumValues)

	conditionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies the configuration settings for the condition to apply to the conditional policy combination effect settings.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s` or `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATARULESEFFECTSETTINGSDTOTYPE_CONDITIONAL_DENY_ELSE_PERMIT), string(authorize.ENUMAUTHORIZEEDITORDATARULESEFFECTSETTINGSDTOTYPE_CONDITIONAL_PERMIT_ELSE_DENY)))

	attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description:         typeDescription.Description,
			MarkdownDescription: typeDescription.MarkdownDescription,
			Required:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataRulesEffectSettingsDTOTypeEnumValues)...),
			},
		},

		// type == "UNCONDITIONAL_PERMIT"
		// (same as base object)

		// type == "UNCONDITIONAL_DENY"
		// (same as base object)

		// type == "CONDITIONAL_PERMIT_ELSE_DENY", type == "CONDITIONAL_DENY_ELSE_PERMIT"
		"condition": schema.SingleNestedAttribute{
			Description:         conditionDescription.Description,
			MarkdownDescription: conditionDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATARULESEFFECTSETTINGSDTOTYPE_CONDITIONAL_DENY_ELSE_PERMIT)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATARULESEFFECTSETTINGSDTOTYPE_CONDITIONAL_PERMIT_ELSE_DENY)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				objectvalidatorinternal.ConflictsIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATARULESEFFECTSETTINGSDTOTYPE_UNCONDITIONAL_DENY)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				objectvalidatorinternal.ConflictsIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATARULESEFFECTSETTINGSDTOTYPE_UNCONDITIONAL_PERMIT)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: dataConditionObjectSchemaAttributesIteration(1),
		},
	}

	return attributes
}

type editorDataRulesEffectSettingsResourceModel struct {
	Type      types.String `tfsdk:"type"`
	Condition types.Object `tfsdk:"condition"`
}

var editorDataRulesEffectSettingsTFObjectTypes = map[string]attr.Type{
	"type":      types.StringType,
	"condition": types.ObjectType{AttrTypes: editorDataConditionTFObjectTypes},
}

func expandEditorDataRulesEffectSettings(ctx context.Context, rulesEffectSettings basetypes.ObjectValue) (rulesEffectSettingsObject *authorize.AuthorizeEditorDataRulesEffectSettingsDTO, diags diag.Diagnostics) {
	var plan *editorDataRulesEffectSettingsResourceModel
	diags.Append(rulesEffectSettings.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return
	}

	rulesEffectSettingsObject, d := plan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	return
}

func (p *editorDataRulesEffectSettingsResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataRulesEffectSettingsDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorize.AuthorizeEditorDataRulesEffectSettingsDTO{}

	switch authorize.EnumAuthorizeEditorDataRulesEffectSettingsDTOType(p.Type.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATARULESEFFECTSETTINGSDTOTYPE_CONDITIONAL_DENY_ELSE_PERMIT:
		data.AuthorizeEditorDataRulesEffectSettingsConditionalDenyElsePermitDTO, d = p.expandConditionalDenyElsePermitRulesEffectSettings(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATARULESEFFECTSETTINGSDTOTYPE_CONDITIONAL_PERMIT_ELSE_DENY:
		data.AuthorizeEditorDataRulesEffectSettingsConditionalPermitElseDenyDTO, d = p.expandConditionalPermitElseDenyRulesEffectSettings(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATARULESEFFECTSETTINGSDTOTYPE_UNCONDITIONAL_DENY:
		data.AuthorizeEditorDataRulesEffectSettingsUnconditionalDenyDTO = p.expandUnconditionalDenyRulesEffectSettings()
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATARULESEFFECTSETTINGSDTOTYPE_UNCONDITIONAL_PERMIT:
		data.AuthorizeEditorDataRulesEffectSettingsUnconditionalPermitDTO = p.expandUnconditionalPermitRulesEffectSettings()
		diags.Append(d...)
	default:
		diags.AddError(
			"Invalid Rules Effect Settings type",
			fmt.Sprintf("The rules effect settings type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &data, diags
}

func (p *editorDataRulesEffectSettingsResourceModel) expandConditionalDenyElsePermitRulesEffectSettings(ctx context.Context) (*authorize.AuthorizeEditorDataRulesEffectSettingsConditionalDenyElsePermitDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	condition, d := expandEditorDataCondition(ctx, p.Condition)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataRulesEffectSettingsConditionalDenyElsePermitDTO(
		authorize.EnumAuthorizeEditorDataRulesEffectSettingsDTOType(p.Type.ValueString()),
		*condition,
	)

	return data, diags
}

func (p *editorDataRulesEffectSettingsResourceModel) expandConditionalPermitElseDenyRulesEffectSettings(ctx context.Context) (*authorize.AuthorizeEditorDataRulesEffectSettingsConditionalPermitElseDenyDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	condition, d := expandEditorDataCondition(ctx, p.Condition)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataRulesEffectSettingsConditionalPermitElseDenyDTO(
		authorize.EnumAuthorizeEditorDataRulesEffectSettingsDTOType(p.Type.ValueString()),
		*condition,
	)

	return data, diags
}

func (p *editorDataRulesEffectSettingsResourceModel) expandUnconditionalDenyRulesEffectSettings() *authorize.AuthorizeEditorDataRulesEffectSettingsUnconditionalDenyDTO {

	data := authorize.NewAuthorizeEditorDataRulesEffectSettingsUnconditionalDenyDTO(
		authorize.EnumAuthorizeEditorDataRulesEffectSettingsDTOType(p.Type.ValueString()),
	)

	return data
}

func (p *editorDataRulesEffectSettingsResourceModel) expandUnconditionalPermitRulesEffectSettings() *authorize.AuthorizeEditorDataRulesEffectSettingsUnconditionalPermitDTO {

	data := authorize.NewAuthorizeEditorDataRulesEffectSettingsUnconditionalPermitDTO(
		authorize.EnumAuthorizeEditorDataRulesEffectSettingsDTOType(p.Type.ValueString()),
	)

	return data
}

func editorDataRulesEffectSettingsOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataRulesEffectSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorize.AuthorizeEditorDataRulesEffectSettingsDTO{}) {
		return types.ObjectNull(editorDataRulesEffectSettingsTFObjectTypes), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case *authorize.AuthorizeEditorDataRulesEffectSettingsConditionalDenyElsePermitDTO:

		conditionResp, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTF(ctx, conditionResp, ok)
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"type":      framework.EnumOkToTF(t.GetTypeOk()),
			"condition": condition,
		}

	case *authorize.AuthorizeEditorDataRulesEffectSettingsConditionalPermitElseDenyDTO:

		conditionResp, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTF(ctx, conditionResp, ok)
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"type":      framework.EnumOkToTF(t.GetTypeOk()),
			"condition": condition,
		}

	case *authorize.AuthorizeEditorDataRulesEffectSettingsUnconditionalDenyDTO:

		attributeMap = map[string]attr.Value{
			"type": framework.EnumOkToTF(t.GetTypeOk()),
		}

	case *authorize.AuthorizeEditorDataRulesEffectSettingsUnconditionalPermitDTO:

		attributeMap = map[string]attr.Value{
			"type": framework.EnumOkToTF(t.GetTypeOk()),
		}

	default:
		tflog.Error(ctx, "Invalid Rules Effect Settings type", map[string]interface{}{
			"Rules Effect Settings type": t,
		})
		diags.AddError(
			"Invalid Rules Effect Settings type",
			"The rules effect settings type is not supported.  Please raise an issue with the provider maintainers.",
		)
	}

	attributeMap = editorDataRulesEffectSettingsConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(editorDataRulesEffectSettingsTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorDataRulesEffectSettingsConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
	nullMap := map[string]attr.Value{
		"type":      types.StringNull(),
		"condition": types.ObjectNull(editorDataConditionTFObjectTypes),
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}
