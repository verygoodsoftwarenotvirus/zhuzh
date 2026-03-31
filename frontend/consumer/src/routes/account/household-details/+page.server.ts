import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { getActiveAccount, getSelf, updateAccount } from '$lib/grpc/clients';

const ACCOUNT_ADMIN_ROLE = 'account_admin';

export const load: PageServerLoad = async ({ locals, url }) => {
  const token = locals.oauthToken;
  if (!token) {
    return { account: null, isAdmin: false, error: null, updated: false };
  }

  try {
    const activeRes = await getActiveAccount(token);
    const account = activeRes.result ?? null;

    if (!account) {
      return {
        account: null,
        isAdmin: false,
        error: null,
        updated: false,
      };
    }

    const selfRes = await getSelf(token);
    const currentUserId = selfRes.result?.id ?? '';

    let isAdmin = false;
    for (const m of account.members ?? []) {
      const userId = m.belongsToUser?.id ?? '';
      if (userId === currentUserId && m.accountRole === ACCOUNT_ADMIN_ROLE) {
        isAdmin = true;
        break;
      }
    }

    if (!isAdmin) {
      throw redirect(302, '/account/settings');
    }

    const error = url.searchParams.get('error');
    const updated = url.searchParams.get('updated') === '1';

    return { account, isAdmin, error, updated };
  } catch (e) {
    if (e && typeof e === 'object' && 'status' in e && (e as { status: number }).status === 302) {
      throw e;
    }
    return {
      account: null,
      isAdmin: false,
      error: 'server',
      updated: false,
    };
  }
};

export const actions: Actions = {
  update: async ({ request, locals }) => {
    const token = locals.oauthToken;
    if (!token) throw redirect(302, '/login');

    const activeRes = await getActiveAccount(token);
    const account = activeRes.result;
    if (!account) {
      throw redirect(302, '/account/household-details?error=server');
    }

    const formData = await request.formData();
    const name = (formData.get('name') as string)?.trim() ?? '';
    const contactPhone = (formData.get('contact_phone') as string)?.trim() ?? '';
    const addressLine1 = (formData.get('address_line_1') as string)?.trim() ?? '';
    const addressLine2 = (formData.get('address_line_2') as string)?.trim() ?? '';
    const city = (formData.get('city') as string)?.trim() ?? '';
    const state = (formData.get('state') as string)?.trim() ?? '';
    const zipCode = (formData.get('zip_code') as string)?.trim() ?? '';
    const country = (formData.get('country') as string)?.trim() ?? '';

    if (!name) {
      throw redirect(302, '/account/household-details?error=invalid_name');
    }

    try {
      await updateAccount(token, {
        accountId: account.id,
        input: {
          name,
          contactPhone: contactPhone || undefined,
          addressLine1: addressLine1 || undefined,
          addressLine2: addressLine2 || undefined,
          city: city || undefined,
          state: state || undefined,
          zipCode: zipCode || undefined,
          country: country || undefined,
          belongsToUser: account.belongsToUser,
        },
      });
      throw redirect(302, '/account/household-details?updated=1');
    } catch (e) {
      if (e && typeof e === 'object' && 'status' in e && (e as { status: number }).status === 302) {
        throw e;
      }
      throw redirect(302, '/account/household-details?error=update_failed');
    }
  },
};
