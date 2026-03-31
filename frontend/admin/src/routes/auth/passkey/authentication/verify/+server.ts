import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { finishPasskeyAuthentication } from '$lib/grpc/clients';
import { encodeSession, getCookieOptions } from '$lib/auth/session';

export const POST: RequestHandler = async ({ request, cookies }) => {
  let body: { challenge?: string; username?: string; assertionResponse?: unknown };
  try {
    body = await request.json();
  } catch {
    return json({ error: 'invalid request' }, { status: 400 });
  }

  const challenge = (body.challenge ?? '').trim();
  const username = (body.username ?? '').trim();
  const assertionResponse = body.assertionResponse;

  if (!challenge || !assertionResponse) {
    return json({ error: 'assertion_response and challenge are required' }, { status: 400 });
  }

  const assertionBytes =
    typeof assertionResponse === 'string'
      ? new TextEncoder().encode(assertionResponse)
      : new TextEncoder().encode(JSON.stringify(assertionResponse));

  try {
    const tokenRes = (await finishPasskeyAuthentication({
      challenge,
      username,
      assertionResponse: assertionBytes,
    })) as { result?: { accessToken?: string } };
    const accessToken = tokenRes.result?.accessToken;
    if (!accessToken) {
      return json({ error: 'no access token' }, { status: 500 });
    }

    const encoded = encodeSession({ accessToken });
    const opts = getCookieOptions();
    cookies.set(opts.name, encoded, {
      path: opts.path,
      httpOnly: opts.httpOnly,
      secure: opts.secure,
      sameSite: opts.sameSite,
      maxAge: opts.maxAge,
    });

    return json({ success: true, redirect: '/' });
  } catch {
    return json({ error: 'authentication failed' }, { status: 401 });
  }
};
