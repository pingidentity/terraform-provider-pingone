resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_trust_framework_processor" "my_awesome_processor" {
  environment_id = pingone_environment.my_environment.id
  name           = "Account Number"
  description    = "My awesome Account Number processor"

  processor = {
    name = "Extract Account Number"
    type = "JSON_PATH"

    expression = "$.accountNo"
    value_type = {
      type = "STRING"
    }
  }
}