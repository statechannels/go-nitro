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

// computed
export const fileUrl = proxyUrl + fileRelativePath;
