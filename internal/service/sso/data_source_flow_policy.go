// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/davincitypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type FlowPolicyDataSource serviceClientType

type FlowPolicyDataSourceModel struct {
	Id                 davincitypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId      pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	FlowPolicyId       davincitypes.ResourceIDValue `tfsdk:"flow_policy_id"`
	Name               types.String                 `tfsdk:"name"`
	Enabled            types.Bool                   `tfsdk:"enabled"`
	DaVinciApplication types.Object                 `tfsdk:"davinci_application"`
	Trigger            types.Object                 `tfsdk:"trigger"`
}

var (
	dvApplicationTFObjectTypes = map[string]attr.Type{
		"id":   davincitypes.ResourceIDType{},
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
			"id": framework.Attr_IDCustomType(davincitypes.ResourceIDType{}),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the DaVinci flow policy."),
			),

			"flow_policy_id": schema.StringAttribute{
				Description: "The ID of the DaVinci flow policy.",
				Optional:    true,

				CustomType: davincitypes.ResourceIDType{},
			},

			"name": schema.StringAttribute{
				Description: "The name of the DaVinci flow policy.",
				Computed:    true,
			},

			"enabled": schema.BoolAttribute{
				Description: "A boolean to specify whether the flow policy is enabled in the environment or not.",
				Computed:    true,
			},

			"davinci_application": schema.SingleNestedAttribute{
				Description: "A single object that describes the DaVinci application that contains the flow policy.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "A string that specifies the ID of the DaVinci application to which the flow policy is assigned.",
						Computed:    true,

						CustomType: davincitypes.ResourceIDType{},
					},
					"name": schema.StringAttribute{
						Description: "A string that specifies the name of the DaVinci application to which the flow policy is assigned.",
						Computed:    true,
					},
				},
			},

			"trigger": schema.SingleNestedAttribute{
				Description: "A single object that describes the configured DaVinci flow policy trigger.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "A string that specifies the type of the DaVinci flow policy.",
						Computed:    true,
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *FlowPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *FlowPolicyDataSourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	var response *management.FlowPolicy
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.FlowPoliciesApi.ReadOneFlowPolicy(ctx, data.EnvironmentId.ValueString(), data.FlowPolicyId.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneFlowPolicy",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response)...)
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

	p.Id = framework.DaVinciResourceIDToTF(apiObject.GetId())
	p.FlowPolicyId = framework.DaVinciResourceIDToTF(apiObject.GetId())
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

func toStateDavinciApplication(davinciApplication *management.FlowPolicyApplication, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || davinciApplication == nil {
		return types.ObjectNull(dvApplicationTFObjectTypes), diags
	}

	dvApplicationMap := map[string]attr.Value{
		"id":   framework.DaVinciResourceIDOkToTF(davinciApplication.GetIdOk()),
		"name": framework.StringOkToTF(davinciApplication.GetNameOk()),
	}

	returnVar, d := types.ObjectValue(dvApplicationTFObjectTypes, dvApplicationMap)
	diags.Append(d...)

	return returnVar, diags

}

func toStateFlowTrigger(davinciApplication *management.FlowPolicyTrigger, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || davinciApplication == nil {
		return types.ObjectNull(flowTriggerTFObjectTypes), diags
	}

	dvApplicationMap := map[string]attr.Value{}
	if v, ok := davinciApplication.GetTypeOk(); ok {

		dvApplicationMap["type"] = framework.StringToTF(string(*v))

	} else {

		dvApplicationMap["type"] = types.StringNull()

	}

	returnVar, d := types.ObjectValue(flowTriggerTFObjectTypes, dvApplicationMap)
	diags.Append(d...)

	return returnVar, diags

}
