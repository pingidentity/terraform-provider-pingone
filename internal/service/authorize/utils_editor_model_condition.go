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
	listvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/listvalidator"
	objectvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

func dataConditionObjectSchemaAttributes() (attributes map[string]schema.Attribute) {
	const initialIteration = 1
	return dataConditionObjectSchemaAttributesIteration(initialIteration)
}
func dataConditionObjectSchemaAttributesIteration(iteration int32) (attributes map[string]schema.Attribute) {

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the resource's condition type.",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataConditionDTOTypeEnumValues)

	conditionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The list of conditions to apply in the given order.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_CHAIN)))

	predicateDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER)))

	conditionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_TRANSFORM)))

	expressionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`, `%s` or `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH)))

	valueTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`, `%s` or `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH)))

	conditionRefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_REFERENCE)))

	if iteration == 10 {
		attributes = map[string]schema.Attribute{}
		return
	}

	attributes = map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A user-friendly authorization condition name. The value must be unique.").Description,
			Required:    true,
		},

		"type": schema.StringAttribute{
			Description:         typeDescription.Description,
			MarkdownDescription: typeDescription.MarkdownDescription,
			Required:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataConditionDTOTypeEnumValues)...),
			},
		},

		// type == "CHAIN"
		"conditions": schema.ListNestedAttribute{
			Description:         conditionsDescription.Description,
			MarkdownDescription: conditionsDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.List{
				listvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_CHAIN)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			NestedObject: schema.NestedAttributeObject{
				Attributes: dataConditionObjectSchemaAttributesIteration(iteration + 1),
			},
		},

		// type == "COLLECTION_FILTER"
		"predicate": schema.SingleNestedAttribute{
			Description:         predicateDescription.Description,
			MarkdownDescription: predicateDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: dataConditionObjectSchemaAttributesIteration(iteration + 1),
		},

		// type == "COLLECTION_TRANSFORM"
		"condition": schema.SingleNestedAttribute{
			Description:         conditionDescription.Description,
			MarkdownDescription: conditionDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_TRANSFORM)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: dataConditionObjectSchemaAttributesIteration(iteration + 1),
		},

		// type == "JSON_PATH", type == "SPEL", type == "XPATH"
		"expression": schema.StringAttribute{
			Description:         expressionDescription.Description,
			MarkdownDescription: expressionDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.String{
				stringvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				stringvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				stringvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},
		},

		"value_type": schema.SingleNestedAttribute{
			Description: valueTypeDescription.Description,
			Optional:    true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: valueTypeObjectSchemaAttributes(),
		},

		// type == "REFERENCE"
		"condition_ref": schema.SingleNestedAttribute{
			Description:         conditionRefDescription.Description,
			MarkdownDescription: conditionRefDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_REFERENCE)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Required:    true,

					CustomType: pingonetypes.ResourceIDType{},
				},
			},
		},
	}

	return
}

type editorDataConditionResourceModel struct {
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	Conditions   types.List   `tfsdk:"conditions"`
	Predicate    types.Object `tfsdk:"predicate"`
	Condition    types.Object `tfsdk:"condition"`
	Expression   types.String `tfsdk:"expression"`
	ValueType    types.Object `tfsdk:"value_type"`
	ConditionRef types.Object `tfsdk:"condition_ref"`
}

var editorDataConditionTFObjectTypes = initializeEditorDataConditionTFObjectTypes()

func initializeEditorDataConditionTFObjectTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"type": types.StringType,
		"conditions": types.ListType{
			ElemType: types.ObjectType{AttrTypes: nil}, // Temporarily set to nil
		},
		"predicate":     types.ObjectType{AttrTypes: nil}, // Temporarily set to nil
		"condition":     types.ObjectType{AttrTypes: nil}, // Temporarily set to nil
		"expression":    types.StringType,
		"value_type":    types.ObjectType{AttrTypes: editorValueTypeTFObjectTypes},
		"condition_ref": types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
	}
}

