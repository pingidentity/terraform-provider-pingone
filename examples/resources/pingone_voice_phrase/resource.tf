resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_voice_phrase" "my_verify_voice_phrase" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Verify Voice Phrase for my Verify Policy"
}
