package authorize

import (
	"context"
	"fmt"
	"strings"

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

func dataResolverObjectSchemaAttributes() (attributes map[string]schema.Attribute) {

	resolversTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the resolver type.",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataResolverDTOTypeEnumValues)

	valueRefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies configuration settings for the authorization attribute (if `type` is `%s`) or the authorization service (if `type` is `%s`) to use as the data value.", string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_ATTRIBUTE), string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_SERVICE)),
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s` or `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_ATTRIBUTE), string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_SERVICE)))

	valueTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that describes configuration settings for the output value type when using a constant value.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_CONSTANT)))

	valueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies a constant text value to use as the resulting data value.  If `type` is `%s`, the options are `%s`.  If `type` is `%s`, any value can be configured.", string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_SYSTEM), strings.Join(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataAttributeResolversSystemResolverDTOValueEnumValues), "`, `"), string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_CONSTANT)),
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s` or `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_CONSTANT), string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_SYSTEM)))

	queryDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for the query to use to resolve the data value.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_USER)))

	attributes = map[string]schema.Attribute{
		"condition": schema.ListNestedAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings for an authorization condition to apply to the resolver.").Description,
			Optional:    true,

			NestedObject: schema.NestedAttributeObject{
				Attributes: dataConditionObjectSchemaAttributes(),
			},
		},

		"name": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a name to apply to the resolver.").Description,
			Optional:    true,
		},

		"processor": schema.ListNestedAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings for a processor to apply to the resolver.").Description,
			Optional:    true,

			NestedObject: schema.NestedAttributeObject{
				Attributes: dataProcessorObjectSchemaAttributes(),
			},
		},

		"type": schema.StringAttribute{
			Description:         resolversTypeDescription.Description,
			MarkdownDescription: resolversTypeDescription.MarkdownDescription,
			Required:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataResolverDTOTypeEnumValues)...),
			},
		},

		// type == ATTRIBUTE, type == SERVICE
		"value_ref": schema.SingleNestedAttribute{
			Description:         valueRefDescription.Description,
			MarkdownDescription: valueRefDescription.MarkdownDescription,
			Optional:            true,

			Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A string that specifies the ID of the authorization attribute (if `type` is `%s`) or the authorization service (if `type` is `%s`) in the trust framework.", string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_ATTRIBUTE), string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_SERVICE)))),

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_ATTRIBUTE)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_SERVICE)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},
		},

		// type == CONSTANT
		"value_type": schema.SingleNestedAttribute{
			Description:         valueTypeDescription.Description,
			MarkdownDescription: valueTypeDescription.MarkdownDescription,
			Optional:            true,

			Attributes: valueTypeObjectSchemaAttributes(),

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_CONSTANT)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},
		},

		// type == CURRENT_REPETITION_VALUE
		// (same as base object)

		// type == CURRENT_USER_ID
		// (same as base object)

		// type == REQUEST
		// (same as base object)

		// type == CONSTANT, type == SYSTEM
		"value": schema.StringAttribute{
			Description:         valueDescription.Description,
			MarkdownDescription: valueDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.String{
				stringvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_CONSTANT)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				stringvalidator.All(
					stringvalidatorinternal.IsRequiredIfMatchesPathValue(
						types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_SYSTEM)),
						path.MatchRelative().AtParent().AtName("type"),
					),
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataAttributeResolversSystemResolverDTOValueEnumValues)...),
				),
			},
		},

		// type == USER
		"query": schema.SingleNestedAttribute{
			Description:         queryDescription.Description,
			MarkdownDescription: queryDescription.MarkdownDescription,
			Optional:            true,

			Attributes: dataResolverQueryObjectSchemaAttributes(),

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_USER)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},
		},
	}

	return attributes
}

type editorDataResolverResourceModel struct {
	Condition types.Object `tfsdk:"condition"`
	Name      types.String `tfsdk:"name"`
	Processor types.Object `tfsdk:"processor"`
	Type      types.String `tfsdk:"type"`
	ValueRef  types.Object `tfsdk:"value_ref"`
	ValueType types.Object `tfsdk:"value_type"`
	Value     types.String `tfsdk:"value"`
	Query     types.Object `tfsdk:"query"`
}

var (
	editorDataResolverTFObjectTypes = map[string]attr.Type{
		"condition":  types.ObjectType{AttrTypes: editorDataConditionTFObjectTypes},
		"name":       types.StringType,
		"processor":  types.ObjectType{AttrTypes: editorDataProcessorTFObjectTypes},
		"type":       types.StringType,
		"value_ref":  types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
		"value_type": types.ObjectType{AttrTypes: editorValueTypeTFObjectTypes},
		"value":      types.StringType,
		"query":      types.ObjectType{AttrTypes: editorDataResolverQueryTFObjectTypes},
	}
)

