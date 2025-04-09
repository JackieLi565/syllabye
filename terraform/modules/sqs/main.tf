resource "aws_sqs_queue" "this" {
  name = "WebhookQueue"
}

resource "aws_lambda_permission" "this" {
  action        = "lambda:InvokeFunction"
  function_name = var.function_name
  principal     = "sqs.amazonaws.com"
  source_arn    = aws_sqs_queue.this.arn
}

resource "aws_lambda_event_source_mapping" "this" {
  event_source_arn                   = aws_sqs_queue.this.arn
  function_name                      = var.function_name
  batch_size                         = 5
  maximum_batching_window_in_seconds = 60 * 5
}
