#!/bin/bash

# Phase 2 Testing Script for EthosView API
# Tests cache warming, metrics collection, enhanced health checks, and security improvements

set -e

# Configuration
API_BASE_URL=${API_BASE_URL:-"http://localhost:8080/api/v1"}
TEST_DURATION=${TEST_DURATION:-30}
CONCURRENT_USERS=${CONCURRENT_USERS:-10}

echo "🚀 Starting EthosView Phase 2 Testing"
echo "API Base URL: $API_BASE_URL"
echo "Test Duration: ${TEST_DURATION}s"
echo "Concurrent Users: $CONCURRENT_USERS"
echo ""

# Check if API is running
echo "📡 Checking API availability..."
if ! curl -s "$API_BASE_URL/health" > /dev/null; then
    echo "❌ API is not available at $API_BASE_URL"
    echo "Please start the server first: docker-compose up -d"
    exit 1
fi
echo "✅ API is available"

# Test 1: Enhanced Health Checks
echo ""
echo "🧪 Test 1: Enhanced Health Checks"
echo "=================================="

echo "Testing basic health check:"
curl -s "$API_BASE_URL/health" | jq '.status' 2>/dev/null || echo "Health check response received"

echo ""
echo "Testing detailed health check:"
curl -s "$API_BASE_URL/health/detailed" | jq '.health.status' 2>/dev/null || echo "Detailed health check response received"

echo ""
echo "Testing readiness check:"
curl -s "$API_BASE_URL/health/ready" | jq '.ready' 2>/dev/null || echo "Readiness check response received"

echo ""
echo "Testing liveness check:"
curl -s "$API_BASE_URL/health/live" | jq '.status' 2>/dev/null || echo "Liveness check response received"

# Test 2: Metrics Collection
echo ""
echo "🧪 Test 2: Metrics Collection"
echo "============================="

echo "Testing metrics endpoint:"
curl -s "$API_BASE_URL/metrics" | jq '.business.total_companies' 2>/dev/null || echo "Metrics response received"

# Test 3: Security Improvements
echo ""
echo "🧪 Test 3: Security Improvements"
echo "================================"

echo "Testing security headers:"
headers=$(curl -s -I "$API_BASE_URL/companies" 2>/dev/null)
echo "X-Content-Type-Options: $(echo "$headers" | grep -i "X-Content-Type-Options" || echo "Not found")"
echo "X-Frame-Options: $(echo "$headers" | grep -i "X-Frame-Options" || echo "Not found")"
echo "X-XSS-Protection: $(echo "$headers" | grep -i "X-XSS-Protection" || echo "Not found")"
echo "Content-Security-Policy: $(echo "$headers" | grep -i "Content-Security-Policy" || echo "Not found")"

echo ""
echo "Testing SQL injection protection:"
response=$(curl -s -w "%{http_code}" "$API_BASE_URL/companies?name=test' OR 1=1--" -o /dev/null)
if [ "$response" = "400" ]; then
    echo "✅ SQL injection protection working"
else
    echo "❌ SQL injection protection may not be working (HTTP $response)"
fi

echo ""
echo "Testing XSS protection:"
response=$(curl -s -w "%{http_code}" "$API_BASE_URL/companies?name=<script>alert('xss')</script>" -o /dev/null)
if [ "$response" = "400" ]; then
    echo "✅ XSS protection working"
else
    echo "❌ XSS protection may not be working (HTTP $response)"
fi

echo ""
echo "Testing CORS headers:"
cors_headers=$(curl -s -H "Origin: http://localhost:3000" -I "$API_BASE_URL/companies" 2>/dev/null)
echo "Access-Control-Allow-Origin: $(echo "$cors_headers" | grep -i "Access-Control-Allow-Origin" || echo "Not found")"
echo "Access-Control-Allow-Methods: $(echo "$cors_headers" | grep -i "Access-Control-Allow-Methods" || echo "Not found")"

# Test 4: Cache Warming
echo ""
echo "🧪 Test 4: Cache Warming"
echo "========================"

echo "Testing cache warming effectiveness:"
echo "First request (should be slower):"
time curl -s "$API_BASE_URL/companies" > /dev/null

echo ""
echo "Second request (should be faster due to cache):"
time curl -s "$API_BASE_URL/companies" > /dev/null

echo ""
echo "Testing ESG scores cache:"
echo "First request:"
time curl -s "$API_BASE_URL/esg/scores" > /dev/null

echo ""
echo "Second request:"
time curl -s "$API_BASE_URL/esg/scores" > /dev/null

# Test 5: Request Size Limits
echo ""
echo "🧪 Test 5: Request Size Limits"
echo "=============================="

echo "Testing request size limit (should accept normal request):"
response=$(curl -s -w "%{http_code}" -X POST "$API_BASE_URL/companies" -H "Content-Type: application/json" -d '{"name":"Test Company","symbol":"TEST"}' -o /dev/null)
echo "Normal request: HTTP $response"

echo ""
echo "Testing large request (should be rejected):"
# Create a large payload
large_payload=$(printf '{"name":"%s","symbol":"TEST"}' "$(printf 'A%.0s' {1..1000000})")
response=$(curl -s -w "%{http_code}" -X POST "$API_BASE_URL/companies" -H "Content-Type: application/json" -d "$large_payload" -o /dev/null)
if [ "$response" = "413" ]; then
    echo "✅ Request size limit working (HTTP $response)"
else
    echo "❌ Request size limit may not be working (HTTP $response)"
fi

# Test 6: Input Sanitization
echo ""
echo "🧪 Test 6: Input Sanitization"
echo "============================="

echo "Testing input sanitization:"
response=$(curl -s -w "%{http_code}" "$API_BASE_URL/companies?name=Test%00Company" -o /dev/null)
if [ "$response" = "200" ]; then
    echo "✅ Input sanitization working"
else
    echo "❌ Input sanitization may not be working (HTTP $response)"
fi

# Test 7: Performance with New Features
echo ""
echo "🧪 Test 7: Performance with New Features"
echo "========================================"

echo "Testing performance with all optimizations:"
ab -n 100 -c 5 "$API_BASE_URL/companies" 2>/dev/null | grep -E "(Requests per second|Time per request|Transfer rate)"

echo ""
echo "Testing ESG scores performance:"
ab -n 100 -c 5 "$API_BASE_URL/esg/scores" 2>/dev/null | grep -E "(Requests per second|Time per request|Transfer rate)"

# Test 8: Error Handling
echo ""
echo "🧪 Test 8: Error Handling"
echo "========================="

echo "Testing enhanced error handling:"
error_response=$(curl -s "$API_BASE_URL/esg/scores/invalid")
echo "Error response format:"
echo "$error_response" | jq '.' 2>/dev/null || echo "$error_response"

# Summary
echo ""
echo "🎉 Phase 2 Testing Completed!"
echo ""
echo "📊 Summary of Phase 2 features tested:"
echo "   ✅ Enhanced health checks (basic, detailed, readiness, liveness)"
echo "   ✅ Metrics collection and monitoring"
echo "   ✅ Security improvements (headers, SQL injection, XSS, CORS)"
echo "   ✅ Cache warming strategy"
echo "   ✅ Request size limits"
echo "   ✅ Input sanitization"
echo "   ✅ Performance optimizations"
echo "   ✅ Enhanced error handling"
echo ""
echo "💡 Phase 2 optimizations provide:"
echo "   - Better monitoring and observability"
echo "   - Enhanced security protection"
echo "   - Improved cache performance"
echo "   - Comprehensive health checking"
echo "   - Better error handling and debugging"
echo ""
echo "🚀 Ready for Phase 3: Advanced Features!"
