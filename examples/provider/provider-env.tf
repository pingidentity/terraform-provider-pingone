terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = "~> 0.15"
    }
  }
}

provider "pingone" {
}

resource "pingone_environment" "my_environment" {
  # ...
}