package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type AlertChannelResource serviceClientType

type AlertChannelResourceModel struct {
	Addresses         types.Set                    `tfsdk:"addresses"`
	AlertName         types.String                 `tfsdk:"alert_name"`
	ChannelType       types.String                 `tfsdk:"channel_type"`
	EnvironmentId     pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ExcludeAlertTypes types.Set                    `tfsdk:"exclude_alert_types"`
	Id                pingonetypes.ResourceIDValue `tfsdk:"id"`
	IncludeAlertTypes types.Set                    `tfsdk:"include_alert_types"`
	IncludeSeverities types.Set                    `tfsdk:"include_severities"`
}

// Framework interfaces
var (
	_ resource.Resource                = &AlertChannelResource{}
	_ resource.ResourceWithConfigure   = &AlertChannelResource{}
	_ resource.ResourceWithImportState = &AlertChannelResource{}
)

// New Object
func NewAlertChannelResource() resource.Resource {
	return &AlertChannelResource{}
}

// Metadata
func (r *AlertChannelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_channel"
}

// Schema
func (r *AlertChannelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	channelTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the alert channel type. Currently, this must be `EMAIL`.",
	).AllowedValuesEnum(utils.EnumSliceToStringSlice(management.AllowedEnumAlertChannelTypeEnumValues))

	excludedAlertTypesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of strings that specifies the list of alert types that administrators will not be emailed alerts for. If left undefined, no alert types are excluded.",
	).AllowedValuesEnum(utils.EnumSliceToStringSlice(management.AllowedEnumAlertChannelAlertTypeEnumValues))

	includedAlertTypesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of strings that specifies the list of alert types that administrators will be emailed alerts for.",
	).AllowedValuesEnum(utils.EnumSliceToStringSlice(management.AllowedEnumAlertChannelAlertTypeEnumValues))

	includeSeveritiesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of strings that specifies the severity to filters alerts by.",
	).AllowedValuesEnum(utils.EnumSliceToStringSlice(management.AllowedEnumAlertChannelSeverityEnumValues))

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage alert channels in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage an alert channel for."),
			),

			"addresses": schema.SetAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of strings that specifies the administrator email addresses to send the alerts to.").Description,
				Required:    true,

				ElementType: types.StringType,

				Validators: []validator.Set{
					// setvalidator.ValueStringsAre(
					// 	stringvalidator.RegexMatches(verify.IsEmail, "must be a valid email address"),
					// ),
					setvalidator.SizeAtLeast(attrMinLength),
				},
			},

			"alert_name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the name to assign to the alert channel.").Description,
				Optional:    true,
			},

			"channel_type": schema.StringAttribute{
				Description:         channelTypeDescription.Description,
				MarkdownDescription: channelTypeDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumAlertChannelTypeEnumValues)...),
				},
			},

			"exclude_alert_types": schema.SetAttribute{
				Description:         excludedAlertTypesDescription.Description,
				MarkdownDescription: excludedAlertTypesDescription.MarkdownDescription,
				Optional:            true,

				ElementType: types.StringType,

				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumAlertChannelAlertTypeEnumValues)...),
					),
					setvalidator.SizeAtLeast(attrMinLength),
				},
			},

			"include_alert_types": schema.SetAttribute{
				Description:         includedAlertTypesDescription.Description,
				MarkdownDescription: includedAlertTypesDescription.MarkdownDescription,
				Required:            true,

				ElementType: types.StringType,

				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumAlertChannelAlertTypeEnumValues)...),
					),
					setvalidator.SizeAtLeast(attrMinLength),
				},
			},

			"include_severities": schema.SetAttribute{
				Description:         includeSeveritiesDescription.Description,
				MarkdownDescription: includeSeveritiesDescription.MarkdownDescription,
				Required:            true,

				ElementType: types.StringType,

				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumAlertChannelSeverityEnumValues)...),
					),
				},
			},
		},
	}
}

