package base

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	stringdefaultinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringdefaultinternal"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type EnvironmentResource struct {
	serviceClientType
	region  model.RegionMapping
	options client.GlobalOptions
}

type environmentResourceModel struct {
	Id             pingonetypes.ResourceIDValue `tfsdk:"id"`
	Name           types.String                 `tfsdk:"name"`
	Description    types.String                 `tfsdk:"description"`
	Type           types.String                 `tfsdk:"type"`
	Region         types.String                 `tfsdk:"region"`
	LicenseId      pingonetypes.ResourceIDValue `tfsdk:"license_id"`
	OrganizationId pingonetypes.ResourceIDValue `tfsdk:"organization_id"`
	Solution       types.String                 `tfsdk:"solution"`
	Services       types.Set                    `tfsdk:"services"`
}

type environmentServiceModel struct {
	Type       types.String `tfsdk:"type"`
	ConsoleUrl types.String `tfsdk:"console_url"`
	Bookmarks  types.Set    `tfsdk:"bookmarks"`
	Tags       types.Set    `tfsdk:"tags"`
}

type environmentServiceBookmarkModel struct {
	Name types.String `tfsdk:"name"`
	Url  types.String `tfsdk:"url"`
}

var (
	environmentServiceTFObjectTypes = map[string]attr.Type{
		"type":        types.StringType,
		"console_url": types.StringType,
		"bookmarks":   types.SetType{ElemType: types.ObjectType{AttrTypes: environmentServiceBookmarkTFObjectTypes}},
		"tags":        types.SetType{ElemType: types.StringType},
	}

	environmentServiceBookmarkTFObjectTypes = map[string]attr.Type{
		"name": types.StringType,
		"url":  types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                   = &EnvironmentResource{}
	_ resource.ResourceWithConfigure      = &EnvironmentResource{}
	_ resource.ResourceWithImportState    = &EnvironmentResource{}
	_ resource.ResourceWithModifyPlan     = &EnvironmentResource{}
	_ resource.ResourceWithValidateConfig = &EnvironmentResource{}
)

// New Object
func NewEnvironmentResource() resource.Resource {
	return &EnvironmentResource{}
}

// Metadata
func (r *EnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

// Schema
func (r *EnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const emailAddressMaxLength = 5

	const maximumServiceBookmarks = 5
	const maximumServices = 13
	const minimumServices = 1

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of the environment to create.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMENVIRONMENTTYPE_SANDBOX):    "for a development/testing environment",
		string(management.ENUMENVIRONMENTTYPE_PRODUCTION): "for environments that require protection from deletion",
	}).AppendMarkdownString("Once an environment has been set as `PRODUCTION` type, it cannot be reset back to `SANDBOX` within Terraform.  Administrators must log in to the web admin console to override the data protection features of `PRODUCTION` environments.").DefaultValue(string(management.ENUMENVIRONMENTTYPE_SANDBOX))

	regionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the region to create the environment in.  Should be consistent with the PingOne organisation region.",
	).AllowedValuesEnum(management.AllowedEnumRegionCodeEnumValues).AppendMarkdownString("Will default to the region specified in the provider configuration if not specified, or can be set with the `PINGONE_REGION_CODE` environment variable.")

	solutionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the solution context of the environment.  Leave blank for a custom, non-workforce solution context.  Valid options are `%s`, or no value for custom solution context.  Workforce solution environments are not yet supported in this provider resource, but can be fetched using the `pingone_environment` datasource.", string(management.ENUMSOLUTIONTYPE_CUSTOMER)),
	).RequiresReplace()

	servicesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of objects that specify the services to enable in the environment.",
	)

	serviceTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the service type to enable in the environment.",
	).AllowedValuesEnum(model.ProductsSelectableList())

	serviceConsoleUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the custom console URL to set.  Generally used with services that are deployed separately to the PingOne SaaS service, such as `PingFederate`, `PingAccess`, `PingDirectory`, `PingAuthorize` and `PingCentral`.",
	)

	daVinciService, err := model.FindProductByAPICode(management.ENUMPRODUCTTYPE_ONE_DAVINCI)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot find DaVinci product",
			"In compiling the schema, the DaVinci product could not be found.  This is always a bug in the provider.  Please report this issue to the provider maintainers.",
		)

		return
	}

	serviceTagsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A set of string tags to apply upon environment creation.  Only configurable when the service `type` is `%s`.", daVinciService.ProductCode),
	).AllowedValuesComplex(
		map[string]string{
			string(management.ENUMBILLOFMATERIALSPRODUCTTAGS_DAVINCI_MINIMAL): "allows for a creation of an environment without example/demo configuration in the DaVinci service",
		},
	).RequiresReplace()

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne environments.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the name of the environment.").Description,

				Required: true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description to apply to the environment.").Description,

				Optional: true,
			},

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,

				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(string(management.ENUMENVIRONMENTTYPE_SANDBOX)),

				Validators: []validator.String{
					stringvalidator.OneOf(func() []string {
						strings := make([]string, 0)
						for _, v := range management.AllowedEnumEnvironmentTypeEnumValues {
							strings = append(strings, string(v))
						}
						return strings
					}()...),
				},
			},

			"region": schema.StringAttribute{
				Description:         regionDescription.Description,
				MarkdownDescription: regionDescription.MarkdownDescription,

				Optional: true,
				Computed: true,

				Default: stringdefaultinternal.StaticStringUnknownable(func() basetypes.StringValue {

					if v := os.Getenv("PINGONE_TERRAFORM_REGION_OVERRIDE"); v != "" {
						return framework.StringToTF(v)
					}

					if v := os.Getenv("PINGONE_REGION_CODE"); v != "" {
						return framework.StringToTF(v)
					}

					if v := os.Getenv("PINGONE_REGION"); v != "" {
						return framework.StringToTF(v)
					}

					if r.region.APICode != "" {
						return types.StringValue(string(r.region.APICode))
					}

					return types.StringUnknown()
				}()),

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.OneOf(
						func() []string {
							if v := os.Getenv("PINGONE_TERRAFORM_REGION_OVERRIDE"); v != "" {
								return []string{
									v,
								}
							}

							return utils.EnumSliceToStringSlice(management.AllowedEnumRegionCodeEnumValues)
						}()...),
				},
			},

			"license_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of a valid license to apply to the environment.  Must be a valid PingOne resource ID.").Description,

				Required: true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"organization_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that represents the ID of the PingOne organization tenant to which the environment belongs.").Description,

				Computed: true,

				CustomType: pingonetypes.ResourceIDType{},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"solution": schema.StringAttribute{
				Description:         solutionDescription.Description,
				MarkdownDescription: solutionDescription.MarkdownDescription,

				Optional: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"services": schema.SetNestedAttribute{
				Description:         servicesDescription.Description,
				MarkdownDescription: servicesDescription.MarkdownDescription,
				Required:            true,

				NestedObject: schema.NestedAttributeObject{

					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description:         serviceTypeDescription.Description,
							MarkdownDescription: serviceTypeDescription.MarkdownDescription,

							Required: true,

							Validators: []validator.String{
								stringvalidator.OneOf(model.ProductsSelectableList()...),
							},
						},

						"console_url": schema.StringAttribute{
							Description:         serviceConsoleUrlDescription.Description,
							MarkdownDescription: serviceConsoleUrlDescription.MarkdownDescription,

							Optional: true,
						},

						"tags": schema.SetAttribute{
							Description:         serviceTagsDescription.Description,
							MarkdownDescription: serviceTagsDescription.MarkdownDescription,

							ElementType: types.StringType,

							Optional: true,

							PlanModifiers: []planmodifier.Set{
								setplanmodifier.RequiresReplace(),
							},

							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumBillOfMaterialsProductTagsEnumValues)...),
								),
							},
						},

						"bookmarks": schema.SetNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of objects that specify custom bookmark links for the service.").Description,
							Optional:    true,

							NestedObject: schema.NestedAttributeObject{

								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the bookmark name.").Description,

										Required: true,

										Validators: []validator.String{
											stringvalidator.LengthAtLeast(attrMinLength),
										},
									},

									"url": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that represents the bookmark URL.").Description,

										Required: true,

										Validators: []validator.String{
											stringvalidator.LengthAtLeast(attrMinLength),
										},
									},
								},
							},

							Validators: []validator.Set{
								setvalidator.SizeAtMost(maximumServiceBookmarks),
							},
						},
					},
				},

				Validators: []validator.Set{
					setvalidator.SizeAtMost(maximumServices),
					setvalidator.SizeAtLeast(minimumServices),
				},
			},
		},
	}
}

