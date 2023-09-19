// literals
export const QUERY_KEY = "rpcUrl";
export const costPerByte = 1;

// TODO: This is temporary until https://github.com/statechannels/go-nitro/pull/1754
// We just use the largest deployed file size for now
export const dataSize = 2221266;
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

const ENV_VAR_SPLIT_CHAR = ";";

export const files: { fileName: string; url: string }[] =
  import.meta.env.VITE_FILE_PATHS.split(ENV_VAR_SPLIT_CHAR).map(
    (filePath: string) => ({
      url: proxyUrl + filePath,
      fileName: filePath.split("/").pop() || filePath,
    })
  );
