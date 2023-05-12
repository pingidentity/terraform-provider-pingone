package credentials

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type CredentialIssuanceRuleDataSource struct {
	client *credentials.APIClient
	region model.RegionMapping
}

type CredentialIssuanceRuleDataSourceModel struct {
	Id                         types.String `tfsdk:"id"`
	EnvironmentId              types.String `tfsdk:"environment_id"`
	CredentialTypeId           types.String `tfsdk:"credential_type_id"`
	CredentialIssuanceRuleId   types.String `tfsdk:"credential_issuance_rule_id"`
	DigitalWalletApplicationId types.String `tfsdk:"digital_wallet_application_id"`
	Automation                 types.Object `tfsdk:"automation"`
	Filter                     types.Object `tfsdk:"filter"`
	Notification               types.Object `tfsdk:"notification"`
	Status                     types.String `tfsdk:"status"`
}

type FilterDataSourceModel struct {
	GroupIds      types.Set    `tfsdk:"group_ids"`
	PopulationIds types.Set    `tfsdk:"population_ids"`
	Scim          types.String `tfsdk:"scim"`
}

type AutomationDataSourceModel struct {
	Issue  types.String `tfsdk:"issue"`
	Revoke types.String `tfsdk:"revoke"`
	Update types.String `tfsdk:"update"`
}

type NotificationDataSourceModel struct {
	Methods  types.Set    `tfsdk:"methods"`
	Template types.Object `tfsdk:"template"`
}

type NotificationTemplateDataSourceModel struct {
	Locale  types.String `tfsdk:"locale"`
	Variant types.String `tfsdk:"variant"`
}

var (
	filterDataSourceServiceTFObjectTypes = map[string]attr.Type{ // todo: make naming consistent with Tfobjecttype
		"group_ids":      types.SetType{ElemType: types.StringType},
		"population_ids": types.SetType{ElemType: types.StringType},
		"scim":           types.StringType,
	}

	automationDataSourceServiceTFObjectTypes = map[string]attr.Type{ // todo: make naming consistent with tfobjecttypes
		"issue":  types.StringType,
		"revoke": types.StringType,
		"update": types.StringType,
	}

	notificationDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"methods":  types.SetType{ElemType: types.StringType},
		"template": types.ObjectType{AttrTypes: notificationTemplateDataSourceServiceTFObjectTypes},
	}

	notificationTemplateDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"locale":  types.StringType,
		"variant": types.StringType,
	}
)

// Framework interfaces
var (
	_ datasource.DataSource = &CredentialIssuanceRuleDataSource{}
)

// New Object
func NewCredentialIssuanceRuleDataSource() datasource.DataSource {
	return &CredentialIssuanceRuleDataSource{}
}

// Metadata
func (r *CredentialIssuanceRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_issuance_rule"
}

func (r *CredentialIssuanceRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Credentials credential issuance rules.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment to create the credential type in."},
			),

			"credential_type_id": framework.Attr_LinkIDWithValidators(framework.SchemaDescription{
				Description: "The ID of the credential type with which this credential issuance rule is associated.",
			},
				[]validator.String{
					verify.P1ResourceIDValidator(),
				},
			),

			"credential_issuance_rule_id": framework.Attr_LinkIDWithValidators(framework.SchemaDescription{
				Description: "The ID of the credential issuance rule assigned to the credential type.",
			},
				[]validator.String{
					verify.P1ResourceIDValidator(),
				},
			),

			"digital_wallet_application_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the digital wallet application that will interact with the user's Digital Wallet",
				Computed:            true,
			},

			"status": schema.StringAttribute{
				MarkdownDescription: "ACTIVE or DISABLED status of the credential issuance rule.",
				Computed:            true,
			},

			"filter": schema.SingleNestedAttribute{
				MarkdownDescription: "",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"group_ids": schema.SetAttribute{
						ElementType:         types.StringType,
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
					},
					"population_ids": schema.SetAttribute{
						ElementType:         types.StringType,
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
					},
					"scim": schema.StringAttribute{
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
					},
				},
			},

			"automation": schema.SingleNestedAttribute{
				MarkdownDescription: "",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"issue": schema.StringAttribute{
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
					},
					"revoke": schema.StringAttribute{
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
					},
					"update": schema.StringAttribute{
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
					},
				},
			},

			"notification": schema.SingleNestedAttribute{
				MarkdownDescription: "",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"methods": schema.SetAttribute{
						ElementType:         types.StringType,
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
					},
					"template": schema.SingleNestedAttribute{
						MarkdownDescription: "",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"locale": schema.StringAttribute{
								Description:         "",
								MarkdownDescription: "",
								Computed:            true,
							},
							"variant": schema.StringAttribute{
								Description:         "",
								MarkdownDescription: "",
								Computed:            true,
							},
						},
					},
				},
			},
		},
	}
}

