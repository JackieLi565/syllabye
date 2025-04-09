resource "aws_s3_bucket" "this" {
  bucket = var.bucket
}

resource "aws_lambda_permission" "syllabus" {
  action        = "lambda:InvokeFunction"
  function_name = var.syllabus_lambda_name
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.this.arn
}

resource "aws_lambda_permission" "thumbnail" {
  action        = "lambda:InvokeFunction"
  function_name = var.thumbnail_lambda_name
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.this.arn
}

resource "aws_s3_bucket_notification" "this" {
  bucket = aws_s3_bucket.this.id

  dynamic "lambda_function" {
    for_each = var.lambda_arns

    content {
      lambda_function_arn = lambda_function.value
      events              = ["s3:ObjectCreated:*"]
    }
  }

  depends_on = [aws_lambda_permission.thumbnail, aws_lambda_permission.syllabus]
}
