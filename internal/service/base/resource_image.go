package base

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ImageResource serviceClientType

type ImageResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ImageFileB64  types.String                 `tfsdk:"image_file_base64"`
	UploadedImage types.Object                 `tfsdk:"uploaded_image"`
}

type ImageUploadedImageResourceModel struct {
	Width  types.Int64  `tfsdk:"width"`
	Height types.Int64  `tfsdk:"height"`
	Type   types.String `tfsdk:"type"`
	Href   types.String `tfsdk:"href"`
}

var (
	imageUploadedImageTFObjectTypes = map[string]attr.Type{
		"width":  types.Int64Type,
		"height": types.Int64Type,
		"type":   types.StringType,
		"href":   types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &ImageResource{}
	_ resource.ResourceWithConfigure   = &ImageResource{}
	_ resource.ResourceWithImportState = &ImageResource{}
)

// New Object
func NewImageResource() resource.Resource {
	return &ImageResource{}
}

// Metadata
func (r *ImageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

// Schema
func (r *ImageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	uploadedImageTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of format used for the image. Options are `jpg`, `png`, and `gif`.",
	)

	const attrMinLength = 2

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: framework.SchemaDescriptionFromMarkdown("Resource to create and manage PingOne images in an environment.").Description,

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the image in."),
			),

			"image_file_base64": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A base64 encoded image file to import.  Only PNG, GIF and JPG images are supported.  This field is immutable and will trigger a replace plan if changed.").Description,
				Required:    true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidatorinternal.IsBase64Encoded(),
					stringvalidatorinternal.IsB64ContentType("image/jpeg", "image/gif", "image/png"),
				},
			},

			"uploaded_image": schema.SingleNestedAttribute{
				Description: "A single object that specifies the processed image details.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"width": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("The width of the image (in pixels).").Description,
						Computed:    true,
					},

					"height": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("The height of the image (in pixels).").Description,
						Computed:    true,
					},

					"type": schema.StringAttribute{
						Description:         uploadedImageTypeDescription.Description,
						MarkdownDescription: uploadedImageTypeDescription.MarkdownDescription,
						Computed:            true,
					},

					"href": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the URL or fully qualified path to the image source file.").Description,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r *ImageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(framework.ResourceType)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected the provider client, got: %T. Please report this issue to the provider maintainers.", req.ProviderData),
		)

		return
	}

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *ImageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ImageResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	archive, fileName, contentType, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.Image
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ImagesApi.CreateImage(ctx, plan.EnvironmentId.ValueString()).ContentType(*contentType).ContentDisposition(fmt.Sprintf("attachment; filename=%s", *fileName)).File(archive).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateImage",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ImageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ImageResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.Image
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ImagesApi.ReadImage(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadImage",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ImageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ImageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ImageResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.ImagesApi.DeleteImage(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteImage",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ImageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "image_id",
			Regexp:    verify.P1ResourceIDRegexp,
			PrimaryID: true,
		},
	}

	attributes, err := framework.ParseImportID(req.ID, idComponents...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			err.Error(),
		)
		return
	}

	for _, idComponent := range idComponents {
		pathKey := idComponent.Label

		if idComponent.PrimaryID {
			pathKey = "id"
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(pathKey), attributes[idComponent.Label])...)
	}
}

func (p *ImageResourceModel) expand() (*[]byte, *string, *string, diag.Diagnostics) {
	var diags diag.Diagnostics

	var archive []byte

	var err error
	archive, err = base64.StdEncoding.DecodeString(p.ImageFileB64.ValueString())
	if err != nil {
		diags.AddError(
			"Cannot base64 decode provided image file.",
			fmt.Sprintf("The file cannot be base64 decoded: %s", err),
		)

		return nil, nil, nil, diags
	}

	chars := 12
	generatedName, err := utils.RandStringFromCharSet(chars, "abcdefghijklmnopqrstuvwxyz012346789")
	if err != nil {
		diags.AddError(
			"Cannot generate a filename to use",
			fmt.Sprintf("A random filename cannot be generated: %s", err),
		)

		return nil, nil, nil, diags
	}

	extensionMapping := map[string]string{
		"image/jpeg": "jpg",
		"image/gif":  "gif",
		"image/png":  "png",
	}

	contentType := http.DetectContentType(archive)
	if _, ok := extensionMapping[contentType]; !ok {
		diags.AddError(
			"Cannot determine the content type of the image.  Ensure the file is a jpg, gif or png format.",
			fmt.Sprintf("The file type has been determined to be `%s`, which is not supported.", contentType),
		)

		return nil, nil, nil, diags
	}

	fileName := fmt.Sprintf("%s.%s", generatedName, extensionMapping[contentType])

	return &archive, &fileName, &contentType, diags
}

func (p *ImageResourceModel) toState(apiObject *management.Image) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)

	var d diag.Diagnostics
	p.UploadedImage, d = toStateImageTarget(apiObject.GetTargetsOk())
	diags.Append(d...)

	return diags
}

func toStateImageTarget(v *management.ImageTargets, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	original, originalOk := v.GetOriginalOk()

	if !ok || v == nil || !originalOk || original == nil {
		return types.ObjectNull(imageUploadedImageTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"width":  framework.Int32OkToTF(original.GetWidthOk()),
		"height": framework.Int32OkToTF(original.GetHeightOk()),
		"type":   framework.EnumOkToTF(original.GetTypeOk()),
		"href":   framework.StringOkToTF(original.GetHrefOk()),
	}

	returnVar, d := types.ObjectValue(imageUploadedImageTFObjectTypes, objMap)
	diags.Append(d...)

	return returnVar, diags

}
