#!/bin/bash

# Run load tests with k6
# Usage: ./run-load-tests.sh [test-type] [base-url]

set -e

TEST_TYPE=${1:-"api-load-test"}
BASE_URL=${2:-"http://localhost:8080"}

echo "🚀 Running k6 load tests..."
echo "Test: $TEST_TYPE"
echo "Base URL: $BASE_URL"
echo ""

# Check if k6 is installed
if ! command -v k6 &> /dev/null; then
    echo "❌ k6 is not installed. Install it from: https://k6.io/docs/getting-started/installation/"
    exit 1
fi

# Check if the application is running
echo "🔍 Checking if application is running..."
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo "❌ Application is not running at $BASE_URL"
    echo "Start it with: docker-compose up"
    exit 1
fi

echo "✅ Application is running"
echo ""

# Run the test
echo "📊 Starting load test..."
export BASE_URL=$BASE_URL
k6 run "loadtests/$TEST_TYPE.js"

echo ""
echo "✅ Load test completed!"
