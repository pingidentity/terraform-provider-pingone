resource "pingone_davinci_connector_instance" "microsoftDynamicsCustomerInsightsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "microsoftDynamicsCustomerInsightsConnector"
  }
  name = "My awesome microsoftDynamicsCustomerInsightsConnector"
  property {
    name  = "baseURL"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_base_url
  }
  property {
    name  = "body"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_body
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_client_secret
  }
  property {
    name  = "contactEmail"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_contact_email
  }
  property {
    name  = "contactFirstName"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_contact_first_name
  }
  property {
    name  = "contactId"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_contact_id
  }
  property {
    name  = "contactLastName"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_contact_last_name
  }
  property {
    name  = "contactMobile"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_contact_mobile
  }
  property {
    name  = "contactNewEmail"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_contact_new_email
  }
  property {
    name  = "contactProperties"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_contact_properties
  }
  property {
    name  = "contactUpdateProperties"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_contact_update_properties
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_endpoint
  }
  property {
    name  = "environmentName"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_environment_name
  }
  property {
    name  = "grantType"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_grant_type
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_headers
  }
  property {
    name  = "leadEmail"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_lead_email
  }
  property {
    name  = "leadFirstName"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_lead_first_name
  }
  property {
    name  = "leadId"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_lead_id
  }
  property {
    name  = "leadLastName"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_lead_last_name
  }
  property {
    name  = "leadMobile"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_lead_mobile
  }
  property {
    name  = "leadNewEmail"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_lead_new_email
  }
  property {
    name  = "leadProperties"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_lead_properties
  }
  property {
    name  = "leadUpdateProperties"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_lead_update_properties
  }
  property {
    name  = "method"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_method
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_query_parameters
  }
  property {
    name  = "tenant"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_tenant
  }
  property {
    name  = "version"
    type  = "string"
    value = var.microsoftdynamicscustomerinsightsconnector_property_version
  }
}
