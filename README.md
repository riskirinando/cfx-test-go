# Go Web App on AWS EKS

A containerized Go web application deployed to Amazon EKS with automated CI/CD using Jenkins.

## 🌐 Live Application

**Public URL**: https://api.rinando.my.id

## 🏗️ Architecture

```
GitHub → Jenkins → AWS ECR → Amazon EKS → ALB → Public Internet
```

## 🚀 Features

- **REST API**: JSON endpoints with health checks
- **Containerized**: Docker-based deployment
- **Auto-scaling**: 3 replicas with resource limits
- **CI/CD**: Automated Jenkins pipeline
- **Custom Domain**: ALB ingress with SSL
- **Monitoring**: Built-in health and readiness probes

## 📱 API Endpoints

| Endpoint | Description | Example Response |
|----------|-------------|------------------|
| `/` | HTML landing page | Web interface |
| `/api/hello` | JSON API | `{"message": "Hello from Go!", "host": "pod-123"}` |
| `/health` | Health check | `{"status": "healthy"}` |
| `/ready` | Readiness probe | `{"status": "ready"}` |

## 🔧 Local Development

**Prerequisites**: Go 1.19+

```bash
# Clone and run
git clone <repo-url>
cd <repo-name>
go mod download  # Install dependencies from go.mod
go run main.go

# Access locally
open http://localhost:8080
```

## 🐳 Docker

```bash
# Build image
docker build -t cfx-test-go .

# Run container
docker run -p 8080:8080 cfx-test-go
```

## ☸️ Kubernetes Deployment

**Current Setup**:
- **Deployment**: `go-web-app` (3 replicas)
- **Service**: `go-web-app-service` (ClusterIP)
- **Ingress**: ALB with domain `api.rinando.my.id`
- **Resources**: 128Mi memory, 100m CPU per pod

**Files**:
- `k8s/deployment.yml` - Deployment + Service
- `k8s/ingress.yml` - ALB ingress configuration

## 🔄 CI/CD Pipeline (Jenkins)

**Automatic deployment on git push**:

1. **Build**: Docker image creation
2. **Push**: Upload to AWS ECR
3. **Deploy**: Rolling update to EKS cluster

**Environment**:
- **ECR Repository**: `cfx-test-go`
- **EKS Cluster**: `test-project-eks-cluster`
- **AWS Region**: `us-east-1`

## ⚠️ Configuration Issues to Fix

**Your Jenkins pipeline needs updates to match your Kubernetes resources:**

1. **Deployment name mismatch**:
   - Jenkinsfile uses: `cfx-go-app`
   - Actual deployment: `go-web-app`

2. **Service name mismatch**:
   - Jenkinsfile uses: `cfx-go-service`
   - Actual service: `go-web-app-service`

3. **Ingress service reference**:
   - Current: `go-api-service` (doesn't exist)
   - Should be: `go-web-app-service`

4. **Image placeholder**:
   - Replace hardcoded image in deployment.yml with `IMAGE_PLACEHOLDER`

## 🔧 Quick Fixes

**Fix 1: Update Jenkinsfile deployment commands**:
```groovy
// Change from:
kubectl set image deployment/cfx-go-app cfx-go-app=${env.FULL_IMAGE_URI}

// To:
kubectl set image deployment/go-web-app go-web-app=${env.FULL_IMAGE_URI}
```

**Fix 2: Update ingress.yml service name**:
```yaml
backend:
  service:
    name: go-web-app-service  # Fix this
    port:
      number: 80
```

**Fix 3: Update deployment.yml image**:
```yaml
# Change from:
image: 112113402575.dkr.ecr.us-east-1.amazonaws.com/cfx-test-go:latest

# To:
image: IMAGE_PLACEHOLDER
```

## 🛠️ Manual Commands

**Deploy manually**:
```bash
# Build and push
docker build -t cfx-test-go .
docker tag cfx-test-go 112113402575.dkr.ecr.us-east-1.amazonaws.com/cfx-test-go:latest
docker push 112113402575.dkr.ecr.us-east-1.amazonaws.com/cfx-test-go:latest

# Update deployment
kubectl set image deployment/go-web-app go-web-app=112113402575.dkr.ecr.us-east-1.amazonaws.com/cfx-test-go:latest
```

**Check status**:
```bash
# View pods
kubectl get pods -l app=go-web-app

# Check logs
kubectl logs -l app=go-web-app

# Test locally
kubectl port-forward svc/go-web-app-service 8080:80
```

## 📊 Monitoring

**Health checks configured**:
- **Liveness**: `/health` every 10s
- **Readiness**: `/ready` every 5s
- **ALB Health**: `/health` every 30s

**Resource limits**:
- **Memory**: 64Mi request, 128Mi limit
- **CPU**: 50m request, 100m limit

## 🚨 Troubleshooting

**Common issues**:

```bash
# Pod not starting
kubectl describe pod <pod-name>

# Service not accessible
kubectl get svc go-web-app-service
kubectl describe ingress multi-app-ingress

# Check ALB status
kubectl get ingress -o wide
```

**Pipeline fails**:
- Check AWS credentials in Jenkins
- Verify ECR repository exists
- Ensure EKS cluster is accessible

## 📝 Project Structure

```
├── main.go              # Go application
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
├── Dockerfile           # Container configuration
├── Jenkinsfile          # CI/CD pipeline
├── k8s/
│   ├── deployment.yml   # Kubernetes deployment + service
│   └── ingress.yml      # ALB ingress configuration
└── README.md           # This file
```

## 🤝 Contributing

1. Fork repository
2. Make changes
3. Test locally: `go run main.go`
4. Push to trigger Jenkins pipeline
5. Check deployment at https://api.rinando.my.id
