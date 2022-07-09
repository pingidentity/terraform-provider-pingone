terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = "~> 0.1"
    }
  }
}

provider "pingone" {
}

resource "pingone_environment" "development" {
  # ...
}