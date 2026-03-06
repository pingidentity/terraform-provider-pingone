resource "pingone_davinci_connector_instance" "authenticIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "authenticIdConnector"
  }
  name = "My awesome authenticIdConnector"
  properties = jsonencode({
    "accountAccessKey" = var.authenticidconnector_property_account_access_key
    "androidSDKLicenseKey" = var.authenticidconnector_property_android_sdk_license_key
    "apiUrl" = var.authenticidconnector_property_api_url
    "baseUrl" = var.authenticidconnector_property_base_url
    "clientCertificate" = var.authenticidconnector_property_client_certificate
    "clientKey" = var.authenticidconnector_property_client_key
    "iOSSDKLicenseKey" = var.authenticidconnector_property_ios_sdk_license_key
    "passphrase" = var.authenticidconnector_property_passphrase
    "secretToken" = var.authenticidconnector_property_secret_token
  })
}
