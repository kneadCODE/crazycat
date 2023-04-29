locals {
  caf_root_id = "crazycat"

  caf_corp_mg_default_archetype = {
    archetype_id = "es_corp"
    parameters = {
      Deny-Resource-Locations = {
        listOfAllowedLocations = [
          local.location,
        ]
      }
      Deny-RSG-Locations = {
        listOfAllowedLocations = [
          local.location,
        ]
      }
    }
    access_control = {}
  }
  caf_online_mg_default_archetype = {
    archetype_id = "es_corp"
    parameters = {
      Deny-Resource-Locations = {
        listOfAllowedLocations = [
          local.location,
        ]
      }
      Deny-RSG-Locations = {
        listOfAllowedLocations = [
          local.location,
        ]
      }
    }
    access_control = {}
  }
}

module "enterprise_scale" {
  source  = "Azure/caf-enterprise-scale/azurerm"
  version = "3.3.0"

  providers = {
    azurerm              = azurerm.caf
    azurerm.management   = azurerm.management
    azurerm.connectivity = azurerm.caf # TODO: Change this to azurerm.connectivity.
  }

  root_parent_id   = var.caf_mg_root_parent_id
  root_id          = local.caf_root_id
  root_name        = local.caf_root_id
  library_path     = "${path.root}/caflib"
  default_location = local.location

  deploy_core_landing_zones   = true
  deploy_corp_landing_zones   = true
  deploy_demo_landing_zones   = false
  deploy_online_landing_zones = true
  deploy_sap_landing_zones    = false

  deploy_management_resources = true
  subscription_id_management  = var.caf_subscription_id_management
  configure_management_resources = {
    settings = {
      log_analytics = {
        enabled = true
        config = {
          retention_in_days                                 = 30    # this is the lowest number allowed for this resource
          enable_monitoring_for_vm                          = false # no money, too poor
          enable_monitoring_for_vmss                        = false # no money, too poor
          enable_solution_for_agent_health_assessment       = false # no money, too poor
          enable_solution_for_anti_malware                  = false # no money, too poor
          enable_solution_for_change_tracking               = false # no money, too poor
          enable_solution_for_service_map                   = false # no money, too poor
          enable_solution_for_sql_assessment                = false # no money, too poor
          enable_solution_for_sql_vulnerability_assessment  = false # no money, too poor
          enable_solution_for_sql_advanced_threat_detection = false # no money, too poor
          enable_solution_for_updates                       = false # no money, too poor
          enable_solution_for_vm_insights                   = false # no money, too poor
          enable_sentinel                                   = false # no money, too poor
        }
      }
      security_center = {
        enabled = true
        config = {
          email_security_contact             = var.caf_security_email_contact
          enable_defender_for_app_services   = false # too expensive for a demo
          enable_defender_for_arm            = false # too expensive for a demo
          enable_defender_for_containers     = false # too expensive for a demo
          enable_defender_for_dns            = false # too expensive for a demo
          enable_defender_for_key_vault      = false # too expensive for a demo
          enable_defender_for_oss_databases  = false # too expensive for a demo
          enable_defender_for_servers        = false # too expensive for a demo
          enable_defender_for_sql_servers    = false # too expensive for a demo
          enable_defender_for_sql_server_vms = false # too expensive for a demo
          enable_defender_for_storage        = false # too expensive for a demo
        }
      }
    }
    advanced = {
      custom_settings_by_resource_type = {
        azurerm_resource_group = {
          management = {
            name = "rg-mgmt-001" # Need to override it else it doesn't follow a proper convention
          }
        }
        azurerm_log_analytics_workspace = {
          management = {
            name = "la-mgmt-001" # Need to override it else it doesn't follow a proper convention
          }
        }
        azurerm_automation_account = {
          management = {
            name     = "aa-mgmt-001" # Need to override it else it doesn't follow a proper convention
            sku_name = "Free"
          }
        }
      }
    }
  }

  deploy_identity_resources = true
  subscription_id_identity  = var.caf_subscription_id_identity
  configure_identity_resources = {
    settings = {
      identity = {
        enabled = true # Enabling this more for demo. No resources are deployed, so should not cost money.
        config = {
          enable_deny_public_ip             = true
          enable_deny_rdp_from_internet     = true
          enable_deny_subnet_without_nsg    = true
          enable_deploy_azure_backup_on_vms = true # Since we are not really deploying anything in identity, it is ok to turn this on as it won't incur any cost.
        }
      }
    }
  }

  deploy_connectivity_resources = false

  custom_landing_zones = {
    "${local.caf_root_id}-tool" = {
      display_name               = "Tooling"
      parent_management_group_id = "${local.caf_root_id}-corp"
      archetype_config           = local.caf_corp_mg_default_archetype
      subscription_ids           = []
    }
    "${local.caf_root_id}-data" = {
      display_name               = "Data"
      parent_management_group_id = "${local.caf_root_id}-corp"
      archetype_config           = local.caf_corp_mg_default_archetype
      subscription_ids           = []
    }
    "${local.caf_root_id}-compute" = {
      display_name               = "Compute"
      parent_management_group_id = "${local.caf_root_id}-corp"
      archetype_config           = local.caf_corp_mg_default_archetype
      subscription_ids           = []
    }
    "${local.caf_root_id}-web" = {
      display_name               = "Web"
      parent_management_group_id = "${local.caf_root_id}-online"
      archetype_config           = local.caf_online_mg_default_archetype
      subscription_ids           = []
    }
  }
}