// ModifyPlan
func (r *EnvironmentResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

	// Destruction plan
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state environmentResourceModel
	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(resp.Plan.Get(ctx, &plan)...)

	if !req.State.Raw.IsNull() {
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Region.IsUnknown() {

		if r.region.Region == "" {
			resp.Diagnostics.AddError(
				"Cannot determine the default region",
				"The PingOne region default value cannot be determined.  This is always a bug in the provider.  Please report this issue to the provider maintainers.",
			)
			return
		}

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("region"), types.StringValue(r.region.Region))...)
	}

	if !req.State.Raw.IsNull() && !state.Type.IsNull() && state.Type.Equal(types.StringValue(string(management.ENUMENVIRONMENTTYPE_PRODUCTION))) && !state.Type.Equal(plan.Type) {
		resp.Diagnostics.AddError(
			"Data protection notice - The environment type cannot be changed from PRODUCTION to SANDBOX",
			fmt.Sprintf("The plan for environment %[1]s is to change the environment type away from \"PRODUCTION\".  This may result in the loss of user data.  The environment cannot be changed away from a `PRODUCTION` type in the Terraform provider and must be completed as a manual activity in the admin console.", plan.Id.ValueString()),
		)
		return
	}

	var servicePlan []environmentServiceModel
	resp.Diagnostics.Append(plan.Services.ElementsAs(ctx, &servicePlan, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(servicePlan) == 0 {

		serviceDefaultMap := map[string]attr.Value{
			"type":        framework.StringToTF("SSO"),
			"console_url": types.StringNull(),
			"bookmark":    types.SetNull(types.ObjectType{AttrTypes: environmentServiceBookmarkTFObjectTypes}),
			"tags":        types.SetNull(types.StringType),
		}

		serviceDefault, d := types.SetValue(
			types.ObjectType{AttrTypes: environmentServiceTFObjectTypes},
			append(
				make([]attr.Value, 0),
				types.ObjectValueMust(environmentServiceTFObjectTypes, serviceDefaultMap),
			),
		)

		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("service"), serviceDefault)...)
	}

}

