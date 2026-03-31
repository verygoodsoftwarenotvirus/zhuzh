import path from 'node:path';
import { includeIgnoreFile } from '@eslint/compat';
import { defineConfig } from 'eslint/config';
import { createSvelteKitConfig } from '@zhuzh/eslint-config-custom';
import svelteConfig from './svelte.config.js';

const gitignorePath = path.resolve(import.meta.dirname, '.gitignore');

export default defineConfig(
  { ignores: ['.svelte-kit/**', 'build/**'] },
  includeIgnoreFile(gitignorePath),
  ...createSvelteKitConfig(svelteConfig),
);
