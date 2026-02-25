// Copyright © 2026 Ping Identity Corporation

//go:build beta

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
	"github.com/patrickcping/pingone-go-sdk-v2/authorizeeditor"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	objectvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

const conditionNestedIterationMaxDepth = 4

var leafConditionTypes = []authorizeeditor.EnumAuthorizeEditorDataConditionDTOType{
	"COMPARISON",
	"EMPTY",
	"REFERENCE",
}

func dataConditionObjectSchemaAttributes() (attributes map[string]schema.Attribute) {
	const initialIteration = 1
	return dataConditionObjectSchemaAttributesIteration(initialIteration)
}

func dataConditionObjectSchemaAttributesIteration(iteration int32) (attributes map[string]schema.Attribute) {

	supportedTypes := authorizeeditor.AllowedEnumAuthorizeEditorDataConditionDTOTypeEnumValues

	if iteration >= conditionNestedIterationMaxDepth {
		supportedTypes = leafConditionTypes
	}

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the authorization condition type.",
	).AllowedValuesEnum(supportedTypes)

	comparatorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the comparison operator used to evaluate the authorization condition.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON))).AllowedValuesEnum(authorizeeditor.AllowedEnumAuthorizeEditorDataConditionsComparisonConditionDTOComparatorEnumValues)

	leftDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings that apply to the left side of the authorization condition statement.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON)))

	rightDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings that apply to the right side of the authorization condition statement.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON)))

	conditionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of objects that specifies configuration settings for multiple authorization conditions to evaluate.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s` or `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_AND), string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_OR)))

	conditionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for a single authorization condition to evaluate.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_NOT)))

	referenceDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for the authorization condition reference to evaluate.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE)))

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
	if slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON) {
		attributes["comparator"] = schema.StringAttribute{
			Description:         comparatorDescription.Description,
			MarkdownDescription: comparatorDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorizeeditor.AllowedEnumAuthorizeEditorDataConditionsComparisonConditionDTOComparatorEnumValues)...),
				stringvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON)),
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
					types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: dataConditionComparandObjectLeftSchemaAttributes(),
		}

		attributes["right"] = schema.SingleNestedAttribute{
			Description:         rightDescription.Description,
			MarkdownDescription: rightDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: dataConditionComparandObjectRightSchemaAttributes(),
		}
	}

	// type == "AND", type == "OR"
	if slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_AND) ||
		slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_OR) {
		attributes["conditions"] = schema.SetNestedAttribute{
			Description:         conditionsDescription.Description,
			MarkdownDescription: conditionsDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Set{
				setvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_AND)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				setvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_OR)),
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
	if slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_NOT) {
		attributes["condition"] = schema.SingleNestedAttribute{
			Description:         conditionDescription.Description,
			MarkdownDescription: conditionDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_NOT)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: dataConditionObjectSchemaAttributesIteration(iteration + 1),
		}
	}

	// type == "REFERENCE"
	if slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE) {
		attributes["reference"] = schema.SingleNestedAttribute{
			Description:         referenceDescription.Description,
			MarkdownDescription: referenceDescription.MarkdownDescription,
			Optional:            true,

			Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the authorization condition reference in the trust framework.")),

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},
		}
	}

	return attributes
}

type editorDataConditionLeafResourceModel struct {
	Type       types.String `tfsdk:"type"`
	Comparator types.String `tfsdk:"comparator"`
	Left       types.Object `tfsdk:"left"`
	Right      types.Object `tfsdk:"right"`
	Reference  types.Object `tfsdk:"reference"`
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

func initializeEditorDataConditionTFObjectTypes(iteration int32) map[string]attr.Type {

	supportedTypes := authorizeeditor.AllowedEnumAuthorizeEditorDataConditionDTOTypeEnumValues

	if iteration >= conditionNestedIterationMaxDepth {
		supportedTypes = leafConditionTypes
	}

	attrMap := map[string]attr.Type{
		"type": types.StringType,
	}

	if slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON) {
		attrMap["comparator"] = types.StringType
		attrMap["left"] = types.ObjectType{AttrTypes: editorDataConditionComparandLeftTFObjectTypes}
		attrMap["right"] = types.ObjectType{AttrTypes: editorDataConditionComparandRightTFObjectTypes}
	}

	if slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_AND) ||
		slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_OR) {
		attrMap["conditions"] = types.SetType{
			ElemType: types.ObjectType{AttrTypes: initializeEditorDataConditionTFObjectTypes(iteration + 1)},
		}
	}

	if slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_NOT) {
		attrMap["condition"] = types.ObjectType{AttrTypes: initializeEditorDataConditionTFObjectTypes(iteration + 1)}
	}

	if slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE) {
		attrMap["reference"] = types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes}
	}

	return attrMap
}

