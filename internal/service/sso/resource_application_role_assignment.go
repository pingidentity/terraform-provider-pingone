package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceApplicationRoleAssignment() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne admin role assignments to administrator defined applications.",

		CreateContext: resourcePingOneApplicationRoleAssignmentCreate,
		ReadContext:   resourcePingOneApplicationRoleAssignmentRead,
		DeleteContext: resourcePingOneApplicationRoleAssignmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePingOneApplicationRoleAssignmentImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"application_id": {
				Description:      "The ID of an application to assign an admin role to.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"role_id": {
				Description:      "The ID of an admin role to assign to the application.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"scope_organization_id": {
				Description:      "Limit the scope of the admin role assignment to the specified organisation ID.",
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ExactlyOneOf:     []string{"scope_organization_id", "scope_environment_id", "scope_population_id"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"scope_environment_id": {
				Description:      "Limit the scope of the admin role assignment to the specified environment ID.",
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ExactlyOneOf:     []string{"scope_organization_id", "scope_environment_id", "scope_population_id"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"scope_population_id": {
				Description:      "Limit the scope of the admin role assignment to the specified population ID.",
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ExactlyOneOf:     []string{"scope_organization_id", "scope_environment_id", "scope_population_id"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"read_only": {
				Description: "A flag to show whether the admin role assignment is read only or can be changed.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func resourcePingOneApplicationRoleAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	//d.Get("scope_id").(string)

	scopeID := ""
	scopeType := ""
	if organisationID, ok := d.GetOk("scope_organization_id"); ok {
		scopeID = organisationID.(string)
		scopeType = "ORGANIZATION"

	} else if environmentID, ok := d.GetOk("scope_environment_id"); ok {
		scopeID = environmentID.(string)
		scopeType = "ENVIRONMENT"

	} else if populationID, ok := d.GetOk("scope_population_id"); ok {
		scopeID = populationID.(string)
		scopeType = "POPULATION"

	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "One of scope_organization_id, scope_environment_id or scope_population_id must be set",
			Detail:   "One of scope_organization_id, scope_environment_id or scope_population_id must be set",
		})

		return diags
	}

	applicationRoleAssignmentRole := *management.NewRoleAssignmentRole(d.Get("role_id").(string))
	applicationRoleAssignmentScope := *management.NewRoleAssignmentScope(scopeID, management.EnumRoleAssignmentScopeType(scopeType))
	applicationRoleAssignment := *management.NewRoleAssignment(applicationRoleAssignmentRole, applicationRoleAssignmentScope) // ApplicationRoleAssignment |  (optional)

	applicationOk, diags := checkApplicationTypeForRoleAssignment(ctx, apiClient, d.Get("environment_id").(string), d.Get("application_id").(string))
	if diags.HasError() {
		return diags
	}
	if !applicationOk {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid parameter value - Unmappable application type",
			Detail:   fmt.Sprintf("The application ID provided (%s) relates to an application that is neither `OPENID_CONNECT` or `SAML` type.  Roles cannot be mapped to this application.", d.Get("application_id").(string)),
		})

		return diags
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ApplicationRoleAssignmentsApi.CreateApplicationRoleAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string)).RoleAssignment(applicationRoleAssignment).Execute()
		},
		"CreateApplicationRoleAssignment",
		func(error model.P1Error) diag.Diagnostics {

			// Invalid role/scope combination
			if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
				if target, ok := details[0].GetTargetOk(); ok && *target == "scope" {
					diags = diag.FromErr(fmt.Errorf("Incompatible role and scope combination. Role: %s / Scope: %s", applicationRoleAssignmentRole.GetId(), applicationRoleAssignmentScope.GetType()))

					return diags
				}
			}

			return nil
		},
		sdk.RoleAssignmentRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.RoleAssignment)

	d.SetId(respObject.GetId())

	return resourcePingOneApplicationRoleAssignmentRead(ctx, d, meta)
}

func resourcePingOneApplicationRoleAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	applicationOk, diags := checkApplicationTypeForRoleAssignment(ctx, apiClient, d.Get("environment_id").(string), d.Get("application_id").(string))
	if diags.HasError() {
		return diags
	}
	if !applicationOk {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid parameter value - Unmappable application type",
			Detail:   fmt.Sprintf("The application ID provided (%s) relates to an application that is neither `OPENID_CONNECT` or `SAML` type.  Roles cannot be mapped to this application.", d.Get("application_id").(string)),
		})

		d.SetId("")

		return diags
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ApplicationRoleAssignmentsApi.ReadOneApplicationRoleAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
		},
		"ReadOneApplicationRoleAssignment",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	respObject := resp.(*management.RoleAssignment)

	d.Set("role_id", respObject.GetRole().Id)
	d.Set("read_only", respObject.GetReadOnly())

	if respObject.GetScope().Type == "ORGANIZATION" {
		d.Set("scope_organization_id", respObject.GetScope().Id)
		d.Set("scope_environment_id", nil)
		d.Set("scope_population_id", nil)

	} else if respObject.GetScope().Type == "ENVIRONMENT" {
		d.Set("scope_organization_id", nil)
		d.Set("scope_environment_id", respObject.GetScope().Id)
		d.Set("scope_population_id", nil)

	} else if respObject.GetScope().Type == "POPULATION" {
		d.Set("scope_organization_id", nil)
		d.Set("scope_environment_id", nil)
		d.Set("scope_population_id", respObject.GetScope().Id)
	}

	return diags
}

func resourcePingOneApplicationRoleAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	if d.Get("read_only").(bool) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Role assignment %s cannot be deleted as it is read only", d.Id()),
		})

		return diags
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			r, err := apiClient.ApplicationRoleAssignmentsApi.DeleteApplicationRoleAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteApplicationRoleAssignment",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourcePingOneApplicationRoleAssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

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
			Label:  "role_assignment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
	}

	attributes, err := framework.ParseImportID(d.Id(), idComponents...)
	if err != nil {
		return nil, err
	}

	d.Set("environment_id", attributes["environment_id"])
	d.Set("application_id", attributes["application_id"])
	d.SetId(attributes["role_assignment_id"])

	resourcePingOneApplicationRoleAssignmentRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func checkApplicationTypeForRoleAssignment(ctx context.Context, apiClient *management.APIClient, environmentId, applicationId string) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	resp, d := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ApplicationsApi.ReadOneApplication(ctx, environmentId, applicationId).Execute()
		},
		"ReadOneApplication",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	diags = append(diags, d...)
	if diags.HasError() {
		return false, diags
	}

	respObject := resp.(*management.ReadOneApplication200Response)

	if respObject.ApplicationOIDC != nil && respObject.ApplicationOIDC.GetId() != "" {
		return true, diags
	}

	if respObject.ApplicationSAML != nil && respObject.ApplicationSAML.GetId() != "" {
		return true, diags
	}

	return false, diags
}
