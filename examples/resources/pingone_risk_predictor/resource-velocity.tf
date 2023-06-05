resource "pingone_risk_predictor" "my_awesome_velocity_predictor_by_ip" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Velocity Predictor By IP"
  compact_name   = "myAwesomeVelocityPredictorByIp"

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
  compact_name   = "myAwesomeVelocityPredictorByUser"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_velocity = {
    of = "$${event.ip}"
  }
}