package credentials

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	customobjectvalidator "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type CredentialIssuanceRuleResource struct {
	client *credentials.APIClient
	region model.RegionMapping
}

type CredentialIssuanceRuleResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	EnvironmentId              types.String `tfsdk:"environment_id"`
	CredentialTypeId           types.String `tfsdk:"credential_type_id"`
	DigitalWalletApplicationId types.String `tfsdk:"digital_wallet_application_id"`
	Automation                 types.Object `tfsdk:"automation"`
	Filter                     types.Object `tfsdk:"filter"`
	Notification               types.Object `tfsdk:"notification"`
	Status                     types.String `tfsdk:"status"`
}

type FilterModel struct {
	GroupIds      types.Set    `tfsdk:"group_ids"`
	PopulationIds types.Set    `tfsdk:"population_ids"`
	Scim          types.String `tfsdk:"scim"`
}

type AutomationModel struct {
	Issue  types.String `tfsdk:"issue"`
	Revoke types.String `tfsdk:"revoke"`
	Update types.String `tfsdk:"update"`
}

type NotificationModel struct {
	Methods  types.Set    `tfsdk:"methods"`
	Template types.Object `tfsdk:"template"`
}

type NotificationTemplateModel struct {
	Locale  types.String `tfsdk:"locale"`
	Variant types.String `tfsdk:"variant"`
}

