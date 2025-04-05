output "sqs_webhook_queue_url" {
  description = "URL of the WebhookQueue"
  value       = aws_sqs_queue.webhook.id
}
