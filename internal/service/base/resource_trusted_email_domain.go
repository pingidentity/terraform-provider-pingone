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

func ResourceTrustedEmailDomain() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Trusted Email Domains.",

		CreateContext: resourceTrustedEmailDomainCreate,
		ReadContext:   resourceTrustedEmailDomainRead,
		DeleteContext: resourceTrustedEmailDomainDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceTrustedEmailDomainImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the certificate in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"domain_name": {
				Type:             schema.TypeString,
				Description:      "A string that specifies the domain name to use, which must be provided and must be unique within an environment (for example, `demo.bxretail.org`).",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
		},
	}
}

func resourceTrustedEmailDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var resp interface{}

	emailDomain := *management.NewEmailDomain(d.Get("domain_name").(string))

	resp, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.TrustedEmailDomainsApi.CreateTrustedEmailDomain(ctx, d.Get("environment_id").(string)).EmailDomain(emailDomain).Execute()
		},
		"CreateTrustedEmailDomain",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.EmailDomain)

	d.SetId(respObject.GetId())

	return resourceTrustedEmailDomainRead(ctx, d, meta)
}

func resourceTrustedEmailDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.TrustedEmailDomainsApi.ReadOneTrustedEmailDomain(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneTrustedEmailDomain",
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

	respObject := resp.(*management.EmailDomain)

	d.Set("domain_name", respObject.GetDomainName())

	return diags
}

func resourceTrustedEmailDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.TrustedEmailDomainsApi.DeleteTrustedEmailDomain(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteTrustedEmailDomain",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceTrustedEmailDomainImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/trustedEmailDomainID\"", d.Id())
	}

	environmentID, trustedEmailDomainID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)

	d.SetId(trustedEmailDomainID)

	resourceTrustedEmailDomainRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
