# Deployment Status

## ✅ Completed Tasks

### AWS Infrastructure
- **EC2 Instance**: t3.micro (Free Tier) in ap-southeast-2
- **Instance ID**: i-0eff58566b8170bf1
- **Elastic IP**: 3.104.45.98
- **SSH Key**: `credentials/study-key.pem` (saved locally)
- **Security Group**: Configured for ports 22, 80, 443, 8080-8089
- **Old Resources**: All previous instances terminated and elastic IPs released

### Server Configuration
- **OS**: Ubuntu 22.04 LTS
- **Docker**: v28.5.1 installed
- **Docker Compose**: v2.40.1 installed
- **Nginx**: Configured as reverse proxy for `study.khoipn.id.vn` → `localhost:8080`
- **Repository**: Cloned with `.env` file deployed
- **SSH**: Working and accessible

### GitHub CI/CD Pipeline
- **Workflow**: `.github/workflows/deploy.yml`
- **Action Plugin**: `appleboy/ssh-action@v1.2.2` (latest version)
- **Secrets Configured**:
  - `EC2_HOST`: 3.104.45.98
  - `EC2_USERNAME`: ubuntu
  - `EC2_SSH_KEY`: Private key content
  - `DEPLOY_TOKEN`: GitHub PAT for cloning
  - `AWS_ACCESS_KEY_ID` & `AWS_SECRET_ACCESS_KEY`: For future AWS operations

### Deployment Features
- Automatic deployment on push to master branch
- Manual deployment via workflow_dispatch
- 30-minute timeout for builds
- Reduced parallelism (`COMPOSE_PARALLEL_LIMIT=2`) optimized for t3.micro
- Automatic health check after deployment
- Docker image cleanup after deployment

## ⚠️ Current Issue

**Build Performance on t3.micro**:
The t3.micro instance (1 vCPU, 1GB RAM) struggles to build all 10 microservices in parallel. The build process is currently running but may take 15-30 minutes to complete.

### Build Status
- Initial build started manually on EC2
- Services: postgres, redis, minio (infrastructure) started
- Application services: Building in progress (slow due to limited resources)

## 🔧 Solutions

### Option 1: Wait for Current Build (Recommended for Free Tier)
```bash
# SSH to server and monitor build progress
ssh -i credentials/study-key.pem ubuntu@3.104.45.98

# Check build status
cd ~/study-platform
docker compose ps

# View build logs
docker compose logs -f --tail=100

# Once completed, test API
curl http://3.104.45.98:8080/health
```

### Option 2: Build Images in GitHub Actions (Faster Deployments)
Modify workflow to build Docker images in GitHub Actions (which has more resources) and push to Docker Hub or AWS ECR, then just pull on EC2.

### Option 3: Temporarily Use Larger Instance
During initial build, temporarily switch to t3.small or t3.medium, then downgrade back to t3.micro after images are built.

## 📋 Next Steps

### Immediate
1. **Wait for build to complete** (~15-30 minutes)
2. **Point DNS**: Add A record for `study.khoipn.id.vn` → `3.104.45.98`
3. **Test deployment**:
   ```bash
   curl http://3.104.45.98:8080/health
   curl http://study.khoipn.id.vn/health
   ```

### Optional Improvements
1. **SSL Certificate**: Add Let's Encrypt for HTTPS
2. **Docker Registry**: Pre-build images to speed up deployments
3. **Health Checks**: Add liveness/readiness probes
4. **Monitoring**: Set up CloudWatch or Datadog
5. **Backup**: Configure automated database backups

## 🚀 Deployment Workflow

Every push to master branch automatically:
1. Connects to EC2 via SSH
2. Pulls latest code from GitHub
3. Stops running containers
4. Builds and starts services with reduced parallelism
5. Waits for services to stabilize
6. Tests health endpoint
7. Cleans up old Docker images

## 📝 Access Information

**Server**: `ssh -i credentials/study-key.pem ubuntu@3.104.45.98`
**API**: `http://3.104.45.98:8080` or `http://study.khoipn.id.vn` (after DNS)
**GitHub Actions**: https://github.com/khoipn21/study-platform/actions

## ⚙️ Environment Variables

All production environment variables are deployed from `.env` file including:
- Database connections (Supabase)
- Redis configuration
- AWS S3 credentials
- API keys (OpenAI, Cloudflare, Stripe, Lemon Squeezy)
- Service URLs and ports

**Note**: Domain CORS has been updated to include `study.khoipn.id.vn`
