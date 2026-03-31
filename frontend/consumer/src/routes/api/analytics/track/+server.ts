import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { trackEvent, trackAnonymousEvent } from '$lib/analytics';
import { randomUUID } from 'node:crypto';

const ANONYMOUS_ID_COOKIE = 'analytics_anonymous_id';

export const POST: RequestHandler = async ({ request, cookies }) => {
  const body = await request.json().catch(() => null);
  if (!body || typeof body.event !== 'string') {
    return json({ error: 'event is required' }, { status: 400 });
  }

  const event = body.event as string;
  const properties =
    typeof body.properties === 'object' && body.properties !== null
      ? Object.fromEntries(
          Object.entries(body.properties)
            .filter(([, v]) => typeof v === 'string' || typeof v === 'number' || typeof v === 'boolean')
            .map(([k, v]) => [k, String(v)]),
        )
      : {};

  const userId = body.userId;
  if (userId && typeof userId === 'string') {
    await trackEvent(event, { ...properties, userId });
  } else {
    let anonymousId = body.anonymousId ?? cookies.get(ANONYMOUS_ID_COOKIE);
    if (!anonymousId || typeof anonymousId !== 'string') {
      anonymousId = randomUUID();
      cookies.set(ANONYMOUS_ID_COOKIE, anonymousId, {
        path: '/',
        maxAge: 60 * 60 * 24 * 365,
        httpOnly: true,
        secure: true,
        sameSite: 'lax',
      });
    }
    await trackAnonymousEvent(event, anonymousId, properties);
  }

  return json({ ok: true });
};
