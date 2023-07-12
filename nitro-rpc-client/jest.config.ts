import type { Config } from "@jest/types";
// Sync object
const config: Config.InitialOptions = {
  testPathIgnorePatterns: ["<rootDir>/node_modules/", "<rootDir>/dist/"],
  preset: "ts-jest",
  testEnvironment: "node",
  verbose: true,
};
export default config;
