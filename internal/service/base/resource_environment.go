package base

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	stringdefaultinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringdefaultinternal"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type EnvironmentResource struct {
	Client      *pingone.Client
	region      model.RegionMapping
	forceDelete bool
}

type environmentResourceModel struct {
	Id                  types.String   `tfsdk:"id"`
	Name                types.String   `tfsdk:"name"`
	Description         types.String   `tfsdk:"description"`
	Type                types.String   `tfsdk:"type"`
	Region              types.String   `tfsdk:"region"`
	LicenseId           types.String   `tfsdk:"license_id"`
	OrganizationId      types.String   `tfsdk:"organization_id"`
	Solution            types.String   `tfsdk:"solution"`
	DefaultPopulationId types.String   `tfsdk:"default_population_id"` // Deprecated
	DefaultPopulation   types.List     `tfsdk:"default_population"`    // Deprecated
	Services            types.Set      `tfsdk:"service"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}

type environmentDefaultPopulationModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type environmentServiceModel struct {
	Type       types.String `tfsdk:"type"`
	ConsoleUrl types.String `tfsdk:"console_url"`
	Bookmarks  types.Set    `tfsdk:"bookmark"`
}

type environmentServiceBookmarkModel struct {
	Name types.String `tfsdk:"name"`
	Url  types.String `tfsdk:"url"`
}

var (

	///////////////////
	// Deprecated start
	environmentDefaultPopulationTFObjectTypes = map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
	}
	// Deprecated end
	///////////////////

	environmentServiceTFObjectTypes = map[string]attr.Type{
		"type":        types.StringType,
		"console_url": types.StringType,
		"bookmark":    types.SetType{ElemType: types.ObjectType{AttrTypes: environmentServiceBookmarkTFObjectTypes}},
	}

	environmentServiceBookmarkTFObjectTypes = map[string]attr.Type{
		"name": types.StringType,
		"url":  types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &EnvironmentResource{}
	_ resource.ResourceWithConfigure   = &EnvironmentResource{}
	_ resource.ResourceWithImportState = &EnvironmentResource{}
	_ resource.ResourceWithModifyPlan  = &EnvironmentResource{}
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
		fmt.Sprintf("The type of the environment to create.  Options are `%s` for a development/testing environment and `%s` for environments that require protection from deletion. Defaults to `%s`.", management.ENUMENVIRONMENTTYPE_SANDBOX, management.ENUMENVIRONMENTTYPE_PRODUCTION, management.ENUMENVIRONMENTTYPE_SANDBOX),
	)

	regionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The region to create the environment in.  Should be consistent with the PingOne organisation region.  Valid options are `AsiaPacific` `Canada` `Europe` and `NorthAmerica`.  Default can be set with the `PINGONE_REGION` environment variable.",
	)

	solutionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The solution context of the environment.  Leave blank for a custom, non-workforce solution context.  Valid options are `%s`, or no value for custom solution context.  Workforce solution environments are not yet supported in this provider resource, but can be fetched using the `pingone_environment` datasource.", string(management.ENUMSOLUTIONTYPE_CUSTOMER)),
	)

	defaultPopulationIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"**Deprecation Message** The `default_population_id` attribute has been deprecated.  Default population functionality has moved to the `pingone_population_default` resource.  This attribute will be removed in the next major version of the provider.  The ID of the environment's default population.  This attribute is only populated when also using the `default_population` block to define a default population.",
	)

	defaultPopulationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"**Deprecation Message** The `default_population` block has been deprecated.  Default population functionality has moved to the `pingone_population_default` resource.  This attribute will be removed in the next major version of the provider.  The environment's default population.",
	)

	defaultPopulationNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the environment's default population.",
	).DefaultValue("Default")

	serviceDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The services to enable in the environment.",
	).DefaultValue("SSO")

	serviceTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The service type to enable in the environment.  Valid options are `%s`.  Defaults to `SSO`.", strings.Join(model.ProductsSelectableList(), "`, `")),
	)

	serviceConsoleUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A custom console URL to set.  Generally used with services that are deployed separately to the PingOne SaaS service, such as `PingFederate`, `PingAccess`, `PingDirectory`, `PingAuthorize` and `PingCentral`.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne environments.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"name": schema.StringAttribute{
				Description: "The name of the environment.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: "A description of the environment.",
				Optional:    true,
			},

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(string(management.ENUMENVIRONMENTTYPE_SANDBOX)),
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
				Optional:            true,
				Computed:            true,
				Default: stringdefaultinternal.StaticStringUnknownable(func() basetypes.StringValue {

					if v := os.Getenv("PINGONE_REGION"); v != "" {
						return framework.StringToTF(v)
					}

					if r.region.Region != "" {
						return types.StringValue(r.region.Region)
					}

					return types.StringUnknown()
				}()),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(model.RegionsAvailableList()...),
				},
			},

			"license_id": schema.StringAttribute{
				Description: "An ID of a valid license to apply to the environment.  Must be a valid PingOne resource ID.",
				Required:    true,
				Validators: []validator.String{
					verify.P1ResourceIDValidator(),
				},
			},

			"organization_id": schema.StringAttribute{
				Description: "The ID of the PingOne organization tenant to which the environment belongs.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"solution": schema.StringAttribute{
				Description:         solutionDescription.Description,
				MarkdownDescription: solutionDescription.MarkdownDescription,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			///////////////////
			// Deprecated start
			"default_population_id": schema.StringAttribute{
				Description:         defaultPopulationIdDescription.Description,
				MarkdownDescription: defaultPopulationIdDescription.MarkdownDescription,
				Computed:            true,
				DeprecationMessage:  "The `default_population_id` block has been deprecated.  Default population functionality has moved to the `pingone_population_default` resource.  This attribute will be removed in the next major version of the provider.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// Deprecated end
			///////////////////
		},

		Blocks: map[string]schema.Block{
			///////////////////
			// Deprecated start
			"default_population": schema.ListNestedBlock{
				Description:         defaultPopulationDescription.Description,
				MarkdownDescription: defaultPopulationDescription.MarkdownDescription,
				DeprecationMessage:  "The `default_population` block has been deprecated.  Default population functionality has moved to the `pingone_population_default` resource.  This block will be removed in the next major version of the provider.",

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         defaultPopulationNameDescription.Description,
							MarkdownDescription: defaultPopulationNameDescription.MarkdownDescription,
							// DeprecationMessage: "The `default_population.name` attribute has been deprecated.  Default population functionality has moved to the `pingone_population_default` resource.  This parameter will be removed in the next major version of the provider.",
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString("Default"),
						},

						"description": schema.StringAttribute{
							Description: "A description to apply to the environment's default population.",
							Optional:    true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.SizeAtLeast(1),
				},
			},
			// Deprecated end
			///////////////////

			"service": schema.SetNestedBlock{
				Description:         serviceDescription.Description,
				MarkdownDescription: serviceDescription.MarkdownDescription,

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description:         serviceTypeDescription.Description,
							MarkdownDescription: serviceTypeDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString("SSO"),
							Validators: []validator.String{
								stringvalidator.OneOf(model.ProductsSelectableList()...),
							},
						},

						"console_url": schema.StringAttribute{
							Description:         serviceConsoleUrlDescription.Description,
							MarkdownDescription: serviceConsoleUrlDescription.MarkdownDescription,
							Optional:            true,
						},
					},

					Blocks: map[string]schema.Block{
						"bookmark": schema.SetNestedBlock{
							Description: "Custom bookmark links for the service.",

							NestedObject: schema.NestedBlockObject{

								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Bookmark name.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(attrMinLength),
										},
									},

									"url": schema.StringAttribute{
										Description: "Bookmark URL.",
										Required:    true,
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

			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
			}),
		},
	}
}

// ModifyPlan
func (r *EnvironmentResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

	// Destruction plan
	if req.Plan.Raw.IsNull() {
		return
	}

	///////////////////
	// Deprecated start
	var plan []environmentDefaultPopulationModel
	resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root("default_population"), &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state []environmentDefaultPopulationModel
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("default_population"), &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(state) > 0 && len(plan) == 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("default_population"),
			"State change warning",
			"The default population configuration on the \"pingone_environment\" resource will be removed from state, but will not be removed from the platform.",
		)

		resp.Diagnostics.AddAttributeWarning(
			path.Root("default_population_id"),
			"State change warning",
			"The default population configuration on the \"pingone_environment\" resource will be removed from state, the \"default_population_id\" will no longer carry the default population's ID.",
		)
	}

	if len(plan) == 0 {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("default_population_id"), types.StringNull())...)
	}
	// Deprecated end
	///////////////////

	var regionPlan types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("region"), &regionPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if regionPlan.IsUnknown() {

		if r.region.Region == "" {
			resp.Diagnostics.AddError(
				"Cannot determine the default region",
				"The PingOne region default value cannot be determined.  This is always a bug in the provider.  Please report this issue to the provider maintainers.",
			)
			return
		}

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("region"), types.StringValue(r.region.Region))...)
	}

	var servicePlan []environmentServiceModel
	resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root("service"), &servicePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(servicePlan) == 0 {

		serviceDefaultMap := map[string]attr.Value{
			"type":        framework.StringToTF("SSO"),
			"console_url": types.StringNull(),
			"bookmark":    types.SetNull(types.ObjectType{AttrTypes: environmentServiceBookmarkTFObjectTypes}),
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
}

func (r *EnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state environmentResourceModel

	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	defaultTimeout := 20 * time.Minute
	createTimeout, d := plan.Timeouts.Create(ctx, defaultTimeout)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	environment, population, d := plan.expand(ctx)
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

	///////////////////
	// Deprecated start
	var populationResponse *management.Population = nil
	if population != nil {
		// Seed a default population.  The platform does this implicitly but we see latencies.  This ensures we have a quick environment provision.
		defaultPopulationObj := *management.NewPopulation("Default")
		defaultPopulationObj.SetDescription("Automatically created population.")
		defaultPopulationObj.SetDefault(true)

	defaultPopulationResponse, _ := sso.PingOnePopulationCreate(ctx, r.Client.ManagementAPIClient, environmentResponse.GetId(), defaultPopulationObj)
	if defaultPopulationResponse == nil {
		resp.Diagnostics.AddWarning(
			"Cannot seed the default population",
			"The default population cannot be seeded explicitly by the provider.  Relying on the implicit environment bootstrapping service for creation.",
		)
	}

	///////////////////
	// Deprecated start
	// Population

		var defaultPopulation *management.Population

	// Save some processing if we already seeded our population
	if defaultPopulationResponse != nil {
		defaultPopulation = defaultPopulationResponse
	} else {
		defaultPopulation, d = sso.FetchDefaultPopulationWithTimeout(ctx, r.Client.ManagementAPIClient, environmentResponse.GetId(), createTimeout)
		resp.Diagnostics.Append(d...)
	}

	var populationResponse *management.Population = nil

		if population != nil {

			if defaultPopulation == nil {
				resp.Diagnostics.AddError(
					"Default population not found.",
					"A default population was expected to be found in the environment after creation, but none was found.  Please report this issue to the provider maintainers.")
			}

			resp.Diagnostics.Append(framework.ParseResponse(
				ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.UpdatePopulation(ctx, environmentResponse.GetId(), defaultPopulation.GetId()).Population(*population).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environmentResponse.GetId(), fO, fR, fErr)
			},
			"UpdatePopulation",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&populationResponse,
		)...)
	}
	// Deprecated end
	///////////////////

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(environmentResponse, billOfMaterials, populationResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *EnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *environmentResourceModel

	if r.Client.ManagementAPIClient == nil {
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

	///////////////////
	// Deprecated start
	// The default population
	var populationResponse *management.Population = nil
	if !data.DefaultPopulationId.IsNull() {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.ReadOnePopulation(ctx, data.Id.ValueString(), data.DefaultPopulationId.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.Id.ValueString(), fO, fR, fErr)
			},
			"ReadOnePopulation",
			framework.CustomErrorResourceNotFoundWarning,
			sdk.DefaultCreateReadRetryable,
			&populationResponse,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	// Deprecated end
	///////////////////

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(environmentResponse, billOfMaterialsResponse, populationResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state environmentResourceModel

	if r.Client == nil {
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
	environment, population, d := plan.expand(ctx)
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

	///////////////////
	// Deprecated start
	var populationResponse *management.Population = nil

	if plan.DefaultPopulation.IsNull() && !state.DefaultPopulation.IsNull() && population == nil {
		resp.Diagnostics.AddWarning(
			"Default population removed from state",
			"The default population has been removed from state, but has not been removed from the platform.  Please use the `pingone_population_default` resource to manage the default population going forward.",
		)
	}

	if !plan.DefaultPopulation.Equal(state.DefaultPopulation) && population != nil {

		var populationId string
		if state.DefaultPopulationId.IsNull() {
			defaultPopulation, d := sso.FetchDefaultPopulation(ctx, r.Client.ManagementAPIClient, plan.Id.ValueString())
			resp.Diagnostics.Append(d...)

			if defaultPopulation == nil {
				resp.Diagnostics.AddError(
					"Default population not found.",
					"A default population was expected to be found in the environment after update, but none was found.  Please report this issue to the provider maintainers.")
				return
			}

			populationId = defaultPopulation.GetId()
		} else {
			populationId = state.DefaultPopulationId.ValueString()
		}

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.UpdatePopulation(ctx, plan.Id.ValueString(), populationId).Population(*population).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.Id.ValueString(), fO, fR, fErr)
			},
			"UpdatePopulation",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&populationResponse,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if populationResponse == nil && population != nil && !state.DefaultPopulationId.IsNull() {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.ReadOnePopulation(ctx, state.Id.ValueString(), state.DefaultPopulationId.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.Id.ValueString(), fO, fR, fErr)
			},
			"ReadOnePopulation",
			framework.CustomErrorResourceNotFoundWarning,
			sdk.DefaultCreateReadRetryable,
			&populationResponse,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	// Deprecated end
	///////////////////

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(environmentResponse, billOfMaterialsResponse, populationResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *EnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *environmentResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	resp.Diagnostics.Append(deleteEnvironment(ctx, r.Client.ManagementAPIClient, data.Id.ValueString(), r.forceDelete)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

func (r *EnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	compiledRegexString := fmt.Sprintf(`%s|^%s\/%s$`, verify.P1ResourceIDRegexpFullString.String(), verify.P1ResourceIDRegexp.String(), verify.P1ResourceIDRegexp.String())
	m, err := regexp.MatchString(compiledRegexString, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Cannot verify import ID regex: %s", err),
		)
		return
	}

	if !m {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Invalid import ID specified (\"%s\").  The ID should be in the format \"environment_id\" or the deprecated form \"environment_id/population_id\" and must match regex: %s", req.ID, compiledRegexString),
		)
		return
	}

	maxSplitLength := 2 // deprecated
	minSplitLength := 1
	attributes := strings.SplitN(req.ID, "/", maxSplitLength)

	if len(attributes) < minSplitLength || len(attributes) > maxSplitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id\" where the import should find and import the population marked as Default in the environment, or \"environment_id/population_id\" where the user specifies a population as the environment default.", req.ID),
		)
		return
	}

	///////////////////
	// Deprecated start
	if len(attributes) == 2 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("default_population_id"), attributes[1])...)
	}

	if len(attributes) == 1 {
		population, d := sso.FetchDefaultPopulation(ctx, r.Client.ManagementAPIClient, attributes[0])
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		if population == nil {
			resp.Diagnostics.AddError(
				"Default population not found",
				"The Default population is not found in the environment.  Either ensure a population is configured to be the default, or you use the \"environment_id/population_id\" import ID pattern.",
			)
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("default_population_id"), population.GetId())...)
	}

	// Deprecated end
	///////////////////

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[0])...)
}

func deleteEnvironment(ctx context.Context, apiClient *management.APIClient, environmentId string, forceDelete bool) diag.Diagnostics {
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
		return diags
	}

	// If we have a production environment, it won't destroy successfully without a switch to "SANDBOX".  We check our provider config for a force delete flag before we do this
	if environmentResponse.GetType() == management.ENUMENVIRONMENTTYPE_PRODUCTION && forceDelete {

		updateEnvironmentTypeRequest := *management.NewUpdateEnvironmentTypeRequest()
		updateEnvironmentTypeRequest.SetType("SANDBOX")
		diags.Append(framework.ParseResponse(
			ctx,
			func() (any, *http.Response, error) {
				fO, fR, fErr := apiClient.EnvironmentsApi.UpdateEnvironmentType(ctx, environmentId).UpdateEnvironmentTypeRequest(updateEnvironmentTypeRequest).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentId, fO, fR, fErr)
			},
			"UpdateEnvironmentType",
			framework.CustomErrorResourceNotFoundWarning,
			nil,
			nil,
		)...)
		if diags.HasError() {
			return diags
		}

	}

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

	return diags
}

func (p *environmentResourceModel) expand(ctx context.Context) (*management.Environment, *management.Population, diag.Diagnostics) {
	var diags diag.Diagnostics

	var environmentLicense management.EnvironmentLicense
	if !p.LicenseId.IsNull() {
		environmentLicense = *management.NewEnvironmentLicense(p.LicenseId.ValueString())
	}

	environment := management.NewEnvironment(environmentLicense, p.Name.ValueString(), model.FindRegionByName(p.Region.ValueString()).APICode, management.EnumEnvironmentType(p.Type.ValueString()))

	if !p.Description.IsNull() {
		environment.SetDescription(p.Description.ValueString())
	}

	if !p.Services.IsNull() {

		var servicesPlan []environmentServiceModel
		diags.Append(p.Services.ElementsAs(ctx, &servicesPlan, false)...)
		if diags.HasError() {
			return nil, nil, diags
		}

		bomServices := make([]management.BillOfMaterialsProductsInner, 0)
		for _, v := range servicesPlan {

			service, d := v.expand(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, nil, diags
			}

			bomServices = append(bomServices, *service)
		}

		billOfMaterials := *management.NewBillOfMaterials(bomServices)

		if !p.Solution.IsNull() {
			billOfMaterials.SetSolutionType(management.EnumSolutionType(p.Solution.ValueString()))
		}

		environment.SetBillOfMaterials(billOfMaterials)
	}

	///////////////////
	// Deprecated start
	var population *management.Population = nil

	if !p.DefaultPopulation.IsNull() {

		var populationPlan []environmentDefaultPopulationModel
		diags.Append(p.DefaultPopulation.ElementsAs(ctx, &populationPlan, false)...)
		if diags.HasError() {
			return nil, nil, diags
		}

		var d diag.Diagnostics
		population, d = populationPlan[0].expand()
		diags.Append(d...)
		if diags.HasError() {
			return nil, nil, diags
		}
	}
	// Deprecated end
	///////////////////

	return environment, population, diags
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

// expand extends the environmentDefaultPopulationModel, which returns a *management.Population pointer object of the model.
//
// Deprecated: default population configuration is replaced by a separate TF resource, `pingone_population_default`
func (p *environmentDefaultPopulationModel) expand() (*management.Population, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.Name.IsNull() {
		diags.AddError(
			"Required parameter missing",
			"The \"name\" parameters is required for a default population.",
		)

		return nil, diags
	}

	population := management.NewPopulation(p.Name.ValueString())

	if !p.Description.IsNull() {
		population.SetDescription(p.Description.ValueString())
	}

	population.SetDefault(true)

	return population, diags
}

func (p *environmentResourceModel) toState(environmentApiObject *management.Environment, servicesApiObject *management.BillOfMaterials, populationApiObject *management.Population) diag.Diagnostics {
	var diags diag.Diagnostics

	if environmentApiObject == nil || servicesApiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(environmentApiObject.GetIdOk())
	p.Name = framework.StringOkToTF(environmentApiObject.GetNameOk())
	p.Description = framework.StringOkToTF(environmentApiObject.GetDescriptionOk())
	p.Type = framework.EnumOkToTF(environmentApiObject.GetTypeOk())
	p.Region = enumRegionCodeOkToTF(environmentApiObject.GetRegionOk())

	if v, ok := environmentApiObject.GetLicenseOk(); ok {
		p.LicenseId = framework.StringOkToTF(v.GetIdOk())
	}

	if v, ok := environmentApiObject.GetOrganizationOk(); ok {
		p.OrganizationId = framework.StringOkToTF(v.GetIdOk())
	} else {
		p.OrganizationId = types.StringNull()
	}

	p.Solution = framework.EnumOkToTF(servicesApiObject.GetSolutionTypeOk())

	services, d := toStateEnvironmentServices(servicesApiObject.GetProducts())
	diags.Append(d...)
	p.Services = services

	///////////////////
	// Deprecated start
	if populationApiObject != nil {
		p.DefaultPopulationId = framework.StringOkToTF(populationApiObject.GetIdOk())

		defaultPopulation, d := toStateEnvironmentDefaultPopulation(populationApiObject)
		diags.Append(d...)
		p.DefaultPopulation = defaultPopulation

	} else {
		p.DefaultPopulationId = types.StringNull()
		p.DefaultPopulation = types.ListNull(types.ObjectType{AttrTypes: environmentDefaultPopulationTFObjectTypes})
	}
	// Deprecated end
	///////////////////

	return diags
}

// toStateEnvironmentDefaultPopulation takes a population object and converts it to a TF state object
//
// Deprecated: default population configuration is replaced by a separate TF resource, `pingone_population_default`
func toStateEnvironmentDefaultPopulation(population *management.Population) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: environmentDefaultPopulationTFObjectTypes}

	if population == nil {
		return types.ListNull(types.ObjectType{AttrTypes: environmentDefaultPopulationTFObjectTypes}), diags
	}

	defaultPopulation := map[string]attr.Value{
		"name":        framework.StringOkToTF(population.GetNameOk()),
		"description": framework.StringOkToTF(population.GetDescriptionOk()),
	}

	flattenedObj, d := types.ObjectValue(environmentDefaultPopulationTFObjectTypes, defaultPopulation)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags

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

		bookmarks, d := toStateEnvironmentServicesBookmark(v.GetBookmarks())
		diags.Append(d...)
		service["bookmark"] = bookmarks

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
		return types.SetValueMust(tfObjType, []attr.Value{}), diags
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

func enumRegionCodeOkToTF(v *management.EnumRegionCode, ok bool) basetypes.StringValue {
	if !ok || v == nil {
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
