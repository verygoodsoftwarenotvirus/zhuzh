import type { PageServerLoad } from './$types';
import { verifyEmailAddress } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ url }) => {
  const token = url.searchParams.get('t') ?? '';

  if (!token) {
    return {
      success: false,
      message:
        'This verification link is invalid. Please check your email for the correct link or sign in to request a new one.',
    };
  }

  try {
    await verifyEmailAddress({ token });
    return {
      success: true,
      message: 'Your email has been verified successfully.',
    };
  } catch {
    return {
      success: false,
      message: 'This verification link is invalid or has expired. Please sign in to request a new verification email.',
    };
  }
};
