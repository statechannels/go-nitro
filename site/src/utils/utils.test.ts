import { prettyPrintWei } from ".";

describe("Pretty printing wei", () => {
  it.each`
    input          | output
    ${18n}         | ${"18.0 wei"}
    ${1001n}       | ${"1.0 kwei"}
    ${1599n}       | ${"1.6 kwei"}
    ${1234567n}    | ${"1.2 Mwei"}
    ${12345678n}   | ${"12.3 Mwei"}
    ${123456789n}  | ${"123.5 Mwei"}
    ${1234567890n} | ${"1.2 Gwei"}
    ${0n}          | ${"0 wei"}
    ${10n ** 18n}  | ${"1.0 ether"}
  `("prettyPrintWei(bigNumberify($input)) = $output", ({ input, output }) => {
    expect(prettyPrintWei(input)).toStrictEqual(output);
  });
});
