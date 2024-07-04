package sso

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

type passwordPolicyResourceModelV0 struct {
	Id                           pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId                pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                         types.String                 `tfsdk:"name"`
	Description                  types.String                 `tfsdk:"description"`
	EnvironmentDefault           types.Bool                   `tfsdk:"environment_default"`
	BypassPolicy                 types.Bool                   `tfsdk:"bypass_policy"`
	ExcludeCommonlyUsedPasswords types.Bool                   `tfsdk:"exclude_commonly_used_passwords"`
	ExcludeProfileData           types.Bool                   `tfsdk:"exclude_profile_data"`
	PasswordHistory              types.List                   `tfsdk:"password_history"`
	PasswordLength               types.List                   `tfsdk:"password_length"`
	AccountLockout               types.List                   `tfsdk:"account_lockout"`
	MinCharacters                types.List                   `tfsdk:"min_characters"`
	PasswordAge                  types.List                   `tfsdk:"password_age"`
	MaxRepeatedCharacters        types.Int64                  `tfsdk:"max_repeated_characters"`
	MinComplexity                types.Int64                  `tfsdk:"min_complexity"`
	MinUniqueCharacters          types.Int64                  `tfsdk:"min_unique_characters"`
	NotSimilarToCurrent          types.Bool                   `tfsdk:"not_similar_to_current"`
	PopulationCount              types.Int64                  `tfsdk:"population_count"`
}

type passwordPolicyPasswordHistoryResourceModelV0 struct {
	PriorPasswordCount types.Int64 `tfsdk:"prior_password_count"`
	RetentionDays      types.Int64 `tfsdk:"retention_days"`
}

type passwordPolicyPasswordLengthResourceModelV0 passwordPolicyLengthResourceModelV1

type passwordPolicyAccountLockoutResourceModelV0 struct {
	DurationSeconds types.Int64 `tfsdk:"duration_seconds"`
	FailCount       types.Int64 `tfsdk:"fail_count"`
}

type passwordPolicyMinCharactersResourceModelV0 passwordPolicyMinCharactersResourceModelV1

type passwordPolicyPasswordAgeResourceModelV0 struct {
	Max types.Int64 `tfsdk:"max"`
	Min types.Int64 `tfsdk:"min"`
}

func (r *PasswordPolicyResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {

	const passwordLengthMaxDefault = 255
	const passwordLengthMinDefault = 8

	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": framework.Attr_ID(),

					"environment_id": framework.Attr_LinkID(
						framework.SchemaAttributeDescriptionFromMarkdown(""),
					),

					"name": schema.StringAttribute{
						Required: true,
					},

					"description": schema.StringAttribute{
						Optional: true,
					},

					"environment_default": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(false),
					},

					"bypass_policy": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(false),
					},

					"exclude_commonly_used_passwords": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(false),
					},

					"exclude_profile_data": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(false),
					},

					"max_repeated_characters": schema.Int64Attribute{
						Optional: true,
					},

					"min_complexity": schema.Int64Attribute{
						Optional: true,
					},

					"min_unique_characters": schema.Int64Attribute{
						Optional: true,
					},

					"not_similar_to_current": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(false),
					},

					"population_count": schema.Int64Attribute{
						Computed: true,
					},
				},

				Blocks: map[string]schema.Block{
					"password_history": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"prior_password_count": schema.Int64Attribute{
									Optional: true,
								},

								"retention_days": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},

					"password_length": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"max": schema.Int64Attribute{
									Optional: true,
									Computed: true,

									Default: int64default.StaticInt64(passwordLengthMaxDefault),
								},

								"min": schema.Int64Attribute{
									Optional: true,
									Computed: true,

									Default: int64default.StaticInt64(passwordLengthMinDefault),
								},
							},
						},
					},

					"account_lockout": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"duration_seconds": schema.Int64Attribute{
									Optional: true,
								},

								"fail_count": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},

					"min_characters": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"alphabetical_uppercase": schema.Int64Attribute{
									Optional: true,
								},

								"alphabetical_lowercase": schema.Int64Attribute{
									Optional: true,
								},

								"numeric": schema.Int64Attribute{
									Optional: true,
								},

								"special_characters": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},

					"password_age": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"max": schema.Int64Attribute{
									Optional: true,
								},

								"min": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var d diag.Diagnostics
				var priorStateData passwordPolicyResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				history, d := priorStateData.schemaUpgradeHistoryV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				length, d := priorStateData.schemaUpgradeLengthV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				lockout, d := priorStateData.schemaUpgradeLockoutV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				minCharacters, d := priorStateData.schemaUpgradeMinCharactersV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				passwordAgeMax, d := priorStateData.schemaUpgradePasswordAgeMaxV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				passwordAgeMin, d := priorStateData.schemaUpgradePasswordAgeMinV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				upgradedStateData := passwordPolicyResourceModelV1{
					Id:                            priorStateData.Id,
					EnvironmentId:                 priorStateData.EnvironmentId,
					Name:                          priorStateData.Name,
					Description:                   priorStateData.Description,
					Default:                       priorStateData.EnvironmentDefault,
					ExcludesCommonlyUsedPasswords: priorStateData.ExcludeCommonlyUsedPasswords,
					ExcludesProfileData:           priorStateData.ExcludeProfileData,
					History:                       history,
					Length:                        length,
					Lockout:                       lockout,
					MinCharacters:                 minCharacters,
					PasswordAgeMax:                passwordAgeMax,
					PasswordAgeMin:                passwordAgeMin,
					MaxRepeatedCharacters:         priorStateData.MaxRepeatedCharacters,
					MinComplexity:                 priorStateData.MinComplexity,
					MinUniqueCharacters:           priorStateData.MinUniqueCharacters,
					NotSimilarToCurrent:           priorStateData.NotSimilarToCurrent,
					PopulationCount:               priorStateData.PopulationCount,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}

