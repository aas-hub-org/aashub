import { FlatCompat } from '@eslint/eslintrc';
import ESLint from '@eslint/js';
import Oxlint from 'eslint-plugin-oxlint';
import Vue from 'eslint-plugin-vue';
import globals from 'globals';
import { dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const compat = new FlatCompat({ baseDirectory: __dirname });

export default [
  {
    ignores: ['{dist,public}/**/*'],
  },
  ESLint.configs.recommended,
  ...Vue.configs['flat/recommended'],
  ...compat.extends('@vue/eslint-config-typescript/recommended'),
  Oxlint.configs['flat/recommended'],
  ...compat.extends('@vue/eslint-config-prettier/skip-formatting'),
  {
    files: ['**/*.{js,mjs,cjs,jsx,vue,ts,mts,cts,tsx}'],
    linterOptions: {
      reportUnusedDisableDirectives: true,
    },
    languageOptions: {
      globals: {
        ...globals.node,
        ...globals.browser,
        ...globals.es2021,
      },
    },
    plugins: {},
    rules: {
      '@typescript-eslint/no-explicit-any': ['warn', { ignoreRestArgs: true }],
      '@typescript-eslint/ban-types': [
        'warn',
        {
          types: {
            '{}': {
              message: 'Consider using a more specific type instead of `{}`.',
            },
          },
          extendDefaults: true,
        },
      ],
    },
  },
];
