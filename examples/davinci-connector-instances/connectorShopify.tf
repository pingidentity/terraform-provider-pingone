resource "pingone_davinci_connector_instance" "connectorShopify" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorShopify"
  }
  name = "My awesome connectorShopify"
  property {
    name  = "accessToken"
    type  = "string"
    value = var.connectorshopify_property_access_token
  }
  property {
    name  = "address"
    type  = "string"
    value = var.connectorshopify_property_address
  }
  property {
    name  = "address2"
    type  = "string"
    value = var.connectorshopify_property_address2
  }
  property {
    name  = "apiVersion"
    type  = "string"
    value = var.connectorshopify_property_api_version
  }
  property {
    name  = "body"
    type  = "string"
    value = var.connectorshopify_property_body
  }
  property {
    name  = "city"
    type  = "string"
    value = var.connectorshopify_property_city
  }
  property {
    name  = "country"
    type  = "string"
    value = var.connectorshopify_property_country
  }
  property {
    name  = "customerID"
    type  = "string"
    value = var.connectorshopify_property_customer_id
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connectorshopify_property_email
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectorshopify_property_endpoint
  }
  property {
    name  = "fname"
    type  = "string"
    value = var.connectorshopify_property_fname
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectorshopify_property_headers
  }
  property {
    name  = "listonlyCreatedAfter"
    type  = "string"
    value = var.connectorshopify_property_listonly_created_after
  }
  property {
    name  = "listonlyCreatedBefore"
    type  = "string"
    value = var.connectorshopify_property_listonly_created_before
  }
  property {
    name  = "listonlyLimit"
    type  = "string"
    value = var.connectorshopify_property_listonly_limit
  }
  property {
    name  = "listonlyUpdatedAfter"
    type  = "string"
    value = var.connectorshopify_property_listonly_updated_after
  }
  property {
    name  = "listonlyUpdatedBefore"
    type  = "string"
    value = var.connectorshopify_property_listonly_updated_before
  }
  property {
    name  = "lname"
    type  = "string"
    value = var.connectorshopify_property_lname
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectorshopify_property_method
  }
  property {
    name  = "multipassSecret"
    type  = "string"
    value = var.connectorshopify_property_multipass_secret
  }
  property {
    name  = "multipassStoreDomain"
    type  = "string"
    value = var.connectorshopify_property_multipass_store_domain
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.connectorshopify_property_phone
  }
  property {
    name  = "province"
    type  = "string"
    value = var.connectorshopify_property_province
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectorshopify_property_query_parameters
  }
  property {
    name  = "remoteIp"
    type  = "string"
    value = var.connectorshopify_property_remote_ip
  }
  property {
    name  = "returnTo"
    type  = "string"
    value = var.connectorshopify_property_return_to
  }
  property {
    name  = "verified"
    type  = "string"
    value = var.connectorshopify_property_verified
  }
  property {
    name  = "yourStoreName"
    type  = "string"
    value = var.connectorshopify_property_your_store_name
  }
  property {
    name  = "zip"
    type  = "string"
    value = var.connectorshopify_property_zip
  }
}
