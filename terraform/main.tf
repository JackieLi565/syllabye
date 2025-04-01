locals {
  handler_name = "lambda_function.lambda_handler"

  python_lambda_src = "${path.module}/lambda/python/"
  python_lambda_zip = "${path.module}/lambda/out/python.zip"
}

provider "aws" {
  access_key = var.aws_access_key_id
  secret_key = var.aws_secret_access_key
  region     = var.aws_region

  # Only required required for localstack development
  s3_use_path_style           = var.env == "development"
  skip_credentials_validation = var.env == "development"
  skip_metadata_api_check     = var.env == "development"
  skip_requesting_account_id  = var.env == "development"

  endpoints {
    s3     = var.aws_s3_endpoint
    iam    = var.aws_iam_endpoint
    lambda = var.aws_lambda_endpoint
  }
}

resource "aws_s3_bucket" "syllabi_bucket" {
  bucket = var.aws_s3_syllabi_bucket
}

resource "aws_iam_role" "lambda_execute" {
  name = "AWSLambdaExecute"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_execute_basic" {
  role       = aws_iam_role.lambda_execute.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "archive_file" "zip_python_lambda" {
  type        = "zip"
  source_dir  = local.python_lambda_src
  output_path = local.python_lambda_zip
}

resource "aws_lambda_function" "syllabus_trigger" {
  filename      = local.python_lambda_zip
  function_name = "s3_syllabus_trigger" # Name functions based on programming language naming convention e.g. Python snake_case
  handler       = local.handler_name
  runtime       = "python3.9"
  role          = aws_iam_role.lambda_execute.arn
}

resource "aws_lambda_permission" "syllabus_trigger_allow_s3" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.syllabus_trigger.function_name
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.syllabi_bucket.arn
}

resource "aws_s3_bucket_notification" "syllabus_trigger" {
  bucket = aws_s3_bucket.syllabi_bucket.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.syllabus_trigger.arn
    events              = ["s3:ObjectCreated:*"]
  }

  depends_on = [aws_lambda_permission.syllabus_trigger_allow_s3]
}
