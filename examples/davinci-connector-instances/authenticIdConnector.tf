resource "pingone_davinci_connector_instance" "authenticIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "authenticIdConnector"
  }
  name = "My awesome authenticIdConnector"
  property {
    name  = "accountAccessKey"
    type  = "string"
    value = var.authenticidconnector_property_account_access_key
  }
  property {
    name  = "accountCode"
    type  = "string"
    value = var.authenticidconnector_property_account_code
  }
  property {
    name  = "androidSDKLicenseKey"
    type  = "string"
    value = var.authenticidconnector_property_android_sdk_license_key
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.authenticidconnector_property_api_url
  }
  property {
    name  = "applicationLandscapeErrorMessage"
    type  = "string"
    value = var.authenticidconnector_property_application_landscape_error_message
  }
  property {
    name  = "authenticIdDocumentType"
    type  = "string"
    value = var.authenticidconnector_property_authentic_id_document_type
  }
  property {
    name  = "authenticIdDocumentTypeV2"
    type  = "string"
    value = var.authenticidconnector_property_authentic_id_document_type_v2
  }
  property {
    name  = "backCaptureAttempt"
    type  = "number"
    value = var.authenticidconnector_property_back_capture_attempt
  }
  property {
    name  = "backCaptureMode"
    type  = "string"
    value = var.authenticidconnector_property_back_capture_mode
  }
  property {
    name  = "backDocument"
    type  = "string"
    value = var.authenticidconnector_property_back_document
  }
  property {
    name  = "backImageData"
    type  = "string"
    value = var.authenticidconnector_property_back_image_data
  }
  property {
    name  = "backIsBarcodeDetectedEnable"
    type  = "string"
    value = var.authenticidconnector_property_back_is_barcode_detected_enable
  }
  property {
    name  = "backOverlayColor"
    type  = "string"
    value = var.authenticidconnector_property_back_overlay_color
  }
  property {
    name  = "backOverlayTextAuto"
    type  = "string"
    value = var.authenticidconnector_property_back_overlay_text_auto
  }
  property {
    name  = "backOverlayTextManual"
    type  = "string"
    value = var.authenticidconnector_property_back_overlay_text_manual
  }
  property {
    name  = "backSetManualTimeout"
    type  = "number"
    value = var.authenticidconnector_property_back_set_manual_timeout
  }
  property {
    name  = "barcodeDetectedLabel"
    type  = "string"
    value = var.authenticidconnector_property_barcode_detected_label
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.authenticidconnector_property_base_url
  }
  property {
    name  = "bodyHeaderText"
    type  = "string"
    value = var.authenticidconnector_property_body_header_text
  }
  property {
    name  = "buttonAccept"
    type  = "string"
    value = var.authenticidconnector_property_button_accept
  }
  property {
    name  = "buttonContinue"
    type  = "string"
    value = var.authenticidconnector_property_button_continue
  }
  property {
    name  = "buttonDecline"
    type  = "string"
    value = var.authenticidconnector_property_button_decline
  }
  property {
    name  = "buttonGo"
    type  = "string"
    value = var.authenticidconnector_property_button_go
  }
  property {
    name  = "buttonGotIt"
    type  = "string"
    value = var.authenticidconnector_property_button_got_it
  }
  property {
    name  = "buttonProceed"
    type  = "string"
    value = var.authenticidconnector_property_button_proceed
  }
  property {
    name  = "buttonRetake"
    type  = "string"
    value = var.authenticidconnector_property_button_retake
  }
  property {
    name  = "callPostbackSecurely"
    type  = "string"
    value = var.authenticidconnector_property_call_postback_securely
  }
  property {
    name  = "cancelChooseDocument"
    type  = "string"
    value = var.authenticidconnector_property_cancel_choose_document
  }
  property {
    name  = "captureCriteriaMissedMessage"
    type  = "string"
    value = var.authenticidconnector_property_capture_criteria_missed_message
  }
  property {
    name  = "captureCriteriaMissedTitle"
    type  = "string"
    value = var.authenticidconnector_property_capture_criteria_missed_title
  }
  property {
    name  = "captureFarSelfieContent"
    type  = "string"
    value = var.authenticidconnector_property_capture_far_selfie_content
  }
  property {
    name  = "captureLicenseBackContent"
    type  = "string"
    value = var.authenticidconnector_property_capture_license_back_content
  }
  property {
    name  = "captureLicenseFrontContent"
    type  = "string"
    value = var.authenticidconnector_property_capture_license_front_content
  }
  property {
    name  = "captureNearSelfieContent"
    type  = "string"
    value = var.authenticidconnector_property_capture_near_selfie_content
  }
  property {
    name  = "capturePassportCardBackContent"
    type  = "string"
    value = var.authenticidconnector_property_capture_passport_card_back_content
  }
  property {
    name  = "capturePassportCardFrontContent"
    type  = "string"
    value = var.authenticidconnector_property_capture_passport_card_front_content
  }
  property {
    name  = "capturePassportFrontContent"
    type  = "string"
    value = var.authenticidconnector_property_capture_passport_front_content
  }
  property {
    name  = "channel"
    type  = "string"
    value = var.authenticidconnector_property_channel
  }
  property {
    name  = "channelResponse"
    type  = "string"
    value = var.authenticidconnector_property_channel_response
  }
  property {
    name  = "chooseDocument"
    type  = "string"
    value = var.authenticidconnector_property_choose_document
  }
  property {
    name  = "clientCertificate"
    type  = "string"
    value = var.authenticidconnector_property_client_certificate
  }
  property {
    name  = "clientKey"
    type  = "string"
    value = var.authenticidconnector_property_client_key
  }
  property {
    name  = "consentAgreementStatement"
    type  = "string"
    value = var.authenticidconnector_property_consent_agreement_statement
  }
  property {
    name  = "consentCheckbox"
    type  = "string"
    value = var.authenticidconnector_property_consent_checkbox
  }
  property {
    name  = "consentNotice"
    type  = "string"
    value = var.authenticidconnector_property_consent_notice
  }
  property {
    name  = "consentNoticeHeading"
    type  = "string"
    value = var.authenticidconnector_property_consent_notice_heading
  }
  property {
    name  = "consentNoticeLinkHeading"
    type  = "string"
    value = var.authenticidconnector_property_consent_notice_link_heading
  }
  property {
    name  = "customColor"
    type  = "string"
    value = var.authenticidconnector_property_custom_color
  }
  property {
    name  = "declineConfirmText"
    type  = "string"
    value = var.authenticidconnector_property_decline_confirm_text
  }
  property {
    name  = "declineModalCancelButton"
    type  = "string"
    value = var.authenticidconnector_property_decline_modal_cancel_button
  }
  property {
    name  = "declineModalConfirmButton"
    type  = "string"
    value = var.authenticidconnector_property_decline_modal_confirm_button
  }
  property {
    name  = "declineModalTitle"
    type  = "string"
    value = var.authenticidconnector_property_decline_modal_title
  }
  property {
    name  = "deviceDoesNotHaveACamera"
    type  = "string"
    value = var.authenticidconnector_property_device_does_not_have_acamera
  }
  property {
    name  = "documentLicense"
    type  = "string"
    value = var.authenticidconnector_property_document_license
  }
  property {
    name  = "documentPassport"
    type  = "string"
    value = var.authenticidconnector_property_document_passport
  }
  property {
    name  = "documentPassportCard"
    type  = "string"
    value = var.authenticidconnector_property_document_passport_card
  }
  property {
    name  = "documentUploadFailed"
    type  = "string"
    value = var.authenticidconnector_property_document_upload_failed
  }
  property {
    name  = "doneButton"
    type  = "string"
    value = var.authenticidconnector_property_done_button
  }
  property {
    name  = "donePageTitle"
    type  = "string"
    value = var.authenticidconnector_property_done_page_title
  }
  property {
    name  = "email"
    type  = "string"
    value = var.authenticidconnector_property_email
  }
  property {
    name  = "emailContent"
    type  = "string"
    value = var.authenticidconnector_property_email_content
  }
  property {
    name  = "emailFrom"
    type  = "string"
    value = var.authenticidconnector_property_email_from
  }
  property {
    name  = "emailSubject"
    type  = "string"
    value = var.authenticidconnector_property_email_subject
  }
  property {
    name  = "enableFarSelfie"
    type  = "string"
    value = var.authenticidconnector_property_enable_far_selfie
  }
  property {
    name  = "enableLocationDetection"
    type  = "string"
    value = var.authenticidconnector_property_enable_location_detection
  }
  property {
    name  = "enableSelfieCapture"
    type  = "string"
    value = var.authenticidconnector_property_enable_selfie_capture
  }
  property {
    name  = "errorHeader"
    type  = "string"
    value = var.authenticidconnector_property_error_header
  }
  property {
    name  = "errorPageTitle"
    type  = "string"
    value = var.authenticidconnector_property_error_page_title
  }
  property {
    name  = "faceDetectedLabel"
    type  = "string"
    value = var.authenticidconnector_property_face_detected_label
  }
  property {
    name  = "fixCameraMessage"
    type  = "string"
    value = var.authenticidconnector_property_fix_camera_message
  }
  property {
    name  = "fixCameraTitle"
    type  = "string"
    value = var.authenticidconnector_property_fix_camera_title
  }
  property {
    name  = "focusBack"
    type  = "number"
    value = var.authenticidconnector_property_focus_back
  }
  property {
    name  = "focusFront"
    type  = "number"
    value = var.authenticidconnector_property_focus_front
  }
  property {
    name  = "focusLabel"
    type  = "string"
    value = var.authenticidconnector_property_focus_label
  }
  property {
    name  = "frontCaptureAttempt"
    type  = "number"
    value = var.authenticidconnector_property_front_capture_attempt
  }
  property {
    name  = "frontCaptureMode"
    type  = "string"
    value = var.authenticidconnector_property_front_capture_mode
  }
  property {
    name  = "frontDocument"
    type  = "string"
    value = var.authenticidconnector_property_front_document
  }
  property {
    name  = "frontEnableFaceDetection"
    type  = "string"
    value = var.authenticidconnector_property_front_enable_face_detection
  }
  property {
    name  = "frontImageData"
    type  = "string"
    value = var.authenticidconnector_property_front_image_data
  }
  property {
    name  = "frontOverlayColor"
    type  = "string"
    value = var.authenticidconnector_property_front_overlay_color
  }
  property {
    name  = "frontOverlayTextAuto"
    type  = "string"
    value = var.authenticidconnector_property_front_overlay_text_auto
  }
  property {
    name  = "frontOverlayTextManual"
    type  = "string"
    value = var.authenticidconnector_property_front_overlay_text_manual
  }
  property {
    name  = "frontSetManualTimeout"
    type  = "number"
    value = var.authenticidconnector_property_front_set_manual_timeout
  }
  property {
    name  = "glareBack"
    type  = "number"
    value = var.authenticidconnector_property_glare_back
  }
  property {
    name  = "glareFront"
    type  = "number"
    value = var.authenticidconnector_property_glare_front
  }
  property {
    name  = "glareLabel"
    type  = "string"
    value = var.authenticidconnector_property_glare_label
  }
  property {
    name  = "homeScreen"
    type  = "string"
    value = var.authenticidconnector_property_home_screen
  }
  property {
    name  = "iOSSDKLicenseKey"
    type  = "string"
    value = var.authenticidconnector_property_ios_sdk_license_key
  }
  property {
    name  = "identityVerification"
    type  = "string"
    value = var.authenticidconnector_property_identity_verification
  }
  property {
    name  = "informationalHeader"
    type  = "string"
    value = var.authenticidconnector_property_informational_header
  }
  property {
    name  = "internalError"
    type  = "string"
    value = var.authenticidconnector_property_internal_error
  }
  property {
    name  = "invalidLocationCode"
    type  = "string"
    value = var.authenticidconnector_property_invalid_location_code
  }
  property {
    name  = "livenessSelfie"
    type  = "string"
    value = var.authenticidconnector_property_liveness_selfie
  }
  property {
    name  = "logo"
    type  = "string"
    value = var.authenticidconnector_property_logo
  }
  property {
    name  = "mailHost"
    type  = "string"
    value = var.authenticidconnector_property_mail_host
  }
  property {
    name  = "mailPassword"
    type  = "string"
    value = var.authenticidconnector_property_mail_password
  }
  property {
    name  = "mailPort"
    type  = "string"
    value = var.authenticidconnector_property_mail_port
  }
  property {
    name  = "mailSmtpAuth"
    type  = "string"
    value = var.authenticidconnector_property_mail_smtp_auth
  }
  property {
    name  = "mailSmtpStarttlsEnable"
    type  = "string"
    value = var.authenticidconnector_property_mail_smtp_starttls_enable
  }
  property {
    name  = "mailUsername"
    type  = "string"
    value = var.authenticidconnector_property_mail_username
  }
  property {
    name  = "nextEvent"
    type  = "string"
    value = var.authenticidconnector_property_next_event
  }
  property {
    name  = "noTokenFound"
    type  = "string"
    value = var.authenticidconnector_property_no_token_found
  }
  property {
    name  = "okButton"
    type  = "string"
    value = var.authenticidconnector_property_ok_button
  }
  property {
    name  = "passphrase"
    type  = "string"
    value = var.authenticidconnector_property_passphrase
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.authenticidconnector_property_phone
  }
  property {
    name  = "previewLicenseBackContent"
    type  = "string"
    value = var.authenticidconnector_property_preview_license_back_content
  }
  property {
    name  = "previewLicenseFrontContent"
    type  = "string"
    value = var.authenticidconnector_property_preview_license_front_content
  }
  property {
    name  = "previewPassportCardBackContent"
    type  = "string"
    value = var.authenticidconnector_property_preview_passport_card_back_content
  }
  property {
    name  = "previewPassportCardFrontContent"
    type  = "string"
    value = var.authenticidconnector_property_preview_passport_card_front_content
  }
  property {
    name  = "previewPassportFrontContent"
    type  = "string"
    value = var.authenticidconnector_property_preview_passport_front_content
  }
  property {
    name  = "processingSubtitle"
    type  = "string"
    value = var.authenticidconnector_property_processing_subtitle
  }
  property {
    name  = "processingTitle"
    type  = "string"
    value = var.authenticidconnector_property_processing_title
  }
  property {
    name  = "redirectURL"
    type  = "string"
    value = var.authenticidconnector_property_redirect_url
  }
  property {
    name  = "requestCompleted"
    type  = "string"
    value = var.authenticidconnector_property_request_completed
  }
  property {
    name  = "requestDeclined"
    type  = "string"
    value = var.authenticidconnector_property_request_declined
  }
  property {
    name  = "requestExpiryTimeInMin"
    type  = "number"
    value = var.authenticidconnector_property_request_expiry_time_in_min
  }
  property {
    name  = "requestTerminated"
    type  = "string"
    value = var.authenticidconnector_property_request_terminated
  }
  property {
    name  = "resultPageRescanMessage"
    type  = "string"
    value = var.authenticidconnector_property_result_page_rescan_message
  }
  property {
    name  = "resultPageRescanTitle"
    type  = "string"
    value = var.authenticidconnector_property_result_page_rescan_title
  }
  property {
    name  = "resultPageSuccessMessage"
    type  = "string"
    value = var.authenticidconnector_property_result_page_success_message
  }
  property {
    name  = "resultPageSuccessMessageWithRedirect"
    type  = "string"
    value = var.authenticidconnector_property_result_page_success_message_with_redirect
  }
  property {
    name  = "retryCount"
    type  = "string"
    value = var.authenticidconnector_property_retry_count
  }
  property {
    name  = "reviewScreenBack"
    type  = "string"
    value = var.authenticidconnector_property_review_screen_back
  }
  property {
    name  = "reviewScreenFront"
    type  = "string"
    value = var.authenticidconnector_property_review_screen_front
  }
  property {
    name  = "secretToken"
    type  = "string"
    value = var.authenticidconnector_property_secret_token
  }
  property {
    name  = "segmentResultFailed"
    type  = "string"
    value = var.authenticidconnector_property_segment_result_failed
  }
  property {
    name  = "selfieCaptureAttempt"
    type  = "number"
    value = var.authenticidconnector_property_selfie_capture_attempt
  }
  property {
    name  = "selfieCaptureMode"
    type  = "string"
    value = var.authenticidconnector_property_selfie_capture_mode
  }
  property {
    name  = "selfieEnableFaceDetection"
    type  = "string"
    value = var.authenticidconnector_property_selfie_enable_face_detection
  }
  property {
    name  = "selfieImageData"
    type  = "string"
    value = var.authenticidconnector_property_selfie_image_data
  }
  property {
    name  = "selfieOrientationErrorText"
    type  = "string"
    value = var.authenticidconnector_property_selfie_orientation_error_text
  }
  property {
    name  = "selfieOverlayColor"
    type  = "string"
    value = var.authenticidconnector_property_selfie_overlay_color
  }
  property {
    name  = "selfieOverlayTextAuto"
    type  = "string"
    value = var.authenticidconnector_property_selfie_overlay_text_auto
  }
  property {
    name  = "selfieOverlayTextManual"
    type  = "string"
    value = var.authenticidconnector_property_selfie_overlay_text_manual
  }
  property {
    name  = "selfieSetManualTimeout"
    type  = "number"
    value = var.authenticidconnector_property_selfie_set_manual_timeout
  }
  property {
    name  = "selfieUploadFailed"
    type  = "string"
    value = var.authenticidconnector_property_selfie_upload_failed
  }
  property {
    name  = "selfieUseBackCamera"
    type  = "string"
    value = var.authenticidconnector_property_selfie_use_back_camera
  }
  property {
    name  = "showConsent"
    type  = "string"
    value = var.authenticidconnector_property_show_consent
  }
  property {
    name  = "startButtonText"
    type  = "string"
    value = var.authenticidconnector_property_start_button_text
  }
  property {
    name  = "step1Content"
    type  = "string"
    value = var.authenticidconnector_property_step1_content
  }
  property {
    name  = "step1Title"
    type  = "string"
    value = var.authenticidconnector_property_step1_title
  }
  property {
    name  = "step2Content"
    type  = "string"
    value = var.authenticidconnector_property_step2_content
  }
  property {
    name  = "step2Title"
    type  = "string"
    value = var.authenticidconnector_property_step2_title
  }
  property {
    name  = "step3Content"
    type  = "string"
    value = var.authenticidconnector_property_step3_content
  }
  property {
    name  = "step3Title"
    type  = "string"
    value = var.authenticidconnector_property_step3_title
  }
  property {
    name  = "title"
    type  = "string"
    value = var.authenticidconnector_property_title
  }
  property {
    name  = "transactionAttempts"
    type  = "number"
    value = var.authenticidconnector_property_transaction_attempts
  }
  property {
    name  = "transactionExpiredMessage"
    type  = "string"
    value = var.authenticidconnector_property_transaction_expired_message
  }
  property {
    name  = "transactionExpiryTimeInMin"
    type  = "number"
    value = var.authenticidconnector_property_transaction_expiry_time_in_min
  }
  property {
    name  = "transactionID"
    type  = "string"
    value = var.authenticidconnector_property_transaction_id
  }
  property {
    name  = "twilioAccountSid"
    type  = "string"
    value = var.authenticidconnector_property_twilio_account_sid
  }
  property {
    name  = "twilioAuthToken"
    type  = "string"
    value = var.authenticidconnector_property_twilio_auth_token
  }
  property {
    name  = "twilioFromNumber"
    type  = "string"
    value = var.authenticidconnector_property_twilio_from_number
  }
  property {
    name  = "twilioSmsContent"
    type  = "string"
    value = var.authenticidconnector_property_twilio_sms_content
  }
  property {
    name  = "unknownError"
    type  = "string"
    value = var.authenticidconnector_property_unknown_error
  }
  property {
    name  = "welcomeMessage"
    type  = "string"
    value = var.authenticidconnector_property_welcome_message
  }
}
