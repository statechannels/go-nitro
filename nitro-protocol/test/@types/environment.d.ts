declare global {
  namespace NodeJS {
    interface ProcessEnv {
      DEFAULT_GAS: string;
      DEFAULT_GAS_PRICE: string;
      GANACHE_HOST: string;
      GANACHE_PORT: string;
      CHAIN_NETWORK_ID: string;
      DEV_HTTP_SERVER_PORT: string;

      // These contract addresses get defined in the global jest setup
      NITRO_ADJUDICATOR_ADDRESS: string;
      COUNTING_APP_ADDRESS: string;
      CONSENSUS_APP_ADDRESS: string;
      VIRTUAL_PAYMENT_APP_ADDRESS: string;
      HASH_LOCK_ADDRESS: string;
      EMBEDDED_APPLICATION_ADDRESS: string;
      SINGLE_ASSET_PAYMENTS_ADDRESS: string;
      TRIVIAL_APP_ADDRESS: string;
      TEST_NITRO_UTILS_ADDRESS: string;
      TEST_STRICT_TURN_TAKING_ADDRESS: string;
      TEST_CONSENSUS_ADDRESS: string;
      TEST_FORCE_MOVE_ADDRESS: string;
      TEST_NITRO_ADJUDICATOR_ADDRESS: string;
      TEST_TOKEN_ADDRESS: string;
      BAD_TOKEN_ADDRESS: string;
      BATCH_OPERATOR_ADDRESS: string;
    }
  }
}

// If this file has no import/export statements (i.e. is a script)
// convert it into a module by adding an empty export statement.
export {};
