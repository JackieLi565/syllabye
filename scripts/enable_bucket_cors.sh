#!/bin/bash

# enable_bucket_cors.sh enables CORS to all local buckets.
# Intended use to be for local development only. 

if [ $ENV != "development" ]; then
  echo "apply_lambda_env.sh can only be applied in development"
  exit 1
fi

BUCKETS=$(
  aws s3api list-buckets \
    --query "Buckets[].Name" \
    --output text \
    --endpoint-url $AWS_S3_ENDPOINT
)

for BUCKET in $BUCKETS; do
  echo "enabling CORS for bucket $BUCKET"
  aws s3api put-bucket-cors \
    --bucket $BUCKET \
    --cors-configuration '{
      "CORSRules": [
        {
          "AllowedHeaders": ["*"],
          "AllowedMethods": ["GET", "PUT", "POST", "DELETE"],
          "AllowedOrigins": ["*"],
          "ExposeHeaders": [],
          "MaxAgeSeconds": 3000
        }
      ]
    }' \
    --endpoint $AWS_S3_ENDPOINT
done