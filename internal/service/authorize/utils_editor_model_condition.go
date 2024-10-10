package authorize

import (
	"context"
	"fmt"
	"slices"

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
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

const conditionNestedIterationMaxDepth = 2

var leafConditionTypes = []authorize.EnumAuthorizeEditorDataConditionDTOType{
	"COMPARISON",
	"EMPTY",
	"REFERENCE",
}

func dataConditionObjectSchemaAttributes() (attributes map[string]schema.Attribute) {
	const initialIteration = 1
	return dataConditionObjectSchemaAttributesIteration(initialIteration)
}

func dataConditionObjectSchemaAttributesIteration(iteration int32) (attributes map[string]schema.Attribute) {

	supportedTypes := authorize.AllowedEnumAuthorizeEditorDataConditionDTOTypeEnumValues

	if iteration >= conditionNestedIterationMaxDepth {
		supportedTypes = leafConditionTypes
	}

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the authorization condition type.",
	).AllowedValuesEnum(supportedTypes)

	comparatorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the comparison operator used to evaluate the authorization condition.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON))).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataConditionsComparisonConditionDTOComparatorEnumValues)

	leftDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings that apply to the left side of the authorization condition statement.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON)))

	rightDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings that apply to the right side of the authorization condition statement.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON)))

	conditionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of objects that specifies configuration settings for multiple authorization conditions to evaluate.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s` or `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_AND), string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_OR)))

	conditionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for a single authorization condition to evaluate.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_NOT)))

	referenceDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for the authorization condition reference to evaluate.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE)))

	attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description:         typeDescription.Description,
			MarkdownDescription: typeDescription.MarkdownDescription,
			Required:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(supportedTypes)...),
			},
		},
	}

	// type == "COMPARISON"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON) {
		attributes["comparator"] = schema.StringAttribute{
			Description:         comparatorDescription.Description,
			MarkdownDescription: comparatorDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataConditionsComparisonConditionDTOComparatorEnumValues)...),
				stringvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},
		}

		attributes["left"] = schema.SingleNestedAttribute{
			Description:         leftDescription.Description,
			MarkdownDescription: leftDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: dataConditionComparandObjectSchemaAttributes(),
		}

		attributes["right"] = schema.SingleNestedAttribute{
			Description:         rightDescription.Description,
			MarkdownDescription: rightDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: dataConditionComparandObjectSchemaAttributes(),
		}
	}

	// type == "AND", type == "OR"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_AND) ||
		slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_OR) {
		attributes["conditions"] = schema.SetNestedAttribute{
			Description:         conditionsDescription.Description,
			MarkdownDescription: conditionsDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Set{
				setvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_AND)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				setvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_OR)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			NestedObject: schema.NestedAttributeObject{
				Attributes: dataConditionObjectSchemaAttributesIteration(iteration + 1),
			},
		}
	}

	// type == "EMPTY"
	// (same as base object)

	// type == "NOT"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_NOT) {
		attributes["condition"] = schema.SingleNestedAttribute{
			Description:         conditionDescription.Description,
			MarkdownDescription: conditionDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_NOT)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: dataConditionObjectSchemaAttributesIteration(iteration + 1),
		}
	}

	// type == "REFERENCE"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE) {
		attributes["reference"] = schema.SingleNestedAttribute{
			Description:         referenceDescription.Description,
			MarkdownDescription: referenceDescription.MarkdownDescription,
			Optional:            true,

			Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the authorization condition reference in the trust framework.")),

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},
		}
	}

	return attributes
}

type editorDataConditionResourceModel struct {
	Type       types.String `tfsdk:"type"`
	Comparator types.String `tfsdk:"comparator"`
	Left       types.Object `tfsdk:"left"`
	Right      types.Object `tfsdk:"right"`
	Conditions types.Set    `tfsdk:"conditions"`
	Condition  types.Object `tfsdk:"condition"`
	Reference  types.Object `tfsdk:"reference"`
}

var editorDataConditionTFObjectTypes = initializeEditorDataConditionTFObjectTypes(1)

