resource "aws_iam_role" "this" {
  name = "WebhookLambdaRole"
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

resource "aws_iam_role_policy_attachment" "this" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
}

resource "aws_lambda_function" "this" {
  filename      = var.filename
  function_name = "webhook"
  handler       = "lambda.handler"
  runtime       = "python3.11"
  role          = aws_iam_role.this.arn
}