func (r *EnvironmentResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data environmentResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	resp.Diagnostics.Append(environmentServicesValidateTags(ctx, data.Services)...)
}

func (r *EnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.region = resourceConfig.Client.API.Region

	if resourceConfig.Client.GlobalOptions != nil {
		r.options = *resourceConfig.Client.GlobalOptions
	}
}

func (r *EnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state environmentResourceModel

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

	// Build the model for the API
	environment, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var environmentResponse *management.Environment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.ManagementAPIClient.EnvironmentsApi.CreateEnvironmentActiveLicense(ctx).Environment(*environment).Execute()
		},
		"CreateEnvironmentActiveLicense",
		environmentCreateCustomErrorHandler,
		sdk.DefaultCreateReadRetryable,
		&environmentResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var billOfMaterials (*management.BillOfMaterials) = nil
	if v, ok := environmentResponse.GetBillOfMaterialsOk(); ok {
		billOfMaterials = v
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(environmentResponse, billOfMaterials)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *EnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *environmentResourceModel

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
	var environmentResponse *management.Environment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.EnvironmentsApi.ReadOneEnvironment(ctx, data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.Id.ValueString(), fO, fR, fErr)
		},
		"ReadOneEnvironment",
		framework.CustomErrorResourceNotFoundWarning,
		retryEnvironmentDefault,
		&environmentResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if environmentResponse == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// The bill of materials
	var billOfMaterialsResponse *management.BillOfMaterials
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.BillOfMaterialsBOMApi.ReadOneBillOfMaterials(ctx, data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.Id.ValueString(), fO, fR, fErr)
		},
		"ReadOneBillOfMaterials",
		framework.CustomErrorResourceNotFoundWarning,
		retryEnvironmentDefault,
		&billOfMaterialsResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(environmentResponse, billOfMaterialsResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state environmentResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	environment, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If environment type has changed from SANDBOX -> PRODUCTION and vice versa we need a separate API call
	if !plan.Type.Equal(state.Type) {
		updateEnvironmentTypeRequest := *management.NewUpdateEnvironmentTypeRequest()

		updateEnvironmentTypeRequest.SetType(management.EnumEnvironmentType(plan.Type.ValueString()))
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,
			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.EnvironmentsApi.UpdateEnvironmentType(ctx, plan.Id.ValueString()).UpdateEnvironmentTypeRequest(updateEnvironmentTypeRequest).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.Id.ValueString(), fO, fR, fErr)
			},
			"UpdateEnvironmentType",
			framework.DefaultCustomError,
			nil,
			nil,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Run the API call
	var environmentResponse *management.Environment
	if !plan.Name.Equal(state.Name) ||
		!plan.Description.Equal(state.Description) ||
		!plan.LicenseId.Equal(state.LicenseId) {

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.EnvironmentsApi.UpdateEnvironment(ctx, plan.Id.ValueString()).Environment(*environment).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.Id.ValueString(), fO, fR, fErr)
			},
			"UpdateEnvironment",
			environmentCreateCustomErrorHandler,
			sdk.DefaultCreateReadRetryable,
			&environmentResponse,
		)...)

	} else {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.EnvironmentsApi.ReadOneEnvironment(ctx, plan.Id.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.Id.ValueString(), fO, fR, fErr)
			},
			"ReadOneEnvironment",
			framework.CustomErrorResourceNotFoundWarning,
			retryEnvironmentDefault,
			&environmentResponse,
		)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// The bill of materials
	var billOfMaterialsResponse *management.BillOfMaterials
	if !plan.Services.Equal(state.Services) {

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.BillOfMaterialsBOMApi.UpdateBillOfMaterials(ctx, plan.Id.ValueString()).BillOfMaterials(*environment.BillOfMaterials).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.Id.ValueString(), fO, fR, fErr)
			},
			"UpdateBillOfMaterials",
			framework.CustomErrorResourceNotFoundWarning,
			retryEnvironmentDefault,
			&billOfMaterialsResponse,
		)...)

	} else {

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.BillOfMaterialsBOMApi.ReadOneBillOfMaterials(ctx, plan.Id.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.Id.ValueString(), fO, fR, fErr)
			},
			"ReadOneBillOfMaterials",
			framework.CustomErrorResourceNotFoundWarning,
			retryEnvironmentDefault,
			&billOfMaterialsResponse,
		)...)

	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(environmentResponse, billOfMaterialsResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *EnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *environmentResourceModel

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
	deletedEnv, d := deleteEnvironment(ctx, r.Client.ManagementAPIClient, data.Id.ValueString())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if deletedEnv {
		deleteStateConf := &retry.StateChangeConf{
			Pending: []string{
				"200",
				"403",
			},
			Target: []string{
				"404",
			},
			Refresh: func() (interface{}, string, error) {
				resp, r, _ := r.Client.ManagementAPIClient.EnvironmentsApi.ReadOneEnvironment(ctx, data.Id.ValueString()).Execute()

				base := 10
				return resp, strconv.FormatInt(int64(r.StatusCode), base), nil
			},
			Timeout:                   20 * time.Minute,
			Delay:                     1 * time.Second,
			MinTimeout:                500 * time.Millisecond,
			ContinuousTargetOccurence: 2,
		}
		_, err := deleteStateConf.WaitForStateContext(ctx)
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Environment Delete Timeout",
				fmt.Sprintf("Error waiting for environment (%s) to be deleted: %s", data.Id.ValueString(), err),
			)

			return
		}
	}

}

