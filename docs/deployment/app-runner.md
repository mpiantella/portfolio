# AWS App Runner Deployment

AWS App Runner is a fully managed service that makes it easy to deploy containerized web applications and APIs at scale.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                    AWS App Runner Architecture                       │
└─────────────────────────────────────────────────────────────────────┘

                         ┌──────────────────┐
                         │  Developer       │
                         │  Pushes Code     │
                         └────────┬─────────┘
                                  │
                                  ▼
                         ┌──────────────────┐
                         │  GitHub Actions  │
                         │  (CI/CD)         │
                         └────────┬─────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │  1. Build Docker Image    │
                    │  2. Push to ECR           │
                    └─────────────┬─────────────┘
                                  │
                                  ▼
                    ┌──────────────────────────┐
                    │   Amazon ECR             │
                    │   (Container Registry)   │
                    └──────────────┬───────────┘
                                   │
                                   │ Auto-deploy on new image
                                   │
                                   ▼
                    ┌──────────────────────────┐
                    │   AWS App Runner         │
                    │   ┌────────────────┐     │
                    │   │  Auto-scaling  │     │
                    │   │  Load Balancer │     │
                    │   │  SSL/TLS       │     │
                    │   └────────────────┘     │
                    │                          │
                    │   ┌──────┐  ┌──────┐    │
                    │   │Task 1│  │Task 2│    │
                    │   │:8080 │  │:8080 │    │
                    │   └──────┘  └──────┘    │
                    └──────────────┬───────────┘
                                   │
                                   ▼
                    ┌──────────────────────────┐
                    │  HTTPS URL               │
                    │  (Auto-provisioned)      │
                    │  *.awsapprunner.com      │
                    └──────────────────────────┘
                                   │
                                   ▼
                         ┌──────────────────┐
                         │  End Users       │
                         └──────────────────┘
```

## Resources Required

### AWS Resources
1. **Amazon ECR Repository**
   - Stores Docker images
   - Name: `portfolio-website`

2. **AWS App Runner Service**
   - Pulls images from ECR
   - Handles all infrastructure
   - Auto-provisions HTTPS endpoint

3. **IAM Roles (Auto-created)**
   - **Instance Role**: Allows App Runner to pull from ECR
   - **Service Role**: Allows App Runner to manage resources

### GitHub Secrets
```bash
AWS_ACCESS_KEY_ID          # IAM user access key
AWS_SECRET_ACCESS_KEY      # IAM user secret key
AWS_REGION                 # e.g., us-east-1
ECR_REPOSITORY            # e.g., portfolio-website
APP_RUNNER_SERVICE_ARN    # ARN of App Runner service (after creation)
```

## Step-by-Step Setup

### 1. Create IAM User (One-time)

```bash
# Create IAM user for GitHub Actions
aws iam create-user --user-name github-actions-portfolio

# Attach required policies
aws iam attach-user-policy \
  --user-name github-actions-portfolio \
  --policy-arn arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser

# Create access key
aws iam create-access-key --user-name github-actions-portfolio
```

**Required Permissions:**
- `ecr:*` - ECR operations
- `apprunner:*` - App Runner operations (if auto-deploying)

### 2. Create ECR Repository

```bash
# Create repository
aws ecr create-repository \
  --repository-name portfolio-website \
  --region us-east-1

# Output will include repository URI
```

### 3. Configure GitHub Secrets

In your GitHub repository, go to Settings → Secrets and variables → Actions:

```
AWS_ACCESS_KEY_ID: AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_REGION: us-east-1
ECR_REPOSITORY: portfolio-website
```

### 4. Update GitHub Actions Workflow

Create `.github/workflows/deploy-app-runner.yml`:

```yaml
name: Deploy to AWS App Runner

on:
  push:
    branches: [ main ]
  workflow_dispatch: {}

env:
  AWS_REGION: ${{ secrets.AWS_REGION }}
  ECR_REPOSITORY: ${{ secrets.ECR_REPOSITORY }}

