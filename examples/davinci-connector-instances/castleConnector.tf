resource "pingone_davinci_connector_instance" "castleConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "castleConnector"
  }
  name = "My awesome castleConnector"
  property {
    name  = "accept"
    type  = "string"
    value = var.castleconnector_property_accept
  }
  property {
    name  = "acceptEncoding"
    type  = "string"
    value = var.castleconnector_property_accept_encoding
  }
  property {
    name  = "acceptLanguage"
    type  = "string"
    value = var.castleconnector_property_accept_language
  }
  property {
    name  = "address_city"
    type  = "string"
    value = var.castleconnector_property_address_city
  }
  property {
    name  = "address_countrycode"
    type  = "string"
    value = var.castleconnector_property_address_countrycode
  }
  property {
    name  = "address_fingerprint"
    type  = "string"
    value = var.castleconnector_property_address_fingerprint
  }
  property {
    name  = "address_line1"
    type  = "string"
    value = var.castleconnector_property_address_line1
  }
  property {
    name  = "address_line2"
    type  = "string"
    value = var.castleconnector_property_address_line2
  }
  property {
    name  = "address_postalCode"
    type  = "string"
    value = var.castleconnector_property_address_postal_code
  }
  property {
    name  = "address_regionCode"
    type  = "string"
    value = var.castleconnector_property_address_region_code
  }
  property {
    name  = "amountType"
    type  = "string"
    value = var.castleconnector_property_amount_type
  }
  property {
    name  = "amountValue"
    type  = "string"
    value = var.castleconnector_property_amount_value
  }
  property {
    name  = "apiSecret"
    type  = "string"
    value = var.castleconnector_property_api_secret
  }
  property {
    name  = "auth_email"
    type  = "string"
    value = var.castleconnector_property_auth_email
  }
  property {
    name  = "auth_emailFilter"
    type  = "string"
    value = var.castleconnector_property_auth_email_filter
  }
  property {
    name  = "auth_emailFilter_chg"
    type  = "string"
    value = var.castleconnector_property_auth_email_filter_chg
  }
  property {
    name  = "auth_emailFilter_reg"
    type  = "string"
    value = var.castleconnector_property_auth_email_filter_reg
  }
  property {
    name  = "auth_email_chg"
    type  = "string"
    value = var.castleconnector_property_auth_email_chg
  }
  property {
    name  = "auth_email_reg"
    type  = "string"
    value = var.castleconnector_property_auth_email_reg
  }
  property {
    name  = "auth_phone"
    type  = "string"
    value = var.castleconnector_property_auth_phone
  }
  property {
    name  = "auth_phoneFilter"
    type  = "string"
    value = var.castleconnector_property_auth_phone_filter
  }
  property {
    name  = "auth_phoneFilter_chg"
    type  = "string"
    value = var.castleconnector_property_auth_phone_filter_chg
  }
  property {
    name  = "auth_phoneFilter_reg"
    type  = "string"
    value = var.castleconnector_property_auth_phone_filter_reg
  }
  property {
    name  = "auth_phone_chg"
    type  = "string"
    value = var.castleconnector_property_auth_phone_chg
  }
  property {
    name  = "auth_phone_reg"
    type  = "string"
    value = var.castleconnector_property_auth_phone_reg
  }
  property {
    name  = "auth_type"
    type  = "string"
    value = var.castleconnector_property_auth_type
  }
  property {
    name  = "auth_typeFilter"
    type  = "string"
    value = var.castleconnector_property_auth_type_filter
  }
  property {
    name  = "auth_typeFilter_chg"
    type  = "string"
    value = var.castleconnector_property_auth_type_filter_chg
  }
  property {
    name  = "auth_typeFilter_reg"
    type  = "string"
    value = var.castleconnector_property_auth_type_filter_reg
  }
  property {
    name  = "auth_type_chg"
    type  = "string"
    value = var.castleconnector_property_auth_type_chg
  }
  property {
    name  = "auth_type_reg"
    type  = "string"
    value = var.castleconnector_property_auth_type_reg
  }
  property {
    name  = "baseAmount"
    type  = "string"
    value = var.castleconnector_property_base_amount
  }
  property {
    name  = "billingAddressCcode"
    type  = "string"
    value = var.castleconnector_property_billing_address_ccode
  }
  property {
    name  = "billingAddressCity"
    type  = "string"
    value = var.castleconnector_property_billing_address_city
  }
  property {
    name  = "billingAddressFingeprint"
    type  = "string"
    value = var.castleconnector_property_billing_address_fingeprint
  }
  property {
    name  = "billingAddressLine1"
    type  = "string"
    value = var.castleconnector_property_billing_address_line1
  }
  property {
    name  = "billingAddressLine2"
    type  = "string"
    value = var.castleconnector_property_billing_address_line2
  }
  property {
    name  = "billingAddressPostalCode"
    type  = "string"
    value = var.castleconnector_property_billing_address_postal_code
  }
  property {
    name  = "billingAddressRegionCode"
    type  = "string"
    value = var.castleconnector_property_billing_address_region_code
  }
  property {
    name  = "cardBin"
    type  = "string"
    value = var.castleconnector_property_card_bin
  }
  property {
    name  = "cardExpMonth"
    type  = "string"
    value = var.castleconnector_property_card_exp_month
  }
  property {
    name  = "cardExpyear"
    type  = "string"
    value = var.castleconnector_property_card_expyear
  }
  property {
    name  = "cardFunding"
    type  = "string"
    value = var.castleconnector_property_card_funding
  }
  property {
    name  = "cardLast4"
    type  = "string"
    value = var.castleconnector_property_card_last4
  }
  property {
    name  = "cardNetwork"
    type  = "string"
    value = var.castleconnector_property_card_network
  }
  property {
    name  = "changeAuthenticatorMethodFrom"
    type  = "string"
    value = var.castleconnector_property_change_authenticator_method_from
  }
  property {
    name  = "changeAuthenticatorMethodFrom_reset"
    type  = "string"
    value = var.castleconnector_property_change_authenticator_method_from_reset
  }
  property {
    name  = "changeAuthenticatorMethodTo"
    type  = "string"
    value = var.castleconnector_property_change_authenticator_method_to
  }
  property {
    name  = "changeAuthenticatorMethodTo_reset"
    type  = "string"
    value = var.castleconnector_property_change_authenticator_method_to_reset
  }
  property {
    name  = "changeNameFrom"
    type  = "string"
    value = var.castleconnector_property_change_name_from
  }
  property {
    name  = "changeNameFrom_reset"
    type  = "string"
    value = var.castleconnector_property_change_name_from_reset
  }
  property {
    name  = "changeNameTo"
    type  = "string"
    value = var.castleconnector_property_change_name_to
  }
  property {
    name  = "changeNameTo_reset"
    type  = "string"
    value = var.castleconnector_property_change_name_to_reset
  }
  property {
    name  = "changesetPassword"
    type  = "string"
    value = var.castleconnector_property_changeset_password
  }
  property {
    name  = "changesetPassword_reset"
    type  = "string"
    value = var.castleconnector_property_changeset_password_reset
  }
  property {
    name  = "city"
    type  = "string"
    value = var.castleconnector_property_city
  }
  property {
    name  = "country_code"
    type  = "string"
    value = var.castleconnector_property_country_code
  }
  property {
    name  = "created_at"
    type  = "string"
    value = var.castleconnector_property_created_at
  }
  property {
    name  = "csetAuthtype"
    type  = "string"
    value = var.castleconnector_property_cset_authtype
  }
  property {
    name  = "csetAuthtype_reset"
    type  = "string"
    value = var.castleconnector_property_cset_authtype_reset
  }
  property {
    name  = "csetEmail"
    type  = "string"
    value = var.castleconnector_property_cset_email
  }
  property {
    name  = "csetEmail_reset"
    type  = "string"
    value = var.castleconnector_property_cset_email_reset
  }
  property {
    name  = "csetName"
    type  = "string"
    value = var.castleconnector_property_cset_name
  }
  property {
    name  = "csetName_reset"
    type  = "string"
    value = var.castleconnector_property_cset_name_reset
  }
  property {
    name  = "csetPhone"
    type  = "string"
    value = var.castleconnector_property_cset_phone
  }
  property {
    name  = "csetPhone_reset"
    type  = "string"
    value = var.castleconnector_property_cset_phone_reset
  }
  property {
    name  = "customName"
    type  = "string"
    value = var.castleconnector_property_custom_name
  }
  property {
    name  = "customNameFilter"
    type  = "string"
    value = var.castleconnector_property_custom_name_filter
  }
  property {
    name  = "descriptionAuth_method"
    type  = "string"
    value = var.castleconnector_property_description_auth_method
  }
  property {
    name  = "descriptionAuth_methodDetails"
    type  = "string"
    value = var.castleconnector_property_description_auth_method_details
  }
  property {
    name  = "descriptionChangeset"
    type  = "string"
    value = var.castleconnector_property_description_changeset
  }
  property {
    name  = "descriptionChangesetDetail"
    type  = "string"
    value = var.castleconnector_property_description_changeset_detail
  }
  property {
    name  = "descriptionCustom"
    type  = "string"
    value = var.castleconnector_property_description_custom
  }
  property {
    name  = "descriptionHeader"
    type  = "string"
    value = var.castleconnector_property_description_header
  }
  property {
    name  = "descriptionSession"
    type  = "string"
    value = var.castleconnector_property_description_session
  }
  property {
    name  = "descriptionTransaction"
    type  = "string"
    value = var.castleconnector_property_description_transaction
  }
  property {
    name  = "descriptionTransactionAmount"
    type  = "string"
    value = var.castleconnector_property_description_transaction_amount
  }
  property {
    name  = "descriptionTransactionBilling"
    type  = "string"
    value = var.castleconnector_property_description_transaction_billing
  }
  property {
    name  = "descriptionTransactionCard"
    type  = "string"
    value = var.castleconnector_property_description_transaction_card
  }
  property {
    name  = "descriptionTransactionMerchant"
    type  = "string"
    value = var.castleconnector_property_description_transaction_merchant
  }
  property {
    name  = "descriptionTransactionMerchantAddrress"
    type  = "string"
    value = var.castleconnector_property_description_transaction_merchant_addrress
  }
  property {
    name  = "descriptionTransactionPayment"
    type  = "string"
    value = var.castleconnector_property_description_transaction_payment
  }
  property {
    name  = "descriptionTransactionshipping"
    type  = "string"
    value = var.castleconnector_property_description_transactionshipping
  }
  property {
    name  = "descriptionUser"
    type  = "string"
    value = var.castleconnector_property_description_user
  }
  property {
    name  = "descriptionUserAddress"
    type  = "string"
    value = var.castleconnector_property_description_user_address
  }
  property {
    name  = "email"
    type  = "string"
    value = var.castleconnector_property_email
  }
  property {
    name  = "emailFrom"
    type  = "string"
    value = var.castleconnector_property_email_from
  }
  property {
    name  = "emailFrom_reset"
    type  = "string"
    value = var.castleconnector_property_email_from_reset
  }
  property {
    name  = "emailTo"
    type  = "string"
    value = var.castleconnector_property_email_to
  }
  property {
    name  = "emailTo_reset"
    type  = "string"
    value = var.castleconnector_property_email_to_reset
  }
  property {
    name  = "fingerprint"
    type  = "string"
    value = var.castleconnector_property_fingerprint
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.castleconnector_property_headers
  }
  property {
    name  = "host"
    type  = "string"
    value = var.castleconnector_property_host
  }
  property {
    name  = "ip"
    type  = "string"
    value = var.castleconnector_property_ip
  }
  property {
    name  = "line1"
    type  = "string"
    value = var.castleconnector_property_line1
  }
  property {
    name  = "line2"
    type  = "string"
    value = var.castleconnector_property_line2
  }
  property {
    name  = "matchingUserId"
    type  = "string"
    value = var.castleconnector_property_matching_user_id
  }
  property {
    name  = "merchantCode"
    type  = "string"
    value = var.castleconnector_property_merchant_code
  }
  property {
    name  = "merchantDesc"
    type  = "string"
    value = var.castleconnector_property_merchant_desc
  }
  property {
    name  = "merchantId"
    type  = "string"
    value = var.castleconnector_property_merchant_id
  }
  property {
    name  = "merchantName"
    type  = "string"
    value = var.castleconnector_property_merchant_name
  }
  property {
    name  = "payment_methodBankName"
    type  = "string"
    value = var.castleconnector_property_payment_method_bank_name
  }
  property {
    name  = "payment_methodCountryCode"
    type  = "string"
    value = var.castleconnector_property_payment_method_country_code
  }
  property {
    name  = "payment_methodFingerprint"
    type  = "string"
    value = var.castleconnector_property_payment_method_fingerprint
  }
  property {
    name  = "payment_methodHoldeName"
    type  = "string"
    value = var.castleconnector_property_payment_method_holde_name
  }
  property {
    name  = "payment_methodType"
    type  = "string"
    value = var.castleconnector_property_payment_method_type
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.castleconnector_property_phone
  }
  property {
    name  = "phoneFrom"
    type  = "string"
    value = var.castleconnector_property_phone_from
  }
  property {
    name  = "phoneFrom_reset"
    type  = "string"
    value = var.castleconnector_property_phone_from_reset
  }
  property {
    name  = "phoneTo"
    type  = "string"
    value = var.castleconnector_property_phone_to
  }
  property {
    name  = "phoneTo_reset"
    type  = "string"
    value = var.castleconnector_property_phone_to_reset
  }
  property {
    name  = "postal_code"
    type  = "string"
    value = var.castleconnector_property_postal_code
  }
  property {
    name  = "productId"
    type  = "string"
    value = var.castleconnector_property_product_id
  }
  property {
    name  = "queryParameter"
    type  = "string"
    value = var.castleconnector_property_query_parameter
  }
  property {
    name  = "queryParameters_reset"
    type  = "string"
    value = var.castleconnector_property_query_parameters_reset
  }
  property {
    name  = "queryParameters_update"
    type  = "string"
    value = var.castleconnector_property_query_parameters_update
  }
  property {
    name  = "region_code"
    type  = "string"
    value = var.castleconnector_property_region_code
  }
  property {
    name  = "registered_at"
    type  = "string"
    value = var.castleconnector_property_registered_at
  }
  property {
    name  = "request_token"
    type  = "string"
    value = var.castleconnector_property_request_token
  }
  property {
    name  = "sessionId"
    type  = "string"
    value = var.castleconnector_property_session_id
  }
  property {
    name  = "session_created_at"
    type  = "string"
    value = var.castleconnector_property_session_created_at
  }
  property {
    name  = "shipping_address_city"
    type  = "string"
    value = var.castleconnector_property_shipping_address_city
  }
  property {
    name  = "shipping_address_countrycode"
    type  = "string"
    value = var.castleconnector_property_shipping_address_countrycode
  }
  property {
    name  = "shipping_address_fingerprint"
    type  = "string"
    value = var.castleconnector_property_shipping_address_fingerprint
  }
  property {
    name  = "shipping_address_line1"
    type  = "string"
    value = var.castleconnector_property_shipping_address_line1
  }
  property {
    name  = "shipping_address_line2"
    type  = "string"
    value = var.castleconnector_property_shipping_address_line2
  }
  property {
    name  = "shipping_address_postalCode"
    type  = "string"
    value = var.castleconnector_property_shipping_address_postal_code
  }
  property {
    name  = "shipping_address_regionCode"
    type  = "string"
    value = var.castleconnector_property_shipping_address_region_code
  }
  property {
    name  = "skip_context"
    type  = "string"
    value = var.castleconnector_property_skip_context
  }
  property {
    name  = "skip_context_validation_filter"
    type  = "string"
    value = var.castleconnector_property_skip_context_validation_filter
  }
  property {
    name  = "skip_request_token"
    type  = "string"
    value = var.castleconnector_property_skip_request_token
  }
  property {
    name  = "skip_request_token_validation_filter"
    type  = "string"
    value = var.castleconnector_property_skip_request_token_validation_filter
  }
  property {
    name  = "statusFilter_challenge"
    type  = "string"
    value = var.castleconnector_property_status_filter_challenge
  }
  property {
    name  = "statusFilter_login"
    type  = "string"
    value = var.castleconnector_property_status_filter_login
  }
  property {
    name  = "statusFilter_reg"
    type  = "string"
    value = var.castleconnector_property_status_filter_reg
  }
  property {
    name  = "statusFilter_reset"
    type  = "string"
    value = var.castleconnector_property_status_filter_reset
  }
  property {
    name  = "status_challenge"
    type  = "string"
    value = var.castleconnector_property_status_challenge
  }
  property {
    name  = "status_login"
    type  = "string"
    value = var.castleconnector_property_status_login
  }
  property {
    name  = "status_logout"
    type  = "string"
    value = var.castleconnector_property_status_logout
  }
  property {
    name  = "status_reg"
    type  = "string"
    value = var.castleconnector_property_status_reg
  }
  property {
    name  = "status_reset"
    type  = "string"
    value = var.castleconnector_property_status_reset
  }
  property {
    name  = "status_transaction"
    type  = "string"
    value = var.castleconnector_property_status_transaction
  }
  property {
    name  = "status_update"
    type  = "string"
    value = var.castleconnector_property_status_update
  }
  property {
    name  = "traits"
    type  = "string"
    value = var.castleconnector_property_traits
  }
  property {
    name  = "transactionCurrency"
    type  = "string"
    value = var.castleconnector_property_transaction_currency
  }
  property {
    name  = "transactionId"
    type  = "string"
    value = var.castleconnector_property_transaction_id
  }
  property {
    name  = "transactionType"
    type  = "string"
    value = var.castleconnector_property_transaction_type
  }
  property {
    name  = "type"
    type  = "string"
    value = var.castleconnector_property_type
  }
  property {
    name  = "typeFilter"
    type  = "string"
    value = var.castleconnector_property_type_filter
  }
  property {
    name  = "userAgent"
    type  = "string"
    value = var.castleconnector_property_user_agent
  }
  property {
    name  = "userID"
    type  = "string"
    value = var.castleconnector_property_user_id
  }
  property {
    name  = "username"
    type  = "string"
    value = var.castleconnector_property_username
  }
  property {
    name  = "variant"
    type  = "string"
    value = var.castleconnector_property_variant
  }
  property {
    name  = "variantFilter"
    type  = "string"
    value = var.castleconnector_property_variant_filter
  }
  property {
    name  = "variantFilter_chg"
    type  = "string"
    value = var.castleconnector_property_variant_filter_chg
  }
  property {
    name  = "variantFilter_reg"
    type  = "string"
    value = var.castleconnector_property_variant_filter_reg
  }
  property {
    name  = "variant_chg"
    type  = "string"
    value = var.castleconnector_property_variant_chg
  }
  property {
    name  = "variant_reg"
    type  = "string"
    value = var.castleconnector_property_variant_reg
  }
}
