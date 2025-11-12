# CI/CD Pipeline Setup Guide

This guide will help you set up the complete CI/CD pipeline for your Stream application using GitHub Actions and Docker.

## Overview

The CI/CD pipeline automates:
- **Frontend testing**: Linting, type checking, and building
- **Backend testing**: Go tests and building
- **Docker image building**: Creates and pushes images to GitHub Container Registry
- **Deployment**: Automatically deploys to production server
- **Health checks**: Verifies deployment success

## Prerequisites

1. **GitHub Repository** with your code
2. **Production Server** with Docker and Docker Compose installed
3. **Domain name** pointing to your server
4. **SSH access** to your production server

## Step 1: Server Setup

### Install Docker and Docker Compose
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt install docker-compose-plugin -y

# Add user to docker group
sudo usermod -aG docker $USER
```

### Setup Project Directory
```bash
# Create project directory
sudo mkdir -p /opt/stream
cd /opt/stream

# Clone your repository
git clone https://github.com/yourusername/your-repo.git .

# Make sure deploy.sh is executable
chmod +x deploy.sh
```

## Step 2: GitHub Secrets Configuration

Navigate to **Settings** → **Secrets and variables** → **Actions** in your GitHub repository and add these secrets:

### Required Secrets
- `PROD_HOST`: Your server IP (e.g., `192.168.1.100`)
- `PROD_USER`: SSH username (e.g., `root` or `ubuntu`)
- `PROD_SSH_KEY`: Private SSH key content
- `REDIS_PASSWORD`: Secure Redis password (16+ chars)
- `FEATURED_CHANNES`: Featured channels (e.g., `channel1,channel2,channel3`)
- `DOMAIN`: Your domain (e.g., `freetvchannels.online`)
- `SSL_EMAIL`: Email for SSL certificates

### SSH Key Setup
```bash
# Generate SSH key on your local machine
ssh-keygen -t rsa -b 4096 -C "github-actions"

# Copy public key to server
ssh-copy-id -i ~/.ssh/id_rsa.pub user@your-server

# Copy private key content to GitHub secret
cat ~/.ssh/id_rsa
```

## Step 3: Environment Files

### Production Environment
Create `.env` on your server:
```bash
# Redis Configuration
REDIS_PASSWORD=your_secure_redis_password_here
REDIS_DB=0

# Featured Channels
FEATURED_CHANNES=channel1,channel2,channel3

# Domain Configuration
DOMAIN=freetvchannels.online

# SSL Email
SSL_EMAIL=admin@freetvchannels.online
```

## Step 4: First Time Deployment

### SSL Certificate Setup
```bash
# Navigate to project directory
cd /opt/stream

# Generate SSL certificates (first time only)
sudo docker-compose -f docker-compose.prod.yml run --rm certbot certonly \
  --webroot --webroot-path=/var/www/certbot \
  --email admin@freetvchannels.online \
  --agree-tos --no-eff-email \
  -d freetvchannels.online
```

### Initial Deployment
```bash
# Deploy the application
./deploy.sh
```

## Step 5: CI/CD Pipeline

The pipeline triggers automatically on:
- **Push to main branch**: Full deployment
- **Pull requests**: Testing only
- **Push to develop branch**: Testing and building

### Pipeline Stages
1. **Frontend Test**: Linting, type checking, building
2. **Backend Test**: Go tests, building
3. **Build & Push**: Create Docker images, push to GHCR
4. **Deploy**: Deploy to production server
5. **Health Check**: Verify deployment success

## Step 6: Monitoring

### Check Pipeline Status
- Go to **Actions** tab in your GitHub repository
- View workflow runs and their status

### Check Deployment
```bash
# Check running containers
docker ps

# Check logs
docker-compose -f docker-compose.ci.yml logs -f

# Check health endpoint
curl https://freetvchannels.online/api/health
```

## Troubleshooting

### Common Issues

1. **SSH Connection Failed**
   - Verify SSH key is correctly added to secrets
   - Check server firewall allows SSH
   - Ensure user has proper permissions

2. **Docker Build Failed**
   - Check Dockerfile syntax
   - Verify all dependencies are installed
   - Check build logs for specific errors

3. **Deployment Failed**
   - Verify environment variables are set
   - Check Docker Compose configuration
   - Ensure server has enough resources

4. **SSL Certificate Issues**
   - Verify domain points to server IP
   - Check port 80/443 are open
   - Review certbot logs

### Manual Rollback
```bash
# Stop current deployment
docker-compose -f docker-compose.ci.yml down

# Switch to previous version
git checkout previous-commit-hash

# Redeploy
./deploy.sh
```

## Security Best Practices

1. **Rotate secrets regularly** - Update SSH keys and passwords
2. **Use strong passwords** - At least 16 characters for Redis
3. **Limit access** - Restrict SSH to specific IPs
4. **Monitor logs** - Regularly check deployment and access logs
5. **Update dependencies** - Keep Docker images and packages updated

## Performance Optimization

1. **Enable caching** - Docker layer caching in CI/CD
2. **Use CDN** - For static assets
3. **Monitor resources** - CPU, memory, and disk usage
4. **Optimize images** - Multi-stage builds, smaller base images

## Support

For issues with:
- **GitHub Actions**: Check [GitHub Actions Documentation](https://docs.github.com/en/actions)
- **Docker**: Check [Docker Documentation](https://docs.docker.com/)
- **Deployment**: Review logs and check configuration files
