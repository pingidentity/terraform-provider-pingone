package authorize

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	var valueConflictingPathKeys = []string{
		"name",
		"description",
		"enabled",
		// "statements",
		"condition",
		"combining_algorithm",
		"repetition_settings",
	}

	if iteration < policyNestedIterationMaxDepth {
		valueConflictingPathKeys = append(valueConflictingPathKeys, "children")
	}

	valueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that defines a relationship to a child policy.",
	).ConflictsWith(valueConflictingPathKeys)

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a user-friendly name to apply to the authorization policy.  The value must be unique.",
	).AppendMarkdownString("Also requires `name` and `combining_algorithm`.").ConflictsWith([]string{"value"})

	descriptionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a description to apply to the policy.",
	).AppendMarkdownString("Also requires `name` and `combining_algorithm`.").ConflictsWith([]string{"value"})

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the policy is enabled, and whether the policy is evaluated.",
	).DefaultValue(true).AppendMarkdownString("Also requires `name` and `combining_algorithm`.").ConflictsWith([]string{"value"})

	conditionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for an authorization condition to apply to the policy.",
	).AppendMarkdownString("Also requires `name` and `combining_algorithm`.").ConflictsWith([]string{"value"})

	combiningAlgorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings that determine how rules are combined to produce an authorization decision.",
	).AppendMarkdownString("Also requires `name` and `combining_algorithm`.").ConflictsWith([]string{"value"})

	repetitionSettingsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings that appies the policy to each item of the specific attribute, filtered by decision.",
	).AppendMarkdownString("Also requires `name` and `combining_algorithm`.").ConflictsWith([]string{"value"})

	childrenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An ordered list of objects that specifies child policies or policy sets. Policies can either be specified by reference using the `value` field, or by inline definition.",
	).AppendMarkdownString("Also requires `name` and `combining_algorithm`.").ConflictsWith([]string{"value"})

	valueConflictsWithExpressions := make([]path.Expression, 0, len(valueConflictingPathKeys))

	for _, key := range valueConflictingPathKeys {
		valueConflictsWithExpressions = append(valueConflictsWithExpressions, path.MatchRelative().AtParent().AtName(key))
	}

	attributes = map[string]schema.Attribute{
		"value": schema.SingleNestedAttribute{
			Description:         valueDescription.Description,
			MarkdownDescription: valueDescription.MarkdownDescription,
			// Optional:            true,
			Computed: true,

			// Validators: []validator.Object{
			// 	objectvalidator.ConflictsWith(valueConflictsWithExpressions...),
			// },

			Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of a child policy.")),
		},

		"name": schema.StringAttribute{
			Description:         nameDescription.Description,
			MarkdownDescription: nameDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.String{
				stringvalidator.LengthAtLeast(attrMinLength),
				stringvalidator.ConflictsWith(
					path.MatchRelative().AtParent().AtName("value"),
				),
				stringvalidator.AlsoRequires(
					path.MatchRelative().AtParent().AtName("combining_algorithm"),
					path.MatchRelative().AtParent().AtName("name"),
				),
			},
		},

		"type": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the type of the policy.").Description,
			Computed:    true,

			Default: stringdefault.StaticString("POLICY"),
		},

		"description": schema.StringAttribute{
			Description:         descriptionDescription.Description,
			MarkdownDescription: descriptionDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.String{
				stringvalidator.ConflictsWith(
					path.MatchRelative().AtParent().AtName("value"),
				),
				stringvalidator.AlsoRequires(
					path.MatchRelative().AtParent().AtName("combining_algorithm"),
					path.MatchRelative().AtParent().AtName("name"),
				),
			},
		},

		"enabled": schema.BoolAttribute{
			Description:         enabledDescription.Description,
			MarkdownDescription: enabledDescription.MarkdownDescription,
			Optional:            true,
			Computed:            true,

			Default: booldefault.StaticBool(true),

			Validators: []validator.Bool{
				boolvalidator.ConflictsWith(
					path.MatchRelative().AtParent().AtName("value"),
				),
				boolvalidator.AlsoRequires(
					path.MatchRelative().AtParent().AtName("combining_algorithm"),
					path.MatchRelative().AtParent().AtName("name"),
				),
			},
		},

		// "statements": schema.ListNestedAttribute{
		// 	Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
		// 	Optional:    true,

		// 	NestedObject: schema.NestedAttributeObject{
		// 		Attributes: map[string]schema.Attribute{},
		// 	},

		// Validators: []validator.List{
		// 	listvalidator.ConflictsWith(
		// 		path.MatchRelative().AtParent().AtName("value"),
		// 	),
		// listvalidator.AlsoRequires(
		// 	path.MatchRelative().AtParent().AtName("combining_algorithm"),
		// 	path.MatchRelative().AtParent().AtName("name"),
		// ),
		// },
		// },

		"condition": schema.SingleNestedAttribute{
			Description:         conditionDescription.Description,
			MarkdownDescription: conditionDescription.MarkdownDescription,
			Optional:            true,

			Attributes: dataConditionObjectSchemaAttributes(),

			Validators: []validator.Object{
				objectvalidator.ConflictsWith(
					path.MatchRelative().AtParent().AtName("value"),
				),
				objectvalidator.AlsoRequires(
					path.MatchRelative().AtParent().AtName("combining_algorithm"),
					path.MatchRelative().AtParent().AtName("name"),
				),
			},
		},

		"combining_algorithm": schema.SingleNestedAttribute{
			Description:         combiningAlgorithmDescription.Description,
			MarkdownDescription: combiningAlgorithmDescription.MarkdownDescription,
			Optional:            true,

			Attributes: combiningAlgorithmObjectSchemaAttributes(),

			Validators: []validator.Object{
				objectvalidator.ConflictsWith(
					path.MatchRelative().AtParent().AtName("value"),
				),
				objectvalidator.AlsoRequires(
					path.MatchRelative().AtParent().AtName("combining_algorithm"),
					path.MatchRelative().AtParent().AtName("name"),
				),
			},
		},

		"repetition_settings": schema.SingleNestedAttribute{
			Description:         repetitionSettingsDescription.Description,
			MarkdownDescription: repetitionSettingsDescription.MarkdownDescription,
			Optional:            true,

			Attributes: repetitionSettingsObjectSchemaAttributes(),

			Validators: []validator.Object{
				objectvalidator.ConflictsWith(
					path.MatchRelative().AtParent().AtName("value"),
				),
				objectvalidator.AlsoRequires(
					path.MatchRelative().AtParent().AtName("combining_algorithm"),
					path.MatchRelative().AtParent().AtName("name"),
				),
			},
		},
	}

	if iteration < policyNestedIterationMaxDepth {
		attributes["children"] = schema.ListNestedAttribute{
			Description:         childrenDescription.Description,
			MarkdownDescription: childrenDescription.MarkdownDescription,
			Optional:            true,

			NestedObject: schema.NestedAttributeObject{
				Attributes: dataPolicyObjectSchemaAttributesIteration(iteration + 1),
			},

			Validators: []validator.List{
				listvalidator.ConflictsWith(
					path.MatchRelative().AtParent().AtName("value"),
				),
				listvalidator.AlsoRequires(
					path.MatchRelative().AtParent().AtName("combining_algorithm"),
					path.MatchRelative().AtParent().AtName("name"),
				),
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
	Value              types.Object `tfsdk:"value"`
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
	Value              types.Object `tfsdk:"value"`
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
		"value": types.ObjectType{
			AttrTypes: editorReferenceObjectTFObjectTypes,
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
	return expandChildPolicy(ctx, p.Name, p.Type, p.Description, p.Enabled, p.CombiningAlgorithm, p.Condition, p.RepetitionSettings, p.Value)
}

func (p *editorDataPolicyResourceModel) expand(ctx context.Context, iteration int32) (*authorize.AuthorizeEditorDataPoliciesPolicyChild, diag.Diagnostics) {
	var diags diag.Diagnostics

	data, d := expandChildPolicy(ctx, p.Name, p.Type, p.Description, p.Enabled, p.CombiningAlgorithm, p.Condition, p.RepetitionSettings, p.Value)
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

func expandChildPolicy(ctx context.Context, name, policyType, description basetypes.StringValue, enabled basetypes.BoolValue, combiningAlgorithm, condition, repetitionSettings, refValue basetypes.ObjectValue) (*authorize.AuthorizeEditorDataPoliciesPolicyChild, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.NewAuthorizeEditorDataPoliciesPolicyChild(
		policyType.ValueString(),
	)

	if !refValue.IsNull() && !refValue.IsUnknown() {
		refValueObj, d := expandEditorReferenceData(ctx, refValue)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		data.SetValue(*refValueObj)
	}

	if !name.IsNull() && !name.IsUnknown() {
		data.SetName(name.ValueString())
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

	if !combiningAlgorithm.IsNull() && !combiningAlgorithm.IsUnknown() {
		var plan *policyManagementPolicyCombiningAlgorithmResourceModel
		diags.Append(combiningAlgorithm.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		combiningAlgorithmExp := plan.expand()

		data.SetCombiningAlgorithm(*combiningAlgorithmExp)
	}

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

	value, d := editorDataReferenceObjectOkToTF(apiObject.GetValueOk())
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
		"value":               value,
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
