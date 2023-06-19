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

func ResourceApplicationPushCredential() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage push credentials for a mobile MFA application configured in PingOne.",

		CreateContext: resourcePingOneApplicationPushCredentialCreate,
		ReadContext:   resourcePingOneApplicationPushCredentialRead,
		UpdateContext: resourcePingOneApplicationPushCredentialUpdate,
		DeleteContext: resourcePingOneApplicationPushCredentialDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePingOneApplicationPushCredentialImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the application push notification credential in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"application_id": {
				Description:      "The ID of the application to create the push notification credential for.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"fcm": {
				Description:  "A block that specifies the credential settings for the Firebase Cloud Messaging service.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"fcm", "apns", "hms"},
				ForceNew:     true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Description:  "A string that represents the server key of the Firebase cloud messaging service.  One of `key` or `google_service_account_credentials` must be specified.",
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							Deprecated:   "This field is deprecated and will be removed in a future release.  Use `google_service_account_credentials` instead.",
							ExactlyOneOf: []string{"fcm.0.key", "fcm.0.google_service_account_credentials"},
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return old == "DUMMY_SUPPRESS_VALUE"
							},
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"google_service_account_credentials": {
							Description:  "A string in JSON format that represents the service account credentials of Firebase cloud messaging service.  One of `key` or `google_service_account_credentials` must be specified.",
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ExactlyOneOf: []string{"fcm.0.key", "fcm.0.google_service_account_credentials"},
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return old == "DUMMY_SUPPRESS_VALUE"
							},
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
						},
					},
				},
			},
			"apns": {
				Description:  "A block that specifies the credential settings for the Apple Push Notification Service.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"fcm", "apns", "hms"},
				ForceNew:     true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Description: "A string that Apple uses as an identifier to identify an authentication key.",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return old == "DUMMY_SUPPRESS_VALUE"
							},
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"team_id": {
							Description: "A string that Apple uses as an identifier to identify teams.",
							Type:        schema.TypeString,
							Required:    true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return old == "DUMMY_SUPPRESS_VALUE"
							},
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"token_signing_key": {
							Description: "A string that Apple uses as the authentication token signing key to securely connect to APNS. This is the contents of a p8 file with a private key format.",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return old == "DUMMY_SUPPRESS_VALUE"
							},
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
					},
				},
			},
			"hms": {
				Description:  "A block that specifies the credential settings for Huawei Moble Service push messaging.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"fcm", "apns", "hms"},
				ForceNew:     true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Description: "A string that represents the OAuth 2.0 Client ID from the Huawei Developers API console.",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return old == "DUMMY_SUPPRESS_VALUE"
							},
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"client_secret": {
							Description: "A string that represents the client secret associated with the OAuth 2.0 Client ID.",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return old == "DUMMY_SUPPRESS_VALUE"
							},
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
					},
				},
			},
		},
	}
}

func resourcePingOneApplicationPushCredentialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	mfaPushCredentialRequest, diags := expandPushCredentialRequest(d)
	if diags.HasError() {
		return diags
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ApplicationsApplicationMFAPushCredentialsApi.CreateMFAPushCredential(ctx, d.Get("environment_id").(string), d.Get("application_id").(string)).CreateMFAPushCredentialRequest(*mfaPushCredentialRequest).Execute()
		},
		"CreateMFAPushCredential",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*mfa.MFAPushCredentialResponse)

	d.SetId(respObject.GetId())

	return resourcePingOneApplicationPushCredentialRead(ctx, d, meta)
}

func expandPushCredentialRequest(d *schema.ResourceData) (*mfa.CreateMFAPushCredentialRequest, diag.Diagnostics) {

	mfaPushCredentialRequest := &mfa.CreateMFAPushCredentialRequest{}
	var diags diag.Diagnostics

	if v, ok := d.GetOk("fcm"); ok {
		mfaPushCredentialRequest.MFAPushCredentialFCM, mfaPushCredentialRequest.MFAPushCredentialFCMHTTPV1, diags = expandPushCredentialRequestFCM(v)
	}

	if v, ok := d.GetOk("apns"); ok {
		mfaPushCredentialRequest.MFAPushCredentialAPNS, diags = expandPushCredentialRequestAPNS(v)
	}

	if v, ok := d.GetOk("hms"); ok {
		mfaPushCredentialRequest.MFAPushCredentialHMS, diags = expandPushCredentialRequestHMS(v)
	}

	return mfaPushCredentialRequest, diags
}

