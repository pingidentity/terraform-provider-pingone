data "pingone_verify_voice_phrase" "find_by_id_example" {
  environment_id  = var.environment_id
  voice_phrase_id = var.voice_phrase_id
}

data "pingone_verify_voice_phrase" "find_by_name_example" {
  environment_id = var.environment_id
  name           = "foo"
}
