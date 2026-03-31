resource "pingone_davinci_connector_instance" "iovationConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "iovationConnector"
  }
  name = "My awesome iovationConnector"
  property {
    name  = "achRoutingNumber"
    type  = "string"
    value = var.iovationconnector_property_ach_routing_number
  }
  property {
    name  = "alternateIp"
    type  = "string"
    value = var.iovationconnector_property_alternate_ip
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.iovationconnector_property_api_url
  }
  property {
    name  = "billingCity"
    type  = "string"
    value = var.iovationconnector_property_billing_city
  }
  property {
    name  = "billingCountry"
    type  = "string"
    value = var.iovationconnector_property_billing_country
  }
  property {
    name  = "billingPostalCode"
    type  = "string"
    value = var.iovationconnector_property_billing_postal_code
  }
  property {
    name  = "billingRegion"
    type  = "string"
    value = var.iovationconnector_property_billing_region
  }
  property {
    name  = "billingShippingMismatch"
    type  = "string"
    value = var.iovationconnector_property_billing_shipping_mismatch
  }
  property {
    name  = "billingStreet"
    type  = "string"
    value = var.iovationconnector_property_billing_street
  }
  property {
    name  = "blackBox"
    type  = "string"
    value = var.iovationconnector_property_black_box
  }
  property {
    name  = "creditCardBin"
    type  = "string"
    value = var.iovationconnector_property_credit_card_bin
  }
  property {
    name  = "email"
    type  = "string"
    value = var.iovationconnector_property_email
  }
  property {
    name  = "emailVerified"
    type  = "string"
    value = var.iovationconnector_property_email_verified
  }
  property {
    name  = "eventId"
    type  = "string"
    value = var.iovationconnector_property_event_id
  }
  property {
    name  = "homePhoneNumber"
    type  = "string"
    value = var.iovationconnector_property_home_phone_number
  }
  property {
    name  = "homePhoneSmsEnabled"
    type  = "string"
    value = var.iovationconnector_property_home_phone_sms_enabled
  }
  property {
    name  = "homePhoneVerified"
    type  = "string"
    value = var.iovationconnector_property_home_phone_verified
  }
  property {
    name  = "integrationPoint"
    type  = "string"
    value = var.iovationconnector_property_integration_point
  }
  property {
    name  = "javascriptCdnUrl"
    type  = "string"
    value = var.iovationconnector_property_javascript_cdn_url
  }
  property {
    name  = "loadingText"
    type  = "string"
    value = var.iovationconnector_property_loading_text
  }
  property {
    name  = "mobilePhoneNumber"
    type  = "string"
    value = var.iovationconnector_property_mobile_phone_number
  }
  property {
    name  = "mobilePhoneSmsEnabled"
    type  = "string"
    value = var.iovationconnector_property_mobile_phone_sms_enabled
  }
  property {
    name  = "mobilePhoneVerified"
    type  = "string"
    value = var.iovationconnector_property_mobile_phone_verified
  }
  property {
    name  = "nextEvent"
    type  = "string"
    value = var.iovationconnector_property_next_event
  }
  property {
    name  = "officePhoneNumber"
    type  = "string"
    value = var.iovationconnector_property_office_phone_number
  }
  property {
    name  = "officePhoneSmsEnabled"
    type  = "string"
    value = var.iovationconnector_property_office_phone_sms_enabled
  }
  property {
    name  = "officePhoneVerified"
    type  = "string"
    value = var.iovationconnector_property_office_phone_verified
  }
  property {
    name  = "onlineId"
    type  = "string"
    value = var.iovationconnector_property_online_id
  }
  property {
    name  = "referrerUrl"
    type  = "string"
    value = var.iovationconnector_property_referrer_url
  }
  property {
    name  = "securityPin"
    type  = "string"
    value = var.iovationconnector_property_security_pin
  }
  property {
    name  = "securityPinType"
    type  = "string"
    value = var.iovationconnector_property_security_pin_type
  }
  property {
    name  = "serviceId"
    type  = "string"
    value = var.iovationconnector_property_service_id
  }
  property {
    name  = "sessionAlias"
    type  = "string"
    value = var.iovationconnector_property_session_alias
  }
  property {
    name  = "shippingCity"
    type  = "string"
    value = var.iovationconnector_property_shipping_city
  }
  property {
    name  = "shippingCountry"
    type  = "string"
    value = var.iovationconnector_property_shipping_country
  }
  property {
    name  = "shippingPostalCode"
    type  = "string"
    value = var.iovationconnector_property_shipping_postal_code
  }
  property {
    name  = "shippingRegion"
    type  = "string"
    value = var.iovationconnector_property_shipping_region
  }
  property {
    name  = "shippingStreet"
    type  = "string"
    value = var.iovationconnector_property_shipping_street
  }
  property {
    name  = "sku"
    type  = "string"
    value = var.iovationconnector_property_sku
  }
  property {
    name  = "subKey"
    type  = "string"
    value = var.iovationconnector_property_sub_key
  }
  property {
    name  = "subscriberAccount"
    type  = "string"
    value = var.iovationconnector_property_subscriber_account
  }
  property {
    name  = "subscriberId"
    type  = "string"
    value = var.iovationconnector_property_subscriber_id
  }
  property {
    name  = "subscriberPasscode"
    type  = "string"
    value = var.iovationconnector_property_subscriber_passcode
  }
  property {
    name  = "tenantId"
    type  = "string"
    value = var.iovationconnector_property_tenant_id
  }
  property {
    name  = "upc"
    type  = "string"
    value = var.iovationconnector_property_upc
  }
  property {
    name  = "userAccountCode"
    type  = "string"
    value = var.iovationconnector_property_user_account_code
  }
  property {
    name  = "userIp"
    type  = "string"
    value = var.iovationconnector_property_user_ip
  }
  property {
    name  = "userReference"
    type  = "string"
    value = var.iovationconnector_property_user_reference
  }
  property {
    name  = "valueAmount"
    type  = "string"
    value = var.iovationconnector_property_value_amount
  }
  property {
    name  = "valueCurrency"
    type  = "string"
    value = var.iovationconnector_property_value_currency
  }
  property {
    name  = "version"
    type  = "string"
    value = var.iovationconnector_property_version
  }
}
