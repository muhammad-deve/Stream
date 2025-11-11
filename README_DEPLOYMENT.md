# 🚀 Stream Application - Complete Deployment Setup

## ✅ What's Been Done

Your frontend is now fully configured to work with your domain `freetvchannels.online` and all API calls will go through your domain instead of localhost.

### Changes Made

#### 1. **Frontend API Configuration** ✅
- Updated `front/src/lib/channels.ts` to use environment-based API URLs
- All hardcoded `http://127.0.0.1:8090` calls replaced with `API_BASE_URL` variable
- Created `.env.production` with production API URL

#### 2. **Docker Setup** ✅
- Created `front/Dockerfile` - Builds and serves React app
- Created `docker-compose.prod.yml` - Complete stack with all services
- Services included:
  - Frontend (React on :3000)
  - Backend (PocketBase on :8090)
  - Redis cache
  - Nginx reverse proxy
  - Certbot SSL manager

#### 3. **Nginx Configuration** ✅
- Created `nginx/nginx.conf` - Main nginx configuration
- Created `nginx/conf.d/default.conf` - Site configuration
- Features:
  - ✅ HTTP → HTTPS redirect
  - ✅ Frontend served at `/`
  - ✅ API proxied at `/api/*`
  - ✅ SSL/TLS with auto-renewal
  - ✅ Security headers
  - ✅ Gzip compression
  - ✅ Static file caching

#### 4. **Documentation** ✅
- `DEPLOYMENT.md` - Complete deployment guide
- `SETUP_SUMMARY.md` - Quick reference
- `deploy.sh` - Automated deployment script
- `.env.example` - Environment template

---

## 🎯 How It Works

### Request Flow

```
┌─────────────────────────────────────────────────────────────┐
│                  User's Browser                             │
│            https://freetvchannels.online                    │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
        ┌────────────────────────────────┐
        │    Nginx Reverse Proxy         │
        │  (Ports 80 → 443 SSL)          │
        └────────────────────────────────┘
                    │              │
        ┌───────────┘              └──────────┐
        │                                     │
        ▼                                     ▼
    Frontend                            API Requests
  (React App)                         (Backend API)
   Port 3000                          Port 8090
        │                                     │
        │                                     │
   Serves UI                          Returns Data
   Handles Routing                    Stream Info
   Makes API Calls                    Categories
        │                                     │
        └─────────────┬──────────────────────┘
                      │
            All requests go through
            https://freetvchannels.online
```

### API URL Examples

**Before (Hardcoded):**
```
http://127.0.0.1:8090/api/v1/stream/categories
http://127.0.0.1:8090/api/v1/stream/featured
```

**After (Domain-based):**
```
https://freetvchannels.online/api/v1/stream/categories
https://freetvchannels.online/api/v1/stream/featured
```

---

## 🚀 Deployment Instructions

### Prerequisites
- Docker & Docker Compose installed
- Domain `freetvchannels.online` pointing to your server IP
- SSH access to your server

### Step 1: Prepare Environment

```bash
cd ~/Stream

# Create .env file
cat > .env << 'EOF'
REDIS_PASSWORD=your_secure_password_here_at_least_16_chars
REDIS_DB=0
FEATURED_CHANNES=channel1,channel2,channel3
EOF
```

### Step 2: Generate SSL Certificates (First Time Only)

```bash
# Create necessary directories
mkdir -p certbot/conf certbot/www nginx/conf.d

# Generate initial SSL certificate
docker run --rm \
  -v $(pwd)/certbot/conf:/etc/letsencrypt \
  -v $(pwd)/certbot/www:/var/www/certbot \
  certbot/certbot certonly --webroot \
  -w /var/www/certbot \
  -d freetvchannels.online \
  -d www.freetvchannels.online \
  --agree-tos \
  --no-eff-email \
  -m your-email@example.com
```

### Step 3: Deploy

```bash
# Build all services
docker-compose -f docker-compose.prod.yml build

# Start all services
docker-compose -f docker-compose.prod.yml up -d

# Check status
docker-compose -f docker-compose.prod.yml ps
```

### Step 4: Verify

```bash
# Test frontend
curl https://freetvchannels.online

# Test API
curl https://freetvchannels.online/api/v1/stream/categories

# View logs
docker-compose -f docker-compose.prod.yml logs -f
```

---

## 📁 File Structure

