# Quick Start - Redis URL Obfuscation

## What Was Implemented

✅ **Redis-based URL obfuscation system** that:
- Generates UUID tokens for each channel URL
- Stores token → URL mapping in Redis with 2-hour expiration
- Returns tokens instead of actual URLs to frontend
- Provides a `/play` endpoint to resolve tokens to actual URLs
- Automatically expires tokens after 2 hours

## Setup Steps

### 1. Install Redis

**Windows:**
```powershell
# Using Chocolatey
choco install redis-64

# Or download from: https://github.com/tporadowski/redis/releases
```

**Start Redis:**
```bash
redis-server
```

Verify Redis is running:
```bash
redis-cli ping
# Should return: PONG
```

### 2. Configure Environment

Copy `.env.example` to `.env` and update Redis settings:
```bash
cd back/app
cp .env.example .env
```

Default Redis configuration in `.env`:
```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 3. Install Dependencies

```bash
cd back/app
go mod tidy
```

This will download the Redis client library (`github.com/redis/go-redis/v9`)

### 4. Run Your Application

```bash
go run cmd/main.go serve
```

## Frontend Changes Needed

### Before (Old Flow):
```javascript
// Channel data had actual URL
const channelData = {
  url: "https://actual-stream-url.com/channel.m3u8"
};

// Play directly
videoPlayer.src = channelData.url;
```

### After (New Flow with Redis):
```javascript
// Step 1: Get channel data (URL is now a token)
const channelData = {
  url: "550e8400-e29b-41d4-a716-446655440000"  // UUID token
};

// Step 2: When user clicks play, resolve token
async function playChannel(token) {
  const response = await fetch('/api/v1/stream/play', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ token })
  });
  
  if (response.ok) {
    const { url } = await response.json();
    videoPlayer.src = url;
    videoPlayer.play();
  } else {
    // Handle expired/invalid token
    alert('Token expired. Please refresh the channel list.');
  }
}

// Step 3: Call playChannel when user clicks play
playButton.addEventListener('click', () => {
  playChannel(channelData.url);
});
```

## API Changes

### New Endpoint: `/api/v1/stream/play`

**Request:**
```json
POST /api/v1/stream/play
{
  "token": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Success Response:**
```json
{
  "url": "https://actual-stream-url.com/channel.m3u8"
}
```

**Error Response (Expired/Invalid Token):**
```json
{
  "error": "invalid or expired token"
}
```

### Modified Endpoints (Now Return Tokens Instead of URLs):
- `POST /api/v1/stream/watch` - Returns token in `url` field
- `GET /api/v1/stream/featured` - Returns tokens in `url` fields
- `GET /api/v1/stream/channel/:name` - Returns token in `url` field
- `POST /api/v1/stream/category` - Returns tokens in `url` fields
- `POST /api/v1/stream/recommend` - Returns tokens in `url` fields
- `POST /api/v1/stream/all` - Returns tokens in `url` fields
- `POST /api/v1/stream/search` - Returns tokens in `url` fields

## Testing

### Test Token Generation
```bash
curl -X POST http://localhost:8090/api/v1/stream/watch \
  -H "Content-Type: application/json" \
  -d '{"channel_id": "your_channel_id"}'
```

You should see a UUID in the `url` field instead of an actual URL.

### Test Token Resolution
```bash
curl -X POST http://localhost:8090/api/v1/stream/play \
  -H "Content-Type: application/json" \
  -d '{"token": "paste_token_from_above"}'
```

You should get the actual stream URL back.

### Test Token Expiration
Wait 2 hours, then try to use the same token. You should get an error:
```json
{
  "error": "invalid or expired token"
}
```

## Security Benefits

🔒 **Before:** Anyone inspecting the browser could see and copy actual stream URLs
🔐 **After:** 
- Browser only sees temporary UUID tokens
- Tokens expire after 2 hours
- Actual URLs are never exposed to the client
- Users cannot share permanent stream URLs

## Troubleshooting

### Error: "Failed to connect to Redis"
- Make sure Redis server is running: `redis-server`
- Check Redis connection: `redis-cli ping`
- Verify REDIS_HOST and REDIS_PORT in `.env`

### Error: "invalid or expired token"
- Token has expired (2 hours)
- Frontend should refresh channel list to get new tokens
- Implement token refresh logic in frontend

### Check Redis Contents
```bash
redis-cli
> KEYS *  # List all tokens
> TTL <token>  # Check time remaining for a token (in seconds)
> GET <token>  # Get the actual URL for a token
```

## Production Checklist

- [ ] Redis server running with persistence enabled
- [ ] Set strong REDIS_PASSWORD in production
- [ ] Consider Redis Cluster for high availability
- [ ] Implement frontend token refresh logic
- [ ] Monitor Redis memory usage
- [ ] Set up Redis backup/restore procedures
- [ ] Consider Redis over TLS for enhanced security
- [ ] Implement rate limiting on `/play` endpoint

## Need Help?

See `REDIS_URL_OBFUSCATION.md` for detailed documentation.
