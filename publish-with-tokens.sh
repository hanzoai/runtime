#!/bin/bash
# Quick publish script for Hanzo Runtime packages
# Run this after setting your tokens as environment variables

set -e

echo "🚀 Publishing Hanzo Runtime Packages"
echo "===================================="
echo ""

# Check tokens
if [ -z "$NPM_OTP" ]; then
    echo "⚠️  NPM_OTP not set. You'll need to enter OTP manually for each npm package."
    echo "   Set with: export NPM_OTP=123456"
    echo ""
fi

if [ -z "$PYPI_TOKEN" ]; then
    echo "❌ PYPI_TOKEN not set. Please set it first:"
    echo "   export PYPI_TOKEN=pypi-..."
    exit 1
fi

# NPM Packages
echo "📦 Publishing NPM packages..."
echo ""

# API Client
echo "Publishing @hanzo/api-client..."
cd libs/api-client
if [ -n "$NPM_OTP" ]; then
    npm publish --access public --otp="$NPM_OTP"
else
    npm publish --access public
fi
cd ../..

# Runner API Client
echo "Publishing @hanzo/runner-api-client..."
cd libs/runner-api-client
if [ -n "$NPM_OTP" ]; then
    npm publish --access public --otp="$NPM_OTP"
else
    npm publish --access public
fi
cd ../..

# TypeScript SDK
echo "Publishing @hanzo/runtime..."
cd libs/sdk-typescript
if [ -n "$NPM_OTP" ]; then
    npm publish --access public --otp="$NPM_OTP"
else
    npm publish --access public
fi
cd ../..

# Python Packages
echo ""
echo "🐍 Publishing Python packages..."
echo ""

# Main SDK
echo "Publishing hanzo-runtime..."
cd libs/sdk-python
python -m build
python -m twine upload dist/* --username __token__ --password "$PYPI_TOKEN"
cd ../..

# API Client
echo "Publishing hanzo_runtime_api_client..."
cd libs/api-client-python
python -m build
python -m twine upload dist/* --username __token__ --password "$PYPI_TOKEN"
cd ../..

# Async API Client
echo "Publishing hanzo_runtime_api_client_async..."
cd libs/api-client-python-async
python -m build
python -m twine upload dist/* --username __token__ --password "$PYPI_TOKEN"
cd ../..

echo ""
echo "✅ All done! Verify with:"
echo "   npm view @hanzo/runtime"
echo "   pip index versions hanzo-runtime"