func (r *EnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:     "environment_id",
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

func deleteEnvironment(ctx context.Context, apiClient *management.APIClient, environmentId string) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	var environmentResponse *management.Environment
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, environmentId).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentId, fO, fR, fErr)
		},
		"ReadOneEnvironment-Delete",
		framework.CustomErrorResourceNotFoundWarning,
		retryEnvironmentDefault,
		&environmentResponse,
	)...)
	if diags.HasError() {
		return false, diags
	}

	var deletedEnv bool
	// If we have a production environment, it won't destroy successfully without a switch to "SANDBOX".
	if environmentResponse.GetType() == management.ENUMENVIRONMENTTYPE_PRODUCTION {
		diags.AddWarning(
			"Data protection notice",
			"The environment being destroyed is marked as a `PRODUCTION` type, which is protected to prevent accidental data loss.  The environment has been removed from Terraform state and is no longer managed by Terraform, but has been left in place in the PingOne service.",
		)
		deletedEnv = false
	} else {
		diags.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fR, fErr := apiClient.EnvironmentsApi.DeleteEnvironment(ctx, environmentId).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentId, nil, fR, fErr)
			},
			"DeleteEnvironment",
			framework.CustomErrorResourceNotFoundWarning,
			sdk.DefaultCreateReadRetryable,
			nil,
		)...)
		deletedEnv = true
	}

	return deletedEnv, diags
}

