terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = "~> 0.14"
    }
  }
}

provider "pingone" {
}

resource "pingone_environment" "development" {
  # ...
}