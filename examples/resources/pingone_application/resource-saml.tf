resource "pingone_application" "my_awesome_saml_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome SAML App"
  enabled        = true

  saml_options {
    acs_urls           = ["https://my-saas-app.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:localhost"

    sp_verification_certificate_ids = [var.sp_verification_certificate_id]
  }
}