func initializeEditorDataConditionTFObjectTypes(iteration int) map[string]attr.Type {

	supportedTypes := authorize.AllowedEnumAuthorizeEditorDataConditionDTOTypeEnumValues

	if iteration >= conditionNestedIterationMaxDepth {
		supportedTypes = leafConditionTypes
	}

	attrMap := map[string]attr.Type{
		"type": types.StringType,
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON) {
		attrMap["comparator"] = types.StringType
		attrMap["left"] = types.ObjectType{AttrTypes: editorDataConditionComparandTFObjectTypes}
		attrMap["right"] = types.ObjectType{AttrTypes: editorDataConditionComparandTFObjectTypes}
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_AND) ||
		slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_OR) {
		attrMap["conditions"] = types.ListType{
			ElemType: types.ObjectType{AttrTypes: initializeEditorDataConditionTFObjectTypes(iteration + 1)},
		}
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_NOT) {
		attrMap["condition"] = types.ObjectType{AttrTypes: initializeEditorDataConditionTFObjectTypes(iteration + 1)}
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE) {
		attrMap["reference"] = types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes}
	}

	return attrMap
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
	case authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_AND:
		data.AuthorizeEditorDataConditionsAndConditionDTO, d = p.expandAndCondition(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON:
		data.AuthorizeEditorDataConditionsComparisonConditionDTO, d = p.expandComparisonCondition(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_EMPTY:
		data.AuthorizeEditorDataConditionsEmptyConditionDTO = p.expandEmptyCondition()
	case authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_NOT:
		data.AuthorizeEditorDataConditionsNotConditionDTO, d = p.expandNotCondition(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_OR:
		data.AuthorizeEditorDataConditionsOrConditionDTO = p.expandOrCondition( /*ctx*/ )
		// diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE:
		data.AuthorizeEditorDataConditionsReferenceConditionDTO, d = p.expandReferenceCondition(ctx)
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

func (p *editorDataConditionResourceModel) expandAndCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsAndConditionDTO, diag.Diagnostics) {
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

	data := authorize.NewAuthorizeEditorDataConditionsAndConditionDTO(
		authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		conditions,
	)

	return data, diags
}

func (p *editorDataConditionResourceModel) expandComparisonCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsComparisonConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	left, d := expandEditorDataConditionComparand(ctx, p.Left)
	diags.Append(d...)
	right, d := expandEditorDataConditionComparand(ctx, p.Right)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataConditionsComparisonConditionDTO(
		authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		*left,
		*right,
		authorize.EnumAuthorizeEditorDataConditionsComparisonConditionDTOComparator(p.Comparator.ValueString()),
	)

	return data, diags
}

func (p *editorDataConditionResourceModel) expandEmptyCondition() *authorize.AuthorizeEditorDataConditionsEmptyConditionDTO {

	data := authorize.NewAuthorizeEditorDataConditionsEmptyConditionDTO(
		authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
	)

	return data
}

func (p *editorDataConditionResourceModel) expandNotCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsNotConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	condition, d := expandEditorDataCondition(ctx, p.Condition)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataConditionsNotConditionDTO(
		authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		*condition,
	)

	return data, diags
}

func (p *editorDataConditionResourceModel) expandOrCondition( /*ctx context.Context*/ ) *authorize.AuthorizeEditorDataConditionsOrConditionDTO { //, diag.Diagnostics) {
	// var diags diag.Diagnostics

	// var plan []editorDataConditionResourceModel
	// diags.Append(p.Conditions.ElementsAs(ctx, &plan, false)...)
	// if diags.HasError() {
	// 	return nil, diags
	// }

	// conditions := make([]authorize.AuthorizeEditorDataConditionDTO, 0, len(plan))
	// for _, conditionPlan := range plan {

	// 	conditionObject, d := conditionPlan.expand(ctx)
	// 	diags.Append(d...)
	// 	if diags.HasError() {
	// 		return nil, diags
	// 	}

	// 	conditions = append(conditions, *conditionObject)
	// }

	data := authorize.NewAuthorizeEditorDataConditionsOrConditionDTO(
		authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		// conditions,
	)

	return data
}

func (p *editorDataConditionResourceModel) expandReferenceCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsReferenceConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	reference, d := expandEditorReferenceData(ctx, p.Reference)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataConditionsReferenceConditionDTO(
		authorize.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		*reference,
	)

	return data, diags
}

