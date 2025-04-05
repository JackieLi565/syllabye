locals {
  lambda_handler    = "lambda.handler"
  python_lambda_src = "${path.module}/lambda/python/"
  python_lambda_zip = "${path.module}/lambda/out/python.zip"
  python_rt         = "python3.11"

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

resource "aws_s3_bucket" "syllabi_bucket" {
  bucket = var.aws_s3_syllabi_bucket
}

resource "aws_iam_role" "lambda_s3_execute" {
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

resource "aws_iam_role_policy_attachment" "lambda_s3_execute" {
  role       = aws_iam_role.lambda_s3_execute.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_function" "syllabus_trigger" {
  filename      = local.python_lambda_zip
  function_name = "syllabus_trigger"
  handler       = local.lambda_handler
  runtime       = local.python_rt
  role          = aws_iam_role.lambda_s3_execute.arn
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

resource "aws_sqs_queue" "webhook" {
  name = "WebhookQueue"
}

resource "aws_iam_role" "lambda_sqs_execute" {
  name = "AWSLambdaSQSExecute"
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

resource "aws_iam_role_policy_attachment" "lambda_sqs_execute" {
  role       = aws_iam_role.lambda_sqs_execute.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
}

resource "aws_lambda_function" "webhook" {
  filename      = local.python_lambda_zip
  function_name = "webhook"
  handler       = local.lambda_handler
  runtime       = local.python_rt
  role          = aws_iam_role.lambda_sqs_execute.arn
}

resource "aws_lambda_permission" "webhook_allow_sqs" {
  statement_id  = "AllowExecutionFromSQS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.webhook.function_name
  principal     = "sqs.amazonaws.com"
  source_arn    = aws_sqs_queue.webhook.arn
}

resource "aws_lambda_event_source_mapping" "sqs_webhook_event" {
  event_source_arn = aws_sqs_queue.webhook.arn
  function_name    = aws_lambda_function.webhook.arn
  batch_size       = 5
  # 5 second delay in dev
  maximum_batching_window_in_seconds = local.is_dev ? 5 : 60 * 5
}
