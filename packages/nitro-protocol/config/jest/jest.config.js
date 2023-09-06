// eslint-disable-next-line no-undef
const {resolve} = require('path');

/* eslint-disable no-undef */
const {configureEnvVariables} = require('@statechannels/devtools');

configureEnvVariables();
// eslint-disable-next-line no-undef
const root = resolve(__dirname, '../../');

// eslint-disable-next-line no-undef
module.exports = {
  globals: {
    'ts-jest': {
      tsconfig: './test/tsconfig.json',
    },
  },
  rootDir: root,
  collectCoverageFrom: ['**/*.{js,jsx,ts,tsx}'],
  reporters: ['default'],
  testMatch: ['<rootDir>/test/src/**/?(*.)test.ts?(x)', '<rootDir>/test/?(*.)test.ts?(x)'],
  testEnvironment: 'node',
  testURL: 'http://localhost',
  preset: 'ts-jest',
};
