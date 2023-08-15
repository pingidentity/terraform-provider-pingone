package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type WebhookResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type WebhookResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	EnvironmentId         types.String `tfsdk:"environment_id"`
	Name                  types.String `tfsdk:"name"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	HttpEndpointUrl       types.String `tfsdk:"http_endpoint_url"`
	HttpEndpointHeaders   types.Map    `tfsdk:"http_endpoint_headers"`
	VerifyTLSCertificates types.Bool   `tfsdk:"verify_tls_certificates"`
	Format                types.String `tfsdk:"format"`
	FilterOptions         types.List   `tfsdk:"filter_options"`
}

type WebookFilterOptionsModel struct {
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
		"included_application_ids": types.SetType{ElemType: types.StringType},
		"included_population_ids":  types.SetType{ElemType: types.StringType},
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
	const attrFilterOptionsLimit = 1

	resp.Schema = schema.Schema{
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
					stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Must be a valid HTTPS URL"),
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

			"format": schema.StringAttribute{
				Description:         formatDescription.Description,
				MarkdownDescription: formatDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumSubscriptionFormatEnumValues)...),
				},
			},
		},

		Blocks: map[string]schema.Block{

			"filter_options": schema.ListNestedBlock{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A block that specifies the PingOne platform event filters to be included to trigger this webhook.").Description,

				NestedObject: schema.NestedBlockObject{

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

							ElementType: types.StringType,

							Validators: []validator.Set{
								setvalidator.SizeAtMost(attrFilterOptionsIncludedIDsMaxLength),
								setvalidator.ValueStringsAre(
									verify.P1ResourceIDValidator(),
								),
							},
						},

						"included_population_ids": schema.SetAttribute{
							Description:         filterOptionsIncludedPopulationIDsDescription.Description,
							MarkdownDescription: filterOptionsIncludedPopulationIDsDescription.MarkdownDescription,
							Optional:            true,

							ElementType: types.StringType,

							Validators: []validator.Set{
								setvalidator.SizeAtMost(attrFilterOptionsIncludedIDsMaxLength),
								setvalidator.ValueStringsAre(
									verify.P1ResourceIDValidator(),
								),
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

				Validators: []validator.List{
					listvalidator.SizeAtLeast(attrFilterOptionsLimit),
					listvalidator.SizeAtMost(attrFilterOptionsLimit),
					listvalidator.IsRequired(),
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

func (r *WebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state WebhookResourceModel

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
	subscription, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.Subscription
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.SubscriptionsWebhooksApi.CreateSubscription(ctx, plan.EnvironmentId.ValueString()).Subscription(*subscription).Execute()
		},
		"CreateSubscription",
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

func (r *WebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *WebhookResourceModel

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
	var response *management.Subscription
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.SubscriptionsWebhooksApi.ReadOneSubscription(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneSubscription",
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

func (r *WebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state WebhookResourceModel

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
	subscription, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.Subscription
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.SubscriptionsWebhooksApi.UpdateSubscription(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Subscription(*subscription).Execute()
		},
		"UpdateSubscription",
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

func (r *WebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *WebhookResourceModel

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
			r, err := r.client.SubscriptionsWebhooksApi.DeleteSubscription(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteSubscription",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *WebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/webhook_subscription_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *WebhookResourceModel) expand(ctx context.Context) (*management.Subscription, diag.Diagnostics) {
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

		httpEndpoint.SetHeaders(obj)
	}

	var filterOptions *management.SubscriptionFilterOptions
	var d diag.Diagnostics
	if len(filterOptionsPlan) == 1 {
		filterOptions, d = filterOptionsPlan[0].expand(ctx)
		diags.Append(d...)
	} else {
		d.AddError("Invalid webhook filter options", "Exactly one filter options block must be specified")
	}
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

	return data, diags
}

func (p *WebookFilterOptionsModel) expand(ctx context.Context) (*management.SubscriptionFilterOptions, diag.Diagnostics) {
	var diags diag.Diagnostics

	var includedActionTypes []string
	diags.Append(p.IncludedActionTypes.ElementsAs(ctx, &includedActionTypes, false)...)
	if diags.HasError() {
		return nil, diags
	}

	data := management.NewSubscriptionFilterOptions(includedActionTypes)

	if !p.IncludedApplicationIds.IsNull() && !p.IncludedApplicationIds.IsUnknown() {
		var typePlan []string
		diags.Append(p.IncludedApplicationIds.ElementsAs(ctx, &typePlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		objList := make([]management.SubscriptionFilterOptionsIncludedApplicationsInner, 0)
		for _, v := range typePlan {
			objList = append(objList, *management.NewSubscriptionFilterOptionsIncludedApplicationsInner(v))
		}

		data.SetIncludedApplications(objList)
	}

	if !p.IncludedPopulationIds.IsNull() && !p.IncludedPopulationIds.IsUnknown() {
		var typePlan []string
		diags.Append(p.IncludedPopulationIds.ElementsAs(ctx, &typePlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		objList := make([]management.SubscriptionFilterOptionsIncludedApplicationsInner, 0)
		for _, v := range typePlan {
			objList = append(objList, *management.NewSubscriptionFilterOptionsIncludedApplicationsInner(v))
		}

		data.SetIncludedPopulations(objList)
	}

	if !p.IncludedTags.IsNull() && !p.IncludedTags.IsUnknown() {
		var typePlan []string
		diags.Append(p.IncludedTags.ElementsAs(ctx, &typePlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		objList := make([]management.EnumSubscriptionFilterIncludedTags, 0)
		for _, v := range typePlan {
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

func (p *WebhookResourceModel) toState(apiObject *management.Subscription) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
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
	p.Format = framework.EnumOkToTF(apiObject.GetFormatOk())

	var d diag.Diagnostics
	p.FilterOptions, d = toStateWebhookFilterOptions(apiObject.GetFilterOptionsOk())
	diags.Append(d...)

	return diags
}

func toStateWebhookFilterOptions(v *management.SubscriptionFilterOptions, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: webhookFilterOptionsTFObjectTypes}

	if !ok || v == nil {
		return types.ListNull(types.ObjectType{AttrTypes: webhookFilterOptionsTFObjectTypes}), diags
	}

	var applicationIDSet, populationIDSet basetypes.SetValue

	applicationIDList := make([]string, 0)
	if items, ok := v.GetIncludedApplicationsOk(); ok {
		for _, application := range items {
			applicationIDList = append(applicationIDList, application.GetId())
		}

		applicationIDSet = framework.StringSetToTF(applicationIDList)
	} else {
		applicationIDSet = types.SetNull(types.StringType)
	}

	populationIDList := make([]string, 0)
	if items, ok := v.GetIncludedPopulationsOk(); ok {
		for _, population := range items {
			populationIDList = append(populationIDList, population.GetId())
		}

		populationIDSet = framework.StringSetToTF(populationIDList)
	} else {
		populationIDSet = types.SetNull(types.StringType)
	}

	objMap := map[string]attr.Value{
		"included_action_types":    framework.StringSetOkToTF(v.GetIncludedActionTypesOk()),
		"included_application_ids": applicationIDSet,
		"included_population_ids":  populationIDSet,
		"included_tags":            framework.EnumSetOkToTF(v.GetIncludedTagsOk()),
		"ip_address_exposed":       framework.BoolOkToTF(v.GetIpAddressExposedOk()),
		"useragent_exposed":        framework.BoolOkToTF(v.GetUserAgentExposedOk()),
	}

	flattenedObj, d := types.ObjectValue(webhookFilterOptionsTFObjectTypes, objMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags

}
