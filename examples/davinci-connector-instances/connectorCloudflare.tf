resource "pingone_davinci_connector_instance" "connectorCloudflare" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorCloudflare"
  }
  name = "My awesome connectorCloudflare"
  property {
    name  = "accountId"
    type  = "string"
    value = var.connectorcloudflare_property_account_id
  }
  property {
    name  = "apiToken"
    type  = "string"
    value = var.connectorcloudflare_property_api_token
  }
  property {
    name  = "body"
    type  = "string"
    value = var.connectorcloudflare_property_body
  }
  property {
    name  = "domain"
    type  = "string"
    value = var.connectorcloudflare_property_domain
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectorcloudflare_property_endpoint
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectorcloudflare_property_headers
  }
  property {
    name  = "ip"
    type  = "string"
    value = var.connectorcloudflare_property_ip
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectorcloudflare_property_method
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectorcloudflare_property_query_parameters
  }
}
