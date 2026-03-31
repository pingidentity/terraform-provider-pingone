resource "pingone_davinci_connector_instance" "challengeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "challengeConnector"
  }
  name = "My awesome challengeConnector"
  property {
    name  = "challenge"
    type  = "string"
    value = var.challengeconnector_property_challenge
  }
  property {
    name  = "challengeStatus"
    type  = "string"
    value = var.challengeconnector_property_challenge_status
  }
  property {
    name  = "challengeTimeout"
    type  = "number"
    value = var.challengeconnector_property_challenge_timeout
  }
  property {
    name  = "claimsNameValuePairs"
    type  = "string"
    value = var.challengeconnector_property_claims_name_value_pairs
  }
  property {
    name  = "isChallengeComplete"
    type  = "string"
    value = var.challengeconnector_property_is_challenge_complete
  }
  property {
    name  = "pollInterval"
    type  = "string"
    value = var.challengeconnector_property_poll_interval
  }
  property {
    name  = "pollRetries"
    type  = "string"
    value = var.challengeconnector_property_poll_retries
  }
  property {
    name  = "updatedByFlowId"
    type  = "string"
    value = var.challengeconnector_property_updated_by_flow_id
  }
}
