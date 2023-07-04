package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type FlowPolicyDataSource struct {
	client *management.APIClient
	region model.RegionMapping
}

type FlowPolicyDataSourceModel struct {
	Id                 types.String `tfsdk:"id"`
	EnvironmentId      types.String `tfsdk:"environment_id"`
	FlowPolicyId       types.String `tfsdk:"flow_policy_id"`
	Name               types.String `tfsdk:"name"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	DaVinciApplication types.List   `tfsdk:"davinci_application"`
	Trigger            types.List   `tfsdk:"trigger"`
}

var (
	dvApplicationTFObjectTypes = map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}

	flowTriggerTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
	}
)

// Framework interfaces
var (
	_ datasource.DataSource = &FlowPolicyDataSource{}
)

// New Object
func NewFlowPolicyDataSource() datasource.DataSource {
	return &FlowPolicyDataSource{}
}

// Metadata
func (r *FlowPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flow_policy"
}

// Schema
func (r *FlowPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a PingOne DaVinci flow policy.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the DaVinci flow policy."),
			),

			"flow_policy_id": schema.StringAttribute{
				Description: "The ID of the DaVinci flow policy.",
				Optional:    true,
				Validators: []validator.String{
					verify.P1DVResourceIDValidator(),
				},
			},

			"name": schema.StringAttribute{
				Description: "The name of the DaVinci flow policy.",
				Computed:    true,
			},

			"enabled": schema.BoolAttribute{
				Description: "A boolean to specify whether the flow policy is enabled in the environment or not.",
				Computed:    true,
			},
		},

		Blocks: map[string]schema.Block{
			"davinci_application": schema.ListNestedBlock{
				Description: "A block that describes the DaVinci application that contains the flow policy.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "A string that specifies the ID of the DaVinci application to which the flow policy is assigned.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "A string that specifies the name of the DaVinci application to which the flow policy is assigned.",
							Computed:    true,
						},
					},
				},
			},

			"trigger": schema.ListNestedBlock{
				Description: "A block that describes the configured DaVinci flow policy trigger.",

				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "A string that specifies the type of the DaVinci flow policy.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *FlowPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *FlowPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *FlowPolicyDataSourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.FlowPoliciesApi.ReadOneFlowPolicy(ctx, data.EnvironmentId.ValueString(), data.FlowPolicyId.ValueString()).Execute()
		},
		"ReadOneFlowPolicy",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response.(*management.FlowPolicy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *FlowPolicyDataSourceModel) toState(apiObject *management.FlowPolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.FlowPolicyId = framework.StringToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())

	davinciApplication, d := toStateDavinciApplication(apiObject.GetApplicationOk())
	diags.Append(d...)
	p.DaVinciApplication = davinciApplication

	trigger, d := toStateFlowTrigger(apiObject.GetTriggerOk())
	diags.Append(d...)
	p.Trigger = trigger

	return diags
}

func toStateDavinciApplication(davinciApplication *management.FlowPolicyApplication, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: dvApplicationTFObjectTypes}

	if !ok || davinciApplication == nil {
		return types.ListValueMust(tfObjType, []attr.Value{}), diags
	}

	dvApplicationMap := map[string]attr.Value{
		"id":   framework.StringOkToTF(davinciApplication.GetIdOk()),
		"name": framework.StringOkToTF(davinciApplication.GetNameOk()),
	}

	flattenedObj, d := types.ObjectValue(dvApplicationTFObjectTypes, dvApplicationMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags

}

func toStateFlowTrigger(davinciApplication *management.FlowPolicyTrigger, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: flowTriggerTFObjectTypes}

	if !ok || davinciApplication == nil {
		return types.ListValueMust(tfObjType, []attr.Value{}), diags
	}

	dvApplicationMap := map[string]attr.Value{}
	if v, ok := davinciApplication.GetTypeOk(); ok {

		dvApplicationMap["type"] = framework.StringToTF(string(*v))

	} else {

		dvApplicationMap["type"] = types.StringNull()

	}

	flattenedObj, d := types.ObjectValue(flowTriggerTFObjectTypes, dvApplicationMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags

}
