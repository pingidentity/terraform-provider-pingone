package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func DatasourceEnvironment() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne environment data",

		ReadContext: datasourcePingOneEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"name", "license_id"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:   "The name of the environment.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"environment_id", "license_id"},
			},
			"description": {
				Description: "A description of the environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "The type of the environment.  Options are SANDBOX for a development/testing environment and PRODUCTION for environments that require protection from deletion.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"region": {
				Description: "The region the environment is created in.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"license_id": {
				Description: "An ID of a valid license to apply to the environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			// "solution": {
			// 	Description: "The solution context of the environment.  Blank values indicate a non-workforce solution context.  Valid options are `WORKFORCE` and `CUSTOMER`",
			// 	Type:        schema.TypeString,
			// 	Computed:    true,
			// },
			"service": {
				Description: "The services enabled in the environment.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "The service type.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"console_url": {
							Description: "A custom console URL.  Generally used with services that are deployed separately to the PingOne SaaS service.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"bookmark": {
							Description: "Custom bookmark links for the service.",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description: "Bookmark name.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"url": {
										Description: "Bookmark URL.",
										Type:        schema.TypeString,
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

func datasourcePingOneEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var resp management.Environment

	if v, ok := d.GetOk("name"); ok {

		respList, diags := sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.EnvironmentsApi.ReadAllEnvironments(ctx).Execute()
			},
			"ReadAllEnvironments",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		if diags.HasError() {
			return diags
		}

		respObject := respList.(*management.EntityArray)

		if environments, ok := respObject.Embedded.GetEnvironmentsOk(); ok {

			found := false
			for _, environment := range environments {

				if environment.GetName() == v.(string) {
					resp = environment
					found = true
					break
				}
			}

			if !found {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Cannot find environment %s", v),
				})

				return diags
			}

		}

	} else if v, ok2 := d.GetOk("environment_id"); ok2 {

		environmentResp, diags := sdk.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, v.(string)).Execute()
			},
			"ReadOneEnvironment",
			sdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		if diags.HasError() {
			return diags
		}

		resp = *environmentResp.(*management.Environment)

	} else {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Neither environment_id or name are set",
			Detail:   "Neither environment_id or name are set",
		})

		return diags

	}

	d.SetId(resp.GetId())
	d.Set("environment_id", resp.GetId())
	d.Set("name", resp.GetName())

	if v, ok := resp.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	d.Set("type", resp.GetType())
	d.Set("region", model.FindRegionByAPICode(resp.GetRegion()).Region)
	d.Set("license_id", resp.GetLicense().Id)

	// The bill of materials

	servicesResp, diags := sdk.ParseResponse(
		ctx,
		func() (interface{}, *http.Response, error) {
			return apiClient.BillOfMaterialsBOMApi.ReadOneBillOfMaterials(ctx, resp.GetId()).Execute()
		},
		"ReadOneBillOfMaterials",
		sdk.DefaultCustomError,
		retryEnvironmentDefault,
	)
	if diags.HasError() {
		return diags
	}

	// d.Set("solution", servicesResp.SolutionType)
	productBOMItems, err := flattenBOMProducts(servicesResp.(*management.BillOfMaterials))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error mapping platform services with the configured services``: %v", err),
			Detail:   fmt.Sprintf("Platform services: %v\n", servicesResp.(*management.BillOfMaterials)),
		})

		return diags
	}
	d.Set("service", productBOMItems)

	return diags
}
