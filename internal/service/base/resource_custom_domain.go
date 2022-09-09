package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceCustomDomain() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Custom Domains.",

		CreateContext: resourceCustomDomainCreate,
		ReadContext:   resourceCustomDomainRead,
		DeleteContext: resourceCustomDomainDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCustomDomainImport,
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
				Description:      "A string that specifies the domain name to use, which must be provided and must be unique within an environment (for example, `auth.bxretail.org`).",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"status": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("A string that specifies the status of the custom domain. Options are `%s`, `%s` and `%s`.", string(management.ENUMCUSTOMDOMAINSTATUS_ACTIVE), string(management.ENUMCUSTOMDOMAINSTATUS_VERIFICATION_REQUIRED), string(management.ENUMCUSTOMDOMAINSTATUS_SSL_CERTIFICATE_REQUIRED)),
				Computed:    true,
			},
			"canonical_name": {
				Type:        schema.TypeString,
				Description: "A string that specifies the domain name that should be used as the value of the CNAME record in the customer’s DNS.",
				Computed:    true,
			},
			"certificate_expires_at": {
				Type:        schema.TypeString,
				Description: "The time when the certificate expires.  If this property is not present, it indicates that an SSL certificate has not been setup for this custom domain.",
				Computed:    true,
			},
		},
	}
}

func resourceCustomDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var resp interface{}

	customDomain := *management.NewCustomDomain(d.Get("domain_name").(string))

	resp, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.CustomDomainsApi.CreateDomain(ctx, d.Get("environment_id").(string)).CustomDomain(customDomain).Execute()
		},
		"CreateDomain",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.CustomDomain)

	d.SetId(respObject.GetId())

	return resourceCustomDomainRead(ctx, d, meta)
}

func resourceCustomDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.CustomDomainsApi.ReadOneDomain(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneDomain",
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

	respObject := resp.(*management.CustomDomain)

	d.Set("domain_name", respObject.GetDomainName())
	d.Set("status", string(respObject.GetStatus()))

	if v, ok := respObject.GetCanonicalNameOk(); ok {
		d.Set("canonical_name", v)
	} else {
		d.Set("canonical_name", nil)
	}

	if v, ok := respObject.GetCertificateOk(); ok {
		d.Set("certificate_expires_at", v.GetExpiresAt().Format(time.RFC3339))
	} else {
		d.Set("certificate_expires_at", nil)
	}

	return diags
}

func resourceCustomDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.CustomDomainsApi.DeleteDomain(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteDomain",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceCustomDomainImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/customDomainID\"", d.Id())
	}

	environmentID, customDomainID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)

	d.SetId(customDomainID)

	resourceCustomDomainRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
