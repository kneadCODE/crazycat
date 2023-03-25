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
      version = "~> 3.49"
    }
  }
}

provider "azurerm" {
  features {}
}
