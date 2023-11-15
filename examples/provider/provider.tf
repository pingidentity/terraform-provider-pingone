terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = "~> 0.24"
    }
  }
}

provider "pingone" {
  client_id      = var.client_id
  client_secret  = var.client_secret
  environment_id = var.environment_id
  region         = var.region

  force_delete_production_type = false
}

resource "pingone_environment" "my_environment" {
  # ...
}