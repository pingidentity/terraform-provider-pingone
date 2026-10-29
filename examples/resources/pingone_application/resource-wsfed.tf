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

resource "pingone_application" "my_awesome_wsfed" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome WS-Fed Application"
  enabled        = true

  login_page_url = "https://portal.office.com"

  wsfed_options = {
    domain_name                    = "my.office365.domain"
    home_page_url                  = "https://www.microsoft.com/en-ca/microsoft-365"
    reply_url                      = "https://login.microsoftonline.com/login.srf"
    slo_endpoint                   = "https://example.com/slo"
    subject_name_identifier_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"

    # Entra ID hybrid join (limited access feature)
    hybrid_join_device_attribute  = "deviceAuthenticatorId"
    sts_token_type                = "urn:oasis:names:tc:SAML:1.0:assertion"
    use_mex13                     = true
    user_attribute_username_match = "userName"

    # Windows Login passwordless (PingID PRT) support
    user_attribute_x509_match = "mail"

    valid_issuance_certificates = [
      {
        certificate_id = pingone_key.my_awesome_issuance_key.id
      },
    ]

    idp_signing_key = {
      key_id    = pingone_key.my_awesome_key.id
      algorithm = pingone_key.my_awesome_key.signature_algorithm
    }
  }
}

resource "pingone_key" "my_awesome_issuance_key" {
  environment_id = pingone_environment.my_environment.id

  name                = "Example Issuance Key"
  algorithm           = "RSA"
  key_length          = 2048
  signature_algorithm = "SHA256withRSA"
  usage_type          = "ISSUANCE"
  subject_dn          = "CN=Example Issuance Key, OU=BX Retail, O=BX Retail, L=, ST=, C=US"
  validity_period     = 365
}
