# ✅ AWS Deployment Checklist

Use this checklist to track your deployment progress.

---

## 📋 Pre-Deployment

### AWS Account Setup
- [ ] Create AWS account at aws.amazon.com
- [ ] Verify email address
- [ ] Add payment method
- [ ] Choose "Basic Support - Free"
- [ ] **CRITICAL:** Set up billing alert ($10 threshold)
- [ ] Create "Zero spend budget" for immediate alerts

### IAM Security
- [ ] Create IAM user: `realty-wizard-admin`
- [ ] Attach `AdministratorAccess` policy
- [ ] Download credentials CSV
- [ ] **Never use root account for daily work**

### Local Setup
- [ ] Install AWS CLI: `brew install awscli`
- [ ] Configure AWS CLI: `aws configure`
  - [ ] Enter Access Key ID
  - [ ] Enter Secret Access Key
  - [ ] Set region: `us-east-1`
  - [ ] Set output: `json`
- [ ] Test CLI: `aws sts get-caller-identity`

---

## 🔧 Backend Deployment

### Build Lambda Package
- [ ] `cd backend`
- [ ] `./deploy.sh`
- [ ] Verify `lambda.zip` created (~10-15 MB)

### Create Lambda Function
- [ ] Go to Lambda Console
- [ ] Create function: `realty-wizard-api`
- [ ] Runtime: "Provide your own bootstrap on Amazon Linux 2023"
- [ ] Architecture: x86_64
- [ ] Upload `lambda.zip`

### Configure Lambda
- [ ] Memory: 512 MB
- [ ] Timeout: 30 seconds
- [ ] Environment variable: `DB_PATH=/tmp/realty-wizard.db`

### Create API Gateway
- [ ] Go to API Gateway Console
- [ ] Create HTTP API (not REST API)
- [ ] Name: `realty-wizard-api`
- [ ] Integration: Lambda → `realty-wizard-api`
- [ ] Route: `ANY /{proxy+}`
- [ ] Auto-deploy: Enabled

### Test Backend
- [ ] Copy API Gateway Invoke URL
- [ ] Test: `curl https://YOUR-API-URL/api/health`
- [ ] Should return: `OK`
- [ ] Save API URL for frontend config

**API URL:** _______________________________________

---

## 🎨 Frontend Deployment

### Build Frontend
- [ ] `cd frontend`
- [ ] Create `.env.production`:
  ```
  VITE_API_URL=https://YOUR-API-URL
  ```
- [ ] `npm run build`
- [ ] Verify `dist/` directory created

### Create S3 Bucket
- [ ] Choose unique bucket name: `realty-wizard-XXXXXX`
- [ ] Create bucket: `aws s3 mb s3://BUCKET-NAME`
- [ ] Enable static website hosting
- [ ] Set bucket policy (public read)

**Bucket Name:** _______________________________________

### Upload Files
- [ ] `aws s3 sync dist/ s3://BUCKET-NAME/ --delete`
- [ ] Set cache headers for assets
- [ ] Test: `http://BUCKET-NAME.s3-website-us-east-1.amazonaws.com`

### Create CloudFront Distribution
- [ ] Go to CloudFront Console
- [ ] Create distribution
- [ ] Origin: Your S3 bucket
- [ ] Viewer protocol: Redirect HTTP to HTTPS
- [ ] Default root object: `index.html`
- [ ] Wait 5-10 minutes for deployment

### Configure Error Pages (SPA Routing)
- [ ] Add custom error response: 403 → `/index.html` (200)
- [ ] Add custom error response: 404 → `/index.html` (200)

**CloudFront URL:** _______________________________________

---

## 🔒 Security & Configuration

### CORS Configuration
- [ ] Update Lambda CORS to allow CloudFront domain
- [ ] Edit `backend/cmd/lambda/main.go`
- [ ] Change `AllowedOrigins` from `*` to CloudFront URL
- [ ] Rebuild and redeploy Lambda

### Test End-to-End
- [ ] Open CloudFront URL in browser
- [ ] Create new transaction
- [ ] Add property details
- [ ] Complete disclosure
- [ ] Verify API calls work (check browser console)

---

## 📊 Monitoring & Costs

### Set Up Monitoring
- [ ] Enable CloudWatch Logs for Lambda
- [ ] View logs: `aws logs tail /aws/lambda/realty-wizard-api --follow`
- [ ] Set up CloudWatch alarm for Lambda errors

### Check Costs
- [ ] View current costs in AWS Billing Dashboard
- [ ] Run: `aws ce get-cost-and-usage --time-period Start=2024-10-01,End=2024-10-31 --granularity MONTHLY --metrics BlendedCost`
- [ ] Verify within free tier limits

**Current Monthly Cost:** $_______

---

## 🚀 Optional Enhancements

### Custom Domain
- [ ] Register domain (Route 53 or elsewhere)
- [ ] Request SSL certificate in ACM (us-east-1)
- [ ] Add alternate domain to CloudFront
- [ ] Create Route 53 A record → CloudFront

**Custom Domain:** _______________________________________

### Persistent Database (EFS)
- [ ] Create EFS file system
- [ ] Create mount targets in Lambda VPC
- [ ] Attach EFS to Lambda at `/mnt/efs`
- [ ] Update `DB_PATH=/mnt/efs/realty-wizard.db`
- [ ] Test data persists across cold starts

**Additional Cost:** ~$0.30/GB/month

### CI/CD Pipeline
- [ ] Create GitHub Actions workflow
- [ ] Set up AWS credentials as secrets
- [ ] Auto-deploy on push to `main` branch

---

## 📝 Notes & Troubleshooting

### Important URLs
- AWS Console: https://console.aws.amazon.com
- Lambda: https://console.aws.amazon.com/lambda
- API Gateway: https://console.aws.amazon.com/apigateway
- S3: https://console.aws.amazon.com/s3
- CloudFront: https://console.aws.amazon.com/cloudfront
- Billing: https://console.aws.amazon.com/billing

### Common Issues
**CORS errors:**
- Update Lambda CORS configuration
- Verify API Gateway CORS settings
- Check browser console for specific error

**Database resets:**
- Expected with `/tmp` storage
- Upgrade to EFS for persistence

**CloudFront not updating:**
- Create invalidation: `aws cloudfront create-invalidation --distribution-id ID --paths "/*"`
- Wait 5-10 minutes for propagation

**Lambda timeout:**
- Increase timeout in Configuration
- Check CloudWatch Logs for slow queries

### Update Commands

**Update Backend:**
```bash
cd backend
./deploy.sh
aws lambda update-function-code --function-name realty-wizard-api --zip-file fileb://lambda.zip
```

**Update Frontend:**
```bash
cd frontend
npm run build
aws s3 sync dist/ s3://BUCKET-NAME/ --delete
aws cloudfront create-invalidation --distribution-id DISTRIBUTION-ID --paths "/*"
```

---

## 🎉 Deployment Complete!

- [ ] App is live and accessible
- [ ] All tests passing
- [ ] Monitoring configured
- [ ] Billing alerts active
- [ ] Documentation updated with URLs

**Next Steps:**
1. Share app URL with users
2. Monitor CloudWatch metrics
3. Review costs weekly
4. Plan next features
5. Set up CI/CD

---

**Deployment Date:** _______________

**Deployed By:** _______________

**Notes:**
_______________________________________
_______________________________________
_______________________________________
