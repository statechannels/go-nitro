/* global module */
module.exports = {
  env: {
    browser: true,
    es2020: true,
  },
  extends: ["plugin:react-hooks/recommended", "plugin:storybook/recommended"],
  parserOptions: {
    ecmaVersion: "latest",
    sourceType: "module",
  },
  plugins: ["react-refresh"],
  rules: {
    "react-refresh/only-export-components": "warn",
  },
};
