resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_language_translation" "my_language_translation" {
  environment_id = pingone_environment.my_environment.id

  key             = "myKey"
  locale          = "en"
  translated_test = "This is translated text"
}
