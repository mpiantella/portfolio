# Implementation Guide - Data Ingestion Pipeline

This guide provides step-by-step instructions for implementing and deploying the data ingestion pipeline.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Local Development Setup](#local-development-setup)
3. [Building the Application](#building-the-application)
4. [Infrastructure Deployment](#infrastructure-deployment)
5. [Lambda Deployment](#lambda-deployment)
6. [Testing](#testing)
7. [Monitoring](#monitoring)
8. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Tools

```bash
# Go 1.21 or higher
go version

# AWS CLI configured
aws --version
aws configure list

# Terraform or AWS CDK
terraform version
# or
cdk --version

# Docker (for local testing)
docker --version

# Make (optional but recommended)
make --version
```

### AWS Account Setup

1. **Create AWS Account** (if not already exists)
2. **Configure IAM User** with appropriate permissions
3. **Set up AWS CLI** with credentials:

```bash
aws configure
# AWS Access Key ID: [your-access-key]
# AWS Secret Access Key: [your-secret-key]
# Default region name: us-east-1
# Default output format: json
```

---

## Local Development Setup

### 1. Clone Repository

```bash
git clone https://github.com/your-org/ingestion-pipeline.git
cd ingestion-pipeline
```

### 2. Install Go Dependencies

```bash
# From the project root
go mod download

# Verify dependencies
go mod verify
```

### 3. Set Up Local PostgreSQL

```bash
# Using Docker
docker run --name postgres-dev \
  -e POSTGRES_PASSWORD=devpassword \
  -e POSTGRES_DB=ingestion \
  -p 5432:5432 \
  -d postgres:15

# Run migrations
export DATABASE_URL="postgres://postgres:devpassword@localhost:5432/ingestion?sslmode=disable"
goose -dir migrations postgres "$DATABASE_URL" up
```

### 4. Set Up Local S3 (using LocalStack)

```bash
# Start LocalStack
docker run --name localstack-dev \
  -p 4566:4566 \
  -e SERVICES=s3,dynamodb,secretsmanager \
  -d localstack/localstack

# Create S3 buckets
aws --endpoint-url=http://localhost:4566 s3 mb s3://ingestion-landing-local
aws --endpoint-url=http://localhost:4566 s3 mb s3://ingestion-processing-local
aws --endpoint-url=http://localhost:4566 s3 mb s3://ingestion-archive-local
```

### 5. Environment Configuration

Create `.env.local` file:

```bash
# Database
DATABASE_URL=postgres://postgres:devpassword@localhost:5432/ingestion?sslmode=disable

# AWS (LocalStack)
AWS_ENDPOINT_URL=http://localhost:4566
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test

# S3 Buckets
S3_LANDING_BUCKET=ingestion-landing-local
S3_PROCESSING_BUCKET=ingestion-processing-local
S3_ARCHIVE_BUCKET=ingestion-archive-local

# Logging
LOG_LEVEL=debug

# Quality Scoring Weights
QUALITY_WEIGHT_COMPLETENESS=20
QUALITY_WEIGHT_ACCURACY=25
QUALITY_WEIGHT_CONSISTENCY=20
QUALITY_WEIGHT_TIMELINESS=10
QUALITY_WEIGHT_UNIQUENESS=15
QUALITY_WEIGHT_VALIDITY=10
```

Load environment:

```bash
source .env.local
# or
export $(cat .env.local | xargs)
```

---

## Building the Application

### Build All Lambda Functions

```bash
# Create build directory
mkdir -p build/lambda

# Build for AWS Lambda (Linux ARM64)
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o build/lambda/parser ./cmd/lambda/parser
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o build/lambda/normalizer ./cmd/lambda/normalizer
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o build/lambda/validator ./cmd/lambda/validator
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o build/lambda/scorer ./cmd/lambda/scorer
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o build/lambda/loader ./cmd/lambda/loader

# Create deployment packages
cd build/lambda
zip parser.zip parser
zip normalizer.zip normalizer
zip validator.zip validator
zip scorer.zip scorer
zip loader.zip loader
cd ../..
```

### Using Makefile (Recommended)

Create `Makefile`:

```makefile
.PHONY: build test clean deploy

# Build all Lambda functions
build:
	@echo "Building Lambda functions..."
	mkdir -p build/lambda
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o build/lambda/parser ./cmd/lambda/parser
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o build/lambda/normalizer ./cmd/lambda/normalizer
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o build/lambda/validator ./cmd/lambda/validator
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o build/lambda/scorer ./cmd/lambda/scorer
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o build/lambda/loader ./cmd/lambda/loader
	cd build/lambda && \
		zip parser.zip parser && \
		zip normalizer.zip normalizer && \
		zip validator.zip validator && \
		zip scorer.zip scorer && \
		zip loader.zip loader

# Run tests
test:
	go test -v -cover ./...

# Clean build artifacts
clean:
	rm -rf build/

# Deploy to AWS
deploy: build
	terraform apply -auto-approve
```

Usage:

```bash
make build   # Build Lambda functions
make test    # Run tests
make deploy  # Deploy to AWS
```

---

## Infrastructure Deployment

### Using Terraform

#### 1. Initialize Terraform

Create `terraform/main.tf`:

```hcl
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  
  backend "s3" {
    bucket = "your-terraform-state-bucket"
    key    = "ingestion-pipeline/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = var.aws_region
}

# VPC Module
module "vpc" {
  source = "./modules/vpc"
  
  vpc_cidr = "10.0.0.0/16"
  azs      = ["us-east-1a", "us-east-1b"]
}

# RDS Module
module "rds" {
  source = "./modules/rds"
  
  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.vpc.private_subnet_ids
  instance_class  = "db.r6g.xlarge"
  allocated_storage = 500
}

# Lambda Module
module "lambda" {
  source = "./modules/lambda"
  
  vpc_id         = module.vpc.vpc_id
  subnet_ids     = module.vpc.private_subnet_ids
  db_endpoint    = module.rds.endpoint
  s3_buckets     = module.s3.bucket_names
}

# S3 Module
module "s3" {
  source = "./modules/s3"
  
  bucket_names = [
    "ingestion-landing-${var.environment}",
    "ingestion-processing-${var.environment}",
    "ingestion-archive-${var.environment}"
  ]
}

# Step Functions Module
module "step_functions" {
  source = "./modules/step_functions"
  
  lambda_arns = module.lambda.function_arns
}
```

#### 2. Deploy Infrastructure

```bash
cd terraform

# Initialize
terraform init

# Plan deployment
terraform plan -out=tfplan

# Apply changes
terraform apply tfplan
```

### Using AWS CDK

#### 1. Initialize CDK Project

```bash
mkdir cdk && cd cdk
cdk init app --language=typescript
```

#### 2. Define Stacks

```typescript
// lib/ingestion-pipeline-stack.ts
import * as cdk from 'aws-cdk-lib';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as rds from 'aws-cdk-lib/aws-rds';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as sfn from 'aws-cdk-lib/aws-stepfunctions';

export class IngestionPipelineStack extends cdk.Stack {
  constructor(scope: cdk.App, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // VPC
    const vpc = new ec2.Vpc(this, 'VPC', {
      maxAzs: 2,
      natGateways: 1,
    });

    // S3 Buckets
    const landingBucket = new s3.Bucket(this, 'LandingBucket', {
      encryption: s3.BucketEncryption.KMS,
      versioned: true,
    });

    // RDS PostgreSQL
    const database = new rds.DatabaseInstance(this, 'Database', {
      engine: rds.DatabaseInstanceEngine.postgres({
        version: rds.PostgresEngineVersion.VER_15,
      }),
      instanceType: ec2.InstanceType.of(
        ec2.InstanceClass.R6G,
        ec2.InstanceSize.XLARGE
      ),
      vpc,
      multiAz: true,
    });

    // Lambda Functions
    const parserFunction = new lambda.Function(this, 'ParserFunction', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      handler: 'bootstrap',
      code: lambda.Code.fromAsset('../build/lambda/parser.zip'),
      vpc,
      memorySize: 3008,
      timeout: cdk.Duration.minutes(15),
    });

    // Step Functions State Machine
    // ... define state machine
  }
}
```

#### 3. Deploy CDK

```bash
# Bootstrap CDK (first time only)
cdk bootstrap

# Deploy stack
cdk deploy IngestionPipelineStack
```

---

## Lambda Deployment

### Deploy Specific Function

```bash
# Update Lambda function code
aws lambda update-function-code \
  --function-name ingestion-parser \
  --zip-file fileb://build/lambda/parser.zip

# Wait for update to complete
aws lambda wait function-updated \
  --function-name ingestion-parser

# Update configuration if needed
aws lambda update-function-configuration \
  --function-name ingestion-parser \
  --memory-size 3008 \
  --timeout 900
```

### Deploy All Functions

```bash
# Using script
./scripts/deploy-lambdas.sh

# Or using AWS SAM
sam deploy --template-file template.yaml \
  --stack-name ingestion-pipeline \
  --capabilities CAPABILITY_IAM
```

---

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Integration Tests

```bash
# Run integration tests (requires LocalStack)
go test -tags=integration ./test/integration/...
```

### End-to-End Tests

```bash
# Upload test file
aws s3 cp test-data/sample.xlsx s3://ingestion-landing-prod/test/

# Monitor Step Functions execution
aws stepfunctions list-executions \
  --state-machine-arn arn:aws:states:us-east-1:ACCOUNT:stateMachine:ingestion-pipeline \
  --max-items 1

# Get execution details
aws stepfunctions describe-execution \
  --execution-arn <execution-arn>
```

---

## Monitoring

### CloudWatch Dashboards

```bash
# Create custom dashboard
aws cloudwatch put-dashboard \
  --dashboard-name ingestion-pipeline \
  --dashboard-body file://cloudwatch-dashboard.json
```

### Alarms

```bash
# Create alarm for Lambda errors
aws cloudwatch put-metric-alarm \
  --alarm-name ingestion-parser-errors \
  --alarm-description "Alert on parser errors" \
  --metric-name Errors \
  --namespace AWS/Lambda \
  --statistic Sum \
  --period 300 \
  --threshold 10 \
  --comparison-operator GreaterThanThreshold \
  --dimensions Name=FunctionName,Value=ingestion-parser \
  --evaluation-periods 1 \
  --alarm-actions arn:aws:sns:us-east-1:ACCOUNT:alerts
```

### Logs

```bash
# Tail Lambda logs
aws logs tail /aws/lambda/ingestion-parser --follow

# Query logs
aws logs filter-log-events \
  --log-group-name /aws/lambda/ingestion-parser \
  --filter-pattern "ERROR" \
  --start-time $(date -u -d '1 hour ago' +%s)000
```

---

## Troubleshooting

### Common Issues

#### 1. Lambda Timeout

**Symptom**: Lambda function times out before completing

**Solution**:
```bash
# Increase timeout
aws lambda update-function-configuration \
  --function-name ingestion-parser \
  --timeout 900  # 15 minutes
```

#### 2. Database Connection Limit

**Symptom**: "too many connections" error

**Solution**:
- Implement connection pooling
- Use RDS Proxy
- Increase max_connections parameter

#### 3. Out of Memory

**Symptom**: Lambda OOM errors

**Solution**:
```bash
# Increase memory allocation
aws lambda update-function-configuration \
  --function-name ingestion-parser \
  --memory-size 10240  # 10GB
```

#### 4. S3 Access Denied

**Symptom**: "Access Denied" when reading S3

**Solution**:
- Check Lambda execution role has S3 permissions
- Verify bucket policy
- Check VPC endpoint configuration

### Debug Mode

Enable debug logging:

```bash
# Update Lambda environment variables
aws lambda update-function-configuration \
  --function-name ingestion-parser \
  --environment Variables={LOG_LEVEL=debug}
```

---

## Scalability Considerations

### Lambda Concurrency

```bash
# Set reserved concurrency
aws lambda put-function-concurrency \
  --function-name ingestion-parser \
  --reserved-concurrent-executions 100

# Set provisioned concurrency
aws lambda put-provisioned-concurrency-config \
  --function-name ingestion-parser \
  --provisioned-concurrent-executions 10 \
  --qualifier $LATEST
```

### RDS Scaling

```bash
# Modify instance class
aws rds modify-db-instance \
  --db-instance-identifier ingestion-db \
  --db-instance-class db.r6g.2xlarge \
  --apply-immediately
```

---

## Security Checklist

- [ ] All S3 buckets have encryption enabled
- [ ] RDS has encryption at rest enabled
- [ ] Secrets stored in AWS Secrets Manager
- [ ] Lambda functions in VPC with private subnets
- [ ] Security groups follow least privilege
- [ ] IAM roles have minimal permissions
- [ ] CloudTrail logging enabled
- [ ] VPC Flow Logs enabled
- [ ] GuardDuty enabled
- [ ] Security Hub enabled

---

## Maintenance

### Database Backups

```bash
# Create manual snapshot
aws rds create-db-snapshot \
  --db-instance-identifier ingestion-db \
  --db-snapshot-identifier ingestion-db-$(date +%Y%m%d)
```

### Cost Optimization

```bash
# Check S3 storage costs
aws s3api list-objects-v2 \
  --bucket ingestion-archive-prod \
  --query 'sum(Contents[].Size)' \
  --output text | awk '{print $1/1024/1024/1024 " GB"}'

# Apply lifecycle policy
aws s3api put-bucket-lifecycle-configuration \
  --bucket ingestion-archive-prod \
  --lifecycle-configuration file://s3-lifecycle.json
```

---

## Additional Resources

- [AWS Lambda Best Practices](https://docs.aws.amazon.com/lambda/latest/dg/best-practices.html)
- [PostgreSQL Performance Tuning](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [Go Lambda Optimization](https://aws.amazon.com/blogs/compute/optimizing-aws-lambda-function-performance-for-go/)
- [Clean Architecture in Go](https://github.com/bxcodec/go-clean-arch)

---

**Last Updated**: January 2026  
**Maintainer**: Architecture Team
