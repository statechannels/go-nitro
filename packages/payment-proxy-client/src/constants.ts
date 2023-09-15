// literals
export const QUERY_KEY = "rpcUrl";
export const CHANNEL_ID_KEY = "channelId";

// env vars
export const provider = import.meta.env.VITE_PROVIDER;
export const hub = import.meta.env.VITE_HUB;
export const defaultNitroRPCUrl = import.meta.env.VITE_NITRO_RPC_URL;
export const defaultFileUrl = import.meta.env.VITE_FILE_URL;
export const initialChannelBalance = parseInt(
  import.meta.env.VITE_INITIAL_CHANNEL_BALANCE,
  10
);
