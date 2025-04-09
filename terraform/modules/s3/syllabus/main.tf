resource "aws_s3_bucket" "this" {
  bucket = var.bucket
}

resource "aws_lambda_permission" "this" {
  for_each = { for lambda in var.lambdas : lambda.name => lambda }

  action        = "lambda:InvokeFunction"
  function_name = each.value.name
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.this.arn
}

resource "aws_s3_bucket_notification" "this" {
  bucket = aws_s3_bucket.this.id

  dynamic "lambda_function" {
    for_each = var.lambdas

    content {
      lambda_function_arn = lambda_function.value.arn
      events              = ["s3:ObjectCreated:*"]
    }
  }
}
