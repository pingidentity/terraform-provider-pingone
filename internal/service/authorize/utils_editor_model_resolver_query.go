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
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

func dataResolverQueryObjectSchemaAttributes() (attributes map[string]schema.Attribute) {

	queryTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataAttributeResolversUserQueryDTOTypeEnumValues)

	queryUserIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAATTRIBUTERESOLVERSUSERQUERYDTOTYPE_USER_ID)))

	attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description:         queryTypeDescription.Description,
			MarkdownDescription: queryTypeDescription.MarkdownDescription,
			Required:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataAttributeResolversUserQueryDTOTypeEnumValues)...),
			},
		},

		"user_id": schema.StringAttribute{
			Description:         queryUserIdDescription.Description,
			MarkdownDescription: queryUserIdDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.String{
				stringvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAATTRIBUTERESOLVERSUSERQUERYDTOTYPE_USER_ID)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},
		},
	}

	return attributes
}

type editorDataResolverQueryResourceModel struct {
	Type   types.String `tfsdk:"type"`
	UserId types.String `tfsdk:"user_id"`
}

var (
	editorDataResolverQueryTFObjectTypes = map[string]attr.Type{
		"type":    types.StringType,
		"user_id": types.StringType,
	}
)

func expandEditorResolverQuery(ctx context.Context, query basetypes.ObjectValue) (queryObject *authorize.AuthorizeEditorDataAttributeResolversUserQueryDTO, diags diag.Diagnostics) {
	var plan *editorDataResolverQueryResourceModel
	diags.Append(query.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return
	}

	queryObject, d := plan.expand()
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	return
}

func (p *editorDataResolverQueryResourceModel) expand() (*authorize.AuthorizeEditorDataAttributeResolversUserQueryDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.AuthorizeEditorDataAttributeResolversUserQueryDTO{}

	switch authorize.EnumAuthorizeEditorDataAttributeResolversUserQueryDTOType(p.Type.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATAATTRIBUTERESOLVERSUSERQUERYDTOTYPE_USER_ID:
		data.AuthorizeEditorDataAttributeResolversUserQueryUserIdQueryDTO = p.expandUserResolverQueryUserID()
	default:
		diags.AddError(
			"Invalid resolver query type",
			fmt.Sprintf("The resolver query type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &data, diags
}

func (p *editorDataResolverQueryResourceModel) expandUserResolverQueryUserID() *authorize.AuthorizeEditorDataAttributeResolversUserQueryUserIdQueryDTO {

	data := authorize.NewAuthorizeEditorDataAttributeResolversUserQueryUserIdQueryDTO(
		authorize.EnumAuthorizeEditorDataAttributeResolversUserQueryDTOType(p.Type.ValueString()),
		p.UserId.ValueString(),
	)

	return data
}

func editorDataResolverQueryOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataAttributeResolversUserQueryDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorize.AuthorizeEditorDataAttributeResolversUserQueryDTO{}) {
		return types.ObjectNull(editorDataResolverQueryTFObjectTypes), diags
	}

	var attributeMap = map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case *authorize.AuthorizeEditorDataAttributeResolversUserQueryUserIdQueryDTO:

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["user_id"] = framework.StringOkToTF(t.GetUserIdOk())

	default:
		tflog.Error(ctx, "Invalid resolver query type", map[string]interface{}{
			"resolver query type": t,
		})
		diags.AddError(
			"Invalid resolver query type",
			"The resolver query type is not supported.  Please raise an issue with the provider maintainers.",
		)
		return types.ObjectNull(editorDataResolverQueryTFObjectTypes), diags
	}

	attributeMap = editorDataResolverConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(editorDataResolverTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorDataResolverQueryConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
	nullMap := map[string]attr.Value{
		"type":    types.StringNull(),
		"user_id": types.StringNull(),
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}
