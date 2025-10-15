# 🚀 AWS Deployment Guide - Realty Wizard

This guide walks you through deploying Realty Wizard to AWS using serverless architecture.

## 📋 Prerequisites Checklist

- ✅ AWS Account created
- ✅ Billing alerts configured ($10 threshold)
- ✅ IAM user created with AdministratorAccess
- ✅ AWS CLI installed and configured (`aws configure`)
- ✅ Node.js 18+ installed
- ✅ Go 1.21+ installed

## 🏗️ Architecture Overview

```
Internet
   │
   ├─→ CloudFront CDN ──→ S3 Bucket (React Frontend)
   │
   └─→ API Gateway ──→ Lambda Function (Go Backend)
                          │
                          └─→ /tmp/realty-wizard.db (SQLite)
```

**Cost Estimate:**
- Free Tier (12 months): ~$0-2/month
- After Free Tier: ~$5-10/month for light usage

---

## Phase 1: Deploy Backend (Lambda + API Gateway)

### Step 1: Build Lambda Package

```bash
cd backend
./deploy.sh
```

This creates `lambda.zip` (~10-15MB) with your Go binary.

### Step 2: Create Lambda Function via AWS Console

1. **Go to Lambda Console:** https://console.aws.amazon.com/lambda
2. **Click:** "Create function"
3. **Configuration:**
   - **Function name:** `realty-wizard-api`
   - **Runtime:** "Provide your own bootstrap on Amazon Linux 2023"
   - **Architecture:** x86_64
   - **Permissions:** Create new role with basic Lambda permissions
4. **Click:** "Create function"

### Step 3: Upload Lambda Code

1. **In the Lambda function page:**
   - **Code tab** → "Upload from" → ".zip file"
   - Select `backend/lambda.zip`
   - Click "Save"
2. **Wait for upload** (may take 30-60 seconds)

### Step 4: Configure Lambda Settings

1. **Configuration tab** → "General configuration" → "Edit"
   - **Memory:** 512 MB (sufficient for SQLite)
   - **Timeout:** 30 seconds
   - **Ephemeral storage:** 512 MB (default)
   - Click "Save"

2. **Configuration tab** → "Environment variables" → "Edit"
   - Add variable:
     - **Key:** `DB_PATH`
     - **Value:** `/tmp/realty-wizard.db`
   - Click "Save"

⚠️ **Note:** Using `/tmp` means database resets on cold starts (every ~15 min of inactivity). For persistence, see "Optional: Add EFS Storage" below.

### Step 5: Create API Gateway

1. **Go to API Gateway:** https://console.aws.amazon.com/apigateway
2. **Click:** "Create API"
3. **Choose:** "HTTP API" (cheaper than REST API)
4. **Integrations:**
   - Click "Add integration"
   - **Type:** Lambda
   - **Lambda function:** `realty-wizard-api`
   - **Version:** 2.0
   - **API name:** `realty-wizard-api`
5. **Routes:** Configure routes
   - **Method:** ANY
   - **Resource path:** `/{proxy+}`
6. **Stages:**
   - **Stage name:** `$default`
   - **Auto-deploy:** ✅ Enabled
7. **Click:** "Create"

### Step 6: Test the API

1. **Copy the Invoke URL** from API Gateway (looks like: `https://abc123.execute-api.us-east-1.amazonaws.com`)
2. **Test health endpoint:**

```bash
curl https://YOUR-API-ID.execute-api.us-east-1.amazonaws.com/api/health
# Should return: OK
```

3. **Test creating a project:**

```bash
curl -X POST https://YOUR-API-ID.execute-api.us-east-1.amazonaws.com/api/projects \
  -H "Content-Type: application/json" \
  -d '{
    "property_address": "123 Main St, Austin, TX",
    "seller_names": ["John Doe"],
    "seller_email": "john@example.com"
  }'
```

✅ **Backend is now deployed!**

**Save your API URL** - you'll need it for the frontend.

---

## Phase 2: Deploy Frontend (S3 + CloudFront)

### Step 1: Build React App

First, update the API URL in your frontend:

```bash
cd frontend

# Create production environment file
cat > .env.production << EOF
VITE_API_URL=https://YOUR-API-ID.execute-api.us-east-1.amazonaws.com
EOF

# Build for production
npm run build
```

