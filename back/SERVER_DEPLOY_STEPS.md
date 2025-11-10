# 🚀 Quick Server Deployment

## On Your Local Machine (Windows)

```bash
# 1. Commit and push changes
cd c:\Users\user\Desktop\Stream
git add .
git commit -m "Fix Docker config reading and add deploy script"
git push origin master
```

## On Your Server (root@vmi2877544)

```bash
# 1. Navigate to project
cd ~/Stream/back

# 2. Pull latest changes
git pull origin master

# 3. Create .env.prod file
nano .env.prod
```

**Add this to .env.prod:**
```env
REDIS_PASSWORD=your_strong_password_here
REDIS_DB=0
FEATURED_CHANNES=9m3fy1v50tlelwe,75l42yzl07l4948,3v539ya9sxjr8b6,fk9p603ke1k9p01,9nytgzajem68r9p,234c5688eyn2n29
```

Press `Ctrl+X`, then `Y`, then `Enter` to save.

**Generate a strong password:**
```bash
openssl rand -base64 32
```

```bash
# 4. Make deploy script executable
chmod +x deploy-server.sh

# 5. Deploy!
./deploy-server.sh
```

---

## What the Script Does

1. ✅ Checks if `.env.prod` exists
2. ✅ Loads environment variables
3. ✅ Stops old container (if exists)
4. ✅ Uses your existing Redis container (`telegram_bot_redis`)
5. ✅ Builds the application
6. ✅ Starts the container with correct environment variables
7. ✅ Shows logs

---

## After Deployment

### Check if running:
```bash
docker ps | grep stream-backend
```

### View logs:
```bash
docker logs -f stream-backend
```

You should see:
```
✅ Environment variables loaded
Configuration:
  REDIS_HOST: redis
  REDIS_PORT: 6379
  REDIS_DB: 0
  REDIS_PASSWORD: ab****xy
  FEATURED_CHANNES: ✓
✅ Successfully connected to Redis
```

### Test API:
```bash
curl http://localhost:8090/api/v1/health
curl http://localhost:8090/api/v1/stream/featured
```

### Restart:
```bash
docker restart stream-backend
```

### Update after code changes:
```bash
cd ~/Stream/back
git pull
./deploy-server.sh
```

---

## ✅ Expected Result

Container should start successfully with:
- ✅ Environment variables loaded from `.env.prod`
- ✅ Connected to Redis (your existing `telegram_bot_redis`)
- ✅ App running on port 8090
- ✅ No more "connection refused" errors

---

## 🐛 Troubleshooting

### If container exits immediately:
```bash
docker logs stream-backend
```

### If Redis connection fails:
```bash
# Check if Redis is running
docker ps | grep redis

# Test Redis connection
docker exec -it telegram_bot_redis redis-cli -a YOUR_PASSWORD ping
```

### If environment variables not loaded:
```bash
# Check .env.prod content
cat .env.prod

# Make sure no spaces around =
# Wrong: REDIS_PASSWORD = password
# Right: REDIS_PASSWORD=password
```

### Start fresh:
```bash
docker stop stream-backend
docker rm stream-backend
./deploy-server.sh
```

---

## 📞 Quick Commands Reference

```bash
# Deploy/Update
./deploy-server.sh

# View logs
docker logs -f stream-backend

# Restart
docker restart stream-backend

# Stop
docker stop stream-backend

# Check status
docker ps | grep stream

# Test API
curl http://localhost:8090/api/v1/stream/featured
```
