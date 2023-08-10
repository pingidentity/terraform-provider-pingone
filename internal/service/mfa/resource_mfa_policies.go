package mfa

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type MFAPoliciesResource struct {
	client *mfa.APIClient
	region model.RegionMapping
}

type MFAPoliciesResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	MigrateData   types.Set    `tfsdk:"migrate_data"`
}

type MFAPoliciesMigrateDataResourceModel struct {
	DeviceAuthenticationPolicyId types.String `tfsdk:"device_authentication_policy_id"`
	Fido2PolicyId                types.String `tfsdk:"fido2_policy_id"`
}

// Framework interfaces
var (
	_ resource.Resource              = &MFAPoliciesResource{}
	_ resource.ResourceWithConfigure = &MFAPoliciesResource{}
)

// New Object
func NewMFAPoliciesResource() resource.Resource {
	return &MFAPoliciesResource{}
}

// Metadata
func (r *MFAPoliciesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_policies"
}

func (r *MFAPoliciesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage bulk settings of MFA device policies in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure MFA device policies in."),
			),

			"migrate_data": schema.SetNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of objects that describe MFA Device policies to migrate.").Description,
				Required:    true,

				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_authentication_policy_id": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID of an MFA Device policy to migrate.").Description,
							Required:    true,

							Validators: []validator.String{
								verify.P1ResourceIDValidator(),
							},
						},

						"fido2_policy_id": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID of a FIDO2 policy to assign to the new FIDO2 device type.").Description,
							Optional:    true,

							Validators: []validator.String{
								verify.P1ResourceIDValidator(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *MFAPoliciesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	preparedClient, err := PrepareClient(ctx, resourceConfig)
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

func (r *MFAPoliciesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state MFAPoliciesResourceModel

	if r.client == nil {
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
	mfaPolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *mfa.DeviceAuthenticationPolicyPostResponse
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.DeviceAuthenticationPolicyApi.CreateDeviceAuthenticationPolicies(ctx, plan.EnvironmentId.ValueString()).ContentType(mfa.ENUMDEVICEAUTHENTICATIONPOLICYPOSTCONTENTTYPE_VND_PINGIDENTITY_DEVICE_AUTHENTICATION_POLICY_FIDO2_MIGRATEJSON).DeviceAuthenticationPolicyPost(*mfaPolicy).Execute()
		},
		"CreateDeviceAuthenticationPolicies",
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
	resp.Diagnostics.Append(state.toState(response.EntityArray)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *MFAPoliciesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *MFAPoliciesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *MFAPoliciesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (p *MFAPoliciesResourceModel) expand(ctx context.Context) (*mfa.DeviceAuthenticationPolicyPost, diag.Diagnostics) {
	var diags diag.Diagnostics

	migrateData := make([]mfa.DeviceAuthenticationPolicyMigrateData, 0)
	if !p.MigrateData.IsNull() && !p.MigrateData.IsUnknown() {

		var migrateDataPlan []MFAPoliciesMigrateDataResourceModel
		diags.Append(p.MigrateData.ElementsAs(ctx, &migrateDataPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, migrateDataItemPlan := range migrateDataPlan {

			migrateDataItem := *mfa.NewDeviceAuthenticationPolicyMigrateData(
				migrateDataItemPlan.DeviceAuthenticationPolicyId.ValueString(),
			)

			if !migrateDataItemPlan.Fido2PolicyId.IsNull() && !migrateDataItemPlan.Fido2PolicyId.IsUnknown() {
				migrateDataItem.SetFido2PolicyId(migrateDataItemPlan.Fido2PolicyId.ValueString())
			}

			migrateData = append(migrateData, migrateDataItem)
		}
	}

	data := mfa.DeviceAuthenticationPolicyPost{}

	// Main object
	data.DeviceAuthenticationPolicyMigrate = mfa.NewDeviceAuthenticationPolicyMigrate(migrateData)

	return &data, diags
}

func (p *MFAPoliciesResourceModel) toState(apiObject *mfa.EntityArray) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	if p.Id.IsNull() || p.Id.IsUnknown() {
		p.Id = framework.StringToTF(uuid.New().String())
	}

	p.EnvironmentId = framework.StringToTF(p.EnvironmentId.ValueString())

	return diags
}
