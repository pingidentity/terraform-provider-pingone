resource "pingone_davinci_connector_instance" "samlConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "samlConnector"
  }
  name = "My awesome samlConnector"
  property {
    name  = "assertionLifeTimeInSeconds"
    type  = "string"
    value = var.samlconnector_property_assertion_life_time_in_seconds
  }
  property {
    name  = "authnContextClassRef"
    type  = "string"
    value = var.samlconnector_property_authn_context_class_ref
  }
  property {
    name  = "claimsMapping"
    type  = "string"
    value = var.samlconnector_property_claims_mapping
  }
  property {
    name  = "claimsNameValuePairs"
    type  = "string"
    value = var.samlconnector_property_claims_name_value_pairs
  }
  property {
    name  = "connectionId"
    type  = "string"
    value = var.samlconnector_property_connection_id
  }
  property {
    name  = "digestAlgorithm"
    type  = "string"
    value = var.samlconnector_property_digest_algorithm
  }
  property {
    name  = "encryptAssertion"
    type  = "string"
    value = var.samlconnector_property_encrypt_assertion
  }
  property {
    name  = "encryptionAlgorithm"
    type  = "string"
    value = var.samlconnector_property_encryption_algorithm
  }
  property {
    name  = "includeAttributeNameFormat"
    type  = "string"
    value = var.samlconnector_property_include_attribute_name_format
  }
  property {
    name  = "keyEncryptionAlgorithm"
    type  = "string"
    value = var.samlconnector_property_key_encryption_algorithm
  }
  property {
    name  = "nameId"
    type  = "string"
    value = var.samlconnector_property_name_id
  }
  property {
    name  = "nameIdFormat"
    type  = "string"
    value = var.samlconnector_property_name_id_format
  }
  property {
    name  = "nameIdNameQualifier"
    type  = "string"
    value = var.samlconnector_property_name_id_name_qualifier
  }
  property {
    name  = "nameIdSPNameQualifier"
    type  = "string"
    value = var.samlconnector_property_name_id_spname_qualifier
  }
  property {
    name  = "nameIdSPProvidedID"
    type  = "string"
    value = var.samlconnector_property_name_id_spprovided_id
  }
  property {
    name  = "samlResponseStatus"
    type  = "string"
    value = var.samlconnector_property_saml_response_status
  }
  property {
    name  = "samlResponseStatusMessage"
    type  = "string"
    value = var.samlconnector_property_saml_response_status_message
  }
  property {
    name  = "sessionIndex"
    type  = "string"
    value = var.samlconnector_property_session_index
  }
  property {
    name  = "signResponse"
    type  = "string"
    value = var.samlconnector_property_sign_response
  }
  property {
    name  = "signatureAlgorithm"
    type  = "string"
    value = var.samlconnector_property_signature_algorithm
  }
  property {
    name  = "signatureNamespacePrefix"
    type  = "string"
    value = var.samlconnector_property_signature_namespace_prefix
  }
  property {
    name  = "spCertForEncryption"
    type  = "string"
    value = var.samlconnector_property_sp_cert_for_encryption
  }
  property {
    name  = "typedAttributes"
    type  = "string"
    value = var.samlconnector_property_typed_attributes
  }
}
