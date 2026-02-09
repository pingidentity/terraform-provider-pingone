// Copyright © 2026 Ping Identity Corporation

package base

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

type notificationSettingsEmailResourceModelV0 struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Host          types.String                 `tfsdk:"host"`
	Port          types.Int32                  `tfsdk:"port"`
	Protocol      types.String                 `tfsdk:"protocol"`
	Username      types.String                 `tfsdk:"username"`
	Password      types.String                 `tfsdk:"password"`
	From          types.List                   `tfsdk:"from"`
	ReplyTo       types.List                   `tfsdk:"reply_to"`
}

type emailSourceModelV0 emailSourceModelV1

func (r *NotificationSettingsEmailResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": framework.Attr_ID(),

					"environment_id": framework.Attr_LinkID(
						framework.SchemaAttributeDescriptionFromMarkdown(""),
					),

					"host": schema.StringAttribute{
						Required: true,
					},

					"port": schema.Int32Attribute{
						Required: true,
					},

					"protocol": schema.StringAttribute{
						Computed: true,
					},

					"username": schema.StringAttribute{
						Required: true,
					},

					"password": schema.StringAttribute{
						Required:  true,
						Sensitive: true,
					},
				},

				Blocks: map[string]schema.Block{
					"from": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Optional: true,
								},
								"email_address": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
					"reply_to": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Optional: true,
								},
								"email_address": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var d diag.Diagnostics
				var priorStateData notificationSettingsEmailResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				from, d := priorStateData.schemaUpgradeFromV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				replyTo, d := priorStateData.schemaUpgradeReplyToV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				upgradedStateData := notificationSettingsEmailResourceModelV1{
					Id:            priorStateData.Id,
					EnvironmentId: priorStateData.EnvironmentId,
					Host:          priorStateData.Host,
					Port:          priorStateData.Port,
					Protocol:      priorStateData.Protocol,
					Username:      priorStateData.Username,
					Password:      priorStateData.Password,
					From:          from,
					ReplyTo:       replyTo,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}

func (p *notificationSettingsEmailResourceModelV0) schemaUpgradeEmailSourceV0toV1(ctx context.Context, planAttribute types.List) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := emailSourceTFObjectTypes

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []emailSourceModelV0
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

func (p *notificationSettingsEmailResourceModelV0) schemaUpgradeFromV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return p.schemaUpgradeEmailSourceV0toV1(ctx, p.From)
}

func (p *notificationSettingsEmailResourceModelV0) schemaUpgradeReplyToV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return p.schemaUpgradeEmailSourceV0toV1(ctx, p.ReplyTo)
}
