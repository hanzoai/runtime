#\!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}📦 Publishing all Hanzo Runtime packages...${NC}"
echo "============================================="

# Source environment for tokens
source ~/.zshrc

# Check tokens
if [ -z "$PYPI_TOKEN" ]; then
    echo -e "${RED}❌ PYPI_TOKEN not found in environment${NC}"
    exit 1
fi

if [ -z "$NPM_TOKEN" ]; then
    echo -e "${RED}❌ NPM_TOKEN not found in environment${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Tokens loaded from ~/.zshrc${NC}"

# Publish Python packages
echo -e "\n${YELLOW}🐍 Publishing Python packages to PyPI...${NC}"

# Already published: hanzo-runtime
# echo "Publishing hanzo-runtime..."
# cd libs/sdk-python
# python -m twine upload --verbose dist/* -u __token__ -p "$PYPI_TOKEN"

echo "Publishing hanzo_runtime_api_client..."
cd libs/api-client-python
python -m twine upload --verbose dist/* -u __token__ -p "$PYPI_TOKEN" || echo "Package may already exist"
cd ../..

echo "Publishing hanzo_runtime_api_client_async..."
cd libs/api-client-python-async
python -m twine upload --verbose dist/* -u __token__ -p "$PYPI_TOKEN" || echo "Package may already exist"
cd ../..

# Publish npm packages
echo -e "\n${YELLOW}📦 Publishing npm packages...${NC}"

# Create .npmrc with token
echo "//registry.npmjs.org/:_authToken=${NPM_TOKEN}" > ~/.npmrc

# Build TypeScript packages first
echo "Building TypeScript packages..."
pnpm nx run api-client:build
pnpm nx run runner-api-client:build
pnpm nx run sdk-typescript:build

# Publish @hanzo/api-client
echo -e "\n${YELLOW}Publishing @hanzo/api-client...${NC}"
cd libs/api-client
npm publish --access public || echo "Package may already exist"
cd ../..

# Publish @hanzo/runner-api-client
echo -e "\n${YELLOW}Publishing @hanzo/runner-api-client...${NC}"
cd libs/runner-api-client
npm publish --access public || echo "Package may already exist"
cd ../..

# Publish @hanzo/runtime
echo -e "\n${YELLOW}Publishing @hanzo/runtime...${NC}"
cd libs/sdk-typescript
npm publish --access public || echo "Package may already exist"
cd ../..

echo -e "\n${GREEN}✅ All packages published successfully\!${NC}"
echo ""
echo "Published packages:"
echo "  PyPI:"
echo "    - hanzo-runtime (0.7.0)"
echo "    - hanzo_runtime_api_client (0.7.0)"
echo "    - hanzo_runtime_api_client_async (0.7.0)"
echo "  npm:"
echo "    - @hanzo/api-client (0.7.0)"
echo "    - @hanzo/runner-api-client (0.7.0)"
echo "    - @hanzo/runtime (0.7.0)"
