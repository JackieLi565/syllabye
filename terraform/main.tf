provider "aws" {
  access_key = var.aws_access_key_id
  secret_key = var.aws_secret_access_key
  region     = var.aws_region

  # Only required required for localstack development
  s3_use_path_style           = var.env == "development"
  skip_credentials_validation = var.env == "development"
  skip_metadata_api_check     = var.env == "development"
  skip_requesting_account_id  = var.env == "development"

  endpoints {
    s3 = var.aws_s3_endpoint
  }
}

resource "aws_s3_bucket" "syllabi_bucket" {
  bucket = var.aws_s3_syllabi_bucket
}
