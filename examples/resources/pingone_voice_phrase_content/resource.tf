resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_voice_phrase" "my_verify_voice_phrase" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Verify Voice Phrase for my Verify Policy"
}

resource "pingone_voice_phrase_content" {
  environment_id  = pingone_environment.my_environment.id
  voice_phrase_id = pingone_voice_phrase.my_verify_voice_phrase.id
  locale          = "en"
  content         = "My voice content to be used in voice enrollment or verification."
}