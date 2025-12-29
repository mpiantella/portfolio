# Portfolio Website Documentation

Welcome to the documentation for the Maria Lucena Portfolio Website.

## 📚 Documentation Index

### Deployment Guides

Comprehensive guides for deploying your portfolio website to AWS.

- **[Deployment Overview](./deployment/README.md)** - Start here for an overview of all options
- **[Comparison Guide](./deployment/comparison.md)** - Detailed side-by-side comparison of all deployment options

#### Deployment Options (Ordered by Cost)

1. **[GitHub Pages](./deployment/github-pages.md)** 💰 **Free Option**
   - Cost: $0/month
   - Static site hosting
   - Requires conversion
   - **Best for zero-cost hosting**

2. **[AWS Lightsail](./deployment/lightsail.md)** ⭐ **Simplest Dynamic**
   - Fixed pricing: $7-25/month
   - All-in-one solution
   - Perfect for portfolios
   - **Recommended for most users**

3. **[AWS App Runner](./deployment/app-runner.md)** ⭐ **Best for Growth**
   - Pay-per-use: $5-40/month
   - Auto-scaling included
   - Minimal configuration
   - **Recommended for production**

4. **[ECS Fargate Simplified](./deployment/ecs-simplified.md)**
   - Fixed cost: $15-30/month
   - No load balancer
   - Good for learning ECS
   - More AWS control

5. **[ECS Fargate Full](./deployment/ecs-full.md)**
   - Production-grade: $85+/month
   - High availability
   - Complete AWS setup
   - **For enterprise needs**

## 🚀 Quick Start

### First Time Deployment

1. **Choose your deployment option** from the guides above
2. **Set up AWS account** (if you don't have one)
3. **Follow the step-by-step guide** for your chosen option
4. **Configure GitHub secrets** for CI/CD
5. **Push to main branch** to deploy automatically

### Recommended Path

```
Want zero cost?
    └─→ Use GitHub Pages (Free, but static)

New to AWS?
    └─→ Start with Lightsail ($7/mo)

Need auto-scaling?
    └─→ Use App Runner ($30-40/mo)

Need full control?
    └─→ Use ECS Full ($85+/mo)
```

## 📖 Documentation Structure

```
docs/
├── README.md (this file)
└── deployment/
    ├── README.md           # Deployment overview
    ├── comparison.md       # Detailed comparison
    ├── app-runner.md       # AWS App Runner guide
    ├── lightsail.md        # AWS Lightsail guide
    ├── ecs-simplified.md   # ECS without ALB
    └── ecs-full.md         # Full production ECS
```

## 🎯 Common Tasks

### Deploy to Production
```bash
# Push to main branch
git add .
git commit -m "Deploy to production"
git push origin main

# GitHub Actions will automatically deploy
```

### View Deployment Status
- GitHub Actions tab in your repository
- Check AWS console for your chosen service
- Review CloudWatch logs for errors

### Update Configuration
1. Modify application code
2. Test locally: `go run cmd/server/main.go`
3. Commit and push changes
4. Automatic deployment via GitHub Actions

## 💰 Cost Estimates

| Deployment Type | Monthly Cost | Best For |
|----------------|--------------|----------|
| GitHub Pages | $0 | Static portfolio (no server-side) |
| Lightsail Nano | $7 | Personal portfolio |
| Lightsail Small | $25 | Professional site |
| App Runner (low traffic) | $10-15 | Development |
| App Runner (production) | $30-40 | Growing business |
| ECS Simplified | $15-30 | Learning AWS |
| ECS Full | $85+ | Enterprise |

## 🛠️ Technologies Used

- **Backend**: Go 1.21
- **Frontend**: HTML, TailwindCSS, HTMX
- **Deployment**: AWS (multiple options)
- **CI/CD**: GitHub Actions
- **Container**: Docker

## 📝 Application Structure

```
portfolio-website/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── domain/          # Business logic
│   ├── infrastructure/  # Data, config, logging
│   └── interfaces/      # HTTP handlers
├── web/
│   ├── static/          # CSS, images
│   └── templates/       # HTML templates
├── docs/                # Documentation
└── .github/workflows/   # CI/CD pipelines
```

## 🔧 Local Development

### Prerequisites
- Go 1.21+
- Docker (for containerization)
- Node.js (for TailwindCSS)

### Run Locally
```bash
# Install dependencies
go mod download
npm install

# Build CSS
npm run build:css

# Run server
go run cmd/server/main.go

# Access at http://localhost:8080
```

### Build Docker Image
```bash
docker build -t portfolio:latest .
docker run -p 8080:8080 portfolio:latest
```

## 📊 Monitoring & Logs

### CloudWatch Logs
All deployment options use CloudWatch for logging:
```bash
# View logs
aws logs tail /ecs/portfolio --follow

# Filter errors
aws logs filter-log-events \
  --log-group-name /ecs/portfolio \
  --filter-pattern "ERROR"
```

### Metrics
Monitor these key metrics:
- Request count
- Response time
- Error rate (4xx, 5xx)
- CPU/Memory utilization

## 🔒 Security

### Best Practices
- Use IAM roles with minimal permissions
- Enable HTTPS for production
- Rotate access keys regularly
- Keep dependencies updated
- Use security groups to restrict access

### Secrets Management
Never commit secrets to Git. Use:
- GitHub Secrets for CI/CD
- AWS Systems Manager Parameter Store
- AWS Secrets Manager (for production)

## 🆘 Troubleshooting

### Common Issues

**Deployment Fails**
- Check GitHub Actions logs
- Verify AWS credentials
- Ensure ECR repository exists
- Check service quotas

**Application Won't Start**
- Verify health check endpoint
- Check CloudWatch logs
- Ensure correct port (8080)
- Verify environment variables

**High Costs**
- Review CloudWatch billing
- Check auto-scaling settings
- Consider smaller instance sizes
- Use App Runner pay-per-use model

## 📚 Additional Resources

### AWS Documentation
- [AWS App Runner](https://docs.aws.amazon.com/apprunner/)
- [AWS Lightsail](https://docs.aws.amazon.com/lightsail/)
- [Amazon ECS](https://docs.aws.amazon.com/ecs/)
- [Elastic Container Registry](https://docs.aws.amazon.com/ecr/)

### Tools & Services
- [GitHub Actions](https://docs.github.com/actions)
- [Docker](https://docs.docker.com/)
- [Go](https://go.dev/doc/)
- [HTMX](https://htmx.org/)

## 🤝 Contributing

Improvements to documentation are welcome!

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📧 Support

For questions or issues:
1. Check the deployment guides
2. Review troubleshooting section
3. Check AWS documentation
4. Create a GitHub issue

## 📄 License

This project is for portfolio purposes.

---

**Ready to deploy?** Start with the [Deployment Overview](./deployment/README.md) →
