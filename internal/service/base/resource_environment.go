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

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/planmodifiers"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type EnvironmentResource struct {
	client      *management.APIClient
	region      model.RegionMapping
	forceDelete bool
}

type environmentResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	Type                types.String `tfsdk:"type"`
	Region              types.String `tfsdk:"region"`
	LicenseId           types.String `tfsdk:"license_id"`
	OrganizationId      types.String `tfsdk:"organization_id"`
	Solution            types.String `tfsdk:"solution"`
	DefaultPopulationId types.String `tfsdk:"default_population_id"`
	DefaultPopulation   types.List   `tfsdk:"default_population"`
	Services            types.Set    `tfsdk:"service"`
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
	environmentDefaultPopulationTFObjectTypes = map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
	}

	environmentServiceTFObjectTypes = map[string]attr.Type{
		"type":        types.StringType,
		"console_url": types.StringType,
		"bookmark":    types.SetType{},
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

	typeDescriptionFmt := fmt.Sprintf("The type of the environment to create.  Options are `%s` for a development/testing environment and `%s` for environments that require protection from deletion.", management.ENUMENVIRONMENTTYPE_SANDBOX, management.ENUMENVIRONMENTTYPE_PRODUCTION)
	typeDescription := framework.SchemaDescription{
		MarkdownDescription: typeDescriptionFmt,
		Description:         strings.ReplaceAll(typeDescriptionFmt, "`", "\""),
	}

	regionDescriptionFmt := "The region to create the environment in.  Should be consistent with the PingOne organisation region.  Valid options are `AsiaPacific` `Canada` `Europe` and `NorthAmerica`."
	regionDescription := framework.SchemaDescription{
		MarkdownDescription: regionDescriptionFmt,
		Description:         strings.ReplaceAll(regionDescriptionFmt, "`", "\""),
	}

	solutionDescriptionFmt := fmt.Sprintf("The solution context of the environment.  Leave blank for a custom, non-workforce solution context.  Valid options are `%s`, or no value for custom solution context.  Workforce solution environments are not yet supported in this provider resource, but can be fetched using the `pingone_environment` datasource.", string(management.ENUMSOLUTIONTYPE_CUSTOMER))
	solutionDescription := framework.SchemaDescription{
		MarkdownDescription: solutionDescriptionFmt,
		Description:         strings.ReplaceAll(solutionDescriptionFmt, "`", "\""),
	}

	serviceTypeDescriptionFmt := fmt.Sprintf("The service type to enable in the environment.  Valid options are `%s`.", strings.Join(model.ProductsSelectableList(), "`"))
	serviceTypeDescription := framework.SchemaDescription{
		MarkdownDescription: serviceTypeDescriptionFmt,
		Description:         strings.ReplaceAll(serviceTypeDescriptionFmt, "`", "\""),
	}

	serviceConsoleUrlDescriptionFmt := "A custom console URL to set.  Generally used with services that are deployed separately to the PingOne SaaS service, such as `PingFederate`, `PingAccess`, `PingDirectory`, `PingAuthorize` and `PingCentral`."
	serviceConsoleUrlDescription := framework.SchemaDescription{
		MarkdownDescription: serviceConsoleUrlDescriptionFmt,
		Description:         strings.ReplaceAll(serviceConsoleUrlDescriptionFmt, "`", "\""),
	}

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
				PlanModifiers: []planmodifier.String{
					planmodifiers.StringDefaultValue(framework.StringToTF(string(management.ENUMENVIRONMENTTYPE_SANDBOX))),
				},
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
				PlanModifiers: []planmodifier.String{
					planmodifiers.StringDefaultValue(framework.StringToTF(os.Getenv("PINGONE_REGION"))),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(model.RegionsAvailableList()...),
				},
			},

			"license_id": schema.StringAttribute{
				Description: "An ID of a valid license to apply to the environment.",
				Required:    true,
				Validators: []validator.String{
					verify.P1ResourceIDValidator(),
				},
			},

			"organization_id": schema.StringAttribute{
				Description: "The ID of the PingOne organization tenant to which the environment belongs.",
				Computed:    true,
			},

			"solution": schema.StringAttribute{
				Description:         solutionDescription.Description,
				MarkdownDescription: solutionDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"default_population_id": schema.StringAttribute{
				Description: "The ID of the environment's default population.",
				Computed:    true,
			},
		},

		Blocks: map[string]schema.Block{
			"default_population": schema.ListNestedBlock{
				Description: "The environment's default population.",

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the environment's default population.",
							Optional:    true,
							PlanModifiers: []planmodifier.String{
								planmodifiers.StringDefaultValue(framework.StringToTF("Default")),
							},
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(attrMinLength),
							},
						},

						"description": schema.StringAttribute{
							Description: "A description to apply to the environment's default population.",
							Optional:    true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},

			"service": schema.SetNestedBlock{
				Description: "The services to enable in the environment.",

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description:         serviceTypeDescription.Description,
							MarkdownDescription: serviceTypeDescription.MarkdownDescription,
							Optional:            true,
							PlanModifiers: []planmodifier.String{
								planmodifiers.StringDefaultValue(framework.StringToTF("SSO")),
							},
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
								setvalidator.SizeAtMost(5),
							},
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(13),
				},
			},
		},
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

