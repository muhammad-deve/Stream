# 🚀 Deployment Checklist - freetvchannels.online

## ✅ What's Already Done

- [x] **Docker configuration**
  - `Dockerfile` - App containerization
  - `docker-compose.prod.yml` - Production orchestration with Redis & Nginx
  - `.dockerignore` - Optimized build context

- [x] **Nginx & SSL Setup**
  - `nginx/nginx.conf` - Main nginx configuration
  - `nginx/conf.d/freetvchannels.conf` - Domain config with HTTPS
  - SSL certificates will be obtained via Let's Encrypt

- [x] **Redis Integration**
  - Redis client configured in Go application
  - Token-based URL obfuscation implemented
  - 2-hour token expiration
  - Redis persistence configured

- [x] **CI/CD Pipeline**
  - `.gitlab-ci.yml` - Automated build & deploy
  - Docker registry integration

- [x] **Documentation**
  - `DEPLOYMENT_GUIDE.md` - Complete step-by-step guide
  - `QUICK_DEPLOY.md` - Quick reference
  - `QUICK_START_REDIS.md` - Redis setup guide
  - `DEPLOYMENT_CHECKLIST.md` - This file

- [x] **Automation**
  - `deploy.sh` - Automated deployment script
  - SSL auto-renewal configured

---

## 🎯 What You Need to Do Now

### 1. DNS Configuration (5 minutes)
- [ ] Log into GoDaddy
- [ ] Add DNS records for `freetvchannels.online`:
  ```
  Type: AAAA
  Name: @
  Value: 2a02:c207:2287:7544::1
  TTL: 600
  
  Type: AAAA
  Name: www
  Value: 2a02:c207:2287:7544::1
  TTL: 600
  ```
- [ ] Wait 10-30 mins for DNS propagation
- [ ] Verify: `nslookup freetvchannels.online`

### 2. Server Preparation (10 minutes)
- [ ] SSH into server: `ssh root@vmi2877544`
- [ ] Update system: `apt update && apt upgrade -y`
- [ ] Install Docker:
  ```bash
  curl -fsSL https://get.docker.com -o get-docker.sh
  sh get-docker.sh
  apt install docker-compose-plugin -y
  ```
- [ ] Configure firewall:
  ```bash
  ufw allow 22/tcp
  ufw allow 80/tcp
  ufw allow 443/tcp
  ufw enable
  ```

### 3. Deploy Application (10 minutes)
- [ ] Clone repository:
  ```bash
  cd /opt
  git clone <your-repo-url> stream
  cd stream/back
  ```
- [ ] Create environment file:
  ```bash
  cp .env.prod.example .env.prod
  nano .env.prod
  ```
- [ ] Generate Redis password: `openssl rand -base64 32`
- [ ] Update `.env.prod` with:
  - Strong Redis password
  - Your featured channels

### 4. Run Deployment Script (5 minutes)
- [ ] Make script executable: `chmod +x deploy.sh`
- [ ] Run deployment: `./deploy.sh`
- [ ] Follow prompts
- [ ] Enter email for SSL certificate
- [ ] Wait for SSL certificate generation

### 5. Verification (5 minutes)
- [ ] Visit: https://freetvchannels.online
- [ ] Check SSL: Green lock in browser
- [ ] Test API: `curl https://freetvchannels.online/api/v1/stream/featured`
- [ ] Check Redis: `docker compose -f docker-compose.prod.yml exec redis redis-cli -a PASSWORD ping`
- [ ] View logs: `docker compose -f docker-compose.prod.yml logs -f`

---

## 📦 Files Summary

### Configuration Files (Already Created)
```
back/
├── docker-compose.prod.yml       # Main production config
├── Dockerfile                     # App container
├── .dockerignore                  # Build optimization
├── .env.prod.example              # Environment template
├── deploy.sh                      # Automated deployment
├── nginx/
│   ├── nginx.conf                # Nginx main config
│   └── conf.d/
│       └── freetvchannels.conf   # Domain + SSL config
└── Documentation/
    ├── DEPLOYMENT_GUIDE.md       # Complete guide
    ├── QUICK_DEPLOY.md           # Quick reference
    ├── QUICK_START_REDIS.md      # Redis info
    └── DEPLOYMENT_CHECKLIST.md   # This file
```

