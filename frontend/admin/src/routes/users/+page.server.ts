import type { PageServerLoad } from './$types';
import { getUsers } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ locals }) => {
  const token = locals.accessToken;
  if (!token) {
    return { users: [], error: 'Not authenticated' };
  }
  try {
    const res = (await getUsers(token, { filter: { maxResponseSize: 100 } })) as {
      results?: Array<{ id?: string; username?: string; firstName?: string; lastName?: string }>;
    };
    return { users: res?.results ?? [] };
  } catch (e) {
    return {
      users: [],
      error: e instanceof Error ? e.message : 'Failed to load users',
    };
  }
};