func (r *EnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state environmentResourceModel

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
	environment, population, d := plan.expand(ctx, r.region.APICode)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	environmentResponse, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.EnvironmentsApi.CreateEnvironmentActiveLicense(ctx).Environment(*environment).Execute()
		},
		"CreateEnvironmentActiveLicense",
		environmentCreateCustomErrorHandler,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The bill of materials
	billOfMaterialsResponse, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.BillOfMaterialsBOMApi.ReadOneBillOfMaterials(ctx, environment.GetId()).Execute()
		},
		"ReadOneBillOfMaterials",
		framework.CustomErrorResourceNotFoundWarning,
		retryEnvironmentDefault,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(deleteEnvironment(ctx, r.client, environment.GetId(), r.forceDelete)...)
		return
	}

	// Population

	defaultPopulation, d := sso.FetchDefaultPopulation(ctx, r.client, environment.GetId())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(deleteEnvironment(ctx, r.client, environment.GetId(), r.forceDelete)...)
		return
	}

	if defaultPopulation == nil {
		resp.Diagnostics.AddError(
			"Default population not found on environment creation",
			"A default population was expected to be found in the environment after creation, but none were found.  Please report this issue to the provider maintainers.")
		resp.Diagnostics.Append(deleteEnvironment(ctx, r.client, environment.GetId(), r.forceDelete)...)
		return
	}

	populationResponse, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.PopulationsApi.UpdatePopulation(ctx, environment.GetId(), defaultPopulation.GetId()).Population(*population).Execute()
		},
		"UpdatePopulation",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(deleteEnvironment(ctx, r.client, environment.GetId(), r.forceDelete)...)
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(environmentResponse.(*management.Environment), billOfMaterialsResponse.(*management.BillOfMaterials), populationResponse.(*management.Population))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *EnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *environmentResourceModel

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
	environmentResponse, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.EnvironmentsApi.ReadOneEnvironment(ctx, data.Id.ValueString()).Execute()
		},
		"ReadOneEnvironment",
		framework.CustomErrorResourceNotFoundWarning,
		retryEnvironmentDefault,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if environmentResponse == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// The bill of materials
	billOfMaterialsResponse, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.BillOfMaterialsBOMApi.ReadOneBillOfMaterials(ctx, data.Id.ValueString()).Execute()
		},
		"ReadOneBillOfMaterials",
		framework.CustomErrorResourceNotFoundWarning,
		retryEnvironmentDefault,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The default population
	populationResponse, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.PopulationsApi.ReadOnePopulation(ctx, data.Id.ValueString(), data.DefaultPopulationId.ValueString()).Execute()
		},
		"ReadOnePopulation",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(environmentResponse.(*management.Environment), billOfMaterialsResponse.(*management.BillOfMaterials), populationResponse.(*management.Population))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state environmentResourceModel

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
	environment, population, d := plan.expand(ctx, r.region.APICode)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	environmentResponse, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.EnvironmentsApi.UpdateEnvironment(ctx, plan.Id.ValueString()).Environment(*environment).Execute()
		},
		"UpdateEnvironment",
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
	resp.Diagnostics.Append(state.toState(response.(*management.Environment))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *EnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *environmentResourceModel

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
	resp.Diagnostics.Append(deleteEnvironment(ctx, r.client, data.Id.ValueString(), r.forceDelete)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteStateConf := &sdkv2resource.StateChangeConf{
		Pending: []string{
			"200",
			"403",
		},
		Target: []string{
			"404",
		},
		Refresh: func() (interface{}, string, error) {
			resp, r, _ := r.client.EnvironmentsApi.ReadOneEnvironment(ctx, data.Id.ValueString()).Execute()

			base := 10
			return resp, strconv.FormatInt(int64(r.StatusCode), base), nil
		},
		Timeout:                   20 * time.Minute,
		Delay:                     1 * time.Second,
		MinTimeout:                500 * time.Millisecond,
		ContinuousTargetOccurence: 2,
	}
	_, err := deleteStateConf.WaitForState()
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Environment Delete Timeout",
			fmt.Sprintf("Error waiting for environment (%s) to be deleted: %s", data.Id.ValueString(), err),
		)

		return
	}

}

