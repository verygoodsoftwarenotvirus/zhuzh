import type { PageServerLoad } from './$types';
import { getSelf, getAccountsForUser } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ locals }) => {
  const token = locals.oauthToken;
  if (!token) {
    return { hasAccount: false };
  }

  try {
    const selfRes = await getSelf(token);
    const user = selfRes.result;
    const userId = user?.id ?? '';

    if (!userId) {
      return { hasAccount: false };
    }

    const accountsRes = await getAccountsForUser(token, {
      userId,
      filter: { maxResponseSize: 1 },
    });
    const hasAccount = (accountsRes.results?.length ?? 0) > 0;
    return { hasAccount };
  } catch {
    return { hasAccount: false };
  }
};
