variable "aws_access_key_id" {}

variable "aws_secret_access_key" {}

variable "aws_region" {}

variable "aws_s3_endpoint" {}

variable "aws_s3_syllabi_bucket" {}

variable "aws_iam_endpoint" {}

variable "aws_lambda_endpoint" {}

variable "aws_sqs_endpoint" {}

variable "aws_s3_thumbnail_bucket" {
  type        = string
  description = "Name of thumbnail bucket"
}

variable "env" {
  default = "development"
}
