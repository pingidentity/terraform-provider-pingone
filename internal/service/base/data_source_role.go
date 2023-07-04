package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
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

	var diags diag.Diagnostics

	var resp management.Role

	respList, diags := sdk.ParseResponse(
		ctx,
		func() (interface{}, *http.Response, error) {
			return apiClient.RolesApi.ReadAllRoles(ctx).Execute()
		},
		"ReadAllRoles",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	if roles, ok := respList.(*management.EntityArray).Embedded.GetRolesOk(); ok {

		found := false
		for _, role := range roles {

			if role.GetName() == management.EnumRoleName(d.Get("name").(string)) {
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
