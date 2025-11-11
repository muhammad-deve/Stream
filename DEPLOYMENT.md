# Deployment Guide for freetvchannels.online

This guide explains how to deploy the entire Stream application (frontend + backend) using Docker.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    freetvchannels.online                    │
│                      (Nginx Reverse Proxy)                  │
├──────────────────────┬──────────────────────────────────────┤
│                      │                                      │
│   Frontend (React)   │        Backend (PocketBase)          │
│   Port: 3000         │        Port: 8090                    │
│                      │                                      │
│   - Serves UI        │   - API Endpoints                    │
│   - Handles routing  │   - Stream management                │
│                      │   - Redis caching                    │
└──────────────────────┴──────────────────────────────────────┘
```

## Prerequisites

- Docker & Docker Compose installed
- Domain: `freetvchannels.online` pointing to your server IP
- SSL certificates (auto-managed by Let's Encrypt/Certbot)
- Environment variables configured

## Setup Instructions

### 1. Clone/Navigate to Project Root

```bash
cd ~/Stream
```

### 2. Create Environment File

Create `.env` file in the root directory:

```bash
cat > .env << 'EOF'
REDIS_PASSWORD=your_secure_redis_password_here
REDIS_DB=0
FEATURED_CHANNES=channel1,channel2,channel3
EOF
```

### 3. Create Necessary Directories

```bash
mkdir -p certbot/conf certbot/www nginx/conf.d
```

### 4. Initial SSL Certificate Setup (First Time Only)

Before running the full stack, you need to generate SSL certificates:

```bash
# Stop any existing nginx
docker-compose -f docker-compose.prod.yml down

# Run certbot to generate initial certificate
docker run --rm -v $(pwd)/certbot/conf:/etc/letsencrypt -v $(pwd)/certbot/www:/var/www/certbot certbot/certbot certonly --webroot -w /var/www/certbot -d freetvchannels.online -d www.freetvchannels.online --agree-tos --no-eff-email -m your-email@example.com
```

### 5. Build and Deploy

```bash
# Build all services
docker-compose -f docker-compose.prod.yml build

# Start all services
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f
```

### 6. Verify Deployment

```bash
# Check running containers
docker ps

# Test frontend
curl https://freetvchannels.online

# Test API
curl https://freetvchannels.online/api/v1/stream/categories
```

## File Structure

```
Stream/
├── docker-compose.prod.yml      # Main compose file
├── .env                         # Environment variables
├── nginx/
│   ├── nginx.conf              # Nginx main config
│   └── conf.d/
│       └── default.conf        # Site configuration
├── certbot/
│   ├── conf/                   # SSL certificates
│   └── www/                    # Certbot validation
├── front/
│   ├── Dockerfile              # Frontend build
│   ├── docker-compose.yml      # Frontend only (dev)
│   ├── .env.production         # Frontend env vars
│   └── src/
│       └── lib/
│           └── channels.ts     # API configuration
└── back/
    ├── Dockerfile              # Backend build
    └── docker-compose.prod.yml # Backend only (legacy)
```

## API Configuration

The frontend automatically uses the correct API base URL based on environment:

- **Development**: `http://localhost:8090`
- **Production**: `https://freetvchannels.online/api`

This is configured in:
- `front/.env.production` - Sets `VITE_API_BASE_URL`
- `front/src/lib/channels.ts` - Uses the environment variable

## Nginx Configuration

The nginx reverse proxy:
1. **Redirects HTTP → HTTPS** for security
2. **Serves frontend** from `http://front-app:3000`
3. **Proxies API requests** from `/api/*` to `http://back-app:8090/api/*`
4. **Manages SSL certificates** with auto-renewal

### URL Routing

| URL | Destination |
|-----|-------------|
| `https://freetvchannels.online/` | Frontend (React app) |
| `https://freetvchannels.online/api/v1/stream/*` | Backend API |
| `https://freetvchannels.online/api/admin/*` | PocketBase Admin |

## Common Commands

```bash
# View logs
docker-compose -f docker-compose.prod.yml logs -f

# View specific service logs
docker-compose -f docker-compose.prod.yml logs -f front-app
docker-compose -f docker-compose.prod.yml logs -f back-app

# Restart services
docker-compose -f docker-compose.prod.yml restart

# Stop all services
docker-compose -f docker-compose.prod.yml down

# Rebuild and restart
docker-compose -f docker-compose.prod.yml up -d --build

# Check service status
docker-compose -f docker-compose.prod.yml ps
```

## Troubleshooting

### Frontend not loading
```bash
# Check frontend logs
docker-compose -f docker-compose.prod.yml logs front-app

# Verify nginx is proxying correctly
docker-compose -f docker-compose.prod.yml logs nginx
```

### API calls failing
```bash
# Check backend logs
docker-compose -f docker-compose.prod.yml logs back-app

# Test API directly
curl http://localhost:8090/api/v1/stream/categories
```

### SSL certificate issues
```bash
# Check certificate status
docker-compose -f docker-compose.prod.yml logs certbot

# Manually renew certificate
docker-compose -f docker-compose.prod.yml exec certbot certbot renew --force-renewal
```

### Redis connection issues
```bash
# Check Redis logs
docker-compose -f docker-compose.prod.yml logs back-redis

# Test Redis connection
docker-compose -f docker-compose.prod.yml exec back-redis redis-cli -a $REDIS_PASSWORD ping
```

## Performance Optimization

The nginx configuration includes:
- **Gzip compression** for text/JSON responses
- **Static file caching** (1 year expiry for assets)
- **HTTP/2 support** for faster connections
- **Security headers** (HSTS, X-Frame-Options, etc.)

## Monitoring

Check container health:
```bash
docker-compose -f docker-compose.prod.yml ps

# Should show all services as "Up"
```

## Updating the Application

### Update Frontend
```bash
cd front
git pull
docker-compose -f docker-compose.prod.yml up -d --build front-app
```

### Update Backend
```bash
cd back
git pull
docker-compose -f docker-compose.prod.yml up -d --build back-app
```

## Security Notes

- All traffic is redirected to HTTPS
- SSL certificates auto-renew via Certbot
- Redis password is required (set in `.env`)
- API requests are proxied through nginx (no direct backend access)
- Security headers are configured in nginx

## Support

For issues or questions, check:
1. Container logs: `docker-compose logs -f`
2. Nginx configuration: `nginx/conf.d/default.conf`
3. Frontend environment: `front/.env.production`
4. Backend environment: `.env`
