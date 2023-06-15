package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// Types
type PhoneDeliverySettingsResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type PhoneDeliverySettingsResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	EnvironmentId           types.String `tfsdk:"environment_id"`
	ProviderType            types.String `tfsdk:"provider_type"`
	ProviderCustom          types.Object `tfsdk:"provider_custom"`
	ProviderCustomTwilio    types.Object `tfsdk:"provider_custom_twilio"`
	ProviderCustomSyniverse types.Object `tfsdk:"provider_custom_syniverse"`
	CreatedAt               types.String `tfsdk:"created_at"`
	UpdatedAt               types.String `tfsdk:"updated_at"`
}

type PhoneDeliverySettingsProviderCustomResourceModel struct {
	Authentication types.Object `tfsdk:"authentication"`
	Name           types.String `tfsdk:"name"`
	Numbers        types.Set    `tfsdk:"numbers"`
	Requests       types.Set    `tfsdk:"requests"`
}

type PhoneDeliverySettingsProviderCustomAuthenticationResourceModel struct {
	Method   types.String `tfsdk:"method"`
	Password types.String `tfsdk:"password"`
	Token    types.String `tfsdk:"token"`
	Username types.String `tfsdk:"username"`
}

type PhoneDeliverySettingsProviderCustomNumbersResourceModel struct {
	SupportedCountries types.Set    `tfsdk:"supported_countries"`
	Type               types.String `tfsdk:"type"`
	Selected           types.Bool   `tfsdk:"selected"`
	Available          types.Bool   `tfsdk:"available"`
	Number             types.String `tfsdk:"number"`
	Capabilities       types.Set    `tfsdk:"capabilities"`
}

type PhoneDeliverySettingsProviderCustomRequestsResourceModel struct {
	DeliveryMethod    types.String `tfsdk:"delivery_method"`
	Url               types.String `tfsdk:"url"`
	Method            types.String `tfsdk:"method"`
	Body              types.String `tfsdk:"body"`
	Headers           types.Map    `tfsdk:"headers"`
	BeforeTag         types.String `tfsdk:"before_tag"`
	AfterTag          types.String `tfsdk:"after_tag"`
	PhoneNumberFormat types.String `tfsdk:"phone_number_format"`
}

type PhoneDeliverySettingsProviderCustomTwilioResourceModel struct {
	Sid       types.String `tfsdk:"sid"`
	AuthToken types.String `tfsdk:"auth_token"`
}

type PhoneDeliverySettingsProviderCustomSyniverseResourceModel struct {
	AuthToken types.String `tfsdk:"auth_token"`
}

