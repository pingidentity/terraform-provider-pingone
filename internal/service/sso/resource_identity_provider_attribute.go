package sso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"golang.org/x/exp/slices"
)

// Types
type IdentityProviderAttributeResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type IdentityProviderAttributeResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	EnvironmentId      types.String `tfsdk:"environment_id"`
	IdentityProviderId types.String `tfsdk:"identity_provider_id"`
	Name               types.String `tfsdk:"name"`
	Update             types.String `tfsdk:"update"`
	Value              types.String `tfsdk:"value"`
	MappingType        types.String `tfsdk:"mapping_type"`
}

// Framework interfaces
var (
	_ resource.Resource                = &IdentityProviderAttributeResource{}
	_ resource.ResourceWithConfigure   = &IdentityProviderAttributeResource{}
	_ resource.ResourceWithImportState = &IdentityProviderAttributeResource{}
)

// New Object
func NewIdentityProviderAttributeResource() resource.Resource {
	return &IdentityProviderAttributeResource{}
}

// Metadata
func (r *IdentityProviderAttributeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_provider_attribute"
}

// Schema.
func (r *IdentityProviderAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	reservedNames := []string{"account", "id", "created", "updated", "lifecycle", "mfaEnabled", "enabled"}
	nameDescriptionFmt := fmt.Sprintf("The user attribute, which is unique per provider. The attribute must not be defined as read only from the user schema or of type `COMPLEX` based on the user schema. Valid examples `username`, and `name.first`. The following attributes may not be used: `%s`.", strings.Join(reservedNames, "`, `"))
	nameDescription := framework.SchemaDescription{
		MarkdownDescription: nameDescriptionFmt,
		Description:         strings.ReplaceAll(nameDescriptionFmt, "`", "\""),
	}

	updateDescriptionFmt := fmt.Sprintf("Indicates whether to update the user attribute in the directory with the non-empty mapped value from the IdP. Options are `%s` (only update the user attribute if it has an empty value); `%s` (always update the user attribute value).", string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_EMPTY_ONLY), string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_ALWAYS))
	updateDescription := framework.SchemaDescription{
		MarkdownDescription: updateDescriptionFmt,
		Description:         strings.ReplaceAll(updateDescriptionFmt, "`", "\""),
	}

	valueDescriptionFmt := "A placeholder referring to the attribute (or attributes) from the provider. Placeholders must be valid for the attributes returned by the IdP type and use the `${}` syntax (for example, `${email}`). For SAML, any placeholder is acceptable, and it is mapped against the attributes available in the SAML assertion after authentication. The `${samlAssertion.subject}` placeholder is a special reserved placeholder used to refer to the subject name ID in the SAML assertion response."
	valueDescription := framework.SchemaDescription{
		MarkdownDescription: valueDescriptionFmt,
		Description:         strings.ReplaceAll(valueDescriptionFmt, "`", "\""),
	}

	mappingTypeDescriptionFmt := fmt.Sprintf("The mapping type. Options are `%s` (This attribute is required by the schema and cannot be removed. The name and update properties cannot be changed.) or `%s` (All user-created attributes are of this type.)", string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGTYPE_CORE), string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGTYPE_CUSTOM))
	mappingTypeDescription := framework.SchemaDescription{
		MarkdownDescription: mappingTypeDescriptionFmt,
		Description:         strings.ReplaceAll(mappingTypeDescriptionFmt, "`", "\""),
	}

	coreValueNames := []string{
		"username",
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an attribute mapping for identity providers configured in PingOne.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment to create the identity provider attribute in."},
			),

			"identity_provider_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the identity provider to create the attribute mapping for."},
			),

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(
						func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
							isCoreValueName := struct {
								state bool
								plan  bool
							}{
								state: slices.Contains(coreValueNames, req.StateValue.ValueString()),
								plan:  slices.Contains(coreValueNames, req.PlanValue.ValueString()),
							}

							resp.RequiresReplace = isCoreValueName.state != isCoreValueName.plan
						},
						"The resource must be replaced if changing between core and custom attribute types.",
						"The resource must be replaced if changing between core and custom attribute types.",
					),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.NoneOf(reservedNames...),
				},
			},

			"update": schema.StringAttribute{
				Description:         updateDescription.Description,
				MarkdownDescription: updateDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_EMPTY_ONLY), string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_ALWAYS)),
				},
			},

			"value": schema.StringAttribute{
				Description:         valueDescription.Description,
				MarkdownDescription: valueDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"mapping_type": schema.StringAttribute{
				Description:         mappingTypeDescription.Description,
				MarkdownDescription: mappingTypeDescription.MarkdownDescription,
				Computed:            true,
			},
		},
	}
}

func (r *IdentityProviderAttributeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IdentityProviderAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state IdentityProviderAttributeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	createIdentityProviderAttribute := plan.expand()

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.IdentityProviderAttributesApi.CreateIdentityProviderAttribute(ctx, plan.EnvironmentId.ValueString(), plan.IdentityProviderId.ValueString()).IdentityProviderAttribute(*createIdentityProviderAttribute).Execute()
		},
		"CreateIdentityProviderAttribute",
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
	resp.Diagnostics.Append(state.toState(response.(*management.IdentityProviderAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IdentityProviderAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *IdentityProviderAttributeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
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
			return r.client.IdentityProviderAttributesApi.ReadOneIdentityProviderAttribute(ctx, data.EnvironmentId.ValueString(), data.IdentityProviderId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneIdentityProviderAttribute",
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
	resp.Diagnostics.Append(data.toState(response.(*management.IdentityProviderAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IdentityProviderAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state IdentityProviderAttributeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	identityProviderAttributeMapping := plan.expand()

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.IdentityProviderAttributesApi.UpdateIdentityProviderAttribute(ctx, plan.EnvironmentId.ValueString(), plan.IdentityProviderId.ValueString(), plan.Id.ValueString()).IdentityProviderAttribute(*identityProviderAttributeMapping).Execute()
		},
		"UpdateIdentityProviderAttribute",
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
	resp.Diagnostics.Append(state.toState(response.(*management.IdentityProviderAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IdentityProviderAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *IdentityProviderAttributeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
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
			r, err := r.client.IdentityProviderAttributesApi.DeleteIdentityProviderAttribute(ctx, data.EnvironmentId.ValueString(), data.IdentityProviderId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteIdentityProviderAttribute",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *IdentityProviderAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/identity_provider_id/identity_provider_attribute_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identity_provider_id"), attributes[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *IdentityProviderAttributeResourceModel) expand() *management.IdentityProviderAttribute {

	data := management.NewIdentityProviderAttribute(p.Name.ValueString(), p.Value.ValueString(), management.EnumIdentityProviderAttributeMappingUpdate(p.Update.ValueString()))

	return data
}

func (p *IdentityProviderAttributeResourceModel) toState(apiObject *management.IdentityProviderAttribute) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Update = IdentityProviderAttributeMappingUpdateOkToTF(apiObject.GetUpdateOk())
	p.Value = framework.StringOkToTF(apiObject.GetValueOk())
	p.MappingType = IdentityProviderAttributeMappingTypeOkToTF(apiObject.GetMappingTypeOk())

	return diags
}

func IdentityProviderAttributeMappingTypeOkToTF(v *management.EnumIdentityProviderAttributeMappingType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}

func IdentityProviderAttributeMappingUpdateOkToTF(v *management.EnumIdentityProviderAttributeMappingUpdate, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
