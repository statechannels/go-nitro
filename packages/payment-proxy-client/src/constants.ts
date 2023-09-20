// literals
export const QUERY_KEY = "rpcUrl";
export const costPerByte = 1;

// env vars
export const CHUNK_SIZE = parseInt(import.meta.env.VITE_CHUNK_SIZE, 10);
export const proxyUrl = import.meta.env.VITE_PROXY_URL;
export const fileRelativePath = import.meta.env.VITE_FILE_RELATIVE_PATH;

export const provider = import.meta.env.VITE_PROVIDER;
export const hub = import.meta.env.VITE_HUB;
export const defaultNitroRPCUrl = import.meta.env.VITE_NITRO_RPC_URL;
export const initialChannelBalance = parseInt(
  import.meta.env.VITE_INITIAL_CHANNEL_BALANCE,
  10
);
const ENV_VAR_SPLIT_CHAR = ";";

const fileSizes = import.meta.env.VITE_FILE_SIZES.split(ENV_VAR_SPLIT_CHAR).map(
  (size: string) => parseInt(size, 10)
);
const fileNames = import.meta.env.VITE_FILE_NAMES.split(ENV_VAR_SPLIT_CHAR);

export interface AvailableFile {
  fileName: string;
  url: string;
  size: number;
}
export const files: AvailableFile[] = import.meta.env.VITE_FILE_PATHS.split(
  ENV_VAR_SPLIT_CHAR
).map((filePath: string, index: number) => {
  return {
    url: proxyUrl + filePath,
    fileName: fileNames[index],
    size: fileSizes[index],
  };
});