func (r *AlertChannelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AlertChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state AlertChannelResourceModel

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
	alertChannel, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.AlertChannel
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AlertingApi.CreateAlertChannel(ctx, plan.EnvironmentId.ValueString()).AlertChannel(*alertChannel).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateAlertChannel",
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

func (r *AlertChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AlertChannelResourceModel

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
	var listResponse *management.EntityArray
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AlertingApi.ReadAllAlertChannels(ctx, data.EnvironmentId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadAllAlertChannels",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&listResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if listResponse == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Find the resource in the list
	var response management.AlertChannel
	found := false
	if embedded, ok := listResponse.GetEmbeddedOk(); ok {
		if alertChannels, ok := embedded.GetAlertChannelsOk(); ok {
			for _, alertChannel := range alertChannels {
				if alertChannel.GetId() == data.Id.ValueString() {
					response = alertChannel
					found = true
					break
				}
			}
		}
	}

	// Remove from state if resource is not found
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AlertChannelResourceModel

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
	alertChannel, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.AlertChannel
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AlertingApi.UpdateAlertChannel(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).AlertChannel(*alertChannel).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAlertChannel",
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

func (r *AlertChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AlertChannelResourceModel

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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.AlertingApi.DeleteAlertChannel(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteAlertChannel",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AlertChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "alert_channel_id",
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

func (p *AlertChannelResourceModel) expand(ctx context.Context) (*management.AlertChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	var addressesPlan []types.String
	diags.Append(p.Addresses.ElementsAs(ctx, &addressesPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	addresses, d := framework.TFTypeStringSliceToStringSlice(addressesPlan, path.Root("addresses"))
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := management.NewAlertChannel(
		management.EnumAlertChannelType(p.ChannelType.ValueString()),
		addresses,
	)

	if !p.AlertName.IsNull() && !p.AlertName.IsUnknown() {
		data.SetAlertName(p.AlertName.ValueString())
	}

	if !p.ExcludeAlertTypes.IsNull() && !p.ExcludeAlertTypes.IsUnknown() {

		var excludeAlertTypesPlan []types.String
		diags.Append(p.ExcludeAlertTypes.ElementsAs(ctx, &excludeAlertTypesPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		excludeAlertTypesPlanStr, d := framework.TFTypeStringSliceToStringSlice(excludeAlertTypesPlan, path.Root("exclude_alert_types"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		excludeAlertTypes := make([]management.EnumAlertChannelAlertType, len(excludeAlertTypesPlan))
		for i, v := range excludeAlertTypesPlanStr {
			excludeAlertTypes[i] = management.EnumAlertChannelAlertType(v)
		}

		data.SetExcludeAlertTypes(excludeAlertTypes)
	}

	if !p.IncludeAlertTypes.IsNull() && !p.IncludeAlertTypes.IsUnknown() {

		var includeAlertTypesPlan []types.String
		diags.Append(p.IncludeAlertTypes.ElementsAs(ctx, &includeAlertTypesPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		includeAlertTypesPlanStr, d := framework.TFTypeStringSliceToStringSlice(includeAlertTypesPlan, path.Root("include_alert_types"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		includeAlertTypes := make([]management.EnumAlertChannelAlertType, len(includeAlertTypesPlan))
		for i, v := range includeAlertTypesPlanStr {
			includeAlertTypes[i] = management.EnumAlertChannelAlertType(v)
		}

		data.SetIncludeAlertTypes(includeAlertTypes)
	}

	if !p.IncludeSeverities.IsNull() && !p.IncludeSeverities.IsUnknown() {

		var includeSeveritiesPlan []types.String
		diags.Append(p.IncludeSeverities.ElementsAs(ctx, &includeSeveritiesPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		includeSeveritiesStr, d := framework.TFTypeStringSliceToStringSlice(includeSeveritiesPlan, path.Root("include_severities"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		includeSeverities := make([]management.EnumAlertChannelSeverity, len(includeSeveritiesPlan))
		for i, v := range includeSeveritiesStr {
			includeSeverities[i] = management.EnumAlertChannelSeverity(v)
		}

		data.SetIncludeSeverities(includeSeverities)
	}

	return data, diags
}

func (p *AlertChannelResourceModel) toState(apiObject *management.AlertChannel) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.Addresses = framework.StringSetOkToTF(apiObject.GetAddressesOk())
	p.AlertName = framework.StringOkToTF(apiObject.GetAlertNameOk())
	p.ChannelType = framework.EnumOkToTF(apiObject.GetChannelTypeOk())
	p.ExcludeAlertTypes = framework.EnumSetOkToTF(apiObject.GetExcludeAlertTypesOk())
	p.IncludeAlertTypes = framework.EnumSetOkToTF(apiObject.GetIncludeAlertTypesOk())
	p.IncludeSeverities = framework.EnumSetOkToTF(apiObject.GetIncludeSeveritiesOk())

	return diags
}
