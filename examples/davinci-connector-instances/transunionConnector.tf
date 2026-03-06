resource "pingone_davinci_connector_instance" "transunionConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "transunionConnector"
  }
  name = "My awesome transunionConnector"
  properties = jsonencode({
    "apiUrl" = var.transunionconnector_property_api_url
    "docVerificationPassword" = var.transunionconnector_property_doc_verification_password
    "docVerificationPublicKey" = var.transunionconnector_property_doc_verification_public_key
    "docVerificationSecret" = var.transunionconnector_property_doc_verification_secret
    "docVerificationSiteId" = var.transunionconnector_property_doc_verification_site_id
    "docVerificationUsername" = var.transunionconnector_property_doc_verification_username
    "idVerificationPassword" = var.transunionconnector_property_id_verification_password
    "idVerificationPublicKey" = var.transunionconnector_property_id_verification_public_key
    "idVerificationSecret" = var.transunionconnector_property_id_verification_secret
    "idVerificationSiteId" = var.transunionconnector_property_id_verification_site_id
    "idVerificationUsername" = var.transunionconnector_property_id_verification_username
    "kbaPassword" = var.transunionconnector_property_kba_password
    "kbaPublicKey" = var.transunionconnector_property_kba_public_key
    "kbaSecret" = var.transunionconnector_property_kba_secret
    "kbaSiteId" = var.transunionconnector_property_kba_site_id
    "kbaUsername" = var.transunionconnector_property_kba_username
    "otpPassword" = var.transunionconnector_property_otp_password
    "otpPublicKey" = var.transunionconnector_property_otp_public_key
    "otpSecret" = var.transunionconnector_property_otp_secret
    "otpSiteId" = var.transunionconnector_property_otp_site_id
    "otpUsername" = var.transunionconnector_property_otp_username
  })
}
