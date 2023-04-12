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
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
	"golang.org/x/exp/slices"
)

// Types
type ResourceAttributeResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type ResourceAttributeResourceModel struct {
	Id              types.String `tfsdk:"id"`
	EnvironmentId   types.String `tfsdk:"environment_id"`
	ResourceId      types.String `tfsdk:"resource_id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Value           types.String `tfsdk:"value"`
	IDTokenEnabled  types.Bool   `tfsdk:"id_token_enabled"`
	UserinfoEnabled types.Bool   `tfsdk:"userinfo_enabled"`
}

// Framework interfaces
var (
	_ resource.Resource                = &ResourceAttributeResource{}
	_ resource.ResourceWithConfigure   = &ResourceAttributeResource{}
	_ resource.ResourceWithImportState = &ResourceAttributeResource{}
)

// New Object
func NewResourceAttributeResource() resource.Resource {
	return &ResourceAttributeResource{}
}

// Metadata
func (r *ResourceAttributeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_attribute"
}

// Schema.
func (r *ResourceAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	nameDescriptionFmt := fmt.Sprintf("A string that specifies the name of the custom resource attribute to be included in the access token. When the resource's type property is `OPENID_CONNECT`, the following are reserved names and cannot be used: %s.  When the resource's type property is `OPENID_CONNECT`, using the following names will override the default configured values, rather than creating new attributes: %s.", verify.IllegalOIDCAttributeNameString(), verify.OverrideOIDCAttributeNameString())
	nameDescription := framework.SchemaDescription{
		MarkdownDescription: nameDescriptionFmt,
		Description:         strings.ReplaceAll(nameDescriptionFmt, "`", "\""),
	}

	valueDescriptionFmt := "A string that specifies the value of the custom resource attribute. This value can be a placeholder that references an attribute in the user schema, expressed as “${user.path.to.value}”, or it can be a static string. Placeholders must be valid, enabled attributes in the environment’s user schema. Examples of valid values are: `${user.email}`, `${user.name.family}`, and `myClaimValueString`."
	valueDescription := framework.SchemaDescription{
		MarkdownDescription: valueDescriptionFmt,
		Description:         strings.ReplaceAll(valueDescriptionFmt, "`", "\""),
	}

	typeDescriptionFmt := "A string that specifies the type of resource attribute. Options are: `CORE` (The claim is required and cannot not be removed), `CUSTOM` (The claim is not a CORE attribute. All created attributes are of this type), `PREDEFINED` (A designation for predefined OIDC resource attributes such as given_name. These attributes cannot be removed; however, they can be modified)."
	typeDescription := framework.SchemaDescription{
		MarkdownDescription: typeDescriptionFmt,
		Description:         strings.ReplaceAll(typeDescriptionFmt, "`", "\""),
	}

	idTokenEnabledDescriptionFmt := "A boolean that specifies whether the attribute mapping should be available in the ID Token.  Only applies to resources that are of type `OPENID_CONNECT` and the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to false. Defaults to `true`."
	idTokenEnabledDescription := framework.SchemaDescription{
		MarkdownDescription: idTokenEnabledDescriptionFmt,
		Description:         strings.ReplaceAll(idTokenEnabledDescriptionFmt, "`", "\""),
	}

	userinfoEnabledDescriptionFmt := "A boolean that specifies whether the attribute mapping should be available through the /as/userinfo endpoint.  Only applies to resources that are of type `OPENID_CONNECT` and the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to false. Defaults to `true`."
	userinfoEnabledDescription := framework.SchemaDescription{
		MarkdownDescription: userinfoEnabledDescriptionFmt,
		Description:         strings.ReplaceAll(userinfoEnabledDescriptionFmt, "`", "\""),
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage resource attributes in PingOne.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment to create the resource attribute in."},
			),

			"resource_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the resource to assign the resource attribute to."},
			),

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
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

			"id_token_enabled": schema.BoolAttribute{
				Description:         idTokenEnabledDescription.Description,
				MarkdownDescription: idTokenEnabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
			},

			"userinfo_enabled": schema.BoolAttribute{
				Description:         userinfoEnabledDescription.Description,
				MarkdownDescription: userinfoEnabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
			},

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Computed:            true,
			},
		},
	}
}

