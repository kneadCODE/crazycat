module "enterprise_scale" {
  source  = "Azure/caf-enterprise-scale/azurerm"
  version = "3.3.0"

  providers = {
    azurerm              = azurerm
    azurerm.connectivity = azurerm
    azurerm.management   = azurerm
  }

  root_parent_id   = var.caf_mg_root_parent_id
  root_id          = "crazycat"
  root_name        = "crazycat"
  library_path     = "${path.root}/caflib"
  default_location = local.location

  deploy_core_landing_zones = true
  deploy_corp_landing_zones = false
  deploy_demo_landing_zones = false
  deploy_sap_landing_zones  = false

  deploy_management_resources = false
  subscription_id_management  = var.caf_subscription_id_management
#  configure_management_resources = {
#    settings = {
#      log_analytics = {
#        enabled = true
#        config = {
#          retention_in_days                                 = 30    # this is the lowest number allowed for this resource
#          enable_monitoring_for_vm                          = false # no money, too poor
#          enable_monitoring_for_vmss                        = false # no money, too poor
#          enable_solution_for_agent_health_assessment       = false # no money, too poor
#          enable_solution_for_anti_malware                  = false # no money, too poor
#          enable_solution_for_change_tracking               = false # no money, too poor
#          enable_solution_for_service_map                   = false # no money, too poor
#          enable_solution_for_sql_assessment                = false # no money, too poor
#          enable_solution_for_sql_vulnerability_assessment  = false # no money, too poor
#          enable_solution_for_sql_advanced_threat_detection = false # no money, too poor
#          enable_solution_for_updates                       = false # no money, too poor
#          enable_solution_for_vm_insights                   = false # no money, too poor
#          enable_sentinel                                   = false # no money, too poor
#        }
#      }
#      security_center = {
#        enabled = true
#        config = {
#          email_security_contact             = var.caf_security_email_contact
#          enable_defender_for_app_services   = false # too expensive for a demo
#          enable_defender_for_arm            = false # too expensive for a demo
#          enable_defender_for_containers     = false # too expensive for a demo
#          enable_defender_for_dns            = false # too expensive for a demo
#          enable_defender_for_key_vault      = false # too expensive for a demo
#          enable_defender_for_oss_databases  = false # too expensive for a demo
#          enable_defender_for_servers        = false # too expensive for a demo
#          enable_defender_for_sql_servers    = false # too expensive for a demo
#          enable_defender_for_sql_server_vms = false # too expensive for a demo
#          enable_defender_for_storage        = false # too expensive for a demo
#        }
#      }
#    }
#    advanced = {
#      custom_settings_by_resource_type = {
#        azurerm_resource_group = {
#          management = {
#            name = "rg-mgmt-001" # Need to override it else it doesn't follow a proper convention
#          }
#        }
#        azurerm_log_analytics_workspace = {
#          management = {
#            name = "la-mgmt-001" # Need to override it else it doesn't follow a proper convention
#          }
#        }
#        azurerm_automation_account = {
#          management = {
#            name     = "aa-mgmt-001" # Need to override it else it doesn't follow a proper convention
#            sku_name = "Free"
#          }
#        }
#      }
#    }
#  }

  deploy_connectivity_resources = false
  deploy_identity_resources     = false
}
