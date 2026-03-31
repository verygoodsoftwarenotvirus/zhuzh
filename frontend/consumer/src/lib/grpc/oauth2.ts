/**
 * OAuth2 authorization code flow. Replicates backend/pkg/client/client.go WithOAuth2Credentials.
 * Uses JWT (from LoginForToken) as Bearer token to obtain OAuth2 access token for gRPC.
 */

import { randomBytes } from 'node:crypto';

export interface OAuth2TokenResult {
  accessToken: string;
  refreshToken?: string;
  expiresIn?: number;
}

/**
 * Exchanges JWT (from LoginForToken) for OAuth2 access token via authorization code flow.
 * 1. GET /oauth2/authorize with Bearer JWT -> redirect with ?code=...
 * 2. POST /oauth2/token with code -> OAuth2 access token
 */
export async function exchangeJwtForOAuth2Token(
  authServerUrl: string,
  clientId: string,
  clientSecret: string,
  jwt: string,
): Promise<OAuth2TokenResult> {
  const state = randomBytes(32).toString('base64url');
  const authUrl = new URL('/oauth2/authorize', authServerUrl);
  authUrl.searchParams.set('response_type', 'code');
  authUrl.searchParams.set('client_id', clientId);
  authUrl.searchParams.set('redirect_uri', authServerUrl);
  authUrl.searchParams.set('state', state);
  authUrl.searchParams.set('code_challenge_method', 'plain');

  const authRes = await fetch(authUrl.toString(), {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${jwt}`,
    },
    redirect: 'manual',
  });

  const location = authRes.headers.get('location');
  if (!location) {
    throw new Error('No redirect location from oauth2 authorize');
  }

  const redirectUrl = new URL(location, authServerUrl);
  const code = redirectUrl.searchParams.get('code');
  if (!code) {
    throw new Error('Code not returned from oauth2 redirect');
  }

  const tokenUrl = `${authServerUrl.replace(/\/$/, '')}/oauth2/token`;
  const tokenRes = await fetch(tokenUrl, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body: new URLSearchParams({
      grant_type: 'authorization_code',
      code,
      redirect_uri: authServerUrl,
      client_id: clientId,
      client_secret: clientSecret,
    }).toString(),
  });

  if (!tokenRes.ok) {
    const text = await tokenRes.text();
    throw new Error(`OAuth2 token exchange failed: ${tokenRes.status} ${text}`);
  }

  const tokenData = (await tokenRes.json()) as {
    access_token: string;
    refresh_token?: string;
    expires_in?: number;
  };

  return {
    accessToken: tokenData.access_token,
    refreshToken: tokenData.refresh_token,
    expiresIn: tokenData.expires_in,
  };
}
