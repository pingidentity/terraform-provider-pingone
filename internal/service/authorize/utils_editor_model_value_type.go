package authorize

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

func valueTypeObjectSchemaAttributes(customValidators ...stringvalidatorinternal.CustomStringValidatorModel) (attributes map[string]schema.Attribute) {

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of the value.",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataValueTypeDTOEnumValues)

	validators := []validator.String{
		stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataValueTypeDTOEnumValues)...),
	}

	if len(customValidators) > 0 {
		for _, customValidator := range customValidators {
			validators = append(validators, customValidator.Validators...)
			typeDescription = typeDescription.Append(customValidator.Description)
		}

	}

	attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description:         typeDescription.Description,
			MarkdownDescription: typeDescription.MarkdownDescription,
			Required:            true,

			Validators: validators,
		},
	}

	return
}

func dataSourceValueTypeObjectSchemaAttributes(customValidators ...stringvalidatorinternal.CustomStringValidatorModel) (attributes map[string]dsschema.Attribute) {

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of the value.",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataValueTypeDTOEnumValues)

	if len(customValidators) > 0 {
		for _, customValidator := range customValidators {
			typeDescription = typeDescription.Append(customValidator.Description)
		}

	}

	attributes = map[string]dsschema.Attribute{
		"type": schema.StringAttribute{
			Description:         typeDescription.Description,
			MarkdownDescription: typeDescription.MarkdownDescription,
			Computed:            true,
		},
	}

	return
}

type editorValueTypeResourceModel struct {
	Type types.String `tfsdk:"type"`
}

var editorValueTypeTFObjectTypes = map[string]attr.Type{
	"type": types.StringType,
}

func expandEditorValueType(ctx context.Context, valueType basetypes.ObjectValue) (valueTypeObject *authorize.AuthorizeEditorDataValueTypeDTO, diags diag.Diagnostics) {
	var plan *editorValueTypeResourceModel
	diags.Append(valueType.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return
	}

	valueTypeObject = plan.expand()

	return
}

func (p *editorValueTypeResourceModel) expand() *authorize.AuthorizeEditorDataValueTypeDTO {
	return authorize.NewAuthorizeEditorDataValueTypeDTO(authorize.EnumAuthorizeEditorDataValueTypeDTO(p.Type.ValueString()))
}

func editorValueTypeOkToTF(apiObject *authorize.AuthorizeEditorDataValueTypeDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorValueTypeTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorValueTypeTFObjectTypes, map[string]attr.Value{
		"type": framework.EnumOkToTF(apiObject.GetTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
