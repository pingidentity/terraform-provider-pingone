resource "pingone_davinci_connector_instance" "connector-oai-pfadminapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-pfadminapi"
  }
  name = "My awesome connector-oai-pfadminapi"
  property {
    name  = "_export_body_ExportRequestSigningSettingsSigningKeyPairRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property__export_body_export_request_signing_settings_signing_key_pair_ref_id
  }
  property {
    name  = "_export_body_ExportRequestSigningSettingsSigningKeyPairRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property__export_body_export_request_signing_settings_signing_key_pair_ref_location
  }
  property {
    name  = "_export_body_ExportRequestSigningSettings_algorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property__export_body_export_request_signing_settings_algorithm
  }
  property {
    name  = "_export_body_ExportRequestSigningSettings_includeCertInSignature"
    type  = "string"
    value = var.connector-oai-pfadminapi_property__export_body_export_request_signing_settings_include_cert_in_signature
  }
  property {
    name  = "_export_body_ExportRequestSigningSettings_includeRawKeyInSignature"
    type  = "string"
    value = var.connector-oai-pfadminapi_property__export_body_export_request_signing_settings_include_raw_key_in_signature
  }
  property {
    name  = "_export_body_ExportRequest_connectionId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property__export_body_export_request_connection_id
  }
  property {
    name  = "_export_body_ExportRequest_connectionType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property__export_body_export_request_connection_type
  }
  property {
    name  = "_export_body_ExportRequest_useSecondaryPortForSoap"
    type  = "string"
    value = var.connector-oai-pfadminapi_property__export_body_export_request_use_secondary_port_for_soap
  }
  property {
    name  = "_export_body_ExportRequest_virtualHostName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property__export_body_export_request_virtual_host_name
  }
  property {
    name  = "_export_body_ExportRequest_virtualServerId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property__export_body_export_request_virtual_server_id
  }
  property {
    name  = "addAccount_body_GetAccounts200ResponseItemsInner_active"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_account_body_get_accounts200_response_items_inner_active
  }
  property {
    name  = "addAccount_body_GetAccounts200ResponseItemsInner_auditor"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_account_body_get_accounts200_response_items_inner_auditor
  }
  property {
    name  = "addAccount_body_GetAccounts200ResponseItemsInner_department"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_account_body_get_accounts200_response_items_inner_department
  }
  property {
    name  = "addAccount_body_GetAccounts200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_account_body_get_accounts200_response_items_inner_description
  }
  property {
    name  = "addAccount_body_GetAccounts200ResponseItemsInner_emailAddress"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_account_body_get_accounts200_response_items_inner_email_address
  }
  property {
    name  = "addAccount_body_GetAccounts200ResponseItemsInner_encryptedPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_account_body_get_accounts200_response_items_inner_encrypted_password
  }
  property {
    name  = "addAccount_body_GetAccounts200ResponseItemsInner_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_account_body_get_accounts200_response_items_inner_password
  }
  property {
    name  = "addAccount_body_GetAccounts200ResponseItemsInner_phoneNumber"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_account_body_get_accounts200_response_items_inner_phone_number
  }
  property {
    name  = "addAccount_body_GetAccounts200ResponseItemsInner_roles"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_account_body_get_accounts200_response_items_inner_roles
  }
  property {
    name  = "addAccount_body_GetAccounts200ResponseItemsInner_username"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_account_body_get_accounts200_response_items_inner_username
  }
  property {
    name  = "addCommonScopeGroup_body_GetAuthorizationServerSettings200ResponseScopeGroupsInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_common_scope_group_body_get_authorization_server_settings200_response_scope_groups_inner_description
  }
  property {
    name  = "addCommonScopeGroup_body_GetAuthorizationServerSettings200ResponseScopeGroupsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_common_scope_group_body_get_authorization_server_settings200_response_scope_groups_inner_name
  }
  property {
    name  = "addCommonScopeGroup_body_GetAuthorizationServerSettings200ResponseScopeGroupsInner_scopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_common_scope_group_body_get_authorization_server_settings200_response_scope_groups_inner_scopes
  }
  property {
    name  = "addCommonScope_body_GetAuthorizationServerSettings200ResponseScopesInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_common_scope_body_get_authorization_server_settings200_response_scopes_inner_description
  }
  property {
    name  = "addCommonScope_body_GetAuthorizationServerSettings200ResponseScopesInner_dynamic"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_common_scope_body_get_authorization_server_settings200_response_scopes_inner_dynamic
  }
  property {
    name  = "addCommonScope_body_GetAuthorizationServerSettings200ResponseScopesInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_common_scope_body_get_authorization_server_settings200_response_scopes_inner_name
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_crypto_provider
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_expires"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_expires
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_id
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_issuerDN"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_issuer_dn
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_keyAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_key_algorithm
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_keySize"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_key_size
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_serialNumber"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_serial_number
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_sha1Fingerprint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_sha1_fingerprint
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_sha256Fingerprint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_sha256_fingerprint
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_signatureAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_signature_algorithm
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_status"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_status
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_subjectAlternativeNames"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_subject_alternative_names
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_subjectDN"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_subject_dn
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_validFrom"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_valid_from
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_version"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_cert_view_version
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerX509File_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_x509_file_crypto_provider
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerX509File_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_x509_file_file_data
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInnerX509File_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_x509_file_id
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInner_activeVerificationCert"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_active_verification_cert
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInner_encryptionCert"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_encryption_cert
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInner_primaryVerificationCert"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_primary_verification_cert
  }
  property {
    name  = "addConnectionCert1_body_ConvertRequestTemplateConnectionCredentialsCertsInner_secondaryVerificationCert"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_body_convert_request_template_connection_credentials_certs_inner_secondary_verification_cert
  }
  property {
    name  = "addConnectionCert1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert1_id
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_crypto_provider
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_expires"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_expires
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_id
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_issuerDN"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_issuer_dn
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_keyAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_key_algorithm
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_keySize"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_key_size
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_serialNumber"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_serial_number
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_sha1Fingerprint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_sha1_fingerprint
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_sha256Fingerprint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_sha256_fingerprint
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_signatureAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_signature_algorithm
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_status"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_status
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_subjectAlternativeNames"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_subject_alternative_names
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_subjectDN"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_subject_dn
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_validFrom"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_valid_from
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerCertView_version"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_cert_view_version
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerX509File_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_x509_file_crypto_provider
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerX509File_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_x509_file_file_data
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInnerX509File_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_x509_file_id
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInner_activeVerificationCert"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_active_verification_cert
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInner_encryptionCert"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_encryption_cert
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInner_primaryVerificationCert"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_primary_verification_cert
  }
  property {
    name  = "addConnectionCert2_body_ConvertRequestTemplateConnectionCredentialsCertsInner_secondaryVerificationCert"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_body_convert_request_template_connection_credentials_certs_inner_secondary_verification_cert
  }
  property {
    name  = "addConnectionCert2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_connection_cert2_id
  }
  property {
    name  = "addExclusiveScopeGroup_body_GetAuthorizationServerSettings200ResponseScopeGroupsInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_exclusive_scope_group_body_get_authorization_server_settings200_response_scope_groups_inner_description
  }
  property {
    name  = "addExclusiveScopeGroup_body_GetAuthorizationServerSettings200ResponseScopeGroupsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_exclusive_scope_group_body_get_authorization_server_settings200_response_scope_groups_inner_name
  }
  property {
    name  = "addExclusiveScopeGroup_body_GetAuthorizationServerSettings200ResponseScopeGroupsInner_scopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_exclusive_scope_group_body_get_authorization_server_settings200_response_scope_groups_inner_scopes
  }
  property {
    name  = "addExclusiveScope_body_GetAuthorizationServerSettings200ResponseScopesInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_exclusive_scope_body_get_authorization_server_settings200_response_scopes_inner_description
  }
  property {
    name  = "addExclusiveScope_body_GetAuthorizationServerSettings200ResponseScopesInner_dynamic"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_exclusive_scope_body_get_authorization_server_settings200_response_scopes_inner_dynamic
  }
  property {
    name  = "addExclusiveScope_body_GetAuthorizationServerSettings200ResponseScopesInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_exclusive_scope_body_get_authorization_server_settings200_response_scopes_inner_name
  }
  property {
    name  = "addIssuer_body_GetIssuerById200Response_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_issuer_body_get_issuer_by_id200_response_description
  }
  property {
    name  = "addIssuer_body_GetIssuerById200Response_host"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_issuer_body_get_issuer_by_id200_response_host
  }
  property {
    name  = "addIssuer_body_GetIssuerById200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_issuer_body_get_issuer_by_id200_response_id
  }
  property {
    name  = "addIssuer_body_GetIssuerById200Response_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_issuer_body_get_issuer_by_id200_response_name
  }
  property {
    name  = "addIssuer_body_GetIssuerById200Response_path"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_issuer_body_get_issuer_by_id200_response_path
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_crypto_provider
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_expires"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_expires
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_id
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_issuerDN"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_issuer_dn
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_keyAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_key_algorithm
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_keySize"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_key_size
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_serialNumber"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_serial_number
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_sha1Fingerprint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_sha1_fingerprint
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_sha256Fingerprint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_sha256_fingerprint
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_signatureAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_signature_algorithm
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_status"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_status
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_subjectAlternativeNames"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_subject_alternative_names
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_subjectDN"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_subject_dn
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_validFrom"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_valid_from
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_version"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_version
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerX509File_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_x509_file_crypto_provider
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerX509File_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_x509_file_file_data
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerX509File_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_x509_file_id
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_id
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_name
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInner_url"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_url
  }
  property {
    name  = "addMetadataUrl_body_GetMetadataUrls200ResponseItemsInner_validateSignature"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_add_metadata_url_body_get_metadata_urls200_response_items_inner_validate_signature
  }
  property {
    name  = "authPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_auth_password
  }
  property {
    name  = "authUsername"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_auth_username
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_base_path
  }
  property {
    name  = "changePassword_body_ResetPasswordRequest_currentPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_change_password_body_reset_password_request_current_password
  }
  property {
    name  = "changePassword_body_ResetPasswordRequest_newPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_change_password_body_reset_password_request_new_password
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionAdditionalAllowedEntitiesConfiguration_additionalAllowedEntities"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_additional_allowed_entities_configuration_additional_allowed_entities
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionAdditionalAllowedEntitiesConfiguration_allowAdditionalEntities"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_additional_allowed_entities_configuration_allow_additional_entities
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionAdditionalAllowedEntitiesConfiguration_allowAllEntities"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_additional_allowed_entities_configuration_allow_all_entities
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionContactInfo_company"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_contact_info_company
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionContactInfo_email"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_contact_info_email
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionContactInfo_firstName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_contact_info_first_name
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionContactInfo_lastName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_contact_info_last_name
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionContactInfo_phone"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_contact_info_phone
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsDecryptionKeyPairRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_decryption_key_pair_ref_id
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsDecryptionKeyPairRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_decryption_key_pair_ref_location
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsInboundBackChannelAuth_certs"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_inbound_back_channel_auth_certs
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsInboundBackChannelAuth_digitalSignature"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_inbound_back_channel_auth_digital_signature
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsInboundBackChannelAuth_requireSsl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_inbound_back_channel_auth_require_ssl
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsInboundBackChannelAuth_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_inbound_back_channel_auth_type
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsInboundBackChannelAuth_verificationIssuerDN"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_inbound_back_channel_auth_verification_issuer_dn
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsInboundBackChannelAuth_verificationSubjectDN"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_inbound_back_channel_auth_verification_subject_dn
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsOutboundBackChannelAuthAllOf1SslAuthKeyPairRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_outbound_back_channel_auth_all_of1_ssl_auth_key_pair_ref_id
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsOutboundBackChannelAuthAllOf1SslAuthKeyPairRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_outbound_back_channel_auth_all_of1_ssl_auth_key_pair_ref_location
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsOutboundBackChannelAuthAllOfHttpBasicCredentials_encryptedPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_outbound_back_channel_auth_all_of_http_basic_credentials_encrypted_password
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsOutboundBackChannelAuthAllOfHttpBasicCredentials_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_outbound_back_channel_auth_all_of_http_basic_credentials_password
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsOutboundBackChannelAuthAllOfHttpBasicCredentials_username"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_outbound_back_channel_auth_all_of_http_basic_credentials_username
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsOutboundBackChannelAuth_digitalSignature"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_outbound_back_channel_auth_digital_signature
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsOutboundBackChannelAuth_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_outbound_back_channel_auth_type
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsOutboundBackChannelAuth_validatePartnerCert"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_outbound_back_channel_auth_validate_partner_cert
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsSecondaryDecryptionKeyPairRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_secondary_decryption_key_pair_ref_id
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsSecondaryDecryptionKeyPairRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_secondary_decryption_key_pair_ref_location
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsSigningSettings_algorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_signing_settings_algorithm
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsSigningSettings_alternativeSigningKeyPairRefs"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_signing_settings_alternative_signing_key_pair_refs
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsSigningSettings_includeCertInSignature"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_signing_settings_include_cert_in_signature
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentialsSigningSettings_includeRawKeyInSignature"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_signing_settings_include_raw_key_in_signature
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentials_blockEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_block_encryption_algorithm
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentials_certs"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_certs
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentials_keyTransportAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_key_transport_algorithm
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentials_verificationIssuerDN"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_verification_issuer_dn
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionCredentials_verificationSubjectDN"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_credentials_verification_subject_dn
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionMetadataReloadSettingsMetadataUrlRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_metadata_reload_settings_metadata_url_ref_id
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionMetadataReloadSettingsMetadataUrlRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_metadata_reload_settings_metadata_url_ref_location
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnectionMetadataReloadSettings_enableAutoMetadataUpdate"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_metadata_reload_settings_enable_auto_metadata_update
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnection_active"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_active
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnection_baseUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_base_url
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnection_defaultVirtualEntityId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_default_virtual_entity_id
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnection_entityId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_entity_id
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnection_extendedProperties"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_extended_properties
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnection_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_id
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnection_licenseConnectionGroup"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_license_connection_group
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnection_loggingMode"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_logging_mode
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnection_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_name
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnection_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_type
  }
  property {
    name  = "convert_body_ConvertRequestTemplateConnection_virtualEntityIds"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_template_connection_virtual_entity_ids
  }
  property {
    name  = "convert_body_ConvertRequest_connectionType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_connection_type
  }
  property {
    name  = "convert_body_ConvertRequest_expectedEntityId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_expected_entity_id
  }
  property {
    name  = "convert_body_ConvertRequest_expectedProtocol"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_expected_protocol
  }
  property {
    name  = "convert_body_ConvertRequest_samlMetadata"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_saml_metadata
  }
  property {
    name  = "convert_body_ConvertRequest_verificationCertificate"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_convert_request_verification_certificate
  }
  property {
    name  = "convert_body_ExportRequestSigningSettingsSigningKeyPairRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_export_request_signing_settings_signing_key_pair_ref_id
  }
  property {
    name  = "convert_body_ExportRequestSigningSettingsSigningKeyPairRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_convert_body_export_request_signing_settings_signing_key_pair_ref_location
  }
  property {
    name  = "createApcMapping_body_GetApcMappings200ResponseItemsInnerAuthenticationPolicyContractRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_mapping_body_get_apc_mappings200_response_items_inner_authentication_policy_contract_ref_id
  }
  property {
    name  = "createApcMapping_body_GetApcMappings200ResponseItemsInnerAuthenticationPolicyContractRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_mapping_body_get_apc_mappings200_response_items_inner_authentication_policy_contract_ref_location
  }
  property {
    name  = "createApcMapping_body_GetApcMappings200ResponseItemsInner_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_mapping_body_get_apc_mappings200_response_items_inner_attribute_contract_fulfillment
  }
  property {
    name  = "createApcMapping_body_GetApcMappings200ResponseItemsInner_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_mapping_body_get_apc_mappings200_response_items_inner_attribute_sources
  }
  property {
    name  = "createApcMapping_body_GetApcMappings200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_mapping_body_get_apc_mappings200_response_items_inner_id
  }
  property {
    name  = "createApcMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_mapping_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "createApcMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_mapping_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "createApcMapping_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_mapping_x_bypass_external_validation
  }
  property {
    name  = "createApcToSpAdapterMapping_body_GetApcToSpAdapterMappings200ResponseItemsInner_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_to_sp_adapter_mapping_body_get_apc_to_sp_adapter_mappings200_response_items_inner_attribute_contract_fulfillment
  }
  property {
    name  = "createApcToSpAdapterMapping_body_GetApcToSpAdapterMappings200ResponseItemsInner_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_to_sp_adapter_mapping_body_get_apc_to_sp_adapter_mappings200_response_items_inner_attribute_sources
  }
  property {
    name  = "createApcToSpAdapterMapping_body_GetApcToSpAdapterMappings200ResponseItemsInner_defaultTargetResource"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_to_sp_adapter_mapping_body_get_apc_to_sp_adapter_mappings200_response_items_inner_default_target_resource
  }
  property {
    name  = "createApcToSpAdapterMapping_body_GetApcToSpAdapterMappings200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_to_sp_adapter_mapping_body_get_apc_to_sp_adapter_mappings200_response_items_inner_id
  }
  property {
    name  = "createApcToSpAdapterMapping_body_GetApcToSpAdapterMappings200ResponseItemsInner_licenseConnectionGroupAssignment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_to_sp_adapter_mapping_body_get_apc_to_sp_adapter_mappings200_response_items_inner_license_connection_group_assignment
  }
  property {
    name  = "createApcToSpAdapterMapping_body_GetApcToSpAdapterMappings200ResponseItemsInner_sourceId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_to_sp_adapter_mapping_body_get_apc_to_sp_adapter_mappings200_response_items_inner_source_id
  }
  property {
    name  = "createApcToSpAdapterMapping_body_GetApcToSpAdapterMappings200ResponseItemsInner_targetId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_to_sp_adapter_mapping_body_get_apc_to_sp_adapter_mappings200_response_items_inner_target_id
  }
  property {
    name  = "createApcToSpAdapterMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_to_sp_adapter_mapping_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "createApcToSpAdapterMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_to_sp_adapter_mapping_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "createApcToSpAdapterMapping_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_apc_to_sp_adapter_mapping_x_bypass_external_validation
  }
  property {
    name  = "createApplication_body_GetAuthenticationApiApplications200ResponseItemsInnerClientForRedirectlessModeRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_application_body_get_authentication_api_applications200_response_items_inner_client_for_redirectless_mode_ref_id
  }
  property {
    name  = "createApplication_body_GetAuthenticationApiApplications200ResponseItemsInnerClientForRedirectlessModeRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_application_body_get_authentication_api_applications200_response_items_inner_client_for_redirectless_mode_ref_location
  }
  property {
    name  = "createApplication_body_GetAuthenticationApiApplications200ResponseItemsInner_additionalAllowedOrigins"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_application_body_get_authentication_api_applications200_response_items_inner_additional_allowed_origins
  }
  property {
    name  = "createApplication_body_GetAuthenticationApiApplications200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_application_body_get_authentication_api_applications200_response_items_inner_description
  }
  property {
    name  = "createApplication_body_GetAuthenticationApiApplications200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_application_body_get_authentication_api_applications200_response_items_inner_id
  }
  property {
    name  = "createApplication_body_GetAuthenticationApiApplications200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_application_body_get_authentication_api_applications200_response_items_inner_name
  }
  property {
    name  = "createApplication_body_GetAuthenticationApiApplications200ResponseItemsInner_url"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_application_body_get_authentication_api_applications200_response_items_inner_url
  }
  property {
    name  = "createAuthenticationPolicyContract_body_GetAuthenticationPolicyContracts200ResponseItemsInner_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_authentication_policy_contract_body_get_authentication_policy_contracts200_response_items_inner_core_attributes
  }
  property {
    name  = "createAuthenticationPolicyContract_body_GetAuthenticationPolicyContracts200ResponseItemsInner_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_authentication_policy_contract_body_get_authentication_policy_contracts200_response_items_inner_extended_attributes
  }
  property {
    name  = "createAuthenticationPolicyContract_body_GetAuthenticationPolicyContracts200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_authentication_policy_contract_body_get_authentication_policy_contracts200_response_items_inner_id
  }
  property {
    name  = "createAuthenticationPolicyContract_body_GetAuthenticationPolicyContracts200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_authentication_policy_contract_body_get_authentication_policy_contracts200_response_items_inner_name
  }
  property {
    name  = "createAuthenticationSelector_body_CreateAuthenticationSelectorRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_authentication_selector_body_create_authentication_selector_request_id
  }
  property {
    name  = "createAuthenticationSelector_body_CreateAuthenticationSelectorRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_authentication_selector_body_create_authentication_selector_request_name
  }
  property {
    name  = "createAuthenticationSelector_body_GetAuthenticationSelectors200ResponseItemsInnerAllOfAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_authentication_selector_body_get_authentication_selectors200_response_items_inner_all_of_attribute_contract_extended_attributes
  }
  property {
    name  = "createAuthenticationSelector_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_authentication_selector_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "createAuthenticationSelector_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_authentication_selector_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "createAuthenticationSelector_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_authentication_selector_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "createAuthenticationSelector_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_authentication_selector_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "createAuthenticationSelector_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_authentication_selector_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "createAuthenticationSelector_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_authentication_selector_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "createClient_body_GetClient200ResponseClientAuth_clientCertIssuerDn"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_client_auth_client_cert_issuer_dn
  }
  property {
    name  = "createClient_body_GetClient200ResponseClientAuth_clientCertSubjectDn"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_client_auth_client_cert_subject_dn
  }
  property {
    name  = "createClient_body_GetClient200ResponseClientAuth_encryptedSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_client_auth_encrypted_secret
  }
  property {
    name  = "createClient_body_GetClient200ResponseClientAuth_enforceReplayPrevention"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_client_auth_enforce_replay_prevention
  }
  property {
    name  = "createClient_body_GetClient200ResponseClientAuth_secondarySecrets"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_client_auth_secondary_secrets
  }
  property {
    name  = "createClient_body_GetClient200ResponseClientAuth_secret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_client_auth_secret
  }
  property {
    name  = "createClient_body_GetClient200ResponseClientAuth_tokenEndpointAuthSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_client_auth_token_endpoint_auth_signing_algorithm
  }
  property {
    name  = "createClient_body_GetClient200ResponseClientAuth_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_client_auth_type
  }
  property {
    name  = "createClient_body_GetClient200ResponseDefaultAccessTokenManagerRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_default_access_token_manager_ref_id
  }
  property {
    name  = "createClient_body_GetClient200ResponseDefaultAccessTokenManagerRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_default_access_token_manager_ref_location
  }
  property {
    name  = "createClient_body_GetClient200ResponseJwksSettings_jwks"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_jwks_settings_jwks
  }
  property {
    name  = "createClient_body_GetClient200ResponseJwksSettings_jwksUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_jwks_settings_jwks_url
  }
  property {
    name  = "createClient_body_GetClient200ResponseOidcPolicyPolicyGroup_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_oidc_policy_policy_group_id
  }
  property {
    name  = "createClient_body_GetClient200ResponseOidcPolicyPolicyGroup_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_oidc_policy_policy_group_location
  }
  property {
    name  = "createClient_body_GetClient200ResponseOidcPolicy_grantAccessSessionRevocationApi"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_oidc_policy_grant_access_session_revocation_api
  }
  property {
    name  = "createClient_body_GetClient200ResponseOidcPolicy_grantAccessSessionSessionManagementApi"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_oidc_policy_grant_access_session_session_management_api
  }
  property {
    name  = "createClient_body_GetClient200ResponseOidcPolicy_idTokenContentEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_oidc_policy_id_token_content_encryption_algorithm
  }
  property {
    name  = "createClient_body_GetClient200ResponseOidcPolicy_idTokenEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_oidc_policy_id_token_encryption_algorithm
  }
  property {
    name  = "createClient_body_GetClient200ResponseOidcPolicy_idTokenSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_oidc_policy_id_token_signing_algorithm
  }
  property {
    name  = "createClient_body_GetClient200ResponseOidcPolicy_logoutUris"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_oidc_policy_logout_uris
  }
  property {
    name  = "createClient_body_GetClient200ResponseOidcPolicy_pairwiseIdentifierUserType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_oidc_policy_pairwise_identifier_user_type
  }
  property {
    name  = "createClient_body_GetClient200ResponseOidcPolicy_pingAccessLogoutCapable"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_oidc_policy_ping_access_logout_capable
  }
  property {
    name  = "createClient_body_GetClient200ResponseOidcPolicy_sectorIdentifierUri"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_oidc_policy_sector_identifier_uri
  }
  property {
    name  = "createClient_body_GetClient200ResponseRequestPolicyRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_request_policy_ref_id
  }
  property {
    name  = "createClient_body_GetClient200ResponseRequestPolicyRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_request_policy_ref_location
  }
  property {
    name  = "createClient_body_GetClient200ResponseTokenExchangeProcessorPolicyRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_token_exchange_processor_policy_ref_id
  }
  property {
    name  = "createClient_body_GetClient200ResponseTokenExchangeProcessorPolicyRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_token_exchange_processor_policy_ref_location
  }
  property {
    name  = "createClient_body_GetClient200Response_allowAuthenticationApiInit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_allow_authentication_api_init
  }
  property {
    name  = "createClient_body_GetClient200Response_bypassActivationCodeConfirmationOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_bypass_activation_code_confirmation_override
  }
  property {
    name  = "createClient_body_GetClient200Response_bypassApprovalPage"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_bypass_approval_page
  }
  property {
    name  = "createClient_body_GetClient200Response_cibaDeliveryMode"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_ciba_delivery_mode
  }
  property {
    name  = "createClient_body_GetClient200Response_cibaNotificationEndpoint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_ciba_notification_endpoint
  }
  property {
    name  = "createClient_body_GetClient200Response_cibaPollingInterval"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_ciba_polling_interval
  }
  property {
    name  = "createClient_body_GetClient200Response_cibaRequestObjectSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_ciba_request_object_signing_algorithm
  }
  property {
    name  = "createClient_body_GetClient200Response_cibaRequireSignedRequests"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_ciba_require_signed_requests
  }
  property {
    name  = "createClient_body_GetClient200Response_cibaUserCodeSupported"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_ciba_user_code_supported
  }
  property {
    name  = "createClient_body_GetClient200Response_clientId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_client_id
  }
  property {
    name  = "createClient_body_GetClient200Response_clientSecretChangedTime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_client_secret_changed_time
  }
  property {
    name  = "createClient_body_GetClient200Response_clientSecretRetentionPeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_client_secret_retention_period
  }
  property {
    name  = "createClient_body_GetClient200Response_clientSecretRetentionPeriodType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_client_secret_retention_period_type
  }
  property {
    name  = "createClient_body_GetClient200Response_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_description
  }
  property {
    name  = "createClient_body_GetClient200Response_deviceFlowSettingType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_device_flow_setting_type
  }
  property {
    name  = "createClient_body_GetClient200Response_devicePollingIntervalOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_device_polling_interval_override
  }
  property {
    name  = "createClient_body_GetClient200Response_enabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_enabled
  }
  property {
    name  = "createClient_body_GetClient200Response_exclusiveScopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_exclusive_scopes
  }
  property {
    name  = "createClient_body_GetClient200Response_extendedParameters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_extended_parameters
  }
  property {
    name  = "createClient_body_GetClient200Response_grantTypes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_grant_types
  }
  property {
    name  = "createClient_body_GetClient200Response_jwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_jwt_secured_authorization_response_mode_content_encryption_algorithm
  }
  property {
    name  = "createClient_body_GetClient200Response_jwtSecuredAuthorizationResponseModeEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_jwt_secured_authorization_response_mode_encryption_algorithm
  }
  property {
    name  = "createClient_body_GetClient200Response_jwtSecuredAuthorizationResponseModeSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_jwt_secured_authorization_response_mode_signing_algorithm
  }
  property {
    name  = "createClient_body_GetClient200Response_logoUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_logo_url
  }
  property {
    name  = "createClient_body_GetClient200Response_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_name
  }
  property {
    name  = "createClient_body_GetClient200Response_pendingAuthorizationTimeoutOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_pending_authorization_timeout_override
  }
  property {
    name  = "createClient_body_GetClient200Response_persistentGrantExpirationTime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_persistent_grant_expiration_time
  }
  property {
    name  = "createClient_body_GetClient200Response_persistentGrantExpirationTimeUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_persistent_grant_expiration_time_unit
  }
  property {
    name  = "createClient_body_GetClient200Response_persistentGrantExpirationType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_persistent_grant_expiration_type
  }
  property {
    name  = "createClient_body_GetClient200Response_persistentGrantIdleTimeout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_persistent_grant_idle_timeout
  }
  property {
    name  = "createClient_body_GetClient200Response_persistentGrantIdleTimeoutTimeUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_persistent_grant_idle_timeout_time_unit
  }
  property {
    name  = "createClient_body_GetClient200Response_persistentGrantIdleTimeoutType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_persistent_grant_idle_timeout_type
  }
  property {
    name  = "createClient_body_GetClient200Response_persistentGrantReuseGrantTypes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_persistent_grant_reuse_grant_types
  }
  property {
    name  = "createClient_body_GetClient200Response_persistentGrantReuseType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_persistent_grant_reuse_type
  }
  property {
    name  = "createClient_body_GetClient200Response_redirectUris"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_redirect_uris
  }
  property {
    name  = "createClient_body_GetClient200Response_refreshRolling"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_refresh_rolling
  }
  property {
    name  = "createClient_body_GetClient200Response_refreshTokenRollingGracePeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_refresh_token_rolling_grace_period
  }
  property {
    name  = "createClient_body_GetClient200Response_refreshTokenRollingGracePeriodType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_refresh_token_rolling_grace_period_type
  }
  property {
    name  = "createClient_body_GetClient200Response_refreshTokenRollingInterval"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_refresh_token_rolling_interval
  }
  property {
    name  = "createClient_body_GetClient200Response_refreshTokenRollingIntervalType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_refresh_token_rolling_interval_type
  }
  property {
    name  = "createClient_body_GetClient200Response_requestObjectSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_request_object_signing_algorithm
  }
  property {
    name  = "createClient_body_GetClient200Response_requireJwtSecuredAuthorizationResponseMode"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_require_jwt_secured_authorization_response_mode
  }
  property {
    name  = "createClient_body_GetClient200Response_requireProofKeyForCodeExchange"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_require_proof_key_for_code_exchange
  }
  property {
    name  = "createClient_body_GetClient200Response_requirePushedAuthorizationRequests"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_require_pushed_authorization_requests
  }
  property {
    name  = "createClient_body_GetClient200Response_requireSignedRequests"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_require_signed_requests
  }
  property {
    name  = "createClient_body_GetClient200Response_restrictScopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_restrict_scopes
  }
  property {
    name  = "createClient_body_GetClient200Response_restrictToDefaultAccessTokenManager"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_restrict_to_default_access_token_manager
  }
  property {
    name  = "createClient_body_GetClient200Response_restrictedResponseTypes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_restricted_response_types
  }
  property {
    name  = "createClient_body_GetClient200Response_restrictedScopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_restricted_scopes
  }
  property {
    name  = "createClient_body_GetClient200Response_tokenIntrospectionContentEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_token_introspection_content_encryption_algorithm
  }
  property {
    name  = "createClient_body_GetClient200Response_tokenIntrospectionEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_token_introspection_encryption_algorithm
  }
  property {
    name  = "createClient_body_GetClient200Response_tokenIntrospectionSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_token_introspection_signing_algorithm
  }
  property {
    name  = "createClient_body_GetClient200Response_userAuthorizationUrlOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_user_authorization_url_override
  }
  property {
    name  = "createClient_body_GetClient200Response_validateUsingAllEligibleAtms"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_client_body_get_client200_response_validate_using_all_eligible_atms
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseAttributeQueryPolicy_encryptNameId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_attribute_query_policy_encrypt_name_id
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseAttributeQueryPolicy_maskAttributeValues"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_attribute_query_policy_mask_attribute_values
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseAttributeQueryPolicy_requireEncryptedAssertion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_attribute_query_policy_require_encrypted_assertion
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseAttributeQueryPolicy_requireSignedAssertion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_attribute_query_policy_require_signed_assertion
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseAttributeQueryPolicy_requireSignedResponse"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_attribute_query_policy_require_signed_response
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseAttributeQueryPolicy_signAttributeQuery"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_attribute_query_policy_sign_attribute_query
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseAttributeQuery_nameMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_attribute_query_name_mappings
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseAttributeQuery_url"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_attribute_query_url
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoArtifact_lifetime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_artifact_lifetime
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoArtifact_resolverLocations"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_artifact_resolver_locations
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoArtifact_sourceId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_artifact_source_id
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_attribute_contract_core_attributes
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_attribute_contract_extended_attributes
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoDecryptionPolicy_assertionEncrypted"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_decryption_policy_assertion_encrypted
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoDecryptionPolicy_attributesEncrypted"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_decryption_policy_attributes_encrypted
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoDecryptionPolicy_sloEncryptSubjectNameID"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_decryption_policy_slo_encrypt_subject_name_id
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoDecryptionPolicy_sloSubjectNameIDEncrypted"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_decryption_policy_slo_subject_name_idencrypted
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoDecryptionPolicy_subjectNameIdEncrypted"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_decryption_policy_subject_name_id_encrypted
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoJitProvisioningUserAttributes_attributeContract"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_jit_provisioning_user_attributes_attribute_contract
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoJitProvisioningUserAttributes_doAttributeQuery"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_jit_provisioning_user_attributes_do_attribute_query
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoJitProvisioningUserRepository_jitRepositoryAttributeMapping"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_jit_provisioning_user_repository_jit_repository_attribute_mapping
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoJitProvisioningUserRepository_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_jit_provisioning_user_repository_type
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoJitProvisioning_errorHandling"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_jit_provisioning_error_handling
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoJitProvisioning_eventTrigger"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_jit_provisioning_event_trigger
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoOauthAuthenticationPolicyContractRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_oauth_authentication_policy_contract_ref_id
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoOauthAuthenticationPolicyContractRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_oauth_authentication_policy_contract_ref_location
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_authenticationScheme"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_authentication_scheme
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_authenticationSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_authentication_signing_algorithm
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_authorizationEndpoint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_authorization_endpoint
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_enablePKCE"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_enable_pkce
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_jwksURL"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_jwks_url
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_loginType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_login_type
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_requestParameters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_request_parameters
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_requestSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_request_signing_algorithm
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_scopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_scopes
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_tokenEndpoint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_token_endpoint
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_userInfoEndpoint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_user_info_endpoint
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoSsoOAuthMapping_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_sso_oauth_mapping_attribute_contract_fulfillment
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSsoSsoOAuthMapping_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_sso_oauth_mapping_attribute_sources
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_adapterMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_adapter_mappings
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_alwaysSignArtifactResponse"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_always_sign_artifact_response
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_assertionsSigned"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_assertions_signed
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_authenticationPolicyContractMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_authentication_policy_contract_mappings
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_authnContextMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_authn_context_mappings
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_defaultTargetUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_default_target_url
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_enabledProfiles"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_enabled_profiles
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_idpIdentityMapping"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_idp_identity_mapping
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_incomingBindings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_incoming_bindings
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_messageCustomizations"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_message_customizations
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_protocol"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_protocol
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_signAuthnRequests"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_sign_authn_requests
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_sloServiceEndpoints"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_slo_service_endpoints
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_ssoServiceEndpoints"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_sso_service_endpoints
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpBrowserSso_urlWhitelistEntries"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_browser_sso_url_whitelist_entries
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpOAuthGrantAttributeMappingIdpOAuthAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_oauth_grant_attribute_mapping_idp_oauth_attribute_contract_core_attributes
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpOAuthGrantAttributeMappingIdpOAuthAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_oauth_grant_attribute_mapping_idp_oauth_attribute_contract_extended_attributes
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseIdpOAuthGrantAttributeMapping_accessTokenManagerMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_idp_oauth_grant_attribute_mapping_access_token_manager_mappings
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseInboundProvisioningCustomSchema_attributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_inbound_provisioning_custom_schema_attributes
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseInboundProvisioningCustomSchema_namespace"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_inbound_provisioning_custom_schema_namespace
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseInboundProvisioningGroupsReadGroups_attributeFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_inbound_provisioning_groups_read_groups_attribute_fulfillment
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseInboundProvisioningGroupsReadGroups_attributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_inbound_provisioning_groups_read_groups_attributes
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseInboundProvisioningGroupsWriteGroups_attributeFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_inbound_provisioning_groups_write_groups_attribute_fulfillment
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseInboundProvisioningUserRepository_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_inbound_provisioning_user_repository_type
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseInboundProvisioningUsersReadUsersAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_inbound_provisioning_users_read_users_attribute_contract_core_attributes
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseInboundProvisioningUsersReadUsersAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_inbound_provisioning_users_read_users_attribute_contract_extended_attributes
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseInboundProvisioningUsersReadUsers_attributeFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_inbound_provisioning_users_read_users_attribute_fulfillment
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseInboundProvisioningUsersReadUsers_attributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_inbound_provisioning_users_read_users_attributes
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseInboundProvisioningUsersWriteUsers_attributeFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_inbound_provisioning_users_write_users_attribute_fulfillment
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseInboundProvisioning_actionOnDelete"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_inbound_provisioning_action_on_delete
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseInboundProvisioning_groupSupport"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_inbound_provisioning_group_support
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseOidcClientCredentials_clientId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_oidc_client_credentials_client_id
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseOidcClientCredentials_clientSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_oidc_client_credentials_client_secret
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseOidcClientCredentials_encryptedSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_oidc_client_credentials_encrypted_secret
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseWsTrustAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_ws_trust_attribute_contract_core_attributes
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseWsTrustAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_ws_trust_attribute_contract_extended_attributes
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseWsTrust_generateLocalToken"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_ws_trust_generate_local_token
  }
  property {
    name  = "createConnection1_body_GetConnection1200ResponseWsTrust_tokenGeneratorMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_connection1200_response_ws_trust_token_generator_mappings
  }
  property {
    name  = "createConnection1_body_GetMappings200ResponseItemsInnerAttributeSourcesInnerDataStoreRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_mappings200_response_items_inner_attribute_sources_inner_data_store_ref_id
  }
  property {
    name  = "createConnection1_body_GetMappings200ResponseItemsInnerAttributeSourcesInnerDataStoreRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_mappings200_response_items_inner_attribute_sources_inner_data_store_ref_location
  }
  property {
    name  = "createConnection1_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "createConnection1_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "createConnection1_body_UpdateConnection1Request_errorPageMsgId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_body_update_connection1_request_error_page_msg_id
  }
  property {
    name  = "createConnection1_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection1_x_bypass_external_validation
  }
  property {
    name  = "createConnection2_body_GetConnection1200ResponseIdpBrowserSsoArtifact_lifetime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection1200_response_idp_browser_sso_artifact_lifetime
  }
  property {
    name  = "createConnection2_body_GetConnection1200ResponseIdpBrowserSsoArtifact_resolverLocations"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection1200_response_idp_browser_sso_artifact_resolver_locations
  }
  property {
    name  = "createConnection2_body_GetConnection1200ResponseIdpBrowserSsoArtifact_sourceId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection1200_response_idp_browser_sso_artifact_source_id
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseAttributeQueryPolicy_encryptAssertion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_attribute_query_policy_encrypt_assertion
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseAttributeQueryPolicy_requireEncryptedNameId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_attribute_query_policy_require_encrypted_name_id
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseAttributeQueryPolicy_requireSignedAttributeQuery"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_attribute_query_policy_require_signed_attribute_query
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseAttributeQueryPolicy_signAssertion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_attribute_query_policy_sign_assertion
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseAttributeQueryPolicy_signResponse"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_attribute_query_policy_sign_response
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseAttributeQuery_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_attribute_query_attribute_contract_fulfillment
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseAttributeQuery_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_attribute_query_attribute_sources
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseAttributeQuery_attributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_attribute_query_attributes
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseOutboundProvisionCustomSchema_attributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_outbound_provision_custom_schema_attributes
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseOutboundProvisionCustomSchema_namespace"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_outbound_provision_custom_schema_namespace
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseOutboundProvision_channels"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_outbound_provision_channels
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseOutboundProvision_targetSettings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_outbound_provision_target_settings
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseOutboundProvision_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_outbound_provision_type
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSsoAssertionLifetime_minutesAfter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_assertion_lifetime_minutes_after
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSsoAssertionLifetime_minutesBefore"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_assertion_lifetime_minutes_before
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSsoAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_attribute_contract_core_attributes
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSsoAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_attribute_contract_extended_attributes
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSsoEncryptionPolicy_encryptAssertion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_encryption_policy_encrypt_assertion
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSsoEncryptionPolicy_encryptSloSubjectNameId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_encryption_policy_encrypt_slo_subject_name_id
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSsoEncryptionPolicy_encryptedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_encryption_policy_encrypted_attributes
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSsoEncryptionPolicy_sloSubjectNameIDEncrypted"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_encryption_policy_slo_subject_name_idencrypted
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_adapterMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_adapter_mappings
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_alwaysSignArtifactResponse"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_always_sign_artifact_response
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_authenticationPolicyContractAssertionMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_authentication_policy_contract_assertion_mappings
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_defaultTargetUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_default_target_url
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_enabledProfiles"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_enabled_profiles
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_incomingBindings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_incoming_bindings
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_messageCustomizations"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_message_customizations
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_protocol"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_protocol
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_requireSignedAuthnRequests"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_require_signed_authn_requests
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_signAssertions"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_sign_assertions
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_signResponseAsRequired"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_sign_response_as_required
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_sloServiceEndpoints"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_slo_service_endpoints
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_spSamlIdentityMapping"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_sp_saml_identity_mapping
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_spWsFedIdentityMapping"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_sp_ws_fed_identity_mapping
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_ssoServiceEndpoints"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_sso_service_endpoints
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_urlWhitelistEntries"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_url_whitelist_entries
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_wsFedTokenType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_ws_fed_token_type
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseSpBrowserSso_wsTrustVersion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_sp_browser_sso_ws_trust_version
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrustAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_attribute_contract_core_attributes
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrustAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_attribute_contract_extended_attributes
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrustRequestContractRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_request_contract_ref_id
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrustRequestContractRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_request_contract_ref_location
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrust_abortIfNotFulfilledFromRequest"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_abort_if_not_fulfilled_from_request
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrust_defaultTokenType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_default_token_type
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrust_encryptSaml2Assertion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_encrypt_saml2_assertion
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrust_generateKey"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_generate_key
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrust_messageCustomizations"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_message_customizations
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrust_minutesAfter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_minutes_after
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrust_minutesBefore"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_minutes_before
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrust_oAuthAssertionProfiles"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_o_auth_assertion_profiles
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrust_partnerServiceIds"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_partner_service_ids
  }
  property {
    name  = "createConnection2_body_GetConnection2200ResponseWsTrust_tokenProcessorMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_connection2200_response_ws_trust_token_processor_mappings
  }
  property {
    name  = "createConnection2_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "createConnection2_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "createConnection2_body_UpdateConnection2Request_applicationIconUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_update_connection2_request_application_icon_url
  }
  property {
    name  = "createConnection2_body_UpdateConnection2Request_applicationName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_update_connection2_request_application_name
  }
  property {
    name  = "createConnection2_body_UpdateConnection2Request_connectionTargetType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_body_update_connection2_request_connection_target_type
  }
  property {
    name  = "createConnection2_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_connection2_x_bypass_external_validation
  }
  property {
    name  = "createDataStore_body_CreateDataStoreRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_data_store_body_create_data_store_request_id
  }
  property {
    name  = "createDataStore_body_CreateDataStoreRequest_maskAttributeValues"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_data_store_body_create_data_store_request_mask_attribute_values
  }
  property {
    name  = "createDataStore_body_CreateDataStoreRequest_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_data_store_body_create_data_store_request_type
  }
  property {
    name  = "createDataStore_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_data_store_x_bypass_external_validation
  }
  property {
    name  = "createDynamicClientRegistrationPolicy_body_CreateDynamicClientRegistrationPolicyRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_dynamic_client_registration_policy_body_create_dynamic_client_registration_policy_request_id
  }
  property {
    name  = "createDynamicClientRegistrationPolicy_body_CreateDynamicClientRegistrationPolicyRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_dynamic_client_registration_policy_body_create_dynamic_client_registration_policy_request_name
  }
  property {
    name  = "createDynamicClientRegistrationPolicy_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_dynamic_client_registration_policy_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "createDynamicClientRegistrationPolicy_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_dynamic_client_registration_policy_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "createDynamicClientRegistrationPolicy_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_dynamic_client_registration_policy_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "createDynamicClientRegistrationPolicy_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_dynamic_client_registration_policy_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "createDynamicClientRegistrationPolicy_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_dynamic_client_registration_policy_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "createDynamicClientRegistrationPolicy_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_dynamic_client_registration_policy_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "createFragment_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerRootNodeAction_context"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_fragment_body_get_default_authentication_policy200_response_authn_selection_trees_inner_root_node_action_context
  }
  property {
    name  = "createFragment_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerRootNodeAction_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_fragment_body_get_default_authentication_policy200_response_authn_selection_trees_inner_root_node_action_type
  }
  property {
    name  = "createFragment_body_GetFragment200ResponseInputs_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_fragment_body_get_fragment200_response_inputs_id
  }
  property {
    name  = "createFragment_body_GetFragment200ResponseInputs_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_fragment_body_get_fragment200_response_inputs_location
  }
  property {
    name  = "createFragment_body_GetFragment200ResponseOutputs_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_fragment_body_get_fragment200_response_outputs_id
  }
  property {
    name  = "createFragment_body_GetFragment200ResponseOutputs_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_fragment_body_get_fragment200_response_outputs_location
  }
  property {
    name  = "createFragment_body_GetFragment200ResponseRootNode_children"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_fragment_body_get_fragment200_response_root_node_children
  }
  property {
    name  = "createFragment_body_GetFragment200Response_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_fragment_body_get_fragment200_response_description
  }
  property {
    name  = "createFragment_body_GetFragment200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_fragment_body_get_fragment200_response_id
  }
  property {
    name  = "createFragment_body_GetFragment200Response_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_fragment_body_get_fragment200_response_name
  }
  property {
    name  = "createFragment_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_fragment_x_bypass_external_validation
  }
  property {
    name  = "createGroup_body_GetGroup200Response_generatorMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_group_body_get_group200_response_generator_mappings
  }
  property {
    name  = "createGroup_body_GetGroup200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_group_body_get_group200_response_id
  }
  property {
    name  = "createGroup_body_GetGroup200Response_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_group_body_get_group200_response_name
  }
  property {
    name  = "createGroup_body_GetGroup200Response_resourceUris"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_group_body_get_group200_response_resource_uris
  }
  property {
    name  = "createGroup_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_group_bypass_external_validation
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerApcId_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_apc_id_id
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerApcId_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_apc_id_location
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerAuthSourceUpdatePolicy_retainAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_auth_source_update_policy_retain_attributes
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerAuthSourceUpdatePolicy_storeAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_auth_source_update_policy_store_attributes
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerAuthSourceUpdatePolicy_updateAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_auth_source_update_policy_update_attributes
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerAuthSourceUpdatePolicy_updateInterval"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_auth_source_update_policy_update_interval
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerDataStoreConfig_dataStoreMapping"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_data_store_config_data_store_mapping
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerDataStoreConfig_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_data_store_config_type
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfigNotificationPublisherRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_notification_publisher_ref_id
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfigNotificationPublisherRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_notification_publisher_ref_location
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_allowedOtpCharacterSet"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_allowed_otp_character_set
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_emailVerificationEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_email_verification_enabled
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_emailVerificationErrorTemplateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_email_verification_error_template_name
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_emailVerificationOtpTemplateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_email_verification_otp_template_name
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_emailVerificationSentTemplateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_email_verification_sent_template_name
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_emailVerificationSuccessTemplateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_email_verification_success_template_name
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_emailVerificationType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_email_verification_type
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_fieldForEmailToVerify"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_field_for_email_to_verify
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_fieldStoringVerificationStatus"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_field_storing_verification_status
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_otlTimeToLive"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_otl_time_to_live
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_otpLength"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_otp_length
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_otpRetryAttempts"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_otp_retry_attempts
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_otpTimeToLive"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_otp_time_to_live
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_requireVerifiedEmail"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_require_verified_email
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_requireVerifiedEmailTemplateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_require_verified_email_template_name
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_verifyEmailTemplateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_verify_email_template_name
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerFieldConfig_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_field_config_fields
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerFieldConfig_stripSpaceFromUniqueField"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_field_config_strip_space_from_unique_field
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerProfileConfig_deleteIdentityEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_profile_config_delete_identity_enabled
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerProfileConfig_templateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_profile_config_template_name
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfigRegistrationWorkflow_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_registration_workflow_id
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfigRegistrationWorkflow_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_registration_workflow_location
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfig_captchaEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_captcha_enabled
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfig_createAuthnSessionAfterRegistration"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_create_authn_session_after_registration
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfig_executeWorkflow"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_execute_workflow
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfig_templateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_template_name
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfig_thisIsMyDeviceEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_this_is_my_device_enabled
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfig_usernameField"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_username_field
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInner_authSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_auth_sources
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_id
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_name
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInner_profileEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_profile_enabled
  }
  property {
    name  = "createIdentityProfile_body_GetIdentityProfiles200ResponseItemsInner_registrationEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_identity_profiles200_response_items_inner_registration_enabled
  }
  property {
    name  = "createIdentityProfile_body_GetMappings200ResponseItemsInnerAttributeSourcesInnerDataStoreRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_mappings200_response_items_inner_attribute_sources_inner_data_store_ref_id
  }
  property {
    name  = "createIdentityProfile_body_GetMappings200ResponseItemsInnerAttributeSourcesInnerDataStoreRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_body_get_mappings200_response_items_inner_attribute_sources_inner_data_store_ref_location
  }
  property {
    name  = "createIdentityProfile_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_profile_x_bypass_external_validation
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInnerAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_attribute_contract_core_attributes
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInnerAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_attribute_contract_extended_attributes
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInnerAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_attribute_contract_inherited
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInnerGroupAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_group_attribute_contract_core_attributes
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInnerGroupAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_group_attribute_contract_extended_attributes
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInnerGroupAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_group_attribute_contract_inherited
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_id
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_name
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "createIdentityStoreProvisioner_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_identity_store_provisioner_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "createIdpAdapterMapping_body_GetIdpAdapterMappings200ResponseItemsInnerIdpAdapterRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_mapping_body_get_idp_adapter_mappings200_response_items_inner_idp_adapter_ref_id
  }
  property {
    name  = "createIdpAdapterMapping_body_GetIdpAdapterMappings200ResponseItemsInnerIdpAdapterRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_mapping_body_get_idp_adapter_mappings200_response_items_inner_idp_adapter_ref_location
  }
  property {
    name  = "createIdpAdapterMapping_body_GetIdpAdapterMappings200ResponseItemsInner_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_mapping_body_get_idp_adapter_mappings200_response_items_inner_attribute_contract_fulfillment
  }
  property {
    name  = "createIdpAdapterMapping_body_GetIdpAdapterMappings200ResponseItemsInner_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_mapping_body_get_idp_adapter_mappings200_response_items_inner_attribute_sources
  }
  property {
    name  = "createIdpAdapterMapping_body_GetIdpAdapterMappings200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_mapping_body_get_idp_adapter_mappings200_response_items_inner_id
  }
  property {
    name  = "createIdpAdapterMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_mapping_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "createIdpAdapterMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_mapping_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "createIdpAdapterMapping_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_mapping_x_bypass_external_validation
  }
  property {
    name  = "createIdpAdapter_body_CreateIdpAdapterRequest_authnCtxClassRef"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_create_idp_adapter_request_authn_ctx_class_ref
  }
  property {
    name  = "createIdpAdapter_body_CreateIdpAdapterRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_create_idp_adapter_request_id
  }
  property {
    name  = "createIdpAdapter_body_CreateIdpAdapterRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_create_idp_adapter_request_name
  }
  property {
    name  = "createIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_contract_core_attributes
  }
  property {
    name  = "createIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_contract_extended_attributes
  }
  property {
    name  = "createIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_contract_inherited
  }
  property {
    name  = "createIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeContract_maskOgnlValues"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_contract_mask_ognl_values
  }
  property {
    name  = "createIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeContract_uniqueUserKeyAttribute"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_contract_unique_user_key_attribute
  }
  property {
    name  = "createIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeMapping_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_mapping_attribute_contract_fulfillment
  }
  property {
    name  = "createIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeMapping_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_mapping_attribute_sources
  }
  property {
    name  = "createIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeMapping_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_mapping_inherited
  }
  property {
    name  = "createIdpAdapter_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "createIdpAdapter_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "createIdpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "createIdpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "createIdpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "createIdpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "createIdpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "createIdpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "createIdpAdapter_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_adapter_x_bypass_external_validation
  }
  property {
    name  = "createIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_applicationIconUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_application_icon_url
  }
  property {
    name  = "createIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_applicationName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_application_name
  }
  property {
    name  = "createIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_attribute_contract_fulfillment
  }
  property {
    name  = "createIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_attribute_sources
  }
  property {
    name  = "createIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_defaultTargetResource"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_default_target_resource
  }
  property {
    name  = "createIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_id
  }
  property {
    name  = "createIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_licenseConnectionGroupAssignment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_license_connection_group_assignment
  }
  property {
    name  = "createIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_sourceId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_source_id
  }
  property {
    name  = "createIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_targetId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_target_id
  }
  property {
    name  = "createIdpToSpAdapterMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_to_sp_adapter_mapping_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "createIdpToSpAdapterMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_to_sp_adapter_mapping_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "createIdpToSpAdapterMapping_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_idp_to_sp_adapter_mapping_x_bypass_external_validation
  }
  property {
    name  = "createKerberosRealm_body_GetKerberosRealms200ResponseItemsInnerLdapGatewayDataStoreRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_kerberos_realm_body_get_kerberos_realms200_response_items_inner_ldap_gateway_data_store_ref_id
  }
  property {
    name  = "createKerberosRealm_body_GetKerberosRealms200ResponseItemsInnerLdapGatewayDataStoreRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_kerberos_realm_body_get_kerberos_realms200_response_items_inner_ldap_gateway_data_store_ref_location
  }
  property {
    name  = "createKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_connectionType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_kerberos_realm_body_get_kerberos_realms200_response_items_inner_connection_type
  }
  property {
    name  = "createKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_kerberos_realm_body_get_kerberos_realms200_response_items_inner_id
  }
  property {
    name  = "createKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_kerberosEncryptedPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_kerberos_realm_body_get_kerberos_realms200_response_items_inner_kerberos_encrypted_password
  }
  property {
    name  = "createKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_kerberosPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_kerberos_realm_body_get_kerberos_realms200_response_items_inner_kerberos_password
  }
  property {
    name  = "createKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_kerberosRealmName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_kerberos_realm_body_get_kerberos_realms200_response_items_inner_kerberos_realm_name
  }
  property {
    name  = "createKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_kerberosUsername"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_kerberos_realm_body_get_kerberos_realms200_response_items_inner_kerberos_username
  }
  property {
    name  = "createKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_keyDistributionCenters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_kerberos_realm_body_get_kerberos_realms200_response_items_inner_key_distribution_centers
  }
  property {
    name  = "createKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_keySets"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_kerberos_realm_body_get_kerberos_realms200_response_items_inner_key_sets
  }
  property {
    name  = "createKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_retainPreviousKeysOnPasswordChange"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_kerberos_realm_body_get_kerberos_realms200_response_items_inner_retain_previous_keys_on_password_change
  }
  property {
    name  = "createKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_suppressDomainNameConcatenation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_kerberos_realm_body_get_kerberos_realms200_response_items_inner_suppress_domain_name_concatenation
  }
  property {
    name  = "createKerberosRealm_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_kerberos_realm_x_bypass_external_validation
  }
  property {
    name  = "createKeyPair1_body_CreateKeyPair1Request_city"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair1_body_create_key_pair1_request_city
  }
  property {
    name  = "createKeyPair1_body_CreateKeyPair1Request_commonName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair1_body_create_key_pair1_request_common_name
  }
  property {
    name  = "createKeyPair1_body_CreateKeyPair1Request_country"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair1_body_create_key_pair1_request_country
  }
  property {
    name  = "createKeyPair1_body_CreateKeyPair1Request_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair1_body_create_key_pair1_request_crypto_provider
  }
  property {
    name  = "createKeyPair1_body_CreateKeyPair1Request_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair1_body_create_key_pair1_request_id
  }
  property {
    name  = "createKeyPair1_body_CreateKeyPair1Request_keyAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair1_body_create_key_pair1_request_key_algorithm
  }
  property {
    name  = "createKeyPair1_body_CreateKeyPair1Request_keySize"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair1_body_create_key_pair1_request_key_size
  }
  property {
    name  = "createKeyPair1_body_CreateKeyPair1Request_organization"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair1_body_create_key_pair1_request_organization
  }
  property {
    name  = "createKeyPair1_body_CreateKeyPair1Request_organizationUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair1_body_create_key_pair1_request_organization_unit
  }
  property {
    name  = "createKeyPair1_body_CreateKeyPair1Request_signatureAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair1_body_create_key_pair1_request_signature_algorithm
  }
  property {
    name  = "createKeyPair1_body_CreateKeyPair1Request_state"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair1_body_create_key_pair1_request_state
  }
  property {
    name  = "createKeyPair1_body_CreateKeyPair1Request_subjectAlternativeNames"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair1_body_create_key_pair1_request_subject_alternative_names
  }
  property {
    name  = "createKeyPair1_body_CreateKeyPair1Request_validDays"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair1_body_create_key_pair1_request_valid_days
  }
  property {
    name  = "createKeyPair2_body_CreateKeyPair1Request_city"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair2_body_create_key_pair1_request_city
  }
  property {
    name  = "createKeyPair2_body_CreateKeyPair1Request_commonName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair2_body_create_key_pair1_request_common_name
  }
  property {
    name  = "createKeyPair2_body_CreateKeyPair1Request_country"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair2_body_create_key_pair1_request_country
  }
  property {
    name  = "createKeyPair2_body_CreateKeyPair1Request_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair2_body_create_key_pair1_request_crypto_provider
  }
  property {
    name  = "createKeyPair2_body_CreateKeyPair1Request_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair2_body_create_key_pair1_request_id
  }
  property {
    name  = "createKeyPair2_body_CreateKeyPair1Request_keyAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair2_body_create_key_pair1_request_key_algorithm
  }
  property {
    name  = "createKeyPair2_body_CreateKeyPair1Request_keySize"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair2_body_create_key_pair1_request_key_size
  }
  property {
    name  = "createKeyPair2_body_CreateKeyPair1Request_organization"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair2_body_create_key_pair1_request_organization
  }
  property {
    name  = "createKeyPair2_body_CreateKeyPair1Request_organizationUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair2_body_create_key_pair1_request_organization_unit
  }
  property {
    name  = "createKeyPair2_body_CreateKeyPair1Request_signatureAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair2_body_create_key_pair1_request_signature_algorithm
  }
  property {
    name  = "createKeyPair2_body_CreateKeyPair1Request_state"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair2_body_create_key_pair1_request_state
  }
  property {
    name  = "createKeyPair2_body_CreateKeyPair1Request_subjectAlternativeNames"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair2_body_create_key_pair1_request_subject_alternative_names
  }
  property {
    name  = "createKeyPair2_body_CreateKeyPair1Request_validDays"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair2_body_create_key_pair1_request_valid_days
  }
  property {
    name  = "createKeyPair3_body_CreateKeyPair1Request_city"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair3_body_create_key_pair1_request_city
  }
  property {
    name  = "createKeyPair3_body_CreateKeyPair1Request_commonName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair3_body_create_key_pair1_request_common_name
  }
  property {
    name  = "createKeyPair3_body_CreateKeyPair1Request_country"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair3_body_create_key_pair1_request_country
  }
  property {
    name  = "createKeyPair3_body_CreateKeyPair1Request_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair3_body_create_key_pair1_request_crypto_provider
  }
  property {
    name  = "createKeyPair3_body_CreateKeyPair1Request_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair3_body_create_key_pair1_request_id
  }
  property {
    name  = "createKeyPair3_body_CreateKeyPair1Request_keyAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair3_body_create_key_pair1_request_key_algorithm
  }
  property {
    name  = "createKeyPair3_body_CreateKeyPair1Request_keySize"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair3_body_create_key_pair1_request_key_size
  }
  property {
    name  = "createKeyPair3_body_CreateKeyPair1Request_organization"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair3_body_create_key_pair1_request_organization
  }
  property {
    name  = "createKeyPair3_body_CreateKeyPair1Request_organizationUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair3_body_create_key_pair1_request_organization_unit
  }
  property {
    name  = "createKeyPair3_body_CreateKeyPair1Request_signatureAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair3_body_create_key_pair1_request_signature_algorithm
  }
  property {
    name  = "createKeyPair3_body_CreateKeyPair1Request_state"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair3_body_create_key_pair1_request_state
  }
  property {
    name  = "createKeyPair3_body_CreateKeyPair1Request_subjectAlternativeNames"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair3_body_create_key_pair1_request_subject_alternative_names
  }
  property {
    name  = "createKeyPair3_body_CreateKeyPair1Request_validDays"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_pair3_body_create_key_pair1_request_valid_days
  }
  property {
    name  = "createKeySet_body_GetKeySet200ResponseSigningKeys_p256PublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_key_set200_response_signing_keys_p256_publish_x5c_parameter
  }
  property {
    name  = "createKeySet_body_GetKeySet200ResponseSigningKeys_p384PublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_key_set200_response_signing_keys_p384_publish_x5c_parameter
  }
  property {
    name  = "createKeySet_body_GetKeySet200ResponseSigningKeys_p521PublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_key_set200_response_signing_keys_p521_publish_x5c_parameter
  }
  property {
    name  = "createKeySet_body_GetKeySet200ResponseSigningKeys_rsaPublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_key_set200_response_signing_keys_rsa_publish_x5c_parameter
  }
  property {
    name  = "createKeySet_body_GetKeySet200Response_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_key_set200_response_description
  }
  property {
    name  = "createKeySet_body_GetKeySet200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_key_set200_response_id
  }
  property {
    name  = "createKeySet_body_GetKeySet200Response_issuers"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_key_set200_response_issuers
  }
  property {
    name  = "createKeySet_body_GetKeySet200Response_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_key_set200_response_name
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseP256ActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_p256_active_cert_ref_id
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseP256ActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_p256_active_cert_ref_location
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseP256PreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_p256_previous_cert_ref_id
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseP256PreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_p256_previous_cert_ref_location
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseP384ActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_p384_active_cert_ref_id
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseP384ActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_p384_active_cert_ref_location
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseP384PreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_p384_previous_cert_ref_id
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseP384PreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_p384_previous_cert_ref_location
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseP521ActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_p521_active_cert_ref_id
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseP521ActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_p521_active_cert_ref_location
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseP521PreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_p521_previous_cert_ref_id
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseP521PreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_p521_previous_cert_ref_location
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseRsaActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_rsa_active_cert_ref_id
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseRsaActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_rsa_active_cert_ref_location
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseRsaPreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_rsa_previous_cert_ref_id
  }
  property {
    name  = "createKeySet_body_GetOauthOidcKeysSettings200ResponseRsaPreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_key_set_body_get_oauth_oidc_keys_settings200_response_rsa_previous_cert_ref_location
  }
  property {
    name  = "createMapping_body_GetMappings200ResponseItemsInnerAccessTokenManagerRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_mapping_body_get_mappings200_response_items_inner_access_token_manager_ref_id
  }
  property {
    name  = "createMapping_body_GetMappings200ResponseItemsInnerAccessTokenManagerRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_mapping_body_get_mappings200_response_items_inner_access_token_manager_ref_location
  }
  property {
    name  = "createMapping_body_GetMappings200ResponseItemsInnerContextContextRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_mapping_body_get_mappings200_response_items_inner_context_context_ref_id
  }
  property {
    name  = "createMapping_body_GetMappings200ResponseItemsInnerContextContextRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_mapping_body_get_mappings200_response_items_inner_context_context_ref_location
  }
  property {
    name  = "createMapping_body_GetMappings200ResponseItemsInnerContext_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_mapping_body_get_mappings200_response_items_inner_context_type
  }
  property {
    name  = "createMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_mapping_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "createMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_mapping_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "createMapping_body_GetMappings200ResponseItemsInner_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_mapping_body_get_mappings200_response_items_inner_attribute_contract_fulfillment
  }
  property {
    name  = "createMapping_body_GetMappings200ResponseItemsInner_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_mapping_body_get_mappings200_response_items_inner_attribute_sources
  }
  property {
    name  = "createMapping_body_GetMappings200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_mapping_body_get_mappings200_response_items_inner_id
  }
  property {
    name  = "createMapping_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_mapping_x_bypass_external_validation
  }
  property {
    name  = "createNotificationPublisher_body_CreateNotificationPublisherRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_notification_publisher_body_create_notification_publisher_request_id
  }
  property {
    name  = "createNotificationPublisher_body_CreateNotificationPublisherRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_notification_publisher_body_create_notification_publisher_request_name
  }
  property {
    name  = "createNotificationPublisher_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_notification_publisher_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "createNotificationPublisher_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_notification_publisher_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "createNotificationPublisher_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_notification_publisher_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "createNotificationPublisher_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_notification_publisher_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "createNotificationPublisher_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_notification_publisher_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "createNotificationPublisher_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_notification_publisher_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "createOOBAuthenticator_body_CreateOOBAuthenticatorRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_oobauthenticator_body_create_oobauthenticator_request_id
  }
  property {
    name  = "createOOBAuthenticator_body_CreateOOBAuthenticatorRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_oobauthenticator_body_create_oobauthenticator_request_name
  }
  property {
    name  = "createOOBAuthenticator_body_GetOOBAuthenticators200ResponseItemsInnerAllOfAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_oobauthenticator_body_get_oobauthenticators200_response_items_inner_all_of_attribute_contract_core_attributes
  }
  property {
    name  = "createOOBAuthenticator_body_GetOOBAuthenticators200ResponseItemsInnerAllOfAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_oobauthenticator_body_get_oobauthenticators200_response_items_inner_all_of_attribute_contract_extended_attributes
  }
  property {
    name  = "createOOBAuthenticator_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_oobauthenticator_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "createOOBAuthenticator_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_oobauthenticator_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "createOOBAuthenticator_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_oobauthenticator_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "createOOBAuthenticator_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_oobauthenticator_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "createOOBAuthenticator_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_oobauthenticator_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "createOOBAuthenticator_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_oobauthenticator_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "createPasswordCredentialValidator_body_CreatePasswordCredentialValidatorRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_password_credential_validator_body_create_password_credential_validator_request_id
  }
  property {
    name  = "createPasswordCredentialValidator_body_CreatePasswordCredentialValidatorRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_password_credential_validator_body_create_password_credential_validator_request_name
  }
  property {
    name  = "createPasswordCredentialValidator_body_GetPasswordCredentialValidators200ResponseItemsInnerAllOfAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_password_credential_validator_body_get_password_credential_validators200_response_items_inner_all_of_attribute_contract_core_attributes
  }
  property {
    name  = "createPasswordCredentialValidator_body_GetPasswordCredentialValidators200ResponseItemsInnerAllOfAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_password_credential_validator_body_get_password_credential_validators200_response_items_inner_all_of_attribute_contract_extended_attributes
  }
  property {
    name  = "createPasswordCredentialValidator_body_GetPasswordCredentialValidators200ResponseItemsInnerAllOfAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_password_credential_validator_body_get_password_credential_validators200_response_items_inner_all_of_attribute_contract_inherited
  }
  property {
    name  = "createPasswordCredentialValidator_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_password_credential_validator_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "createPasswordCredentialValidator_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_password_credential_validator_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "createPasswordCredentialValidator_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_password_credential_validator_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "createPasswordCredentialValidator_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_password_credential_validator_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "createPasswordCredentialValidator_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_password_credential_validator_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "createPasswordCredentialValidator_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_password_credential_validator_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_active"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_active
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_creationDate"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_creation_date
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_credential"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_credential
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_credentialId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_credential_id
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_description
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_encryptedCredential"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_encrypted_credential
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_environmentId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_environment_id
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_id
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_name
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_organizationName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_organization_name
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_pingOneAuthenticationApiEndpoint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_ping_one_authentication_api_endpoint
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_pingOneConnectionId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_ping_one_connection_id
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_pingOneManagementApiEndpoint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_ping_one_management_api_endpoint
  }
  property {
    name  = "createPingOneConnection_body_GetPingOneConnections200ResponseItemsInner_region"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_body_get_ping_one_connections200_response_items_inner_region
  }
  property {
    name  = "createPingOneConnection_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_ping_one_connection_x_bypass_external_validation
  }
  property {
    name  = "createPolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerAuthenticationApiApplicationRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_authentication_api_application_ref_id
  }
  property {
    name  = "createPolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerAuthenticationApiApplicationRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_authentication_api_application_ref_location
  }
  property {
    name  = "createPolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerRootNodeAction_context"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_root_node_action_context
  }
  property {
    name  = "createPolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerRootNodeAction_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_root_node_action_type
  }
  property {
    name  = "createPolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerRootNode_children"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_root_node_children
  }
  property {
    name  = "createPolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_description
  }
  property {
    name  = "createPolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInner_enabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_enabled
  }
  property {
    name  = "createPolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInner_handleFailuresLocally"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_handle_failures_locally
  }
  property {
    name  = "createPolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_id
  }
  property {
    name  = "createPolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_name
  }
  property {
    name  = "createPolicy1_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy1_x_bypass_external_validation
  }
  property {
    name  = "createPolicy2_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "createPolicy2_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInnerAuthenticatorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_authenticator_ref_id
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInnerAuthenticatorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_authenticator_ref_location
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInnerIdentityHintContractFulfillment_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_identity_hint_contract_fulfillment_attribute_contract_fulfillment
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInnerIdentityHintContractFulfillment_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_identity_hint_contract_fulfillment_attribute_sources
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInnerIdentityHintContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_identity_hint_contract_core_attributes
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInnerIdentityHintContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_identity_hint_contract_extended_attributes
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInnerIdentityHintMapping_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_identity_hint_mapping_attribute_contract_fulfillment
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInnerIdentityHintMapping_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_identity_hint_mapping_attribute_sources
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInnerUserCodePcvRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_user_code_pcv_ref_id
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInnerUserCodePcvRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_user_code_pcv_ref_location
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInner_allowUnsignedLoginHintToken"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_allow_unsigned_login_hint_token
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInner_alternativeLoginHintTokenIssuers"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_alternative_login_hint_token_issuers
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_id
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_name
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInner_requireTokenForIdentityHint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_require_token_for_identity_hint
  }
  property {
    name  = "createPolicy2_body_GetPolicies1200ResponseItemsInner_transactionLifetime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_body_get_policies1200_response_items_inner_transaction_lifetime
  }
  property {
    name  = "createPolicy2_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy2_x_bypass_external_validation
  }
  property {
    name  = "createPolicy3_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "createPolicy3_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInnerAccessTokenManagerRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_access_token_manager_ref_id
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInnerAccessTokenManagerRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_access_token_manager_ref_location
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInnerAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_attribute_contract_core_attributes
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInnerAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_attribute_contract_extended_attributes
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInnerAttributeMapping_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_attribute_mapping_attribute_contract_fulfillment
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInnerAttributeMapping_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_attribute_mapping_attribute_sources
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_id
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInner_idTokenLifetime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_id_token_lifetime
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInner_includeSHashInIdToken"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_include_shash_in_id_token
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInner_includeSriInIdToken"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_include_sri_in_id_token
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInner_includeUserInfoInIdToken"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_include_user_info_in_id_token
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_name
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInner_reissueIdTokenInHybridFlow"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_reissue_id_token_in_hybrid_flow
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInner_returnIdTokenOnRefreshGrant"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_return_id_token_on_refresh_grant
  }
  property {
    name  = "createPolicy3_body_GetPolicies2200ResponseItemsInner_scopeAttributeMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_body_get_policies2200_response_items_inner_scope_attribute_mappings
  }
  property {
    name  = "createPolicy3_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy3_x_bypass_external_validation
  }
  property {
    name  = "createPolicy4_body_GetPolicies3200ResponseItemsInnerAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy4_body_get_policies3200_response_items_inner_attribute_contract_core_attributes
  }
  property {
    name  = "createPolicy4_body_GetPolicies3200ResponseItemsInnerAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy4_body_get_policies3200_response_items_inner_attribute_contract_extended_attributes
  }
  property {
    name  = "createPolicy4_body_GetPolicies3200ResponseItemsInner_actorTokenRequired"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy4_body_get_policies3200_response_items_inner_actor_token_required
  }
  property {
    name  = "createPolicy4_body_GetPolicies3200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy4_body_get_policies3200_response_items_inner_id
  }
  property {
    name  = "createPolicy4_body_GetPolicies3200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy4_body_get_policies3200_response_items_inner_name
  }
  property {
    name  = "createPolicy4_body_GetPolicies3200ResponseItemsInner_processorMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy4_body_get_policies3200_response_items_inner_processor_mappings
  }
  property {
    name  = "createPolicy4_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_policy4_bypass_external_validation
  }
  property {
    name  = "createResourceOwnerCredentialsMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_resource_owner_credentials_mapping_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "createResourceOwnerCredentialsMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_resource_owner_credentials_mapping_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "createResourceOwnerCredentialsMapping_body_GetResourceOwnerCredentialsMappings200ResponseItemsInnerPasswordValidatorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_resource_owner_credentials_mapping_body_get_resource_owner_credentials_mappings200_response_items_inner_password_validator_ref_id
  }
  property {
    name  = "createResourceOwnerCredentialsMapping_body_GetResourceOwnerCredentialsMappings200ResponseItemsInnerPasswordValidatorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_resource_owner_credentials_mapping_body_get_resource_owner_credentials_mappings200_response_items_inner_password_validator_ref_location
  }
  property {
    name  = "createResourceOwnerCredentialsMapping_body_GetResourceOwnerCredentialsMappings200ResponseItemsInner_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_resource_owner_credentials_mapping_body_get_resource_owner_credentials_mappings200_response_items_inner_attribute_contract_fulfillment
  }
  property {
    name  = "createResourceOwnerCredentialsMapping_body_GetResourceOwnerCredentialsMappings200ResponseItemsInner_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_resource_owner_credentials_mapping_body_get_resource_owner_credentials_mappings200_response_items_inner_attribute_sources
  }
  property {
    name  = "createResourceOwnerCredentialsMapping_body_GetResourceOwnerCredentialsMappings200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_resource_owner_credentials_mapping_body_get_resource_owner_credentials_mappings200_response_items_inner_id
  }
  property {
    name  = "createResourceOwnerCredentialsMapping_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_resource_owner_credentials_mapping_x_bypass_external_validation
  }
  property {
    name  = "createSecretManager_body_CreateSecretManagerRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_secret_manager_body_create_secret_manager_request_id
  }
  property {
    name  = "createSecretManager_body_CreateSecretManagerRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_secret_manager_body_create_secret_manager_request_name
  }
  property {
    name  = "createSecretManager_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_secret_manager_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "createSecretManager_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_secret_manager_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "createSecretManager_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_secret_manager_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "createSecretManager_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_secret_manager_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "createSecretManager_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_secret_manager_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "createSecretManager_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_secret_manager_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "createSourcePolicy_body_GetDefaultAuthenticationPolicy200ResponseDefaultAuthenticationSourcesInnerSourceRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_source_policy_body_get_default_authentication_policy200_response_default_authentication_sources_inner_source_ref_id
  }
  property {
    name  = "createSourcePolicy_body_GetDefaultAuthenticationPolicy200ResponseDefaultAuthenticationSourcesInnerSourceRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_source_policy_body_get_default_authentication_policy200_response_default_authentication_sources_inner_source_ref_location
  }
  property {
    name  = "createSourcePolicy_body_GetSourcePolicies200ResponseItemsInnerAuthenticationSource_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_source_policy_body_get_source_policies200_response_items_inner_authentication_source_type
  }
  property {
    name  = "createSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_authnContextSensitive"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_source_policy_body_get_source_policies200_response_items_inner_authn_context_sensitive
  }
  property {
    name  = "createSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_enableSessions"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_source_policy_body_get_source_policies200_response_items_inner_enable_sessions
  }
  property {
    name  = "createSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_source_policy_body_get_source_policies200_response_items_inner_id
  }
  property {
    name  = "createSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_idleTimeoutMins"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_source_policy_body_get_source_policies200_response_items_inner_idle_timeout_mins
  }
  property {
    name  = "createSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_maxTimeoutMins"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_source_policy_body_get_source_policies200_response_items_inner_max_timeout_mins
  }
  property {
    name  = "createSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_persistent"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_source_policy_body_get_source_policies200_response_items_inner_persistent
  }
  property {
    name  = "createSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_timeoutDisplayUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_source_policy_body_get_source_policies200_response_items_inner_timeout_display_unit
  }
  property {
    name  = "createSpAdapter_body_CreateSpAdapterRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_create_sp_adapter_request_id
  }
  property {
    name  = "createSpAdapter_body_CreateSpAdapterRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_create_sp_adapter_request_name
  }
  property {
    name  = "createSpAdapter_body_GetConnection1200ResponseIdpBrowserSsoAdapterMappingsInnerAdapterOverrideSettingsAllOfAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_get_connection1200_response_idp_browser_sso_adapter_mappings_inner_adapter_override_settings_all_of_attribute_contract_core_attributes
  }
  property {
    name  = "createSpAdapter_body_GetConnection1200ResponseIdpBrowserSsoAdapterMappingsInnerAdapterOverrideSettingsAllOfAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_get_connection1200_response_idp_browser_sso_adapter_mappings_inner_adapter_override_settings_all_of_attribute_contract_extended_attributes
  }
  property {
    name  = "createSpAdapter_body_GetConnection1200ResponseIdpBrowserSsoAdapterMappingsInnerAdapterOverrideSettingsAllOfAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_get_connection1200_response_idp_browser_sso_adapter_mappings_inner_adapter_override_settings_all_of_attribute_contract_inherited
  }
  property {
    name  = "createSpAdapter_body_GetConnection1200ResponseIdpBrowserSsoAdapterMappingsInnerAdapterOverrideSettingsAllOfTargetApplicationInfo_applicationIconUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_get_connection1200_response_idp_browser_sso_adapter_mappings_inner_adapter_override_settings_all_of_target_application_info_application_icon_url
  }
  property {
    name  = "createSpAdapter_body_GetConnection1200ResponseIdpBrowserSsoAdapterMappingsInnerAdapterOverrideSettingsAllOfTargetApplicationInfo_applicationName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_get_connection1200_response_idp_browser_sso_adapter_mappings_inner_adapter_override_settings_all_of_target_application_info_application_name
  }
  property {
    name  = "createSpAdapter_body_GetConnection1200ResponseIdpBrowserSsoAdapterMappingsInnerAdapterOverrideSettingsAllOfTargetApplicationInfo_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_get_connection1200_response_idp_browser_sso_adapter_mappings_inner_adapter_override_settings_all_of_target_application_info_inherited
  }
  property {
    name  = "createSpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "createSpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "createSpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "createSpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "createSpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "createSpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sp_adapter_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "createStsRequestParamContract_body_GetStsRequestParamContracts200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sts_request_param_contract_body_get_sts_request_param_contracts200_response_items_inner_id
  }
  property {
    name  = "createStsRequestParamContract_body_GetStsRequestParamContracts200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sts_request_param_contract_body_get_sts_request_param_contracts200_response_items_inner_name
  }
  property {
    name  = "createStsRequestParamContract_body_GetStsRequestParamContracts200ResponseItemsInner_parameters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_sts_request_param_contract_body_get_sts_request_param_contracts200_response_items_inner_parameters
  }
  property {
    name  = "createTokenGeneratorMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_mapping_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "createTokenGeneratorMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_mapping_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "createTokenGeneratorMapping_body_GetTokenGeneratorMappings200ResponseItemsInner_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_mapping_body_get_token_generator_mappings200_response_items_inner_attribute_contract_fulfillment
  }
  property {
    name  = "createTokenGeneratorMapping_body_GetTokenGeneratorMappings200ResponseItemsInner_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_mapping_body_get_token_generator_mappings200_response_items_inner_attribute_sources
  }
  property {
    name  = "createTokenGeneratorMapping_body_GetTokenGeneratorMappings200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_mapping_body_get_token_generator_mappings200_response_items_inner_id
  }
  property {
    name  = "createTokenGeneratorMapping_body_GetTokenGeneratorMappings200ResponseItemsInner_licenseConnectionGroupAssignment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_mapping_body_get_token_generator_mappings200_response_items_inner_license_connection_group_assignment
  }
  property {
    name  = "createTokenGeneratorMapping_body_GetTokenGeneratorMappings200ResponseItemsInner_sourceId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_mapping_body_get_token_generator_mappings200_response_items_inner_source_id
  }
  property {
    name  = "createTokenGeneratorMapping_body_GetTokenGeneratorMappings200ResponseItemsInner_targetId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_mapping_body_get_token_generator_mappings200_response_items_inner_target_id
  }
  property {
    name  = "createTokenGeneratorMapping_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_mapping_x_bypass_external_validation
  }
  property {
    name  = "createTokenGenerator_body_CreateTokenGeneratorRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_body_create_token_generator_request_id
  }
  property {
    name  = "createTokenGenerator_body_CreateTokenGeneratorRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_body_create_token_generator_request_name
  }
  property {
    name  = "createTokenGenerator_body_GetTokenGenerators200ResponseItemsInnerAllOfAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_body_get_token_generators200_response_items_inner_all_of_attribute_contract_core_attributes
  }
  property {
    name  = "createTokenGenerator_body_GetTokenGenerators200ResponseItemsInnerAllOfAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_body_get_token_generators200_response_items_inner_all_of_attribute_contract_extended_attributes
  }
  property {
    name  = "createTokenGenerator_body_GetTokenGenerators200ResponseItemsInnerAllOfAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_body_get_token_generators200_response_items_inner_all_of_attribute_contract_inherited
  }
  property {
    name  = "createTokenGenerator_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "createTokenGenerator_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "createTokenGenerator_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "createTokenGenerator_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "createTokenGenerator_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "createTokenGenerator_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_generator_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "createTokenManager_body_CreateTokenManagerRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_create_token_manager_request_id
  }
  property {
    name  = "createTokenManager_body_CreateTokenManagerRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_create_token_manager_request_name
  }
  property {
    name  = "createTokenManager_body_CreateTokenManagerRequest_sequenceNumber"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_create_token_manager_request_sequence_number
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AccessControlSettings_allowedClients"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_access_control_settings_allowed_clients
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AccessControlSettings_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_access_control_settings_inherited
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AccessControlSettings_restrictClients"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_access_control_settings_restrict_clients
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_attribute_contract_core_attributes
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AttributeContract_defaultSubjectAttribute"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_attribute_contract_default_subject_attribute
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_attribute_contract_extended_attributes
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_attribute_contract_inherited
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SelectionSettings_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_selection_settings_inherited
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SelectionSettings_resourceUris"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_selection_settings_resource_uris
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SessionValidationSettings_checkSessionRevocationStatus"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_session_validation_settings_check_session_revocation_status
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SessionValidationSettings_checkValidAuthnSession"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_session_validation_settings_check_valid_authn_session
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SessionValidationSettings_includeSessionId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_session_validation_settings_include_session_id
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SessionValidationSettings_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_session_validation_settings_inherited
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SessionValidationSettings_updateAuthnSessionActivity"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of1_session_validation_settings_update_authn_session_activity
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "createTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_manager_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "createTokenProcessor_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_processor_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "createTokenProcessor_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_processor_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "createTokenProcessor_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_processor_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "createTokenProcessor_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_processor_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "createTokenProcessor_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_processor_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "createTokenProcessor_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_processor_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "createTokenProcessor_body_GetTokenProcessor200ResponseAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_processor_body_get_token_processor200_response_attribute_contract_core_attributes
  }
  property {
    name  = "createTokenProcessor_body_GetTokenProcessor200ResponseAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_processor_body_get_token_processor200_response_attribute_contract_extended_attributes
  }
  property {
    name  = "createTokenProcessor_body_GetTokenProcessor200ResponseAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_processor_body_get_token_processor200_response_attribute_contract_inherited
  }
  property {
    name  = "createTokenProcessor_body_GetTokenProcessor200ResponseAttributeContract_maskOgnlValues"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_processor_body_get_token_processor200_response_attribute_contract_mask_ognl_values
  }
  property {
    name  = "createTokenProcessor_body_UpdateTokenProcessorRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_processor_body_update_token_processor_request_id
  }
  property {
    name  = "createTokenProcessor_body_UpdateTokenProcessorRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_processor_body_update_token_processor_request_name
  }
  property {
    name  = "createTokenToTokenMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_to_token_mapping_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "createTokenToTokenMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_to_token_mapping_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "createTokenToTokenMapping_body_GetTokenToTokenMappingById200Response_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_to_token_mapping_body_get_token_to_token_mapping_by_id200_response_attribute_contract_fulfillment
  }
  property {
    name  = "createTokenToTokenMapping_body_GetTokenToTokenMappingById200Response_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_to_token_mapping_body_get_token_to_token_mapping_by_id200_response_attribute_sources
  }
  property {
    name  = "createTokenToTokenMapping_body_GetTokenToTokenMappingById200Response_defaultTargetResource"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_to_token_mapping_body_get_token_to_token_mapping_by_id200_response_default_target_resource
  }
  property {
    name  = "createTokenToTokenMapping_body_GetTokenToTokenMappingById200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_to_token_mapping_body_get_token_to_token_mapping_by_id200_response_id
  }
  property {
    name  = "createTokenToTokenMapping_body_GetTokenToTokenMappingById200Response_licenseConnectionGroupAssignment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_to_token_mapping_body_get_token_to_token_mapping_by_id200_response_license_connection_group_assignment
  }
  property {
    name  = "createTokenToTokenMapping_body_GetTokenToTokenMappingById200Response_sourceId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_to_token_mapping_body_get_token_to_token_mapping_by_id200_response_source_id
  }
  property {
    name  = "createTokenToTokenMapping_body_GetTokenToTokenMappingById200Response_targetId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_to_token_mapping_body_get_token_to_token_mapping_by_id200_response_target_id
  }
  property {
    name  = "createTokenToTokenMapping_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_create_token_to_token_mapping_x_bypass_external_validation
  }
  property {
    name  = "deleteAccount_username"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_account_username
  }
  property {
    name  = "deleteApcMapping_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_apc_mapping_id
  }
  property {
    name  = "deleteApcToSpAdapterMappingById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_apc_to_sp_adapter_mapping_by_id_id
  }
  property {
    name  = "deleteApplication_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_application_id
  }
  property {
    name  = "deleteAuthenticationPolicyContract_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_authentication_policy_contract_id
  }
  property {
    name  = "deleteAuthenticationSelector_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_authentication_selector_id
  }
  property {
    name  = "deleteCertificateFromGroup_group_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_certificate_from_group_group_name
  }
  property {
    name  = "deleteCertificateFromGroup_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_certificate_from_group_id
  }
  property {
    name  = "deleteCertificate_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_certificate_id
  }
  property {
    name  = "deleteClient_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_client_id
  }
  property {
    name  = "deleteConnection1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_connection1_id
  }
  property {
    name  = "deleteConnection2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_connection2_id
  }
  property {
    name  = "deleteDataStore_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_data_store_id
  }
  property {
    name  = "deleteDynamicClientRegistrationPolicy_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_dynamic_client_registration_policy_id
  }
  property {
    name  = "deleteFragment_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_fragment_id
  }
  property {
    name  = "deleteGroup_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_group_id
  }
  property {
    name  = "deleteIdentityProfile_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_identity_profile_id
  }
  property {
    name  = "deleteIdentityStoreProvisioner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_identity_store_provisioner_id
  }
  property {
    name  = "deleteIdpAdapterMapping_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_idp_adapter_mapping_id
  }
  property {
    name  = "deleteIdpAdapter_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_idp_adapter_id
  }
  property {
    name  = "deleteIdpToSpAdapterMappingsById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_idp_to_sp_adapter_mappings_by_id_id
  }
  property {
    name  = "deleteIssuer_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_issuer_id
  }
  property {
    name  = "deleteKerberosRealm_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_kerberos_realm_id
  }
  property {
    name  = "deleteKeyPair1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_key_pair1_id
  }
  property {
    name  = "deleteKeyPair2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_key_pair2_id
  }
  property {
    name  = "deleteKeyPair3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_key_pair3_id
  }
  property {
    name  = "deleteKeyPairRotationSettings_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_key_pair_rotation_settings_id
  }
  property {
    name  = "deleteKeySet_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_key_set_id
  }
  property {
    name  = "deleteMapping_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_mapping_id
  }
  property {
    name  = "deleteMetadataUrl_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_metadata_url_id
  }
  property {
    name  = "deleteNotificationPublisher_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_notification_publisher_id
  }
  property {
    name  = "deleteOOBAuthenticator_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_oobauthenticator_id
  }
  property {
    name  = "deleteOcspCertificateById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_ocsp_certificate_by_id_id
  }
  property {
    name  = "deletePasswordCredentialValidator_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_password_credential_validator_id
  }
  property {
    name  = "deletePingOneConnection_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_ping_one_connection_id
  }
  property {
    name  = "deletePolicy1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_policy1_id
  }
  property {
    name  = "deletePolicy2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_policy2_id
  }
  property {
    name  = "deletePolicy3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_policy3_id
  }
  property {
    name  = "deletePolicy4_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_policy4_id
  }
  property {
    name  = "deleteResourceOwnerCredentialsMapping_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_resource_owner_credentials_mapping_id
  }
  property {
    name  = "deleteSecretManager_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_secret_manager_id
  }
  property {
    name  = "deleteSetting_bundle"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_setting_bundle
  }
  property {
    name  = "deleteSetting_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_setting_id
  }
  property {
    name  = "deleteSourcePolicy_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_source_policy_id
  }
  property {
    name  = "deleteSpAdapter_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_sp_adapter_id
  }
  property {
    name  = "deleteStsRequestParamContractById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_sts_request_param_contract_by_id_id
  }
  property {
    name  = "deleteTokenGeneratorMappingById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_token_generator_mapping_by_id_id
  }
  property {
    name  = "deleteTokenGenerator_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_token_generator_id
  }
  property {
    name  = "deleteTokenManager_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_token_manager_id
  }
  property {
    name  = "deleteTokenProcessor_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_token_processor_id
  }
  property {
    name  = "deleteTokenToTokenMappingById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_token_to_token_mapping_by_id_id
  }
  property {
    name  = "deleteTrustedCA_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_delete_trusted_ca_id
  }
  property {
    name  = "exportCertificateFile1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_certificate_file1_id
  }
  property {
    name  = "exportCertificateFile2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_certificate_file2_id
  }
  property {
    name  = "exportCertificateFile3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_certificate_file3_id
  }
  property {
    name  = "exportCertificateFile4_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_certificate_file4_id
  }
  property {
    name  = "exportConfiguration_include_external_resources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_configuration_include_external_resources
  }
  property {
    name  = "exportCsr1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_csr1_id
  }
  property {
    name  = "exportCsr2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_csr2_id
  }
  property {
    name  = "exportCsr3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_csr3_id
  }
  property {
    name  = "exportPEMFile1_body_ExportPKCS12File1Request_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_pemfile1_body_export_pkcs12_file1_request_password
  }
  property {
    name  = "exportPEMFile1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_pemfile1_id
  }
  property {
    name  = "exportPEMFile2_body_ExportPKCS12File1Request_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_pemfile2_body_export_pkcs12_file1_request_password
  }
  property {
    name  = "exportPEMFile2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_pemfile2_id
  }
  property {
    name  = "exportPEMFile3_body_ExportPKCS12File1Request_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_pemfile3_body_export_pkcs12_file1_request_password
  }
  property {
    name  = "exportPEMFile3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_pemfile3_id
  }
  property {
    name  = "exportPKCS12File1_body_ExportPKCS12File1Request_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_pkcs12_file1_body_export_pkcs12_file1_request_password
  }
  property {
    name  = "exportPKCS12File1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_pkcs12_file1_id
  }
  property {
    name  = "exportPKCS12File2_body_ExportPKCS12File1Request_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_pkcs12_file2_body_export_pkcs12_file1_request_password
  }
  property {
    name  = "exportPKCS12File2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_pkcs12_file2_id
  }
  property {
    name  = "exportPKCS12File3_body_ExportPKCS12File1Request_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_pkcs12_file3_body_export_pkcs12_file1_request_password
  }
  property {
    name  = "exportPKCS12File3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_export_pkcs12_file3_id
  }
  property {
    name  = "getAccount_username"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_account_username
  }
  property {
    name  = "getAction1_action_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_action1_action_id
  }
  property {
    name  = "getAction1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_action1_id
  }
  property {
    name  = "getAction2_action_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_action2_action_id
  }
  property {
    name  = "getAction2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_action2_id
  }
  property {
    name  = "getAction3_action_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_action3_action_id
  }
  property {
    name  = "getAction3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_action3_id
  }
  property {
    name  = "getAction4_action_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_action4_action_id
  }
  property {
    name  = "getAction4_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_action4_id
  }
  property {
    name  = "getAction5_action_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_action5_action_id
  }
  property {
    name  = "getAction5_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_action5_id
  }
  property {
    name  = "getAction6_action_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_action6_action_id
  }
  property {
    name  = "getAction6_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_action6_id
  }
  property {
    name  = "getActions1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_actions1_id
  }
  property {
    name  = "getActions2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_actions2_id
  }
  property {
    name  = "getActions3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_actions3_id
  }
  property {
    name  = "getActions4_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_actions4_id
  }
  property {
    name  = "getActions5_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_actions5_id
  }
  property {
    name  = "getActions6_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_actions6_id
  }
  property {
    name  = "getApcMapping_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_apc_mapping_id
  }
  property {
    name  = "getApcToSpAdapterMappingById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_apc_to_sp_adapter_mapping_by_id_id
  }
  property {
    name  = "getApplication_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_application_id
  }
  property {
    name  = "getAuthenticationPolicyContract_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_authentication_policy_contract_id
  }
  property {
    name  = "getAuthenticationPolicyContracts_filter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_authentication_policy_contracts_filter
  }
  property {
    name  = "getAuthenticationPolicyContracts_number_per_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_authentication_policy_contracts_number_per_page
  }
  property {
    name  = "getAuthenticationPolicyContracts_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_authentication_policy_contracts_page
  }
  property {
    name  = "getAuthenticationSelectorDescriptorsById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_authentication_selector_descriptors_by_id_id
  }
  property {
    name  = "getAuthenticationSelector_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_authentication_selector_id
  }
  property {
    name  = "getAuthenticationSelectors_filter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_authentication_selectors_filter
  }
  property {
    name  = "getAuthenticationSelectors_number_per_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_authentication_selectors_number_per_page
  }
  property {
    name  = "getAuthenticationSelectors_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_authentication_selectors_page
  }
  property {
    name  = "getCert_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_cert_id
  }
  property {
    name  = "getCertificateFromGroup_group_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_certificate_from_group_group_name
  }
  property {
    name  = "getCertificateFromGroup_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_certificate_from_group_id
  }
  property {
    name  = "getCertificatesForGroup_group_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_certificates_for_group_group_name
  }
  property {
    name  = "getClientSecret_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_client_secret_id
  }
  property {
    name  = "getClient_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_client_id
  }
  property {
    name  = "getClients_filter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_clients_filter
  }
  property {
    name  = "getClients_number_per_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_clients_number_per_page
  }
  property {
    name  = "getClients_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_clients_page
  }
  property {
    name  = "getCommonScopeGroup_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_common_scope_group_name
  }
  property {
    name  = "getCommonScope_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_common_scope_name
  }
  property {
    name  = "getConnection1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_connection1_id
  }
  property {
    name  = "getConnection2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_connection2_id
  }
  property {
    name  = "getConnectionCerts1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_connection_certs1_id
  }
  property {
    name  = "getConnectionCerts2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_connection_certs2_id
  }
  property {
    name  = "getConnections1_entity_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_connections1_entity_id
  }
  property {
    name  = "getConnections1_filter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_connections1_filter
  }
  property {
    name  = "getConnections1_number_per_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_connections1_number_per_page
  }
  property {
    name  = "getConnections1_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_connections1_page
  }
  property {
    name  = "getConnections2_entity_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_connections2_entity_id
  }
  property {
    name  = "getConnections2_filter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_connections2_filter
  }
  property {
    name  = "getConnections2_number_per_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_connections2_number_per_page
  }
  property {
    name  = "getConnections2_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_connections2_page
  }
  property {
    name  = "getCredentialStatus_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_credential_status_id
  }
  property {
    name  = "getCustomDataStoreDescriptor_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_custom_data_store_descriptor_id
  }
  property {
    name  = "getDataStore_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_data_store_id
  }
  property {
    name  = "getDecryptionKeys1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_decryption_keys1_id
  }
  property {
    name  = "getDecryptionKeys2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_decryption_keys2_id
  }
  property {
    name  = "getDynamicClientRegistrationDescriptor_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_dynamic_client_registration_descriptor_id
  }
  property {
    name  = "getDynamicClientRegistrationPolicy_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_dynamic_client_registration_policy_id
  }
  property {
    name  = "getExclusiveScopeGroup_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_exclusive_scope_group_name
  }
  property {
    name  = "getExclusiveScope_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_exclusive_scope_name
  }
  property {
    name  = "getFragment_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_fragment_id
  }
  property {
    name  = "getFragments_filter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_fragments_filter
  }
  property {
    name  = "getFragments_number_per_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_fragments_number_per_page
  }
  property {
    name  = "getFragments_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_fragments_page
  }
  property {
    name  = "getGroup_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_group_id
  }
  property {
    name  = "getIdentityProfile_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_identity_profile_id
  }
  property {
    name  = "getIdentityProfiles_apc_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_identity_profiles_apc_id
  }
  property {
    name  = "getIdentityProfiles_filter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_identity_profiles_filter
  }
  property {
    name  = "getIdentityProfiles_number_per_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_identity_profiles_number_per_page
  }
  property {
    name  = "getIdentityProfiles_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_identity_profiles_page
  }
  property {
    name  = "getIdentityStoreProvisionerDescriptorById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_identity_store_provisioner_descriptor_by_id_id
  }
  property {
    name  = "getIdentityStoreProvisioner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_identity_store_provisioner_id
  }
  property {
    name  = "getIdpAdapterDescriptorsById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_idp_adapter_descriptors_by_id_id
  }
  property {
    name  = "getIdpAdapterMapping_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_idp_adapter_mapping_id
  }
  property {
    name  = "getIdpAdapter_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_idp_adapter_id
  }
  property {
    name  = "getIdpAdapters_filter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_idp_adapters_filter
  }
  property {
    name  = "getIdpAdapters_number_per_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_idp_adapters_number_per_page
  }
  property {
    name  = "getIdpAdapters_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_idp_adapters_page
  }
  property {
    name  = "getIdpConnectorDescriptorById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_idp_connector_descriptor_by_id_id
  }
  property {
    name  = "getIdpToSpAdapterMappingsById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_idp_to_sp_adapter_mappings_by_id_id
  }
  property {
    name  = "getIssuerById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_issuer_by_id_id
  }
  property {
    name  = "getKerberosRealm_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_kerberos_realm_id
  }
  property {
    name  = "getKeyPair1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_key_pair1_id
  }
  property {
    name  = "getKeyPair2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_key_pair2_id
  }
  property {
    name  = "getKeyPair3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_key_pair3_id
  }
  property {
    name  = "getKeySet_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_key_set_id
  }
  property {
    name  = "getMapping_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_mapping_id
  }
  property {
    name  = "getMetadataUrl_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_metadata_url_id
  }
  property {
    name  = "getNotificationPublisherPluginDescriptor_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_notification_publisher_plugin_descriptor_id
  }
  property {
    name  = "getNotificationPublisher_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_notification_publisher_id
  }
  property {
    name  = "getOOBAuthPluginDescriptor_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_oobauth_plugin_descriptor_id
  }
  property {
    name  = "getOOBAuthenticator_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_oobauthenticator_id
  }
  property {
    name  = "getOcspCertificateById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_ocsp_certificate_by_id_id
  }
  property {
    name  = "getPasswordCredentialValidatorDescriptor_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_password_credential_validator_descriptor_id
  }
  property {
    name  = "getPasswordCredentialValidator_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_password_credential_validator_id
  }
  property {
    name  = "getPingOneConnectionAssociations_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_ping_one_connection_associations_id
  }
  property {
    name  = "getPingOneConnectionEnvironments_filter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_ping_one_connection_environments_filter
  }
  property {
    name  = "getPingOneConnectionEnvironments_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_ping_one_connection_environments_id
  }
  property {
    name  = "getPingOneConnectionEnvironments_number_per_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_ping_one_connection_environments_number_per_page
  }
  property {
    name  = "getPingOneConnectionEnvironments_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_ping_one_connection_environments_page
  }
  property {
    name  = "getPingOneConnectionUsages_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_ping_one_connection_usages_id
  }
  property {
    name  = "getPingOneConnection_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_ping_one_connection_id
  }
  property {
    name  = "getPolicy1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_policy1_id
  }
  property {
    name  = "getPolicy2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_policy2_id
  }
  property {
    name  = "getPolicy3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_policy3_id
  }
  property {
    name  = "getPolicy4_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_policy4_id
  }
  property {
    name  = "getResourceOwnerCredentialsMapping_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_resource_owner_credentials_mapping_id
  }
  property {
    name  = "getRotationSettings_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_rotation_settings_id
  }
  property {
    name  = "getSecretManagerPluginDescriptor_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_secret_manager_plugin_descriptor_id
  }
  property {
    name  = "getSecretManager_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_secret_manager_id
  }
  property {
    name  = "getSetting_bundle"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_setting_bundle
  }
  property {
    name  = "getSetting_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_setting_id
  }
  property {
    name  = "getSettings3_bundle"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_settings3_bundle
  }
  property {
    name  = "getSigningSettings1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_signing_settings1_id
  }
  property {
    name  = "getSigningSettings3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_signing_settings3_id
  }
  property {
    name  = "getSourcePolicy_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_source_policy_id
  }
  property {
    name  = "getSpAdapterDescriptorsById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_sp_adapter_descriptors_by_id_id
  }
  property {
    name  = "getSpAdapter_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_sp_adapter_id
  }
  property {
    name  = "getSpAdapters_filter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_sp_adapters_filter
  }
  property {
    name  = "getSpAdapters_number_per_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_sp_adapters_number_per_page
  }
  property {
    name  = "getSpAdapters_page"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_sp_adapters_page
  }
  property {
    name  = "getStsRequestParamContractById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_sts_request_param_contract_by_id_id
  }
  property {
    name  = "getTokenGeneratorDescriptorsById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_token_generator_descriptors_by_id_id
  }
  property {
    name  = "getTokenGeneratorMappingById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_token_generator_mapping_by_id_id
  }
  property {
    name  = "getTokenGenerator_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_token_generator_id
  }
  property {
    name  = "getTokenManagerDescriptor_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_token_manager_descriptor_id
  }
  property {
    name  = "getTokenManager_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_token_manager_id
  }
  property {
    name  = "getTokenProcessorDescriptorsById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_token_processor_descriptors_by_id_id
  }
  property {
    name  = "getTokenProcessor_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_token_processor_id
  }
  property {
    name  = "getTokenToTokenMappingById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_token_to_token_mapping_by_id_id
  }
  property {
    name  = "getTrustedCert_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_get_trusted_cert_id
  }
  property {
    name  = "importCertificate_body_ImportFeatureCertRequest_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_certificate_body_import_feature_cert_request_crypto_provider
  }
  property {
    name  = "importCertificate_body_ImportFeatureCertRequest_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_certificate_body_import_feature_cert_request_file_data
  }
  property {
    name  = "importCertificate_body_ImportFeatureCertRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_certificate_body_import_feature_cert_request_id
  }
  property {
    name  = "importConfigArchive_file"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_config_archive_file
  }
  property {
    name  = "importConfigArchive_force_import"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_config_archive_force_import
  }
  property {
    name  = "importConfigArchive_force_unsupported_import"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_config_archive_force_unsupported_import
  }
  property {
    name  = "importConfigArchive_reencrypt_data"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_config_archive_reencrypt_data
  }
  property {
    name  = "importConfiguration_body_ExportConfiguration200ResponseMetadata_pfVersion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_configuration_body_export_configuration200_response_metadata_pf_version
  }
  property {
    name  = "importConfiguration_body_ExportConfiguration200Response_operations"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_configuration_body_export_configuration200_response_operations
  }
  property {
    name  = "importConfiguration_fail_fast"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_configuration_fail_fast
  }
  property {
    name  = "importConfiguration_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_configuration_x_bypass_external_validation
  }
  property {
    name  = "importCsrResponse1_body_ImportCsrResponse1Request_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_csr_response1_body_import_csr_response1_request_file_data
  }
  property {
    name  = "importCsrResponse1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_csr_response1_id
  }
  property {
    name  = "importCsrResponse2_body_ImportCsrResponse1Request_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_csr_response2_body_import_csr_response1_request_file_data
  }
  property {
    name  = "importCsrResponse2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_csr_response2_id
  }
  property {
    name  = "importCsrResponse3_body_ImportCsrResponse1Request_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_csr_response3_body_import_csr_response1_request_file_data
  }
  property {
    name  = "importCsrResponse3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_csr_response3_id
  }
  property {
    name  = "importFeatureCert_body_ImportFeatureCertRequest_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_feature_cert_body_import_feature_cert_request_crypto_provider
  }
  property {
    name  = "importFeatureCert_body_ImportFeatureCertRequest_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_feature_cert_body_import_feature_cert_request_file_data
  }
  property {
    name  = "importFeatureCert_body_ImportFeatureCertRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_feature_cert_body_import_feature_cert_request_id
  }
  property {
    name  = "importFeatureCert_group_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_feature_cert_group_name
  }
  property {
    name  = "importKeyPair1_body_ImportKeyPair1Request_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair1_body_import_key_pair1_request_crypto_provider
  }
  property {
    name  = "importKeyPair1_body_ImportKeyPair1Request_encryptedPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair1_body_import_key_pair1_request_encrypted_password
  }
  property {
    name  = "importKeyPair1_body_ImportKeyPair1Request_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair1_body_import_key_pair1_request_file_data
  }
  property {
    name  = "importKeyPair1_body_ImportKeyPair1Request_format"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair1_body_import_key_pair1_request_format
  }
  property {
    name  = "importKeyPair1_body_ImportKeyPair1Request_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair1_body_import_key_pair1_request_id
  }
  property {
    name  = "importKeyPair1_body_ImportKeyPair1Request_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair1_body_import_key_pair1_request_password
  }
  property {
    name  = "importKeyPair2_body_ImportKeyPair1Request_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair2_body_import_key_pair1_request_crypto_provider
  }
  property {
    name  = "importKeyPair2_body_ImportKeyPair1Request_encryptedPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair2_body_import_key_pair1_request_encrypted_password
  }
  property {
    name  = "importKeyPair2_body_ImportKeyPair1Request_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair2_body_import_key_pair1_request_file_data
  }
  property {
    name  = "importKeyPair2_body_ImportKeyPair1Request_format"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair2_body_import_key_pair1_request_format
  }
  property {
    name  = "importKeyPair2_body_ImportKeyPair1Request_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair2_body_import_key_pair1_request_id
  }
  property {
    name  = "importKeyPair2_body_ImportKeyPair1Request_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair2_body_import_key_pair1_request_password
  }
  property {
    name  = "importKeyPair3_body_ImportKeyPair1Request_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair3_body_import_key_pair1_request_crypto_provider
  }
  property {
    name  = "importKeyPair3_body_ImportKeyPair1Request_encryptedPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair3_body_import_key_pair1_request_encrypted_password
  }
  property {
    name  = "importKeyPair3_body_ImportKeyPair1Request_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair3_body_import_key_pair1_request_file_data
  }
  property {
    name  = "importKeyPair3_body_ImportKeyPair1Request_format"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair3_body_import_key_pair1_request_format
  }
  property {
    name  = "importKeyPair3_body_ImportKeyPair1Request_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair3_body_import_key_pair1_request_id
  }
  property {
    name  = "importKeyPair3_body_ImportKeyPair1Request_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_key_pair3_body_import_key_pair1_request_password
  }
  property {
    name  = "importOcspCertificate_body_ImportFeatureCertRequest_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_ocsp_certificate_body_import_feature_cert_request_crypto_provider
  }
  property {
    name  = "importOcspCertificate_body_ImportFeatureCertRequest_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_ocsp_certificate_body_import_feature_cert_request_file_data
  }
  property {
    name  = "importOcspCertificate_body_ImportFeatureCertRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_ocsp_certificate_body_import_feature_cert_request_id
  }
  property {
    name  = "importTrustedCA_body_ImportFeatureCertRequest_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_trusted_ca_body_import_feature_cert_request_crypto_provider
  }
  property {
    name  = "importTrustedCA_body_ImportFeatureCertRequest_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_trusted_ca_body_import_feature_cert_request_file_data
  }
  property {
    name  = "importTrustedCA_body_ImportFeatureCertRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_import_trusted_ca_body_import_feature_cert_request_id
  }
  property {
    name  = "invokeActionWithOptions1_action_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options1_action_id
  }
  property {
    name  = "invokeActionWithOptions1_body_InvokeActionWithOptions1Request_parameters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options1_body_invoke_action_with_options1_request_parameters
  }
  property {
    name  = "invokeActionWithOptions1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options1_id
  }
  property {
    name  = "invokeActionWithOptions2_action_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options2_action_id
  }
  property {
    name  = "invokeActionWithOptions2_body_InvokeActionWithOptions1Request_parameters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options2_body_invoke_action_with_options1_request_parameters
  }
  property {
    name  = "invokeActionWithOptions2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options2_id
  }
  property {
    name  = "invokeActionWithOptions3_action_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options3_action_id
  }
  property {
    name  = "invokeActionWithOptions3_body_InvokeActionWithOptions1Request_parameters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options3_body_invoke_action_with_options1_request_parameters
  }
  property {
    name  = "invokeActionWithOptions3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options3_id
  }
  property {
    name  = "invokeActionWithOptions4_action_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options4_action_id
  }
  property {
    name  = "invokeActionWithOptions4_body_InvokeActionWithOptions1Request_parameters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options4_body_invoke_action_with_options1_request_parameters
  }
  property {
    name  = "invokeActionWithOptions4_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options4_id
  }
  property {
    name  = "invokeActionWithOptions5_action_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options5_action_id
  }
  property {
    name  = "invokeActionWithOptions5_body_InvokeActionWithOptions1Request_parameters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options5_body_invoke_action_with_options1_request_parameters
  }
  property {
    name  = "invokeActionWithOptions5_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options5_id
  }
  property {
    name  = "invokeActionWithOptions6_action_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options6_action_id
  }
  property {
    name  = "invokeActionWithOptions6_body_InvokeActionWithOptions1Request_parameters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options6_body_invoke_action_with_options1_request_parameters
  }
  property {
    name  = "invokeActionWithOptions6_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_invoke_action_with_options6_id
  }
  property {
    name  = "movePolicy_body_MovePolicyRequest_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_move_policy_body_move_policy_request_location
  }
  property {
    name  = "movePolicy_body_MovePolicyRequest_moveToId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_move_policy_body_move_policy_request_move_to_id
  }
  property {
    name  = "movePolicy_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_move_policy_id
  }
  property {
    name  = "removeCommonScopeGroup_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_remove_common_scope_group_name
  }
  property {
    name  = "removeCommonScope_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_remove_common_scope_name
  }
  property {
    name  = "removeExclusiveScopeGroup_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_remove_exclusive_scope_group_name
  }
  property {
    name  = "removeExclusiveScope_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_remove_exclusive_scope_name
  }
  property {
    name  = "resetPassword_body_ResetPasswordRequest_currentPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_reset_password_body_reset_password_request_current_password
  }
  property {
    name  = "resetPassword_body_ResetPasswordRequest_newPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_reset_password_body_reset_password_request_new_password
  }
  property {
    name  = "resetPassword_username"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_reset_password_username
  }
  property {
    name  = "sslVerification"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_ssl_verification
  }
  property {
    name  = "updateAccount_body_GetAccounts200ResponseItemsInner_active"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_account_body_get_accounts200_response_items_inner_active
  }
  property {
    name  = "updateAccount_body_GetAccounts200ResponseItemsInner_auditor"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_account_body_get_accounts200_response_items_inner_auditor
  }
  property {
    name  = "updateAccount_body_GetAccounts200ResponseItemsInner_department"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_account_body_get_accounts200_response_items_inner_department
  }
  property {
    name  = "updateAccount_body_GetAccounts200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_account_body_get_accounts200_response_items_inner_description
  }
  property {
    name  = "updateAccount_body_GetAccounts200ResponseItemsInner_emailAddress"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_account_body_get_accounts200_response_items_inner_email_address
  }
  property {
    name  = "updateAccount_body_GetAccounts200ResponseItemsInner_encryptedPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_account_body_get_accounts200_response_items_inner_encrypted_password
  }
  property {
    name  = "updateAccount_body_GetAccounts200ResponseItemsInner_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_account_body_get_accounts200_response_items_inner_password
  }
  property {
    name  = "updateAccount_body_GetAccounts200ResponseItemsInner_phoneNumber"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_account_body_get_accounts200_response_items_inner_phone_number
  }
  property {
    name  = "updateAccount_body_GetAccounts200ResponseItemsInner_roles"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_account_body_get_accounts200_response_items_inner_roles
  }
  property {
    name  = "updateAccount_body_GetAccounts200ResponseItemsInner_username"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_account_body_get_accounts200_response_items_inner_username
  }
  property {
    name  = "updateAccount_username"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_account_username
  }
  property {
    name  = "updateApcMapping_body_GetApcMappings200ResponseItemsInnerAuthenticationPolicyContractRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_mapping_body_get_apc_mappings200_response_items_inner_authentication_policy_contract_ref_id
  }
  property {
    name  = "updateApcMapping_body_GetApcMappings200ResponseItemsInnerAuthenticationPolicyContractRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_mapping_body_get_apc_mappings200_response_items_inner_authentication_policy_contract_ref_location
  }
  property {
    name  = "updateApcMapping_body_GetApcMappings200ResponseItemsInner_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_mapping_body_get_apc_mappings200_response_items_inner_attribute_contract_fulfillment
  }
  property {
    name  = "updateApcMapping_body_GetApcMappings200ResponseItemsInner_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_mapping_body_get_apc_mappings200_response_items_inner_attribute_sources
  }
  property {
    name  = "updateApcMapping_body_GetApcMappings200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_mapping_body_get_apc_mappings200_response_items_inner_id
  }
  property {
    name  = "updateApcMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_mapping_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "updateApcMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_mapping_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "updateApcMapping_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_mapping_id
  }
  property {
    name  = "updateApcMapping_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_mapping_x_bypass_external_validation
  }
  property {
    name  = "updateApcToSpAdapterMappingById_body_GetApcToSpAdapterMappings200ResponseItemsInner_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_to_sp_adapter_mapping_by_id_body_get_apc_to_sp_adapter_mappings200_response_items_inner_attribute_contract_fulfillment
  }
  property {
    name  = "updateApcToSpAdapterMappingById_body_GetApcToSpAdapterMappings200ResponseItemsInner_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_to_sp_adapter_mapping_by_id_body_get_apc_to_sp_adapter_mappings200_response_items_inner_attribute_sources
  }
  property {
    name  = "updateApcToSpAdapterMappingById_body_GetApcToSpAdapterMappings200ResponseItemsInner_defaultTargetResource"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_to_sp_adapter_mapping_by_id_body_get_apc_to_sp_adapter_mappings200_response_items_inner_default_target_resource
  }
  property {
    name  = "updateApcToSpAdapterMappingById_body_GetApcToSpAdapterMappings200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_to_sp_adapter_mapping_by_id_body_get_apc_to_sp_adapter_mappings200_response_items_inner_id
  }
  property {
    name  = "updateApcToSpAdapterMappingById_body_GetApcToSpAdapterMappings200ResponseItemsInner_licenseConnectionGroupAssignment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_to_sp_adapter_mapping_by_id_body_get_apc_to_sp_adapter_mappings200_response_items_inner_license_connection_group_assignment
  }
  property {
    name  = "updateApcToSpAdapterMappingById_body_GetApcToSpAdapterMappings200ResponseItemsInner_sourceId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_to_sp_adapter_mapping_by_id_body_get_apc_to_sp_adapter_mappings200_response_items_inner_source_id
  }
  property {
    name  = "updateApcToSpAdapterMappingById_body_GetApcToSpAdapterMappings200ResponseItemsInner_targetId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_to_sp_adapter_mapping_by_id_body_get_apc_to_sp_adapter_mappings200_response_items_inner_target_id
  }
  property {
    name  = "updateApcToSpAdapterMappingById_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_to_sp_adapter_mapping_by_id_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "updateApcToSpAdapterMappingById_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_to_sp_adapter_mapping_by_id_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "updateApcToSpAdapterMappingById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_to_sp_adapter_mapping_by_id_id
  }
  property {
    name  = "updateApcToSpAdapterMappingById_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_apc_to_sp_adapter_mapping_by_id_x_bypass_external_validation
  }
  property {
    name  = "updateApplicationPolicy_body_GetApplicationPolicy200Response_idleTimeoutMins"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_application_policy_body_get_application_policy200_response_idle_timeout_mins
  }
  property {
    name  = "updateApplicationPolicy_body_GetApplicationPolicy200Response_maxTimeoutMins"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_application_policy_body_get_application_policy200_response_max_timeout_mins
  }
  property {
    name  = "updateApplication_body_GetAuthenticationApiApplications200ResponseItemsInnerClientForRedirectlessModeRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_application_body_get_authentication_api_applications200_response_items_inner_client_for_redirectless_mode_ref_id
  }
  property {
    name  = "updateApplication_body_GetAuthenticationApiApplications200ResponseItemsInnerClientForRedirectlessModeRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_application_body_get_authentication_api_applications200_response_items_inner_client_for_redirectless_mode_ref_location
  }
  property {
    name  = "updateApplication_body_GetAuthenticationApiApplications200ResponseItemsInner_additionalAllowedOrigins"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_application_body_get_authentication_api_applications200_response_items_inner_additional_allowed_origins
  }
  property {
    name  = "updateApplication_body_GetAuthenticationApiApplications200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_application_body_get_authentication_api_applications200_response_items_inner_description
  }
  property {
    name  = "updateApplication_body_GetAuthenticationApiApplications200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_application_body_get_authentication_api_applications200_response_items_inner_id
  }
  property {
    name  = "updateApplication_body_GetAuthenticationApiApplications200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_application_body_get_authentication_api_applications200_response_items_inner_name
  }
  property {
    name  = "updateApplication_body_GetAuthenticationApiApplications200ResponseItemsInner_url"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_application_body_get_authentication_api_applications200_response_items_inner_url
  }
  property {
    name  = "updateApplication_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_application_id
  }
  property {
    name  = "updateAuthenticationApiSettings_body_GetAuthenticationApiSettings200ResponseDefaultApplicationRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_api_settings_body_get_authentication_api_settings200_response_default_application_ref_id
  }
  property {
    name  = "updateAuthenticationApiSettings_body_GetAuthenticationApiSettings200ResponseDefaultApplicationRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_api_settings_body_get_authentication_api_settings200_response_default_application_ref_location
  }
  property {
    name  = "updateAuthenticationApiSettings_body_GetAuthenticationApiSettings200Response_apiEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_api_settings_body_get_authentication_api_settings200_response_api_enabled
  }
  property {
    name  = "updateAuthenticationApiSettings_body_GetAuthenticationApiSettings200Response_enableApiDescriptions"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_api_settings_body_get_authentication_api_settings200_response_enable_api_descriptions
  }
  property {
    name  = "updateAuthenticationApiSettings_body_GetAuthenticationApiSettings200Response_includeRequestContext"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_api_settings_body_get_authentication_api_settings200_response_include_request_context
  }
  property {
    name  = "updateAuthenticationApiSettings_body_GetAuthenticationApiSettings200Response_restrictAccessToRedirectlessMode"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_api_settings_body_get_authentication_api_settings200_response_restrict_access_to_redirectless_mode
  }
  property {
    name  = "updateAuthenticationPolicyContract_body_GetAuthenticationPolicyContracts200ResponseItemsInner_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_policy_contract_body_get_authentication_policy_contracts200_response_items_inner_core_attributes
  }
  property {
    name  = "updateAuthenticationPolicyContract_body_GetAuthenticationPolicyContracts200ResponseItemsInner_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_policy_contract_body_get_authentication_policy_contracts200_response_items_inner_extended_attributes
  }
  property {
    name  = "updateAuthenticationPolicyContract_body_GetAuthenticationPolicyContracts200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_policy_contract_body_get_authentication_policy_contracts200_response_items_inner_id
  }
  property {
    name  = "updateAuthenticationPolicyContract_body_GetAuthenticationPolicyContracts200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_policy_contract_body_get_authentication_policy_contracts200_response_items_inner_name
  }
  property {
    name  = "updateAuthenticationPolicyContract_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_policy_contract_id
  }
  property {
    name  = "updateAuthenticationSelector_body_CreateAuthenticationSelectorRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_selector_body_create_authentication_selector_request_id
  }
  property {
    name  = "updateAuthenticationSelector_body_CreateAuthenticationSelectorRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_selector_body_create_authentication_selector_request_name
  }
  property {
    name  = "updateAuthenticationSelector_body_GetAuthenticationSelectors200ResponseItemsInnerAllOfAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_selector_body_get_authentication_selectors200_response_items_inner_all_of_attribute_contract_extended_attributes
  }
  property {
    name  = "updateAuthenticationSelector_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_selector_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "updateAuthenticationSelector_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_selector_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "updateAuthenticationSelector_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_selector_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "updateAuthenticationSelector_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_selector_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "updateAuthenticationSelector_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_selector_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "updateAuthenticationSelector_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_selector_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "updateAuthenticationSelector_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authentication_selector_id
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200ResponseAdminWebServicePcvRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_admin_web_service_pcv_ref_id
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200ResponseAdminWebServicePcvRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_admin_web_service_pcv_ref_location
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200ResponsePersistentGrantContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_persistent_grant_contract_core_attributes
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200ResponsePersistentGrantContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_persistent_grant_contract_extended_attributes
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_activationCodeCheckMode"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_activation_code_check_mode
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_allowUnidentifiedClientExtensionGrants"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_allow_unidentified_client_extension_grants
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_allowUnidentifiedClientROCreds"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_allow_unidentified_client_rocreds
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_allowedOrigins"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_allowed_origins
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_approvedScopesAttribute"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_approved_scopes_attribute
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_atmIdForOAuthGrantManagement"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_atm_id_for_oauth_grant_management
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_authorizationCodeEntropy"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_authorization_code_entropy
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_authorizationCodeTimeout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_authorization_code_timeout
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_bypassActivationCodeConfirmation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_bypass_activation_code_confirmation
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_bypassAuthorizationForApprovedGrants"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_bypass_authorization_for_approved_grants
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_clientSecretRetentionPeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_client_secret_retention_period
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_defaultScopeDescription"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_default_scope_description
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_devicePollingInterval"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_device_polling_interval
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_disallowPlainPKCE"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_disallow_plain_pkce
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_exclusiveScopeGroups"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_exclusive_scope_groups
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_exclusiveScopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_exclusive_scopes
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_includeIssuerInAuthorizationResponse"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_include_issuer_in_authorization_response
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_jwtSecuredAuthorizationResponseModeLifetime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_jwt_secured_authorization_response_mode_lifetime
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_parReferenceLength"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_par_reference_length
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_parReferenceTimeout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_par_reference_timeout
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_parStatus"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_par_status
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_pendingAuthorizationTimeout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_pending_authorization_timeout
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_persistentGrantIdleTimeout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_persistent_grant_idle_timeout
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_persistentGrantIdleTimeoutTimeUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_persistent_grant_idle_timeout_time_unit
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_persistentGrantLifetime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_persistent_grant_lifetime
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_persistentGrantLifetimeUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_persistent_grant_lifetime_unit
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_persistentGrantReuseGrantTypes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_persistent_grant_reuse_grant_types
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_refreshRollingInterval"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_refresh_rolling_interval
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_refreshTokenLength"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_refresh_token_length
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_refreshTokenRollingGracePeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_refresh_token_rolling_grace_period
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_registeredAuthorizationPath"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_registered_authorization_path
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_rollRefreshTokenValues"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_roll_refresh_token_values
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_scopeForOAuthGrantManagement"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_scope_for_oauth_grant_management
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_scopeGroups"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_scope_groups
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_scopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_scopes
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_tokenEndpointBaseUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_token_endpoint_base_url
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_trackUserSessionsForLogout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_track_user_sessions_for_logout
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_userAuthorizationConsentAdapter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_user_authorization_consent_adapter
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_userAuthorizationConsentPageSetting"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_user_authorization_consent_page_setting
  }
  property {
    name  = "updateAuthorizationServerSettings_body_GetAuthorizationServerSettings200Response_userAuthorizationUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_authorization_server_settings_body_get_authorization_server_settings200_response_user_authorization_url
  }
  property {
    name  = "updateCaptchaSettings_body_GetCaptchaSettings200Response_encryptedSecretKey"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_captcha_settings_body_get_captcha_settings200_response_encrypted_secret_key
  }
  property {
    name  = "updateCaptchaSettings_body_GetCaptchaSettings200Response_secretKey"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_captcha_settings_body_get_captcha_settings200_response_secret_key
  }
  property {
    name  = "updateCaptchaSettings_body_GetCaptchaSettings200Response_siteKey"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_captcha_settings_body_get_captcha_settings200_response_site_key
  }
  property {
    name  = "updateClientSecret_body_GetClientSecret200Response_encryptedSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_secret_body_get_client_secret200_response_encrypted_secret
  }
  property {
    name  = "updateClientSecret_body_GetClientSecret200Response_secondarySecrets"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_secret_body_get_client_secret200_response_secondary_secrets
  }
  property {
    name  = "updateClientSecret_body_GetClientSecret200Response_secret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_secret_body_get_client_secret200_response_secret
  }
  property {
    name  = "updateClientSecret_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_secret_id
  }
  property {
    name  = "updateClientSettings_body_GetClient200ResponseDefaultAccessTokenManagerRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client200_response_default_access_token_manager_ref_id
  }
  property {
    name  = "updateClientSettings_body_GetClient200ResponseDefaultAccessTokenManagerRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client200_response_default_access_token_manager_ref_location
  }
  property {
    name  = "updateClientSettings_body_GetClient200ResponseOidcPolicyPolicyGroup_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client200_response_oidc_policy_policy_group_id
  }
  property {
    name  = "updateClientSettings_body_GetClient200ResponseOidcPolicyPolicyGroup_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client200_response_oidc_policy_policy_group_location
  }
  property {
    name  = "updateClientSettings_body_GetClient200ResponseRequestPolicyRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client200_response_request_policy_ref_id
  }
  property {
    name  = "updateClientSettings_body_GetClient200ResponseRequestPolicyRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client200_response_request_policy_ref_location
  }
  property {
    name  = "updateClientSettings_body_GetClient200ResponseTokenExchangeProcessorPolicyRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client200_response_token_exchange_processor_policy_ref_id
  }
  property {
    name  = "updateClientSettings_body_GetClient200ResponseTokenExchangeProcessorPolicyRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client200_response_token_exchange_processor_policy_ref_location
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistrationClientCertIssuerRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_client_cert_issuer_ref_id
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistrationClientCertIssuerRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_client_cert_issuer_ref_location
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistrationOidcPolicy_idTokenContentEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_oidc_policy_id_token_content_encryption_algorithm
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistrationOidcPolicy_idTokenEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_oidc_policy_id_token_encryption_algorithm
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistrationOidcPolicy_idTokenSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_oidc_policy_id_token_signing_algorithm
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_allowClientDelete"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_allow_client_delete
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_allowedExclusiveScopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_allowed_exclusive_scopes
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_bypassActivationCodeConfirmationOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_bypass_activation_code_confirmation_override
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_cibaPollingInterval"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_ciba_polling_interval
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_cibaRequireSignedRequests"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_ciba_require_signed_requests
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_clientCertIssuerType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_client_cert_issuer_type
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_clientSecretRetentionPeriodOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_client_secret_retention_period_override
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_clientSecretRetentionPeriodType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_client_secret_retention_period_type
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_deviceFlowSettingType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_device_flow_setting_type
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_devicePollingIntervalOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_device_polling_interval_override
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_disableRegistrationAccessTokens"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_disable_registration_access_tokens
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_enforceReplayPrevention"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_enforce_replay_prevention
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_initialAccessTokenScope"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_initial_access_token_scope
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_pendingAuthorizationTimeoutOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_pending_authorization_timeout_override
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_persistentGrantExpirationTime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_persistent_grant_expiration_time
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_persistentGrantExpirationTimeUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_persistent_grant_expiration_time_unit
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_persistentGrantExpirationType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_persistent_grant_expiration_type
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_persistentGrantIdleTimeout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_persistent_grant_idle_timeout
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_persistentGrantIdleTimeoutTimeUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_persistent_grant_idle_timeout_time_unit
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_persistentGrantIdleTimeoutType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_persistent_grant_idle_timeout_type
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_policyRefs"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_policy_refs
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_refreshRolling"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_refresh_rolling
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_refreshTokenRollingGracePeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_refresh_token_rolling_grace_period
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_refreshTokenRollingGracePeriodType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_refresh_token_rolling_grace_period_type
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_refreshTokenRollingInterval"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_refresh_token_rolling_interval
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_refreshTokenRollingIntervalType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_refresh_token_rolling_interval_type
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_requireJwtSecuredAuthorizationResponseMode"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_require_jwt_secured_authorization_response_mode
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_requireProofKeyForCodeExchange"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_require_proof_key_for_code_exchange
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_requireSignedRequests"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_require_signed_requests
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_restrictCommonScopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_restrict_common_scopes
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_restrictToDefaultAccessTokenManager"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_restrict_to_default_access_token_manager
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_restrictedCommonScopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_restricted_common_scopes
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_retainClientSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_retain_client_secret
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_rotateClientSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_rotate_client_secret
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_rotateRegistrationAccessToken"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_rotate_registration_access_token
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200ResponseDynamicClientRegistration_userAuthorizationUrlOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_dynamic_client_registration_user_authorization_url_override
  }
  property {
    name  = "updateClientSettings_body_GetClientSettings200Response_clientMetadata"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_settings_body_get_client_settings200_response_client_metadata
  }
  property {
    name  = "updateClient_body_GetClient200ResponseClientAuth_clientCertIssuerDn"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_client_auth_client_cert_issuer_dn
  }
  property {
    name  = "updateClient_body_GetClient200ResponseClientAuth_clientCertSubjectDn"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_client_auth_client_cert_subject_dn
  }
  property {
    name  = "updateClient_body_GetClient200ResponseClientAuth_encryptedSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_client_auth_encrypted_secret
  }
  property {
    name  = "updateClient_body_GetClient200ResponseClientAuth_enforceReplayPrevention"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_client_auth_enforce_replay_prevention
  }
  property {
    name  = "updateClient_body_GetClient200ResponseClientAuth_secondarySecrets"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_client_auth_secondary_secrets
  }
  property {
    name  = "updateClient_body_GetClient200ResponseClientAuth_secret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_client_auth_secret
  }
  property {
    name  = "updateClient_body_GetClient200ResponseClientAuth_tokenEndpointAuthSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_client_auth_token_endpoint_auth_signing_algorithm
  }
  property {
    name  = "updateClient_body_GetClient200ResponseClientAuth_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_client_auth_type
  }
  property {
    name  = "updateClient_body_GetClient200ResponseDefaultAccessTokenManagerRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_default_access_token_manager_ref_id
  }
  property {
    name  = "updateClient_body_GetClient200ResponseDefaultAccessTokenManagerRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_default_access_token_manager_ref_location
  }
  property {
    name  = "updateClient_body_GetClient200ResponseJwksSettings_jwks"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_jwks_settings_jwks
  }
  property {
    name  = "updateClient_body_GetClient200ResponseJwksSettings_jwksUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_jwks_settings_jwks_url
  }
  property {
    name  = "updateClient_body_GetClient200ResponseOidcPolicyPolicyGroup_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_oidc_policy_policy_group_id
  }
  property {
    name  = "updateClient_body_GetClient200ResponseOidcPolicyPolicyGroup_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_oidc_policy_policy_group_location
  }
  property {
    name  = "updateClient_body_GetClient200ResponseOidcPolicy_grantAccessSessionRevocationApi"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_oidc_policy_grant_access_session_revocation_api
  }
  property {
    name  = "updateClient_body_GetClient200ResponseOidcPolicy_grantAccessSessionSessionManagementApi"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_oidc_policy_grant_access_session_session_management_api
  }
  property {
    name  = "updateClient_body_GetClient200ResponseOidcPolicy_idTokenContentEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_oidc_policy_id_token_content_encryption_algorithm
  }
  property {
    name  = "updateClient_body_GetClient200ResponseOidcPolicy_idTokenEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_oidc_policy_id_token_encryption_algorithm
  }
  property {
    name  = "updateClient_body_GetClient200ResponseOidcPolicy_idTokenSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_oidc_policy_id_token_signing_algorithm
  }
  property {
    name  = "updateClient_body_GetClient200ResponseOidcPolicy_logoutUris"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_oidc_policy_logout_uris
  }
  property {
    name  = "updateClient_body_GetClient200ResponseOidcPolicy_pairwiseIdentifierUserType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_oidc_policy_pairwise_identifier_user_type
  }
  property {
    name  = "updateClient_body_GetClient200ResponseOidcPolicy_pingAccessLogoutCapable"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_oidc_policy_ping_access_logout_capable
  }
  property {
    name  = "updateClient_body_GetClient200ResponseOidcPolicy_sectorIdentifierUri"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_oidc_policy_sector_identifier_uri
  }
  property {
    name  = "updateClient_body_GetClient200ResponseRequestPolicyRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_request_policy_ref_id
  }
  property {
    name  = "updateClient_body_GetClient200ResponseRequestPolicyRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_request_policy_ref_location
  }
  property {
    name  = "updateClient_body_GetClient200ResponseTokenExchangeProcessorPolicyRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_token_exchange_processor_policy_ref_id
  }
  property {
    name  = "updateClient_body_GetClient200ResponseTokenExchangeProcessorPolicyRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_token_exchange_processor_policy_ref_location
  }
  property {
    name  = "updateClient_body_GetClient200Response_allowAuthenticationApiInit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_allow_authentication_api_init
  }
  property {
    name  = "updateClient_body_GetClient200Response_bypassActivationCodeConfirmationOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_bypass_activation_code_confirmation_override
  }
  property {
    name  = "updateClient_body_GetClient200Response_bypassApprovalPage"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_bypass_approval_page
  }
  property {
    name  = "updateClient_body_GetClient200Response_cibaDeliveryMode"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_ciba_delivery_mode
  }
  property {
    name  = "updateClient_body_GetClient200Response_cibaNotificationEndpoint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_ciba_notification_endpoint
  }
  property {
    name  = "updateClient_body_GetClient200Response_cibaPollingInterval"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_ciba_polling_interval
  }
  property {
    name  = "updateClient_body_GetClient200Response_cibaRequestObjectSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_ciba_request_object_signing_algorithm
  }
  property {
    name  = "updateClient_body_GetClient200Response_cibaRequireSignedRequests"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_ciba_require_signed_requests
  }
  property {
    name  = "updateClient_body_GetClient200Response_cibaUserCodeSupported"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_ciba_user_code_supported
  }
  property {
    name  = "updateClient_body_GetClient200Response_clientId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_client_id
  }
  property {
    name  = "updateClient_body_GetClient200Response_clientSecretChangedTime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_client_secret_changed_time
  }
  property {
    name  = "updateClient_body_GetClient200Response_clientSecretRetentionPeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_client_secret_retention_period
  }
  property {
    name  = "updateClient_body_GetClient200Response_clientSecretRetentionPeriodType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_client_secret_retention_period_type
  }
  property {
    name  = "updateClient_body_GetClient200Response_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_description
  }
  property {
    name  = "updateClient_body_GetClient200Response_deviceFlowSettingType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_device_flow_setting_type
  }
  property {
    name  = "updateClient_body_GetClient200Response_devicePollingIntervalOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_device_polling_interval_override
  }
  property {
    name  = "updateClient_body_GetClient200Response_enabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_enabled
  }
  property {
    name  = "updateClient_body_GetClient200Response_exclusiveScopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_exclusive_scopes
  }
  property {
    name  = "updateClient_body_GetClient200Response_extendedParameters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_extended_parameters
  }
  property {
    name  = "updateClient_body_GetClient200Response_grantTypes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_grant_types
  }
  property {
    name  = "updateClient_body_GetClient200Response_jwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_jwt_secured_authorization_response_mode_content_encryption_algorithm
  }
  property {
    name  = "updateClient_body_GetClient200Response_jwtSecuredAuthorizationResponseModeEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_jwt_secured_authorization_response_mode_encryption_algorithm
  }
  property {
    name  = "updateClient_body_GetClient200Response_jwtSecuredAuthorizationResponseModeSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_jwt_secured_authorization_response_mode_signing_algorithm
  }
  property {
    name  = "updateClient_body_GetClient200Response_logoUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_logo_url
  }
  property {
    name  = "updateClient_body_GetClient200Response_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_name
  }
  property {
    name  = "updateClient_body_GetClient200Response_pendingAuthorizationTimeoutOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_pending_authorization_timeout_override
  }
  property {
    name  = "updateClient_body_GetClient200Response_persistentGrantExpirationTime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_persistent_grant_expiration_time
  }
  property {
    name  = "updateClient_body_GetClient200Response_persistentGrantExpirationTimeUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_persistent_grant_expiration_time_unit
  }
  property {
    name  = "updateClient_body_GetClient200Response_persistentGrantExpirationType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_persistent_grant_expiration_type
  }
  property {
    name  = "updateClient_body_GetClient200Response_persistentGrantIdleTimeout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_persistent_grant_idle_timeout
  }
  property {
    name  = "updateClient_body_GetClient200Response_persistentGrantIdleTimeoutTimeUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_persistent_grant_idle_timeout_time_unit
  }
  property {
    name  = "updateClient_body_GetClient200Response_persistentGrantIdleTimeoutType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_persistent_grant_idle_timeout_type
  }
  property {
    name  = "updateClient_body_GetClient200Response_persistentGrantReuseGrantTypes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_persistent_grant_reuse_grant_types
  }
  property {
    name  = "updateClient_body_GetClient200Response_persistentGrantReuseType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_persistent_grant_reuse_type
  }
  property {
    name  = "updateClient_body_GetClient200Response_redirectUris"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_redirect_uris
  }
  property {
    name  = "updateClient_body_GetClient200Response_refreshRolling"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_refresh_rolling
  }
  property {
    name  = "updateClient_body_GetClient200Response_refreshTokenRollingGracePeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_refresh_token_rolling_grace_period
  }
  property {
    name  = "updateClient_body_GetClient200Response_refreshTokenRollingGracePeriodType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_refresh_token_rolling_grace_period_type
  }
  property {
    name  = "updateClient_body_GetClient200Response_refreshTokenRollingInterval"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_refresh_token_rolling_interval
  }
  property {
    name  = "updateClient_body_GetClient200Response_refreshTokenRollingIntervalType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_refresh_token_rolling_interval_type
  }
  property {
    name  = "updateClient_body_GetClient200Response_requestObjectSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_request_object_signing_algorithm
  }
  property {
    name  = "updateClient_body_GetClient200Response_requireJwtSecuredAuthorizationResponseMode"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_require_jwt_secured_authorization_response_mode
  }
  property {
    name  = "updateClient_body_GetClient200Response_requireProofKeyForCodeExchange"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_require_proof_key_for_code_exchange
  }
  property {
    name  = "updateClient_body_GetClient200Response_requirePushedAuthorizationRequests"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_require_pushed_authorization_requests
  }
  property {
    name  = "updateClient_body_GetClient200Response_requireSignedRequests"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_require_signed_requests
  }
  property {
    name  = "updateClient_body_GetClient200Response_restrictScopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_restrict_scopes
  }
  property {
    name  = "updateClient_body_GetClient200Response_restrictToDefaultAccessTokenManager"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_restrict_to_default_access_token_manager
  }
  property {
    name  = "updateClient_body_GetClient200Response_restrictedResponseTypes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_restricted_response_types
  }
  property {
    name  = "updateClient_body_GetClient200Response_restrictedScopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_restricted_scopes
  }
  property {
    name  = "updateClient_body_GetClient200Response_tokenIntrospectionContentEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_token_introspection_content_encryption_algorithm
  }
  property {
    name  = "updateClient_body_GetClient200Response_tokenIntrospectionEncryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_token_introspection_encryption_algorithm
  }
  property {
    name  = "updateClient_body_GetClient200Response_tokenIntrospectionSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_token_introspection_signing_algorithm
  }
  property {
    name  = "updateClient_body_GetClient200Response_userAuthorizationUrlOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_user_authorization_url_override
  }
  property {
    name  = "updateClient_body_GetClient200Response_validateUsingAllEligibleAtms"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_body_get_client200_response_validate_using_all_eligible_atms
  }
  property {
    name  = "updateClient_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_client_id
  }
  property {
    name  = "updateCommonScopeGroup_body_GetAuthorizationServerSettings200ResponseScopeGroupsInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_common_scope_group_body_get_authorization_server_settings200_response_scope_groups_inner_description
  }
  property {
    name  = "updateCommonScopeGroup_body_GetAuthorizationServerSettings200ResponseScopeGroupsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_common_scope_group_body_get_authorization_server_settings200_response_scope_groups_inner_name
  }
  property {
    name  = "updateCommonScopeGroup_body_GetAuthorizationServerSettings200ResponseScopeGroupsInner_scopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_common_scope_group_body_get_authorization_server_settings200_response_scope_groups_inner_scopes
  }
  property {
    name  = "updateCommonScopeGroup_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_common_scope_group_name
  }
  property {
    name  = "updateCommonScope_body_GetAuthorizationServerSettings200ResponseScopesInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_common_scope_body_get_authorization_server_settings200_response_scopes_inner_description
  }
  property {
    name  = "updateCommonScope_body_GetAuthorizationServerSettings200ResponseScopesInner_dynamic"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_common_scope_body_get_authorization_server_settings200_response_scopes_inner_dynamic
  }
  property {
    name  = "updateCommonScope_body_GetAuthorizationServerSettings200ResponseScopesInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_common_scope_body_get_authorization_server_settings200_response_scopes_inner_name
  }
  property {
    name  = "updateCommonScope_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_common_scope_name
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseAttributeQueryPolicy_encryptNameId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_attribute_query_policy_encrypt_name_id
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseAttributeQueryPolicy_maskAttributeValues"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_attribute_query_policy_mask_attribute_values
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseAttributeQueryPolicy_requireEncryptedAssertion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_attribute_query_policy_require_encrypted_assertion
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseAttributeQueryPolicy_requireSignedAssertion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_attribute_query_policy_require_signed_assertion
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseAttributeQueryPolicy_requireSignedResponse"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_attribute_query_policy_require_signed_response
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseAttributeQueryPolicy_signAttributeQuery"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_attribute_query_policy_sign_attribute_query
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseAttributeQuery_nameMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_attribute_query_name_mappings
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseAttributeQuery_url"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_attribute_query_url
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoArtifact_lifetime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_artifact_lifetime
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoArtifact_resolverLocations"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_artifact_resolver_locations
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoArtifact_sourceId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_artifact_source_id
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_attribute_contract_core_attributes
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_attribute_contract_extended_attributes
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoDecryptionPolicy_assertionEncrypted"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_decryption_policy_assertion_encrypted
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoDecryptionPolicy_attributesEncrypted"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_decryption_policy_attributes_encrypted
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoDecryptionPolicy_sloEncryptSubjectNameID"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_decryption_policy_slo_encrypt_subject_name_id
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoDecryptionPolicy_sloSubjectNameIDEncrypted"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_decryption_policy_slo_subject_name_idencrypted
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoDecryptionPolicy_subjectNameIdEncrypted"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_decryption_policy_subject_name_id_encrypted
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoJitProvisioningUserAttributes_attributeContract"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_jit_provisioning_user_attributes_attribute_contract
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoJitProvisioningUserAttributes_doAttributeQuery"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_jit_provisioning_user_attributes_do_attribute_query
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoJitProvisioningUserRepository_jitRepositoryAttributeMapping"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_jit_provisioning_user_repository_jit_repository_attribute_mapping
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoJitProvisioningUserRepository_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_jit_provisioning_user_repository_type
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoJitProvisioning_errorHandling"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_jit_provisioning_error_handling
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoJitProvisioning_eventTrigger"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_jit_provisioning_event_trigger
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoOauthAuthenticationPolicyContractRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_oauth_authentication_policy_contract_ref_id
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoOauthAuthenticationPolicyContractRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_oauth_authentication_policy_contract_ref_location
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_authenticationScheme"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_authentication_scheme
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_authenticationSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_authentication_signing_algorithm
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_authorizationEndpoint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_authorization_endpoint
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_enablePKCE"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_enable_pkce
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_jwksURL"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_jwks_url
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_loginType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_login_type
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_requestParameters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_request_parameters
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_requestSigningAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_request_signing_algorithm
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_scopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_scopes
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_tokenEndpoint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_token_endpoint
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoOidcProviderSettings_userInfoEndpoint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_oidc_provider_settings_user_info_endpoint
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoSsoOAuthMapping_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_sso_oauth_mapping_attribute_contract_fulfillment
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSsoSsoOAuthMapping_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_sso_oauth_mapping_attribute_sources
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_adapterMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_adapter_mappings
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_alwaysSignArtifactResponse"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_always_sign_artifact_response
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_assertionsSigned"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_assertions_signed
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_authenticationPolicyContractMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_authentication_policy_contract_mappings
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_authnContextMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_authn_context_mappings
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_defaultTargetUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_default_target_url
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_enabledProfiles"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_enabled_profiles
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_idpIdentityMapping"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_idp_identity_mapping
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_incomingBindings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_incoming_bindings
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_messageCustomizations"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_message_customizations
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_protocol"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_protocol
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_signAuthnRequests"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_sign_authn_requests
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_sloServiceEndpoints"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_slo_service_endpoints
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_ssoServiceEndpoints"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_sso_service_endpoints
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpBrowserSso_urlWhitelistEntries"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_browser_sso_url_whitelist_entries
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpOAuthGrantAttributeMappingIdpOAuthAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_oauth_grant_attribute_mapping_idp_oauth_attribute_contract_core_attributes
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpOAuthGrantAttributeMappingIdpOAuthAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_oauth_grant_attribute_mapping_idp_oauth_attribute_contract_extended_attributes
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseIdpOAuthGrantAttributeMapping_accessTokenManagerMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_idp_oauth_grant_attribute_mapping_access_token_manager_mappings
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseInboundProvisioningCustomSchema_attributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_inbound_provisioning_custom_schema_attributes
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseInboundProvisioningCustomSchema_namespace"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_inbound_provisioning_custom_schema_namespace
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseInboundProvisioningGroupsReadGroups_attributeFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_inbound_provisioning_groups_read_groups_attribute_fulfillment
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseInboundProvisioningGroupsReadGroups_attributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_inbound_provisioning_groups_read_groups_attributes
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseInboundProvisioningGroupsWriteGroups_attributeFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_inbound_provisioning_groups_write_groups_attribute_fulfillment
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseInboundProvisioningUserRepository_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_inbound_provisioning_user_repository_type
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseInboundProvisioningUsersReadUsersAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_inbound_provisioning_users_read_users_attribute_contract_core_attributes
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseInboundProvisioningUsersReadUsersAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_inbound_provisioning_users_read_users_attribute_contract_extended_attributes
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseInboundProvisioningUsersReadUsers_attributeFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_inbound_provisioning_users_read_users_attribute_fulfillment
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseInboundProvisioningUsersReadUsers_attributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_inbound_provisioning_users_read_users_attributes
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseInboundProvisioningUsersWriteUsers_attributeFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_inbound_provisioning_users_write_users_attribute_fulfillment
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseInboundProvisioning_actionOnDelete"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_inbound_provisioning_action_on_delete
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseInboundProvisioning_groupSupport"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_inbound_provisioning_group_support
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseOidcClientCredentials_clientId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_oidc_client_credentials_client_id
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseOidcClientCredentials_clientSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_oidc_client_credentials_client_secret
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseOidcClientCredentials_encryptedSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_oidc_client_credentials_encrypted_secret
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseWsTrustAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_ws_trust_attribute_contract_core_attributes
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseWsTrustAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_ws_trust_attribute_contract_extended_attributes
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseWsTrust_generateLocalToken"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_ws_trust_generate_local_token
  }
  property {
    name  = "updateConnection1_body_GetConnection1200ResponseWsTrust_tokenGeneratorMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_connection1200_response_ws_trust_token_generator_mappings
  }
  property {
    name  = "updateConnection1_body_GetMappings200ResponseItemsInnerAttributeSourcesInnerDataStoreRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_mappings200_response_items_inner_attribute_sources_inner_data_store_ref_id
  }
  property {
    name  = "updateConnection1_body_GetMappings200ResponseItemsInnerAttributeSourcesInnerDataStoreRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_mappings200_response_items_inner_attribute_sources_inner_data_store_ref_location
  }
  property {
    name  = "updateConnection1_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "updateConnection1_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "updateConnection1_body_UpdateConnection1Request_errorPageMsgId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_body_update_connection1_request_error_page_msg_id
  }
  property {
    name  = "updateConnection1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_id
  }
  property {
    name  = "updateConnection1_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection1_x_bypass_external_validation
  }
  property {
    name  = "updateConnection2_body_GetConnection1200ResponseIdpBrowserSsoArtifact_lifetime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection1200_response_idp_browser_sso_artifact_lifetime
  }
  property {
    name  = "updateConnection2_body_GetConnection1200ResponseIdpBrowserSsoArtifact_resolverLocations"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection1200_response_idp_browser_sso_artifact_resolver_locations
  }
  property {
    name  = "updateConnection2_body_GetConnection1200ResponseIdpBrowserSsoArtifact_sourceId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection1200_response_idp_browser_sso_artifact_source_id
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseAttributeQueryPolicy_encryptAssertion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_attribute_query_policy_encrypt_assertion
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseAttributeQueryPolicy_requireEncryptedNameId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_attribute_query_policy_require_encrypted_name_id
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseAttributeQueryPolicy_requireSignedAttributeQuery"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_attribute_query_policy_require_signed_attribute_query
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseAttributeQueryPolicy_signAssertion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_attribute_query_policy_sign_assertion
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseAttributeQueryPolicy_signResponse"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_attribute_query_policy_sign_response
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseAttributeQuery_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_attribute_query_attribute_contract_fulfillment
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseAttributeQuery_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_attribute_query_attribute_sources
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseAttributeQuery_attributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_attribute_query_attributes
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseOutboundProvisionCustomSchema_attributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_outbound_provision_custom_schema_attributes
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseOutboundProvisionCustomSchema_namespace"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_outbound_provision_custom_schema_namespace
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseOutboundProvision_channels"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_outbound_provision_channels
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseOutboundProvision_targetSettings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_outbound_provision_target_settings
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseOutboundProvision_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_outbound_provision_type
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSsoAssertionLifetime_minutesAfter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_assertion_lifetime_minutes_after
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSsoAssertionLifetime_minutesBefore"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_assertion_lifetime_minutes_before
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSsoAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_attribute_contract_core_attributes
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSsoAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_attribute_contract_extended_attributes
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSsoEncryptionPolicy_encryptAssertion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_encryption_policy_encrypt_assertion
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSsoEncryptionPolicy_encryptSloSubjectNameId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_encryption_policy_encrypt_slo_subject_name_id
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSsoEncryptionPolicy_encryptedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_encryption_policy_encrypted_attributes
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSsoEncryptionPolicy_sloSubjectNameIDEncrypted"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_encryption_policy_slo_subject_name_idencrypted
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_adapterMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_adapter_mappings
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_alwaysSignArtifactResponse"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_always_sign_artifact_response
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_authenticationPolicyContractAssertionMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_authentication_policy_contract_assertion_mappings
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_defaultTargetUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_default_target_url
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_enabledProfiles"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_enabled_profiles
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_incomingBindings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_incoming_bindings
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_messageCustomizations"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_message_customizations
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_protocol"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_protocol
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_requireSignedAuthnRequests"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_require_signed_authn_requests
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_signAssertions"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_sign_assertions
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_signResponseAsRequired"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_sign_response_as_required
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_sloServiceEndpoints"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_slo_service_endpoints
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_spSamlIdentityMapping"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_sp_saml_identity_mapping
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_spWsFedIdentityMapping"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_sp_ws_fed_identity_mapping
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_ssoServiceEndpoints"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_sso_service_endpoints
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_urlWhitelistEntries"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_url_whitelist_entries
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_wsFedTokenType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_ws_fed_token_type
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseSpBrowserSso_wsTrustVersion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_sp_browser_sso_ws_trust_version
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrustAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_attribute_contract_core_attributes
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrustAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_attribute_contract_extended_attributes
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrustRequestContractRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_request_contract_ref_id
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrustRequestContractRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_request_contract_ref_location
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrust_abortIfNotFulfilledFromRequest"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_abort_if_not_fulfilled_from_request
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrust_defaultTokenType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_default_token_type
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrust_encryptSaml2Assertion"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_encrypt_saml2_assertion
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrust_generateKey"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_generate_key
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrust_messageCustomizations"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_message_customizations
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrust_minutesAfter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_minutes_after
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrust_minutesBefore"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_minutes_before
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrust_oAuthAssertionProfiles"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_o_auth_assertion_profiles
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrust_partnerServiceIds"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_partner_service_ids
  }
  property {
    name  = "updateConnection2_body_GetConnection2200ResponseWsTrust_tokenProcessorMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_connection2200_response_ws_trust_token_processor_mappings
  }
  property {
    name  = "updateConnection2_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "updateConnection2_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "updateConnection2_body_UpdateConnection2Request_applicationIconUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_update_connection2_request_application_icon_url
  }
  property {
    name  = "updateConnection2_body_UpdateConnection2Request_applicationName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_update_connection2_request_application_name
  }
  property {
    name  = "updateConnection2_body_UpdateConnection2Request_connectionTargetType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_body_update_connection2_request_connection_target_type
  }
  property {
    name  = "updateConnection2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_id
  }
  property {
    name  = "updateConnection2_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection2_x_bypass_external_validation
  }
  property {
    name  = "updateConnectionCerts1_body_GetConnectionCerts1200Response_items"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection_certs1_body_get_connection_certs1200_response_items
  }
  property {
    name  = "updateConnectionCerts1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection_certs1_id
  }
  property {
    name  = "updateConnectionCerts2_body_GetConnectionCerts1200Response_items"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection_certs2_body_get_connection_certs1200_response_items
  }
  property {
    name  = "updateConnectionCerts2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_connection_certs2_id
  }
  property {
    name  = "updateDataStore_body_CreateDataStoreRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_data_store_body_create_data_store_request_id
  }
  property {
    name  = "updateDataStore_body_CreateDataStoreRequest_maskAttributeValues"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_data_store_body_create_data_store_request_mask_attribute_values
  }
  property {
    name  = "updateDataStore_body_CreateDataStoreRequest_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_data_store_body_create_data_store_request_type
  }
  property {
    name  = "updateDataStore_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_data_store_id
  }
  property {
    name  = "updateDataStore_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_data_store_x_bypass_external_validation
  }
  property {
    name  = "updateDecryptionKeys1_body_GetDecryptionKeys1200ResponsePrimaryKeyRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_decryption_keys1_body_get_decryption_keys1200_response_primary_key_ref_id
  }
  property {
    name  = "updateDecryptionKeys1_body_GetDecryptionKeys1200ResponsePrimaryKeyRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_decryption_keys1_body_get_decryption_keys1200_response_primary_key_ref_location
  }
  property {
    name  = "updateDecryptionKeys1_body_GetDecryptionKeys1200ResponseSecondaryKeyPairRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_decryption_keys1_body_get_decryption_keys1200_response_secondary_key_pair_ref_id
  }
  property {
    name  = "updateDecryptionKeys1_body_GetDecryptionKeys1200ResponseSecondaryKeyPairRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_decryption_keys1_body_get_decryption_keys1200_response_secondary_key_pair_ref_location
  }
  property {
    name  = "updateDecryptionKeys1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_decryption_keys1_id
  }
  property {
    name  = "updateDecryptionKeys2_body_GetDecryptionKeys1200ResponsePrimaryKeyRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_decryption_keys2_body_get_decryption_keys1200_response_primary_key_ref_id
  }
  property {
    name  = "updateDecryptionKeys2_body_GetDecryptionKeys1200ResponsePrimaryKeyRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_decryption_keys2_body_get_decryption_keys1200_response_primary_key_ref_location
  }
  property {
    name  = "updateDecryptionKeys2_body_GetDecryptionKeys1200ResponseSecondaryKeyPairRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_decryption_keys2_body_get_decryption_keys1200_response_secondary_key_pair_ref_id
  }
  property {
    name  = "updateDecryptionKeys2_body_GetDecryptionKeys1200ResponseSecondaryKeyPairRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_decryption_keys2_body_get_decryption_keys1200_response_secondary_key_pair_ref_location
  }
  property {
    name  = "updateDecryptionKeys2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_decryption_keys2_id
  }
  property {
    name  = "updateDefaultAuthenticationPolicy_body_GetDefaultAuthenticationPolicy200Response_authnSelectionTrees"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_default_authentication_policy_body_get_default_authentication_policy200_response_authn_selection_trees
  }
  property {
    name  = "updateDefaultAuthenticationPolicy_body_GetDefaultAuthenticationPolicy200Response_defaultAuthenticationSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_default_authentication_policy_body_get_default_authentication_policy200_response_default_authentication_sources
  }
  property {
    name  = "updateDefaultAuthenticationPolicy_body_GetDefaultAuthenticationPolicy200Response_failIfNoSelection"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_default_authentication_policy_body_get_default_authentication_policy200_response_fail_if_no_selection
  }
  property {
    name  = "updateDefaultAuthenticationPolicy_body_GetDefaultAuthenticationPolicy200Response_trackedHttpParameters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_default_authentication_policy_body_get_default_authentication_policy200_response_tracked_http_parameters
  }
  property {
    name  = "updateDefaultAuthenticationPolicy_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_default_authentication_policy_x_bypass_external_validation
  }
  property {
    name  = "updateDefaultUrlSettings_body_GetDefaultUrl200Response_confirmIdpSlo"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_default_url_settings_body_get_default_url200_response_confirm_idp_slo
  }
  property {
    name  = "updateDefaultUrlSettings_body_GetDefaultUrl200Response_idpErrorMsg"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_default_url_settings_body_get_default_url200_response_idp_error_msg
  }
  property {
    name  = "updateDefaultUrlSettings_body_GetDefaultUrl200Response_idpSloSuccessUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_default_url_settings_body_get_default_url200_response_idp_slo_success_url
  }
  property {
    name  = "updateDefaultUrls_body_GetDefaultUrls200Response_confirmSlo"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_default_urls_body_get_default_urls200_response_confirm_slo
  }
  property {
    name  = "updateDefaultUrls_body_GetDefaultUrls200Response_sloSuccessUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_default_urls_body_get_default_urls200_response_slo_success_url
  }
  property {
    name  = "updateDefaultUrls_body_GetDefaultUrls200Response_ssoSuccessUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_default_urls_body_get_default_urls200_response_sso_success_url
  }
  property {
    name  = "updateDynamicClientRegistrationPolicy_body_CreateDynamicClientRegistrationPolicyRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_dynamic_client_registration_policy_body_create_dynamic_client_registration_policy_request_id
  }
  property {
    name  = "updateDynamicClientRegistrationPolicy_body_CreateDynamicClientRegistrationPolicyRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_dynamic_client_registration_policy_body_create_dynamic_client_registration_policy_request_name
  }
  property {
    name  = "updateDynamicClientRegistrationPolicy_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_dynamic_client_registration_policy_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "updateDynamicClientRegistrationPolicy_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_dynamic_client_registration_policy_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "updateDynamicClientRegistrationPolicy_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_dynamic_client_registration_policy_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "updateDynamicClientRegistrationPolicy_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_dynamic_client_registration_policy_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "updateDynamicClientRegistrationPolicy_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_dynamic_client_registration_policy_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "updateDynamicClientRegistrationPolicy_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_dynamic_client_registration_policy_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "updateDynamicClientRegistrationPolicy_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_dynamic_client_registration_policy_id
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_emailServer"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_email_server
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_enableUtf8MessageHeaders"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_enable_utf8_message_headers
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_encryptedPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_encrypted_password
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_password
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_port"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_port
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_retryAttempts"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_retry_attempts
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_retryDelay"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_retry_delay
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_sourceAddr"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_source_addr
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_sslPort"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_ssl_port
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_timeout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_timeout
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_useDebugging"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_use_debugging
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_useSSL"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_use_ssl
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_useTLS"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_use_tls
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_username"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_username
  }
  property {
    name  = "updateEmailServerSettings_body_GetEmailServerSettings200Response_verifyHostname"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_body_get_email_server_settings200_response_verify_hostname
  }
  property {
    name  = "updateEmailServerSettings_validate_only"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_validate_only
  }
  property {
    name  = "updateEmailServerSettings_validation_email"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_email_server_settings_validation_email
  }
  property {
    name  = "updateExclusiveScopeGroups_body_GetAuthorizationServerSettings200ResponseScopeGroupsInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_exclusive_scope_groups_body_get_authorization_server_settings200_response_scope_groups_inner_description
  }
  property {
    name  = "updateExclusiveScopeGroups_body_GetAuthorizationServerSettings200ResponseScopeGroupsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_exclusive_scope_groups_body_get_authorization_server_settings200_response_scope_groups_inner_name
  }
  property {
    name  = "updateExclusiveScopeGroups_body_GetAuthorizationServerSettings200ResponseScopeGroupsInner_scopes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_exclusive_scope_groups_body_get_authorization_server_settings200_response_scope_groups_inner_scopes
  }
  property {
    name  = "updateExclusiveScopeGroups_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_exclusive_scope_groups_name
  }
  property {
    name  = "updateExclusiveScope_body_GetAuthorizationServerSettings200ResponseScopesInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_exclusive_scope_body_get_authorization_server_settings200_response_scopes_inner_description
  }
  property {
    name  = "updateExclusiveScope_body_GetAuthorizationServerSettings200ResponseScopesInner_dynamic"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_exclusive_scope_body_get_authorization_server_settings200_response_scopes_inner_dynamic
  }
  property {
    name  = "updateExclusiveScope_body_GetAuthorizationServerSettings200ResponseScopesInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_exclusive_scope_body_get_authorization_server_settings200_response_scopes_inner_name
  }
  property {
    name  = "updateExclusiveScope_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_exclusive_scope_name
  }
  property {
    name  = "updateExtendedProperties_body_GetExtendedProperties200Response_items"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_extended_properties_body_get_extended_properties200_response_items
  }
  property {
    name  = "updateFragment_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerRootNodeAction_context"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_fragment_body_get_default_authentication_policy200_response_authn_selection_trees_inner_root_node_action_context
  }
  property {
    name  = "updateFragment_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerRootNodeAction_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_fragment_body_get_default_authentication_policy200_response_authn_selection_trees_inner_root_node_action_type
  }
  property {
    name  = "updateFragment_body_GetFragment200ResponseInputs_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_fragment_body_get_fragment200_response_inputs_id
  }
  property {
    name  = "updateFragment_body_GetFragment200ResponseInputs_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_fragment_body_get_fragment200_response_inputs_location
  }
  property {
    name  = "updateFragment_body_GetFragment200ResponseOutputs_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_fragment_body_get_fragment200_response_outputs_id
  }
  property {
    name  = "updateFragment_body_GetFragment200ResponseOutputs_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_fragment_body_get_fragment200_response_outputs_location
  }
  property {
    name  = "updateFragment_body_GetFragment200ResponseRootNode_children"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_fragment_body_get_fragment200_response_root_node_children
  }
  property {
    name  = "updateFragment_body_GetFragment200Response_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_fragment_body_get_fragment200_response_description
  }
  property {
    name  = "updateFragment_body_GetFragment200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_fragment_body_get_fragment200_response_id
  }
  property {
    name  = "updateFragment_body_GetFragment200Response_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_fragment_body_get_fragment200_response_name
  }
  property {
    name  = "updateFragment_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_fragment_id
  }
  property {
    name  = "updateFragment_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_fragment_x_bypass_external_validation
  }
  property {
    name  = "updateGeneralSettings_body_GetGeneralSettings200Response_datastoreValidationIntervalSecs"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_general_settings_body_get_general_settings200_response_datastore_validation_interval_secs
  }
  property {
    name  = "updateGeneralSettings_body_GetGeneralSettings200Response_disableAutomaticConnectionValidation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_general_settings_body_get_general_settings200_response_disable_automatic_connection_validation
  }
  property {
    name  = "updateGeneralSettings_body_GetGeneralSettings200Response_idpConnectionTransactionLoggingOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_general_settings_body_get_general_settings200_response_idp_connection_transaction_logging_override
  }
  property {
    name  = "updateGeneralSettings_body_GetGeneralSettings200Response_requestHeaderForCorrelationId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_general_settings_body_get_general_settings200_response_request_header_for_correlation_id
  }
  property {
    name  = "updateGeneralSettings_body_GetGeneralSettings200Response_spConnectionTransactionLoggingOverride"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_general_settings_body_get_general_settings200_response_sp_connection_transaction_logging_override
  }
  property {
    name  = "updateGlobalPolicy_body_GetGlobalPolicy200Response_enableSessions"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_global_policy_body_get_global_policy200_response_enable_sessions
  }
  property {
    name  = "updateGlobalPolicy_body_GetGlobalPolicy200Response_hashUniqueUserKeyAttribute"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_global_policy_body_get_global_policy200_response_hash_unique_user_key_attribute
  }
  property {
    name  = "updateGlobalPolicy_body_GetGlobalPolicy200Response_idleTimeoutDisplayUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_global_policy_body_get_global_policy200_response_idle_timeout_display_unit
  }
  property {
    name  = "updateGlobalPolicy_body_GetGlobalPolicy200Response_idleTimeoutMins"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_global_policy_body_get_global_policy200_response_idle_timeout_mins
  }
  property {
    name  = "updateGlobalPolicy_body_GetGlobalPolicy200Response_maxTimeoutDisplayUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_global_policy_body_get_global_policy200_response_max_timeout_display_unit
  }
  property {
    name  = "updateGlobalPolicy_body_GetGlobalPolicy200Response_maxTimeoutMins"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_global_policy_body_get_global_policy200_response_max_timeout_mins
  }
  property {
    name  = "updateGlobalPolicy_body_GetGlobalPolicy200Response_persistentSessions"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_global_policy_body_get_global_policy200_response_persistent_sessions
  }
  property {
    name  = "updateGroup_body_GetGroup200Response_generatorMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_group_body_get_group200_response_generator_mappings
  }
  property {
    name  = "updateGroup_body_GetGroup200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_group_body_get_group200_response_id
  }
  property {
    name  = "updateGroup_body_GetGroup200Response_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_group_body_get_group200_response_name
  }
  property {
    name  = "updateGroup_body_GetGroup200Response_resourceUris"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_group_body_get_group200_response_resource_uris
  }
  property {
    name  = "updateGroup_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_group_bypass_external_validation
  }
  property {
    name  = "updateGroup_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_group_id
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerApcId_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_apc_id_id
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerApcId_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_apc_id_location
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerAuthSourceUpdatePolicy_retainAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_auth_source_update_policy_retain_attributes
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerAuthSourceUpdatePolicy_storeAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_auth_source_update_policy_store_attributes
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerAuthSourceUpdatePolicy_updateAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_auth_source_update_policy_update_attributes
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerAuthSourceUpdatePolicy_updateInterval"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_auth_source_update_policy_update_interval
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerDataStoreConfig_dataStoreMapping"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_data_store_config_data_store_mapping
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerDataStoreConfig_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_data_store_config_type
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfigNotificationPublisherRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_notification_publisher_ref_id
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfigNotificationPublisherRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_notification_publisher_ref_location
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_allowedOtpCharacterSet"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_allowed_otp_character_set
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_emailVerificationEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_email_verification_enabled
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_emailVerificationErrorTemplateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_email_verification_error_template_name
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_emailVerificationOtpTemplateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_email_verification_otp_template_name
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_emailVerificationSentTemplateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_email_verification_sent_template_name
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_emailVerificationSuccessTemplateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_email_verification_success_template_name
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_emailVerificationType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_email_verification_type
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_fieldForEmailToVerify"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_field_for_email_to_verify
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_fieldStoringVerificationStatus"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_field_storing_verification_status
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_otlTimeToLive"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_otl_time_to_live
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_otpLength"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_otp_length
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_otpRetryAttempts"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_otp_retry_attempts
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_otpTimeToLive"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_otp_time_to_live
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_requireVerifiedEmail"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_require_verified_email
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_requireVerifiedEmailTemplateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_require_verified_email_template_name
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfig_verifyEmailTemplateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_email_verification_config_verify_email_template_name
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerFieldConfig_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_field_config_fields
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerFieldConfig_stripSpaceFromUniqueField"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_field_config_strip_space_from_unique_field
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerProfileConfig_deleteIdentityEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_profile_config_delete_identity_enabled
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerProfileConfig_templateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_profile_config_template_name
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfigRegistrationWorkflow_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_registration_workflow_id
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfigRegistrationWorkflow_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_registration_workflow_location
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfig_captchaEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_captcha_enabled
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfig_createAuthnSessionAfterRegistration"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_create_authn_session_after_registration
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfig_executeWorkflow"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_execute_workflow
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfig_templateName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_template_name
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfig_thisIsMyDeviceEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_this_is_my_device_enabled
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInnerRegistrationConfig_usernameField"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_registration_config_username_field
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInner_authSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_auth_sources
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_id
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_name
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInner_profileEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_profile_enabled
  }
  property {
    name  = "updateIdentityProfile_body_GetIdentityProfiles200ResponseItemsInner_registrationEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_identity_profiles200_response_items_inner_registration_enabled
  }
  property {
    name  = "updateIdentityProfile_body_GetMappings200ResponseItemsInnerAttributeSourcesInnerDataStoreRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_mappings200_response_items_inner_attribute_sources_inner_data_store_ref_id
  }
  property {
    name  = "updateIdentityProfile_body_GetMappings200ResponseItemsInnerAttributeSourcesInnerDataStoreRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_body_get_mappings200_response_items_inner_attribute_sources_inner_data_store_ref_location
  }
  property {
    name  = "updateIdentityProfile_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_id
  }
  property {
    name  = "updateIdentityProfile_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_profile_x_bypass_external_validation
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInnerAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_attribute_contract_core_attributes
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInnerAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_attribute_contract_extended_attributes
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInnerAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_attribute_contract_inherited
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInnerGroupAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_group_attribute_contract_core_attributes
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInnerGroupAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_group_attribute_contract_extended_attributes
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInnerGroupAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_group_attribute_contract_inherited
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_id
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetIdentityStoreProvisioners200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_identity_store_provisioners200_response_items_inner_name
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "updateIdentityStoreProvisioner_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "updateIdentityStoreProvisioner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_identity_store_provisioner_id
  }
  property {
    name  = "updateIdpAdapterMapping_body_GetIdpAdapterMappings200ResponseItemsInnerIdpAdapterRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_mapping_body_get_idp_adapter_mappings200_response_items_inner_idp_adapter_ref_id
  }
  property {
    name  = "updateIdpAdapterMapping_body_GetIdpAdapterMappings200ResponseItemsInnerIdpAdapterRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_mapping_body_get_idp_adapter_mappings200_response_items_inner_idp_adapter_ref_location
  }
  property {
    name  = "updateIdpAdapterMapping_body_GetIdpAdapterMappings200ResponseItemsInner_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_mapping_body_get_idp_adapter_mappings200_response_items_inner_attribute_contract_fulfillment
  }
  property {
    name  = "updateIdpAdapterMapping_body_GetIdpAdapterMappings200ResponseItemsInner_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_mapping_body_get_idp_adapter_mappings200_response_items_inner_attribute_sources
  }
  property {
    name  = "updateIdpAdapterMapping_body_GetIdpAdapterMappings200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_mapping_body_get_idp_adapter_mappings200_response_items_inner_id
  }
  property {
    name  = "updateIdpAdapterMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_mapping_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "updateIdpAdapterMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_mapping_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "updateIdpAdapterMapping_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_mapping_id
  }
  property {
    name  = "updateIdpAdapterMapping_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_mapping_x_bypass_external_validation
  }
  property {
    name  = "updateIdpAdapter_body_CreateIdpAdapterRequest_authnCtxClassRef"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_create_idp_adapter_request_authn_ctx_class_ref
  }
  property {
    name  = "updateIdpAdapter_body_CreateIdpAdapterRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_create_idp_adapter_request_id
  }
  property {
    name  = "updateIdpAdapter_body_CreateIdpAdapterRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_create_idp_adapter_request_name
  }
  property {
    name  = "updateIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_contract_core_attributes
  }
  property {
    name  = "updateIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_contract_extended_attributes
  }
  property {
    name  = "updateIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_contract_inherited
  }
  property {
    name  = "updateIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeContract_maskOgnlValues"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_contract_mask_ognl_values
  }
  property {
    name  = "updateIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeContract_uniqueUserKeyAttribute"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_contract_unique_user_key_attribute
  }
  property {
    name  = "updateIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeMapping_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_mapping_attribute_contract_fulfillment
  }
  property {
    name  = "updateIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeMapping_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_mapping_attribute_sources
  }
  property {
    name  = "updateIdpAdapter_body_GetIdpAdapters200ResponseItemsInnerAllOfAttributeMapping_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_idp_adapters200_response_items_inner_all_of_attribute_mapping_inherited
  }
  property {
    name  = "updateIdpAdapter_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "updateIdpAdapter_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "updateIdpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "updateIdpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "updateIdpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "updateIdpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "updateIdpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "updateIdpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "updateIdpAdapter_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_id
  }
  property {
    name  = "updateIdpAdapter_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_adapter_x_bypass_external_validation
  }
  property {
    name  = "updateIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_applicationIconUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_application_icon_url
  }
  property {
    name  = "updateIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_applicationName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_application_name
  }
  property {
    name  = "updateIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_attribute_contract_fulfillment
  }
  property {
    name  = "updateIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_attribute_sources
  }
  property {
    name  = "updateIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_defaultTargetResource"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_default_target_resource
  }
  property {
    name  = "updateIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_id
  }
  property {
    name  = "updateIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_licenseConnectionGroupAssignment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_license_connection_group_assignment
  }
  property {
    name  = "updateIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_sourceId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_source_id
  }
  property {
    name  = "updateIdpToSpAdapterMapping_body_GetIdpToSpAdapterMappingsById200Response_targetId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_to_sp_adapter_mapping_body_get_idp_to_sp_adapter_mappings_by_id200_response_target_id
  }
  property {
    name  = "updateIdpToSpAdapterMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_to_sp_adapter_mapping_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "updateIdpToSpAdapterMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_to_sp_adapter_mapping_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "updateIdpToSpAdapterMapping_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_to_sp_adapter_mapping_id
  }
  property {
    name  = "updateIdpToSpAdapterMapping_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_idp_to_sp_adapter_mapping_x_bypass_external_validation
  }
  property {
    name  = "updateIncomingProxySettings_body_GetIncomingProxySettings200Response_clientCertChainSSLHeaderName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_incoming_proxy_settings_body_get_incoming_proxy_settings200_response_client_cert_chain_sslheader_name
  }
  property {
    name  = "updateIncomingProxySettings_body_GetIncomingProxySettings200Response_clientCertSSLHeaderName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_incoming_proxy_settings_body_get_incoming_proxy_settings200_response_client_cert_sslheader_name
  }
  property {
    name  = "updateIncomingProxySettings_body_GetIncomingProxySettings200Response_forwardedHostHeaderIndex"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_incoming_proxy_settings_body_get_incoming_proxy_settings200_response_forwarded_host_header_index
  }
  property {
    name  = "updateIncomingProxySettings_body_GetIncomingProxySettings200Response_forwardedHostHeaderName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_incoming_proxy_settings_body_get_incoming_proxy_settings200_response_forwarded_host_header_name
  }
  property {
    name  = "updateIncomingProxySettings_body_GetIncomingProxySettings200Response_forwardedIpAddressHeaderIndex"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_incoming_proxy_settings_body_get_incoming_proxy_settings200_response_forwarded_ip_address_header_index
  }
  property {
    name  = "updateIncomingProxySettings_body_GetIncomingProxySettings200Response_forwardedIpAddressHeaderName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_incoming_proxy_settings_body_get_incoming_proxy_settings200_response_forwarded_ip_address_header_name
  }
  property {
    name  = "updateIncomingProxySettings_body_GetIncomingProxySettings200Response_proxyTerminatesHttpsConns"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_incoming_proxy_settings_body_get_incoming_proxy_settings200_response_proxy_terminates_https_conns
  }
  property {
    name  = "updateIssuer_body_GetIssuerById200Response_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_issuer_body_get_issuer_by_id200_response_description
  }
  property {
    name  = "updateIssuer_body_GetIssuerById200Response_host"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_issuer_body_get_issuer_by_id200_response_host
  }
  property {
    name  = "updateIssuer_body_GetIssuerById200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_issuer_body_get_issuer_by_id200_response_id
  }
  property {
    name  = "updateIssuer_body_GetIssuerById200Response_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_issuer_body_get_issuer_by_id200_response_name
  }
  property {
    name  = "updateIssuer_body_GetIssuerById200Response_path"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_issuer_body_get_issuer_by_id200_response_path
  }
  property {
    name  = "updateIssuer_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_issuer_id
  }
  property {
    name  = "updateKerberosRealm_body_GetKerberosRealms200ResponseItemsInnerLdapGatewayDataStoreRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_body_get_kerberos_realms200_response_items_inner_ldap_gateway_data_store_ref_id
  }
  property {
    name  = "updateKerberosRealm_body_GetKerberosRealms200ResponseItemsInnerLdapGatewayDataStoreRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_body_get_kerberos_realms200_response_items_inner_ldap_gateway_data_store_ref_location
  }
  property {
    name  = "updateKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_connectionType"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_body_get_kerberos_realms200_response_items_inner_connection_type
  }
  property {
    name  = "updateKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_body_get_kerberos_realms200_response_items_inner_id
  }
  property {
    name  = "updateKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_kerberosEncryptedPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_body_get_kerberos_realms200_response_items_inner_kerberos_encrypted_password
  }
  property {
    name  = "updateKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_kerberosPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_body_get_kerberos_realms200_response_items_inner_kerberos_password
  }
  property {
    name  = "updateKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_kerberosRealmName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_body_get_kerberos_realms200_response_items_inner_kerberos_realm_name
  }
  property {
    name  = "updateKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_kerberosUsername"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_body_get_kerberos_realms200_response_items_inner_kerberos_username
  }
  property {
    name  = "updateKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_keyDistributionCenters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_body_get_kerberos_realms200_response_items_inner_key_distribution_centers
  }
  property {
    name  = "updateKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_keySets"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_body_get_kerberos_realms200_response_items_inner_key_sets
  }
  property {
    name  = "updateKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_retainPreviousKeysOnPasswordChange"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_body_get_kerberos_realms200_response_items_inner_retain_previous_keys_on_password_change
  }
  property {
    name  = "updateKerberosRealm_body_GetKerberosRealms200ResponseItemsInner_suppressDomainNameConcatenation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_body_get_kerberos_realms200_response_items_inner_suppress_domain_name_concatenation
  }
  property {
    name  = "updateKerberosRealm_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_id
  }
  property {
    name  = "updateKerberosRealm_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_kerberos_realm_x_bypass_external_validation
  }
  property {
    name  = "updateKeySet_body_GetKeySet200ResponseSigningKeys_p256PublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_key_set200_response_signing_keys_p256_publish_x5c_parameter
  }
  property {
    name  = "updateKeySet_body_GetKeySet200ResponseSigningKeys_p384PublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_key_set200_response_signing_keys_p384_publish_x5c_parameter
  }
  property {
    name  = "updateKeySet_body_GetKeySet200ResponseSigningKeys_p521PublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_key_set200_response_signing_keys_p521_publish_x5c_parameter
  }
  property {
    name  = "updateKeySet_body_GetKeySet200ResponseSigningKeys_rsaPublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_key_set200_response_signing_keys_rsa_publish_x5c_parameter
  }
  property {
    name  = "updateKeySet_body_GetKeySet200Response_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_key_set200_response_description
  }
  property {
    name  = "updateKeySet_body_GetKeySet200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_key_set200_response_id
  }
  property {
    name  = "updateKeySet_body_GetKeySet200Response_issuers"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_key_set200_response_issuers
  }
  property {
    name  = "updateKeySet_body_GetKeySet200Response_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_key_set200_response_name
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseP256ActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_p256_active_cert_ref_id
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseP256ActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_p256_active_cert_ref_location
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseP256PreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_p256_previous_cert_ref_id
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseP256PreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_p256_previous_cert_ref_location
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseP384ActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_p384_active_cert_ref_id
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseP384ActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_p384_active_cert_ref_location
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseP384PreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_p384_previous_cert_ref_id
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseP384PreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_p384_previous_cert_ref_location
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseP521ActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_p521_active_cert_ref_id
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseP521ActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_p521_active_cert_ref_location
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseP521PreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_p521_previous_cert_ref_id
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseP521PreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_p521_previous_cert_ref_location
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseRsaActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_rsa_active_cert_ref_id
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseRsaActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_rsa_active_cert_ref_location
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseRsaPreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_rsa_previous_cert_ref_id
  }
  property {
    name  = "updateKeySet_body_GetOauthOidcKeysSettings200ResponseRsaPreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_body_get_oauth_oidc_keys_settings200_response_rsa_previous_cert_ref_location
  }
  property {
    name  = "updateKeySet_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_key_set_id
  }
  property {
    name  = "updateLicenseAgreement_body_GetLicenseAgreement200Response_accepted"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_license_agreement_body_get_license_agreement200_response_accepted
  }
  property {
    name  = "updateLicenseAgreement_body_GetLicenseAgreement200Response_licenseAgreementUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_license_agreement_body_get_license_agreement200_response_license_agreement_url
  }
  property {
    name  = "updateLicense_body_UpdateLicenseRequest_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_license_body_update_license_request_file_data
  }
  property {
    name  = "updateLifetimeSettings_body_GetLifetimeSettings200Response_cacheDuration"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_lifetime_settings_body_get_lifetime_settings200_response_cache_duration
  }
  property {
    name  = "updateLifetimeSettings_body_GetLifetimeSettings200Response_reloadDelay"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_lifetime_settings_body_get_lifetime_settings200_response_reload_delay
  }
  property {
    name  = "updateMapping_body_GetMappings200ResponseItemsInnerAccessTokenManagerRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_mapping_body_get_mappings200_response_items_inner_access_token_manager_ref_id
  }
  property {
    name  = "updateMapping_body_GetMappings200ResponseItemsInnerAccessTokenManagerRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_mapping_body_get_mappings200_response_items_inner_access_token_manager_ref_location
  }
  property {
    name  = "updateMapping_body_GetMappings200ResponseItemsInnerContextContextRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_mapping_body_get_mappings200_response_items_inner_context_context_ref_id
  }
  property {
    name  = "updateMapping_body_GetMappings200ResponseItemsInnerContextContextRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_mapping_body_get_mappings200_response_items_inner_context_context_ref_location
  }
  property {
    name  = "updateMapping_body_GetMappings200ResponseItemsInnerContext_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_mapping_body_get_mappings200_response_items_inner_context_type
  }
  property {
    name  = "updateMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_mapping_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "updateMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_mapping_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "updateMapping_body_GetMappings200ResponseItemsInner_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_mapping_body_get_mappings200_response_items_inner_attribute_contract_fulfillment
  }
  property {
    name  = "updateMapping_body_GetMappings200ResponseItemsInner_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_mapping_body_get_mappings200_response_items_inner_attribute_sources
  }
  property {
    name  = "updateMapping_body_GetMappings200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_mapping_body_get_mappings200_response_items_inner_id
  }
  property {
    name  = "updateMapping_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_mapping_id
  }
  property {
    name  = "updateMapping_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_mapping_x_bypass_external_validation
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_crypto_provider
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_expires"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_expires
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_id
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_issuerDN"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_issuer_dn
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_keyAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_key_algorithm
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_keySize"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_key_size
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_serialNumber"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_serial_number
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_sha1Fingerprint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_sha1_fingerprint
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_sha256Fingerprint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_sha256_fingerprint
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_signatureAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_signature_algorithm
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_status"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_status
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_subjectAlternativeNames"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_subject_alternative_names
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_subjectDN"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_subject_dn
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_validFrom"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_valid_from
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerCertView_version"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_cert_view_version
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerX509File_cryptoProvider"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_x509_file_crypto_provider
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerX509File_fileData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_x509_file_file_data
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInnerX509File_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_x509_file_id
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_id
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_name
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInner_url"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_url
  }
  property {
    name  = "updateMetadataUrl_body_GetMetadataUrls200ResponseItemsInner_validateSignature"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_body_get_metadata_urls200_response_items_inner_validate_signature
  }
  property {
    name  = "updateMetadataUrl_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_metadata_url_id
  }
  property {
    name  = "updateNotificationPublisher_body_CreateNotificationPublisherRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_publisher_body_create_notification_publisher_request_id
  }
  property {
    name  = "updateNotificationPublisher_body_CreateNotificationPublisherRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_publisher_body_create_notification_publisher_request_name
  }
  property {
    name  = "updateNotificationPublisher_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_publisher_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "updateNotificationPublisher_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_publisher_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "updateNotificationPublisher_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_publisher_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "updateNotificationPublisher_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_publisher_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "updateNotificationPublisher_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_publisher_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "updateNotificationPublisher_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_publisher_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "updateNotificationPublisher_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_publisher_id
  }
  property {
    name  = "updateNotificationSettings_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfigNotificationPublisherRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_settings_body_get_identity_profiles200_response_items_inner_email_verification_config_notification_publisher_ref_id
  }
  property {
    name  = "updateNotificationSettings_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfigNotificationPublisherRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_settings_body_get_identity_profiles200_response_items_inner_email_verification_config_notification_publisher_ref_location
  }
  property {
    name  = "updateNotificationSettings_body_GetNotificationSettings200Response_notifyAdminUserPasswordChanges"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_settings_body_get_notification_settings200_response_notify_admin_user_password_changes
  }
  property {
    name  = "updateNotificationSettings_body_GetServerSettings200ResponseNotificationsAccountChangesNotificationPublisherRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_settings_body_get_server_settings200_response_notifications_account_changes_notification_publisher_ref_id
  }
  property {
    name  = "updateNotificationSettings_body_GetServerSettings200ResponseNotificationsAccountChangesNotificationPublisherRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_settings_body_get_server_settings200_response_notifications_account_changes_notification_publisher_ref_location
  }
  property {
    name  = "updateNotificationSettings_body_GetServerSettings200ResponseNotificationsCertificateExpirations_emailAddress"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_settings_body_get_server_settings200_response_notifications_certificate_expirations_email_address
  }
  property {
    name  = "updateNotificationSettings_body_GetServerSettings200ResponseNotificationsCertificateExpirations_finalWarningPeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_settings_body_get_server_settings200_response_notifications_certificate_expirations_final_warning_period
  }
  property {
    name  = "updateNotificationSettings_body_GetServerSettings200ResponseNotificationsCertificateExpirations_initialWarningPeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_settings_body_get_server_settings200_response_notifications_certificate_expirations_initial_warning_period
  }
  property {
    name  = "updateNotificationSettings_body_GetServerSettings200ResponseNotificationsLicenseEvents_emailAddress"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_settings_body_get_server_settings200_response_notifications_license_events_email_address
  }
  property {
    name  = "updateNotificationSettings_body_GetServerSettings200ResponseNotificationsMetadataNotificationSettings_emailAddress"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_notification_settings_body_get_server_settings200_response_notifications_metadata_notification_settings_email_address
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP256ActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p256_active_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP256ActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p256_active_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP256DecryptionActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p256_decryption_active_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP256DecryptionActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p256_decryption_active_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP256DecryptionPreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p256_decryption_previous_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP256DecryptionPreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p256_decryption_previous_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP256PreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p256_previous_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP256PreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p256_previous_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP384ActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p384_active_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP384ActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p384_active_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP384DecryptionActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p384_decryption_active_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP384DecryptionActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p384_decryption_active_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP384DecryptionPreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p384_decryption_previous_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP384DecryptionPreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p384_decryption_previous_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP384PreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p384_previous_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP384PreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p384_previous_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP521ActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p521_active_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP521ActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p521_active_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP521DecryptionActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p521_decryption_active_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP521DecryptionActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p521_decryption_active_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP521DecryptionPreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p521_decryption_previous_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP521DecryptionPreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p521_decryption_previous_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP521PreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p521_previous_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseP521PreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p521_previous_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseRsaActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_rsa_active_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseRsaActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_rsa_active_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseRsaDecryptionActiveCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_rsa_decryption_active_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseRsaDecryptionActiveCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_rsa_decryption_active_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseRsaDecryptionPreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_rsa_decryption_previous_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseRsaDecryptionPreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_rsa_decryption_previous_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseRsaPreviousCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_rsa_previous_cert_ref_id
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200ResponseRsaPreviousCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_rsa_previous_cert_ref_location
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200Response_p256DecryptionPublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p256_decryption_publish_x5c_parameter
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200Response_p256PublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p256_publish_x5c_parameter
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200Response_p384DecryptionPublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p384_decryption_publish_x5c_parameter
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200Response_p384PublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p384_publish_x5c_parameter
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200Response_p521DecryptionPublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p521_decryption_publish_x5c_parameter
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200Response_p521PublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_p521_publish_x5c_parameter
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200Response_rsaDecryptionPublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_rsa_decryption_publish_x5c_parameter
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200Response_rsaPublishX5cParameter"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_rsa_publish_x5c_parameter
  }
  property {
    name  = "updateOAuthOidcKeysSettings_body_GetOauthOidcKeysSettings200Response_staticJwksEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oauth_oidc_keys_settings_body_get_oauth_oidc_keys_settings200_response_static_jwks_enabled
  }
  property {
    name  = "updateOOBAuthenticator_body_CreateOOBAuthenticatorRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oobauthenticator_body_create_oobauthenticator_request_id
  }
  property {
    name  = "updateOOBAuthenticator_body_CreateOOBAuthenticatorRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oobauthenticator_body_create_oobauthenticator_request_name
  }
  property {
    name  = "updateOOBAuthenticator_body_GetOOBAuthenticators200ResponseItemsInnerAllOfAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oobauthenticator_body_get_oobauthenticators200_response_items_inner_all_of_attribute_contract_core_attributes
  }
  property {
    name  = "updateOOBAuthenticator_body_GetOOBAuthenticators200ResponseItemsInnerAllOfAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oobauthenticator_body_get_oobauthenticators200_response_items_inner_all_of_attribute_contract_extended_attributes
  }
  property {
    name  = "updateOOBAuthenticator_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oobauthenticator_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "updateOOBAuthenticator_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oobauthenticator_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "updateOOBAuthenticator_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oobauthenticator_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "updateOOBAuthenticator_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oobauthenticator_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "updateOOBAuthenticator_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oobauthenticator_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "updateOOBAuthenticator_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oobauthenticator_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "updateOOBAuthenticator_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_oobauthenticator_id
  }
  property {
    name  = "updateOutBoundProvisioningSettings_body_GetMappings200ResponseItemsInnerAttributeSourcesInnerDataStoreRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_out_bound_provisioning_settings_body_get_mappings200_response_items_inner_attribute_sources_inner_data_store_ref_id
  }
  property {
    name  = "updateOutBoundProvisioningSettings_body_GetMappings200ResponseItemsInnerAttributeSourcesInnerDataStoreRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_out_bound_provisioning_settings_body_get_mappings200_response_items_inner_attribute_sources_inner_data_store_ref_location
  }
  property {
    name  = "updateOutBoundProvisioningSettings_body_GetOutBoundProvisioningSettings200Response_synchronizationFrequency"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_out_bound_provisioning_settings_body_get_out_bound_provisioning_settings200_response_synchronization_frequency
  }
  property {
    name  = "updatePasswordCredentialValidator_body_CreatePasswordCredentialValidatorRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_password_credential_validator_body_create_password_credential_validator_request_id
  }
  property {
    name  = "updatePasswordCredentialValidator_body_CreatePasswordCredentialValidatorRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_password_credential_validator_body_create_password_credential_validator_request_name
  }
  property {
    name  = "updatePasswordCredentialValidator_body_GetPasswordCredentialValidators200ResponseItemsInnerAllOfAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_password_credential_validator_body_get_password_credential_validators200_response_items_inner_all_of_attribute_contract_core_attributes
  }
  property {
    name  = "updatePasswordCredentialValidator_body_GetPasswordCredentialValidators200ResponseItemsInnerAllOfAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_password_credential_validator_body_get_password_credential_validators200_response_items_inner_all_of_attribute_contract_extended_attributes
  }
  property {
    name  = "updatePasswordCredentialValidator_body_GetPasswordCredentialValidators200ResponseItemsInnerAllOfAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_password_credential_validator_body_get_password_credential_validators200_response_items_inner_all_of_attribute_contract_inherited
  }
  property {
    name  = "updatePasswordCredentialValidator_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_password_credential_validator_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "updatePasswordCredentialValidator_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_password_credential_validator_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "updatePasswordCredentialValidator_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_password_credential_validator_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "updatePasswordCredentialValidator_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_password_credential_validator_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "updatePasswordCredentialValidator_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_password_credential_validator_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "updatePasswordCredentialValidator_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_password_credential_validator_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "updatePasswordCredentialValidator_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_password_credential_validator_id
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_active"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_active
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_creationDate"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_creation_date
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_credential"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_credential
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_credentialId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_credential_id
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_description
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_encryptedCredential"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_encrypted_credential
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_environmentId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_environment_id
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_id
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_name
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_organizationName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_organization_name
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_pingOneAuthenticationApiEndpoint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_ping_one_authentication_api_endpoint
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_pingOneConnectionId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_ping_one_connection_id
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_pingOneManagementApiEndpoint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_ping_one_management_api_endpoint
  }
  property {
    name  = "updatePingOneConnection_body_GetPingOneConnections200ResponseItemsInner_region"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_body_get_ping_one_connections200_response_items_inner_region
  }
  property {
    name  = "updatePingOneConnection_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_id
  }
  property {
    name  = "updatePingOneConnection_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_connection_x_bypass_external_validation
  }
  property {
    name  = "updatePingOneSettings_body_UpdatePingOneForEnterpriseIdentityRepository200ResponsePingOneSsoConnection_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_settings_body_update_ping_one_for_enterprise_identity_repository200_response_ping_one_sso_connection_id
  }
  property {
    name  = "updatePingOneSettings_body_UpdatePingOneForEnterpriseIdentityRepository200ResponsePingOneSsoConnection_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_settings_body_update_ping_one_for_enterprise_identity_repository200_response_ping_one_sso_connection_location
  }
  property {
    name  = "updatePingOneSettings_body_UpdatePingOneForEnterpriseIdentityRepository200Response_companyName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_settings_body_update_ping_one_for_enterprise_identity_repository200_response_company_name
  }
  property {
    name  = "updatePingOneSettings_body_UpdatePingOneForEnterpriseIdentityRepository200Response_connectedToPingOneForEnterprise"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_settings_body_update_ping_one_for_enterprise_identity_repository200_response_connected_to_ping_one_for_enterprise
  }
  property {
    name  = "updatePingOneSettings_body_UpdatePingOneForEnterpriseIdentityRepository200Response_currentAuthnKeyCreationTime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_settings_body_update_ping_one_for_enterprise_identity_repository200_response_current_authn_key_creation_time
  }
  property {
    name  = "updatePingOneSettings_body_UpdatePingOneForEnterpriseIdentityRepository200Response_enableAdminConsoleSso"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_settings_body_update_ping_one_for_enterprise_identity_repository200_response_enable_admin_console_sso
  }
  property {
    name  = "updatePingOneSettings_body_UpdatePingOneForEnterpriseIdentityRepository200Response_enableMonitoring"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_settings_body_update_ping_one_for_enterprise_identity_repository200_response_enable_monitoring
  }
  property {
    name  = "updatePingOneSettings_body_UpdatePingOneForEnterpriseIdentityRepository200Response_identityRepositoryUpdateRequired"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_settings_body_update_ping_one_for_enterprise_identity_repository200_response_identity_repository_update_required
  }
  property {
    name  = "updatePingOneSettings_body_UpdatePingOneForEnterpriseIdentityRepository200Response_previousAuthnKeyCreationTime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ping_one_settings_body_update_ping_one_for_enterprise_identity_repository200_response_previous_authn_key_creation_time
  }
  property {
    name  = "updatePolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerAuthenticationApiApplicationRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_authentication_api_application_ref_id
  }
  property {
    name  = "updatePolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerAuthenticationApiApplicationRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_authentication_api_application_ref_location
  }
  property {
    name  = "updatePolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerRootNodeAction_context"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_root_node_action_context
  }
  property {
    name  = "updatePolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerRootNodeAction_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_root_node_action_type
  }
  property {
    name  = "updatePolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInnerRootNode_children"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_root_node_children
  }
  property {
    name  = "updatePolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInner_description"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_description
  }
  property {
    name  = "updatePolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInner_enabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_enabled
  }
  property {
    name  = "updatePolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInner_handleFailuresLocally"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_handle_failures_locally
  }
  property {
    name  = "updatePolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_id
  }
  property {
    name  = "updatePolicy1_body_GetDefaultAuthenticationPolicy200ResponseAuthnSelectionTreesInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy1_body_get_default_authentication_policy200_response_authn_selection_trees_inner_name
  }
  property {
    name  = "updatePolicy1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy1_id
  }
  property {
    name  = "updatePolicy1_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy1_x_bypass_external_validation
  }
  property {
    name  = "updatePolicy2_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "updatePolicy2_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInnerAuthenticatorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_authenticator_ref_id
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInnerAuthenticatorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_authenticator_ref_location
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInnerIdentityHintContractFulfillment_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_identity_hint_contract_fulfillment_attribute_contract_fulfillment
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInnerIdentityHintContractFulfillment_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_identity_hint_contract_fulfillment_attribute_sources
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInnerIdentityHintContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_identity_hint_contract_core_attributes
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInnerIdentityHintContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_identity_hint_contract_extended_attributes
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInnerIdentityHintMapping_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_identity_hint_mapping_attribute_contract_fulfillment
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInnerIdentityHintMapping_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_identity_hint_mapping_attribute_sources
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInnerUserCodePcvRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_user_code_pcv_ref_id
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInnerUserCodePcvRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_user_code_pcv_ref_location
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInner_allowUnsignedLoginHintToken"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_allow_unsigned_login_hint_token
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInner_alternativeLoginHintTokenIssuers"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_alternative_login_hint_token_issuers
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_id
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_name
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInner_requireTokenForIdentityHint"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_require_token_for_identity_hint
  }
  property {
    name  = "updatePolicy2_body_GetPolicies1200ResponseItemsInner_transactionLifetime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_body_get_policies1200_response_items_inner_transaction_lifetime
  }
  property {
    name  = "updatePolicy2_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_id
  }
  property {
    name  = "updatePolicy2_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy2_x_bypass_external_validation
  }
  property {
    name  = "updatePolicy3_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "updatePolicy3_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInnerAccessTokenManagerRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_access_token_manager_ref_id
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInnerAccessTokenManagerRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_access_token_manager_ref_location
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInnerAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_attribute_contract_core_attributes
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInnerAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_attribute_contract_extended_attributes
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInnerAttributeMapping_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_attribute_mapping_attribute_contract_fulfillment
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInnerAttributeMapping_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_attribute_mapping_attribute_sources
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_id
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInner_idTokenLifetime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_id_token_lifetime
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInner_includeSHashInIdToken"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_include_shash_in_id_token
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInner_includeSriInIdToken"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_include_sri_in_id_token
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInner_includeUserInfoInIdToken"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_include_user_info_in_id_token
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_name
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInner_reissueIdTokenInHybridFlow"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_reissue_id_token_in_hybrid_flow
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInner_returnIdTokenOnRefreshGrant"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_return_id_token_on_refresh_grant
  }
  property {
    name  = "updatePolicy3_body_GetPolicies2200ResponseItemsInner_scopeAttributeMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_body_get_policies2200_response_items_inner_scope_attribute_mappings
  }
  property {
    name  = "updatePolicy3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_id
  }
  property {
    name  = "updatePolicy3_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy3_x_bypass_external_validation
  }
  property {
    name  = "updatePolicy4_body_GetPolicies3200ResponseItemsInnerAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy4_body_get_policies3200_response_items_inner_attribute_contract_core_attributes
  }
  property {
    name  = "updatePolicy4_body_GetPolicies3200ResponseItemsInnerAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy4_body_get_policies3200_response_items_inner_attribute_contract_extended_attributes
  }
  property {
    name  = "updatePolicy4_body_GetPolicies3200ResponseItemsInner_actorTokenRequired"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy4_body_get_policies3200_response_items_inner_actor_token_required
  }
  property {
    name  = "updatePolicy4_body_GetPolicies3200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy4_body_get_policies3200_response_items_inner_id
  }
  property {
    name  = "updatePolicy4_body_GetPolicies3200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy4_body_get_policies3200_response_items_inner_name
  }
  property {
    name  = "updatePolicy4_body_GetPolicies3200ResponseItemsInner_processorMappings"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy4_body_get_policies3200_response_items_inner_processor_mappings
  }
  property {
    name  = "updatePolicy4_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy4_bypass_external_validation
  }
  property {
    name  = "updatePolicy4_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_policy4_id
  }
  property {
    name  = "updateRedirectValidationSettings_body_GetRedirectValidationSettings200ResponseRedirectValidationLocalSettings_enableInErrorResourceValidation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_redirect_validation_settings_body_get_redirect_validation_settings200_response_redirect_validation_local_settings_enable_in_error_resource_validation
  }
  property {
    name  = "updateRedirectValidationSettings_body_GetRedirectValidationSettings200ResponseRedirectValidationLocalSettings_enableTargetResourceValidationForIdpDiscovery"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_redirect_validation_settings_body_get_redirect_validation_settings200_response_redirect_validation_local_settings_enable_target_resource_validation_for_idp_discovery
  }
  property {
    name  = "updateRedirectValidationSettings_body_GetRedirectValidationSettings200ResponseRedirectValidationLocalSettings_enableTargetResourceValidationForSLO"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_redirect_validation_settings_body_get_redirect_validation_settings200_response_redirect_validation_local_settings_enable_target_resource_validation_for_slo
  }
  property {
    name  = "updateRedirectValidationSettings_body_GetRedirectValidationSettings200ResponseRedirectValidationLocalSettings_enableTargetResourceValidationForSSO"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_redirect_validation_settings_body_get_redirect_validation_settings200_response_redirect_validation_local_settings_enable_target_resource_validation_for_sso
  }
  property {
    name  = "updateRedirectValidationSettings_body_GetRedirectValidationSettings200ResponseRedirectValidationLocalSettings_whiteList"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_redirect_validation_settings_body_get_redirect_validation_settings200_response_redirect_validation_local_settings_white_list
  }
  property {
    name  = "updateRedirectValidationSettings_body_GetRedirectValidationSettings200ResponseRedirectValidationPartnerSettings_enableWreplyValidationSLO"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_redirect_validation_settings_body_get_redirect_validation_settings200_response_redirect_validation_partner_settings_enable_wreply_validation_slo
  }
  property {
    name  = "updateResourceOwnerCredentialsMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_resource_owner_credentials_mapping_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "updateResourceOwnerCredentialsMapping_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_resource_owner_credentials_mapping_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "updateResourceOwnerCredentialsMapping_body_GetResourceOwnerCredentialsMappings200ResponseItemsInnerPasswordValidatorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_resource_owner_credentials_mapping_body_get_resource_owner_credentials_mappings200_response_items_inner_password_validator_ref_id
  }
  property {
    name  = "updateResourceOwnerCredentialsMapping_body_GetResourceOwnerCredentialsMappings200ResponseItemsInnerPasswordValidatorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_resource_owner_credentials_mapping_body_get_resource_owner_credentials_mappings200_response_items_inner_password_validator_ref_location
  }
  property {
    name  = "updateResourceOwnerCredentialsMapping_body_GetResourceOwnerCredentialsMappings200ResponseItemsInner_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_resource_owner_credentials_mapping_body_get_resource_owner_credentials_mappings200_response_items_inner_attribute_contract_fulfillment
  }
  property {
    name  = "updateResourceOwnerCredentialsMapping_body_GetResourceOwnerCredentialsMappings200ResponseItemsInner_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_resource_owner_credentials_mapping_body_get_resource_owner_credentials_mappings200_response_items_inner_attribute_sources
  }
  property {
    name  = "updateResourceOwnerCredentialsMapping_body_GetResourceOwnerCredentialsMappings200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_resource_owner_credentials_mapping_body_get_resource_owner_credentials_mappings200_response_items_inner_id
  }
  property {
    name  = "updateResourceOwnerCredentialsMapping_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_resource_owner_credentials_mapping_id
  }
  property {
    name  = "updateResourceOwnerCredentialsMapping_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_resource_owner_credentials_mapping_x_bypass_external_validation
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseCrlSettings_nextRetryMinsWhenNextUpdateInPast"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_crl_settings_next_retry_mins_when_next_update_in_past
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseCrlSettings_nextRetryMinsWhenResolveFailed"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_crl_settings_next_retry_mins_when_resolve_failed
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseCrlSettings_treatNonRetrievableCrlAsRevoked"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_crl_settings_treat_non_retrievable_crl_as_revoked
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseCrlSettings_verifyCrlSignature"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_crl_settings_verify_crl_signature
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseOcspSettingsResponderCertReference_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_ocsp_settings_responder_cert_reference_id
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseOcspSettingsResponderCertReference_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_ocsp_settings_responder_cert_reference_location
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseOcspSettings_actionOnResponderUnavailable"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_ocsp_settings_action_on_responder_unavailable
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseOcspSettings_actionOnStatusUnknown"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_ocsp_settings_action_on_status_unknown
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseOcspSettings_actionOnUnsuccessfulResponse"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_ocsp_settings_action_on_unsuccessful_response
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseOcspSettings_currentUpdateGracePeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_ocsp_settings_current_update_grace_period
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseOcspSettings_nextUpdateGracePeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_ocsp_settings_next_update_grace_period
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseOcspSettings_requesterAddNonce"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_ocsp_settings_requester_add_nonce
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseOcspSettings_responderTimeout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_ocsp_settings_responder_timeout
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseOcspSettings_responderUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_ocsp_settings_responder_url
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseOcspSettings_responseCachePeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_ocsp_settings_response_cache_period
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseProxySettings_host"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_proxy_settings_host
  }
  property {
    name  = "updateRevocationSettings_body_GetRevocationSettings200ResponseProxySettings_port"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_revocation_settings_body_get_revocation_settings200_response_proxy_settings_port
  }
  property {
    name  = "updateRotationSettings_body_GetRotationSettings200Response_activationBufferDays"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_rotation_settings_body_get_rotation_settings200_response_activation_buffer_days
  }
  property {
    name  = "updateRotationSettings_body_GetRotationSettings200Response_creationBufferDays"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_rotation_settings_body_get_rotation_settings200_response_creation_buffer_days
  }
  property {
    name  = "updateRotationSettings_body_GetRotationSettings200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_rotation_settings_body_get_rotation_settings200_response_id
  }
  property {
    name  = "updateRotationSettings_body_GetRotationSettings200Response_keyAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_rotation_settings_body_get_rotation_settings200_response_key_algorithm
  }
  property {
    name  = "updateRotationSettings_body_GetRotationSettings200Response_keySize"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_rotation_settings_body_get_rotation_settings200_response_key_size
  }
  property {
    name  = "updateRotationSettings_body_GetRotationSettings200Response_signatureAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_rotation_settings_body_get_rotation_settings200_response_signature_algorithm
  }
  property {
    name  = "updateRotationSettings_body_GetRotationSettings200Response_validDays"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_rotation_settings_body_get_rotation_settings200_response_valid_days
  }
  property {
    name  = "updateRotationSettings_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_rotation_settings_id
  }
  property {
    name  = "updateSecretManager_body_CreateSecretManagerRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_secret_manager_body_create_secret_manager_request_id
  }
  property {
    name  = "updateSecretManager_body_CreateSecretManagerRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_secret_manager_body_create_secret_manager_request_name
  }
  property {
    name  = "updateSecretManager_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_secret_manager_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "updateSecretManager_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_secret_manager_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "updateSecretManager_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_secret_manager_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "updateSecretManager_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_secret_manager_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "updateSecretManager_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_secret_manager_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "updateSecretManager_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_secret_manager_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "updateSecretManager_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_secret_manager_id
  }
  property {
    name  = "updateServerSettings_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfigNotificationPublisherRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_identity_profiles200_response_items_inner_email_verification_config_notification_publisher_ref_id
  }
  property {
    name  = "updateServerSettings_body_GetIdentityProfiles200ResponseItemsInnerEmailVerificationConfigNotificationPublisherRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_identity_profiles200_response_items_inner_email_verification_config_notification_publisher_ref_location
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseCaptchaSettings_encryptedSecretKey"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_captcha_settings_encrypted_secret_key
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseCaptchaSettings_secretKey"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_captcha_settings_secret_key
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseCaptchaSettings_siteKey"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_captcha_settings_site_key
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseContactInfo_company"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_contact_info_company
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseContactInfo_email"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_contact_info_email
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseContactInfo_firstName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_contact_info_first_name
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseContactInfo_lastName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_contact_info_last_name
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseContactInfo_phone"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_contact_info_phone
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_emailServer"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_email_server
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_enableUtf8MessageHeaders"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_enable_utf8_message_headers
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_encryptedPassword"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_encrypted_password
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_password"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_password
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_port"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_port
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_retryAttempts"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_retry_attempts
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_retryDelay"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_retry_delay
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_sourceAddr"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_source_addr
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_sslPort"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_ssl_port
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_timeout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_timeout
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_useDebugging"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_use_debugging
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_useSSL"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_use_ssl
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_useTLS"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_use_tls
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_username"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_username
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseEmailServer_verifyHostname"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_email_server_verify_hostname
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseFederationInfo_autoConnectEntityId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_federation_info_auto_connect_entity_id
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseFederationInfo_baseUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_federation_info_base_url
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseFederationInfo_saml1xIssuerId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_federation_info_saml1x_issuer_id
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseFederationInfo_saml1xSourceId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_federation_info_saml1x_source_id
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseFederationInfo_saml2EntityId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_federation_info_saml2_entity_id
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseFederationInfo_wsfedRealm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_federation_info_wsfed_realm
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseNotificationsAccountChangesNotificationPublisherRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_notifications_account_changes_notification_publisher_ref_id
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseNotificationsAccountChangesNotificationPublisherRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_notifications_account_changes_notification_publisher_ref_location
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseNotificationsCertificateExpirations_emailAddress"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_notifications_certificate_expirations_email_address
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseNotificationsCertificateExpirations_finalWarningPeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_notifications_certificate_expirations_final_warning_period
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseNotificationsCertificateExpirations_initialWarningPeriod"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_notifications_certificate_expirations_initial_warning_period
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseNotificationsLicenseEvents_emailAddress"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_notifications_license_events_email_address
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseNotificationsMetadataNotificationSettings_emailAddress"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_notifications_metadata_notification_settings_email_address
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseNotifications_notifyAdminUserPasswordChanges"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_notifications_notify_admin_user_password_changes
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsIdpRoleAllOf1Saml20Profile_enable"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_idp_role_all_of1_saml20_profile_enable
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsIdpRoleAllOf1Saml20Profile_enableAutoConnect"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_idp_role_all_of1_saml20_profile_enable_auto_connect
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsIdpRole_enable"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_idp_role_enable
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsIdpRole_enableOutboundProvisioning"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_idp_role_enable_outbound_provisioning
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsIdpRole_enableSaml10"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_idp_role_enable_saml10
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsIdpRole_enableSaml11"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_idp_role_enable_saml11
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsIdpRole_enableWsFed"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_idp_role_enable_ws_fed
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsIdpRole_enableWsTrust"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_idp_role_enable_ws_trust
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsOauthRole_enableOauth"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_oauth_role_enable_oauth
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsOauthRole_enableOpenIdConnect"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_oauth_role_enable_open_id_connect
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsSpRoleAllOfSaml20Profile_enable"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_sp_role_all_of_saml20_profile_enable
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsSpRoleAllOfSaml20Profile_enableAutoConnect"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_sp_role_all_of_saml20_profile_enable_auto_connect
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsSpRoleAllOfSaml20Profile_enableXASP"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_sp_role_all_of_saml20_profile_enable_xasp
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsSpRole_enable"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_sp_role_enable
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsSpRole_enableInboundProvisioning"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_sp_role_enable_inbound_provisioning
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsSpRole_enableOpenIDConnect"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_sp_role_enable_open_idconnect
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsSpRole_enableSaml10"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_sp_role_enable_saml10
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsSpRole_enableSaml11"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_sp_role_enable_saml11
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsSpRole_enableWsFed"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_sp_role_enable_ws_fed
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocolsSpRole_enableWsTrust"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_sp_role_enable_ws_trust
  }
  property {
    name  = "updateServerSettings_body_GetServerSettings200ResponseRolesAndProtocols_enableIdpDiscovery"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_server_settings_body_get_server_settings200_response_roles_and_protocols_enable_idp_discovery
  }
  property {
    name  = "updateServiceAuthentication_body_GetServiceAuthentication200ResponseAttributeQuery_encryptedSharedSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_service_authentication_body_get_service_authentication200_response_attribute_query_encrypted_shared_secret
  }
  property {
    name  = "updateServiceAuthentication_body_GetServiceAuthentication200ResponseAttributeQuery_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_service_authentication_body_get_service_authentication200_response_attribute_query_id
  }
  property {
    name  = "updateServiceAuthentication_body_GetServiceAuthentication200ResponseAttributeQuery_sharedSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_service_authentication_body_get_service_authentication200_response_attribute_query_shared_secret
  }
  property {
    name  = "updateServiceAuthentication_body_GetServiceAuthentication200ResponseConnectionManagement_encryptedSharedSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_service_authentication_body_get_service_authentication200_response_connection_management_encrypted_shared_secret
  }
  property {
    name  = "updateServiceAuthentication_body_GetServiceAuthentication200ResponseConnectionManagement_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_service_authentication_body_get_service_authentication200_response_connection_management_id
  }
  property {
    name  = "updateServiceAuthentication_body_GetServiceAuthentication200ResponseConnectionManagement_sharedSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_service_authentication_body_get_service_authentication200_response_connection_management_shared_secret
  }
  property {
    name  = "updateServiceAuthentication_body_GetServiceAuthentication200ResponseJmx_encryptedSharedSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_service_authentication_body_get_service_authentication200_response_jmx_encrypted_shared_secret
  }
  property {
    name  = "updateServiceAuthentication_body_GetServiceAuthentication200ResponseJmx_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_service_authentication_body_get_service_authentication200_response_jmx_id
  }
  property {
    name  = "updateServiceAuthentication_body_GetServiceAuthentication200ResponseJmx_sharedSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_service_authentication_body_get_service_authentication200_response_jmx_shared_secret
  }
  property {
    name  = "updateServiceAuthentication_body_GetServiceAuthentication200ResponseSsoDirectoryService_encryptedSharedSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_service_authentication_body_get_service_authentication200_response_sso_directory_service_encrypted_shared_secret
  }
  property {
    name  = "updateServiceAuthentication_body_GetServiceAuthentication200ResponseSsoDirectoryService_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_service_authentication_body_get_service_authentication200_response_sso_directory_service_id
  }
  property {
    name  = "updateServiceAuthentication_body_GetServiceAuthentication200ResponseSsoDirectoryService_sharedSecret"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_service_authentication_body_get_service_authentication200_response_sso_directory_service_shared_secret
  }
  property {
    name  = "updateSessionSettings_body_GetSessionSettings200Response_revokeUserSessionOnLogout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_session_settings_body_get_session_settings200_response_revoke_user_session_on_logout
  }
  property {
    name  = "updateSessionSettings_body_GetSessionSettings200Response_sessionRevocationLifetime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_session_settings_body_get_session_settings200_response_session_revocation_lifetime
  }
  property {
    name  = "updateSessionSettings_body_GetSessionSettings200Response_trackAdapterSessionsForLogout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_session_settings_body_get_session_settings200_response_track_adapter_sessions_for_logout
  }
  property {
    name  = "updateSetting_body_GetSetting200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_setting_body_get_setting200_response_id
  }
  property {
    name  = "updateSetting_body_GetSetting200Response_listValue"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_setting_body_get_setting200_response_list_value
  }
  property {
    name  = "updateSetting_body_GetSetting200Response_mapValue"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_setting_body_get_setting200_response_map_value
  }
  property {
    name  = "updateSetting_body_GetSetting200Response_stringValue"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_setting_body_get_setting200_response_string_value
  }
  property {
    name  = "updateSetting_body_GetSetting200Response_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_setting_body_get_setting200_response_type
  }
  property {
    name  = "updateSetting_bundle"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_setting_bundle
  }
  property {
    name  = "updateSetting_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_setting_id
  }
  property {
    name  = "updateSettings1_body_GetSettings1200ResponseDefaultAccessTokenManagerRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings1_body_get_settings1200_response_default_access_token_manager_ref_id
  }
  property {
    name  = "updateSettings1_body_GetSettings1200ResponseDefaultAccessTokenManagerRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings1_body_get_settings1200_response_default_access_token_manager_ref_location
  }
  property {
    name  = "updateSettings2_body_GetSettings2200Response_enableIdpAuthnSelection"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings2_body_get_settings2200_response_enable_idp_authn_selection
  }
  property {
    name  = "updateSettings2_body_GetSettings2200Response_enableSpAuthnSelection"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings2_body_get_settings2200_response_enable_sp_authn_selection
  }
  property {
    name  = "updateSettings3_body_GetKerberosRealmSettings200Response_debugLogOutput"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings3_body_get_kerberos_realm_settings200_response_debug_log_output
  }
  property {
    name  = "updateSettings3_body_GetKerberosRealmSettings200Response_forceTcp"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings3_body_get_kerberos_realm_settings200_response_force_tcp
  }
  property {
    name  = "updateSettings3_body_GetKerberosRealmSettings200Response_kdcRetries"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings3_body_get_kerberos_realm_settings200_response_kdc_retries
  }
  property {
    name  = "updateSettings3_body_GetKerberosRealmSettings200Response_kdcTimeout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings3_body_get_kerberos_realm_settings200_response_kdc_timeout
  }
  property {
    name  = "updateSettings3_body_GetKerberosRealmSettings200Response_keySetRetentionPeriodMins"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings3_body_get_kerberos_realm_settings200_response_key_set_retention_period_mins
  }
  property {
    name  = "updateSettings4_body_GetSettings4200ResponseDefaultNotificationPublisherRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings4_body_get_settings4200_response_default_notification_publisher_ref_id
  }
  property {
    name  = "updateSettings4_body_GetSettings4200ResponseDefaultNotificationPublisherRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings4_body_get_settings4200_response_default_notification_publisher_ref_location
  }
  property {
    name  = "updateSettings5_body_GetSettings5200ResponseDefaultRequestPolicyRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings5_body_get_settings5200_response_default_request_policy_ref_id
  }
  property {
    name  = "updateSettings5_body_GetSettings5200ResponseDefaultRequestPolicyRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings5_body_get_settings5200_response_default_request_policy_ref_location
  }
  property {
    name  = "updateSettings5_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings5_x_bypass_external_validation
  }
  property {
    name  = "updateSettings6_body_GetSettings6200ResponseDefaultPolicyRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings6_body_get_settings6200_response_default_policy_ref_id
  }
  property {
    name  = "updateSettings6_body_GetSettings6200ResponseDefaultPolicyRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings6_body_get_settings6200_response_default_policy_ref_location
  }
  property {
    name  = "updateSettings6_body_GetSettings6200ResponseSessionSettings_revokeUserSessionOnLogout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings6_body_get_settings6200_response_session_settings_revoke_user_session_on_logout
  }
  property {
    name  = "updateSettings6_body_GetSettings6200ResponseSessionSettings_sessionRevocationLifetime"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings6_body_get_settings6200_response_session_settings_session_revocation_lifetime
  }
  property {
    name  = "updateSettings6_body_GetSettings6200ResponseSessionSettings_trackUserSessionsForLogout"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings6_body_get_settings6200_response_session_settings_track_user_sessions_for_logout
  }
  property {
    name  = "updateSettings7_body_GetSettings7200ResponseAdminConsoleCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings7_body_get_settings7200_response_admin_console_cert_ref_id
  }
  property {
    name  = "updateSettings7_body_GetSettings7200ResponseAdminConsoleCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings7_body_get_settings7200_response_admin_console_cert_ref_location
  }
  property {
    name  = "updateSettings7_body_GetSettings7200ResponseRuntimeServerCertRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings7_body_get_settings7200_response_runtime_server_cert_ref_id
  }
  property {
    name  = "updateSettings7_body_GetSettings7200ResponseRuntimeServerCertRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings7_body_get_settings7200_response_runtime_server_cert_ref_location
  }
  property {
    name  = "updateSettings7_body_GetSettings7200Response_activeAdminConsoleCerts"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings7_body_get_settings7200_response_active_admin_console_certs
  }
  property {
    name  = "updateSettings7_body_GetSettings7200Response_activeRuntimeServerCerts"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings7_body_get_settings7200_response_active_runtime_server_certs
  }
  property {
    name  = "updateSettings8_body_GetSettings8200ResponseDefaultGeneratorGroupRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings8_body_get_settings8200_response_default_generator_group_ref_id
  }
  property {
    name  = "updateSettings8_body_GetSettings8200ResponseDefaultGeneratorGroupRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings8_body_get_settings8200_response_default_generator_group_ref_location
  }
  property {
    name  = "updateSettings8_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings8_bypass_external_validation
  }
  property {
    name  = "updateSettings9_body_GetSettings9200ResponseDefaultProcessorPolicyRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings9_body_get_settings9200_response_default_processor_policy_ref_id
  }
  property {
    name  = "updateSettings9_body_GetSettings9200ResponseDefaultProcessorPolicyRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings9_body_get_settings9200_response_default_processor_policy_ref_location
  }
  property {
    name  = "updateSettings9_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_settings9_bypass_external_validation
  }
  property {
    name  = "updateSigningSettings1_body_ExportRequestSigningSettingsSigningKeyPairRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings1_body_export_request_signing_settings_signing_key_pair_ref_id
  }
  property {
    name  = "updateSigningSettings1_body_ExportRequestSigningSettingsSigningKeyPairRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings1_body_export_request_signing_settings_signing_key_pair_ref_location
  }
  property {
    name  = "updateSigningSettings1_body_GetSigningSettings1200Response_algorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings1_body_get_signing_settings1200_response_algorithm
  }
  property {
    name  = "updateSigningSettings1_body_GetSigningSettings1200Response_alternativeSigningKeyPairRefs"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings1_body_get_signing_settings1200_response_alternative_signing_key_pair_refs
  }
  property {
    name  = "updateSigningSettings1_body_GetSigningSettings1200Response_includeCertInSignature"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings1_body_get_signing_settings1200_response_include_cert_in_signature
  }
  property {
    name  = "updateSigningSettings1_body_GetSigningSettings1200Response_includeRawKeyInSignature"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings1_body_get_signing_settings1200_response_include_raw_key_in_signature
  }
  property {
    name  = "updateSigningSettings1_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings1_id
  }
  property {
    name  = "updateSigningSettings2_body_GetSigningSettings2200ResponseSigningKeyRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings2_body_get_signing_settings2200_response_signing_key_ref_id
  }
  property {
    name  = "updateSigningSettings2_body_GetSigningSettings2200ResponseSigningKeyRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings2_body_get_signing_settings2200_response_signing_key_ref_location
  }
  property {
    name  = "updateSigningSettings2_body_GetSigningSettings2200Response_signatureAlgorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings2_body_get_signing_settings2200_response_signature_algorithm
  }
  property {
    name  = "updateSigningSettings3_body_ExportRequestSigningSettingsSigningKeyPairRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings3_body_export_request_signing_settings_signing_key_pair_ref_id
  }
  property {
    name  = "updateSigningSettings3_body_ExportRequestSigningSettingsSigningKeyPairRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings3_body_export_request_signing_settings_signing_key_pair_ref_location
  }
  property {
    name  = "updateSigningSettings3_body_GetSigningSettings1200Response_algorithm"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings3_body_get_signing_settings1200_response_algorithm
  }
  property {
    name  = "updateSigningSettings3_body_GetSigningSettings1200Response_alternativeSigningKeyPairRefs"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings3_body_get_signing_settings1200_response_alternative_signing_key_pair_refs
  }
  property {
    name  = "updateSigningSettings3_body_GetSigningSettings1200Response_includeCertInSignature"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings3_body_get_signing_settings1200_response_include_cert_in_signature
  }
  property {
    name  = "updateSigningSettings3_body_GetSigningSettings1200Response_includeRawKeyInSignature"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings3_body_get_signing_settings1200_response_include_raw_key_in_signature
  }
  property {
    name  = "updateSigningSettings3_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_signing_settings3_id
  }
  property {
    name  = "updateSourcePolicy_body_GetDefaultAuthenticationPolicy200ResponseDefaultAuthenticationSourcesInnerSourceRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_source_policy_body_get_default_authentication_policy200_response_default_authentication_sources_inner_source_ref_id
  }
  property {
    name  = "updateSourcePolicy_body_GetDefaultAuthenticationPolicy200ResponseDefaultAuthenticationSourcesInnerSourceRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_source_policy_body_get_default_authentication_policy200_response_default_authentication_sources_inner_source_ref_location
  }
  property {
    name  = "updateSourcePolicy_body_GetSourcePolicies200ResponseItemsInnerAuthenticationSource_type"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_source_policy_body_get_source_policies200_response_items_inner_authentication_source_type
  }
  property {
    name  = "updateSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_authnContextSensitive"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_source_policy_body_get_source_policies200_response_items_inner_authn_context_sensitive
  }
  property {
    name  = "updateSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_enableSessions"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_source_policy_body_get_source_policies200_response_items_inner_enable_sessions
  }
  property {
    name  = "updateSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_source_policy_body_get_source_policies200_response_items_inner_id
  }
  property {
    name  = "updateSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_idleTimeoutMins"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_source_policy_body_get_source_policies200_response_items_inner_idle_timeout_mins
  }
  property {
    name  = "updateSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_maxTimeoutMins"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_source_policy_body_get_source_policies200_response_items_inner_max_timeout_mins
  }
  property {
    name  = "updateSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_persistent"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_source_policy_body_get_source_policies200_response_items_inner_persistent
  }
  property {
    name  = "updateSourcePolicy_body_GetSourcePolicies200ResponseItemsInner_timeoutDisplayUnit"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_source_policy_body_get_source_policies200_response_items_inner_timeout_display_unit
  }
  property {
    name  = "updateSourcePolicy_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_source_policy_id
  }
  property {
    name  = "updateSpAdapter_body_CreateSpAdapterRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_create_sp_adapter_request_id
  }
  property {
    name  = "updateSpAdapter_body_CreateSpAdapterRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_create_sp_adapter_request_name
  }
  property {
    name  = "updateSpAdapter_body_GetConnection1200ResponseIdpBrowserSsoAdapterMappingsInnerAdapterOverrideSettingsAllOfAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_get_connection1200_response_idp_browser_sso_adapter_mappings_inner_adapter_override_settings_all_of_attribute_contract_core_attributes
  }
  property {
    name  = "updateSpAdapter_body_GetConnection1200ResponseIdpBrowserSsoAdapterMappingsInnerAdapterOverrideSettingsAllOfAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_get_connection1200_response_idp_browser_sso_adapter_mappings_inner_adapter_override_settings_all_of_attribute_contract_extended_attributes
  }
  property {
    name  = "updateSpAdapter_body_GetConnection1200ResponseIdpBrowserSsoAdapterMappingsInnerAdapterOverrideSettingsAllOfAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_get_connection1200_response_idp_browser_sso_adapter_mappings_inner_adapter_override_settings_all_of_attribute_contract_inherited
  }
  property {
    name  = "updateSpAdapter_body_GetConnection1200ResponseIdpBrowserSsoAdapterMappingsInnerAdapterOverrideSettingsAllOfTargetApplicationInfo_applicationIconUrl"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_get_connection1200_response_idp_browser_sso_adapter_mappings_inner_adapter_override_settings_all_of_target_application_info_application_icon_url
  }
  property {
    name  = "updateSpAdapter_body_GetConnection1200ResponseIdpBrowserSsoAdapterMappingsInnerAdapterOverrideSettingsAllOfTargetApplicationInfo_applicationName"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_get_connection1200_response_idp_browser_sso_adapter_mappings_inner_adapter_override_settings_all_of_target_application_info_application_name
  }
  property {
    name  = "updateSpAdapter_body_GetConnection1200ResponseIdpBrowserSsoAdapterMappingsInnerAdapterOverrideSettingsAllOfTargetApplicationInfo_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_get_connection1200_response_idp_browser_sso_adapter_mappings_inner_adapter_override_settings_all_of_target_application_info_inherited
  }
  property {
    name  = "updateSpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "updateSpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "updateSpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "updateSpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "updateSpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "updateSpAdapter_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "updateSpAdapter_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sp_adapter_id
  }
  property {
    name  = "updateStsRequestParamContractById_body_GetStsRequestParamContracts200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sts_request_param_contract_by_id_body_get_sts_request_param_contracts200_response_items_inner_id
  }
  property {
    name  = "updateStsRequestParamContractById_body_GetStsRequestParamContracts200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sts_request_param_contract_by_id_body_get_sts_request_param_contracts200_response_items_inner_name
  }
  property {
    name  = "updateStsRequestParamContractById_body_GetStsRequestParamContracts200ResponseItemsInner_parameters"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sts_request_param_contract_by_id_body_get_sts_request_param_contracts200_response_items_inner_parameters
  }
  property {
    name  = "updateStsRequestParamContractById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_sts_request_param_contract_by_id_id
  }
  property {
    name  = "updateSystemKeys_body_GetSystemKeys200ResponseCurrent_creationDate"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_system_keys_body_get_system_keys200_response_current_creation_date
  }
  property {
    name  = "updateSystemKeys_body_GetSystemKeys200ResponseCurrent_encryptedKeyData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_system_keys_body_get_system_keys200_response_current_encrypted_key_data
  }
  property {
    name  = "updateSystemKeys_body_GetSystemKeys200ResponseCurrent_keyData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_system_keys_body_get_system_keys200_response_current_key_data
  }
  property {
    name  = "updateSystemKeys_body_GetSystemKeys200ResponsePending_creationDate"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_system_keys_body_get_system_keys200_response_pending_creation_date
  }
  property {
    name  = "updateSystemKeys_body_GetSystemKeys200ResponsePending_encryptedKeyData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_system_keys_body_get_system_keys200_response_pending_encrypted_key_data
  }
  property {
    name  = "updateSystemKeys_body_GetSystemKeys200ResponsePending_keyData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_system_keys_body_get_system_keys200_response_pending_key_data
  }
  property {
    name  = "updateSystemKeys_body_GetSystemKeys200ResponsePrevious_creationDate"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_system_keys_body_get_system_keys200_response_previous_creation_date
  }
  property {
    name  = "updateSystemKeys_body_GetSystemKeys200ResponsePrevious_encryptedKeyData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_system_keys_body_get_system_keys200_response_previous_encrypted_key_data
  }
  property {
    name  = "updateSystemKeys_body_GetSystemKeys200ResponsePrevious_keyData"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_system_keys_body_get_system_keys200_response_previous_key_data
  }
  property {
    name  = "updateTokenGeneratorMappingById_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_mapping_by_id_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "updateTokenGeneratorMappingById_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_mapping_by_id_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "updateTokenGeneratorMappingById_body_GetTokenGeneratorMappings200ResponseItemsInner_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_mapping_by_id_body_get_token_generator_mappings200_response_items_inner_attribute_contract_fulfillment
  }
  property {
    name  = "updateTokenGeneratorMappingById_body_GetTokenGeneratorMappings200ResponseItemsInner_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_mapping_by_id_body_get_token_generator_mappings200_response_items_inner_attribute_sources
  }
  property {
    name  = "updateTokenGeneratorMappingById_body_GetTokenGeneratorMappings200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_mapping_by_id_body_get_token_generator_mappings200_response_items_inner_id
  }
  property {
    name  = "updateTokenGeneratorMappingById_body_GetTokenGeneratorMappings200ResponseItemsInner_licenseConnectionGroupAssignment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_mapping_by_id_body_get_token_generator_mappings200_response_items_inner_license_connection_group_assignment
  }
  property {
    name  = "updateTokenGeneratorMappingById_body_GetTokenGeneratorMappings200ResponseItemsInner_sourceId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_mapping_by_id_body_get_token_generator_mappings200_response_items_inner_source_id
  }
  property {
    name  = "updateTokenGeneratorMappingById_body_GetTokenGeneratorMappings200ResponseItemsInner_targetId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_mapping_by_id_body_get_token_generator_mappings200_response_items_inner_target_id
  }
  property {
    name  = "updateTokenGeneratorMappingById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_mapping_by_id_id
  }
  property {
    name  = "updateTokenGeneratorMappingById_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_mapping_by_id_x_bypass_external_validation
  }
  property {
    name  = "updateTokenGenerator_body_CreateTokenGeneratorRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_body_create_token_generator_request_id
  }
  property {
    name  = "updateTokenGenerator_body_CreateTokenGeneratorRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_body_create_token_generator_request_name
  }
  property {
    name  = "updateTokenGenerator_body_GetTokenGenerators200ResponseItemsInnerAllOfAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_body_get_token_generators200_response_items_inner_all_of_attribute_contract_core_attributes
  }
  property {
    name  = "updateTokenGenerator_body_GetTokenGenerators200ResponseItemsInnerAllOfAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_body_get_token_generators200_response_items_inner_all_of_attribute_contract_extended_attributes
  }
  property {
    name  = "updateTokenGenerator_body_GetTokenGenerators200ResponseItemsInnerAllOfAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_body_get_token_generators200_response_items_inner_all_of_attribute_contract_inherited
  }
  property {
    name  = "updateTokenGenerator_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "updateTokenGenerator_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "updateTokenGenerator_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "updateTokenGenerator_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "updateTokenGenerator_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "updateTokenGenerator_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "updateTokenGenerator_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_generator_id
  }
  property {
    name  = "updateTokenManager_body_CreateTokenManagerRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_create_token_manager_request_id
  }
  property {
    name  = "updateTokenManager_body_CreateTokenManagerRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_create_token_manager_request_name
  }
  property {
    name  = "updateTokenManager_body_CreateTokenManagerRequest_sequenceNumber"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_create_token_manager_request_sequence_number
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AccessControlSettings_allowedClients"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_access_control_settings_allowed_clients
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AccessControlSettings_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_access_control_settings_inherited
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AccessControlSettings_restrictClients"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_access_control_settings_restrict_clients
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_attribute_contract_core_attributes
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AttributeContract_defaultSubjectAttribute"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_attribute_contract_default_subject_attribute
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_attribute_contract_extended_attributes
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1AttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_attribute_contract_inherited
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SelectionSettings_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_selection_settings_inherited
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SelectionSettings_resourceUris"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_selection_settings_resource_uris
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SessionValidationSettings_checkSessionRevocationStatus"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_session_validation_settings_check_session_revocation_status
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SessionValidationSettings_checkValidAuthnSession"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_session_validation_settings_check_valid_authn_session
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SessionValidationSettings_includeSessionId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_session_validation_settings_include_session_id
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SessionValidationSettings_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_session_validation_settings_inherited
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOf1SessionValidationSettings_updateAuthnSessionActivity"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of1_session_validation_settings_update_authn_session_activity
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "updateTokenManager_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "updateTokenManager_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_manager_id
  }
  property {
    name  = "updateTokenProcessor_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_fields"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_processor_body_get_token_managers200_response_items_inner_all_of_configuration_fields
  }
  property {
    name  = "updateTokenProcessor_body_GetTokenManagers200ResponseItemsInnerAllOfConfiguration_tables"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_processor_body_get_token_managers200_response_items_inner_all_of_configuration_tables
  }
  property {
    name  = "updateTokenProcessor_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_processor_body_get_token_managers200_response_items_inner_all_of_parent_ref_id
  }
  property {
    name  = "updateTokenProcessor_body_GetTokenManagers200ResponseItemsInnerAllOfParentRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_processor_body_get_token_managers200_response_items_inner_all_of_parent_ref_location
  }
  property {
    name  = "updateTokenProcessor_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_processor_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_id
  }
  property {
    name  = "updateTokenProcessor_body_GetTokenManagers200ResponseItemsInnerAllOfPluginDescriptorRef_location"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_processor_body_get_token_managers200_response_items_inner_all_of_plugin_descriptor_ref_location
  }
  property {
    name  = "updateTokenProcessor_body_GetTokenProcessor200ResponseAttributeContract_coreAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_processor_body_get_token_processor200_response_attribute_contract_core_attributes
  }
  property {
    name  = "updateTokenProcessor_body_GetTokenProcessor200ResponseAttributeContract_extendedAttributes"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_processor_body_get_token_processor200_response_attribute_contract_extended_attributes
  }
  property {
    name  = "updateTokenProcessor_body_GetTokenProcessor200ResponseAttributeContract_inherited"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_processor_body_get_token_processor200_response_attribute_contract_inherited
  }
  property {
    name  = "updateTokenProcessor_body_GetTokenProcessor200ResponseAttributeContract_maskOgnlValues"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_processor_body_get_token_processor200_response_attribute_contract_mask_ognl_values
  }
  property {
    name  = "updateTokenProcessor_body_UpdateTokenProcessorRequest_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_processor_body_update_token_processor_request_id
  }
  property {
    name  = "updateTokenProcessor_body_UpdateTokenProcessorRequest_name"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_processor_body_update_token_processor_request_name
  }
  property {
    name  = "updateTokenProcessor_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_processor_id
  }
  property {
    name  = "updateTokenToTokenMappingById_body_GetMappings200ResponseItemsInnerIssuanceCriteria_conditionalCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_to_token_mapping_by_id_body_get_mappings200_response_items_inner_issuance_criteria_conditional_criteria
  }
  property {
    name  = "updateTokenToTokenMappingById_body_GetMappings200ResponseItemsInnerIssuanceCriteria_expressionCriteria"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_to_token_mapping_by_id_body_get_mappings200_response_items_inner_issuance_criteria_expression_criteria
  }
  property {
    name  = "updateTokenToTokenMappingById_body_GetTokenToTokenMappingById200Response_attributeContractFulfillment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_to_token_mapping_by_id_body_get_token_to_token_mapping_by_id200_response_attribute_contract_fulfillment
  }
  property {
    name  = "updateTokenToTokenMappingById_body_GetTokenToTokenMappingById200Response_attributeSources"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_to_token_mapping_by_id_body_get_token_to_token_mapping_by_id200_response_attribute_sources
  }
  property {
    name  = "updateTokenToTokenMappingById_body_GetTokenToTokenMappingById200Response_defaultTargetResource"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_to_token_mapping_by_id_body_get_token_to_token_mapping_by_id200_response_default_target_resource
  }
  property {
    name  = "updateTokenToTokenMappingById_body_GetTokenToTokenMappingById200Response_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_to_token_mapping_by_id_body_get_token_to_token_mapping_by_id200_response_id
  }
  property {
    name  = "updateTokenToTokenMappingById_body_GetTokenToTokenMappingById200Response_licenseConnectionGroupAssignment"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_to_token_mapping_by_id_body_get_token_to_token_mapping_by_id200_response_license_connection_group_assignment
  }
  property {
    name  = "updateTokenToTokenMappingById_body_GetTokenToTokenMappingById200Response_sourceId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_to_token_mapping_by_id_body_get_token_to_token_mapping_by_id200_response_source_id
  }
  property {
    name  = "updateTokenToTokenMappingById_body_GetTokenToTokenMappingById200Response_targetId"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_to_token_mapping_by_id_body_get_token_to_token_mapping_by_id200_response_target_id
  }
  property {
    name  = "updateTokenToTokenMappingById_id"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_to_token_mapping_by_id_id
  }
  property {
    name  = "updateTokenToTokenMappingById_x_bypass_external_validation"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_token_to_token_mapping_by_id_x_bypass_external_validation
  }
  property {
    name  = "updateUrlMappings1_body_GetUrlMappings1200Response_items"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_url_mappings1_body_get_url_mappings1200_response_items
  }
  property {
    name  = "updateUrlMappings2_body_GetUrlMappings2200Response_items"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_url_mappings2_body_get_url_mappings2200_response_items
  }
  property {
    name  = "updateVirtualHostNamesSettings_body_GetVirtualHostNamesSettings200Response_virtualHostNames"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_virtual_host_names_settings_body_get_virtual_host_names_settings200_response_virtual_host_names
  }
  property {
    name  = "updateWsTrustStsSettings_body_GetWsTrustStsSettings200Response_basicAuthnEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ws_trust_sts_settings_body_get_ws_trust_sts_settings200_response_basic_authn_enabled
  }
  property {
    name  = "updateWsTrustStsSettings_body_GetWsTrustStsSettings200Response_clientCertAuthnEnabled"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ws_trust_sts_settings_body_get_ws_trust_sts_settings200_response_client_cert_authn_enabled
  }
  property {
    name  = "updateWsTrustStsSettings_body_GetWsTrustStsSettings200Response_issuerCerts"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ws_trust_sts_settings_body_get_ws_trust_sts_settings200_response_issuer_certs
  }
  property {
    name  = "updateWsTrustStsSettings_body_GetWsTrustStsSettings200Response_restrictByIssuerCert"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ws_trust_sts_settings_body_get_ws_trust_sts_settings200_response_restrict_by_issuer_cert
  }
  property {
    name  = "updateWsTrustStsSettings_body_GetWsTrustStsSettings200Response_restrictBySubjectDn"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ws_trust_sts_settings_body_get_ws_trust_sts_settings200_response_restrict_by_subject_dn
  }
  property {
    name  = "updateWsTrustStsSettings_body_GetWsTrustStsSettings200Response_subjectDns"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ws_trust_sts_settings_body_get_ws_trust_sts_settings200_response_subject_dns
  }
  property {
    name  = "updateWsTrustStsSettings_body_GetWsTrustStsSettings200Response_users"
    type  = "string"
    value = var.connector-oai-pfadminapi_property_update_ws_trust_sts_settings_body_get_ws_trust_sts_settings200_response_users
  }
}
