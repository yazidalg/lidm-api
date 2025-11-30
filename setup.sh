#!/bin/bash

# Setup script to create necessary directories and set permissions

echo "=========================================="
echo "Setting up LIDM API directories"
echo "=========================================="
echo ""

# Create uploads directories
echo "Creating upload directories..."
mkdir -p uploads/icons
mkdir -p uploads/profiles
mkdir -p uploads/images

# Set permissions
echo "Setting permissions..."
chmod 755 uploads
chmod 755 uploads/icons
chmod 755 uploads/profiles
chmod 755 uploads/images

echo ""
echo "✅ Directories created:"
ls -la uploads/

echo ""
echo "=========================================="
echo "Setup completed successfully!"
echo "=========================================="
echo ""
echo "You can now run the application:"
echo "  ./main"
echo ""
