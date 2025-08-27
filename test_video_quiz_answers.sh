#!/bin/bash

# Video Quiz Answer Testing Script
# Usage: ./test_video_quiz_answers.sh [JWT_TOKEN]

# Configuration
BASE_URL="http://localhost:8080"
JWT_TOKEN="${1:-YOUR_JWT_TOKEN_HERE}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_section() {
    echo -e "\n${BLUE}=====================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}=====================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Check if JWT token is provided
if [ "$JWT_TOKEN" = "YOUR_JWT_TOKEN_HERE" ]; then
    print_error "Please provide a valid JWT token as argument or edit the script"
    echo "Usage: ./test_video_quiz_answers.sh YOUR_JWT_TOKEN"
    exit 1
fi

print_section "VIDEO QUIZ ANSWER TESTING"
print_info "Base URL: $BASE_URL"
print_info "Using JWT Token: ${JWT_TOKEN:0:20}..."

# ============================================
# TEST MODULE 2 VIDEO QUIZ
# ============================================

print_section "MODULE 2 - VIDEO QUIZ TESTING"

print_info "Testing Quiz 1 (Module 2) - Correct Answer (A)"
response=$(curl -s -X POST "$BASE_URL/api/video-quiz/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "video_quiz_id": 1,
    "selected_answer": "A",
    "response_time": 15
  }')

if echo "$response" | grep -q "successfully"; then
    print_success "Module 2 Quiz 1 - Correct answer submitted"
    echo "$response" | jq .
else
    print_error "Module 2 Quiz 1 - Failed to submit answer"
    echo "$response"
fi

sleep 2

print_info "Testing Quiz 1 (Module 2) - Wrong Answer (B) for comparison"
response=$(curl -s -X POST "$BASE_URL/api/video-quiz/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "video_quiz_id": 1,
    "selected_answer": "B",
    "response_time": 20
  }')

if echo "$response" | grep -q "successfully"; then
    print_success "Module 2 Quiz 1 - Wrong answer submitted (for testing)"
    echo "$response" | jq .
else
    print_error "Module 2 Quiz 1 - Failed to submit wrong answer"
    echo "$response"
fi

# ============================================
# TEST MODULE 4 VIDEO QUIZZES
# ============================================

print_section "MODULE 4 - VIDEO QUIZ TESTING"

# Quiz 1 Module 4
print_info "Testing Quiz 1 (Module 4) - Question about photosynthesis requirements"
response=$(curl -s -X POST "$BASE_URL/api/video-quiz/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "video_quiz_id": 2,
    "selected_answer": "A",
    "response_time": 12
  }')

if echo "$response" | grep -q "successfully"; then
    print_success "Module 4 Quiz 1 - Answer submitted"
    echo "$response" | jq .
else
    print_error "Module 4 Quiz 1 - Failed to submit answer"
    echo "$response"
fi

sleep 2

# Quiz 2 Module 4
print_info "Testing Quiz 2 (Module 4) - Question about photosynthesis results"
response=$(curl -s -X POST "$BASE_URL/api/video-quiz/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "video_quiz_id": 3,
    "selected_answer": "B",
    "response_time": 18
  }')

if echo "$response" | grep -q "successfully"; then
    print_success "Module 4 Quiz 2 - Answer submitted"
    echo "$response" | jq .
else
    print_error "Module 4 Quiz 2 - Failed to submit answer"
    echo "$response"
fi

sleep 2

# Quiz 3 Module 4
print_info "Testing Quiz 3 (Module 4) - Question about photosynthesis importance"
response=$(curl -s -X POST "$BASE_URL/api/video-quiz/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "video_quiz_id": 4,
    "selected_answer": "B",
    "response_time": 25
  }')

if echo "$response" | grep -q "successfully"; then
    print_success "Module 4 Quiz 3 - Answer submitted"
    echo "$response" | jq .
else
    print_error "Module 4 Quiz 3 - Failed to submit answer"
    echo "$response"
fi

# ============================================
# GET USER ANSWERS
# ============================================

print_section "RETRIEVING USER ANSWERS"

print_info "Getting all user video quiz answers"
response=$(curl -s -X GET "$BASE_URL/api/video-quiz/user-answers" \
  -H "Authorization: Bearer $JWT_TOKEN")

if echo "$response" | grep -q "successfully"; then
    print_success "All user answers retrieved"
    echo "$response" | jq .
else
    print_error "Failed to retrieve user answers"
    echo "$response"
fi

sleep 2

print_info "Getting user answers for Module 2 (Video Material ID: 3)"
response=$(curl -s -X GET "$BASE_URL/api/video-quiz/user-answers/3" \
  -H "Authorization: Bearer $JWT_TOKEN")

if echo "$response" | grep -q "successfully"; then
    print_success "Module 2 answers retrieved"
    echo "$response" | jq .
else
    print_error "Failed to retrieve Module 2 answers"
    echo "$response"
fi

sleep 2

print_info "Getting user answers for Module 4 (Video Material ID: 5)"
response=$(curl -s -X GET "$BASE_URL/api/video-quiz/user-answers/5" \
  -H "Authorization: Bearer $JWT_TOKEN")

if echo "$response" | grep -q "successfully"; then
    print_success "Module 4 answers retrieved"
    echo "$response" | jq .
else
    print_error "Failed to retrieve Module 4 answers"
    echo "$response"
fi

print_section "VIDEO QUIZ TESTING COMPLETED"
print_success "All tests completed! Check the responses above for results."
