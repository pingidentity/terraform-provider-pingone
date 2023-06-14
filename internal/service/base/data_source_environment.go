package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type EnvironmentDataSource struct {
	client *management.APIClient
	region model.RegionMapping
}

type EnvironmentDataSourceModel struct {
	Id             types.String `tfsdk:"id"`
	EnvironmentId  types.String `tfsdk:"environment_id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Type           types.String `tfsdk:"type"`
	Region         types.String `tfsdk:"region"`
	LicenseId      types.String `tfsdk:"license_id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	Solution       types.String `tfsdk:"solution"`
	Services       types.Set    `tfsdk:"service"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &EnvironmentDataSource{}
)

// New Object
func NewEnvironmentDataSource() datasource.DataSource {
	return &EnvironmentDataSource{}
}

// Metadata
func (r *EnvironmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

// Schema
func (r *EnvironmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	nameLength := 1

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The type of the environment.  Options are `%s` for a development/testing environment and `%s` for environments that require protection from deletion.", management.ENUMENVIRONMENTTYPE_SANDBOX, management.ENUMENVIRONMENTTYPE_PRODUCTION),
	)

	regionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The region the environment is created in.  Valid options are `AsiaPacific` `Canada` `Europe` and `NorthAmerica`.",
	)

	solutionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The solution context of the environment.  Blank or null values indicate a custom, non-workforce solution context.  Valid options are `%s`, `%s` or no value for custom solution context.", string(management.ENUMSOLUTIONTYPE_CUSTOMER), string(management.ENUMSOLUTIONTYPE_WORKFORCE)),
	)

	serviceTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The service type applied to the environment.  Valid options are `%s`.", strings.Join(model.ProductsSelectableList(), "`, `")),
	)

	serviceConsoleUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A custom console URL set for the service.  Generally used with services that are deployed separately to the PingOne SaaS service, such as `PingFederate`, `PingAccess`, `PingDirectory`, `PingAuthorize` and `PingCentral`.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve details of a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": schema.StringAttribute{
				Description: "The ID of the environment to retrieve. Either `environment_id`, or `name` can be used to retrieve the environment, but cannot be set together.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
					verify.P1ResourceIDValidator(),
				},
			},

			"name": schema.StringAttribute{
				Description: "A string that specifies the name of the environment to retrieve. Either `environment_id`, or `name` can be used to retrieve the environment, but cannot be set together.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("environment_id")),
					stringvalidator.LengthAtLeast(nameLength),
				},
			},

			"description": schema.StringAttribute{
				Description: "A string that specifies the description of the environment.",
				Computed:    true,
			},

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Computed:            true,
			},

			"region": schema.StringAttribute{
				Description:         regionDescription.Description,
				MarkdownDescription: regionDescription.MarkdownDescription,
				Computed:            true,
			},

			"license_id": schema.StringAttribute{
				Description: "An ID of a valid license applied to the environment.",
				Computed:    true,
			},

			"organization_id": schema.StringAttribute{
				Description: "The ID of the PingOne organization tenant to which the environment belongs.",
				Computed:    true,
			},

			"solution": schema.StringAttribute{
				Description:         solutionDescription.Description,
				MarkdownDescription: solutionDescription.MarkdownDescription,
				Computed:            true,
			},
		},

		Blocks: map[string]schema.Block{
			"service": schema.SetNestedBlock{
				Description: "The services that are enabled in the environment.",

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description:         serviceTypeDescription.Description,
							MarkdownDescription: serviceTypeDescription.MarkdownDescription,
							Computed:            true,
						},

						"console_url": schema.StringAttribute{
							Description:         serviceConsoleUrlDescription.Description,
							MarkdownDescription: serviceConsoleUrlDescription.MarkdownDescription,
							Computed:            true,
						},
					},

					Blocks: map[string]schema.Block{
						"bookmark": schema.SetNestedBlock{
							Description: "Custom bookmark links for the service.",

							NestedObject: schema.NestedBlockObject{

								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Bookmark name.",
										Computed:    true,
									},

									"url": schema.StringAttribute{
										Description: "Bookmark URL.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *EnvironmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *EnvironmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *EnvironmentDataSourceModel

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
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var environment management.Environment

	if !data.Name.IsNull() {

		// Run the API call
		response, diags := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.EnvironmentsApi.ReadAllEnvironments(ctx).Execute()
			},
			"ReadAllEnvironments",
			framework.DefaultCustomError,
			retryEnvironmentDefault,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		entityArray := response.(*management.EntityArray)

		if environments, ok := entityArray.Embedded.GetEnvironmentsOk(); ok {

			found := false
			for _, environmentItem := range environments {

				if environmentItem.GetName() == data.Name.ValueString() {
					environment = environmentItem
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Cannot find environment from name",
					fmt.Sprintf("The environment %s cannot be found", data.Name.String()),
				)
				return
			}

		}

	} else if !data.EnvironmentId.IsNull() {

		// Run the API call
		response, diags := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.EnvironmentsApi.ReadOneEnvironment(ctx, data.EnvironmentId.ValueString()).Execute()
			},
			"ReadOneEnvironment",
			framework.DefaultCustomError,
			retryEnvironmentDefault,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		environment = *response.(*management.Environment)
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested environment. environment_id or name must be set.",
		)
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
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&environment, billOfMaterialsResponse.(*management.BillOfMaterials))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *EnvironmentDataSourceModel) toState(environmentApiObject *management.Environment, servicesApiObject *management.BillOfMaterials) diag.Diagnostics {
	var diags diag.Diagnostics

	if environmentApiObject == nil || servicesApiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(environmentApiObject.GetIdOk())
	p.EnvironmentId = framework.StringOkToTF(environmentApiObject.GetIdOk())
	p.Name = framework.StringOkToTF(environmentApiObject.GetNameOk())
	p.Description = framework.StringOkToTF(environmentApiObject.GetDescriptionOk())
	p.Type = framework.EnumOkToTF(environmentApiObject.GetTypeOk())
	p.Region = framework.EnumOkToTF(environmentApiObject.GetRegionOk())

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

	return diags
}
