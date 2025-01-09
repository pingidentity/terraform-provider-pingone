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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type PolicyManagementPolicyResource serviceClientType

type policyManagementPolicyResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	// Type          types.String                 `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	// Statements         types.List                   `tfsdk:"statements"`
	Condition          types.Object `tfsdk:"condition"`
	CombiningAlgorithm types.Object `tfsdk:"combining_algorithm"`
	// Children           types.List   `tfsdk:"children"`
	RepetitionSettings types.Object `tfsdk:"repetition_settings"`
	// ManagedEntity      types.Object `tfsdk:"managed_entity"`
	Version types.String `tfsdk:"version"`
}

type policyManagementPolicyCombiningAlgorithmResourceModel struct {
	Algorithm types.String `tfsdk:"algorithm"`
}

// type policyManagementPolicyChildrenResourceModel struct{}

type policyManagementPolicyRepetitionSettingsResourceModel struct {
	Source   types.Object `tfsdk:"source"`
	Decision types.String `tfsdk:"decision"`
}

var (
	policyManagementPolicyCombiningAlgorithmTFObjectTypes = map[string]attr.Type{
		"algorithm": types.StringType,
	}

	policyManagementPolicyRepetitionSettingsTFObjectTypes = map[string]attr.Type{
		"source":   types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
		"decision": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &PolicyManagementPolicyResource{}
	_ resource.ResourceWithConfigure   = &PolicyManagementPolicyResource{}
	_ resource.ResourceWithImportState = &PolicyManagementPolicyResource{}
)

// New Object
func NewPolicyManagementPolicyResource() resource.Resource {
	return &PolicyManagementPolicyResource{}
}

// Metadata
func (r *PolicyManagementPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_policy_management_policy"
}

func (r *PolicyManagementPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the policy is enabled, and whether the policy is evaluated.",
	).DefaultValue(true)

	combiningAlgorithmAlgorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the algorithm that determines how rules are combined to produce an authorization decision.",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataPoliciesCombiningAlgorithmDTOAlgorithmEnumValues)

	repetitionSettingsDecisionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the decision filter.",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataPoliciesRepetitionSettingsDTODecisionEnumValues)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an authorization policy for the PingOne Authorize Policy Manager in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the Authorize editor policy in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a user-friendly name to apply to the authorization policy.  The value must be unique.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			// "type": schema.StringAttribute{
			// 	Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the type of the policy.").Description,
			// 	Optional:    true,
			// },

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a description to apply to the policy.").Description,
				Required:    true,
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
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings for an authorization condition to apply to the policy.").Description,
				Optional:    true,

				Attributes: dataConditionObjectSchemaAttributes(),
			},

			"combining_algorithm": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings that determine how rules are combined to produce an authorization decision.").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"algorithm": schema.StringAttribute{
						Description:         combiningAlgorithmAlgorithmDescription.Description,
						MarkdownDescription: combiningAlgorithmAlgorithmDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataPoliciesCombiningAlgorithmDTOAlgorithmEnumValues)...),
						},
					},
				},
			},

			// "children": schema.ListNestedAttribute{
			// 	Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
			// 	Optional:    true,

			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{},
			// 	},
			// },

			"repetition_settings": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings that appies the policy to each item of the specific attribute, filtered by decision.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"source": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings for the source associated with the policy rule.").Description,
						Required:    true,

						Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the authorization policy repetition source in the policy manager.")),
					},

					"decision": schema.StringAttribute{
						Description:         repetitionSettingsDecisionDescription.Description,
						MarkdownDescription: repetitionSettingsDecisionDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataPoliciesRepetitionSettingsDTODecisionEnumValues)...),
						},
					},
				},
			},

			// "managed_entity": schema.SingleNestedAttribute{
			// 	Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
			// 	Optional:    true,

			// 	Attributes: managedEntityObjectSchemaAttributes(),
			// },

			"version": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a random ID generated by the system for concurrency control purposes.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *PolicyManagementPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PolicyManagementPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state policyManagementPolicyResourceModel

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
	policyManagementPolicy, d := plan.expandCreate(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorPoliciesApi.CreatePolicy(ctx, plan.EnvironmentId.ValueString()).AuthorizeEditorDataPoliciesPolicyDTO(*policyManagementPolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreatePolicy",
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

func (r *PolicyManagementPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *policyManagementPolicyResourceModel

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
	var response *authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorPoliciesApi.GetPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetPolicy",
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

func (r *PolicyManagementPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state policyManagementPolicyResourceModel

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
	var getResponse *authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorPoliciesApi.GetPolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetPolicy-Update",
		framework.DefaultCustomError,
		retryAuthorizeEditorCreateUpdate,
		&getResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
	version := getResponse.GetVersion()

	// Build the model for the API
	policyManagementPolicy, d := plan.expandUpdate(ctx, version)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorPoliciesApi.UpdatePolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).AuthorizeEditorDataPoliciesReferenceablePolicyDTO(*policyManagementPolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdatePolicy",
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

func (r *PolicyManagementPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *policyManagementPolicyResourceModel

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
					fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorPoliciesApi.DeletePolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
				},
				"DeletePolicy",
				framework.CustomErrorResourceNotFoundWarning,
				nil,
				nil,
			)...)
			if resp.Diagnostics.HasError() {
				return nil, "ERROR", fmt.Errorf("Error deleting authorize policy (%s)", data.Id.ValueString())
			}

			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorPoliciesApi.GetPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
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
			"Authorize Policy Delete Timeout",
			fmt.Sprintf("Error waiting for authorize policy (%s) to be deleted: %s", data.Id.ValueString(), err),
		)

		return
	}
}