```
Stream/
├── docker-compose.prod.yml          ← Main deployment file
├── .env                             ← Your configuration (create this)
├── .env.example                     ← Template
├── deploy.sh                        ← Automated deployment script
├── DEPLOYMENT.md                    ← Full deployment guide
├── SETUP_SUMMARY.md                 ← Quick reference
├── README_DEPLOYMENT.md             ← This file
│
├── nginx/
│   ├── nginx.conf                   ← Main nginx config
│   └── conf.d/
│       └── default.conf             ← Site configuration
│
├── certbot/
│   ├── conf/                        ← SSL certificates (auto-generated)
│   └── www/                         ← Certbot validation
│
├── front/
│   ├── Dockerfile                   ← Frontend build
│   ├── docker-compose.yml           ← Frontend only (dev)
│   ├── .env.production              ← Production env vars
│   ├── src/
│   │   └── lib/
│   │       └── channels.ts          ← MODIFIED: Uses env-based API URL
│   └── ... (other frontend files)
│
└── back/
    ├── Dockerfile                   ← Backend build
    ├── docker-compose.prod.yml      ← Backend only (legacy)
    └── ... (other backend files)
```

---

## 🔧 Common Commands

```bash
# View all logs
docker-compose -f docker-compose.prod.yml logs -f

# View specific service logs
docker-compose -f docker-compose.prod.yml logs -f front-app
docker-compose -f docker-compose.prod.yml logs -f back-app
docker-compose -f docker-compose.prod.yml logs -f nginx

# Restart all services
docker-compose -f docker-compose.prod.yml restart

# Restart specific service
docker-compose -f docker-compose.prod.yml restart front-app

# Stop all services
docker-compose -f docker-compose.prod.yml down

# Rebuild and restart
docker-compose -f docker-compose.prod.yml up -d --build

# Check service status
docker-compose -f docker-compose.prod.yml ps
```

---

## 🔐 Security Features

✅ **HTTPS/SSL**
- Automatic SSL certificate generation via Let's Encrypt
- Auto-renewal every 12 hours
- HTTP → HTTPS redirect

✅ **Security Headers**
- HSTS (HTTP Strict Transport Security)
- X-Frame-Options (Clickjacking protection)
- X-Content-Type-Options (MIME sniffing protection)
- X-XSS-Protection (XSS protection)

✅ **Network Security**
- Services communicate via internal Docker network
- No direct backend access (only through nginx)
- Redis password protected

✅ **Performance**
- Gzip compression for responses
- Static file caching (1 year)
- HTTP/2 support

---

## 🐛 Troubleshooting

### Frontend not loading
```bash
# Check frontend logs
docker-compose -f docker-compose.prod.yml logs front-app

# Check nginx logs
docker-compose -f docker-compose.prod.yml logs nginx

# Test frontend directly
curl http://localhost:3000
```

### API calls failing
```bash
# Check backend logs
docker-compose -f docker-compose.prod.yml logs back-app

# Test API directly
curl http://localhost:8090/api/v1/stream/categories

# Check Redis
docker-compose -f docker-compose.prod.yml logs back-redis
```

### SSL certificate issues
```bash
# Check certificate status
docker-compose -f docker-compose.prod.yml logs certbot

# Manually renew
docker-compose -f docker-compose.prod.yml exec certbot \
  certbot renew --force-renewal

# Check certificate expiry
docker-compose -f docker-compose.prod.yml exec certbot \
  certbot certificates
```

### Services won't start
```bash
# Check all logs
docker-compose -f docker-compose.prod.yml logs

# Verify ports are available
netstat -tlnp | grep -E ':(80|443|3000|8090)'

# Check .env file
cat .env
```

---

## 📊 What Changed in Frontend Code

### Before (Hardcoded)
```typescript
// front/src/lib/channels.ts
export const fetchFeaturedChannels = async () => {
  const response = await fetch('http://127.0.0.1:8090/api/v1/stream/featured');
  // ...
};
```

### After (Environment-based)
```typescript
// front/src/lib/channels.ts
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:8090';

export const fetchFeaturedChannels = async () => {
  const response = await fetch(`${API_BASE_URL}/v1/stream/featured`);
  // ...
};
```

**Benefits:**
- ✅ Works with any domain
- ✅ No code changes needed for different environments
- ✅ Supports both direct backend and proxied URLs
- ✅ Easy to test locally or in production

---

## 🎯 Next Steps

1. **Update `.env`** with your Redis password
2. **Generate SSL certificates** (first time only)
3. **Deploy** using docker-compose
4. **Monitor logs** to ensure everything works
5. **Test** at `https://freetvchannels.online`

---

## 📞 Support

For detailed information, see:
- `DEPLOYMENT.md` - Complete deployment guide
- `SETUP_SUMMARY.md` - Quick reference
- Docker logs: `docker-compose -f docker-compose.prod.yml logs -f`

---

**Status**: ✅ All files created and configured
**Ready to Deploy**: Yes
**Next Action**: Create `.env` and run deployment