func (r *CredentialIssuanceRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *CredentialIssuanceRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CredentialIssuanceRuleDataSourceModel

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

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.CredentialIssuanceRulesApi.ReadOneCredentialIssuanceRule(ctx, data.EnvironmentId.ValueString(), data.CredentialTypeId.ValueString(), data.CredentialIssuanceRuleId.ValueString()).Execute()
		},
		"ReadOneCredentialIssuanceRule",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Error if not found
	if response == nil {
		resp.Diagnostics.AddError(
			"Cannot find credential issuance rule",
			fmt.Sprintf("The credential issuance rule %s for environment %s cannot be found.", data.CredentialTypeId.String(), data.EnvironmentId.String()),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response.(*credentials.CredentialIssuanceRule))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *CredentialIssuanceRuleDataSourceModel) toState(apiObject *credentials.CredentialIssuanceRule) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	// core issuance rule attributes
	p.Id = framework.StringToTF(apiObject.GetId())
	p.EnvironmentId = framework.StringToTF(apiObject.GetEnvironment().Id)
	p.DigitalWalletApplicationId = framework.StringToTF(apiObject.GetDigitalWalletApplication().Id)
	p.CredentialTypeId = framework.StringToTF(apiObject.CredentialType.Id)
	p.CredentialIssuanceRuleId = framework.StringToTF(apiObject.GetId())
	p.Status = enumCredentialIssuanceStatusDataSourceOkToTF(apiObject.GetStatusOk())

	// automation object
	automation, d := toStateAutomationDataSource(apiObject.GetAutomationOk())
	diags.Append(d...)
	p.Automation = automation

	// filter object
	filter, d := toStateFilterDataSource(apiObject.GetFilterOk())
	diags.Append(d...)
	p.Filter = filter

	// notification object
	notificationMethodState := enumCredentialIssuanceRuleNotificationMethodDataSourceOkToTF(apiObject.Notification.GetMethodsOk())

	if notificationMethodState.IsNull() {
		// todo: not sure how to handle this at the moment...

	} else {
		notification, d := toStateNotificationDataSource(apiObject.GetNotificationOk())
		diags.Append(d...)

		p.Notification = notification
	}

	return diags
}

func toStateAutomationDataSource(automation *credentials.CredentialIssuanceRuleAutomation, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	automationMap := map[string]attr.Value{
		"issue":  enumCredentialIssuanceRuleAutomationDataSourceOkToTF(automation.GetIssueOk()),
		"revoke": enumCredentialIssuanceRuleAutomationDataSourceOkToTF(automation.GetRevokeOk()),
		"update": enumCredentialIssuanceRuleAutomationDataSourceOkToTF(automation.GetUpdateOk()),
	}
	flattenedObj, d := types.ObjectValue(automationDataSourceServiceTFObjectTypes, automationMap)
	diags.Append(d...)

	return flattenedObj, diags
}

func toStateFilterDataSource(filter *credentials.CredentialIssuanceRuleFilter, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	filterMap := map[string]attr.Value{
		"population_ids": framework.StringSetOkToTF(filter.GetPopulationIdsOk()),
		"group_ids":      framework.StringSetOkToTF(filter.GetGroupIdsOk()),
		"scim":           framework.StringOkToTF(filter.GetScimOk()),
	}
	flattenedObj, d := types.ObjectValue(filterDataSourceServiceTFObjectTypes, filterMap)
	diags.Append(d...)

	return flattenedObj, diags
}

func toStateNotificationDataSource(notification *credentials.CredentialIssuanceRuleNotification, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	notificationTemplate := map[string]attr.Value{
		"locale":  framework.StringOkToTF(notification.Template.GetLocaleOk()),
		"variant": framework.StringOkToTF(notification.Template.GetVariantOk()),
	}

	flattenedTemplate, d := types.ObjectValue(notificationTemplateDataSourceServiceTFObjectTypes, notificationTemplate)
	diags.Append(d...)

	notificationMap := map[string]attr.Value{
		"methods":  enumCredentialIssuanceRuleNotificationMethodDataSourceOkToTF(notification.GetMethodsOk()),
		"template": flattenedTemplate,
	}

	flattenedObj, d := types.ObjectValue(notificationDataSourceServiceTFObjectTypes, notificationMap)
	diags.Append(d...)

	return flattenedObj, diags
}

func enumCredentialIssuanceRuleNotificationMethodDataSourceOkToTF(v []credentials.EnumCredentialIssuanceRuleNotificationMethod, ok bool) basetypes.SetValue {
	if !ok || v == nil {
		return types.SetNull(types.StringType)
	} else {
		list := make([]attr.Value, 0)
		for _, item := range v {
			method := types.StringValue(string(item))
			list = append(list, method)
		}

		return types.SetValueMust(types.StringType, list)
	}
}

func enumCredentialIssuanceRuleAutomationDataSourceOkToTF(v *credentials.EnumCredentialIssuanceRuleAutomationMethod, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}

func enumCredentialIssuanceStatusDataSourceOkToTF(v *credentials.EnumCredentialIssuanceRuleStatus, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