This creates `frontend/dist/` with static files.

### Step 2: Create S3 Bucket for Hosting

```bash
# Replace with a unique bucket name
BUCKET_NAME="realty-wizard-frontend-$(date +%s)"

# Create bucket
aws s3 mb s3://$BUCKET_NAME --region us-east-1

# Enable static website hosting
aws s3 website s3://$BUCKET_NAME \
  --index-document index.html \
  --error-document index.html

echo "Bucket created: $BUCKET_NAME"
```

### Step 3: Upload Frontend Files

```bash
cd frontend

# Upload all files
aws s3 sync dist/ s3://$BUCKET_NAME/ --delete

# Set cache headers for assets
aws s3 cp s3://$BUCKET_NAME/assets s3://$BUCKET_NAME/assets \
  --recursive \
  --metadata-directive REPLACE \
  --cache-control "public, max-age=31536000, immutable"
```

### Step 4: Configure S3 Bucket Policy

```bash
# Allow public read access (CloudFront will access this)
cat > bucket-policy.json << EOF
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
  --bucket $BUCKET_NAME \
  --policy file://bucket-policy.json

rm bucket-policy.json
```

### Step 5: Create CloudFront Distribution

1. **Go to CloudFront:** https://console.aws.amazon.com/cloudfront
2. **Click:** "Create distribution"
3. **Origin settings:**
   - **Origin domain:** Select your S3 bucket (from dropdown)
   - **Origin path:** (leave blank)
   - **Name:** (auto-filled)
4. **Default cache behavior:**
   - **Viewer protocol policy:** "Redirect HTTP to HTTPS"
   - **Allowed HTTP methods:** GET, HEAD, OPTIONS
   - **Cache policy:** "CachingOptimized"
5. **Settings:**
   - **Price class:** "Use only North America and Europe" (cheapest)
   - **Default root object:** `index.html`
6. **Click:** "Create distribution"

⏳ **Wait 5-10 minutes** for CloudFront to deploy (status: "Deploying" → "Enabled")

### Step 6: Configure Custom Error Pages (SPA Routing)

1. **In your CloudFront distribution:**
   - **Error pages tab** → "Create custom error response"
   - **HTTP error code:** 403 Forbidden
   - **Customize error response:** Yes
   - **Response page path:** `/index.html`
   - **HTTP response code:** 200 OK
2. **Repeat for 404 Not Found**

### Step 7: Test Frontend

1. **Copy CloudFront Domain Name** (looks like: `d123abc456def.cloudfront.net`)
2. **Open in browser:** `https://d123abc456def.cloudfront.net`

✅ **Frontend is now deployed!**

---

## Phase 3: Update Frontend to Use Production API

Your frontend needs to know where the API is:

### Option A: Environment Variable (Recommended)

Update `frontend/src/api.ts` (or wherever you configure axios):

```typescript
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const api = axios.create({
  baseURL: API_URL,
});
```

Then rebuild and redeploy:

```bash
cd frontend
npm run build
aws s3 sync dist/ s3://$BUCKET_NAME/ --delete

# Invalidate CloudFront cache
aws cloudfront create-invalidation \
  --distribution-id YOUR-DISTRIBUTION-ID \
  --paths "/*"
```

### Option B: Hardcode (Quick & Dirty)

Just update the API URL directly in your code before building.

---

## 🎉 You're Live!

**Your app is now running on:**
- **Frontend:** `https://YOUR-CLOUDFRONT-DOMAIN.cloudfront.net`
- **Backend API:** `https://YOUR-API-ID.execute-api.us-east-1.amazonaws.com`

---

## 📊 Cost Breakdown

### Free Tier (First 12 months):
- **Lambda:** 1M requests/month + 400,000 GB-seconds compute (FREE)
- **API Gateway:** 1M requests/month (FREE)
- **S3:** 5GB storage + 20,000 GET requests (FREE)
- **CloudFront:** 1TB data transfer + 10M requests (FREE)

**Expected cost in free tier:** ~$0-2/month

### After Free Tier:
- **Lambda:** $0.20 per 1M requests + $0.0000166667/GB-second
- **API Gateway:** $1.00 per 1M requests
- **S3:** $0.023/GB/month
- **CloudFront:** $0.085/GB data transfer

