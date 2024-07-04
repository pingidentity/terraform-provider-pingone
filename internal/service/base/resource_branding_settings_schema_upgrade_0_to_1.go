package base

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
)

type brandingSettingsResourceModelV0 struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	CompanyName   types.String                 `tfsdk:"company_name"`
	LogoImage     types.List                   `tfsdk:"logo_image"`
}

func (r *BrandingSettingsResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": framework.Attr_ID(),

					"environment_id": framework.Attr_LinkID(
						framework.SchemaAttributeDescriptionFromMarkdown(""),
					),

					"company_name": schema.StringAttribute{
						Optional: true,
						Computed: true,

						Default: stringdefault.StaticString(""),
					},
				},

				Blocks: map[string]schema.Block{
					"logo_image": schema.ListNestedBlock{
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
				var priorStateData brandingSettingsResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				icon, d := priorStateData.schemaUpgradeIconV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				upgradedStateData := brandingSettingsResourceModelV1{
					Id:            priorStateData.Id,
					EnvironmentId: priorStateData.EnvironmentId,
					CompanyName:   priorStateData.CompanyName,
					LogoImage:     icon,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}

func (p *brandingSettingsResourceModelV0) schemaUpgradeIconV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return service.ImageListToObjectSchemaUpgrade(ctx, p.LogoImage)
}