func expandEditorDataCondition(ctx context.Context, condition basetypes.ObjectValue) (conditionObject *authorizeeditor.AuthorizeEditorDataConditionDTO, diags diag.Diagnostics) {
	const initialIteration = 1
	return expandEditorDataConditionIteration(ctx, condition, initialIteration)
}

func expandEditorDataConditionIteration(ctx context.Context, condition basetypes.ObjectValue, iteration int32) (conditionObject *authorizeeditor.AuthorizeEditorDataConditionDTO, diags diag.Diagnostics) {
	var d diag.Diagnostics

	leaf := iteration >= conditionNestedIterationMaxDepth

	if leaf {
		var plan *editorDataConditionLeafResourceModel
		diags.Append(condition.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		conditionObject, d = plan.expand(ctx)
		diags.Append(d...)
	} else {
		var plan *editorDataConditionResourceModel
		diags.Append(condition.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		conditionObject, d = plan.expand(ctx, iteration)
		diags.Append(d...)
	}
	if diags.HasError() {
		return nil, diags
	}

	return conditionObject, diags
}

func expandEditorDataConditions(ctx context.Context, conditions basetypes.SetValue) (conditionSet []authorizeeditor.AuthorizeEditorDataConditionDTO, diags diag.Diagnostics) {
	const initialIteration = 1
	return expandEditorDataConditionsIteration(ctx, conditions, initialIteration)
}

func expandEditorDataConditionsIteration(ctx context.Context, conditions basetypes.SetValue, iteration int32) (conditionSet []authorizeeditor.AuthorizeEditorDataConditionDTO, diags diag.Diagnostics) {

	leaf := iteration >= conditionNestedIterationMaxDepth

	if leaf {
		var plan []editorDataConditionLeafResourceModel
		diags.Append(conditions.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		conditionSet = make([]authorizeeditor.AuthorizeEditorDataConditionDTO, 0, len(plan))
		for _, conditionPlan := range plan {

			conditionObject, d := conditionPlan.expand(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			conditionSet = append(conditionSet, *conditionObject)
		}

	} else {

		var plan []editorDataConditionResourceModel
		diags.Append(conditions.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		conditionSet = make([]authorizeeditor.AuthorizeEditorDataConditionDTO, 0, len(plan))
		for _, conditionPlan := range plan {

			conditionObject, d := conditionPlan.expand(ctx, iteration)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			conditionSet = append(conditionSet, *conditionObject)
		}

	}

	return conditionSet, diags
}

func (p *editorDataConditionResourceModel) expand(ctx context.Context, iteration int32) (*authorizeeditor.AuthorizeEditorDataConditionDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorizeeditor.AuthorizeEditorDataConditionDTO{}

	switch authorizeeditor.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()) {
	case authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_AND:
		data.AuthorizeEditorDataConditionsAndConditionDTO, d = p.expandAndCondition(ctx, iteration)
		diags.Append(d...)
	case authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON:
		data.AuthorizeEditorDataConditionsComparisonConditionDTO, d = p.expandComparisonCondition(ctx)
		diags.Append(d...)
	case authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_EMPTY:
		data.AuthorizeEditorDataConditionsEmptyConditionDTO = p.expandEmptyCondition()
	case authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_NOT:
		data.AuthorizeEditorDataConditionsNotConditionDTO, d = p.expandNotCondition(ctx, iteration)
		diags.Append(d...)
	case authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_OR:
		data.AuthorizeEditorDataConditionsOrConditionDTO, d = p.expandOrCondition(ctx, iteration)
		diags.Append(d...)
	case authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE:
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

func (p *editorDataConditionLeafResourceModel) expand(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataConditionDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorizeeditor.AuthorizeEditorDataConditionDTO{}

	switch authorizeeditor.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()) {
	case authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_AND:
		diags.AddError(
			"Invalid leaf condition type",
			fmt.Sprintf("The leaf condition type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	case authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON:
		data.AuthorizeEditorDataConditionsComparisonConditionDTO, d = p.expandComparisonCondition(ctx)
		diags.Append(d...)
	case authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_EMPTY:
		data.AuthorizeEditorDataConditionsEmptyConditionDTO = p.expandEmptyCondition()
	case authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_NOT:
		diags.AddError(
			"Invalid leaf condition type",
			fmt.Sprintf("The leaf condition type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	case authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_OR:
		diags.AddError(
			"Invalid leaf condition type",
			fmt.Sprintf("The leaf condition type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	case authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE:
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

func (p *editorDataConditionResourceModel) expandAndCondition(ctx context.Context, iteration int32) (*authorizeeditor.AuthorizeEditorDataConditionsAndConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	conditions, d := expandEditorDataConditionsIteration(ctx, p.Conditions, iteration+1)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorizeeditor.NewAuthorizeEditorDataConditionsAndConditionDTO(
		authorizeeditor.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		conditions,
	)

	return data, diags
}

func (p *editorDataConditionResourceModel) expandComparisonCondition(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataConditionsComparisonConditionDTO, diag.Diagnostics) {
	return expandComparisonCondition(ctx, p.Left, p.Right, p.Comparator)
}

func (p *editorDataConditionLeafResourceModel) expandComparisonCondition(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataConditionsComparisonConditionDTO, diag.Diagnostics) {
	return expandComparisonCondition(ctx, p.Left, p.Right, p.Comparator)
}

func expandComparisonCondition(ctx context.Context, leftComparand, rightComparand basetypes.ObjectValue, comparator basetypes.StringValue) (*authorizeeditor.AuthorizeEditorDataConditionsComparisonConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	left, d := expandEditorDataConditionLeftComparand(ctx, leftComparand)
	diags.Append(d...)
	right, d := expandEditorDataConditionRightComparand(ctx, rightComparand)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	data := authorizeeditor.NewAuthorizeEditorDataConditionsComparisonConditionDTO(
		authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON,
		*left,
		*right,
		authorizeeditor.EnumAuthorizeEditorDataConditionsComparisonConditionDTOComparator(comparator.ValueString()),
	)

	return data, diags
}

func (p *editorDataConditionResourceModel) expandEmptyCondition() *authorizeeditor.AuthorizeEditorDataConditionsEmptyConditionDTO {
	return expandEmptyCondition()
}

func (p *editorDataConditionLeafResourceModel) expandEmptyCondition() *authorizeeditor.AuthorizeEditorDataConditionsEmptyConditionDTO {
	return expandEmptyCondition()
}

func expandEmptyCondition() *authorizeeditor.AuthorizeEditorDataConditionsEmptyConditionDTO {

	data := authorizeeditor.NewAuthorizeEditorDataConditionsEmptyConditionDTO(
		authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_EMPTY,
	)

	return data
}

func (p *editorDataConditionResourceModel) expandNotCondition(ctx context.Context, iteration int32) (*authorizeeditor.AuthorizeEditorDataConditionsNotConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	condition, d := expandEditorDataConditionIteration(ctx, p.Condition, iteration+1)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorizeeditor.NewAuthorizeEditorDataConditionsNotConditionDTO(
		authorizeeditor.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		*condition,
	)

	return data, diags
}

func (p *editorDataConditionResourceModel) expandOrCondition(ctx context.Context, iteration int32) (*authorizeeditor.AuthorizeEditorDataConditionsOrConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	conditions, d := expandEditorDataConditionsIteration(ctx, p.Conditions, iteration+1)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorizeeditor.NewAuthorizeEditorDataConditionsOrConditionDTO(
		authorizeeditor.EnumAuthorizeEditorDataConditionDTOType(p.Type.ValueString()),
		conditions,
	)

	return data, diags
}

func (p *editorDataConditionResourceModel) expandReferenceCondition(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataConditionsReferenceConditionDTO, diag.Diagnostics) {
	return expandReferenceCondition(ctx, p.Reference)
}

func (p *editorDataConditionLeafResourceModel) expandReferenceCondition(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataConditionsReferenceConditionDTO, diag.Diagnostics) {
	return expandReferenceCondition(ctx, p.Reference)
}

func expandReferenceCondition(ctx context.Context, referenceObj basetypes.ObjectValue) (*authorizeeditor.AuthorizeEditorDataConditionsReferenceConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	reference, d := expandEditorReferenceData(ctx, referenceObj)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorizeeditor.NewAuthorizeEditorDataConditionsReferenceConditionDTO(
		authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE,
		*reference,
	)

	return data, diags
}

func editorDataConditionsOkToSetTFIteration(ctx context.Context, iteration int32, apiObject []authorizeeditor.AuthorizeEditorDataConditionDTO, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: initializeEditorDataConditionTFObjectTypes(iteration)}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := editorDataConditionOkToTFIteration(ctx, iteration, &v, true)
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

func editorDataConditionOkToTF(ctx context.Context, apiObject *authorizeeditor.AuthorizeEditorDataConditionDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	const initialIteration = 1
	return editorDataConditionOkToTFIteration(ctx, initialIteration, apiObject, ok)
}

func editorDataConditionOkToTFIteration(ctx context.Context, iteration int32, apiObject *authorizeeditor.AuthorizeEditorDataConditionDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorizeeditor.AuthorizeEditorDataConditionDTO{}) {
		return types.ObjectNull(initializeEditorDataConditionTFObjectTypes(iteration)), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case *authorizeeditor.AuthorizeEditorDataConditionsAndConditionDTO:

		conditionsResp, ok := t.GetConditionsOk()
		conditions, d := editorDataConditionsOkToSetTFIteration(ctx, iteration+1, conditionsResp, ok)
		diags.Append(d...)

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["conditions"] = conditions

	case *authorizeeditor.AuthorizeEditorDataConditionsComparisonConditionDTO:

		leftResp, ok := t.GetLeftOk()
		left, d := editorDataConditionComparandLeftOkToTF(ctx, leftResp, ok)
		diags.Append(d...)

		rightResp, ok := t.GetRightOk()
		right, d := editorDataConditionComparandRightOkToTF(ctx, rightResp, ok)
		diags.Append(d...)

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["comparator"] = framework.EnumOkToTF(t.GetComparatorOk())
		attributeMap["left"] = left
		attributeMap["right"] = right

	case *authorizeeditor.AuthorizeEditorDataConditionsEmptyConditionDTO:

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())

	case *authorizeeditor.AuthorizeEditorDataConditionsNotConditionDTO:

		conditionResp, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTFIteration(ctx, iteration+1, conditionResp, ok)
		diags.Append(d...)

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["condition"] = condition

	case *authorizeeditor.AuthorizeEditorDataConditionsOrConditionDTO:

		conditionsResp, ok := t.GetConditionsOk()
		conditions, d := editorDataConditionsOkToSetTFIteration(ctx, iteration+1, conditionsResp, ok)
		diags.Append(d...)

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["conditions"] = conditions

	case *authorizeeditor.AuthorizeEditorDataConditionsReferenceConditionDTO:

		reference, d := editorDataReferenceObjectOkToTF(t.GetReferenceOk())
		diags.Append(d...)

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["reference"] = reference

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

	objValue, d := types.ObjectValue(initializeEditorDataConditionTFObjectTypes(iteration), attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorDataConditionConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value, iteration int32) map[string]attr.Value {

	supportedTypes := authorizeeditor.AllowedEnumAuthorizeEditorDataConditionDTOTypeEnumValues

	if iteration >= conditionNestedIterationMaxDepth {
		supportedTypes = leafConditionTypes
	}

	nullMap := map[string]attr.Value{
		"type": types.StringNull(),
	}

	if slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_COMPARISON) {
		nullMap["comparator"] = types.StringNull()
		nullMap["left"] = types.ObjectNull(editorDataConditionComparandLeftTFObjectTypes)
		nullMap["right"] = types.ObjectNull(editorDataConditionComparandRightTFObjectTypes)
	}

	if slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_AND) ||
		slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_OR) {
		nullMap["conditions"] = types.SetNull(types.ObjectType{AttrTypes: initializeEditorDataConditionTFObjectTypes(iteration + 1)})
	}

	if slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_NOT) {
		nullMap["condition"] = types.ObjectNull(initializeEditorDataConditionTFObjectTypes(iteration + 1))
	}

	if slices.Contains(supportedTypes, authorizeeditor.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_REFERENCE) {
		nullMap["reference"] = types.ObjectNull(editorReferenceObjectTFObjectTypes)
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}