func (p *passwordPolicyResourceModelV0) schemaUpgradeHistoryV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := passwordPolicyHistoryTFObjectTypes
	planAttribute := p.PasswordHistory

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []passwordPolicyPasswordHistoryResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := passwordPolicyHistoryResourceModelV1{
			Count:         priorStateData[0].PriorPasswordCount,
			RetentionDays: priorStateData[0].RetentionDays,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *passwordPolicyResourceModelV0) schemaUpgradeLengthV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := passwordPolicyLengthTFObjectTypes
	planAttribute := p.PasswordLength

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []passwordPolicyPasswordLengthResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := passwordPolicyLengthResourceModelV1{
			Max: priorStateData[0].Max,
			Min: priorStateData[0].Min,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *passwordPolicyResourceModelV0) schemaUpgradeLockoutV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := passwordPolicyLockoutTFObjectTypes
	planAttribute := p.AccountLockout

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []passwordPolicyAccountLockoutResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := passwordPolicyLockoutResourceModelV1{
			DurationSeconds: priorStateData[0].DurationSeconds,
			FailureCount:    priorStateData[0].FailCount,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *passwordPolicyResourceModelV0) schemaUpgradeMinCharactersV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := passwordPolicyMinCharactersTFObjectTypes
	planAttribute := p.MinCharacters

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []passwordPolicyMinCharactersResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := passwordPolicyMinCharactersResourceModelV1{
			AlphabeticalUppercase: priorStateData[0].AlphabeticalUppercase,
			AlphabeticalLowercase: priorStateData[0].AlphabeticalLowercase,
			Numeric:               priorStateData[0].Numeric,
			SpecialCharacters:     priorStateData[0].SpecialCharacters,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *passwordPolicyResourceModelV0) schemaUpgradePasswordAgeMaxV0toV1(ctx context.Context) (types.Int64, diag.Diagnostics) {
	var diags diag.Diagnostics

	planAttribute := p.PasswordAge

	if planAttribute.IsNull() {
		return types.Int64Null(), diags
	} else if planAttribute.IsUnknown() {
		return types.Int64Unknown(), diags
	} else {
		var priorStateData []passwordPolicyPasswordAgeResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.Int64Null(), diags
		}

		if len(priorStateData) == 0 {
			return types.Int64Null(), diags
		}

		returnVar := priorStateData[0].Max

		return returnVar, diags
	}
}

func (p *passwordPolicyResourceModelV0) schemaUpgradePasswordAgeMinV0toV1(ctx context.Context) (types.Int64, diag.Diagnostics) {
	var diags diag.Diagnostics

	planAttribute := p.PasswordAge

	if planAttribute.IsNull() {
		return types.Int64Null(), diags
	} else if planAttribute.IsUnknown() {
		return types.Int64Unknown(), diags
	} else {
		var priorStateData []passwordPolicyPasswordAgeResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.Int64Null(), diags
		}

		if len(priorStateData) == 0 {
			return types.Int64Null(), diags
		}

		returnVar := priorStateData[0].Min

		return returnVar, diags
	}
}
