import prettier from 'eslint-config-prettier';
import js from '@eslint/js';
import svelte from 'eslint-plugin-svelte';
import globals from 'globals';
import ts from 'typescript-eslint';

/**
 * Shared ESLint flat config for SvelteKit apps. Pass the app's svelte.config so Svelte files parse correctly.
 * @param {import('@sveltejs/kit').Config} svelteConfig - The app's svelte.config.js default export
 * @returns {import('eslint').Linter.FlatConfig[]}
 */
export function createSvelteKitConfig(svelteConfig) {
  return [
    js.configs.recommended,
    ts.configs.recommended,
    svelte.configs.recommended,
    prettier,
    svelte.configs.prettier,
    {
      languageOptions: {
        globals: { ...globals.browser, ...globals.node },
      },
      rules: {
        'no-undef': 'off',
        'svelte/require-each-key': 'warn',
        'svelte/no-navigation-without-resolve': 'warn',
        'svelte/no-at-html-tags': 'warn',
        'svelte/prefer-svelte-reactivity': 'warn',
        '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }],
      },
    },
    {
      files: ['**/*.svelte', '**/*.svelte.ts', '**/*.svelte.js'],
      languageOptions: {
        parserOptions: {
          projectService: true,
          extraFileExtensions: ['.svelte'],
          parser: ts.parser,
          svelteConfig,
        },
      },
    },
  ];
}
