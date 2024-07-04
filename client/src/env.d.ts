/// <reference types="astro/client" />

interface ImportMetaEnv {
  readonly PUBLIC_SERVER_PORT: number;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
