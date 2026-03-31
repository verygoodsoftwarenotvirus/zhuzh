import type { PageServerLoad } from './$types';
import { getUsers, getAccounts } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ locals }) => {
  const token = locals.accessToken;
  if (!token) {
    return { userCount: '-', accountCount: '-', error: 'Not authenticated' };
  }

  let userCount = '-';
  let accountCount = '-';

  try {
    const usersRes = (await getUsers(token, {})) as { results?: unknown[] };
    if (usersRes?.results) userCount = String(usersRes.results.length);
  } catch {
    // leave as '-'
  }
  try {
    const accountsRes = (await getAccounts(token, {})) as { results?: unknown[] };
    if (accountsRes?.results) accountCount = String(accountsRes.results.length);
  } catch {
    // leave as '-'
  }
  return { userCount, accountCount };
};
