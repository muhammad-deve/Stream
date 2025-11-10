# Quick Deploy Guide - TL;DR

For complete instructions, see `DEPLOYMENT_GUIDE.md`

## 🎯 Quick Setup (On Your Server)

```bash
# 1. Install Docker & Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
apt install docker-compose-plugin -y

# 2. Clone repository
cd /opt
git clone <your-repo> stream
cd stream/back

# 3. Create environment file
cp .env.prod.example .env.prod
nano .env.prod  # Edit with your settings

# 4. Set DNS (GoDaddy)
# Point freetvchannels.online to: 2a02:c207:2287:7544::1

# 5. Run deployment script
chmod +x deploy.sh
./deploy.sh

# 6. Follow prompts to get SSL certificate
```

## 📋 Environment Variables (.env.prod)

```env
REDIS_PASSWORD=your_strong_password_here
REDIS_DB=0
FEATURED_CHANNES=9m3fy1v50tlelwe,75l42yzl07l4948,3v539ya9sxjr8b6,fk9p603ke1k9p01,9nytgzajem68r9p,234c5688eyn2n29
```

Generate password:
```bash
openssl rand -base64 32
```

## 🌐 DNS Configuration (GoDaddy)

| Type | Name | Value                      |
|------|------|----------------------------|
| AAAA | @    | 2a02:c207:2287:7544::1    |
| AAAA | www  | 2a02:c207:2287:7544::1    |

## ✅ Verify Deployment

```bash
# Check services
docker compose -f docker-compose.prod.yml ps

# View logs
docker compose -f docker-compose.prod.yml logs -f

# Test API
curl https://freetvchannels.online/api/v1/stream/featured

# Test Redis
docker compose -f docker-compose.prod.yml exec redis redis-cli -a YOUR_PASSWORD ping
```

## 🔧 Useful Commands

```bash
# View logs
docker compose -f docker-compose.prod.yml logs -f [service]

# Restart service
docker compose -f docker-compose.prod.yml restart [service]

# Update application
git pull && ./deploy.sh

# Stop everything
docker compose -f docker-compose.prod.yml down

# Backup Redis
docker compose -f docker-compose.prod.yml exec redis redis-cli -a PASSWORD SAVE
```

## 🔐 Security Checklist

- [x] HTTPS with Let's Encrypt
- [x] Strong Redis password
- [x] Firewall configured (ports 80, 443)
- [ ] SSH key authentication only
- [ ] Fail2ban installed
- [ ] Regular backups scheduled

## 🐛 Troubleshooting

| Issue | Solution |
|-------|----------|
| Site not loading | Check DNS propagation: `nslookup freetvchannels.online` |
| 502 Bad Gateway | Check app logs: `docker compose logs app` |
| SSL errors | Re-run certbot: See DEPLOYMENT_GUIDE.md |
| Redis connection | Check password in .env.prod |

## 📞 Quick Help

```bash
# Service not starting?
docker compose -f docker-compose.prod.yml logs [service]

# Need to rebuild?
docker compose -f docker-compose.prod.yml build --no-cache

# Reset everything?
docker compose -f docker-compose.prod.yml down -v
./deploy.sh
```

---

**Need detailed steps?** → See `DEPLOYMENT_GUIDE.md`  
**Redis info?** → See `QUICK_START_REDIS.md`
