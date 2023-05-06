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

func DatasourceUsers() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve multiple PingOne user IDs selected by a SCIM filter.",

		ReadContext: datasourcePingOneUsersRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment that contains the users to filter.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"scim_filter": {
				Description:  "A SCIM filter to apply to the user selection.  A SCIM filter offers the greatest flexibility in filtering users.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"scim_filter", "data_filter"},
			},
			"data_filter": {
				Description:  "Individual data filters to apply to the user selection.",
				Type:         schema.TypeSet,
				Optional:     true,
				ExactlyOneOf: []string{"scim_filter", "data_filter"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description:      "The attribute name to filter on.  Options are `accountId`, `address.streetAddress`, `address.locality`, `address.region`, `address.postalCode`, `address.countryCode`, `email`, `enabled`, `externalId`, `locale`, `mobilePhone`, `name.formatted`, `name.given`, `name.middle`, `name.family`, `name.honorificPrefix`, `name.honorificSuffix`, `nickname`, `population.id`, `photo.href`, `preferredLanguage`, `primaryPhone`, `timezone`, `title`, `type`, `username`, `memberOfGroups.id`.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"accountId", "address.streetAddress", "address.locality", "address.region", "address.postalCode", "address.countryCode", "email", "enabled", "externalId", "locale", "mobilePhone", "name.formatted", "name.given", "name.middle", "name.family", "name.honorificPrefix", "name.honorificSuffix", "nickname", "population.id", "photo.href", "preferredLanguage", "primaryPhone", "timezone", "title", "type", "username", "memberOfGroups.id"}, false)),
						},
						"values": {
							Description: "The possible values (case sensitive) of the attribute defined in the `name` parameter to filter.",
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
				Description: "The list of resulting IDs of users that have been successfully retrieved and filtered.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func datasourcePingOneUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var filterFunction sdk.SDKInterfaceFunc

	if v, ok := d.GetOk("scim_filter"); ok {

		filterFunction = func() (interface{}, *http.Response, error) {
			return apiClient.UsersApi.ReadAllUsers(ctx, d.Get("environment_id").(string)).Filter(v.(string)).Execute()
		}

	}

	if v, ok := d.GetOk("data_filter"); ok {

		attributeMapping := map[string]string{
			"memberOfGroups.id": `memberOfGroups[id eq "%s"]`,
		}

		scimFilter := filter.BuildScimFilter(v.(*schema.Set).List(), attributeMapping)

		filterFunction = func() (interface{}, *http.Response, error) {
			return apiClient.UsersApi.ReadAllUsers(ctx, d.Get("environment_id").(string)).Filter(scimFilter).Execute()
		}

	}

	resp, diags := sdk.ParseResponse(
		ctx,
		filterFunction,
		"ReadAllUsers",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.EntityArray)

	d.SetId(d.Get("environment_id").(string))

	usersList := respObject.GetEmbedded().Users

	idList := make([]string, 0)
	for _, v := range usersList {
		idList = append(idList, v.GetId())
	}

	d.Set("ids", idList)

	return diags
}
