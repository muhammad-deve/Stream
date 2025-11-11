#!/bin/bash

# Stream Application Deployment Script
# This script helps deploy the entire application stack

set -e

echo "🚀 Stream Application Deployment"
echo "=================================="

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Create necessary directories
echo "📁 Creating necessary directories..."
mkdir -p certbot/conf certbot/www nginx/conf.d

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Creating template..."
    cat > .env << 'EOF'
REDIS_PASSWORD=change_me_to_secure_password
REDIS_DB=0
FEATURED_CHANNES=channel1,channel2,channel3
EOF
    echo "📝 Please edit .env file with your configuration"
    exit 1
fi

# Check if SSL certificates exist
if [ ! -f "certbot/conf/live/freetvchannels.online/fullchain.pem" ]; then
    echo "🔐 SSL certificates not found. Generating initial certificates..."
    echo "⚠️  Make sure your domain is pointing to this server's IP address"
    
    docker run --rm \
        -v "$(pwd)/certbot/conf:/etc/letsencrypt" \
        -v "$(pwd)/certbot/www:/var/www/certbot" \
        certbot/certbot certonly --webroot \
        -w /var/www/certbot \
        -d freetvchannels.online \
        -d www.freetvchannels.online \
        --agree-tos \
        --no-eff-email \
        -m admin@freetvchannels.online \
        --force-renewal
fi

# Build services
echo "🔨 Building Docker images..."
docker-compose -f docker-compose.prod.yml build

# Start services
echo "🚀 Starting services..."
docker-compose -f docker-compose.prod.yml up -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 5

# Check service status
echo "✅ Checking service status..."
docker-compose -f docker-compose.prod.yml ps

echo ""
echo "🎉 Deployment complete!"
echo ""
echo "📍 Access your application at:"
echo "   Frontend: https://freetvchannels.online"
echo "   API: https://freetvchannels.online/api/v1/stream/categories"
echo "   Admin: https://freetvchannels.online/api/admin"
echo ""
echo "📋 Useful commands:"
echo "   View logs: docker-compose -f docker-compose.prod.yml logs -f"
echo "   Stop: docker-compose -f docker-compose.prod.yml down"
echo "   Restart: docker-compose -f docker-compose.prod.yml restart"
echo ""
