module "subscription_budget_caf" {
  source = "../tfmodules/budget"
  providers = {
    azurerm = azurerm.caf
  }
  amount                = 10.0
  alert_forecast_amount = 8.0
  alert_actual_amount   = 5.0
  contact_emails        = [var.caf_security_email_contact]
}

module "subscription_budget_management" {
  source = "../tfmodules/budget"
  providers = {
    azurerm = azurerm.management
  }
  amount                = 10.0
  alert_forecast_amount = 8.0
  alert_actual_amount   = 5.0
  contact_emails        = [var.caf_security_email_contact]
}
