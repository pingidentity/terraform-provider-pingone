resource "pingone_davinci_connector_instance" "connector-oai-pingaccessadministrativeapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-pingaccessadministrativeapi"
  }
  name = "My awesome connector-oai-pingaccessadministrativeapi"
  property {
    name  = "accessTokenValidatorsGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_get_filter
  }
  property {
    name  = "accessTokenValidatorsGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_get_name
  }
  property {
    name  = "accessTokenValidatorsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_get_number_per_page
  }
  property {
    name  = "accessTokenValidatorsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_get_order
  }
  property {
    name  = "accessTokenValidatorsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_get_page
  }
  property {
    name  = "accessTokenValidatorsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_get_sort_key
  }
  property {
    name  = "accessTokenValidatorsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_id_delete_id
  }
  property {
    name  = "accessTokenValidatorsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_id_get_id
  }
  property {
    name  = "accessTokenValidatorsIdPut_accessTokenValidatorsGet200ResponseItemsInner_AccessTokenValidatorsGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_id_put_access_token_validators_get200_response_items_inner_access_token_validators_get200_response_items_inner_class_name
  }
  property {
    name  = "accessTokenValidatorsIdPut_accessTokenValidatorsGet200ResponseItemsInner_AccessTokenValidatorsGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_id_put_access_token_validators_get200_response_items_inner_access_token_validators_get200_response_items_inner_configuration
  }
  property {
    name  = "accessTokenValidatorsIdPut_accessTokenValidatorsGet200ResponseItemsInner_AccessTokenValidatorsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_id_put_access_token_validators_get200_response_items_inner_access_token_validators_get200_response_items_inner_id
  }
  property {
    name  = "accessTokenValidatorsIdPut_accessTokenValidatorsGet200ResponseItemsInner_AccessTokenValidatorsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_id_put_access_token_validators_get200_response_items_inner_access_token_validators_get200_response_items_inner_name
  }
  property {
    name  = "accessTokenValidatorsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_id_put_id
  }
  property {
    name  = "accessTokenValidatorsPost_accessTokenValidatorsGet200ResponseItemsInner_AccessTokenValidatorsGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_post_access_token_validators_get200_response_items_inner_access_token_validators_get200_response_items_inner_class_name
  }
  property {
    name  = "accessTokenValidatorsPost_accessTokenValidatorsGet200ResponseItemsInner_AccessTokenValidatorsGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_post_access_token_validators_get200_response_items_inner_access_token_validators_get200_response_items_inner_configuration
  }
  property {
    name  = "accessTokenValidatorsPost_accessTokenValidatorsGet200ResponseItemsInner_AccessTokenValidatorsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_post_access_token_validators_get200_response_items_inner_access_token_validators_get200_response_items_inner_id
  }
  property {
    name  = "accessTokenValidatorsPost_accessTokenValidatorsGet200ResponseItemsInner_AccessTokenValidatorsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_access_token_validators_post_access_token_validators_get200_response_items_inner_access_token_validators_get200_response_items_inner_name
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsAcmeCertificateRequestIdDelete_acme_account_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_acme_certificate_request_id_delete_acme_account_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsAcmeCertificateRequestIdDelete_acme_certificate_request_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_acme_certificate_request_id_delete_acme_certificate_request_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsAcmeCertificateRequestIdDelete_acme_server_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_acme_certificate_request_id_delete_acme_server_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsAcmeCertificateRequestIdGet_acme_account_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_acme_certificate_request_id_get_acme_account_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsAcmeCertificateRequestIdGet_acme_certificate_request_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_acme_certificate_request_id_get_acme_certificate_request_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsAcmeCertificateRequestIdGet_acme_server_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_acme_certificate_request_id_get_acme_server_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet_acme_account_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get_acme_account_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet_acme_server_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get_acme_server_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet_key_pair_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get_key_pair_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get_number_per_page
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get_order
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get_page
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get_sort_key
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsPost_acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200Response_AcmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200ResponseAcmeCertStatus_problems"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_post_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_acme_cert_status_problems
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsPost_acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200Response_AcmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200ResponseAcmeCertStatus_state"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_post_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_acme_cert_status_state
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsPost_acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200Response_AcmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200Response_acmeAccountId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_post_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_acme_account_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsPost_acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200Response_AcmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200Response_acmeServerId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_post_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_acme_server_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsPost_acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200Response_AcmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200Response_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_post_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsPost_acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200Response_AcmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200Response_keyPairId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_post_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_key_pair_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsPost_acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200Response_AcmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsGet200Response_url"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_post_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_get200_response_url
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsPost_acme_account_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_post_acme_account_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdCertificateRequestsPost_acme_server_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_certificate_requests_post_acme_server_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdDelete_acme_account_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_delete_acme_account_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdDelete_acme_server_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_delete_acme_server_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdGet_acme_account_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_get_acme_account_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsAcmeAccountIdGet_acme_server_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_acme_account_id_get_acme_server_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsGet_acme_server_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_get_acme_server_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_get_number_per_page
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_get_order
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_get_page
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_get_sort_key
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsPost_acmeServersAcmeServerIdAccountsGet200Response_AcmeServersAcmeServerIdAccountsGet200ResponsePrivateKey_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_post_acme_servers_acme_server_id_accounts_get200_response_acme_servers_acme_server_id_accounts_get200_response_private_key_encrypted_value
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsPost_acmeServersAcmeServerIdAccountsGet200Response_AcmeServersAcmeServerIdAccountsGet200ResponsePrivateKey_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_post_acme_servers_acme_server_id_accounts_get200_response_acme_servers_acme_server_id_accounts_get200_response_private_key_value
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsPost_acmeServersAcmeServerIdAccountsGet200Response_AcmeServersAcmeServerIdAccountsGet200ResponsePublicKey_created"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_post_acme_servers_acme_server_id_accounts_get200_response_acme_servers_acme_server_id_accounts_get200_response_public_key_created
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsPost_acmeServersAcmeServerIdAccountsGet200Response_AcmeServersAcmeServerIdAccountsGet200ResponsePublicKey_jwk"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_post_acme_servers_acme_server_id_accounts_get200_response_acme_servers_acme_server_id_accounts_get200_response_public_key_jwk
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsPost_acmeServersAcmeServerIdAccountsGet200Response_AcmeServersAcmeServerIdAccountsGet200Response_acmeServerId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_post_acme_servers_acme_server_id_accounts_get200_response_acme_servers_acme_server_id_accounts_get200_response_acme_server_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsPost_acmeServersAcmeServerIdAccountsGet200Response_AcmeServersAcmeServerIdAccountsGet200Response_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_post_acme_servers_acme_server_id_accounts_get200_response_acme_servers_acme_server_id_accounts_get200_response_id
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsPost_acmeServersAcmeServerIdAccountsGet200Response_AcmeServersAcmeServerIdAccountsGet200Response_keyAlgorithm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_post_acme_servers_acme_server_id_accounts_get200_response_acme_servers_acme_server_id_accounts_get200_response_key_algorithm
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsPost_acmeServersAcmeServerIdAccountsGet200Response_AcmeServersAcmeServerIdAccountsGet200Response_url"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_post_acme_servers_acme_server_id_accounts_get200_response_acme_servers_acme_server_id_accounts_get200_response_url
  }
  property {
    name  = "acmeServersAcmeServerIdAccountsPost_acme_server_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_accounts_post_acme_server_id
  }
  property {
    name  = "acmeServersAcmeServerIdDelete_acme_server_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_delete_acme_server_id
  }
  property {
    name  = "acmeServersAcmeServerIdGet_acme_server_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_acme_server_id_get_acme_server_id
  }
  property {
    name  = "acmeServersDefaultPut_acmeServersGet200ResponseItemsInnerAcmeAccountsInner_AcmeServersGet200ResponseItemsInnerAcmeAccountsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_default_put_acme_servers_get200_response_items_inner_acme_accounts_inner_acme_servers_get200_response_items_inner_acme_accounts_inner_id
  }
  property {
    name  = "acmeServersDefaultPut_acmeServersGet200ResponseItemsInnerAcmeAccountsInner_AcmeServersGet200ResponseItemsInnerAcmeAccountsInner_location"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_default_put_acme_servers_get200_response_items_inner_acme_accounts_inner_acme_servers_get200_response_items_inner_acme_accounts_inner_location
  }
  property {
    name  = "acmeServersGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_get_filter
  }
  property {
    name  = "acmeServersGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_get_name
  }
  property {
    name  = "acmeServersGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_get_number_per_page
  }
  property {
    name  = "acmeServersGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_get_order
  }
  property {
    name  = "acmeServersGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_get_page
  }
  property {
    name  = "acmeServersGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_get_sort_key
  }
  property {
    name  = "acmeServersPost_acmeServersGet200ResponseItemsInner_AcmeServersGet200ResponseItemsInner_acmeAccounts"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_post_acme_servers_get200_response_items_inner_acme_servers_get200_response_items_inner_acme_accounts
  }
  property {
    name  = "acmeServersPost_acmeServersGet200ResponseItemsInner_AcmeServersGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_post_acme_servers_get200_response_items_inner_acme_servers_get200_response_items_inner_id
  }
  property {
    name  = "acmeServersPost_acmeServersGet200ResponseItemsInner_AcmeServersGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_post_acme_servers_get200_response_items_inner_acme_servers_get200_response_items_inner_name
  }
  property {
    name  = "acmeServersPost_acmeServersGet200ResponseItemsInner_AcmeServersGet200ResponseItemsInner_url"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_acme_servers_post_acme_servers_get200_response_items_inner_acme_servers_get200_response_items_inner_url
  }
  property {
    name  = "adminConfigPut_adminConfigDelete200Response_AdminConfigDelete200Response_hostPort"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_put_admin_config_delete200_response_admin_config_delete200_response_host_port
  }
  property {
    name  = "adminConfigPut_adminConfigDelete200Response_AdminConfigDelete200Response_httpProxyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_put_admin_config_delete200_response_admin_config_delete200_response_http_proxy_id
  }
  property {
    name  = "adminConfigPut_adminConfigDelete200Response_AdminConfigDelete200Response_httpsProxyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_put_admin_config_delete200_response_admin_config_delete200_response_https_proxy_id
  }
  property {
    name  = "adminConfigReplicaAdminsIdConfigPost_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_config_post_id
  }
  property {
    name  = "adminConfigReplicaAdminsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_delete_id
  }
  property {
    name  = "adminConfigReplicaAdminsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_get_id
  }
  property {
    name  = "adminConfigReplicaAdminsIdPut_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInnerCertificateHash_algorithm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_put_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_certificate_hash_algorithm
  }
  property {
    name  = "adminConfigReplicaAdminsIdPut_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInnerCertificateHash_hexValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_put_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_certificate_hash_hex_value
  }
  property {
    name  = "adminConfigReplicaAdminsIdPut_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_configReplicationEnabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_put_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_config_replication_enabled
  }
  property {
    name  = "adminConfigReplicaAdminsIdPut_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_put_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_description
  }
  property {
    name  = "adminConfigReplicaAdminsIdPut_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_hostPort"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_put_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_host_port
  }
  property {
    name  = "adminConfigReplicaAdminsIdPut_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_httpProxyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_put_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_http_proxy_id
  }
  property {
    name  = "adminConfigReplicaAdminsIdPut_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_httpsProxyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_put_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_https_proxy_id
  }
  property {
    name  = "adminConfigReplicaAdminsIdPut_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_put_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_id
  }
  property {
    name  = "adminConfigReplicaAdminsIdPut_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_keys"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_put_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_keys
  }
  property {
    name  = "adminConfigReplicaAdminsIdPut_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_put_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_name
  }
  property {
    name  = "adminConfigReplicaAdminsIdPut_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_selectedCertificateId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_put_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_selected_certificate_id
  }
  property {
    name  = "adminConfigReplicaAdminsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_id_put_id
  }
  property {
    name  = "adminConfigReplicaAdminsPost_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInnerCertificateHash_algorithm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_post_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_certificate_hash_algorithm
  }
  property {
    name  = "adminConfigReplicaAdminsPost_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInnerCertificateHash_hexValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_post_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_certificate_hash_hex_value
  }
  property {
    name  = "adminConfigReplicaAdminsPost_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_configReplicationEnabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_post_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_config_replication_enabled
  }
  property {
    name  = "adminConfigReplicaAdminsPost_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_post_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_description
  }
  property {
    name  = "adminConfigReplicaAdminsPost_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_hostPort"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_post_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_host_port
  }
  property {
    name  = "adminConfigReplicaAdminsPost_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_httpProxyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_post_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_http_proxy_id
  }
  property {
    name  = "adminConfigReplicaAdminsPost_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_httpsProxyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_post_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_https_proxy_id
  }
  property {
    name  = "adminConfigReplicaAdminsPost_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_post_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_id
  }
  property {
    name  = "adminConfigReplicaAdminsPost_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_keys"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_post_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_keys
  }
  property {
    name  = "adminConfigReplicaAdminsPost_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_post_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_name
  }
  property {
    name  = "adminConfigReplicaAdminsPost_adminConfigReplicaAdminsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInner_selectedCertificateId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_admin_config_replica_admins_post_admin_config_replica_admins_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_selected_certificate_id
  }
  property {
    name  = "agentsAgentIdConfigSharedSecretIdGet_agent_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_agent_id_config_shared_secret_id_get_agent_id
  }
  property {
    name  = "agentsAgentIdConfigSharedSecretIdGet_shared_secret_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_agent_id_config_shared_secret_id_get_shared_secret_id
  }
  property {
    name  = "agentsCertificatesGet_alias"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_certificates_get_alias
  }
  property {
    name  = "agentsCertificatesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_certificates_get_filter
  }
  property {
    name  = "agentsCertificatesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_certificates_get_number_per_page
  }
  property {
    name  = "agentsCertificatesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_certificates_get_order
  }
  property {
    name  = "agentsCertificatesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_certificates_get_page
  }
  property {
    name  = "agentsCertificatesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_certificates_get_sort_key
  }
  property {
    name  = "agentsCertificatesIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_certificates_id_get_id
  }
  property {
    name  = "agentsGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_get_filter
  }
  property {
    name  = "agentsGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_get_name
  }
  property {
    name  = "agentsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_get_number_per_page
  }
  property {
    name  = "agentsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_get_order
  }
  property {
    name  = "agentsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_get_page
  }
  property {
    name  = "agentsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_get_sort_key
  }
  property {
    name  = "agentsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_delete_id
  }
  property {
    name  = "agentsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_get_id
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInnerCertificateHash_algorithm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_certificate_hash_algorithm
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInnerCertificateHash_hexValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_certificate_hash_hex_value
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInnerIpSource_fallbackToLastHopIp"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_ip_source_fallback_to_last_hop_ip
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInnerIpSource_headerNameList"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_ip_source_header_name_list
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInnerIpSource_listValueLocation"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_ip_source_list_value_location
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_description
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_failedRetryTimeout"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_failed_retry_timeout
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_failoverHosts"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_failover_hosts
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_hostname"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_hostname
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_id
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_maxRetries"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_max_retries
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_name
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_overrideIpSource"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_override_ip_source
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_port
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_selectedCertificateId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_selected_certificate_id
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_sharedSecretIds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_shared_secret_ids
  }
  property {
    name  = "agentsIdPut_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_unknownResourceMode"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_agents_get200_response_items_inner_agents_get200_response_items_inner_unknown_resource_mode
  }
  property {
    name  = "agentsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_id_put_id
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInnerCertificateHash_algorithm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_certificate_hash_algorithm
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInnerCertificateHash_hexValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_certificate_hash_hex_value
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInnerIpSource_fallbackToLastHopIp"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_ip_source_fallback_to_last_hop_ip
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInnerIpSource_headerNameList"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_ip_source_header_name_list
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInnerIpSource_listValueLocation"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_ip_source_list_value_location
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_description
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_failedRetryTimeout"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_failed_retry_timeout
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_failoverHosts"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_failover_hosts
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_hostname"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_hostname
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_id
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_maxRetries"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_max_retries
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_name
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_overrideIpSource"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_override_ip_source
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_port
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_selectedCertificateId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_selected_certificate_id
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_sharedSecretIds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_shared_secret_ids
  }
  property {
    name  = "agentsPost_agentsGet200ResponseItemsInner_AgentsGet200ResponseItemsInner_unknownResourceMode"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_agents_post_agents_get200_response_items_inner_agents_get200_response_items_inner_unknown_resource_mode
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdDelete_application_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_delete_application_id
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdDelete_resource_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_delete_resource_id
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdGet_application_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_get_application_id
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdGet_resource_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_get_resource_id
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_application_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_application_id
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInnerQueryParamConfig_matchesNoParams"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_query_param_config_matches_no_params
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInnerQueryParamConfig_params"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_query_param_config_params
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInnerResourceTypeConfigurationResponseGenerator_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_resource_type_configuration_response_generator_class_name
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInnerResourceTypeConfigurationResponseGenerator_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_resource_type_configuration_response_generator_configuration
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_anonymous"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_anonymous
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_applicationId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_application_id
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_auditLevel"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_audit_level
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_authenticationChallengePolicyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_authentication_challenge_policy_id
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_defaultAuthTypeOverride"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_default_auth_type_override
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_enabled
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_id
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_methods"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_methods
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_name
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_pathPatterns"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_path_patterns
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_pathPrefixes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_path_prefixes
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_policy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_policy
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_resourceType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_resource_type
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_rootResource"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_root_resource
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_unprotected"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_unprotected
  }
  property {
    name  = "applicationsApplicationIdResourcesResourceIdPut_resource_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_application_id_resources_resource_id_put_resource_id
  }
  property {
    name  = "applicationsGet_agent_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_get_agent_id
  }
  property {
    name  = "applicationsGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_get_filter
  }
  property {
    name  = "applicationsGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_get_name
  }
  property {
    name  = "applicationsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_get_number_per_page
  }
  property {
    name  = "applicationsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_get_order
  }
  property {
    name  = "applicationsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_get_page
  }
  property {
    name  = "applicationsGet_rule_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_get_rule_id
  }
  property {
    name  = "applicationsGet_ruleset_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_get_ruleset_id
  }
  property {
    name  = "applicationsGet_sideband_client_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_get_sideband_client_id
  }
  property {
    name  = "applicationsGet_site_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_get_site_id
  }
  property {
    name  = "applicationsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_get_sort_key
  }
  property {
    name  = "applicationsGet_virtual_host_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_get_virtual_host_id
  }
  property {
    name  = "applicationsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_delete_id
  }
  property {
    name  = "applicationsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_get_id
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_accessValidatorId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_access_validator_id
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_agentCacheInvalidatedExpiration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_agent_cache_invalidated_expiration
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_agentCacheInvalidatedResponseDuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_agent_cache_invalidated_response_duration
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_agentId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_agent_id
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_allowEmptyPathSegments"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_allow_empty_path_segments
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_applicationType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_application_type
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_authenticationChallengePolicyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_authentication_challenge_policy_id
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_caseSensitivePath"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_case_sensitive_path
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_contextRoot"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_context_root
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_defaultAuthType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_default_auth_type
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_description
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_destination"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_destination
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_enabled
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_fallbackPostEncoding"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_fallback_post_encoding
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_id
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_identityMappingIds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_identity_mapping_ids
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_issuer"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_issuer
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_lastModified"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_last_modified
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_manualOrderingEnabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_manual_ordering_enabled
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_name
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_policy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_policy
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_realm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_realm
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_requireHTTPS"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_require_https
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_resourceOrder"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_resource_order
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_sidebandClientId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_sideband_client_id
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_siteId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_site_id
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_spaSupportEnabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_spa_support_enabled
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_virtualHostIds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_virtual_host_ids
  }
  property {
    name  = "applicationsIdPut_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_webSessionId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_applications_get200_response_items_inner_applications_get200_response_items_inner_web_session_id
  }
  property {
    name  = "applicationsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_put_id
  }
  property {
    name  = "applicationsIdResourceMatchingEvaluationOrderGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resource_matching_evaluation_order_get_id
  }
  property {
    name  = "applicationsIdResourcesAutoOrderGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_auto_order_get_id
  }
  property {
    name  = "applicationsIdResourcesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_get_filter
  }
  property {
    name  = "applicationsIdResourcesGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_get_id
  }
  property {
    name  = "applicationsIdResourcesGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_get_name
  }
  property {
    name  = "applicationsIdResourcesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_get_number_per_page
  }
  property {
    name  = "applicationsIdResourcesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_get_order
  }
  property {
    name  = "applicationsIdResourcesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_get_page
  }
  property {
    name  = "applicationsIdResourcesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_get_sort_key
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInnerQueryParamConfig_matchesNoParams"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_query_param_config_matches_no_params
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInnerQueryParamConfig_params"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_query_param_config_params
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInnerResourceTypeConfigurationResponseGenerator_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_resource_type_configuration_response_generator_class_name
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInnerResourceTypeConfigurationResponseGenerator_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_resource_type_configuration_response_generator_configuration
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_anonymous"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_anonymous
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_applicationId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_application_id
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_auditLevel"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_audit_level
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_authenticationChallengePolicyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_authentication_challenge_policy_id
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_defaultAuthTypeOverride"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_default_auth_type_override
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_enabled
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_id
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_methods"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_methods
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_name
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_pathPatterns"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_path_patterns
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_pathPrefixes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_path_prefixes
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_policy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_policy
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_resourceType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_resource_type
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_rootResource"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_root_resource
  }
  property {
    name  = "applicationsIdResourcesPost_applicationsResourcesGet200ResponseItemsInner_ApplicationsResourcesGet200ResponseItemsInner_unprotected"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_applications_resources_get200_response_items_inner_applications_resources_get200_response_items_inner_unprotected
  }
  property {
    name  = "applicationsIdResourcesPost_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_id_resources_post_id
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_accessValidatorId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_access_validator_id
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_agentCacheInvalidatedExpiration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_agent_cache_invalidated_expiration
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_agentCacheInvalidatedResponseDuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_agent_cache_invalidated_response_duration
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_agentId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_agent_id
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_allowEmptyPathSegments"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_allow_empty_path_segments
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_applicationType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_application_type
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_authenticationChallengePolicyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_authentication_challenge_policy_id
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_caseSensitivePath"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_case_sensitive_path
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_contextRoot"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_context_root
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_defaultAuthType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_default_auth_type
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_description
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_destination"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_destination
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_enabled
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_fallbackPostEncoding"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_fallback_post_encoding
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_id
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_identityMappingIds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_identity_mapping_ids
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_issuer"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_issuer
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_lastModified"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_last_modified
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_manualOrderingEnabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_manual_ordering_enabled
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_name
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_policy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_policy
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_realm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_realm
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_requireHTTPS"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_require_https
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_resourceOrder"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_resource_order
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_sidebandClientId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_sideband_client_id
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_siteId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_site_id
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_spaSupportEnabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_spa_support_enabled
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_virtualHostIds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_virtual_host_ids
  }
  property {
    name  = "applicationsPost_applicationsGet200ResponseItemsInner_ApplicationsGet200ResponseItemsInner_webSessionId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_post_applications_get200_response_items_inner_applications_get200_response_items_inner_web_session_id
  }
  property {
    name  = "applicationsReservedPut_applicationsReservedDelete200Response_ApplicationsReservedDelete200Response_contextRoot"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_reserved_put_applications_reserved_delete200_response_applications_reserved_delete200_response_context_root
  }
  property {
    name  = "applicationsResourcesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_resources_get_filter
  }
  property {
    name  = "applicationsResourcesGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_resources_get_name
  }
  property {
    name  = "applicationsResourcesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_resources_get_number_per_page
  }
  property {
    name  = "applicationsResourcesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_resources_get_order
  }
  property {
    name  = "applicationsResourcesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_resources_get_page
  }
  property {
    name  = "applicationsResourcesGet_rule_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_resources_get_rule_id
  }
  property {
    name  = "applicationsResourcesGet_ruleset_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_resources_get_ruleset_id
  }
  property {
    name  = "applicationsResourcesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_resources_get_sort_key
  }
  property {
    name  = "applicationsResourcesResponseGeneratorsDescriptorsResponseGeneratorTypeGet_response_generator_type"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_applications_resources_response_generators_descriptors_response_generator_type_get_response_generator_type
  }
  property {
    name  = "authBasicPut_authBasicDeleteRequest_AuthBasicDeleteRequest_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_basic_put_auth_basic_delete_request_auth_basic_delete_request_enabled
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseAccessTokenValidator_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_access_token_validator_class_name
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseAccessTokenValidator_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_access_token_validator_configuration
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseClientCredentialsClientSecret_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_client_credentials_client_secret_encrypted_value
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseClientCredentialsClientSecret_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_client_credentials_client_secret_value
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseClientCredentials_clientId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_client_credentials_client_id
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseClientCredentials_credentialsType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_client_credentials_credentials_type
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseClientCredentials_keyPairId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_client_credentials_key_pair_id
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseClientSecret_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_client_secret_encrypted_value
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseClientSecret_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_client_secret_value
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseRoleMappingAdministrator_attributes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_role_mapping_administrator_attributes
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseRoleMappingAuditor_attributes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_role_mapping_auditor_attributes
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseRoleMappingAuditor_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_role_mapping_auditor_enabled
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseRoleMappingPlatformAdmin_attributes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_role_mapping_platform_admin_attributes
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseRoleMappingPlatformAdmin_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_role_mapping_platform_admin_enabled
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200ResponseRoleMapping_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_role_mapping_enabled
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200Response_clientId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_client_id
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200Response_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_enabled
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200Response_scope"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_scope
  }
  property {
    name  = "authOauthPut_authOauthDelete200Response_AuthOauthDelete200Response_subjectAttributeName"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oauth_put_auth_oauth_delete200_response_auth_oauth_delete200_response_subject_attribute_name
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOauthDelete200ResponseClientCredentialsClientSecret_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oauth_delete200_response_client_credentials_client_secret_encrypted_value
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOauthDelete200ResponseClientCredentialsClientSecret_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oauth_delete200_response_client_credentials_client_secret_value
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOauthDelete200ResponseRoleMappingAdministrator_attributes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oauth_delete200_response_role_mapping_administrator_attributes
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOauthDelete200ResponseRoleMappingAuditor_attributes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oauth_delete200_response_role_mapping_auditor_attributes
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOauthDelete200ResponseRoleMappingAuditor_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oauth_delete200_response_role_mapping_auditor_enabled
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOauthDelete200ResponseRoleMappingPlatformAdmin_attributes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oauth_delete200_response_role_mapping_platform_admin_attributes
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOauthDelete200ResponseRoleMappingPlatformAdmin_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oauth_delete200_response_role_mapping_platform_admin_enabled
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOauthDelete200ResponseRoleMapping_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oauth_delete200_response_role_mapping_enabled
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200ResponseOidcConfigurationClientCredentials_clientId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_oidc_configuration_client_credentials_client_id
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200ResponseOidcConfigurationClientCredentials_credentialsType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_oidc_configuration_client_credentials_credentials_type
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200ResponseOidcConfigurationClientCredentials_keyPairId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_oidc_configuration_client_credentials_key_pair_id
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200ResponseOidcConfiguration_cacheUserAttributes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_oidc_configuration_cache_user_attributes
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200ResponseOidcConfiguration_enableRefreshUser"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_oidc_configuration_enable_refresh_user
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200ResponseOidcConfiguration_oidcLoginType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_oidc_configuration_oidc_login_type
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200ResponseOidcConfiguration_pfsessionStateCacheInSeconds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_oidc_configuration_pfsession_state_cache_in_seconds
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200ResponseOidcConfiguration_pkceChallengeType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_oidc_configuration_pkce_challenge_type
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200ResponseOidcConfiguration_refreshUserInfoClaimsInterval"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_oidc_configuration_refresh_user_info_claims_interval
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200ResponseOidcConfiguration_scopes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_oidc_configuration_scopes
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200ResponseOidcConfiguration_sendRequestedUrlToProvider"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_oidc_configuration_send_requested_url_to_provider
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200ResponseOidcConfiguration_validateSessionIsAlive"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_oidc_configuration_validate_session_is_alive
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200Response_authnReqListId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_authn_req_list_id
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200Response_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_enabled
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200Response_useSlo"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_use_slo
  }
  property {
    name  = "authOidcPut_authOidcDelete200Response_AuthOidcDelete200Response_usernameAttributeName"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_put_auth_oidc_delete200_response_auth_oidc_delete200_response_username_attribute_name
  }
  property {
    name  = "authOidcScopesGet_client_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_oidc_scopes_get_client_id
  }
  property {
    name  = "authPassword"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_password
  }
  property {
    name  = "authTokenManagementKeySetPut_authTokenManagementKeySetGet200Response_AuthTokenManagementKeySetGet200Response_keySet"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_token_management_key_set_put_auth_token_management_key_set_get200_response_auth_token_management_key_set_get200_response_key_set
  }
  property {
    name  = "authTokenManagementKeySetPut_authTokenManagementKeySetGet200Response_AuthTokenManagementKeySetGet200Response_nonce"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_token_management_key_set_put_auth_token_management_key_set_get200_response_auth_token_management_key_set_get200_response_nonce
  }
  property {
    name  = "authTokenManagementPut_authTokenManagementDelete200Response_AuthTokenManagementDelete200Response_issuer"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_token_management_put_auth_token_management_delete200_response_auth_token_management_delete200_response_issuer
  }
  property {
    name  = "authTokenManagementPut_authTokenManagementDelete200Response_AuthTokenManagementDelete200Response_keyRollEnabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_token_management_put_auth_token_management_delete200_response_auth_token_management_delete200_response_key_roll_enabled
  }
  property {
    name  = "authTokenManagementPut_authTokenManagementDelete200Response_AuthTokenManagementDelete200Response_keyRollPeriodInHours"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_token_management_put_auth_token_management_delete200_response_auth_token_management_delete200_response_key_roll_period_in_hours
  }
  property {
    name  = "authTokenManagementPut_authTokenManagementDelete200Response_AuthTokenManagementDelete200Response_signingAlgorithm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_token_management_put_auth_token_management_delete200_response_auth_token_management_delete200_response_signing_algorithm
  }
  property {
    name  = "authTokenProviderPut_authTokenProviderDelete200Response_AuthTokenProviderDelete200Response_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_token_provider_put_auth_token_provider_delete200_response_auth_token_provider_delete200_response_description
  }
  property {
    name  = "authTokenProviderPut_authTokenProviderDelete200Response_AuthTokenProviderDelete200Response_issuer"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_token_provider_put_auth_token_provider_delete200_response_auth_token_provider_delete200_response_issuer
  }
  property {
    name  = "authTokenProviderPut_authTokenProviderDelete200Response_AuthTokenProviderDelete200Response_sslCiphers"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_token_provider_put_auth_token_provider_delete200_response_auth_token_provider_delete200_response_ssl_ciphers
  }
  property {
    name  = "authTokenProviderPut_authTokenProviderDelete200Response_AuthTokenProviderDelete200Response_sslProtocols"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_token_provider_put_auth_token_provider_delete200_response_auth_token_provider_delete200_response_ssl_protocols
  }
  property {
    name  = "authTokenProviderPut_authTokenProviderDelete200Response_AuthTokenProviderDelete200Response_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_token_provider_put_auth_token_provider_delete200_response_auth_token_provider_delete200_response_trusted_certificate_group_id
  }
  property {
    name  = "authTokenProviderPut_authTokenProviderDelete200Response_AuthTokenProviderDelete200Response_useProxy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_token_provider_put_auth_token_provider_delete200_response_auth_token_provider_delete200_response_use_proxy
  }
  property {
    name  = "authUsername"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_username
  }
  property {
    name  = "authWebSessionPut_authWebSessionDelete200Response_AuthWebSessionDelete200Response_audience"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_web_session_put_auth_web_session_delete200_response_auth_web_session_delete200_response_audience
  }
  property {
    name  = "authWebSessionPut_authWebSessionDelete200Response_AuthWebSessionDelete200Response_cookieDomain"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_web_session_put_auth_web_session_delete200_response_auth_web_session_delete200_response_cookie_domain
  }
  property {
    name  = "authWebSessionPut_authWebSessionDelete200Response_AuthWebSessionDelete200Response_cookieType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_web_session_put_auth_web_session_delete200_response_auth_web_session_delete200_response_cookie_type
  }
  property {
    name  = "authWebSessionPut_authWebSessionDelete200Response_AuthWebSessionDelete200Response_expirationWarningInMinutes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_web_session_put_auth_web_session_delete200_response_auth_web_session_delete200_response_expiration_warning_in_minutes
  }
  property {
    name  = "authWebSessionPut_authWebSessionDelete200Response_AuthWebSessionDelete200Response_idleTimeoutInMinutes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_web_session_put_auth_web_session_delete200_response_auth_web_session_delete200_response_idle_timeout_in_minutes
  }
  property {
    name  = "authWebSessionPut_authWebSessionDelete200Response_AuthWebSessionDelete200Response_sessionPollIntervalInSeconds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_web_session_put_auth_web_session_delete200_response_auth_web_session_delete200_response_session_poll_interval_in_seconds
  }
  property {
    name  = "authWebSessionPut_authWebSessionDelete200Response_AuthWebSessionDelete200Response_sessionTimeoutInMinutes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_auth_web_session_put_auth_web_session_delete200_response_auth_web_session_delete200_response_session_timeout_in_minutes
  }
  property {
    name  = "authenticationChallengePoliciesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_get_filter
  }
  property {
    name  = "authenticationChallengePoliciesGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_get_name
  }
  property {
    name  = "authenticationChallengePoliciesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_get_number_per_page
  }
  property {
    name  = "authenticationChallengePoliciesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_get_order
  }
  property {
    name  = "authenticationChallengePoliciesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_get_page
  }
  property {
    name  = "authenticationChallengePoliciesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_get_sort_key
  }
  property {
    name  = "authenticationChallengePoliciesIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_id_delete_id
  }
  property {
    name  = "authenticationChallengePoliciesIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_id_get_id
  }
  property {
    name  = "authenticationChallengePoliciesIdPut_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInnerChallengeResponseChainInnerChallengeResponseFilter_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_id_put_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_challenge_response_chain_inner_challenge_response_filter_class_name
  }
  property {
    name  = "authenticationChallengePoliciesIdPut_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInnerChallengeResponseChainInnerChallengeResponseFilter_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_id_put_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_challenge_response_chain_inner_challenge_response_filter_configuration
  }
  property {
    name  = "authenticationChallengePoliciesIdPut_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInnerChallengeResponseChainInnerChallengeResponseGenerator_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_id_put_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_challenge_response_chain_inner_challenge_response_generator_class_name
  }
  property {
    name  = "authenticationChallengePoliciesIdPut_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInnerChallengeResponseChainInnerChallengeResponseGenerator_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_id_put_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_challenge_response_chain_inner_challenge_response_generator_configuration
  }
  property {
    name  = "authenticationChallengePoliciesIdPut_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInner_challengeResponseChain"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_id_put_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_challenge_response_chain
  }
  property {
    name  = "authenticationChallengePoliciesIdPut_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_id_put_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_description
  }
  property {
    name  = "authenticationChallengePoliciesIdPut_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_id_put_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_id
  }
  property {
    name  = "authenticationChallengePoliciesIdPut_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_id_put_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_name
  }
  property {
    name  = "authenticationChallengePoliciesIdPut_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInner_system"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_id_put_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_system
  }
  property {
    name  = "authenticationChallengePoliciesIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_id_put_id
  }
  property {
    name  = "authenticationChallengePoliciesPost_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInnerChallengeResponseChainInnerChallengeResponseFilter_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_post_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_challenge_response_chain_inner_challenge_response_filter_class_name
  }
  property {
    name  = "authenticationChallengePoliciesPost_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInnerChallengeResponseChainInnerChallengeResponseFilter_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_post_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_challenge_response_chain_inner_challenge_response_filter_configuration
  }
  property {
    name  = "authenticationChallengePoliciesPost_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInnerChallengeResponseChainInnerChallengeResponseGenerator_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_post_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_challenge_response_chain_inner_challenge_response_generator_class_name
  }
  property {
    name  = "authenticationChallengePoliciesPost_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInnerChallengeResponseChainInnerChallengeResponseGenerator_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_post_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_challenge_response_chain_inner_challenge_response_generator_configuration
  }
  property {
    name  = "authenticationChallengePoliciesPost_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInner_challengeResponseChain"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_post_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_challenge_response_chain
  }
  property {
    name  = "authenticationChallengePoliciesPost_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_post_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_description
  }
  property {
    name  = "authenticationChallengePoliciesPost_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_post_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_id
  }
  property {
    name  = "authenticationChallengePoliciesPost_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_post_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_name
  }
  property {
    name  = "authenticationChallengePoliciesPost_authenticationChallengePoliciesGet200ResponseItemsInner_AuthenticationChallengePoliciesGet200ResponseItemsInner_system"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_post_authentication_challenge_policies_get200_response_items_inner_authentication_challenge_policies_get200_response_items_inner_system
  }
  property {
    name  = "authenticationChallengePoliciesRequestMatchersDescriptorsRequestMatcherTypeGet_request_matcher_type"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_request_matchers_descriptors_request_matcher_type_get_request_matcher_type
  }
  property {
    name  = "authenticationChallengePoliciesResponseFiltersDescriptorsResponseFilterTypeGet_response_filter_type"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_response_filters_descriptors_response_filter_type_get_response_filter_type
  }
  property {
    name  = "authenticationChallengePoliciesResponseGeneratorsDescriptorsResponseGeneratorTypeGet_response_generator_type"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authentication_challenge_policies_response_generators_descriptors_response_generator_type_get_response_generator_type
  }
  property {
    name  = "authnReqListsGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_get_filter
  }
  property {
    name  = "authnReqListsGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_get_name
  }
  property {
    name  = "authnReqListsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_get_number_per_page
  }
  property {
    name  = "authnReqListsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_get_order
  }
  property {
    name  = "authnReqListsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_get_page
  }
  property {
    name  = "authnReqListsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_get_sort_key
  }
  property {
    name  = "authnReqListsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_id_delete_id
  }
  property {
    name  = "authnReqListsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_id_get_id
  }
  property {
    name  = "authnReqListsIdPut_authnReqListsGet200ResponseItemsInner_AuthnReqListsGet200ResponseItemsInner_authnReqs"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_id_put_authn_req_lists_get200_response_items_inner_authn_req_lists_get200_response_items_inner_authn_reqs
  }
  property {
    name  = "authnReqListsIdPut_authnReqListsGet200ResponseItemsInner_AuthnReqListsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_id_put_authn_req_lists_get200_response_items_inner_authn_req_lists_get200_response_items_inner_id
  }
  property {
    name  = "authnReqListsIdPut_authnReqListsGet200ResponseItemsInner_AuthnReqListsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_id_put_authn_req_lists_get200_response_items_inner_authn_req_lists_get200_response_items_inner_name
  }
  property {
    name  = "authnReqListsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_id_put_id
  }
  property {
    name  = "authnReqListsPost_authnReqListsGet200ResponseItemsInner_AuthnReqListsGet200ResponseItemsInner_authnReqs"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_post_authn_req_lists_get200_response_items_inner_authn_req_lists_get200_response_items_inner_authn_reqs
  }
  property {
    name  = "authnReqListsPost_authnReqListsGet200ResponseItemsInner_AuthnReqListsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_post_authn_req_lists_get200_response_items_inner_authn_req_lists_get200_response_items_inner_id
  }
  property {
    name  = "authnReqListsPost_authnReqListsGet200ResponseItemsInner_AuthnReqListsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_authn_req_lists_post_authn_req_lists_get200_response_items_inner_authn_req_lists_get200_response_items_inner_name
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_base_path
  }
  property {
    name  = "certificatesGet_alias"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_get_alias
  }
  property {
    name  = "certificatesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_get_filter
  }
  property {
    name  = "certificatesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_get_number_per_page
  }
  property {
    name  = "certificatesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_get_order
  }
  property {
    name  = "certificatesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_get_page
  }
  property {
    name  = "certificatesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_get_sort_key
  }
  property {
    name  = "certificatesIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_id_delete_id
  }
  property {
    name  = "certificatesIdFileGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_id_file_get_id
  }
  property {
    name  = "certificatesIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_id_get_id
  }
  property {
    name  = "certificatesIdPut_certificatesGetRequest_CertificatesGetRequest_alias"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_id_put_certificates_get_request_certificates_get_request_alias
  }
  property {
    name  = "certificatesIdPut_certificatesGetRequest_CertificatesGetRequest_fileData"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_id_put_certificates_get_request_certificates_get_request_file_data
  }
  property {
    name  = "certificatesIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_id_put_id
  }
  property {
    name  = "certificatesPost_certificatesGetRequest_CertificatesGetRequest_alias"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_post_certificates_get_request_certificates_get_request_alias
  }
  property {
    name  = "certificatesPost_certificatesGetRequest_CertificatesGetRequest_fileData"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_certificates_post_certificates_get_request_certificates_get_request_file_data
  }
  property {
    name  = "configExportWorkflowsIdDataGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_config_export_workflows_id_data_get_id
  }
  property {
    name  = "configExportWorkflowsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_config_export_workflows_id_get_id
  }
  property {
    name  = "configImportPost_body"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_config_import_post_body
  }
  property {
    name  = "configImportWorkflowsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_config_import_workflows_id_get_id
  }
  property {
    name  = "configImportWorkflowsPost_configExportGet200Response_ConfigExportGet200ResponseMasterKeys_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_config_import_workflows_post_config_export_get200_response_config_export_get200_response_master_keys_encrypted_value
  }
  property {
    name  = "configImportWorkflowsPost_configExportGet200Response_ConfigExportGet200ResponseMasterKeys_keyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_config_import_workflows_post_config_export_get200_response_config_export_get200_response_master_keys_key_id
  }
  property {
    name  = "configImportWorkflowsPost_configExportGet200Response_ConfigExportGet200Response_data"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_config_import_workflows_post_config_export_get200_response_config_export_get200_response_data
  }
  property {
    name  = "configImportWorkflowsPost_configExportGet200Response_ConfigExportGet200Response_encryptionKey"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_config_import_workflows_post_config_export_get200_response_config_export_get200_response_encryption_key
  }
  property {
    name  = "configImportWorkflowsPost_configExportGet200Response_ConfigExportGet200Response_version"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_config_import_workflows_post_config_export_get200_response_config_export_get200_response_version
  }
  property {
    name  = "configImportWorkflowsPost_fail_fast"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_config_import_workflows_post_fail_fast
  }
  property {
    name  = "defaultsEntitiesApplicationPut_defaultsEntitiesApplicationDelete200Response_DefaultsEntitiesApplicationDelete200Response_defaultAuthnChallengePolicyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_defaults_entities_application_put_defaults_entities_application_delete200_response_defaults_entities_application_delete200_response_default_authn_challenge_policy_id
  }
  property {
    name  = "engineListenersGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_get_filter
  }
  property {
    name  = "engineListenersGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_get_name
  }
  property {
    name  = "engineListenersGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_get_number_per_page
  }
  property {
    name  = "engineListenersGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_get_order
  }
  property {
    name  = "engineListenersGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_get_page
  }
  property {
    name  = "engineListenersGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_get_sort_key
  }
  property {
    name  = "engineListenersIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_id_delete_id
  }
  property {
    name  = "engineListenersIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_id_get_id
  }
  property {
    name  = "engineListenersIdPut_engineListenersGet200ResponseItemsInner_EngineListenersGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_id_put_engine_listeners_get200_response_items_inner_engine_listeners_get200_response_items_inner_id
  }
  property {
    name  = "engineListenersIdPut_engineListenersGet200ResponseItemsInner_EngineListenersGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_id_put_engine_listeners_get200_response_items_inner_engine_listeners_get200_response_items_inner_name
  }
  property {
    name  = "engineListenersIdPut_engineListenersGet200ResponseItemsInner_EngineListenersGet200ResponseItemsInner_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_id_put_engine_listeners_get200_response_items_inner_engine_listeners_get200_response_items_inner_port
  }
  property {
    name  = "engineListenersIdPut_engineListenersGet200ResponseItemsInner_EngineListenersGet200ResponseItemsInner_secure"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_id_put_engine_listeners_get200_response_items_inner_engine_listeners_get200_response_items_inner_secure
  }
  property {
    name  = "engineListenersIdPut_engineListenersGet200ResponseItemsInner_EngineListenersGet200ResponseItemsInner_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_id_put_engine_listeners_get200_response_items_inner_engine_listeners_get200_response_items_inner_trusted_certificate_group_id
  }
  property {
    name  = "engineListenersIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_id_put_id
  }
  property {
    name  = "engineListenersPost_engineListenersGet200ResponseItemsInner_EngineListenersGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_post_engine_listeners_get200_response_items_inner_engine_listeners_get200_response_items_inner_id
  }
  property {
    name  = "engineListenersPost_engineListenersGet200ResponseItemsInner_EngineListenersGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_post_engine_listeners_get200_response_items_inner_engine_listeners_get200_response_items_inner_name
  }
  property {
    name  = "engineListenersPost_engineListenersGet200ResponseItemsInner_EngineListenersGet200ResponseItemsInner_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_post_engine_listeners_get200_response_items_inner_engine_listeners_get200_response_items_inner_port
  }
  property {
    name  = "engineListenersPost_engineListenersGet200ResponseItemsInner_EngineListenersGet200ResponseItemsInner_secure"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_post_engine_listeners_get200_response_items_inner_engine_listeners_get200_response_items_inner_secure
  }
  property {
    name  = "engineListenersPost_engineListenersGet200ResponseItemsInner_EngineListenersGet200ResponseItemsInner_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engine_listeners_post_engine_listeners_get200_response_items_inner_engine_listeners_get200_response_items_inner_trusted_certificate_group_id
  }
  property {
    name  = "enginesCertificatesGet_alias"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_certificates_get_alias
  }
  property {
    name  = "enginesCertificatesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_certificates_get_filter
  }
  property {
    name  = "enginesCertificatesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_certificates_get_number_per_page
  }
  property {
    name  = "enginesCertificatesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_certificates_get_order
  }
  property {
    name  = "enginesCertificatesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_certificates_get_page
  }
  property {
    name  = "enginesCertificatesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_certificates_get_sort_key
  }
  property {
    name  = "enginesCertificatesIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_certificates_id_get_id
  }
  property {
    name  = "enginesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_get_filter
  }
  property {
    name  = "enginesGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_get_name
  }
  property {
    name  = "enginesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_get_number_per_page
  }
  property {
    name  = "enginesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_get_order
  }
  property {
    name  = "enginesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_get_page
  }
  property {
    name  = "enginesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_get_sort_key
  }
  property {
    name  = "enginesIdConfigPost_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_config_post_id
  }
  property {
    name  = "enginesIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_delete_id
  }
  property {
    name  = "enginesIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_get_id
  }
  property {
    name  = "enginesIdPut_enginesGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInnerCertificateHash_algorithm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_put_engines_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_certificate_hash_algorithm
  }
  property {
    name  = "enginesIdPut_enginesGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInnerCertificateHash_hexValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_put_engines_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_certificate_hash_hex_value
  }
  property {
    name  = "enginesIdPut_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_configReplicationEnabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_put_engines_get200_response_items_inner_engines_get200_response_items_inner_config_replication_enabled
  }
  property {
    name  = "enginesIdPut_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_put_engines_get200_response_items_inner_engines_get200_response_items_inner_description
  }
  property {
    name  = "enginesIdPut_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_httpProxyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_put_engines_get200_response_items_inner_engines_get200_response_items_inner_http_proxy_id
  }
  property {
    name  = "enginesIdPut_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_httpsProxyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_put_engines_get200_response_items_inner_engines_get200_response_items_inner_https_proxy_id
  }
  property {
    name  = "enginesIdPut_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_put_engines_get200_response_items_inner_engines_get200_response_items_inner_id
  }
  property {
    name  = "enginesIdPut_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_keys"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_put_engines_get200_response_items_inner_engines_get200_response_items_inner_keys
  }
  property {
    name  = "enginesIdPut_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_put_engines_get200_response_items_inner_engines_get200_response_items_inner_name
  }
  property {
    name  = "enginesIdPut_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_selectedCertificateId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_put_engines_get200_response_items_inner_engines_get200_response_items_inner_selected_certificate_id
  }
  property {
    name  = "enginesIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_id_put_id
  }
  property {
    name  = "enginesPost_enginesGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInnerCertificateHash_algorithm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_post_engines_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_certificate_hash_algorithm
  }
  property {
    name  = "enginesPost_enginesGet200ResponseItemsInner_AdminConfigReplicaAdminsGet200ResponseItemsInnerCertificateHash_hexValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_post_engines_get200_response_items_inner_admin_config_replica_admins_get200_response_items_inner_certificate_hash_hex_value
  }
  property {
    name  = "enginesPost_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_configReplicationEnabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_post_engines_get200_response_items_inner_engines_get200_response_items_inner_config_replication_enabled
  }
  property {
    name  = "enginesPost_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_post_engines_get200_response_items_inner_engines_get200_response_items_inner_description
  }
  property {
    name  = "enginesPost_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_httpProxyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_post_engines_get200_response_items_inner_engines_get200_response_items_inner_http_proxy_id
  }
  property {
    name  = "enginesPost_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_httpsProxyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_post_engines_get200_response_items_inner_engines_get200_response_items_inner_https_proxy_id
  }
  property {
    name  = "enginesPost_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_post_engines_get200_response_items_inner_engines_get200_response_items_inner_id
  }
  property {
    name  = "enginesPost_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_keys"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_post_engines_get200_response_items_inner_engines_get200_response_items_inner_keys
  }
  property {
    name  = "enginesPost_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_post_engines_get200_response_items_inner_engines_get200_response_items_inner_name
  }
  property {
    name  = "enginesPost_enginesGet200ResponseItemsInner_EnginesGet200ResponseItemsInner_selectedCertificateId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_post_engines_get200_response_items_inner_engines_get200_response_items_inner_selected_certificate_id
  }
  property {
    name  = "enginesRegistrationTokenPost_enginesRegistrationTokenPostRequest_EnginesRegistrationTokenPostRequest_expirationSeconds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_registration_token_post_engines_registration_token_post_request_engines_registration_token_post_request_expiration_seconds
  }
  property {
    name  = "enginesRegistrationTokenPost_enginesRegistrationTokenPostRequest_EnginesRegistrationTokenPostRequest_httpProxyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_registration_token_post_engines_registration_token_post_request_engines_registration_token_post_request_http_proxy_id
  }
  property {
    name  = "enginesRegistrationTokenPost_enginesRegistrationTokenPostRequest_EnginesRegistrationTokenPostRequest_httpsProxyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_registration_token_post_engines_registration_token_post_request_engines_registration_token_post_request_https_proxy_id
  }
  property {
    name  = "enginesRegistrationTokenPost_enginesRegistrationTokenPostRequest_EnginesRegistrationTokenPostRequest_selectedCertificateId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_engines_registration_token_post_engines_registration_token_post_request_engines_registration_token_post_request_selected_certificate_id
  }
  property {
    name  = "environmentPut_environmentDelete200Response_EnvironmentDelete200Response_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_environment_put_environment_delete200_response_environment_delete200_response_name
  }
  property {
    name  = "globalUnprotectedResourcesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_get_filter
  }
  property {
    name  = "globalUnprotectedResourcesGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_get_name
  }
  property {
    name  = "globalUnprotectedResourcesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_get_number_per_page
  }
  property {
    name  = "globalUnprotectedResourcesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_get_order
  }
  property {
    name  = "globalUnprotectedResourcesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_get_page
  }
  property {
    name  = "globalUnprotectedResourcesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_get_sort_key
  }
  property {
    name  = "globalUnprotectedResourcesIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_id_delete_id
  }
  property {
    name  = "globalUnprotectedResourcesIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_id_get_id
  }
  property {
    name  = "globalUnprotectedResourcesIdPut_globalUnprotectedResourcesGet200ResponseItemsInner_GlobalUnprotectedResourcesGet200ResponseItemsInner_auditLevel"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_id_put_global_unprotected_resources_get200_response_items_inner_global_unprotected_resources_get200_response_items_inner_audit_level
  }
  property {
    name  = "globalUnprotectedResourcesIdPut_globalUnprotectedResourcesGet200ResponseItemsInner_GlobalUnprotectedResourcesGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_id_put_global_unprotected_resources_get200_response_items_inner_global_unprotected_resources_get200_response_items_inner_description
  }
  property {
    name  = "globalUnprotectedResourcesIdPut_globalUnprotectedResourcesGet200ResponseItemsInner_GlobalUnprotectedResourcesGet200ResponseItemsInner_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_id_put_global_unprotected_resources_get200_response_items_inner_global_unprotected_resources_get200_response_items_inner_enabled
  }
  property {
    name  = "globalUnprotectedResourcesIdPut_globalUnprotectedResourcesGet200ResponseItemsInner_GlobalUnprotectedResourcesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_id_put_global_unprotected_resources_get200_response_items_inner_global_unprotected_resources_get200_response_items_inner_id
  }
  property {
    name  = "globalUnprotectedResourcesIdPut_globalUnprotectedResourcesGet200ResponseItemsInner_GlobalUnprotectedResourcesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_id_put_global_unprotected_resources_get200_response_items_inner_global_unprotected_resources_get200_response_items_inner_name
  }
  property {
    name  = "globalUnprotectedResourcesIdPut_globalUnprotectedResourcesGet200ResponseItemsInner_GlobalUnprotectedResourcesGet200ResponseItemsInner_wildcardPath"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_id_put_global_unprotected_resources_get200_response_items_inner_global_unprotected_resources_get200_response_items_inner_wildcard_path
  }
  property {
    name  = "globalUnprotectedResourcesIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_id_put_id
  }
  property {
    name  = "globalUnprotectedResourcesPost_globalUnprotectedResourcesGet200ResponseItemsInner_GlobalUnprotectedResourcesGet200ResponseItemsInner_auditLevel"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_post_global_unprotected_resources_get200_response_items_inner_global_unprotected_resources_get200_response_items_inner_audit_level
  }
  property {
    name  = "globalUnprotectedResourcesPost_globalUnprotectedResourcesGet200ResponseItemsInner_GlobalUnprotectedResourcesGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_post_global_unprotected_resources_get200_response_items_inner_global_unprotected_resources_get200_response_items_inner_description
  }
  property {
    name  = "globalUnprotectedResourcesPost_globalUnprotectedResourcesGet200ResponseItemsInner_GlobalUnprotectedResourcesGet200ResponseItemsInner_enabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_post_global_unprotected_resources_get200_response_items_inner_global_unprotected_resources_get200_response_items_inner_enabled
  }
  property {
    name  = "globalUnprotectedResourcesPost_globalUnprotectedResourcesGet200ResponseItemsInner_GlobalUnprotectedResourcesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_post_global_unprotected_resources_get200_response_items_inner_global_unprotected_resources_get200_response_items_inner_id
  }
  property {
    name  = "globalUnprotectedResourcesPost_globalUnprotectedResourcesGet200ResponseItemsInner_GlobalUnprotectedResourcesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_post_global_unprotected_resources_get200_response_items_inner_global_unprotected_resources_get200_response_items_inner_name
  }
  property {
    name  = "globalUnprotectedResourcesPost_globalUnprotectedResourcesGet200ResponseItemsInner_GlobalUnprotectedResourcesGet200ResponseItemsInner_wildcardPath"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_global_unprotected_resources_post_global_unprotected_resources_get200_response_items_inner_global_unprotected_resources_get200_response_items_inner_wildcard_path
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesDescriptorsAvailabilityProfileTypeGet_availability_profile_type"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_descriptors_availability_profile_type_get_availability_profile_type
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_get_filter
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_get_name
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_get_number_per_page
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_get_order
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_get_page
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_get_sort_key
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_id_delete_id
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_id_get_id
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesIdPut_highAvailabilityAvailabilityProfilesGet200ResponseItemsInner_HighAvailabilityAvailabilityProfilesGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_id_put_high_availability_availability_profiles_get200_response_items_inner_high_availability_availability_profiles_get200_response_items_inner_class_name
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesIdPut_highAvailabilityAvailabilityProfilesGet200ResponseItemsInner_HighAvailabilityAvailabilityProfilesGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_id_put_high_availability_availability_profiles_get200_response_items_inner_high_availability_availability_profiles_get200_response_items_inner_configuration
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesIdPut_highAvailabilityAvailabilityProfilesGet200ResponseItemsInner_HighAvailabilityAvailabilityProfilesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_id_put_high_availability_availability_profiles_get200_response_items_inner_high_availability_availability_profiles_get200_response_items_inner_id
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesIdPut_highAvailabilityAvailabilityProfilesGet200ResponseItemsInner_HighAvailabilityAvailabilityProfilesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_id_put_high_availability_availability_profiles_get200_response_items_inner_high_availability_availability_profiles_get200_response_items_inner_name
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_id_put_id
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesPost_highAvailabilityAvailabilityProfilesGet200ResponseItemsInner_HighAvailabilityAvailabilityProfilesGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_post_high_availability_availability_profiles_get200_response_items_inner_high_availability_availability_profiles_get200_response_items_inner_class_name
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesPost_highAvailabilityAvailabilityProfilesGet200ResponseItemsInner_HighAvailabilityAvailabilityProfilesGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_post_high_availability_availability_profiles_get200_response_items_inner_high_availability_availability_profiles_get200_response_items_inner_configuration
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesPost_highAvailabilityAvailabilityProfilesGet200ResponseItemsInner_HighAvailabilityAvailabilityProfilesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_post_high_availability_availability_profiles_get200_response_items_inner_high_availability_availability_profiles_get200_response_items_inner_id
  }
  property {
    name  = "highAvailabilityAvailabilityProfilesPost_highAvailabilityAvailabilityProfilesGet200ResponseItemsInner_HighAvailabilityAvailabilityProfilesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_availability_profiles_post_high_availability_availability_profiles_get200_response_items_inner_high_availability_availability_profiles_get200_response_items_inner_name
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesDescriptorsLoadBalancingStrategyTypeGet_load_balancing_strategy_type"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_descriptors_load_balancing_strategy_type_get_load_balancing_strategy_type
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_get_filter
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_get_name
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_get_number_per_page
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_get_order
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_get_page
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_get_sort_key
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_id_delete_id
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_id_get_id
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesIdPut_highAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_HighAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_id_put_high_availability_load_balancing_strategies_get200_response_items_inner_high_availability_load_balancing_strategies_get200_response_items_inner_class_name
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesIdPut_highAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_HighAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_id_put_high_availability_load_balancing_strategies_get200_response_items_inner_high_availability_load_balancing_strategies_get200_response_items_inner_configuration
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesIdPut_highAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_HighAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_id_put_high_availability_load_balancing_strategies_get200_response_items_inner_high_availability_load_balancing_strategies_get200_response_items_inner_id
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesIdPut_highAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_HighAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_id_put_high_availability_load_balancing_strategies_get200_response_items_inner_high_availability_load_balancing_strategies_get200_response_items_inner_name
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_id_put_id
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesPost_highAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_HighAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_post_high_availability_load_balancing_strategies_get200_response_items_inner_high_availability_load_balancing_strategies_get200_response_items_inner_class_name
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesPost_highAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_HighAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_post_high_availability_load_balancing_strategies_get200_response_items_inner_high_availability_load_balancing_strategies_get200_response_items_inner_configuration
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesPost_highAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_HighAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_post_high_availability_load_balancing_strategies_get200_response_items_inner_high_availability_load_balancing_strategies_get200_response_items_inner_id
  }
  property {
    name  = "highAvailabilityLoadBalancingStrategiesPost_highAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_HighAvailabilityLoadBalancingStrategiesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_high_availability_load_balancing_strategies_post_high_availability_load_balancing_strategies_get200_response_items_inner_high_availability_load_balancing_strategies_get200_response_items_inner_name
  }
  property {
    name  = "hsmProvidersGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_get_filter
  }
  property {
    name  = "hsmProvidersGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_get_name
  }
  property {
    name  = "hsmProvidersGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_get_number_per_page
  }
  property {
    name  = "hsmProvidersGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_get_order
  }
  property {
    name  = "hsmProvidersGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_get_page
  }
  property {
    name  = "hsmProvidersGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_get_sort_key
  }
  property {
    name  = "hsmProvidersIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_id_delete_id
  }
  property {
    name  = "hsmProvidersIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_id_get_id
  }
  property {
    name  = "hsmProvidersIdPut_hsmProvidersGet200Response_HsmProvidersGet200Response_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_id_put_hsm_providers_get200_response_hsm_providers_get200_response_class_name
  }
  property {
    name  = "hsmProvidersIdPut_hsmProvidersGet200Response_HsmProvidersGet200Response_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_id_put_hsm_providers_get200_response_hsm_providers_get200_response_configuration
  }
  property {
    name  = "hsmProvidersIdPut_hsmProvidersGet200Response_HsmProvidersGet200Response_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_id_put_hsm_providers_get200_response_hsm_providers_get200_response_id
  }
  property {
    name  = "hsmProvidersIdPut_hsmProvidersGet200Response_HsmProvidersGet200Response_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_id_put_hsm_providers_get200_response_hsm_providers_get200_response_name
  }
  property {
    name  = "hsmProvidersIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_id_put_id
  }
  property {
    name  = "hsmProvidersPost_hsmProvidersGet200Response_HsmProvidersGet200Response_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_post_hsm_providers_get200_response_hsm_providers_get200_response_class_name
  }
  property {
    name  = "hsmProvidersPost_hsmProvidersGet200Response_HsmProvidersGet200Response_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_post_hsm_providers_get200_response_hsm_providers_get200_response_configuration
  }
  property {
    name  = "hsmProvidersPost_hsmProvidersGet200Response_HsmProvidersGet200Response_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_post_hsm_providers_get200_response_hsm_providers_get200_response_id
  }
  property {
    name  = "hsmProvidersPost_hsmProvidersGet200Response_HsmProvidersGet200Response_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_hsm_providers_post_hsm_providers_get200_response_hsm_providers_get200_response_name
  }
  property {
    name  = "httpConfigMonitoringPut_httpConfigMonitoringDelete200Response_HttpConfigMonitoringDelete200Response_auditLevel"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_http_config_monitoring_put_http_config_monitoring_delete200_response_http_config_monitoring_delete200_response_audit_level
  }
  property {
    name  = "httpConfigRequestHostSourcePut_httpConfigRequestHostSourceDelete200Response_HttpConfigRequestHostSourceDelete200Response_headerNameList"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_http_config_request_host_source_put_http_config_request_host_source_delete200_response_http_config_request_host_source_delete200_response_header_name_list
  }
  property {
    name  = "httpConfigRequestHostSourcePut_httpConfigRequestHostSourceDelete200Response_HttpConfigRequestHostSourceDelete200Response_listValueLocation"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_http_config_request_host_source_put_http_config_request_host_source_delete200_response_http_config_request_host_source_delete200_response_list_value_location
  }
  property {
    name  = "httpConfigRequestIpSourcePut_httpConfigRequestIpSourceDelete200Response_HttpConfigRequestIpSourceDelete200Response_fallbackToLastHopIp"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_http_config_request_ip_source_put_http_config_request_ip_source_delete200_response_http_config_request_ip_source_delete200_response_fallback_to_last_hop_ip
  }
  property {
    name  = "httpConfigRequestIpSourcePut_httpConfigRequestIpSourceDelete200Response_HttpConfigRequestIpSourceDelete200Response_headerNameList"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_http_config_request_ip_source_put_http_config_request_ip_source_delete200_response_http_config_request_ip_source_delete200_response_header_name_list
  }
  property {
    name  = "httpConfigRequestIpSourcePut_httpConfigRequestIpSourceDelete200Response_HttpConfigRequestIpSourceDelete200Response_listValueLocation"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_http_config_request_ip_source_put_http_config_request_ip_source_delete200_response_http_config_request_ip_source_delete200_response_list_value_location
  }
  property {
    name  = "httpConfigRequestProtocolSourcePut_httpConfigRequestProtocolSourceDelete200Response_HttpConfigRequestProtocolSourceDelete200Response_headerName"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_http_config_request_protocol_source_put_http_config_request_protocol_source_delete200_response_http_config_request_protocol_source_delete200_response_header_name
  }
  property {
    name  = "httpsListenersGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_https_listeners_get_order
  }
  property {
    name  = "httpsListenersGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_https_listeners_get_sort_key
  }
  property {
    name  = "httpsListenersIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_https_listeners_id_get_id
  }
  property {
    name  = "httpsListenersIdPut_httpsListenersGet200ResponseItemsInner_HttpsListenersGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_https_listeners_id_put_https_listeners_get200_response_items_inner_https_listeners_get200_response_items_inner_id
  }
  property {
    name  = "httpsListenersIdPut_httpsListenersGet200ResponseItemsInner_HttpsListenersGet200ResponseItemsInner_keyPairId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_https_listeners_id_put_https_listeners_get200_response_items_inner_https_listeners_get200_response_items_inner_key_pair_id
  }
  property {
    name  = "httpsListenersIdPut_httpsListenersGet200ResponseItemsInner_HttpsListenersGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_https_listeners_id_put_https_listeners_get200_response_items_inner_https_listeners_get200_response_items_inner_name
  }
  property {
    name  = "httpsListenersIdPut_httpsListenersGet200ResponseItemsInner_HttpsListenersGet200ResponseItemsInner_restartRequired"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_https_listeners_id_put_https_listeners_get200_response_items_inner_https_listeners_get200_response_items_inner_restart_required
  }
  property {
    name  = "httpsListenersIdPut_httpsListenersGet200ResponseItemsInner_HttpsListenersGet200ResponseItemsInner_useServerCipherSuiteOrder"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_https_listeners_id_put_https_listeners_get200_response_items_inner_https_listeners_get200_response_items_inner_use_server_cipher_suite_order
  }
  property {
    name  = "httpsListenersIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_https_listeners_id_put_id
  }
  property {
    name  = "identityMappingsDescriptorsIdentityMappingTypeGet_identity_mapping_type"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_descriptors_identity_mapping_type_get_identity_mapping_type
  }
  property {
    name  = "identityMappingsGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_get_filter
  }
  property {
    name  = "identityMappingsGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_get_name
  }
  property {
    name  = "identityMappingsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_get_number_per_page
  }
  property {
    name  = "identityMappingsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_get_order
  }
  property {
    name  = "identityMappingsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_get_page
  }
  property {
    name  = "identityMappingsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_get_sort_key
  }
  property {
    name  = "identityMappingsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_id_delete_id
  }
  property {
    name  = "identityMappingsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_id_get_id
  }
  property {
    name  = "identityMappingsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_id_put_id
  }
  property {
    name  = "identityMappingsIdPut_identityMappingsGet200ResponseItemsInner_IdentityMappingsGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_id_put_identity_mappings_get200_response_items_inner_identity_mappings_get200_response_items_inner_class_name
  }
  property {
    name  = "identityMappingsIdPut_identityMappingsGet200ResponseItemsInner_IdentityMappingsGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_id_put_identity_mappings_get200_response_items_inner_identity_mappings_get200_response_items_inner_configuration
  }
  property {
    name  = "identityMappingsIdPut_identityMappingsGet200ResponseItemsInner_IdentityMappingsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_id_put_identity_mappings_get200_response_items_inner_identity_mappings_get200_response_items_inner_id
  }
  property {
    name  = "identityMappingsIdPut_identityMappingsGet200ResponseItemsInner_IdentityMappingsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_id_put_identity_mappings_get200_response_items_inner_identity_mappings_get200_response_items_inner_name
  }
  property {
    name  = "identityMappingsPost_identityMappingsGet200ResponseItemsInner_IdentityMappingsGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_post_identity_mappings_get200_response_items_inner_identity_mappings_get200_response_items_inner_class_name
  }
  property {
    name  = "identityMappingsPost_identityMappingsGet200ResponseItemsInner_IdentityMappingsGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_post_identity_mappings_get200_response_items_inner_identity_mappings_get200_response_items_inner_configuration
  }
  property {
    name  = "identityMappingsPost_identityMappingsGet200ResponseItemsInner_IdentityMappingsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_post_identity_mappings_get200_response_items_inner_identity_mappings_get200_response_items_inner_id
  }
  property {
    name  = "identityMappingsPost_identityMappingsGet200ResponseItemsInner_IdentityMappingsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_identity_mappings_post_identity_mappings_get200_response_items_inner_identity_mappings_get200_response_items_inner_name
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_alias"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_alias
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_city"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_city
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_commonName"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_common_name
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_country"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_country
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_hsmProviderId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_hsm_provider_id
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_id
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_keyAlgorithm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_key_algorithm
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_keySize"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_key_size
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_organization"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_organization
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_organizationUnit"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_organization_unit
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_signatureAlgorithm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_signature_algorithm
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_state"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_state
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_subjectAlternativeNames"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_subject_alternative_names
  }
  property {
    name  = "keyPairsGeneratePost_keyPairsGeneratePostRequest_KeyPairsGeneratePostRequest_validDays"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_generate_post_key_pairs_generate_post_request_key_pairs_generate_post_request_valid_days
  }
  property {
    name  = "keyPairsGet_alias"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_get_alias
  }
  property {
    name  = "keyPairsGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_get_filter
  }
  property {
    name  = "keyPairsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_get_number_per_page
  }
  property {
    name  = "keyPairsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_get_order
  }
  property {
    name  = "keyPairsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_get_page
  }
  property {
    name  = "keyPairsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_get_sort_key
  }
  property {
    name  = "keyPairsIdCertificateGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_certificate_get_id
  }
  property {
    name  = "keyPairsIdCsrGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_csr_get_id
  }
  property {
    name  = "keyPairsIdCsrPost_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_csr_post_id
  }
  property {
    name  = "keyPairsIdCsrPost_keyPairsIdCsrGetRequest_KeyPairsIdCsrGetRequest_chainCertificates"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_csr_post_key_pairs_id_csr_get_request_key_pairs_id_csr_get_request_chain_certificates
  }
  property {
    name  = "keyPairsIdCsrPost_keyPairsIdCsrGetRequest_KeyPairsIdCsrGetRequest_fileData"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_csr_post_key_pairs_id_csr_get_request_key_pairs_id_csr_get_request_file_data
  }
  property {
    name  = "keyPairsIdCsrPost_keyPairsIdCsrGetRequest_KeyPairsIdCsrGetRequest_hsmProviderId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_csr_post_key_pairs_id_csr_get_request_key_pairs_id_csr_get_request_hsm_provider_id
  }
  property {
    name  = "keyPairsIdCsrPost_keyPairsIdCsrGetRequest_KeyPairsIdCsrGetRequest_trustedCertGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_csr_post_key_pairs_id_csr_get_request_key_pairs_id_csr_get_request_trusted_cert_group_id
  }
  property {
    name  = "keyPairsIdCsrPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_csr_put_id
  }
  property {
    name  = "keyPairsIdCsrPut_keyPairsIdCsrGetRequest_KeyPairsIdCsrGetRequest_chainCertificates"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_csr_put_key_pairs_id_csr_get_request_key_pairs_id_csr_get_request_chain_certificates
  }
  property {
    name  = "keyPairsIdCsrPut_keyPairsIdCsrGetRequest_KeyPairsIdCsrGetRequest_fileData"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_csr_put_key_pairs_id_csr_get_request_key_pairs_id_csr_get_request_file_data
  }
  property {
    name  = "keyPairsIdCsrPut_keyPairsIdCsrGetRequest_KeyPairsIdCsrGetRequest_hsmProviderId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_csr_put_key_pairs_id_csr_get_request_key_pairs_id_csr_get_request_hsm_provider_id
  }
  property {
    name  = "keyPairsIdCsrPut_keyPairsIdCsrGetRequest_KeyPairsIdCsrGetRequest_trustedCertGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_csr_put_key_pairs_id_csr_get_request_key_pairs_id_csr_get_request_trusted_cert_group_id
  }
  property {
    name  = "keyPairsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_delete_id
  }
  property {
    name  = "keyPairsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_get_id
  }
  property {
    name  = "keyPairsIdPatch_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_patch_id
  }
  property {
    name  = "keyPairsIdPatch_keyPairsIdDeleteRequest_KeyPairsIdDeleteRequest_addChainCertificates"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_patch_key_pairs_id_delete_request_key_pairs_id_delete_request_add_chain_certificates
  }
  property {
    name  = "keyPairsIdPemPost_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_pem_post_id
  }
  property {
    name  = "keyPairsIdPemPost_keyPairsIdPemPostRequest_KeyPairsIdPemPostRequestPassword_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_pem_post_key_pairs_id_pem_post_request_key_pairs_id_pem_post_request_password_encrypted_value
  }
  property {
    name  = "keyPairsIdPemPost_keyPairsIdPemPostRequest_KeyPairsIdPemPostRequestPassword_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_pem_post_key_pairs_id_pem_post_request_key_pairs_id_pem_post_request_password_value
  }
  property {
    name  = "keyPairsIdPkcs12Post_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_pkcs12_post_id
  }
  property {
    name  = "keyPairsIdPkcs12Post_keyPairsIdPemPostRequest_KeyPairsIdPemPostRequestPassword_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_pkcs12_post_key_pairs_id_pem_post_request_key_pairs_id_pem_post_request_password_encrypted_value
  }
  property {
    name  = "keyPairsIdPkcs12Post_keyPairsIdPemPostRequest_KeyPairsIdPemPostRequestPassword_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_pkcs12_post_key_pairs_id_pem_post_request_key_pairs_id_pem_post_request_password_value
  }
  property {
    name  = "keyPairsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_put_id
  }
  property {
    name  = "keyPairsIdPut_keyPairsImportPostRequest_KeyPairsImportPostRequestPassword_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_put_key_pairs_import_post_request_key_pairs_import_post_request_password_encrypted_value
  }
  property {
    name  = "keyPairsIdPut_keyPairsImportPostRequest_KeyPairsImportPostRequestPassword_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_put_key_pairs_import_post_request_key_pairs_import_post_request_password_value
  }
  property {
    name  = "keyPairsIdPut_keyPairsImportPostRequest_KeyPairsImportPostRequest_alias"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_put_key_pairs_import_post_request_key_pairs_import_post_request_alias
  }
  property {
    name  = "keyPairsIdPut_keyPairsImportPostRequest_KeyPairsImportPostRequest_chainCertificates"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_put_key_pairs_import_post_request_key_pairs_import_post_request_chain_certificates
  }
  property {
    name  = "keyPairsIdPut_keyPairsImportPostRequest_KeyPairsImportPostRequest_fileData"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_put_key_pairs_import_post_request_key_pairs_import_post_request_file_data
  }
  property {
    name  = "keyPairsIdPut_keyPairsImportPostRequest_KeyPairsImportPostRequest_hsmProviderId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_id_put_key_pairs_import_post_request_key_pairs_import_post_request_hsm_provider_id
  }
  property {
    name  = "keyPairsImportPost_keyPairsImportPostRequest_KeyPairsImportPostRequestPassword_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_import_post_key_pairs_import_post_request_key_pairs_import_post_request_password_encrypted_value
  }
  property {
    name  = "keyPairsImportPost_keyPairsImportPostRequest_KeyPairsImportPostRequestPassword_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_import_post_key_pairs_import_post_request_key_pairs_import_post_request_password_value
  }
  property {
    name  = "keyPairsImportPost_keyPairsImportPostRequest_KeyPairsImportPostRequest_alias"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_import_post_key_pairs_import_post_request_key_pairs_import_post_request_alias
  }
  property {
    name  = "keyPairsImportPost_keyPairsImportPostRequest_KeyPairsImportPostRequest_chainCertificates"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_import_post_key_pairs_import_post_request_key_pairs_import_post_request_chain_certificates
  }
  property {
    name  = "keyPairsImportPost_keyPairsImportPostRequest_KeyPairsImportPostRequest_fileData"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_import_post_key_pairs_import_post_request_key_pairs_import_post_request_file_data
  }
  property {
    name  = "keyPairsImportPost_keyPairsImportPostRequest_KeyPairsImportPostRequest_hsmProviderId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_import_post_key_pairs_import_post_request_key_pairs_import_post_request_hsm_provider_id
  }
  property {
    name  = "keyPairsKeyPairIdChainCertificatesChainCertificateIdDelete_chain_certificate_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_key_pair_id_chain_certificates_chain_certificate_id_delete_chain_certificate_id
  }
  property {
    name  = "keyPairsKeyPairIdChainCertificatesChainCertificateIdDelete_key_pair_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_key_pairs_key_pair_id_chain_certificates_chain_certificate_id_delete_key_pair_id
  }
  property {
    name  = "licensePost_licenseGetRequest_LicenseGetRequest_fileData"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_license_post_license_get_request_license_get_request_file_data
  }
  property {
    name  = "licensePut_licenseGetRequest_LicenseGetRequest_fileData"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_license_put_license_get_request_license_get_request_file_data
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_AuthOauthDelete200ResponseClientCredentialsClientSecret_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_auth_oauth_delete200_response_client_credentials_client_secret_encrypted_value
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_AuthOauthDelete200ResponseClientCredentialsClientSecret_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_auth_oauth_delete200_response_client_credentials_client_secret_value
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_AuthOidcDelete200ResponseOidcConfigurationClientCredentials_clientId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_auth_oidc_delete200_response_oidc_configuration_client_credentials_client_id
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_AuthOidcDelete200ResponseOidcConfigurationClientCredentials_credentialsType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_auth_oidc_delete200_response_oidc_configuration_client_credentials_credentials_type
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_AuthOidcDelete200ResponseOidcConfigurationClientCredentials_keyPairId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_auth_oidc_delete200_response_oidc_configuration_client_credentials_key_pair_id
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_OauthAuthServerDelete200Response_auditLevel"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_oauth_auth_server_delete200_response_audit_level
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_OauthAuthServerDelete200Response_cacheTokens"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_oauth_auth_server_delete200_response_cache_tokens
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_OauthAuthServerDelete200Response_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_oauth_auth_server_delete200_response_description
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_OauthAuthServerDelete200Response_introspectionEndpoint"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_oauth_auth_server_delete200_response_introspection_endpoint
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_OauthAuthServerDelete200Response_secure"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_oauth_auth_server_delete200_response_secure
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_OauthAuthServerDelete200Response_sendAudience"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_oauth_auth_server_delete200_response_send_audience
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_OauthAuthServerDelete200Response_subjectAttributeName"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_oauth_auth_server_delete200_response_subject_attribute_name
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_OauthAuthServerDelete200Response_targets"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_oauth_auth_server_delete200_response_targets
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_OauthAuthServerDelete200Response_tokenEndpoint"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_oauth_auth_server_delete200_response_token_endpoint
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_OauthAuthServerDelete200Response_tokenTimeToLiveSeconds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_oauth_auth_server_delete200_response_token_time_to_live_seconds
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_OauthAuthServerDelete200Response_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_oauth_auth_server_delete200_response_trusted_certificate_group_id
  }
  property {
    name  = "oauthAuthServerPut_oauthAuthServerDelete200Response_OauthAuthServerDelete200Response_useProxy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_auth_server_put_oauth_auth_server_delete200_response_oauth_auth_server_delete200_response_use_proxy
  }
  property {
    name  = "oauthKeyManagementKeySetPut_authTokenManagementKeySetGet200Response_AuthTokenManagementKeySetGet200Response_keySet"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_key_management_key_set_put_auth_token_management_key_set_get200_response_auth_token_management_key_set_get200_response_key_set
  }
  property {
    name  = "oauthKeyManagementKeySetPut_authTokenManagementKeySetGet200Response_AuthTokenManagementKeySetGet200Response_nonce"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_key_management_key_set_put_auth_token_management_key_set_get200_response_auth_token_management_key_set_get200_response_nonce
  }
  property {
    name  = "oauthKeyManagementPut_oauthKeyManagementDelete200Response_OauthKeyManagementDelete200Response_keyRollEnabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_key_management_put_oauth_key_management_delete200_response_oauth_key_management_delete200_response_key_roll_enabled
  }
  property {
    name  = "oauthKeyManagementPut_oauthKeyManagementDelete200Response_OauthKeyManagementDelete200Response_keyRollPeriodInHours"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_key_management_put_oauth_key_management_delete200_response_oauth_key_management_delete200_response_key_roll_period_in_hours
  }
  property {
    name  = "oauthKeyManagementPut_oauthKeyManagementDelete200Response_OauthKeyManagementDelete200Response_signingAlgorithm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oauth_key_management_put_oauth_key_management_delete200_response_oauth_key_management_delete200_response_signing_algorithm
  }
  property {
    name  = "oidcProviderDescriptorsPluginTypeGet_plugin_type"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oidc_provider_descriptors_plugin_type_get_plugin_type
  }
  property {
    name  = "oidcProviderPut_oidcProviderDelete200Response_OidcProviderDelete200ResponsePlugin_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oidc_provider_put_oidc_provider_delete200_response_oidc_provider_delete200_response_plugin_class_name
  }
  property {
    name  = "oidcProviderPut_oidcProviderDelete200Response_OidcProviderDelete200ResponsePlugin_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oidc_provider_put_oidc_provider_delete200_response_oidc_provider_delete200_response_plugin_configuration
  }
  property {
    name  = "oidcProviderPut_oidcProviderDelete200Response_OidcProviderDelete200Response_auditLevel"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oidc_provider_put_oidc_provider_delete200_response_oidc_provider_delete200_response_audit_level
  }
  property {
    name  = "oidcProviderPut_oidcProviderDelete200Response_OidcProviderDelete200Response_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oidc_provider_put_oidc_provider_delete200_response_oidc_provider_delete200_response_description
  }
  property {
    name  = "oidcProviderPut_oidcProviderDelete200Response_OidcProviderDelete200Response_issuer"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oidc_provider_put_oidc_provider_delete200_response_oidc_provider_delete200_response_issuer
  }
  property {
    name  = "oidcProviderPut_oidcProviderDelete200Response_OidcProviderDelete200Response_queryParameters"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oidc_provider_put_oidc_provider_delete200_response_oidc_provider_delete200_response_query_parameters
  }
  property {
    name  = "oidcProviderPut_oidcProviderDelete200Response_OidcProviderDelete200Response_requestSupportedScopesOnly"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oidc_provider_put_oidc_provider_delete200_response_oidc_provider_delete200_response_request_supported_scopes_only
  }
  property {
    name  = "oidcProviderPut_oidcProviderDelete200Response_OidcProviderDelete200Response_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oidc_provider_put_oidc_provider_delete200_response_oidc_provider_delete200_response_trusted_certificate_group_id
  }
  property {
    name  = "oidcProviderPut_oidcProviderDelete200Response_OidcProviderDelete200Response_useProxy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oidc_provider_put_oidc_provider_delete200_response_oidc_provider_delete200_response_use_proxy
  }
  property {
    name  = "oidcProviderPut_oidcProviderDelete200Response_OidcProviderDelete200Response_useSlo"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_oidc_provider_put_oidc_provider_delete200_response_oidc_provider_delete200_response_use_slo
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_AuthOauthDelete200ResponseClientCredentialsClientSecret_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_auth_oauth_delete200_response_client_credentials_client_secret_encrypted_value
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_AuthOauthDelete200ResponseClientCredentialsClientSecret_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_auth_oauth_delete200_response_client_credentials_client_secret_value
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_PingfederateAccessTokensDelete200ResponseClientCredentials_clientId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_pingfederate_access_tokens_delete200_response_client_credentials_client_id
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_PingfederateAccessTokensDelete200ResponseClientCredentials_credentialsType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_pingfederate_access_tokens_delete200_response_client_credentials_credentials_type
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_PingfederateAccessTokensDelete200ResponseClientCredentials_keyPairId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_pingfederate_access_tokens_delete200_response_client_credentials_key_pair_id
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_PingfederateAccessTokensDelete200ResponseClientSecret_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_pingfederate_access_tokens_delete200_response_client_secret_encrypted_value
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_PingfederateAccessTokensDelete200ResponseClientSecret_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_pingfederate_access_tokens_delete200_response_client_secret_value
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_PingfederateAccessTokensDelete200Response_accessValidatorId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_pingfederate_access_tokens_delete200_response_access_validator_id
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_PingfederateAccessTokensDelete200Response_cacheTokens"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_pingfederate_access_tokens_delete200_response_cache_tokens
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_PingfederateAccessTokensDelete200Response_clientId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_pingfederate_access_tokens_delete200_response_client_id
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_PingfederateAccessTokensDelete200Response_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_pingfederate_access_tokens_delete200_response_name
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_PingfederateAccessTokensDelete200Response_sendAudience"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_pingfederate_access_tokens_delete200_response_send_audience
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_PingfederateAccessTokensDelete200Response_subjectAttributeName"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_pingfederate_access_tokens_delete200_response_subject_attribute_name
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_PingfederateAccessTokensDelete200Response_tokenTimeToLiveSeconds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_pingfederate_access_tokens_delete200_response_token_time_to_live_seconds
  }
  property {
    name  = "pingfederateAccessTokensPut_pingfederateAccessTokensDelete200Response_PingfederateAccessTokensDelete200Response_useTokenIntrospection"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_access_tokens_put_pingfederate_access_tokens_delete200_response_pingfederate_access_tokens_delete200_response_use_token_introspection
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_AuthOauthDelete200ResponseClientCredentialsClientSecret_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_auth_oauth_delete200_response_client_credentials_client_secret_encrypted_value
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_AuthOauthDelete200ResponseClientCredentialsClientSecret_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_auth_oauth_delete200_response_client_credentials_client_secret_value
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200ResponseAdminPassword_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_admin_password_encrypted_value
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200ResponseAdminPassword_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_admin_password_value
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200ResponseOAuthAuthenticationConfigClientCredentials_clientId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_oauth_authentication_config_client_credentials_client_id
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200ResponseOAuthAuthenticationConfigClientCredentials_credentialsType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_oauth_authentication_config_client_credentials_credentials_type
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200ResponseOAuthAuthenticationConfigClientCredentials_keyPairId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_oauth_authentication_config_client_credentials_key_pair_id
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200ResponseOAuthAuthenticationConfig_configuredAuthorizationServerType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_oauth_authentication_config_configured_authorization_server_type
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200ResponseOAuthAuthenticationConfig_scopes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_oauth_authentication_config_scopes
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200Response_adminUsername"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_admin_username
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200Response_auditLevel"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_audit_level
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200Response_authenticationType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_authentication_type
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200Response_basePath"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_base_path
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200Response_expectedHostname"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_expected_hostname
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200Response_host"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_host
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200Response_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_port
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200Response_secure"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_secure
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200Response_skipHostnameVerification"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_skip_hostname_verification
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200Response_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_trusted_certificate_group_id
  }
  property {
    name  = "pingfederateAdminPut_pingfederateAdminDelete200Response_PingfederateAdminDelete200Response_useProxy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_admin_put_pingfederate_admin_delete200_response_pingfederate_admin_delete200_response_use_proxy
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200ResponseApplication_additionalVirtualHostIds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_application_additional_virtual_host_ids
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200ResponseApplication_caseSensitive"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_application_case_sensitive
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200ResponseApplication_clientCertHeaderNames"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_application_client_cert_header_names
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200ResponseApplication_contextRoot"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_application_context_root
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200ResponseApplication_policy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_application_policy
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200ResponseApplication_primaryVirtualHostId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_application_primary_virtual_host_id
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_auditLevel"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_audit_level
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_availabilityProfileId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_availability_profile_id
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_backChannelBasePath"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_back_channel_base_path
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_backChannelSecure"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_back_channel_secure
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_basePath"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_base_path
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_expectedHostname"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_expected_hostname
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_host"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_host
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_loadBalancingStrategyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_load_balancing_strategy_id
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_port
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_secure"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_secure
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_skipHostnameVerification"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_skip_hostname_verification
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_targets"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_targets
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_trusted_certificate_group_id
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_useProxy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_use_proxy
  }
  property {
    name  = "pingfederatePut_pingfederateDelete200Response_PingfederateDelete200Response_useSlo"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_put_pingfederate_delete200_response_pingfederate_delete200_response_use_slo
  }
  property {
    name  = "pingfederateRuntimePut_pingfederateRuntimeDelete200Response_PingfederateRuntimeDelete200Response_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_runtime_put_pingfederate_runtime_delete200_response_pingfederate_runtime_delete200_response_description
  }
  property {
    name  = "pingfederateRuntimePut_pingfederateRuntimeDelete200Response_PingfederateRuntimeDelete200Response_issuer"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_runtime_put_pingfederate_runtime_delete200_response_pingfederate_runtime_delete200_response_issuer
  }
  property {
    name  = "pingfederateRuntimePut_pingfederateRuntimeDelete200Response_PingfederateRuntimeDelete200Response_skipHostnameVerification"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_runtime_put_pingfederate_runtime_delete200_response_pingfederate_runtime_delete200_response_skip_hostname_verification
  }
  property {
    name  = "pingfederateRuntimePut_pingfederateRuntimeDelete200Response_PingfederateRuntimeDelete200Response_stsTokenExchangeEndpoint"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_runtime_put_pingfederate_runtime_delete200_response_pingfederate_runtime_delete200_response_sts_token_exchange_endpoint
  }
  property {
    name  = "pingfederateRuntimePut_pingfederateRuntimeDelete200Response_PingfederateRuntimeDelete200Response_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_runtime_put_pingfederate_runtime_delete200_response_pingfederate_runtime_delete200_response_trusted_certificate_group_id
  }
  property {
    name  = "pingfederateRuntimePut_pingfederateRuntimeDelete200Response_PingfederateRuntimeDelete200Response_useProxy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_runtime_put_pingfederate_runtime_delete200_response_pingfederate_runtime_delete200_response_use_proxy
  }
  property {
    name  = "pingfederateRuntimePut_pingfederateRuntimeDelete200Response_PingfederateRuntimeDelete200Response_useSlo"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingfederate_runtime_put_pingfederate_runtime_delete200_response_pingfederate_runtime_delete200_response_use_slo
  }
  property {
    name  = "pingoneCustomersPut_pingoneCustomersDelete200Response_PingoneCustomersDelete200Response_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingone_customers_put_pingone_customers_delete200_response_pingone_customers_delete200_response_description
  }
  property {
    name  = "pingoneCustomersPut_pingoneCustomersDelete200Response_PingoneCustomersDelete200Response_issuer"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingone_customers_put_pingone_customers_delete200_response_pingone_customers_delete200_response_issuer
  }
  property {
    name  = "pingoneCustomersPut_pingoneCustomersDelete200Response_PingoneCustomersDelete200Response_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingone_customers_put_pingone_customers_delete200_response_pingone_customers_delete200_response_trusted_certificate_group_id
  }
  property {
    name  = "pingoneCustomersPut_pingoneCustomersDelete200Response_PingoneCustomersDelete200Response_useProxy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_pingone_customers_put_pingone_customers_delete200_response_pingone_customers_delete200_response_use_proxy
  }
  property {
    name  = "proxiesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_get_filter
  }
  property {
    name  = "proxiesGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_get_name
  }
  property {
    name  = "proxiesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_get_number_per_page
  }
  property {
    name  = "proxiesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_get_order
  }
  property {
    name  = "proxiesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_get_page
  }
  property {
    name  = "proxiesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_get_sort_key
  }
  property {
    name  = "proxiesIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_id_delete_id
  }
  property {
    name  = "proxiesIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_id_get_id
  }
  property {
    name  = "proxiesIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_id_put_id
  }
  property {
    name  = "proxiesIdPut_proxiesGet200Response_ProxiesGet200ResponsePassword_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_id_put_proxies_get200_response_proxies_get200_response_password_encrypted_value
  }
  property {
    name  = "proxiesIdPut_proxiesGet200Response_ProxiesGet200ResponsePassword_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_id_put_proxies_get200_response_proxies_get200_response_password_value
  }
  property {
    name  = "proxiesIdPut_proxiesGet200Response_ProxiesGet200Response_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_id_put_proxies_get200_response_proxies_get200_response_description
  }
  property {
    name  = "proxiesIdPut_proxiesGet200Response_ProxiesGet200Response_host"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_id_put_proxies_get200_response_proxies_get200_response_host
  }
  property {
    name  = "proxiesIdPut_proxiesGet200Response_ProxiesGet200Response_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_id_put_proxies_get200_response_proxies_get200_response_id
  }
  property {
    name  = "proxiesIdPut_proxiesGet200Response_ProxiesGet200Response_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_id_put_proxies_get200_response_proxies_get200_response_name
  }
  property {
    name  = "proxiesIdPut_proxiesGet200Response_ProxiesGet200Response_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_id_put_proxies_get200_response_proxies_get200_response_port
  }
  property {
    name  = "proxiesIdPut_proxiesGet200Response_ProxiesGet200Response_requiresAuthentication"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_id_put_proxies_get200_response_proxies_get200_response_requires_authentication
  }
  property {
    name  = "proxiesIdPut_proxiesGet200Response_ProxiesGet200Response_username"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_id_put_proxies_get200_response_proxies_get200_response_username
  }
  property {
    name  = "proxiesPost_proxiesGet200Response_ProxiesGet200ResponsePassword_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_post_proxies_get200_response_proxies_get200_response_password_encrypted_value
  }
  property {
    name  = "proxiesPost_proxiesGet200Response_ProxiesGet200ResponsePassword_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_post_proxies_get200_response_proxies_get200_response_password_value
  }
  property {
    name  = "proxiesPost_proxiesGet200Response_ProxiesGet200Response_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_post_proxies_get200_response_proxies_get200_response_description
  }
  property {
    name  = "proxiesPost_proxiesGet200Response_ProxiesGet200Response_host"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_post_proxies_get200_response_proxies_get200_response_host
  }
  property {
    name  = "proxiesPost_proxiesGet200Response_ProxiesGet200Response_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_post_proxies_get200_response_proxies_get200_response_id
  }
  property {
    name  = "proxiesPost_proxiesGet200Response_ProxiesGet200Response_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_post_proxies_get200_response_proxies_get200_response_name
  }
  property {
    name  = "proxiesPost_proxiesGet200Response_ProxiesGet200Response_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_post_proxies_get200_response_proxies_get200_response_port
  }
  property {
    name  = "proxiesPost_proxiesGet200Response_ProxiesGet200Response_requiresAuthentication"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_post_proxies_get200_response_proxies_get200_response_requires_authentication
  }
  property {
    name  = "proxiesPost_proxiesGet200Response_ProxiesGet200Response_username"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_proxies_post_proxies_get200_response_proxies_get200_response_username
  }
  property {
    name  = "redirectsGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_get_filter
  }
  property {
    name  = "redirectsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_get_number_per_page
  }
  property {
    name  = "redirectsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_get_order
  }
  property {
    name  = "redirectsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_get_page
  }
  property {
    name  = "redirectsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_get_sort_key
  }
  property {
    name  = "redirectsGet_source"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_get_source
  }
  property {
    name  = "redirectsGet_target"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_get_target
  }
  property {
    name  = "redirectsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_id_delete_id
  }
  property {
    name  = "redirectsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_id_get_id
  }
  property {
    name  = "redirectsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_id_put_id
  }
  property {
    name  = "redirectsIdPut_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInnerSource_host"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_id_put_redirects_get200_response_items_inner_redirects_get200_response_items_inner_source_host
  }
  property {
    name  = "redirectsIdPut_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInnerSource_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_id_put_redirects_get200_response_items_inner_redirects_get200_response_items_inner_source_port
  }
  property {
    name  = "redirectsIdPut_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInnerTarget_host"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_id_put_redirects_get200_response_items_inner_redirects_get200_response_items_inner_target_host
  }
  property {
    name  = "redirectsIdPut_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInnerTarget_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_id_put_redirects_get200_response_items_inner_redirects_get200_response_items_inner_target_port
  }
  property {
    name  = "redirectsIdPut_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInnerTarget_secure"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_id_put_redirects_get200_response_items_inner_redirects_get200_response_items_inner_target_secure
  }
  property {
    name  = "redirectsIdPut_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInner_auditLevel"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_id_put_redirects_get200_response_items_inner_redirects_get200_response_items_inner_audit_level
  }
  property {
    name  = "redirectsIdPut_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_id_put_redirects_get200_response_items_inner_redirects_get200_response_items_inner_id
  }
  property {
    name  = "redirectsIdPut_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInner_responseCode"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_id_put_redirects_get200_response_items_inner_redirects_get200_response_items_inner_response_code
  }
  property {
    name  = "redirectsPost_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInnerSource_host"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_post_redirects_get200_response_items_inner_redirects_get200_response_items_inner_source_host
  }
  property {
    name  = "redirectsPost_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInnerSource_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_post_redirects_get200_response_items_inner_redirects_get200_response_items_inner_source_port
  }
  property {
    name  = "redirectsPost_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInnerTarget_host"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_post_redirects_get200_response_items_inner_redirects_get200_response_items_inner_target_host
  }
  property {
    name  = "redirectsPost_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInnerTarget_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_post_redirects_get200_response_items_inner_redirects_get200_response_items_inner_target_port
  }
  property {
    name  = "redirectsPost_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInnerTarget_secure"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_post_redirects_get200_response_items_inner_redirects_get200_response_items_inner_target_secure
  }
  property {
    name  = "redirectsPost_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInner_auditLevel"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_post_redirects_get200_response_items_inner_redirects_get200_response_items_inner_audit_level
  }
  property {
    name  = "redirectsPost_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_post_redirects_get200_response_items_inner_redirects_get200_response_items_inner_id
  }
  property {
    name  = "redirectsPost_redirectsGet200ResponseItemsInner_RedirectsGet200ResponseItemsInner_responseCode"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_redirects_post_redirects_get200_response_items_inner_redirects_get200_response_items_inner_response_code
  }
  property {
    name  = "rejectionHandlersDescriptorsRejectionHandlerTypeGet_rejection_handler_type"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_descriptors_rejection_handler_type_get_rejection_handler_type
  }
  property {
    name  = "rejectionHandlersGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_get_filter
  }
  property {
    name  = "rejectionHandlersGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_get_name
  }
  property {
    name  = "rejectionHandlersGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_get_number_per_page
  }
  property {
    name  = "rejectionHandlersGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_get_order
  }
  property {
    name  = "rejectionHandlersGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_get_page
  }
  property {
    name  = "rejectionHandlersGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_get_sort_key
  }
  property {
    name  = "rejectionHandlersIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_id_delete_id
  }
  property {
    name  = "rejectionHandlersIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_id_get_id
  }
  property {
    name  = "rejectionHandlersIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_id_put_id
  }
  property {
    name  = "rejectionHandlersIdPut_rejectionHandlersGet200ResponseItemsInner_RejectionHandlersGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_id_put_rejection_handlers_get200_response_items_inner_rejection_handlers_get200_response_items_inner_class_name
  }
  property {
    name  = "rejectionHandlersIdPut_rejectionHandlersGet200ResponseItemsInner_RejectionHandlersGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_id_put_rejection_handlers_get200_response_items_inner_rejection_handlers_get200_response_items_inner_configuration
  }
  property {
    name  = "rejectionHandlersIdPut_rejectionHandlersGet200ResponseItemsInner_RejectionHandlersGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_id_put_rejection_handlers_get200_response_items_inner_rejection_handlers_get200_response_items_inner_id
  }
  property {
    name  = "rejectionHandlersIdPut_rejectionHandlersGet200ResponseItemsInner_RejectionHandlersGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_id_put_rejection_handlers_get200_response_items_inner_rejection_handlers_get200_response_items_inner_name
  }
  property {
    name  = "rejectionHandlersPost_rejectionHandlersGet200ResponseItemsInner_RejectionHandlersGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_post_rejection_handlers_get200_response_items_inner_rejection_handlers_get200_response_items_inner_class_name
  }
  property {
    name  = "rejectionHandlersPost_rejectionHandlersGet200ResponseItemsInner_RejectionHandlersGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_post_rejection_handlers_get200_response_items_inner_rejection_handlers_get200_response_items_inner_configuration
  }
  property {
    name  = "rejectionHandlersPost_rejectionHandlersGet200ResponseItemsInner_RejectionHandlersGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_post_rejection_handlers_get200_response_items_inner_rejection_handlers_get200_response_items_inner_id
  }
  property {
    name  = "rejectionHandlersPost_rejectionHandlersGet200ResponseItemsInner_RejectionHandlersGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rejection_handlers_post_rejection_handlers_get200_response_items_inner_rejection_handlers_get200_response_items_inner_name
  }
  property {
    name  = "rulesDescriptorsRuleTypeGet_rule_type"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_descriptors_rule_type_get_rule_type
  }
  property {
    name  = "rulesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_get_filter
  }
  property {
    name  = "rulesGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_get_name
  }
  property {
    name  = "rulesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_get_number_per_page
  }
  property {
    name  = "rulesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_get_order
  }
  property {
    name  = "rulesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_get_page
  }
  property {
    name  = "rulesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_get_sort_key
  }
  property {
    name  = "rulesIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_id_delete_id
  }
  property {
    name  = "rulesIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_id_get_id
  }
  property {
    name  = "rulesIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_id_put_id
  }
  property {
    name  = "rulesIdPut_rulesGet200ResponseItemsInner_RulesGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_id_put_rules_get200_response_items_inner_rules_get200_response_items_inner_class_name
  }
  property {
    name  = "rulesIdPut_rulesGet200ResponseItemsInner_RulesGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_id_put_rules_get200_response_items_inner_rules_get200_response_items_inner_configuration
  }
  property {
    name  = "rulesIdPut_rulesGet200ResponseItemsInner_RulesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_id_put_rules_get200_response_items_inner_rules_get200_response_items_inner_id
  }
  property {
    name  = "rulesIdPut_rulesGet200ResponseItemsInner_RulesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_id_put_rules_get200_response_items_inner_rules_get200_response_items_inner_name
  }
  property {
    name  = "rulesIdPut_rulesGet200ResponseItemsInner_RulesGet200ResponseItemsInner_supportedDestinations"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_id_put_rules_get200_response_items_inner_rules_get200_response_items_inner_supported_destinations
  }
  property {
    name  = "rulesPost_rulesGet200ResponseItemsInner_RulesGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_post_rules_get200_response_items_inner_rules_get200_response_items_inner_class_name
  }
  property {
    name  = "rulesPost_rulesGet200ResponseItemsInner_RulesGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_post_rules_get200_response_items_inner_rules_get200_response_items_inner_configuration
  }
  property {
    name  = "rulesPost_rulesGet200ResponseItemsInner_RulesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_post_rules_get200_response_items_inner_rules_get200_response_items_inner_id
  }
  property {
    name  = "rulesPost_rulesGet200ResponseItemsInner_RulesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_post_rules_get200_response_items_inner_rules_get200_response_items_inner_name
  }
  property {
    name  = "rulesPost_rulesGet200ResponseItemsInner_RulesGet200ResponseItemsInner_supportedDestinations"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rules_post_rules_get200_response_items_inner_rules_get200_response_items_inner_supported_destinations
  }
  property {
    name  = "rulesetsGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_get_filter
  }
  property {
    name  = "rulesetsGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_get_name
  }
  property {
    name  = "rulesetsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_get_number_per_page
  }
  property {
    name  = "rulesetsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_get_order
  }
  property {
    name  = "rulesetsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_get_page
  }
  property {
    name  = "rulesetsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_get_sort_key
  }
  property {
    name  = "rulesetsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_id_delete_id
  }
  property {
    name  = "rulesetsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_id_get_id
  }
  property {
    name  = "rulesetsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_id_put_id
  }
  property {
    name  = "rulesetsIdPut_rulesetsGet200ResponseItemsInner_RulesetsGet200ResponseItemsInner_elementType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_id_put_rulesets_get200_response_items_inner_rulesets_get200_response_items_inner_element_type
  }
  property {
    name  = "rulesetsIdPut_rulesetsGet200ResponseItemsInner_RulesetsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_id_put_rulesets_get200_response_items_inner_rulesets_get200_response_items_inner_id
  }
  property {
    name  = "rulesetsIdPut_rulesetsGet200ResponseItemsInner_RulesetsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_id_put_rulesets_get200_response_items_inner_rulesets_get200_response_items_inner_name
  }
  property {
    name  = "rulesetsIdPut_rulesetsGet200ResponseItemsInner_RulesetsGet200ResponseItemsInner_policy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_id_put_rulesets_get200_response_items_inner_rulesets_get200_response_items_inner_policy
  }
  property {
    name  = "rulesetsIdPut_rulesetsGet200ResponseItemsInner_RulesetsGet200ResponseItemsInner_successCriteria"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_id_put_rulesets_get200_response_items_inner_rulesets_get200_response_items_inner_success_criteria
  }
  property {
    name  = "rulesetsPost_rulesetsGet200ResponseItemsInner_RulesetsGet200ResponseItemsInner_elementType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_post_rulesets_get200_response_items_inner_rulesets_get200_response_items_inner_element_type
  }
  property {
    name  = "rulesetsPost_rulesetsGet200ResponseItemsInner_RulesetsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_post_rulesets_get200_response_items_inner_rulesets_get200_response_items_inner_id
  }
  property {
    name  = "rulesetsPost_rulesetsGet200ResponseItemsInner_RulesetsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_post_rulesets_get200_response_items_inner_rulesets_get200_response_items_inner_name
  }
  property {
    name  = "rulesetsPost_rulesetsGet200ResponseItemsInner_RulesetsGet200ResponseItemsInner_policy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_post_rulesets_get200_response_items_inner_rulesets_get200_response_items_inner_policy
  }
  property {
    name  = "rulesetsPost_rulesetsGet200ResponseItemsInner_RulesetsGet200ResponseItemsInner_successCriteria"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_rulesets_post_rulesets_get200_response_items_inner_rulesets_get200_response_items_inner_success_criteria
  }
  property {
    name  = "sharedSecretsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_shared_secrets_get_order
  }
  property {
    name  = "sharedSecretsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_shared_secrets_get_sort_key
  }
  property {
    name  = "sharedSecretsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_shared_secrets_id_delete_id
  }
  property {
    name  = "sharedSecretsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_shared_secrets_id_get_id
  }
  property {
    name  = "sharedSecretsPost_sharedSecretsGet200ResponseItemsInner_SharedSecretsGet200ResponseItemsInnerSecret_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_shared_secrets_post_shared_secrets_get200_response_items_inner_shared_secrets_get200_response_items_inner_secret_encrypted_value
  }
  property {
    name  = "sharedSecretsPost_sharedSecretsGet200ResponseItemsInner_SharedSecretsGet200ResponseItemsInnerSecret_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_shared_secrets_post_shared_secrets_get200_response_items_inner_shared_secrets_get200_response_items_inner_secret_value
  }
  property {
    name  = "sharedSecretsPost_sharedSecretsGet200ResponseItemsInner_SharedSecretsGet200ResponseItemsInner_created"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_shared_secrets_post_shared_secrets_get200_response_items_inner_shared_secrets_get200_response_items_inner_created
  }
  property {
    name  = "sharedSecretsPost_sharedSecretsGet200ResponseItemsInner_SharedSecretsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_shared_secrets_post_shared_secrets_get200_response_items_inner_shared_secrets_get200_response_items_inner_id
  }
  property {
    name  = "sidebandClientsGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_get_filter
  }
  property {
    name  = "sidebandClientsGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_get_name
  }
  property {
    name  = "sidebandClientsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_get_number_per_page
  }
  property {
    name  = "sidebandClientsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_get_order
  }
  property {
    name  = "sidebandClientsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_get_page
  }
  property {
    name  = "sidebandClientsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_get_sort_key
  }
  property {
    name  = "sidebandClientsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_id_delete_id
  }
  property {
    name  = "sidebandClientsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_id_get_id
  }
  property {
    name  = "sidebandClientsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_id_put_id
  }
  property {
    name  = "sidebandClientsIdPut_sidebandClientsGet200ResponseItemsInner_SidebandClientsGet200ResponseItemsInner_clientCredentials"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_id_put_sideband_clients_get200_response_items_inner_sideband_clients_get200_response_items_inner_client_credentials
  }
  property {
    name  = "sidebandClientsIdPut_sidebandClientsGet200ResponseItemsInner_SidebandClientsGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_id_put_sideband_clients_get200_response_items_inner_sideband_clients_get200_response_items_inner_description
  }
  property {
    name  = "sidebandClientsIdPut_sidebandClientsGet200ResponseItemsInner_SidebandClientsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_id_put_sideband_clients_get200_response_items_inner_sideband_clients_get200_response_items_inner_id
  }
  property {
    name  = "sidebandClientsIdPut_sidebandClientsGet200ResponseItemsInner_SidebandClientsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_id_put_sideband_clients_get200_response_items_inner_sideband_clients_get200_response_items_inner_name
  }
  property {
    name  = "sidebandClientsPost_sidebandClientsGet200ResponseItemsInner_SidebandClientsGet200ResponseItemsInner_clientCredentials"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_post_sideband_clients_get200_response_items_inner_sideband_clients_get200_response_items_inner_client_credentials
  }
  property {
    name  = "sidebandClientsPost_sidebandClientsGet200ResponseItemsInner_SidebandClientsGet200ResponseItemsInner_description"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_post_sideband_clients_get200_response_items_inner_sideband_clients_get200_response_items_inner_description
  }
  property {
    name  = "sidebandClientsPost_sidebandClientsGet200ResponseItemsInner_SidebandClientsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_post_sideband_clients_get200_response_items_inner_sideband_clients_get200_response_items_inner_id
  }
  property {
    name  = "sidebandClientsPost_sidebandClientsGet200ResponseItemsInner_SidebandClientsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sideband_clients_post_sideband_clients_get200_response_items_inner_sideband_clients_get200_response_items_inner_name
  }
  property {
    name  = "siteAuthenticatorsDescriptorsSiteAuthenticatorTypeGet_site_authenticator_type"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_descriptors_site_authenticator_type_get_site_authenticator_type
  }
  property {
    name  = "siteAuthenticatorsGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_get_filter
  }
  property {
    name  = "siteAuthenticatorsGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_get_name
  }
  property {
    name  = "siteAuthenticatorsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_get_number_per_page
  }
  property {
    name  = "siteAuthenticatorsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_get_order
  }
  property {
    name  = "siteAuthenticatorsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_get_page
  }
  property {
    name  = "siteAuthenticatorsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_get_sort_key
  }
  property {
    name  = "siteAuthenticatorsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_id_delete_id
  }
  property {
    name  = "siteAuthenticatorsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_id_get_id
  }
  property {
    name  = "siteAuthenticatorsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_id_put_id
  }
  property {
    name  = "siteAuthenticatorsIdPut_siteAuthenticatorsGet200ResponseItemsInner_SiteAuthenticatorsGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_id_put_site_authenticators_get200_response_items_inner_site_authenticators_get200_response_items_inner_class_name
  }
  property {
    name  = "siteAuthenticatorsIdPut_siteAuthenticatorsGet200ResponseItemsInner_SiteAuthenticatorsGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_id_put_site_authenticators_get200_response_items_inner_site_authenticators_get200_response_items_inner_configuration
  }
  property {
    name  = "siteAuthenticatorsIdPut_siteAuthenticatorsGet200ResponseItemsInner_SiteAuthenticatorsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_id_put_site_authenticators_get200_response_items_inner_site_authenticators_get200_response_items_inner_id
  }
  property {
    name  = "siteAuthenticatorsIdPut_siteAuthenticatorsGet200ResponseItemsInner_SiteAuthenticatorsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_id_put_site_authenticators_get200_response_items_inner_site_authenticators_get200_response_items_inner_name
  }
  property {
    name  = "siteAuthenticatorsPost_siteAuthenticatorsGet200ResponseItemsInner_SiteAuthenticatorsGet200ResponseItemsInner_className"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_post_site_authenticators_get200_response_items_inner_site_authenticators_get200_response_items_inner_class_name
  }
  property {
    name  = "siteAuthenticatorsPost_siteAuthenticatorsGet200ResponseItemsInner_SiteAuthenticatorsGet200ResponseItemsInner_configuration"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_post_site_authenticators_get200_response_items_inner_site_authenticators_get200_response_items_inner_configuration
  }
  property {
    name  = "siteAuthenticatorsPost_siteAuthenticatorsGet200ResponseItemsInner_SiteAuthenticatorsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_post_site_authenticators_get200_response_items_inner_site_authenticators_get200_response_items_inner_id
  }
  property {
    name  = "siteAuthenticatorsPost_siteAuthenticatorsGet200ResponseItemsInner_SiteAuthenticatorsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_site_authenticators_post_site_authenticators_get200_response_items_inner_site_authenticators_get200_response_items_inner_name
  }
  property {
    name  = "sitesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_get_filter
  }
  property {
    name  = "sitesGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_get_name
  }
  property {
    name  = "sitesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_get_number_per_page
  }
  property {
    name  = "sitesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_get_order
  }
  property {
    name  = "sitesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_get_page
  }
  property {
    name  = "sitesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_get_sort_key
  }
  property {
    name  = "sitesIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_delete_id
  }
  property {
    name  = "sitesIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_get_id
  }
  property {
    name  = "sitesIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_id
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_availabilityProfileId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_availability_profile_id
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_expectedHostname"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_expected_hostname
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_id
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_keepAliveTimeout"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_keep_alive_timeout
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_loadBalancingStrategyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_load_balancing_strategy_id
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_maxConnections"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_max_connections
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_maxWebSocketConnections"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_max_web_socket_connections
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_name
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_secure"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_secure
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_sendPaCookie"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_send_pa_cookie
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_siteAuthenticatorIds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_site_authenticator_ids
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_skipHostnameVerification"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_skip_hostname_verification
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_targets"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_targets
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_trusted_certificate_group_id
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_useProxy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_use_proxy
  }
  property {
    name  = "sitesIdPut_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_useTargetHostHeader"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_id_put_sites_get200_response_items_inner_sites_get200_response_items_inner_use_target_host_header
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_availabilityProfileId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_availability_profile_id
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_expectedHostname"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_expected_hostname
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_id
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_keepAliveTimeout"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_keep_alive_timeout
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_loadBalancingStrategyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_load_balancing_strategy_id
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_maxConnections"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_max_connections
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_maxWebSocketConnections"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_max_web_socket_connections
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_name
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_secure"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_secure
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_sendPaCookie"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_send_pa_cookie
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_siteAuthenticatorIds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_site_authenticator_ids
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_skipHostnameVerification"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_skip_hostname_verification
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_targets"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_targets
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_trusted_certificate_group_id
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_useProxy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_use_proxy
  }
  property {
    name  = "sitesPost_sitesGet200ResponseItemsInner_SitesGet200ResponseItemsInner_useTargetHostHeader"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_sites_post_sites_get200_response_items_inner_sites_get200_response_items_inner_use_target_host_header
  }
  property {
    name  = "sslVerification"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_ssl_verification
  }
  property {
    name  = "thirdPartyServicesGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_get_filter
  }
  property {
    name  = "thirdPartyServicesGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_get_name
  }
  property {
    name  = "thirdPartyServicesGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_get_number_per_page
  }
  property {
    name  = "thirdPartyServicesGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_get_order
  }
  property {
    name  = "thirdPartyServicesGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_get_page
  }
  property {
    name  = "thirdPartyServicesGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_get_sort_key
  }
  property {
    name  = "thirdPartyServicesIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_delete_id
  }
  property {
    name  = "thirdPartyServicesIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_get_id
  }
  property {
    name  = "thirdPartyServicesIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_put_id
  }
  property {
    name  = "thirdPartyServicesIdPut_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_availabilityProfileId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_put_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_availability_profile_id
  }
  property {
    name  = "thirdPartyServicesIdPut_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_expectedHostname"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_put_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_expected_hostname
  }
  property {
    name  = "thirdPartyServicesIdPut_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_hostValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_put_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_host_value
  }
  property {
    name  = "thirdPartyServicesIdPut_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_put_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_id
  }
  property {
    name  = "thirdPartyServicesIdPut_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_loadBalancingStrategyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_put_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_load_balancing_strategy_id
  }
  property {
    name  = "thirdPartyServicesIdPut_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_maxConnections"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_put_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_max_connections
  }
  property {
    name  = "thirdPartyServicesIdPut_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_put_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_name
  }
  property {
    name  = "thirdPartyServicesIdPut_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_secure"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_put_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_secure
  }
  property {
    name  = "thirdPartyServicesIdPut_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_skipHostnameVerification"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_put_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_skip_hostname_verification
  }
  property {
    name  = "thirdPartyServicesIdPut_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_targets"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_put_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_targets
  }
  property {
    name  = "thirdPartyServicesIdPut_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_put_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_trusted_certificate_group_id
  }
  property {
    name  = "thirdPartyServicesIdPut_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_useProxy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_id_put_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_use_proxy
  }
  property {
    name  = "thirdPartyServicesPost_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_availabilityProfileId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_post_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_availability_profile_id
  }
  property {
    name  = "thirdPartyServicesPost_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_expectedHostname"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_post_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_expected_hostname
  }
  property {
    name  = "thirdPartyServicesPost_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_hostValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_post_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_host_value
  }
  property {
    name  = "thirdPartyServicesPost_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_post_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_id
  }
  property {
    name  = "thirdPartyServicesPost_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_loadBalancingStrategyId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_post_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_load_balancing_strategy_id
  }
  property {
    name  = "thirdPartyServicesPost_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_maxConnections"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_post_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_max_connections
  }
  property {
    name  = "thirdPartyServicesPost_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_post_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_name
  }
  property {
    name  = "thirdPartyServicesPost_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_secure"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_post_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_secure
  }
  property {
    name  = "thirdPartyServicesPost_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_skipHostnameVerification"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_post_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_skip_hostname_verification
  }
  property {
    name  = "thirdPartyServicesPost_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_targets"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_post_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_targets
  }
  property {
    name  = "thirdPartyServicesPost_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_post_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_trusted_certificate_group_id
  }
  property {
    name  = "thirdPartyServicesPost_thirdPartyServicesGet200ResponseItemsInner_ThirdPartyServicesGet200ResponseItemsInner_useProxy"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_third_party_services_post_third_party_services_get200_response_items_inner_third_party_services_get200_response_items_inner_use_proxy
  }
  property {
    name  = "tokenProviderSettingsPut_tokenProviderSettingsDelete200Response_TokenProviderSettingsDelete200Response_type"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_token_provider_settings_put_token_provider_settings_delete200_response_token_provider_settings_delete200_response_type
  }
  property {
    name  = "tokenProviderSettingsPut_tokenProviderSettingsDelete200Response_TokenProviderSettingsDelete200Response_useThirdParty"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_token_provider_settings_put_token_provider_settings_delete200_response_token_provider_settings_delete200_response_use_third_party
  }
  property {
    name  = "trustedCertificateGroupsGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_get_filter
  }
  property {
    name  = "trustedCertificateGroupsGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_get_name
  }
  property {
    name  = "trustedCertificateGroupsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_get_number_per_page
  }
  property {
    name  = "trustedCertificateGroupsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_get_order
  }
  property {
    name  = "trustedCertificateGroupsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_get_page
  }
  property {
    name  = "trustedCertificateGroupsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_get_sort_key
  }
  property {
    name  = "trustedCertificateGroupsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_delete_id
  }
  property {
    name  = "trustedCertificateGroupsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_get_id
  }
  property {
    name  = "trustedCertificateGroupsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_put_id
  }
  property {
    name  = "trustedCertificateGroupsIdPut_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInnerRevocationChecking_crlChecking"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_put_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_revocation_checking_crl_checking
  }
  property {
    name  = "trustedCertificateGroupsIdPut_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInnerRevocationChecking_denyRevocationStatusUnknown"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_put_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_revocation_checking_deny_revocation_status_unknown
  }
  property {
    name  = "trustedCertificateGroupsIdPut_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInnerRevocationChecking_ocsp"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_put_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_revocation_checking_ocsp
  }
  property {
    name  = "trustedCertificateGroupsIdPut_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInnerRevocationChecking_skipTrustAnchors"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_put_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_revocation_checking_skip_trust_anchors
  }
  property {
    name  = "trustedCertificateGroupsIdPut_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInnerRevocationChecking_supportDisorderedChain"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_put_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_revocation_checking_support_disordered_chain
  }
  property {
    name  = "trustedCertificateGroupsIdPut_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_certIds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_put_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_cert_ids
  }
  property {
    name  = "trustedCertificateGroupsIdPut_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_put_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_id
  }
  property {
    name  = "trustedCertificateGroupsIdPut_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_ignoreAllCertificateErrors"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_put_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_ignore_all_certificate_errors
  }
  property {
    name  = "trustedCertificateGroupsIdPut_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_put_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_name
  }
  property {
    name  = "trustedCertificateGroupsIdPut_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_skipCertificateDateCheck"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_put_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_skip_certificate_date_check
  }
  property {
    name  = "trustedCertificateGroupsIdPut_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_systemGroup"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_put_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_system_group
  }
  property {
    name  = "trustedCertificateGroupsIdPut_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_useJavaTrustStore"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_id_put_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_use_java_trust_store
  }
  property {
    name  = "trustedCertificateGroupsPost_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInnerRevocationChecking_crlChecking"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_post_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_revocation_checking_crl_checking
  }
  property {
    name  = "trustedCertificateGroupsPost_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInnerRevocationChecking_denyRevocationStatusUnknown"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_post_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_revocation_checking_deny_revocation_status_unknown
  }
  property {
    name  = "trustedCertificateGroupsPost_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInnerRevocationChecking_ocsp"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_post_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_revocation_checking_ocsp
  }
  property {
    name  = "trustedCertificateGroupsPost_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInnerRevocationChecking_skipTrustAnchors"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_post_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_revocation_checking_skip_trust_anchors
  }
  property {
    name  = "trustedCertificateGroupsPost_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInnerRevocationChecking_supportDisorderedChain"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_post_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_revocation_checking_support_disordered_chain
  }
  property {
    name  = "trustedCertificateGroupsPost_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_certIds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_post_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_cert_ids
  }
  property {
    name  = "trustedCertificateGroupsPost_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_post_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_id
  }
  property {
    name  = "trustedCertificateGroupsPost_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_ignoreAllCertificateErrors"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_post_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_ignore_all_certificate_errors
  }
  property {
    name  = "trustedCertificateGroupsPost_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_post_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_name
  }
  property {
    name  = "trustedCertificateGroupsPost_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_skipCertificateDateCheck"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_post_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_skip_certificate_date_check
  }
  property {
    name  = "trustedCertificateGroupsPost_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_systemGroup"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_post_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_system_group
  }
  property {
    name  = "trustedCertificateGroupsPost_trustedCertificateGroupsGet200ResponseItemsInner_TrustedCertificateGroupsGet200ResponseItemsInner_useJavaTrustStore"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_trusted_certificate_groups_post_trusted_certificate_groups_get200_response_items_inner_trusted_certificate_groups_get200_response_items_inner_use_java_trust_store
  }
  property {
    name  = "unknownResourcesSettingsPut_unknownResourcesSettingsDelete200Response_UnknownResourcesSettingsDelete200Response_agentDefaultCacheTTL"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_unknown_resources_settings_put_unknown_resources_settings_delete200_response_unknown_resources_settings_delete200_response_agent_default_cache_ttl
  }
  property {
    name  = "unknownResourcesSettingsPut_unknownResourcesSettingsDelete200Response_UnknownResourcesSettingsDelete200Response_agentDefaultMode"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_unknown_resources_settings_put_unknown_resources_settings_delete200_response_unknown_resources_settings_delete200_response_agent_default_mode
  }
  property {
    name  = "unknownResourcesSettingsPut_unknownResourcesSettingsDelete200Response_UnknownResourcesSettingsDelete200Response_auditLevel"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_unknown_resources_settings_put_unknown_resources_settings_delete200_response_unknown_resources_settings_delete200_response_audit_level
  }
  property {
    name  = "unknownResourcesSettingsPut_unknownResourcesSettingsDelete200Response_UnknownResourcesSettingsDelete200Response_errorContentType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_unknown_resources_settings_put_unknown_resources_settings_delete200_response_unknown_resources_settings_delete200_response_error_content_type
  }
  property {
    name  = "unknownResourcesSettingsPut_unknownResourcesSettingsDelete200Response_UnknownResourcesSettingsDelete200Response_errorStatusCode"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_unknown_resources_settings_put_unknown_resources_settings_delete200_response_unknown_resources_settings_delete200_response_error_status_code
  }
  property {
    name  = "unknownResourcesSettingsPut_unknownResourcesSettingsDelete200Response_UnknownResourcesSettingsDelete200Response_errorTemplateFile"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_unknown_resources_settings_put_unknown_resources_settings_delete200_response_unknown_resources_settings_delete200_response_error_template_file
  }
  property {
    name  = "usersGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_get_filter
  }
  property {
    name  = "usersGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_get_number_per_page
  }
  property {
    name  = "usersGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_get_order
  }
  property {
    name  = "usersGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_get_page
  }
  property {
    name  = "usersGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_get_sort_key
  }
  property {
    name  = "usersGet_username"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_get_username
  }
  property {
    name  = "usersIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_id_get_id
  }
  property {
    name  = "usersIdPasswordPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_id_password_put_id
  }
  property {
    name  = "usersIdPasswordPut_usersIdPasswordPutRequest_UsersIdPasswordPutRequest_currentPassword"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_id_password_put_users_id_password_put_request_users_id_password_put_request_current_password
  }
  property {
    name  = "usersIdPasswordPut_usersIdPasswordPutRequest_UsersIdPasswordPutRequest_newPassword"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_id_password_put_users_id_password_put_request_users_id_password_put_request_new_password
  }
  property {
    name  = "usersIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_id_put_id
  }
  property {
    name  = "usersIdPut_usersGet200ResponseItemsInner_UsersGet200ResponseItemsInner_email"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_id_put_users_get200_response_items_inner_users_get200_response_items_inner_email
  }
  property {
    name  = "usersIdPut_usersGet200ResponseItemsInner_UsersGet200ResponseItemsInner_firstLogin"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_id_put_users_get200_response_items_inner_users_get200_response_items_inner_first_login
  }
  property {
    name  = "usersIdPut_usersGet200ResponseItemsInner_UsersGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_id_put_users_get200_response_items_inner_users_get200_response_items_inner_id
  }
  property {
    name  = "usersIdPut_usersGet200ResponseItemsInner_UsersGet200ResponseItemsInner_showTutorial"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_id_put_users_get200_response_items_inner_users_get200_response_items_inner_show_tutorial
  }
  property {
    name  = "usersIdPut_usersGet200ResponseItemsInner_UsersGet200ResponseItemsInner_slaAccepted"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_id_put_users_get200_response_items_inner_users_get200_response_items_inner_sla_accepted
  }
  property {
    name  = "usersIdPut_usersGet200ResponseItemsInner_UsersGet200ResponseItemsInner_username"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_users_id_put_users_get200_response_items_inner_users_get200_response_items_inner_username
  }
  property {
    name  = "virtualhostsGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_get_filter
  }
  property {
    name  = "virtualhostsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_get_number_per_page
  }
  property {
    name  = "virtualhostsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_get_order
  }
  property {
    name  = "virtualhostsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_get_page
  }
  property {
    name  = "virtualhostsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_get_sort_key
  }
  property {
    name  = "virtualhostsGet_virtual_host"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_get_virtual_host
  }
  property {
    name  = "virtualhostsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_id_delete_id
  }
  property {
    name  = "virtualhostsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_id_get_id
  }
  property {
    name  = "virtualhostsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_id_put_id
  }
  property {
    name  = "virtualhostsIdPut_virtualhostsGet200ResponseItemsInner_VirtualhostsGet200ResponseItemsInner_agentResourceCacheTTL"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_id_put_virtualhosts_get200_response_items_inner_virtualhosts_get200_response_items_inner_agent_resource_cache_ttl
  }
  property {
    name  = "virtualhostsIdPut_virtualhostsGet200ResponseItemsInner_VirtualhostsGet200ResponseItemsInner_host"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_id_put_virtualhosts_get200_response_items_inner_virtualhosts_get200_response_items_inner_host
  }
  property {
    name  = "virtualhostsIdPut_virtualhostsGet200ResponseItemsInner_VirtualhostsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_id_put_virtualhosts_get200_response_items_inner_virtualhosts_get200_response_items_inner_id
  }
  property {
    name  = "virtualhostsIdPut_virtualhostsGet200ResponseItemsInner_VirtualhostsGet200ResponseItemsInner_keyPairId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_id_put_virtualhosts_get200_response_items_inner_virtualhosts_get200_response_items_inner_key_pair_id
  }
  property {
    name  = "virtualhostsIdPut_virtualhostsGet200ResponseItemsInner_VirtualhostsGet200ResponseItemsInner_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_id_put_virtualhosts_get200_response_items_inner_virtualhosts_get200_response_items_inner_port
  }
  property {
    name  = "virtualhostsIdPut_virtualhostsGet200ResponseItemsInner_VirtualhostsGet200ResponseItemsInner_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_id_put_virtualhosts_get200_response_items_inner_virtualhosts_get200_response_items_inner_trusted_certificate_group_id
  }
  property {
    name  = "virtualhostsPost_virtualhostsGet200ResponseItemsInner_VirtualhostsGet200ResponseItemsInner_agentResourceCacheTTL"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_post_virtualhosts_get200_response_items_inner_virtualhosts_get200_response_items_inner_agent_resource_cache_ttl
  }
  property {
    name  = "virtualhostsPost_virtualhostsGet200ResponseItemsInner_VirtualhostsGet200ResponseItemsInner_host"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_post_virtualhosts_get200_response_items_inner_virtualhosts_get200_response_items_inner_host
  }
  property {
    name  = "virtualhostsPost_virtualhostsGet200ResponseItemsInner_VirtualhostsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_post_virtualhosts_get200_response_items_inner_virtualhosts_get200_response_items_inner_id
  }
  property {
    name  = "virtualhostsPost_virtualhostsGet200ResponseItemsInner_VirtualhostsGet200ResponseItemsInner_keyPairId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_post_virtualhosts_get200_response_items_inner_virtualhosts_get200_response_items_inner_key_pair_id
  }
  property {
    name  = "virtualhostsPost_virtualhostsGet200ResponseItemsInner_VirtualhostsGet200ResponseItemsInner_port"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_post_virtualhosts_get200_response_items_inner_virtualhosts_get200_response_items_inner_port
  }
  property {
    name  = "virtualhostsPost_virtualhostsGet200ResponseItemsInner_VirtualhostsGet200ResponseItemsInner_trustedCertificateGroupId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_virtualhosts_post_virtualhosts_get200_response_items_inner_virtualhosts_get200_response_items_inner_trusted_certificate_group_id
  }
  property {
    name  = "webSessionManagementKeySetPut_authTokenManagementKeySetGet200Response_AuthTokenManagementKeySetGet200Response_keySet"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_session_management_key_set_put_auth_token_management_key_set_get200_response_auth_token_management_key_set_get200_response_key_set
  }
  property {
    name  = "webSessionManagementKeySetPut_authTokenManagementKeySetGet200Response_AuthTokenManagementKeySetGet200Response_nonce"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_session_management_key_set_put_auth_token_management_key_set_get200_response_auth_token_management_key_set_get200_response_nonce
  }
  property {
    name  = "webSessionManagementOidcScopesGet_client_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_session_management_oidc_scopes_get_client_id
  }
  property {
    name  = "webSessionManagementPut_webSessionManagementDelete200Response_WebSessionManagementDelete200Response_cookieName"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_session_management_put_web_session_management_delete200_response_web_session_management_delete200_response_cookie_name
  }
  property {
    name  = "webSessionManagementPut_webSessionManagementDelete200Response_WebSessionManagementDelete200Response_encryptionAlgorithm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_session_management_put_web_session_management_delete200_response_web_session_management_delete200_response_encryption_algorithm
  }
  property {
    name  = "webSessionManagementPut_webSessionManagementDelete200Response_WebSessionManagementDelete200Response_issuer"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_session_management_put_web_session_management_delete200_response_web_session_management_delete200_response_issuer
  }
  property {
    name  = "webSessionManagementPut_webSessionManagementDelete200Response_WebSessionManagementDelete200Response_keyRollEnabled"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_session_management_put_web_session_management_delete200_response_web_session_management_delete200_response_key_roll_enabled
  }
  property {
    name  = "webSessionManagementPut_webSessionManagementDelete200Response_WebSessionManagementDelete200Response_keyRollPeriodInHours"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_session_management_put_web_session_management_delete200_response_web_session_management_delete200_response_key_roll_period_in_hours
  }
  property {
    name  = "webSessionManagementPut_webSessionManagementDelete200Response_WebSessionManagementDelete200Response_nonceCookieTimeToLiveInMinutes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_session_management_put_web_session_management_delete200_response_web_session_management_delete200_response_nonce_cookie_time_to_live_in_minutes
  }
  property {
    name  = "webSessionManagementPut_webSessionManagementDelete200Response_WebSessionManagementDelete200Response_sessionStateCookieName"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_session_management_put_web_session_management_delete200_response_web_session_management_delete200_response_session_state_cookie_name
  }
  property {
    name  = "webSessionManagementPut_webSessionManagementDelete200Response_WebSessionManagementDelete200Response_signingAlgorithm"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_session_management_put_web_session_management_delete200_response_web_session_management_delete200_response_signing_algorithm
  }
  property {
    name  = "webSessionManagementPut_webSessionManagementDelete200Response_WebSessionManagementDelete200Response_updateTokenWindowInSeconds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_session_management_put_web_session_management_delete200_response_web_session_management_delete200_response_update_token_window_in_seconds
  }
  property {
    name  = "webSessionsGet_filter"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_get_filter
  }
  property {
    name  = "webSessionsGet_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_get_name
  }
  property {
    name  = "webSessionsGet_number_per_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_get_number_per_page
  }
  property {
    name  = "webSessionsGet_order"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_get_order
  }
  property {
    name  = "webSessionsGet_page"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_get_page
  }
  property {
    name  = "webSessionsGet_sort_key"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_get_sort_key
  }
  property {
    name  = "webSessionsIdDelete_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_delete_id
  }
  property {
    name  = "webSessionsIdGet_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_get_id
  }
  property {
    name  = "webSessionsIdPut_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_id
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_AuthOauthDelete200ResponseClientCredentialsClientSecret_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_auth_oauth_delete200_response_client_credentials_client_secret_encrypted_value
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_AuthOauthDelete200ResponseClientCredentialsClientSecret_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_auth_oauth_delete200_response_client_credentials_client_secret_value
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_AuthOidcDelete200ResponseOidcConfigurationClientCredentials_clientId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_auth_oidc_delete200_response_oidc_configuration_client_credentials_client_id
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_AuthOidcDelete200ResponseOidcConfigurationClientCredentials_credentialsType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_auth_oidc_delete200_response_oidc_configuration_client_credentials_credentials_type
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_AuthOidcDelete200ResponseOidcConfigurationClientCredentials_keyPairId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_auth_oidc_delete200_response_oidc_configuration_client_credentials_key_pair_id
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_audience"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_audience
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_cacheUserAttributes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_cache_user_attributes
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_cookieDomain"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_cookie_domain
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_cookieType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_cookie_type
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_enableRefreshUser"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_enable_refresh_user
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_failOnUnsupportedPreservationContentType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_fail_on_unsupported_preservation_content_type
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_httpOnlyCookie"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_http_only_cookie
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_id
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_idleTimeoutInMinutes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_idle_timeout_in_minutes
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_name
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_oidcLoginType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_oidc_login_type
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_pfsessionStateCacheInSeconds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_pfsession_state_cache_in_seconds
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_pkceChallengeType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_pkce_challenge_type
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_refreshUserInfoClaimsInterval"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_refresh_user_info_claims_interval
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_requestPreservationType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_request_preservation_type
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_requestProfile"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_request_profile
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_sameSite"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_same_site
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_scopes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_scopes
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_secureCookie"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_secure_cookie
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_sendRequestedUrlToProvider"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_send_requested_url_to_provider
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_sessionTimeoutInMinutes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_session_timeout_in_minutes
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_validateSessionIsAlive"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_validate_session_is_alive
  }
  property {
    name  = "webSessionsIdPut_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_webStorageType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_id_put_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_web_storage_type
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_AuthOauthDelete200ResponseClientCredentialsClientSecret_encryptedValue"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_auth_oauth_delete200_response_client_credentials_client_secret_encrypted_value
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_AuthOauthDelete200ResponseClientCredentialsClientSecret_value"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_auth_oauth_delete200_response_client_credentials_client_secret_value
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_AuthOidcDelete200ResponseOidcConfigurationClientCredentials_clientId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_auth_oidc_delete200_response_oidc_configuration_client_credentials_client_id
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_AuthOidcDelete200ResponseOidcConfigurationClientCredentials_credentialsType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_auth_oidc_delete200_response_oidc_configuration_client_credentials_credentials_type
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_AuthOidcDelete200ResponseOidcConfigurationClientCredentials_keyPairId"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_auth_oidc_delete200_response_oidc_configuration_client_credentials_key_pair_id
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_audience"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_audience
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_cacheUserAttributes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_cache_user_attributes
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_cookieDomain"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_cookie_domain
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_cookieType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_cookie_type
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_enableRefreshUser"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_enable_refresh_user
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_failOnUnsupportedPreservationContentType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_fail_on_unsupported_preservation_content_type
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_httpOnlyCookie"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_http_only_cookie
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_id"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_id
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_idleTimeoutInMinutes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_idle_timeout_in_minutes
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_name"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_name
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_oidcLoginType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_oidc_login_type
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_pfsessionStateCacheInSeconds"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_pfsession_state_cache_in_seconds
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_pkceChallengeType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_pkce_challenge_type
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_refreshUserInfoClaimsInterval"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_refresh_user_info_claims_interval
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_requestPreservationType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_request_preservation_type
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_requestProfile"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_request_profile
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_sameSite"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_same_site
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_scopes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_scopes
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_secureCookie"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_secure_cookie
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_sendRequestedUrlToProvider"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_send_requested_url_to_provider
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_sessionTimeoutInMinutes"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_session_timeout_in_minutes
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_validateSessionIsAlive"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_validate_session_is_alive
  }
  property {
    name  = "webSessionsPost_webSessionsGet200ResponseItemsInner_WebSessionsGet200ResponseItemsInner_webStorageType"
    type  = "string"
    value = var.connector-oai-pingaccessadministrativeapi_property_web_sessions_post_web_sessions_get200_response_items_inner_web_sessions_get200_response_items_inner_web_storage_type
  }
}