jobs:
  deploy:
    name: Deploy to App Runner
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ${{ secrets.AWS_REGION }}

      - name: Login to Amazon ECR
        id: ecr_login
        uses: aws-actions/amazon-ecr-login@v1

      - name: Build, tag, and push image to ECR
        id: build-image
        env:
          ECR_REGISTRY: ${{ steps.ecr_login.outputs.registry }}
          IMAGE_TAG: ${{ github.sha }}
        run: |
          docker build -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG .
          docker tag $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG $ECR_REGISTRY/$ECR_REPOSITORY:latest
          docker push $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG
          docker push $ECR_REGISTRY/$ECR_REPOSITORY:latest
          echo "image=$ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG" >> $GITHUB_OUTPUT

      - name: Deploy to App Runner
        id: deploy-apprunner
        uses: awslabs/amazon-app-runner-deploy@main
        with:
          service: portfolio-website
          image: ${{ steps.build-image.outputs.image }}
          access-role-arn: ${{ secrets.APP_RUNNER_ROLE_ARN }}
          region: ${{ secrets.AWS_REGION }}
          cpu: 1
          memory: 2
          port: 8080
          wait-for-service-stability-seconds: 600

      - name: App Runner URL
        run: echo "App Runner URL ${{ steps.deploy-apprunner.outputs.service-url }}"
```

### 5. Create App Runner Service (One-time)

#### Option A: AWS Console (Easiest)

1. Go to AWS App Runner console
2. Click "Create service"
3. Choose "Container registry" → "Amazon ECR"
4. Select your ECR repository
5. Configure:
   - **Service name**: portfolio-website
   - **Port**: 8080
   - **CPU**: 1 vCPU
   - **Memory**: 2 GB
   - **Auto-scaling**: 1-10 instances
6. Click "Create & deploy"

#### Option B: AWS CLI

```bash
# First, create an IAM role for App Runner
aws iam create-role \
  --role-name AppRunnerECRAccessRole \
  --assume-role-policy-document file://app-runner-trust-policy.json

# Attach ECR read policy
aws iam attach-role-policy \
  --role-name AppRunnerECRAccessRole \
  --policy-arn arn:aws:iam::aws:policy/service-role/AWSAppRunnerServicePolicyForECRAccess

# Create App Runner service
aws apprunner create-service \
  --service-name portfolio-website \
  --source-configuration file://app-runner-config.json \
  --instance-configuration Cpu=1024,Memory=2048 \
  --region us-east-1
```

**app-runner-trust-policy.json:**
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "build.apprunner.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

**app-runner-config.json:**
```json
{
  "ImageRepository": {
    "ImageIdentifier": "YOUR_ECR_URI:latest",
    "ImageRepositoryType": "ECR",
    "ImageConfiguration": {
      "Port": "8080"
    }
  },
  "AutoDeploymentsEnabled": true,
  "AuthenticationConfiguration": {
    "AccessRoleArn": "arn:aws:iam::YOUR_ACCOUNT:role/AppRunnerECRAccessRole"
  }
}
```

### 6. Deploy Application

```bash
# Push to main branch
git add .
git commit -m "Deploy to App Runner"
git push origin main

# GitHub Actions will automatically:
# 1. Build Docker image
# 2. Push to ECR
# 3. Deploy to App Runner
```

## Configuration Options

### Auto-Scaling Configuration

```json
{
  "MaxConcurrency": 100,
  "MaxSize": 10,
  "MinSize": 1
}
```

### Health Check Configuration

```json
{
  "HealthCheckConfiguration": {
    "Protocol": "HTTP",
    "Path": "/api/health",
    "Interval": 10,
    "Timeout": 5,
    "HealthyThreshold": 1,
    "UnhealthyThreshold": 5
  }
}
```

### Environment Variables

```bash
aws apprunner update-service \
  --service-arn YOUR_SERVICE_ARN \
  --source-configuration '{
    "ImageRepository": {
      "RuntimeEnvironmentVariables": {
        "LOG_LEVEL": "info",
        "PORT": "8080"
      }
    }
  }'
