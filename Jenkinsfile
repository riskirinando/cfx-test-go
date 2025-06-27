pipeline {
    agent any
    
    environment {
        AWS_REGION = 'us-east-1'
        ECR_REPOSITORY = 'cfx-test-go'
        EKS_CLUSTER_NAME = 'test-project-eks-cluster'
        IMAGE_TAG = "${BUILD_NUMBER}"
        ECR_REGISTRY = "${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"
        CURRENT_STAGE = 'NONE'
        KUBECONFIG = '/tmp/kubeconfig'
        GO_VERSION = '1.21'  // Adjust to your preferred Go version
        GOOS = 'linux'
        GOARCH = 'amd64'
        CGO_ENABLED = '0'
    }
    
    stages {
        stage('Checkout') {
            steps {
                script {
                    env.CURRENT_STAGE = 'CHECKOUT'
                    echo "=== STAGE: ${env.CURRENT_STAGE} ==="
                    
                    try {
                        checkout scm
                        echo "✅ SCM checkout successful"
                        
                        // Get Git commit hash
                        try {
                            env.GIT_COMMIT = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
                            echo "✅ Git commit: ${env.GIT_COMMIT}"
                        } catch (Exception e) {
                            echo "⚠️ Git operation failed: ${e.getMessage()}"
                            env.GIT_COMMIT = "unknown-${env.BUILD_NUMBER}"
                        }
                        
                        // List files to verify checkout
                        sh 'ls -la'
                        
                        echo "✅ CHECKOUT stage completed successfully"
                        
                    } catch (Exception e) {
                        env.CURRENT_STAGE = 'CHECKOUT_FAILED'
                        echo "❌ CHECKOUT stage failed: ${e.getMessage()}"
                        throw e
                    }
                }
            }
        }
        
        stage('Prerequisites Check') {
            steps {
                script {
                    env.CURRENT_STAGE = 'PREREQUISITES'
                    echo "=== STAGE: ${env.CURRENT_STAGE} ==="
                    
                    try {
                        // Check for required files
                        echo "Checking for required files..."
                        
                        if (!fileExists('Dockerfile')) {
                            error "❌ Dockerfile not found! Please ensure Dockerfile exists in the repository."
                        }
                        echo "✅ Dockerfile found"
                        
                        if (fileExists('go.mod')) {
                            echo "✅ go.mod found"
                            sh 'cat go.mod'
                        } else {
                            echo "⚠️ go.mod not found - ensure this is a Go module"
                        }
                        
                        if (fileExists('go.sum')) {
                            echo "✅ go.sum found"
                        } else {
                            echo "⚠️ go.sum not found - will be created during build"
                        }
                        
                        // Check required tools
                        sh 'docker --version'
                        sh 'aws --version'
                        sh 'kubectl version --client'
                        
                        // Check Go version
                        try {
                            sh 'go version'
                            echo "✅ Go is available"
                        } catch (Exception e) {
                            echo "⚠️ Go not found in PATH: ${e.getMessage()}"
                            echo "💡 Go will be installed in the build stage if needed"
                        }
                        
                        echo "✅ PREREQUISITES stage completed successfully"
                        
                    } catch (Exception e) {
                        env.CURRENT_STAGE = 'PREREQUISITES_FAILED'
                        echo "❌ PREREQUISITES stage failed: ${e.getMessage()}"
                        throw e
                    }
                }
            }
        }
        
        stage('Setup Go Environment') {
            steps {
                script {
                    env.CURRENT_STAGE = 'GO_SETUP'
                    echo "=== STAGE: ${env.CURRENT_STAGE} ==="
                    
                    try {
                        // Setup Go environment
                        sh """
                            # Check if Go is available and correct version
                            if command -v go &> /dev/null; then
                                CURRENT_GO_VERSION=\$(go version | awk '{print \$3}' | sed 's/go//')
                                echo "Current Go version: \$CURRENT_GO_VERSION"
                            else
                                echo "Go not found, will use Docker for build"
                            fi
                            
                            # Set Go environment variables
                            export GOOS=${env.GOOS}
                            export GOARCH=${env.GOARCH}
                            export CGO_ENABLED=${env.CGO_ENABLED}
                            export GO111MODULE=on
                            
                            echo "Go environment configured:"
                            echo "GOOS=${env.GOOS}"
                            echo "GOARCH=${env.GOARCH}"
                            echo "CGO_ENABLED=${env.CGO_ENABLED}"
                        """
                        
                        echo "✅ GO_SETUP stage completed successfully"
                        
                    } catch (Exception e) {
                        env.CURRENT_STAGE = 'GO_SETUP_FAILED'
                        echo "❌ GO_SETUP stage failed: ${e.getMessage()}"
                        throw e
                    }
                }
            }
        }
        
        stage('Build & Test Go Application') {
            steps {
                script {
                    env.CURRENT_STAGE = 'GO_BUILD_TEST'
                    echo "=== STAGE: ${env.CURRENT_STAGE} ==="
                    
                    try {
                        // Validate go.mod
                        if (fileExists('go.mod')) {
                            sh 'cat go.mod'
                            echo "✅ go.mod validated"
                        }
                        
                        // Download dependencies and run tests
                        sh """
                            # Check if we can use local Go or need Docker
                            if command -v go &> /dev/null && go version | grep -q "go${env.GO_VERSION}"; then
                                echo "Using local Go installation"
                                
                                # Download dependencies
                                go mod download
                                go mod verify
                                go mod tidy
                                
                                # Run tests (skip if no test files exist)
                                if find . -name '*_test.go' | grep -q .; then
                                    echo "Running Go tests..."
                                    go test ./... -v
                                else
                                    echo "No test files found, skipping tests"
                                fi
                                
                                # Lint the code (optional)
                                if command -v golangci-lint &> /dev/null; then
                                    echo "Running linter..."
                                    golangci-lint run
                                else
                                    echo "golangci-lint not found, using go vet"
                                    go vet ./...
                                fi
                                
                                # Build the application (matching your Dockerfile)
                                export GOOS=${env.GOOS}
                                export GOARCH=${env.GOARCH}
                                export CGO_ENABLED=${env.CGO_ENABLED}
                                
                                go build -a -installsuffix cgo -o main .
                                
                                # Check binary
                                ls -la main
                                file main
                                
                            else
                                echo "Using Docker for Go build and test"
                                echo "Will rely on multi-stage Dockerfile for build process"
                                echo "Validating Dockerfile syntax..."
                                
                                # Just validate that Docker build works
                                docker build --target builder -t go-build-test .
                                docker rmi go-build-test
                            fi
                        """
                        
                        echo "✅ GO_BUILD_TEST stage completed successfully"
                        
                    } catch (Exception e) {
                        env.CURRENT_STAGE = 'GO_BUILD_TEST_FAILED'
                        echo "❌ GO_BUILD_TEST stage failed: ${e.getMessage()}"
                        throw e
                    }
                }
            }
        }
        
        stage('Get AWS Account ID') {
            steps {
                script {
                    env.CURRENT_STAGE = 'AWS_ACCOUNT_ID'
                    echo "=== STAGE: ${env.CURRENT_STAGE} ==="
                    
                    try {
                        // Check AWS CLI configuration first
                        sh 'aws sts get-caller-identity || echo "AWS CLI not configured or no permissions"'
                        
                        // Get AWS Account ID with better error handling
                        def awsAccountResult = sh(
                            returnStdout: true,
                            script: 'aws sts get-caller-identity --query Account --output text 2>/dev/null || echo "FAILED"'
                        ).trim()
                        
                        if (awsAccountResult == "FAILED" || awsAccountResult == "") {
                            echo "⚠️ Could not get AWS Account ID from AWS CLI, using fallback..."
                            // Fallback: try to extract from AWS CLI config or use a default
                            env.AWS_ACCOUNT_ID = "123456789012" // Replace with your actual account ID
                            echo "⚠️ Using fallback AWS Account ID: ${env.AWS_ACCOUNT_ID}"
                            echo "⚠️ Please ensure AWS CLI is properly configured with credentials"
                        } else {
                            env.AWS_ACCOUNT_ID = awsAccountResult
                            echo "✅ AWS Account ID: ${env.AWS_ACCOUNT_ID}"
                        }
                        
                        // Update ECR registry URL - force string interpolation
                        def awsAccountId = env.AWS_ACCOUNT_ID
                        def awsRegion = env.AWS_REGION
                        def ecrRepo = env.ECR_REPOSITORY
                        def imageTag = env.IMAGE_TAG
                        
                        env.ECR_REGISTRY = "${awsAccountId}.dkr.ecr.${awsRegion}.amazonaws.com"
                        env.FULL_IMAGE_URI = "${awsAccountId}.dkr.ecr.${awsRegion}.amazonaws.com/${ecrRepo}:${imageTag}"
                        
                        echo "✅ ECR Registry: ${env.ECR_REGISTRY}"
                        echo "✅ Full Image URI: ${env.FULL_IMAGE_URI}"
                        
                        echo "✅ AWS_ACCOUNT_ID stage completed successfully"
                        
                    } catch (Exception e) {
                        env.CURRENT_STAGE = 'AWS_ACCOUNT_ID_FAILED'
                        echo "❌ AWS_ACCOUNT_ID stage failed: ${e.getMessage()}"
                        throw e
                    }
                }
            }
        }
        
        stage('Docker Build') {
            steps {
                script {
                    env.CURRENT_STAGE = 'DOCKER_BUILD'
                    echo "=== STAGE: ${env.CURRENT_STAGE} ==="
                    
                    try {
                        echo "Building Docker image..."
                        echo "Image URI: ${env.FULL_IMAGE_URI}"
                        
                        // Check Docker daemon access and current user
                        sh 'whoami'
                        sh 'groups'
                        sh 'ls -la /var/run/docker.sock'
                        
                        // Try docker without sudo first, then with sudo as fallback
                        def dockerCmd = ""
                        try {
                            sh 'docker version'
                            dockerCmd = "docker"
                            echo "✅ Docker accessible without sudo"
                        } catch (Exception e) {
                            echo "⚠️ Docker not accessible without sudo, trying with sudo..."
                            sh 'sudo docker version'
                            dockerCmd = "sudo docker"
                            echo "✅ Docker accessible with sudo"
                        }
                        
                        // Build Docker image
                        sh """
                            ${dockerCmd} build -t ${env.ECR_REPOSITORY}:${env.IMAGE_TAG} .
                            ${dockerCmd} tag ${env.ECR_REPOSITORY}:${env.IMAGE_TAG} ${env.FULL_IMAGE_URI}
                            ${dockerCmd} tag ${env.ECR_REPOSITORY}:${env.IMAGE_TAG} ${env.ECR_REGISTRY}/${env.ECR_REPOSITORY}:latest
                        """
                        
                        // List Docker images
                        sh "${dockerCmd} images | grep cfx-test-go || echo 'No cfx-test-go images found'"
                        
                        // Store docker command for later stages
                        env.DOCKER_CMD = dockerCmd
                        
                        echo "✅ DOCKER_BUILD stage completed successfully"
                        
                    } catch (Exception e) {
                        env.CURRENT_STAGE = 'DOCKER_BUILD_FAILED'
                        echo "❌ DOCKER_BUILD stage failed: ${e.getMessage()}"
                        echo "💡 To fix permanently, run: sudo usermod -aG docker jenkins && sudo systemctl restart jenkins"
                        throw e
                    }
                }
            }
        }
        
        stage('ECR Login & Push') {
            steps {
                script {
                    env.CURRENT_STAGE = 'ECR_PUSH'
                    echo "=== STAGE: ${env.CURRENT_STAGE} ==="
        
                    // 1. Get AWS account ID dynamically
                    env.AWS_ACCOUNT_ID = sh(
                        script: "aws sts get-caller-identity --query Account --output text",
                        returnStdout: true
                    ).trim()
        
                    // 2. Now set ECR registry and image URI using updated AWS_ACCOUNT_ID
                    def ECR_REGISTRY = "${env.AWS_ACCOUNT_ID}.dkr.ecr.${env.AWS_REGION}.amazonaws.com"
                    def FULL_IMAGE_URI = "${ECR_REGISTRY}/${env.ECR_REPOSITORY}:${env.IMAGE_TAG}"
        
                    echo "✅ AWS_ACCOUNT_ID: ${env.AWS_ACCOUNT_ID}"
                    echo "✅ ECR_REGISTRY: ${ECR_REGISTRY}"
                    echo "✅ FULL_IMAGE_URI: ${FULL_IMAGE_URI}"
        
                    // 3. Login to ECR
                    sh """
                        aws ecr get-login-password --region ${env.AWS_REGION} | docker login --username AWS --password-stdin ${ECR_REGISTRY}
                    """
        
                    echo "✅ ECR login successful"
        
                    // 4. Check if repository exists or create it
                    sh """
                        aws ecr describe-repositories --repository-names ${env.ECR_REPOSITORY} --region ${env.AWS_REGION} || \
                        aws ecr create-repository --repository-name ${env.ECR_REPOSITORY} --region ${env.AWS_REGION}
                    """
        
                    echo "✅ ECR repository verified/created"
        
                    // 5. Tag and push Docker image
                    def dockerCmd = env.DOCKER_CMD ?: "docker" // or sudo docker if needed
                    def imageName = "${env.ECR_REPOSITORY}"
                    def imageTag = "${env.BUILD_NUMBER}"
                    def fullImageUri = "${ECR_REGISTRY}/${imageName}:${imageTag}"
                    def latestImageUri = "${ECR_REGISTRY}/${imageName}:latest"
                    
                    sh """
                        # Build the image with both tags
                        ${dockerCmd} build -t ${imageName}:${imageTag} -t ${imageName}:latest .
                    
                        # Tag both for ECR
                        ${dockerCmd} tag ${imageName}:${imageTag} ${fullImageUri}
                        ${dockerCmd} tag ${imageName}:latest ${latestImageUri}
                    
                        # Push both tags to ECR
                        ${dockerCmd} push ${fullImageUri}
                        ${dockerCmd} push ${latestImageUri}
                    """
        
                    echo "✅ Image pushed to ECR successfully"
                }
            }
        }

        
        stage('Configure Kubectl') {
            steps {
                script {
                    env.CURRENT_STAGE = 'KUBECTL_CONFIG'
                    echo "=== STAGE: ${env.CURRENT_STAGE} ==="
                    
                    try {
                        echo "Configuring kubectl for EKS cluster..."
                        
                        // Update kubeconfig for EKS
                       withEnv(["KUBECONFIG=${env.KUBECONFIG}"]) {
                        sh """
                            echo "== PATH =="
                            echo \$PATH
                            echo "== AWS CLI =="
                            which aws && aws --version
                            echo "== Kubeconfig =="
                            aws sts get-caller-identity
                            aws eks --region us-east-1 update-kubeconfig --name ${env.EKS_CLUSTER_NAME} --kubeconfig ${env.KUBECONFIG}
                            echo "== Kubeconfig current context =="
                            kubectl config current-context
                            echo "== Nodes =="
                            kubectl get nodes
                        """
                    }
                        echo "✅ kubectl configured successfully"
                        echo "✅ KUBECTL_CONFIG stage completed successfully"
                        
                    } catch (Exception e) {
                        env.CURRENT_STAGE = 'KUBECTL_CONFIG_FAILED'
                        echo "❌ KUBECTL_CONFIG stage failed: ${e.getMessage()}"
                        throw e
                    }
                }
            }
        }
        
        stage('Deploy to EKS') {
            steps {
                script {
                    env.CURRENT_STAGE = 'EKS_DEPLOY'
                    echo "=== STAGE: ${env.CURRENT_STAGE} ==="
                    
                    try {
                        echo "Deploying to EKS cluster..."
                        
                        // Create Kubernetes deployment manifest
                        writeFile file: 'k8s-deployment.yaml', text: """
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cfx-go-app
  labels:
    app: cfx-go-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: cfx-go-app
  template:
    metadata:
      labels:
        app: cfx-go-app
    spec:
      containers:
      - name: cfx-go-app
        image: ${env.FULL_IMAGE_URI}
        ports:
        - containerPort: 8080
        env:
        - name: BUILD_NUMBER
          value: "${env.BUILD_NUMBER}"
        - name: GIT_COMMIT
          value: "${env.GIT_COMMIT}"
        - name: GO_ENV
          value: "production"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: cfx-go-service
spec:
  selector:
    app: cfx-go-app
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
"""
                        
                        // Apply deployment
                        sh """
                            export KUBECONFIG=${env.KUBECONFIG}
                            kubectl apply -f k8s-deployment.yaml
                            kubectl rollout status deployment/cfx-go-app --timeout=300s
                            kubectl get deployments
                            kubectl get services
                            kubectl get pods
                        """
                        
                        // Get service URL
                        try {
                            def serviceUrl = sh(
                                returnStdout: true,
                                script: """
                                    export KUBECONFIG=${env.KUBECONFIG}
                                    kubectl get service cfx-go-service -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
                                """
                            ).trim()
                            
                            if (serviceUrl) {
                                echo "🌐 Application URL: http://${serviceUrl}"
                                env.APP_URL = "http://${serviceUrl}"
                            }
                        } catch (Exception e) {
                            echo "⚠️ Could not get service URL: ${e.getMessage()}"
                        }
                        
                        echo "✅ EKS_DEPLOY stage completed successfully"
                        env.CURRENT_STAGE = 'ALL_STAGES_COMPLETED'
                        
                    } catch (Exception e) {
                        env.CURRENT_STAGE = 'EKS_DEPLOY_FAILED'
                        echo "❌ EKS_DEPLOY stage failed: ${e.getMessage()}"
                        throw e
                    }
                }
            }
        }
    }
    
    post {
        always {
            script {
                try {
                    echo "=== BUILD SUMMARY ==="
                    echo "Build Number: ${env.BUILD_NUMBER}"
                    echo "Git Commit: ${env.GIT_COMMIT}"
                    echo "Image Tag: ${env.IMAGE_TAG}"
                    echo "ECR Repository: ${env.ECR_REPOSITORY}"
                    echo "EKS Cluster: ${env.EKS_CLUSTER_NAME}"
                    echo "AWS Region: ${env.AWS_REGION}"
                    echo "Go Version: ${env.GO_VERSION}"
                    echo "Build Status: ${currentBuild.currentResult}"
                    echo "Final Stage: ${env.CURRENT_STAGE}"
                    if (env.APP_URL) {
                        echo "Application URL: ${env.APP_URL}"
                    }
                    echo "========================="
                    
                    // Cleanup
                    def dockerCmd = env.DOCKER_CMD ?: "sudo docker"
                    sh "${dockerCmd} system prune -f || true"
                    sh 'rm -f ${env.KUBECONFIG} || true'
                    sh 'rm -f main || true'  // Clean up Go binary (matches your Dockerfile)
                    
                } catch (Exception e) {
                    echo "Error in post always: ${e.getMessage()}"
                }
            }
        }
        
        success {
            script {
                try {
                    echo "🎉 Go application deployment completed successfully!"
                    echo "✅ Application: ${env.ECR_REPOSITORY}"
                    echo "✅ Version: ${env.IMAGE_TAG}"
                    echo "✅ Cluster: ${env.EKS_CLUSTER_NAME}"
                    echo "✅ Region: ${env.AWS_REGION}"
                    echo "✅ Go Version: ${env.GO_VERSION}"
                    if (env.APP_URL) {
                        echo "🌐 Access your Go application at: ${env.APP_URL}"
                    }
                } catch (Exception e) {
                    echo "Error in post success: ${e.getMessage()}"
                }
            }
        }
        
        failure {
            script {
                try {
                    echo "❌ Go application deployment failed!"
                    echo "💥 FAILURE SUMMARY:"
                    echo "- Application: ${env.ECR_REPOSITORY}"
                    echo "- Version: ${env.IMAGE_TAG}"
                    echo "- Cluster: ${env.EKS_CLUSTER_NAME}"
                    echo "- Region: ${env.AWS_REGION}"
                    echo "- Go Version: ${env.GO_VERSION}"
                    echo "- Build URL: ${env.BUILD_URL}"
                    echo "- Failed Stage: ${env.CURRENT_STAGE}"
                    
                } catch (Exception e) {
                    echo "Error in post failure: ${e.getMessage()}"
                }
            }
        }
    }
}
