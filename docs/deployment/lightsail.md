# AWS Lightsail Containers Deployment

AWS Lightsail Containers is a simplified container service with predictable, fixed monthly pricing. Perfect for developers who want simplicity without the complexity of ECS.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│              AWS Lightsail Container Architecture                    │
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
                    │  2. Push to Lightsail     │
                    └─────────────┬─────────────┘
                                  │
                                  ▼
        ┌────────────────────────────────────────────────┐
        │   AWS Lightsail Container Service              │
        │                                                 │
        │   ┌─────────────────────────────────────┐     │
        │   │  Built-in Load Balancer             │     │
        │   │  ✓ SSL Certificate (Free)           │     │
        │   │  ✓ HTTPS Endpoint                   │     │
        │   │  ✓ Health Checks                    │     │
        │   └──────────────┬──────────────────────┘     │
        │                  │                             │
        │   ┌──────────────┴──────────────┐             │
        │   │                              │             │
        │   ▼                              ▼             │
        │  ┌────────────┐           ┌────────────┐      │
        │  │ Container  │           │ Container  │      │
        │  │ Instance 1 │           │ Instance 2 │      │
        │  │  Port 8080 │           │  Port 8080 │      │
        │  └────────────┘           └────────────┘      │
        │   (Scale: 1-20)                                │
        └────────────────────┬───────────────────────────┘
                             │
                             ▼
              ┌──────────────────────────────┐
              │  Public HTTPS Endpoint       │
              │  *.lightsail.aws.com         │
              │  OR                          │
              │  custom-domain.com           │
              └──────────────┬───────────────┘
                             │
                             ▼
                   ┌──────────────────┐
                   │  End Users       │
                   └──────────────────┘
```

## Resources Required

### AWS Resources
1. **Lightsail Container Service**
   - All-in-one: containers, load balancer, SSL, storage
   - No separate ECR, ALB, or networking needed

### GitHub Secrets
```bash
AWS_ACCESS_KEY_ID          # IAM user access key
AWS_SECRET_ACCESS_KEY      # IAM user secret key
AWS_REGION                 # e.g., us-east-1
LIGHTSAIL_SERVICE_NAME    # e.g., portfolio-website
```

## Pricing Tiers

| Tier | vCPU | RAM | Price/Month | Best For |
|------|------|-----|-------------|----------|
| Nano | 0.25 | 512 MB | $7 | Development |
| Micro | 0.5 | 1 GB | $10 | Low traffic |
| Small | 1 | 2 GB | $25 | Medium traffic |
| Medium | 2 | 4 GB | $50 | High traffic |
| Large | 4 | 8 GB | $100 | Very high traffic |

**Includes**: Compute, storage, load balancer, SSL certificate, data transfer

## Step-by-Step Setup

### 1. Create IAM User (One-time)

```bash
# Create IAM user for GitHub Actions
aws iam create-user --user-name github-actions-lightsail

# Create inline policy
cat > lightsail-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "lightsail:*"
      ],
      "Resource": "*"
    }
  ]
}
EOF

aws iam put-user-policy \
  --user-name github-actions-lightsail \
  --policy-name LightsailFullAccess \
  --policy-document file://lightsail-policy.json

# Create access key
aws iam create-access-key --user-name github-actions-lightsail
```

### 2. Create Lightsail Container Service

#### Option A: AWS Console (Recommended)

1. Go to [Lightsail Console](https://lightsail.aws.amazon.com/)
2. Click "Containers" → "Create container service"
3. Choose your region
4. Select capacity:
   - **Nano** ($7/month) for development
   - **Small** ($25/month) for production
5. Choose number of nodes (1-3 for high availability)
6. Name your service: `portfolio-website`
7. Click "Create container service"

#### Option B: AWS CLI

```bash
# Create container service
aws lightsail create-container-service \
  --service-name portfolio-website \
  --power nano \
  --scale 1 \
  --region us-east-1

# Wait for service to be ready (takes 2-3 minutes)
aws lightsail get-container-services \
  --service-name portfolio-website
```

### 3. GitHub Actions Workflow

Create `.github/workflows/deploy-lightsail.yml`:

```yaml
name: Deploy to Lightsail

on:
  push:
    branches: [ main ]
  workflow_dispatch: {}

env:
  AWS_REGION: us-east-1
  LIGHTSAIL_SERVICE_NAME: portfolio-website

