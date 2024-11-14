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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

// Types
type EnvironmentDataSource serviceClientType

type EnvironmentDataSourceModel struct {
	Id             pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId  pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name           types.String                 `tfsdk:"name"`
	Description    types.String                 `tfsdk:"description"`
	Type           types.String                 `tfsdk:"type"`
	Region         types.String                 `tfsdk:"region"`
	LicenseId      pingonetypes.ResourceIDValue `tfsdk:"license_id"`
	OrganizationId pingonetypes.ResourceIDValue `tfsdk:"organization_id"`
	Solution       types.String                 `tfsdk:"solution"`
	Services       types.Set                    `tfsdk:"services"`
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

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the name of the environment to retrieve. Either `environment_id`, or `name` can be used to retrieve the environment, but cannot be set together.")

	environmentIdDescription := framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to retrieve. Either `environment_id`, or `name` can be used to retrieve the environment, but cannot be set together.")

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The type of the environment.  Options are `%s` for a development/testing environment and `%s` for environments that require protection from deletion.", management.ENUMENVIRONMENTTYPE_SANDBOX, management.ENUMENVIRONMENTTYPE_PRODUCTION),
	)

	regionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The region the environment is created in.",
	).AllowedValuesEnum(management.AllowedEnumRegionCodeEnumValues)

	solutionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The solution context of the environment.  Blank or null values indicate a custom, non-workforce solution context.  Valid options are `%s`, `%s` or no value for custom solution context.", string(management.ENUMSOLUTIONTYPE_CUSTOMER), string(management.ENUMSOLUTIONTYPE_WORKFORCE)),
	)

	serviceTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("The service type applied to the environment.  Valid options are `%s`.", strings.Join(model.ProductsSelectableList(), "`, `")),
	)

	serviceConsoleUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A custom console URL set for the service.  Generally used with services that are deployed separately to the PingOne SaaS service, such as `PingFederate`, `PingAccess`, `PingDirectory`, `PingAuthorize` and `PingCentral`.",
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
		fmt.Sprintf("A set of tags applied upon environment creation.  Only configurable when the service `type` is `%s`.", daVinciService.ProductCode),
	).AllowedValuesEnum(management.AllowedEnumBillOfMaterialsProductTagsEnumValues)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve details of a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": schema.StringAttribute{
				Description:         environmentIdDescription.Description,
				MarkdownDescription: environmentIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
				},
			},

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("environment_id")),
					stringvalidator.LengthAtLeast(nameLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description of the environment.").Description,
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
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of a valid license applied to the environment.").Description,
				Computed:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"organization_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the PingOne organization tenant to which the environment belongs.").Description,
				Computed:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"solution": schema.StringAttribute{
				Description:         solutionDescription.Description,
				MarkdownDescription: solutionDescription.MarkdownDescription,
				Computed:            true,
			},

			"services": schema.SetNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of objects that specify the services that are enabled in the environment.").Description,
				Computed:    true,

				NestedObject: schema.NestedAttributeObject{

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

						"tags": schema.SetAttribute{
							Description:         serviceTagsDescription.Description,
							MarkdownDescription: serviceTagsDescription.MarkdownDescription,

							ElementType: types.StringType,

							Computed: true,
						},

						"bookmarks": schema.SetNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of objects that specify custom bookmark links for the service.").Description,
							Computed:    true,

							NestedObject: schema.NestedAttributeObject{

								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the bookmark name.").Description,
										Computed:    true,
									},

									"url": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the bookmark URL.").Description,
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *EnvironmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *EnvironmentDataSourceModel

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

	var environment *management.Environment

	if !data.Name.IsNull() {

		scimFilter := fmt.Sprintf("name sw \"%s\"", data.Name.ValueString())

		// Run the API call
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.EnvironmentsApi.ReadAllEnvironments(ctx).Filter(scimFilter).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return nil, pageCursor.HTTPResponse, err
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if environments, ok := pageCursor.EntityArray.Embedded.GetEnvironmentsOk(); ok {
						for _, environmentItem := range environments {

							if environmentItem.GetName() == data.Name.ValueString() {
								return &environmentItem, initialHttpResponse, nil
							}
						}

					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllEnvironments",
			framework.DefaultCustomError,
			retryEnvironmentDefault,
			&environment,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if environment == nil {
			resp.Diagnostics.AddError(
				"Cannot find environment from name",
				fmt.Sprintf("The environment %s cannot be found", data.Name.String()),
			)
			return
		}

	} else if !data.EnvironmentId.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.EnvironmentsApi.ReadOneEnvironment(ctx, data.EnvironmentId.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneEnvironment",
			framework.DefaultCustomError,
			retryEnvironmentDefault,
			&environment,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested environment. environment_id or name must be set.",
		)
		return
	}

	// The bill of materials
	var billOfMaterialsResponse *management.BillOfMaterials
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.BillOfMaterialsBOMApi.ReadOneBillOfMaterials(ctx, environment.GetId()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environment.GetId(), fO, fR, fErr)
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
	resp.Diagnostics.Append(data.toState(environment, billOfMaterialsResponse)...)
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

	p.Id = framework.PingOneResourceIDOkToTF(environmentApiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDOkToTF(environmentApiObject.GetIdOk())
	p.Name = framework.StringOkToTF(environmentApiObject.GetNameOk())
	p.Description = framework.StringOkToTF(environmentApiObject.GetDescriptionOk())
	p.Type = framework.EnumOkToTF(environmentApiObject.GetTypeOk())
	p.Region = framework.EnumOkToTF(environmentApiObject.GetRegionOk())

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
