package credentials

import (
	"context"
	"fmt"
	"net/http"

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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	customobjectvalidator "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type CredentialIssuanceRuleResource serviceClientType

type CredentialIssuanceRuleResourceModel struct {
	Id                         pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId              pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	CredentialTypeId           pingonetypes.ResourceIDValue `tfsdk:"credential_type_id"`
	DigitalWalletApplicationId pingonetypes.ResourceIDValue `tfsdk:"digital_wallet_application_id"`
	Automation                 types.Object                 `tfsdk:"automation"`
	Filter                     types.Object                 `tfsdk:"filter"`
	Notification               types.Object                 `tfsdk:"notification"`
	Status                     types.String                 `tfsdk:"status"`
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

	statusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Status of the credential issuance rule. Can be `ACTIVE` or `DISABLED`.",
	)

	filterDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Contains one and only one filter (`group_ids`, `population_ids`, or `scim`) that selects the users to which the credential issuance rule applies. A filter must be defined if the issuance rule `status` is `ACTIVE`.",
	)

	automationOptionPhraseFmt := "Can be `PERIODIC` or `ON_DEMAND`." // I'm following the documentation here.
	automationIssueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The method the service uses to issue credentials with the credential issuance rule. %s", automationOptionPhraseFmt),
	)

	automationRevokeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The method the service uses to revoke credentials with the credential issuance rule. %s", automationOptionPhraseFmt),
	)

	automationUpdateDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The method the service uses to update credentials with the credential issuance rule. %s", automationOptionPhraseFmt),
	)

	notificationMethodsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Array of methods for notifying the user; can be `EMAIL`, `SMS`, or both.",
	)

	notificationTemplateLocaleDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ISO 2-character language code used for the notification; for example, `en`.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create, read, and update rules for issuing, updating, and revoking credentials by credential type.\n\n" +
			"An issuance rule is defined for a specific `credential_type` and `digital_wallet_application`, and the `filter` determines the targeted list of users allowed to receive the specific credential type.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the credential issuance rule exists."),
			),

			"credential_type_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("Identifier (UUID) of the credential type with which this credential issuance rule is associated."),
			),

			"digital_wallet_application_id": schema.StringAttribute{
				Description: "Identifier (UUID) of the customer's Digital Wallet App that will interact with the user's Digital Wallet. If present, digital wallet pairing automatically starts when a user matches the credential issuance rule.",
				Optional:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},

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
				Validators:  []validator.Object{},
				Attributes: map[string]schema.Attribute{
					"methods": schema.SetAttribute{
						ElementType:         types.StringType,
						Description:         notificationMethodsDescription.Description,
						MarkdownDescription: notificationMethodsDescription.MarkdownDescription,
						Required:            true,
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
							objectvalidator.AlsoRequires(
								path.MatchRelative().AtParent().AtName("methods"),
							),
						},
						Attributes: map[string]schema.Attribute{
							"locale": schema.StringAttribute{
								Description:         notificationTemplateLocaleDescription.Description,
								MarkdownDescription: notificationTemplateLocaleDescription.MarkdownDescription,
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf(verify.FullIsoList()...),
								},
							},
							"variant": schema.StringAttribute{
								Description: "The unique user-defined name for the content variant that contains the message text used for the notification.",
								Required:    true,
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *CredentialIssuanceRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state CredentialIssuanceRuleResourceModel

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
	CredentialIssuanceRule, d := plan.expand(ctx, r)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *credentials.CredentialIssuanceRule
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.CredentialsAPIClient.CredentialIssuanceRulesApi.CreateCredentialIssuanceRule(ctx, plan.EnvironmentId.ValueString(), plan.CredentialTypeId.ValueString()).CredentialIssuanceRule(*CredentialIssuanceRule).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateCredentialIssuanceRule",
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

func (r *CredentialIssuanceRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CredentialIssuanceRuleResourceModel

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
	var response *credentials.CredentialIssuanceRule
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.CredentialsAPIClient.CredentialIssuanceRulesApi.ReadOneCredentialIssuanceRule(ctx, data.EnvironmentId.ValueString(), data.CredentialTypeId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneCredentialIssuanceRule",
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

func (r *CredentialIssuanceRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state CredentialIssuanceRuleResourceModel

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
	CredentialIssuanceRule, d := plan.expand(ctx, r)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *credentials.CredentialIssuanceRule
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.CredentialsAPIClient.CredentialIssuanceRulesApi.UpdateCredentialIssuanceRule(ctx, plan.EnvironmentId.ValueString(), plan.CredentialTypeId.ValueString(), plan.Id.ValueString()).CredentialIssuanceRule(*CredentialIssuanceRule).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateCredentialIssuanceRule",
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

func (r *CredentialIssuanceRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CredentialIssuanceRuleResourceModel

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
			fR, fErr := r.Client.CredentialsAPIClient.CredentialIssuanceRulesApi.DeleteCredentialIssuanceRule(ctx, data.EnvironmentId.ValueString(), data.CredentialTypeId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteCredentialIssuanceRule",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CredentialIssuanceRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "credential_type_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "credential_issuance_rule_id",
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

func (p *CredentialIssuanceRuleResourceModel) expand(ctx context.Context, r *CredentialIssuanceRuleResource) (*credentials.CredentialIssuanceRule, diag.Diagnostics) {
	// The P1 Credentials service automatically sets the Issuance Rule to disabled in the backend if the Credential Type associated with it has a management.mode of `MANAGED`.
	// Perform check to prevent an out of plan / drift condition.
	diags := checkCredentialTypeManagementMode(ctx, r, p.EnvironmentId.ValueString(), p.CredentialTypeId.ValueString())
	if diags.HasError() {
		return nil, diags
	}

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

	// build issuance rule object with required attributes
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
	if !p.DigitalWalletApplicationId.IsNull() && !p.DigitalWalletApplicationId.IsUnknown() {
		application := credentials.NewCredentialIssuanceRuleDigitalWalletApplication(p.DigitalWalletApplicationId.ValueString())
		data.SetDigitalWalletApplication(*application)
	}

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
	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.CredentialTypeId = framework.PingOneResourceIDToTF(apiObject.CredentialType.GetId())
	p.Status = framework.EnumOkToTF(apiObject.GetStatusOk())

	if v, ok := apiObject.GetDigitalWalletApplicationOk(); ok {
		p.DigitalWalletApplicationId = framework.PingOneResourceIDToTF(v.GetId())
	}

	// automation object
	if v, ok := apiObject.GetAutomationOk(); ok {
		automation, d := toStateAutomation(v)
		diags.Append(d...)
		p.Automation = automation
	}

	// filter object
	if v, ok := apiObject.GetFilterOk(); ok {
		if v.HasGroupIds() || v.HasPopulationIds() || v.HasScim() { // check because values are optional
			filter, d := toStateFilter(v)
			diags.Append(d...)
			p.Filter = filter
		}
	}

	// notification object
	if v, ok := apiObject.GetNotificationOk(); ok {
		if v.HasMethods() || v.HasTemplate() { // check because values are optional
			notification, d := toStateNotification(v)
			diags.Append(d...)
			p.Notification = notification
		}
	}

	return diags
}

func toStateAutomation(automation *credentials.CredentialIssuanceRuleAutomation) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	automationMap := map[string]attr.Value{
		"issue":  framework.EnumOkToTF(automation.GetIssueOk()),
		"revoke": framework.EnumOkToTF(automation.GetRevokeOk()),
		"update": framework.EnumOkToTF(automation.GetUpdateOk()),
	}
	flattenedObj, d := types.ObjectValue(automationServiceTFObjectTypes, automationMap)
	diags.Append(d...)

	return flattenedObj, diags
}

func toStateFilter(filter *credentials.CredentialIssuanceRuleFilter) (types.Object, diag.Diagnostics) {
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

func toStateNotification(notification *credentials.CredentialIssuanceRuleNotification) (types.Object, diag.Diagnostics) {
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

func checkCredentialTypeManagementMode(ctx context.Context, r *CredentialIssuanceRuleResource, environmentId, credentialTypeId string) diag.Diagnostics {
	var diags diag.Diagnostics

	// Run the API call
	var respObject *credentials.CredentialType
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.CredentialsAPIClient.CredentialTypesApi.ReadOneCredentialType(ctx, environmentId, credentialTypeId).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environmentId, fO, fR, fErr)
		},
		"ReadOneCredentialType",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&respObject,
	)...)
	if diags.HasError() {
		return diags
	}

	if respObject == nil {
		diags.AddError(
			"Credential Type Id Invalid or Missing",
			"Creantial Type referenced in `credential_type.id` does not exist",
		)
		return diags
	}

	if v, ok := respObject.GetManagementOk(); ok {
		managementMode, managementModeOk := v.GetModeOk()
		if !managementModeOk {
			diags.AddError(
				"Credential Type referenced in `credential_type.id` does not have a management mode defined.",
				fmt.Sprintf("Credential Type Id %s does not contain a `management.mode` value, or the value could not be found. Please report this to the provider maintainers.", credentialTypeId),
			)
			return diags
		}

		if *managementMode == credentials.ENUMCREDENTIALTYPEMANAGEMENTMODE_MANAGED {
			diags.AddError(
				fmt.Sprintf("A Credential Issuance Rule cannot be assigned to a Credential Type that has a management mode of %s.", string(credentials.ENUMCREDENTIALTYPEMANAGEMENTMODE_MANAGED)),
				fmt.Sprintf("The Credential Type Id %s associated with the configured Issuance Rule is set to a `management.mode` of %s.  The Issuance Rule must be removed, or the Credential Type updated.", credentialTypeId, string(credentials.ENUMCREDENTIALTYPEMANAGEMENTMODE_MANAGED)),
			)
			return diags
		}
	}

	return diags
}
