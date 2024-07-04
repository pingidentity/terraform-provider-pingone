resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_population" "my_population" {
  environment_id = pingone_environment.my_environment.id
  # ...
}

resource "pingone_application" "my_awesome_native_app" {
  environment_id = pingone_environment.my_environment.id
  # ...
}

resource "pingone_digital_wallet_application" "my_digital_wallet_app" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_native_app.id
  # ...
}

resource "pingone_credential_type" "my_credential" {
  environment_id = pingone_environment.my_environment.id
  # ...
}

resource "pingone_credential_issuance_rule" "my_credential_issuance_rule" {
  environment_id                = pingone_environment.my_environment.id
  digital_wallet_application_id = pingone_digital_wallet_application.my_digital_wallet_app.id
  credential_type_id            = pingone_credential_type.my_credential.id

  status = "ACTIVE"

  filter = {
    population_ids = [pingone_population.my_population.id]
  }

  automation = {
    issue  = "ON_DEMAND"
    revoke = "ON_DEMAND"
    update = "PERIODIC"
  }

  notification = {
    methods = ["EMAIL", "SMS"]
    template = {
      locale  = "en"
      variant = "credential_issued_template_B"
    }
  }
}