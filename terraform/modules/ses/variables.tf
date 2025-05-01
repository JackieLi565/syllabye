variable "domain" {}

variable "is_dev" {
  type        = bool
  description = "Indicator to create domain identity in local env"
}

# TODO? Move each template to its own module
variable "welcome_template_name" {}

variable "upload_success_template_name" {}

variable "upload_error_template_name" {}
