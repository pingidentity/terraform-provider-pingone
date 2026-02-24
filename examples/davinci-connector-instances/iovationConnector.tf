resource "pingone_davinci_connector_instance" "iovationConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "iovationConnector"
  }
  name = "My awesome iovationConnector"
  properties = jsonencode({
    "apiUrl" = var.iovationconnector_property_api_url
    "javascriptCdnUrl" = var.iovationconnector_property_javascript_cdn_url
    "subKey" = var.iovationconnector_property_sub_key
    "subscriberAccount" = var.iovationconnector_property_subscriber_account
    "subscriberId" = var.iovationconnector_property_subscriber_id
    "subscriberPasscode" = var.iovationconnector_property_subscriber_passcode
    "version" = var.iovationconnector_property_version
  })
}
