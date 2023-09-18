import { NitroRpcClient } from "@statechannels/nitro-rpc-client";
import { Voucher } from "@statechannels/nitro-rpc-client/src/types";

export async function fetchFile(
  url: string,
  paymentAmount: number,
  channelId: string,
  nitroClient: NitroRpcClient
): Promise<File> {
  console.time("Create Payment Vouncher");
  const voucher = await nitroClient.CreateVoucher(channelId, paymentAmount);
  console.timeEnd("Create Payment Vouncher");
  console.time("Fetch file");
  const response = await fetch(addVoucherToUrl(url, voucher));
  if (response.status == 402) {
    throw new Error(`402 ${await response.text()}`);
  }
  console.timeEnd("Fetch file");

  const fileName = parseFileNameFromUrl(response.url);

  return new File([await response.blob()], fileName);
}
export async function fetchFileInChunks(
  chunkSize: number,
  url: string,
  costPerByte: number,
  channelId: string,
  nitroClient: NitroRpcClient,
  updateProgress: (progress: number) => void
): Promise<File> {
  updateProgress(0);

  const firstChunk = await fetchChunk(
    0,
    chunkSize - 1,
    url,
    costPerByte,
    channelId,
    nitroClient
  );

  const { contentLength, fileName } = firstChunk;
  let remainingContentLength = contentLength - chunkSize;

  console.log(
    `Fetched the first chunk of the file using a chunk size of ${chunkSize} bytes`
  );

  console.log(`The file ${fileName} is ${contentLength} bytes in size`);

  const fileContents = new Uint8Array(contentLength);
  fileContents.set(firstChunk.data);

  if (remainingContentLength <= 0) {
    console.log("We have fetched the entire file in 1 chunk");
    updateProgress(100);
    return new File([fileContents], fileName);
  }

  while (remainingContentLength > chunkSize) {
    updateProgress(100 - (remainingContentLength / contentLength) * 100);
    console.log(
      `We have ${remainingContentLength} bytes to fetch in ${Math.ceil(
        remainingContentLength / chunkSize
      )} chunks`
    );

    const start = contentLength - remainingContentLength;
    const stop = start + chunkSize - 1;

    const { data } = await fetchChunk(
      start,
      stop,
      url,
      costPerByte,
      channelId,
      nitroClient
    );

    fileContents.set(data, start);
    remainingContentLength -= chunkSize;
  }

  if (remainingContentLength > 0) {
    const start = contentLength - remainingContentLength;
    const stop = contentLength - 1;
    const { data } = await fetchChunk(
      start,
      stop,
      url,
      costPerByte,
      channelId,
      nitroClient
    );
    fileContents.set(data, start);

    console.log(`Fetched final chunk of size ${remainingContentLength} bytes`);
  }
  updateProgress(100);
  console.log("Finished fetching all chunks");
  return new File([fileContents], fileName);
}

async function fetchChunk(
  start: number,
  stop: number,
  url: string,
  costPerByte: number,
  channelId: string,
  nitroClient: NitroRpcClient
): Promise<{ data: Uint8Array; contentLength: number; fileName: string }> {
  const dataLength = stop - start + 1; // +1 because stop is inclusive
  const chunkCost = dataLength * costPerByte;

  const voucher = await nitroClient.CreateVoucher(channelId, chunkCost);

  const req = new Request(addVoucherToUrl(url, voucher));
  req.headers.set("Range", `bytes=${start}-${stop}`);

  const response = await fetch(req);
  if (response.status == 402) {
    throw new Error(`402 ${await response.text()}`);
  }
  return {
    data: await getChunkData(response),
    contentLength: parseTotalSizeFromContentRange(
      response.headers.get("Content-Range")
    ),
    fileName: parseFileNameFromUrl(response.url),
  };
}

function addVoucherToUrl(url: string, voucher: Voucher): string {
  return `${url}?channelId=${voucher.ChannelId}&amount=${voucher.Amount}&signature=${voucher.Signature}`;
}

async function getChunkData(res: Response): Promise<Uint8Array> {
  if (!res.body) {
    throw new Error("Response body is null");
  }

  if (!res.ok) {
    throw new Error(`Response status ${res.status}`);
  }
  const result = await res.body.getReader().read();
  return result.value || new Uint8Array();
}

function parseTotalSizeFromContentRange(
  contentRange: string | null | undefined
): number {
  if (!contentRange) {
    throw new Error("Content range is null or undefined");
  }
  const match = /^.*\/([0-9]*)$/.exec(contentRange);
  if (!match) {
    throw new Error(`Could not parse content range ${contentRange}`);
  }
  return parseInt(match[1]);
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
