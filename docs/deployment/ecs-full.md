# ECS Fargate Deployment (Full Production Setup)

Complete production-ready deployment with Application Load Balancer, auto-scaling, HTTPS, and high availability.

## Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│                 ECS Fargate Full Production Architecture                  │
└──────────────────────────────────────────────────────────────────────────┘

                            ┌──────────────────┐
                            │  Developer       │
                            │  Pushes Code     │
                            └────────┬─────────┘
                                     │
                                     ▼
                            ┌──────────────────┐
                            │  GitHub Actions  │
                            │  (CI/CD Pipeline)│
                            └────────┬─────────┘
                                     │
                       ┌─────────────┴─────────────┐
                       │  1. Build Docker Image    │
                       │  2. Push to ECR           │
                       │  3. Update ECS Service    │
                       └─────────────┬─────────────┘
                                     │
                                     ▼
                       ┌──────────────────────────┐
                       │   Amazon ECR             │
                       │   (Container Registry)   │
                       └──────────────┬───────────┘
                                      │
┌─────────────────────────────────────┼────────────────────────────────────┐
│                                     │                                     │
│  AWS VPC (Virtual Private Cloud)   │                                     │
│                                     ▼                                     │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                  Application Load Balancer (ALB)                 │    │
│  │                                                                   │    │
│  │  ┌──────────────────┐         ┌────────────────────────────┐    │    │
│  │  │  Listener :80    │────────▶│ Redirect to HTTPS          │    │    │
│  │  │  (HTTP)          │         └────────────────────────────┘    │    │
│  │  └──────────────────┘                                            │    │
│  │                                                                   │    │
│  │  ┌──────────────────┐         ┌────────────────────────────┐    │    │
│  │  │  Listener :443   │────────▶│ Forward to Target Group    │    │    │
│  │  │  (HTTPS + SSL)   │         │ Health Check: /api/health  │    │    │
│  │  └──────────────────┘         └────────────────────────────┘    │    │
│  └──────────────────────────────────┬───────────────────────────────┘    │
│                                     │                                     │
│  ┌──────────────────────────────────┼───────────────────────────────┐    │
│  │         ECS Cluster               ▼                               │    │
│  │  ┌────────────────────────────────────────────────────────────┐  │    │
│  │  │                    ECS Service                              │  │    │
│  │  │  (Auto-scaling: 1-10 tasks)                                │  │    │
│  │  │                                                             │  │    │
│  │  │  ┌─────────────────┐        ┌─────────────────┐            │  │    │
│  │  │  │ Availability    │        │ Availability    │            │  │    │
│  │  │  │ Zone A          │        │ Zone B          │            │  │    │
│  │  │  │                 │        │                 │            │  │    │
│  │  │  │ ┌─────────────┐ │        │ ┌─────────────┐ │            │  │    │
│  │  │  │ │  Task 1     │ │        │ │  Task 2     │ │            │  │    │
│  │  │  │ │  Container  │ │        │ │  Container  │ │            │  │    │
│  │  │  │ │  :8080      │ │        │ │  :8080      │ │            │  │    │
│  │  │  │ └─────────────┘ │        │ └─────────────┘ │            │  │    │
│  │  │  │                 │        │                 │            │  │    │
│  │  │  │ Private Subnet  │        │ Private Subnet  │            │  │    │
│  │  │  └─────────────────┘        └─────────────────┘            │  │    │
│  │  │                                                             │  │    │
│  │  └─────────────────────────────────────────────────────────────┘  │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  Security Groups                                                 │    │
│  │  ┌──────────────────┐         ┌────────────────────────────┐    │    │
│  │  │  ALB SG          │         │  ECS Task SG               │    │    │
│  │  │  In: 80, 443     │────────▶│  In: 8080 from ALB SG only │    │    │
│  │  │  from 0.0.0.0/0  │         │  Out: All                  │    │    │
│  │  └──────────────────┘         └────────────────────────────┘    │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  NAT Gateway (for private subnets to access internet)           │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└──────────────────────────────────────┬───────────────────────────────────┘
                                       │
                                       ▼
                        ┌──────────────────────────────┐
                        │  Route 53 (Optional)         │
                        │  portfolio.example.com       │
                        │  → ALB DNS                   │
                        └──────────────┬───────────────┘
                                       │
                                       ▼
                             ┌──────────────────┐
                             │  End Users       │
                             │  HTTPS Access    │
                             └──────────────────┘
