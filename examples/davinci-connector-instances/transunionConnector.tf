resource "pingone_davinci_connector_instance" "transunionConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "transunionConnector"
  }
  name = "My awesome transunionConnector"
  property {
    name  = "AccountNumber"
    type  = "string"
    value = var.transunionconnector_property_account_number
  }
  property {
    name  = "Address1"
    type  = "string"
    value = var.transunionconnector_property_address1
  }
  property {
    name  = "Address2"
    type  = "string"
    value = var.transunionconnector_property_address2
  }
  property {
    name  = "Address3"
    type  = "string"
    value = var.transunionconnector_property_address3
  }
  property {
    name  = "AddressType"
    type  = "string"
    value = var.transunionconnector_property_address_type
  }
  property {
    name  = "AuthURL"
    type  = "string"
    value = var.transunionconnector_property_auth_url
  }
  property {
    name  = "BINNumber"
    type  = "string"
    value = var.transunionconnector_property_binnumber
  }
  property {
    name  = "CaseId"
    type  = "string"
    value = var.transunionconnector_property_case_id
  }
  property {
    name  = "CaseStatusType"
    type  = "string"
    value = var.transunionconnector_property_case_status_type
  }
  property {
    name  = "CaseType"
    type  = "string"
    value = var.transunionconnector_property_case_type
  }
  property {
    name  = "City"
    type  = "string"
    value = var.transunionconnector_property_city
  }
  property {
    name  = "Comment"
    type  = "string"
    value = var.transunionconnector_property_comment
  }
  property {
    name  = "CountryCode"
    type  = "string"
    value = var.transunionconnector_property_country_code
  }
  property {
    name  = "Currency"
    type  = "string"
    value = var.transunionconnector_property_currency
  }
  property {
    name  = "CustomerId"
    type  = "string"
    value = var.transunionconnector_property_customer_id
  }
  property {
    name  = "DateOfBirth"
    type  = "string"
    value = var.transunionconnector_property_date_of_birth
  }
  property {
    name  = "Emails"
    type  = "string"
    value = var.transunionconnector_property_emails
  }
  property {
    name  = "FirstName"
    type  = "string"
    value = var.transunionconnector_property_first_name
  }
  property {
    name  = "GeoLocation"
    type  = "string"
    value = var.transunionconnector_property_geo_location
  }
  property {
    name  = "IdentityConsentId"
    type  = "string"
    value = var.transunionconnector_property_identity_consent_id
  }
  property {
    name  = "ItemValue"
    type  = "number"
    value = var.transunionconnector_property_item_value
  }
  property {
    name  = "LastName"
    type  = "string"
    value = var.transunionconnector_property_last_name
  }
  property {
    name  = "LocationConsentId"
    type  = "string"
    value = var.transunionconnector_property_location_consent_id
  }
  property {
    name  = "LongTermAccessToken"
    type  = "string"
    value = var.transunionconnector_property_long_term_access_token
  }
  property {
    name  = "LongTermAccessTokenExpiry"
    type  = "string"
    value = var.transunionconnector_property_long_term_access_token_expiry
  }
  property {
    name  = "MarketType"
    type  = "string"
    value = var.transunionconnector_property_market_type
  }
  property {
    name  = "Method"
    type  = "string"
    value = var.transunionconnector_property_method
  }
  property {
    name  = "OTPLanguageEnum"
    type  = "string"
    value = var.transunionconnector_property_otplanguage_enum
  }
  property {
    name  = "OTPStatus"
    type  = "string"
    value = var.transunionconnector_property_otpstatus
  }
  property {
    name  = "Passcode"
    type  = "string"
    value = var.transunionconnector_property_passcode
  }
  property {
    name  = "PaymentType"
    type  = "string"
    value = var.transunionconnector_property_payment_type
  }
  property {
    name  = "PhoneDeliveryType"
    type  = "string"
    value = var.transunionconnector_property_phone_delivery_type
  }
  property {
    name  = "PhoneNumber"
    type  = "string"
    value = var.transunionconnector_property_phone_number
  }
  property {
    name  = "PostalCode"
    type  = "string"
    value = var.transunionconnector_property_postal_code
  }
  property {
    name  = "Secret"
    type  = "string"
    value = var.transunionconnector_property_secret
  }
  property {
    name  = "ShortTermAccessToken"
    type  = "string"
    value = var.transunionconnector_property_short_term_access_token
  }
  property {
    name  = "ShortTermAccessTokenExpiry"
    type  = "string"
    value = var.transunionconnector_property_short_term_access_token_expiry
  }
  property {
    name  = "SocialId"
    type  = "number"
    value = var.transunionconnector_property_social_id
  }
  property {
    name  = "SocialNetworkType"
    type  = "string"
    value = var.transunionconnector_property_social_network_type
  }
  property {
    name  = "SocialSecurityNumber"
    type  = "string"
    value = var.transunionconnector_property_social_security_number
  }
  property {
    name  = "State"
    type  = "string"
    value = var.transunionconnector_property_state
  }
  property {
    name  = "TimeToFulfilment"
    type  = "string"
    value = var.transunionconnector_property_time_to_fulfilment
  }
  property {
    name  = "TotalTransactionValue"
    type  = "number"
    value = var.transunionconnector_property_total_transaction_value
  }
  property {
    name  = "TransactionItemName"
    type  = "string"
    value = var.transunionconnector_property_transaction_item_name
  }
  property {
    name  = "TransactionItemQuantity"
    type  = "number"
    value = var.transunionconnector_property_transaction_item_quantity
  }
  property {
    name  = "TransactionProcessingTime"
    type  = "string"
    value = var.transunionconnector_property_transaction_processing_time
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.transunionconnector_property_api_url
  }
  property {
    name  = "customCSS"
    type  = "string"
    value = var.transunionconnector_property_custom_css
  }
  property {
    name  = "customHTML"
    type  = "string"
    value = var.transunionconnector_property_custom_html
  }
  property {
    name  = "customScript"
    type  = "string"
    value = var.transunionconnector_property_custom_script
  }
  property {
    name  = "docVerificationPassword"
    type  = "string"
    value = var.transunionconnector_property_doc_verification_password
  }
  property {
    name  = "docVerificationPublicKey"
    type  = "string"
    value = var.transunionconnector_property_doc_verification_public_key
  }
  property {
    name  = "docVerificationSecret"
    type  = "string"
    value = var.transunionconnector_property_doc_verification_secret
  }
  property {
    name  = "docVerificationSiteId"
    type  = "string"
    value = var.transunionconnector_property_doc_verification_site_id
  }
  property {
    name  = "docVerificationUsername"
    type  = "string"
    value = var.transunionconnector_property_doc_verification_username
  }
  property {
    name  = "documentBackCaptureAttempts"
    type  = "number"
    value = var.transunionconnector_property_document_back_capture_attempts
  }
  property {
    name  = "documentEnableFaceDetection"
    type  = "boolean"
    value = var.transunionconnector_property_document_enable_face_detection
  }
  property {
    name  = "documentFrontCaptureAttempts"
    type  = "number"
    value = var.transunionconnector_property_document_front_capture_attempts
  }
  property {
    name  = "documentGoodImageFoundADA"
    type  = "string"
    value = var.transunionconnector_property_document_good_image_found_ada
  }
  property {
    name  = "documentIsBarcodeDetectedEnabled"
    type  = "boolean"
    value = var.transunionconnector_property_document_is_barcode_detected_enabled
  }
  property {
    name  = "documentManualTimeout"
    type  = "number"
    value = var.transunionconnector_property_document_manual_timeout
  }
  property {
    name  = "documentNotGoodImageADA"
    type  = "string"
    value = var.transunionconnector_property_document_not_good_image_ada
  }
  property {
    name  = "documentNotificationADA"
    type  = "string"
    value = var.transunionconnector_property_document_notification_ada
  }
  property {
    name  = "documentOverlayColor"
    type  = "string"
    value = var.transunionconnector_property_document_overlay_color
  }
  property {
    name  = "documentOverlayText"
    type  = "string"
    value = var.transunionconnector_property_document_overlay_text
  }
  property {
    name  = "documentOverlayTextAuto"
    type  = "string"
    value = var.transunionconnector_property_document_overlay_text_auto
  }
  property {
    name  = "documentOverlayTextAutoADA"
    type  = "string"
    value = var.transunionconnector_property_document_overlay_text_auto_ada
  }
  property {
    name  = "documentOverlayTextManualADA"
    type  = "string"
    value = var.transunionconnector_property_document_overlay_text_manual_ada
  }
  property {
    name  = "frontCaptureMode"
    type  = "string"
    value = var.transunionconnector_property_front_capture_mode
  }
  property {
    name  = "htmlConfig"
    type  = "string"
    value = var.transunionconnector_property_html_config
  }
  property {
    name  = "idVerificationPassword"
    type  = "string"
    value = var.transunionconnector_property_id_verification_password
  }
  property {
    name  = "idVerificationPayloadUrl"
    type  = "string"
    value = var.transunionconnector_property_id_verification_payload_url
  }
  property {
    name  = "idVerificationPublicKey"
    type  = "string"
    value = var.transunionconnector_property_id_verification_public_key
  }
  property {
    name  = "idVerificationSecret"
    type  = "string"
    value = var.transunionconnector_property_id_verification_secret
  }
  property {
    name  = "idVerificationSiteId"
    type  = "string"
    value = var.transunionconnector_property_id_verification_site_id
  }
  property {
    name  = "idVerificationUsername"
    type  = "string"
    value = var.transunionconnector_property_id_verification_username
  }
  property {
    name  = "kbaPassword"
    type  = "string"
    value = var.transunionconnector_property_kba_password
  }
  property {
    name  = "kbaPublicKey"
    type  = "string"
    value = var.transunionconnector_property_kba_public_key
  }
  property {
    name  = "kbaSecret"
    type  = "string"
    value = var.transunionconnector_property_kba_secret
  }
  property {
    name  = "kbaSiteId"
    type  = "string"
    value = var.transunionconnector_property_kba_site_id
  }
  property {
    name  = "kbaUsername"
    type  = "string"
    value = var.transunionconnector_property_kba_username
  }
  property {
    name  = "licenseFaceDetectionProportionMax"
    type  = "string"
    value = var.transunionconnector_property_license_face_detection_proportion_max
  }
  property {
    name  = "licenseFaceDetectionProportionMin"
    type  = "number"
    value = var.transunionconnector_property_license_face_detection_proportion_min
  }
  property {
    name  = "nativeBackFocusThreshold"
    type  = "number"
    value = var.transunionconnector_property_native_back_focus_threshold
  }
  property {
    name  = "nativeBackGlareThreshold"
    type  = "number"
    value = var.transunionconnector_property_native_back_glare_threshold
  }
  property {
    name  = "nativeFrontFocusThreshold"
    type  = "number"
    value = var.transunionconnector_property_native_front_focus_threshold
  }
  property {
    name  = "nativeFrontGlareThreshold"
    type  = "number"
    value = var.transunionconnector_property_native_front_glare_threshold
  }
  property {
    name  = "otpPassword"
    type  = "string"
    value = var.transunionconnector_property_otp_password
  }
  property {
    name  = "otpPublicKey"
    type  = "string"
    value = var.transunionconnector_property_otp_public_key
  }
  property {
    name  = "otpSecret"
    type  = "string"
    value = var.transunionconnector_property_otp_secret
  }
  property {
    name  = "otpSiteId"
    type  = "string"
    value = var.transunionconnector_property_otp_site_id
  }
  property {
    name  = "otpUsername"
    type  = "string"
    value = var.transunionconnector_property_otp_username
  }
  property {
    name  = "passportFaceDetectionProportionMax"
    type  = "string"
    value = var.transunionconnector_property_passport_face_detection_proportion_max
  }
  property {
    name  = "passportFaceDetectionProportionMin"
    type  = "number"
    value = var.transunionconnector_property_passport_face_detection_proportion_min
  }
  property {
    name  = "passportFrontCaptureAttempts"
    type  = "number"
    value = var.transunionconnector_property_passport_front_capture_attempts
  }
  property {
    name  = "screen0Config"
    type  = "string"
    value = var.transunionconnector_property_screen0_config
  }
  property {
    name  = "selfieCaptureAttempts"
    type  = "number"
    value = var.transunionconnector_property_selfie_capture_attempts
  }
  property {
    name  = "selfieCaptureMode"
    type  = "string"
    value = var.transunionconnector_property_selfie_capture_mode
  }
  property {
    name  = "selfieEnableFaceDetection"
    type  = "boolean"
    value = var.transunionconnector_property_selfie_enable_face_detection
  }
  property {
    name  = "selfieGoodImageFoundADA"
    type  = "string"
    value = var.transunionconnector_property_selfie_good_image_found_ada
  }
  property {
    name  = "selfieManualTimeout"
    type  = "number"
    value = var.transunionconnector_property_selfie_manual_timeout
  }
  property {
    name  = "selfieNotGoodImageADA"
    type  = "string"
    value = var.transunionconnector_property_selfie_not_good_image_ada
  }
  property {
    name  = "selfieNotificationADA"
    type  = "string"
    value = var.transunionconnector_property_selfie_notification_ada
  }
  property {
    name  = "selfieOrientationErrorText"
    type  = "string"
    value = var.transunionconnector_property_selfie_orientation_error_text
  }
  property {
    name  = "selfieOverlayText"
    type  = "string"
    value = var.transunionconnector_property_selfie_overlay_text
  }
  property {
    name  = "selfieOverlayTextAuto"
    type  = "string"
    value = var.transunionconnector_property_selfie_overlay_text_auto
  }
  property {
    name  = "selfieOverlayTextAutoADA"
    type  = "string"
    value = var.transunionconnector_property_selfie_overlay_text_auto_ada
  }
  property {
    name  = "selfieOverlayTextManualADA"
    type  = "string"
    value = var.transunionconnector_property_selfie_overlay_text_manual_ada
  }
  property {
    name  = "selfieUseBackCamera"
    type  = "boolean"
    value = var.transunionconnector_property_selfie_use_back_camera
  }
  property {
    name  = "transunionAuthType"
    type  = "string"
    value = var.transunionconnector_property_transunion_auth_type
  }
  property {
    name  = "useCustomScreens"
    type  = "string"
    value = var.transunionconnector_property_use_custom_screens
  }
  property {
    name  = "validationBadFocusText"
    type  = "string"
    value = var.transunionconnector_property_validation_bad_focus_text
  }
  property {
    name  = "validationBadGlareText"
    type  = "string"
    value = var.transunionconnector_property_validation_bad_glare_text
  }
  property {
    name  = "validationDefaultText"
    type  = "string"
    value = var.transunionconnector_property_validation_default_text
  }
  property {
    name  = "validationNoBarcodeText"
    type  = "string"
    value = var.transunionconnector_property_validation_no_barcode_text
  }
  property {
    name  = "validationNoFacesFoundText"
    type  = "string"
    value = var.transunionconnector_property_validation_no_faces_found_text
  }
  property {
    name  = "validationOkButtonColor"
    type  = "string"
    value = var.transunionconnector_property_validation_ok_button_color
  }
  property {
    name  = "validationOkButtonText"
    type  = "string"
    value = var.transunionconnector_property_validation_ok_button_text
  }
}
