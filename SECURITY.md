# Study Platform - Security Guide

## 🔒 Critical Security Overview

This document outlines the comprehensive security measures implemented in the Study Platform. **READ THIS ENTIRELY** before deploying to production.

## 🚨 Emergency Security Setup

### Immediate Actions Required

1. **Generate Secure Environment**
   ```bash
   cd Study-Platform/
   ./scripts/setup-secure-env.sh
   ```

2. **Replace Placeholder Credentials**
   - Edit `.env` file and replace ALL `REPLACE_WITH_REAL_*` values
   - Use proper secret management in production

3. **Validate Security**
   ```bash
   # Check environment variables are set
   grep -E "REPLACE_WITH_REAL|your-secret-key-here|password123|admin" .env
   # Should return NO matches
   ```

## 🛡️ Security Architecture

### 1. Secrets Management

#### Environment Variables (.env)
```bash
# CRITICAL: Never commit .env to version control
# Strong credentials generated automatically
JWT_SECRET=<64-character-secure-string>
POSTGRES_PASSWORD=<24-character-secure-password>
```

#### Docker Compose Security
- All hardcoded credentials removed
- Environment variable substitution implemented
- Secure defaults with fallbacks

#### Production Secret Management
```bash
# AWS Secrets Manager (Recommended)
aws secretsmanager create-secret --name "study-platform/jwt-secret" --secret-string "..."

# Docker Swarm Secrets
echo "secure-password" | docker secret create db_password -

# Kubernetes Secrets
kubectl create secret generic study-platform-secrets --from-env-file=.env
```

### 2. Authentication & Authorization

#### JWT Security
- **256-bit minimum** secret keys
- HS256 algorithm with secure implementation
- Token expiry: 24 hours (configurable)
- Refresh token support: 168 hours
- Algorithm confusion attack prevention

#### Password Security
- bcrypt with 12 rounds (configurable)
- Salted hashes
- Password strength validation
- Account lockout after 5 failed attempts
- 15-minute lockout duration

#### Multi-Factor Authentication (MFA)
- TOTP support ready
- Backup codes implementation
- MFA enforcement for admin accounts

### 3. API Security

#### Input Validation Middleware
```go
// Protects against:
- SQL injection
- XSS attacks
- Command injection
- Path traversal
- LDAP injection
```

#### Rate Limiting
- **100 requests/second** per IP (configurable)
- **200 burst capacity**
- IP-based tracking
- Automatic cleanup of tracking data

#### Security Headers
```http
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
Content-Security-Policy: default-src 'self'
```

#### CORS Security
- Whitelist-based origin validation
- Credentials support with strict origin checking
- Preflight request handling

### 4. Database Security

#### Connection Security
- Encrypted connections (SSL/TLS)
- Connection pooling with limits
- Prepared statements (SQL injection prevention)

#### Access Control
- Role-based database users
- Minimal privilege principle
- Row-level security (RLS) enabled
- Audit logging for all operations

#### Data Protection
```sql
-- User data encryption
- Password hashing with salt
- Sensitive data encryption at rest
- Audit trail for all changes
```

### 5. Inter-Service Communication

#### Service Authentication
- JWT tokens for service-to-service calls
- Shared secret validation
- Service whitelist enforcement

#### Network Security
- Internal Docker network isolation
- Only API Gateway exposed externally (port 8080)
- TLS for gRPC communications

#### Service Discovery
- Health checking system
- Service registry with TTL
- Automatic failover support

### 6. Infrastructure Security

#### Docker Security
```yaml
# Non-root users in containers
# Read-only file systems where possible
# Minimal attack surface
# Resource limits enforced
```

#### Network Isolation
```bash
# Only port 8080 exposed externally
# All services communicate via internal Docker network
# Database accessible only from application services
```

## 🔧 Security Configuration

### Environment-Specific Settings

#### Development
```bash
ENVIRONMENT=development
DEBUG_MODE=false  # Never true in any environment
TLS_ENABLED=false  # Use HTTPS proxy in production
```

#### Staging
```bash
ENVIRONMENT=staging
TLS_ENABLED=true
RATE_LIMIT_REQUESTS=50  # Lower for testing
```

#### Production
```bash
ENVIRONMENT=production
TLS_ENABLED=true
RATE_LIMIT_REQUESTS=100
SECURITY_HEADERS_ENABLED=true
LOG_LEVEL=warn
```

### Security Monitoring

#### Audit Logging
```sql
-- All authentication events logged
-- Failed login attempts tracked
-- Suspicious activity alerts
-- Session management tracked
```

