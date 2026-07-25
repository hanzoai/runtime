/**
 * Hanzo Runtime API Middleware for IAM Integration
 * 
 * This middleware validates tokens from Hanzo IAM before allowing
 * access to Runtime API endpoints.
 */

import { Request, Response, NextFunction } from 'express';
import axios from 'axios';

interface IAMConfig {
  baseUrl: string;
  clientId: string;
  clientSecret: string;
  cacheTTL?: number; // Token validation cache TTL in seconds
}

interface TokenInfo {
  active: boolean;
  scope?: string;
  client_id?: string;
  username?: string;
  exp?: number;
  sub?: string;
  aud?: string[];
}

class IAMTokenValidator {
  private cache: Map<string, { info: TokenInfo; expires: number }> = new Map();
  
  constructor(private config: IAMConfig) {
    // Clean up expired cache entries every minute
    setInterval(() => this.cleanupCache(), 60000);
  }

  /**
   * Validate token with Hanzo IAM
   */
  async validateToken(token: string): Promise<TokenInfo> {
    // Check cache first
    const cached = this.cache.get(token);
    if (cached && cached.expires > Date.now()) {
      return cached.info;
    }

    try {
      // Introspect token with IAM
      const response = await axios.post(
        `${this.config.baseUrl}/v1/iam/oauth/introspect`,
        new URLSearchParams({
          token: token,
          token_type_hint: 'access_token'
        }),
        {
          auth: {
            username: this.config.clientId,
            password: this.config.clientSecret
          },
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
          },
          timeout: 5000 // 5 second timeout
        }
      );

      const tokenInfo = response.data as TokenInfo;
      
      // Cache the result
      if (tokenInfo.active) {
        const cacheTTL = this.config.cacheTTL || 300; // Default 5 minutes
        this.cache.set(token, {
          info: tokenInfo,
          expires: Date.now() + (cacheTTL * 1000)
        });
      }

      return tokenInfo;
    } catch (error) {
      console.error('Token validation error:', error);
      return { active: false };
    }
  }

  /**
   * Clean up expired cache entries
   */
  private cleanupCache(): void {
    const now = Date.now();
    for (const [token, entry] of this.cache.entries()) {
      if (entry.expires <= now) {
        this.cache.delete(token);
      }
    }
  }
}

/**
 * Create Express middleware for IAM authentication
 */
export function createIAMAuthMiddleware(config: IAMConfig) {
  const validator = new IAMTokenValidator(config);

  return async (req: Request, res: Response, next: NextFunction) => {
    try {
      // Extract token from Authorization header
      const authHeader = req.headers.authorization;
      if (!authHeader || !authHeader.startsWith('Bearer ')) {
        return res.status(401).json({
          error: 'Unauthorized',
          message: 'Missing or invalid authorization header'
        });
      }

      const token = authHeader.substring(7); // Remove 'Bearer ' prefix

      // Validate token with IAM
      const tokenInfo = await validator.validateToken(token);

      if (!tokenInfo.active) {
        return res.status(401).json({
          error: 'Unauthorized',
          message: 'Invalid or expired token'
        });
      }

      // Check required scopes if specified
      const requiredScope = (req as any).requiredScope;
      if (requiredScope && tokenInfo.scope) {
        const tokenScopes = tokenInfo.scope.split(' ');
        if (!tokenScopes.includes(requiredScope)) {
          return res.status(403).json({
            error: 'Forbidden',
            message: `Insufficient scope. Required: ${requiredScope}`
          });
        }
      }

      // Attach token info to request for use in route handlers
      (req as any).auth = {
        tokenInfo,
        clientId: tokenInfo.client_id,
        username: tokenInfo.username,
        subject: tokenInfo.sub
      };

      next();
    } catch (error) {
      console.error('Auth middleware error:', error);
      res.status(500).json({
        error: 'Internal Server Error',
        message: 'Authentication processing failed'
      });
    }
  };
}

/**
 * Scope checking middleware
 */
export function requireScope(scope: string) {
  return (req: Request, res: Response, next: NextFunction) => {
    (req as any).requiredScope = scope;
    next();
  };
}

/**
 * Example: Using the middleware in Runtime API
 */
export function setupRuntimeAPIWithIAM(app: any) {
  // Configure IAM integration
  const iamConfig: IAMConfig = {
    baseUrl: process.env.IAM_ENDPOINT || 'https://iam.hanzo.ai',
    clientId: process.env.IAM_CLIENT_ID!,
    clientSecret: process.env.IAM_CLIENT_SECRET!,
    cacheTTL: 300 // Cache validations for 5 minutes
  };

  // Apply IAM auth middleware to all API routes
  const iamAuth = createIAMAuthMiddleware(iamConfig);
  app.use('/api', iamAuth);

  // Example routes with scope requirements
  app.get('/api/sandboxes', requireScope('runtime:read'), (req: Request, res: Response) => {
    // Access authenticated user info
    const auth = (req as any).auth;
    console.log(`User ${auth.username} listing sandboxes`);
    
    // ... sandbox listing logic
  });

  app.post('/api/sandboxes', requireScope('runtime:write'), (req: Request, res: Response) => {
    const auth = (req as any).auth;
    console.log(`User ${auth.username} creating sandbox`);
    
    // ... sandbox creation logic
  });

  // Admin routes require special scope
  app.delete('/api/sandboxes/:id', requireScope('runtime:admin'), (req: Request, res: Response) => {
    // ... sandbox deletion logic
  });

  // Health check endpoint (no auth required)
  app.get('/health', (req: Request, res: Response) => {
    res.json({ status: 'ok', iam: 'enabled' });
  });
}

/**
 * Helper to extract user context from IAM token
 */
export interface UserContext {
  id: string;
  username: string;
  email?: string;
  organizations?: string[];
  permissions?: string[];
}

export async function getUserContext(tokenInfo: TokenInfo): Promise<UserContext> {
  // Extract user information from token
  // In a real implementation, you might need to call IAM's user info endpoint
  return {
    id: tokenInfo.sub || '',
    username: tokenInfo.username || '',
    // Additional fields would be populated from IAM user info endpoint
  };
}