func (r *EnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	maxSplitLength := 2
	minSplitLength := 1
	attributes := strings.SplitN(req.ID, "/", maxSplitLength)

	if len(attributes) < minSplitLength || len(attributes) > maxSplitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id\" where the import should find and import the environment default population, or \"environment_id/population_id\" where the import should import a user specified population.", req.ID),
		)
		return
	}

	if len(attributes) == 1 {
		// find the default population and set to populationId
		defaultPopulation, d := sso.FetchDefaultPopulation(ctx, r.client, attributes[0])
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		if defaultPopulation == nil {
			resp.Diagnostics.AddError(
				"Default population not found on environment import.",
				"A default population was expected to be found in the environment during import, but none were found.  Consider using the import syntax \"environment_id/population_id\".")
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("default_population_id"), defaultPopulation.GetId())...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("default_population_id"), attributes[1])...)
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[0])...)
}

func deleteEnvironment(ctx context.Context, apiClient *management.APIClient, environmentId string, forceDelete bool) diag.Diagnostics {
	var diags diag.Diagnostics

	readResponse, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, environmentId).Execute()
		},
		"ReadOneEnvironment-Delete",
		framework.CustomErrorResourceNotFoundWarning,
		retryEnvironmentDefault,
	)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	// If we have a production environment, it won't destroy successfully without a switch to "SANDBOX".  We check our provider config for a force delete flag before we do this
	if readResponse.(*management.Environment).GetType() == management.ENUMENVIRONMENTTYPE_PRODUCTION && forceDelete {

		updateEnvironmentTypeRequest := *management.NewUpdateEnvironmentTypeRequest()
		updateEnvironmentTypeRequest.SetType("SANDBOX")
		_, d := framework.ParseResponse(
			ctx,
			func() (interface{}, *http.Response, error) {
				return apiClient.EnvironmentsApi.UpdateEnvironmentType(ctx, environmentId).UpdateEnvironmentTypeRequest(updateEnvironmentTypeRequest).Execute()
			},
			"UpdateEnvironmentType",
			framework.CustomErrorResourceNotFoundWarning,
			sdk.DefaultRetryable,
		)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

	}

	_, d = framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.EnvironmentsApi.DeleteEnvironment(ctx, environmentId).Execute()
			return nil, r, err
		},
		"DeleteEnvironment",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)

	return d
}

