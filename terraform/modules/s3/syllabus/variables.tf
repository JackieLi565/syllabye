variable "bucket" {
  type        = string
  description = "Syllabus bucket name"
}

variable "syllabus_lambda_name" {
  type        = string
  description = "Name of syllabus lambda"
}

variable "thumbnail_lambda_name" {
  type        = string
  description = "Name of thumbnail lambda"
}

variable "lambda_arns" {
  type        = list(string)
  description = "List of lambda ARNs receiving the trigger"
}

