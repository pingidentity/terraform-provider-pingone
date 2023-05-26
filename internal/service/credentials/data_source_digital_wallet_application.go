package credentials

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type DigitalWalletApplicationDataSource struct {
	client *credentials.APIClient
	region model.RegionMapping
}

type DigitalWalletApplicationDataSourceModel struct {
	Id              types.String `tfsdk:"id"`
	EnvironmentId   types.String `tfsdk:"environment_id"`
	DigitalWalletId types.String `tfsdk:"digital_wallet_id"`
	ApplicationId   types.String `tfsdk:"application_id"`
	AppOpenUrl      types.String `tfsdk:"app_open_url"`
	Name            types.String `tfsdk:"name"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &DigitalWalletApplicationDataSource{}
)

// New Object
func NewDigitalWalletApplicationDataSource() datasource.DataSource {
	return &DigitalWalletApplicationDataSource{}
}

// Metadata
func (r *DigitalWalletApplicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_digital_wallet_application"
}

func (r *DigitalWalletApplicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	// schema definition
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a PingOne Credentials digital wallet application.\n\n" +
			"The service controls the relationship between the customer's digital wallet application, which communicates with users' digital wallets, and a customer's PingOne application.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the credential digital wallet app exists."),
			),

			"digital_wallet_id": schema.StringAttribute{
				Description: "Identifier (UUID) associated with the credential digital wallet application.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("name"),
						path.MatchRelative().AtParent().AtName("application_id"),
					),
					verify.P1ResourceIDValidator(),
				},
			},

			"application_id": schema.StringAttribute{
				Description: "The identifier (UUID) of the PingOne application associated with the digital wallet application.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("name"),
						path.MatchRelative().AtParent().AtName("digital_wallet_id"),
					),
					verify.P1ResourceIDValidator(),
				},
			},

			"name": schema.StringAttribute{
				Description: "The name associated with the digital wallet application.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("application_id"),
						path.MatchRelative().AtParent().AtName("digital_wallet_id"),
					),
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"app_open_url": schema.StringAttribute{
				Description: "The URL sent in notifications to the user to communicate with the service.",
				Computed:    true,
			},
		},
	}
}

func (r *DigitalWalletApplicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *DigitalWalletApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *DigitalWalletApplicationDataSourceModel

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
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var digitalWalletApp credentials.DigitalWalletApplication

	if !data.DigitalWalletId.IsNull() {
		// Run the API call
		response, diags := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.DigitalWalletAppsApi.ReadOneDigitalWalletApp(ctx, data.EnvironmentId.ValueString(), data.DigitalWalletId.ValueString()).Execute()
			},
			"ReadOneDigitalWalletApplication",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		digitalWalletApp = *response.(*credentials.DigitalWalletApplication)

	} else if !data.ApplicationId.IsNull() {

		// Run the API call
		response, diags := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.DigitalWalletAppsApi.ReadAllDigitalWalletApps(ctx, data.EnvironmentId.ValueString()).Execute()
			},
			"ReadAllDigitalWalletApplication",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		entityArray := response.(*credentials.EntityArray)
		if digitalWalletApps, ok := entityArray.Embedded.GetDigitalWalletApplicationsOk(); ok {

			found := false
			for _, digitalWalletAppItem := range digitalWalletApps {

				if *digitalWalletAppItem.GetApplication().Id == data.ApplicationId.ValueString() {
					digitalWalletApp = digitalWalletAppItem
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Cannot find digital wallet application from application_id",
					fmt.Sprintf("The application %s for environment %s cannot be found", data.ApplicationId.String(), data.EnvironmentId.String()),
				)
				return
			}

		}

	} else if !data.Name.IsNull() {

		// Run the API call
		response, diags := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.DigitalWalletAppsApi.ReadAllDigitalWalletApps(ctx, data.EnvironmentId.ValueString()).Execute()
			},
			"ReadAllDigitalWalletApplication",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		entityArray := response.(*credentials.EntityArray)
		if digitalWalletApps, ok := entityArray.Embedded.GetDigitalWalletApplicationsOk(); ok {

			found := false
			for _, digitalWalletAppItem := range digitalWalletApps {

				if digitalWalletAppItem.GetName() == data.Name.ValueString() {
					digitalWalletApp = digitalWalletAppItem
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Cannot find digital wallet application from name",
					fmt.Sprintf("The name %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
				)
				return
			}

		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested PingOne Credentials Digital Wallet Application: digital_wallet_id, application_id or name must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&digitalWalletApp)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *DigitalWalletApplicationDataSourceModel) toState(apiObject *credentials.DigitalWalletApplication) diag.Diagnostics {
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
	p.DigitalWalletId = framework.StringToTF(apiObject.GetId())
	p.ApplicationId = framework.StringToTF(*apiObject.GetApplication().Id)
	p.AppOpenUrl = framework.StringToTF(apiObject.GetAppOpenUrl())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	return diags
}
