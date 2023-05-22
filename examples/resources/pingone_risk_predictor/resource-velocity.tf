resource "pingone_risk_predictor" "my_awesome_velocity_predictor_by_ip" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Velocity Predictor By IP"
  compact_name   = "my_awesome_velocity_predictor_by_ip"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_velocity = {
    of = "$${event.user.id}"
  }
}

resource "pingone_risk_predictor" "my_awesome_velocity_predictor_by_user" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Velocity Predictor By User"
  compact_name   = "my_awesome_velocity_predictor_by_user"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_velocity = {
    of = "$${event.ip}"
  }
}