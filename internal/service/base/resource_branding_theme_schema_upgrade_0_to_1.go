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
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
)

type brandingThemeResourceModelV0 struct {
	Id                   pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId        pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                 types.String                 `tfsdk:"name"`
	Template             types.String                 `tfsdk:"template"`
	Default              types.Bool                   `tfsdk:"default"`
	Logo                 types.List                   `tfsdk:"logo"`
	BackgroundImage      types.List                   `tfsdk:"background_image"`
	BackgroundColor      types.String                 `tfsdk:"background_color"`
	UseDefaultBackground types.Bool                   `tfsdk:"use_default_background"`
	BodyTextColor        types.String                 `tfsdk:"body_text_color"`
	ButtonColor          types.String                 `tfsdk:"button_color"`
	ButtonTextColor      types.String                 `tfsdk:"button_text_color"`
	CardColor            types.String                 `tfsdk:"card_color"`
	FooterText           types.String                 `tfsdk:"footer_text"`
	HeadingTextColor     types.String                 `tfsdk:"heading_text_color"`
	LinkTextColor        types.String                 `tfsdk:"link_text_color"`
}

func (r *BrandingThemeResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
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

					"template": schema.StringAttribute{
						Required: true,
					},

					"default": schema.BoolAttribute{
						Computed: true,
					},

					"background_color": schema.StringAttribute{
						Optional: true,
					},

					"use_default_background": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(false),
					},

					"body_text_color": schema.StringAttribute{
						Required: true,
					},

					"button_color": schema.StringAttribute{
						Required: true,
					},

					"button_text_color": schema.StringAttribute{
						Required: true,
					},

					"card_color": schema.StringAttribute{
						Required: true,
					},

					"footer_text": schema.StringAttribute{
						Optional: true,
					},

					"heading_text_color": schema.StringAttribute{
						Required: true,
					},

					"link_text_color": schema.StringAttribute{
						Required: true,
					},
				},

				Blocks: map[string]schema.Block{
					"logo": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Required: true,

									CustomType: pingonetypes.ResourceIDType{},
								},
								"href": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},

					"background_image": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Required: true,

									CustomType: pingonetypes.ResourceIDType{},
								},
								"href": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var d diag.Diagnostics
				var priorStateData brandingThemeResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				logo, d := priorStateData.schemaUpgradeLogoV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				backgroundImage, d := priorStateData.schemaUpgradeBackgroundImageV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				upgradedStateData := brandingThemeResourceModelV1{
					Id:                   priorStateData.Id,
					EnvironmentId:        priorStateData.EnvironmentId,
					Name:                 priorStateData.Name,
					Template:             priorStateData.Template,
					Default:              priorStateData.Default,
					Logo:                 logo,
					BackgroundImage:      backgroundImage,
					BackgroundColor:      priorStateData.BackgroundColor,
					UseDefaultBackground: priorStateData.UseDefaultBackground,
					BodyTextColor:        priorStateData.BodyTextColor,
					ButtonColor:          priorStateData.ButtonColor,
					ButtonTextColor:      priorStateData.ButtonTextColor,
					CardColor:            priorStateData.CardColor,
					FooterText:           priorStateData.FooterText,
					HeadingTextColor:     priorStateData.HeadingTextColor,
					LinkTextColor:        priorStateData.LinkTextColor,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}

func (p *brandingThemeResourceModelV0) schemaUpgradeLogoV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return service.ImageListToObjectSchemaUpgrade(ctx, p.Logo)
}

func (p *brandingThemeResourceModelV0) schemaUpgradeBackgroundImageV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return service.ImageListToObjectSchemaUpgrade(ctx, p.BackgroundImage)
}
