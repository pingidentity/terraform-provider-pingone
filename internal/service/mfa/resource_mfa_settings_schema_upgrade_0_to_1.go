// Copyright © 2025 Ping Identity Corporation

package mfa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

type mFASettingsResourceModelV0 struct {
	Id                     pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId          pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	PhoneExtensionsEnabled types.Bool                   `tfsdk:"phone_extensions_enabled"`
	Pairing                types.List                   `tfsdk:"pairing"`
	Lockout                types.List                   `tfsdk:"lockout"`
	Authentication         types.List                   `tfsdk:"authentication"`
}

type mFASettingsLockoutResourceModelV0 mFASettingsLockoutResourceModelV1

type mFASettingsPairingResourceModelV0 mFASettingsPairingResourceModelV1

func (r *MFASettingsResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {

	const pairingMaxAllowedDevicesDescription = 5
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": framework.Attr_ID(),

					"environment_id": framework.Attr_LinkID(
						framework.SchemaAttributeDescriptionFromMarkdown(""),
					),

					"phone_extensions_enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(false),
					},
				},

				Blocks: map[string]schema.Block{

					"pairing": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"max_allowed_devices": schema.Int32Attribute{
									Optional: true,
									Computed: true,

									Default: int32default.StaticInt32(pairingMaxAllowedDevicesDescription),
								},

								"pairing_key_format": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},

					"lockout": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"failure_count": schema.Int32Attribute{
									Required: true,
								},

								"duration_seconds": schema.Int32Attribute{
									Optional: true,
								},
							},
						},
					},

					"authentication": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"device_selection": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var d diag.Diagnostics
				var priorStateData mFASettingsResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				lockout, d := priorStateData.schemaUpgradeLockoutV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				pairing, d := priorStateData.schemaUpgradePairingV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				phoneExtentions, d := priorStateData.schemaUpgradePhoneExtensionsV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				upgradedStateData := mFASettingsResourceModelV1{
					EnvironmentId:   priorStateData.EnvironmentId,
					Lockout:         lockout,
					Pairing:         pairing,
					PhoneExtensions: phoneExtentions,
					Users:           types.ObjectNull(MFASettingsUsersTFObjectTypes),
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}

func (p *mFASettingsResourceModelV0) schemaUpgradeLockoutV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := MFASettingsLockoutTFObjectTypes
	planAttribute := p.Lockout

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []mFASettingsLockoutResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := priorStateData[0]

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *mFASettingsResourceModelV0) schemaUpgradePairingV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := MFASettingsPairingTFObjectTypes
	planAttribute := p.Pairing

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []mFASettingsPairingResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := priorStateData[0]

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *mFASettingsResourceModelV0) schemaUpgradePhoneExtensionsV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := MFASettingsPhoneExtensionsTFObjectTypes
	planAttribute := p.PhoneExtensionsEnabled

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		upgradedStateData := mFASettingsPhoneExtensionsResourceModelV1{
			Enabled: planAttribute,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}
