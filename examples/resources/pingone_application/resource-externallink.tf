resource "pingone_application" "my_awesome_saml_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome External App"
  enabled        = true

  external_link_options {
    home_page_url = "https://demo.bxretail.org/"
  }
}
