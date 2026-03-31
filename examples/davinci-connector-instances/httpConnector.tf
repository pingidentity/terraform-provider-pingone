resource "pingone_davinci_connector_instance" "httpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "httpConnector"
  }
  name = "My awesome httpConnector"
  property {
    name  = "additionalFieldsName"
    type  = "string"
    value = var.httpconnector_property_additional_fields_name
  }
  property {
    name  = "blockRedirects"
    type  = "string"
    value = var.httpconnector_property_block_redirects
  }
  property {
    name  = "bodyHeaderText"
    type  = "string"
    value = var.httpconnector_property_body_header_text
  }
  property {
    name  = "bodyParams"
    type  = "string"
    value = var.httpconnector_property_body_params
  }
  property {
    name  = "button"
    type  = "string"
    value = var.httpconnector_property_button
  }
  property {
    name  = "challenge"
    type  = "string"
    value = var.httpconnector_property_challenge
  }
  property {
    name  = "claimsNameValuePairs"
    type  = "string"
    value = var.httpconnector_property_claims_name_value_pairs
  }
  property {
    name  = "connectionId"
    type  = "string"
    value = var.httpconnector_property_connection_id
  }
  property {
    name  = "connectionInstanceId"
    type  = "string"
    value = var.httpconnector_property_connection_instance_id
  }
  property {
    name  = "contentSource"
    type  = "string"
    value = var.httpconnector_property_content_source
  }
  property {
    name  = "contentType"
    type  = "string"
    value = var.httpconnector_property_content_type
  }
  property {
    name  = "contentTypeError"
    type  = "string"
    value = var.httpconnector_property_content_type_error
  }
  property {
    name  = "customCSS"
    type  = "string"
    value = var.httpconnector_property_custom_css
  }
  property {
    name  = "customCSSByCdn"
    type  = "string"
    value = var.httpconnector_property_custom_cssby_cdn
  }
  property {
    name  = "customHTML"
    type  = "string"
    value = var.httpconnector_property_custom_html
  }
  property {
    name  = "customHTMLByCdn"
    type  = "string"
    value = var.httpconnector_property_custom_htmlby_cdn
  }
  property {
    name  = "customScript"
    type  = "string"
    value = var.httpconnector_property_custom_script
  }
  property {
    name  = "customScriptByCdn"
    type  = "string"
    value = var.httpconnector_property_custom_script_by_cdn
  }
  property {
    name  = "customUrl"
    type  = "string"
    value = var.httpconnector_property_custom_url
  }
  property {
    name  = "delayTime"
    type  = "string"
    value = var.httpconnector_property_delay_time
  }
  property {
    name  = "description"
    type  = "string"
    value = var.httpconnector_property_description
  }
  property {
    name  = "enablePolling"
    type  = "string"
    value = var.httpconnector_property_enable_polling
  }
  property {
    name  = "fieldValidation"
    type  = "string"
    value = var.httpconnector_property_field_validation
  }
  property {
    name  = "formFieldsList"
    type  = "string"
    value = var.httpconnector_property_form_fields_list
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.httpconnector_property_headers
  }
  property {
    name  = "httpBody"
    type  = "string"
    value = var.httpconnector_property_http_body
  }
  property {
    name  = "httpHeaders"
    type  = "string"
    value = var.httpconnector_property_http_headers
  }
  property {
    name  = "httpMethod"
    type  = "string"
    value = var.httpconnector_property_http_method
  }
  property {
    name  = "httpStatusCode"
    type  = "string"
    value = var.httpconnector_property_http_status_code
  }
  property {
    name  = "inputSchema"
    type  = "string"
    value = var.httpconnector_property_input_schema
  }
  property {
    name  = "jsonString"
    type  = "string"
    value = var.httpconnector_property_json_string
  }
  property {
    name  = "keepOutputIfNotValid"
    type  = "string"
    value = var.httpconnector_property_keep_output_if_not_valid
  }
  property {
    name  = "message"
    type  = "string"
    value = var.httpconnector_property_message
  }
  property {
    name  = "messageIcon"
    type  = "string"
    value = var.httpconnector_property_message_icon
  }
  property {
    name  = "messageIconHeight"
    type  = "string"
    value = var.httpconnector_property_message_icon_height
  }
  property {
    name  = "messageTitle"
    type  = "string"
    value = var.httpconnector_property_message_title
  }
  property {
    name  = "nextButtonText"
    type  = "string"
    value = var.httpconnector_property_next_button_text
  }
  property {
    name  = "nextEvent"
    type  = "string"
    value = var.httpconnector_property_next_event
  }
  property {
    name  = "outboundMtlsKey"
    type  = "string"
    value = var.httpconnector_property_outbound_mtls_key
  }
  property {
    name  = "outputSchema"
    type  = "string"
    value = var.httpconnector_property_output_schema
  }
  property {
    name  = "outputSchemaError"
    type  = "string"
    value = var.httpconnector_property_output_schema_error
  }
  property {
    name  = "pollChallengeStatus"
    type  = "string"
    value = var.httpconnector_property_poll_challenge_status
  }
  property {
    name  = "pollInterval"
    type  = "string"
    value = var.httpconnector_property_poll_interval
  }
  property {
    name  = "pollRetries"
    type  = "string"
    value = var.httpconnector_property_poll_retries
  }
  property {
    name  = "queryParams"
    type  = "string"
    value = var.httpconnector_property_query_params
  }
  property {
    name  = "raw"
    type  = "string"
    value = var.httpconnector_property_raw
  }
  property {
    name  = "recaptchaSecretKey"
    type  = "string"
    value = var.httpconnector_property_recaptcha_secret_key
  }
  property {
    name  = "recaptchaSiteKey"
    type  = "string"
    value = var.httpconnector_property_recaptcha_site_key
  }
  property {
    name  = "returnRequestParameters"
    type  = "string"
    value = var.httpconnector_property_return_request_parameters
  }
  property {
    name  = "returnSuccess"
    type  = "string"
    value = var.httpconnector_property_return_success
  }
  property {
    name  = "showContinueButton"
    type  = "string"
    value = var.httpconnector_property_show_continue_button
  }
  property {
    name  = "showFooter"
    type  = "string"
    value = var.httpconnector_property_show_footer
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.httpconnector_property_show_powered_by
  }
  property {
    name  = "signResponse"
    type  = "string"
    value = var.httpconnector_property_sign_response
  }
  property {
    name  = "sktemplate"
    type  = "string"
    value = var.httpconnector_property_sktemplate
  }
  property {
    name  = "timeout"
    type  = "string"
    value = var.httpconnector_property_timeout
  }
  property {
    name  = "title"
    type  = "string"
    value = var.httpconnector_property_title
  }
  property {
    name  = "unsafeIgnoreTLSErrors"
    type  = "string"
    value = var.httpconnector_property_unsafe_ignore_tlserrors
  }
  property {
    name  = "url"
    type  = "string"
    value = var.httpconnector_property_url
  }
  property {
    name  = "useRecaptcha"
    type  = "string"
    value = var.httpconnector_property_use_recaptcha
  }
  property {
    name  = "validationRules"
    type  = "string"
    value = var.httpconnector_property_validation_rules
  }
  property {
    name  = "validity"
    type  = "string"
    value = var.httpconnector_property_validity
  }
  property {
    name  = "whiteList"
    type  = "string"
    value = var.httpconnector_property_white_list
  }
}