jobs:
  deploy:
    name: Build and Deploy
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ${{ env.AWS_REGION }}

      - name: Install Lightsail CLI Plugin
        run: |
          curl "https://s3.us-west-2.amazonaws.com/lightsailctl/latest/linux-amd64/lightsailctl" -o "/usr/local/bin/lightsailctl"
          sudo chmod +x /usr/local/bin/lightsailctl

      - name: Build Docker image
        run: docker build -t portfolio:${{ github.sha }} .

      - name: Push image to Lightsail
        run: |
          aws lightsail push-container-image \
            --service-name ${{ env.LIGHTSAIL_SERVICE_NAME }} \
            --label portfolio \
            --image portfolio:${{ github.sha }}

      - name: Get image tag
        id: image
        run: |
          IMAGE_TAG=$(aws lightsail get-container-images --service-name ${{ env.LIGHTSAIL_SERVICE_NAME }} | jq -r '.containerImages[0].image')
          echo "image=$IMAGE_TAG" >> $GITHUB_OUTPUT

      - name: Create deployment
        run: |
          cat > deployment.json <<EOF
          {
            "containers": {
              "app": {
                "image": "${{ steps.image.outputs.image }}",
                "ports": {
                  "8080": "HTTP"
                },
                "environment": {
                  "PORT": "8080"
                }
              }
            },
            "publicEndpoint": {
              "containerName": "app",
              "containerPort": 8080,
              "healthCheck": {
                "path": "/api/health",
                "intervalSeconds": 10
              }
            }
          }
          EOF

          aws lightsail create-container-service-deployment \
            --service-name ${{ env.LIGHTSAIL_SERVICE_NAME }} \
            --cli-input-json file://deployment.json

      - name: Wait for deployment
        run: |
          echo "Waiting for deployment to complete..."
          sleep 30

          # Check deployment status
          aws lightsail get-container-services \
            --service-name ${{ env.LIGHTSAIL_SERVICE_NAME }} \
            --query 'containerServices[0].url' \
            --output text

      - name: Get service URL
        run: |
          URL=$(aws lightsail get-container-services \
            --service-name ${{ env.LIGHTSAIL_SERVICE_NAME }} \
            --query 'containerServices[0].url' \
            --output text)
          echo "🚀 Deployed to: https://$URL"
```

### 4. Configure GitHub Secrets

In GitHub repository settings:

```
AWS_ACCESS_KEY_ID: AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_REGION: us-east-1
LIGHTSAIL_SERVICE_NAME: portfolio-website
```

### 5. Deploy Application

```bash
# Push to main branch
git add .
git commit -m "Deploy to Lightsail"
git push origin main

# GitHub Actions will automatically deploy
```

## Configuration

### Scaling Configuration

```bash
# Scale to 2 instances for high availability
aws lightsail update-container-service \
  --service-name portfolio-website \
  --scale 2
```

### Update Service Power

```bash
# Upgrade to Small tier
aws lightsail update-container-service \
  --service-name portfolio-website \
  --power small
```

### Health Check Configuration

```json
{
  "healthCheck": {
    "path": "/api/health",
    "intervalSeconds": 10,
    "timeoutSeconds": 5,
    "successCodes": "200-299",
    "healthyThreshold": 2,
    "unhealthyThreshold": 2
  }
}
```

### Environment Variables

```json
{
  "containers": {
    "app": {
      "environment": {
        "PORT": "8080",
        "LOG_LEVEL": "info",
        "NODE_ENV": "production"
      }
    }
  }
}
```

## Custom Domain Setup

### 1. Create SSL Certificate

```bash
# Certificate is automatically created when you add a domain
aws lightsail create-certificate \
  --certificate-name portfolio-cert \
  --domain-name portfolio.example.com \
  --subject-alternative-names www.portfolio.example.com
```

### 2. Validate Certificate

Add DNS records provided by Lightsail to your domain registrar.

### 3. Attach Domain

```bash
# Attach custom domain
aws lightsail attach-certificate-to-distribution \
  --distribution-name portfolio-website \
  --certificate-name portfolio-cert
```

### 4. Update DNS

```
Type: A
Name: portfolio.example.com
Value: [Your Lightsail IP from console]

Type: CNAME
Name: www.portfolio.example.com
Value: portfolio.example.com
```

## Monitoring & Logs

### View Logs

```bash
# Get container logs
aws lightsail get-container-log \
  --service-name portfolio-website \
  --container-name app \
  --start-time $(date -u -d '1 hour ago' '+%s') \
  --end-time $(date -u '+%s')

# Follow logs (requires lightsailctl)
lightsailctl logs follow \
  --service-name portfolio-website \
  --container-name app
