// eslint-disable-next-line no-undef
var config = require('./jest.config');

config.testMatch = ['<rootDir>/test/contracts/**/*.test.ts'];
config.reporters = ['default'];
config.globalSetup = '<rootDir>/jest/contract-test-setup.ts';
config.globalTeardown = '<rootDir>/jest/contract-test-teardown.ts';
config.testTimeout = 90_000;

// eslint-disable-next-line no-undef
module.exports = config;
