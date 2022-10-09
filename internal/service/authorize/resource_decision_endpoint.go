package authorize

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceDecisionEndpoint() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Authorize decision endpoints.",

		CreateContext: resourceDecisionEndpointCreate,
		ReadContext:   resourceDecisionEndpointRead,
		UpdateContext: resourceDecisionEndpointUpdate,
		DeleteContext: resourceDecisionEndpointDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDecisionEndpointImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the decision endpoint in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"name": {
				Description:      "A string that specifies the policy decision resource name.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description:      "A string that specifies the description of the policy decision resource.",
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"owned": {
				Description: "A boolean that when true restricts modifications of the endpoint to PingOne-owned clients.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"record_recent_requests": {
				Description: "A boolean that specifies whether to record a limited history of recent decision requests and responses, which can be queried through a separate API.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"alternate_id": {
				Description: "A string that specifies alternative unique identifier for the endpoint, which provides a method for locating the resource by a known, fixed identifier.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"authorization_version_id": {
				Description:      "A string that specifies the ID of the Authorization Version deployed to this endpoint. Versioning allows independent development and deployment of policies. If omitted, the endpoint always uses the latest policy version available from the policy editor service.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
		},
	}
}

func resourceDecisionEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.AuthorizeAPIClient
	ctx = context.WithValue(ctx, authorize.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	decisionEndpoint := *authorize.NewDecisionEndpoint(d.Get("description").(string), d.Get("name").(string), d.Get("record_recent_requests").(bool)) // DecisionEndpoint |  (optional)

	if v, ok := d.GetOk("alternate_id"); ok {
		decisionEndpoint.SetAlternateId(v.(string))
	}

	if v, ok := d.GetOk("authorization_version_id"); ok {
		authzVersion := *authorize.NewDecisionEndpointAuthorizationVersion()
		authzVersion.SetId(v.(string))

		decisionEndpoint.SetAuthorizationVersion(authzVersion)
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.PolicyDecisionManagementApi.CreateDecisionEndpoint(ctx, d.Get("environment_id").(string)).DecisionEndpoint(decisionEndpoint).Execute()
		},
		"CreateDecisionEndpoint",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*authorize.DecisionEndpoint)

	d.SetId(respObject.GetId())

	return resourceDecisionEndpointRead(ctx, d, meta)
}

func resourceDecisionEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.AuthorizeAPIClient
	ctx = context.WithValue(ctx, authorize.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.PolicyDecisionManagementApi.ReadOneDecisionEndpoint(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneDecisionEndpoint",
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

	respObject := resp.(*authorize.DecisionEndpoint)

	d.Set("name", respObject.GetName())

	if v, ok := respObject.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	if v, ok := respObject.GetOwnedOk(); ok {
		d.Set("owned", v)
	} else {
		d.Set("owned", nil)
	}

	if v, ok := respObject.GetRecordRecentRequestsOk(); ok {
		d.Set("record_recent_requests", v)
	} else {
		d.Set("record_recent_requests", nil)
	}

	if v, ok := respObject.GetAlternateIdOk(); ok {
		d.Set("alternate_id", v)
	} else {
		d.Set("alternate_id", nil)
	}

	if v, ok := respObject.GetAuthorizationVersionOk(); ok {
		d.Set("authorization_version_id", v.GetId())
	} else {
		d.Set("authorization_version_id", nil)
	}

	return diags
}

func resourceDecisionEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.AuthorizeAPIClient
	ctx = context.WithValue(ctx, authorize.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	decisionEndpoint := *authorize.NewDecisionEndpoint(d.Get("description").(string), d.Get("name").(string), d.Get("record_recent_requests").(bool)) // DecisionEndpoint |  (optional)

	if v, ok := d.GetOk("alternate_id"); ok {
		decisionEndpoint.SetAlternateId(v.(string))
	}

	if v, ok := d.GetOk("authorization_version_id"); ok {
		authzVersion := *authorize.NewDecisionEndpointAuthorizationVersion()
		authzVersion.SetId(v.(string))

		decisionEndpoint.SetAuthorizationVersion(authzVersion)
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.PolicyDecisionManagementApi.UpdateDecisionEndpoint(ctx, d.Get("environment_id").(string), d.Id()).DecisionEndpoint(decisionEndpoint).Execute()
		},
		"UpdateDecisionEndpoint",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceDecisionEndpointRead(ctx, d, meta)
}

func resourceDecisionEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.AuthorizeAPIClient
	ctx = context.WithValue(ctx, authorize.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.PolicyDecisionManagementApi.DeleteDecisionEndpoint(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteDecisionEndpoint",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceDecisionEndpointImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/decisionEndpointID\"", d.Id())
	}

	environmentID, decisionEndpointID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(decisionEndpointID)

	resourceDecisionEndpointRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