```

## Resources Required

### AWS Resources (12 total)

1. **VPC & Networking**
   - VPC with CIDR block
   - 2 Public Subnets (different AZs)
   - 2 Private Subnets (different AZs)
   - Internet Gateway
   - NAT Gateway (optional, for private subnets)
   - Route Tables

2. **Security Groups**
   - ALB Security Group
   - ECS Task Security Group

3. **Load Balancing**
   - Application Load Balancer
   - Target Group
   - ALB Listeners (HTTP & HTTPS)

4. **Container Service**
   - ECR Repository
   - ECS Cluster
   - ECS Task Definition
   - ECS Service

5. **IAM Roles**
   - Task Execution Role
   - Task Role (optional)

6. **Monitoring**
   - CloudWatch Log Group
   - CloudWatch Alarms (optional)

7. **SSL/TLS (Optional)**
   - ACM Certificate
   - Route 53 Hosted Zone

### GitHub Secrets
```bash
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
AWS_REGION
ECR_REPOSITORY
ECS_CLUSTER
ECS_SERVICE
ECS_TASK_DEFINITION
```

## Quick Setup with CloudFormation

For the fastest setup, use this CloudFormation template:

```bash
# Download template
curl -O https://raw.githubusercontent.com/aws-samples/ecs-refarch-cloudformation/master/master.yaml

# Deploy stack
aws cloudformation create-stack \
  --stack-name portfolio-infrastructure \
  --template-body file://master.yaml \
  --parameters \
    ParameterKey=EnvironmentName,ParameterValue=portfolio \
  --capabilities CAPABILITY_IAM

# Wait for completion (~10 minutes)
aws cloudformation wait stack-create-complete \
  --stack-name portfolio-infrastructure
```

## Manual Step-by-Step Setup

### 1. Create VPC and Subnets

```bash
# Create VPC
VPC_ID=$(aws ec2 create-vpc \
  --cidr-block 10.0.0.0/16 \
  --query 'Vpc.VpcId' \
  --output text)

aws ec2 create-tags \
  --resources $VPC_ID \
  --tags Key=Name,Value=portfolio-vpc

# Enable DNS
aws ec2 modify-vpc-attribute \
  --vpc-id $VPC_ID \
  --enable-dns-hostnames

# Create Internet Gateway
IGW_ID=$(aws ec2 create-internet-gateway \
  --query 'InternetGateway.InternetGatewayId' \
  --output text)

aws ec2 attach-internet-gateway \
  --vpc-id $VPC_ID \
  --internet-gateway-id $IGW_ID

# Create Public Subnets (2 AZs)
PUBLIC_SUBNET_1=$(aws ec2 create-subnet \
  --vpc-id $VPC_ID \
  --cidr-block 10.0.1.0/24 \
  --availability-zone us-east-1a \
  --query 'Subnet.SubnetId' \
  --output text)

PUBLIC_SUBNET_2=$(aws ec2 create-subnet \
  --vpc-id $VPC_ID \
  --cidr-block 10.0.2.0/24 \
  --availability-zone us-east-1b \
  --query 'Subnet.SubnetId' \
  --output text)

# Create Private Subnets (2 AZs)
PRIVATE_SUBNET_1=$(aws ec2 create-subnet \
  --vpc-id $VPC_ID \
  --cidr-block 10.0.10.0/24 \
  --availability-zone us-east-1a \
  --query 'Subnet.SubnetId' \
  --output text)

PRIVATE_SUBNET_2=$(aws ec2 create-subnet \
  --vpc-id $VPC_ID \
  --cidr-block 10.0.11.0/24 \
  --availability-zone us-east-1b \
  --query 'Subnet.SubnetId' \
  --output text)

# Create Route Table for Public Subnets
PUBLIC_RT=$(aws ec2 create-route-table \
  --vpc-id $VPC_ID \
  --query 'RouteTable.RouteTableId' \
  --output text)

aws ec2 create-route \
  --route-table-id $PUBLIC_RT \
  --destination-cidr-block 0.0.0.0/0 \
  --gateway-id $IGW_ID

aws ec2 associate-route-table \
  --subnet-id $PUBLIC_SUBNET_1 \
  --route-table-id $PUBLIC_RT

aws ec2 associate-route-table \
  --subnet-id $PUBLIC_SUBNET_2 \
  --route-table-id $PUBLIC_RT
