---
page_title: "pingone_branding_theme Resource - terraform-provider-pingone"
subcategory: "Platform"
description: |-
  Resource to create and manage PingOne branding themes for an environment.
---

# pingone_branding_theme (Resource)

Resource to create and manage PingOne branding themes for an environment.

## Example Usage

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_image" "company_logo" {
  environment_id = pingone_environment.my_environment.id

  image_file_base64 = filebase64("../path/to/image.jpg")
}

resource "pingone_image" "theme_background" {
  environment_id = pingone_environment.my_environment.id

  image_file_base64 = filebase64("../path/to/background-image.jpg")
}

resource "pingone_branding_theme" "my_awesome_theme" {
  environment_id = pingone_environment.my_environment.id

  name     = "My Awesome Theme"
  template = "split"

  logo = {
    id   = pingone_image.company_logo.id
    href = pingone_image.company_logo.uploaded_image.href
  }

  background_image = {
    id   = pingone_image.theme_background.id
    href = pingone_image.theme_background.uploaded_image.href
  }

  button_text_color  = "#FFFFFF"
  heading_text_color = "#686F77"
  card_color         = "#FCFCFC"
  body_text_color    = "#263956"
  link_text_color    = "#263956"
  button_color       = "#263956"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `body_text_color` (String) The body text color for the theme. It must be a valid hexadecimal color code.
- `button_color` (String) The button color for the theme. It must be a valid hexadecimal color code.
- `button_text_color` (String) The button text color for the branding theme. It must be a valid hexadecimal color code.
- `card_color` (String) The card color for the branding theme. It must be a valid hexadecimal color code.
- `environment_id` (String) The ID of the environment to set branding settings for.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `heading_text_color` (String) The heading text color for the branding theme. It must be a valid hexadecimal color code.
- `link_text_color` (String) The hyperlink text color for the branding theme. It must be a valid hexadecimal color code.
- `name` (String) A string that specifies the unique name of the branding theme.
- `template` (String) The template name of the branding theme associated with the environment.  Options are `default`, `focus`, `mural`, `slate`, `split`.

### Optional

- `background_color` (String) The background color for the theme. It must be a valid hexadecimal color code.  Exactly one of the following must be defined: `background_image`, `background_color`, `use_default_background`.
- `background_image` (Attributes) A single object that specifies the HREF and ID for the background image.  Exactly one of the following must be defined: `background_image`, `background_color`, `use_default_background`. (see [below for nested schema](#nestedatt--background_image))
- `footer_text` (String) The text to be displayed in the footer of the branding theme.
- `logo` (Attributes) A single object that specifies the HREF and ID for the company logo, for this branding template.  If not set, the environment's default logo (set with the `pingone_branding_settings` resource) will be applied. (see [below for nested schema](#nestedatt--logo))
- `use_default_background` (Boolean) A boolean to specify that the background should be set to the theme template's default.  Exactly one of the following must be defined: `background_image`, `background_color`, `use_default_background`.

### Read-Only

- `default` (Boolean) Specifies whether this theme is the environment's default branding configuration.
- `id` (String) The ID of this resource.

<a id="nestedatt--background_image"></a>
### Nested Schema for `background_image`

Required:

- `href` (String) The URL or fully qualified path to the background image file used for branding.  This can be retrieved from the `uploaded_image.href` parameter of the `pingone_image` resource.
- `id` (String) The ID of the background image.  This can be retrieved from the `id` parameter of the `pingone_image` resource.  Must be a valid PingOne resource ID.


<a id="nestedatt--logo"></a>
### Nested Schema for `logo`

Required:

- `href` (String) The URL or fully qualified path to the logo file used for branding.  This can be retrieved from the `uploaded_image.href` parameter of the `pingone_image` resource.
- `id` (String) The ID of the logo image.  This can be retrieved from the `id` parameter of the `pingone_image` resource.  Must be a valid PingOne resource ID.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_branding_theme.example <environment_id>/<branding_theme_id>
```