#### Metrics and Alerting
- Failed authentication rates
- Rate limit violations
- Unusual traffic patterns
- Error rate monitoring

## 🚀 Deployment Security

### Pre-Deployment Checklist

- [ ] All placeholder credentials replaced
- [ ] JWT secrets are 256-bit minimum
- [ ] Database passwords are 20+ characters
- [ ] CORS origins updated for production domains
- [ ] TLS certificates configured
- [ ] Firewall rules configured
- [ ] Monitoring and alerting set up
- [ ] Backup strategies implemented
- [ ] Security scanning completed

### Production Deployment

#### 1. Infrastructure Setup
```bash
# Configure firewall
ufw allow 443/tcp  # HTTPS
ufw allow 22/tcp   # SSH (restrict to admin IPs)
ufw deny 8080/tcp  # Block direct access to API gateway

# Set up reverse proxy with TLS
nginx/traefik/cloudflare for TLS termination
```

#### 2. Container Security
```bash
# Use specific versions (not latest)
# Scan images for vulnerabilities
docker scan study-platform/api-gateway:v1.0.0

# Run containers as non-root
USER 1001:1001 in Dockerfiles
```

#### 3. Database Security
```bash
# Enable encryption at rest
# Set up automated backups with encryption
# Configure read replicas for scaling
# Set up monitoring and alerting
```

### 4. Monitoring Setup
```bash
# Application Performance Monitoring
# Security Information and Event Management (SIEM)
# Log aggregation and analysis
# Real-time alerting for security events
```

## 🔍 Security Testing

### Automated Testing
```bash
# Run security tests
go test ./internal/middleware/security_middleware_test.go
go test ./internal/security/service_auth_test.go

# Static security analysis
gosec ./...

# Dependency vulnerability scanning
govulncheck ./...
```

### Manual Security Testing

#### Authentication Testing
- [ ] JWT token validation
- [ ] Password strength enforcement
- [ ] Account lockout mechanism
- [ ] Session management

#### API Security Testing
- [ ] Input validation (SQL injection, XSS)
- [ ] Rate limiting effectiveness
- [ ] CORS policy enforcement
- [ ] Security headers present

#### Infrastructure Testing
- [ ] Network isolation
- [ ] Port exposure verification
- [ ] Container security scanning
- [ ] TLS configuration validation

## 📋 Security Maintenance

### Regular Tasks

#### Daily
- Monitor security alerts
- Review failed authentication attempts
- Check system resource usage

#### Weekly
- Review audit logs
- Update dependencies
- Security patch assessment

#### Monthly
- Rotate JWT secrets
- Review user access permissions
- Security configuration review
- Vulnerability scanning

#### Quarterly
- Full security audit
- Penetration testing
- Business continuity testing
- Security training updates

### Incident Response

#### Security Incident Procedure
1. **Immediate Response**
   - Isolate affected systems
   - Preserve evidence
   - Notify stakeholders

2. **Investigation**
   - Analyze logs and evidence
   - Determine scope and impact
   - Document findings

3. **Remediation**
   - Apply security patches
   - Rotate compromised credentials
   - Update security controls

4. **Recovery**
   - Restore normal operations
   - Monitor for recurring issues
   - Update procedures

## 📞 Security Contacts

### Emergency Response
- **Security Team**: security@studyplatform.com
- **On-Call**: +1-XXX-XXX-XXXX
- **Incident Response**: incidents@studyplatform.com

### Reporting Vulnerabilities
- **Security Email**: security@studyplatform.com
- **Bug Bounty Program**: hackerone.com/studyplatform
- **PGP Key**: [Public key for encrypted communications]

## 📚 Additional Resources

### Security Standards
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [OWASP Application Security Verification Standard](https://owasp.org/www-project-application-security-verification-standard/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)

### Security Tools
- [OWASP ZAP](https://owasp.org/www-project-zap/) - Security testing
- [Snyk](https://snyk.io/) - Vulnerability scanning
- [SonarQube](https://www.sonarqube.org/) - Code security analysis

### Training Resources
- [OWASP WebGoat](https://owasp.org/www-project-webgoat/) - Security training
- [PortSwigger Academy](https://portswigger.net/web-security) - Web security education

---

## ⚠️ Critical Reminders

1. **Never commit `.env` files** to version control
2. **Always use HTTPS** in production
3. **Rotate secrets regularly** (quarterly minimum)
4. **Monitor security logs** continuously
5. **Keep dependencies updated** weekly
6. **Test security measures** regularly
7. **Train development team** on security best practices

---

**Last Updated**: [Current Date]  
**Version**: 1.0.0  
**Next Review**: [Date + 3 months]