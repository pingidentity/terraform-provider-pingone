package authorize

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

const policyNestedIterationMaxDepth = 5

func dataPolicyObjectSchemaAttributes() (attributes map[string]schema.Attribute) {
	const initialIteration = 1
	return dataPolicyObjectSchemaAttributesIteration(initialIteration)
}

func dataPolicyObjectSchemaAttributesIteration(iteration int32) (attributes map[string]schema.Attribute) {

	const attrMinLength = 1

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the policy is enabled, and whether the policy is evaluated.",
	).DefaultValue(true)

	attributes = map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a user-friendly name to apply to the authorization policy.  The value must be unique.").Description,
			Required:    true,

			Validators: []validator.String{
				stringvalidator.LengthAtLeast(attrMinLength),
			},
		},

		"type": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the type of the policy.").Description,
			Optional:    true,
			Computed:    true,

			Default: stringdefault.StaticString("POLICY"),
		},

		"description": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a description to apply to the policy.").Description,
			Optional:    true,
		},

		"enabled": schema.BoolAttribute{
			Description:         enabledDescription.Description,
			MarkdownDescription: enabledDescription.MarkdownDescription,
			Optional:            true,
			Computed:            true,

			Default: booldefault.StaticBool(true),
		},

		// "statements": schema.ListNestedAttribute{
		// 	Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
		// 	Optional:    true,

		// 	NestedObject: schema.NestedAttributeObject{
		// 		Attributes: map[string]schema.Attribute{},
		// 	},
		// },

		"condition": schema.SingleNestedAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings for an authorization condition to apply to the policy.").Description,
			Optional:    true,

			Attributes: dataConditionObjectSchemaAttributes(),
		},

		"combining_algorithm": schema.SingleNestedAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings that determine how rules are combined to produce an authorization decision.").Description,
			Required:    true,

			Attributes: combiningAlgorithmObjectSchemaAttributes(),
		},

		"repetition_settings": schema.SingleNestedAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings that appies the policy to each item of the specific attribute, filtered by decision.").Description,
			Optional:    true,

			Attributes: repetitionSettingsObjectSchemaAttributes(),
		},
	}

	if iteration < policyNestedIterationMaxDepth {
		attributes["children"] = schema.ListNestedAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("An ordered list of objects that specifies child policies or policy sets.").Description,
			Optional:    true,

			NestedObject: schema.NestedAttributeObject{
				Attributes: dataPolicyObjectSchemaAttributesIteration(iteration + 1),
			},
		}
	}

	return attributes
}

type editorDataPolicyLeafResourceModel struct {
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	// Statements         types.List                   `tfsdk:"statements"`
	Condition          types.Object `tfsdk:"condition"`
	CombiningAlgorithm types.Object `tfsdk:"combining_algorithm"`
	RepetitionSettings types.Object `tfsdk:"repetition_settings"`
}

type editorDataPolicyResourceModel struct {
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	// Statements         types.List                   `tfsdk:"statements"`
	Condition          types.Object `tfsdk:"condition"`
	CombiningAlgorithm types.Object `tfsdk:"combining_algorithm"`
	Children           types.List   `tfsdk:"children"`
	RepetitionSettings types.Object `tfsdk:"repetition_settings"`
}

var editorDataPolicyTFObjectTypes = initializeEditorDataPolicyTFObjectTypes(1)

func initializeEditorDataPolicyTFObjectTypes(iteration int32) map[string]attr.Type {

	attrMap := map[string]attr.Type{
		"type":        types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
		"enabled":     types.BoolType,
		// "statements": types.BoolType,
		"condition": types.ObjectType{
			AttrTypes: editorDataConditionTFObjectTypes,
		},
		"combining_algorithm": types.ObjectType{
			AttrTypes: policyManagementPolicyCombiningAlgorithmTFObjectTypes,
		},
		"repetition_settings": types.ObjectType{
			AttrTypes: policyManagementPolicyRepetitionSettingsTFObjectTypes,
		},
	}

	if iteration < policyNestedIterationMaxDepth {
		attrMap["children"] = types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: initializeEditorDataPolicyTFObjectTypes(iteration + 1),
			},
		}
	}

	return attrMap
}

func expandEditorDataPolicyChildren(ctx context.Context, policyChildren basetypes.ListValue) (policyObjects []authorize.AuthorizeEditorDataPoliciesPolicyChild, diags diag.Diagnostics) {
	const initialIteration = 1
	return expandEditorDataPolicyChildrenIteration(ctx, policyChildren, initialIteration)
}

func expandEditorDataPolicyChildrenIteration(ctx context.Context, policyChildren basetypes.ListValue, iteration int32) (policyObjects []authorize.AuthorizeEditorDataPoliciesPolicyChild, diags diag.Diagnostics) {

	leaf := iteration >= policyNestedIterationMaxDepth

	returnPolicies := make([]authorize.AuthorizeEditorDataPoliciesPolicyChild, 0)

	if leaf {
		var plan []editorDataPolicyLeafResourceModel
		diags.Append(policyChildren.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return
		}

		for _, policyPlan := range plan {
			policyObject, d := policyPlan.expand(ctx)
			diags.Append(d...)
			if diags.HasError() {
				continue
			}
			returnPolicies = append(returnPolicies, *policyObject)
		}
	} else {
		var plan []editorDataPolicyResourceModel
		diags.Append(policyChildren.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return
		}

		for _, policyPlan := range plan {
			policyObject, d := policyPlan.expand(ctx, iteration)
			diags.Append(d...)
			if diags.HasError() {
				continue
			}
			returnPolicies = append(returnPolicies, *policyObject)
		}
	}

	policyObjects = returnPolicies

	return
}

