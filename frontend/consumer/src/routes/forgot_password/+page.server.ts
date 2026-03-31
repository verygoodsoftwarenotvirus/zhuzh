import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { requestPasswordResetToken } from '$lib/grpc/clients';

export const load: PageServerLoad = async () => {
  return {};
};

export const actions: Actions = {
  default: async ({ request }) => {
    const formData = await request.formData();
    const email = (formData.get('email') as string)?.trim() ?? '';

    if (!email) {
      return fail(400, {
        error: 'Email is required',
        success: false,
      });
    }

    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!emailRegex.test(email)) {
      return fail(400, {
        error: 'Please enter a valid email address',
        success: false,
      });
    }

    try {
      await requestPasswordResetToken({ emailAddress: email });
      return { success: true, error: null };
    } catch {
      // Always show success to avoid email enumeration
      return { success: true, error: null };
    }
  },
};
