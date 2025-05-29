import eslintConfigPrettier from 'eslint-config-prettier';
import globals from 'globals';
import eslint from '@eslint/js';
import { globalIgnores } from 'eslint/config';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  eslint.configs.recommended,
  tseslint.configs.strict,
  tseslint.configs.recommendedTypeCheckedOnly,
  tseslint.configs.stylisticTypeChecked,
  globalIgnores(['**/dist/', '*.config.mjs']),
  {
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      globals: {
        ...globals.node, // todo - configure globals to split node vs browser
      },
      ecmaVersion: 'latest',
      sourceType: 'module',
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
  },
  {
    rules: {
      complexity: 'error',
      '@typescript-eslint/return-await': ['error', 'always'],
    },
  },
  eslintConfigPrettier,
);
