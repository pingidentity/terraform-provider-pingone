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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	boolvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/boolvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type SystemApplicationResource serviceClientType

type systemApplicationResourceModel struct {
	Id                        pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId             pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Type                      types.String                 `tfsdk:"type"`
	Name                      types.String                 `tfsdk:"name"`
	Enabled                   types.Bool                   `tfsdk:"enabled"`
	AccessControlRoleType     types.String                 `tfsdk:"access_control_role_type"`
	AccessControlGroupOptions types.Object                 `tfsdk:"access_control_group_options"`
	ApplyDefaultTheme         types.Bool                   `tfsdk:"apply_default_theme"`
	EnableDefaultThemeFooter  types.Bool                   `tfsdk:"enable_default_theme_footer"`
}

type applicationAccessControlGroupOptionsResourceModel struct {
	Type   types.String `tfsdk:"type"`
	Groups types.Set    `tfsdk:"groups"`
}

var (
	applicationAccessControlGroupOptionsTFObjectTypes = map[string]attr.Type{
		"type":   types.StringType,
		"groups": types.SetType{ElemType: pingonetypes.ResourceIDType{}},
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
	).AllowedValues(
		string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL),
		string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE),
	)

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies the enabled/disabled status of the application.",
	).DefaultValue("true")

	accessControlGroupType := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user role required to access the application. A user is an admin user if they have one or more of the following roles assigned: `Organization Admin`, `Environment Admin`, `Identity Data Admin`, or `Client Application Developer`.",
	).AllowedValues(
		string(management.ENUMAPPLICATIONACCESSCONTROLTYPE_ADMIN_USERS_ONLY),
	)

	accessControlGroupOptionsType := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the group type required to access the application.",
	).AllowedValuesComplex(map[string]string{
		"ANY_GROUP":  "the actor must belong to at least one group listed in the `groups` property",
		"ALL_GROUPS": "the actor must belong to all groups listed in the `groups` property",
	})

	applyDefaultThemeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether to apply the default theme to the Self-Service or PingOne Portal application.",
	).DefaultValue(false)

	enableDefaultThemeFooterDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether to show the default theme footer on the self-service application. Configurable only when the `type` is `PING_ONE_SELF_SERVICE` and `apply_default_theme` is also `true`.",
	)

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

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

						ElementType: pingonetypes.ResourceIDType{},
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

			"apply_default_theme": schema.BoolAttribute{
				Description:         applyDefaultThemeDescription.Description,
				MarkdownDescription: applyDefaultThemeDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"enable_default_theme_footer": schema.BoolAttribute{
				Description:         enableDefaultThemeFooterDescription.Description,
				MarkdownDescription: enableDefaultThemeFooterDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Bool{
					boolvalidatorinternal.ConflictsIfMatchesPathValue(
						types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL)),
						path.MatchRoot("type"),
					),
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *SystemApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state systemApplicationResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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

	if !plan.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL))) && !plan.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE))) {
		resp.Diagnostics.AddError(
			"Invalid application type",
			fmt.Sprintf("Application type not supported.  Type found: %s, expected one of: %s, %s.", plan.Type.ValueString(), string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL), string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE)),
		)

		return
	}

	// Build the model for the API
	updateSystemApplication, applicationId, d := plan.expand(ctx, r.Client.ManagementAPIClient)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.ReadOneApplication200Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.UpdateApplication(ctx, plan.EnvironmentId.ValueString(), *applicationId).UpdateApplicationRequest(*updateSystemApplication).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateApplication",
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

func (r *SystemApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *systemApplicationResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	var response *management.ReadOneApplication200Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.ReadOneApplication(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneApplication",
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

	if !data.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL))) && !data.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE))) {
		resp.Diagnostics.AddError(
			"Invalid application type",
			fmt.Sprintf("Application type not supported.  Type found: %s, expected one of: %s, %s.", data.Type.ValueString(), string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL), string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE)),
		)

		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state systemApplicationResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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

	if !plan.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL))) && !plan.Type.Equal(types.StringValue(string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE))) {
		resp.Diagnostics.AddError(
			"Invalid application type",
			fmt.Sprintf("Application type not supported.  Type found: %s, expected one of: %s, %s.", plan.Type.ValueString(), string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL), string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE)),
		)

		resp.State.RemoveResource(ctx)

		return
	}

	// Build the model for the API
	updateSystemApplication, _, d := plan.expand(ctx, r.Client.ManagementAPIClient)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.ReadOneApplication200Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.UpdateApplication(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).UpdateApplicationRequest(*updateSystemApplication).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateApplication",
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

