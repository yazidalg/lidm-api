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
export INSTANCE_CONNECTION_NAME="project:region:instance"  # Your Cloud SQL instance connection name
```

### 2. **IMPORTANT: Connect Cloud SQL to Cloud Run Service**

**This is the most critical step!** You must connect your Cloud SQL instance to your Cloud Run service:

```bash
# First, find your Cloud SQL instance connection name
gcloud sql instances describe <instance-name> --format="value(connectionName)"

# Then connect it to your Cloud Run service
gcloud run services update ${SERVICE_NAME} \
  --region ${REGION} \
  --add-cloudsql-instances ${INSTANCE_CONNECTION_NAME}
```

**Note:** The `--add-cloudsql-instances` flag is what enables Unix socket connections. Without this, TCP connections will time out!

### 4. Deploy with Environment Variables

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
  --add-cloudsql-instances ${INSTANCE_CONNECTION_NAME} \
  --set-env-vars "ENV=production,PORT=8080" \
  --set-env-vars "DB_USER=${DB_USER},DB_PASSWORD=${DB_PASSWORD},DB_NAME=${DB_NAME}" \
  --set-env-vars "INSTANCE_CONNECTION_NAME=${INSTANCE_CONNECTION_NAME}"
```

**Key Points:**
- **DO NOT set `DB_HOST`** if you're using Unix socket (recommended)
- **DO set `INSTANCE_CONNECTION_NAME`** environment variable
- **DO include `--add-cloudsql-instances`** flag in deployment

### 5. Test Deployment
```bash
# Health check
curl https://${SERVICE_NAME}-${PROJECT_ID}.a.run.app/health

# Ready check
curl https://${SERVICE_NAME}-${PROJECT_ID}.a.run.app/ready
```

## Troubleshooting

### Common Issues

1. **Container fails to start**: Check environment variables are set correctly
2. **Database connection timeout (connection timed out)**: 
   - **This is the most common issue!**
   - Make sure you've connected Cloud SQL to Cloud Run using `--add-cloudsql-instances`
   - Remove `DB_HOST` environment variable if set (it forces TCP connection)
   - Ensure `INSTANCE_CONNECTION_NAME` is set correctly
   - Verify the connection name format: `project:region:instance`
3. **Port binding issues**: Ensure PORT=8080 is set
4. **Memory issues**: Increase memory allocation if needed

### Fixing Connection Timeout Error

If you see errors like:
```
dial tcp 34.101.93.176:3306: connect: connection timed out
```

**Solution:**
1. Remove `DB_HOST` environment variable from your Cloud Run service:
   ```bash
   gcloud run services update ${SERVICE_NAME} \
     --region ${REGION} \
     --update-env-vars "DB_HOST="
   ```

2. Make sure `INSTANCE_CONNECTION_NAME` is set:
   ```bash
   gcloud run services update ${SERVICE_NAME} \
     --region ${REGION} \
     --update-env-vars "INSTANCE_CONNECTION_NAME=${INSTANCE_CONNECTION_NAME}"
   ```

3. Connect Cloud SQL instance to Cloud Run (if not already done):
   ```bash
   gcloud run services update ${SERVICE_NAME} \
     --region ${REGION} \
     --add-cloudsql-instances ${INSTANCE_CONNECTION_NAME}
   ```

4. Redeploy your service

### Debug Commands
```bash
# View logs
gcloud run services logs read ${SERVICE_NAME} --region=${REGION}

# Check service status and environment variables
gcloud run services describe ${SERVICE_NAME} --region=${REGION}

# List environment variables
gcloud run services describe ${SERVICE_NAME} --region=${REGION} --format="value(spec.template.spec.containers[0].env)"

# Check Cloud SQL connection
gcloud run services describe ${SERVICE_NAME} --region=${REGION} --format="value(spec.template.spec.containers[0].env)" | grep CLOUDSQL
```

## Health Check Endpoints

- `/health` - Basic health check (always returns 200)
- `/ready` - Readiness check (includes database status)
- `/healthy` - Liveness check (simple database ping)
