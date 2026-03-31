import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { decodeSession, getCookieName } from '$lib/auth/session';
import { revokeCurrentSession } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ cookies }) => {
  const cookieName = getCookieName();
  const cookieValue = cookies.get(cookieName);

  if (cookieValue) {
    try {
      const { accessToken } = decodeSession(cookieValue);
      await revokeCurrentSession(accessToken);
    } catch {
      /* best-effort: proceed with logout even if revocation fails */
    }
  }

  cookies.delete(cookieName, { path: '/' });
  throw redirect(302, '/login');
};
