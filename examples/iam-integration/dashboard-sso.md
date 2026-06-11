# Hanzo Runtime Dashboard SSO Integration

This guide shows how to integrate the Runtime Dashboard with Hanzo IAM for Single Sign-On (SSO).

## Overview

The Runtime Dashboard uses OIDC (OpenID Connect) for authentication, which integrates seamlessly with Hanzo IAM (Casdoor).

## Configuration Steps

### 1. Configure IAM Application for Dashboard

In Hanzo IAM, create or update the Runtime Dashboard application:

```json
{
  "name": "runtime-dashboard",
  "displayName": "Hanzo Runtime Dashboard",
  "logo": "https://runtime.hanzo.ai/logo.png",
  "homepageUrl": "https://runtime.hanzo.ai",
  "redirectUris": [
    "https://runtime.hanzo.ai",
    "https://runtime.hanzo.ai/callback",
    "http://localhost:5173", // For local development
    "http://localhost:5173/callback"
  ],
  "grantTypes": ["authorization_code", "refresh_token"],
  "responseTypes": ["code", "id_token"],
  "tokenFormat": "JWT-Standard",
  "expireInHours": 24,
  "refreshExpireInHours": 168 // 7 days
}
```

### 2. Configure Dashboard Environment

Create `.env.local` file in the dashboard directory:

```bash
# OIDC Configuration
VITE_OIDC_DOMAIN=https://iam.hanzo.ai
VITE_OIDC_CLIENT_ID=runtime-dashboard
VITE_OIDC_AUDIENCE=https://api.hanzo.ai

# API Configuration
VITE_API_URL=https://api.hanzo.ai
VITE_APP_TITLE=Hanzo Runtime

# Optional: Billing API
VITE_BILLING_API_URL=https://billing.hanzo.ai
```

### 3. IAM Login Page Customization

Add a "Runtime Dashboard" button to IAM's application list or home page:

```html
<!-- In IAM's application template -->
<div class="app-card" onclick="loginToRuntime()">
  <img src="/static/img/runtime-logo.png" alt="Runtime">
  <h3>Runtime Dashboard</h3>
  <p>Manage your AI sandboxes and environments</p>
</div>

<script>
function loginToRuntime() {
  // Direct OAuth flow to Runtime Dashboard
  const params = new URLSearchParams({
    client_id: 'runtime-dashboard',
    redirect_uri: 'https://runtime.hanzo.ai',
    response_type: 'code',
    scope: 'openid profile email offline_access',
    state: JSON.stringify({ returnTo: '/dashboard' })
  });
  
  window.location.href = `/oauth/authorize?${params}`;
}
</script>
```

### 4. Update Dashboard Authentication

The dashboard already uses OIDC, but we can enhance it for better IAM integration:

```typescript
// apps/dashboard/src/auth/iam-enhanced-config.ts
import { AuthProviderProps } from 'react-oidc-context'
import { RoutePath } from '@/enums/RoutePath'

export const iamEnhancedConfig: AuthProviderProps = {
  authority: import.meta.env.VITE_OIDC_DOMAIN,
  client_id: import.meta.env.VITE_OIDC_CLIENT_ID,
  scope: 'openid profile email offline_access runtime:full',
  redirect_uri: window.location.origin,
  
  // Enhanced token handling
  automaticSilentRenew: true,
  revokeTokensOnSignout: true,
  
  // Custom metadata for IAM
  metadata: {
    // IAM-specific endpoints
    end_session_endpoint: `${import.meta.env.VITE_OIDC_DOMAIN}/oauth/logout`,
    revocation_endpoint: `${import.meta.env.VITE_OIDC_DOMAIN}/oauth/revoke`,
  },
  
  // Handle post-login redirect
  onSigninCallback: (user) => {
    // Check if user has access to Runtime
    if (!user?.profile?.organizations?.includes('runtime')) {
      window.location.href = '/no-access'
      return
    }
    
    const state = user?.state as { returnTo?: string }
    const targetUrl = state?.returnTo || RoutePath.DASHBOARD
    window.history.replaceState({}, '', targetUrl)
  },
  
  // Handle logout
  post_logout_redirect_uri: `${import.meta.env.VITE_OIDC_DOMAIN}/logout?redirect=${window.location.origin}`,
}
```

### 5. API Authentication Middleware

Update the Runtime API to validate IAM tokens:

