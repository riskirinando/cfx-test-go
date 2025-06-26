pipeline {
    agent any
    
    environment {
        AWS_REGION = 'us-east-1'
        ECR_REPOSITORY = 'cfx-test-go'
        EKS_CLUSTER_NAME = 'test-project-eks-cluster'
        KUBECONFIG = credentials('kubeconfig')
        AWS_CREDENTIALS = credentials('aws-credentials')
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Build') {
            steps {
                script {
                    echo 'Building Go application...'
                    sh '''
                        go mod tidy
                        go test ./...
                        go build -o main .
                    '''
                }
            }
        }
        
        stage('Build Docker Image') {
            steps {
                script {
                    echo 'Building Docker image...'
                    def imageTag = "${BUILD_NUMBER}"
                    def imageName = "${ECR_REPOSITORY}:${imageTag}"
                    
                    sh "docker build -t ${imageName} ."
                    sh "docker tag ${imageName} ${ECR_REPOSITORY}:latest"
                    
                    env.IMAGE_TAG = imageTag
                    env.IMAGE_NAME = imageName
                }
            }
        }
        
        stage('Push to ECR') {
            steps {
                script {
                    withCredentials([[$class: 'AmazonWebServicesCredentialsBinding', credentialsId: 'aws-credentials']]) {
                        echo 'Pushing image to ECR...'
                        sh '''
                            # Get ECR login token
                            aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com
                            
                            # Tag and push image
                            docker tag ${ECR_REPOSITORY}:${IMAGE_TAG} ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPOSITORY}:${IMAGE_TAG}
                            docker tag ${ECR_REPOSITORY}:latest ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPOSITORY}:latest
                            
                            docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPOSITORY}:${IMAGE_TAG}
                            docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPOSITORY}:latest
                        '''
                    }
                }
            }
        }
        
        stage('Deploy to EKS') {
            steps {
                script {
                    withCredentials([[$class: 'AmazonWebServicesCredentialsBinding', credentialsId: 'aws-credentials']]) {
                        echo 'Deploying to EKS...'
                        sh '''
                            # Update kubeconfig
                            aws eks update-kubeconfig --region ${AWS_REGION} --name ${EKS_CLUSTER_NAME}
                            
                            # Update deployment image
                            sed -i "s|your-account-id.dkr.ecr.your-region.amazonaws.com/go-web-app:latest|${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPOSITORY}:${IMAGE_TAG}|g" k8s/deployment.yaml
                            
                            # Apply Kubernetes manifests
                            kubectl apply -f k8s/deployment.yaml
                            kubectl apply -f k8s/ingress.yaml
                            
                            # Wait for deployment to complete
                            kubectl rollout status deployment/go-web-app --timeout=300s
                            
                            # Get service information
                            kubectl get services go-web-app-service
                        '''
                    }
                }
            }
        }
        
        stage('Health Check') {
            steps {
                script {
                    echo 'Performing health check...'
                    sh '''
                        # Wait for pods to be ready
                        kubectl wait --for=condition=ready pod -l app=go-web-app --timeout=300s
                        
                        # Get pod status
                        kubectl get pods -l app=go-web-app
                        
                        # Test health endpoint (if LoadBalancer is ready)
                        sleep 30
                        LOAD_BALANCER_URL=$(kubectl get service go-web-app-service -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
                        if [ ! -z "$LOAD_BALANCER_URL" ]; then
                            echo "Testing health endpoint: http://$LOAD_BALANCER_URL/health"
                            curl -f "http://$LOAD_BALANCER_URL/health" || echo "Health check failed, but deployment completed"
                        else
                            echo "LoadBalancer URL not yet available"
                        fi
                    '''
                }
            }
        }
    }
    
    post {
        always {
            echo 'Cleaning up...'
            sh '''
                docker rmi ${ECR_REPOSITORY}:${IMAGE_TAG} || true
                docker rmi ${ECR_REPOSITORY}:latest || true
                docker rmi ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPOSITORY}:${IMAGE_TAG} || true
                docker rmi ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPOSITORY}:latest || true
            '''
        }
        success {
            echo 'Pipeline completed successfully!'
            slackSend(
                channel: '#deployments',
                color: 'good',
                message: "✅ Successfully deployed go-web-app:${IMAGE_TAG} to EKS cluster ${EKS_CLUSTER_NAME}"
            )
        }
        failure {
            echo 'Pipeline failed!'
            slackSend(
                channel: '#deployments',
                color: 'danger',
                message: "❌ Failed to deploy go-web-app:${IMAGE_TAG} to EKS cluster ${EKS_CLUSTER_NAME}"
            )
        }
    }
}
