/* eslint-disable no-undef */
var config = require('./jest.config');

config.testMatch = ['<rootDir>/gas-benchmarks/**/*.test.ts'];
config.reporters = ['default'];
config.setupFilesAfterEnv = ['<rootDir>/gas-benchmarks/jestSetup.ts'];
config.maxConcurrency = 1;

/* eslint-disable no-undef */
module.exports = config;
