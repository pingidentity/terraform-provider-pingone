// Copyright © 2025 Ping Identity Corporation

package credentials

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type DigitalWalletApplicationResource serviceClientType

type DigitalWalletApplicationResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ApplicationId pingonetypes.ResourceIDValue `tfsdk:"application_id"`
	AppOpenUrl    types.String                 `tfsdk:"app_open_url"`
	Name          types.String                 `tfsdk:"name"`
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

	// schema descriptions and validation settings
	const attrMinLength = 1

	// schema definition
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Credentials digital wallet applications.\n\n" +
			"The service controls the relationship between the customer's digital wallet application, which communicates with users' digital wallets, and a customer's PingOne application.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the credential digital wallet application is created and managed."),
			),

			"application_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The identifier (UUID) of the PingOne application associated with the digital wallet application."),
			),

			"app_open_url": schema.StringAttribute{
				Description: "The URL enables deep-linking to the digital wallet application, and is sent in notifications to the user to communicate with the service.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`),
						"Expected value to have a url with scheme of \"https\". A scheme of \"http\" is allowed but not recommended."),
				},
			},

			"name": schema.StringAttribute{
				Description: "The name associated with the digital wallet application.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *DigitalWalletApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state DigitalWalletApplicationResourceModel

	if r.Client == nil || r.Client.CredentialsAPIClient == nil {
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
	digitalWalletApplication, d := plan.expand(ctx, r)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *credentials.DigitalWalletApplication
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.CredentialsAPIClient.DigitalWalletAppsApi.CreateDigitalWalletApp(ctx, plan.EnvironmentId.ValueString()).DigitalWalletApplication(*digitalWalletApplication).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateDigitalWalletApplication",
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

func (r *DigitalWalletApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DigitalWalletApplicationResourceModel

	if r.Client == nil || r.Client.CredentialsAPIClient == nil {
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
	var response *credentials.DigitalWalletApplication
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.CredentialsAPIClient.DigitalWalletAppsApi.ReadOneDigitalWalletApp(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneDigitalWalletApplication",
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

func (r *DigitalWalletApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state DigitalWalletApplicationResourceModel

	if r.Client == nil || r.Client.CredentialsAPIClient == nil {
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
	digitalWalletApplication, d := plan.expand(ctx, r)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *credentials.DigitalWalletApplication
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.CredentialsAPIClient.DigitalWalletAppsApi.UpdateDigitalWalletApp(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).DigitalWalletApplication(*digitalWalletApplication).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateDigitalWalletApplication",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DigitalWalletApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DigitalWalletApplicationResourceModel

	if r.Client == nil || r.Client.CredentialsAPIClient == nil {
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
			fR, fErr := r.Client.CredentialsAPIClient.DigitalWalletAppsApi.DeleteDigitalWalletApp(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteDigitalWalletApplication",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DigitalWalletApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "digital_wallet_application_id",
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

func (p *DigitalWalletApplicationResourceModel) expand(ctx context.Context, r *DigitalWalletApplicationResource) (*credentials.DigitalWalletApplication, diag.Diagnostics) {

	// a digital wallet application is correlated to a Native Application - make sure it exists and is configured properly
	application, diags := confirmParentAppExistsAndIsNative(ctx, r, p.EnvironmentId.ValueString(), p.ApplicationId.ValueString())
	if diags.HasError() {
		return nil, diags
	}

	data := credentials.NewDigitalWalletApplication(*application, p.AppOpenUrl.ValueString(), p.Name.ValueString())
	return data, diags
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

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.ApplicationId = framework.PingOneResourceIDToTF(*apiObject.GetApplication().Id)
	p.AppOpenUrl = framework.StringToTF(apiObject.GetAppOpenUrl())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	return diags
}

func confirmParentAppExistsAndIsNative(ctx context.Context, r *DigitalWalletApplicationResource, environmentId, applicationId string) (*credentials.ObjectApplication, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Run the API call
	var respObject *management.ReadOneApplication200Response
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.ReadOneApplication(ctx, environmentId, applicationId).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environmentId, fO, fR, fErr)
		},
		"ReadOneApplication",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&respObject,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	if respObject == nil {
		diags.AddError(
			"Digital Wallet Parent Application Missing",
			"Application referenced in `application.id` does not exist",
		)
		return nil, diags
	}

	// check if oidc
	if respObject.ApplicationOIDC == nil {
		diags.AddError(
			"Application referenced in `application.id` is not of type OIDC",
			"To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with a package or bundle set.",
		)
		return nil, diags
	}

	// check if native
	if respObject.ApplicationOIDC.GetType() != management.ENUMAPPLICATIONTYPE_NATIVE_APP && respObject.ApplicationOIDC.GetType() != management.ENUMAPPLICATIONTYPE_CUSTOM_APP {
		diags.AddError(
			"Application referenced in `application.id` is OIDC, but is not the required `Native` OIDC application type",
			"To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with a package or bundle set.",
		)
		return nil, diags
	}

	// check if mobile set and package/bundle set
	if _, ok := respObject.ApplicationOIDC.GetMobileOk(); !ok {
		diags.AddError(
			"Application referenced in `application.id` does not contain mobile application configuration",
			"To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with a package or bundle set.",
		)
		return nil, diags
	}

	if v, ok := respObject.ApplicationOIDC.GetMobileOk(); ok {

		_, bundleIDOk := v.GetBundleIdOk()
		_, packageNameOk := v.GetPackageNameOk()

		if !bundleIDOk && !packageNameOk {
			diags.AddError(
				"Application referenced in `application.id` does not contain mobile application configuration",
				"To configure a mobile application in PingOne, the application must be an OIDC application of type `Native`, with a package or bundle set.",
			)
			return nil, diags
		}
	}

	// checks complete - return app object the wallet want
	applicationObject := credentials.NewObjectApplication()
	applicationObject.SetId(applicationId)

	return applicationObject, diags
}
