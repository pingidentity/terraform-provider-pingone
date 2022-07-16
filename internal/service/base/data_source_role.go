package base

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func DatasourceRole() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne admin role data",

		ReadContext: datasourcePingOneRoleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description:      "The name of the role.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "A description of the role.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func datasourcePingOneRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var resp management.Role

	respList, r, err := apiClient.RolesApi.ReadAllRoles(ctx).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `RolesApi.ReadAllRoles``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	if roles, ok := respList.Embedded.GetRolesOk(); ok {

		found := false
		for _, role := range roles {

			if role.GetName() == d.Get("name").(string) {
				resp = role
				found = true
				break
			}
		}

		if !found {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Cannot find role %s", d.Get("name")),
			})

			return diags
		}

	}

	d.SetId(resp.GetId())
	d.Set("name", resp.GetName())
	d.Set("description", resp.GetDescription())

	return diags
}