var (
	filterServiceTFObjectTypes = map[string]attr.Type{
		"group_ids":      types.SetType{ElemType: types.StringType},
		"population_ids": types.SetType{ElemType: types.StringType},
		"scim":           types.StringType,
	}

	automationServiceTFObjectTypes = map[string]attr.Type{
		"issue":  types.StringType,
		"revoke": types.StringType,
		"update": types.StringType,
	}

	notificationServiceTFObjectTypes = map[string]attr.Type{
		"methods":  types.SetType{ElemType: types.StringType},
		"template": types.ObjectType{AttrTypes: notificationTemplateServiceTFObjectTypes},
	}

	notificationTemplateServiceTFObjectTypes = map[string]attr.Type{
		"locale":  types.StringType,
		"variant": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &CredentialIssuanceRuleResource{}
	_ resource.ResourceWithConfigure   = &CredentialIssuanceRuleResource{}
	_ resource.ResourceWithImportState = &CredentialIssuanceRuleResource{}
)

// New Object
func NewCredentialIssuanceRuleResource() resource.Resource {
	return &CredentialIssuanceRuleResource{}
}

// Metadata
func (r *CredentialIssuanceRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_issuance_rule"
}

func (r *CredentialIssuanceRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	statusdDescriptionFmt := "Status of the credential issuance rule. Can be `ACTIVE` or `DISABLED`."
	statusDescription := framework.SchemaDescription{
		MarkdownDescription: statusdDescriptionFmt,
		Description:         strings.ReplaceAll(statusdDescriptionFmt, "`", "\""),
	}

	filterDescriptionFmt := "Contains one and only one filter (`group_ids`, `population_ids`, or `scim`) that selects the users to which the credential issuance rule applies."
	filterDescription := framework.SchemaDescription{
		MarkdownDescription: filterDescriptionFmt,
		Description:         strings.ReplaceAll(filterDescriptionFmt, "`", "\""),
	}

	automationOptionPhraseFmt := "Can be `PERIODIC` or `ON_DEMAND`." // I'm following the documentation here.
	automationIssueDescriptionFmt := fmt.Sprintf("The method the service uses to issue credentials with the credential issuance rule. %s", automationOptionPhraseFmt)
	automationIssueDescription := framework.SchemaDescription{
		MarkdownDescription: automationIssueDescriptionFmt,
		Description:         strings.ReplaceAll(automationIssueDescriptionFmt, "`", "\""),
	}

	automationRevokeDescriptionFmt := fmt.Sprintf("The method the service uses to revoke credentials with the credential issuance rule. %s", automationOptionPhraseFmt)
	automationRevokeDescription := framework.SchemaDescription{
		MarkdownDescription: automationRevokeDescriptionFmt,
		Description:         strings.ReplaceAll(automationRevokeDescriptionFmt, "`", "\""),
	}

	automationUpdateDescriptionFmt := fmt.Sprintf("The method the service uses to update credentials with the credential issuance rule. %s", automationOptionPhraseFmt)
	automationUpdateDescription := framework.SchemaDescription{
		MarkdownDescription: automationUpdateDescriptionFmt,
		Description:         strings.ReplaceAll(automationUpdateDescriptionFmt, "`", "\""),
	}

	notificationMethodsDescriptionFmt := "Array of methods for notifying the user; can be `EMAIL`, `SMS`, or both."
	notificationMethodsDescription := framework.SchemaDescription{
		MarkdownDescription: notificationMethodsDescriptionFmt,
		Description:         strings.ReplaceAll(notificationMethodsDescriptionFmt, "`", "\""),
	}

	notificationTemplateLocaleDescriptionFmt := "The ISO 2-character language code used for the notification; for example, `en`."
	notificationTemplateLocaleDescription := framework.SchemaDescription{
		MarkdownDescription: notificationTemplateLocaleDescriptionFmt,
		Description:         strings.ReplaceAll(notificationTemplateLocaleDescriptionFmt, "`", "\""),
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Credentials credential issuance rules.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "PingOne environment identifier (UUID) in which the credential issuance rule exists."},
			),

			"credential_type_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "Identifier (UUID) of the credential type with which this credential issuance rule is associated."},
			),

			"digital_wallet_application_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "Identifier (UUID) of the customer's Digital Wallet App that will interact with the user's Digital Wallet."},
			),

			"status": schema.StringAttribute{
				Description:         statusDescription.Description,
				MarkdownDescription: statusDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(credentials.ENUMCREDENTIALISSUANCERULESTATUS_ACTIVE),
						string(credentials.ENUMCREDENTIALISSUANCERULESTATUS_DISABLED)),
				},
			},

			"filter": schema.SingleNestedAttribute{
				Description:         filterDescription.Description,
				MarkdownDescription: filterDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.Object{
					customobjectvalidator.IsRequiredIfMatchesPathValue(
						basetypes.NewStringValue(string(credentials.ENUMCREDENTIALISSUANCERULESTATUS_ACTIVE)),
						path.MatchRelative().AtParent().AtName("status"),
					),
				},
				Attributes: map[string]schema.Attribute{
					"group_ids": schema.SetAttribute{
						ElementType: types.StringType,
						Description: "Array of one or more identifiers (UUIDs) of groups, any of which a user must belong for the credential issuance rule to apply.",
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("population_ids")),
							setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("scim")),
							setvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("population_ids"),
								path.MatchRelative().AtParent().AtName("scim"),
							),
							setvalidator.SizeAtLeast(attrMinLength),
						},
					},
					"population_ids": schema.SetAttribute{
						ElementType: types.StringType,
						Description: "Array of one or more identifiers (UUIDs) of populations, any of which a user must belong for the credential issuance rule to apply. ",
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("group_ids")),
							setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("scim")),
							setvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("group_ids"),
								path.MatchRelative().AtParent().AtName("scim"),
							),
							setvalidator.SizeAtLeast(attrMinLength),
						},
					},

					"scim": schema.StringAttribute{
						Description: "A SCIM query that selects users to which the credential issuance rule applies.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("group_ids")),
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("population_ids")),
							stringvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("group_ids"),
								path.MatchRelative().AtParent().AtName("population_ids"),
							),
						},
					},
				},
			},

			"automation": schema.SingleNestedAttribute{
				Description: "Contains a list of actions, as key names, and the update method for each action.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"issue": schema.StringAttribute{
						Description:         automationIssueDescription.Description,
						MarkdownDescription: automationIssueDescription.MarkdownDescription,
						Required:            true,
					},
					"revoke": schema.StringAttribute{
						Description:         automationRevokeDescription.Description,
						MarkdownDescription: automationRevokeDescription.MarkdownDescription,
						Required:            true,
					},
					"update": schema.StringAttribute{
						Description:         automationUpdateDescription.Description,
						MarkdownDescription: automationUpdateDescription.MarkdownDescription,
						Required:            true,
					},
				},
			},

			"notification": schema.SingleNestedAttribute{
				Description: "Contains notification information. When this property is supplied, the information within is used to create a custom notification.",
				Optional:    true,
				Validators: []validator.Object{
					objectvalidator.Any(
						objectvalidator.AlsoRequires(path.MatchRelative().AtName("methods")),
						objectvalidator.AlsoRequires(path.MatchRelative().AtName("template")),
					),
				},
				Attributes: map[string]schema.Attribute{
					"methods": schema.SetAttribute{
						ElementType:         types.StringType,
						Description:         notificationMethodsDescription.Description,
						MarkdownDescription: notificationMethodsDescription.MarkdownDescription,
						Optional:            true,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf(
									string(credentials.ENUMCREDENTIALISSUANCERULENOTIFICATIONMETHOD_EMAIL),
									string(credentials.ENUMCREDENTIALISSUANCERULENOTIFICATIONMETHOD_SMS),
								),
							),
							setvalidator.SizeAtLeast(attrMinLength),
						},
					},
					"template": schema.SingleNestedAttribute{
						Description: "Contains template parameters.",
						Optional:    true,
						Validators: []validator.Object{
							objectvalidator.Any(
								objectvalidator.AlsoRequires(path.MatchRelative().AtName("locale")),
								objectvalidator.AlsoRequires(path.MatchRelative().AtName("variant")),
							),
							objectvalidator.AlsoRequires(
								path.MatchRelative().AtParent().AtName("methods"),
							),
						},
						Attributes: map[string]schema.Attribute{
							"locale": schema.StringAttribute{
								Description:         notificationTemplateLocaleDescription.Description,
								MarkdownDescription: notificationTemplateLocaleDescription.MarkdownDescription,
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.OneOf(verify.FullIsoList()...),
								},
							},
							"variant": schema.StringAttribute{
								Description: "The unique user-defined name for the content variant that contains the message text used for the notification.",
								Optional:    true,
							},
						},
					},
				},
			},
		},
	}
}

