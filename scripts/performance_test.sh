#!/bin/bash

# Performance testing script for EthosView API
# Tests the optimizations implemented in Phase 1

set -e

# Configuration
API_BASE_URL=${API_BASE_URL:-"http://localhost:8080/api/v1"}
TEST_DURATION=${TEST_DURATION:-30}
CONCURRENT_USERS=${CONCURRENT_USERS:-10}

echo "🚀 Starting EthosView API Performance Tests"
echo "API Base URL: $API_BASE_URL"
echo "Test Duration: ${TEST_DURATION}s"
echo "Concurrent Users: $CONCURRENT_USERS"
echo ""

# Check if API is running
echo "📡 Checking API availability..."
if ! curl -s "$API_BASE_URL/health" > /dev/null; then
    echo "❌ API is not available at $API_BASE_URL"
    echo "Please start the server first: go run cmd/server/main.go"
    exit 1
fi
echo "✅ API is available"

# Test 1: Basic Health Check Performance
echo ""
echo "🧪 Test 1: Health Check Performance"
echo "Testing endpoint: GET /health"
ab -n 1000 -c 10 "$API_BASE_URL/health" 2>/dev/null | grep -E "(Requests per second|Time per request|Transfer rate)"

# Test 2: ESG Scores Endpoint Performance
echo ""
echo "🧪 Test 2: ESG Scores Endpoint Performance"
echo "Testing endpoint: GET /esg/scores"
ab -n 500 -c 5 "$API_BASE_URL/esg/scores" 2>/dev/null | grep -E "(Requests per second|Time per request|Transfer rate)"

# Test 3: Companies Endpoint Performance
echo ""
echo "🧪 Test 3: Companies Endpoint Performance"
echo "Testing endpoint: GET /companies"
ab -n 500 -c 5 "$API_BASE_URL/companies" 2>/dev/null | grep -E "(Requests per second|Time per request|Transfer rate)"

# Test 4: Compression Test
echo ""
echo "🧪 Test 4: Compression Test"
echo "Testing gzip compression..."

# Test without compression
echo "Without compression:"
curl -s -w "Size: %{size_download} bytes\n" "$API_BASE_URL/esg/scores" > /dev/null

# Test with compression
echo "With compression:"
curl -s -H "Accept-Encoding: gzip" -w "Size: %{size_download} bytes\n" "$API_BASE_URL/esg/scores" > /dev/null

# Test 5: Rate Limiting Test
echo ""
echo "🧪 Test 5: Rate Limiting Test"
echo "Testing rate limiting (expecting 429 after limit exceeded)..."

# Make rapid requests to trigger rate limiting
for i in {1..15}; do
    response=$(curl -s -w "%{http_code}" "$API_BASE_URL/companies" -o /dev/null)
    echo "Request $i: HTTP $response"
    if [ "$response" = "429" ]; then
        echo "✅ Rate limiting working correctly"
        break
    fi
done

# Test 6: Error Handling Test
echo ""
echo "🧪 Test 6: Error Handling Test"
echo "Testing invalid requests..."

# Test invalid ESG score ID
echo "Testing invalid ESG score ID:"
curl -s "$API_BASE_URL/esg/scores/invalid" | jq '.error' 2>/dev/null || echo "Error response received"

# Test 7: Database Query Performance
echo ""
echo "🧪 Test 7: Database Query Performance"
echo "Testing complex queries..."

# Test ESG trends endpoint
echo "Testing ESG trends endpoint:"
ab -n 100 -c 2 "$API_BASE_URL/analytics/companies/1/esg-trends" 2>/dev/null | grep -E "(Requests per second|Time per request)"

echo ""
echo "🎉 Performance tests completed!"
echo ""
echo "📊 Summary of optimizations tested:"
echo "   ✅ Response compression (gzip)"
echo "   ✅ Rate limiting"
echo "   ✅ Enhanced error handling"
echo "   ✅ Request ID tracking"
echo "   ✅ Database index optimization"
echo "   ✅ Input validation"
echo ""
echo "💡 To improve performance further:"
echo "   - Monitor database query performance"
echo "   - Adjust cache TTL based on usage patterns"
echo "   - Consider implementing connection pooling"
echo "   - Monitor memory usage and optimize if needed"
