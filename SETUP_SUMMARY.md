# Setup Summary - Frontend + Backend Deployment

## What Was Created

### 1. **Frontend Configuration**
- ✅ `front/.env.production` - Production API URL configuration
- ✅ `front/Dockerfile` - Multi-stage build for React app
- ✅ `front/docker-compose.yml` - Standalone frontend compose (for dev)
- ✅ `front/src/lib/channels.ts` - Updated to use environment-based API URLs

### 2. **Docker Compose Stack** (Root Level)
- ✅ `docker-compose.prod.yml` - Complete production stack with:
  - Frontend (React) on port 3000
  - Backend (PocketBase) on port 8090
  - Redis cache
  - Nginx reverse proxy
  - Certbot SSL management

### 3. **Nginx Configuration**
- ✅ `nginx/nginx.conf` - Main nginx configuration
- ✅ `nginx/conf.d/default.conf` - Site configuration with:
  - HTTP → HTTPS redirect
  - Frontend proxy to React app
  - API proxy to backend
  - SSL/TLS configuration
  - Security headers
  - Gzip compression
  - Static file caching

### 4. **Deployment Helpers**
- ✅ `DEPLOYMENT.md` - Complete deployment guide
- ✅ `deploy.sh` - Automated deployment script
- ✅ `SETUP_SUMMARY.md` - This file

## How It Works

### API URL Resolution

**Frontend** (`front/src/lib/channels.ts`):
```javascript
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:8090';
```

**Environment Variables**:
- Development: `VITE_API_BASE_URL=http://localhost:8090`
- Production: `VITE_API_BASE_URL=https://freetvchannels.online/api`

### Request Flow

```
User Browser
    ↓
https://freetvchannels.online
    ↓
Nginx (Port 443)
    ├─ / → Frontend (React app on :3000)
    └─ /api/* → Backend (PocketBase on :8090)
    ↓
Frontend loads and makes API calls to:
https://freetvchannels.online/api/v1/stream/*
    ↓
Nginx proxies to backend:
http://back-app:8090/api/v1/stream/*
```

## Quick Start

### Step 1: Prepare Environment
```bash
cd ~/Stream
cat > .env << 'EOF'
REDIS_PASSWORD=your_secure_password
REDIS_DB=0
FEATURED_CHANNES=channel1,channel2,channel3
EOF
```

### Step 2: Generate SSL Certificates (First Time)
```bash
mkdir -p certbot/conf certbot/www nginx/conf.d

docker run --rm \
  -v $(pwd)/certbot/conf:/etc/letsencrypt \
  -v $(pwd)/certbot/www:/var/www/certbot \
  certbot/certbot certonly --webroot \
  -w /var/www/certbot \
  -d freetvchannels.online \
  -d www.freetvchannels.online \
  --agree-tos --no-eff-email \
  -m your-email@example.com
```

### Step 3: Deploy
```bash
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d
```

### Step 4: Verify
```bash
# Check services
docker-compose -f docker-compose.prod.yml ps

# Test frontend
curl https://freetvchannels.online

# Test API
curl https://freetvchannels.online/api/v1/stream/categories
```

## Key Changes Made

### Frontend Code Changes
1. **`front/src/lib/channels.ts`** - All API calls now use `API_BASE_URL` variable
   - Supports environment-based configuration
   - Falls back to localhost for development
   - Works with both direct backend and proxied URLs

### New Files Created
- Frontend Dockerfile with production build
- Docker Compose for full stack
- Nginx configuration for reverse proxy
- Deployment documentation and scripts

## Testing the Setup

### Local Development
```bash
cd front
npm install
VITE_API_BASE_URL=http://localhost:8090 npm run dev
```

### Production Deployment
```bash
# From root directory
docker-compose -f docker-compose.prod.yml up -d

# Access at https://freetvchannels.online
```

## Important Notes

1. **Domain Setup**: Ensure `freetvchannels.online` DNS points to your server IP
2. **SSL Certificates**: Auto-renewed by Certbot (runs in background)
3. **API Proxying**: All `/api/*` requests go through nginx to backend
4. **Network**: Services communicate via internal Docker network
5. **Redis**: Secured with password from `.env`

## File Structure After Setup

```
Stream/
├── docker-compose.prod.yml
├── deploy.sh
├── DEPLOYMENT.md
├── SETUP_SUMMARY.md
├── .env
├── nginx/
│   ├── nginx.conf
│   └── conf.d/default.conf
├── certbot/
│   ├── conf/
│   └── www/
├── front/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   ├── .env.production
│   └── src/lib/channels.ts (modified)
└── back/
    └── (existing backend files)
```

## Troubleshooting

### Frontend shows "Cannot reach API"
- Check nginx logs: `docker-compose -f docker-compose.prod.yml logs nginx`
- Verify backend is running: `docker-compose -f docker-compose.prod.yml ps`
- Check API URL in browser console

### SSL certificate errors
- Ensure domain DNS is pointing to server
- Check certbot logs: `docker-compose -f docker-compose.prod.yml logs certbot`
- Manually renew: `docker-compose -f docker-compose.prod.yml exec certbot certbot renew --force-renewal`

### Services won't start
- Check logs: `docker-compose -f docker-compose.prod.yml logs`
- Ensure ports 80, 443 are available
- Check .env file is properly configured

## Next Steps

1. ✅ Review and update `.env` with your settings
2. ✅ Generate SSL certificates
3. ✅ Deploy using `docker-compose up -d`
4. ✅ Monitor logs: `docker-compose logs -f`
5. ✅ Access at `https://freetvchannels.online`

## Support Commands

```bash
# View all logs
docker-compose -f docker-compose.prod.yml logs -f

# View specific service
docker-compose -f docker-compose.prod.yml logs -f front-app

# Restart services
docker-compose -f docker-compose.prod.yml restart

# Stop all
docker-compose -f docker-compose.prod.yml down

# Rebuild and restart
docker-compose -f docker-compose.prod.yml up -d --build
```

---

**Status**: ✅ All files created and configured
**Next Action**: Deploy using the docker-compose command
