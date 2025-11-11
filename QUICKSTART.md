# ⚡ Quick Start - Deploy in 5 Minutes

## TL;DR

```bash
cd ~/Stream

# 1. Create environment file
cat > .env << 'EOF'
REDIS_PASSWORD=your_secure_password_here
REDIS_DB=0
FEATURED_CHANNES=channel1,channel2,channel3
EOF

# 2. Generate SSL (first time only)
mkdir -p certbot/conf certbot/www nginx/conf.d
docker run --rm -v $(pwd)/certbot/conf:/etc/letsencrypt -v $(pwd)/certbot/www:/var/www/certbot \
  certbot/certbot certonly --webroot -w /var/www/certbot \
  -d freetvchannels.online -d www.freetvchannels.online \
  --agree-tos --no-eff-email -m your-email@example.com

# 3. Deploy
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d

# 4. Verify
docker-compose -f docker-compose.prod.yml ps
curl https://freetvchannels.online
```

## What You Get

✅ Frontend at `https://freetvchannels.online`
✅ API at `https://freetvchannels.online/api/v1/stream/*`
✅ SSL/HTTPS with auto-renewal
✅ All services in Docker containers

## Key Files

| File | Purpose |
|------|---------|
| `docker-compose.prod.yml` | Main deployment file |
| `.env` | Your configuration (create this) |
| `nginx/conf.d/default.conf` | Nginx routing |
| `front/Dockerfile` | Frontend build |
| `DEPLOYMENT.md` | Full guide |

## Useful Commands

```bash
# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Restart
docker-compose -f docker-compose.prod.yml restart

# Stop
docker-compose -f docker-compose.prod.yml down

# Check status
docker-compose -f docker-compose.prod.yml ps
```

## Troubleshooting

**Frontend not loading?**
```bash
docker-compose -f docker-compose.prod.yml logs front-app
```

**API not working?**
```bash
docker-compose -f docker-compose.prod.yml logs back-app
```

**SSL issues?**
```bash
docker-compose -f docker-compose.prod.yml logs certbot
```

---

**For detailed info**: See `DEPLOYMENT.md` or `README_DEPLOYMENT.md`
