#!/bin/bash

# Example API usage script for the Golang Backend
# Make sure the server is running on localhost:8080

BASE_URL="http://localhost:8080"

echo "=== Golang Backend API Test ==="
echo ""

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s "$BASE_URL/health" | jq .
echo ""

# Create a user
echo "2. Creating a user..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/users/" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123"
  }')
echo "$CREATE_RESPONSE" | jq .

# Extract user ID from response
USER_ID=$(echo "$CREATE_RESPONSE" | jq -r '.user.id')
echo ""

# Get all users
echo "3. Getting all users..."
curl -s "$BASE_URL/api/users/" | jq .
echo ""

# Get specific user
echo "4. Getting user with ID $USER_ID..."
curl -s "$BASE_URL/api/users/$USER_ID" | jq .
echo ""

# Update user
echo "5. Updating user with ID $USER_ID..."
curl -s -X PUT "$BASE_URL/api/users/$USER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_updated",
    "email": "john_updated@example.com"
  }' | jq .
echo ""

# Update password
echo "6. Updating password for user with ID $USER_ID..."
curl -s -X PATCH "$BASE_URL/api/users/$USER_ID/password" \
  -H "Content-Type: application/json" \
  -d '{
    "password": "newpassword123"
  }' | jq .
echo ""

# Get updated user
echo "7. Getting updated user with ID $USER_ID..."
curl -s "$BASE_URL/api/users/$USER_ID" | jq .
echo ""

# Delete user
echo "8. Deleting user with ID $USER_ID..."
curl -s -X DELETE "$BASE_URL/api/users/$USER_ID" | jq .
echo ""

# Verify user is deleted
echo "9. Verifying user is deleted..."
curl -s "$BASE_URL/api/users/$USER_ID" | jq .
echo ""

echo "=== API Test Complete ==="