var (
	customTFObjectTypes = map[string]attr.Type{
		"authentication": types.ObjectType{
			AttrTypes: customAuthenticationTFObjectTypes,
		},
		"name": types.StringType,
		"numbers": types.SetType{ElemType: types.ObjectType{
			AttrTypes: customNumbersTFObjectTypes,
		}},
		"requests": types.SetType{ElemType: types.ObjectType{
			AttrTypes: customRequestsTFObjectTypes,
		}},
	}

	customAuthenticationTFObjectTypes = map[string]attr.Type{
		"method":   types.StringType,
		"password": types.StringType,
		"token":    types.StringType,
		"username": types.StringType,
	}

	customNumbersTFObjectTypes = map[string]attr.Type{
		"available":           types.BoolType,
		"capabilities":        types.SetType{ElemType: types.StringType},
		"number":              types.StringType,
		"selected":            types.BoolType,
		"supported_countries": types.SetType{ElemType: types.StringType},
		"type":                types.StringType,
	}

	customRequestsTFObjectTypes = map[string]attr.Type{
		"after_tag":           types.StringType,
		"before_tag":          types.StringType,
		"body":                types.StringType,
		"delivery_method":     types.StringType,
		"headers":             types.MapType{ElemType: types.StringType},
		"method":              types.StringType,
		"phone_number_format": types.StringType,
		"url":                 types.StringType,
	}

	twilioTFObjectTypes = map[string]attr.Type{
		"auth_token": types.StringType,
		"sid":        types.StringType,
	}

	syniverseTFObjectTypes = map[string]attr.Type{
		"auth_token": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &PhoneDeliverySettingsResource{}
	_ resource.ResourceWithConfigure   = &PhoneDeliverySettingsResource{}
	_ resource.ResourceWithImportState = &PhoneDeliverySettingsResource{}
)

// New Object
func NewPhoneDeliverySettingsResource() resource.Resource {
	return &PhoneDeliverySettingsResource{}
}

// Metadata
func (r *PhoneDeliverySettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_phone_delivery_settings"
}

// Schema.
func (r *PhoneDeliverySettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	providerTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of the phone delivery service.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsProviderEnumValues)

	// Custom provider
	providerCustomDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Required when the `provider` parameter is set to `%s`.  A nested attribute with attributes that describe custom phone delivery settings.", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSPROVIDER_CUSTOM_PROVIDER),
	)

	providerCustomNumbersCapabilitiesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	providerCustomNumbersSupportedCountriesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	providerCustomRequestsAfterTagDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	providerCustomRequestsBeforeTagDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	providerCustomRequestsBodyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	providerCustomRequestsDeliveryMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	providerCustomRequestsMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	providerCustomRequestsPhoneNumberFormatDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	providerCustomRequestsUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	providerCustomNumbersTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The type of phone number.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomNumbersTypeEnumValues)

	// Twilio provider
	providerCustomTwilioDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Required when the `provider` parameter is set to `%s`.  A nested attribute with attributes that describe phone delivery settings for a custom Twilio account.", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSPROVIDER_CUSTOM_TWILIO),
	)

	providerCustomTwilioAuthTokenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The secret key of the Twilio account.",
	).RequiresReplace()

	providerCustomTwilioSidDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The public ID of the Twilio account.",
	).RequiresReplace()

	// Syniverse provider
	providerCustomSyniverseDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Required when the `provider` parameter is set to `%s`.  A nested attribute with attributes that describe phone delivery settings for a custom syniverse account.", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSPROVIDER_CUSTOM_SYNIVERSE),
	)

	providerCustomSyniverseAuthTokenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The secret key of the Syniverse account.",
	).RequiresReplace()

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage SMS/voice delivery settings in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure SMS/voice settings for."),
			),

			"provider_type": schema.StringAttribute{
				Description:         providerTypeDescription.Description,
				MarkdownDescription: providerTypeDescription.MarkdownDescription,
				Required:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsProviderEnumValues)...),
				},
			},

			"provider_custom": schema.SingleNestedAttribute{
				Description:         providerCustomDescription.Description,
				MarkdownDescription: providerCustomDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("The string that specifies the name of the custom provider used to identify in the PingOne platform.").Description,
						Required:    true,
					},

					// "authentication": schema.SingleNestedAttribute{},

					"numbers": schema.SetNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("One or more objects that describe the numbers to use for phone delivery.").Description,
						Optional:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"available": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Optional:    true,
								},

								"capabilities": schema.SetAttribute{
									Description:         providerCustomNumbersCapabilitiesDescription.Description,
									MarkdownDescription: providerCustomNumbersCapabilitiesDescription.MarkdownDescription,
									Required:            true,

									ElementType: types.StringType,

									Validators: []validator.Set{},
								},

								"number": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the phone number, toll-free number or short code.").Description,
									Required:    true,
								},

								"selected": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Optional:    true,
								},

								"supported_countries": schema.SetAttribute{
									Description:         providerCustomNumbersSupportedCountriesDescription.Description,
									MarkdownDescription: providerCustomNumbersSupportedCountriesDescription.MarkdownDescription,
									Required:            true,

									ElementType: types.StringType,

									Validators: []validator.Set{},
								},

								"type": schema.StringAttribute{
									Description:         providerCustomNumbersTypeDescription.Description,
									MarkdownDescription: providerCustomNumbersTypeDescription.MarkdownDescription,
									Required:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomNumbersTypeEnumValues)...),
									},
								},
							},
						},
					},

					"requests": schema.SetNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("One or more objects that describe the outbound custom notification requests.").Description,
						Required:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"after_tag": schema.StringAttribute{
									Description:         providerCustomRequestsAfterTagDescription.Description,
									MarkdownDescription: providerCustomRequestsAfterTagDescription.MarkdownDescription,
									Optional:            true,
								},

								"before_tag": schema.StringAttribute{
									Description:         providerCustomRequestsBeforeTagDescription.Description,
									MarkdownDescription: providerCustomRequestsBeforeTagDescription.MarkdownDescription,
									Optional:            true,
								},

								"body": schema.StringAttribute{
									Description:         providerCustomRequestsBodyDescription.Description,
									MarkdownDescription: providerCustomRequestsBodyDescription.MarkdownDescription,
									Required:            true,
								},

								"delivery_method": schema.StringAttribute{
									Description:         providerCustomRequestsDeliveryMethodDescription.Description,
									MarkdownDescription: providerCustomRequestsDeliveryMethodDescription.MarkdownDescription,
									Required:            true,

									// Validators: []validator.String{
									// 	stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomRequestsDeliveryMethodEnumValues)...),
									// },
								},

								// "headers": types.MapType{ElemType: types.StringType},

								"method": schema.StringAttribute{
									Description:         providerCustomRequestsMethodDescription.Description,
									MarkdownDescription: providerCustomRequestsMethodDescription.MarkdownDescription,
									Required:            true,

									// Validators: []validator.String{
									// 	stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomRequestsDeliveryMethodEnumValues)...),
									// },
								},

								"phone_number_format": schema.StringAttribute{
									Description:         providerCustomRequestsPhoneNumberFormatDescription.Description,
									MarkdownDescription: providerCustomRequestsPhoneNumberFormatDescription.MarkdownDescription,
									Required:            true,

									// Validators: []validator.String{
									// 	stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomRequestsDeliveryMethodEnumValues)...),
									// },
								},

								"url": schema.StringAttribute{
									Description:         providerCustomRequestsUrlDescription.Description,
									MarkdownDescription: providerCustomRequestsUrlDescription.MarkdownDescription,
									Required:            true,

									// Validators: []validator.String{
									// 	stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomRequestsDeliveryMethodEnumValues)...),
									// },
								},
							},
						},
					},
				},
			},

			"provider_custom_twilio": schema.SingleNestedAttribute{
				Description:         providerCustomTwilioDescription.Description,
				MarkdownDescription: providerCustomTwilioDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"auth_token": schema.StringAttribute{
						Description:         providerCustomTwilioAuthTokenDescription.Description,
						MarkdownDescription: providerCustomTwilioAuthTokenDescription.MarkdownDescription,
						Required:            true,
						Sensitive:           true,

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},

					"sid": schema.StringAttribute{
						Description:         providerCustomTwilioSidDescription.Description,
						MarkdownDescription: providerCustomTwilioSidDescription.MarkdownDescription,
						Required:            true,

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},

			"provider_custom_syniverse": schema.SingleNestedAttribute{
				Description:         providerCustomSyniverseDescription.Description,
				MarkdownDescription: providerCustomSyniverseDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"auth_token": schema.StringAttribute{
						Description:         providerCustomSyniverseAuthTokenDescription.Description,
						MarkdownDescription: providerCustomSyniverseAuthTokenDescription.MarkdownDescription,
						Required:            true,
						Sensitive:           true,

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},

			"created_at": schema.StringAttribute{
				Description: "A string that specifies the time the resource was created.",
				Computed:    true,
			},

			"updated_at": schema.StringAttribute{
				Description: "A string that specifies the time the resource was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *PhoneDeliverySettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PhoneDeliverySettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state PhoneDeliverySettingsResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	phoneDeliverySettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.PhoneDeliverySettingsApi.CreatePhoneDeliverySettings(ctx, plan.EnvironmentId.ValueString()).NotificationsSettingsPhoneDeliverySettings(*phoneDeliverySettings).Execute()
		},
		"CreatePhoneDeliverySettings",
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
	resp.Diagnostics.Append(state.toState(response.(*management.NotificationsSettingsPhoneDeliverySettings))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *PhoneDeliverySettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PhoneDeliverySettingsResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
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
			return r.client.PhoneDeliverySettingsApi.ReadOnePhoneDeliverySettings(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOnePhoneDeliverySettings",
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
	resp.Diagnostics.Append(data.toState(response.(*management.NotificationsSettingsPhoneDeliverySettings))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PhoneDeliverySettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PhoneDeliverySettingsResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	phoneDeliverySettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.PhoneDeliverySettingsApi.UpdatePhoneDeliverySettings(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).NotificationsSettingsPhoneDeliverySettings(*phoneDeliverySettings).Execute()
		},
		"UpdatePhoneDeliverySettings",
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
	resp.Diagnostics.Append(state.toState(response.(*management.NotificationsSettingsPhoneDeliverySettings))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *PhoneDeliverySettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PhoneDeliverySettingsResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
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
			r, err := r.client.PhoneDeliverySettingsApi.DeletePhoneDeliverySettings(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeletePhoneDeliverySettings",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *PhoneDeliverySettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/agreement_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *PhoneDeliverySettingsResourceModel) expand(ctx context.Context) (*management.NotificationsSettingsPhoneDeliverySettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NotificationsSettingsPhoneDeliverySettings{
		NotificationsSettingsPhoneDeliverySettingsCustom:          nil,
		NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse: nil,
	}

	if !p.ProviderCustom.IsNull() && !p.ProviderCustom.IsUnknown() {
		var providerPlan PhoneDeliverySettingsProviderCustomResourceModel
		diags.Append(p.ProviderCustom.As(ctx, &providerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		// Expand authentication
		var authenticationPlan PhoneDeliverySettingsProviderCustomAuthenticationResourceModel
		diags.Append(providerPlan.Authentication.As(ctx, &authenticationPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		authentication := management.NewNotificationsSettingsPhoneDeliverySettingsCustomAllOfAuthentication(
			management.EnumNotificationsSettingsPhoneDeliverySettingsCustomAuthMethod(authenticationPlan.Method.ValueString()),
		)

		if !authenticationPlan.Password.IsNull() && !authenticationPlan.Password.IsUnknown() {
			authentication.SetPassword(authenticationPlan.Password.ValueString())
		}

		if !authenticationPlan.Username.IsNull() && !authenticationPlan.Username.IsUnknown() {
			authentication.SetUsername(authenticationPlan.Username.ValueString())
		}

		if !authenticationPlan.Token.IsNull() && !authenticationPlan.Token.IsUnknown() {
			authentication.SetToken(authenticationPlan.Token.ValueString())
		}

		// Expand requests
		requests := make([]management.NotificationsSettingsPhoneDeliverySettingsCustomRequest, 0)

		if !providerPlan.Requests.IsNull() && !providerPlan.Requests.IsUnknown() {
			var requestsPlan []PhoneDeliverySettingsProviderCustomRequestsResourceModel
			diags.Append(providerPlan.Requests.ElementsAs(ctx, &requestsPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			for _, requestPlan := range requestsPlan {

				var headers map[string]string
				diags.Append(requestPlan.Headers.ElementsAs(ctx, &headers, false)...)
				if diags.HasError() {
					return nil, diags
				}

				request := management.NewNotificationsSettingsPhoneDeliverySettingsCustomRequest(
					management.EnumNotificationsSettingsPhoneDeliverySettingsCustomDeliveryMethod(requestPlan.DeliveryMethod.ValueString()),
					requestPlan.Url.ValueString(),
					requestPlan.Body.ValueString(),
					headers,
					management.EnumNotificationsSettingsPhoneDeliverySettingsCustomRequestMethod(requestPlan.Method.ValueString()),
					management.EnumNotificationsSettingsPhoneDeliverySettingsCustomNumberFormat(requestPlan.PhoneNumberFormat.ValueString()),
				)

				requests = append(requests, *request)
			}
		}

		providerData := management.NewNotificationsSettingsPhoneDeliverySettingsCustom(
			management.EnumNotificationsSettingsPhoneDeliverySettingsProvider(p.ProviderType.ValueString()),
			providerPlan.Name.ValueString(),
			requests,
			*authentication,
		)

		if !providerPlan.Numbers.IsNull() && !providerPlan.Numbers.IsUnknown() {
			numbers := make([]management.NotificationsSettingsPhoneDeliverySettingsCustomNumbers, 0)

			providerData.SetNumbers(numbers)
		}

		data.NotificationsSettingsPhoneDeliverySettingsCustom = providerData
	}

	if !p.ProviderCustomTwilio.IsNull() && !p.ProviderCustomTwilio.IsUnknown() {
		var providerPlan PhoneDeliverySettingsProviderCustomTwilioResourceModel
		diags.Append(p.ProviderCustomTwilio.As(ctx, &providerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		numbers := make([]management.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverseAllOfNumbers, 0)

		providerData := management.NewNotificationsSettingsPhoneDeliverySettingsTwilioSyniverse(
			management.EnumNotificationsSettingsPhoneDeliverySettingsProvider(p.ProviderType.ValueString()),
			providerPlan.Sid.ValueString(),
			providerPlan.AuthToken.ValueString(),
			numbers,
		)

		data.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse = providerData
	}

	if !p.ProviderCustomSyniverse.IsNull() && !p.ProviderCustomSyniverse.IsUnknown() {
		var providerPlan PhoneDeliverySettingsProviderCustomSyniverseResourceModel
		diags.Append(p.ProviderCustomSyniverse.As(ctx, &providerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		numbers := make([]management.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverseAllOfNumbers, 0)

		providerData := management.NewNotificationsSettingsPhoneDeliverySettingsTwilioSyniverse(
			management.EnumNotificationsSettingsPhoneDeliverySettingsProvider(p.ProviderType.ValueString()),
			"",
			providerPlan.AuthToken.ValueString(),
			numbers,
		)

		data.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse = providerData
	}

	return &data, diags
}

func (p *PhoneDeliverySettingsResourceModel) toState(apiObject *management.NotificationsSettingsPhoneDeliverySettings) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	apiObjectCommon := management.NotificationsSettingsPhoneDeliverySettingsCommon{}

	if v := apiObject.NotificationsSettingsPhoneDeliverySettingsCustom; v != nil {
		apiObjectCommon = management.NotificationsSettingsPhoneDeliverySettingsCommon{
			Id:          v.Id,
			Environment: v.Environment,
			Provider:    v.Provider,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}

	if v := apiObject.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse; v != nil {
		apiObjectCommon = management.NotificationsSettingsPhoneDeliverySettingsCommon{
			Id:          v.Id,
			Environment: v.Environment,
			Provider:    v.Provider,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}

	p.Id = framework.StringOkToTF(apiObjectCommon.GetIdOk())
	p.EnvironmentId = framework.StringToTF(*apiObjectCommon.GetEnvironment().Id)
	p.ProviderType = framework.EnumOkToTF(apiObjectCommon.GetProviderOk())
	p.CreatedAt = framework.TimeOkToTF(apiObjectCommon.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObjectCommon.GetUpdatedAtOk())

	var d diag.Diagnostics

	p.ProviderCustom, d = p.toStatePhoneDeliverySettingsProviderCustom(apiObject.NotificationsSettingsPhoneDeliverySettingsCustom)
	diags.Append(d...)

	p.ProviderCustomTwilio, d = p.toStatePhoneDeliverySettingsProviderCustomTwilio(apiObject.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse)
	diags.Append(d...)

	p.ProviderCustomSyniverse, d = p.toStatePhoneDeliverySettingsProviderCustomSyniverse(apiObject.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse)
	diags.Append(d...)

	return diags
}

func (p *PhoneDeliverySettingsResourceModel) toStatePhoneDeliverySettingsProviderCustom(apiObject *management.NotificationsSettingsPhoneDeliverySettingsCustom) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(customTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"name": framework.StringOkToTF(apiObject.GetNameOk()),
	}

	var d diag.Diagnostics

	objMap["authentication"], d = phoneDeliverySettingsCustomAuthenticationOkToTF(apiObject.GetAuthenticationOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(customTFObjectTypes), diags
	}

	objMap["numbers"], d = phoneDeliverySettingsCustomNumbersOkToTF(apiObject.GetNumbersOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(customTFObjectTypes), diags
	}

	objMap["requests"], d = phoneDeliverySettingsCustomRequestsOkToTF(apiObject.GetRequestsOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(customTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(customTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func phoneDeliverySettingsCustomAuthenticationOkToTF(apiObject *management.NotificationsSettingsPhoneDeliverySettingsCustomAllOfAuthentication, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(customAuthenticationTFObjectTypes), diags
	}

	returnVar, d := types.ObjectValue(customAuthenticationTFObjectTypes, map[string]attr.Value{
		"method":   framework.EnumOkToTF(apiObject.GetMethodOk()),
		"password": framework.StringOkToTF(apiObject.GetPasswordOk()),
		"token":    framework.StringOkToTF(apiObject.GetTokenOk()),
		"username": framework.StringOkToTF(apiObject.GetUsernameOk()),
	})
	diags.Append(d...)

	return returnVar, diags
}

func phoneDeliverySettingsCustomNumbersOkToTF(apiObject []management.NotificationsSettingsPhoneDeliverySettingsCustomNumbers, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: customNumbersTFObjectTypes}

	if !ok || len(apiObject) == 0 {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		objMap := map[string]attr.Value{
			"supported_countries": framework.StringSetOkToTF(v.GetSupportedCountriesOk()),
			"type":                framework.EnumOkToTF(v.GetTypeOk()),
			"selected":            framework.BoolOkToTF(v.GetSelectedOk()),
			"available":           framework.BoolOkToTF(v.GetAvailableOk()),
			"number":              framework.StringOkToTF(v.GetNumberOk()),
			"capabilities":        framework.EnumSetOkToTF(v.GetCapabilitiesOk()),
		}

		flattenedObj, d := types.ObjectValue(customNumbersTFObjectTypes, objMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func phoneDeliverySettingsCustomRequestsOkToTF(apiObject []management.NotificationsSettingsPhoneDeliverySettingsCustomRequest, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: customRequestsTFObjectTypes}

	if !ok || len(apiObject) == 0 {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		objMap := map[string]attr.Value{
			"delivery_method":     framework.EnumOkToTF(v.GetDeliveryMethodOk()),
			"url":                 framework.StringOkToTF(v.GetUrlOk()),
			"method":              framework.EnumOkToTF(v.GetMethodOk()),
			"body":                framework.StringOkToTF(v.GetBodyOk()),
			"headers":             framework.StringMapOkToTF(v.GetHeadersOk()),
			"before_tag":          framework.StringOkToTF(v.GetBeforeTagOk()),
			"after_tag":           framework.StringOkToTF(v.GetAfterTagOk()),
			"phone_number_format": framework.EnumOkToTF(v.GetPhoneNumberFormatOk()),
		}

		flattenedObj, d := types.ObjectValue(customRequestsTFObjectTypes, objMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func (p *PhoneDeliverySettingsResourceModel) toStatePhoneDeliverySettingsProviderCustomTwilio(apiObject *management.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(twilioTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(twilioTFObjectTypes, map[string]attr.Value{
		"sid":        framework.StringOkToTF(apiObject.GetSidOk()),
		"auth_token": framework.StringOkToTF(apiObject.GetAuthTokenOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *PhoneDeliverySettingsResourceModel) toStatePhoneDeliverySettingsProviderCustomSyniverse(apiObject *management.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(syniverseTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(syniverseTFObjectTypes, map[string]attr.Value{
		"auth_token": framework.StringOkToTF(apiObject.GetAuthTokenOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
