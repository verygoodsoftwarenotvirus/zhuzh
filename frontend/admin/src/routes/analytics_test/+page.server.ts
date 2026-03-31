import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';
import { trackEvent } from '$lib/grpc/clients';

export const actions: Actions = {
  default: async ({ request, locals }) => {
    const token = locals.accessToken;
    if (!token) {
      return fail(401, { error: 'Not authenticated' });
    }
    const formData = await request.formData();
    const eventName = (formData.get('event_name') as string)?.trim() ?? 'admin_analytics_test';
    try {
      await trackEvent(token, {
        source: 'admin',
        event: eventName,
        properties: {},
      });
      return { success: true, message: `Event "${eventName}" sent` };
    } catch (e) {
      return fail(500, {
        error: e instanceof Error ? e.message : 'Analytics test failed',
      });
    }
  },
};