func (p *environmentResourceModel) expand(ctx context.Context, region management.EnumRegionCode) (*management.Environment, *management.Population, diag.Diagnostics) {
	var diags diag.Diagnostics

	var environmentLicense management.EnvironmentLicense
	if !p.LicenseId.IsNull() {
		environmentLicense = *management.NewEnvironmentLicense(p.LicenseId.ValueString())
	}

	if !p.Region.IsNull() {
		region = model.FindRegionByName(p.Region.ValueString()).APICode
	}

	environment := management.NewEnvironment(environmentLicense, p.Name.ValueString(), region, management.EnumEnvironmentType(p.Type.ValueString()))

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

	var population *management.Population

	if !p.DefaultPopulation.IsNull() {

		var populationPlan []environmentDefaultPopulationModel
		diags.Append(p.DefaultPopulation.ElementsAs(ctx, &populationPlan, false)...)
		if diags.HasError() {
			return nil, nil, diags
		}

		var d diag.Diagnostics
		population, d = populationPlan[0].expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, nil, diags
		}
	}

	return environment, population, diags
}

func (p *environmentServiceModel) expand(ctx context.Context) (*management.BillOfMaterialsProductsInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	product, err := model.FindProductByName(p.Type.ValueString())
	if err != nil {
		diags.AddError(
			"Invalid parameter",
			fmt.Sprintf("Cannot retrieve the service from the service code: %w", err))
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

			bookmark, d := v.expand(ctx)
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

func (p *environmentServiceBookmarkModel) expand(ctx context.Context) (*management.BillOfMaterialsProductsInnerBookmarksInner, diag.Diagnostics) {
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

func (p *environmentDefaultPopulationModel) expand(ctx context.Context) (*management.Population, diag.Diagnostics) {
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

	return population, diags
}

func (p *environmentResourceModel) toState(environmentApiObject *management.Environment, servicesApiObject *management.BillOfMaterials, populationApiObject *management.Population) diag.Diagnostics {
	var diags diag.Diagnostics

	if environmentApiObject == nil || servicesApiObject == nil || populationApiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(environmentApiObject.GetIdOk())
	p.Name = framework.StringOkToTF(environmentApiObject.GetNameOk())
	p.Description = framework.StringOkToTF(environmentApiObject.GetDescriptionOk())
	p.Type = enumEnvironmentTypeOkToTF(environmentApiObject.GetTypeOk())
	p.Region = enumRegionCodeOkToTF(environmentApiObject.GetRegionOk())

	if v, ok := environmentApiObject.GetLicenseOk(); ok {
		p.LicenseId = framework.StringOkToTF(v.GetIdOk())
	}

	if v, ok := environmentApiObject.GetLicenseOk(); ok {
		p.OrganizationId = framework.StringOkToTF(v.GetIdOk())
	} else {
		p.OrganizationId = types.StringNull()
	}

	p.Solution = enumSolutionTypeOkToTF(servicesApiObject.GetSolutionTypeOk())
	p.DefaultPopulationId = framework.StringOkToTF(populationApiObject.GetIdOk())

	defaultPopulation, d := toStateEnvironmentDefaultPopulation(populationApiObject)
	diags.Append(d...)
	p.DefaultPopulation = defaultPopulation

	services, d := toStateEnvironmentServices(servicesApiObject.GetProducts())
	diags.Append(d...)
	p.Services = services

	return diags
}

func toStateEnvironmentDefaultPopulation(population *management.Population) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: environmentDefaultPopulationTFObjectTypes}

	if population == nil {
		return types.ListValueMust(tfObjType, []attr.Value{}), diags
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
			service["type"] = framework.StringToTF(string(*c))
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

		flattenedObj, d := types.ObjectValue(environmentServiceTFObjectTypes, bookmark)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags

}

func enumEnvironmentTypeOkToTF(v *management.EnumEnvironmentType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}

func enumRegionCodeOkToTF(v *management.EnumRegionCode, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}

func enumSolutionTypeOkToTF(v *management.EnumSolutionType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
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
