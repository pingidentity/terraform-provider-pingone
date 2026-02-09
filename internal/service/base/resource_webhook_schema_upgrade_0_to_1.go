// Copyright © 2026 Ping Identity Corporation

package base

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

type webhookResourceModelV0 struct {
	Id                     pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId          pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                   types.String                 `tfsdk:"name"`
	Enabled                types.Bool                   `tfsdk:"enabled"`
	HttpEndpointUrl        types.String                 `tfsdk:"http_endpoint_url"`
	HttpEndpointHeaders    types.Map                    `tfsdk:"http_endpoint_headers"`
	VerifyTLSCertificates  types.Bool                   `tfsdk:"verify_tls_certificates"`
	TLSClientAuthKeyPairId pingonetypes.ResourceIDValue `tfsdk:"tls_client_auth_key_pair_id"`
	Format                 types.String                 `tfsdk:"format"`
	FilterOptions          types.List                   `tfsdk:"filter_options"`
}

type webhookFilterOptionsResourceModelV0 webhookFilterOptionsResourceModelV1

func (r *WebhookResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
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

					"enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(false),
					},

					"http_endpoint_url": schema.StringAttribute{
						Required: true,
					},

					"http_endpoint_headers": schema.MapAttribute{
						Optional: true,

						ElementType: types.StringType,
					},

					"verify_tls_certificates": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(true),
					},

					"tls_client_auth_key_pair_id": schema.StringAttribute{
						Optional: true,

						CustomType: pingonetypes.ResourceIDType{},
					},

					"format": schema.StringAttribute{
						Required: true,
					},
				},

				Blocks: map[string]schema.Block{

					"filter_options": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"included_action_types": schema.SetAttribute{
									Required: true,

									ElementType: types.StringType,
								},

								"included_application_ids": schema.SetAttribute{
									Optional: true,

									ElementType: pingonetypes.ResourceIDType{},
								},

								"included_population_ids": schema.SetAttribute{
									Optional: true,

									ElementType: pingonetypes.ResourceIDType{},
								},

								"included_tags": schema.SetAttribute{
									Optional: true,

									ElementType: types.StringType,
								},

								"ip_address_exposed": schema.BoolAttribute{
									Optional: true,
									Computed: true,

									Default: booldefault.StaticBool(false),
								},

								"useragent_exposed": schema.BoolAttribute{
									Optional: true,
									Computed: true,

									Default: booldefault.StaticBool(false),
								},
							},
						},
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var d diag.Diagnostics
				var priorStateData webhookResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				filterOptions, d := priorStateData.schemaUpgradeFilterOptionsV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				upgradedStateData := webhookResourceModelV1{
					Id:                     priorStateData.Id,
					EnvironmentId:          priorStateData.EnvironmentId,
					Name:                   priorStateData.Name,
					Enabled:                priorStateData.Enabled,
					HttpEndpointUrl:        priorStateData.HttpEndpointUrl,
					HttpEndpointHeaders:    priorStateData.HttpEndpointHeaders,
					VerifyTLSCertificates:  priorStateData.VerifyTLSCertificates,
					TLSClientAuthKeyPairId: priorStateData.TLSClientAuthKeyPairId,
					Format:                 priorStateData.Format,
					FilterOptions:          filterOptions,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}

func (p *webhookResourceModelV0) schemaUpgradeFilterOptionsV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := webhookFilterOptionsTFObjectTypes
	planAttribute := p.FilterOptions

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []webhookFilterOptionsResourceModelV0
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
