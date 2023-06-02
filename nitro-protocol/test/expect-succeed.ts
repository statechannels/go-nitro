export async function expectSupportedState(fn: () => void) {
  const txResult = (await fn()) as any;

  // `.stateIsSupported` returns a (bool, string) tuple
  expect(txResult.length).toBe(2);
}
