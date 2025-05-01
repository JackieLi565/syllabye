variable "aws_access_key_id" {}

variable "aws_secret_access_key" {}

variable "aws_region" {}

variable "aws_s3_endpoint" {}

variable "aws_s3_syllabi_bucket" {
  type        = string
  description = "Name of syllabi bucket"
}

variable "aws_iam_endpoint" {}

variable "aws_lambda_endpoint" {}

variable "aws_sqs_endpoint" {}

variable "aws_ses_endpoint" {}

variable "welcome_template_name" {
  type        = string
  description = "Name for welcome template"
}

variable "upload_success_template_name" {
  type        = string
  description = "Name for upload success template"
}

variable "upload_error_template_name" {
  type        = string
  description = "Name for upload error template"
}

variable "aws_s3_thumbnail_bucket" {
  type        = string
  description = "Name of thumbnail bucket"
}

variable "domain" {
  type        = string
  description = "Domain name"
}

variable "env" {
  default = "development"
}
