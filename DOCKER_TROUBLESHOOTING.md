# Docker Build Troubleshooting

## 🔧 **Fixed Issues**

### **"Unknown blob" Error**
- **Problem**: Docker registry cache corruption during build/push
- **Solutions Applied**:
  1. **Updated Backend Dockerfile**:
     - Switched from `golang:1.23-bullseye` to `golang:1.23-alpine` (smaller, more reliable)
     - Removed complex cache mounts that can cause issues
     - Added proper `go.sum` copy
     - Used standard Go build flags
     - Added ca-certificates to Alpine image

  2. **Enhanced CI/CD Workflow**:
     - Added `docker/setup-buildx-action@v3` for better buildx setup
     - Added GitHub Actions cache (`type=gha`) for more reliable caching
     - Improved error handling and build process

## 🚀 **Improvements Made**

### **Backend Dockerfile Changes**
```dockerfile
# Before (problematic)
FROM golang:1.23-bullseye AS build
COPY ./app/go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
RUN go build -ldflags="-linkmode external -extldflags -static" -tags netgo -o /app/main ./cmd/main.go

# After (fixed)
FROM golang:1.23-alpine AS build
COPY app/go.mod app/go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go
```

### **CI/CD Workflow Enhancements**
```yaml
# Added Buildx setup
- name: Set up Docker Buildx
  uses: docker/setup-buildx-action@v3

# Added reliable caching
cache-from: type=gha
cache-to: type=gha,mode=max
```

## 📋 **If Issues Persist**

### **Manual Fixes**
1. **Clear Docker Cache**:
   ```bash
   docker builder prune -a
   docker system prune -a
   ```

2. **Restart Docker Service**:
   ```bash
   sudo systemctl restart docker
   ```

3. **Re-authenticate with Registry**:
   ```bash
   docker logout ghcr.io
   echo $GITHUB_TOKEN | docker login ghcr.io -u $GITHUB_ACTOR --password-stdin
   ```

### **GitHub Actions Debugging**
1. **Enable Debug Logging**:
   ```yaml
   env:
     DOCKER_BUILDKIT: 1
     BUILDKIT_PROGRESS: plain
   ```

2. **Add Retry Logic**:
   ```yaml
   - name: Build and push backend image
     uses: docker/build-push-action@v5
     continue-on-error: true
     # ... rest of config
   ```

## ✅ **Expected Results**

With these fixes, your Docker builds should:
- ✅ Build reliably without "unknown blob" errors
- ✅ Use efficient caching with GitHub Actions
- ✅ Push successfully to GitHub Container Registry
- ✅ Deploy correctly to production

## 🔄 **Next Steps**

1. Commit and push the Dockerfile changes
2. Monitor the GitHub Actions build
3. If successful, the deployment should proceed automatically
4. Check the deployment at https://freetvchannels.online

## 📞 **Additional Support**

If issues continue:
1. Check GitHub Actions logs for specific error messages
2. Verify registry permissions in GitHub repository settings
3. Ensure Docker buildx is properly configured
4. Consider using GitHub Container Registry package limits
