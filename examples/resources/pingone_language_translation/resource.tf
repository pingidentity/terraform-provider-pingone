resource "pingone_environment" "my_environment" {
  # ...
}

locals {
  language_locales = {
    "en" = "Create new Account",      # English
    "es" = "Crear nueva cuenta",      # Spanish
    "fr" = "Créer un nouveau compte", # French
    "it" = "Crea un nuovo account",   # Italian
  }
}

resource "pingone_language_translation" "my_language_translation" {
  for_each       = local.language_locales
  environment_id = pingone_environment.my_environment.id

  key             = "flow-ui.button.createNewAccount"
  locale          = each.key
  translated_test = each.value
}