func (p *editorDataPolicyLeafResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataPoliciesPolicyChild, diag.Diagnostics) {
	return expandChildPolicy(ctx, p.Name, p.Type, p.Description, p.Enabled, p.CombiningAlgorithm, p.Condition, p.RepetitionSettings)
}

func (p *editorDataPolicyResourceModel) expand(ctx context.Context, iteration int32) (*authorize.AuthorizeEditorDataPoliciesPolicyChild, diag.Diagnostics) {
	var diags diag.Diagnostics

	data, d := expandChildPolicy(ctx, p.Name, p.Type, p.Description, p.Enabled, p.CombiningAlgorithm, p.Condition, p.RepetitionSettings)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if !p.Children.IsNull() && !p.Children.IsUnknown() {
		children, d := expandEditorDataPolicyChildrenIteration(ctx, p.Children, iteration+1)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetChildren(children)
	}

	return data, diags
}

func expandChildPolicy(ctx context.Context, name, policyType, description basetypes.StringValue, enabled basetypes.BoolValue, combiningAlgorithm, condition, repetitionSettings basetypes.ObjectValue) (*authorize.AuthorizeEditorDataPoliciesPolicyChild, diag.Diagnostics) {
	var diags diag.Diagnostics

	var plan *policyManagementPolicyCombiningAlgorithmResourceModel
	diags.Append(combiningAlgorithm.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	combiningAlgorithmExp := plan.expand()

	data := authorize.NewAuthorizeEditorDataPoliciesPolicyChild(
		name.ValueString(),
		*combiningAlgorithmExp,
	)

	if !policyType.IsNull() && !policyType.IsUnknown() {
		data.SetType(policyType.ValueString())
	}

	if !description.IsNull() && !description.IsUnknown() {
		data.SetDescription(description.ValueString())
	}

	if !enabled.IsNull() && !enabled.IsUnknown() {
		data.SetEnabled(enabled.ValueBool())
	}

	// if !p.Statements.IsNull() && !p.Statements.IsUnknown() {
	// 	var plan []policyManagementPolicyStatementResourceModel
	// 	diags.Append(p.Statements.ElementsAs(ctx, &plan, false)...)
	// 	if diags.HasError() {
	// 		return nil, diags
	// 	}

	// 	statements := make([]map[string]interface{}, 0)
	// 	for _, planItem := range plan {
	// 		statements = append(statements, planItem.expand())
	// 	}

	// 	data.SetStatements(statements)
	// }

	if !condition.IsNull() && !condition.IsUnknown() {
		conditionExp, d := expandEditorDataCondition(ctx, condition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetCondition(*conditionExp)
	}

	if !repetitionSettings.IsNull() && !repetitionSettings.IsUnknown() {
		var plan *policyManagementPolicyRepetitionSettingsResourceModel
		diags.Append(repetitionSettings.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		repetitionSettingsExp, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetRepetitionSettings(*repetitionSettingsExp)
	}

	return data, diags
}

func editorDataPolicysOkToListTF(ctx context.Context, apiObject []authorize.AuthorizeEditorDataPoliciesPolicyChild, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	return editorDataPolicysOkToListTFIteration(ctx, 1, apiObject, ok)
}

func editorDataPolicysOkToListTFIteration(ctx context.Context, iteration int32, apiObject []authorize.AuthorizeEditorDataPoliciesPolicyChild, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: initializeEditorDataPolicyTFObjectTypes(iteration)}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := editorDataPolicyOkToTFIteration(ctx, iteration, &v, true)
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

func editorDataPolicyOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataPoliciesPolicyChild, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	const initialIteration = 1
	return editorDataPolicyOkToTFIteration(ctx, initialIteration, apiObject, ok)
}

func editorDataPolicyOkToTFIteration(ctx context.Context, iteration int32, apiObject *authorize.AuthorizeEditorDataPoliciesPolicyChild, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(initializeEditorDataPolicyTFObjectTypes(iteration)), diags
	}

	conditionVal, ok := apiObject.GetConditionOk()
	condition, d := editorDataConditionOkToTF(ctx, conditionVal, ok)
	diags.Append(d...)

	combiningAlgorithm, d := policyManagementPolicyCombiningAlgorithmOkToTF(apiObject.GetCombiningAlgorithmOk())
	diags.Append(d...)

	repetitionSettings, d := policyManagementPolicyRepetitionSettingsOkToTF(apiObject.GetRepetitionSettingsOk())
	diags.Append(d...)

	if diags.HasError() {
		return types.ObjectNull(initializeEditorDataPolicyTFObjectTypes(iteration)), diags
	}

	attrMap := map[string]attr.Value{
		"type":        framework.EnumOkToTF(apiObject.GetTypeOk()),
		"name":        framework.StringOkToTF(apiObject.GetNameOk()),
		"description": framework.StringOkToTF(apiObject.GetDescriptionOk()),
		"enabled":     framework.BoolOkToTF(apiObject.GetEnabledOk()),
		// "statements": framework.ListOkToTF(apiObject.GetStatementsOk()),
		"condition":           condition,
		"combining_algorithm": combiningAlgorithm,
		"repetition_settings": repetitionSettings,
	}

	if iteration < policyNestedIterationMaxDepth {
		childrenPolicies, ok := apiObject.GetChildrenOk()
		children, d := editorDataPolicysOkToListTFIteration(ctx, iteration+1, childrenPolicies, ok)
		diags.Append(d...)

		attrMap["children"] = children
	}

	objValue, d := types.ObjectValue(initializeEditorDataPolicyTFObjectTypes(iteration), attrMap)
	diags.Append(d...)

	return objValue, diags
}
