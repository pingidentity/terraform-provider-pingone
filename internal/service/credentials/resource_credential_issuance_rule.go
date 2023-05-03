package credentials

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
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
	Automation                 types.List   `tfsdk:"automation"`
	Filter                     types.List   `tfsdk:"filter"`
	Notification               types.List   `tfsdk:"notification"`
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
	Methods  types.List `tfsdk:"methods"`
	Template types.List `tfsdk:"template"`
}

type NotificationTemplateModel struct {
	Locale    types.String `tfsdk:"locale"`
	Variables types.List   `tfsdk:"variables"`
	Variant   types.String `tfsdk:"variant"`
}

var (
	filterTypes = map[string]attr.Type{
		"group_ids":      types.SetType{ElemType: types.StringType},
		"population_ids": types.SetType{ElemType: types.StringType},
		"scim":           types.StringType,
	}

	automationTypes = map[string]attr.Type{
		"issue":  types.StringType,
		"revoke": types.StringType,
		"update": types.StringType,
	}

	notificationTypes = map[string]attr.Type{
		"methods":  types.StringType,
		"template": types.ListType{ElemType: types.StringType}, // todo - test & review
	}

	notificationTemplate = map[string]attr.Type{
		"locale":    types.StringType,
		"variables": types.ObjectType{},
		"variant":   types.StringType,
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

			"digital_wallet_application_id": framework.Attr_LinkIDWithValidators(framework.SchemaDescription{
				Description: "The ID of the digital wallet application that will interact with the user's Digital Wallet",
			},
				[]validator.String{
					verify.P1ResourceIDValidator(),
				},
			),

			"status": schema.StringAttribute{
				MarkdownDescription: "ACTIVE or DISABLED status of the credential issuance rule.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(credentials.ENUMCREDENTIALISSUANCERULESTATUS_ACTIVE),
						string(credentials.ENUMCREDENTIALISSUANCERULESTATUS_DISABLED)),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description:         "",
				MarkdownDescription: "",
				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"group_ids": schema.SetAttribute{
							ElementType:         types.StringType,
							Description:         "",
							MarkdownDescription: "",
							Optional:            true,
							Validators: []validator.Set{
								setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("population_ids")),
								setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("scim")),
							},
						},
						"population_ids": schema.SetAttribute{
							ElementType:         types.StringType,
							Description:         "",
							MarkdownDescription: "",
							Optional:            true,
							Validators: []validator.Set{
								setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("group_ids")),
								setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("scim")),
							},
						},
						"scim": schema.StringAttribute{
							Description:         "",
							MarkdownDescription: "",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("group_ids")),
								stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("population_ids")),
							},
						},
					},
				},
			},
			"automation": schema.ListNestedBlock{
				Description:         "",
				MarkdownDescription: "",
				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"issue": schema.StringAttribute{
							Description:         "",
							MarkdownDescription: "",
							Required:            true,
						},
						"revoke": schema.StringAttribute{
							Description:         "",
							MarkdownDescription: "",
							Required:            true,
						},
						"update": schema.StringAttribute{
							Description:         "",
							MarkdownDescription: "",
							Required:            true,
						},
					},
				},
			},
			"notification": schema.ListNestedBlock{
				Description:         "",
				MarkdownDescription: "",
				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"methods": schema.StringAttribute{
							Description:         "",
							MarkdownDescription: "",
							Required:            true,
						},
						"template": schema.ObjectAttribute{
							Description:         "",
							MarkdownDescription: "",
							Required:            true,
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

	/*management.ReadOneApplicationRequest(ctx, )
	if CredentialIssuanceRule.GetApplication().Id{
			// make sure it exists

		}

	    t.GetOk("oidc_options"); ok {
		    var application *management.ApplicationOIDC
		    application, diags = expandApplicationOIDC(d)
		    if diags.HasError() {
		        return diags
		    }
		    applicationRequest.ApplicationOIDC = application
		} */

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
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/credentials_credential_type_id/\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *CredentialIssuanceRuleResourceModel) expand(ctx context.Context) (*credentials.CredentialIssuanceRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	// expand automation
	// todo: move to function
	var automationRules []AutomationModel
	diags.Append(p.Automation.ElementsAs(ctx, &automationRules, false)...)
	if diags.HasError() {
		return nil, diags
	}
	automation := credentials.NewCredentialIssuanceRuleAutomationWithDefaults()
	for _, v := range automationRules {
		if !v.Issue.IsNull() && !v.Issue.IsUnknown() {
			automation.SetIssue(credentials.EnumCredentialIssuanceRuleAutomationMethod(v.Issue.ValueString()))
		}

		if !v.Revoke.IsNull() && !v.Revoke.IsUnknown() {
			automation.SetRevoke(credentials.EnumCredentialIssuanceRuleAutomationMethod(v.Revoke.ValueString()))
		}

		if !v.Update.IsNull() && !v.Update.IsUnknown() {
			automation.SetUpdate(credentials.EnumCredentialIssuanceRuleAutomationMethod(v.Update.ValueString()))
		}
	}

	// expand filter
	// todo: move to function
	var filterRules []FilterModel
	diags.Append(p.Filter.ElementsAs(ctx, &filterRules, false)...)
	if diags.HasError() {
		return nil, diags
	}
	filter := credentials.NewCredentialIssuanceRuleFilterWithDefaults()
	for _, v := range filterRules {
		if !v.PopulationIds.IsNull() && !v.PopulationIds.IsUnknown() {
			diags.Append(v.PopulationIds.ElementsAs(ctx, &filter.PopulationIds, false)...)
			if diags.HasError() {
				return nil, diags
			}
			filter.SetPopulationIds(filter.PopulationIds)
		}

		if !v.GroupIds.IsNull() && !v.GroupIds.IsUnknown() {
			diags.Append(v.GroupIds.ElementsAs(ctx, &filter.GroupIds, false)...)
			if diags.HasError() {
				return nil, diags
			}
			filter.SetGroupIds(filter.GroupIds)
		}

		if !v.Scim.IsNull() && !v.Scim.IsUnknown() {
			filter.SetScim(v.Scim.ValueString())
		}
	}

	// much to do...
	data := credentials.NewCredentialIssuanceRule(*automation, credentials.EnumCredentialIssuanceRuleStatus(p.Status.ValueString()))

	// set digital wallet if present - what if not present
	application := credentials.NewCredentialIssuanceRuleDigitalWalletApplication(p.DigitalWalletApplicationId.ValueString())
	data.SetDigitalWalletApplication(*application)

	// set filter
	data.SetFilter(*filter)

	return data, diags
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

	p.Id = framework.StringToTF(apiObject.GetId())
	p.EnvironmentId = framework.StringToTF(apiObject.GetEnvironment().Id)
	p.DigitalWalletApplicationId = framework.StringToTF((apiObject.GetDigitalWalletApplication().Id))
	p.CredentialTypeId = framework.StringToTF((apiObject.CredentialType.Id))

	// automation
	// move to function
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	tfObjType := types.ObjectType{AttrTypes: automationTypes}
	automationMap := map[string]attr.Value{
		"issue":  framework.StringToTF(string(apiObject.GetAutomation().Issue)),
		"revoke": framework.StringToTF(string(apiObject.GetAutomation().Revoke)),
		"update": framework.StringToTF(string(apiObject.GetAutomation().Update)),
	}
	flattenedObj, d := types.ObjectValue(automationTypes, automationMap)
	diags.Append(d...)

	automation, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)
	p.Automation = automation

	// fields
	// move to function
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	tfObjType = types.ObjectType{AttrTypes: filterTypes}
	filterMap := map[string]attr.Value{
		"population_ids": framework.StringSetOkToTF(apiObject.Filter.GetPopulationIdsOk()),
		"group_ids":      framework.StringSetOkToTF(apiObject.Filter.GetGroupIdsOk()),
		"scim":           framework.StringOkToTF(apiObject.Filter.GetScimOk()),
	}
	flattenedObj, d = types.ObjectValue(filterTypes, filterMap)
	diags.Append(d...)

	filter, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)
	p.Filter = filter

	return diags
}
