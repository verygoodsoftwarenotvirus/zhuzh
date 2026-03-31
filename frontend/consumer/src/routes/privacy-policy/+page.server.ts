import type { PageServerLoad } from './$types';
import { marked } from 'marked';

const privacyRaw = await import('$lib/content/privacy.md?raw').then((m) => m.default);

export const load: PageServerLoad = async () => {
  const html = marked.parse(privacyRaw, { async: false }) as string;
  return { html };
};
