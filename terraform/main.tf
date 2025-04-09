locals {
  python_lambda_src = "${path.module}/lambda/python/"
  python_lambda_zip = "${path.module}/lambda/out/python.zip"

  is_dev = var.env == "development"
}

provider "aws" {
  access_key = var.aws_access_key_id
  secret_key = var.aws_secret_access_key
  region     = var.aws_region

  # Only required required for localstack development
  s3_use_path_style           = local.is_dev
  skip_credentials_validation = local.is_dev
  skip_metadata_api_check     = local.is_dev
  skip_requesting_account_id  = local.is_dev

  endpoints {
    s3     = var.aws_s3_endpoint
    iam    = var.aws_iam_endpoint
    lambda = var.aws_lambda_endpoint
    sqs    = var.aws_sqs_endpoint
  }
}

data "archive_file" "zip_python_lambda" {
  type        = "zip"
  source_dir  = local.python_lambda_src
  output_path = local.python_lambda_zip
}

module "syllabus_lambda" {
  source   = "./modules/lambda/syllabus"
  filename = local.python_lambda_zip
}

module "webhook_lambda" {
  source   = "./modules/lambda/webhook"
  filename = local.python_lambda_zip
}

module "thumbnail_lambda" {
  source   = "./modules/lambda/thumbnail"
  filename = local.python_lambda_zip
}

module "syllabus_bucket" {
  source                = "./modules/s3/syllabus"
  bucket                = var.aws_s3_syllabi_bucket
  syllabus_lambda_name  = module.syllabus_lambda.function_name
  thumbnail_lambda_name = module.thumbnail_lambda.function_name
  lambda_arns           = [module.syllabus_lambda.function_arn, module.syllabus_lambda.function_arn]
}

module "thumbnail_bucket" {
  source = "./modules/s3/thumbnail"
  bucket = var.aws_s3_thumbnail_bucket
}

module "webhook_queue" {
  source        = "./modules/sqs"
  function_name = module.thumbnail_lambda.function_name
}
