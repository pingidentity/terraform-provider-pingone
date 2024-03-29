---
page_title: "pingone_verify_voice_phrase_content Resource - terraform-provider-pingone"
subcategory: "Neo (Verify & Credentials)"
description: |-
  Resource to configure the phrases to speak during voice verification enrollment or validation.
---

# pingone_verify_voice_phrase_content (Resource)

Resource to configure the phrases to speak during voice verification enrollment or validation.

## Example Usage

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_verify_voice_phrase" "my_verify_voice_phrase" {
  environment_id = pingone_environment.my_environment.id
  display_name   = "My Awesome Verify Voice Phrase for my Verify Policy"
}

resource "pingone_verify_voice_phrase_content" "my_verify_voice_phrase_content" {
  environment_id  = pingone_environment.my_environment.id
  voice_phrase_id = pingone_verify_voice_phrase.my_verify_voice_phrase.id
  locale          = "en"
  content         = "My voice content to be used in voice enrollment or verification."
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `content` (String) The phrase a user must speak as part of the voice enrollment or verification. The phrase must be written in the language and character set required by the language specified in the `locale` property.
- `environment_id` (String) PingOne environment identifier (UUID) in which the verify voice phrase exists.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `locale` (String) Language localization requirement for the voice phrase contents.
- `voice_phrase_id` (String) The identifier (UUID) of the `voice_phrase` associated with the `voice_phrase_content` configuration.  This field is immutable and will trigger a replace plan if changed.

### Read-Only

- `created_at` (String) Date and time the verify phrase content was created.
- `id` (String) The ID of this resource.
- `updated_at` (String) Date and time the verify phrase content was updated. Can be null.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
$ terraform import pingone_verify_voice_phrase_content.example <environment_id>/<voice_phrase_id>/<voice_phrase_content_id>
```
