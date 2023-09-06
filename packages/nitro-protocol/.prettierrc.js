module.exports = {
  trailingComma: "es5",
  arrowParens: "avoid",
  bracketSpacing: false,
  printWidth: 100,
  singleQuote: true,
  tabWidth: 2,
  overrides: [
    {
      files: '*.sol',
      options: {
        tabWidth: 4
      }
    }
  ]
};
