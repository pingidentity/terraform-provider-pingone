package credentials

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

// Types
type CredentialIssuerProfileDataSource struct {
	client *credentials.APIClient
	region model.RegionMapping
}

type CredentialIssuerProfileDataSourceModel struct {
	Id                    types.String `tfsdk:"id"`
	EnvironmentId         types.String `tfsdk:"environment_id"`
	ApplicationInstanceId types.String `tfsdk:"application_instance_id"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
	Name                  types.String `tfsdk:"name"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &CredentialIssuerProfileDataSource{}
)

// New Object
func NewCredentialIssuerProfileDataSource() datasource.DataSource {
	return &CredentialIssuerProfileDataSource{}
}

// Metadata
func (r *CredentialIssuerProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_issuer_profile"
}

// Schema
func (r *CredentialIssuerProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// schema definition
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a PingOne Credentials issuer profile.  A PingOne environment contains only one Credential Issuer Profile.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that contains the credential issuer."),
			),

			"application_instance_id": schema.StringAttribute{
				Description: "Identifier (UUID) of the application instance registered with the PingOne platform service. This enables the client to send messages to the service.",
				Computed:    true,
			},

			"created_at": schema.StringAttribute{
				Description: "Date and time the issuer profile was created.",
				Computed:    true,
			},

			"updated_at": schema.StringAttribute{
				Description: "Date and time the issuer profile was last updated.",
				Computed:    true,
			},

			"name": schema.StringAttribute{
				Description: "The name of the credential issuer. The name is included in the metadata of an issued verifiable credential.",
				Computed:    true,
			},
		},
	}
}

func (r *CredentialIssuerProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *CredentialIssuerProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CredentialIssuerProfileDataSourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
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
		framework.DefaultCustomError,
		credentialIssuerDataSourceRetryConditions,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response.(*credentials.CredentialIssuerProfile))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *CredentialIssuerProfileDataSourceModel) toState(apiObject *credentials.CredentialIssuerProfile) diag.Diagnostics {
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
	p.ApplicationInstanceId = framework.StringToTF(apiObject.ApplicationInstance.GetId())
	p.CreatedAt = framework.TimeOkToTF(apiObject.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	return diags
}

func credentialIssuerDataSourceRetryConditions(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

	var err error

	if p1error != nil {

		// Credential Issuer Profile's keys may not have propagated after initial environment setup.
		// Rare, but possible.
		if m, _ := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
			tflog.Warn(ctx, "Insufficient PingOne privileges detected")
			return true
		}
		if err != nil {
			tflog.Warn(ctx, "Cannot match error string for retry")
			return false
		}

	}

	return false
}
