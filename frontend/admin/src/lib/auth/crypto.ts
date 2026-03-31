/**
 * AES-256-GCM cookie encryption/decryption for admin session.
 * Uses a single base64-encoded 32-byte key (COOKIE_ENCRYPTION_KEY).
 */

import { env } from '$env/dynamic/private';
import { createCipheriv, createDecipheriv, randomBytes } from 'node:crypto';

const IV_LENGTH = 12;
const AUTH_TAG_LENGTH = 16;
const KEY_LENGTH = 32;

function getKey(): Buffer {
  const keyB64 = env.COOKIE_ENCRYPTION_KEY;
  if (!keyB64) {
    throw new Error('COOKIE_ENCRYPTION_KEY is required');
  }
  const key = Buffer.from(keyB64, 'base64');
  if (key.length !== KEY_LENGTH) {
    throw new Error(`COOKIE_ENCRYPTION_KEY must decode to ${KEY_LENGTH} bytes`);
  }
  return key;
}

/**
 * Encrypt and base64-encode a JSON-serializable value.
 * Format: base64(iv || ciphertext || authTag)
 */
export function encrypt(value: unknown): string {
  const key = getKey();
  const iv = randomBytes(IV_LENGTH);
  const cipher = createCipheriv('aes-256-gcm', key, iv, { authTagLength: AUTH_TAG_LENGTH });

  const plaintext = JSON.stringify(value);
  const encrypted = Buffer.concat([cipher.update(plaintext, 'utf8'), cipher.final(), cipher.getAuthTag()]);

  return Buffer.concat([iv, encrypted]).toString('base64');
}

/**
 * Decrypt a value produced by encrypt().
 */
export function decrypt<T = unknown>(encoded: string): T {
  const key = getKey();
  const buf = Buffer.from(encoded, 'base64');
  if (buf.length < IV_LENGTH + AUTH_TAG_LENGTH) {
    throw new Error('Invalid encrypted payload');
  }

  const iv = buf.subarray(0, IV_LENGTH);
  const ciphertext = buf.subarray(IV_LENGTH, buf.length - AUTH_TAG_LENGTH);
  const authTag = buf.subarray(buf.length - AUTH_TAG_LENGTH);

  const decipher = createDecipheriv('aes-256-gcm', key, iv, { authTagLength: AUTH_TAG_LENGTH });
  decipher.setAuthTag(authTag);

  const plaintext = decipher.update(ciphertext) + decipher.final('utf8');
  return JSON.parse(plaintext) as T;
}
