package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceEnvironment() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne environments.",

		CreateContext: resourcePingOneEnvironmentCreate,
		ReadContext:   resourcePingOneEnvironmentRead,
		UpdateContext: resourcePingOneEnvironmentUpdate,
		DeleteContext: resourcePingOneEnvironmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePingOneEnvironmentImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description:      "The name of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "A description of the environment.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description:      "The type of the environment to create.  Options are `SANDBOX` for a development/testing environment and `PRODUCTION` for environments that require protection from deletion.",
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "SANDBOX",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SANDBOX", "PRODUCTION"}, false)),
			},
			"region": {
				Description:      "The region to create the environment in.  Should be consistent with the PingOne organisation region.  Valid options are `AsiaPacific` `Canada` `Europe` and `NorthAmerica`.",
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc("PINGONE_REGION", nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(model.RegionsAvailableList(), false)),
				ForceNew:         true,
			},
			"license_id": {
				Description:      "An ID of a valid license to apply to the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"organization_id": {
				Description: "The ID of the PingOne organization tenant to which the environment belongs.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"solution": {
				Description:      fmt.Sprintf("The solution context of the environment.  Leave blank for a custom, non-workforce solution context.  Valid options are `%s`, or no value for custom solution context.  Workforce solution environments are not yet supported in this provider resource, but can be fetched using the `pingone_environment` datasource.", string(management.ENUMSOLUTIONTYPE_CUSTOMER)),
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMSOLUTIONTYPE_CUSTOMER)}, false)),
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
			},
			"default_population_id": {
				Description: "The ID of the environment's default population.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"default_population": {
				Description: "The environment's default population.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description:      "The name of the environment's default population.",
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "Default",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"description": {
							Description: "A description to apply to the environment's default population.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"service": {
				Description: "The services to enable in the environment.",
				Type:        schema.TypeSet,
				MaxItems:    13, // total services that exist
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:      "The service type to enable in the environment.  Valid options are `SSO`, `MFA`, `Risk`, `Verify`, `Credentials`, `APIIntelligence`, `Authorize`, `PingID`, `PingFederate`, `PingAccess`, `PingDirectory`, `PingAuthorize` and `PingCentral`.",
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(model.ProductsSelectableList(), false)),
							Optional:         true,
							Default:          "SSO",
						},
						"console_url": {
							Description: "A custom console URL to set.  Generally used with services that are deployed separately to the PingOne SaaS service, such as `PingFederate`, `PingAccess`, `PingDirectory`, `PingAuthorize` and `PingCentral`.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"bookmark": {
							Description: "Custom bookmark links for the service.",
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    5,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description:      "Bookmark name.",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"url": {
										Description:      "Bookmark URL.",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
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

func resourcePingOneEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	var diags diag.Diagnostics

	// Environment creation

	var environmentLicense management.EnvironmentLicense
	if v, ok := d.GetOk("license_id"); ok {
		environmentLicense = *management.NewEnvironmentLicense(v.(string))
	}

	region := p1Client.API.Region.APICode

	if v, ok := d.GetOk("region"); ok && v != "" {
		region = model.FindRegionByName(v.(string)).APICode
	}

	environment := *management.NewEnvironment(environmentLicense, d.Get("name").(string), region, management.EnumEnvironmentType(d.Get("type").(string))) // Environment |  (optional)

	if v, ok := d.GetOk("description"); ok {
		environment.SetDescription(v.(string))
	}

	if services, ok := d.GetOk("service"); ok {
		productBOMItems, err := expandBOMProducts(services.(*schema.Set))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error mapping configured services with the platform services``: %v", err),
				Detail:   fmt.Sprintf("Configured services: %v\n", services),
			})

			return diags
		}

		billOfMaterials := *management.NewBillOfMaterials(productBOMItems)

		if v, ok := d.GetOk("solution"); ok {
			billOfMaterials.SetSolutionType(management.EnumSolutionType(v.(string)))
		}

		environment.SetBillOfMaterials(billOfMaterials)
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.EnvironmentsApi.CreateEnvironmentActiveLicense(ctx).Environment(environment).Execute()
		},
		"CreateEnvironmentActiveLicense",
		func(error model.P1Error) diag.Diagnostics {

			// Invalid region
			if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
				if target, ok := details[0].GetTargetOk(); ok && *target == "region" {
					allowedRegions := make([]string, 0)
					for _, allowedRegion := range details[0].GetInnerError().AllowedValues {
						allowedRegions = append(allowedRegions, model.FindRegionByAPICode(management.EnumRegionCode(allowedRegion)).Region)
					}
					diags = diag.FromErr(fmt.Errorf("Incompatible environment region for the organization tenant.  Expecting regions %v, region provided: %s", allowedRegions, model.FindRegionByAPICode(region).Region))

					return diags
				}
			}

			// DV FF
			m, err := regexp.MatchString("^Organization does not have Ping One DaVinci FF enabled", error.GetMessage())
			if err != nil {
				diags = diag.FromErr(fmt.Errorf("Invalid regexp: DV FF error"))
				return diags
			}
			if m {
				diags = diag.FromErr(fmt.Errorf("The PingOne DaVinci service is not enabled in this organization tenant."))

				return diags
			}

			return nil
		},
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.Environment)

	// Set the default population
	// We have to create a default population because the API must require one population in the environment. If we don't do this we have a problem with the 'destroy all' routine

	population := *management.NewPopulation("Default") // Population |  (optional)

	if defaultPopulation, defaultPopulationOk := d.GetOk("default_population"); defaultPopulationOk {

		population.SetName(defaultPopulation.([]interface{})[0].(map[string]interface{})["name"].(string))
		description := defaultPopulation.([]interface{})[0].(map[string]interface{})["description"]
		if description != nil && description.(string) != "" {
			population.SetDescription(description.(string))
		}

	}

	populationResp, diags := sso.PingOnePopulationCreate(ctx, apiClient, respObject.GetId(), population)
	if diags.HasError() {
		return diags
	}

	d.SetId(respObject.GetId())
	d.Set("default_population_id", populationResp.GetId())

	return resourcePingOneEnvironmentRead(ctx, d, meta)
}

func resourcePingOneEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	environmentID := d.Id()
	populationID := d.Get("default_population_id").(string)

	// The environment
	resp, diags := sdk.ParseResponse(
		ctx,
		func() (interface{}, *http.Response, error) {
			return apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, environmentID).Execute()
		},
		"ReadOneEnvironment",
		sdk.CustomErrorResourceNotFoundWarning,
		retryEnvironmentDefault,
	)
	if diags.HasError() {
		return diags
	}

	if resp == nil {
		d.SetId("")
		return nil
	}
	respObject := resp.(*management.Environment)

	d.Set("name", respObject.GetName())

	if v, ok := respObject.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	d.Set("type", respObject.GetType())
	d.Set("region", model.FindRegionByAPICode(respObject.GetRegion()).Region)
	d.Set("license_id", respObject.GetLicense().Id)
	d.Set("organization_id", respObject.GetOrganization().Id)

	// The bill of materials

	servicesResp, diags := sdk.ParseResponse(
		ctx,
		func() (interface{}, *http.Response, error) {
			return apiClient.BillOfMaterialsBOMApi.ReadOneBillOfMaterials(ctx, environmentID).Execute()
		},
		"ReadOneBillOfMaterials",
		sdk.DefaultCustomError,
		retryEnvironmentDefault,
	)
	if diags.HasError() {
		return diags
	}

	bomObject := servicesResp.(*management.BillOfMaterials)

	if v, ok := bomObject.GetProductsOk(); ok {
		productBOMItems, err := flattenBOMProducts(v)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error mapping platform services with the configured services``: %v", err),
				Detail:   fmt.Sprintf("Platform services: %v\n", v),
			})

			return diags
		}

		d.Set("service", productBOMItems)
	} else {
		d.Set("service", nil)
	}

	if v, ok := bomObject.GetSolutionTypeOk(); ok {

		if *v == management.ENUMSOLUTIONTYPE_WORKFORCE {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "The configured environment has a WORKFORCE solution context.  Workforce solution context environments are not yet supported in this resource.",
			})

			return diags
		}

		d.Set("solution", string(*v))
	} else {
		d.Set("solution", nil)
	}

	// The population

	populationResp, diags := sso.PingOnePopulationRead(ctx, apiClient, environmentID, populationID)
	if diags.HasError() {
		return diags
	}

	populationConfigs := []interface{}{}

	if v, ok := populationResp.GetDescriptionOk(); ok {
		populationConfigs = append(populationConfigs, map[string]interface{}{
			"name":        populationResp.GetName(),
			"description": v,
		})
	} else {
		populationConfigs = append(populationConfigs, map[string]interface{}{
			"name":        populationResp.GetName(),
			"description": nil,
		})
	}

	d.Set("default_population", populationConfigs)

	return diags
}

func resourcePingOneEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	environmentID := d.Id()
	populationID := d.Get("default_population_id").(string)

	// The environment

	var environmentLicense management.EnvironmentLicense
	if v, ok := d.GetOk("license_id"); ok {
		environmentLicense = *management.NewEnvironmentLicense(v.(string))
	}

	region := p1Client.API.Region.APICode

	if v, ok := d.GetOk("region"); ok {
		region = model.FindRegionByName(v.(string)).APICode
	}

	environment := *management.NewEnvironment(environmentLicense, d.Get("name").(string), region, management.EnumEnvironmentType(d.Get("type").(string))) // Environment |  (optional)
	if v, ok := d.GetOk("description"); ok {
		environment.SetDescription(v.(string))
	}

	// Check if we have to change the environment type

	if change := d.HasChange("type"); change {
		//If type has changed from SANDBOX -> PRODUCTION and vice versa we need a separate API call
		updateEnvironmentTypeRequest := *management.NewUpdateEnvironmentTypeRequest()
		newType := d.Get("type")
		updateEnvironmentTypeRequest.SetType(management.EnumEnvironmentType(newType.(string)))
		_, diags = sdk.ParseResponse(
			ctx,
			func() (interface{}, *http.Response, error) {
				return apiClient.EnvironmentsApi.UpdateEnvironmentType(ctx, environmentID).UpdateEnvironmentTypeRequest(updateEnvironmentTypeRequest).Execute()
			},
			"UpdateEnvironmentType",
			sdk.DefaultCustomError,
			sdk.DefaultRetryable,
		)
		if diags.HasError() {
			return diags
		}
	}

	_, diags = sdk.ParseResponse(
		ctx,
		func() (interface{}, *http.Response, error) {
			return apiClient.EnvironmentsApi.UpdateEnvironment(ctx, environmentID).Environment(environment).Execute()
		},
		"UpdateEnvironment",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	// The bill of materials

	if services, ok := d.GetOk("service"); ok {
		productBOMItems, err := expandBOMProducts(services.(*schema.Set))

		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error mapping configured services with the platform services``: %v", err),
				Detail:   fmt.Sprintf("Configured services: %v\n", services),
			})

			return diags
		}

		billOfMaterials := *management.NewBillOfMaterials(productBOMItems)

		_, diags = sdk.ParseResponse(
			ctx,
			func() (interface{}, *http.Response, error) {
				return apiClient.BillOfMaterialsBOMApi.UpdateBillOfMaterials(ctx, environmentID).BillOfMaterials(billOfMaterials).Execute()
			},
			"UpdateBillOfMaterials",
			sdk.DefaultCustomError,
			sdk.DefaultRetryable,
		)
		if diags.HasError() {
			return diags
		}
	}

	// Default Population

	population := *management.NewPopulation("Default") // Population |  (optional)

	if defaultPopulation, defaultPopulationOk := d.GetOk("default_population"); defaultPopulationOk {

		population.SetName(defaultPopulation.([]interface{})[0].(map[string]interface{})["name"].(string))
		description := defaultPopulation.([]interface{})[0].(map[string]interface{})["description"].(string)
		if description != "" {
			population.SetDescription(description)
		}

	}

	_, diags = sso.PingOnePopulationUpdate(ctx, apiClient, environmentID, populationID, population)
	if diags.HasError() {
		return diags
	}

	return resourcePingOneEnvironmentRead(ctx, d, meta)
}

func resourcePingOneEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	// If we have a production environment, it won't destroy successfully without a switch to "SANDBOX".  We check our provider config for a force delete flag before we do this
	if d.Get("type").(string) == "PRODUCTION" && p1Client.ForceDelete {

		updateEnvironmentTypeRequest := *management.NewUpdateEnvironmentTypeRequest()
		updateEnvironmentTypeRequest.SetType("SANDBOX")
		_, diags = sdk.ParseResponse(
			ctx,
			func() (interface{}, *http.Response, error) {
				return apiClient.EnvironmentsApi.UpdateEnvironmentType(ctx, d.Id()).UpdateEnvironmentTypeRequest(updateEnvironmentTypeRequest).Execute()
			},
			"UpdateEnvironmentType",
			sdk.CustomErrorResourceNotFoundWarning,
			sdk.DefaultRetryable,
		)
		if diags.HasError() {
			return diags
		}

	}

	_, diags = sdk.ParseResponse(
		ctx,
		func() (interface{}, *http.Response, error) {
			r, err := apiClient.EnvironmentsApi.DeleteEnvironment(ctx, d.Id()).Execute()
			return nil, r, err
		},
		"DeleteEnvironment",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	deleteStateConf := &resource.StateChangeConf{
		Pending: []string{
			"200",
			"403",
		},
		Target: []string{
			"404",
		},
		Refresh: func() (interface{}, string, error) {
			resp, r, _ := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, d.Id()).Execute()

			base := 10
			return resp, strconv.FormatInt(int64(r.StatusCode), base), nil
		},
		Timeout:                   d.Timeout(schema.TimeoutDelete) - time.Minute,
		Delay:                     1 * time.Second,
		MinTimeout:                500 * time.Millisecond,
		ContinuousTargetOccurence: 2,
	}
	_, err := deleteStateConf.WaitForState()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Error waiting for environment (%s) to be deleted: %s", d.Id(), err),
		})

		return diags
	}

	return nil
}

func resourcePingOneEnvironmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/populationID\"", d.Id())
	}

	environmentID, populationID := attributes[0], attributes[1]

	d.SetId(environmentID)
	d.Set("default_population_id", populationID)

	resourcePingOneEnvironmentRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

var (
	retryEnvironmentDefault = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

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
)

func expandBOMProducts(items *schema.Set) ([]management.BillOfMaterialsProductsInner, error) {
	var productBOMItems []management.BillOfMaterialsProductsInner

	for _, item := range items.List() {

		v, err := model.FindProductByName(item.(map[string]interface{})["type"].(string))
		if err != nil {
			return nil, fmt.Errorf("Cannot retrieve the service from the service code: %w", err)
		}

		productBOM := management.NewBillOfMaterialsProductsInner(v.APICode)

		if (item.(map[string]interface{})["console_url"] != nil) && (item.(map[string]interface{})["console_url"] != "") {
			productBOMItemConsole := management.NewBillOfMaterialsProductsInnerConsole(item.(map[string]interface{})["console_url"].(string))

			productBOM.SetConsole(*productBOMItemConsole)
		}

		var productBOMBookmarkItems []management.BillOfMaterialsProductsInnerBookmarksInner

		for _, bookmarkItem := range item.(map[string]interface{})["bookmark"].(*schema.Set).List() {

			productBOMBookmark := management.NewBillOfMaterialsProductsInnerBookmarksInner(bookmarkItem.(map[string]interface{})["name"].(string), bookmarkItem.(map[string]interface{})["url"].(string))

			productBOMBookmarkItems = append(productBOMBookmarkItems, *productBOMBookmark)
		}

		productBOM.SetBookmarks(productBOMBookmarkItems)

		productBOMItems = append(productBOMItems, *productBOM)
	}

	return productBOMItems, nil
}

func flattenBOMProducts(products []management.BillOfMaterialsProductsInner) ([]interface{}, error) {
	productItems := make([]interface{}, 0)

	for _, product := range products {

		v, err := model.FindProductByAPICode(product.GetType())
		if err != nil {
			return nil, fmt.Errorf("Cannot retrieve the service from the service code: %w", err)
		}

		productItemsMap := map[string]interface{}{
			"type": v.ProductCode,
		}

		if v, ok := product.Console.GetHrefOk(); ok {
			productItemsMap["console_url"] = v
		}

		if v, ok := product.GetBookmarksOk(); ok {
			productItemsMap["bookmark"] = flattenBOMProductsBookmarkList(v)
		}

		productItems = append(productItems, productItemsMap)

	}

	return productItems, nil
}

func flattenBOMProductsBookmarkList(bookmarkList []management.BillOfMaterialsProductsInnerBookmarksInner) []interface{} {
	bookmarkItems := make([]interface{}, 0, len(bookmarkList))
	for _, bookmark := range bookmarkList {

		bookmarkName := ""
		if _, ok := bookmark.GetNameOk(); ok {
			bookmarkName = bookmark.GetName()
		}
		bookmarkHref := ""
		if _, ok := bookmark.GetHrefOk(); ok {
			bookmarkHref = bookmark.GetHref()
		}

		bookmarkItems = append(bookmarkItems, map[string]interface{}{
			"name": bookmarkName,
			"url":  bookmarkHref,
		})
	}
	return bookmarkItems
}
