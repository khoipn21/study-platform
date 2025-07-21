# Study Platform Deployment Guide

## Overview

This guide covers deploying the Study Platform in various environments, from development to production.

## Table of Contents

1. [Development Environment](#development-environment)
2. [Production Environment](#production-environment)
3. [Docker Deployment](#docker-deployment)
4. [Kubernetes Deployment](#kubernetes-deployment)
5. [Environment Variables](#environment-variables)
6. [Database Setup](#database-setup)
7. [Monitoring & Logging](#monitoring--logging)
8. [Security Configuration](#security-configuration)
9. [Performance Tuning](#performance-tuning)
10. [Troubleshooting](#troubleshooting)

## Development Environment

### Prerequisites

- Docker & Docker Compose
- Go 1.23+
- PostgreSQL 13+
- Redis 7+
- Git

### Quick Start

```bash
# Clone repository
git clone <repository-url>
cd Study-Platform

# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

### Development with Live Reload

```bash
# Start only dependencies
docker-compose up -d postgres redis minio

# Run services locally
cd auth-service && go run cmd/main.go &
cd course-service && go run cmd/main.go &
cd progress-service && go run cmd/main.go &
cd api-gateway && go run cmd/main.go &
```

## Production Environment

### System Requirements

**Minimum:**
- 2 CPU cores
- 4 GB RAM
- 50 GB storage
- Ubuntu 20.04+ / CentOS 8+

**Recommended:**
- 4 CPU cores
- 8 GB RAM
- 100 GB SSD storage
- Load balancer
- Monitoring system

### Production Docker Compose

Create a `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    networks:
      - backend
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    networks:
      - backend
    restart: unless-stopped

  minio:
    image: minio/minio
    environment:
      MINIO_ROOT_USER: ${MINIO_USER}
      MINIO_ROOT_PASSWORD: ${MINIO_PASSWORD}
    command: server /data --console-address ":9001"
    volumes:
      - minio_data:/data
    networks:
      - backend
    restart: unless-stopped

  auth-service:
    build:
      context: .
      dockerfile: Dockerfile.auth
    environment:
      DATABASE_URL: postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable
      JWT_SECRET: ${JWT_SECRET}
      GRPC_PORT: 8080
    depends_on:
      - postgres
    networks:
      - backend
    restart: unless-stopped
    deploy:
      replicas: 2

  course-service:
    build:
      context: .
      dockerfile: Dockerfile.course
    environment:
      DATABASE_URL: postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable
      GRPC_PORT: 8080
    depends_on:
      - postgres
    networks:
      - backend
    restart: unless-stopped
    deploy:
      replicas: 2

  progress-service:
    build:
      context: .
      dockerfile: Dockerfile.progress
    environment:
      DATABASE_URL: postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable
      GRPC_PORT: 8080
    depends_on:
      - postgres
    networks:
      - backend
    restart: unless-stopped
    deploy:
      replicas: 2

  api-gateway:
    build:
      context: .
      dockerfile: Dockerfile.gateway
    environment:
      AUTH_SERVICE_URL: auth-service:8080
      COURSE_SERVICE_URL: course-service:8080
      PROGRESS_SERVICE_URL: progress-service:8080
      JWT_SECRET: ${JWT_SECRET}
      HTTP_PORT: 8080
    depends_on:
      - auth-service
      - course-service
      - progress-service
    networks:
      - backend
      - frontend
    restart: unless-stopped
    deploy:
      replicas: 3

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - api-gateway
    networks:
      - frontend
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
  minio_data:

networks:
  backend:
    driver: bridge
  frontend:
    driver: bridge
```

### Environment Variables

Create `.env` file:

```bash
# Database
DB_NAME=studyplatform
DB_USER=studyplatform_user
DB_PASSWORD=your_secure_password

# JWT
JWT_SECRET=your_jwt_secret_key_here

# MinIO
MINIO_USER=admin
MINIO_PASSWORD=your_minio_password

# OAuth (if using)
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
```

### SSL/TLS Configuration

Create `nginx.conf`:

```nginx
events {
    worker_connections 1024;
}

http {
    upstream api_gateway {
        server api-gateway:8080;
    }

    server {
        listen 80;
        server_name your-domain.com;
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name your-domain.com;

        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers HIGH:!aNULL:!MD5;

        location / {
            proxy_pass http://api_gateway;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            
            # Rate limiting
            limit_req zone=api burst=50 nodelay;
        }

        location /health {
            proxy_pass http://api_gateway;
            access_log off;
        }
    }

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
}
```

### Deployment

```bash
# Create .env file with production values
cp .env.example .env
vim .env

# Deploy
docker-compose -f docker-compose.prod.yml up -d

# Check status
docker-compose -f docker-compose.prod.yml ps

# Scale services
docker-compose -f docker-compose.prod.yml up -d --scale api-gateway=3
```

## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster (1.20+)
- kubectl configured
- Helm 3+ (optional)

### Kubernetes Manifests

Create `k8s/namespace.yaml`:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: study-platform
```

Create `k8s/configmap.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: study-platform-config
  namespace: study-platform
data:
  DB_NAME: studyplatform
  GRPC_PORT: "8080"
  HTTP_PORT: "8080"
```

Create `k8s/secret.yaml`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: study-platform-secret
  namespace: study-platform
type: Opaque
data:
  DB_PASSWORD: <base64-encoded-password>
  JWT_SECRET: <base64-encoded-jwt-secret>
```

Create `k8s/postgres.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: study-platform
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:16
        env:
        - name: POSTGRES_DB
          valueFrom:
            configMapKeyRef:
              name: study-platform-config
              key: DB_NAME
        - name: POSTGRES_USER
          value: postgres
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: study-platform-secret
              key: DB_PASSWORD
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: study-platform
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: study-platform
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

Create `k8s/api-gateway.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
  namespace: study-platform
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
    spec:
      containers:
      - name: api-gateway
        image: study-platform/api-gateway:latest
        env:
        - name: AUTH_SERVICE_URL
          value: "auth-service:8080"
        - name: COURSE_SERVICE_URL
          value: "course-service:8080"
        - name: PROGRESS_SERVICE_URL
          value: "progress-service:8080"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: study-platform-secret
              key: JWT_SECRET
        - name: HTTP_PORT
          valueFrom:
            configMapKeyRef:
              name: study-platform-config
              key: HTTP_PORT
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: api-gateway
  namespace: study-platform
spec:
  selector:
    app: api-gateway
  ports:
  - port: 8080
    targetPort: 8080
  type: LoadBalancer
```

### Deploy to Kubernetes

```bash
# Apply manifests
kubectl apply -f k8s/

# Check deployments
kubectl get deployments -n study-platform

# Check pods
kubectl get pods -n study-platform

# Check services
kubectl get services -n study-platform

# View logs
kubectl logs -f deployment/api-gateway -n study-platform
```

## Database Setup

### Production Database

For production, use managed PostgreSQL services:

- **AWS RDS**
- **Google Cloud SQL**
- **Azure Database for PostgreSQL**
- **DigitalOcean Managed Databases**

### Database Configuration

```sql
-- Create database
CREATE DATABASE studyplatform;

-- Create user
CREATE USER studyplatform_user WITH PASSWORD 'your_secure_password';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE studyplatform TO studyplatform_user;

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
```

### Connection Pooling

Use connection pooling for better performance:

```yaml
# PgBouncer configuration
auth-service:
  environment:
    DATABASE_URL: postgres://studyplatform_user:password@pgbouncer:5432/studyplatform
```

## Monitoring & Logging

### Prometheus Metrics

Add to services:

```yaml
# docker-compose.yml
prometheus:
  image: prom/prometheus
  ports:
    - "9090:9090"
  volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml

grafana:
  image: grafana/grafana
  ports:
    - "3000:3000"
  environment:
    - GF_SECURITY_ADMIN_PASSWORD=admin
```

### Logging Configuration

Use structured logging:

```yaml
services:
  app:
    environment:
      LOG_LEVEL: info
      LOG_FORMAT: json
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### Health Checks

All services include health checks:

```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

## Security Configuration

### SSL/TLS

- Use Let's Encrypt for free SSL certificates
- Configure HTTPS redirects
- Enable HSTS headers
- Use TLS 1.2+ only

### Security Headers

```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
```

### Network Security

- Use private networks for internal communication
- Implement firewall rules
- Enable VPC security groups
- Use secrets management

### Authentication

- Use strong JWT secrets
- Implement token rotation
- Set appropriate token expiration
- Use secure OAuth configurations

## Performance Tuning

### Database Optimization

```sql
-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_courses_category ON courses(category);
CREATE INDEX idx_enrollments_user_id ON enrollments(user_id);
CREATE INDEX idx_progress_user_course ON user_progress(user_id, course_id);
```

### Connection Pooling

```yaml
auth-service:
  environment:
    DB_MAX_OPEN_CONNS: 25
    DB_MAX_IDLE_CONNS: 5
    DB_CONN_MAX_LIFETIME: 300s
```

### Caching

Add Redis caching:

```yaml
redis:
  image: redis:7-alpine
  command: redis-server --maxmemory 512mb --maxmemory-policy allkeys-lru
```

### Load Balancing

Use multiple service instances:

```yaml
api-gateway:
  deploy:
    replicas: 3
    resources:
      limits:
        cpus: "0.5"
        memory: 512M
      reservations:
        cpus: "0.25"
        memory: 256M
```

## Backup & Recovery

### Database Backups

```bash
# Create backup
docker exec postgres pg_dump -U postgres studyplatform > backup.sql

# Restore backup
docker exec -i postgres psql -U postgres studyplatform < backup.sql
```

### Automated Backups

```yaml
# Cron job for backups
backup:
  image: postgres:16
  volumes:
    - ./backups:/backups
  command: |
    sh -c "
    while true; do
      pg_dump -h postgres -U postgres studyplatform > /backups/backup-$(date +%Y%m%d-%H%M%S).sql
      sleep 86400
    done"
```

## Troubleshooting

### Common Issues

**Services won't start:**
```bash
# Check logs
docker-compose logs service-name

# Check resource usage
docker stats

# Check network connectivity
docker exec service-name ping other-service
```

**Database connection issues:**
```bash
# Check database status
docker exec postgres pg_isready -U postgres

# Check database logs
docker logs postgres

# Test connection
docker exec -it postgres psql -U postgres -d studyplatform
```

**High memory usage:**
```bash
# Check memory usage
docker stats

# Limit memory usage
docker-compose.yml:
  services:
    app:
      mem_limit: 512m
```

### Debugging

Enable debug logging:

```yaml
environment:
  LOG_LEVEL: debug
```

Use health check endpoints:

```bash
# Check service health
curl http://localhost:8080/api/v1/health

# Check circuit breakers
curl http://localhost:8080/api/v1/health/circuit-breakers
```

### Performance Issues

Monitor key metrics:
- Response times
- Error rates
- CPU/Memory usage
- Database query performance
- Connection pool stats

## Scaling

### Horizontal Scaling

```bash
# Scale specific services
docker-compose up -d --scale api-gateway=3 --scale auth-service=2

# Kubernetes scaling
kubectl scale deployment api-gateway --replicas=5 -n study-platform
```

### Vertical Scaling

```yaml
# Increase resource limits
services:
  api-gateway:
    deploy:
      resources:
        limits:
          cpus: "1.0"
          memory: 1G
```

### Auto-scaling

Kubernetes HPA:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-gateway-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api-gateway
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

## Maintenance

### Updates

```bash
# Update images
docker-compose pull

# Recreate containers
docker-compose up -d

# Rolling update in Kubernetes
kubectl rollout restart deployment/api-gateway -n study-platform
```

### Database Migrations

```bash
# Run migrations
docker exec service-name /app/migrate up

# Check migration status
docker exec service-name /app/migrate version
```

This deployment guide provides comprehensive instructions for deploying the Study Platform in various environments with proper security, monitoring, and scaling configurations.