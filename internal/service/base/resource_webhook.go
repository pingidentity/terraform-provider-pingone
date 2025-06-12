// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type WebhookResource serviceClientType

type webhookResourceModelV1 struct {
	Id                     pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId          pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                   types.String                 `tfsdk:"name"`
	Enabled                types.Bool                   `tfsdk:"enabled"`
	HttpEndpointUrl        types.String                 `tfsdk:"http_endpoint_url"`
	HttpEndpointHeaders    types.Map                    `tfsdk:"http_endpoint_headers"`
	VerifyTLSCertificates  types.Bool                   `tfsdk:"verify_tls_certificates"`
	TLSClientAuthKeyPairId pingonetypes.ResourceIDValue `tfsdk:"tls_client_auth_key_pair_id"`
	Format                 types.String                 `tfsdk:"format"`
	FilterOptions          types.Object                 `tfsdk:"filter_options"`
}

type webhookFilterOptionsResourceModelV1 struct {
	IncludedActionTypes    types.Set  `tfsdk:"included_action_types"`
	IncludedApplicationIds types.Set  `tfsdk:"included_application_ids"`
	IncludedPopulationIds  types.Set  `tfsdk:"included_population_ids"`
	IncludedTags           types.Set  `tfsdk:"included_tags"`
	IPAddressExposed       types.Bool `tfsdk:"ip_address_exposed"`
	UseragentExposed       types.Bool `tfsdk:"useragent_exposed"`
}

