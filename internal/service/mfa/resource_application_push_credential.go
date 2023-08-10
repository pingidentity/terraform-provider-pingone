package mfa

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	listplanmodifierinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/listplanmodifier"
	stringplanmodifierinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringplanmodifier"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type ApplicationPushCredentialResource struct {
	client *mfa.APIClient
	region model.RegionMapping
}

type applicationPushCredentialResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	ApplicationId types.String `tfsdk:"application_id"`
	Fcm           types.List   `tfsdk:"fcm"`
	Apns          types.List   `tfsdk:"apns"`
	Hms           types.List   `tfsdk:"hms"`
}

type applicationPushCredentialFcmResourceModel struct {
	Key                             types.String `tfsdk:"key"`
	GoogleServiceAccountCredentials types.String `tfsdk:"google_service_account_credentials"`
}

type applicationPushCredentialApnsResourceModel struct {
	Key             types.String `tfsdk:"key"`
	TeamId          types.String `tfsdk:"team_id"`
	TokenSigningKey types.String `tfsdk:"token_signing_key"`
}

type applicationPushCredentialHmsResourceModel struct {
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

// Framework interfaces
var (
	_ resource.Resource                = &ApplicationPushCredentialResource{}
	_ resource.ResourceWithConfigure   = &ApplicationPushCredentialResource{}
	_ resource.ResourceWithImportState = &ApplicationPushCredentialResource{}
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

	fcmKeyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that represents the server key of the Firebase cloud messaging service.",
	).ExactlyOneOf([]string{
		"key",
		"google_service_account_credentials",
	})

	fcmGoogleServiceAccountCredentialsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string in JSON format that represents the service account credentials of Firebase cloud messaging service.",
	).ExactlyOneOf([]string{
		"key",
		"google_service_account_credentials",
	})

	resp.Schema = schema.Schema{
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
		},

		Blocks: map[string]schema.Block{

			"fcm": schema.ListNestedBlock{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies the credential settings for the Firebase Cloud Messaging service.").Description,

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description:         fcmKeyDescription.Description,
							MarkdownDescription: fcmKeyDescription.MarkdownDescription,
							Optional:            true,
							Sensitive:           true,
							DeprecationMessage:  "This field is deprecated and will be removed in a future release.  Use `google_service_account_credentials` instead.",

							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIf(
									stringplanmodifierinternal.RequiresReplaceIfNowNull(),
									"The attribute has been previously defined.  To nullify the attribute, this will change the credential type and it must be replaced.",
									"The attribute has been previously defined.  To nullify the attribute, this will change the credential type and it must be replaced.",
								),
							},

							Validators: []validator.String{
								stringvalidator.LengthAtLeast(attrMinLength),
								stringvalidator.ExactlyOneOf(
									path.MatchRelative().AtParent().AtName("key"),
									path.MatchRelative().AtParent().AtName("google_service_account_credentials"),
								),
							},
						},

						"google_service_account_credentials": schema.StringAttribute{
							Description:         fcmGoogleServiceAccountCredentialsDescription.Description,
							MarkdownDescription: fcmGoogleServiceAccountCredentialsDescription.MarkdownDescription,
							Optional:            true,
							Sensitive:           true,

							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIf(
									stringplanmodifierinternal.RequiresReplaceIfNowNull(),
									"The attribute has been previously defined.  To nullify the attribute, it must be replaced.",
									"The attribute has been previously defined.  To nullify the attribute, it must be replaced.",
								),
							},

							Validators: []validator.String{
								stringvalidatorinternal.IsParseableJSON(),
								stringvalidator.ExactlyOneOf(
									path.MatchRelative().AtParent().AtName("key"),
									path.MatchRelative().AtParent().AtName("google_service_account_credentials"),
								),
							},
						},
					},
				},

				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIf(
						listplanmodifierinternal.RequiresReplaceIfNowNull(),
						"The attribute has been previously defined.  To nullify the attribute, this will change the credential type and it must be replaced.",
						"The attribute has been previously defined.  To nullify the attribute, this will change the credential type and it must be replaced.",
					),
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(attrMinLength),
					listvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("fcm"),
						path.MatchRelative().AtParent().AtName("apns"),
						path.MatchRelative().AtParent().AtName("hms"),
					),
				},
			},

			"apns": schema.ListNestedBlock{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies the credential settings for the Apple Push Notification Service.").Description,

				NestedObject: schema.NestedBlockObject{

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
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that Apple uses as the authentication token signing key to securely connect to APNS. This is the contents of a p8 file with a private key format.").Description,
							Required:    true,
							Sensitive:   true,

							Validators: []validator.String{
								stringvalidator.LengthAtLeast(attrMinLength),
							},
						},
					},
				},

				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIf(
						listplanmodifierinternal.RequiresReplaceIfNowNull(),
						"The attribute has been previously defined.  To nullify the attribute, this will change the credential type and it must be replaced.",
						"The attribute has been previously defined.  To nullify the attribute, this will change the credential type and it must be replaced.",
					),
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(attrMinLength),
					listvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("fcm"),
						path.MatchRelative().AtParent().AtName("apns"),
						path.MatchRelative().AtParent().AtName("hms"),
					),
				},
			},

			"hms": schema.ListNestedBlock{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies the credential settings for Huawei Moble Service push messaging.").Description,

				NestedObject: schema.NestedBlockObject{

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
				},

				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIf(
						listplanmodifierinternal.RequiresReplaceIfNowNull(),
						"The attribute has been previously defined.  To nullify the attribute, this will change the credential type and it must be replaced.",
						"The attribute has been previously defined.  To nullify the attribute, this will change the credential type and it must be replaced.",
					),
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(attrMinLength),
					listvalidator.ExactlyOneOf(
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

	preparedClient, err := PrepareClient(ctx, resourceConfig)
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

func (r *ApplicationPushCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state applicationPushCredentialResourceModel

	if r.client == nil {
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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.ApplicationsApplicationMFAPushCredentialsApi.CreateMFAPushCredential(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).MFAPushCredentialRequest(*applicationPushCredential).Execute()
		},
		"CreateMFAPushCredential",
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

func (r *ApplicationPushCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *applicationPushCredentialResourceModel

	if r.client == nil {
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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.ApplicationsApplicationMFAPushCredentialsApi.ReadOneMFAPushCredential(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneMFAPushCredential",
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

func (r *ApplicationPushCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state applicationPushCredentialResourceModel

	if r.client == nil {
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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.ApplicationsApplicationMFAPushCredentialsApi.UpdateMFAPushCredential(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), plan.Id.ValueString()).MFAPushCredentialRequest(*applicationPushCredential).Execute()
		},
		"UpdateMFAPushCredential",
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

func (r *ApplicationPushCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *applicationPushCredentialResourceModel

	if r.client == nil {
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
			r, err := r.client.ApplicationsApplicationMFAPushCredentialsApi.DeleteMFAPushCredential(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteMFAPushCredential",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApplicationPushCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/application_id/push_credential_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), attributes[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *applicationPushCredentialResourceModel) expand(ctx context.Context) (*mfa.MFAPushCredentialRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := &mfa.MFAPushCredentialRequest{}

	if !p.Fcm.IsNull() && !p.Fcm.IsUnknown() {
		var plan []applicationPushCredentialFcmResourceModel
		diags.Append(p.Fcm.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		fcmPlan := plan[0]

		if !fcmPlan.Key.IsNull() && !fcmPlan.Key.IsUnknown() {
			data.MFAPushCredentialFCM = mfa.NewMFAPushCredentialFCM(
				mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_FCM,
				fcmPlan.Key.ValueString(),
			)
		}

		if !fcmPlan.GoogleServiceAccountCredentials.IsNull() && !fcmPlan.GoogleServiceAccountCredentials.IsUnknown() {
			data.MFAPushCredentialFCMHTTPV1 = mfa.NewMFAPushCredentialFCMHTTPV1(
				mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_FCM_HTTP_V1,
				fcmPlan.GoogleServiceAccountCredentials.ValueString(),
			)
		}
	}

	if !p.Apns.IsNull() && !p.Apns.IsUnknown() {
		var plan []applicationPushCredentialApnsResourceModel
		diags.Append(p.Apns.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		apnsPlan := plan[0]

		data.MFAPushCredentialAPNS = mfa.NewMFAPushCredentialAPNS(
			mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_APNS,
			apnsPlan.Key.ValueString(),
			apnsPlan.TeamId.ValueString(),
			apnsPlan.TokenSigningKey.ValueString(),
		)
	}

	if !p.Hms.IsNull() && !p.Hms.IsUnknown() {
		var plan []applicationPushCredentialHmsResourceModel
		diags.Append(p.Hms.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		hmsPlan := plan[0]

		data.MFAPushCredentialHMS = mfa.NewMFAPushCredentialHMS(
			mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_HMS,
			hmsPlan.ClientId.ValueString(),
			hmsPlan.ClientSecret.ValueString(),
		)
	}

	return data, diags
}

func (p *applicationPushCredentialResourceModel) toState(apiObject *mfa.MFAPushCredentialResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())

	// The rest are credentials not returned from the API and passed through as-is

	return diags
}