func editorDataConditionsOkToSetTF(ctx context.Context, apiObject []authorize.AuthorizeEditorDataConditionDTO, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: editorDataConditionTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := editorDataConditionOkToTF(ctx, &v, true)
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

func editorDataConditionOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataConditionDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	const initialIteration = 1
	return editorDataConditionOkToTFIteration(ctx, initialIteration, apiObject, ok)
}

func editorDataConditionOkToTFIteration(ctx context.Context, iteration int, apiObject *authorize.AuthorizeEditorDataConditionDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorize.AuthorizeEditorDataConditionDTO{}) {
		return types.ObjectNull(initializeEditorDataConditionTFObjectTypes(iteration)), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case *authorize.AuthorizeEditorDataConditionsAndConditionDTO:

		conditionsResp, ok := t.GetConditionsOk()
		conditions, d := editorDataConditionsOkToSetTF(ctx, conditionsResp, ok)
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"type":       framework.EnumOkToTF(t.GetTypeOk()),
			"conditions": conditions,
		}

	case *authorize.AuthorizeEditorDataConditionsComparisonConditionDTO:

		leftResp, ok := t.GetLeftOk()
		left, d := editorDataConditionComparandOkToTF(ctx, leftResp, ok)
		diags.Append(d...)

		rightResp, ok := t.GetRightOk()
		right, d := editorDataConditionComparandOkToTF(ctx, rightResp, ok)
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"type":       framework.EnumOkToTF(t.GetTypeOk()),
			"comparator": framework.EnumOkToTF(t.GetComparatorOk()),
			"left":       left,
			"right":      right,
		}

	case *authorize.AuthorizeEditorDataConditionsEmptyConditionDTO:

		attributeMap = map[string]attr.Value{
			"type": framework.EnumOkToTF(t.GetTypeOk()),
		}

	case *authorize.AuthorizeEditorDataConditionsNotConditionDTO:

		conditionResp, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTFIteration(ctx, iteration+1, conditionResp, ok)
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"type":      framework.EnumOkToTF(t.GetTypeOk()),
			"condition": condition,
		}

	case *authorize.AuthorizeEditorDataConditionsOrConditionDTO:

		// conditionsResp, ok := t.GetConditionsOk()
		// conditions, d := editorDataConditionsOkToSetTF(ctx, conditionsResp, ok)
		// diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"type": framework.EnumOkToTF(t.GetTypeOk()),
			// "conditions": conditions,
		}

	case *authorize.AuthorizeEditorDataConditionsReferenceConditionDTO:

		reference, d := editorDataReferenceObjectOkToTF(t.GetReferenceOk())
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"type":      framework.EnumOkToTF(t.GetTypeOk()),
			"reference": reference,
		}

	default:
		tflog.Error(ctx, "Invalid condition type", map[string]interface{}{
			"condition type": t,
		})
		diags.AddError(
			"Invalid condition type",
			"The condition type is not supported.  Please raise an issue with the provider maintainers.",
		)
		return types.ObjectNull(initializeEditorDataConditionTFObjectTypes(iteration)), diags
	}

	attributeMap = editorDataConditionConvertEmptyValuesToTFNulls(attributeMap, iteration)

	objValue, d := types.ObjectValue(editorDataConditionTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorDataConditionConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value, iteration int) map[string]attr.Value {

	supportedTypes := authorize.AllowedEnumAuthorizeEditorDataConditionDTOTypeEnumValues

	if iteration >= conditionNestedIterationMaxDepth {
		supportedTypes = leafConditionTypes
	}

	nullMap := map[string]attr.Value{
		"type": types.StringNull(),
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON) {
		nullMap["comparator"] = types.StringNull()
		nullMap["left"] = types.ObjectNull(editorDataConditionComparandTFObjectTypes)
		nullMap["right"] = types.ObjectNull(editorDataConditionComparandTFObjectTypes)
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_AND) ||
		slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_OR) {
		nullMap["conditions"] = types.ListNull(types.ObjectType{AttrTypes: initializeEditorDataConditionTFObjectTypes(iteration + 1)})
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_NOT) {
		nullMap["condition"] = types.ObjectNull(initializeEditorDataConditionTFObjectTypes(iteration + 1))
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE) {
		nullMap["reference"] = types.ObjectNull(editorReferenceObjectTFObjectTypes)
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}
