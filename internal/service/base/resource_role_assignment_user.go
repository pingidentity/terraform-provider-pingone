package base

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	pingone "github.com/patrickcping/pingone-go/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func ResourceRoleAssignmentUser() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne admin role assignments to users.",

		CreateContext: resourcePingOneRoleAssignmentUserCreate,
		ReadContext:   resourcePingOneRoleAssignmentUserRead,
		DeleteContext: resourcePingOneRoleAssignmentUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePingOneRoleAssignmentUserImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"user_id": {
				Description:      "The ID of a user to assign an admin role to.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"role_id": {
				Description:      "The ID of an admin role to assign to the user.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"scope_organization_id": {
				Description:   "Limit the scope of the admin role assignment to the specified organisation ID.",
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"scope_environment_id", "scope_population_id"},
			},
			"scope_environment_id": {
				Description:   "Limit the scope of the admin role assignment to the specified environment ID.",
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"scope_organization_id", "scope_population_id"},
			},
			"scope_population_id": {
				Description:   "Limit the scope of the admin role assignment to the specified population ID.",
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"scope_organization_id", "scope_environment_id"},
			},
			"read_only": {
				Description: "A flag to show whether the admin role assignment is read only or can be changed.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func resourcePingOneRoleAssignmentUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
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

	userRoleAssignmentRole := *pingone.NewRoleAssignmentRole(d.Get("role_id").(string))
	userRoleAssignmentScope := *pingone.NewRoleAssignmentScope(scopeID, scopeType)
	userRoleAssignment := *pingone.NewRoleAssignment(userRoleAssignmentRole, userRoleAssignmentScope) // UserRoleAssignment |  (optional)

	resp, r, err := apiClient.UsersUserRoleAssignmentsApi.CreateUserRoleAssignment(ctx, d.Get("environment_id").(string), d.Get("user_id").(string)).RoleAssignment(userRoleAssignment).Execute()
	if (err != nil) || (r.StatusCode != 201) {

		response := &pingone.P1Error{}
		errDecode := json.NewDecoder(r.Body).Decode(response)
		if errDecode == nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Cannot decode error response: %v", errDecode),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})
		}

		if r.StatusCode == 400 && response.GetDetails()[0].GetTarget() == "scope" {
			diags = diag.FromErr(fmt.Errorf("Incompatible role and scope combination. Role: %s / Scope: %s", userRoleAssignmentRole.GetId(), userRoleAssignmentScope.GetType()))

			return diags
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `UsersUserRoleAssignmentsApi.CreateUserRoleAssignment``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.SetId(resp.GetId())

	return resourcePingOneRoleAssignmentUserRead(ctx, d, meta)
}

func resourcePingOneRoleAssignmentUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := apiClient.UsersUserRoleAssignmentsApi.ReadOneUserRoleAssignment(ctx, d.Get("environment_id").(string), d.Get("user_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Role Assignment %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `UsersUserRoleAssignmentsApi.ReadOneRoleAssignment``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.Set("role_id", resp.GetRole().Id)
	d.Set("read_only", resp.GetReadOnly())

	if resp.GetScope().Type == "ORGANIZATION" {
		d.Set("scope_organization_id", resp.GetScope().Id)

	} else if resp.GetScope().Type == "ENVIRONMENT" {
		d.Set("scope_environment_id", resp.GetScope().Id)

	} else if resp.GetScope().Type == "POPULATION" {
		d.Set("scope_population_id", resp.GetScope().Id)
	}

	return diags
}

func resourcePingOneRoleAssignmentUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	if d.Get("read_only").(bool) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Role assignment %s cannot be deleted as it is read only", d.Id()),
		})

		return diags
	}

	_, err := apiClient.UsersUserRoleAssignmentsApi.DeleteUserRoleAssignment(ctx, d.Get("environment_id").(string), d.Get("user_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `UsersUserRoleAssignmentsApi.DeleteUserRoleAssignment``: %v", err),
		})

		return diags
	}

	return nil
}

func resourcePingOneRoleAssignmentUserImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 3)

	if len(attributes) != 3 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/userID/roleAssignmentID\"", d.Id())
	}

	environmentID, userID, roleAssignmentID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("user_id", userID)
	d.SetId(roleAssignmentID)

	resourcePingOneRoleAssignmentUserRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
