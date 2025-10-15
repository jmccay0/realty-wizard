# 🔑 AWS Credentials Setup - Visual Guide

This guide shows you **exactly** where to get your AWS credentials and how to use them.

---

## 📖 What Are AWS Credentials?

Think of AWS credentials as your **username and password** for the command line.

- **Access Key ID** = Your username (public, like "AKIAI44QH8DHBEXAMPLE")
- **Secret Access Key** = Your password (private, like a 40-character random string)

---

## 🎯 Step-by-Step: Getting Your Credentials

### Step 1: Create IAM User (in AWS Console)

1. **Sign in to AWS Console:** https://console.aws.amazon.com
2. **Search bar:** Type "IAM" → Click "IAM"
3. **Left sidebar:** Click "Users"
4. **Click:** "Create user"

### Step 2: Configure User

5. **User name:** `realty-wizard-admin`
6. **Check:** ☑️ "Provide user access to the AWS Management Console"
7. **Select:** "I want to create an IAM user"
8. **Password:** Create a strong password (save it!)
9. **Uncheck:** "User must create a new password at next sign-in"
10. **Click:** "Next"

### Step 3: Set Permissions

11. **Select:** "Attach policies directly"
12. **Search:** Type "AdministratorAccess"
13. **Check:** ☑️ AdministratorAccess
14. **Click:** "Next" → "Create user"

### Step 4: Download Credentials

15. **Click:** "View user" (or go to Users → `realty-wizard-admin`)
16. **Tab:** "Security credentials"
17. **Section:** "Access keys"
18. **Click:** "Create access key"
19. **Use case:** Select "Command Line Interface (CLI)"
20. **Check:** ☑️ "I understand the above recommendation"
21. **Click:** "Next"
22. **Description:** "Local development machine"
23. **Click:** "Create access key"

### Step 5: Save Your Credentials

24. **You'll see a screen with:**
    ```
    Access key ID: AKIAIOSFODNN7EXAMPLE
    Secret access key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
    ```

25. **IMPORTANT - Choose ONE:**
    - **Option A:** Click "Download .csv file" (saves to Downloads folder)
    - **Option B:** Copy both values to a text file

26. **Click:** "Done"

⚠️ **WARNING:** You can NEVER see the secret access key again after closing this screen!

---

## 💾 What the CSV File Looks Like

If you downloaded the CSV, open `Downloads/realty-wizard-admin_credentials.csv`:

```csv
User name,Access key ID,Secret access key
realty-wizard-admin,AKIAIOSFODNN7EXAMPLE,wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

**Visual breakdown:**
```
Column 1: User name          → realty-wizard-admin
Column 2: Access key ID      → AKIAIOSFODNN7EXAMPLE
Column 3: Secret access key  → wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
                                 ↑ THIS IS WHAT YOU NEED ↑
```

---

## 🖥️ Using Credentials with AWS CLI

### Install AWS CLI (if not already installed)

```bash
# Mac
brew install awscli

# Verify installation
aws --version
# Should show: aws-cli/2.x.x
```

### Configure AWS CLI

Open Terminal and run:

```bash
aws configure
```

**You'll be prompted 4 times. Here's what to enter:**

#### Prompt 1: Access Key ID
```
AWS Access Key ID [None]:
```
👉 **Paste:** `AKIAIOSFODNN7EXAMPLE` (from CSV column 2)

#### Prompt 2: Secret Access Key
```
AWS Secret Access Key [None]:
```
👉 **Paste:** `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY` (from CSV column 3)

#### Prompt 3: Region
```
Default region name [None]:
```
👉 **Type:** `us-east-1` (or `us-west-2` if you prefer West Coast)

**Common regions:**
- `us-east-1` - Virginia (most common, cheapest)
- `us-west-2` - Oregon
- `eu-west-1` - Ireland
- `ap-southeast-1` - Singapore

#### Prompt 4: Output Format
```
Default output format [None]:
```
👉 **Type:** `json`

---

## ✅ Full Example Session

Here's what it looks like in your terminal:

```bash
$ aws configure
AWS Access Key ID [None]: AKIAI44QH8DHBEXAMPLE
AWS Secret Access Key [None]: je7MtGbClwBF/2Zp9Utk/h3yCo8nvbEXAMPLEKEY
Default region name [None]: us-east-1
Default output format [None]: json
```

**That's it!** No confirmation message - it just returns to the prompt.

---

## 🧪 Test Your Configuration

### Test 1: Check Identity
```bash
aws sts get-caller-identity
```

**Expected output:**
```json
{
    "UserId": "AIDAI23HXX2EXAMPLE",
    "Account": "123456789012",
    "Arn": "arn:aws:iam::123456789012:user/realty-wizard-admin"
}
```

### Test 2: List Regions
```bash
aws ec2 describe-regions --output table
```

**Expected output:** A table of AWS regions

### Test 3: Check S3 Buckets
```bash
aws s3 ls
```

**Expected output:** Empty (or list of buckets if you have any)

---

## 📁 Where Are Credentials Stored?

AWS CLI saves your credentials in:

**Mac/Linux:**
```
~/.aws/credentials
~/.aws/config
```

**Windows:**
```
C:\Users\YourName\.aws\credentials
C:\Users\YourName\.aws\config
```

### View Your Stored Credentials

```bash
# Mac/Linux
cat ~/.aws/credentials
```

**Contents:**
```ini
[default]
aws_access_key_id = AKIAIOSFODNN7EXAMPLE
aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

