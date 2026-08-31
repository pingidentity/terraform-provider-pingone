resource "pingone_davinci_connector_instance" "socureriskosConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "socureriskosConnector"
  }
  name = "My awesome socureriskosConnector"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.socureriskosconnector_property_api_key
  }
  property {
    name  = "attachments"
    type  = "string"
    value = var.socureriskosconnector_property_attachments
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.socureriskosconnector_property_base_url
  }
  property {
    name  = "customer_metadata"
    type  = "string"
    value = var.socureriskosconnector_property_customer_metadata
  }
  property {
    name  = "data"
    type  = "string"
    value = var.socureriskosconnector_property_data
  }
  property {
    name  = "decision"
    type  = "string"
    value = var.socureriskosconnector_property_decision
  }
  property {
    name  = "eval_id"
    type  = "string"
    value = var.socureriskosconnector_property_eval_id
  }
  property {
    name  = "fraud_label"
    type  = "string"
    value = var.socureriskosconnector_property_fraud_label
  }
  property {
    name  = "fraud_type"
    type  = "string"
    value = var.socureriskosconnector_property_fraud_type
  }
  property {
    name  = "id"
    type  = "string"
    value = var.socureriskosconnector_property_id
  }
  property {
    name  = "notes"
    type  = "string"
    value = var.socureriskosconnector_property_notes
  }
  property {
    name  = "queue"
    type  = "string"
    value = var.socureriskosconnector_property_queue
  }
  property {
    name  = "recorded_at"
    type  = "string"
    value = var.socureriskosconnector_property_recorded_at
  }
  property {
    name  = "recorded_by"
    type  = "string"
    value = var.socureriskosconnector_property_recorded_by
  }
  property {
    name  = "status"
    type  = "string"
    value = var.socureriskosconnector_property_status
  }
  property {
    name  = "sub_status"
    type  = "string"
    value = var.socureriskosconnector_property_sub_status
  }
  property {
    name  = "tags"
    type  = "string"
    value = var.socureriskosconnector_property_tags
  }
  property {
    name  = "timestamp"
    type  = "string"
    value = var.socureriskosconnector_property_timestamp
  }
  property {
    name  = "workflow"
    type  = "string"
    value = var.socureriskosconnector_property_workflow
  }
}
