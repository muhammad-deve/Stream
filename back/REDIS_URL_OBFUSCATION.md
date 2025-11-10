# Redis URL Obfuscation System

## Overview
This system implements a Redis-based URL obfuscation mechanism to hide actual live stream channel URLs from browser inspection. URLs are replaced with temporary UUID tokens that expire after 2 hours.

## How It Works

### 1. Token Generation
When a user requests channel information (via any stream endpoint), the system:
- Retrieves the actual channel URL from the database
- Generates a unique UUID token
- Stores the mapping `token -> actual_URL` in Redis with a 2-hour TTL (Time To Live)
- Returns the token instead of the actual URL to the frontend

### 2. Token Resolution (Playing Stream)
When a user clicks play:
- Frontend sends the token to `/api/v1/stream/play` endpoint
- Backend validates the token in Redis
- If valid and not expired, returns the actual URL
- If invalid or expired, returns an error

### 3. Automatic Expiration
- Tokens automatically expire after 2 hours
- Redis handles cleanup automatically via TTL
- No manual deletion needed

## API Endpoints

### Get Channel Information (Returns Token)
**POST** `/api/v1/stream/watch`
```json
{
  "channel_id": "channel_id_here"
}
```

**Response:**
```json
{
  "channel": "channel_name",
  "title": "Channel Title",
  "url": "550e8400-e29b-41d4-a716-446655440000",  // UUID token, not actual URL
  "quality": "1080p",
  "logo": {...},
  "category": {...},
  "country": {...},
  "language": {...}
}
```

### Play Stream (Resolve Token to URL)
**POST** `/api/v1/stream/play`
```json
{
  "token": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response (Success):**
```json
{
  "url": "https://actual-stream-url.com/channel.m3u8"
}
```

**Response (Error - Invalid/Expired Token):**
```json
{
  "error": "invalid or expired token"
}
```

## Configuration

Add the following environment variables to your `.env` file:

```env
# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### Default Values
- `REDIS_HOST`: localhost
- `REDIS_PORT`: 6379
- `REDIS_PASSWORD`: (empty string)
- `REDIS_DB`: 0

## Frontend Integration

### 1. Get Channel List
When fetching channels (featured, category, search, etc.), the `url` field will contain a UUID token instead of the actual stream URL.

### 2. Play Channel
When user clicks play:
```javascript
// Step 1: User clicks play, you have the token from channel data
const token = channelData.url; // This is the UUID token

// Step 2: Request actual URL from backend
const response = await fetch('/api/v1/stream/play', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ token })
});

// Step 3: Get actual URL and play
const { url } = await response.json();
videoPlayer.src = url;
videoPlayer.play();
```

## Security Benefits

1. **Hidden URLs**: Actual stream URLs are never exposed in browser DevTools or network inspection
2. **Time-Limited Access**: Tokens expire after 2 hours, limiting unauthorized access
3. **Single-Use Tokens**: Each channel request generates a new token
4. **No URL Sharing**: Users cannot share actual stream URLs as they only have temporary tokens

## Installation & Setup

### 1. Install Redis
```bash
# Windows (using Chocolatey)
choco install redis-64

# Or download from: https://redis.io/download

# Start Redis server
redis-server
```

### 2. Install Dependencies
```bash
cd back/app
go mod tidy
```

### 3. Configure Environment
Create or update `.env` file:
```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 4. Run Application
```bash
go run cmd/main.go
```

## Testing

### Test Token Generation
```bash
curl -X POST http://localhost:8090/api/v1/stream/watch \
  -H "Content-Type: application/json" \
  -d '{"channel_id": "your_channel_id"}'
```

### Test Token Resolution
```bash
curl -X POST http://localhost:8090/api/v1/stream/play \
  -H "Content-Type: application/json" \
  -d '{"token": "your_token_here"}'
```

## Monitoring

### Check Redis Connection
```bash
redis-cli ping
# Should return: PONG
```

### View Stored Tokens
```bash
redis-cli
> KEYS *
> TTL <token>  # Check remaining time for a token
> GET <token>  # Get URL for a token
```

## Troubleshooting

### Redis Connection Failed
- Ensure Redis server is running: `redis-cli ping`
- Check Redis host and port in `.env`
- Verify firewall settings

### Token Expired Error
- Tokens expire after 2 hours
- User needs to refresh channel list to get new tokens
- Frontend should handle token expiration gracefully

### Performance Considerations
- Redis is in-memory, very fast
- Each token ~100 bytes in Redis
- Can handle millions of concurrent tokens
- Consider Redis persistence settings for production

## Production Recommendations

1. **Redis Persistence**: Enable RDB or AOF for data persistence
2. **Redis Cluster**: Use Redis Cluster for high availability
3. **Connection Pooling**: Already handled by go-redis client
4. **Monitoring**: Use Redis monitoring tools (RedisInsight, Prometheus)
5. **Security**: Use Redis AUTH (password) in production
6. **Network**: Consider Redis over TLS for added security
