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
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
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
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				Required:         true,
			},
			// TODO: Full schema as-and-when needed
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	user := *management.NewUser(d.Get("email").(string), d.Get("username").(string))

	population := *management.NewUserPopulation(d.Get("population_id").(string))
	user.SetPopulation(population)

	// Create user

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.UsersApi.CreateUser(ctx, d.Get("environment_id").(string)).User(user).Execute()
		},
		"CreateUser",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.User)

	// Set status

	userEnabled := *management.NewUserEnabled() // UserEnabled |  (optional)
	if d.Get("status").(string) == "ENABLED" {
		userEnabled.SetEnabled(true)
	} else {
		userEnabled.SetEnabled(false)
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.EnableUsersApi.UpdateUserEnabled(ctx, d.Get("environment_id").(string), respObject.GetId()).UserEnabled(userEnabled).Execute()
		},
		"UpdateUserEnabled",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	d.SetId(respObject.GetId())

	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.UsersApi.ReadUser(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadUser",
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

	respObject := resp.(*management.User)

	d.Set("username", respObject.GetUsername())
	d.Set("email", respObject.GetEmail())
	if respObject.GetEnabled() {
		d.Set("status", "ENABLED")
	} else {
		d.Set("status", "DISABLED")
	}
	d.Set("population_id", respObject.GetPopulation().Id)

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	// The user
	user := *management.NewUser(d.Get("email").(string), d.Get("username").(string))

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.UsersApi.UpdateUserPut(ctx, d.Get("environment_id").(string), d.Id()).User(user).Execute()
		},
		"UpdateUserPut",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	// Set status

	userEnabled := *management.NewUserEnabled() // UserEnabled |  (optional)
	if d.Get("status").(string) == "ENABLED" {
		userEnabled.SetEnabled(true)
	} else {
		userEnabled.SetEnabled(false)
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.EnableUsersApi.UpdateUserEnabled(ctx, d.Get("environment_id").(string), d.Id()).UserEnabled(userEnabled).Execute()
		},
		"UpdateUserEnabled",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	// Set population

	population := *management.NewUserPopulation(d.Get("population_id").(string))

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.UserPopulationsApi.UpdateUserPopulation(ctx, d.Get("environment_id").(string), d.Id()).UserPopulation(population).Execute()
		},
		"UpdateUserPopulation",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return resourceUserRead(ctx, d, meta)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.UsersApi.DeleteUser(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteUser",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceUserImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/userID\"", d.Id())
	}

	environmentID, userID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(userID)

	resourceUserRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
