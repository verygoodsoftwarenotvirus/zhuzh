import type { PageServerLoad } from './$types';
import { getOAuth2Clients } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ locals }) => {
  const token = locals.accessToken;
  if (!token) {
    return { clients: [], error: 'Not authenticated' };
  }
  try {
    const res = (await getOAuth2Clients(token, { filter: { maxResponseSize: 100 } })) as {
      results?: Array<{ id?: string; clientId?: string; name?: string }>;
    };
    return { clients: res?.results ?? [] };
  } catch (e) {
    return {
      clients: [],
      error: e instanceof Error ? e.message : 'Failed to load OAuth2 clients',
    };
  }
};
