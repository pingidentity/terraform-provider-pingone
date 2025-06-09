// Copyright © 2025 Ping Identity Corporation

package mfa

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	objectplanmodifierinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectplanmodifier"
	stringplanmodifierinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringplanmodifier"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationPushCredentialResource serviceClientType

type applicationPushCredentialResourceModelV1 struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ApplicationId pingonetypes.ResourceIDValue `tfsdk:"application_id"`
	Fcm           types.Object                 `tfsdk:"fcm"`
	Apns          types.Object                 `tfsdk:"apns"`
	Hms           types.Object                 `tfsdk:"hms"`
}

type applicationPushCredentialFcmResourceModelV1 struct {
	GoogleServiceAccountCredentials jsontypes.Normalized `tfsdk:"google_service_account_credentials"`
}

type applicationPushCredentialApnsResourceModelV1 struct {
	Key             types.String `tfsdk:"key"`
	TeamId          types.String `tfsdk:"team_id"`
	TokenSigningKey types.String `tfsdk:"token_signing_key"`
}

type applicationPushCredentialHmsResourceModelV1 struct {
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

// Framework interfaces
var (
	_ resource.Resource                 = &ApplicationPushCredentialResource{}
	_ resource.ResourceWithConfigure    = &ApplicationPushCredentialResource{}
	_ resource.ResourceWithImportState  = &ApplicationPushCredentialResource{}
	_ resource.ResourceWithUpgradeState = &ApplicationPushCredentialResource{}
)

// New Object
func NewApplicationPushCredentialResource() resource.Resource {
	return &ApplicationPushCredentialResource{}
}

// Metadata
func (r *ApplicationPushCredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_application_push_credential"
}

// Schema.
func (r *ApplicationPushCredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	fcmGoogleServiceAccountCredentialsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string in JSON format that represents the service account credentials of Firebase cloud messaging service.",
	)

	fcmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the credential settings for the Firebase Cloud Messaging service.",
	).RequiresReplaceNestedAttributes()

	apnsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the credential settings for the Apple Push Notification Service.",
	).RequiresReplaceNestedAttributes()

	hmsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the credential settings for Huawei Moble Service push messaging.",
	).RequiresReplaceNestedAttributes()

	resp.Schema = schema.Schema{

		Version: 1,

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage push credentials for a mobile MFA application configured in PingOne.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the application push notification credential in."),
			),

			"application_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the application to create the push notification credential for."),
			),

			"fcm": schema.SingleNestedAttribute{
				Description:         fcmDescription.Description,
				MarkdownDescription: fcmDescription.MarkdownDescription,

				Optional: true,

				Attributes: map[string]schema.Attribute{
					"google_service_account_credentials": schema.StringAttribute{
						Description:         fcmGoogleServiceAccountCredentialsDescription.Description,
						MarkdownDescription: fcmGoogleServiceAccountCredentialsDescription.MarkdownDescription,
						Required:            true,
						Sensitive:           true,

						CustomType: jsontypes.NormalizedType{},

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplaceIf(
								stringplanmodifierinternal.RequiresReplaceIfNowNull(),
								"The attribute has been previously defined.  To nullify the attribute, it must be replaced.",
								"The attribute has been previously defined.  To nullify the attribute, it must be replaced.",
							),
						},
					},
				},

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("fcm"),
						path.MatchRelative().AtParent().AtName("apns"),
						path.MatchRelative().AtParent().AtName("hms"),
					),
				},
			},

			"apns": schema.SingleNestedAttribute{
				Description:         apnsDescription.Description,
				MarkdownDescription: apnsDescription.MarkdownDescription,

				Optional: true,

				Attributes: map[string]schema.Attribute{
					"key": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that Apple uses as an identifier to identify an authentication key.").Description,
						Required:    true,
						Sensitive:   true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"team_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that Apple uses as an identifier to identify teams.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"token_signing_key": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that Apple uses as the authentication token signing key to securely connect to APNS. This is the contents of a pkcs8 file with a private key format.").Description,
						Required:    true,
						Sensitive:   true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},
				},

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("fcm"),
						path.MatchRelative().AtParent().AtName("apns"),
						path.MatchRelative().AtParent().AtName("hms"),
					),
				},
			},

			"hms": schema.SingleNestedAttribute{
				Description:         hmsDescription.Description,
				MarkdownDescription: hmsDescription.MarkdownDescription,

				Optional: true,

				Attributes: map[string]schema.Attribute{
					"client_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that represents the OAuth 2.0 Client ID from the Huawei Developers API console.").Description,
						Required:    true,
						Sensitive:   true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"client_secret": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that represents the client secret associated with the OAuth 2.0 Client ID.").Description,
						Required:    true,
						Sensitive:   true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},
				},

				PlanModifiers: []planmodifier.Object{
					objectplanmodifierinternal.RequiresReplaceIfExistenceChanges(),
				},

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("fcm"),
						path.MatchRelative().AtParent().AtName("apns"),
						path.MatchRelative().AtParent().AtName("hms"),
					),
				},
			},
		},
	}
}

