resource "pingone_risk_predictor" "my_awesome_ip_reputation_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome IP Reputation Predictor"
  compact_name   = "my_awesome_ip_reputation_predictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_ip_reputation = {
    allowed_cidr_list = ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/24"]
  }
}