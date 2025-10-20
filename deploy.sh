#!/bin/bash

# Cloud Run Deployment Script for LIDM Backend
# Make sure to set these variables before running

# Required variables - SET THESE BEFORE RUNNING
PROJECT_ID="lively-oxide-475504-p8"  # Your GCP Project ID
REGION="us-central1"                  # Your preferred region
SERVICE_NAME="lidm-backend"           # Your Cloud Run service name

# Database configuration - SET THESE BEFORE RUNNING
DB_USER="your_db_user"                # Your database username
DB_PASSWORD="your_db_password"         # Your database password
DB_NAME="your_database_name"          # Your database name
INSTANCE_CONNECTION_NAME="your-project:region:instance-name"  # Your Cloud SQL instance

echo "🚀 Deploying LIDM Backend to Cloud Run..."

# Build the container
echo "📦 Building container..."
gcloud builds submit --tag gcr.io/$PROJECT_ID/$SERVICE_NAME

# Deploy to Cloud Run
echo "🚀 Deploying to Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --image gcr.io/$PROJECT_ID/$SERVICE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --port 8080 \
  --memory 1Gi \
  --cpu 1 \
  --timeout 300 \
  --max-instances 10 \
  --set-env-vars "ENV=production,DB_USER=$DB_USER,DB_PASSWORD=$DB_PASSWORD,DB_NAME=$DB_NAME,INSTANCE_CONNECTION_NAME=$INSTANCE_CONNECTION_NAME"

echo "✅ Deployment complete!"
echo "🌐 Service URL: https://$SERVICE_NAME-$REGION-$PROJECT_ID.a.run.app"
echo "🔍 Check logs: gcloud run services logs read $SERVICE_NAME --region=$REGION"
