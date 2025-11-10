#!/bin/bash
set -e

echo "🚀 Deploying to freetvchannels.online"
echo "======================================"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if Redis is running
echo "🔍 Checking Redis..."
if docker ps | grep -q telegram_bot_redis; then
    echo "✅ Redis is already running (telegram_bot_redis)"
    REDIS_CONTAINER="telegram_bot_redis"
    REDIS_HOST="telegram_bot_redis"
elif docker ps | grep -q redis; then
    echo "✅ Redis is already running"
    REDIS_CONTAINER="redis"
    REDIS_HOST="redis"
else
    echo -e "${YELLOW}⚠️  No Redis container found, starting new one...${NC}"
    docker run -d \
        --name stream-redis \
        --restart unless-stopped \
        -p 127.0.0.1:6380:6379 \
        redis:7-alpine \
        redis-server --appendonly yes --requirepass "$REDIS_PASSWORD"
    REDIS_CONTAINER="stream-redis"
    REDIS_HOST="stream-redis"
    echo "✅ Redis started"
    sleep 5
fi

# Run the backend container
echo ""
echo "🚀 Starting application..."
docker run -d \
    --name stream-backend \
    --restart unless-stopped \
    -p 8090:8090 \
    --link "$REDIS_CONTAINER:redis" \
    -e REDIS_HOST="$REDIS_HOST" \
    -e REDIS_PORT=6379 \
    -e REDIS_PASSWORD="$REDIS_PASSWORD" \
    -e REDIS_DB="${REDIS_DB:-0}" \
    -e FEATURED_CHANNES="$FEATURED_CHANNES" \
    stream-backend


echo ""
echo "⏳ Waiting for application to start..."
sleep 5

# Check logs
echo ""
echo "📋 Application logs:"
echo "===================="
docker logs stream-backend

echo ""
echo "======================================"
echo -e "${GREEN}✅ Deployment Complete!${NC}"
echo ""
echo "🔍 Check status:"
echo "   docker ps | grep stream-backend"
echo ""
echo "📋 View logs:"
echo "   docker logs -f stream-backend"
echo ""
echo "🔄 Restart:"
echo "   docker restart stream-backend"
echo ""
echo "🛑 Stop:"
echo "   docker stop stream-backend"
echo "======================================"