func (p *environmentResourceModel) expand(ctx context.Context) (*management.Environment, diag.Diagnostics) {
	var diags diag.Diagnostics

	var environmentLicense management.EnvironmentLicense
	if !p.LicenseId.IsNull() {
		environmentLicense = *management.NewEnvironmentLicense(p.LicenseId.ValueString())
	}

	var region management.EnvironmentRegion
	if v := os.Getenv("PINGONE_TERRAFORM_REGION_OVERRIDE"); v != "" {
		region = management.EnvironmentRegion{
			String: &v,
		}
	} else {
		regionCode := model.FindRegionByName(p.Region.ValueString()).APICode
		region = management.EnvironmentRegion{
			EnumRegionCode: &regionCode,
		}
	}

	environment := management.NewEnvironment(
		environmentLicense,
		p.Name.ValueString(),
		region,
		management.EnumEnvironmentType(p.Type.ValueString()),
	)

	if !p.Description.IsNull() {
		environment.SetDescription(p.Description.ValueString())
	}

	diags.Append(environmentServicesValidateTags(ctx, p.Services)...)
	if diags.HasError() {
		return nil, diags
	}

	if !p.Services.IsNull() {

		var servicesPlan []environmentServiceModel
		diags.Append(p.Services.ElementsAs(ctx, &servicesPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		bomServices := make([]management.BillOfMaterialsProductsInner, 0)
		for _, v := range servicesPlan {

			service, d := v.expand(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			bomServices = append(bomServices, *service)
		}

		billOfMaterials := *management.NewBillOfMaterials(bomServices)

		if !p.Solution.IsNull() {
			billOfMaterials.SetSolutionType(management.EnumSolutionType(p.Solution.ValueString()))
		}

		environment.SetBillOfMaterials(billOfMaterials)
	}

	return environment, diags
}

func (p *environmentServiceModel) expand(ctx context.Context) (*management.BillOfMaterialsProductsInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	product, err := model.FindProductByName(p.Type.ValueString())
	if err != nil {
		diags.AddError(
			"Invalid parameter",
			fmt.Sprintf("Cannot retrieve the service from the service code: %s", err))
		return nil, diags
	}

	bomService := management.NewBillOfMaterialsProductsInner(product.APICode)

	if !p.ConsoleUrl.IsNull() {
		productBOMItemConsole := management.NewBillOfMaterialsProductsInnerConsole(p.ConsoleUrl.ValueString())

		bomService.SetConsole(*productBOMItemConsole)
	}

	if !p.Bookmarks.IsNull() {

		var servicesBookmarksPlan []environmentServiceBookmarkModel
		diags.Append(p.Bookmarks.ElementsAs(ctx, &servicesBookmarksPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		bookmarks := make([]management.BillOfMaterialsProductsInnerBookmarksInner, 0)
		for _, v := range servicesBookmarksPlan {

			bookmark, d := v.expand()
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			bookmarks = append(bookmarks, *bookmark)

		}

		bomService.SetBookmarks(bookmarks)
	}

	if !p.Tags.IsNull() {

		var servicesTagsPlan []string
		diags.Append(p.Tags.ElementsAs(ctx, &servicesTagsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		servicesTagsEnum := make([]management.EnumBillOfMaterialsProductTags, 0)
		for _, v := range servicesTagsPlan {
			servicesTagsEnum = append(servicesTagsEnum, management.EnumBillOfMaterialsProductTags(v))
		}

		bomService.SetTags(servicesTagsEnum)
	}

	return bomService, diags
}

func (p *environmentServiceBookmarkModel) expand() (*management.BillOfMaterialsProductsInnerBookmarksInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.Name.IsNull() || p.Url.IsNull() {
		diags.AddError(
			"Required parameter missing",
			"The \"name\" and \"url\" parameters are required for a service bookmark.",
		)

		return nil, diags
	}

	return management.NewBillOfMaterialsProductsInnerBookmarksInner(p.Name.ValueString(), p.Url.ValueString()), diags
}

func (p *environmentResourceModel) toState(environmentApiObject *management.Environment, servicesApiObject *management.BillOfMaterials) diag.Diagnostics {
	var diags diag.Diagnostics

	if environmentApiObject == nil || servicesApiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(environmentApiObject.GetIdOk())
	p.Name = framework.StringOkToTF(environmentApiObject.GetNameOk())
	p.Description = framework.StringOkToTF(environmentApiObject.GetDescriptionOk())
	p.Type = framework.EnumOkToTF(environmentApiObject.GetTypeOk())

	if v, ok := environmentApiObject.GetRegionOk(); ok {
		if v.EnumRegionCode != nil {
			p.Region = enumRegionCodeToTF(v.EnumRegionCode)
		}

		if v.String != nil {
			p.Region = framework.StringToTF(*v.String)
		}
	}

	if v, ok := environmentApiObject.GetLicenseOk(); ok {
		p.LicenseId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	}

	if v, ok := environmentApiObject.GetOrganizationOk(); ok {
		p.OrganizationId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	} else {
		p.OrganizationId = pingonetypes.NewResourceIDNull()
	}

	p.Solution = framework.EnumOkToTF(servicesApiObject.GetSolutionTypeOk())

	services, d := toStateEnvironmentServices(servicesApiObject.GetProducts())
	diags.Append(d...)
	p.Services = services

	return diags
}

func toStateEnvironmentServices(services []management.BillOfMaterialsProductsInner) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: environmentServiceTFObjectTypes}

	if len(services) == 0 {
		return types.SetValueMust(tfObjType, []attr.Value{}), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range services {

		service := map[string]attr.Value{}

		if c, ok := v.GetTypeOk(); ok {
			mapping, err := model.FindProductByAPICode(*c)
			if err != nil {
				diags.AddError(
					"Cannot find PingOne product/service from code",
					fmt.Sprintf("Cannot find the PingOne product/service from the provided code %s.  Please report this error to the provider maintainers.", string(*c)),
				)
				service["type"] = types.StringNull()
			} else {
				service["type"] = framework.StringToTF(mapping.ProductCode)
			}
		} else {
			service["type"] = types.StringNull()
		}

		if c, ok := v.GetConsoleOk(); ok {
			service["console_url"] = framework.StringOkToTF(c.GetHrefOk())
		} else {
			service["console_url"] = types.StringNull()
		}

		service["tags"] = framework.EnumSetOkToTF(v.GetTagsOk())

		bookmarks, d := toStateEnvironmentServicesBookmark(v.GetBookmarks())
		diags.Append(d...)
		service["bookmarks"] = bookmarks

		flattenedObj, d := types.ObjectValue(environmentServiceTFObjectTypes, service)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags

}

func toStateEnvironmentServicesBookmark(bookmarks []management.BillOfMaterialsProductsInnerBookmarksInner) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: environmentServiceBookmarkTFObjectTypes}

	if len(bookmarks) == 0 {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range bookmarks {

		bookmark := map[string]attr.Value{
			"name": framework.StringOkToTF(v.GetNameOk()),
			"url":  framework.StringOkToTF(v.GetHrefOk()),
		}

		flattenedObj, d := types.ObjectValue(environmentServiceBookmarkTFObjectTypes, bookmark)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags

}

func enumRegionCodeToTF(v *management.EnumRegionCode) basetypes.StringValue {
	if v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(model.FindRegionByAPICode(*v).Region)
	}
}

func environmentCreateCustomErrorHandler(error model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	// Invalid region
	if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
		if target, ok := details[0].GetTargetOk(); ok && *target == "region" {
			allowedRegions := make([]string, 0)
			for _, allowedRegion := range details[0].GetInnerError().AllowedValues {
				allowedRegions = append(allowedRegions, model.FindRegionByAPICode(management.EnumRegionCode(allowedRegion)).Region)
			}

			diags.AddError(
				fmt.Sprintf("Incompatible environment region for the organization tenant.  Allowed regions: %v.", allowedRegions),
				"Ensure the region parameter is correctly set.  If the region parameter is correctly set in the resource creation, please raise an issue with the provider maintainers.",
			)

			return diags
		}
	}

	// DV FF
	m, _ := regexp.MatchString("^Organization does not have Ping One DaVinci FF enabled", error.GetMessage())

	if m {
		diags.AddError(
			"The PingOne DaVinci service is not enabled in this organization tenant.",
			"To enable PingOne DaVinci, the service needs to be enabled in the organization by addition to the license or by enabling the feature flag.  Please contact your Ping customer account manager.",
		)

		return diags
	}

	return nil
}

var retryEnvironmentDefault = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

	var err error

	if p1error != nil {

		// Permissions may not have propagated by this point
		if m, err := regexp.MatchString("^The request could not be completed. You do not have access to this resource.", p1error.GetMessage()); err == nil && m {
			tflog.Warn(ctx, "Insufficient PingOne privileges detected")
			return true
		}
		if err != nil {
			tflog.Warn(ctx, "Cannot match error string for retry")
			return false
		}

	}

	return false
}

func environmentServicesValidateTags(ctx context.Context, services basetypes.SetValue) diag.Diagnostics {
	var diags diag.Diagnostics

	if !services.IsNull() && !services.IsUnknown() {

		var servicesPlan []environmentServiceModel
		diags.Append(services.ElementsAs(ctx, &servicesPlan, false)...)
		if diags.HasError() {
			return diags
		}

		if len(servicesPlan) > 0 {
			daVinciService, err := model.FindProductByAPICode(management.ENUMPRODUCTTYPE_ONE_DAVINCI)
			if err != nil {
				diags.AddAttributeError(
					path.Root("service").AtName("tags"),
					"Cannot find DaVinci product",
					"In validating the configuration, the DaVinci product could not be found.  This is always a bug in the provider.  Please report this issue to the provider maintainers.",
				)

				return diags
			}

			for _, service := range servicesPlan {
				if !service.Type.Equal(types.StringValue(daVinciService.ProductCode)) {
					if !service.Tags.IsNull() {
						diags.AddAttributeError(
							path.Root("service").AtName("tags"),
							"Invalid configuration",
							fmt.Sprintf("The `tags` parameter is only configurable where the `type` is set to `%s`.  Please unset the `tags` to an empty set or remove the `tags` parameter for the service.", daVinciService.ProductCode),
						)
					}
				}
			}
		}
	}

	return diags
}
