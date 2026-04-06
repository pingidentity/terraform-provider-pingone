resource "pingone_davinci_connector_instance" "mitekConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "mitekConnector"
  }
  name = "My awesome mitekConnector"
  property {
    name  = "clientId"
    type  = "string"
    value = var.mitekconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.mitekconnector_property_client_secret
  }
  property {
    name  = "hostURL"
    type  = "string"
    value = var.mitekconnector_property_host_url
  }
  property {
    name  = "mitekAttachmentNames"
    type  = "string"
    value = var.mitekconnector_property_mitek_attachment_names
  }
  property {
    name  = "mitekDocumentTypes"
    type  = "string"
    value = var.mitekconnector_property_mitek_document_types
  }
  property {
    name  = "mitekEmail"
    type  = "string"
    value = var.mitekconnector_property_mitek_email
  }
  property {
    name  = "mitekEndDate"
    type  = "string"
    value = var.mitekconnector_property_mitek_end_date
  }
  property {
    name  = "mitekEnvironment"
    type  = "string"
    value = var.mitekconnector_property_mitek_environment
  }
  property {
    name  = "mitekLanguage"
    type  = "string"
    value = var.mitekconnector_property_mitek_language
  }
  property {
    name  = "mitekName"
    type  = "string"
    value = var.mitekconnector_property_mitek_name
  }
  property {
    name  = "mitekNote"
    type  = "string"
    value = var.mitekconnector_property_mitek_note
  }
  property {
    name  = "mitekPhone"
    type  = "string"
    value = var.mitekconnector_property_mitek_phone
  }
  property {
    name  = "mitekReference"
    type  = "string"
    value = var.mitekconnector_property_mitek_reference
  }
  property {
    name  = "mitekScopeAttachments"
    type  = "string"
    value = var.mitekconnector_property_mitek_scope_attachments
  }
  property {
    name  = "mitekScopeDocuments"
    type  = "string"
    value = var.mitekconnector_property_mitek_scope_documents
  }
  property {
    name  = "mitekScopeSelfie"
    type  = "string"
    value = var.mitekconnector_property_mitek_scope_selfie
  }
  property {
    name  = "mitekSendChoice"
    type  = "string"
    value = var.mitekconnector_property_mitek_send_choice
  }
  property {
    name  = "mitekTabsOrder"
    type  = "string"
    value = var.mitekconnector_property_mitek_tabs_order
  }
  property {
    name  = "mitekUseDefaults"
    type  = "string"
    value = var.mitekconnector_property_mitek_use_defaults
  }
  property {
    name  = "requstAPIVersion"
    type  = "string"
    value = var.mitekconnector_property_requst_apiversion
  }
  property {
    name  = "skWebhookUri"
    type  = "string"
    value = var.mitekconnector_property_sk_webhook_uri
  }
}
