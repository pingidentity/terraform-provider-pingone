resource "pingone_davinci_connector_instance" "pingOneVerifyConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneVerifyConnector"
  }
  name = "My awesome pingOneVerifyConnector"
  property {
    name  = "aspectHeight"
    type  = "number"
    value = var.pingoneverifyconnector_property_aspect_height
  }
  property {
    name  = "aspectWidth"
    type  = "number"
    value = var.pingoneverifyconnector_property_aspect_width
  }
  property {
    name  = "attempt"
    type  = "string"
    value = var.pingoneverifyconnector_property_attempt
  }
  property {
    name  = "biographic"
    type  = "string"
    value = var.pingoneverifyconnector_property_biographic
  }
  property {
    name  = "challenge"
    type  = "string"
    value = var.pingoneverifyconnector_property_challenge
  }
  property {
    name  = "challengeId"
    type  = "string"
    value = var.pingoneverifyconnector_property_challenge_id
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.pingone_worker_app_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.pingone_worker_app_client_secret
  }
  property {
    name  = "colorPicker"
    type  = "string"
    value = var.pingoneverifyconnector_property_color_picker
  }
  property {
    name  = "deviceOsType"
    type  = "string"
    value = var.pingoneverifyconnector_property_device_os_type
  }
  property {
    name  = "documentId"
    type  = "string"
    value = var.pingoneverifyconnector_property_document_id
  }
  property {
    name  = "documentTypeName"
    type  = "string"
    value = var.pingoneverifyconnector_property_document_type_name
  }
  property {
    name  = "documentTypes"
    type  = "string"
    value = var.pingoneverifyconnector_property_document_types
  }
  property {
    name  = "documentValue"
    type  = "string"
    value = var.pingoneverifyconnector_property_document_value
  }
  property {
    name  = "envId"
    type  = "string"
    value = var.pingone_worker_app_environment_id
  }
  property {
    name  = "isLastClientStep"
    type  = "string"
    value = var.pingoneverifyconnector_property_is_last_client_step
  }
  property {
    name  = "limit"
    type  = "number"
    value = var.pingoneverifyconnector_property_limit
  }
  property {
    name  = "metadataType"
    type  = "string"
    value = var.pingoneverifyconnector_property_metadata_type
  }
  property {
    name  = "notifyEmail"
    type  = "string"
    value = var.pingoneverifyconnector_property_notify_email
  }
  property {
    name  = "notifyPhone"
    type  = "string"
    value = var.pingoneverifyconnector_property_notify_phone
  }
  property {
    name  = "probeBiographic"
    type  = "string"
    value = var.pingoneverifyconnector_property_probe_biographic
  }
  property {
    name  = "redirectMessage"
    type  = "string"
    value = var.pingoneverifyconnector_property_redirect_message
  }
  property {
    name  = "redirectUrl"
    type  = "string"
    value = var.pingoneverifyconnector_property_redirect_url
  }
  property {
    name  = "referenceImage"
    type  = "string"
    value = var.pingoneverifyconnector_property_reference_image
  }
  property {
    name  = "region"
    type  = "string"
    value = var.pingoneverifyconnector_property_region
  }
  property {
    name  = "selfieId"
    type  = "string"
    value = var.pingoneverifyconnector_property_selfie_id
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.pingoneverifyconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.pingoneverifyconnector_property_skip_button_press
  }
  property {
    name  = "transactionId"
    type  = "string"
    value = var.pingoneverifyconnector_property_transaction_id
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.pingoneverifyconnector_property_user_id
  }
  property {
    name  = "verifiedType"
    type  = "string"
    value = var.pingoneverifyconnector_property_verified_type
  }
  property {
    name  = "verifyEmail"
    type  = "string"
    value = var.pingoneverifyconnector_property_verify_email
  }
  property {
    name  = "verifyPhone"
    type  = "string"
    value = var.pingoneverifyconnector_property_verify_phone
  }
  property {
    name  = "verifyPolicy"
    type  = "string"
    value = var.pingoneverifyconnector_property_verify_policy
  }
  property {
    name  = "verifyPolicyId"
    type  = "string"
    value = var.pingoneverifyconnector_property_verify_policy_id
  }
  property {
    name  = "verifyPolicyIdSelect"
    type  = "string"
    value = var.pingoneverifyconnector_property_verify_policy_id_select
  }
  property {
    name  = "verifyPolicySelect"
    type  = "string"
    value = var.pingoneverifyconnector_property_verify_policy_select
  }
  property {
    name  = "verifyStatus"
    type  = "string"
    value = var.pingoneverifyconnector_property_verify_status
  }
  property {
    name  = "voiceSampleIndex"
    type  = "number"
    value = var.pingoneverifyconnector_property_voice_sample_index
  }
  property {
    name  = "webVerificationUrl"
    type  = "string"
    value = var.pingoneverifyconnector_property_web_verification_url
  }
}
