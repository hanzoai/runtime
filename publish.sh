#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Hanzo Runtime Publishing Script${NC}"
echo "================================"

# Check if tokens are provided as arguments or in environment
if [ -n "$1" ] && [ -n "$2" ]; then
    export NPM_TOKEN="$1"
    export PYPI_TOKEN="$2"
    echo -e "${GREEN}✓ Using tokens from command line arguments${NC}"
elif [ -z "$NPM_TOKEN" ] || [ -z "$PYPI_TOKEN" ]; then
    echo -e "${RED}Error: Tokens not found${NC}"
    echo ""
    echo "Usage:"
    echo "  1. Set environment variables:"
    echo "     export NPM_TOKEN='npm_...'"
    echo "     export PYPI_TOKEN='pypi-...'"
    echo "     ./publish.sh"
    echo ""
    echo "  2. Or pass as arguments:"
    echo "     ./publish.sh 'npm_...' 'pypi-...'"
    echo ""
    exit 1
else
    echo -e "${GREEN}✓ Using tokens from environment${NC}"
fi

# Run make publish
echo ""
echo "Running make publish..."
make publish

echo ""
echo -e "${GREEN}✅ Publishing complete!${NC}"
echo ""
echo "Packages published:"
echo "  - Python: https://pypi.org/project/hanzo-runtime/"
echo "  - npm: https://www.npmjs.com/package/@hanzo/runtime"