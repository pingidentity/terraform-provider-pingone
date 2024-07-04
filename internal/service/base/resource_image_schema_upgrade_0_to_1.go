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

type imageResourceModelV0 struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ImageFileB64  types.String                 `tfsdk:"image_file_base64"`
	UploadedImage types.List                   `tfsdk:"uploaded_image"`
}

type imageUploadedImageResourceModelV0 imageUploadedImageResourceModelV1

func (r *ImageResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {

	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": framework.Attr_ID(),

					"environment_id": framework.Attr_LinkID(
						framework.SchemaAttributeDescriptionFromMarkdown(""),
					),

					"image_file_base64": schema.StringAttribute{
						Required: true,
					},
				},

				Blocks: map[string]schema.Block{
					"uploaded_image": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"width": schema.Int64Attribute{
									Computed: true,
								},
								"height": schema.Int64Attribute{
									Computed: true,
								},
								"type": schema.StringAttribute{
									Computed: true,
								},
								"href": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var d diag.Diagnostics
				var priorStateData imageResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				uploadedImage, d := priorStateData.schemaUpgradeUploadedImageV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				upgradedStateData := imageResourceModelV1{
					Id:            priorStateData.Id,
					EnvironmentId: priorStateData.EnvironmentId,
					ImageFileB64:  priorStateData.ImageFileB64,
					UploadedImage: uploadedImage,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}

func (p *imageResourceModelV0) schemaUpgradeUploadedImageV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := imageUploadedImageTFObjectTypes
	planAttribute := p.UploadedImage

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []imageUploadedImageResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := imageUploadedImageResourceModelV1{
			Width:  priorStateData[0].Width,
			Height: priorStateData[0].Height,
			Type:   priorStateData[0].Type,
			Href:   priorStateData[0].Href,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}
