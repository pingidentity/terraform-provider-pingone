package sso

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	pingone "github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func ResourceUser() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne users",

		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceUserImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the user in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				ForceNew:         true,
			},
			"username": {
				Description:      "The username of the user.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty), // TODO: validation per API docs
				/*
					pattern: '^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$'
					minLength: 1
					maxLength: 128
				*/
			},
			"email": {
				Description:      "The email address of the user.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty), // TODO: Email RFC format
			},
			"status": {
				Description:      "The enabled status of the user.  Possible values are `ENABLED` or `DISABLED`.",
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ENABLED", "DISABLED"}, false)),
				Default:          "ENABLED",
				Optional:         true,
			},
			"population_id": {
				Description:      "The population ID to add the user to.",
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				Required:         true,
			},
			// TODO: Full schema as-and-when needed
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	user := *pingone.NewUser(d.Get("email").(string), d.Get("username").(string))

	population := *pingone.NewUserPopulation(d.Get("population_id").(string))
	user.SetPopulation(population)

	// Create user

	resp, r, err := apiClient.UsersUsersApi.CreateUser(ctx, d.Get("environment_id").(string)).User(user).Execute()
	if (err != nil) || (r.StatusCode != 201) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `UsersUsersApi.CreateUser``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	// Set status

	userEnabled := *pingone.NewUserEnabled() // UserEnabled |  (optional)
	if d.Get("status").(string) == "ENABLED" {
		userEnabled.SetEnabled(true)
	} else {
		userEnabled.SetEnabled(false)
	}

	_, r, err = apiClient.UsersEnableUsersApi.UpdateUserEnabled(ctx, d.Get("environment_id").(string), resp.GetId()).UserEnabled(userEnabled).Execute()
	if (err != nil) || (r.StatusCode != 200) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `UsersEnableUsersApi.UpdateUserEnabled``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.SetId(resp.GetId())

	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := apiClient.UsersUsersApi.ReadUser(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne User %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `UsersUsersApi.ReadUser``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.Set("username", resp.GetUsername())
	d.Set("email", resp.GetEmail())
	if resp.GetEnabled() {
		d.Set("status", "ENABLED")
	} else {
		d.Set("status", "DISABLED")
	}
	d.Set("population_id", resp.GetPopulation().Id)

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	// The user
	user := *pingone.NewUser(d.Get("email").(string), d.Get("username").(string))

	_, r, err := apiClient.UsersUsersApi.UpdateUserPut(ctx, d.Get("environment_id").(string), d.Id()).User(user).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `UsersUsersApi.UpdateUserPut``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	// Set status

	userEnabled := *pingone.NewUserEnabled() // UserEnabled |  (optional)
	if d.Get("status").(string) == "ENABLED" {
		userEnabled.SetEnabled(true)
	} else {
		userEnabled.SetEnabled(false)
	}

	_, r, err = apiClient.UsersEnableUsersApi.UpdateUserEnabled(ctx, d.Get("environment_id").(string), d.Id()).UserEnabled(userEnabled).Execute()
	if (err != nil) || (r.StatusCode != 200) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `UsersEnableUsersApi.UpdateUserEnabled``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	// Set population

	population := *pingone.NewUserPopulation(d.Get("population_id").(string))

	_, r, err = apiClient.UsersUserPopulationsApi.UpdateUserPopulation(ctx, d.Get("environment_id").(string), d.Id()).UserPopulation(population).Execute()
	if (err != nil) || (r.StatusCode != 200) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `UsersUserPopulationsApi.UpdateUserPopulation``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	return resourceUserRead(ctx, d, meta)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	_, err := apiClient.UsersUsersApi.DeleteUser(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `UsersUsersApi.DeleteUser``: %v", err),
		})

		return diags
	}

	return nil
}

func resourceUserImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/userID\"", d.Id())
	}

	environmentID, userID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(userID)

	resourceUserRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
