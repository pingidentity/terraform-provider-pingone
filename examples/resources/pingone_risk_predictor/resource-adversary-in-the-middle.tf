resource "pingone_risk_predictor" "my_awesome_adversary_in_the_middle_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Adversary In The Middle Predictor"
  compact_name   = "myAwesomeAitMPredictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_adversary_in_the_middle = {
    allowed_domain_list = ["domain1.com", "domain2.com", "domain3.com"]
  }
}