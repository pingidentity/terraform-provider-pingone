package sso

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/filter"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func DatasourceUser() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne user data",

		ReadContext: datasourcePingOneUserRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"user_id": {
				Description:      "The ID of the user.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"user_id", "username", "email"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"username": {
				Description:  "The username of the user.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"user_id", "username", "email"},
			},
			"email": {
				Description:  "The email address of the user.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"user_id", "username", "email"},
			},
			"status": {
				Description: "The enabled status of the user.  Possible values are `ENABLED` or `DISABLED`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"population_id": {
				Description: "The population ID the user is assigned to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func datasourcePingOneUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var resp management.User
	var scimFilter string

	if v, ok := d.GetOk("username"); ok {

		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "username",
				"values": []string{v.(string)},
			}), map[string]string{})

	} else if v, ok := d.GetOk("user_id"); ok {

		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "id",
				"values": []string{v.(string)},
			}), map[string]string{})

	} else if v, ok := d.GetOk("email"); ok {

		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "email",
				"values": []string{v.(string)},
			}), map[string]string{})

	} else {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "None of user_id, username or email are set",
		})

		return diags

	}

	respList, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.UsersApi.ReadAllUsers(ctx, d.Get("environment_id").(string)).Filter(scimFilter).Execute()
		},
		"ReadAllUsers",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	if users, ok := respList.(*management.EntityArray).Embedded.GetUsersOk(); ok && len(users) > 0 && users[0].Id != nil {

		resp = users[0]

		d.SetId(resp.GetId())
		d.Set("user_id", resp.GetId())
		d.Set("username", resp.GetUsername())
		d.Set("email", resp.GetEmail())
		d.Set("status", string(*resp.GetAccount().Status))
		d.Set("population_id", resp.GetPopulation().Id)

	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot find user",
		})

		return diags
	}

	return diags
}