var (
	webhookFilterOptionsTFObjectTypes = map[string]attr.Type{
		"included_action_types":    types.SetType{ElemType: types.StringType},
		"included_application_ids": types.SetType{ElemType: pingonetypes.ResourceIDType{}},
		"included_population_ids":  types.SetType{ElemType: pingonetypes.ResourceIDType{}},
		"included_tags":            types.SetType{ElemType: types.StringType},
		"ip_address_exposed":       types.BoolType,
		"useragent_exposed":        types.BoolType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &WebhookResource{}
	_ resource.ResourceWithConfigure   = &WebhookResource{}
	_ resource.ResourceWithImportState = &WebhookResource{}
)

// New Object
func NewWebhookResource() resource.Resource {
	return &WebhookResource{}
}

// Metadata
func (r *WebhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

// Schema
func (r *WebhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether a created or updated webhook should be active or suspended. A suspended state (`\"enabled\":false`) accumulates all matched events, but these events are not delivered until the webhook becomes active again (`\"enabled\":true`). For suspended webhooks, events accumulate for a maximum of two weeks. Events older than two weeks are deleted. Restarted webhooks receive the saved events (up to two weeks from the restart date).",
	).DefaultValue("false")

	httpEndpointHeaders := framework.SchemaAttributeDescriptionFromMarkdown(
		"A map that specifies the headers applied to the outbound request (for example, `Authorization` `Basic usernamepassword`. The purpose of these headers is for the HTTPS endpoint to authenticate the PingOne service, ensuring that the information from PingOne is from a trusted source.",
	)

	verifyTlsCertificatesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether a certificates should be verified. If this property's value is set to `false`, then all certificates are trusted. (Setting this property's value to false introduces a security risk.)",
	).DefaultValue("true")

	tlsClientAuthKeyPairIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the PingOne resource ID of a key to be used for outbound mutual TLS (mTLS) authentication.  This key is used as a client credential to authenticate the webhook.  When using the `pingone_key` resource, the key must have a `usage_type` of `OUTBOUND_MTLS`.  If this property is set, `verify_tls_certificates` must be set to `true`.",
	).AppendMarkdownString("Value must be a valid PingOne resource ID.")

	formatDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies one of the supported webhook formats.",
	).AllowedValuesEnum(management.AllowedEnumSubscriptionFormatEnumValues)

	filterOptionsIncludedActionTypesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A non-empty list that specifies the list of action types that should be matched for the webhook.\n\nRefer to the [PingOne API Reference - Subscription Action Types](https://apidocs.pingidentity.com/pingone/platform/v1/api/#subscription-action-types) documentation for a full list of configurable action types.",
	)

	filterOptionsIncludedApplicationIDsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An array that specifies the list of applications (by ID) whose events are monitored by the webhook (maximum of 10 IDs in the array). If a list of applications is not provided, events are monitored for all applications in the environment.",
	).AppendMarkdownString("Values must be valid PingOne resource IDs.")

	filterOptionsIncludedPopulationIDsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An array that specifies the list of populations (by ID) whose events are monitored by the webhook (maximum of 10 IDs in the array). This property matches events for users in the specified populations, as opposed to events generated in which the user in one of the populations is the actor.",
	).AppendMarkdownString("Values must be valid PingOne resource IDs.")

	filterOptionsIncludedTagsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An array of tags that events must have to be monitored by the webhook. If tags are not specified, all events are monitored.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMSUBSCRIPTIONFILTERINCLUDEDTAGS_ADMIN_IDENTITY_EVENT): "Identifies the event as the action of an administrator on other administrators",
	})

	filterOptionsIPAddressExposedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the IP address of an actor should be present in the source section of the event.",
	).DefaultValue("false")

	filterOptionsUseragentExposedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the User-Agent HTTP header of an event should be present in the source section of the event.",
	).DefaultValue("false")

	const attrMinLength = 1
	const attrFilterOptionsIncludedIDsMaxLength = 10

	resp.Schema = schema.Schema{

		Version: 1,

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Webhooks / Data Subscriptions.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the webhook in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the webhook name.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"enabled": schema.BoolAttribute{
				MarkdownDescription: enabledDescription.MarkdownDescription,
				Description:         enabledDescription.Description,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"http_endpoint_url": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a valid HTTPS URL to which event messages are sent.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^https:\/\/.*`), "Must be a valid HTTPS URL"),
				},
			},

			"http_endpoint_headers": schema.MapAttribute{
				Description:         httpEndpointHeaders.Description,
				MarkdownDescription: httpEndpointHeaders.MarkdownDescription,
				Optional:            true,

				ElementType: types.StringType,
			},

			"verify_tls_certificates": schema.BoolAttribute{
				MarkdownDescription: verifyTlsCertificatesDescription.MarkdownDescription,
				Description:         verifyTlsCertificatesDescription.Description,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(true),
			},

			"tls_client_auth_key_pair_id": schema.StringAttribute{
				Description:         tlsClientAuthKeyPairIdDescription.Description,
				MarkdownDescription: tlsClientAuthKeyPairIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"format": schema.StringAttribute{
				Description:         formatDescription.Description,
				MarkdownDescription: formatDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumSubscriptionFormatEnumValues)...),
				},
			},

			"filter_options": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies the PingOne platform event filters to be included to trigger this webhook.").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"included_action_types": schema.SetAttribute{
						Description:         filterOptionsIncludedActionTypesDescription.Description,
						MarkdownDescription: filterOptionsIncludedActionTypesDescription.MarkdownDescription,
						Required:            true,

						ElementType: types.StringType,

						Validators: []validator.Set{
							setvalidator.SizeAtLeast(attrMinLength),
							setvalidator.ValueStringsAre(
								stringvalidator.LengthAtLeast(attrMinLength),
							),
						},
					},

					"included_application_ids": schema.SetAttribute{
						Description:         filterOptionsIncludedApplicationIDsDescription.Description,
						MarkdownDescription: filterOptionsIncludedApplicationIDsDescription.MarkdownDescription,
						Optional:            true,

						ElementType: pingonetypes.ResourceIDType{},

						Validators: []validator.Set{
							setvalidator.SizeAtMost(attrFilterOptionsIncludedIDsMaxLength),
						},
					},

					"included_population_ids": schema.SetAttribute{
						Description:         filterOptionsIncludedPopulationIDsDescription.Description,
						MarkdownDescription: filterOptionsIncludedPopulationIDsDescription.MarkdownDescription,
						Optional:            true,

						ElementType: pingonetypes.ResourceIDType{},

						Validators: []validator.Set{
							setvalidator.SizeAtMost(attrFilterOptionsIncludedIDsMaxLength),
						},
					},

					"included_tags": schema.SetAttribute{
						Description:         filterOptionsIncludedTagsDescription.Description,
						MarkdownDescription: filterOptionsIncludedTagsDescription.MarkdownDescription,
						Optional:            true,

						ElementType: types.StringType,

						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumSubscriptionFilterIncludedTagsEnumValues)...),
							),
						},
					},

					"ip_address_exposed": schema.BoolAttribute{
						Description:         filterOptionsIPAddressExposedDescription.Description,
						MarkdownDescription: filterOptionsIPAddressExposedDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},

					"useragent_exposed": schema.BoolAttribute{
						Description:         filterOptionsUseragentExposedDescription.Description,
						MarkdownDescription: filterOptionsUseragentExposedDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},
				},
			},
		},
	}
}

func (r *WebhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(legacysdk.ResourceType)
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

func (r *WebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state webhookResourceModelV1

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
	subscription, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.Subscription
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.SubscriptionsWebhooksApi.CreateSubscription(ctx, plan.EnvironmentId.ValueString()).Subscription(*subscription).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateSubscription",
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

func (r *WebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *webhookResourceModelV1

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
	var response *management.Subscription
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.SubscriptionsWebhooksApi.ReadOneSubscription(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneSubscription",
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

func (r *WebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state webhookResourceModelV1

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
	subscription, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.Subscription
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.SubscriptionsWebhooksApi.UpdateSubscription(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Subscription(*subscription).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateSubscription",
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

func (r *WebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *webhookResourceModelV1

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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.SubscriptionsWebhooksApi.DeleteSubscription(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteSubscription",
		legacysdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *WebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "webhook_subscription_id",
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

func (p *webhookResourceModelV1) expand(ctx context.Context) (*management.Subscription, diag.Diagnostics) {
	var diags diag.Diagnostics

	httpEndpoint := *management.NewSubscriptionHttpEndpoint(p.HttpEndpointUrl.ValueString())

	if !p.HttpEndpointHeaders.IsNull() && !p.HttpEndpointHeaders.IsUnknown() {
		var headersPlan map[string]string
		diags.Append(p.HttpEndpointHeaders.ElementsAs(ctx, &headersPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		httpEndpoint.SetHeaders(headersPlan)
	}

	var filterOptionsPlan webhookFilterOptionsResourceModelV1
	diags.Append(p.FilterOptions.As(ctx, &filterOptionsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	var filterOptions *management.SubscriptionFilterOptions
	var d diag.Diagnostics

	filterOptions, d = filterOptionsPlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := management.NewSubscription(
		p.Enabled.ValueBool(),
		*filterOptions,
		management.EnumSubscriptionFormat(p.Format.ValueString()),
		httpEndpoint,
		p.Name.ValueString(),
		p.VerifyTLSCertificates.ValueBool(),
	)

	if !p.TLSClientAuthKeyPairId.IsNull() && !p.TLSClientAuthKeyPairId.IsUnknown() {
		keyPair := management.NewSubscriptionTlsClientAuthKeyPair()
		keyPair.SetId(p.TLSClientAuthKeyPairId.ValueString())

		data.SetTlsClientAuthKeyPair(*keyPair)
	}

	return data, diags
}

func (p *webhookFilterOptionsResourceModelV1) expand(ctx context.Context) (*management.SubscriptionFilterOptions, diag.Diagnostics) {
	var diags diag.Diagnostics

	var includedActionTypesPlan []types.String
	diags.Append(p.IncludedActionTypes.ElementsAs(ctx, &includedActionTypesPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	includedActionTypes, d := framework.TFTypeStringSliceToStringSlice(includedActionTypesPlan, path.Root("filter_options").AtName("included_action_types"))
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := management.NewSubscriptionFilterOptions(includedActionTypes)

	if !p.IncludedApplicationIds.IsNull() && !p.IncludedApplicationIds.IsUnknown() {
		var typePlan []pingonetypes.ResourceIDValue
		diags.Append(p.IncludedApplicationIds.ElementsAs(ctx, &typePlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		typesStr, d := framework.TFTypePingOneResourceIDSliceToStringSlice(typePlan, path.Root("filter_options").AtName("included_application_ids"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		objList := make([]management.SubscriptionFilterOptionsIncludedApplicationsInner, 0)
		for _, v := range typesStr {
			objList = append(objList, *management.NewSubscriptionFilterOptionsIncludedApplicationsInner(v))
		}

		data.SetIncludedApplications(objList)
	}

	if !p.IncludedPopulationIds.IsNull() && !p.IncludedPopulationIds.IsUnknown() {
		var typePlan []pingonetypes.ResourceIDValue
		diags.Append(p.IncludedPopulationIds.ElementsAs(ctx, &typePlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		typesStr, d := framework.TFTypePingOneResourceIDSliceToStringSlice(typePlan, path.Root("filter_options").AtName("included_population_ids"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		objList := make([]management.SubscriptionFilterOptionsIncludedApplicationsInner, 0)
		for _, v := range typesStr {
			objList = append(objList, *management.NewSubscriptionFilterOptionsIncludedApplicationsInner(v))
		}

		data.SetIncludedPopulations(objList)
	}

	if !p.IncludedTags.IsNull() && !p.IncludedTags.IsUnknown() {
		var typePlan []types.String
		diags.Append(p.IncludedTags.ElementsAs(ctx, &typePlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		typesStr, d := framework.TFTypeStringSliceToStringSlice(typePlan, path.Root("filter_options").AtName("included_tags"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		objList := make([]management.EnumSubscriptionFilterIncludedTags, 0)
		for _, v := range typesStr {
			objList = append(objList, management.EnumSubscriptionFilterIncludedTags(v))
		}

		data.SetIncludedTags(objList)
	}

	if !p.IPAddressExposed.IsNull() && !p.IPAddressExposed.IsUnknown() {
		data.SetIpAddressExposed(p.IPAddressExposed.ValueBool())
	}

	if !p.UseragentExposed.IsNull() && !p.UseragentExposed.IsUnknown() {
		data.SetUserAgentExposed(p.UseragentExposed.ValueBool())
	}

	return data, diags
}

func (p *webhookResourceModelV1) toState(apiObject *management.Subscription) diag.Diagnostics {
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
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())

	if v, ok := apiObject.GetHttpEndpointOk(); ok {
		p.HttpEndpointUrl = framework.StringOkToTF(v.GetUrlOk())
		p.HttpEndpointHeaders = framework.StringMapOkToTF(v.GetHeadersOk())
	} else {
		p.HttpEndpointUrl = types.StringNull()
		p.HttpEndpointHeaders = types.MapNull(types.StringType)
	}

	p.VerifyTLSCertificates = framework.BoolOkToTF(apiObject.GetVerifyTlsCertificatesOk())

	p.TLSClientAuthKeyPairId = pingonetypes.NewResourceIDNull()
	if v, ok := apiObject.GetTlsClientAuthKeyPairOk(); ok {
		p.TLSClientAuthKeyPairId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	}

	p.Format = framework.EnumOkToTF(apiObject.GetFormatOk())

	var d diag.Diagnostics
	p.FilterOptions, d = toStateWebhookFilterOptions(apiObject.GetFilterOptionsOk())
	diags.Append(d...)

	return diags
}

func toStateWebhookFilterOptions(v *management.SubscriptionFilterOptions, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || v == nil {
		return types.ObjectNull(webhookFilterOptionsTFObjectTypes), diags
	}

	var applicationIDSet, populationIDSet basetypes.SetValue

	applicationIDList := make([]string, 0)
	if items, ok := v.GetIncludedApplicationsOk(); ok {
		for _, application := range items {
			applicationIDList = append(applicationIDList, application.GetId())
		}

		applicationIDSet = framework.PingOneResourceIDSetToTF(applicationIDList)
	} else {
		applicationIDSet = types.SetNull(pingonetypes.ResourceIDType{})
	}

	populationIDList := make([]string, 0)
	if items, ok := v.GetIncludedPopulationsOk(); ok {
		for _, population := range items {
			populationIDList = append(populationIDList, population.GetId())
		}

		populationIDSet = framework.PingOneResourceIDSetToTF(populationIDList)
	} else {
		populationIDSet = types.SetNull(pingonetypes.ResourceIDType{})
	}

	objMap := map[string]attr.Value{
		"included_action_types":    framework.StringSetOkToTF(v.GetIncludedActionTypesOk()),
		"included_application_ids": applicationIDSet,
		"included_population_ids":  populationIDSet,
		"included_tags":            framework.EnumSetOkToTF(v.GetIncludedTagsOk()),
		"ip_address_exposed":       framework.BoolOkToTF(v.GetIpAddressExposedOk()),
		"useragent_exposed":        framework.BoolOkToTF(v.GetUserAgentExposedOk()),
	}

	returnVar, d := types.ObjectValue(webhookFilterOptionsTFObjectTypes, objMap)
	diags.Append(d...)

	return returnVar, diags

}
