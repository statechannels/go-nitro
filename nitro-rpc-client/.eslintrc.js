module.exports = {
  extends: ['../../.eslintrc.js', 'plugin:storybook/recommended'],
  overrides: [
    {
      files: ['**/*.{ts,tsx}'],
      rules: {
        'jsdoc/require-jsdoc': 0,
        'jsdoc/match-description': 0,
        '@typescript-eslint/prefer-for-of': 0,
      },
    },
  ],
  ignorePatterns: ['!.eslintrc.js', 'build/'],
};
