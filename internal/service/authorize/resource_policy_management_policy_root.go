package authorize

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type PolicyManagementPolicyRootResource serviceClientType

type policyManagementPolicyRootResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	// Type          types.String                 `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	// Statements         types.List                   `tfsdk:"statements"`
	Condition          types.Object `tfsdk:"condition"`
	CombiningAlgorithm types.Object `tfsdk:"combining_algorithm"`
	Children           types.List   `tfsdk:"children"`
	// RepetitionSettings types.Object `tfsdk:"repetition_settings"`
	// ManagedEntity      types.Object `tfsdk:"managed_entity"`
	Version types.String `tfsdk:"version"`
}

// Framework interfaces
var (
	_ resource.Resource                = &PolicyManagementPolicyRootResource{}
	_ resource.ResourceWithConfigure   = &PolicyManagementPolicyRootResource{}
	_ resource.ResourceWithImportState = &PolicyManagementPolicyRootResource{}
)

// New Object
func NewPolicyManagementPolicyRootResource() resource.Resource {
	return &PolicyManagementPolicyRootResource{}
}

// Metadata
func (r *PolicyManagementPolicyRootResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_policy_management_policy_root"
}

func (r *PolicyManagementPolicyRootResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the policy is enabled, and whether the policy is evaluated.",
	).DefaultValue(true)

	// repetitionSettingsDecisionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
	// 	"A string that specifies the decision filter.",
	// ).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataPoliciesRepetitionSettingsDTODecisionEnumValues)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage the root authorization policy for the PingOne Authorize Policy Manager in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the Authorize editor policy in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a user-friendly name to apply to the authorization policy.  The value must be unique.").Description,
				Optional:    true,
				Computed:    true,

				Default: stringdefault.StaticString("Policies"),

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			// "type": schema.StringAttribute{
			// 	Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the type of the policy.").Description,
			// 	Optional:    true,
			// 	Computed: true,

			// 	Default: stringdefault.StaticString("POLICY"),
			// },

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a description to apply to the policy.").Description,
				Optional:    true,
				Computed:    true,
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

				Attributes: combiningAlgorithmObjectSchemaAttributes(),
			},

			"children": schema.ListNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An ordered list of objects that specifies child policies or policy sets.").Description,
				Optional:    true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: dataPolicyObjectSchemaAttributes(),
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

func (r *PolicyManagementPolicyRootResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PolicyManagementPolicyRootResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state policyManagementPolicyRootResourceModel

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
	getResponse, d := r.getRootPolicy(ctx, plan.EnvironmentId.ValueString())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	version := getResponse.GetVersion()

	// Build the model for the API
	policyManagementPolicyRoot, d := plan.expandUpdate(ctx, version)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyManagementPolicyRoot.SetId(getResponse.GetId())

	// Run the API call
	var response *authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorPoliciesApi.UpdatePolicy(ctx, plan.EnvironmentId.ValueString(), getResponse.GetId()).AuthorizeEditorDataPoliciesReferenceablePolicyDTO(*policyManagementPolicyRoot).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdatePolicy-Create",
		framework.DefaultCustomError,
		nil,
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

func (r *PolicyManagementPolicyRootResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *policyManagementPolicyRootResourceModel

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

	response, d := r.getRootPolicy(ctx, data.EnvironmentId.ValueString())
	resp.Diagnostics.Append(d...)
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

func (r *PolicyManagementPolicyRootResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state policyManagementPolicyRootResourceModel

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
	getResponse, d := r.getRootPolicy(ctx, plan.EnvironmentId.ValueString())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	version := getResponse.GetVersion()

	// Build the model for the API
	policyManagementPolicyRoot, d := plan.expandUpdate(ctx, version)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyManagementPolicyRoot.SetId(getResponse.GetId())

	// Run the API call
	var response *authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorPoliciesApi.UpdatePolicy(ctx, plan.EnvironmentId.ValueString(), getResponse.GetId()).AuthorizeEditorDataPoliciesReferenceablePolicyDTO(*policyManagementPolicyRoot).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdatePolicy-Update",
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

func (r *PolicyManagementPolicyRootResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *policyManagementPolicyRootResourceModel

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
	getResponse, d := r.getRootPolicy(ctx, data.EnvironmentId.ValueString())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyManagementPolicyRoot := authorize.NewAuthorizeEditorDataPoliciesReferenceablePolicyDTO(
		getResponse.GetId(),
		"Policies",
		*authorize.NewAuthorizeEditorDataPoliciesCombiningAlgorithmDTO(
			authorize.ENUMAUTHORIZEEDITORDATAPOLICIESCOMBININGALGORITHMDTOALGORITHM_PERMIT_OVERRIDES,
		),
		getResponse.GetVersion(),
	)

	// Run the API call
	var response *authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorPoliciesApi.UpdatePolicy(ctx, data.EnvironmentId.ValueString(), getResponse.GetId()).AuthorizeEditorDataPoliciesReferenceablePolicyDTO(*policyManagementPolicyRoot).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdatePolicy-Delete",
		framework.DefaultCustomError,
		nil,
		&response,
	)...)

}

func (r *PolicyManagementPolicyRootResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
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

func (p *policyManagementPolicyRootResourceModel) expandCreate(ctx context.Context) (*authorize.AuthorizeEditorDataPoliciesPolicyDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	var plan *policyManagementPolicyCombiningAlgorithmResourceModel
	diags.Append(p.CombiningAlgorithm.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	combiningAlgorithm := plan.expand()

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
	// 	var plan []policyManagementPolicyRootStatementResourceModel
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

	if !p.Children.IsNull() && !p.Children.IsUnknown() {
		children, d := expandEditorDataPolicyChildren(ctx, p.Children)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetChildren(children)
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

func (p *policyManagementPolicyRootResourceModel) expandUpdate(ctx context.Context, versionId string) (*authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO, diag.Diagnostics) {
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

// func (p *policyManagementPolicyRootStatementResourceModel) expand() map[string]interface{} {

// 	log.Panicf("Not implemented")

// 	return nil
// }

func (p *policyManagementPolicyRootResourceModel) toState(ctx context.Context, apiObject *authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO) diag.Diagnostics {
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

	childrenPolicies, ok := apiObject.GetChildrenOk()
	p.Children, d = editorDataPolicysOkToListTF(ctx, childrenPolicies, ok)
	diags.Append(d...)

	// p.ManagedEntity, d = editorManagedEntityOkToTF(apiObject.GetManagedEntityOk())
	// diags.Append(d...)

	p.Version = framework.StringOkToTF(apiObject.GetVersionOk())

	return diags
}

func (r *PolicyManagementPolicyRootResource) getRootPolicy(ctx context.Context, environmentId string) (*authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	var response *authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := r.Client.AuthorizeAPIClient.AuthorizeEditorPoliciesApi.ListRootPolicies(ctx, environmentId).Execute()

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environmentId, nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if policies, ok := pageCursor.EntityArray.Embedded.GetAuthorizationPoliciesOk(); ok {

					var policyObj authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO
					for _, policyObj = range policies {
						if v, ok := policyObj.GetManagedEntityOk(); ok {
							if v1, ok := v.GetOwnerOk(); ok {
								if v2, ok := v1.GetServiceOk(); ok && v2.GetName() == "Editor Service" {
									return &policyObj, pageCursor.HTTPResponse, nil
								}
							}
						}
					}
				}
			}

			return nil, initialHttpResponse, nil
		},
		"ListRootPolicies",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	return response, diags
}
