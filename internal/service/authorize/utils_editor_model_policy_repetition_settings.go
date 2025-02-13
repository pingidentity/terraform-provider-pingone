package authorize

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

func repetitionSettingsObjectSchemaAttributes() (attributes map[string]schema.Attribute) {

	repetitionSettingsDecisionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the decision filter.",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataPoliciesRepetitionSettingsDTODecisionEnumValues)

	attributes = map[string]schema.Attribute{
		"source": schema.SingleNestedAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings for the source associated with the policy rule.").Description,
			Required:    true,

			Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the authorization policy repetition source in the policy manager.")),
		},

		"decision": schema.StringAttribute{
			Description:         repetitionSettingsDecisionDescription.Description,
			MarkdownDescription: repetitionSettingsDecisionDescription.MarkdownDescription,
			Required:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataPoliciesRepetitionSettingsDTODecisionEnumValues)...),
			},
		},
	}

	return
}

type policyManagementPolicyRepetitionSettingsResourceModel struct {
	Source   types.Object `tfsdk:"source"`
	Decision types.String `tfsdk:"decision"`
}

var policyManagementPolicyRepetitionSettingsTFObjectTypes = map[string]attr.Type{
	"source":   types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
	"decision": types.StringType,
}

func (p *policyManagementPolicyRepetitionSettingsResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataPoliciesRepetitionSettingsDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	source, d := expandEditorReferenceData(ctx, p.Source)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataPoliciesRepetitionSettingsDTO(
		*source,
		authorize.EnumAuthorizeEditorDataPoliciesRepetitionSettingsDTODecision(p.Decision.ValueString()),
	)

	return data, diags
}

func policyManagementPolicyRepetitionSettingsOkToTF(apiObject *authorize.AuthorizeEditorDataPoliciesRepetitionSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(policyManagementPolicyRepetitionSettingsTFObjectTypes), diags
	}

	source, d := editorDataReferenceObjectOkToTF(apiObject.GetSourceOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(policyManagementPolicyRepetitionSettingsTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(policyManagementPolicyRepetitionSettingsTFObjectTypes, map[string]attr.Value{
		"source":   source,
		"decision": framework.EnumOkToTF(apiObject.GetDecisionOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
