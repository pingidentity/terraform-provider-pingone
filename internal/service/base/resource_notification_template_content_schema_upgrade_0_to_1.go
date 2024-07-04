package base

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

type notificationTemplateContentResourceModelV0 struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	TemplateName  types.String                 `tfsdk:"template_name"`
	Locale        types.String                 `tfsdk:"locale"`
	Default       types.Bool                   `tfsdk:"default"`
	Variant       types.String                 `tfsdk:"variant"`
	Email         types.List                   `tfsdk:"email"`
	Push          types.List                   `tfsdk:"push"`
	Sms           types.List                   `tfsdk:"sms"`
	Voice         types.List                   `tfsdk:"voice"`
}

type notificationTemplateContentEmailResourceModelV0 struct {
	Body         types.String `tfsdk:"body"`
	From         types.List   `tfsdk:"from"`
	Subject      types.String `tfsdk:"subject"`
	ReplyTo      types.List   `tfsdk:"reply_to"`
	CharacterSet types.String `tfsdk:"character_set"`
	ContentType  types.String `tfsdk:"content_type"`
}

type notificationTemplateContentEmailAddressResourceModelV0 notificationTemplateContentEmailAddressResourceModelV1

type notificationTemplateContentPushResourceModelV0 notificationTemplateContentPushResourceModelV1

type notificationTemplateContentSmsResourceModelV0 notificationTemplateContentSmsResourceModelV1

type notificationTemplateContentVoiceResourceModelV0 notificationTemplateContentVoiceResourceModelV1

func (r *NotificationTemplateContentResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": framework.Attr_ID(),

					"environment_id": framework.Attr_LinkID(
						framework.SchemaAttributeDescriptionFromMarkdown(""),
					),

					"template_name": schema.StringAttribute{
						Required: true,
					},

					"locale": schema.StringAttribute{
						Required: true,
					},

					"default": schema.BoolAttribute{
						Computed: true,
					},

					"variant": schema.StringAttribute{
						Optional: true,
					},
				},

				Blocks: map[string]schema.Block{
					"email": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"body": schema.StringAttribute{
									Required: true,
								},

								"subject": schema.StringAttribute{
									Required: true,
								},

								"character_set": schema.StringAttribute{
									Optional: true,
									Computed: true,

									Default: stringdefault.StaticString("UTF-8"),
								},

								"content_type": schema.StringAttribute{
									Optional: true,
									Computed: true,

									Default: stringdefault.StaticString("text/html"),
								},
							},

							Blocks: map[string]schema.Block{
								"from": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Optional: true,
												Computed: true,

												Default: stringdefault.StaticString("PingOne"),
											},

											"address": schema.StringAttribute{
												Optional: true,
												Computed: true,

												Default: stringdefault.StaticString("noreply@pingidentity.com"),
											},
										},
									},
								},

								"reply_to": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Optional: true,
												Computed: true,
											},

											"address": schema.StringAttribute{
												Optional: true,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},

					"push": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"category": schema.StringAttribute{
									Optional: true,
									Computed: true,

									Default: stringdefault.StaticString(string(management.ENUMTEMPLATECONTENTPUSHCATEGORY_BANNER_BUTTONS)),
								},

								"body": schema.StringAttribute{
									Required: true,
								},

								"title": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},

					"sms": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"content": schema.StringAttribute{
									Required: true,
								},

								"sender": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},

					"voice": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"content": schema.StringAttribute{
									Required: true,
								},

								"type": schema.StringAttribute{
									Optional: true,
									Computed: true,

									Default: stringdefault.StaticString("Alice"),
								},
							},
						},
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var d diag.Diagnostics
				var priorStateData notificationTemplateContentResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				email, d := priorStateData.schemaUpgradeEmailV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				push, d := priorStateData.schemaUpgradePushV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				sms, d := priorStateData.schemaUpgradeSmsV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				voice, d := priorStateData.schemaUpgradeVoiceV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				upgradedStateData := notificationTemplateContentResourceModelV1{
					Id:            priorStateData.Id,
					EnvironmentId: priorStateData.EnvironmentId,
					TemplateName:  priorStateData.TemplateName,
					Locale:        priorStateData.Locale,
					Default:       priorStateData.Default,
					Variant:       priorStateData.Variant,
					Email:         email,
					Push:          push,
					Sms:           sms,
					Voice:         voice,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}

func (p *notificationTemplateContentResourceModelV0) schemaUpgradeEmailV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := notificationTemplateContentEmailTFObjectTypes
	planAttribute := p.Email

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []notificationTemplateContentEmailResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		from, d := priorStateData[0].schemaUpgradeFromV0toV1(ctx)
		diags.Append(d...)

		replyTo, d := priorStateData[0].schemaUpgradeReplyToV0toV1(ctx)
		diags.Append(d...)

		upgradedStateData := notificationTemplateContentEmailResourceModelV1{
			Body:         priorStateData[0].Body,
			From:         from,
			Subject:      priorStateData[0].Subject,
			ReplyTo:      replyTo,
			CharacterSet: priorStateData[0].CharacterSet,
			ContentType:  priorStateData[0].ContentType,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *notificationTemplateContentEmailResourceModelV0) schemaUpgradeEmailAddressV0toV1(ctx context.Context, planAttribute types.List) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := notificationTemplateContentEmailAddressTFObjectTypes

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []notificationTemplateContentEmailAddressResourceModelV0
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

func (p *notificationTemplateContentEmailResourceModelV0) schemaUpgradeFromV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return p.schemaUpgradeEmailAddressV0toV1(ctx, p.From)
}

func (p *notificationTemplateContentEmailResourceModelV0) schemaUpgradeReplyToV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return p.schemaUpgradeEmailAddressV0toV1(ctx, p.ReplyTo)
}

func (p *notificationTemplateContentResourceModelV0) schemaUpgradePushV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := notificationTemplateContentPushTFObjectTypes
	planAttribute := p.Push

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []notificationTemplateContentPushResourceModelV0
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

func (p *notificationTemplateContentResourceModelV0) schemaUpgradeSmsV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := notificationTemplateContentSmsTFObjectTypes
	planAttribute := p.Sms

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []notificationTemplateContentSmsResourceModelV0
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

func (p *notificationTemplateContentResourceModelV0) schemaUpgradeVoiceV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := notificationTemplateContentVoiceTFObjectTypes
	planAttribute := p.Voice

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []notificationTemplateContentVoiceResourceModelV0
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
