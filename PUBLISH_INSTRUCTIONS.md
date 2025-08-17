# Hanzo Runtime Package Publishing Instructions

## Overview

All packages have been prepared for publishing with version 0.7.0. Since npm requires OTP authentication, you'll need to publish manually.

## NPM Packages (JavaScript/TypeScript)

All packages use the @hanzo scope and require npm authentication with OTP.

### 1. @hanzo/api-client
```bash
cd libs/api-client
npm publish --access public --otp=YOUR_OTP_CODE
```

### 2. @hanzo/runner-api-client  
```bash
cd libs/runner-api-client
npm publish --access public --otp=YOUR_OTP_CODE
```

### 3. @hanzo/runtime (Main TypeScript SDK)
```bash
cd libs/sdk-typescript
npm publish --access public --otp=YOUR_OTP_CODE
```

## Python Packages (PyPI)

All Python packages require PyPI authentication token (PYPI_TOKEN).

### 1. hanzo-runtime (Main Python SDK)
```bash
cd libs/sdk-python
poetry build
poetry publish --username __token__ --password $PYPI_TOKEN
```

### 2. hanzo_runtime_api_client
```bash
cd libs/api-client-python
poetry build
poetry publish --username __token__ --password $PYPI_TOKEN
```

### 3. hanzo_runtime_api_client_async
```bash
cd libs/api-client-python-async
poetry build
poetry publish --username __token__ --password $PYPI_TOKEN
```

## Quick Publish All (with tokens ready)

### NPM (requires OTP for each package):
```bash
# From project root
export NPM_TOKEN=your_npm_token_here

# Publish each package (you'll be prompted for OTP)
for pkg in api-client runner-api-client sdk-typescript; do
  cd libs/$pkg
  npm publish --access public
  cd ../..
done
```

### Python:
```bash
# From project root
export PYPI_TOKEN=your_pypi_token_here

# Build and publish all Python packages
for pkg in sdk-python api-client-python api-client-python-async; do
  cd libs/$pkg
  poetry build
  poetry publish --username __token__ --password $PYPI_TOKEN
  cd ../..
done
```

## Post-Publishing Verification

After publishing, verify the packages are available:

### NPM:
```bash
npm view @hanzo/api-client
npm view @hanzo/runner-api-client
npm view @hanzo/runtime
```

### PyPI:
```bash
pip index versions hanzo-runtime
pip index versions hanzo-runtime-api-client
pip index versions hanzo-runtime-api-client-async
```

## Integration

Once published, users can install:

### JavaScript/TypeScript:
```bash
npm install @hanzo/runtime @hanzo/api-client
# or
yarn add @hanzo/runtime @hanzo/api-client
# or 
pnpm add @hanzo/runtime @hanzo/api-client
```

### Python:
```bash
pip install hanzo-runtime
# or for API clients specifically:
pip install hanzo-runtime-api-client
pip install hanzo-runtime-api-client-async
```

## Usage Examples

### TypeScript:
```typescript
import { HanzoRuntime } from '@hanzo/runtime';

const client = new HanzoRuntime({ apiKey: 'your-api-key' });
const sandbox = await client.sandbox.getOrCreate('my-sandbox');
```

### Python:
```python
from hanzo_runtime import HanzoRuntime

client = HanzoRuntime(api_key='your-api-key')
sandbox = client.sandbox.get_or_create('my-sandbox')
```