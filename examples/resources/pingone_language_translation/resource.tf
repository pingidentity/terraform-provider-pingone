resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_language" "my_customers_language" {
  environment_id = pingone_environment.my_environment.id

  locale = "sv"
}


resource "pingone_language_translation" "my_customers_language_translation" {
  environment_id = pingone_environment.my_environment.id

  locale = my_customers_language.locale
  translations = [
    {
      key             = "flow-ui.button.createNewAccount"
      translated_text = "Skapa ett nytt konto"
    },
    {
      key             = "flow-ui.label.email"
      translated_text = "E-post"
    },
  ]
}
