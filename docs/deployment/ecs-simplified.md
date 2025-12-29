# ECS Fargate Deployment (Simplified - No Load Balancer)

Deploy to ECS Fargate without an Application Load Balancer, using a public IP address. This is a cost-effective option that provides more control than App Runner while being cheaper than the full ECS setup.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│            ECS Fargate Simplified Architecture                       │
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
                                   │ Pull image
                                   │
                                   ▼
        ┌──────────────────────────────────────────────┐
        │          AWS VPC (Virtual Network)           │
        │                                               │
        │   ┌──────────────────────────────────────┐   │
        │   │      Public Subnet (us-east-1a)      │   │
        │   │                                       │   │
        │   │  ┌────────────────────────────────┐  │   │
        │   │  │    ECS Fargate Task            │  │   │
        │   │  │  ┌──────────────────────────┐  │  │   │
        │   │  │  │  Container               │  │  │   │
        │   │  │  │  Port 8080               │  │  │   │
        │   │  │  │  (Your Go App)           │  │  │   │
        │   │  │  └──────────────────────────┘  │  │   │
        │   │  │                                 │  │   │
        │   │  │  Public IP: 54.123.45.67       │  │   │
        │   │  └────────────────────────────────┘  │   │
        │   │                                       │   │
        │   └───────────────────┬───────────────────┘   │
        │                       │                        │
        │   ┌───────────────────▼───────────────────┐   │
        │   │      Security Group                   │   │
        │   │  Inbound: Port 8080 from 0.0.0.0/0    │   │
        │   │  Outbound: All traffic                │   │
        │   └───────────────────────────────────────┘   │
        │                                               │
        │   ┌───────────────────────────────────────┐   │
        │   │      Internet Gateway                 │   │
        │   └───────────────────────────────────────┘   │
        └───────────────────┬───────────────────────────┘
                            │
                            ▼
                  ┌──────────────────────┐
                  │  Public Access       │
                  │  http://54.123.45.67:8080 │
                  └──────────────────────┘
                            │
                            ▼
                  ┌──────────────────┐
                  │  End Users       │
                  └──────────────────┘
```

## Resources Required

### AWS Resources
1. **Amazon ECR Repository** - Store Docker images
2. **VPC with Public Subnet** - Network infrastructure (use default VPC)
3. **Security Group** - Allow traffic on port 8080
4. **ECS Cluster** - Container orchestration
5. **ECS Task Definition** - Container specifications
6. **ECS Service** - Runs and maintains tasks
7. **IAM Task Execution Role** - Permissions for ECS

### GitHub Secrets
```bash
AWS_ACCESS_KEY_ID          # IAM user access key
AWS_SECRET_ACCESS_KEY      # IAM user secret key
AWS_REGION                 # e.g., us-east-1
ECR_REPOSITORY            # e.g., portfolio-website
ECS_CLUSTER               # e.g., portfolio-cluster
ECS_SERVICE               # e.g., portfolio-service
```

## Step-by-Step Setup

### 1. Create IAM User (One-time)

```bash
# Create IAM user
aws iam create-user --user-name github-actions-ecs

# Attach policies
aws iam attach-user-policy \
  --user-name github-actions-ecs \
  --policy-arn arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser

aws iam attach-user-policy \
  --user-name github-actions-ecs \
  --policy-arn arn:aws:iam::aws:policy/AmazonECS_FullAccess

# Create access key
aws iam create-access-key --user-name github-actions-ecs
```

### 2. Create ECR Repository

```bash
# Create repository
aws ecr create-repository \
  --repository-name portfolio-website \
  --region us-east-1
```

### 3. Create Security Group

```bash
# Get default VPC ID
VPC_ID=$(aws ec2 describe-vpcs \
  --filters "Name=isDefault,Values=true" \
  --query 'Vpcs[0].VpcId' \
  --output text)

# Create security group
SG_ID=$(aws ec2 create-security-group \
  --group-name portfolio-sg \
  --description "Security group for portfolio website" \
  --vpc-id $VPC_ID \
  --query 'GroupId' \
  --output text)

# Allow inbound traffic on port 8080
aws ec2 authorize-security-group-ingress \
  --group-id $SG_ID \
  --protocol tcp \
  --port 8080 \
  --cidr 0.0.0.0/0

echo "Security Group ID: $SG_ID"
```

### 4. Create IAM Task Execution Role

```bash
# Create trust policy
cat > task-execution-trust-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF

