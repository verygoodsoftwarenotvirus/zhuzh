import type { PageServerLoad } from './$types';
import { getOAuth2Client } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ locals, params }) => {
  const token = locals.accessToken;
  const id = params.id;
  if (!token) {
    return { client: null, error: 'Not authenticated' };
  }
  try {
    const res = (await getOAuth2Client(token, { oauth2ClientId: id })) as {
      result?: Record<string, unknown>;
    };
    return { client: res?.result ?? null };
  } catch (e) {
    return {
      client: null,
      error: e instanceof Error ? e.message : 'Failed to load client',
    };
  }
};
