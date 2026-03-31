resource "pingone_davinci_connector_instance" "connectorHubspot" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorHubspot"
  }
  name = "My awesome connectorHubspot"
  property {
    name  = "associationType"
    type  = "string"
    value = var.connectorhubspot_property_association_type
  }
  property {
    name  = "bearerToken"
    type  = "string"
    value = var.connectorhubspot_property_bearer_token
  }
  property {
    name  = "body"
    type  = "string"
    value = var.connectorhubspot_property_body
  }
  property {
    name  = "company"
    type  = "string"
    value = var.connectorhubspot_property_company
  }
  property {
    name  = "contactID"
    type  = "string"
    value = var.connectorhubspot_property_contact_id
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connectorhubspot_property_email
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectorhubspot_property_endpoint
  }
  property {
    name  = "fname"
    type  = "string"
    value = var.connectorhubspot_property_fname
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectorhubspot_property_headers
  }
  property {
    name  = "hubspotOwnerID"
    type  = "string"
    value = var.connectorhubspot_property_hubspot_owner_id
  }
  property {
    name  = "lname"
    type  = "string"
    value = var.connectorhubspot_property_lname
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectorhubspot_property_method
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.connectorhubspot_property_phone
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectorhubspot_property_query_parameters
  }
  property {
    name  = "subject"
    type  = "string"
    value = var.connectorhubspot_property_subject
  }
  property {
    name  = "ticketID"
    type  = "string"
    value = var.connectorhubspot_property_ticket_id
  }
  property {
    name  = "ticketPipeline"
    type  = "string"
    value = var.connectorhubspot_property_ticket_pipeline
  }
  property {
    name  = "ticketPriority"
    type  = "string"
    value = var.connectorhubspot_property_ticket_priority
  }
  property {
    name  = "ticketStage"
    type  = "string"
    value = var.connectorhubspot_property_ticket_stage
  }
  property {
    name  = "toObjectID"
    type  = "string"
    value = var.connectorhubspot_property_to_object_id
  }
  property {
    name  = "toObjectType"
    type  = "string"
    value = var.connectorhubspot_property_to_object_type
  }
  property {
    name  = "website"
    type  = "string"
    value = var.connectorhubspot_property_website
  }
}
