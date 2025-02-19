package authorize

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	listvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/listvalidator"
	objectvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

const policyChildNestedIterationMaxDepth = 5

func dataPolicyChildObjectSchemaAttributes() (attributes map[string]schema.Attribute) {
	const initialIteration = 1
	return dataPolicyChildObjectSchemaAttributesIteration(initialIteration)
}

func dataPolicyChildObjectSchemaAttributesIteration(iteration int32) (attributes map[string]schema.Attribute) {

	const attrMinLength = 1
	var valueConflictingPathKeys = []string{
		"name",
		"description",
		"enabled",
		// "statements",
		"condition",
		"combining_algorithm",
		"repetition_settings",
		"effect_settings",
	}

	if iteration < policyChildNestedIterationMaxDepth {
		valueConflictingPathKeys = append(valueConflictingPathKeys, "children")
	}

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of the policy child.",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataPoliciesPolicyChildCommonTypeEnumValues)

	valueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that defines a relationship to a child policy.",
	).ConflictsWith(valueConflictingPathKeys)

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a user-friendly name to apply to the authorization policy.  The value must be unique.",
	).AppendMarkdownString(fmt.Sprintf("Required when `type` is `%s` or `%s` and `value` is not set.", string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_POLICY), string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE)))

	descriptionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a description to apply to the policy.",
	).AppendMarkdownString(fmt.Sprintf("Optional when `type` is `%s` or `%s` and `value` is not set.", string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_POLICY), string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE))).AppendMarkdownString("Also requires `name`.")

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the policy is enabled, and whether the policy is evaluated.",
	).AppendMarkdownString(fmt.Sprintf("Optional when `type` is `%s` or `%s` and `value` is not set.", string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_POLICY), string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE))).DefaultValue(true).AppendMarkdownString("Also requires `name`.")

	conditionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for an authorization condition to apply to the policy.",
	).AppendMarkdownString(fmt.Sprintf("Optional when `type` is `%s` or `%s` and `value` is not set.", string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_POLICY), string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE))).AppendMarkdownString("Also requires `name`.")

	combiningAlgorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings that determine how rules are combined to produce an authorization decision.",
	).AppendMarkdownString(fmt.Sprintf("Required when `type` is `%s` and cannot be set when `type` is `%s` or `value` is set.", string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_POLICY), string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE))).AppendMarkdownString("Also requires `name`.")

	repetitionSettingsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings that appies the policy to each item of the specific attribute, filtered by decision.",
	).AppendMarkdownString(fmt.Sprintf("Optional when `type` is `%s` and cannot be set when `type` is `%s` or `value` is set.", string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_POLICY), string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE))).AppendMarkdownString("Also requires `name` and `combining_algorithm`.")

	effectSettingsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings that determine how child rules are combined to produce an outcome for the policy.",
	).AppendMarkdownString(fmt.Sprintf("Required when `type` is `%s` and cannot be set when `type` is `%s` or `value` is set.", string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE), string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_POLICY))).AppendMarkdownString("Also requires `name`.")

	childrenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An ordered list of objects that specifies child policies or policy sets. Policies can either be specified by reference using the `value` field, or by inline definition.",
	).AppendMarkdownString(fmt.Sprintf("Required when `type` is `%s` and cannot be set when `type` is `%s` or `value` is set.", string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_POLICY), string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE))).AppendMarkdownString("Also requires `name` and `combining_algorithm`.")

	// valueConflictsWithExpressions := make([]path.Expression, 0, len(valueConflictingPathKeys))

	// for _, key := range valueConflictingPathKeys {
	// 	valueConflictsWithExpressions = append(valueConflictsWithExpressions, path.MatchRelative().AtParent().AtName(key))
	// }

	attributes = map[string]schema.Attribute{
		"value": schema.SingleNestedAttribute{
			Description:         valueDescription.Description,
			MarkdownDescription: valueDescription.MarkdownDescription,
			// Optional:            true,
			Computed: true,

			// Validators: []validator.Object{
			// 	objectvalidator.ConflictsWith(valueConflictsWithExpressions...),
			// },

			Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of a child policyChild.")),
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
				stringvalidator.Any(
					stringvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("combining_algorithm"),
					),
					stringvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("effect_settings"),
					),
				),
			},
		},

		"type": schema.StringAttribute{
			Description:         typeDescription.Description,
			MarkdownDescription: typeDescription.MarkdownDescription,
			Required:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataPoliciesPolicyChildCommonTypeEnumValues)...),
			},
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
					path.MatchRelative().AtParent().AtName("name"),
				),
				stringvalidator.Any(
					stringvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("combining_algorithm"),
					),
					stringvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("effect_settings"),
					),
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
					path.MatchRelative().AtParent().AtName("name"),
				),
				boolvalidator.Any(
					boolvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("combining_algorithm"),
					),
					boolvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("effect_settings"),
					),
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
					path.MatchRelative().AtParent().AtName("name"),
				),
				objectvalidator.Any(
					objectvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("combining_algorithm"),
					),
					objectvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("effect_settings"),
					),
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
				objectvalidatorinternal.ConflictsIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_POLICY)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				objectvalidator.AlsoRequires(
					path.MatchRelative().AtParent().AtName("name"),
				),
				objectvalidator.Any(
					objectvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("combining_algorithm"),
					),
					objectvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("effect_settings"),
					),
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
				objectvalidatorinternal.ConflictsIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				objectvalidator.AlsoRequires(
					path.MatchRelative().AtParent().AtName("combining_algorithm"),
					path.MatchRelative().AtParent().AtName("name"),
				),
			},
		},

		"effect_settings": schema.SingleNestedAttribute{
			Description:         effectSettingsDescription.Description,
			MarkdownDescription: effectSettingsDescription.MarkdownDescription,
			Optional:            true,

			Attributes: dataRulesEffectSettingsObjectSchemaAttributes(),

			Validators: []validator.Object{
				objectvalidator.ConflictsWith(
					path.MatchRelative().AtParent().AtName("value"),
				),
				objectvalidatorinternal.ConflictsIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_POLICY)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				objectvalidator.AlsoRequires(
					path.MatchRelative().AtParent().AtName("effect_settings"),
					path.MatchRelative().AtParent().AtName("name"),
				),
			},
		},
	}

	if iteration < policyChildNestedIterationMaxDepth {
		attributes["children"] = schema.ListNestedAttribute{
			Description:         childrenDescription.Description,
			MarkdownDescription: childrenDescription.MarkdownDescription,
			Optional:            true,

			NestedObject: schema.NestedAttributeObject{
				Attributes: dataPolicyChildObjectSchemaAttributesIteration(iteration + 1),
			},

			Validators: []validator.List{
				listvalidator.ConflictsWith(
					path.MatchRelative().AtParent().AtName("value"),
				),
				listvalidatorinternal.ConflictsIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE)),
					path.MatchRelative().AtParent().AtName("type"),
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

type editorDataPolicyChildLeafResourceModel struct {
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	// Statements         types.List                   `tfsdk:"statements"`
	Condition          types.Object `tfsdk:"condition"`
	CombiningAlgorithm types.Object `tfsdk:"combining_algorithm"`
	RepetitionSettings types.Object `tfsdk:"repetition_settings"`
	EffectSettings     types.Object `tfsdk:"effect_settings"`
	Value              types.Object `tfsdk:"value"`
}

type editorDataPolicyChildResourceModel struct {
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	// Statements         types.List                   `tfsdk:"statements"`
	Condition          types.Object `tfsdk:"condition"`
	CombiningAlgorithm types.Object `tfsdk:"combining_algorithm"`
	Children           types.List   `tfsdk:"children"`
	RepetitionSettings types.Object `tfsdk:"repetition_settings"`
	EffectSettings     types.Object `tfsdk:"effect_settings"`
	Value              types.Object `tfsdk:"value"`
}

var editorDataPolicyChildTFObjectTypes = initializeEditorDataPolicyChildTFObjectTypes(1)

func initializeEditorDataPolicyChildTFObjectTypes(iteration int32) map[string]attr.Type {

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
		"effect_settings": types.ObjectType{
			AttrTypes: editorDataRulesEffectSettingsTFObjectTypes,
		},
		"value": types.ObjectType{
			AttrTypes: editorReferenceObjectTFObjectTypes,
		},
	}

	if iteration < policyChildNestedIterationMaxDepth {
		attrMap["children"] = types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: initializeEditorDataPolicyChildTFObjectTypes(iteration + 1),
			},
		}
	}

	return attrMap
}

