---
page_title: "pingone_credential_type Resource - terraform-provider-pingone"
subcategory: "Neo (Verify & Credentials)"
description: |-
  Resource to create and manage the credential types used by compatible wallet applications.
  ~> You must ensure that any fields used in the card_design_template are defined appropriately in metadata.fields or errors occur when you attempt to create a credential of that type.
---

# pingone_credential_type (Resource)

Resource to create and manage the credential types used by compatible wallet applications.

~> You must ensure that any fields used in the `card_design_template` are defined appropriately in `metadata.fields` or errors occur when you attempt to create a credential of that type.

## Example Usage

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_image" "verifiedemployee-background_image" {
  environment_id    = pingone_environment.my_environment.id
  image_file_base64 = filebase64("./images/verifiedemployee_background.png")
}

resource "pingone_image" "verifiedemployee-logo_image" {
  environment_id    = pingone_environment.my_environment.id
  image_file_base64 = filebase64("./images/verifiedemployee_logo.png")
}

resource "pingone_credential_type" "verifiedemployee" {
  environment_id   = pingone_environment.my_environment.id
  title            = "VerifiedEmployee"
  description      = "Demo Proof of Employment"
  card_type        = "VerifiedEmployee"
  revoke_on_delete = true

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="$${cardColor}" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity="$${bgOpacityPercent}"></rect>
<image href="$${backgroundImage}" opacity="$${bgOpacityPercent}" height="301" rx="10" ry="10" width="589" x="75" y="160"></image>
<image href="$${logoImage}" x="42" y="43" height="90px" width="90px"></image>
<line y2="160" x2="695" y1="160" x1="42.5" stroke="$${textColor}"></line>
<text fill="$${textColor}" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text fill="$${textColor}" font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>  
EOT

  metadata = {
    name               = "VerifiedEmployee"
    description        = "Demo Proof of Employment"
    bg_opacity_percent = 100

    background_image = pingone_image.verifiedemployee-background_image.uploaded_image.href
    logo_image       = pingone_image.verifiedemployee-logo_image.uploaded_image.href

    card_color = "#ffffff"
    text_color = "#000000"

    fields = [
      {
        type       = "Directory Attribute"
        title      = "givenName"
        attribute  = "name.given"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "surname"
        attribute  = "name.family"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "jobTitle"
        attribute  = "title"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "displayName"
        attribute  = "displayName"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "mail"
        attribute  = "email"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "preferredLanguage"
        attribute  = "preferredLanguage"
        is_visible = false
      },
      {
        type         = "Directory Attribute"
        title        = "photo"
        attribute    = "photo"
        file_support = "REFERENCE_FILE"
        is_visible   = true
      },
      {
        type       = "Directory Attribute"
        title      = "id"
        attribute  = "id"
        is_visible = false
      }
    ]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `card_design_template` (String) An SVG formatted image containing placeholders for the credentials fields that need to be displayed in the image.
- `environment_id` (String) PingOne environment identifier (UUID) in which the credential type exists.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `metadata` (Attributes) Contains the names, data types, and other metadata related to the credential. (see [below for nested schema](#nestedatt--metadata))
- `title` (String) Title of the credential. Verification sites are expected to be able to request the issued credential from the compatible wallet app using the title.  This value aligns to `${cardTitle}` in the `card_design_template`.

### Optional

- `card_type` (String) A descriptor of the credential type. Can be non-identity types such as proof of employment or proof of insurance.
- `description` (String) A description of the credential type. This value aligns to `${cardSubtitle}` in the `card_design_template`.
- `management_mode` (String) Specifies the management mode of the credential type.  Options are `AUTOMATED`, `MANAGED`.  Defaults to `AUTOMATED`.
- `revoke_on_delete` (Boolean) A boolean that specifies whether a user's issued verifiable credentials are automatically revoked when a `credential_type`, `user`, or `environment` is deleted.  Defaults to `true`.

### Read-Only

- `created_at` (String) Date and time the object was created.
- `id` (String) The ID of this resource.
- `issuer_id` (String) The identifier (UUID) of the issuer of the credential, which is the `id` of the `credential_issuer_profile` defined in the `environment`.
- `updated_at` (String) Date and time the object was updated. Can be null.

<a id="nestedatt--metadata"></a>
### Nested Schema for `metadata`

Required:

- `fields` (Attributes List) In a credential, the information is stored as key-value pairs where `fields` defines those key-value pairs. Effectively, `fields.title` is the key and its value is `fields.value` or extracted from the PingOne Directory attribute named in `fields.attribute`. (see [below for nested schema](#nestedatt--metadata--fields))

Optional:

- `background_image` (String) The URL or fully qualified path to the image file used for the credential background.  This can be retrieved from the `uploaded_image.href` parameter of the `pingone_image` resource.  Image size must not exceed 50 KB.
- `bg_opacity_percent` (Number) A numnber indicating the percent opacity of the background image in the credential. High percentage opacity may make text on the credential difficult to read.
- `card_color` (String) A string containing a 6-digit hexadecimal color code specifying the color of the credential.
- `columns` (Number) Indicates a number (between 1-3) of columns to display visible fields on the credential.
- `description` (String) Description of the credential.
- `logo_image` (String) The URL or fully qualified path to the image file used for the credential logo.  This can be retrieved from the `uploaded_image.href` parameter of the `pingone_image` resource.  Image size must not exceed 25 KB.
- `name` (String) Name of the credential.
- `text_color` (String) A string containing a 6-digit hexadecimal color code specifying the color of the credential text.

Read-Only:

- `version` (Number) Number version of this credential.

<a id="nestedatt--metadata--fields"></a>
### Nested Schema for `metadata.fields`

Required:

- `type` (String) Specifies the type of data in the credential field.  Options are `Alphanumeric Text`, `Directory Attribute`, `Issued Timestamp`.

Optional:

- `attribute` (String) Name of the PingOne Directory attribute. Present if `field.type` is `Directory Attribute`.
- `file_support` (String) Specifies how an image is stored in the credential field.  Options are `BASE64_STRING`, `INCLUDE_FILE`, `REFERENCE_FILE`.
- `is_visible` (Boolean) Specifies whether the field should be visible to viewers of the credential.
- `required` (Boolean) Specifies whether the field is required for the credential.
- `title` (String) Descriptive text when showing the field.
- `value` (String) The text to appear on the credential for a `field.type` of `Alphanumeric Text`.

Read-Only:

- `id` (String) Identifier of the field formatted as `<fields.type> -> <fields.title>`.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_credential_type.example <environment_id>/<credential_type_id>
```
