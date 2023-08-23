package base

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceImage() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne images.",

		CreateContext: resourceImageCreate,
		ReadContext:   resourceImageRead,
		DeleteContext: resourceImageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceImageImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the image in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"image_file_base64": {
				Description: "A base64 encoded image file to import.  Only PNG, GIF and JPG images are supported.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"uploaded_image": {
				Description: "A block that specifies the processed image details.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"width": {
							Description: "The width of the image (in pixels).",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"height": {
							Description: "The height of the image (in pixels).",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "A string that specifies the type of format used for the image. Options are `jpg`, `png`, and `gif`.",
							Computed:    true,
						},
						"href": {
							Type:        schema.TypeString,
							Description: "A string that specifies the URL or fully qualified path to the image source file.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	var resp interface{}

	var archive []byte

	var err error
	archive, err = base64.StdEncoding.DecodeString(d.Get("image_file_base64").(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot base64 decode provided image file.",
			Detail:   fmt.Sprintf("The file cannot be base64 decoded: %s", err),
		})

		return diags
	}

	chars := 12
	generatedName, err := utils.RandStringFromCharSet(chars, "abcdefghijklmnopqrstuvwxyz012346789")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot generate a filename to use",
			Detail:   fmt.Sprintf("A random filename cannot be generated: %s", err),
		})

		return diags
	}

	extensionMapping := map[string]string{
		"image/jpeg": "jpg",
		"image/gif":  "gif",
		"image/png":  "png",
	}

	contentType := http.DetectContentType(archive)
	if _, ok := extensionMapping[contentType]; !ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot determine the content type of the image.  Ensure the file is a jpg, gif or png format.",
			Detail:   fmt.Sprintf("The file type has been determined to be `%s`, which is not supported.", contentType),
		})

		return diags
	}

	fileName := fmt.Sprintf("%s.%s", generatedName, extensionMapping[contentType])

	resp, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ImagesApi.CreateImage(ctx, d.Get("environment_id").(string)).ContentType(contentType).ContentDisposition(fmt.Sprintf("attachment; filename=%s", fileName)).File(&archive).Execute()
		},
		"CreateImage",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.Image)

	d.SetId(respObject.GetId())

	return resourceImageRead(ctx, d, meta)
}

func resourceImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ImagesApi.ReadImage(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadImage",
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

	respObject := resp.(*management.Image)

	d.Set("uploaded_image", flattenImageTarget(respObject.GetTargets()))

	return diags
}

func resourceImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			r, err := apiClient.ImagesApi.DeleteImage(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteImage",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceImageImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "image_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
	}

	attributes, err := framework.ParseImportID(d.Id(), idComponents...)
	if err != nil {
		return nil, err
	}

	d.Set("environment_id", attributes["environment_id"])
	d.SetId(attributes["image_id"])

	resourceImageRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func flattenImageTarget(v management.ImageTargets) []interface{} {

	original := v.GetOriginal()

	items := make([]interface{}, 0)
	return append(items, map[string]interface{}{
		"width":  original.GetWidth(),
		"height": original.GetHeight(),
		"type":   original.GetType(),
		"href":   original.GetHref(),
	})
}
