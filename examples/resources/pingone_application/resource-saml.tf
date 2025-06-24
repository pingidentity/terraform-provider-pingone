resource "pingone_key" "my_awesome_key" {
  environment_id = pingone_environment.my_environment.id

  name                = "Example Signing Key"
  algorithm           = "RSA"
  key_length          = 4096
  signature_algorithm = "SHA512withRSA"
  subject_dn          = "CN=Example Signing Key, OU=BX Retail, O=BX Retail, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_application" "my_awesome_saml_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome SAML App"
  enabled        = true

  saml_options = {
    acs_urls           = ["https://my-saas-app.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:localhost"

    idp_signing_key = {
      key_id    = pingone_key.my_awesome_key.id
      algorithm = pingone_key.my_awesome_key.signature_algorithm
    }

    sp_verification = {
      certificate_ids      = [var.sp_verification_certificate_id]
      authn_request_signed = true
    }

    virtual_server_id_settings = {
      enabled = true
      virtual_server_ids = [
        {
          vs_id = "virtual-server-1"
        },
        {
          vs_id   = "virtual-server-2"
          default = true
        },
        {
          vs_id = "virtual-server-3"
        },
        {
          vs_id = "virtual-server-4"
        },
      ]
    }
  }
}
