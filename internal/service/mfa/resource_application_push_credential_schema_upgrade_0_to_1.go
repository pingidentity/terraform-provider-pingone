// Copyright © 2026 Ping Identity Corporation

package mfa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

type applicationPushCredentialResourceModelV0 struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ApplicationId pingonetypes.ResourceIDValue `tfsdk:"application_id"`
	Fcm           types.List                   `tfsdk:"fcm"`
	Apns          types.List                   `tfsdk:"apns"`
	Hms           types.List                   `tfsdk:"hms"`
}

type applicationPushCredentialFcmResourceModelV0 struct {
	Key                             types.String `tfsdk:"key"`
	GoogleServiceAccountCredentials types.String `tfsdk:"google_service_account_credentials"`
}

type applicationPushCredentialApnsResourceModelV0 applicationPushCredentialApnsResourceModelV1

type applicationPushCredentialHmsResourceModelV0 applicationPushCredentialHmsResourceModelV1

func (r *ApplicationPushCredentialResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": framework.Attr_ID(),

					"environment_id": framework.Attr_LinkID(
						framework.SchemaAttributeDescriptionFromMarkdown(""),
					),

					"application_id": framework.Attr_LinkID(
						framework.SchemaAttributeDescriptionFromMarkdown(""),
					),
				},

				Blocks: map[string]schema.Block{

					"fcm": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},

								"google_service_account_credentials": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
							},
						},
					},

					"apns": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},

								"team_id": schema.StringAttribute{
									Required: true,
								},

								"token_signing_key": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},
							},
						},
					},

					"hms": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"client_id": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},

								"client_secret": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var d diag.Diagnostics
				var priorStateData applicationPushCredentialResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				fcm, d := priorStateData.schemaUpgradeFCMV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				apns, d := priorStateData.schemaUpgradeAPNSV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				hms, d := priorStateData.schemaUpgradeHMSV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				upgradedStateData := applicationPushCredentialResourceModelV1{
					Id:            priorStateData.Id,
					EnvironmentId: priorStateData.EnvironmentId,
					ApplicationId: priorStateData.ApplicationId,
					Fcm:           fcm,
					Apns:          apns,
					Hms:           hms,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}

func (p *applicationPushCredentialResourceModelV0) schemaUpgradeFCMV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"google_service_account_credentials": jsontypes.NormalizedType{},
	}
	planAttribute := p.Fcm

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationPushCredentialFcmResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := applicationPushCredentialFcmResourceModelV1{
			GoogleServiceAccountCredentials: priorStateData[0].schemaUpgradeGoogleServiceAccountCredentialsV0toV1(),
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationPushCredentialFcmResourceModelV0) schemaUpgradeGoogleServiceAccountCredentialsV0toV1() jsontypes.Normalized {

	planAttribute := p.GoogleServiceAccountCredentials

	if planAttribute.IsNull() {
		return jsontypes.NewNormalizedNull()
	} else if planAttribute.IsUnknown() {
		return jsontypes.NewNormalizedUnknown()
	} else {
		return jsontypes.NewNormalizedValue(planAttribute.ValueString())
	}
}

func (p *applicationPushCredentialResourceModelV0) schemaUpgradeAPNSV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"key":               types.StringType,
		"team_id":           types.StringType,
		"token_signing_key": types.StringType,
	}
	planAttribute := p.Apns

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationPushCredentialApnsResourceModelV0
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

func (p *applicationPushCredentialResourceModelV0) schemaUpgradeHMSV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"client_id":     types.StringType,
		"client_secret": types.StringType,
	}
	planAttribute := p.Hms

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationPushCredentialHmsResourceModelV0
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
