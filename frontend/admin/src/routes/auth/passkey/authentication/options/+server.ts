import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { beginPasskeyAuthentication } from '$lib/grpc/clients';

export const POST: RequestHandler = async ({ request }) => {
  let body: { username?: string };
  try {
    body = await request.json();
  } catch {
    return json({ error: 'invalid request' }, { status: 400 });
  }

  const username = (body.username ?? '').trim();

  try {
    const res = (await beginPasskeyAuthentication({ username })) as {
      challenge: string;
      publicKeyCredentialRequestOptions: Uint8Array;
    };
    const optionsBytes = res.publicKeyCredentialRequestOptions;
    const bytes = optionsBytes instanceof Uint8Array ? optionsBytes : new Uint8Array(optionsBytes as ArrayBuffer);
    const publicKeyCredentialRequestOptions = Buffer.from(bytes).toString('base64');
    return json({
      challenge: res.challenge,
      publicKeyCredentialRequestOptions,
    });
  } catch {
    return json({ error: 'failed to get passkey options' }, { status: 500 });
  }
};
