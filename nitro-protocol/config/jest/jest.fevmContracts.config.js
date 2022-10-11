// eslint-disable-next-line no-undef
var config = require('./jest.config');

config.testMatch = ['<rootDir>/test/contracts/**/fevm/*.test.ts'];
config.reporters = ['default'];
config.testTimeout = 90_000;

// eslint-disable-next-line no-undef
module.exports = config;
