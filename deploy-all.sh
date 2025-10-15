#!/bin/bash
set -e

echo "🚀 Realty Wizard - Complete Deployment Script"
echo "=============================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check prerequisites
echo "📋 Checking prerequisites..."

if ! command -v aws &> /dev/null; then
    echo -e "${RED}❌ AWS CLI not found. Install it first:${NC}"
    echo "   brew install awscli"
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go not found. Install it first:${NC}"
    echo "   brew install go"
    exit 1
fi

if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ Node.js not found. Install it first:${NC}"
    echo "   brew install node"
    exit 1
fi

echo -e "${GREEN}✅ All prerequisites met${NC}"
echo ""

# Get deployment parameters
echo "🔧 Configuration"
echo "================"
echo ""

read -p "AWS Region (default: us-east-1): " AWS_REGION
AWS_REGION=${AWS_REGION:-us-east-1}

read -p "Lambda function name (default: realty-wizard-api): " LAMBDA_NAME
LAMBDA_NAME=${LAMBDA_NAME:-realty-wizard-api}

read -p "S3 bucket name (must be globally unique): " BUCKET_NAME
if [ -z "$BUCKET_NAME" ]; then
    BUCKET_NAME="realty-wizard-$(date +%s)"
    echo "   Using generated bucket name: $BUCKET_NAME"
fi

echo ""
echo -e "${YELLOW}⚠️  This will deploy to AWS and may incur charges${NC}"
read -p "Continue? (y/n): " CONFIRM
if [ "$CONFIRM" != "y" ]; then
    echo "Deployment cancelled"
    exit 0
fi

echo ""
echo "=========================================="
echo "Phase 1: Deploying Backend (Lambda)"
echo "=========================================="
echo ""

cd backend

echo "📦 Building Lambda function..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap cmd/lambda/main.go
zip lambda.zip bootstrap
rm bootstrap

echo -e "${GREEN}✅ Lambda package created${NC}"
echo ""

# Check if Lambda function exists
if aws lambda get-function --function-name "$LAMBDA_NAME" --region "$AWS_REGION" 2>/dev/null; then
    echo "📤 Updating existing Lambda function..."
    aws lambda update-function-code \
        --function-name "$LAMBDA_NAME" \
        --zip-file fileb://lambda.zip \
        --region "$AWS_REGION" > /dev/null
    echo -e "${GREEN}✅ Lambda function updated${NC}"
else
    echo "🆕 Lambda function doesn't exist yet"
    echo "   Please create it manually in AWS Console first:"
    echo "   1. Go to https://console.aws.amazon.com/lambda"
    echo "   2. Create function named: $LAMBDA_NAME"
    echo "   3. Runtime: 'Provide your own bootstrap on Amazon Linux 2023'"
    echo "   4. Upload lambda.zip"
    echo "   5. Set Memory: 512 MB, Timeout: 30 seconds"
    echo "   6. Set environment variable: DB_PATH=/tmp/realty-wizard.db"
    echo ""
    echo "   Then create API Gateway HTTP API pointing to this Lambda"
    echo ""
    exit 1
fi

# Get API Gateway URL
echo ""
echo "🔍 Finding API Gateway URL..."
API_ID=$(aws apigatewayv2 get-apis --region "$AWS_REGION" --query "Items[?Name=='realty-wizard-api'].ApiId" --output text)

if [ -n "$API_ID" ]; then
    API_URL="https://${API_ID}.execute-api.${AWS_REGION}.amazonaws.com"
    echo -e "${GREEN}✅ API URL: $API_URL${NC}"
else
    echo -e "${YELLOW}⚠️  API Gateway not found. Please create it manually and note the URL${NC}"
    read -p "Enter your API Gateway URL: " API_URL
fi

cd ..

echo ""
echo "=========================================="
echo "Phase 2: Deploying Frontend (S3 + CloudFront)"
echo "=========================================="
echo ""

cd frontend

echo "🔧 Configuring API URL..."
cat > .env.production << EOF
VITE_API_URL=$API_URL
EOF

echo "📦 Building React app..."
npm run build

echo ""
echo "☁️  Creating S3 bucket..."
if aws s3 ls "s3://$BUCKET_NAME" 2>/dev/null; then
    echo "   Bucket already exists"
else
    aws s3 mb "s3://$BUCKET_NAME" --region "$AWS_REGION"

    # Enable static website hosting
    aws s3 website "s3://$BUCKET_NAME" \
        --index-document index.html \
        --error-document index.html

    # Set bucket policy for public read
    cat > /tmp/bucket-policy.json << EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::${BUCKET_NAME}/*"
    }
  ]
}
EOF

    aws s3api put-bucket-policy \
        --bucket "$BUCKET_NAME" \
        --policy file:///tmp/bucket-policy.json

    rm /tmp/bucket-policy.json
fi

echo "📤 Uploading frontend files..."
aws s3 sync dist/ "s3://$BUCKET_NAME/" --delete

# Set cache headers for assets
if [ -d "dist/assets" ]; then
    aws s3 cp "s3://$BUCKET_NAME/assets" "s3://$BUCKET_NAME/assets" \
        --recursive \
        --metadata-directive REPLACE \
        --cache-control "public, max-age=31536000, immutable" \
        --region "$AWS_REGION"
fi

echo -e "${GREEN}✅ Frontend uploaded to S3${NC}"

cd ..

echo ""
echo "=========================================="
echo "🎉 Deployment Complete!"
echo "=========================================="
echo ""
echo "📝 Next Steps:"
echo ""
echo "1. Backend API: $API_URL"
echo "   Test: curl $API_URL/api/health"
echo ""
echo "2. Frontend S3: http://$BUCKET_NAME.s3-website-$AWS_REGION.amazonaws.com"
echo ""
echo "3. Set up CloudFront (recommended for production):"
echo "   - Go to https://console.aws.amazon.com/cloudfront"
echo "   - Create distribution with S3 bucket as origin"
echo "   - Configure custom error pages (403->index.html, 404->index.html)"
echo ""
echo "4. Update CORS in Lambda to allow your CloudFront domain"
echo ""
echo "📊 Monitor costs:"
echo "   aws ce get-cost-and-usage --time-period Start=2024-10-01,End=2024-10-31 --granularity MONTHLY --metrics BlendedCost"
echo ""
echo "📚 Full deployment guide: DEPLOYMENT.md"
echo ""
echo -e "${GREEN}Happy coding! 🚀${NC}"
