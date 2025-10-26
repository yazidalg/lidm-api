#!/bin/bash

# Test script for pagination functionality
# This script tests the new /module/admin/all endpoint with pagination

BASE_URL="http://localhost:8080"
ADMIN_TOKEN="your_admin_token_here"

echo "Testing pagination for /module/admin/all endpoint"
echo "================================================"

# Test 1: Default pagination (page=1, limit=10)
echo "Test 1: Default pagination"
curl -X GET "$BASE_URL/module/admin/all" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" | jq '.'

echo -e "\n"

# Test 2: Custom pagination (page=2, limit=5)
echo "Test 2: Custom pagination (page=2, limit=5)"
curl -X GET "$BASE_URL/module/admin/all?page=2&limit=5" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" | jq '.'

echo -e "\n"

# Test 3: Large limit (should be capped at 100)
echo "Test 3: Large limit (should be capped at 100)"
curl -X GET "$BASE_URL/module/admin/all?page=1&limit=200" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" | jq '.'

echo -e "\n"

# Test 4: Invalid parameters (should use defaults)
echo "Test 4: Invalid parameters (should use defaults)"
curl -X GET "$BASE_URL/module/admin/all?page=invalid&limit=abc" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" | jq '.'

echo -e "\n"

echo "Pagination tests completed!"
echo "Expected response format:"
echo "- data: array of modules"
echo "- pagination: object with current_page, per_page, total_count, total_pages, has_next, has_previous"
