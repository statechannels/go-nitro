/// <reference types="vite/client" />
interface ImportMetaEnv {
  readonly VITE_NITRO_RPC_URL: string;
  readonly VITE_PROXY_URL: string;
  readonly VITE_PROVIDER: string;
  readonly VITE_HUB: string;
  readonly VITE_INITIAL_CHANNEL_BALANCE: string;
  readonly VITE_FILE_PATHS: string;
  readonly VITE_FILE_SIZES: string;
  readonly VITE_COST_PER_BYTE: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
