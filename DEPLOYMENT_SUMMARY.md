# Deployment Summary

## ✅ Issues Resolved

### Problem: t3.micro Instance Freezing During Build
**Root Cause**: Building 10 Go microservices in parallel exhausted the 1GB RAM, causing SSH to freeze.

**Solution Implemented**:
- Build Docker images in GitHub Actions (7GB RAM available)
- Push pre-built images to GitHub Container Registry
- EC2 instance just pulls and runs images (fast, no building)

### Problem: SSH Connection Timeout
**Root Cause**: Docker build process consuming all system resources.

**Solution**: 
- Rebooted EC2 instance
- Implemented new deployment strategy that doesn't build on EC2

---

## 🚀 Current Architecture

### GitHub Actions CI/CD Pipeline
**File**: `.github/workflows/build-and-deploy.yml`

**Stage 1: Build Images (Parallel)**
- Builds all 10 microservice images simultaneously
- Uses GitHub Actions runners (plenty of resources)
- Pushes to `ghcr.io/khoipn21/study-platform/{service}:latest`
- Uses `appleboy/ssh-action@v1.2.2` (latest version)

**Stage 2: Deploy to EC2**
- Connects via SSH using `appleboy/ssh-action`
- Pulls latest code from GitHub
- Logs into GitHub Container Registry
- Pulls all pre-built Docker images
- Stops old containers
- Starts new containers
- Runs health check

### Production Docker Compose
**File**: `docker-compose.prod.yml`

**Infrastructure Services** (from Docker Hub):
- PostgreSQL 16
- Redis 7-alpine
- MinIO

**Application Services** (from GHCR):
- api-gateway: `ghcr.io/khoipn21/study-platform/api-gateway:latest`
- auth-service: `ghcr.io/khoipn21/study-platform/auth-service:latest`
- course-service: `ghcr.io/khoipn21/study-platform/course-service:latest`
- progress-service: `ghcr.io/khoipn21/study-platform/progress-service:latest`
- payment-service: `ghcr.io/khoipn21/study-platform/payment-service:latest`
- video-service: `ghcr.io/khoipn21/study-platform/video-service:latest`
- bucket-service: `ghcr.io/khoipn21/study-platform/bucket-service:latest`
- chatbot-service: `ghcr.io/khoipn21/study-platform/chatbot-service:latest`
- forum-service: `ghcr.io/khoipn21/study-platform/forum-service:latest`
- instructor-dashboard-service: `ghcr.io/khoipn21/study-platform/instructor-dashboard-service:latest`

---

## 📋 Infrastructure Details

**EC2 Instance**:
- Type: t3.micro (Free Tier)
- Region: ap-southeast-2
- Instance ID: i-0eff58566b8170bf1
- Public IP: 3.104.45.98 (Elastic IP)
- OS: Ubuntu 22.04 LTS
- Resources: 1 vCPU, 914MB RAM

**Nginx Configuration**:
- Proxy: `study.khoipn.id.vn` → `localhost:8080`
- Status: Running ✅
- Config: `/etc/nginx/sites-available/study-platform`

**SSH Access**:
- Key: `credentials/study-key.pem` (chmod 400)
- Command: `ssh -i credentials/study-key.pem ubuntu@3.104.45.98`
- Status: Working ✅

**GitHub Repository Secrets**:
- `EC2_HOST`: 3.104.45.98
- `EC2_USERNAME`: ubuntu
- `EC2_SSH_KEY`: Private key content
- `DEPLOY_TOKEN`: GitHub PAT for cloning
- `AWS_ACCESS_KEY_ID` & `AWS_SECRET_ACCESS_KEY`: For future AWS operations
- `GITHUB_TOKEN`: Automatically provided for GHCR authentication

---

## 🔧 Deployment Process

### Automatic Deployment
**Trigger**: Push to `master` branch

**Steps**:
1. GitHub Actions builds all images (~5-10 minutes)
2. Pushes images to GitHub Container Registry
3. SSH to EC2
4. Pulls latest code
5. Pulls pre-built Docker images (~2 minutes)
6. Restarts services
7. Health check validation
8. Cleanup old images

**Total Time**: ~10-15 minutes

### Manual Deployment
```bash
# Trigger via GitHub UI
https://github.com/khoipn21/study-platform/actions/workflows/build-and-deploy.yml

# Or via CLI
gh workflow run build-and-deploy.yml --repo khoipn21/study-platform
```

---

## 🧪 Testing & Verification

### SSH Connection Test
```bash
ssh -i credentials/study-key.pem ubuntu@3.104.45.98
```

### Health Check
```bash
# Direct IP
curl http://3.104.45.98:8080/health

# Via domain (after DNS setup)
curl http://study.khoipn.id.vn/health
```

### Container Status
```bash
ssh -i credentials/study-key.pem ubuntu@3.104.45.98 \
  "cd ~/study-platform && docker compose -f docker-compose.prod.yml ps"
```

### Logs
```bash
ssh -i credentials/study-key.pem ubuntu@3.104.45.98 \
  "cd ~/study-platform && docker compose -f docker-compose.prod.yml logs -f --tail=100"
```

---

## 📝 Next Steps

### 1. DNS Configuration
Point your domain to the EC2 instance:
```
Type: A Record
Name: study.khoipn.id.vn
Value: 3.104.45.98
TTL: 300
```

### 2. SSL Certificate (Optional but Recommended)
```bash
ssh -i credentials/study-key.pem ubuntu@3.104.45.98
sudo apt-get install -y certbot python3-certbot-nginx
sudo certbot --nginx -d study.khoipn.id.vn
```

### 3. Monitor First Deployment
```bash
# Watch workflow progress
gh run watch --repo khoipn21/study-platform

# Check deployment logs
ssh -i credentials/study-key.pem ubuntu@3.104.45.98 \
  "cd ~/study-platform && docker compose -f docker-compose.prod.yml logs -f"
```

### 4. Verify All Services
```bash
# Test API Gateway
curl http://3.104.45.98:8080/health

# Test individual services (if exposed)
curl http://3.104.45.98:8084/health  # Video Service
curl http://3.104.45.98:8085/health  # Bucket Service
# etc...
```

---

## 🎯 Benefits of This Approach

1. **No More Freezing**: Building happens on GitHub's servers, not your t3.micro
2. **Faster Deployments**: Pull pre-built images instead of building (2 min vs 30+ min)
3. **Reliable SSH**: EC2 resources not exhausted during deployment
4. **Automatic CI/CD**: Every push deploys automatically
5. **Health Checks**: Automatic validation after deployment
6. **Cost Effective**: Uses free tier EC2 and free GitHub Actions/GHCR for public repos

---

## 📚 Files Modified

- `.github/workflows/build-and-deploy.yml` - New workflow with image building
- `docker-compose.prod.yml` - Production compose with pre-built images
- `DEPLOYMENT_STATUS.md` - Initial deployment documentation
- `DEPLOYMENT_SUMMARY.md` - This file

## 🔗 Links

- **GitHub Repository**: https://github.com/khoipn21/study-platform
- **GitHub Actions**: https://github.com/khoipn21/study-platform/actions
- **API Endpoint**: http://3.104.45.98:8080 (or http://study.khoipn.id.vn after DNS)
- **SSH Command**: `ssh -i credentials/study-key.pem ubuntu@3.104.45.98`