func expandEditorDataPolicyChildren(ctx context.Context, policyChildren basetypes.ListValue) (policyChildObjects []authorize.AuthorizeEditorDataPoliciesPolicyChild, diags diag.Diagnostics) {
	const initialIteration = 1
	return expandEditorDataPolicyChildrenIteration(ctx, policyChildren, initialIteration)
}

func expandEditorDataPolicyChildrenIteration(ctx context.Context, policyChildren basetypes.ListValue, iteration int32) (policyChildObjects []authorize.AuthorizeEditorDataPoliciesPolicyChild, diags diag.Diagnostics) {

	leaf := iteration >= policyChildNestedIterationMaxDepth

	returnPolicies := make([]authorize.AuthorizeEditorDataPoliciesPolicyChild, 0)

	if leaf {
		var plan []editorDataPolicyChildLeafResourceModel
		diags.Append(policyChildren.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, policyChildPlan := range plan {
			policyChildObject, d := policyChildPlan.expand(ctx)
			diags.Append(d...)
			if diags.HasError() {
				continue
			}
			returnPolicies = append(returnPolicies, *policyChildObject)
		}
	} else {
		var plan []editorDataPolicyChildResourceModel
		diags.Append(policyChildren.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, policyChildPlan := range plan {
			policyChildObject, d := policyChildPlan.expand(ctx, iteration)
			diags.Append(d...)
			if diags.HasError() {
				continue
			}
			returnPolicies = append(returnPolicies, *policyChildObject)
		}
	}

	return returnPolicies, diags
}

func (p *editorDataPolicyChildLeafResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataPoliciesPolicyChild, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorize.AuthorizeEditorDataPoliciesPolicyChild{}

	switch authorize.EnumAuthorizeEditorDataPoliciesPolicyChildCommonType(p.Type.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_POLICY:
		data.AuthorizeEditorDataPoliciesPolicyChildPolicy, d = p.expandPolicy(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE:
		data.AuthorizeEditorDataPoliciesPolicyChildRule, d = p.expandRule(ctx)
		diags.Append(d...)
	default:
		diags.AddError(
			"Invalid policy child type",
			fmt.Sprintf("The policy child type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &data, diags
}

func (p *editorDataPolicyChildResourceModel) expand(ctx context.Context, iteration int32) (*authorize.AuthorizeEditorDataPoliciesPolicyChild, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorize.AuthorizeEditorDataPoliciesPolicyChild{}

	switch authorize.EnumAuthorizeEditorDataPoliciesPolicyChildCommonType(p.Type.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_POLICY:
		data.AuthorizeEditorDataPoliciesPolicyChildPolicy, d = p.expandPolicy(ctx, iteration)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPOLICIESPOLICYCHILDCOMMONTYPE_RULE:
		data.AuthorizeEditorDataPoliciesPolicyChildRule, d = p.expandRule(ctx)
		diags.Append(d...)
	default:
		diags.AddError(
			"Invalid policy child type",
			fmt.Sprintf("The policy child type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &data, diags
}

func (p *editorDataPolicyChildResourceModel) expandPolicy(ctx context.Context, iteration int32) (*authorize.AuthorizeEditorDataPoliciesPolicyChildPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	data, d := expandChildPolicyChildPolicy(ctx, p.Name, p.Type, p.Description, p.Enabled, p.CombiningAlgorithm, p.Condition, p.RepetitionSettings, p.Value)
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

func (p *editorDataPolicyChildLeafResourceModel) expandPolicy(ctx context.Context) (*authorize.AuthorizeEditorDataPoliciesPolicyChildPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	data, d := expandChildPolicyChildPolicy(ctx, p.Name, p.Type, p.Description, p.Enabled, p.CombiningAlgorithm, p.Condition, p.RepetitionSettings, p.Value)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return data, diags
}

func expandChildPolicyChildPolicy(ctx context.Context, name, policyChildType, description basetypes.StringValue, enabled basetypes.BoolValue, combiningAlgorithm, condition, repetitionSettings, refValue basetypes.ObjectValue) (*authorize.AuthorizeEditorDataPoliciesPolicyChildPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.NewAuthorizeEditorDataPoliciesPolicyChildPolicy(
		authorize.EnumAuthorizeEditorDataPoliciesPolicyChildCommonType(policyChildType.ValueString()),
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
	// 	var plan []policyChildManagementPolicyChildStatementResourceModel
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

func (p *editorDataPolicyChildResourceModel) expandRule(ctx context.Context) (*authorize.AuthorizeEditorDataPoliciesPolicyChildRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	data, d := expandChildPolicyChildRule(ctx, p.Name, p.Type, p.Description, p.Enabled, p.Condition, p.EffectSettings, p.Value)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return data, diags
}

func (p *editorDataPolicyChildLeafResourceModel) expandRule(ctx context.Context) (*authorize.AuthorizeEditorDataPoliciesPolicyChildRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	data, d := expandChildPolicyChildRule(ctx, p.Name, p.Type, p.Description, p.Enabled, p.Condition, p.EffectSettings, p.Value)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return data, diags
}

func expandChildPolicyChildRule(ctx context.Context, name, policyChildType, description basetypes.StringValue, enabled basetypes.BoolValue, condition, effectSettings, refValue basetypes.ObjectValue) (*authorize.AuthorizeEditorDataPoliciesPolicyChildRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.NewAuthorizeEditorDataPoliciesPolicyChildRule(
		authorize.EnumAuthorizeEditorDataPoliciesPolicyChildCommonType(policyChildType.ValueString()),
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
	// 	var plan []policyChildManagementPolicyChildStatementResourceModel
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

	if !effectSettings.IsNull() && !effectSettings.IsUnknown() {
		var plan *editorDataRulesEffectSettingsResourceModel
		diags.Append(effectSettings.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		effectSettingsExp, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetEffectSettings(*effectSettingsExp)
	}

	if !condition.IsNull() && !condition.IsUnknown() {
		conditionExp, d := expandEditorDataCondition(ctx, condition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetCondition(*conditionExp)
	}

	return data, diags
}

func editorDataPolicyChildrenOkToListTF(ctx context.Context, apiObject []authorize.AuthorizeEditorDataPoliciesPolicyChild, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	return editorDataPolicyChildrenOkToListTFIteration(ctx, 1, apiObject, ok)
}

func editorDataPolicyChildrenOkToListTFIteration(ctx context.Context, iteration int32, apiObject []authorize.AuthorizeEditorDataPoliciesPolicyChild, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: initializeEditorDataPolicyChildTFObjectTypes(iteration)}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := editorDataPolicyChildOkToTFIteration(ctx, iteration, &v, true)
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

func editorDataPolicyChildOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataPoliciesPolicyChild, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	const initialIteration = 1
	return editorDataPolicyChildOkToTFIteration(ctx, initialIteration, apiObject, ok)
}

func editorDataPolicyChildOkToTFIteration(ctx context.Context, iteration int32, apiObject *authorize.AuthorizeEditorDataPoliciesPolicyChild, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorize.AuthorizeEditorDataPoliciesPolicyChild{}) {
		return types.ObjectNull(initializeEditorDataPolicyChildTFObjectTypes(iteration)), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case *authorize.AuthorizeEditorDataPoliciesPolicyChildPolicy:
		childrenVal, ok := t.GetChildrenOk()
		children, d := editorDataPolicyChildrenOkToListTFIteration(ctx, iteration+1, childrenVal, ok)
		diags.Append(d...)

		conditionVal, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTF(ctx, conditionVal, ok)
		diags.Append(d...)

		combiningAlgorithm, d := policyManagementPolicyCombiningAlgorithmOkToTF(t.GetCombiningAlgorithmOk())
		diags.Append(d...)

		repetitionSettings, d := policyManagementPolicyRepetitionSettingsOkToTF(t.GetRepetitionSettingsOk())
		diags.Append(d...)

		value, d := editorDataReferenceObjectOkToTF(t.GetValueOk())
		diags.Append(d...)

		if diags.HasError() {
			return types.ObjectNull(initializeEditorDataPolicyChildTFObjectTypes(iteration)), diags
		}

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["description"] = framework.StringOkToTF(t.GetDescriptionOk())
		attributeMap["enabled"] = framework.BoolOkToTF(t.GetEnabledOk())
		// "statements": framework.ListOkToTF(apiObject.GetStatementsOk())
		attributeMap["children"] = children
		attributeMap["condition"] = condition
		attributeMap["combining_algorithm"] = combiningAlgorithm
		attributeMap["repetition_settings"] = repetitionSettings
		attributeMap["value"] = value

	case *authorize.AuthorizeEditorDataPoliciesPolicyChildRule:
		conditionVal, ok := t.GetConditionOk()
		condition, d := editorDataConditionOkToTF(ctx, conditionVal, ok)
		diags.Append(d...)

		effectSettingsResp, ok := t.GetEffectSettingsOk()
		effectSettings, d := editorDataRulesEffectSettingsOkToTF(ctx, effectSettingsResp, ok)
		diags.Append(d...)

		value, d := editorDataReferenceObjectOkToTF(t.GetValueOk())
		diags.Append(d...)

		if diags.HasError() {
			return types.ObjectNull(initializeEditorDataPolicyChildTFObjectTypes(iteration)), diags
		}

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["description"] = framework.StringOkToTF(t.GetDescriptionOk())
		attributeMap["enabled"] = framework.BoolOkToTF(t.GetEnabledOk())
		// "statements": framework.ListOkToTF(apiObject.GetStatementsOk())
		attributeMap["condition"] = condition
		attributeMap["effect_settings"] = effectSettings
		attributeMap["value"] = value

	default:
		tflog.Error(ctx, "Invalid policy child type", map[string]interface{}{
			"policy child type": t,
		})
		diags.AddError(
			"Invalid policy child type",
			"The policy child type is not supported.  Please raise an issue with the provider maintainers.",
		)
		return types.ObjectNull(initializeEditorDataPolicyChildTFObjectTypes(iteration)), diags
	}

	attributeMap = editorDataPolicyChildConvertEmptyValuesToTFNulls(attributeMap, iteration)

	objValue, d := types.ObjectValue(initializeEditorDataPolicyChildTFObjectTypes(iteration), attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorDataPolicyChildConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value, iteration int32) map[string]attr.Value {
	nullMap := map[string]attr.Value{
		"type":        types.StringNull(),
		"name":        types.StringNull(),
		"description": types.StringNull(),
		"enabled":     types.BoolNull(),
		// "statements": types.BoolType,
		"condition":           types.ObjectNull(editorDataConditionTFObjectTypes),
		"combining_algorithm": types.ObjectNull(policyManagementPolicyCombiningAlgorithmTFObjectTypes),
		"repetition_settings": types.ObjectNull(policyManagementPolicyRepetitionSettingsTFObjectTypes),
		"effect_settings":     types.ObjectNull(editorDataRulesEffectSettingsTFObjectTypes),
		"value":               types.ObjectNull(editorReferenceObjectTFObjectTypes),
	}

	if iteration < policyChildNestedIterationMaxDepth {
		nullMap["children"] = types.ListNull(types.ObjectType{AttrTypes: initializeEditorDataPolicyChildTFObjectTypes(iteration + 1)})
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}
