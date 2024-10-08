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
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

func dataInputObjectSchemaAttributes() (attributes map[string]schema.Attribute) {

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataInputDTOTypeEnumValues)

	attributeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAINPUTDTOTYPE_ATTRIBUTE)))

	valueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAINPUTDTOTYPE_CONSTANT)))

	attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description:         typeDescription.Description,
			MarkdownDescription: typeDescription.MarkdownDescription,
			Required:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataInputDTOTypeEnumValues)...),
			},
		},

		"attribute": schema.SingleNestedAttribute{
			Description:         attributeDescription.Description,
			MarkdownDescription: attributeDescription.MarkdownDescription,
			Optional:            true,

			Attributes: referenceIdObjectSchemaAttributes(),

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAINPUTDTOTYPE_ATTRIBUTE)),
					path.MatchRoot("service_type"),
				),
			},
		},

		"value": schema.StringAttribute{
			Description:         valueDescription.Description,
			MarkdownDescription: valueDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.String{
				stringvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAINPUTDTOTYPE_CONSTANT)),
					path.MatchRoot("service_type"),
				),
			},
		},
	}

	return
}

type editorDataInputResourceModel struct {
	Type      types.String `tfsdk:"type"`
	Attribute types.Object `tfsdk:"attribute"`
	Value     types.String `tfsdk:"value"`
}

var editorDataInputTFObjectTypes = map[string]attr.Type{
	"type": types.StringType,
	"attribute": types.ObjectType{
		AttrTypes: editorReferenceObjectTFObjectTypes,
	},
	"value": types.StringType,
}

func (p *editorDataInputResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataInputDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorize.AuthorizeEditorDataInputDTO{}

	switch authorize.EnumAuthorizeEditorDataInputMappingDTOType(p.Type.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATAINPUTMAPPINGDTOTYPE_ATTRIBUTE:
		data.AuthorizeEditorDataInputsAttributeInputDTO, d = p.expandAttributeInput(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAINPUTMAPPINGDTOTYPE_INPUT:
		data.AuthorizeEditorDataInputsConstantInputDTO = p.expandConstantInput()
	default:
		diags.AddError(
			"Invalid service settings header value input mapping type",
			fmt.Sprintf("The service settings header value input mapping type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &data, diags
}

func (p *editorDataInputResourceModel) expandAttributeInput(ctx context.Context) (*authorize.AuthorizeEditorDataInputsAttributeInputDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	attribute, d := expandEditorReferenceData(ctx, p.Attribute)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataInputsAttributeInputDTO(
		p.Type.ValueString(),
		*attribute,
	)

	return data, diags
}

func (p *editorDataInputResourceModel) expandConstantInput() *authorize.AuthorizeEditorDataInputsConstantInputDTO {
	data := authorize.NewAuthorizeEditorDataInputsConstantInputDTO(
		p.Type.ValueString(),
		p.Value.ValueString(),
	)

	return data
}

func editorDataInputOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataInputDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorize.AuthorizeEditorDataInputDTO{}) {
		return types.ObjectNull(editorDataInputTFObjectTypes), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case authorize.AuthorizeEditorDataInputsAttributeInputDTO:

		attributeResp, ok := t.GetAttributeOk()
		value, d := editorDataReferenceObjectOkToTF(attributeResp, ok)
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"type":      framework.EnumOkToTF(t.GetTypeOk()),
			"attribute": value,
			"value":     types.StringNull(),
		}

	case authorize.AuthorizeEditorDataInputsConstantInputDTO:

		attributeMap = map[string]attr.Value{
			"type":      framework.EnumOkToTF(t.GetTypeOk()),
			"attribute": types.ObjectNull(editorReferenceObjectTFObjectTypes),
			"value":     framework.StringOkToTF(t.GetValueOk()),
		}

	default:
		tflog.Error(ctx, "Invalid data input mapping type", map[string]interface{}{
			"service data input mapping type": t,
		})
		diags.AddError(
			"Invalid data input mapping type",
			"The data input mapping type is not supported.  Please raise an issue with the provider maintainers.",
		)
	}

	objValue, d := types.ObjectValue(editorDataInputTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}
