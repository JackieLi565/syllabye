variable "bucket" {
  type        = string
  description = "Syllabus bucket name"
}

variable "lambdas" {
  type = list(object({
    name = string
    arn  = string
  }))
  description = "List of Lambda names and arns"
}

