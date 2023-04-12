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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type ApplicationAttributeMappingResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type ApplicationAttributeMappingResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	ApplicationId types.String `tfsdk:"application_id"`
	Name          types.String `tfsdk:"name"`
	Required      types.Bool   `tfsdk:"required"`
	Value         types.String `tfsdk:"value"`
	MappingType   types.String `tfsdk:"mapping_type"`
}

// Framework interfaces
var (
	_ resource.Resource                = &ApplicationAttributeMappingResource{}
	_ resource.ResourceWithConfigure   = &ApplicationAttributeMappingResource{}
	_ resource.ResourceWithImportState = &ApplicationAttributeMappingResource{}
)

// New Object
func NewApplicationAttributeMappingResource() resource.Resource {
	return &ApplicationAttributeMappingResource{}
}

// Metadata
func (r *ApplicationAttributeMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_attribute_mapping"
}

// Schema.
func (r *ApplicationAttributeMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	reservedNames := []string{"acr", "amr", "at_hash", "aud", "auth_time", "azp", "client_id", "exp", "iat", "iss", "jti", "nbf", "nonce", "org", "scope", "sid", "sub"}
	nameDescriptionFmt := fmt.Sprintf("A string that specifies the name of attribute and must be unique within an application. For SAML applications, the `samlAssertion.subject` name is a reserved case-insensitive name which indicates the mapping to be used for the subject in an assertion. For OpenID Connect applications, the following names are reserved and cannot be used: `%s`.", strings.Join(reservedNames, "`, `"))
	nameDescription := framework.SchemaDescription{
		MarkdownDescription: nameDescriptionFmt,
		Description:         strings.ReplaceAll(nameDescriptionFmt, "`", "\""),
	}

	requiredDescriptionFmt := "A boolean to specify whether a mapping value is required for this attribute. If `true`, a value must be set and a non-empty value must be available in the SAML assertion or ID token. Defaults to `false`."
	requiredDescription := framework.SchemaDescription{
		MarkdownDescription: requiredDescriptionFmt,
		Description:         strings.ReplaceAll(requiredDescriptionFmt, "`", "\""),
	}

	valueDescriptionFmt := "A string that specifies the string constants or expression for mapping the attribute path against a specific source. The expression format is `${<source>.<attribute_path>}`. The only supported source is user (for example, `${user.id}`)."
	valueDescription := framework.SchemaDescription{
		MarkdownDescription: valueDescriptionFmt,
		Description:         strings.ReplaceAll(valueDescriptionFmt, "`", "\""),
	}

	mappingTypeDescriptionFmt := "A string that specifies the mapping type of the attribute. Options are `CORE`, `SCOPE`, and `CUSTOM`."
	mappingTypeDescription := framework.SchemaDescription{
		MarkdownDescription: mappingTypeDescriptionFmt,
		Description:         strings.ReplaceAll(mappingTypeDescriptionFmt, "`", "\""),
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an attribute mapping for applications configured in PingOne.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment to create the application attribute mapping in."},
			),

			"application_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the application to create the attribute mapping for."},
			),

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"required": schema.BoolAttribute{
				Description:         requiredDescription.Description,
				MarkdownDescription: requiredDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
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

func (r *ApplicationAttributeMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationAttributeMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationAttributeMappingResourceModel

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
	createApplicationAttributeMapping := plan.expand()

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.ApplicationAttributeMappingApi.CreateApplicationAttributeMapping(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).ApplicationAttributeMapping(*createApplicationAttributeMapping).Execute()
		},
		"CreateApplicationAttributeMapping",
		framework.CustomErrorInvalidValue,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*management.ApplicationAttributeMapping))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationAttributeMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationAttributeMappingResourceModel

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
			return r.client.ApplicationAttributeMappingApi.ReadOneApplicationAttributeMapping(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneApplicationAttributeMapping",
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
	resp.Diagnostics.Append(data.toState(response.(*management.ApplicationAttributeMapping))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationAttributeMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationAttributeMappingResourceModel

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
	applicationAttributeMapping := plan.expand()

	tflog.Debug(ctx, "WUT", map[string]interface{}{
		"required-obj":   applicationAttributeMapping.GetRequired(),
		"required-plan":  plan.Required.ValueBool(),
		"required-state": state.Required.ValueBool(),
	})

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.ApplicationAttributeMappingApi.UpdateApplicationAttributeMapping(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString(), plan.Id.ValueString()).ApplicationAttributeMapping(*applicationAttributeMapping).Execute()
		},
		"UpdateApplicationAttributeMapping",
		framework.CustomErrorInvalidValue,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*management.ApplicationAttributeMapping))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationAttributeMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ApplicationAttributeMappingResourceModel

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
			r, err := r.client.ApplicationAttributeMappingApi.DeleteApplicationAttributeMapping(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteApplicationAttributeMapping",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApplicationAttributeMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/application_id/attribute_mapping_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), attributes[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *ApplicationAttributeMappingResourceModel) expand() *management.ApplicationAttributeMapping {

	data := management.NewApplicationAttributeMapping(p.Name.ValueString(), p.Required.ValueBool(), p.Value.ValueString())

	return data
}

func (p *ApplicationAttributeMappingResourceModel) toState(apiObject *management.ApplicationAttributeMapping) diag.Diagnostics {
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
	p.Required = framework.BoolOkToTF(apiObject.GetRequiredOk())
	p.Value = framework.StringOkToTF(apiObject.GetValueOk())
	p.MappingType = ApplicationAttributeMappingMappingTypeOkToTF(apiObject.GetMappingTypeOk())

	return diags
}

func ApplicationAttributeMappingMappingTypeOkToTF(v *management.EnumAttributeMappingType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
