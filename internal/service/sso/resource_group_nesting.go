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

func ResourceGroupNesting() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne group nesting.",

		CreateContext: resourceGroupNestingCreate,
		ReadContext:   resourceGroupNestingRead,
		DeleteContext: resourceGroupNestingDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGroupNestingImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the group in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"group_id": {
				Description:      "The ID of the environment to create the group in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"nested_group_id": {
				Description:      "The ID of the environment to create the group in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"type": {
				Description: "The type of the group nesting.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceGroupNestingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	groupNesting := *management.NewGroupNesting(d.Get("nested_group_id").(string)) // GroupNesting |  (optional)

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.GroupsApi.CreateGroupNesting(ctx, d.Get("environment_id").(string), d.Get("group_id").(string)).GroupNesting(groupNesting).Execute()
		},
		"CreateGroupNesting",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.GroupNesting)

	d.SetId(respObject.GetId())

	return resourceGroupNestingRead(ctx, d, meta)
}

func resourceGroupNestingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.GroupsApi.ReadOneGroupNesting(ctx, d.Get("environment_id").(string), d.Get("group_id").(string), d.Id()).Execute()
		},
		"ReadOneGroupNesting",
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

	respObject := resp.(*management.GroupNesting)

	d.Set("type", respObject.GetType())

	return diags
}

func resourceGroupNestingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.GroupsApi.DeleteGroupNesting(ctx, d.Get("environment_id").(string), d.Get("group_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteGroupNesting",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceGroupNestingImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 3
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/groupID/groupNestingID\"", d.Id())
	}

	environmentID, groupID, groupNestingID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("group_id", groupID)
	d.SetId(groupNestingID)

	resourceGroupNestingRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
