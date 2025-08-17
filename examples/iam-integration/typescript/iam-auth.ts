/**
 * Hanzo Runtime + Hanzo IAM Integration Example
 * 
 * This example shows how to authenticate with Hanzo IAM (Casdoor)
 * and use the token to access Hanzo Runtime services.
 */

import { HanzoRuntime } from '@hanzo/runtime';
import axios from 'axios';

// Configuration
const IAM_BASE_URL = process.env.HANZO_IAM_URL || 'https://iam.hanzo.ai';
const IAM_CLIENT_ID = process.env.HANZO_IAM_CLIENT_ID || 'your-client-id';
const IAM_CLIENT_SECRET = process.env.HANZO_IAM_CLIENT_SECRET || 'your-client-secret';
const RUNTIME_API_URL = process.env.HANZO_RUNTIME_API_URL || 'https://api.hanzo.ai';

/**
 * Authenticate with Hanzo IAM using OAuth2 Client Credentials flow
 */
async function authenticateWithIAM(): Promise<string> {
  try {
    // Get access token from IAM using client credentials
    const tokenResponse = await axios.post(
      `${IAM_BASE_URL}/api/login/oauth/access_token`,
      new URLSearchParams({
        grant_type: 'client_credentials',
        client_id: IAM_CLIENT_ID,
        client_secret: IAM_CLIENT_SECRET,
        scope: 'runtime:full' // Request full access to Runtime API
      }),
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      }
    );

    return tokenResponse.data.access_token;
  } catch (error) {
    console.error('Failed to authenticate with IAM:', error);
    throw error;
  }
}

/**
 * Validate token with IAM (optional - for service-to-service validation)
 */
async function validateToken(token: string): Promise<boolean> {
  try {
    const introspectResponse = await axios.post(
      `${IAM_BASE_URL}/api/login/oauth/introspect`,
      new URLSearchParams({
        token: token,
        token_type_hint: 'access_token'
      }),
      {
        auth: {
          username: IAM_CLIENT_ID,
          password: IAM_CLIENT_SECRET
        },
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      }
    );

    return introspectResponse.data.active === true;
  } catch (error) {
    console.error('Token validation failed:', error);
    return false;
  }
}

/**
 * Create Hanzo Runtime client with IAM authentication
 */
async function createAuthenticatedRuntimeClient(): Promise<HanzoRuntime> {
  // Get token from IAM
  const accessToken = await authenticateWithIAM();
  
  // Optionally validate the token
  const isValid = await validateToken(accessToken);
  if (!isValid) {
    throw new Error('Token validation failed');
  }

  // Create Runtime client with the IAM token
  const client = new HanzoRuntime({
    apiUrl: RUNTIME_API_URL,
    apiKey: accessToken, // Use IAM token as API key
    // Optionally set custom headers
    headers: {
      'Authorization': `Bearer ${accessToken}`,
      'X-Auth-Provider': 'hanzo-iam'
    }
  });

  return client;
}

/**
 * Example: Using Runtime SDK with IAM authentication
 */
async function main() {
  try {
    console.log('🔐 Authenticating with Hanzo IAM...');
    const client = await createAuthenticatedRuntimeClient();
    console.log('✅ Successfully authenticated!');

    // Now use the Runtime SDK as normal
    console.log('🚀 Creating sandbox...');
    const sandbox = await client.sandbox.getOrCreate('iam-test-sandbox');
    
    console.log('📦 Sandbox created:', sandbox.id);
    
    // Execute commands
    const result = await sandbox.commands.run('echo "Hello from IAM-authenticated sandbox!"');
    console.log('Output:', result.output);

    // Work with files
    await sandbox.files.write('/test.txt', 'Authenticated via Hanzo IAM');
    const content = await sandbox.files.read('/test.txt');
    console.log('File content:', content);

    // Clean up
    await sandbox.kill();
    console.log('✅ Sandbox terminated');

  } catch (error) {
    console.error('❌ Error:', error);
  }
}

// Token refresh middleware for long-running operations
class IAMTokenRefresher {
  private token: string | null = null;
  private tokenExpiry: Date | null = null;

  async getToken(): Promise<string> {
    // Check if token needs refresh (5 minutes before expiry)
    if (!this.token || !this.tokenExpiry || 
        new Date().getTime() > this.tokenExpiry.getTime() - 5 * 60 * 1000) {
      await this.refreshToken();
    }
    return this.token!;
  }

  private async refreshToken(): Promise<void> {
    const tokenResponse = await axios.post(
      `${IAM_BASE_URL}/api/login/oauth/access_token`,
      new URLSearchParams({
        grant_type: 'client_credentials',
        client_id: IAM_CLIENT_ID,
        client_secret: IAM_CLIENT_SECRET,
        scope: 'runtime:full'
      }),
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      }
    );

    this.token = tokenResponse.data.access_token;
    // Set expiry based on expires_in (default 1 hour)
    const expiresIn = tokenResponse.data.expires_in || 3600;
    this.tokenExpiry = new Date(Date.now() + expiresIn * 1000);
  }
}

// Export for use in other modules
export {
  authenticateWithIAM,
  validateToken,
  createAuthenticatedRuntimeClient,
  IAMTokenRefresher
};

// Run if called directly
if (require.main === module) {
  main();
}