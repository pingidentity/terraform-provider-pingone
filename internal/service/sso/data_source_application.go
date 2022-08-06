package sso

// import (
// 	"context"
// 	"fmt"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	"github.com/patrickcping/pingone-go-sdk-v2/management"
// 	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
// )

// func DatasourceApplication() *schema.Resource {
// 	return &schema.Resource{

// 		// This description is used by the documentation generator and the language server.
// 		Description: "Datasource to read PingOne schema data",

// 		ReadContext: datasourcePingOneApplicationRead,

// 		Schema: map[string]*schema.Schema{
// 			"environment_id": {
// 				Description: "The ID of the environment.",
// 				Type:        schema.TypeString,
// 				Required:    true,
// 			},
// 			"schema_id": {
// 				Description:   "The ID of the schema.",
// 				Type:          schema.TypeString,
// 				Optional:      true,
// 				ConflictsWith: []string{"name"},
// 			},
// 			"name": {
// 				Description:   "The name of the schema.",
// 				Type:          schema.TypeString,
// 				Optional:      true,
// 				ConflictsWith: []string{"schema_id"},
// 			},
// 			"description": {
// 				Description: "A description of the schema.",
// 				Type:        schema.TypeString,
// 				Computed:    true,
// 			},
// 		},
// 	}
// }

// func datasourcePingOneApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	p1Client := meta.(*client.Client)
// 	apiClient := p1Client.API.ManagementAPIClient
// 	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
// 		"suffix": p1Client.API.Region.URLSuffix,
// 	})
// 	var diags diag.Diagnostics

// 	var resp management.Schema

// 	if v, ok := d.GetOk("name"); ok {

// 		respList, r, err := apiClient.ApplicationsApi.ReadAllApplications(ctx, d.Get("environment_id").(string)).Execute()
// 		if err != nil {
// 			diags = append(diags, diag.Diagnostic{
// 				Severity: diag.Error,
// 				Summary:  fmt.Sprintf("Error when calling `ApplicationsApi.ReadAllApplications``: %v", err),
// 				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
// 			})

// 			return diags
// 		}

// 		if schemas, ok := respList.Embedded.GetApplicationsOk(); ok {

// 			found := false
// 			for _, schema := range schemas {

// 				if schema.GetName() == v.(string) {
// 					resp = schema
// 					found = true
// 					break
// 				}
// 			}

// 			if !found {
// 				diags = append(diags, diag.Diagnostic{
// 					Severity: diag.Error,
// 					Summary:  fmt.Sprintf("Cannot find schema %s", v),
// 				})

// 				return diags
// 			}

// 		}

// 	} else if v, ok2 := d.GetOk("schema_id"); ok2 {

// 		schemaResp, r, err := apiClient.ApplicationsApi.ReadOneApplication(ctx, d.Get("environment_id").(string), v.(string)).Execute()
// 		if err != nil {
// 			diags = append(diags, diag.Diagnostic{
// 				Severity: diag.Error,
// 				Summary:  fmt.Sprintf("Error when calling `ApplicationsApi.ReadOneApplication``: %v", err),
// 				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
// 			})

// 			return diags
// 		}

// 		resp = *schemaResp

// 	} else {

// 		diags = append(diags, diag.Diagnostic{
// 			Severity: diag.Error,
// 			Summary:  "Neither schema_id or name are set",
// 			Detail:   "Neither schema_id or name are set",
// 		})

// 		return diags

// 	}

// 	d.SetId(resp.GetId())
// 	d.Set("schema_id", resp.GetId())
// 	d.Set("name", resp.GetName())
// 	d.Set("description", resp.GetDescription())

// 	return diags
// }
