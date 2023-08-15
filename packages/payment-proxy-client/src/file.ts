import { NitroRpcClient } from "@statechannels/nitro-rpc-client";

export async function fetchFile(
  baseUrl: string,
  costPerByte: number,
  fileSize: number,
  selectedChannel: string,
  nitroClient: NitroRpcClient
): Promise<File> {
  const voucher = await nitroClient.CreateVoucher(
    selectedChannel,
    fileSize * costPerByte
  );

  const response = await fetch(
    `${baseUrl}?channelId=${voucher.ChannelId}&amount=${voucher.Amount}&signature=${voucher.Signature}`
  );

  try {
    const fileName = parseFileNameFromUrl(response.url);
    console.log(fileName);
    return new File([await response.blob()], fileName);
  } catch (e) {
    console.log(e);
    throw e;
  }
}

function parseFileNameFromUrl(url: string): string {
  try {
    const parsedUrl = new URL(url);
    const segments = parsedUrl.pathname.split("/");

    // Get the last segment of the pathname, which should be the file name
    return segments[segments.length - 1];
  } catch (error) {
    // If we can't parse the URL, just return a default file name
    return "fetched-file";
  }
}
