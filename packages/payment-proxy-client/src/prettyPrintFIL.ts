import bigDecimal from "js-big-decimal-esm";

const names = [
  "attoFIL",
  "femtoFIL",
  "picoFIL",
  "nanoFIL",
  "microFIL",
  "milliFIL",
  "FIL",
];
const decimals = [0n, 3n, 6n, 9n, 12n, 15n, 18n];

export const prettyPrintFIL = (wei: bigint | number | undefined): string => {
  if (typeof wei == "number") {
    wei = BigInt(wei);
  }
  const PRECISION = 2;

  if (wei === 0n) {
    return "0 FIL".padStart(3, "0");
  }

  if (wei == undefined) {
    return "-";
  }

  let formattedString = "";
  decimals.forEach((decimal, index) => {
    if (wei == undefined) {
      return "-";
    }
    if (wei >= 10n ** decimal) {
      formattedString = `${bigDecimal.divide(wei, 10n ** decimal, PRECISION)} ${
        names[index]
      }`;
    }
  });
  return formattedString;
};
