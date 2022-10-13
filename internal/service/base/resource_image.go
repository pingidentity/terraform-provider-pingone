package base

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
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
				Description: "A base64 encoded image file to import.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"image_type": {
				Description:      "Image type.  Options are `PNG`, `JPG` or `GIF`.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PNG", "JPG", "GIF"}, true)),
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
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var resp interface{}

	var archive []byte

	var err error
	archive, err = base64.StdEncoding.DecodeString(d.Get("image_file_base64").(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot base64 decode provided image file.",
		})

		return diags
	}

	chars := 12

	fileName := fmt.Sprintf("%s.%s", utils.RandStringFromCharSet(chars, "abcdefghijklmnopqrstuvwxyz012346789"), strings.ToLower(d.Get("image_type").(string)))

	contentTypes := map[string]string{
		"jpg": "image/jpeg",
		"gif": "image/gif",
		"png": "image/png",
	}

	resp, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.ImagesApi.CreateImage(ctx, d.Get("environment_id").(string)).ContentType(contentTypes[strings.ToLower(d.Get("image_type").(string))]).ContentDisposition(fmt.Sprintf("attachment; filename=%s", fileName)).File(&archive).Execute()
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
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
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
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.ImagesApi.DeleteImage(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteImage",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceImageImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/imageID\"", d.Id())
	}

	environmentID, imageID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)

	d.SetId(imageID)

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
