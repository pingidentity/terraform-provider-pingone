package base

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
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"user_id": {
				Description:      "The ID of a user to assign an admin role to.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"role_id": {
				Description:      "The ID of an admin role to assign to the user.",
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
				ConflictsWith:    []string{"scope_environment_id", "scope_population_id"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"scope_environment_id": {
				Description:      "Limit the scope of the admin role assignment to the specified environment ID.",
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"scope_organization_id", "scope_population_id"},
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"scope_population_id": {
				Description:      "Limit the scope of the admin role assignment to the specified population ID.",
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"scope_organization_id", "scope_environment_id"},
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

func resourcePingOneRoleAssignmentUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	userRoleAssignmentRole := *management.NewRoleAssignmentRole(d.Get("role_id").(string))
	userRoleAssignmentScope := *management.NewRoleAssignmentScope(scopeID, management.EnumRoleAssignmentScopeType(scopeType))
	userRoleAssignment := *management.NewRoleAssignment(userRoleAssignmentRole, userRoleAssignmentScope) // UserRoleAssignment |  (optional)

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.UsersUserRoleAssignmentsApi.CreateUserRoleAssignment(ctx, d.Get("environment_id").(string), d.Get("user_id").(string)).RoleAssignment(userRoleAssignment).Execute()
		},
		"CreateUserRoleAssignment",
		func(error management.P1Error) diag.Diagnostics {

			// Invalid role/scope combination
			if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
				if target, ok := details[0].GetTargetOk(); ok && *target == "scope" {
					diags = diag.FromErr(fmt.Errorf("Incompatible role and scope combination. Role: %s / Scope: %s", userRoleAssignmentRole.GetId(), userRoleAssignmentScope.GetType()))

					return diags
				}
			}

			return nil
		},
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.RoleAssignment)

	d.SetId(respObject.GetId())

	return resourcePingOneRoleAssignmentUserRead(ctx, d, meta)
}

func resourcePingOneRoleAssignmentUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.UsersUserRoleAssignmentsApi.ReadOneUserRoleAssignment(ctx, d.Get("environment_id").(string), d.Get("user_id").(string), d.Id()).Execute()
		},
		"ReadOneUserRoleAssignment",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
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

func resourcePingOneRoleAssignmentUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			r, err := apiClient.UsersUserRoleAssignmentsApi.DeleteUserRoleAssignment(ctx, d.Get("environment_id").(string), d.Get("user_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteUserRoleAssignment",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
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
