output "function_arn" {
  description = "ARN of the webhook Lambda"
  value       = aws_lambda_function.this.arn
}

output "function_name" {
  description = "Function name of webhook lambda"
  value       = aws_lambda_function.this.function_name
}
