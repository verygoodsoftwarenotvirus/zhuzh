import { redirect } from '@sveltejs/kit';
import type { Handle } from '@sveltejs/kit';
import { decodeSession, getCookieName } from '$lib/auth/session';

const LOGIN_PATH = '/login';

const PUBLIC_PATHS = [LOGIN_PATH, '/_ops_', '/auth/passkey/authentication'];

function isPublicPath(pathname: string): boolean {
  return PUBLIC_PATHS.some((p) => pathname === p || pathname.startsWith(`${p}/`));
}

export const handle: Handle = async ({ event, resolve }) => {
  if (isPublicPath(event.url.pathname)) {
    return resolve(event);
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

  if (!payload?.accessToken) {
    event.cookies.delete(cookieName, { path: '/' });
    throw redirect(302, LOGIN_PATH);
  }

  event.locals.accessToken = payload.accessToken;
  return resolve(event);
};
