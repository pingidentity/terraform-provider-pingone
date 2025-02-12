package authorize

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type PolicyManagementRuleResource serviceClientType

type policyManagementRuleResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name          types.String                 `tfsdk:"name"`
	Description   types.String                 `tfsdk:"description"`
	Enabled       types.Bool                   `tfsdk:"enabled"`
	// Statements     types.List   `tfsdk:"statements"`
	Condition      types.Object `tfsdk:"condition"`
	EffectSettings types.Object `tfsdk:"effect_settings"`
	Version        types.String `tfsdk:"version"`
}

type policyManagementRuleEffectSettingsResourceModel struct {
	Type types.String `tfsdk:"type"`
}

// Framework interfaces
var (
	_ resource.Resource                = &PolicyManagementRuleResource{}
	_ resource.ResourceWithConfigure   = &PolicyManagementRuleResource{}
	_ resource.ResourceWithImportState = &PolicyManagementRuleResource{}
)

// New Object
func NewPolicyManagementRuleResource() resource.Resource {
	return &PolicyManagementRuleResource{}
}

// Metadata
func (r *PolicyManagementRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_policy_management_rule"
}

func (r *PolicyManagementRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the authorization rule is enabled and is evaluated.",
	).DefaultValue(true)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an authorization rule for the PingOne Authorize Policy Manager in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the Authorize editor rule in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a user-friendly name for the authorization rule.  The value must be unique.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a description to apply to the authorization rule.").Description,
				Optional:    true,
			},

			"enabled": schema.BoolAttribute{
				Description:         enabledDescription.Description,
				MarkdownDescription: enabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(true),
			},

			// "statements": schema.ListNestedAttribute{
			// 	Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
			// 	Optional:    true,

			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{},
			// 	},
			// },

			"condition": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings for conditions to apply to the authorization rule.").Description,
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(func() basetypes.ObjectValue {
					attributeMap := map[string]attr.Value{
						"type": types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATACONDITIONDTOTYPE_EMPTY)),
					}

					attributeMap = editorDataConditionConvertEmptyValuesToTFNulls(attributeMap, 1)

					return types.ObjectValueMust(editorDataConditionTFObjectTypes, attributeMap)
				}()),

				Attributes: dataConditionObjectSchemaAttributes(),
			},

			"effect_settings": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings that determine how child rules are combined to produce an outcome for the policy.").Description,
				Required:    true,

				Attributes: dataRulesEffectSettingsObjectSchemaAttributes(),
			},

			"version": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that describes a random ID generated by the system for concurrency control purposes.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *PolicyManagementRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PolicyManagementRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state policyManagementRuleResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
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
	policyManagementRule, d := plan.expandCreate(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataRulesReferenceableRuleDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorRulesApi.CreateRule(ctx, plan.EnvironmentId.ValueString()).AuthorizeEditorDataRulesRuleDTO(*policyManagementRule).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateRule",
		framework.DefaultCustomError,
		retryAuthorizeEditorCreateUpdate,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	}
}

