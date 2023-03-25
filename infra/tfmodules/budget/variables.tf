variable "amount" {
  description = "The budget amount"
  type        = number
}

variable "alert_forecast_amount" {
  description = "The forecast amount for which alert should be sent"
  type        = number
}

variable "alert_actual_amount" {
  description = "The actual amount for which alert should be sent"
  type        = number
}

variable "contact_emails" {
  description = "The emails the alert should be sent to"
  type        = list(string)
}
