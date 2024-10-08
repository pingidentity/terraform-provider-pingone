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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type EditorStatementResource serviceClientType

type editorStatementResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name          types.String                 `tfsdk:"name"`
	Description   types.String                 `tfsdk:"description"`
	Code          types.String                 `tfsdk:"code"`
	AppliesTo     types.String                 `tfsdk:"applies_to"`
	AppliesIf     types.String                 `tfsdk:"applies_if"`
	Payload       types.String                 `tfsdk:"payload"`
	Obligatory    types.Bool                   `tfsdk:"obligatory"`
	Attributes    types.Set                    `tfsdk:"attributes"`
	Version       types.String                 `tfsdk:"version"`
}

// Framework interfaces
var (
	_ resource.Resource                = &EditorStatementResource{}
	_ resource.ResourceWithConfigure   = &EditorStatementResource{}
	_ resource.ResourceWithImportState = &EditorStatementResource{}
)

// New Object
func NewEditorStatementResource() resource.Resource {
	return &EditorStatementResource{}
}

// Metadata
func (r *EditorStatementResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_editor_statement"
}

func (r *EditorStatementResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an authorization statement for the PingOne Authorize Policy Manager in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the Authorize editor statement in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},

			"code": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"applies_to": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataStatementsReferenceableStatementDTOAppliesToEnumValues)...),
				},
			},

			"applies_if": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataStatementsReferenceableStatementDTOAppliesToEnumValues)...),
				},
			},

			"payload": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"obligatory": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},

			"attributes": schema.ListNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Required:    true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the authorization attribute in the trust framework.")),
				},
			},

			"version": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Computed:    true,
			},
		},
	}
}

func (r *EditorStatementResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EditorStatementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state editorStatementResourceModel

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
	editorStatement, d := plan.expandCreate(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataStatementsReferenceableStatementDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorStatementsApi.CreateStatement(ctx, plan.EnvironmentId.ValueString()).AuthorizeEditorDataStatementsStatementDTO(*editorStatement).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateStatement",
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

func (r *EditorStatementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *editorStatementResourceModel

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
	var response *authorize.AuthorizeEditorDataStatementsReferenceableStatementDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorStatementsApi.GetStatement(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetStatement",
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

func (r *EditorStatementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state editorStatementResourceModel

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
	editorStatement, d := plan.expandUpdate(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataStatementsReferenceableStatementDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorStatementsApi.UpdateStatement(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).AuthorizeEditorDataStatementsReferenceableStatementDTO(*editorStatement).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateStatement",
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

func (r *EditorStatementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *editorStatementResourceModel

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
			fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorStatementsApi.DeleteStatement(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteStatement",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EditorStatementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "authorize_editor_statement_id",
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

func (p *editorStatementResourceModel) expandCreate(ctx context.Context) (*authorize.AuthorizeEditorDataStatementsStatementDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := make([]authorize.AuthorizeEditorDataReferenceObjectDTO, 0)
	var valueAttributesPlan []editorReferenceDataResourceModel
	diags.Append(p.Attributes.ElementsAs(ctx, &valueAttributesPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	for _, attributePlan := range valueAttributesPlan {
		attributes = append(attributes, *attributePlan.expand())
	}

	// Main object
	data := authorize.NewAuthorizeEditorDataStatementsStatementDTO(
		p.Name.ValueString(),
		p.Code.ValueString(),
		authorize.EnumAuthorizeEditorDataStatementsStatementDTOAppliesTo(p.AppliesTo.ValueString()),
		authorize.EnumAuthorizeEditorDataStatementsStatementDTOAppliesIf(p.AppliesIf.ValueString()),
		p.Payload.ValueString(),
		attributes,
	)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.Obligatory.IsNull() && !p.Obligatory.IsUnknown() {
		data.SetObligatory(p.Obligatory.ValueBool())
	}

	return data, diags
}

func (p *editorStatementResourceModel) expandUpdate(ctx context.Context) (*authorize.AuthorizeEditorDataStatementsReferenceableStatementDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	dataCreate, d := p.expandCreate(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Use json.marshall and unmarshal to cast dataCreate to a AuthorizeEditorDataStatementsReferenceableStatementDTO type
	bytes, err := json.Marshal(dataCreate)
	if err != nil {
		diags.AddError("Failed to marshal data", err.Error())
		return nil, diags
	}

	var data *authorize.AuthorizeEditorDataStatementsReferenceableStatementDTO
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

func (p *editorStatementResourceModel) toState(apiObject *authorize.AuthorizeEditorDataStatementsReferenceableStatementDTO) diag.Diagnostics {
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
	p.Code = framework.StringOkToTF(apiObject.GetCodeOk())
	p.AppliesTo = framework.EnumOkToTF(apiObject.GetAppliesToOk())
	p.AppliesIf = framework.EnumOkToTF(apiObject.GetAppliesIfOk())
	p.Payload = framework.StringOkToTF(apiObject.GetPayloadOk())
	p.Obligatory = framework.BoolOkToTF(apiObject.GetObligatoryOk())
	p.Version = framework.StringOkToTF(apiObject.GetVersionOk())

	p.Attributes, d = editorDataReferenceObjectOkToSetTF(apiObject.GetAttributesOk())
	diags = append(diags, d...)

	return diags
}
