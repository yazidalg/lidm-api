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

# Build and push the Docker image
echo "📦 Building Docker image..."
docker build -t ${IMAGE_NAME} .

echo "📤 Pushing image to Google Container Registry..."
docker push ${IMAGE_NAME}

# Deploy to Cloud Run
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
  --set-env-vars "ENV=production" \
  --set-env-vars "PORT=8080"

echo "✅ Deployment completed!"
echo "🌐 Service URL: https://${SERVICE_NAME}-${PROJECT_ID}.a.run.app"

# Test the deployment
echo "🧪 Testing deployment..."
SERVICE_URL="https://${SERVICE_NAME}-${PROJECT_ID}.a.run.app"
curl -f "${SERVICE_URL}/health" || echo "⚠️  Health check failed"

echo "🎉 Deployment script completed!"
