#!/bin/bash
set -e

echo "🚀 Building Lambda function..."

# Build for Linux (Lambda runtime)
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap cmd/lambda/main.go

# Package into zip
echo "📦 Packaging..."
zip lambda.zip bootstrap

# Clean up
rm bootstrap

echo "✅ Lambda package created: lambda.zip"
echo ""
echo "📊 Package size:"
ls -lh lambda.zip
echo ""
echo "Next steps:"
echo "1. Deploy to Lambda using AWS CLI or Console"
echo "2. aws lambda update-function-code --function-name realty-wizard-api --zip-file fileb://lambda.zip"
