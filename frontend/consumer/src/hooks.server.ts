import { env } from '$env/dynamic/private';
import { redirect } from '@sveltejs/kit';
import type { Handle } from '@sveltejs/kit';
import { decodeSession, encodeSession, getCookieName, getCookieOptions } from '$lib/auth/session';
import { exchangeJwtForOAuth2Token } from '$lib/grpc/oauth2';
import { exchangeToken } from '$lib/grpc/clients';
import { initServerOtel } from '$lib/otel/server';
import { recordRequest } from '$lib/otel/server-metrics';
import { ServerTiming, ServerTimingHeaderName } from '$lib/server-timing';

initServerOtel();

const LOGIN_PATH = '/login';

const PUBLIC_PATHS = [
  LOGIN_PATH,
  '/logout',
  '/forgot_password',
  '/reset_password',
  '/verify_email_address',
  '/terms-of-service',
  '/privacy-policy',
  '/accept_invitation',
  '/_ops_',
  '/.well-known',
  '/auth/passkey/authentication',
];

function isPublicPath(pathname: string): boolean {
  return PUBLIC_PATHS.some((p) => pathname === p || pathname.startsWith(`${p}/`));
}

export const handle: Handle = async ({ event, resolve }) => {
  const timing = new ServerTiming();
  const totalEvent = timing.addEvent('total', 'Total request time');

  if (isPublicPath(event.url.pathname)) {
    const response = await resolve(event);
    totalEvent.end();
    response.headers.set(ServerTimingHeaderName, timing.headerValue());
    recordRequest(event.url.pathname, response.status, totalEvent.duration);
    return response;
  }

  const cookieName = getCookieName();
  const cookieValue = event.cookies.get(cookieName);

  if (!cookieValue) {
    throw redirect(302, LOGIN_PATH);
  }

  let payload;
  try {
    payload = decodeSession(cookieValue);
  } catch {
    event.cookies.delete(cookieName, { path: '/' });
    throw redirect(302, LOGIN_PATH);
  }

  const httpApiUrl = env.HTTP_API_SERVER_URL;
  const clientId = env.OAUTH2_CLIENT_ID;
  const clientSecret = env.OAUTH2_CLIENT_SECRET;

  if (!httpApiUrl || !clientId || !clientSecret) {
    throw new Error('HTTP_API_SERVER_URL, OAUTH2_CLIENT_ID, OAUTH2_CLIENT_SECRET are required');
  }

  const authEvent = timing.addEvent('auth', 'OAuth2 token exchange');
  try {
    const { accessToken } = await exchangeJwtForOAuth2Token(httpApiUrl, clientId, clientSecret, payload.accessToken);
    event.locals.oauthToken = accessToken;
  } catch {
    // JWT may be expired — try refreshing with the stored refresh token
    if (payload.refreshToken) {
      try {
        const refreshed = await exchangeToken(payload.refreshToken);
        const newPayload = {
          accessToken: refreshed.accessToken,
          refreshToken: refreshed.refreshToken,
        };
        const opts = getCookieOptions();
        const encoded = encodeSession(newPayload);
        event.cookies.set(cookieName, encoded, {
          path: opts.path,
          httpOnly: opts.httpOnly,
          secure: opts.secure,
          sameSite: opts.sameSite,
          maxAge: opts.maxAge,
        });

        const { accessToken } = await exchangeJwtForOAuth2Token(
          httpApiUrl,
          clientId,
          clientSecret,
          refreshed.accessToken,
        );
        event.locals.oauthToken = accessToken;
      } catch {
        event.cookies.delete(cookieName, { path: '/' });
        throw redirect(302, LOGIN_PATH);
      }
    } else {
      event.cookies.delete(cookieName, { path: '/' });
      throw redirect(302, LOGIN_PATH);
    }
  } finally {
    authEvent.end();
  }

  const response = await resolve(event);
  totalEvent.end();
  response.headers.set(ServerTimingHeaderName, timing.headerValue());
  recordRequest(event.url.pathname, response.status, totalEvent.duration);
  return response;
};