func (r *CredentialIssuanceRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CredentialIssuanceRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state CredentialIssuanceRuleResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	CredentialIssuanceRule, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.CredentialIssuanceRulesApi.CreateCredentialIssuanceRule(ctx, plan.EnvironmentId.ValueString(), plan.CredentialTypeId.ValueString()).CredentialIssuanceRule(*CredentialIssuanceRule).Execute()
		},
		"CreateCredentialIssuanceRule",
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
	resp.Diagnostics.Append(state.toState(response.(*credentials.CredentialIssuanceRule))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CredentialIssuanceRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CredentialIssuanceRuleResourceModel

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
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.CredentialIssuanceRulesApi.ReadOneCredentialIssuanceRule(ctx, data.EnvironmentId.ValueString(), data.CredentialTypeId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneCredentialIssuanceRule",
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
	resp.Diagnostics.Append(data.toState(response.(*credentials.CredentialIssuanceRule))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CredentialIssuanceRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state CredentialIssuanceRuleResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	CredentialIssuanceRule, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.CredentialIssuanceRulesApi.UpdateCredentialIssuanceRule(ctx, plan.EnvironmentId.ValueString(), plan.CredentialTypeId.ValueString(), plan.Id.ValueString()).CredentialIssuanceRule(*CredentialIssuanceRule).Execute()
		},
		"UpdateCredentialIssuanceRule",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*credentials.CredentialIssuanceRule))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CredentialIssuanceRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CredentialIssuanceRuleResourceModel

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
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	_, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := r.client.CredentialIssuanceRulesApi.DeleteCredentialIssuanceRule(ctx, data.EnvironmentId.ValueString(), data.CredentialTypeId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteCredentialIssuanceRule",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CredentialIssuanceRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/credential_type_id/credential_issuance_rule_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("credential_type_id"), attributes[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *CredentialIssuanceRuleResourceModel) expand(ctx context.Context) (*credentials.CredentialIssuanceRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	// expand automation rules
	credentialIssuanceRuleAutomation := credentials.NewCredentialIssuanceRuleAutomationWithDefaults()
	if !p.Automation.IsNull() && !p.Automation.IsUnknown() {
		var automationRules AutomationModel
		d := p.Automation.As(ctx, &automationRules, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		credentialIssuanceRuleAutomation, d = automationRules.expandAutomationModel()
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// expand filter
	credentialIssuanceRuleFilter := credentials.NewCredentialIssuanceRuleFilterWithDefaults()
	if !p.Filter.IsNull() && !p.Filter.IsUnknown() {
		var filterRules FilterModel
		d := p.Filter.As(ctx, &filterRules, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		credentialIssuanceRuleFilter, d = filterRules.expandFilterModel(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// expand notifications
	credentialIssuanceRuleNotification := credentials.NewCredentialIssuanceRuleNotificationWithDefaults()
	if !p.Notification.IsNull() && !p.Notification.IsUnknown() {
		var notificationRules NotificationModel
		d := p.Notification.As(ctx, &notificationRules, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		credentialIssuanceRuleNotification, d = notificationRules.expandNotificationModel(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// buuild issuance rule object with required attributes
	data := credentials.NewCredentialIssuanceRule(*credentialIssuanceRuleAutomation, credentials.EnumCredentialIssuanceRuleStatus(p.Status.ValueString()))

	// set the filter details
	if credentialIssuanceRuleFilter.HasGroupIds() || credentialIssuanceRuleFilter.HasPopulationIds() || credentialIssuanceRuleFilter.HasScim() {
		data.SetFilter(*credentialIssuanceRuleFilter)
	}

	// set the notification details
	if credentialIssuanceRuleNotification.HasMethods() || credentialIssuanceRuleNotification.HasTemplate() {
		data.SetNotification(*credentialIssuanceRuleNotification)
	}

	// set the digital wallet application
	application := credentials.NewCredentialIssuanceRuleDigitalWalletApplication(p.DigitalWalletApplicationId.ValueString())
	data.SetDigitalWalletApplication(*application)

	return data, diags
}

func (p *AutomationModel) expandAutomationModel() (*credentials.CredentialIssuanceRuleAutomation, diag.Diagnostics) {
	var diags diag.Diagnostics

	automation := credentials.NewCredentialIssuanceRuleAutomationWithDefaults()

	if !p.Issue.IsNull() && !p.Issue.IsUnknown() {
		automation.SetIssue(credentials.EnumCredentialIssuanceRuleAutomationMethod(p.Issue.ValueString()))
	}

	if !p.Revoke.IsNull() && !p.Revoke.IsUnknown() {
		automation.SetRevoke(credentials.EnumCredentialIssuanceRuleAutomationMethod(p.Revoke.ValueString()))
	}

	if !p.Update.IsNull() && !p.Update.IsUnknown() {
		automation.SetUpdate(credentials.EnumCredentialIssuanceRuleAutomationMethod(p.Update.ValueString()))
	}

	if automation == nil {
		diags.AddWarning(
			"Unexpected Value",
			"Automation object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}
	return automation, diags

}

func (p *FilterModel) expandFilterModel(ctx context.Context) (*credentials.CredentialIssuanceRuleFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	filter := credentials.NewCredentialIssuanceRuleFilterWithDefaults()

	if !p.PopulationIds.IsNull() && !p.PopulationIds.IsUnknown() {
		diags.Append(p.PopulationIds.ElementsAs(ctx, &filter.PopulationIds, false)...)
		if diags.HasError() {
			return nil, diags
		}
		filter.SetPopulationIds(filter.PopulationIds)
	}

	if !p.GroupIds.IsNull() && !p.GroupIds.IsUnknown() {
		diags.Append(p.GroupIds.ElementsAs(ctx, &filter.GroupIds, false)...)
		if diags.HasError() {
			return nil, diags
		}
		filter.SetGroupIds(filter.GroupIds)
	}

	if !p.Scim.IsNull() && !p.Scim.IsUnknown() {
		filter.SetScim(p.Scim.ValueString())
	}

	return filter, diags

}

func (p *NotificationModel) expandNotificationModel(ctx context.Context) (*credentials.CredentialIssuanceRuleNotification, diag.Diagnostics) {
	var diags diag.Diagnostics

	notification := credentials.NewCredentialIssuanceRuleNotificationWithDefaults()

	// notification methods
	if !p.Methods.IsNull() && !p.Methods.IsUnknown() {
		var slice []string
		diags.Append(p.Methods.ElementsAs(ctx, &slice, false)...)

		enumSlice := make([]credentials.EnumCredentialIssuanceRuleNotificationMethod, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := credentials.NewEnumCredentialIssuanceRuleNotificationMethodFromValue(slice[i])
			if err != nil {
				return nil, diags
			}
			enumSlice[i] = *enumVal
			notification.Methods = append(notification.Methods, *enumVal)
		}
	}

	// notification template
	if !p.Template.IsNull() && !p.Template.IsUnknown() {
		var notificationTemplate NotificationTemplateModel
		d := p.Template.As(ctx, &notificationTemplate, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		template := credentials.NewCredentialIssuanceRuleNotificationTemplate()
		if !notificationTemplate.Locale.IsNull() && !notificationTemplate.Locale.IsUnknown() {
			template.SetLocale(notificationTemplate.Locale.ValueString())
		}

		if !notificationTemplate.Variant.IsNull() && !notificationTemplate.Variant.IsUnknown() {
			template.SetVariant(notificationTemplate.Variant.ValueString())
		}

		notification.SetTemplate(*template)
	}

	return notification, diags

}

func (p *CredentialIssuanceRuleResourceModel) toState(apiObject *credentials.CredentialIssuanceRule) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	// core issuance rule attributes
	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.DigitalWalletApplicationId = framework.StringToTF(apiObject.GetDigitalWalletApplication().Id)
	p.CredentialTypeId = framework.StringToTF(apiObject.CredentialType.GetId())
	p.Status = enumCredentialIssuanceStatusOkToTF(apiObject.GetStatusOk())

	// automation object
	if v, ok := apiObject.GetAutomationOk(); ok {
		automation, d := toStateAutomation(v, ok)
		diags.Append(d...)
		p.Automation = automation
	}

	// filter object
	if v, ok := apiObject.GetFilterOk(); ok {
		if v.HasGroupIds() || v.HasPopulationIds() || v.HasScim() { // check because values are optional
			filter, d := toStateFilter(v, ok)
			diags.Append(d...)
			p.Filter = filter
		}
	}

	// notification object
	if v, ok := apiObject.GetNotificationOk(); ok {
		if v.HasMethods() || v.HasTemplate() { // check because values are optional
			notification, d := toStateNotification(v, ok)
			diags.Append(d...)
			p.Notification = notification
		}
	}

	return diags
}

func toStateAutomation(automation *credentials.CredentialIssuanceRuleAutomation, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	automationMap := map[string]attr.Value{
		"issue":  enumCredentialIssuanceRuleAutomationOkToTF(automation.GetIssueOk()),
		"revoke": enumCredentialIssuanceRuleAutomationOkToTF(automation.GetRevokeOk()),
		"update": enumCredentialIssuanceRuleAutomationOkToTF(automation.GetUpdateOk()),
	}
	flattenedObj, d := types.ObjectValue(automationServiceTFObjectTypes, automationMap)
	diags.Append(d...)

	return flattenedObj, diags
}

func toStateFilter(filter *credentials.CredentialIssuanceRuleFilter, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if filter == nil {
		return types.ObjectNull(filterServiceTFObjectTypes), diags
	}

	filterMap := map[string]attr.Value{}
	if v, ok := filter.GetPopulationIdsOk(); ok {
		filterMap["population_ids"] = framework.StringSetOkToTF(v, ok)
	} else {
		filterMap["population_ids"] = types.SetNull(types.StringType)
	}

	if v, ok := filter.GetGroupIdsOk(); ok {
		filterMap["group_ids"] = framework.StringSetOkToTF(v, ok)
	} else {
		filterMap["group_ids"] = types.SetNull(types.StringType)
	}

	if v, ok := filter.GetScimOk(); ok {
		filterMap["scim"] = framework.StringOkToTF(v, ok)
	} else {
		filterMap["scim"] = types.StringNull()
	}

	flattenedObj, d := types.ObjectValue(filterServiceTFObjectTypes, filterMap)
	diags.Append(d...)

	return flattenedObj, diags
}

func toStateNotification(notification *credentials.CredentialIssuanceRuleNotification, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if notification == nil {
		return types.ObjectNull(notificationServiceTFObjectTypes), diags
	}

	notificationMap := map[string]attr.Value{}

	// notification.methods
	if v, ok := notification.GetMethodsOk(); ok {
		notificationMap["methods"] = enumCredentialIssuanceRuleNotificationMethodOkToTF(v, ok)
	} else {
		notificationMap["methods"] = types.SetNull(types.StringType)
	}

	// notification.template
	if notification.Template == nil {
		notificationMap["template"] = types.ObjectNull(notificationTemplateServiceTFObjectTypes)
	} else {
		notificationTemplate := map[string]attr.Value{
			"locale":  framework.StringOkToTF(notification.Template.GetLocaleOk()),
			"variant": framework.StringOkToTF(notification.Template.GetVariantOk()),
		}

		flattenedTemplate, d := types.ObjectValue(notificationTemplateServiceTFObjectTypes, notificationTemplate)
		diags.Append(d...)

		notificationMap["template"] = flattenedTemplate
	}

	flattenedObj, d := types.ObjectValue(notificationServiceTFObjectTypes, notificationMap)
	diags.Append(d...)

	return flattenedObj, diags
}

func enumCredentialIssuanceRuleNotificationMethodOkToTF(v []credentials.EnumCredentialIssuanceRuleNotificationMethod, ok bool) basetypes.SetValue {
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

func enumCredentialIssuanceRuleAutomationOkToTF(v *credentials.EnumCredentialIssuanceRuleAutomationMethod, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}

func enumCredentialIssuanceStatusOkToTF(v *credentials.EnumCredentialIssuanceRuleStatus, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
