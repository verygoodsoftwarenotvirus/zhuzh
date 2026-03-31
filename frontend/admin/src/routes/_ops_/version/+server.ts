import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const GET: RequestHandler = async () => {
  return json({
    version: env.VERSION || 'unknown',
    commit_hash: env.COMMIT_HASH || 'unknown',
    commit_time: env.COMMIT_TIME || 'unknown',
    build_time: env.BUILD_TIME || 'unknown',
  });
};
