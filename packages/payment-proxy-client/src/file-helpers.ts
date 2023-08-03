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

  const fileName = parseFileNameFromUrl(response.url);

  return new File([await response.blob()], fileName);
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

export async function fetchFileInChunks(
  chunkSize: number,
  baseUrl: string,
  costPerByte: number,
  selectedChannel: string,
  nitroClient: NitroRpcClient,
  updateProgress: (progress: number) => void
): Promise<File> {
  const firstChunk = await fetchFileChunk(
    0,
    chunkSize - 1,
    baseUrl,
    costPerByte,
    selectedChannel,
    nitroClient
  );

  console.log(
    `Fetched the first chunk of the file using a chunk size of ${chunkSize} bytes`
  );

  const { contentLength, fileName } = firstChunk;

  console.log(`The file ${fileName} is ${contentLength} bytes in size`);

  const fileContents = new Uint8Array(contentLength);
  fileContents.set(firstChunk.data);

  let remainingContentLength = contentLength - chunkSize;

  if (remainingContentLength <= 0) {
    console.log("We have fetched the entire file in 1 chunk");
    updateProgress(100);
    return new File([fileContents], fileName);
  }

  console.log(
    `We have ${remainingContentLength} bytes to fetch in ${Math.ceil(
      remainingContentLength / chunkSize
    )} chunks`
  );
  updateProgress(100 - (remainingContentLength / contentLength) * 100);

  while (remainingContentLength > chunkSize) {
    const start = contentLength - remainingContentLength;
    const stop = start + chunkSize - 1;

    const { data } = await fetchFileChunk(
      start,
      stop,
      baseUrl,
      costPerByte,
      selectedChannel,
      nitroClient
    );

    fileContents.set(data, start);
    remainingContentLength -= chunkSize;

    updateProgress(100 - (remainingContentLength / contentLength) * 100);

    console.log(
      `We have ${remainingContentLength} bytes to fetch in ${Math.ceil(
        remainingContentLength / chunkSize
      )} chunks`
    );
  }

  if (remainingContentLength > 0) {
    const start = contentLength - remainingContentLength;
    const stop = contentLength - 1;
    const { data } = await fetchFileChunk(
      start,
      stop,
      baseUrl,
      costPerByte,
      selectedChannel,
      nitroClient
    );
    fileContents.set(data, start);

    console.log(`Fetched final chunk of size ${remainingContentLength} bytes`);
  }
  updateProgress(100);
  console.log("Finished fetching all chunks");
  return new File([fileContents], fileName);
}

export async function fetchFileChunk(
  start: number,
  stop: number,
  baseUrl: string,
  costPerByte: number,
  selectedChannel: string,
  nitroClient: NitroRpcClient
): Promise<{ data: Uint8Array; contentLength: number; fileName: string }> {
  const dataLength = stop - start + 1; // +1 because stop is inclusive

  const chunkCost = dataLength * costPerByte;

  const voucher = await nitroClient.CreateVoucher(selectedChannel, chunkCost);

  const req = new Request(
    `${baseUrl}?channelId=${voucher.ChannelId}&amount=${voucher.Amount}&signature=${voucher.Signature}`
  );
  req.headers.set("Range", `bytes=${start}-${stop}`);

  const response = await fetch(req);
  if (!response.body) {
    throw new Error("Response body is null");
  }

  if (!response.ok) {
    throw new Error(`Response status ${response.status}`);
  }
  const result = await response.body.getReader().read();

  return {
    data: result.value || new Uint8Array(),
    contentLength: parseTotalSizeFromContentRange(
      response.headers.get("Content-Range") || ""
    ),
    fileName: parseFileNameFromUrl(response.url),
  };
}

function parseTotalSizeFromContentRange(contentRange: string): number {
  const match = /^.*\/([0-9]*)$/.exec(contentRange);
  if (!match) {
    throw new Error(`Could not parse content range ${contentRange}`);
  }
  return parseInt(match[1]);
}
