package mfa

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceFIDOPolicy() *schema.Resource {
	return &schema.Resource{

		DeprecationMessage: "This resource is deprecated, please use the `pingone_mfa_fido2_policy` resource going forward.  This resource is no longer configurable for environments created after 19th June 2023, nor environments that have been upgraded to use the latest FIDO2 policies. Existing environments that were created before 19th June 2023 and have not been upgraded can continue to use this resource to facilitate migration.",

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage MFA FIDO Policies in a PingOne Environment.",

		CreateContext: resourceFIDOPolicyCreate,
		ReadContext:   resourceFIDOPolicyRead,
		UpdateContext: resourceFIDOPolicyUpdate,
		DeleteContext: resourceFIDOPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceFIDOPolicyImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the FIDO policy in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:      "The name to use for the FIDO policy.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "Description of the FIDO policy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"attestation_requirements": {
				Description:      fmt.Sprintf("Determines whether attestation is requested from the authenticator, and whether this information is used to restrict authenticator usage. Can take one of these values: `%s` (attestation is not requested), `%s` (Attestation is requested and the information is used for logging purposes, but the information is not used for filtering authenticators), `%s` (all entries in the MDS table can be used for authentication), `%s` (only FIDO-certified authenticators can be used) and `%s` (only specific authenticators can be used. Used in conjunction with allowedAuthenticators.)", string(mfa.ENUMFIDOATTESTATIONREQUIREMENTS_NONE), string(mfa.ENUMFIDOATTESTATIONREQUIREMENTS_AUDIT_ONLY), string(mfa.ENUMFIDOATTESTATIONREQUIREMENTS_GLOBAL), string(mfa.ENUMFIDOATTESTATIONREQUIREMENTS_CERTIFIED), string(mfa.ENUMFIDOATTESTATIONREQUIREMENTS_SPECIFIC)),
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(mfa.ENUMFIDOATTESTATIONREQUIREMENTS_NONE), string(mfa.ENUMFIDOATTESTATIONREQUIREMENTS_AUDIT_ONLY), string(mfa.ENUMFIDOATTESTATIONREQUIREMENTS_GLOBAL), string(mfa.ENUMFIDOATTESTATIONREQUIREMENTS_CERTIFIED), string(mfa.ENUMFIDOATTESTATIONREQUIREMENTS_SPECIFIC)}, false)),
			},
			"allowed_authenticators": {
				Description: "If `attestation_requirements` is set to `SPECIFIC`, this array is used to specify the IDs of the authenticators that you want to allow.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				},
			},
			"enforce_during_authentication": {
				Description: "This parameter is relevant only if you have set `attestation_requirements` to `SPECIFIC` in order to restrict usage to only certain authenticators. If set to `true`, the policy will be applied both during registration and during each authentication attempt. If set to `false`, the policy is applied only during registration.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"default": {
				Description: "Whether this policy shoud serve as the default FIDO policy.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"resident_key_requirement": {
				Description:      fmt.Sprintf("Used to enable resident keys. Value can be `%s` or `%s`.", string(mfa.ENUMFIDORESIDENTKEYREQUIREMENT_DISCOURAGED), string(mfa.ENUMFIDORESIDENTKEYREQUIREMENT_REQUIRED)),
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(mfa.ENUMFIDORESIDENTKEYREQUIREMENT_DISCOURAGED), string(mfa.ENUMFIDORESIDENTKEYREQUIREMENT_REQUIRED)}, false)),
			},
		},
	}
}

func resourceFIDOPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	fidoPolicy := expandFIDOPolicy(d)

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.FIDOPolicyApi.CreateFidoPolicy(ctx, d.Get("environment_id").(string)).FIDOPolicy(*fidoPolicy).Execute()
		},
		"CreateFidoPolicy",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*mfa.FIDOPolicy)

	d.SetId(respObject.GetId())

	return resourceFIDOPolicyRead(ctx, d, meta)
}

func resourceFIDOPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.FIDOPolicyApi.ReadOneFidoPolicy(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneFidoPolicy",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	respObject := resp.(*mfa.FIDOPolicy)

	d.Set("name", respObject.GetName())

	if v, ok := respObject.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	d.Set("attestation_requirements", respObject.GetAttestationRequirements())

	if v, ok := respObject.GetAllowedAuthenticatorsOk(); ok && len(v) > 0 {

		items := make([]string, 0)

		for _, item := range v {
			items = append(items, item.GetId())
		}

		d.Set("allowed_authenticators", items)
	} else {
		d.Set("allowed_authenticators", nil)
	}

	if v, ok := respObject.GetEnforceDuringAuthenticationOk(); ok {
		d.Set("enforce_during_authentication", v)
	} else {
		d.Set("enforce_during_authentication", nil)
	}

	if v, ok := respObject.GetDefaultOk(); ok {
		d.Set("default", v)
	} else {
		d.Set("default", nil)
	}

	d.Set("resident_key_requirement", respObject.GetResidentKeyRequirement())

	return diags
}

func resourceFIDOPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	fidoPolicy := expandFIDOPolicy(d)

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.FIDOPolicyApi.UpdateFIDOPolicy(ctx, d.Get("environment_id").(string), d.Id()).FIDOPolicy(*fidoPolicy).Execute()
		},
		"UpdateFIDOPolicy",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return resourceFIDOPolicyRead(ctx, d, meta)
}

func resourceFIDOPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.FIDOPolicyApi.DeleteFidoPolicy(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteFidoPolicy",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceFIDOPolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/fidoPolicyID\"", d.Id())
	}

	environmentID, fidoPolicyID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(fidoPolicyID)

	resourceFIDOPolicyRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func expandFIDOPolicy(d *schema.ResourceData) *mfa.FIDOPolicy {

	item := mfa.NewFIDOPolicy(
		d.Get("name").(string),
		mfa.EnumFIDOAttestationRequirements(d.Get("attestation_requirements").(string)),
		mfa.EnumFIDOResidentKeyRequirement(d.Get("resident_key_requirement").(string)),
	)

	if v, ok := d.GetOk("description"); ok {
		item.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("allowed_authenticators"); ok {

		if v1, ok := v.(*schema.Set); ok && len(v1.List()) > 0 && v1.List()[0] != "" {
			items := make([]mfa.FIDOPolicyAllowedAuthenticatorsInner, 0)

			for _, item := range v1.List() {
				items = append(items, *mfa.NewFIDOPolicyAllowedAuthenticatorsInner(item.(string)))
			}

			item.SetAllowedAuthenticators(items)
		}
	}

	if v, ok := d.GetOk("enforce_during_authentication"); ok {
		item.SetEnforceDuringAuthentication(v.(bool))
	}

	// if v, ok := d.GetOk("default"); ok {
	// 	item.SetDefault(v.(bool))
	// }

	return item
}