func init() {
	// Now set the correct AttrTypes to break the initialization cycle
	editorDataConditionTFObjectTypes["conditions"] = types.ListType{
		ElemType: types.ObjectType{AttrTypes: editorDataConditionTFObjectTypes},
	}
	editorDataConditionTFObjectTypes["predicate"] = types.ObjectType{AttrTypes: editorDataConditionTFObjectTypes}
	editorDataConditionTFObjectTypes["condition"] = types.ObjectType{AttrTypes: editorDataConditionTFObjectTypes}
}

func expandEditorDataCondition(ctx context.Context, condition basetypes.ObjectValue) (conditionObject *authorize.AuthorizeEditorDataConditionDTO, diags diag.Diagnostics) {
	var plan *editorDataConditionResourceModel
	diags.Append(condition.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return
	}

	conditionObject, d := plan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	return
}

func (p *editorDataConditionResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataConditionDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorize.AuthorizeEditorDataConditionDTO{}

	switch authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_CHAIN:
		data.AuthorizeEditorDataConditionsChainConditionDTO, d = p.expandChainCondition(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER:
		data.AuthorizeEditorDataConditionsCollectionFilterConditionDTO, d = p.expandCollectionFilterCondition(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_TRANSFORM:
		data.AuthorizeEditorDataConditionsCollectionTransformConditionDTO, d = p.expandCollectionTransformCondition(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH:
		data.AuthorizeEditorDataConditionsJsonPathConditionDTO, d = p.expandJsonPathCondition(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_REFERENCE:
		data.AuthorizeEditorDataConditionsReferenceConditionDTO, d = p.expandReferenceCondition(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL:
		data.AuthorizeEditorDataConditionsSpelConditionDTO, d = p.expandSPELCondition(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH:
		data.AuthorizeEditorDataConditionsXPathConditionDTO, d = p.expandXPATHCondition(ctx)
		diags.Append(d...)
	default:
		diags.AddError(
			"Invalid condition type",
			fmt.Sprintf("The condition type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &data, diags
}

func (p *editorDataConditionResourceModel) expandChainCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsChainConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	var plan []editorDataConditionResourceModel
	diags.Append(p.Conditions.ElementsAs(ctx, &plan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	conditions := make([]authorize.AuthorizeEditorDataConditionDTO, 0, len(plan))
	for _, conditionPlan := range plan {

		conditionObject, d := conditionPlan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		conditions = append(conditions, *conditionObject)
	}

	data := authorize.NewAuthorizeEditorDataConditionsChainConditionDTO(
		p.Name.ValueString(),
		authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		conditions,
	)

	return data, diags
}

func (p *editorDataConditionResourceModel) expandCollectionFilterCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsCollectionFilterConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	predicate, d := expandEditorDataCondition(ctx, p.Condition)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataConditionsCollectionFilterConditionDTO(
		p.Name.ValueString(),
		authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		*predicate,
	)

	return data, diags
}

func (p *editorDataConditionResourceModel) expandCollectionTransformCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsCollectionTransformConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	condition, d := expandEditorDataCondition(ctx, p.Condition)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataConditionsCollectionTransformConditionDTO(
		p.Name.ValueString(),
		authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		*condition,
	)

	return data, diags
}

func (p *editorDataConditionResourceModel) expandJsonPathCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsJsonPathConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	valueType, d := expandEditorValueType(ctx, p.ValueType)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataConditionsJsonPathConditionDTO(
		p.Name.ValueString(),
		authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		p.Expression.ValueString(),
		*valueType,
	)

	return data, diags
}

func (p *editorDataConditionResourceModel) expandReferenceCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsReferenceConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	conditionRef, d := expandEditorReferenceData(ctx, p.ConditionRef)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataConditionsReferenceConditionDTO(
		p.Name.ValueString(),
		authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		*conditionRef,
	)

	return data, diags
}

func (p *editorDataConditionResourceModel) expandSPELCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsSpelConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	valueType, d := expandEditorValueType(ctx, p.ValueType)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataConditionsSpelConditionDTO(
		p.Name.ValueString(),
		authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		p.Expression.ValueString(),
		*valueType,
	)

	return data, diags
}

func (p *editorDataConditionResourceModel) expandXPATHCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsXPathConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	valueType, d := expandEditorValueType(ctx, p.ValueType)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataConditionsXPathConditionDTO(
		p.Name.ValueString(),
		authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		p.Expression.ValueString(),
		*valueType,
	)

	return data, diags
}

func editorDataConditionOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataConditionDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorize.AuthorizeEditorDataConditionDTO{}) {
		return types.ObjectNull(editorDataConditionTFObjectTypes), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case authorize.AuthorizeEditorDataConditionsChainConditionDTO:

		attributeMap = map[string]attr.Value{
			"name":       framework.StringOkToTF(t.GetNameOk()),
			"type":       framework.EnumOkToTF(t.GetTypeOk()),
			"conditions": nil,
		}

	case authorize.AuthorizeEditorDataConditionsCollectionFilterConditionDTO:

		predicateResp, ok := t.GetPredicateOk()
		predicate, d := editorDataConditionOkToTF(ctx, predicateResp, ok)
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"name":      framework.StringOkToTF(t.GetNameOk()),
			"type":      framework.EnumOkToTF(t.GetTypeOk()),
			"predicate": predicate,
		}

	case authorize.AuthorizeEditorDataConditionsCollectionTransformConditionDTO:

		conditionResp, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTF(ctx, conditionResp, ok)
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"name":      framework.StringOkToTF(t.GetNameOk()),
			"type":      framework.EnumOkToTF(t.GetTypeOk()),
			"condition": condition,
		}

	case authorize.AuthorizeEditorDataConditionsJsonPathConditionDTO:

		valueType, d := editorValueTypeOkToTF(t.GetValueTypeOk())
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"name":       framework.StringOkToTF(t.GetNameOk()),
			"type":       framework.EnumOkToTF(t.GetTypeOk()),
			"expression": framework.StringOkToTF(t.GetNameOk()),
			"value_type": valueType,
		}

	case authorize.AuthorizeEditorDataConditionsReferenceConditionDTO:

		conditionRef, d := editorDataReferenceObjectOkToTF(t.GetConditionOk())
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"name":          framework.StringOkToTF(t.GetNameOk()),
			"type":          framework.EnumOkToTF(t.GetTypeOk()),
			"condition_ref": conditionRef,
		}

	case authorize.AuthorizeEditorDataConditionsSpelConditionDTO:

		valueType, d := editorValueTypeOkToTF(t.GetValueTypeOk())
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"name":       framework.StringOkToTF(t.GetNameOk()),
			"type":       framework.EnumOkToTF(t.GetTypeOk()),
			"expression": framework.StringOkToTF(t.GetNameOk()),
			"value_type": valueType,
		}

	case authorize.AuthorizeEditorDataConditionsXPathConditionDTO:

		valueType, d := editorValueTypeOkToTF(t.GetValueTypeOk())
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"name":       framework.StringOkToTF(t.GetNameOk()),
			"type":       framework.EnumOkToTF(t.GetTypeOk()),
			"expression": framework.StringOkToTF(t.GetNameOk()),
			"value_type": valueType,
		}

	default:
		tflog.Error(ctx, "Invalid condition type", map[string]interface{}{
			"condition type": t,
		})
		diags.AddError(
			"Invalid condition type",
			"The condition type is not supported.  Please raise an issue with the provider maintainers.",
		)
	}

	attributeMap = editorDataConditionConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(editorDataConditionTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorDataConditionConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
	nullMap := map[string]attr.Value{
		"name":          types.StringNull(),
		"type":          types.StringNull(),
		"conditions":    types.ListNull(types.ObjectType{AttrTypes: editorDataConditionTFObjectTypes}),
		"predicate":     types.ObjectNull(editorDataConditionTFObjectTypes),
		"condition":     types.ObjectNull(editorDataConditionTFObjectTypes),
		"expression":    types.StringNull(),
		"value_type":    types.ObjectNull(editorValueTypeTFObjectTypes),
		"condition_ref": types.ObjectNull(editorReferenceObjectTFObjectTypes),
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}
