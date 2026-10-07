terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = ">= 1.23, < 1.24"
    }
  }
}

provider "pingone" {
}

resource "pingone_environment" "my_environment" {
  # ...
}