```typescript
// apps/api/src/auth/iam-jwt.strategy.ts
import { Injectable } from '@nestjs/common'
import { PassportStrategy } from '@nestjs/passport'
import { ExtractJwt, Strategy } from 'passport-jwt'
import { passportJwtSecret } from 'jwks-rsa'

@Injectable()
export class IamJwtStrategy extends PassportStrategy(Strategy, 'iam-jwt') {
  constructor() {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      
      // Use IAM's JWKS endpoint for key rotation
      secretOrKeyProvider: passportJwtSecret({
        cache: true,
        rateLimit: true,
        jwksRequestsPerMinute: 5,
        jwksUri: `${process.env.IAM_ENDPOINT}/.well-known/jwks`
      }),
      
      // Validate issuer and audience
      issuer: process.env.IAM_ISSUER,
      audience: process.env.HANZO_RUNTIME_API_URL,
      algorithms: ['RS256'],
    })
  }

  async validate(payload: any) {
    return {
      userId: payload.sub,
      username: payload.preferred_username,
      email: payload.email,
      organizations: payload.organizations || [],
      permissions: payload.permissions || [],
    }
  }
}
```

### 6. Enhanced User Experience

Create a seamless flow from IAM to Dashboard:

```typescript
// apps/dashboard/src/components/IamIntegration.tsx
import { useAuth } from 'react-oidc-context'
import { useEffect } from 'react'

export function IamIntegration() {
  const auth = useAuth()
  
  useEffect(() => {
    // Check if coming from IAM
    const urlParams = new URLSearchParams(window.location.search)
    const fromIam = urlParams.get('from') === 'iam'
    
    if (fromIam && !auth.isAuthenticated) {
      // Auto-trigger login if coming from IAM
      auth.signinRedirect({
        state: { fromIam: true }
      })
    }
  }, [auth])
  
  // Show personalized welcome for IAM users
  if (auth.isAuthenticated && auth.user?.state?.fromIam) {
    return (
      <div className="welcome-banner">
        <h2>Welcome from Hanzo IAM!</h2>
        <p>You're now in the Runtime Dashboard. Your sandboxes are ready.</p>
      </div>
    )
  }
  
  return null
}
```

## Complete Integration Flow

1. **User logs into IAM** → Sees available applications
2. **Clicks "Runtime Dashboard"** → OAuth flow initiated
3. **IAM validates credentials** → Issues JWT token
4. **Redirects to Dashboard** → Dashboard validates token
5. **Dashboard loads** → Shows user's sandboxes and resources

## Security Features

- **Token Validation**: All API calls validate JWT signatures
- **Scope Enforcement**: Fine-grained permissions per operation
- **Session Management**: Automatic token refresh
- **Audit Logging**: All actions logged with user context

## Testing the Integration

1. **Local Development**:
   ```bash
   # Start IAM
   cd ~/work/hanzo/iam
   ./run_local.sh
   
   # Start Runtime API
   cd ~/work/hanzo/runtime
   pnpm nx serve api
   
   # Start Dashboard
   pnpm nx serve dashboard
   ```

2. **Test SSO Flow**:
   - Navigate to http://localhost:8000 (IAM)
   - Click "Runtime Dashboard"
   - Should redirect to http://localhost:5173
   - Verify authentication works

3. **Test API Access**:
   ```bash
   # Get token from IAM
   TOKEN=$(curl -X POST http://localhost:8000/oauth/token \
     -d "grant_type=password" \
     -d "username=admin" \
     -d "password=admin" | jq -r .access_token)
   
   # Use with Runtime API
   curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/api/sandboxes
   ```

## Production Deployment

1. **Configure DNS**:
   - `iam.hanzo.ai` → IAM service
   - `runtime.hanzo.ai` → Dashboard
   - `api.runtime.hanzo.ai` → Runtime API

2. **SSL Certificates**:
   - Use wildcard cert for `*.hanzo.ai`
   - Or individual certs per subdomain

3. **Environment Variables**:
   ```yaml
   # docker-compose.prod.yml
   services:
     dashboard:
       environment:
         VITE_OIDC_DOMAIN: https://iam.hanzo.ai
         VITE_API_URL: https://api.runtime.hanzo.ai
   ```

## Benefits

- **Single Sign-On**: One login for all Hanzo services
- **Centralized User Management**: All users managed in IAM
- **Consistent Experience**: Same auth flow everywhere
- **Enhanced Security**: JWT tokens with short expiry
- **Scalable**: Works with multiple Runtime instances