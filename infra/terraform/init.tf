terraform {
  required_version = "~> 1.4"
  cloud {
    organization = "crazycat"
    workspaces {
      name = "crazycat"
    }
  }
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.54"
    }
  }
}

provider "azurerm" {
  alias = "caf"
  features {}
}
provider "azurerm" {
  alias           = "management"
  subscription_id = var.caf_subscription_id_management
  features {}
}
#provider "azurerm" {
#  alias           = "connectivity"
#  subscription_id = var.caf_subscription_id_connectivity
#  features {}
#}
