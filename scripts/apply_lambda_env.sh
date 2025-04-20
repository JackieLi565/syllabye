#!/bin/bash

# apply_lambda_env.sh applies LAMBDA_ prefixed environment variables to each active Lambda.
# Intended use to be for local development only. 

if [ $ENV != "development" ]; then
  echo "apply_lambda_env.sh can only be applied in development"
  exit 1
fi

LAMBDAS=$(
  aws lambda list-functions \
  --query 'Functions[*].FunctionName' \
  --output text \
  --endpoint $AWS_LAMBDA_ENDPOINT
)

VARIABLES=$(printenv | awk -F= '/^LAMBDA_[A-Za-z0-9_]*=/{sub(/^LAMBDA_/, "", $1); printf "%s=%s,", $1, $2}' | sed 's/,$//')

for LAMBDA in $LAMBDAS; do
  echo "applying variables to $LAMBDA"
  aws lambda update-function-configuration \
    --function-name $LAMBDA \
    --environment "Variables={$VARIABLES}" \
    --endpoint $AWS_LAMBDA_ENDPOINT
done