```

### 2. Create Security Groups

```bash
# ALB Security Group
ALB_SG=$(aws ec2 create-security-group \
  --group-name portfolio-alb-sg \
  --description "Security group for ALB" \
  --vpc-id $VPC_ID \
  --query 'GroupId' \
  --output text)

# Allow HTTP
aws ec2 authorize-security-group-ingress \
  --group-id $ALB_SG \
  --protocol tcp \
  --port 80 \
  --cidr 0.0.0.0/0

# Allow HTTPS
aws ec2 authorize-security-group-ingress \
  --group-id $ALB_SG \
  --protocol tcp \
  --port 443 \
  --cidr 0.0.0.0/0

# ECS Task Security Group
TASK_SG=$(aws ec2 create-security-group \
  --group-name portfolio-task-sg \
  --description "Security group for ECS tasks" \
  --vpc-id $VPC_ID \
  --query 'GroupId' \
  --output text)

# Allow traffic from ALB only
aws ec2 authorize-security-group-ingress \
  --group-id $TASK_SG \
  --protocol tcp \
  --port 8080 \
  --source-group $ALB_SG
```

### 3. Create Application Load Balancer

```bash
# Create ALB
ALB_ARN=$(aws elbv2 create-load-balancer \
  --name portfolio-alb \
  --subnets $PUBLIC_SUBNET_1 $PUBLIC_SUBNET_2 \
  --security-groups $ALB_SG \
  --scheme internet-facing \
  --type application \
  --ip-address-type ipv4 \
  --query 'LoadBalancers[0].LoadBalancerArn' \
  --output text)

# Get ALB DNS name
ALB_DNS=$(aws elbv2 describe-load-balancers \
  --load-balancer-arns $ALB_ARN \
  --query 'LoadBalancers[0].DNSName' \
  --output text)

echo "ALB DNS: $ALB_DNS"

# Create Target Group
TG_ARN=$(aws elbv2 create-target-group \
  --name portfolio-tg \
  --protocol HTTP \
  --port 8080 \
  --vpc-id $VPC_ID \
  --target-type ip \
  --health-check-path /api/health \
  --health-check-interval-seconds 30 \
  --health-check-timeout-seconds 5 \
  --healthy-threshold-count 2 \
  --unhealthy-threshold-count 3 \
  --query 'TargetGroups[0].TargetGroupArn' \
  --output text)

# Create HTTP Listener (redirects to HTTPS)
aws elbv2 create-listener \
  --load-balancer-arn $ALB_ARN \
  --protocol HTTP \
  --port 80 \
  --default-actions Type=redirect,RedirectConfig='{Protocol=HTTPS,Port=443,StatusCode=HTTP_301}'

# Create HTTPS Listener (if you have a certificate)
# CERT_ARN="arn:aws:acm:us-east-1:123456789:certificate/abc-123"
# aws elbv2 create-listener \
#   --load-balancer-arn $ALB_ARN \
#   --protocol HTTPS \
#   --port 443 \
#   --certificates CertificateArn=$CERT_ARN \
#   --default-actions Type=forward,TargetGroupArn=$TG_ARN
```

### 4. Create ECS Resources

```bash
# Create cluster
aws ecs create-cluster --cluster-name portfolio-cluster

# Create task definition (see previous section)
# Register task definition
aws ecs register-task-definition --cli-input-json file://task-definition.json

# Create service
aws ecs create-service \
  --cluster portfolio-cluster \
  --service-name portfolio-service \
  --task-definition portfolio-task \
  --desired-count 2 \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[$PRIVATE_SUBNET_1,$PRIVATE_SUBNET_2],securityGroups=[$TASK_SG],assignPublicIp=DISABLED}" \
  --load-balancers "targetGroupArn=$TG_ARN,containerName=portfolio,containerPort=8080"
