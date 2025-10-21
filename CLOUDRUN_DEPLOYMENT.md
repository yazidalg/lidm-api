# Cloud Run Environment Configuration Guide

## Required Environment Variables

### Database Configuration
```bash
# For Cloud SQL (Production)
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_database_name
INSTANCE_CONNECTION_NAME=project:region:instance

# Alternative: Use DB_HOST for TCP connection
# DB_HOST=your-cloud-sql-ip
# DB_PORT=3306
```

### Application Configuration
```bash
ENV=production
PORT=8080
```

### Optional Configuration
```bash
# For quiz seeding (optional)
SEED_QUIZ_MODULES="Module Title 1,Module Title 2"

# For Google Auth (if using)
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
```

## Cloud Run Deployment Commands

### 1. Set Environment Variables
```bash
export PROJECT_ID="your-project-id"
export SERVICE_NAME="lidm-api"
export REGION="asia-southeast1"
```

### 2. Deploy with Environment Variables
```bash
gcloud run deploy ${SERVICE_NAME} \
  --image gcr.io/${PROJECT_ID}/${SERVICE_NAME} \
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
  --set-env-vars "DB_USER=${DB_USER},DB_PASSWORD=${DB_PASSWORD},DB_NAME=${DB_NAME}" \
  --set-env-vars "INSTANCE_CONNECTION_NAME=${INSTANCE_CONNECTION_NAME}"
```

### 3. Test Deployment
```bash
# Health check
curl https://${SERVICE_NAME}-${PROJECT_ID}.a.run.app/health

# Ready check
curl https://${SERVICE_NAME}-${PROJECT_ID}.a.run.app/ready
```

## Troubleshooting

### Common Issues

1. **Container fails to start**: Check environment variables are set correctly
2. **Database connection fails**: Verify Cloud SQL instance connection name
3. **Port binding issues**: Ensure PORT=8080 is set
4. **Memory issues**: Increase memory allocation if needed

### Debug Commands
```bash
# View logs
gcloud run services logs read ${SERVICE_NAME} --region=${REGION}

# Check service status
gcloud run services describe ${SERVICE_NAME} --region=${REGION}
```

## Health Check Endpoints

- `/health` - Basic health check (always returns 200)
- `/ready` - Readiness check (includes database status)
- `/healthy` - Liveness check (simple database ping)
