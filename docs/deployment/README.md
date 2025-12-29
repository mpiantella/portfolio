# AWS Deployment Options

This document provides a comprehensive overview of all available AWS deployment options for the portfolio website application.

## Quick Comparison

| Option | Monthly Cost | Complexity | Setup Time | Best For |
|--------|-------------|------------|------------|----------|
| **GitHub Pages** | $0 | Low | 20 min | Static sites, zero cost |
| **App Runner** | $5-10 | Low | 15 min | Development, low traffic |
| **Lightsail** | $7-25 | Very Low | 10 min | Predictable costs, simple setup |
| **ECS (No ALB)** | $15 | Medium | 30 min | Cost-conscious production |
| **ECS (Full)** | $32-35 | High | 1-2 hours | High availability, auto-scaling |

## Deployment Options

### 1. [GitHub Pages](./github-pages.md) 💰 **Free Option**
Free static site hosting with GitHub's global CDN.

**Pros:**
- Completely free
- Fast global CDN
- Easy deployment
- HTTPS included
- Custom domains supported

**Cons:**
- Requires conversion to static site
- No server-side logic
- Build step required
- External service for forms

[Read Full Documentation →](./github-pages.md)

---

### 2. [AWS App Runner](./app-runner.md) ⭐ **Recommended for Most Users**
Fully managed service that makes it easy to deploy containerized web applications.

**Pros:**
- Simplest setup
- Pay only for usage
- Auto-scaling built-in
- No infrastructure management
- HTTPS included

**Cons:**
- Less control over infrastructure
- Limited customization options

[Read Full Documentation →](./app-runner.md)

---

### 3. [AWS Lightsail Containers](./lightsail.md)
Simplified container service with fixed pricing.

**Pros:**
- Predictable monthly costs
- Very simple to set up
- Built-in load balancing
- Free SSL certificates
- Domain management included

**Cons:**
- Limited scaling options
- Fixed resource allocations

[Read Full Documentation →](./lightsail.md)

---

### 4. [ECS Fargate (Simplified)](./ecs-simplified.md)
Container service without load balancer, using public IP.

**Pros:**
- More control than App Runner
- Lower cost than full ECS
- Uses standard AWS services
- Good for learning ECS

**Cons:**
- Single container instance
- IP changes on restart
- Manual HTTPS setup

[Read Full Documentation →](./ecs-simplified.md)

---

### 5. [ECS Fargate (Full Production)](./ecs-full.md)
Complete production setup with load balancing and auto-scaling.

**Pros:**
- High availability
- Auto-scaling
- Production-ready
- Full control

**Cons:**
- Most expensive
- Complex setup
- Many resources to manage

[Read Full Documentation →](./ecs-full.md)

---

## Architecture Diagrams

### High-Level Comparison

```
┌─────────────────────────────────────────────────────────────────┐
│                     Deployment Options Overview                  │
└─────────────────────────────────────────────────────────────────┘

1. APP RUNNER (Simplest)
   ┌──────────┐    ┌─────────────┐    ┌──────────────┐
   │ GitHub   │───▶│ ECR         │───▶│ App Runner   │
   │ Actions  │    │ (Images)    │    │ Service      │
   └──────────┘    └─────────────┘    └──────┬───────┘
                                              │
                                              ▼
                                         HTTPS URL
                                      (auto-provisioned)

2. LIGHTSAIL (Fixed Cost)
   ┌──────────┐    ┌─────────────────┐    ┌──────────────┐
   │ GitHub   │───▶│ Lightsail       │───▶│ Public URL   │
   │ Actions  │    │ Container       │    │ with SSL     │
   └──────────┘    └─────────────────┘    └──────────────┘

3. ECS SIMPLIFIED (No ALB)
   ┌──────────┐    ┌─────────┐    ┌─────────┐    ┌──────────┐
   │ GitHub   │───▶│ ECR     │───▶│ ECS     │───▶│ Public   │
   │ Actions  │    │         │    │ Fargate │    │ IP:8080  │
   └──────────┘    └─────────┘    └─────────┘    └──────────┘

4. ECS FULL (Production)
   ┌──────────┐    ┌─────────┐    ┌─────────┐    ┌─────┐    ┌──────┐
   │ GitHub   │───▶│ ECR     │───▶│ ECS     │───▶│ ALB │───▶│ HTTPS│
   │ Actions  │    │         │    │ Fargate │    │     │    │ URL  │
   └──────────┘    └─────────┘    └─────────┘    └─────┘    └──────┘
```

## Decision Guide

### Choose **App Runner** if:
- ✅ You want the simplest setup
- ✅ You have variable traffic patterns
- ✅ You want to minimize costs
- ✅ You don't need advanced networking
- ✅ This is a portfolio/personal project

### Choose **Lightsail** if:
- ✅ You want predictable monthly costs
- ✅ You want the absolute simplest setup
- ✅ You're new to AWS
- ✅ You don't need auto-scaling
- ✅ You want everything in one service

### Choose **ECS Simplified** if:
- ✅ You want to learn ECS
- ✅ You need more control than App Runner
- ✅ You're comfortable with AWS
- ✅ Cost is a primary concern
- ✅ You can tolerate downtime during restarts

### Choose **ECS Full** if:
- ✅ This is a production application
- ✅ You need high availability
- ✅ You need auto-scaling
- ✅ You need advanced networking
- ✅ Budget is not the primary concern

## Getting Started

1. Choose your deployment option from the comparison above
2. Follow the detailed guide for your chosen option
3. Complete the one-time AWS setup
4. Configure GitHub Actions secrets
5. Deploy!

## Support & Resources

- [AWS App Runner Documentation](https://docs.aws.amazon.com/apprunner/)
- [AWS Lightsail Documentation](https://docs.aws.amazon.com/lightsail/)
- [AWS ECS Documentation](https://docs.aws.amazon.com/ecs/)
- [GitHub Actions Documentation](https://docs.github.com/actions)

## Next Steps

Start with the recommended option: [AWS App Runner Setup Guide →](./app-runner.md)
