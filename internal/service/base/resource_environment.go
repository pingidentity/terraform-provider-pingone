package base

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	pingone "github.com/patrickcping/pingone-go/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/sso"
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PRODUCTION", "SANDBOX"}, false)),
			},
			"region": {
				Description:      "The region to create the environment in.  Should be consistent with the PingOne organisation region.  Valid options are `NA`, `EU`, `ASIA` and `CA`.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NA", "EU", "ASIA", "CA"}, false)),
				ForceNew:         true,
			},
			"license_id": {
				Description:      "An ID of a valid license to apply to the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				ForceNew:         true,
			},
			// "solution": {
			// 	Description:  "The solution context of the environment.  Leave blank for a custom, non-workforce solution context.  Valid options are `WORKFORCE` and `CUSTOMER`",
			// 	Type:         schema.TypeString,
			// 	ValidateDiagFunc: validation.StringInSlice([]string{"WORKFORCE", "CUSTOMER"}, false),
			// 	Optional:     true,
			// 	ForceNew:     true,
			// },
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
				Type:        schema.TypeList,
				MaxItems:    13, // total services that exist
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:      "The service type to enable in the environment.  Valid options are `SSO`, `MFA`, `Risk`, `Verify`, `Credentials`, `APIIntelligence`, `Authorize`, `Fraud`, `PingID`, `PingFederate`, `PingAccess`, `PingDirectory`, `PingAuthorize` and `PingCentral`.",
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{`SSO`, `MFA`, `Risk`, `Verify`, `Credentials`, `APIIntelligence`, `Authorize`, `Fraud`, `PingID`, `PingFederate`, `PingAccess`, `PingDirectory`, `PingAuthorize`, `PingCentral`}, false)),
							Optional:         true,
							Default:          "SSO",
						},
						"console_url": {
							Description: "A custom console URL to set.  Generally used with services that are deployed separately to the PingOne SaaS service, such as `PING_FEDERATE`, `PING_ACCESS`, `PING_DIRECTORY`, `PING_AUTHORIZE` and `PING_CENTRAL`.",
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
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})

	var diags diag.Diagnostics

	// Environment creation

	var environmentLicense pingone.EnvironmentLicense
	if license, ok := d.GetOk("license_id"); ok {
		environmentLicense = *pingone.NewEnvironmentLicense(license.(string))
	}

	environment := *pingone.NewEnvironment(environmentLicense, d.Get("name").(string), d.Get("region").(string), d.Get("type").(string)) // Environment |  (optional)

	if v, ok := d.GetOk("description"); ok {
		environment.SetDescription(v.(string))
	}

	resp, r, err := apiClient.EnvironmentsApi.CreateEnvironmentActiveLicense(ctx).Environment(environment).Execute()
	if (err != nil) || (r.StatusCode != 201) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `EnvironmentsApi.CreateEnvironmentActiveLicense``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	// Set the Bill of Materials (the services)

	if services, ok := d.GetOk("service"); ok {
		productBOMItems, err := buildBOMProductsCreateRequest(services.([]interface{}))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error mapping configured services with the platform services``: %v", err),
				Detail:   fmt.Sprintf("Configured services: %v\n", services),
			})

			return diags
		}

		billOfMaterials := *pingone.NewBillOfMaterials(productBOMItems)

		// if solution, ok := d.GetOk("solution"); ok {
		// 	billOfMaterials.SetSolutionType(solution.(string))
		// }

		_, r, err := apiClient.BillOfMaterialsBOMApi.UpdateBillOfMaterials(ctx, resp.GetId()).BillOfMaterials(billOfMaterials).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `BillOfMaterialsBOMApi.UpdateBillOfMaterials``: %v", err),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})

			return diags
		}
	}

	// Set the default population
	// We have to create a default population because the API must require one population in the environment. If we don't do this we have a problem with the 'destroy all' routine

	population := *pingone.NewPopulation("Default") // Population |  (optional)

	if defaultPopulation, defaultPopulationOk := d.GetOk("default_population"); defaultPopulationOk {

		tflog.Debug(ctx, fmt.Sprintf("defaultPopulation: %#v", defaultPopulation))

		population.SetName(defaultPopulation.([]interface{})[0].(map[string]interface{})["name"].(string))
		description := defaultPopulation.([]interface{})[0].(map[string]interface{})["description"].(string)
		if description != "" {
			population.SetDescription(description)
		}

	}

	populationResp, _, err := sso.PingOnePopulationCreate(ctx, apiClient, resp.GetId(), population)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.GetId())
	d.Set("default_population_id", populationResp.GetId())

	return resourcePingOneEnvironmentRead(ctx, d, meta)
}

func resourcePingOneEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	environmentID := d.Id()
	populationID := d.Get("default_population_id").(string)

	// The environment

	resp, r, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, environmentID).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Environment no %s longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `EnvironmentsApi.ReadOneEnvironment``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.Set("name", resp.GetName())
	d.Set("description", resp.GetDescription())
	d.Set("type", resp.GetType())
	d.Set("region", resp.GetRegion())
	d.Set("license_id", resp.GetLicense().Id)

	// The bill of materials

	servicesResp, servicesR, servicesErr := apiClient.BillOfMaterialsBOMApi.ReadOneBillOfMaterials(ctx, environmentID).Execute()
	if servicesErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `EnvironmentsApi.ReadOneEnvironment``: %v", servicesErr),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", servicesR),
		})

		return diags
	}

	// d.Set("solution", servicesResp.SolutionType)
	productBOMItems, err := flattenBOMProducts(servicesResp)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error mapping platform services with the configured services``: %v", err),
			Detail:   fmt.Sprintf("Platform services: %v\n", servicesResp),
		})

		return diags
	}

	d.Set("service", productBOMItems)

	// The population

	populationResp, populationR, populationErr := sso.PingOnePopulationRead(ctx, apiClient, environmentID, populationID)
	if populationErr != nil {

		if populationR.StatusCode == 404 {
			log.Printf("[INFO] PingOne Application Default Population no %s longer exists", populationID)
			d.Set("default_population_id", "")
			return diags
		}

		return diag.FromErr(populationErr)
	}

	populationConfigs := []interface{}{}
	populationConfigs = append(populationConfigs, map[string]interface{}{
		"name":        populationResp.GetName(),
		"description": populationResp.GetDescription(),
	})

	tflog.Debug(ctx, fmt.Sprintf("populationConfigs: %#v", populationConfigs))

	d.Set("default_population", populationConfigs)

	return diags
}

func resourcePingOneEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	environmentID := d.Id()
	populationID := d.Get("default_population_id").(string)

	// The environment

	var environmentLicense pingone.EnvironmentLicense
	if v, ok := d.GetOk("license_id"); ok {
		environmentLicense = *pingone.NewEnvironmentLicense(v.(string))
	}

	environment := *pingone.NewEnvironment(environmentLicense, d.Get("name").(string), d.Get("region").(string), d.Get("type").(string)) // Environment |  (optional)
	if v, ok := d.GetOk("description"); ok {
		environment.SetDescription(v.(string))
	} else {
		environment.SetDescription("")
	}

	// Check if we have to change the environment type

	if change := d.HasChange("type"); change {
		//If type has changed from SANDBOX -> PRODUCTION and vice versa we need a separate API call
		updateEnvironmentTypeRequest := *pingone.NewUpdateEnvironmentTypeRequest()
		newType := d.Get("type")
		updateEnvironmentTypeRequest.SetType(newType.(string))
		_, r, err := apiClient.EnvironmentsApi.UpdateEnvironmentType(ctx, environmentID).UpdateEnvironmentTypeRequest(updateEnvironmentTypeRequest).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `EnvironmentsApi.UpdateEnvironmentType``: %v", err),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})

			return diags
		}
	}

	_, r, err := apiClient.EnvironmentsApi.UpdateEnvironment(ctx, environmentID).Environment(environment).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `EnvironmentsApi.UpdateEnvironment``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	// The bill of materials

	if services, ok := d.GetOk("service"); ok {
		productBOMItems, err := buildBOMProductsCreateRequest(services.([]interface{}))

		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error mapping configured services with the platform services``: %v", err),
				Detail:   fmt.Sprintf("Configured services: %v\n", services),
			})

			return diags
		}

		billOfMaterials := *pingone.NewBillOfMaterials(productBOMItems)

		// if solution, ok := d.GetOk("solution"); ok {
		// 	billOfMaterials.SetSolutionType(solution.(string))
		// }

		_, r, err := apiClient.BillOfMaterialsBOMApi.UpdateBillOfMaterials(ctx, environmentID).BillOfMaterials(billOfMaterials).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `BillOfMaterialsBOMApi.UpdateBillOfMaterials``: %v", err),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})

			return diags
		}
	}

	// Default Population

	population := *pingone.NewPopulation("Default") // Population |  (optional)

	if defaultPopulation, defaultPopulationOk := d.GetOk("default_population"); defaultPopulationOk {

		population.SetName(defaultPopulation.([]interface{})[0].(map[string]interface{})["name"].(string))
		description := defaultPopulation.([]interface{})[0].(map[string]interface{})["description"].(string)
		if description != "" {
			population.SetDescription(description)
		}

	}

	_, _, populationErr := sso.PingOnePopulationUpdate(ctx, apiClient, environmentID, populationID, population)
	if populationErr != nil {
		return diag.FromErr(populationErr)
	}

	return resourcePingOneEnvironmentRead(ctx, d, meta)
}

func resourcePingOneEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	// If we have a production environment, it won't destroy successfully without a switch to "SANDBOX".  We check our provider config for a force delete flag before we do this
	if d.Get("type").(string) == "PRODUCTION" && p1Client.ForceDelete {

		updateEnvironmentTypeRequest := *pingone.NewUpdateEnvironmentTypeRequest()
		updateEnvironmentTypeRequest.SetType("SANDBOX")
		_, r, err := apiClient.EnvironmentsApi.UpdateEnvironmentType(ctx, d.Id()).UpdateEnvironmentTypeRequest(updateEnvironmentTypeRequest).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `EnvironmentsApi.UpdateEnvironmentType``: %v", err),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})

			return diags
		}

	}

	_, err := apiClient.EnvironmentsApi.DeleteEnvironment(ctx, d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `EnvironmentsApi.DeleteEnvironment``: %v", err),
		})

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
			resp, r, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, d.Id()).Execute()
			if err != nil {
				return 0, "", err
			}
			base := 10
			return resp, strconv.FormatInt(int64(r.StatusCode), base), nil
		},
		Timeout:                   d.Timeout(schema.TimeoutDelete) - time.Minute,
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 5,
	}
	_, err = deleteStateConf.WaitForState()
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
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/populationID\"", d.Id())
	}

	environmentID, populationID := attributes[0], attributes[1]

	d.SetId(environmentID)
	d.Set("default_population_id", populationID)

	resourcePingOneEnvironmentRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func buildBOMProductsCreateRequest(items []interface{}) ([]pingone.BillOfMaterialsProductsInner, error) {
	var productBOMItems []pingone.BillOfMaterialsProductsInner

	for _, item := range items {

		v, err := service.ServiceFromProviderCode(item.(map[string]interface{})["type"].(string))
		if err != nil {
			return nil, fmt.Errorf("Cannot retrieve the service from the service code: %w", err)
		}

		productBOM := pingone.NewBillOfMaterialsProductsInner(v.PlatformCode)

		if (item.(map[string]interface{})["console_url"] != nil) && (item.(map[string]interface{})["console_url"] != "") {
			productBOMItemConsole := pingone.NewBillOfMaterialsProductsInnerConsole(item.(map[string]interface{})["console_url"].(string))

			productBOM.SetConsole(*productBOMItemConsole)
		}

		var productBOMBookmarkItems []pingone.BillOfMaterialsProductsInnerBookmarksInner

		for _, bookmarkItem := range item.(map[string]interface{})["bookmark"].(*schema.Set).List() {

			productBOMBookmark := pingone.NewBillOfMaterialsProductsInnerBookmarksInner(bookmarkItem.(map[string]interface{})["name"].(string), bookmarkItem.(map[string]interface{})["url"].(string))

			productBOMBookmarkItems = append(productBOMBookmarkItems, *productBOMBookmark)
		}

		productBOM.SetBookmarks(productBOMBookmarkItems)

		productBOMItems = append(productBOMItems, *productBOM)
	}

	return productBOMItems, nil
}

func flattenBOMProducts(items *pingone.BillOfMaterials) ([]interface{}, error) {
	productItems := make([]interface{}, 0)

	if _, ok := items.GetProductsOk(); ok {

		for _, product := range items.GetProducts() {

			v, err := service.ServiceFromPlatformCode(product.GetType())
			if err != nil {
				return nil, fmt.Errorf("Cannot retrieve the service from the service code: %w", err)
			}

			productItems = append(productItems, map[string]interface{}{
				"type":        v.ProviderCode,
				"console_url": product.Console.GetHref(),
				"bookmark":    flattenBOMProductsBookmarkList(product.GetBookmarks()),
			})

		}

	}

	return productItems, nil
}

func flattenBOMProductsBookmarkList(bookmarkList []pingone.BillOfMaterialsProductsInnerBookmarksInner) []interface{} {
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
