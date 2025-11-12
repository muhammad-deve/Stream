# GitHub Secrets Configuration

This document outlines all the secrets that need to be configured in your GitHub repository for the CI/CD pipeline to work properly.

## Required GitHub Secrets

### Production Environment
    - `PROD_HOST`: Your production server IP address or hostname
    - `PROD_USER`: SSH username for your production server
    - `PROD_SSH_KEY`: Private SSH key for accessing the production server
    - `REDIS_PASSWORD`: Secure password for Redis (at least 16 characters)
    - `FEATURED_CHANNES`: Comma-separated list of featured channel names
    - `DOMAIN`: Your domain name (e.g., freetvchannels.online)
    - `SSL_EMAIL`: Email address for SSL certificate notifications

## How to Add Secrets

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add each secret from the list above

## SSH Key Setup

1. Generate a new SSH key pair:
   ```bash
   ssh-keygen -t rsa -b 4096 -C "github-actions"
   ```

2. Add the public key to your server:
   ```bash
   cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys
   ```

3. Copy the private key and add it as `PROD_SSH_KEY` secret:
   ```bash
   cat ~/.ssh/id_rsa
   ```

## Environment Variables

The following environment variables are automatically set by the workflow:
- `REGISTRY`: ghcr.io (GitHub Container Registry)
- `IMAGE_NAME`: Your GitHub repository name

## Security Notes

- Never commit secrets to your repository
- Use strong, unique passwords for Redis
- Rotate SSH keys periodically
- Limit SSH access to specific IP addresses if possible
- Use GitHub's protected branches for main branch
