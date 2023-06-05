package base

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkv2resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	boolvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/boolvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type SystemApplicationResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type systemApplicationResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	EnvironmentId             types.String `tfsdk:"environment_id"`
	Type                      types.String `tfsdk:"type"`
	Name                      types.String `tfsdk:"name"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	AccessControlRoleType     types.String `tfsdk:"access_control_role_type"`
	AccessControlGroupOptions types.Object `tfsdk:"access_control_group_options"`
}

type applicationAccessControlGroupOptionsResourceModel struct {
	Type   types.String `tfsdk:"type"`
	Groups types.Set    `tfsdk:"groups"`
}

var (
	applicationAccessControlGroupOptionsTFObjectTypes = map[string]attr.Type{
		"type":   types.StringType,
		"groups": types.SetType{ElemType: types.StringType},
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &SystemApplicationResource{}
	_ resource.ResourceWithConfigure   = &SystemApplicationResource{}
	_ resource.ResourceWithImportState = &SystemApplicationResource{}
)

// New Object
func NewSystemApplicationResource() resource.Resource {
	return &SystemApplicationResource{}
}

// Metadata
func (r *SystemApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_application"
}

// Schema.
func (r *SystemApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of system application, used as the primary identifier.",
	).AllowedValues([]string{
		string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL),
		string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE),
	})

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies the enabled/disabled status of the application.",
	).DefaultValue("true")

	accessControlGroupType := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user role required to access the application. A user is an admin user if they have one or more of the following roles assigned: `Organization Admin`, `Environment Admin`, `Identity Data Admin`, or `Client Application Developer`.",
	).AllowedValues([]string{
		string(management.ENUMAPPLICATIONACCESSCONTROLTYPE_ADMIN_USERS_ONLY),
	})

	accessControlGroupOptionsType := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the group type required to access the application.",
	).AllowedValuesComplex(map[string]string{
		"ANY_GROUP":  "the actor must belong to at least one group listed in the `groups` property",
		"ALL_GROUPS": "the actor must belong to all groups listed in the `groups` property",
	})

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to manage the built-in system applications (PingOne Self-Service and PingOne Portal) in PingOne.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage built-in system applications in."),
			),

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Required:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.OneOf(
						string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL),
						string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE),
					),
				},
			},

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the name of the system application.").Description,
				Computed:    true,
			},

			"enabled": schema.BoolAttribute{
				Description:         enabledDescription.Description,
				MarkdownDescription: enabledDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.Bool{
					boolvalidatorinternal.MustBeTrueIfPathSetToValue(
						types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE)),
						path.MatchRoot("type"),
					),
				},
			},

			"access_control_role_type": schema.StringAttribute{
				Description:         accessControlGroupType.Description,
				MarkdownDescription: accessControlGroupType.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(
						string(management.ENUMAPPLICATIONACCESSCONTROLTYPE_ADMIN_USERS_ONLY),
					),
				},
			},

			"access_control_group_options": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Group access control settings.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"groups": schema.SetAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A set that specifies the group IDs for the groups the actor must belong to for access to the application.").Description,
						Required:    true,

						ElementType: types.StringType,
					},

					"type": schema.StringAttribute{
						Description:         accessControlGroupOptionsType.Description,
						MarkdownDescription: accessControlGroupOptionsType.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(
								"ANY_GROUP",
								"ALL_GROUPS",
							),
						},
					},
				},
			},
		},
	}
}

func (r *SystemApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SystemApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state systemApplicationResourceModel

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

	if !plan.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL))) && !plan.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE))) {
		resp.Diagnostics.AddError(
			"Invalid application type",
			fmt.Sprintf("Application type not supported.  Type found: %s, expected one of: %s, %s.", plan.Type.ValueString(), string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL), string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE)),
		)

		return
	}

	// Build the model for the API
	updateSystemApplication, applicationId, d := plan.expand(ctx, r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.ApplicationsApi.UpdateApplication(ctx, plan.EnvironmentId.ValueString(), *applicationId).UpdateApplicationRequest(*updateSystemApplication).Execute()
		},
		"UpdateApplication",
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
	resp.Diagnostics.Append(state.toState(response.(*management.ReadOneApplication200Response))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SystemApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *systemApplicationResourceModel

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

	if !data.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL))) && !data.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE))) {
		resp.Diagnostics.AddError(
			"Invalid application type",
			fmt.Sprintf("Application type not supported.  Type found: %s, expected one of: %s, %s.", data.Type.ValueString(), string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL), string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE)),
		)

		resp.State.RemoveResource(ctx)

		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.ApplicationsApi.ReadOneApplication(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneApplication",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response.(*management.ReadOneApplication200Response))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state systemApplicationResourceModel

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

	if !plan.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL))) && !plan.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE))) {
		resp.Diagnostics.AddError(
			"Invalid application type",
			fmt.Sprintf("Application type not supported.  Type found: %s, expected one of: %s, %s.", plan.Type.ValueString(), string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL), string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE)),
		)

		resp.State.RemoveResource(ctx)

		return
	}

	// Build the model for the API
	updateSystemApplication, _, d := plan.expand(ctx, r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.ApplicationsApi.UpdateApplication(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).UpdateApplicationRequest(*updateSystemApplication).Execute()
		},
		"UpdateApplication",
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
	resp.Diagnostics.Append(state.toState(response.(*management.ReadOneApplication200Response))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SystemApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *systemApplicationResourceModel

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

	data.Enabled = types.BoolValue(true)
	data.AccessControlGroupOptions = types.ObjectNull(applicationAccessControlGroupOptionsTFObjectTypes)
	data.AccessControlRoleType = types.StringNull()

	updateSystemApplication, _, d := data.expand(ctx, r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	_, d = framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.ApplicationsApi.UpdateApplication(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).UpdateApplicationRequest(*updateSystemApplication).Execute()
		},
		"UpdateApplication",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *SystemApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/application_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *systemApplicationResourceModel) expand(ctx context.Context, apiClient *management.APIClient) (*management.UpdateApplicationRequest, *string, diag.Diagnostics) {
	var diags diag.Diagnostics

	var applicationId *string
	applicationType := management.EnumApplicationType(p.Type.ValueString())

	adminConsole, portal, selfService, d := FetchSystemApplication(ctx, apiClient, p.EnvironmentId.ValueString(), applicationType)
	diags.Append(d...)
	if diags.HasError() {
		return nil, nil, diags
	}

	accessControl := management.NewApplicationAccessControl()
	setAccessControl := false

	if !p.AccessControlRoleType.IsNull() && !p.AccessControlRoleType.IsUnknown() {
		accessControl.SetRole(
			*management.NewApplicationAccessControlRole(management.EnumApplicationAccessControlType(p.AccessControlRoleType.ValueString())),
		)
		setAccessControl = true
	}

	if !p.AccessControlGroupOptions.IsNull() && !p.AccessControlGroupOptions.IsUnknown() {
		var plan applicationAccessControlGroupOptionsResourceModel
		diags.Append(p.AccessControlGroupOptions.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)

		var groupsPlan []string
		diags.Append(plan.Groups.ElementsAs(ctx, &groupsPlan, false)...)

		groups := make([]management.ApplicationAccessControlGroupGroupsInner, 0)
		for _, v := range groupsPlan {
			groups = append(groups, *management.NewApplicationAccessControlGroupGroupsInner(v))
		}

		accessControl.SetGroup(
			*management.NewApplicationAccessControlGroup(
				plan.Type.ValueString(),
				groups,
			),
		)
		setAccessControl = true
	}

	data := management.UpdateApplicationRequest{}

	if adminConsole != nil {
		diags.AddError(
			"Admin Console Immutable",
			"The PingOne Admin Console is immutable and cannot be changed.  Please report this to the provider maintainers.",
		)

		return nil, nil, diags
	}

	if v := portal; v != nil {

		updateApplication := management.NewApplicationPingOnePortal(
			p.Enabled.ValueBool(),
			v.GetName(),
			v.GetProtocol(),
			v.GetType(),
			v.GetTokenEndpointAuthMethod(),
			v.GetApplyDefaultTheme(),
		)

		updateApplication.SetDescription(v.GetDescription())
		updateApplication.SetIcon(v.GetIcon())
		updateApplication.SetHiddenFromAppPortal(v.GetHiddenFromAppPortal())

		if setAccessControl {
			updateApplication.SetAccessControl(*accessControl)
		}

		data.ApplicationPingOnePortal = updateApplication

		var ok bool
		if applicationId, ok = v.GetIdOk(); !ok {
			diags.AddError(
				"Portal Application ID Cannot be Retrieved",
				"The PingOne Portal ID cannot be retrieved.  Please report this to the provider maintainers.",
			)

			return nil, nil, diags
		}
	}

	if v := selfService; v != nil {

		updateApplication := management.NewApplicationPingOneSelfService(
			p.Enabled.ValueBool(),
			v.GetName(),
			v.GetProtocol(),
			v.GetType(),
			v.GetTokenEndpointAuthMethod(),
			v.GetApplyDefaultTheme(),
		)

		updateApplication.SetDescription(v.GetDescription())
		updateApplication.SetIcon(v.GetIcon())
		updateApplication.SetHiddenFromAppPortal(v.GetHiddenFromAppPortal())

		if setAccessControl {
			updateApplication.SetAccessControl(*accessControl)
		}

		data.ApplicationPingOneSelfService = updateApplication

		var ok bool
		if applicationId, ok = v.GetIdOk(); !ok {
			diags.AddError(
				"Self Service Application ID Cannot be Retrieved",
				"The PingOne Self-Service ID cannot be retrieved.  Please report this to the provider maintainers.",
			)

			return nil, nil, diags
		}
	}

	return &data, applicationId, diags
}

func (p *systemApplicationResourceModel) toState(apiObject *management.ReadOneApplication200Response) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	apiObjectCommon := management.Application{}

	if apiObject.ApplicationPingOnePortal != nil {
		apiObjectCommon = management.Application{
			Id:            apiObject.ApplicationPingOnePortal.Id,
			Name:          apiObject.ApplicationPingOnePortal.Name,
			Environment:   apiObject.ApplicationPingOnePortal.Environment,
			Type:          apiObject.ApplicationPingOnePortal.Type,
			Enabled:       apiObject.ApplicationPingOnePortal.Enabled,
			Description:   apiObject.ApplicationPingOnePortal.Description,
			AccessControl: apiObject.ApplicationPingOnePortal.AccessControl,
		}
	}

	if apiObject.ApplicationPingOneSelfService != nil {
		apiObjectCommon = management.Application{
			Id:            apiObject.ApplicationPingOneSelfService.Id,
			Name:          apiObject.ApplicationPingOneSelfService.Name,
			Environment:   apiObject.ApplicationPingOneSelfService.Environment,
			Type:          apiObject.ApplicationPingOneSelfService.Type,
			Enabled:       apiObject.ApplicationPingOneSelfService.Enabled,
			Description:   apiObject.ApplicationPingOneSelfService.Description,
			AccessControl: apiObject.ApplicationPingOneSelfService.AccessControl,
		}
	}

	p.Id = framework.StringToTF(apiObjectCommon.GetId())
	p.EnvironmentId = framework.StringToTF(*apiObjectCommon.GetEnvironment().Id)
	p.Type = framework.EnumOkToTF(apiObjectCommon.GetTypeOk())
	p.Name = framework.StringOkToTF(apiObjectCommon.GetNameOk())
	p.Enabled = framework.BoolOkToTF(apiObjectCommon.GetEnabledOk())

	p.AccessControlRoleType = types.StringNull()
	p.AccessControlGroupOptions = types.ObjectNull(applicationAccessControlGroupOptionsTFObjectTypes)
	if v, ok := apiObjectCommon.GetAccessControlOk(); ok {
		if v1, ok := v.GetRoleOk(); ok {
			p.AccessControlRoleType = framework.EnumOkToTF(v1.GetTypeOk())
		}

		if v1, ok := v.GetGroupOk(); ok {

			if v2, ok := v1.GetGroupsOk(); ok {

				groupsSlice := make([]string, 0)

				for _, group := range v2 {
					groupsSlice = append(groupsSlice, group.GetId())
				}

				tfGroupsSlice := framework.StringSetToTF(groupsSlice)

				objValue, d := types.ObjectValue(applicationAccessControlGroupOptionsTFObjectTypes, map[string]attr.Value{
					"groups": tfGroupsSlice,
					"type":   framework.EnumOkToTF(v1.GetTypeOk()),
				})
				diags.Append(d...)

				p.AccessControlGroupOptions = objValue
			}
		}
	}

	return diags
}

func FetchSystemApplication(ctx context.Context, apiClient *management.APIClient, environmentID string, systemApplicationType management.EnumApplicationType) (*management.ApplicationPingOneAdminConsole, *management.ApplicationPingOnePortal, *management.ApplicationPingOneSelfService, diag.Diagnostics) {
	var diags diag.Diagnostics
	if systemApplicationType != management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL && systemApplicationType != management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE && systemApplicationType != management.ENUMAPPLICATIONTYPE_PING_ONE_ADMIN_CONSOLE {
		diags.AddError(
			"Invalid system application type",
			fmt.Sprintf("Invalid application type: %s, expecting %s, %s or %s", systemApplicationType, management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL, management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE, management.ENUMAPPLICATIONTYPE_PING_ONE_ADMIN_CONSOLE),
		)
		return nil, nil, nil, diags
	}

	systemApps, d := FetchApplicationsByType(ctx, apiClient, environmentID, systemApplicationType, true)
	diags.Append(d...)
	if diags.HasError() {
		return nil, nil, nil, diags
	}

	if len(*systemApps) == 1 {
		systemApp := (*systemApps)[0]

		if v := systemApp.ApplicationPingOneAdminConsole; v != nil {
			return v, nil, nil, diags
		}

		if v := systemApp.ApplicationPingOnePortal; v != nil {
			return nil, v, nil, diags
		}

		if v := systemApp.ApplicationPingOneSelfService; v != nil {
			return nil, nil, v, diags
		}
	}

	if len(*systemApps) > 1 {
		diags.AddError(
			"Unexpected applications found",
			fmt.Sprintf("More than one application for type %s found.  Please report to the provider maintainers.", systemApplicationType),
		)
	}

	diags.AddError(
		"System application not found",
		fmt.Sprintf("System application type %s not found.  Please report to the provider maintainers.", systemApplicationType),
	)

	return nil, nil, nil, diags
}

func FetchApplicationsByType(ctx context.Context, apiClient *management.APIClient, environmentID string, applicationType management.EnumApplicationType, expectAtLeastOneResult bool) (*[]management.ReadOneApplication200Response, diag.Diagnostics) {
	defaultTimeout := 10 * time.Second
	return FetchApplicationsByTypeWithTimeout(ctx, apiClient, environmentID, applicationType, expectAtLeastOneResult, defaultTimeout)
}

func FetchApplicationsByTypeWithTimeout(ctx context.Context, apiClient *management.APIClient, environmentID string, applicationType management.EnumApplicationType, expectAtLeastOneResult bool, timeout time.Duration) (*[]management.ReadOneApplication200Response, diag.Diagnostics) {
	var diags diag.Diagnostics

	stateConf := &sdkv2resource.StateChangeConf{
		Pending: []string{
			"false",
		},
		Target: []string{
			"true",
			"err",
		},
		Refresh: func() (interface{}, string, error) {

			// Run the API call
			response, d := framework.ParseResponse(
				ctx,

				func() (interface{}, *http.Response, error) {
					return apiClient.ApplicationsApi.ReadAllApplications(ctx, environmentID).Execute()
				},
				"ReadAllApplications",
				framework.DefaultCustomError,
				sdk.DefaultCreateReadRetryable,
			)
			diags.Append(d...)
			if diags.HasError() {
				return nil, "err", fmt.Errorf("Error reading applications")
			}

			entityArray := response.(*management.EntityArray)

			found := false

			var applicationResponse []management.ReadOneApplication200Response

			if applications, ok := entityArray.Embedded.GetApplicationsOk(); ok {

				for _, applicationItem := range applications {

					data, err := json.Marshal(applicationItem)
					if err != nil {
						return nil, "err", fmt.Errorf("Error marshalling application: %s", err)
					}

					var common management.Application
					if err := json.Unmarshal(data, &common); err != nil {
						return nil, "err", fmt.Errorf("Error unmarshalling application: %s", err)
					}

					if common.GetType() == applicationType {
						applicationResponse = append(applicationResponse, applicationItem)
						found = true
					}
				}
			}

			if !found && expectAtLeastOneResult {
				return nil, "false", fmt.Errorf("No applications found for type %s, but at least one is expected", applicationType)
			}

			tflog.Debug(ctx, "Find applications by type attempt", map[string]interface{}{
				"applicationResponse":      applicationResponse,
				"len(applicationResponse)": len(applicationResponse),
				"result":                   strings.ToLower(strconv.FormatBool(found)),
			})

			return applicationResponse, strings.ToLower(strconv.FormatBool(found)), nil
		},
		Timeout:                   timeout,
		Delay:                     1 * time.Second,
		MinTimeout:                2 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	applicationResponse, err := stateConf.WaitForState()

	if err != nil {
		diags.AddError(
			"Cannot find applications by type",
			fmt.Sprintf("The applications by type %s for environment %s cannot be found: %s", applicationType, environmentID, err),
		)

		return nil, diags
	}

	returnVar := applicationResponse.([]management.ReadOneApplication200Response)

	return &returnVar, diags

}
