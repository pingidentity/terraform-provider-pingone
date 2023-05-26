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
type VerifyPolicyResource struct {
	client *verify.APIClient
	region model.RegionMapping
}

type VerifyPolicyResourceModel struct {
	Id               types.String `tfsdk:"id"`
	EnvironmentId    types.String `tfsdk:"environment_id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Default          types.Bool   `tfsdk:"default"`
	GovernmentId     types.Object `tfsdk:"government_id"`
	FacialComparison types.Object `tfsdk:"facial_comparison"`
	Liveness         types.Object `tfsdk:"liveness"`
	Email            types.Object `tfsdk:"email"`
	Phone            types.Object `tfsdk:"phone"`
	Transaction      types.Object `tfsdk:"transaction"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}

type GovernmentIdModel struct {
	Verify types.String `tfsdk:"verify"`
}

type FacialComparisonModel struct {
	Verify    types.String `tfsdk:"verify"`
	Threshold types.String `tfsdk:"threshold"`
}

type LivenessnModel struct {
	Verify    types.String `tfsdk:"verify"`
	Threshold types.String `tfsdk:"threshold"`
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////
// Is there a better way for this blockto handle these nested objects
type EmailModel struct {
	Verify          types.String `tfsdk:"verify"`
	CreateMfaDevice types.Bool   `tfsdk:"create_mfa_device"`
	OTP             types.Object `tfsdk:"otp"`
}

type PhoneModel struct {
	Verify          types.String `tfsdk:"verify"`
	CreateMfaDevice types.Bool   `tfsdk:"create_mfa_device"`
	OTP             types.Object `tfsdk:"otp"`
}

type OTPAttemptsModel struct {
	Count types.Int64 `tfsdk:"count"`
}

type OTPDeliveriessModel struct {
	Count    types.Int64  `tfsdk:"count"`
	Cooldown types.Object `tfsdk:"cooldown"`
}

type OTPDeliveriessCooldownModel struct {
	Duration types.Int64  `tfsdk:"duration"`
	TimeUnit types.Object `tfsdk:"time_unit"`
}

type OTPLifeTmeModel struct {
	Duration types.Int64  `tfsdk:"duration"`
	TimeUnit types.Object `tfsdk:"time_unit"`
}

type OTPNotificationModel struct {
	TemplateName types.Int64  `tfsdk:"template_name"`
	VariantName  types.Object `tfsdk:"variant_name"`
}

type TransactionModel struct {
	Timeout            types.Object `tfsdk:"timeout"`
	DataCollection     types.Object `tfsdk:"data_collection"`
	DataCollectionOnly types.Bool   `tfsdk:"data_collection_only"`
}

type TransactionTimeoutModel struct {
	Duration types.Int64  `tfsdk:"duration"`
	TimeUnit types.Object `tfsdk:"time_unit"`
}

type TransactionDataCollectionModel struct {
	Timeout types.Object `tfsdk:"timeout"`
}

type TransactionDataCollectionTimeoutModel struct {
	Duration types.Int64  `tfsdk:"duration"`
	TimeUnit types.Object `tfsdk:"time_unit"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

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
	_ resource.Resource                = &VerifyPolicyResource{}
	_ resource.ResourceWithConfigure   = &VerifyPolicyResource{}
	_ resource.ResourceWithImportState = &VerifyPolicyResource{}
)

// New Object
func NewVerifyPolicyResource() resource.Resource {
	return &VerifyPolicyResource{}
}

// Metadata
func (r *VerifyPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_issuance_rule"
}

func (r *VerifyPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

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

			"digital_wallet_application_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("Identifier (UUID) of the customer's Digital Wallet App that will interact with the user's Digital Wallet."),
			),

			"status": schema.StringAttribute{
				Description:         statusDescription.Description,
				MarkdownDescription: statusDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(credentials.ENUMVerifyPolicySTATUS_ACTIVE),
						string(credentials.ENUMVerifyPolicySTATUS_DISABLED)),
				},
			},

			"filter": schema.SingleNestedAttribute{
				Description:         filterDescription.Description,
				MarkdownDescription: filterDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.Object{
					customobjectvalidator.IsRequiredIfMatchesPathValue(
						basetypes.NewStringValue(string(credentials.ENUMVerifyPolicySTATUS_ACTIVE)),
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
									string(credentials.ENUMVerifyPolicyNOTIFICATIONMETHOD_EMAIL),
									string(credentials.ENUMVerifyPolicyNOTIFICATIONMETHOD_SMS),
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

func (r *VerifyPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VerifyPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state VerifyPolicyResourceModel

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
	VerifyPolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.VerifyPolicysApi.CreateVerifyPolicy(ctx, plan.EnvironmentId.ValueString(), plan.CredentialTypeId.ValueString()).VerifyPolicy(*VerifyPolicy).Execute()
		},
		"CreateVerifyPolicy",
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
	resp.Diagnostics.Append(state.toState(response.(*credentials.VerifyPolicy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *VerifyPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *VerifyPolicyResourceModel

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
			return r.client.VerifyPolicysApi.ReadOneVerifyPolicy(ctx, data.EnvironmentId.ValueString(), data.CredentialTypeId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneVerifyPolicy",
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
	resp.Diagnostics.Append(data.toState(response.(*credentials.VerifyPolicy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VerifyPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state VerifyPolicyResourceModel

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
	VerifyPolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.VerifyPolicysApi.UpdateVerifyPolicy(ctx, plan.EnvironmentId.ValueString(), plan.CredentialTypeId.ValueString(), plan.Id.ValueString()).VerifyPolicy(*VerifyPolicy).Execute()
		},
		"UpdateVerifyPolicy",
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
	resp.Diagnostics.Append(state.toState(response.(*credentials.VerifyPolicy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *VerifyPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *VerifyPolicyResourceModel

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
			r, err := r.client.VerifyPolicysApi.DeleteVerifyPolicy(ctx, data.EnvironmentId.ValueString(), data.CredentialTypeId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteVerifyPolicy",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *VerifyPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (p *VerifyPolicyResourceModel) expand(ctx context.Context) (*credentials.VerifyPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	// expand automation rules
	VerifyPolicyAutomation := credentials.NewVerifyPolicyAutomationWithDefaults()
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

		VerifyPolicyAutomation, d = automationRules.expandAutomationModel()
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// expand filter
	VerifyPolicyFilter := credentials.NewVerifyPolicyFilterWithDefaults()
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

		VerifyPolicyFilter, d = filterRules.expandFilterModel(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// expand notifications
	VerifyPolicyNotification := credentials.NewVerifyPolicyNotificationWithDefaults()
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

		VerifyPolicyNotification, d = notificationRules.expandNotificationModel(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// buuild issuance rule object with required attributes
	data := credentials.NewVerifyPolicy(*VerifyPolicyAutomation, credentials.EnumVerifyPolicyStatus(p.Status.ValueString()))

	// set the filter details
	if VerifyPolicyFilter.HasGroupIds() || VerifyPolicyFilter.HasPopulationIds() || VerifyPolicyFilter.HasScim() {
		data.SetFilter(*VerifyPolicyFilter)
	}

	// set the notification details
	if VerifyPolicyNotification.HasMethods() || VerifyPolicyNotification.HasTemplate() {
		data.SetNotification(*VerifyPolicyNotification)
	}

	// set the digital wallet application
	application := credentials.NewVerifyPolicyDigitalWalletApplication(p.DigitalWalletApplicationId.ValueString())
	data.SetDigitalWalletApplication(*application)

	return data, diags
}

func (p *AutomationModel) expandAutomationModel() (*credentials.VerifyPolicyAutomation, diag.Diagnostics) {
	var diags diag.Diagnostics

	automation := credentials.NewVerifyPolicyAutomationWithDefaults()

	if !p.Issue.IsNull() && !p.Issue.IsUnknown() {
		automation.SetIssue(credentials.EnumVerifyPolicyAutomationMethod(p.Issue.ValueString()))
	}

	if !p.Revoke.IsNull() && !p.Revoke.IsUnknown() {
		automation.SetRevoke(credentials.EnumVerifyPolicyAutomationMethod(p.Revoke.ValueString()))
	}

	if !p.Update.IsNull() && !p.Update.IsUnknown() {
		automation.SetUpdate(credentials.EnumVerifyPolicyAutomationMethod(p.Update.ValueString()))
	}

	if automation == nil {
		diags.AddWarning(
			"Unexpected Value",
			"Automation object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}
	return automation, diags

}

func (p *FilterModel) expandFilterModel(ctx context.Context) (*credentials.VerifyPolicyFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	filter := credentials.NewVerifyPolicyFilterWithDefaults()

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

func (p *NotificationModel) expandNotificationModel(ctx context.Context) (*credentials.VerifyPolicyNotification, diag.Diagnostics) {
	var diags diag.Diagnostics

	notification := credentials.NewVerifyPolicyNotificationWithDefaults()

	// notification methods
	if !p.Methods.IsNull() && !p.Methods.IsUnknown() {
		var slice []string
		diags.Append(p.Methods.ElementsAs(ctx, &slice, false)...)

		enumSlice := make([]credentials.EnumVerifyPolicyNotificationMethod, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := credentials.NewEnumVerifyPolicyNotificationMethodFromValue(slice[i])
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

		template := credentials.NewVerifyPolicyNotificationTemplate()
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

func (p *VerifyPolicyResourceModel) toState(apiObject *credentials.VerifyPolicy) diag.Diagnostics {
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

func toStateAutomation(automation *credentials.VerifyPolicyAutomation) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	automationMap := map[string]attr.Value{
		"issue":  enumVerifyPolicyAutomationOkToTF(automation.GetIssueOk()),
		"revoke": enumVerifyPolicyAutomationOkToTF(automation.GetRevokeOk()),
		"update": enumVerifyPolicyAutomationOkToTF(automation.GetUpdateOk()),
	}
	flattenedObj, d := types.ObjectValue(automationServiceTFObjectTypes, automationMap)
	diags.Append(d...)

	return flattenedObj, diags
}

func toStateFilter(filter *credentials.VerifyPolicyFilter) (types.Object, diag.Diagnostics) {
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

func toStateNotification(notification *credentials.VerifyPolicyNotification) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if notification == nil {
		return types.ObjectNull(notificationServiceTFObjectTypes), diags
	}

	notificationMap := map[string]attr.Value{}

	// notification.methods
	if v, ok := notification.GetMethodsOk(); ok {
		notificationMap["methods"] = enumVerifyPolicyNotificationMethodOkToTF(v, ok)
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

func enumVerifyPolicyNotificationMethodOkToTF(v []credentials.EnumVerifyPolicyNotificationMethod, ok bool) basetypes.SetValue {
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

func enumVerifyPolicyAutomationOkToTF(v *credentials.EnumVerifyPolicyAutomationMethod, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}

func enumCredentialIssuanceStatusOkToTF(v *credentials.EnumVerifyPolicyStatus, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
