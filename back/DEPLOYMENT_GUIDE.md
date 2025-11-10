# Production Deployment Guide - freetvchannels.online

## 🌐 Server Information
- **Domain**: freetvchannels.online
- **Server IP**: 2a02:c207:2287:7544::1 (IPv6)
- **SSL**: Let's Encrypt (Free, Auto-Renewal)

---

## 📋 Prerequisites

1. **Server Access**: SSH access to root@vmi2877544
2. **Domain DNS**: Point your domain to server IP
3. **Docker & Docker Compose**: Installed on server
4. **Firewall**: Open ports 80, 443, 8090

---

## 🚀 Step-by-Step Deployment

### 1. Configure DNS (GoDaddy)

Log into your GoDaddy account and add these DNS records:

```
Type    Name    Value                           TTL
A       @       (IPv4 if you have one)          600
AAAA    @       2a02:c207:2287:7544::1          600
AAAA    www     2a02:c207:2287:7544::1          600
```

**Note**: Since you have IPv6, make sure to add AAAA records. If your server also has IPv4, add A records too.

Wait 10-30 minutes for DNS propagation. Check with:
```bash
nslookup freetvchannels.online
```

---

### 2. Server Setup

SSH into your server:
```bash
ssh root@vmi2877544
```

Install Docker and Docker Compose (if not already installed):
```bash
# Update system
apt update && apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Install Docker Compose
apt install docker-compose-plugin -y

# Verify installation
docker --version
docker compose version
```

---

### 3. Clone Your Repository

```bash
cd /opt
git clone <your-repo-url> stream
cd stream/back
```

---

### 4. Create Production Environment File

Create `.env.prod` file:
```bash
nano .env.prod
```

Add these variables:
```env
# Redis Configuration
REDIS_PASSWORD=YOUR_STRONG_PASSWORD_HERE
REDIS_DB=0

# Channel Configuration
FEATURED_CHANNES=9m3fy1v50tlelwe,75l42yzl07l4948,3v539ya9sxjr8b6,fk9p603ke1k9p01,9nytgzajem68r9p,234c5688eyn2n29
```

**⚠️ IMPORTANT**: Replace `YOUR_STRONG_PASSWORD_HERE` with a strong password!

Generate a strong password:
```bash
openssl rand -base64 32
```

---

### 5. Initial Nginx Setup (HTTP Only - For SSL Certificate)

First, we need to get SSL certificates. Temporarily modify nginx config:

Create `nginx/conf.d/initial.conf`:
```bash
mkdir -p nginx/conf.d certbot/conf certbot/www
nano nginx/conf.d/initial.conf
```

Add:
```nginx
server {
    listen 80;
    listen [::]:80;
    server_name freetvchannels.online www.freetvchannels.online;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        proxy_pass http://app:8090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

### 6. Start Services (Initial HTTP)

```bash
# Load environment variables
export $(cat .env.prod | xargs)

# Start services
docker compose -f docker-compose.prod.yml up -d
```

Check if services are running:
```bash
docker compose -f docker-compose.prod.yml ps
docker compose -f docker-compose.prod.yml logs -f
```

Test your site:
```bash
curl http://freetvchannels.online
```

---

### 7. Obtain SSL Certificate

Run certbot to get SSL certificate:
```bash
docker compose -f docker-compose.prod.yml run --rm certbot certonly \
  --webroot \
  --webroot-path=/var/www/certbot \
  --email your-email@example.com \
  --agree-tos \
  --no-eff-email \
  -d freetvchannels.online \
  -d www.freetvchannels.online
```

**Replace** `your-email@example.com` with your actual email.

Verify certificates were created:
```bash
ls -la certbot/conf/live/freetvchannels.online/
```

You should see:
- `fullchain.pem`
- `privkey.pem`

---

### 8. Enable HTTPS

Now replace the initial config with the full HTTPS config:
```bash
rm nginx/conf.d/initial.conf
cp nginx/conf.d/freetvchannels.conf nginx/conf.d/freetvchannels.conf
```

The file should already exist from your repository. Restart nginx:
```bash
docker compose -f docker-compose.prod.yml restart nginx
```

---

### 9. Test HTTPS

Visit your site:
- **HTTP**: http://freetvchannels.online (should redirect to HTTPS)
- **HTTPS**: https://freetvchannels.online (should work with green lock)

Test SSL:
```bash
curl -I https://freetvchannels.online
```

Check SSL grade:
https://www.ssllabs.com/ssltest/analyze.html?d=freetvchannels.online

---

### 10. Configure Firewall

```bash
# Allow SSH (if not already)
ufw allow 22/tcp

# Allow HTTP and HTTPS
ufw allow 80/tcp
ufw allow 443/tcp

# Enable firewall
ufw enable
ufw status
```

---

## 🔄 Automatic SSL Renewal

The certbot container is configured to automatically renew certificates every 12 hours. Check renewal status:

```bash
# Test renewal (dry run)
docker compose -f docker-compose.prod.yml run --rm certbot renew --dry-run

# Force renewal (if needed)
docker compose -f docker-compose.prod.yml run --rm certbot renew --force-renewal
docker compose -f docker-compose.prod.yml restart nginx
```

---

## 📊 Monitoring & Maintenance

### Check Logs
```bash
# All services
docker compose -f docker-compose.prod.yml logs -f

