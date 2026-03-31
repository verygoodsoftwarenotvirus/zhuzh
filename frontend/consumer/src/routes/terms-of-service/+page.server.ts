import type { PageServerLoad } from './$types';
import { marked } from 'marked';

const termsRaw = await import('$lib/content/terms.md?raw').then((m) => m.default);

export const load: PageServerLoad = async () => {
  const html = marked.parse(termsRaw, { async: false }) as string;
  return { html };
};
