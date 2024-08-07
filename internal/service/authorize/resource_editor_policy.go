package authorize

import (
	"context"
	"encoding/json"
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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	listvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/listvalidator"
	objectvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
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

type editorPolicyConditionConditionResourceModel struct {
	Type types.String `tfsdk:"type"`
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

const (
	policyConditionTypeAndConditionValue        = ""
	policyConditionTypeComparisonConditionValue = ""
	policyConditionTypeEmptyConditionValue      = ""
	policyConditionTypeNotConditionValue        = ""
	policyConditionTypeOrConditionValue         = ""
	policyConditionTypeReferenceConditionValue  = ""
)

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

				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Required:    true,
					},

					"conditions": schema.ListNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Optional:    true,

						Validators: []validator.List{
							listvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(policyConditionTypeAndConditionValue),
								path.MatchRoot("type"),
							),
							listvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeComparisonConditionValue),
								path.MatchRoot("type"),
							),
							listvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeEmptyConditionValue),
								path.MatchRoot("type"),
							),
							listvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeNotConditionValue),
								path.MatchRoot("type"),
							),
							listvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeOrConditionValue),
								path.MatchRoot("type"),
							),
							listvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeReferenceConditionValue),
								path.MatchRoot("type"),
							),
						},

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
									Required:    true,
								},
							},
						},
					},

					"left": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Optional:    true,

						Validators: []validator.Object{
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeAndConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(policyConditionTypeComparisonConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeEmptyConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeNotConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeOrConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeReferenceConditionValue),
								path.MatchRoot("type"),
							),
						},

						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Required:    true,
							},
						},
					},

					"comparator": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Optional:    true,

						Validators: []validator.String{
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeAndConditionValue),
								path.MatchRoot("type"),
							),
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(policyConditionTypeComparisonConditionValue),
								path.MatchRoot("type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeEmptyConditionValue),
								path.MatchRoot("type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeNotConditionValue),
								path.MatchRoot("type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeOrConditionValue),
								path.MatchRoot("type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeReferenceConditionValue),
								path.MatchRoot("type"),
							),
						},
					},

					"right": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Optional:    true,

						Validators: []validator.Object{
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeAndConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(policyConditionTypeComparisonConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeEmptyConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeNotConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeOrConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeReferenceConditionValue),
								path.MatchRoot("type"),
							),
						},

						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Required:    true,
							},
						},
					},

					"condition": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Optional:    true,

						Validators: []validator.Object{
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeAndConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeComparisonConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeEmptyConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(policyConditionTypeNotConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeOrConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeReferenceConditionValue),
								path.MatchRoot("type"),
							),
						},

						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Required:    true,
							},
						},
					},

					"reference": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Optional:    true,

						Validators: []validator.Object{
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeAndConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeComparisonConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeEmptyConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeNotConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(policyConditionTypeOrConditionValue),
								path.MatchRoot("type"),
							),
							objectvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(policyConditionTypeReferenceConditionValue),
								path.MatchRoot("type"),
							),
						},

						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Required:    true,
							},
						},
					},
				},
			},

			"combining_algorithm": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"algorithm": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Required:    true,
					},
				},
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

				Attributes: map[string]schema.Attribute{
					"source": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Required:    true,

						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Required:    true,
							},
						},
					},

					"decision": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Required:    true,
					},
				},
			},

			"managed_entity": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"owner": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Required:    true,

						Attributes: map[string]schema.Attribute{
							"service": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Required:    true,

								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
										Required:    true,
									},
								},
							},
						},
					},

					"restrictions": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"read_only": schema.BoolAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Optional:    true,
							},

							"disallow_children": schema.BoolAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Optional:    true,
							},
						},
					},

					"reference": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Optional:    true,
							},

							"type": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Optional:    true,
							},

							"name": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Optional:    true,
							},

							"ui_deep_link": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Optional:    true,
							},
						},
					},
				},
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

	data := authorize.NewAuthorizeEditorDataPoliciesPolicyDTO(
		p.Name.ValueString(),
		*combiningAlgorithm,
	)

	if !p.Type.IsNull() && !p.Type.IsUnknown() {
		data.SetType(p.Type.ValueString())
	}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.Enabled.IsNull() && !p.Enabled.IsUnknown() {
		data.SetEnabled(p.Enabled.ValueBool())
	}

	if !p.Statements.IsNull() && !p.Statements.IsUnknown() {
		var plan []editorPolicyStatementResourceModel
		diags.Append(p.Statements.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		statements := make([]map[string]interface{}, 0)
		for _, planItem := range plan {
			statements = append(statements, planItem.expand())
		}

		data.SetStatements(statements)
	}

	if !p.Condition.IsNull() && !p.Condition.IsUnknown() {
		var plan *editorPolicyConditionResourceModel
		diags.Append(p.Condition.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		condition, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetCondition(*condition)
	}

	if !p.CombiningAlgorithm.IsNull() && !p.CombiningAlgorithm.IsUnknown() {
		var plan *editorPolicyCombiningAlgorithmResourceModel
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

	if !p.Children.IsNull() && !p.Children.IsUnknown() {
		var plan []editorPolicyChildrenResourceModel
		diags.Append(p.Children.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		children := make([]map[string]interface{}, 0)
		for _, planItem := range plan {
			children = append(children, planItem.expand())
		}

		data.SetChildren(children)
	}

	if !p.RepetitionSettings.IsNull() && !p.RepetitionSettings.IsUnknown() {
		var plan *editorPolicyRepetitionSettingsResourceModel
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

	if !p.ManagedEntity.IsNull() && !p.ManagedEntity.IsUnknown() {
		var plan *editorPolicyManagedEntityResourceModel
		diags.Append(p.ManagedEntity.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		managedEntity, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetManagedEntity(*managedEntity)
	}

	return data, diags
}

func (p *editorPolicyResourceModel) expandUpdate(ctx context.Context) (*authorize.AuthorizeEditorDataPoliciesReferenceablePolicyDTO, diag.Diagnostics) {
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

	if !p.Version.IsNull() && !p.Version.IsUnknown() {
		data.SetVersion(p.Version.ValueString())
	}

	return data, diags
}

func (p *editorPolicyStatementResourceModel) expand() map[string]interface{} {

	log.Panicf("Not implemented")

	return nil
}

func (p *editorPolicyConditionResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataRulesRuleDTOCondition, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	andCondition, d := p.expandAndCondition(ctx)
	diags.Append(d...)

	comparisonCondition, d := p.expandComparisonCondition(ctx)
	diags.Append(d...)

	notCondition, d := p.expandNotCondition(ctx)
	diags.Append(d...)

	referenceCondition, d := p.expandReferenceCondition(ctx)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	data := authorize.AuthorizeEditorDataRulesRuleDTOCondition{
		AuthorizeEditorDataConditionsAndConditionDTO:        andCondition,
		AuthorizeEditorDataConditionsComparisonConditionDTO: comparisonCondition,
		AuthorizeEditorDataConditionsEmptyConditionDTO:      p.expandEmptyCondition(),
		AuthorizeEditorDataConditionsNotConditionDTO:        notCondition,
		AuthorizeEditorDataConditionsOrConditionDTO:         p.expandOrCondition(),
		AuthorizeEditorDataConditionsReferenceConditionDTO:  referenceCondition,
	}

	return &data, diags
}

func (p *editorPolicyConditionResourceModel) expandAndCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsAndConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.Type.ValueString() != policyConditionTypeAndConditionValue {
		return nil, diags
	}

	var plan []editorPolicyConditionConditionResourceModel
	diags.Append(p.Conditions.ElementsAs(ctx, &plan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	conditions := make([]authorize.AuthorizeEditorDataConditionDTO, 0)
	for _, planItem := range plan {
		conditions = append(conditions, *planItem.expand())
	}

	data := authorize.NewAuthorizeEditorDataConditionsAndConditionDTO(
		conditions,
		p.Type.ValueString(),
	)

	return data, diags
}

func (p *editorPolicyConditionConditionResourceModel) expand() *authorize.AuthorizeEditorDataConditionDTO {

	data := authorize.NewAuthorizeEditorDataConditionDTO(
		p.Type.ValueString(),
	)

	return data
}

func (p *editorPolicyConditionResourceModel) expandComparisonCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsComparisonConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.Type.ValueString() != policyConditionTypeComparisonConditionValue {
		return nil, diags
	}

	var leftPlan, rightPlan *editorPolicyConditionComprandResourceModel

	diags.Append(p.Left.As(ctx, &leftPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	left := leftPlan.expand()

	diags.Append(p.Right.As(ctx, &rightPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	right := rightPlan.expand()

	data := authorize.NewAuthorizeEditorDataConditionsComparisonConditionDTO(
		*left,
		p.Comparator.ValueString(),
		*right,
		p.Type.ValueString(),
	)

	return data, nil
}

func (p *editorPolicyConditionComprandResourceModel) expand() *authorize.AuthorizeEditorDataConditionsComparandDTO {

	data := authorize.NewAuthorizeEditorDataConditionsComparandDTO(
		p.Type.ValueString(),
	)

	return data
}

func (p *editorPolicyConditionResourceModel) expandEmptyCondition() *authorize.AuthorizeEditorDataConditionsEmptyConditionDTO {

	if p.Type.ValueString() != policyConditionTypeEmptyConditionValue {
		return nil
	}

	data := authorize.NewAuthorizeEditorDataConditionsEmptyConditionDTO(
		p.Type.ValueString(),
	)

	return data
}

func (p *editorPolicyConditionResourceModel) expandNotCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsNotConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.Type.ValueString() != policyConditionTypeNotConditionValue {
		return nil, diags
	}

	var plan *editorPolicyConditionConditionResourceModel
	diags.Append(p.Condition.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	condition := plan.expand()

	data := authorize.NewAuthorizeEditorDataConditionsNotConditionDTO(
		*condition,
		p.Type.ValueString(),
	)

	return data, diags
}

func (p *editorPolicyConditionResourceModel) expandOrCondition() *authorize.AuthorizeEditorDataConditionsOrConditionDTO {

	if p.Type.ValueString() != policyConditionTypeOrConditionValue {
		return nil
	}

	data := authorize.NewAuthorizeEditorDataConditionsOrConditionDTO(
		p.Type.ValueString(),
	)

	return data
}

func (p *editorPolicyConditionResourceModel) expandReferenceCondition(ctx context.Context) (*authorize.AuthorizeEditorDataConditionsReferenceConditionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.Type.ValueString() != policyConditionTypeReferenceConditionValue {
		return nil, diags
	}

	var plan *editorPolicyConditionReferenceResourceModel
	diags.Append(p.Reference.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	reference := plan.expand()

	data := authorize.NewAuthorizeEditorDataConditionsReferenceConditionDTO(
		*reference,
		p.Type.ValueString(),
	)

	return data, diags
}

func (p *editorPolicyConditionReferenceResourceModel) expand() *authorize.AuthorizeEditorDataReferenceObjectDTO {

	data := authorize.NewAuthorizeEditorDataReferenceObjectDTO(
		p.Id.ValueString(),
	)

	return data
}

func (p *editorPolicyCombiningAlgorithmResourceModel) expand() *authorize.AuthorizeEditorDataPoliciesCombiningAlgorithmDTO {

	data := authorize.NewAuthorizeEditorDataPoliciesCombiningAlgorithmDTO(
		p.Algorithm.ValueString(),
	)

	return data
}

func (p *editorPolicyChildrenResourceModel) expand() map[string]interface{} {

	log.Panicf("Not implemented")

	return nil
}

func (p *editorPolicyRepetitionSettingsResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataPoliciesRepetitionSettingsDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	var plan *editorPolicyRepetitionSettingsSourceResourceModel
	diags.Append(p.Source.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	source := plan.expand()

	data := authorize.NewAuthorizeEditorDataPoliciesRepetitionSettingsDTO(
		*source,
		p.Decision.ValueString(),
	)

	return data, diags
}

func (p *editorPolicyRepetitionSettingsSourceResourceModel) expand() *authorize.AuthorizeEditorDataReferenceObjectDTO {

	data := authorize.NewAuthorizeEditorDataReferenceObjectDTO(
		p.Id.ValueString(),
	)

	return data
}

func (p *editorPolicyManagedEntityResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataManagedEntityDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	var plan *editorPolicyManagedEntityOwnerResourceModel
	diags.Append(p.Owner.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	owner, d := plan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataManagedEntityDTO(
		*owner,
	)

	if !p.Restrictions.IsNull() && !p.Restrictions.IsUnknown() {
		var plan *editorPolicyManagedEntityRestrictionsResourceModel
		diags.Append(p.Restrictions.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		restrictions := plan.expand()

		data.SetRestrictions(*restrictions)
	}

	if !p.Reference.IsNull() && !p.Reference.IsUnknown() {
		var plan *editorPolicyManagedEntityReferenceResourceModel
		diags.Append(p.Reference.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		reference := plan.expand()

		data.SetReference(*reference)
	}

	return data, diags
}

func (p *editorPolicyManagedEntityOwnerResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataManagedEntityOwnerDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	var plan *editorPolicyManagedEntityOwnerServiceResourceModel
	diags.Append(p.Service.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	service := plan.expand()

	data := authorize.NewAuthorizeEditorDataManagedEntityOwnerDTO(
		*service,
	)

	return data, diags
}

func (p *editorPolicyManagedEntityOwnerServiceResourceModel) expand() *authorize.AuthorizeEditorDataServiceObjectDTO {

	data := authorize.NewAuthorizeEditorDataServiceObjectDTO(
		p.Name.ValueString(),
	)

	return data
}

func (p *editorPolicyManagedEntityRestrictionsResourceModel) expand() *authorize.AuthorizeEditorDataManagedEntityRestrictionsDTO {

	data := authorize.NewAuthorizeEditorDataManagedEntityRestrictionsDTO()

	if !p.ReadOnly.IsNull() && !p.ReadOnly.IsUnknown() {
		data.SetReadOnly(p.ReadOnly.ValueBool())
	}

	if !p.DisallowChildren.IsNull() && !p.DisallowChildren.IsUnknown() {
		data.SetDisallowChildren(p.DisallowChildren.ValueBool())
	}

	return data
}

func (p *editorPolicyManagedEntityReferenceResourceModel) expand() *authorize.AuthorizeEditorDataManagedEntityManagedEntityReferenceDTO {

	data := authorize.NewAuthorizeEditorDataManagedEntityManagedEntityReferenceDTO()

	if !p.Id.IsNull() && !p.Id.IsUnknown() {
		data.SetId(p.Id.ValueString())
	}

	if !p.Type.IsNull() && !p.Type.IsUnknown() {
		data.SetType(p.Type.ValueString())
	}

	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	if !p.UiDeepLink.IsNull() && !p.UiDeepLink.IsUnknown() {
		data.SetUiDeepLink(p.UiDeepLink.ValueString())
	}

	return data
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
