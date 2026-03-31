resource "pingone_davinci_connector_instance" "flowConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "flowConnector"
  }
  name = "My awesome flowConnector"
  property {
    name  = "authenticatedRequest"
    type  = "string"
    value = var.flowconnector_property_authenticated_request
  }
  property {
    name  = "challengeExpiry"
    type  = "string"
    value = var.flowconnector_property_challenge_expiry
  }
  property {
    name  = "challengeLength"
    type  = "string"
    value = var.flowconnector_property_challenge_length
  }
  property {
    name  = "claimsMapping"
    type  = "string"
    value = var.flowconnector_property_claims_mapping
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.flowconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.flowconnector_property_client_secret
  }
  property {
    name  = "connectionId"
    type  = "string"
    value = var.flowconnector_property_connection_id
  }
  property {
    name  = "connectionInstanceId"
    type  = "string"
    value = var.flowconnector_property_connection_instance_id
  }
  property {
    name  = "customLink"
    type  = "string"
    value = var.flowconnector_property_custom_link
  }
  property {
    name  = "enforcedSignedToken"
    type  = "string"
    value = var.flowconnector_property_enforced_signed_token
  }
  property {
    name  = "generateQr"
    type  = "string"
    value = var.flowconnector_property_generate_qr
  }
  property {
    name  = "inputSchema"
    type  = "string"
    value = var.flowconnector_property_input_schema
  }
  property {
    name  = "issuerUrl"
    type  = "string"
    value = var.flowconnector_property_issuer_url
  }
  property {
    name  = "jwksKeys"
    type  = "string"
    value = var.flowconnector_property_jwks_keys
  }
  property {
    name  = "linkSubFlow"
    type  = "string"
    value = var.flowconnector_property_link_sub_flow
  }
  property {
    name  = "pemPublicKey"
    type  = "string"
    value = var.flowconnector_property_pem_public_key
  }
  property {
    name  = "popOutButton"
    type  = "string"
    value = var.flowconnector_property_pop_out_button
  }
  property {
    name  = "subFlowId"
    type  = "string"
    value = var.flowconnector_property_sub_flow_id
  }
  property {
    name  = "subFlowVersionId"
    type  = "string"
    value = var.flowconnector_property_sub_flow_version_id
  }
  property {
    name  = "tokenHint"
    type  = "string"
    value = var.flowconnector_property_token_hint
  }
  property {
    name  = "tokenSigningMethod"
    type  = "string"
    value = var.flowconnector_property_token_signing_method
  }
  property {
    name  = "useCustomLink"
    type  = "string"
    value = var.flowconnector_property_use_custom_link
  }
}
