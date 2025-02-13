package authorize

import (
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

func combiningAlgorithmObjectSchemaAttributes() (attributes map[string]schema.Attribute) {

	combiningAlgorithmAlgorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the algorithm that determines how rules are combined to produce an authorization decision.",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataPoliciesCombiningAlgorithmDTOAlgorithmEnumValues)

	attributes = map[string]schema.Attribute{
		"algorithm": schema.StringAttribute{
			Description:         combiningAlgorithmAlgorithmDescription.Description,
			MarkdownDescription: combiningAlgorithmAlgorithmDescription.MarkdownDescription,
			Required:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataPoliciesCombiningAlgorithmDTOAlgorithmEnumValues)...),
			},
		},
	}

	return
}

type policyManagementPolicyCombiningAlgorithmResourceModel struct {
	Algorithm types.String `tfsdk:"algorithm"`
}

var policyManagementPolicyCombiningAlgorithmTFObjectTypes = map[string]attr.Type{
	"algorithm": types.StringType,
}

func (p *policyManagementPolicyCombiningAlgorithmResourceModel) expand() *authorize.AuthorizeEditorDataPoliciesCombiningAlgorithmDTO {

	data := authorize.NewAuthorizeEditorDataPoliciesCombiningAlgorithmDTO(
		authorize.EnumAuthorizeEditorDataPoliciesCombiningAlgorithmDTOAlgorithm(p.Algorithm.ValueString()),
	)

	return data
}

func policyManagementPolicyCombiningAlgorithmOkToTF(apiObject *authorize.AuthorizeEditorDataPoliciesCombiningAlgorithmDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(policyManagementPolicyCombiningAlgorithmTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(policyManagementPolicyCombiningAlgorithmTFObjectTypes, map[string]attr.Value{
		"algorithm": framework.EnumOkToTF(apiObject.GetAlgorithmOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
