export async function expectSucceed(fn: () => void) {
  const txResult = (await fn()) as any;

  // As 'requireStateSupported' method is constant (view or pure), if it succeeds, it returns an object/array with returned values
  // which in this case should be empty
  expect(txResult.length).toBe(0);
}
