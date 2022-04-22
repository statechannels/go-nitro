import {resolve} from 'path';

import {configureEnvVariables} from '@statechannels/devtools';

configureEnvVariables();
// eslint-disable-next-line no-undef
const root = resolve(__dirname, '../../');

export default {
  globals: {
    'ts-jest': {
      tsconfig: './tsconfig.json',
    },
  },
  rootDir: root,
  collectCoverageFrom: ['**/*.{js,jsx,ts,tsx}'],
  reporters: ['default'],
  testMatch: ['<rootDir>/test/src/**/?(*.)test.ts?(x)'],
  testEnvironment: 'node',
  testURL: 'http://localhost',
  preset: 'ts-jest',
};
