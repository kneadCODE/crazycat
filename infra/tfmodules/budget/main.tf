resource "azurerm_consumption_budget_subscription" "budget" {
  name            = "monthly"
  subscription_id = data.azurerm_subscription.sub.id
  amount          = var.amount
  time_grain      = "Monthly"
  time_period {
    start_date = "2023-03-01T00:00:00Z"
    end_date = "2030-03-01T00:00:00Z"
  }
  notification {
    enabled        = true
    threshold      = var.alert_forecast_amount
    operator       = "GreaterThanOrEqualTo"
    threshold_type = "Forecasted"
    contact_emails = var.contact_emails
  }
  notification {
    enabled        = true
    threshold      = var.alert_actual_amount
    operator       = "GreaterThanOrEqualTo"
    threshold_type = "Actual"
    contact_emails = var.contact_emails
  }
}
