# ⚡ AWS Deployment - Quick Start (10 Minutes)

This is the **fastest path** to get Realty Wizard live on AWS. For detailed instructions, see [DEPLOYMENT.md](DEPLOYMENT.md).

---

## 🎯 What You'll Deploy

```
Your App URL
    │
    ├─→ CloudFront (Frontend) ───→ S3 Bucket
    │
    └─→ API Gateway (Backend) ───→ Lambda Function
```

**Cost:** ~$0-5/month (Free tier covers most usage)

---

## ✅ Step 1: AWS Account Setup (5 min)

### 1.1 Create Account
1. Go to https://aws.amazon.com/ → "Create an AWS Account"
2. Use your email, set password
3. Add payment method (required, but won't be charged in free tier)
4. Choose **"Basic Support - Free"**

### 1.2 Set Billing Alert (CRITICAL!)
1. Sign in → Search "Billing" → "Billing Preferences"
2. Enable "Receive Billing Alerts"
3. Go to "Budgets" → "Create budget"
4. Use template: "Zero spend budget" OR set $10 threshold
5. Enter your email for alerts

### 1.3 Create IAM User
1. Search "IAM" → "Users" → "Create user"
2. Username: `realty-wizard-admin`
3. Enable "AWS Management Console access"
4. Attach policy: `AdministratorAccess`
5. Download credentials CSV

### 1.4 Install & Configure AWS CLI
```bash
# Install (Mac)
brew install awscli

# Configure (will prompt you for 4 answers)
aws configure
# AWS Access Key ID: (paste from credentials CSV file)
# AWS Secret Access Key: (paste from credentials CSV file)
# Default region name: us-east-1
# Default output format: json

# Test
aws sts get-caller-identity
```

---

## 🚀 Step 2: Deploy Backend (2 min)

### 2.1 Build Lambda Package
```bash
cd backend
./deploy.sh
```

This creates `lambda.zip`.

### 2.2 Create Lambda Function (AWS Console)
1. Go to: https://console.aws.amazon.com/lambda
2. **Create function** → "Author from scratch"
3. Settings:
   - **Function name:** `realty-wizard-api`
   - **Runtime:** "Provide your own bootstrap on Amazon Linux 2023"
   - **Architecture:** x86_64
4. **Create function**

### 2.3 Upload Code
1. **Code tab** → "Upload from" → ".zip file"
2. Select `backend/lambda.zip`
3. **Save**

### 2.4 Configure
1. **Configuration** → "General configuration" → "Edit"
   - Memory: 512 MB
   - Timeout: 30 seconds
   - **Save**

2. **Environment variables** → "Edit" → "Add"
   - Key: `DB_PATH`
   - Value: `/tmp/realty-wizard.db`
   - **Save**

### 2.5 Create API Gateway
1. Go to: https://console.aws.amazon.com/apigateway
2. **Create API** → Choose "HTTP API"
3. **Add integration:**
   - Type: Lambda
   - Function: `realty-wizard-api`
   - API name: `realty-wizard-api`
4. **Configure routes:**
   - Method: ANY
   - Path: `/{proxy+}`
5. **Create**

### 2.6 Test Backend
```bash
# Copy your API URL from API Gateway console
# It looks like: https://abc123xyz.execute-api.us-east-1.amazonaws.com

curl https://YOUR-API-URL/api/health
# Should return: OK
```

✅ **Backend is live!** Save your API URL.

---

## 🎨 Step 3: Deploy Frontend (3 min)

### 3.1 Configure & Build
```bash
cd frontend

# Create production config
cat > .env.production << EOF
VITE_API_URL=https://YOUR-API-URL
EOF

# Build
npm run build
```

### 3.2 Create S3 Bucket
```bash
# Use a unique name (bucket names are global)
BUCKET_NAME="realty-wizard-$(date +%s)"

# Create bucket
aws s3 mb s3://$BUCKET_NAME --region us-east-1

# Enable static website hosting
aws s3 website s3://$BUCKET_NAME \
  --index-document index.html \
  --error-document index.html

echo "Your bucket: $BUCKET_NAME"
```

### 3.3 Upload Files
```bash
# Upload all files
aws s3 sync dist/ s3://$BUCKET_NAME/ --delete

# Make bucket public
cat > /tmp/policy.json << EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": "*",
    "Action": "s3:GetObject",
    "Resource": "arn:aws:s3:::${BUCKET_NAME}/*"
  }]
}
EOF

aws s3api put-bucket-policy --bucket $BUCKET_NAME --policy file:///tmp/policy.json
rm /tmp/policy.json
```

### 3.4 Create CloudFront Distribution
1. Go to: https://console.aws.amazon.com/cloudfront
2. **Create distribution**
3. **Origin:**
   - Domain: Select your S3 bucket
4. **Settings:**
   - Viewer protocol: "Redirect HTTP to HTTPS"
   - Default root object: `index.html`
5. **Create**
6. Wait 5-10 minutes for deployment

### 3.5 Configure SPA Routing
1. In CloudFront distribution → **Error pages** tab
2. **Create custom error response:**
   - HTTP code: 403
   - Response page: `/index.html`
   - HTTP response: 200
3. Repeat for HTTP code: 404

---

## 🎉 You're Live!

**Your app is now running at:**
- Frontend: `https://YOUR-CLOUDFRONT-DOMAIN.cloudfront.net`
- Backend API: `https://YOUR-API-GATEWAY-URL`

### Test It:
1. Open CloudFront URL in browser
2. Create a new transaction
3. Fill in property details
4. Verify data persists

---

## 💰 Costs (12-Month Free Tier)

**What's FREE:**
- Lambda: 1M requests/month + 400K GB-seconds
- API Gateway: 1M requests/month
- S3: 5GB storage + 20K GET requests
- CloudFront: 1TB transfer + 10M requests

**Your likely usage:** Well within free tier limits

**Expected cost:** $0-2/month for first year

---

## 🔄 Update Your App

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

# Invalidate CloudFront cache
aws cloudfront create-invalidation \
  --distribution-id YOUR-DISTRIBUTION-ID \
  --paths "/*"
```

---

## 🐛 Common Issues

### "CORS error" in browser:
- Update Lambda CORS to allow your CloudFront domain
- Edit `backend/cmd/lambda/main.go`:
  ```go
  AllowedOrigins: []string{"https://YOUR-CLOUDFRONT.cloudfront.net"},
  ```
- Rebuild and redeploy Lambda

### Database resets frequently:
- Expected with `/tmp` storage
- For persistence, add EFS (see DEPLOYMENT.md)

### Frontend not updating:
- Invalidate CloudFront cache (see update commands above)
- Wait 5-10 minutes for propagation

---

## 📚 Next Steps

- [ ] Set up custom domain name
- [ ] Add authentication (AWS Cognito)
- [ ] Enable CloudWatch monitoring
- [ ] Set up CI/CD (GitHub Actions)
- [ ] Add EFS for persistent database

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed guides on these topics.

---

## 🆘 Need Help?

1. Check CloudWatch Logs:
   ```bash
   aws logs tail /aws/lambda/realty-wizard-api --follow
   ```

2. View costs:
   ```bash
   aws ce get-cost-and-usage \
     --time-period Start=2024-10-01,End=2024-10-31 \
     --granularity MONTHLY \
     --metrics BlendedCost
   ```

3. AWS Support:
   - AWS Forums: https://forums.aws.amazon.com
   - Stack Overflow: Tag questions with `aws-lambda`

---

**🎉 Congratulations! Your app is live on AWS! 🚀**
