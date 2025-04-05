variable "aws_access_key_id" {}

variable "aws_secret_access_key" {}

variable "aws_region" {}

variable "aws_s3_endpoint" {}

variable "aws_s3_syllabi_bucket" {}

variable "aws_iam_endpoint" {}

variable "aws_lambda_endpoint" {}

variable "aws_sqs_endpoint" {}

variable "env" {
  default = "development"
}
