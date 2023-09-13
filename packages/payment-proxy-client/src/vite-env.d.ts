/// <reference types="vite/client" />
interface ImportMetaEnv {
  readonly VITE_NITRO_RPC_URL: string;
  readonly VITE_FILE_URL: string;
  readonly VITE_PROVIDER: string;
  readonly VITE_HUB: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
