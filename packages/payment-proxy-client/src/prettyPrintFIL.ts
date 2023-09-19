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
  const PRECISION = 3;

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

export const prettyPrintPair = (
  wei1: bigint | number | undefined,
  wei2: bigint | number | undefined
): string => {
  if (typeof wei1 == "number") {
    wei1 = BigInt(wei1);
  }
  if (typeof wei2 == "number") {
    wei2 = BigInt(wei2);
  }
  const PRECISION = 3;

  if (wei1 == undefined || wei2 == undefined) {
    return "-";
  }

  const total = wei1 + wei2;

  if (total === 0n) {
    return "-";
  }

  let formattedString = "";
  decimals.forEach((decimal, index) => {
    if (total >= 10n ** decimal) {
      formattedString = `${bigDecimal
        .divide(wei1, 10n ** decimal, PRECISION)
        .padStart(4, "0")} ${names[index]} ${bigDecimal
        .divide(wei2, 10n ** decimal, PRECISION)
        .padStart(4, "0")}`;
    }
  });
  return formattedString;
};