```

## Cost Breakdown

### Monthly Cost (2 Tasks Always Running)

| Component | Specification | Rate | Monthly Cost |
|-----------|--------------|------|--------------|
| Fargate vCPU | 0.5 vCPU × 2 tasks | $0.04048/vCPU-hour | $59.10 |
| Fargate Memory | 1 GB × 2 tasks | $0.004445/GB-hour | $6.49 |
| Application Load Balancer | Always on | $16.20/month + $0.008/LCU | $16.20 |
| Data Transfer | 20 GB | $0.09/GB | $1.80 |
| NAT Gateway (Optional) | Always on | $0.045/hour | $32.85 |
| CloudWatch Logs | 2 GB | $0.50/GB | $1.00 |
| **Total (no NAT)** | | | **~$84.59** |
| **Total (with NAT)** | | | **~$117.44** |

### Cost Optimization Tips

1. **Remove NAT Gateway**: Use public subnets for tasks
2. **Use Spot pricing**: Save up to 70% on compute
3. **Auto-scale**: Scale down during off-hours
4. **Reduce task size**: Use 0.25 vCPU, 512 MB

## Auto-Scaling Configuration

```bash
# Register scalable target
aws application-autoscaling register-scalable-target \
  --service-namespace ecs \
  --resource-id service/portfolio-cluster/portfolio-service \
  --scalable-dimension ecs:service:DesiredCount \
  --min-capacity 1 \
  --max-capacity 10

# CPU-based scaling
aws application-autoscaling put-scaling-policy \
  --service-namespace ecs \
  --resource-id service/portfolio-cluster/portfolio-service \
  --scalable-dimension ecs:service:DesiredCount \
  --policy-name cpu-scaling \
  --policy-type TargetTrackingScaling \
  --target-tracking-scaling-policy-configuration file://scaling-policy.json
```

**scaling-policy.json:**
```json
{
  "TargetValue": 70.0,
  "PredefinedMetricSpecification": {
    "PredefinedMetricType": "ECSServiceAverageCPUUtilization"
  },
  "ScaleInCooldown": 300,
  "ScaleOutCooldown": 60
}
```

## SSL/TLS Certificate Setup

### 1. Request Certificate (ACM)

```bash
# Request certificate
CERT_ARN=$(aws acm request-certificate \
  --domain-name portfolio.example.com \
  --subject-alternative-names www.portfolio.example.com \
  --validation-method DNS \
  --query 'CertificateArn' \
  --output text)

# Get validation records
aws acm describe-certificate \
  --certificate-arn $CERT_ARN \
  --query 'Certificate.DomainValidationOptions'
```

### 2. Add DNS Records for Validation

Add the CNAME records provided to your DNS provider.

### 3. Wait for Validation

```bash
aws acm wait certificate-validated \
  --certificate-arn $CERT_ARN
```

## Monitoring & Alerts

### Create CloudWatch Alarms

```bash
# High CPU alarm
aws cloudwatch put-metric-alarm \
  --alarm-name portfolio-high-cpu \
  --alarm-description "Alert on high CPU" \
  --metric-name CPUUtilization \
  --namespace AWS/ECS \
  --statistic Average \
  --period 300 \
  --threshold 80 \
  --comparison-operator GreaterThanThreshold \
  --dimensions Name=ServiceName,Value=portfolio-service Name=ClusterName,Value=portfolio-cluster \
  --evaluation-periods 2

# Target unhealthy alarm
aws cloudwatch put-metric-alarm \
  --alarm-name portfolio-unhealthy-targets \
  --metric-name UnHealthyHostCount \
  --namespace AWS/ApplicationELB \
  --statistic Average \
  --period 60 \
  --threshold 1 \
  --comparison-operator GreaterThanOrEqualToThreshold \
  --dimensions Name=TargetGroup,Value=targetgroup/portfolio-tg/abc123 Name=LoadBalancer,Value=app/portfolio-alb/xyz789 \
  --evaluation-periods 2
```

## Advantages

✅ **High Availability** - Multi-AZ deployment
✅ **Auto-Scaling** - Handles traffic spikes
✅ **Load Balancing** - Distributes traffic
✅ **HTTPS Support** - Built-in SSL/TLS
✅ **Production Ready** - Enterprise-grade setup
✅ **Health Checks** - Automatic recovery
✅ **Zero Downtime** - Rolling deployments

## Limitations

❌ **High Cost** - Most expensive option
❌ **Complex Setup** - Many resources to manage
❌ **Over-Engineering** - Overkill for low traffic
❌ **Maintenance** - More components to monitor

## When to Use This Option

Use full ECS setup when you need:
- ✅ Production-grade reliability
- ✅ High availability across multiple AZs
- ✅ Auto-scaling for variable traffic
- ✅ HTTPS with custom domain
- ✅ Enterprise compliance requirements
- ✅ 24/7 uptime requirements

## Next Steps

1. [Configure auto-scaling policies](./autoscaling.md)
2. [Set up CI/CD pipeline](./cicd.md)
3. [Configure monitoring and alerts](./monitoring.md)
4. [Implement blue-green deployments](./blue-green.md)
