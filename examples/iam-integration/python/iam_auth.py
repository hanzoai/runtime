"""
Hanzo Runtime + Hanzo IAM Integration Example

This example shows how to authenticate with Hanzo IAM (Casdoor)
and use the token to access Hanzo Runtime services.
"""

import os
import asyncio
from datetime import datetime, timedelta
from typing import Optional
import httpx
from hanzo_runtime import HanzoRuntime
from hanzo_runtime._async import HanzoRuntime as AsyncHanzoRuntime

# Configuration
IAM_BASE_URL = os.getenv('IAM_ENDPOINT', 'https://iam.hanzo.ai')
IAM_CLIENT_ID = os.getenv('IAM_CLIENT_ID', 'your-client-id')
IAM_CLIENT_SECRET = os.getenv('IAM_CLIENT_SECRET', 'your-client-secret')
RUNTIME_API_URL = os.getenv('HANZO_RUNTIME_API_URL', 'https://api.hanzo.ai')


class HanzoIAMAuth:
    """Handle authentication with Hanzo IAM"""
    
    def __init__(self, base_url: str = IAM_BASE_URL, 
                 client_id: str = IAM_CLIENT_ID,
                 client_secret: str = IAM_CLIENT_SECRET):
        self.base_url = base_url
        self.client_id = client_id
        self.client_secret = client_secret
        self.token: Optional[str] = None
        self.token_expiry: Optional[datetime] = None
    
    def authenticate(self) -> str:
        """Authenticate with Hanzo IAM using OAuth2 Client Credentials flow"""
        with httpx.Client() as client:
            response = client.post(
                f"{self.base_url}/v1/iam/oauth/token",
                data={
                    'grant_type': 'client_credentials',
                    'client_id': self.client_id,
                    'client_secret': self.client_secret,
                    'scope': 'runtime:full'  # Request full access to Runtime API
                },
                headers={'Content-Type': 'application/x-www-form-urlencoded'}
            )
            response.raise_for_status()
            
            data = response.json()
            self.token = data['access_token']
            
            # Set token expiry (default 1 hour)
            expires_in = data.get('expires_in', 3600)
            self.token_expiry = datetime.now() + timedelta(seconds=expires_in)
            
            return self.token
    
    def validate_token(self, token: str) -> bool:
        """Validate token with IAM (optional - for service-to-service validation)"""
        with httpx.Client() as client:
            response = client.post(
                f"{self.base_url}/v1/iam/oauth/introspect",
                data={
                    'token': token,
                    'token_type_hint': 'access_token'
                },
                auth=(self.client_id, self.client_secret),
                headers={'Content-Type': 'application/x-www-form-urlencoded'}
            )
            
            if response.status_code != 200:
                return False
            
            data = response.json()
            return data.get('active', False)
    
    def get_token(self) -> str:
        """Get current token, refreshing if needed"""
        # Check if token needs refresh (5 minutes before expiry)
        if (not self.token or not self.token_expiry or 
            datetime.now() > self.token_expiry - timedelta(minutes=5)):
            self.authenticate()
        
        return self.token


def create_authenticated_runtime_client() -> HanzoRuntime:
    """Create Hanzo Runtime client with IAM authentication"""
    # Initialize IAM auth
    iam_auth = HanzoIAMAuth()
    
    # Get token from IAM
    access_token = iam_auth.authenticate()
    
    # Optionally validate the token
    if not iam_auth.validate_token(access_token):
        raise ValueError("Token validation failed")
    
    # Create Runtime client with the IAM token
    client = HanzoRuntime(
        api_url=RUNTIME_API_URL,
        api_key=access_token,  # Use IAM token as API key
        # Optionally set custom headers
        headers={
            'Authorization': f'Bearer {access_token}',
            'X-Auth-Provider': 'hanzo-iam'
        }
    )
    
    return client


async def create_authenticated_async_runtime_client() -> AsyncHanzoRuntime:
    """Create async Hanzo Runtime client with IAM authentication"""
    # Initialize IAM auth
    iam_auth = HanzoIAMAuth()
    
    # Get token from IAM
    access_token = iam_auth.authenticate()
    
    # Create async Runtime client
    client = AsyncHanzoRuntime(
        api_url=RUNTIME_API_URL,
        api_key=access_token,
        headers={
            'Authorization': f'Bearer {access_token}',
            'X-Auth-Provider': 'hanzo-iam'
        }
    )
    
    return client


def main():
    """Example: Using Runtime SDK with IAM authentication"""
    try:
        print("🔐 Authenticating with Hanzo IAM...")
        client = create_authenticated_runtime_client()
        print("✅ Successfully authenticated!")
        
        # Now use the Runtime SDK as normal
        print("🚀 Creating sandbox...")
        sandbox = client.sandbox.get_or_create("iam-test-sandbox")
        
        print(f"📦 Sandbox created: {sandbox.id}")
        
        # Execute commands
        result = sandbox.commands.run('echo "Hello from IAM-authenticated sandbox!"')
        print(f"Output: {result.output}")
        
        # Work with files
        sandbox.files.write("/test.txt", "Authenticated via Hanzo IAM")
        content = sandbox.files.read("/test.txt")
        print(f"File content: {content}")
        
        # Clean up
        sandbox.kill()
        print("✅ Sandbox terminated")
        
    except Exception as e:
        print(f"❌ Error: {e}")


async def async_main():
    """Async example with IAM authentication"""
    try:
        print("🔐 Authenticating with Hanzo IAM (async)...")
        client = await create_authenticated_async_runtime_client()
        print("✅ Successfully authenticated!")
        
        # Use async Runtime SDK
        sandbox = await client.sandbox.get_or_create("iam-test-async")
        print(f"📦 Sandbox created: {sandbox.id}")
        
        # Execute async operations
        result = await sandbox.commands.run('uname -a')
        print(f"System info: {result.output}")
        
        await sandbox.kill()
        print("✅ Sandbox terminated")
        
    except Exception as e:
        print(f"❌ Error: {e}")
    finally:
        await client.close()


# Custom authentication handler for automatic token refresh
class IAMAuthenticatedRuntime(HanzoRuntime):
    """Runtime client with automatic IAM token refresh"""
    
    def __init__(self, iam_auth: HanzoIAMAuth, **kwargs):
        self.iam_auth = iam_auth
        # Get initial token
        token = iam_auth.get_token()
        
        # Initialize parent with token
        super().__init__(
            api_key=token,
            headers={
                'Authorization': f'Bearer {token}',
                'X-Auth-Provider': 'hanzo-iam'
            },
            **kwargs
        )
    
    def _refresh_auth_if_needed(self):
        """Refresh authentication token if needed"""
        token = self.iam_auth.get_token()
        self.api_key = token
        self.headers['Authorization'] = f'Bearer {token}'
    
    # Override methods that make API calls to refresh token first
    def request(self, *args, **kwargs):
        self._refresh_auth_if_needed()
        return super().request(*args, **kwargs)


if __name__ == "__main__":
    # Run sync example
    main()
    
    # Run async example
    # asyncio.run(async_main())