func (r *PolicyManagementRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *policyManagementRuleResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
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
	var response *authorize.AuthorizeEditorDataRulesReferenceableRuleDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorRulesApi.GetRule(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetRule",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
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
	resp.Diagnostics.Append(data.toState(ctx, response)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *PolicyManagementRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state policyManagementRuleResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
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

	// Run the API call
	var getResponse *authorize.AuthorizeEditorDataRulesReferenceableRuleDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorRulesApi.GetRule(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetRule-Update",
		framework.DefaultCustomError,
		retryAuthorizeEditorCreateUpdate,
		&getResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	version := getResponse.GetVersion()

	// Build the model for the API
	policyManagementRule, d := plan.expandUpdate(ctx, version)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataRulesReferenceableRuleDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorRulesApi.UpdateRule(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).AuthorizeEditorDataRulesReferenceableRuleDTO(*policyManagementRule).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateRule",
		framework.DefaultCustomError,
		nil,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	}
}

func (r *PolicyManagementRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *policyManagementRuleResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
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

	deleteStateConf := &retry.StateChangeConf{
		Pending: []string{
			"200",
		},
		Target: []string{
			"404",
			"ERROR",
		},
		Refresh: func() (interface{}, string, error) {
			// Run the API call
			resp.Diagnostics.Append(framework.ParseResponse(
				ctx,

				func() (any, *http.Response, error) {
					fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorRulesApi.DeleteRule(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
				},
				"DeleteRule",
				framework.CustomErrorResourceNotFoundWarning,
				nil,
				nil,
			)...)
			if resp.Diagnostics.HasError() {
				return nil, "ERROR", fmt.Errorf("Error deleting authorize rule (%s)", data.Id.ValueString())
			}

			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorRulesApi.GetRule(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			getResp, r, err := framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)

			if err != nil || r == nil {
				return getResp, "ERROR", err
			}

			base := 10
			return getResp, strconv.FormatInt(int64(r.StatusCode), base), nil
		},
		Timeout:                   20 * time.Minute,
		Delay:                     1 * time.Second,
		MinTimeout:                500 * time.Millisecond,
		ContinuousTargetOccurence: 2,
	}
	_, err := deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Authorize Rule Delete Timeout",
			fmt.Sprintf("Error waiting for authorize rule (%s) to be deleted: %s", data.Id.ValueString(), err),
		)

		return
	}
}

func (r *PolicyManagementRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "authorize_policy_management_rule_id",
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

func (p *policyManagementRuleResourceModel) expandCreate(ctx context.Context) (*authorize.AuthorizeEditorDataRulesRuleDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	effectSettings, d := expandEditorDataRulesEffectSettings(ctx, p.EffectSettings)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Main object
	data := authorize.NewAuthorizeEditorDataRulesRuleDTO(
		p.Name.ValueString(),
		*effectSettings,
	)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.Enabled.IsNull() && !p.Enabled.IsUnknown() {
		data.SetEnabled(p.Enabled.ValueBool())
	}

	// if !p.Statements.IsNull() && !p.Statements.IsUnknown() {
	// 	var plan []policyManagementRuleStatementResourceModel
	// 	diags.Append(p.Statements.ElementsAs(ctx, &plan, false)...)
	// 	if diags.HasError() {
	// 		return nil, diags
	// 	}

	// 	statements := make([]map[string]interface{}, 0, len(plan))
	// 	for _, statementPlan := range plan {
	// 		statement := statementPlan.expand()

	// 		statements = append(statements, statement)
	// 	}

	// 	data.SetStatements(statements)
	// }

	if !p.Condition.IsNull() && !p.Condition.IsUnknown() {
		condition, d := expandEditorDataCondition(ctx, p.Condition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetCondition(*condition)
	}

	return data, diags
}

func (p *policyManagementRuleResourceModel) expandUpdate(ctx context.Context, versionId string) (*authorize.AuthorizeEditorDataRulesReferenceableRuleDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	dataCreate, d := p.expandCreate(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Use json.marshall and unmarshal to cast dataCreate to a AuthorizeEditorDataRulesReferenceableRuleDTO type
	bytes, err := json.Marshal(dataCreate)
	if err != nil {
		diags.AddError("Failed to marshal data", err.Error())
		return nil, diags
	}

	var data *authorize.AuthorizeEditorDataRulesReferenceableRuleDTO
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		diags.AddError("Failed to unmarshal data", err.Error())
		return nil, diags
	}

	if versionId != "" {
		data.SetVersion(versionId)

		if !p.Id.IsNull() && !p.Id.IsUnknown() {
			data.SetId(p.Id.ValueString())
		}
	}

	return data, diags
}

func (p *policyManagementRuleResourceModel) toState(ctx context.Context, apiObject *authorize.AuthorizeEditorDataRulesReferenceableRuleDTO) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	// p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())

	// p.Statements, d = policyManagementRuleStatementsOkToTF(apiObject.GetStatementsOk())
	// diags.Append(d...)

	conditionVal, ok := apiObject.GetConditionOk()
	p.Condition, d = editorDataConditionOkToTF(ctx, conditionVal, ok)
	diags.Append(d...)

	effectSettingsVal, ok := apiObject.GetEffectSettingsOk()
	p.EffectSettings, d = editorDataRulesEffectSettingsOkToTF(ctx, effectSettingsVal, ok)
	diags.Append(d...)

	p.Version = framework.StringOkToTF(apiObject.GetVersionOk())

	return diags
}