func expandEditorResolver(ctx context.Context, resolver basetypes.ObjectValue) (resolverObject *authorize.AuthorizeEditorDataResolverDTO, diags diag.Diagnostics) {
	var plan *editorDataResolverResourceModel
	diags.Append(resolver.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return
	}

	resolverObject, d := plan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	return
}

func (p *editorDataResolverResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataResolverDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorize.AuthorizeEditorDataResolverDTO{}

	switch authorize.EnumAuthorizeEditorDataResolverDTOType(p.Type.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_ATTRIBUTE:
		data.AuthorizeEditorDataAttributeResolversAttributeResolverDTO, d = p.expandAttributeResolver(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_CONSTANT:
		data.AuthorizeEditorDataAttributeResolversConstantResolverDTO, d = p.expandConstantResolver(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_CURRENT_REPETITION_VALUE:
		data.AuthorizeEditorDataAttributeResolversCurrentRepetitionValueResolverDTO, d = p.expandCurrentRepetitionValueResolver(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_CURRENT_USER_ID:
		data.AuthorizeEditorDataAttributeResolversCurrentUserIDResolverDTO, d = p.expandCurrentUserIdResolver(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_REQUEST:
		data.AuthorizeEditorDataAttributeResolversRequestResolverDTO, d = p.expandRequestResolver(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_SERVICE:
		data.AuthorizeEditorDataAttributeResolversServiceResolverDTO, d = p.expandServiceResolver(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_SYSTEM:
		data.AuthorizeEditorDataAttributeResolversSystemResolverDTO, d = p.expandSystemResolver(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATARESOLVERDTOTYPE_USER:
		data.AuthorizeEditorDataAttributeResolversUserResolverDTO, d = p.expandUserResolver(ctx)
		diags.Append(d...)
	default:
		diags.AddError(
			"Invalid resolver type",
			fmt.Sprintf("The resolver type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &data, diags
}

func (p *editorDataResolverResourceModel) expandAttributeResolver(ctx context.Context) (*authorize.AuthorizeEditorDataAttributeResolversAttributeResolverDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	valueRef, d := expandEditorReferenceData(ctx, p.ValueRef)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataAttributeResolversAttributeResolverDTO(
		authorize.EnumAuthorizeEditorDataResolverDTOType(p.Type.ValueString()),
		*valueRef,
	)

	// Condition
	if !p.Condition.IsNull() && !p.Condition.IsUnknown() {

		condition, d := expandEditorDataCondition(ctx, p.Condition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetCondition(*condition)
	}

	// Name
	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	// Processor
	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {

		processor, d := expandEditorDataProcessor(ctx, p.Processor)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetProcessor(*processor)
	}

	return data, diags
}

func (p *editorDataResolverResourceModel) expandConstantResolver(ctx context.Context) (*authorize.AuthorizeEditorDataAttributeResolversConstantResolverDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	valueType, d := expandEditorValueType(ctx, p.ValueType)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataAttributeResolversConstantResolverDTO(
		authorize.EnumAuthorizeEditorDataResolverDTOType(p.Type.ValueString()),
		p.Value.ValueString(),
		*valueType,
	)

	// Condition
	if !p.Condition.IsNull() && !p.Condition.IsUnknown() {

		condition, d := expandEditorDataCondition(ctx, p.Condition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetCondition(*condition)
	}

	// Name
	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	// Processor
	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {

		processor, d := expandEditorDataProcessor(ctx, p.Processor)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetProcessor(*processor)
	}

	return data, diags
}

func (p *editorDataResolverResourceModel) expandCurrentRepetitionValueResolver(ctx context.Context) (*authorize.AuthorizeEditorDataAttributeResolversCurrentRepetitionValueResolverDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.NewAuthorizeEditorDataAttributeResolversCurrentRepetitionValueResolverDTO(
		authorize.EnumAuthorizeEditorDataResolverDTOType(p.Type.ValueString()),
	)

	// Condition
	if !p.Condition.IsNull() && !p.Condition.IsUnknown() {

		condition, d := expandEditorDataCondition(ctx, p.Condition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetCondition(*condition)
	}

	// Name
	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	// Processor
	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {

		processor, d := expandEditorDataProcessor(ctx, p.Processor)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetProcessor(*processor)
	}

	return data, diags
}

func (p *editorDataResolverResourceModel) expandCurrentUserIdResolver(ctx context.Context) (*authorize.AuthorizeEditorDataAttributeResolversCurrentUserIDResolverDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.NewAuthorizeEditorDataAttributeResolversCurrentUserIDResolverDTO(
		authorize.EnumAuthorizeEditorDataResolverDTOType(p.Type.ValueString()),
	)

	// Condition
	if !p.Condition.IsNull() && !p.Condition.IsUnknown() {

		condition, d := expandEditorDataCondition(ctx, p.Condition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetCondition(*condition)
	}

	// Name
	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	// Processor
	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {

		processor, d := expandEditorDataProcessor(ctx, p.Processor)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetProcessor(*processor)
	}

	return data, diags
}

func (p *editorDataResolverResourceModel) expandRequestResolver(ctx context.Context) (*authorize.AuthorizeEditorDataAttributeResolversRequestResolverDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.NewAuthorizeEditorDataAttributeResolversRequestResolverDTO(
		authorize.EnumAuthorizeEditorDataResolverDTOType(p.Type.ValueString()),
	)

	// Condition
	if !p.Condition.IsNull() && !p.Condition.IsUnknown() {

		condition, d := expandEditorDataCondition(ctx, p.Condition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetCondition(*condition)
	}

	// Name
	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	// Processor
	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {

		processor, d := expandEditorDataProcessor(ctx, p.Processor)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetProcessor(*processor)
	}

	return data, diags
}

func (p *editorDataResolverResourceModel) expandServiceResolver(ctx context.Context) (*authorize.AuthorizeEditorDataAttributeResolversServiceResolverDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	valueRef, d := expandEditorReferenceData(ctx, p.ValueRef)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataAttributeResolversServiceResolverDTO(
		authorize.EnumAuthorizeEditorDataResolverDTOType(p.Type.ValueString()),
		*valueRef,
	)

	// Condition
	if !p.Condition.IsNull() && !p.Condition.IsUnknown() {

		condition, d := expandEditorDataCondition(ctx, p.Condition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetCondition(*condition)
	}

	// Name
	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	// Processor
	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {

		processor, d := expandEditorDataProcessor(ctx, p.Processor)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetProcessor(*processor)
	}

	return data, diags
}

func (p *editorDataResolverResourceModel) expandSystemResolver(ctx context.Context) (*authorize.AuthorizeEditorDataAttributeResolversSystemResolverDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.NewAuthorizeEditorDataAttributeResolversSystemResolverDTO(
		authorize.EnumAuthorizeEditorDataResolverDTOType(p.Type.ValueString()),
		authorize.EnumAuthorizeEditorDataAttributeResolversSystemResolverDTOValue(p.Value.ValueString()),
	)

	// Condition
	if !p.Condition.IsNull() && !p.Condition.IsUnknown() {

		condition, d := expandEditorDataCondition(ctx, p.Condition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetCondition(*condition)
	}

	// Name
	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	// Processor
	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {

		processor, d := expandEditorDataProcessor(ctx, p.Processor)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetProcessor(*processor)
	}

	return data, diags
}

func (p *editorDataResolverResourceModel) expandUserResolver(ctx context.Context) (*authorize.AuthorizeEditorDataAttributeResolversUserResolverDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	query, d := expandEditorResolverQuery(ctx, p.Query)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataAttributeResolversUserResolverDTO(
		authorize.EnumAuthorizeEditorDataResolverDTOType(p.Type.ValueString()),
		*query,
	)

	// Condition
	if !p.Condition.IsNull() && !p.Condition.IsUnknown() {

		condition, d := expandEditorDataCondition(ctx, p.Condition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetCondition(*condition)
	}

	// Name
	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	// Processor
	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {

		processor, d := expandEditorDataProcessor(ctx, p.Processor)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetProcessor(*processor)
	}

	return data, diags
}

func editorResolversOkToListTF(ctx context.Context, apiObject []authorize.AuthorizeEditorDataResolverDTO, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: editorDataResolverTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := editorDataResolverOkToTF(ctx, &v, true)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(tfObjType), diags
		}

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.ListValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func editorResolversOkToSetTF(ctx context.Context, apiObject []authorize.AuthorizeEditorDataResolverDTO, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: editorDataResolverTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := editorDataResolverOkToTF(ctx, &v, true)
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

func editorDataResolverOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataResolverDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorize.AuthorizeEditorDataResolverDTO{}) {
		return types.ObjectNull(editorDataResolverTFObjectTypes), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case *authorize.AuthorizeEditorDataAttributeResolversAttributeResolverDTO:

		conditionValue, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTF(ctx, conditionValue, ok)
		diags.Append(d...)

		processorValue, ok := t.GetProcessorOk()
		processor, d := editorDataProcessorOkToTF(ctx, processorValue, ok)
		diags.Append(d...)

		valueRef, d := editorDataReferenceObjectOkToTF(t.GetValueOk())
		diags.Append(d...)

		attributeMap["condition"] = condition
		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["processor"] = processor
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["value_ref"] = valueRef

	case *authorize.AuthorizeEditorDataAttributeResolversConstantResolverDTO:

		conditionValue, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTF(ctx, conditionValue, ok)
		diags.Append(d...)

		processorValue, ok := t.GetProcessorOk()
		processor, d := editorDataProcessorOkToTF(ctx, processorValue, ok)
		diags.Append(d...)

		valueType, d := editorValueTypeOkToTF(t.GetValueTypeOk())
		diags.Append(d...)

		attributeMap["condition"] = condition
		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["processor"] = processor
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["value_type"] = valueType
		attributeMap["value"] = framework.StringOkToTF(t.GetValueOk())

	case *authorize.AuthorizeEditorDataAttributeResolversCurrentRepetitionValueResolverDTO:

		conditionValue, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTF(ctx, conditionValue, ok)
		diags.Append(d...)

		processorValue, ok := t.GetProcessorOk()
		processor, d := editorDataProcessorOkToTF(ctx, processorValue, ok)
		diags.Append(d...)

		attributeMap["condition"] = condition
		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["processor"] = processor
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())

	case *authorize.AuthorizeEditorDataAttributeResolversCurrentUserIDResolverDTO:

		conditionValue, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTF(ctx, conditionValue, ok)
		diags.Append(d...)

		processorValue, ok := t.GetProcessorOk()
		processor, d := editorDataProcessorOkToTF(ctx, processorValue, ok)
		diags.Append(d...)

		attributeMap["condition"] = condition
		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["processor"] = processor
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())

	case *authorize.AuthorizeEditorDataAttributeResolversRequestResolverDTO:

		conditionValue, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTF(ctx, conditionValue, ok)
		diags.Append(d...)

		processorValue, ok := t.GetProcessorOk()
		processor, d := editorDataProcessorOkToTF(ctx, processorValue, ok)
		diags.Append(d...)

		attributeMap["condition"] = condition
		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["processor"] = processor
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())

	case *authorize.AuthorizeEditorDataAttributeResolversServiceResolverDTO:

		conditionValue, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTF(ctx, conditionValue, ok)
		diags.Append(d...)

		processorValue, ok := t.GetProcessorOk()
		processor, d := editorDataProcessorOkToTF(ctx, processorValue, ok)
		diags.Append(d...)

		valueRef, d := editorDataReferenceObjectOkToTF(t.GetValueOk())
		diags.Append(d...)

		attributeMap["condition"] = condition
		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["processor"] = processor
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["value_ref"] = valueRef

	case *authorize.AuthorizeEditorDataAttributeResolversSystemResolverDTO:

		conditionValue, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTF(ctx, conditionValue, ok)
		diags.Append(d...)

		processorValue, ok := t.GetProcessorOk()
		processor, d := editorDataProcessorOkToTF(ctx, processorValue, ok)
		diags.Append(d...)

		attributeMap["condition"] = condition
		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["processor"] = processor
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["value"] = framework.EnumOkToTF(t.GetValueOk())

	case *authorize.AuthorizeEditorDataAttributeResolversUserResolverDTO:

		conditionValue, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTF(ctx, conditionValue, ok)
		diags.Append(d...)

		processorValue, ok := t.GetProcessorOk()
		processor, d := editorDataProcessorOkToTF(ctx, processorValue, ok)
		diags.Append(d...)

		queryValue, ok := t.GetQueryOk()
		query, d := editorDataResolverQueryOkToTF(ctx, queryValue, ok)
		diags.Append(d...)

		attributeMap["condition"] = condition
		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["processor"] = processor
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["query"] = query

	default:
		tflog.Error(ctx, "Invalid resolver type", map[string]interface{}{
			"resolver type": t,
		})
		diags.AddError(
			"Invalid resolver type",
			"The resolver type is not supported.  Please raise an issue with the provider maintainers.",
		)
		return types.ObjectNull(editorDataResolverTFObjectTypes), diags
	}

	attributeMap = editorDataResolverConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(editorDataResolverTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorDataResolverConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
	nullMap := map[string]attr.Value{
		"condition":  types.ObjectNull(editorDataResolverTFObjectTypes),
		"name":       types.StringNull(),
		"processor":  types.ObjectNull(editorDataProcessorTFObjectTypes),
		"type":       types.StringNull(),
		"value_ref":  types.ObjectNull(editorReferenceObjectTFObjectTypes),
		"value_type": types.ObjectNull(editorValueTypeTFObjectTypes),
		"value":      types.StringNull(),
		"query":      types.ObjectNull(editorDataResolverQueryTFObjectTypes),
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}
