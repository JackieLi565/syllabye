#!/bin/bash

# zip_lambda_package.sh build and packages lambda functions

set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <lambda_name>"
  exit 1
fi

LAMBDA_NAME=$1
SRC_DIR="./aws/$LAMBDA_NAME"
LAMBDA="$SRC_DIR/lambda.py"
REQUIREMENTS="$SRC_DIR/requirements.txt"

WORK_DIR="./out"
PKG_DIR="$WORK_DIR/lambda_package"

echo "Building Lambda package for '$LAMBDA_NAME'..."

# Clean directory in case previous attempts failed
rm -rf "$PKG_DIR"
mkdir -p "$PKG_DIR"

# Install optional requirements
if [ -f "$REQUIREMENTS" ]; then
  echo "Installing dependencies from $REQUIREMENTS..."
  python3 -m pip install -r "$REQUIREMENTS" --target "$PKG_DIR"
fi

echo "Copying lambda.py..."
cp "$LAMBDA" "$PKG_DIR"

# Create the zip file
cd "$PKG_DIR"
ZIP_PATH="$OLDPWD/$WORK_DIR/$LAMBDA_NAME.zip"
zip -r "$ZIP_PATH" .

echo "Zip file created for $LAMBDA_NAME"

# Deploy lambda
aws lambda update-function-code \
--function-name $LAMBDA_NAME \
--zip-file "fileb://$ZIP_PATH" \
--endpoint $AWS_LAMBDA_ENDPOINT

echo "Lambda $LAMBDA_NAME successfully deployed"

# Clean up
cd "$OLDPWD"
rm -rf "$PKG_DIR"
echo "Cleaned up working directory."

