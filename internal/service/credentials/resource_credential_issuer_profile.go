package credentials

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type CredentialIssuerProfileResource struct {
	client *credentials.APIClient
	region model.RegionMapping
}

type CredentialIssuerProfileResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Name          types.String `tfsdk:"name"`
}

// Framework interfaces
var (
	_ resource.Resource                = &CredentialIssuerProfileResource{}
	_ resource.ResourceWithConfigure   = &CredentialIssuerProfileResource{}
	_ resource.ResourceWithImportState = &CredentialIssuerProfileResource{}
)

// New Object
func NewCredentialIssuerProfileResource() resource.Resource {
	return &CredentialIssuerProfileResource{}
}

// Metadata
func (r *CredentialIssuerProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_issuer_profile"
}

func (r *CredentialIssuerProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	// schema definition
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to enable the issuance of credentials, and manage the PingOne Credentials issuer profile.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment to create the credential type in."},
			),

			"name": schema.StringAttribute{
				Description: "The name of the credential issuer. This will be included in credentials issued.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			// omitted logo because it is not used - API target is to deprecate
			// placeholder just in case it needs to be enabled
			//"logo": schema.StringAttribute{
			//	Description: "An image containing the brand logo for the issuer. ",
			//	Optional:    true,
			//},
		},
	}
}

func (r *CredentialIssuerProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CredentialIssuerProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state CredentialIssuerProfileResourceModel

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
	CredentialIssuerProfile, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.CredentialIssuersApi.CreateCredentialIssuerProfile(ctx, plan.EnvironmentId.ValueString()).CredentialIssuerProfile(*CredentialIssuerProfile).Execute()
		},
		"CreateCredentialIssuerProfile",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*credentials.CredentialIssuerProfile))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CredentialIssuerProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CredentialIssuerProfileResourceModel

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
			return r.client.CredentialIssuersApi.ReadCredentialIssuerProfile(ctx, data.EnvironmentId.ValueString()).Execute()

		},
		"ReadCredentialIssuerProfile",
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
	resp.Diagnostics.Append(data.toState(response.(*credentials.CredentialIssuerProfile))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CredentialIssuerProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state CredentialIssuerProfileResourceModel

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
	CredentialIssuerProfile, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.CredentialIssuersApi.UpdateCredentialIssuerProfile(ctx, plan.EnvironmentId.ValueString()).CredentialIssuerProfile(*CredentialIssuerProfile).Execute()
		},
		"UpdateCredentialIssuerProfile",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*credentials.CredentialIssuerProfile))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CredentialIssuerProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Deletion of a profile is not allowed.  Only removal of the service.
	// TODO: Proper handling.
}

func (r *CredentialIssuerProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/credential_issuer_profile_id/\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *CredentialIssuerProfileResourceModel) expand(ctx context.Context) (*credentials.CredentialIssuerProfile, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := credentials.NewCredentialIssuerProfile(p.Name.ValueString())
	return data, diags
}

func (p *CredentialIssuerProfileResourceModel) toState(apiObject *credentials.CredentialIssuerProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.Name = framework.StringToTF(apiObject.GetName())

	return diags
}
