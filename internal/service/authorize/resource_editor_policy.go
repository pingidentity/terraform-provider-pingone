package authorize

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type EditorPolicyResource serviceClientType

type editorPolicyResourceModel struct {
	Id                 pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId      pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Type               types.String                 `tfsdk:"type"`
	Name               types.String                 `tfsdk:"name"`
	Description        types.String                 `tfsdk:"description"`
	Enabled            types.Bool                   `tfsdk:"enabled"`
	Statements         types.List                   `tfsdk:"statements"`
	Condition          types.Object                 `tfsdk:"condition"`
	CombiningAlgorithm types.Object                 `tfsdk:"combining_algorithm"`
	Children           types.List                   `tfsdk:"children"`
	RepetitionSettings types.Object                 `tfsdk:"repetition_settings"`
	ManagedEntity      types.Object                 `tfsdk:"managed_entity"`
	Version            types.String                 `tfsdk:"version"`
}

type editorPolicyStatementResourceModel struct{}

type editorPolicyConditionResourceModel struct {
	Type       types.String `tfsdk:"type"`
	Conditions types.List   `tfsdk:"conditions"`
	Left       types.Object `tfsdk:"left"`
	Comparator types.String `tfsdk:"comparator"`
	Right      types.Object `tfsdk:"right"`
	Condition  types.Object `tfsdk:"condition"`
	Reference  types.Object `tfsdk:"reference"`
}

type editorPolicyConditionComprandResourceModel struct {
	Type types.String `tfsdk:"type"`
}

type editorPolicyConditionReferenceResourceModel struct {
	Id types.String `tfsdk:"id"`
}

type editorPolicyCombiningAlgorithmResourceModel struct {
	Algorithm types.String `tfsdk:"algorithm"`
}

type editorPolicyChildrenResourceModel struct{}

type editorPolicyRepetitionSettingsResourceModel struct {
	Source   types.Object `tfsdk:"source"`
	Decision types.String `tfsdk:"decision"`
}

type editorPolicyRepetitionSettingsSourceResourceModel struct {
	Id types.String `tfsdk:"id"`
}

type editorPolicyManagedEntityResourceModel struct {
	Owner        types.Object `tfsdk:"owner"`
	Restrictions types.Object `tfsdk:"restrictions"`
	Reference    types.Object `tfsdk:"reference"`
}

type editorPolicyManagedEntityOwnerResourceModel struct {
	Service types.Object `tfsdk:"service"`
}

type editorPolicyManagedEntityOwnerServiceResourceModel struct {
	Name types.String `tfsdk:"name"`
}

type editorPolicyManagedEntityRestrictionsResourceModel struct {
	ReadOnly         types.Bool `tfsdk:"read_only"`
	DisallowChildren types.Bool `tfsdk:"disallow_children"`
}

type editorPolicyManagedEntityReferenceResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Type       types.String `tfsdk:"type"`
	Name       types.String `tfsdk:"name"`
	UiDeepLink types.String `tfsdk:"ui_deep_link"`
}

// Framework interfaces
var (
	_ resource.Resource                = &EditorPolicyResource{}
	_ resource.ResourceWithConfigure   = &EditorPolicyResource{}
	_ resource.ResourceWithImportState = &EditorPolicyResource{}
)

// New Object
func NewEditorPolicyResource() resource.Resource {
	return &EditorPolicyResource{}
}

// Metadata
func (r *EditorPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_editor_policy"
}

func (r *EditorPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage Authorize editor policies in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the Authorize editor policy in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"type": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},

			"enabled": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},

			"statements": schema.ListNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{},
				},
			},

			"condition": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{},
			},

			"combining_algorithm": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{},
			},

			"children": schema.ListNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{},
				},
			},

			"repetition_settings": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{},
			},

			"managed_entity": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{},
			},

			"version": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},
		},
	}
}

func (r *EditorPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EditorPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state editorPolicyResourceModel

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
	editorPolicy, d := plan.expandCreate(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataPoliciesPolicyDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorPoliciesApi.CreatePolicy(ctx, plan.EnvironmentId.ValueString()).AuthorizeEditorDataPoliciesPolicyDTO(*editorPolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreatePolicy",
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

func (r *EditorPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *editorPolicyResourceModel

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
	var response *authorize.AuthorizeEditorDataPoliciesPolicyDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorPoliciesApi.GetPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetPolicy",
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

func (r *EditorPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state editorPolicyResourceModel

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
	editorPolicy, d := plan.expandUpdate(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataPoliciesPolicyDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorPoliciesApi.UpdatePolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).AuthorizeEditorDataPoliciesReferenceablePolicyDTO(*editorPolicy).Execute()
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
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *EditorPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *editorPolicyResourceModel

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
		return
	}
}

func (r *EditorPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "authorize_editor_policy_id",
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

func (p *editorPolicyResourceModel) expandCreate(ctx context.Context) (*authorize.AuthorizeEditorDataPoliciesPolicyDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	combiningAlgorithm := &authorize.AuthorizeEditorDataPoliciesCombiningAlgorithmDTO{}

	log.Fatalf("Not implemented")

	data := authorize.NewAuthorizeEditorDataPoliciesPolicyDTO(
		p.Name.ValueString(),
		*combiningAlgorithm,
	)

	return data, diags
}

func (p *editorPolicyResourceModel) expandUpdate(ctx context.Context) (*authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	combiningAlgorithm := &authorize.AuthorizeEditorDataPoliciesCombiningAlgorithmDTO{}

	log.Fatalf("Not implemented")

	data := authorize.NewAuthorizeEditorDataPoliciesReferenceablePolicyDTO(
		p.Id.ValueString(),
		p.Name.ValueString(),
		*combiningAlgorithm,
		p.Version.ValueString(),
	)

	return data, diags
}

func (p *editorPolicyResourceModel) toState(apiObject *authorize.AuthorizeEditorDataPoliciesPolicyDTO) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	return diags
}
