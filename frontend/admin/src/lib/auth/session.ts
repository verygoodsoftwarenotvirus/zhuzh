/**
 * Admin session types and cookie helpers.
 */

import { env } from '$env/dynamic/private';
import { encrypt, decrypt } from './crypto';

export interface AuthPayload {
  accessToken: string;
}

const COOKIE_LIFETIME_SECONDS = 180 * 24 * 60 * 60; // 180 days

export function getCookieName(): string {
  return env.COOKIE_NAME ?? 'admin_webapp';
}

export function encodeSession(payload: AuthPayload): string {
  return encrypt(payload);
}

export function decodeSession(encoded: string): AuthPayload {
  return decrypt<AuthPayload>(encoded);
}

export function getCookieOptions(): {
  name: string;
  path: string;
  httpOnly: boolean;
  secure: boolean;
  sameSite: 'lax';
  maxAge: number;
} {
  return {
    name: getCookieName(),
    path: '/',
    httpOnly: true,
    secure: env.NODE_ENV === 'production',
    sameSite: 'lax',
    maxAge: COOKIE_LIFETIME_SECONDS,
  };
}
