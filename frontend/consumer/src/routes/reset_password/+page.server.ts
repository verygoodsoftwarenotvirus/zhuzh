import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { redeemPasswordResetToken } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ url }) => {
  const token = url.searchParams.get('t') ?? '';
  return { token, missingToken: !token };
};

export const actions: Actions = {
  default: async ({ request }) => {
    const formData = await request.formData();
    const token = (formData.get('token') as string)?.trim() ?? '';
    const newPassword = (formData.get('new_password') as string) ?? '';
    const confirmPassword = (formData.get('confirm_password') as string) ?? '';

    if (!token) {
      return fail(400, {
        error: 'Missing reset token. Please use the link from your email.',
        token: '',
      });
    }

    if (!newPassword || newPassword.length < 8) {
      return fail(400, {
        error: 'Password must be at least 8 characters',
        token,
      });
    }

    if (newPassword !== confirmPassword) {
      return fail(400, {
        error: 'Passwords do not match',
        token,
      });
    }

    try {
      await redeemPasswordResetToken({ token, newPassword });
      throw redirect(302, '/login?reset=success');
    } catch (e) {
      if (e && typeof e === 'object' && 'status' in e && (e as { status: number }).status === 302) {
        throw e;
      }
      return fail(400, {
        error: 'Invalid or expired reset link. Please request a new one.',
        token,
      });
    }
  },
};