func (r *ApplicationPushCredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationPushCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state applicationPushCredentialResourceModelV1

	if r.Client.MFAAPIClient == nil {
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
	applicationPushCredential, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *mfa.MFAPushCredentialResponse
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.ApplicationsApplicationMFAPushCredentialsApi.CreateMFAPushCredential(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).MFAPushCredentialRequest(*applicationPushCredential).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateMFAPushCredential",
		legacysdk.DefaultCustomError,
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

func (r *ApplicationPushCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *applicationPushCredentialResourceModelV1

	if r.Client.MFAAPIClient == nil {
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
	var response *mfa.MFAPushCredentialResponse
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.ApplicationsApplicationMFAPushCredentialsApi.ReadOneMFAPushCredential(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneMFAPushCredential",
		legacysdk.CustomErrorResourceNotFoundWarning,
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

func (r *ApplicationPushCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state applicationPushCredentialResourceModelV1

	if r.Client.MFAAPIClient == nil {
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
	applicationPushCredential, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *mfa.MFAPushCredentialResponse
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.ApplicationsApplicationMFAPushCredentialsApi.UpdateMFAPushCredential(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), plan.Id.ValueString()).MFAPushCredentialRequest(*applicationPushCredential).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateMFAPushCredential",
		legacysdk.DefaultCustomError,
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

func (r *ApplicationPushCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *applicationPushCredentialResourceModelV1

	if r.Client.MFAAPIClient == nil {
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.MFAAPIClient.ApplicationsApplicationMFAPushCredentialsApi.DeleteMFAPushCredential(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteMFAPushCredential",
		legacysdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApplicationPushCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "application_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "push_credential_id",
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

func (p *applicationPushCredentialResourceModelV1) expand(ctx context.Context) (*mfa.MFAPushCredentialRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := &mfa.MFAPushCredentialRequest{}

	if !p.Fcm.IsNull() && !p.Fcm.IsUnknown() {
		var plan applicationPushCredentialFcmResourceModelV1
		diags.Append(p.Fcm.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		if !plan.GoogleServiceAccountCredentials.IsNull() && !plan.GoogleServiceAccountCredentials.IsUnknown() {
			data.MFAPushCredentialFCMHTTPV1 = mfa.NewMFAPushCredentialFCMHTTPV1(
				mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_FCM_HTTP_V1,
				plan.GoogleServiceAccountCredentials.ValueString(),
			)
		}
	}

	if !p.Apns.IsNull() && !p.Apns.IsUnknown() {
		var plan applicationPushCredentialApnsResourceModelV1
		diags.Append(p.Apns.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.MFAPushCredentialAPNS = mfa.NewMFAPushCredentialAPNS(
			mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_APNS,
			plan.Key.ValueString(),
			plan.TeamId.ValueString(),
			plan.TokenSigningKey.ValueString(),
		)
	}

	if !p.Hms.IsNull() && !p.Hms.IsUnknown() {
		var plan applicationPushCredentialHmsResourceModelV1
		diags.Append(p.Hms.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.MFAPushCredentialHMS = mfa.NewMFAPushCredentialHMS(
			mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_HMS,
			plan.ClientId.ValueString(),
			plan.ClientSecret.ValueString(),
		)
	}

	return data, diags
}

func (p *applicationPushCredentialResourceModelV1) toState(apiObject *mfa.MFAPushCredentialResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())

	// The rest are credentials not returned from the API and passed through as-is

	return diags
}