func (r *PolicyManagementPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "authorize_policy_management_policy_id",
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

func (p *policyManagementPolicyResourceModel) expandCreate(ctx context.Context) (*authorize.AuthorizeEditorDataPoliciesPolicyDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	combiningAlgorithm := &authorize.AuthorizeEditorDataPoliciesCombiningAlgorithmDTO{}

	data := authorize.NewAuthorizeEditorDataPoliciesPolicyDTO(
		p.Name.ValueString(),
		*combiningAlgorithm,
	)

	// if !p.Type.IsNull() && !p.Type.IsUnknown() {
	// 	data.SetType(p.Type.ValueString())
	// }

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.Enabled.IsNull() && !p.Enabled.IsUnknown() {
		data.SetEnabled(p.Enabled.ValueBool())
	}

	// if !p.Statements.IsNull() && !p.Statements.IsUnknown() {
	// 	var plan []policyManagementPolicyStatementResourceModel
	// 	diags.Append(p.Statements.ElementsAs(ctx, &plan, false)...)
	// 	if diags.HasError() {
	// 		return nil, diags
	// 	}

	// 	statements := make([]map[string]interface{}, 0)
	// 	for _, planItem := range plan {
	// 		statements = append(statements, planItem.expand())
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

	if !p.CombiningAlgorithm.IsNull() && !p.CombiningAlgorithm.IsUnknown() {
		var plan *policyManagementPolicyCombiningAlgorithmResourceModel
		diags.Append(p.CombiningAlgorithm.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		combiningAlgorithm := plan.expand()

		data.SetCombiningAlgorithm(*combiningAlgorithm)
	}

	// if !p.Children.IsNull() && !p.Children.IsUnknown() {
	// 	var plan []policyManagementPolicyChildrenResourceModel
	// 	diags.Append(p.Children.ElementsAs(ctx, &plan, false)...)
	// 	if diags.HasError() {
	// 		return nil, diags
	// 	}

	// 	children := make([]map[string]interface{}, 0)
	// 	for _, planItem := range plan {
	// 		children = append(children, planItem.expand())
	// 	}

	// 	data.SetChildren(children)
	// }

	if !p.RepetitionSettings.IsNull() && !p.RepetitionSettings.IsUnknown() {
		var plan *policyManagementPolicyRepetitionSettingsResourceModel
		diags.Append(p.RepetitionSettings.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		repetitionSettings, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetRepetitionSettings(*repetitionSettings)
	}

	// if !p.ManagedEntity.IsNull() && !p.ManagedEntity.IsUnknown() {

	// 	managedEntity, d := expandEditorManagedEntity(ctx, p.ManagedEntity)
	// 	diags.Append(d...)
	// 	if diags.HasError() {
	// 		return nil, diags
	// 	}

	// 	data.SetManagedEntity(*managedEntity)
	// }

	return data, diags
}

func (p *policyManagementPolicyResourceModel) expandUpdate(ctx context.Context, versionId string) (*authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	dataCreate, d := p.expandCreate(ctx)
	if d.HasError() {
		return nil, d
	}

	// Use json.marshall and unmarshal to cast dataCreate to a AuthorizeEditorDataRulesReferenceableRuleDTO type
	bytes, err := json.Marshal(dataCreate)
	if err != nil {
		diags.AddError("Failed to marshal data", err.Error())
		return nil, diags
	}

	var data *authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO
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

// func (p *policyManagementPolicyStatementResourceModel) expand() map[string]interface{} {

// 	log.Panicf("Not implemented")

// 	return nil
// }

func (p *policyManagementPolicyCombiningAlgorithmResourceModel) expand() *authorize.AuthorizeEditorDataPoliciesCombiningAlgorithmDTO {

	data := authorize.NewAuthorizeEditorDataPoliciesCombiningAlgorithmDTO(
		authorize.EnumAuthorizeEditorDataPoliciesCombiningAlgorithmDTOAlgorithm(p.Algorithm.ValueString()),
	)

	return data
}

// func (p *policyManagementPolicyChildrenResourceModel) expand() map[string]interface{} {

// 	log.Panicf("Not implemented")

// 	return nil
// }

func (p *policyManagementPolicyRepetitionSettingsResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataPoliciesRepetitionSettingsDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	source, d := expandEditorReferenceData(ctx, p.Source)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataPoliciesRepetitionSettingsDTO(
		*source,
		authorize.EnumAuthorizeEditorDataPoliciesRepetitionSettingsDTODecision(p.Decision.ValueString()),
	)

	return data, diags
}

func (p *policyManagementPolicyResourceModel) toState(ctx context.Context, apiObject *authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO) diag.Diagnostics {
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
	// p.Type = framework.EnumOkToTF(apiObject.GetTypeOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())
	// p.Statements = framework.ListOkToTF(apiObject.GetStatementsOk())

	conditionVal, ok := apiObject.GetConditionOk()
	p.Condition, d = editorDataConditionOkToTF(ctx, conditionVal, ok)
	diags.Append(d...)

	p.CombiningAlgorithm, d = policyManagementPolicyCombiningAlgorithmOkToTF(apiObject.GetCombiningAlgorithmOk())
	diags.Append(d...)

	// p.Children = framework.ListOkToTF(apiObject.GetChildrenOk())
	p.RepetitionSettings, d = policyManagementPolicyRepetitionSettingsOkToTF(apiObject.GetRepetitionSettingsOk())
	diags.Append(d...)

	// p.ManagedEntity, d = editorManagedEntityOkToTF(apiObject.GetManagedEntityOk())
	// diags.Append(d...)

	p.Version = framework.StringOkToTF(apiObject.GetVersionOk())

	return diags
}

func policyManagementPolicyCombiningAlgorithmOkToTF(apiObject *authorize.AuthorizeEditorDataPoliciesCombiningAlgorithmDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(policyManagementPolicyCombiningAlgorithmTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(policyManagementPolicyCombiningAlgorithmTFObjectTypes, map[string]attr.Value{
		"algorithm": framework.EnumOkToTF(apiObject.GetAlgorithmOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func policyManagementPolicyRepetitionSettingsOkToTF(apiObject *authorize.AuthorizeEditorDataPoliciesRepetitionSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(policyManagementPolicyRepetitionSettingsTFObjectTypes), diags
	}

	source, d := editorDataReferenceObjectOkToTF(apiObject.GetSourceOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(policyManagementPolicyRepetitionSettingsTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(policyManagementPolicyRepetitionSettingsTFObjectTypes, map[string]attr.Value{
		"source":   source,
		"decision": framework.EnumOkToTF(apiObject.GetDecisionOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
