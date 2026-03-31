resource "pingone_davinci_connector_instance" "microsoftTeamsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "microsoftTeamsConnector"
  }
  name = "My awesome microsoftTeamsConnector"
  property {
    name  = "authType"
    type  = "string"
    value = var.microsoftteamsconnector_property_auth_type
  }
  property {
    name  = "bodyData"
    type  = "string"
    value = var.microsoftteamsconnector_property_body_data
  }
  property {
    name  = "button"
    type  = "string"
    value = var.microsoftteamsconnector_property_button
  }
  property {
    name  = "channelId"
    type  = "string"
    value = var.microsoftteamsconnector_property_channel_id
  }
  property {
    name  = "channels"
    type  = "string"
    value = var.microsoftteamsconnector_property_channels
  }
  property {
    name  = "chatId"
    type  = "string"
    value = var.microsoftteamsconnector_property_chat_id
  }
  property {
    name  = "customAuth"
    type  = "string"
    value = jsonencode({})
  }
  property {
    name  = "displayName"
    type  = "string"
    value = var.microsoftteamsconnector_property_display_name
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.microsoftteamsconnector_property_endpoint
  }
  property {
    name  = "headersForm"
    type  = "string"
    value = var.microsoftteamsconnector_property_headers_form
  }
  property {
    name  = "memberEmail"
    type  = "string"
    value = var.microsoftteamsconnector_property_member_email
  }
  property {
    name  = "memberId"
    type  = "string"
    value = var.microsoftteamsconnector_property_member_id
  }
  property {
    name  = "members"
    type  = "string"
    value = var.microsoftteamsconnector_property_members
  }
  property {
    name  = "membershipId"
    type  = "string"
    value = var.microsoftteamsconnector_property_membership_id
  }
  property {
    name  = "messageBodyData"
    type  = "string"
    value = var.microsoftteamsconnector_property_message_body_data
  }
  property {
    name  = "method"
    type  = "string"
    value = var.microsoftteamsconnector_property_method
  }
  property {
    name  = "paramsForm"
    type  = "string"
    value = var.microsoftteamsconnector_property_params_form
  }
  property {
    name  = "roles"
    type  = "string"
    value = var.microsoftteamsconnector_property_roles
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.microsoftteamsconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.microsoftteamsconnector_property_skip_button_press
  }
  property {
    name  = "teamId"
    type  = "string"
    value = var.microsoftteamsconnector_property_team_id
  }
  property {
    name  = "teams"
    type  = "string"
    value = var.microsoftteamsconnector_property_teams
  }
  property {
    name  = "userAccessToken"
    type  = "string"
    value = var.microsoftteamsconnector_property_user_access_token
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.microsoftteamsconnector_property_user_id
  }
}
