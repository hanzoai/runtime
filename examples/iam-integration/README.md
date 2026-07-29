# Hanzo IAM + Runtime Integration Guide

This guide shows how to integrate Hanzo IAM with Hanzo Runtime for unified authentication and authorization.

## Overview

Hanzo IAM provides centralized identity and access management, while Hanzo Runtime provides AI-powered code execution environments. This integration enables:

- **Single Sign-On (SSO)** - Users authenticate once with IAM and access all Hanzo services
- **Service-to-Service Auth** - Secure communication between microservices
- **Fine-grained Permissions** - Control access to sandboxes, volumes, and operations
- **Token Management** - Automatic token refresh and validation

## Architecture

```
┌─────────────┐      ┌─────────────┐      ┌──────────────┐
│   Client    │─────▶│  Hanzo IAM  │◀────▶│Hanzo Runtime │
│ Application │      │   (OIDC)    │      │     API      │
└─────────────┘      └─────────────┘      └──────────────┘
      │                     │                      │
      │                     ▼                      ▼
      │              ┌─────────────┐       ┌──────────────┐
      └─────────────▶│   OAuth2    │       │   Sandbox    │
                     │   Tokens    │       │  Execution   │
                     └─────────────┘       └──────────────┘
```

## Setup Instructions

### 1. Configure IAM Application

In Hanzo IAM, create an application for Runtime:

1. Log into IAM admin panel
2. Go to Applications → Add
3. Configure:
   ```
   Name: hanzo-runtime
   Client ID: runtime-client
   Client Secret: [generate secure secret]
   Redirect URIs: https://api.hanzo.ai/callback
   Grant Types: authorization_code, client_credentials
   Token Format: JWT-Standard
   ```

### 2. Configure Runtime API

Set environment variables for Runtime API:

```bash
# IAM Integration
IAM_ENDPOINT=https://iam.hanzo.ai
IAM_CLIENT_ID=runtime-client
IAM_CLIENT_SECRET=your-secret-here
IAM_VALIDATE_TOKENS=true

# Token validation cache (optional)
IAM_CACHE_TTL=300  # 5 minutes
```

### 3. Configure Scopes and Permissions

Define scopes in IAM for Runtime operations:

- `runtime:read` - Read sandboxes, volumes, snapshots
- `runtime:write` - Create and modify resources
- `runtime:execute` - Execute code in sandboxes
- `runtime:admin` - Delete resources, manage permissions

## Usage Examples

### TypeScript/JavaScript

```typescript
import { HanzoRuntime } from '@hanzo/runtime';
import { authenticateWithIAM } from './iam-auth';

// Authenticate and create client
const token = await authenticateWithIAM();
const client = new HanzoRuntime({
  apiKey: token,
  headers: {
    'Authorization': `Bearer ${token}`,
    'X-Auth-Provider': 'hanzo-iam'
  }
});

// Use the authenticated client
const sandbox = await client.sandbox.getOrCreate('my-sandbox');
await sandbox.commands.run('echo "Authenticated!"');
```

### Python

```python
from hanzo_runtime import HanzoRuntime
from iam_auth import HanzoIAMAuth

# Initialize IAM authentication
iam = HanzoIAMAuth()
token = iam.authenticate()

# Create authenticated client
client = HanzoRuntime(
    api_key=token,
    headers={
        'Authorization': f'Bearer {token}',
        'X-Auth-Provider': 'hanzo-iam'
    }
)

# Use the client
sandbox = client.sandbox.get_or_create('my-sandbox')
result = sandbox.commands.run('python --version')
print(result.output)
```

### cURL

```bash
# Get token from IAM
TOKEN=$(curl -X POST https://iam.hanzo.ai/v1/iam/oauth/token \
  -d "grant_type=client_credentials" \
  -d "client_id=runtime-client" \
  -d "client_secret=$CLIENT_SECRET" \
  -d "scope=runtime:full" | jq -r .access_token)

# Use token with Runtime API
curl -H "Authorization: Bearer $TOKEN" \
  https://api.hanzo.ai/api/sandboxes
```

## Advanced Integration

### 1. Token Refresh

For long-running operations, implement automatic token refresh:

```typescript
class TokenManager {
  private token: string | null = null;
  private expiry: Date | null = null;

  async getValidToken(): Promise<string> {
    if (!this.token || new Date() >= this.expiry) {
      await this.refresh();
    }
    return this.token;
  }

  private async refresh() {
    const response = await authenticateWithIAM();
    this.token = response.access_token;
    this.expiry = new Date(Date.now() + response.expires_in * 1000);
  }
}
```

### 2. User Context

Access user information from IAM tokens:

```typescript
// In Runtime API middleware
app.use(async (req, res, next) => {
  const token = req.headers.authorization?.replace('Bearer ', '');
  const tokenInfo = await validateWithIAM(token);
  
  req.user = {
    id: tokenInfo.sub,
    username: tokenInfo.username,
    email: tokenInfo.email,
    organizations: tokenInfo.organizations
  };
  
  next();
});
```

### 3. Fine-grained Permissions

Implement resource-level permissions:

```typescript
// Check if user can access specific sandbox
async function canAccessSandbox(userId: string, sandboxId: string): Promise<boolean> {
  // Check ownership
  const sandbox = await getSandbox(sandboxId);
  if (sandbox.ownerId === userId) return true;
  
  // Check organization membership
  const userOrgs = await getUserOrganizations(userId);
  if (userOrgs.includes(sandbox.organizationId)) return true;
  
  // Check explicit permissions
  const permissions = await getUserPermissions(userId);
  return permissions.includes(`sandbox:${sandboxId}:read`);
}
```

## Security Best Practices

1. **Always validate tokens** - Don't trust client-provided tokens without validation
2. **Use HTTPS everywhere** - Protect tokens in transit
3. **Implement token rotation** - Refresh tokens before expiry
4. **Limit token scope** - Request only needed permissions
5. **Cache validations** - Reduce IAM load but keep TTL reasonable
6. **Audit access** - Log all authenticated operations

## Troubleshooting

### Common Issues

1. **401 Unauthorized**
   - Check token is included in Authorization header
   - Verify token hasn't expired
   - Ensure client credentials are correct

2. **403 Forbidden**
   - Check token has required scopes
   - Verify user has necessary permissions
   - Check organization membership

3. **Token validation fails**
   - Ensure IAM is accessible from Runtime API
   - Check network connectivity
   - Verify client credentials for introspection

### Debug Mode

Enable debug logging:

```bash
# Runtime API
IAM_DEBUG=true

# Client SDK
export DEBUG=hanzo:*
```

## Migration Guide

If migrating from API keys to IAM:

1. Create IAM application and get credentials
2. Update client code to authenticate with IAM
3. Map existing API keys to IAM users (optional)
4. Update Runtime API to validate IAM tokens
5. Deprecate and remove old API keys

## Support

- IAM Documentation: https://docs.hanzo.ai/iam
- Runtime Documentation: https://docs.hanzo.ai/runtime
- GitHub Issues: https://github.com/hanzoai/runtime/issues