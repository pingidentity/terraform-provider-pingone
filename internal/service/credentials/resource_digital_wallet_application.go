package credentials

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type DigitalWalletApplicationResource struct {
	client *credentials.APIClient
	region model.RegionMapping
}

type DigitalWalletApplicationResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	EnvironmentId              types.String `tfsdk:"environment_id"`
	ApplicationId              types.String `tfsdk:"application_id"`
	DigitalWalletApplicationId types.String `tfsdk:"digital_wallet_application_id"`
	AppOpenUrl                 types.String `tfsdk:"app_open_url"`
	Name                       types.String `tfsdk:"name"`
}

// Framework interfaces
var (
	_ resource.Resource                = &DigitalWalletApplicationResource{}
	_ resource.ResourceWithConfigure   = &DigitalWalletApplicationResource{}
	_ resource.ResourceWithImportState = &DigitalWalletApplicationResource{}
)

// New Object
func NewDigitalWalletApplicationResource() resource.Resource {
	return &DigitalWalletApplicationResource{}
}

// Metadata
func (r *DigitalWalletApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_digital_wallet_application"
}

func (r *DigitalWalletApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Credentials digital wallet applications.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment to create the digital wallet application in."},
			),

			"application_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the application to associate with the digital wallet application.",
			}),

			"digital_wallet_application_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the digital wallet application.",
			}),

			"app_open_url": schema.StringAttribute{
				MarkdownDescription: "The URL sent in credential service notifications to the user to communicate with the service.",
				Required:            true,
			},

			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the digital wallet application.",
				Required:            true,
			},
		},
	}
}

func (r *DigitalWalletApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	preparedClient, err := prepareClient(ctx, resourceConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			err.Error(),
		)

		return
	}

	r.client = preparedClient
	r.region = resourceConfig.Client.API.Region
}

func (r *DigitalWalletApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state DigitalWalletApplicationResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	digitalWalletApplication := plan.expand()

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.DigitalWalletAppsApi.CreateDigitalWalletApp(ctx, plan.EnvironmentId.ValueString()).DigitalWalletApplication(*digitalWalletApplication).Execute()
		},
		"CreateDigitalWalletApplication",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*credentials.DigitalWalletApplication))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DigitalWalletApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DigitalWalletApplicationResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.DigitalWalletAppsApi.ReadOneDigitalWalletApp(ctx, data.EnvironmentId.ValueString(), data.DigitalWalletApplicationId.ValueString()).Execute()
		},
		"ReadOneDigitalWalletApplication",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response.(*credentials.DigitalWalletApplication))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// TODO: update changes
func (r *DigitalWalletApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *DigitalWalletApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DigitalWalletApplicationResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	_, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := r.client.DigitalWalletAppsApi.DeleteDigitalWalletApp(ctx, data.EnvironmentId.ValueString(), data.DigitalWalletApplicationId.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteDigitalWalletApplication",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DigitalWalletApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/digital_wallet_application_id/\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("digital_wallet_application_id"), attributes[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *DigitalWalletApplicationResourceModel) expand() *credentials.DigitalWalletApplication {

	applicationObject := credentials.NewObjectApplication()
	applicationObject.SetId(p.ApplicationId.ValueString())

	data := credentials.NewDigitalWalletApplication(*applicationObject, p.AppOpenUrl.ValueString(), p.Name.ValueString())

	return data
}

func (p *DigitalWalletApplicationResourceModel) toState(apiObject *credentials.DigitalWalletApplication) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.ApplicationId = framework.StringToTF(*apiObject.GetApplication().Id)
	p.Name = framework.StringToTF(apiObject.GetName())
	p.AppOpenUrl = framework.StringToTF(apiObject.GetAppOpenUrl())
	// p.DigitalWalletApplicationId = framework.StringOkToTF(*apiObject.GetDigitalWalletApplicationId())

	return diags
}
