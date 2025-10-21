#!/bin/bash

# Simple test script to verify the application works locally
echo "🧪 Testing LIDM API locally..."

# Set environment variables for testing
export ENV=development
export PORT=8080
export DB_USER=root
export DB_PASSWORD=""
export DB_NAME=lidm_db
export DB_HOST=127.0.0.1
export DB_PORT=3306

echo "📋 Environment variables set:"
echo "  ENV: $ENV"
echo "  PORT: $PORT"
echo "  DB_USER: $DB_USER"
echo "  DB_NAME: $DB_NAME"
echo "  DB_HOST: $DB_HOST"
echo "  DB_PORT: $DB_PORT"

# Build the application
echo "🔨 Building application..."
go build -o lidm-test ./cmd

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful"

# Start the application in background
echo "🚀 Starting application..."
./lidm-test &
APP_PID=$!

# Wait for the application to start
echo "⏳ Waiting for application to start..."
sleep 5

# Test health endpoint
echo "🏥 Testing health endpoint..."
HEALTH_RESPONSE=$(curl -s http://localhost:8080/health)
if [ $? -eq 0 ]; then
    echo "✅ Health endpoint responded: $HEALTH_RESPONSE"
else
    echo "❌ Health endpoint failed"
    kill $APP_PID 2>/dev/null
    exit 1
fi

# Test root endpoint
echo "🏠 Testing root endpoint..."
ROOT_RESPONSE=$(curl -s http://localhost:8080/)
if [ $? -eq 0 ]; then
    echo "✅ Root endpoint responded: $ROOT_RESPONSE"
else
    echo "❌ Root endpoint failed"
    kill $APP_PID 2>/dev/null
    exit 1
fi

# Test ready endpoint (if available)
echo "🔍 Testing ready endpoint..."
READY_RESPONSE=$(curl -s http://localhost:8080/ready)
if [ $? -eq 0 ]; then
    echo "✅ Ready endpoint responded: $READY_RESPONSE"
else
    echo "⚠️  Ready endpoint not available (this is OK for basic mode)"
fi

# Clean up
echo "🧹 Cleaning up..."
kill $APP_PID 2>/dev/null
rm -f lidm-test

echo "✅ Local test completed successfully!"
echo "🎉 Application is ready for Cloud Run deployment"
