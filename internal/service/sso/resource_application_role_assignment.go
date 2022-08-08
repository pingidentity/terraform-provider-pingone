package sso

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"application_id": {
				Description:      "The ID of an application to assign an admin role to.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"role_id": {
				Description:      "The ID of an admin role to assign to the application.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"scope_organization_id": {
				Description:  "Limit the scope of the admin role assignment to the specified organisation ID.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"scope_organization_id", "scope_environment_id", "scope_population_id"},
			},
			"scope_environment_id": {
				Description:  "Limit the scope of the admin role assignment to the specified environment ID.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"scope_organization_id", "scope_environment_id", "scope_population_id"},
			},
			"scope_population_id": {
				Description:  "Limit the scope of the admin role assignment to the specified population ID.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"scope_organization_id", "scope_environment_id", "scope_population_id"},
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

	resp, r, err := apiClient.ApplicationsApplicationRoleAssignmentsApi.CreateApplicationRoleAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string)).RoleAssignment(applicationRoleAssignment).Execute()
	if (err != nil) || (r.StatusCode != 201) {

		response := &management.P1Error{}
		errDecode := json.NewDecoder(r.Body).Decode(response)
		if errDecode == nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Cannot decode error response: %v", errDecode),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})
		}

		if r.StatusCode == 400 && response.GetDetails()[0].GetTarget() == "scope" {
			diags = diag.FromErr(fmt.Errorf("Incompatible role and scope combination. Role: %s / Scope: %s", applicationRoleAssignmentRole.GetId(), applicationRoleAssignmentScope.GetType()))

			return diags
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationRoleAssignmentsApi.CreateApplicationRoleAssignment``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.SetId(resp.GetId())

	return resourcePingOneApplicationRoleAssignmentRead(ctx, d, meta)
}

func resourcePingOneApplicationRoleAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := apiClient.ApplicationsApplicationRoleAssignmentsApi.ReadOneApplicationRoleAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Role Assignment %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationRoleAssignmentsApi.ReadOneRoleAssignment``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.Set("role_id", resp.GetRole().Id)
	d.Set("read_only", resp.GetReadOnly())

	if resp.GetScope().Type == "ORGANIZATION" {
		d.Set("scope_organization_id", resp.GetScope().Id)
		d.Set("scope_environment_id", nil)
		d.Set("scope_population_id", nil)

	} else if resp.GetScope().Type == "ENVIRONMENT" {
		d.Set("scope_organization_id", nil)
		d.Set("scope_environment_id", resp.GetScope().Id)
		d.Set("scope_population_id", nil)

	} else if resp.GetScope().Type == "POPULATION" {
		d.Set("scope_organization_id", nil)
		d.Set("scope_environment_id", nil)
		d.Set("scope_population_id", resp.GetScope().Id)
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

	_, err := apiClient.ApplicationsApplicationRoleAssignmentsApi.DeleteApplicationRoleAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationRoleAssignmentsApi.DeleteApplicationRoleAssignment``: %v", err),
		})

		return diags
	}

	return nil
}

func resourcePingOneApplicationRoleAssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 3)

	if len(attributes) != 3 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/applicationID/roleAssignmentID\"", d.Id())
	}

	environmentID, applicationID, roleAssignmentID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("application_id", applicationID)
	d.SetId(roleAssignmentID)

	resourcePingOneApplicationRoleAssignmentRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
