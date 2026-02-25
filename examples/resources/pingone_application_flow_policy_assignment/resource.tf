resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_application" {
  # ...
}

resource "davinci_flow" "authentication" {
  # ...
}

resource "davinci_application" "awesome_dv_app" {
  environment_id = pingone_environment.my_environment.id

  name = "Awesome DaVinci Application"
}

resource "davinci_application_flow_policy" "authentication_flow_policy" {
  environment_id = pingone_environment.my_environment.id
  application_id = davinci_application.awesome_dv_app.id

  name   = "Authentication"
  status = "enabled"

  policy_flow {
    flow_id    = davinci_flow.authentication.id
    version_id = -1
    weight     = 100
  }
}

resource "pingone_application_flow_policy_assignment" "authentication_davinci_flow_policy_assignment" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_application.id

  flow_policy_id = davinci_application_flow_policy.authentication_flow_policy.id

  priority = 1
}