func expandPushCredentialRequestUpdate(d *schema.ResourceData) (*mfa.UpdateMFAPushCredentialRequest, diag.Diagnostics) {

	mfaPushCredentialRequest := &mfa.UpdateMFAPushCredentialRequest{}
	var diags diag.Diagnostics

	if v, ok := d.GetOk("fcm"); ok {
		mfaPushCredentialRequest.MFAPushCredentialFCM, mfaPushCredentialRequest.MFAPushCredentialFCMHTTPV1, diags = expandPushCredentialRequestFCM(v)
	}

	if v, ok := d.GetOk("apns"); ok {
		mfaPushCredentialRequest.MFAPushCredentialAPNS, diags = expandPushCredentialRequestAPNS(v)
	}

	if v, ok := d.GetOk("hms"); ok {
		mfaPushCredentialRequest.MFAPushCredentialHMS, diags = expandPushCredentialRequestHMS(v)
	}

	return mfaPushCredentialRequest, diags
}

func expandPushCredentialRequestFCM(c interface{}) (*mfa.MFAPushCredentialFCM, *mfa.MFAPushCredentialFCMHTTPV1, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := c.([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		if credentialKey, ok := vp["key"].(string); ok && credentialKey != "" {
			credential := mfa.NewMFAPushCredentialFCM(
				mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_FCM,
				credentialKey,
			)

			return credential, nil, diags
		}

		if credentialKey, ok := vp["google_service_account_credentials"].(string); ok && credentialKey != "" {
			credential := mfa.NewMFAPushCredentialFCMHTTPV1(
				mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_FCM_HTTP_V1,
				credentialKey,
			)

			return nil, credential, diags
		}

	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `fcm` must be defined when using the FCM push notification type",
	})

	return nil, nil, diags
}

func expandPushCredentialRequestAPNS(c interface{}) (*mfa.MFAPushCredentialAPNS, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := c.([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		credential := mfa.NewMFAPushCredentialAPNS(
			mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_APNS,
			vp["key"].(string),
			vp["team_id"].(string),
			vp["token_signing_key"].(string),
		)

		return credential, diags

	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `apns` must be defined when using the APNS push notification type",
	})

	return nil, diags
}

func expandPushCredentialRequestHMS(c interface{}) (*mfa.MFAPushCredentialHMS, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := c.([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		credential := mfa.NewMFAPushCredentialHMS(
			mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_HMS,
			vp["client_id"].(string),
			vp["client_secret"].(string),
		)

		return credential, diags

	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `hms` must be defined when using the HMS push notification type",
	})

	return nil, diags
}

func resourcePingOneApplicationPushCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ApplicationsApplicationMFAPushCredentialsApi.ReadOneMFAPushCredential(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
		},
		"ReadOneMFAPushCredential",
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

	respObject := resp.(*mfa.MFAPushCredentialResponse)

	d.Set("fcm", nil)
	d.Set("apns", nil)
	d.Set("hms", nil)

	if respObject.GetType() == mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_FCM {
		credential := make([]interface{}, 0)
		d.Set("fcm", append(credential, map[string]string{
			"key": "DUMMY_SUPPRESS_VALUE",
		}))
	}

	if respObject.GetType() == mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_FCM_HTTP_V1 {
		credential := make([]interface{}, 0)
		d.Set("fcm", append(credential, map[string]string{
			"google_service_account_credentials": "DUMMY_SUPPRESS_VALUE",
		}))
	}

	if respObject.GetType() == mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_APNS {
		credential := make([]interface{}, 0)
		d.Set("apns", append(credential, map[string]string{
			"key":               "DUMMY_SUPPRESS_VALUE",
			"team_id":           "DUMMY_SUPPRESS_VALUE",
			"token_signing_key": "DUMMY_SUPPRESS_VALUE",
		}))
	}

	if respObject.GetType() == mfa.ENUMMFAPUSHCREDENTIALATTRTYPE_HMS {
		credential := make([]interface{}, 0)
		d.Set("hms", append(credential, map[string]string{
			"client_id":     "DUMMY_SUPPRESS_VALUE",
			"client_secret": "DUMMY_SUPPRESS_VALUE",
		}))
	}

	return diags
}

func resourcePingOneApplicationPushCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	mfaPushCredentialRequest, diags := expandPushCredentialRequestUpdate(d)
	if diags.HasError() {
		return diags
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ApplicationsApplicationMFAPushCredentialsApi.UpdateMFAPushCredential(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).UpdateMFAPushCredentialRequest(*mfaPushCredentialRequest).Execute()
		},
		"UpdateMFAPushCredential",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourcePingOneApplicationPushCredentialRead(ctx, d, meta)
}

func resourcePingOneApplicationPushCredentialDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.ApplicationsApplicationMFAPushCredentialsApi.DeleteMFAPushCredential(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteMFAPushCredential",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourcePingOneApplicationPushCredentialImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 3
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/applicationID/pushCredentialID\"", d.Id())
	}

	environmentID, applicationID, pushCredentialID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("application_id", applicationID)
	d.SetId(pushCredentialID)

	resourcePingOneApplicationPushCredentialRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
