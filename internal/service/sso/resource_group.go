package sso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceGroup() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne groups",

		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGroupImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the group in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"name": {
				Description:      "The name of the group.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "A description to apply to the group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"population_id": {
				Description:      "The ID of the population that the group should be assigned to.",
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"user_filter": {
				Description: "A SCIM filter to dynamically assign users to the group.  Examples are found in the [PingOne online documentation](https://docs.pingidentity.com/bundle/pingone/page/kti1564020489340.html).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"external_id": {
				Description: "A user defined ID that represents the counterpart group in an external system.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	group := *management.NewGroup(d.Get("name").(string)) // Group |  (optional)

	if v, ok := d.GetOk("description"); ok {
		group.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("population_id"); ok {
		groupPopulation := *management.NewGroupPopulation(v.(string)) // NewGroupPopulation |  (optional)
		group.SetPopulation(groupPopulation)
	}

	if v, ok := d.GetOk("user_filter"); ok {
		group.SetUserFilter(v.(string))
	}

	if v, ok := d.GetOk("external_id"); ok {
		group.SetExternalId(v.(string))
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.GroupsApi.CreateGroup(ctx, d.Get("environment_id").(string)).Group(group).Execute()
		},
		"CreateGroup",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.Group)

	d.SetId(respObject.GetId())

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.GroupsApi.ReadOneGroup(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneGroup",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	respObject := resp.(*management.Group)

	d.Set("name", respObject.GetName())

	if v, ok := respObject.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	if v, ok := respObject.GetPopulationOk(); ok {
		d.Set("population_id", v.GetId())
	} else {
		d.Set("population_id", nil)
	}

	if v, ok := respObject.GetUserFilterOk(); ok {
		d.Set("user_filter", v)
	} else {
		d.Set("user_filter", nil)
	}

	if v, ok := respObject.GetExternalIdOk(); ok {
		d.Set("external_id", v)
	} else {
		d.Set("external_id", nil)
	}

	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	group := *management.NewGroup(d.Get("name").(string)) // Group |  (optional)

	if v, ok := d.GetOk("description"); ok {
		group.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("population_id"); ok {
		groupPopulation := *management.NewGroupPopulation(v.(string)) // NewGroupPopulation |  (optional)
		group.SetPopulation(groupPopulation)
	}

	if v, ok := d.GetOk("user_filter"); ok {
		group.SetUserFilter(v.(string))
	}

	if v, ok := d.GetOk("external_id"); ok {
		group.SetExternalId(v.(string))
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.GroupsApi.UpdateGroup(ctx, d.Get("environment_id").(string), d.Id()).Group(group).Execute()
		},
		"UpdateGroup",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.GroupsApi.DeleteGroup(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteGroup",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceGroupImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/groupID\"", d.Id())
	}

	environmentID, groupID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(groupID)

	resourceGroupRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
