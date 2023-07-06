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

func ResourceApplicationSignOnPolicyAssignment() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage a sign-on policy assignment for administrator defined applications or built-in system applications configured in PingOne.",

		CreateContext: resourcePingOneApplicationSignOnPolicyAssignmentCreate,
		ReadContext:   resourcePingOneApplicationSignOnPolicyAssignmentRead,
		UpdateContext: resourcePingOneApplicationSignOnPolicyAssignmentUpdate,
		DeleteContext: resourcePingOneApplicationSignOnPolicyAssignmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePingOneApplicationSignOnPolicyAssignmentImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the application sign-on policy assignment in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"application_id": {
				Description:      "The ID of the application to create the sign-on policy assignment for.\n\n-> The value for `application_id` may come from the `id` attribute of the `pingone_application` or `pingone_system_application` resources or data sources.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"sign_on_policy_id": {
				Description:      "The ID of the sign-on policy resource to associate.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"priority": {
				Description:      "The order in which the policy referenced by this assignment is evaluated during an authentication flow relative to other policies. An assignment with a lower priority will be evaluated first.",
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
			},
		},
	}
}

func resourcePingOneApplicationSignOnPolicyAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	signOnPolicy := *management.NewSignOnPolicyActionCommonSignOnPolicy(d.Get("sign_on_policy_id").(string))
	applicationSignOnPolicyAssignment := *management.NewSignOnPolicyAssignment(int32(d.Get("priority").(int)), signOnPolicy)

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ApplicationSignOnPolicyAssignmentsApi.CreateSignOnPolicyAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string)).SignOnPolicyAssignment(applicationSignOnPolicyAssignment).Execute()
		},
		"CreateSignOnPolicyAssignment",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.SignOnPolicyAssignment)

	d.SetId(respObject.GetId())

	return resourcePingOneApplicationSignOnPolicyAssignmentRead(ctx, d, meta)
}

func resourcePingOneApplicationSignOnPolicyAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ApplicationSignOnPolicyAssignmentsApi.ReadOneSignOnPolicyAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
		},
		"ReadOneSignOnPolicyAssignment",
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

	respObject := resp.(*management.SignOnPolicyAssignment)

	d.Set("priority", respObject.GetPriority())

	return diags
}

func resourcePingOneApplicationSignOnPolicyAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	signOnPolicy := *management.NewSignOnPolicyActionCommonSignOnPolicy(d.Get("sign_on_policy_id").(string))
	applicationSignOnPolicyAssignment := *management.NewSignOnPolicyAssignment(int32(d.Get("priority").(int)), signOnPolicy)

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ApplicationSignOnPolicyAssignmentsApi.UpdateSignOnPolicyAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).SignOnPolicyAssignment(applicationSignOnPolicyAssignment).Execute()
		},
		"UpdateSignOnPolicyAssignment",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return resourcePingOneApplicationSignOnPolicyAssignmentRead(ctx, d, meta)
}

func resourcePingOneApplicationSignOnPolicyAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			r, err := apiClient.ApplicationSignOnPolicyAssignmentsApi.DeleteSignOnPolicyAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteSignOnPolicyAssignment",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourcePingOneApplicationSignOnPolicyAssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 3
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/applicationID/SignOnPolicyAssignmentID\"", d.Id())
	}

	environmentID, applicationID, SignOnPolicyAssignmentID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("application_id", applicationID)
	d.SetId(SignOnPolicyAssignmentID)

	resourcePingOneApplicationSignOnPolicyAssignmentRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
