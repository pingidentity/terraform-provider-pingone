---
page_title: "pingone_forms_recaptcha_v2 Resource - terraform-provider-pingone"
subcategory: "Platform"
description: |-
  Resource to manage reCAPTCHA v2 configuration in a PingOne environment, for the purpose of including a reCAPTCHA v2 field in form definitions.
---

# pingone_forms_recaptcha_v2 (Resource)

Resource to manage reCAPTCHA v2 configuration in a PingOne environment, for the purpose of including a reCAPTCHA v2 field in form definitions.

~> Before you can define your environment's Recaptcha configuration in PingOne, you must register your site domain with Google to create a Recaptcha configuration. For more information about how to create the Google configuration for your domain, see [Google reCAPTCHA](https://www.google.com/recaptcha/admin/create). After you complete the configuration, Google provides the secret key and site key that you need to set your environment's PingOne Recaptcha configuration.

## Example Usage

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_forms_recaptcha_v2" "my_awesome_recaptcha_config" {
  environment_id = pingone_environment.my_environment.id

  site_key   = var.google_recaptcha_site_key
  secret_key = var.google_recaptcha_secret_key
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to manage the reCAPTCHA v2 configuration in.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `secret_key` (String, Sensitive) A string that specifies the confidential secret key for the Recaptcha configuration provided by Google. This is a required property.
- `site_key` (String) A string that specifies the public site key for the Recaptcha configuration provided by Google. This is a required property.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_forms_recaptcha_v2.example <environment_id>
```