### Files You Need to Create on Server
```
/opt/stream/back/
└── .env.prod                     # YOUR production environment
```

---

## 🔧 Environment Variables (.env.prod)

Create this file with:
```env
# Generate with: openssl rand -base64 32
REDIS_PASSWORD=YOUR_STRONG_PASSWORD_HERE

# Default is fine
REDIS_DB=0

# Your featured channels
FEATURED_CHANNES=9m3fy1v50tlelwe,75l42yzl07l4948,3v539ya9sxjr8b6,fk9p603ke1k9p01,9nytgzajem68r9p,234c5688eyn2n29
```

---

## 🎯 Expected Result

After deployment:
- ✅ App running at: https://freetvchannels.online
- ✅ API accessible: https://freetvchannels.online/api/v1/
- ✅ SSL certificate (A+ rating)
- ✅ Redis caching tokens
- ✅ Auto-renewing certificates
- ✅ All services in Docker containers
- ✅ Production-ready configuration

---

## 🚨 Common Issues & Solutions

| Issue | Cause | Solution |
|-------|-------|----------|
| DNS not resolving | Propagation delay | Wait 30 mins, check `nslookup` |
| SSL cert fails | Domain not pointing to server | Verify DNS first |
| 502 Bad Gateway | App not started | Check `docker compose logs app` |
| Redis connection error | Wrong password | Check `.env.prod` |
| Port already in use | Old services running | `docker compose down` first |

---

## 📊 Post-Deployment Monitoring

### Daily
- [ ] Check uptime: `docker compose ps`
- [ ] Monitor logs: `docker compose logs --tail=100`

### Weekly
- [ ] Check disk space: `df -h`
- [ ] Review Redis memory: `docker compose exec redis redis-cli info memory`
- [ ] Check SSL expiry: `echo | openssl s_client -connect freetvchannels.online:443 2>/dev/null | openssl x509 -noout -dates`

### Monthly
- [ ] System updates: `apt update && apt upgrade`
- [ ] Backup Redis data
- [ ] Review logs for errors

---

## 🔐 Security Hardening (Optional)

After deployment, consider:
- [ ] Setup SSH key authentication, disable password login
- [ ] Install fail2ban: `apt install fail2ban`
- [ ] Configure automated backups
- [ ] Setup monitoring (Prometheus/Grafana)
- [ ] Enable Docker security scanning
- [ ] Review Nginx security headers
- [ ] Setup rate limiting

---

## 📞 Support Commands

```bash
# View all logs
docker compose -f docker-compose.prod.yml logs -f

# Restart specific service
docker compose -f docker-compose.prod.yml restart app

# Check Redis tokens
docker compose -f docker-compose.prod.yml exec redis redis-cli -a PASSWORD
> KEYS *
> TTL <token-id>

# Force SSL renewal
docker compose -f docker-compose.prod.yml run --rm certbot renew --force-renewal

# Rebuild app
git pull
docker compose -f docker-compose.prod.yml build app
docker compose -f docker-compose.prod.yml up -d app

# Nuclear option (reset everything)
docker compose -f docker-compose.prod.yml down -v
./deploy.sh
```

---

## ✅ Final Checklist

Before going live:
- [ ] DNS configured and propagated
- [ ] SSL certificate obtained and working
- [ ] HTTPS redirect working (HTTP → HTTPS)
- [ ] Redis connection working
- [ ] API endpoints responding
- [ ] Featured channels loading
- [ ] Stream playback working
- [ ] Logs showing no errors
- [ ] Firewall configured
- [ ] .env.prod has strong password
- [ ] Backup strategy in place

---

## 🎉 You're Ready!

Once all boxes are checked, your application will be:
- 🌐 Live at https://freetvchannels.online
- 🔒 Secure with HTTPS
- 🚀 Production-ready
- 📊 Monitored
- 🔄 Auto-updating SSL

**Questions?** Check `DEPLOYMENT_GUIDE.md` for detailed instructions!
