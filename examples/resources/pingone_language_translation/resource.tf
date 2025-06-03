resource "pingone_environment" "my_environment" {
  # ...
}

locals {
  language_locales = [
    "en", # English
    "de", # German
    "es", # Spanish
    "fr", # French
    "it", # Italian
  ]
}

resource "pingone_language_translation" "my_language_translation" {
  for_each       = toset(local.language_locales)
  environment_id = pingone_environment.my_environment.id

  key             = "flow-ui.button.createNewAccount"
  locale          = each.key
  translated_test = "This is translated text"
}
