terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = ">= 1.16, < 1.17"
    }
  }
}

provider "pingone" {
}

resource "pingone_environment" "my_environment" {
  # ...
}
