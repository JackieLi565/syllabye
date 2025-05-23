output "function_arn" {
  description = "ARN of syllabus Lambda"
  value       = aws_lambda_function.this.arn
}

output "function_name" {
  description = "Function name of syllabus lambda"
  value       = aws_lambda_function.this.function_name
}