```bash
cat ~/.aws/config
```

**Contents:**
```ini
[default]
region = us-east-1
output = json
```

---

## 🔒 Security Best Practices

### ✅ DO:
- Keep credentials private (never commit to Git)
- Use IAM user (not root account)
- Rotate access keys every 90 days
- Use different keys for different projects
- Delete unused access keys

### ❌ DON'T:
- Share credentials in Slack/email
- Commit credentials to GitHub
- Use root account access keys
- Hardcode credentials in code
- Leave credentials in Docker images

### Add to `.gitignore`:
```bash
# Add to your .gitignore
.aws/
*.csv
credentials.csv
*_credentials.csv
.env.local
.env.production.local
```

---

## 🐛 Troubleshooting

### "Unable to locate credentials"
**Problem:** AWS CLI can't find your credentials

**Solution:**
```bash
# Re-run configuration
aws configure

# Or check if files exist
ls -la ~/.aws/
```

### "The security token included in the request is invalid"
**Problem:** Credentials are wrong or expired

**Solution:**
1. Go to IAM Console → Users → Your user → Security credentials
2. Delete old access key
3. Create new access key
4. Run `aws configure` again with new keys

### "Access Denied"
**Problem:** Your IAM user doesn't have permissions

**Solution:**
1. Go to IAM Console → Users → Your user
2. Add policy: `AdministratorAccess`

---

## 🔄 Multiple AWS Accounts/Profiles

If you work with multiple AWS accounts:

```bash
# Configure a named profile
aws configure --profile personal
aws configure --profile work

# Use a specific profile
aws s3 ls --profile personal
aws s3 ls --profile work

# Set default profile for session
export AWS_PROFILE=personal
```

**Credentials file with profiles:**
```ini
[default]
aws_access_key_id = AKIAIOSFODNN7EXAMPLE
aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY

[personal]
aws_access_key_id = AKIAI44QH8DHBEXAMPLE
aws_secret_access_key = je7MtGbClwBF/2Zp9Utk/h3yCo8nvbEXAMPLEKEY

[work]
aws_access_key_id = AKIAJFHD73EXAMPLE
aws_secret_access_key = xJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

---

## 📚 Quick Reference Card

| What | Where to Find |
|------|---------------|
| **Access Key ID** | IAM Console → Users → Security credentials → Create access key |
| **Secret Key** | Only shown once when creating access key (download CSV!) |
| **Region** | Choose: `us-east-1` (Virginia) or `us-west-2` (Oregon) |
| **Output Format** | Choose: `json` |
| **Stored At** | `~/.aws/credentials` and `~/.aws/config` |
| **Test Command** | `aws sts get-caller-identity` |

---

## 🆘 Still Stuck?

1. **Verify AWS CLI installed:**
   ```bash
   which aws
   aws --version
   ```

2. **Check credentials file exists:**
   ```bash
   cat ~/.aws/credentials
   ```

3. **Re-configure from scratch:**
   ```bash
   rm -rf ~/.aws
   aws configure
   ```

4. **Check IAM user exists:**
   - Go to https://console.aws.amazon.com/iam
   - Click "Users"
   - Verify your user is listed

---

**Ready to deploy?** Head back to [AWS-QUICKSTART.md](AWS-QUICKSTART.md)!
