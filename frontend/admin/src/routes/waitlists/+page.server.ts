import type { PageServerLoad } from './$types';
import { getWaitlists } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ locals }) => {
  const token = locals.accessToken;
  if (!token) {
    return { waitlists: [], error: 'Not authenticated' };
  }
  try {
    const res = (await getWaitlists(token, { filter: { maxResponseSize: 100 } })) as {
      results?: Array<{ id?: string; name?: string }>;
    };
    return { waitlists: res?.results ?? [] };
  } catch (e) {
    return {
      waitlists: [],
      error: e instanceof Error ? e.message : 'Failed to load waitlists',
    };
  }
};
