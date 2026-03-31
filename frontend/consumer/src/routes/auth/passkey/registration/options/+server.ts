import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { beginPasskeyRegistration } from '$lib/grpc/clients';

function toBase64(bytes: Uint8Array): string {
  return Buffer.from(bytes).toString('base64');
}

export const POST: RequestHandler = async ({ locals }) => {
  const token = locals.oauthToken;
  if (!token) {
    return json({ error: 'authentication required' }, { status: 401 });
  }

  try {
    const res = await beginPasskeyRegistration(token);
    const optionsBytes = res.publicKeyCredentialCreationOptions;
    const bytes = optionsBytes instanceof Uint8Array ? optionsBytes : new Uint8Array(optionsBytes as ArrayBuffer);
    const publicKeyCredentialCreationOptions = toBase64(bytes);
    return json({
      challenge: res.challenge,
      publicKeyCredentialCreationOptions,
    });
  } catch {
    return json({ error: 'failed to get passkey options' }, { status: 500 });
  }
};
