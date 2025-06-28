# Go Web App on EKS

A simple Go web application deployed to Amazon EKS with Jenkins CI/CD.

## 🚀 Quick Start

```bash
# Run locally
go run main.go

# Build Docker image
docker build -t cfx-test-go .

# Deploy to Kubernetes
kubectl apply -f k8s/
```

## 📋 What You Need

- Go 1.19+
- Docker
- AWS CLI
- kubectl
- Jenkins
- EKS cluster

## 🌐 Endpoints

- **Live URL**: https://api.rinando.my.id
- **Health Check**: `/health`
- **API**: `/api/hello`
- **Home**: `/`

## 🔧 Fix These Issues First

Your Jenkins pipeline won't work until you fix these naming mismatches:

### 1. Update Jenkinsfile
Change these lines in your Jenkinsfile:
```bash
# FROM:
kubectl set image deployment/cfx-go-app cfx-go-app=...
kubectl patch deployment cfx-go-app ...

# TO:
kubectl set image deployment/go-web-app go-web-app=...
kubectl patch deployment go-web-app ...
```

### 2. Fix ingress.yml
```yaml
# Change service name from:
name: go-api-service

# To:
name: go-web-app-service
```

### 3. Fix deployment.yml
```yaml
# Change hardcoded image:
image: 112113402575.dkr.ecr.us-east-1.amazonaws.com/cfx-test-go:latest

# To:
image: IMAGE_PLACEHOLDER
```

## 🐳 Docker

The app runs on port 8080 and includes health checks.

## ☸️ Kubernetes

- **Replicas**: 3
- **Resources**: 128Mi memory, 100m CPU
- **Service**: ClusterIP on port 80
- **Ingress**: ALB with custom domain

## 🔄 CI/CD Pipeline

Jenkins automatically:
1. Builds Docker image
2. Pushes to ECR
3. Deploys to EKS
4. Updates with rolling deployment

## 🐛 Debug Commands

```bash
# Check pods
kubectl get pods -l app=go-web-app

# View logs
kubectl logs -l app=go-web-app

# Check service
kubectl get svc go-web-app-service

# Port forward to test
kubectl port-forward svc/go-web-app-service 8080:80
```

## 📁 Project Structure

```
├── main.go             
├── go.mod
├── go.sum
├── Dockerfile           
├── Jenkinsfile         
└── k8s/
    ├── deployment.yml   
    └── ingress.yml      
```

## 🔑 Key Features

- RESTful API with JSON responses
- Built-in health checks
- Docker containerized
- Kubernetes ready
- Auto-scaling with 3 replicas
- Custom domain routing
- Automated CI/CD

---

**Next Steps**: Fix the naming issues above, then push your code to trigger the Jenkins pipeline!