**Expected cost with light usage (100 users, 10k requests/month):** ~$5-10/month

---

## 🔧 Common Commands

### Update Backend:
```bash
cd backend
./deploy.sh
aws lambda update-function-code \
  --function-name realty-wizard-api \
  --zip-file fileb://lambda.zip
```

### Update Frontend:
```bash
cd frontend
npm run build
aws s3 sync dist/ s3://YOUR-BUCKET-NAME/ --delete
aws cloudfront create-invalidation \
  --distribution-id YOUR-DISTRIBUTION-ID \
  --paths "/*"
```

### View Lambda Logs:
```bash
aws logs tail /aws/lambda/realty-wizard-api --follow
```

### Check Costs:
```bash
aws ce get-cost-and-usage \
  --time-period Start=2024-10-01,End=2024-10-31 \
  --granularity MONTHLY \
  --metrics BlendedCost
```

---

## 🔐 Security Hardening (Before Going Live)

### 1. Lock Down CORS

Update `backend/cmd/lambda/main.go`:

```go
AllowedOrigins: []string{"https://YOUR-CLOUDFRONT-DOMAIN.cloudfront.net"},
```

Rebuild and redeploy.

### 2. Add Custom Domain (Optional)

1. **Register domain** in Route 53 or elsewhere
2. **Request SSL certificate** in ACM (us-east-1 for CloudFront)
3. **Add alternate domain** to CloudFront distribution
4. **Create Route 53 record** pointing to CloudFront

### 3. Enable CloudFront Access Logs

Track who's accessing your site.

---

## 📈 Monitoring & Alerts

### Set Up CloudWatch Alarms:

```bash
# Alert if Lambda errors spike
aws cloudwatch put-metric-alarm \
  --alarm-name realty-wizard-lambda-errors \
  --alarm-description "Alert on Lambda errors" \
  --metric-name Errors \
  --namespace AWS/Lambda \
  --statistic Sum \
  --period 300 \
  --threshold 5 \
  --comparison-operator GreaterThanThreshold \
  --evaluation-periods 1 \
  --dimensions Name=FunctionName,Value=realty-wizard-api
```

---

## 🐛 Troubleshooting

### Frontend shows CORS error:
- Check Lambda CORS configuration includes your CloudFront domain
- Verify API Gateway has CORS enabled
- Check browser console for specific error

### Lambda timeout:
- Increase timeout in Lambda configuration
- Check CloudWatch Logs for slow queries

### Database resets frequently:
- This is expected with /tmp storage
- Upgrade to EFS for persistence (see below)

### CloudFront not updating:
- Create cache invalidation after S3 upload
- Wait 5-10 minutes for changes to propagate

---

## 🚀 Optional: Add EFS for Persistent Database

If you need your database to persist across Lambda cold starts:

### 1. Create EFS File System:
```bash
aws efs create-file-system \
  --performance-mode generalPurpose \
  --throughput-mode bursting \
  --encrypted \
  --tags Key=Name,Value=realty-wizard-efs
```

### 2. Create Mount Target:
- Go to EFS Console
- Select your file system
- Create mount targets in same VPC/subnets as Lambda

### 3. Attach to Lambda:
- Lambda Console → Configuration → File systems
- Add file system
- Mount path: `/mnt/efs`

### 4. Update Environment Variable:
```bash
DB_PATH=/mnt/efs/realty-wizard.db
```

**Cost:** ~$0.30/GB/month + $0.03/GB read/write

---

## 📚 Next Steps

1. ✅ Set up custom domain
2. ✅ Add authentication (AWS Cognito)
3. ✅ Enable CloudWatch Logs insights
4. ✅ Set up CI/CD with GitHub Actions
5. ✅ Add EFS for database persistence
6. ✅ Implement backups (S3 snapshots)

---

## 🆘 Need Help?

- **AWS Free Tier:** https://aws.amazon.com/free
- **Lambda Pricing:** https://aws.amazon.com/lambda/pricing
- **CloudFront Pricing:** https://aws.amazon.com/cloudfront/pricing
- **Support:** AWS Forums or Stack Overflow

---

**Happy Deploying! 🎉**
