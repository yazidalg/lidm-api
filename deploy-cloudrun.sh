#!/bin/bash

# Cloud Run Deployment Script for LIDM API
# This script helps deploy the application to Google Cloud Run

set -e

# Configuration
PROJECT_ID=${PROJECT_ID:-"your-project-id"}
SERVICE_NAME=${SERVICE_NAME:-"lidm-api"}
REGION=${REGION:-"asia-southeast1"}
IMAGE_NAME="gcr.io/${PROJECT_ID}/${SERVICE_NAME}"

echo "🚀 Deploying LIDM API to Cloud Run"
echo "Project ID: ${PROJECT_ID}"
echo "Service Name: ${SERVICE_NAME}"
echo "Region: ${REGION}"
echo "Image: ${IMAGE_NAME}"

# Validate required environment variables
if [ "$PROJECT_ID" = "your-project-id" ]; then
    echo "❌ Please set PROJECT_ID environment variable"
    echo "   export PROJECT_ID=your-actual-project-id"
    exit 1
fi

# Check if gcloud is authenticated
echo "🔐 Checking gcloud authentication..."
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    echo "❌ Please authenticate with gcloud first:"
    echo "   gcloud auth login"
    exit 1
fi

# Set the project
echo "🎯 Setting project to ${PROJECT_ID}..."
gcloud config set project ${PROJECT_ID}

# Enable required APIs
echo "🔧 Enabling required APIs..."
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com

# Build and push the Docker image
echo "📦 Building Docker image..."
docker build -t ${IMAGE_NAME} .

if [ $? -ne 0 ]; then
    echo "❌ Docker build failed"
    exit 1
fi

echo "📤 Pushing image to Google Container Registry..."
docker push ${IMAGE_NAME}

if [ $? -ne 0 ]; then
    echo "❌ Docker push failed"
    exit 1
fi

# Deploy to Cloud Run with better configuration
echo "🚀 Deploying to Cloud Run..."
gcloud run deploy ${SERVICE_NAME} \
  --image ${IMAGE_NAME} \
  --platform managed \
  --region ${REGION} \
  --allow-unauthenticated \
  --port 8080 \
  --memory 1Gi \
  --cpu 1 \
  --min-instances 0 \
  --max-instances 10 \
  --timeout 300 \
  --concurrency 100 \
  --set-env-vars "ENV=production,PORT=8080" \
  --set-env-vars "DB_USER=${DB_USER:-root},DB_PASSWORD=${DB_PASSWORD:-},DB_NAME=${DB_NAME:-lidm_db}" \
  --set-env-vars "INSTANCE_CONNECTION_NAME=${INSTANCE_CONNECTION_NAME:-}" \
  --set-env-vars "DB_HOST=${DB_HOST:-},DB_PORT=${DB_PORT:-3306}"

if [ $? -ne 0 ]; then
    echo "❌ Cloud Run deployment failed"
    exit 1
fi

echo "✅ Deployment completed!"
echo "🌐 Service URL: https://${SERVICE_NAME}-${PROJECT_ID}.a.run.app"

# Test the deployment
echo "🧪 Testing deployment..."
SERVICE_URL="https://${SERVICE_NAME}-${PROJECT_ID}.a.run.app"

# Wait a bit for the service to be ready
echo "⏳ Waiting for service to be ready..."
sleep 10

# Test health endpoint
echo "🏥 Testing health endpoint..."
if curl -f "${SERVICE_URL}/health"; then
    echo "✅ Health check passed"
else
    echo "⚠️  Health check failed - service may still be starting"
fi

# Test root endpoint
echo "🏠 Testing root endpoint..."
if curl -f "${SERVICE_URL}/"; then
    echo "✅ Root endpoint working"
else
    echo "⚠️  Root endpoint failed"
fi

echo "🎉 Deployment script completed!"
echo ""
echo "📋 Next steps:"
echo "1. Check the Cloud Run console for logs: https://console.cloud.google.com/run"
echo "2. Set up your database connection if not already done"
echo "3. Test your API endpoints"
echo ""
echo "🔗 Service URL: ${SERVICE_URL}"
