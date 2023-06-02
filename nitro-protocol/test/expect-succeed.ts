/**
 * Wrapper for transactions that are expected to succeed with no return values.
 */
export async function expectSucceed(fn: () => void) {
  const txResult = (await fn()) as any;

  expect(txResult.length).toBe(0);
}

/**
 * Wrapper for calls to `stateIsSupported` that are expected to succeed.
 */
export async function expectSupportedState(fn: () => void) {
  const txResult = (await fn()) as any;

  // `.stateIsSupported` returns a (bool, string) tuple
  expect(txResult.length).toBe(2);
}
