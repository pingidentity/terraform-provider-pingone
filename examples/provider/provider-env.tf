terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = "~> 0.11"
    }
  }
}

provider "pingone" {
}

resource "pingone_environment" "development" {
  # ...
}