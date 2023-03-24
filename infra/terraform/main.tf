module "enterprise_scale" {
  source  = "Azure/caf-enterprise-scale/azurerm"
  version = "3.3.0"

  providers = {
    azurerm              = azurerm
    azurerm.connectivity = azurerm
    azurerm.management   = azurerm
  }

  root_parent_id   = data.azurerm_client_config.core.tenant_id
  root_id          = var.root_id
  root_name        = var.root_name
  library_path     = "${path.root}/caflib"
  default_location = local.location

  deploy_core_landing_zones = true
  deploy_corp_landing_zones = false
  deploy_demo_landing_zones = false
  deploy_sap_landing_zones  = false

  deploy_management_resources   = false
  deploy_connectivity_resources = false
  deploy_identity_resources     = false
}
