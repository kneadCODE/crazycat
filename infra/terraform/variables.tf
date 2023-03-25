variable "caf_mg_root_parent_id" {
  description = "The root management group ID for CAF"
  type        = string
}
variable "caf_subscription_id_management" { # Subscription is not created by the CAF module. So create outside then pass it in for CAF to manage.
  description = "The ID of the management subscription"
  type        = string
}
variable "caf_subscription_id_connectivity" { # Subscription is not created by the CAF module. So create outside then pass it in for CAF to manage.
  description = "The ID of the connectivity subscription"
  type        = string
}
variable "caf_security_email_contact" {
  description = "The email contact for security center"
  type        = string
}
