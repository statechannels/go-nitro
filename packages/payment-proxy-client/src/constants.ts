// literals
export const QUERY_KEY = "rpcUrl";
export const costPerByte = 1;

export const dataSize = 2028;

// env vars
export const proxyUrl = import.meta.env.VITE_PROXY_URL;
export const fileRelativePath = import.meta.env.VITE_FILE_RELATIVE_PATH;

export const provider = import.meta.env.VITE_PROVIDER;
export const hub = import.meta.env.VITE_HUB;
export const defaultNitroRPCUrl = import.meta.env.VITE_NITRO_RPC_URL;
export const initialChannelBalance = parseInt(
  import.meta.env.VITE_INITIAL_CHANNEL_BALANCE,
  10
);

export const files: { fileName: string; url: string }[] =
  import.meta.env.VITE_FILE_PATHS.split(",").map((filePath: string) => ({
    url: proxyUrl + filePath,
    fileName: filePath.split("/").pop() || filePath,
  }));