func (r *SystemApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *systemApplicationResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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

	data.Enabled = types.BoolValue(true)
	data.AccessControlGroupOptions = types.ObjectNull(applicationAccessControlGroupOptionsTFObjectTypes)
	data.AccessControlRoleType = types.StringNull()

	updateSystemApplication, _, d := data.expand(ctx, r.Client.ManagementAPIClient)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.UpdateApplication(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).UpdateApplicationRequest(*updateSystemApplication).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateApplication",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *SystemApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "application_id",
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

		var groupsPlan []pingonetypes.ResourceIDValue
		diags.Append(plan.Groups.ElementsAs(ctx, &groupsPlan, false)...)

		groupsStr, d := framework.TFTypePingOneResourceIDSliceToStringSlice(groupsPlan, path.Root("access_control_group_options").AtName("groups"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, nil, diags
		}

		groups := make([]management.ApplicationAccessControlGroupGroupsInner, 0)
		for _, v := range groupsStr {
			groups = append(groups, *management.NewApplicationAccessControlGroupGroupsInner(v))
		}

		accessControl.SetGroup(
			*management.NewApplicationAccessControlGroup(
				management.EnumApplicationAccessControlGroupType(plan.Type.ValueString()),
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

		if !p.ApplyDefaultTheme.IsNull() && !p.ApplyDefaultTheme.IsUnknown() {
			updateApplication.SetApplyDefaultTheme(p.ApplyDefaultTheme.ValueBool())
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

		if !p.ApplyDefaultTheme.IsNull() && !p.ApplyDefaultTheme.IsUnknown() {
			updateApplication.SetApplyDefaultTheme(p.ApplyDefaultTheme.ValueBool())
		}

		if !p.EnableDefaultThemeFooter.IsNull() && !p.EnableDefaultThemeFooter.IsUnknown() {
			updateApplication.SetEnableDefaultThemeFooter(p.EnableDefaultThemeFooter.ValueBool())
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

		p.ApplyDefaultTheme = framework.BoolOkToTF(apiObject.ApplicationPingOnePortal.GetApplyDefaultThemeOk())
		p.EnableDefaultThemeFooter = types.BoolNull()
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

		p.ApplyDefaultTheme = framework.BoolOkToTF(apiObject.ApplicationPingOneSelfService.GetApplyDefaultThemeOk())
		p.EnableDefaultThemeFooter = framework.BoolOkToTF(apiObject.ApplicationPingOneSelfService.GetEnableDefaultThemeFooterOk())
	}

	p.Id = framework.PingOneResourceIDToTF(apiObjectCommon.GetId())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObjectCommon.GetEnvironment().Id)
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

				tfGroupsSlice := framework.PingOneResourceIDSetToTF(groupsSlice)

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

	stateConf := &retry.StateChangeConf{
		Pending: []string{
			"false",
		},
		Target: []string{
			"true",
			"err",
		},
		Refresh: func() (interface{}, string, error) {

			// Run the API call
			var applicationResponse []management.ReadOneApplication200Response
			diags.Append(framework.ParseResponse(
				ctx,

				func() (any, *http.Response, error) {
					pagedIterator := apiClient.ApplicationsApi.ReadAllApplications(ctx, environmentID).Execute()

					var foundApplications []management.ReadOneApplication200Response

					var initialHttpResponse *http.Response

					for pageCursor, err := range pagedIterator {
						if err != nil {
							return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentID, nil, pageCursor.HTTPResponse, err)
						}

						if initialHttpResponse == nil {
							initialHttpResponse = pageCursor.HTTPResponse
						}

						if applications, ok := pageCursor.EntityArray.Embedded.GetApplicationsOk(); ok {

							for _, applicationItem := range applications {

								data, err := json.Marshal(applicationItem)
								if err != nil {
									return nil, pageCursor.HTTPResponse, fmt.Errorf("Error marshalling application: %s", err)
								}

								var common management.Application
								if err := json.Unmarshal(data, &common); err != nil {
									return nil, pageCursor.HTTPResponse, fmt.Errorf("Error unmarshalling application: %s", err)
								}

								if common.GetType() == applicationType {
									foundApplications = append(foundApplications, applicationItem)
								}
							}
						}
					}

					return foundApplications, initialHttpResponse, nil
				},
				"ReadAllApplications",
				framework.DefaultCustomError,
				sdk.DefaultCreateReadRetryable,
				&applicationResponse,
			)...)
			if diags.HasError() {
				return nil, "err", fmt.Errorf("Error reading applications")
			}

			tflog.Debug(ctx, "Find applications by type attempt", map[string]interface{}{
				"applicationResponse":      applicationResponse,
				"len(applicationResponse)": len(applicationResponse),
				"result":                   strings.ToLower(strconv.FormatBool(len(applicationResponse) > 0)),
			})

			if len(applicationResponse) == 0 && expectAtLeastOneResult {
				return nil, "false", nil
			}

			return applicationResponse, strings.ToLower(strconv.FormatBool(len(applicationResponse) > 0)), nil
		},
		Timeout:                   timeout,
		Delay:                     1 * time.Second,
		MinTimeout:                2 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	applicationResponse, err := stateConf.WaitForStateContext(ctx)

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
