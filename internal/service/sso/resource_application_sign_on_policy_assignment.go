package sso

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
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
			fO, fR, fErr := apiClient.ApplicationSignOnPolicyAssignmentsApi.CreateSignOnPolicyAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string)).SignOnPolicyAssignment(applicationSignOnPolicyAssignment).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
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
			fO, fR, fErr := apiClient.ApplicationSignOnPolicyAssignmentsApi.ReadOneSignOnPolicyAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
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

	if v, ok := respObject.GetSignOnPolicyOk(); ok {
		d.Set("sign_on_policy_id", v.GetId())
	} else {
		d.Set("sign_on_policy_id", nil)
	}

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
			fO, fR, fErr := apiClient.ApplicationSignOnPolicyAssignmentsApi.UpdateSignOnPolicyAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).SignOnPolicyAssignment(applicationSignOnPolicyAssignment).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
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
			fR, fErr := apiClient.ApplicationSignOnPolicyAssignmentsApi.DeleteSignOnPolicyAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), nil, fR, fErr)
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

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "application_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "sign_on_policy_assignment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
	}

	attributes, err := framework.ParseImportID(d.Id(), idComponents...)
	if err != nil {
		return nil, err
	}

	d.Set("environment_id", attributes["environment_id"])
	d.Set("application_id", attributes["application_id"])
	d.SetId(attributes["sign_on_policy_assignment_id"])

	resourcePingOneApplicationSignOnPolicyAssignmentRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
