import bigDecimal from "js-big-decimal-esm";

export const prettyPrintWei = (wei: bigint): string => {
  const PRECISION = 1;
  const names = ["wei", "kwei", "Mwei", "Gwei", "szabo", "finney", "ether"];
  const decimals = [0n, 3n, 6n, 9n, 12n, 15n, 18n];

  if (wei === 0n) {
    return "0 wei";
  } else if (!wei) {
    return "unknown";
  }
  let formattedString = "";
  decimals.forEach((decimal, index) => {
    if (wei >= 10n ** decimal) {
      formattedString = `${bigDecimal.divide(wei, 10n ** decimal, PRECISION)} ${
        names[index]
      }`;
    }
  });
  return formattedString;
};