# Create role
aws iam create-role \
  --role-name ecsTaskExecutionRole \
  --assume-role-policy-document file://task-execution-trust-policy.json

# Attach policy
aws iam attach-role-policy \
  --role-name ecsTaskExecutionRole \
  --policy-arn arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy
```

### 5. Create ECS Cluster

```bash
# Create Fargate cluster
aws ecs create-cluster \
  --cluster-name portfolio-cluster \
  --capacity-providers FARGATE \
  --default-capacity-provider-strategy capacityProvider=FARGATE,weight=1
```

### 6. Create Task Definition

```bash
# Get your account ID
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

# Get execution role ARN
EXECUTION_ROLE_ARN="arn:aws:iam::${ACCOUNT_ID}:role/ecsTaskExecutionRole"

# Create task definition
cat > task-definition.json <<EOF
{
  "family": "portfolio-task",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "$EXECUTION_ROLE_ARN",
  "containerDefinitions": [
    {
      "name": "portfolio",
      "image": "${ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com/portfolio-website:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "essential": true,
      "environment": [
        {
          "name": "PORT",
          "value": "8080"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/portfolio",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
EOF

# Register task definition
aws ecs register-task-definition \
  --cli-input-json file://task-definition.json
```

### 7. Create CloudWatch Log Group

```bash
# Create log group
aws logs create-log-group \
  --log-group-name /ecs/portfolio
```

### 8. Create ECS Service

```bash
# Get public subnet IDs
SUBNET_IDS=$(aws ec2 describe-subnets \
  --filters "Name=vpc-id,Values=$VPC_ID" "Name=map-public-ip-on-launch,Values=true" \
  --query 'Subnets[0:2].SubnetId' \
  --output text | tr '\t' ',')

# Create service
aws ecs create-service \
  --cluster portfolio-cluster \
  --service-name portfolio-service \
  --task-definition portfolio-task \
  --desired-count 1 \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[$SUBNET_IDS],securityGroups=[$SG_ID],assignPublicIp=ENABLED}"
```

### 9. GitHub Actions Workflow

Update `.github/workflows/deploy-to-aws.yml`:

```yaml
name: Deploy to ECS (Simplified)

on:
  push:
    branches: [ main ]
  workflow_dispatch: {}

env:
  AWS_REGION: ${{ secrets.AWS_REGION }}
  ECR_REPOSITORY: ${{ secrets.ECR_REPOSITORY }}
  ECS_CLUSTER: ${{ secrets.ECS_CLUSTER }}
  ECS_SERVICE: ${{ secrets.ECS_SERVICE }}

jobs:
  deploy:
    name: Deploy to ECS
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
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

      - name: Force new ECS deployment
        run: |
          aws ecs update-service \
            --cluster ${{ env.ECS_CLUSTER }} \
            --service ${{ env.ECS_SERVICE }} \
            --force-new-deployment

      - name: Wait for deployment to complete
        run: |
          aws ecs wait services-stable \
            --cluster ${{ env.ECS_CLUSTER }} \
            --services ${{ env.ECS_SERVICE }}

      - name: Get public IP
        run: |
          TASK_ARN=$(aws ecs list-tasks \
            --cluster ${{ env.ECS_CLUSTER }} \
            --service-name ${{ env.ECS_SERVICE }} \
            --query 'taskArns[0]' \
            --output text)

          ENI_ID=$(aws ecs describe-tasks \
            --cluster ${{ env.ECS_CLUSTER }} \
            --tasks $TASK_ARN \
            --query 'tasks[0].attachments[0].details[?name==`networkInterfaceId`].value' \
            --output text)

          PUBLIC_IP=$(aws ec2 describe-network-interfaces \
            --network-interface-ids $ENI_ID \
            --query 'NetworkInterfaces[0].Association.PublicIp' \
            --output text)

          echo "🚀 Deployed to: http://$PUBLIC_IP:8080"
```

### 10. Configure GitHub Secrets

Add to GitHub repository settings:

```
AWS_ACCESS_KEY_ID: AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_REGION: us-east-1
ECR_REPOSITORY: portfolio-website
ECS_CLUSTER: portfolio-cluster
ECS_SERVICE: portfolio-service
```

## Cost Breakdown

### Monthly Cost Estimate

| Component | Specification | Rate | Monthly Cost |
|-----------|--------------|------|--------------|
| Fargate (vCPU) | 0.25 vCPU | $0.04048/hour | $29.55 |
| Fargate (Memory) | 512 MB | $0.004445/GB-hour | $1.62 |
| Data Transfer | 10 GB | $0.09/GB | $0.90 |
| CloudWatch Logs | 1 GB | $0.50/GB | $0.50 |
| **Total** | | | **~$32.57** |

### Cost Optimization

```bash
# Use smaller task size
# Update task definition with:
"cpu": "256",
"memory": "512"

# Stop service when not needed
aws ecs update-service \
  --cluster portfolio-cluster \
  --service portfolio-service \
  --desired-count 0
```

## Getting the Public IP

### After Deployment

```bash
# Get task ARN
TASK_ARN=$(aws ecs list-tasks \
  --cluster portfolio-cluster \
  --service-name portfolio-service \
  --query 'taskArns[0]' \
  --output text)

# Get network interface ID
ENI_ID=$(aws ecs describe-tasks \
  --cluster portfolio-cluster \
  --tasks $TASK_ARN \
  --query 'tasks[0].attachments[0].details[?name==`networkInterfaceId`].value' \
  --output text)

# Get public IP
PUBLIC_IP=$(aws ec2 describe-network-interfaces \
  --network-interface-ids $ENI_ID \
  --query 'NetworkInterfaces[0].Association.PublicIp' \
  --output text)

echo "Access your app at: http://$PUBLIC_IP:8080"
```

### Save IP for Easy Access

```bash
# Create parameter in SSM for easy retrieval
aws ssm put-parameter \
  --name /portfolio/public-ip \
  --value $PUBLIC_IP \
  --type String \
  --overwrite

# Retrieve later
aws ssm get-parameter \
  --name /portfolio/public-ip \
  --query 'Parameter.Value' \
  --output text
```

## Monitoring & Logs

### View Logs

```bash
# View recent logs
aws logs tail /ecs/portfolio --follow

# Filter errors
aws logs filter-log-events \
  --log-group-name /ecs/portfolio \
  --filter-pattern "ERROR"
```

### Check Service Status

```bash
# Check service health
aws ecs describe-services \
  --cluster portfolio-cluster \
  --services portfolio-service \
  --query 'services[0].{Status:status,Running:runningCount,Desired:desiredCount}'
```

### CloudWatch Metrics

Available metrics:
- `CPUUtilization`
- `MemoryUtilization`
- `NetworkRxBytes`
- `NetworkTxBytes`

## Troubleshooting

### Task Keeps Stopping

```bash
# Check stopped tasks
aws ecs describe-tasks \
  --cluster portfolio-cluster \
  --tasks $(aws ecs list-tasks --cluster portfolio-cluster --desired-status STOPPED --query 'taskArns[0]' --output text) \
  --query 'tasks[0].stoppedReason'

# Common issues:
# - Image not found in ECR
# - Health check failing
# - Not enough memory
```

### Can't Access via Public IP

```bash
# Verify security group
aws ec2 describe-security-groups \
  --group-ids $SG_ID

# Check if port 8080 is open
# Should see: IpProtocol: tcp, FromPort: 8080, ToPort: 8080
```

### Service Won't Deploy

```bash
# Check service events
aws ecs describe-services \
  --cluster portfolio-cluster \
  --services portfolio-service \
  --query 'services[0].events[:5]'
```

## Limitations

### IP Changes on Restart
- ❌ Public IP changes when task restarts
- Solution: Use Elastic IP (extra cost) or ALB

### No HTTPS
- ❌ No built-in SSL/TLS
- Solution: Use CloudFront or ALB

### Single Instance
- ❌ Only one task running
- ❌ No automatic failover
- Solution: Use ALB with multiple tasks

### Port 8080 in URL
- ❌ Must include `:8080` in URL
- Solution: Use ALB on port 80/443

## Advantages

✅ **Lower Cost** - No ALB fees (~$16/month saved)
✅ **More Control** - Full ECS configuration
✅ **Good for Learning** - Understand ECS fundamentals
✅ **Quick Setup** - Faster than full ECS

## When to Upgrade to Full ECS

Consider full ECS with ALB when:
- You need HTTPS
- You need a custom domain
- You need high availability (multiple tasks)
- You need auto-scaling
- Traffic increases significantly

## Next Steps

1. [Add Elastic IP for stable address](./elastic-ip.md)
2. [Set up CloudFront for HTTPS](./cloudfront.md)
3. [Upgrade to full ECS with ALB](./ecs-full.md)
4. [Configure auto-scaling](./autoscaling.md)