func (r *ResourceAttributeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ResourceAttributeResourceModel

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

	attributeID, resourceType, d := validateAttributeAgainstResourceType(ctx, r.client, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString(), plan.Name.ValueString())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	resourceAttribute := plan.expand(resourceType)

	// Run the API call
	var response interface{}
	if attributeID == nil {
		response, d = framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.ResourceAttributesApi.CreateResourceAttribute(ctx, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString()).ResourceAttribute(*resourceAttribute).Execute()
			},
			"CreateResourceAttribute",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
	} else {
		response, d = framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.ResourceAttributesApi.UpdateResourceAttribute(ctx, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString(), *attributeID).ResourceAttribute(*resourceAttribute).Execute()
			},
			"UpdateResourceAttribute",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*management.ResourceAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ResourceAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ResourceAttributeResourceModel

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
			return r.client.ResourceAttributesApi.ReadOneResourceAttribute(ctx, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneResourceAttribute",
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
	resp.Diagnostics.Append(data.toState(response.(*management.ResourceAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceAttributeResourceModel

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

	_, resourceType, d := validateAttributeAgainstResourceType(ctx, r.client, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString(), plan.Name.ValueString())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	resourceMapping := plan.expand(resourceType)

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.ResourceAttributesApi.UpdateResourceAttribute(ctx, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString(), plan.Id.ValueString()).ResourceAttribute(*resourceMapping).Execute()
		},
		"UpdateResourceAttribute",
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
	resp.Diagnostics.Append(state.toState(response.(*management.ResourceAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ResourceAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ResourceAttributeResourceModel

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
			r, err := r.client.ResourceAttributesApi.DeleteResourceAttribute(ctx, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteResourceAttribute",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/resource_id/resource_attribute_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource_id"), attributes[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func validateAttributeAgainstResourceType(ctx context.Context, apiClient *management.APIClient, environmentID, resourceID, resourceAttributeName string) (*string, management.EnumResourceType, diag.Diagnostics) {
	var diags diag.Diagnostics

	respObject, d := fetchResource_Framework(ctx, apiClient, environmentID, resourceID)
	diags.Append(d...)
	if diags.HasError() {
		return nil, management.ENUMRESOURCETYPE_CUSTOM, diags
	}

	if respObject.GetType() == management.ENUMRESOURCETYPE_OPENID_CONNECT {
		if slices.Contains(verify.IllegalOIDCattributeNamesList(), resourceAttributeName) {
			diags.AddError(
				fmt.Sprintf("Invalid attribute name `%s` for the configured OpenID Connect resource.", resourceAttributeName),
				fmt.Sprintf("The attribute name provided, `%s`, cannot be used for resource ID `%s`, which is of type `OPENID_CONNECT`.", resourceAttributeName, resourceID),
			)
			return nil, respObject.GetType(), diags
		}

		if slices.Contains(verify.OverrideOIDCAttributeNameList(), resourceAttributeName) {

			resourceAttribute, d := fetchResourceAttributeFromName_Framework(ctx, apiClient, environmentID, resourceID, resourceAttributeName)
			diags.Append(d...)
			if diags.HasError() {
				return nil, respObject.GetType(), diags
			}

			return resourceAttribute.Id, respObject.GetType(), diags
		}
	}

	return nil, respObject.GetType(), diags
}

func (p *ResourceAttributeResourceModel) expand(resourceType management.EnumResourceType) *management.ResourceAttribute {

	data := management.NewResourceAttribute(p.Name.ValueString(), p.Value.ValueString())

	if resourceType == management.ENUMRESOURCETYPE_OPENID_CONNECT {
		if !p.IDTokenEnabled.IsNull() {
			data.SetIdToken(p.IDTokenEnabled.ValueBool())
		} else {
			data.SetIdToken(true)
		}

		if !p.UserinfoEnabled.IsNull() {
			data.SetUserInfo(p.UserinfoEnabled.ValueBool())
		} else {
			data.SetUserInfo(true)
		}
	}

	return data
}

func (p *ResourceAttributeResourceModel) toState(apiObject *management.ResourceAttribute) diag.Diagnostics {
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
	p.Value = framework.StringOkToTF(apiObject.GetValueOk())
	p.Type = ResourceAttributeTypeOkToTF(apiObject.GetTypeOk())
	p.IDTokenEnabled = framework.BoolOkToTF(apiObject.GetIdTokenOk())
	p.UserinfoEnabled = framework.BoolOkToTF(apiObject.GetUserInfoOk())

	return diags
}

func ResourceAttributeTypeOkToTF(v *management.EnumResourceAttributeType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
