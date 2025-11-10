#!/bin/bash
set -e

echo "🚀 Deployment Script for freetvchannels.online"
echo "=============================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env.prod exists
if [ ! -f .env.prod ]; then
    echo -e "${RED}❌ .env.prod file not found!${NC}"
    echo "Please create .env.prod from .env.prod.example"
    echo "Run: cp .env.prod.example .env.prod"
    echo "Then edit .env.prod with your configuration"
    exit 1
fi

# Load environment variables
export $(cat .env.prod | grep -v '^#' | xargs)

echo "📋 Checking prerequisites..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Docker installed${NC}"

# Check if Docker Compose is installed
if ! docker compose version &> /dev/null; then
    echo -e "${RED}❌ Docker Compose is not installed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Docker Compose installed${NC}"

# Create necessary directories
echo ""
echo "📁 Creating directories..."
mkdir -p nginx/conf.d certbot/conf certbot/www
echo -e "${GREEN}✓ Directories created${NC}"

# Check if initial or full deployment
if [ ! -f "certbot/conf/live/freetvchannels.online/fullchain.pem" ]; then
    echo ""
    echo -e "${YELLOW}⚠️  SSL certificates not found${NC}"
    echo "This appears to be the first deployment."
    echo ""
    read -p "Do you want to obtain SSL certificates now? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "🔐 Starting initial deployment (HTTP only)..."
        
        # Use initial config
        if [ ! -f "nginx/conf.d/initial.conf" ]; then
            cat > nginx/conf.d/initial.conf << 'EOF'
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
EOF
        fi
        
        # Start services
        docker compose -f docker-compose.prod.yml up -d
        
        echo ""
        echo "⏳ Waiting 10 seconds for services to start..."
        sleep 10
        
        echo ""
        read -p "Enter your email for Let's Encrypt notifications: " EMAIL
        
        echo "🔐 Obtaining SSL certificate..."
        docker compose -f docker-compose.prod.yml run --rm certbot certonly \
          --webroot \
          --webroot-path=/var/www/certbot \
          --email "$EMAIL" \
          --agree-tos \
          --no-eff-email \
          -d freetvchannels.online \
          -d www.freetvchannels.online
        
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ SSL certificate obtained${NC}"
            rm nginx/conf.d/initial.conf
            echo -e "${GREEN}✓ Switched to HTTPS configuration${NC}"
        else
            echo -e "${RED}❌ Failed to obtain SSL certificate${NC}"
            exit 1
        fi
    fi
fi

echo ""
echo "🚀 Starting/Updating services..."
docker compose -f docker-compose.prod.yml up -d --build

echo ""
echo "⏳ Waiting for services to be ready..."
sleep 5

echo ""
echo "📊 Service Status:"
docker compose -f docker-compose.prod.yml ps

echo ""
echo "🔍 Testing services..."

# Test Redis
echo -n "Redis: "
if docker compose -f docker-compose.prod.yml exec -T redis redis-cli -a "$REDIS_PASSWORD" ping > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Running${NC}"
else
    echo -e "${RED}✗ Not responding${NC}"
fi

# Test App
echo -n "Backend: "
if docker compose -f docker-compose.prod.yml exec -T app wget -q -O- http://localhost:8090/api/v1/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Running${NC}"
else
    echo -e "${RED}✗ Not responding${NC}"
fi

# Test Nginx
echo -n "Nginx: "
if docker compose -f docker-compose.prod.yml exec -T nginx nginx -t > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Running${NC}"
else
    echo -e "${RED}✗ Configuration error${NC}"
fi

echo ""
echo "=============================================="
echo -e "${GREEN}✅ Deployment Complete!${NC}"
echo ""
echo "🌐 Your site is available at:"
echo "   https://freetvchannels.online"
echo ""
echo "📝 Useful commands:"
echo "   View logs:     docker compose -f docker-compose.prod.yml logs -f"
echo "   Stop services: docker compose -f docker-compose.prod.yml down"
echo "   Restart:       docker compose -f docker-compose.prod.yml restart"
echo ""
echo "📚 For more information, see DEPLOYMENT_GUIDE.md"
echo "=============================================="