```

### Metrics

Available in Lightsail console:
- CPU utilization
- Memory utilization
- Request count
- HTTP responses (2xx, 4xx, 5xx)
- Active connections

### Set Up Alarms

```bash
# Create CPU alarm
aws lightsail put-alarm \
  --alarm-name portfolio-high-cpu \
  --metric-name CPUUtilization \
  --monitored-resource-name portfolio-website \
  --comparison-operator GreaterThanThreshold \
  --threshold 80 \
  --evaluation-periods 2
```

## Cost Breakdown

### Monthly Cost (Nano Tier)

| Component | Cost |
|-----------|------|
| Compute (0.25 vCPU, 512 MB) | Included |
| Load Balancer | Included |
| SSL Certificate | Included |
| 500 GB Data Transfer | Included |
| **Total** | **$7.00** |

### Monthly Cost (Small Tier)

| Component | Cost |
|-----------|------|
| Compute (1 vCPU, 2 GB) | Included |
| Load Balancer | Included |
| SSL Certificate | Included |
| 500 GB Data Transfer | Included |
| **Total** | **$25.00** |

### Additional Costs

- **Extra Data Transfer**: $0.09/GB over 500 GB
- **Additional Nodes**: Same rate per node

## Troubleshooting

### Service Won't Start

1. Check deployment status:
```bash
aws lightsail get-container-service-deployments \
  --service-name portfolio-website
```

2. View container logs:
```bash
aws lightsail get-container-log \
  --service-name portfolio-website \
  --container-name app
```

3. Verify health check:
```bash
# Test locally first
docker run -p 8080:8080 your-image:latest
curl http://localhost:8080/api/health
```

### Deployment Fails

```bash
# Check service state
aws lightsail get-container-services \
  --service-name portfolio-website \
  --query 'containerServices[0].state'

# If state is "UPDATING" for too long, check logs
aws lightsail get-container-log \
  --service-name portfolio-website \
  --container-name app
```

### High Memory Usage

```bash
# Upgrade to next tier
aws lightsail update-container-service \
  --service-name portfolio-website \
  --power micro  # or small, medium, etc.
```

## Maintenance

### Update Deployment

GitHub Actions automatically deploys on push to main. Manual deployment:

```bash
# Trigger new deployment
aws lightsail create-container-service-deployment \
  --service-name portfolio-website \
  --cli-input-json file://deployment.json
```

### Disable Service (Pause)

```bash
# Lightsail doesn't support pausing, but you can:
# 1. Scale to 0 (not supported)
# 2. Delete and recreate (loses data)
# 3. Keep running (it's cheap!)
```

### Delete Service

```bash
# Delete container service
aws lightsail delete-container-service \
  --service-name portfolio-website
```

## Comparison with Other Options

### vs App Runner
- ✅ **Simpler**: Easier to understand and set up
- ✅ **Cheaper**: Fixed pricing, starts at $7/month
- ✅ **Predictable**: No surprise bills
- ❌ **Less Flexible**: Can't customize as much
- ❌ **Scale Limits**: Max 20 containers

### vs ECS
- ✅ **Much Simpler**: No VPCs, subnets, security groups
- ✅ **Cheaper**: Includes load balancer in price
- ✅ **Faster Setup**: 10 minutes vs 1-2 hours
- ❌ **Less Control**: Can't customize infrastructure
- ❌ **Limited Regions**: Not available everywhere

## Advantages

✅ **Simple Setup** - Easiest AWS deployment option
✅ **Fixed Pricing** - No surprise costs
✅ **Includes SSL** - Free HTTPS certificates
✅ **Includes LB** - Load balancer built-in
✅ **Quick Start** - Deploy in 10 minutes
✅ **Great Console** - User-friendly interface

## Limitations

❌ **Max 20 Containers** - Limited scaling
❌ **Limited Regions** - Fewer than EC2/ECS
❌ **Less Customization** - Can't configure networking
❌ **No Auto-Pause** - Always running (always paying)

## Best Practices

### 1. Start Small, Scale Up
```bash
# Start with Nano for development
# Upgrade to Small for production
aws lightsail update-container-service \
  --service-name portfolio-website \
  --power small
```

### 2. Use Health Checks
Always configure health checks for reliability:
```json
{
  "healthCheck": {
    "path": "/api/health",
    "intervalSeconds": 10
  }
}
```

### 3. Monitor Regularly
Set up CloudWatch alarms for:
- High CPU usage
- High memory usage
- Error rates

### 4. Use Multiple Nodes
For production, use at least 2 nodes for high availability:
```bash
aws lightsail update-container-service \
  --service-name portfolio-website \
  --scale 2
```

## Next Steps

1. [Set up custom domain](./custom-domain.md)
2. [Configure monitoring](./monitoring.md)
3. [Optimize performance](./performance.md)
4. [Set up staging environment](./staging.md)