# Specific service
docker compose -f docker-compose.prod.yml logs -f app
docker compose -f docker-compose.prod.yml logs -f redis
docker compose -f docker-compose.prod.yml logs -f nginx
```

### Check Redis
```bash
# Connect to Redis CLI
docker compose -f docker-compose.prod.yml exec redis redis-cli -a YOUR_REDIS_PASSWORD

# Inside Redis CLI:
> PING
> INFO
> KEYS *
> TTL <some-token>
> exit
```

### Restart Services
```bash
# Restart all
docker compose -f docker-compose.prod.yml restart

# Restart specific service
docker compose -f docker-compose.prod.yml restart app
```

### Update Application
```bash
cd /opt/stream/back
git pull origin master
docker compose -f docker-compose.prod.yml build app
docker compose -f docker-compose.prod.yml up -d app
```

---

## 🔐 Security Checklist

- [x] HTTPS enabled with Let's Encrypt
- [x] Redis password protected
- [x] Redis only accessible from localhost
- [x] Security headers configured in Nginx
- [x] Firewall configured
- [ ] Consider setting up fail2ban for SSH protection
- [ ] Regular backups of Redis data
- [ ] Monitor server resources (CPU, RAM, disk)

### Setup Fail2Ban (Optional but Recommended)
```bash
apt install fail2ban -y
systemctl enable fail2ban
systemctl start fail2ban
```

---

## 🗄️ Backup Strategy

### Backup Redis Data
```bash
# Manual backup
docker compose -f docker-compose.prod.yml exec redis redis-cli -a YOUR_REDIS_PASSWORD SAVE
cp -r redis-data /backup/redis-$(date +%Y%m%d-%H%M%S)

# Automated daily backup (cron)
echo "0 2 * * * cd /opt/stream/back && docker compose -f docker-compose.prod.yml exec redis redis-cli -a YOUR_REDIS_PASSWORD SAVE && cp -r redis-data /backup/redis-\$(date +\%Y\%m\%d)" | crontab -
```

---

## 🐛 Troubleshooting

### Issue: Site not accessible
```bash
# Check if Docker containers are running
docker compose -f docker-compose.prod.yml ps

# Check nginx logs
docker compose -f docker-compose.prod.yml logs nginx

# Check app logs
docker compose -f docker-compose.prod.yml logs app

# Test nginx config
docker compose -f docker-compose.prod.yml exec nginx nginx -t
```

### Issue: SSL certificate errors
```bash
# Check certificate files
ls -la certbot/conf/live/freetvchannels.online/

# Re-obtain certificate
docker compose -f docker-compose.prod.yml run --rm certbot certonly --force-renewal \
  --webroot --webroot-path=/var/www/certbot \
  -d freetvchannels.online -d www.freetvchannels.online
```

### Issue: Redis connection failed
```bash
# Check Redis logs
docker compose -f docker-compose.prod.yml logs redis

# Test Redis connection
docker compose -f docker-compose.prod.yml exec redis redis-cli -a YOUR_REDIS_PASSWORD ping

# Restart Redis
docker compose -f docker-compose.prod.yml restart redis
```

### Issue: App crashes
```bash
# Check app logs for errors
docker compose -f docker-compose.prod.yml logs app | tail -100

# Check app health
curl http://localhost:8090/api/v1/health

# Rebuild and restart app
docker compose -f docker-compose.prod.yml build app
docker compose -f docker-compose.prod.yml up -d app
```

---

## 📱 API Endpoints

After deployment, your API will be available at:

- **Featured Channels**: `https://freetvchannels.online/api/v1/stream/featured`
- **Watch Channel**: `https://freetvchannels.online/api/v1/stream/watch`
- **Play Stream**: `https://freetvchannels.online/api/v1/stream/play`
- **Search**: `https://freetvchannels.online/api/v1/stream/search`

Test API:
```bash
curl https://freetvchannels.online/api/v1/stream/featured
```

---

## 🎯 Performance Optimization

### Nginx Caching (Add to nginx config if needed)
```nginx
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m max_size=1g inactive=60m;
```

### Redis Memory Limit (Modify docker-compose if needed)
```yaml
redis:
  command: redis-server --appendonly yes --requirepass "${REDIS_PASSWORD}" --maxmemory 256mb --maxmemory-policy allkeys-lru
```

---

## ✅ Final Verification

1. Visit: https://freetvchannels.online
2. Check SSL: https://www.ssllabs.com/ssltest/analyze.html?d=freetvchannels.online
3. Test API: `curl https://freetvchannels.online/api/v1/stream/featured`
4. Check Redis: `docker compose exec redis redis-cli -a PASSWORD ping`
5. Monitor logs: `docker compose logs -f`

---

## 📞 Support

If you encounter any issues:
1. Check the logs: `docker compose -f docker-compose.prod.yml logs`
2. Verify environment variables: `cat .env.prod`
3. Test network connectivity: `curl -v http://localhost:8090`
4. Check DNS: `nslookup freetvchannels.online`

---

## 🎉 Deployment Complete!

Your application is now live with:
- ✅ HTTPS enabled
- ✅ Redis for token management
- ✅ Auto-renewing SSL certificates
- ✅ Production-ready configuration
- ✅ Nginx reverse proxy

**URL**: https://freetvchannels.online
