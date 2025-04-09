output "function_arn" {
  description = "ARN of the thumbnail Lambda"
  value       = aws_lambda_function.this.arn
}

output "function_name" {
  description = "Function name of thumbnail lambda"
  value       = aws_lambda_function.this.function_name
}