```

## Cost Breakdown

### Pricing Model
App Runner charges for:
1. **Compute**: Time your application is running
2. **Memory**: Memory allocated to your application
3. **Requests**: Number of requests processed

### Monthly Cost Estimate

**Configuration**: 1 vCPU, 2 GB Memory

| Component | Usage | Rate | Monthly Cost |
|-----------|-------|------|--------------|
| Compute | 730 hours | $0.064/vCPU-hour | $46.72 |
| Memory | 730 hours | $0.007/GB-hour | $10.22 |
| Requests | 100k | $0.40/million | $0.04 |
| **Total** | | | **~$57/month** |

**With Provisioned Concurrency OFF** (recommended for low traffic):
- Pay only when serving requests
- Can be as low as **$5-10/month** for personal portfolio

### Cost Optimization

```bash
# Reduce to minimal configuration
aws apprunner update-service \
  --service-arn YOUR_SERVICE_ARN \
  --instance-configuration Cpu=256,Memory=512

# Monthly cost drops to ~$15
```

## Custom Domain Setup

### 1. Add Custom Domain

```bash
aws apprunner associate-custom-domain \
  --service-arn YOUR_SERVICE_ARN \
  --domain-name portfolio.example.com
```

### 2. Verify Domain Ownership

App Runner will provide validation records. Add them to your DNS:

```
Type: CNAME
Name: _app-runner-validation.portfolio.example.com
Value: [provided-by-app-runner]
```

### 3. Update DNS

```
Type: CNAME
Name: portfolio.example.com
Value: [your-app-runner-url].awsapprunner.com
```

## Monitoring & Logs

### View Logs

```bash
# Get service ARN
aws apprunner list-services

# View logs in CloudWatch
aws logs tail /aws/apprunner/portfolio-website --follow
```

### Metrics

Available in CloudWatch:
- `2xxStatusResponses`
- `4xxStatusResponses`
- `5xxStatusResponses`
- `RequestCount`
- `ActiveInstances`
- `CPUUtilization`
- `MemoryUtilization`

### Create Alarms

```bash
aws cloudwatch put-metric-alarm \
  --alarm-name portfolio-high-error-rate \
  --alarm-description "Alert on high 5xx errors" \
  --metric-name 5xxStatusResponses \
  --namespace AWS/AppRunner \
  --statistic Sum \
  --period 300 \
  --threshold 10 \
  --comparison-operator GreaterThanThreshold
```

## Troubleshooting

### Service Won't Start

1. Check logs in CloudWatch
2. Verify Docker image runs locally
3. Check health check endpoint
4. Verify port 8080 is exposed

```bash
# Test locally
docker run -p 8080:8080 your-image:latest
curl http://localhost:8080/api/health
```

### Deployment Fails

```bash
# Check service status
aws apprunner describe-service --service-arn YOUR_ARN

# View operation logs
aws logs tail /aws/apprunner/portfolio-website/operation --follow
```

### High Costs

1. Disable provisioned concurrency
2. Reduce CPU/Memory allocation
3. Set up auto-pause for development environments

## Maintenance

### Update Service

```bash
# Trigger new deployment
aws apprunner start-deployment --service-arn YOUR_ARN
```

### Pause Service (Stop Costs)

```bash
# Pause service
aws apprunner pause-service --service-arn YOUR_ARN

# Resume service
aws apprunner resume-service --service-arn YOUR_ARN
```

### Delete Service

```bash
aws apprunner delete-service --service-arn YOUR_ARN
```

## Advantages

✅ **Fully Managed** - No infrastructure to manage
✅ **Auto-Scaling** - Scales from 1 to 25 instances automatically
✅ **Built-in HTTPS** - SSL certificates auto-provisioned
✅ **Fast Deploys** - Updates in ~2-3 minutes
✅ **Pay Per Use** - Only pay when serving requests
✅ **Simple Setup** - Minimal configuration required

## Limitations

❌ **Limited Regions** - Not available in all AWS regions
❌ **No VPC Control** - Cannot customize networking
❌ **Max Instance Size** - Limited to 4 vCPU, 12 GB memory
❌ **Cold Starts** - May have latency on first request

## Next Steps

1. [Set up monitoring and alerts](./monitoring.md)
2. [Configure custom domain](./custom-domain.md)
3. [Optimize costs](./cost-optimization.md)
4. [Set up staging environment](./staging.md)
