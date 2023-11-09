package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type IdentityProviderResource serviceClientType

type IdentityProviderResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	EnvironmentId            types.String `tfsdk:"environment_id"`
	Name                     types.String `tfsdk:"name"`
	Description              types.String `tfsdk:"description"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	RegistrationPopulationId types.String `tfsdk:"registration_population_id"`
	LoginButtonIcon          types.List   `tfsdk:"login_button_icon"`
	Icon                     types.List   `tfsdk:"icon"`
	Facebook                 types.List   `tfsdk:"facebook"`
	Google                   types.List   `tfsdk:"google"`
	LinkedIn                 types.List   `tfsdk:"linkedin"`
	Yahoo                    types.List   `tfsdk:"yahoo"`
	Amazon                   types.List   `tfsdk:"amazon"`
	Twitter                  types.List   `tfsdk:"twitter"`
	Apple                    types.List   `tfsdk:"apple"`
	Paypal                   types.List   `tfsdk:"paypal"`
	Microsoft                types.List   `tfsdk:"microsoft"`
	Github                   types.List   `tfsdk:"github"`
	OpenIDConnect            types.List   `tfsdk:"openid_connect"`
	Saml                     types.List   `tfsdk:"saml"`
}

type IdentityProviderClientIdClientSecretResourceModel struct {
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

type IdentityProviderImageResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Href types.String `tfsdk:"href"`
}

type IdentityProviderLoginButtonIcon IdentityProviderImageResourceModel

type IdentityProviderIcon IdentityProviderImageResourceModel

type IdentityProviderFacebookResourceModel struct {
	AppId     types.String `tfsdk:"app_id"`
	AppSecret types.String `tfsdk:"app_secret"`
}

type IdentityProviderGoogleResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderLinkedInResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderYahooResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderAmazonResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderTwitterResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderAppleResourceModel struct {
	TeamId                 types.String `tfsdk:"team_id"`
	KeyId                  types.String `tfsdk:"key_id"`
	ClientId               types.String `tfsdk:"client_id"`
	ClientSecretSigningKey types.String `tfsdk:"client_secret_signing_key"`
}

type IdentityProviderPaypalResourceModel struct {
	ClientId          types.String `tfsdk:"client_id"`
	ClientSecret      types.String `tfsdk:"client_secret"`
	ClientEnvironment types.String `tfsdk:"client_environment"`
}

type IdentityProviderMicrosoftResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderGithubResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderOIDCResourceModel struct {
	AuthorizationEndpoint   types.String `tfsdk:"authorization_endpoint"`
	ClientId                types.String `tfsdk:"client_id"`
	ClientSecret            types.String `tfsdk:"client_secret"`
	DiscoveryEndpoint       types.String `tfsdk:"discovery_endpoint"`
	Issuer                  types.String `tfsdk:"issuer"`
	JwksEndpoint            types.String `tfsdk:"jwks_endpoint"`
	Scopes                  types.Set    `tfsdk:"scopes"`
	TokenEndpoint           types.String `tfsdk:"token_endpoint"`
	TokenEndpointAuthMethod types.String `tfsdk:"token_endpoint_auth_method"`
	UserinfoEndpoint        types.String `tfsdk:"userinfo_endpoint"`
}

type IdentityProviderSAMLResourceModel struct {
	AuthenticationRequestSigned   types.Bool   `tfsdk:"authentication_request_signed"`
	IdpEntityId                   types.String `tfsdk:"idp_entity_id"`
	SpEntityId                    types.String `tfsdk:"sp_entity_id"`
	IdpVerificationCertificateIds types.Set    `tfsdk:"idp_verification_certificate_ids"`
	SpSigningKeyId                types.String `tfsdk:"sp_signing_key_id"`
	SsoBinding                    types.String `tfsdk:"sso_binding"`
	SsoEndpoint                   types.String `tfsdk:"sso_endpoint"`
	SloBinding                    types.String `tfsdk:"slo_binding"`
	SloEndpoint                   types.String `tfsdk:"slo_endpoint"`
	SloResponseEndpoint           types.String `tfsdk:"slo_response_endpoint"`
	SloWindow                     types.Int64  `tfsdk:"slo_window"`
}

// Framework interfaces
var (
	_ resource.Resource                = &IdentityProviderResource{}
	_ resource.ResourceWithConfigure   = &IdentityProviderResource{}
	_ resource.ResourceWithImportState = &IdentityProviderResource{}
)

// New Object
func NewIdentityProviderResource() resource.Resource {
	return &IdentityProviderResource{}
}

// Metadata
func (r *IdentityProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_provider"
}

// Schema.
func (r *IdentityProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	attributes := map[string]schema.Attribute{
		"id": framework.Attr_ID(),

		"environment_id": framework.Attr_LinkID(
			framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that contains the application to assign the admin role to."),
		),

		"application_id": framework.Attr_LinkID(
			framework.SchemaAttributeDescriptionFromMarkdown("The ID of an application to assign an admin role to."),
		),

		"role_id": framework.Attr_LinkID(
			framework.SchemaAttributeDescriptionFromMarkdown("The ID of an admin role to assign to the application."),
		),

		"read_only": schema.BoolAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A flag to show whether the admin role assignment is read only or can be changed.").Description,
			Computed:    true,
		},
	}

	utils.MergeSchemaAttributeMaps(attributes, service.RoleAssignmentScopeSchema(), true)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Identity Providers in an environment.",

		Attributes: attributes,
	}
}

func (r *IdentityProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IdentityProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state IdentityProviderResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	identityProvider, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.IdentityProvider
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.IdentityProvidersApi.CreateIdentityProvider(ctx, plan.EnvironmentId.ValueString()).IdentityProvider(*identityProvider).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateIdentityProvider",
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

func (r *IdentityProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *IdentityProviderResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	var response *management.IdentityProvider
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.IdentityProvidersApi.ReadOneIdentityProvider(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneIdentityProvider",
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

func (r *IdentityProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state IdentityProviderResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	identityProvider, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.IdentityProvider
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.IdentityProvidersApi.UpdateIdentityProvider(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).IdentityProvider(*identityProvider).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateIdentityProvider",
		framework.DefaultCustomError,
		nil,
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

func (r *IdentityProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *IdentityProviderResourceModel

	if r.Client.ManagementAPIClient == nil {
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
			fR, fErr := r.Client.ManagementAPIClient.IdentityProvidersApi.DeleteIdentityProvider(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteIdentityProvider",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *IdentityProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "identity_provider_id",
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

func (p *IdentityProviderResourceModel) expand() (*management.IdentityProvider, diag.Diagnostics) {
	var diags diag.Diagnostics

	common := *management.NewIdentityProviderCommon(p.Enabled.ValueBool(), p.Name.ValueString(), management.ENUMIDENTITYPROVIDEREXT_OPENID_CONNECT)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		common.SetDescription(p.Description.ValueString())
	}

	if !p.RegistrationPopulationId.IsNull() && !p.RegistrationPopulationId.IsUnknown() {
		registrationPopulation := *management.NewIdentityProviderCommonRegistrationPopulation()
		registrationPopulation.SetId(p.RegistrationPopulationId.ValueString())
		registration := *management.NewIdentityProviderCommonRegistration()
		registration.SetPopulation(registrationPopulation)
		common.SetRegistration(registration)
	}

	if !p.LoginButtonIcon.IsNull() && !p.LoginButtonIcon.IsUnknown() {
		if j, okJ := v.([]interface{}); okJ && j != nil && len(j) > 0 {
			attrs := j[0].(map[string]interface{})
			icon := *management.NewIdentityProviderCommonLoginButtonIcon()
			icon.SetId(attrs["id"].(string))
			icon.SetHref(attrs["href"].(string))
			common.SetLoginButtonIcon(icon)
		}
	}

	if !p.Icon.IsNull() && !p.Icon.IsUnknown() {
		if j, okJ := v.([]interface{}); okJ && j != nil && len(j) > 0 {
			attrs := j[0].(map[string]interface{})
			icon := *management.NewIdentityProviderCommonIcon()
			icon.SetId(attrs["id"].(string))
			icon.SetHref(attrs["href"].(string))
			common.SetIcon(icon)
		}
	}

	data := &management.IdentityProvider{}

	processedCount := 0

	if !p.Facebook.IsNull() && !p.Facebook.IsUnknown() {
		data.IdentityProviderFacebook, diags = expandIdPFacebook(v.([]interface{}), common)
		processedCount += 1
	}

	if !p.Google.IsNull() && !p.Google.IsUnknown() {
		data.IdentityProviderClientIDClientSecret, diags = expandIdPGoogle(v.([]interface{}), common)
		processedCount += 1
	}

	if !p.LinkedIn.IsNull() && !p.LinkedIn.IsUnknown() {
		data.IdentityProviderClientIDClientSecret, diags = expandIdPLinkedIn(v.([]interface{}), common)
		processedCount += 1
	}

	if !p.Yahoo.IsNull() && !p.Yahoo.IsUnknown() {
		data.IdentityProviderClientIDClientSecret, diags = expandIdPYahoo(v.([]interface{}), common)
		processedCount += 1
	}

	if !p.Amazon.IsNull() && !p.Amazon.IsUnknown() {
		data.IdentityProviderClientIDClientSecret, diags = expandIdPAmazon(v.([]interface{}), common)
		processedCount += 1
	}

	if !p.Twitter.IsNull() && !p.Twitter.IsUnknown() {
		data.IdentityProviderClientIDClientSecret, diags = expandIdPTwitter(v.([]interface{}), common)
		processedCount += 1
	}

	if !p.Apple.IsNull() && !p.Apple.IsUnknown() {
		data.IdentityProviderApple, diags = expandIdPApple(v.([]interface{}), common)
		processedCount += 1
	}

	if !p.Paypal.IsNull() && !p.Paypal.IsUnknown() {
		data.IdentityProviderPaypal, diags = expandIdPPaypal(v.([]interface{}), common)
		processedCount += 1
	}

	if !p.Microsoft.IsNull() && !p.Microsoft.IsUnknown() {
		data.IdentityProviderClientIDClientSecret, diags = expandIdPMicrosoft(v.([]interface{}), common)
		processedCount += 1
	}

	if !p.Github.IsNull() && !p.Github.IsUnknown() {
		data.IdentityProviderClientIDClientSecret, diags = expandIdPGithub(v.([]interface{}), common)
		processedCount += 1
	}

	if !p.OpenIDConnect.IsNull() && !p.OpenIDConnect.IsUnknown() {
		data.IdentityProviderOIDC, diags = expandIdPOIDC(v.([]interface{}), common)
		processedCount += 1
	}

	if !p.Saml.IsNull() && !p.Saml.IsUnknown() {
		data.IdentityProviderSAML, diags = expandIdPSAML(p.Saml, common)
		processedCount += 1
	}

	if processedCount > 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "More than one identity provider type configured.  This is not supported.",
		})
		return nil, diags
	} else if processedCount == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "No identity provider types configured.  This is not supported.",
		})
		return nil, diags
	}

	return data, diags
}

func (p *IdentityProviderResourceModel) toState(apiObject *management.IdentityProvider) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringOkToTF(apiObject.Environment.GetIdOk())
	p.RoleId = framework.StringOkToTF(apiObject.Role.GetIdOk())
	p.ReadOnly = framework.BoolOkToTF(apiObject.GetReadOnlyOk())

	p.ScopeEnvironmentId, p.ScopeOrganizationId, p.ScopePopulationId = service.RoleAssignmentScopeOkToTF(apiObject.GetScopeOk())

	return diags
}
