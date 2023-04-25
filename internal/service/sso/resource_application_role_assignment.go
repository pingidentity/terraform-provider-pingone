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
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceApplicationRoleAssignment() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne admin role assignments to applications.",

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
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

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

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
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
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
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
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
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

		func() (interface{}, *http.Response, error) {
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
	splitLength := 3
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/applicationID/roleAssignmentID\"", d.Id())
	}

	environmentID, applicationID, roleAssignmentID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("application_id", applicationID)
	d.SetId(roleAssignmentID)

	resourcePingOneApplicationRoleAssignmentRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
