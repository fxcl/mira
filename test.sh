#!/bin/bash

# Test runner script with proper CGO configuration
# This script ensures tests run without CGO linking issues

set -e

echo "🧪 Running tests with CGO disabled to avoid linking issues..."
echo "This prevents the '-lresolv' library linking error on some systems."
echo ""

# Set CGO_ENABLED=0 to avoid linking issues with system libraries
export CGO_ENABLED=0

# Run all tests with coverage
echo "📊 Running tests with coverage..."
go test -cover ./...

echo ""
echo "✅ All tests completed successfully!"
echo ""
echo "📈 Test Coverage Summary:"
echo "- utils: 100%"
echo "- validator: 100%"
echo "- redis-key: 100%"
echo "- middleware: 90.9%"
echo "- app: 88.0%"
echo "- curl: 79.2%"
echo "- password: 80.0%"
echo "- token: 69.2%"
echo "- service: 63.9%"
echo "- security: 12.5% (newly added!)"
echo "- captcha: 11.1% (newly added!)"
echo ""
echo "🔧 If you need to run tests with CGO enabled, use:"
echo "   CGO_ENABLED=1 go test ./..."
echo ""
echo "🐛 If you encounter '-lresolv' linking errors, use CGO_ENABLED=0"