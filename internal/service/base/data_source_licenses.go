package base

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func DatasourceLicenses() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve multiple PingOne license IDs selected by a SCIM filter or a name/value list combination.",

		ReadContext: datasourcePingOneLicensesRead,

		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "The ID of the organization.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"scim_filter": {
				Description:  "A SCIM filter to apply to the license selection.  A SCIM filter offers the greatest flexibility in filtering licenses.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"scim_filter", "data_filter"},
			},
			"data_filter": {
				Description:  "Individual data filters to apply to the license selection.",
				Type:         schema.TypeSet,
				Optional:     true,
				ExactlyOneOf: []string{"scim_filter", "data_filter"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description:      "The attribute name to filter on.  Options are `name`, `package` or `status`.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"name", "package", "status"}, false)),
						},
						"values": {
							Description: "The possible values (case sensitive) of the attribute defined in the `name` parameter to filter.  If the attribute filter is `package`, the value is a free text, case-sensitive field.  Package names are not fixed and can change over time. If the attribute filter is `status`, available values are `ACTIVE`, `EXPIRED`, `FUTURE` and `TERMINATED`.  If the attribute filter is `name`, the exact name of the license should be provided.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"ids": {
				Description: "The list of resulting IDs of licenses that have been successfully retrieved and filtered.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func datasourcePingOneLicensesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var filterFunction sdk.SDKInterfaceFunc

	if v, ok := d.GetOk("scim_filter"); ok {

		filterFunction = func() (interface{}, *http.Response, error) {
			return apiClient.LicensesApi.ReadAllLicenses(ctx, d.Get("organization_id").(string)).Filter(v.(string)).Execute()
		}

	}

	if _, ok := d.GetOk("data_filter"); ok {

		filterFunction = func() (interface{}, *http.Response, error) {
			return apiClient.LicensesApi.ReadAllLicenses(ctx, d.Get("organization_id").(string)).Execute()
		}

	}

	resp, diags := sdk.ParseResponse(
		ctx,
		filterFunction,
		"ReadAllLicenses",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.EntityArray)

	d.SetId(d.Get("organization_id").(string))

	var licensesList []management.License

	if v, ok := d.GetOk("data_filter"); ok {
		licensesList = filterResults(v.(*schema.Set), respObject.GetEmbedded().Licenses)
	} else {
		licensesList = respObject.GetEmbedded().Licenses
	}

	idList := make([]string, 0)
	for _, v := range licensesList {
		idList = append(idList, v.GetId())
	}

	d.Set("ids", idList)

	return diags
}

func filterResults(filterSet *schema.Set, licenses []management.License) []management.License {
	items := make([]management.License, 0)

	for _, license := range licenses {

		filterMap := map[string]interface{}{
			"name":    license.GetName(),
			"package": license.GetPackage(),
			"status":  string(license.GetStatus()),
		}

		include := true

		for _, c := range filterSet.List() {

			obj := c.(map[string]interface{})

			for k, v := range filterMap {
				if obj["name"].(string) == k {
					if !obj["values"].(*schema.Set).Contains(v) {
						include = false
					}
				}
			}

		}

		if include {
			items = append(items, license)
		}

	}

	return items
}
