export const getAddressMock = {
  jsonrpc: "2.0",
  id: 1684516515941,
  result: "0xAAA6628Ec44A8a742987EF3A114dDFE2D4F7aDCE",
  error: null,
};

export const allLedgerChannelMock = {
  jsonrpc: "2.0",
  id: 1684511523089,
  result: [
    {
      ID: "0x411ae0023593f5a2c9fe99c8017ff7c1a78c0071a072dc93ab2acfd7a87f1059",
      Status: "Open",
      Balance: {
        AssetAddress: "0x0000000000000000000000000000000000000000",
        Hub: "0x111a00868581f73ab42feef67d235ca09ca1e8db",
        Client: "0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce",
        HubBalance: 997000,
        ClientBalance: 997000,
      },
    },
    {
      ID: "0x14ddcda18c2db429866ae79c308ba4542ef19d31a531eb6a4283bdafb1efed3b",
      Status: "Complete",
      Balance: {
        AssetAddress: "0x0000000000000000000000000000000000000000",
        Hub: "0x111a00868581f73ab42feef67d235ca09ca1e8db",
        Client: "0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce",
        HubBalance: 1000,
        ClientBalance: 1000,
      },
    },
    {
      ID: "0xacc0aa3b8271d49c28259d41e2ea28bcbb80b0cefb75b0ad0a655b865e48db69",
      Status: "Complete",
      Balance: {
        AssetAddress: "0x0000000000000000000000000000000000000000",
        Hub: "0x111a00868581f73ab42feef67d235ca09ca1e8db",
        Client: "0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce",
        HubBalance: 1000,
        ClientBalance: 1000,
      },
    },
  ],
  error: null,
};

export const getLedgerChannelMock = {
  jsonrpc: "2.0",
  id: 1684516515960,
  result: {
    ID: "0x411ae0023593f5a2c9fe99c8017ff7c1a78c0071a072dc93ab2acfd7a87f1059",
    Status: "Open",
    Balance: {
      AssetAddress: "0x0000000000000000000000000000000000000000",
      Hub: "0x111a00868581f73ab42feef67d235ca09ca1e8db",
      Client: "0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce",
      HubBalance: 9970,
      ClientBalance: 9970,
    },
  },
  error: null,
};

export const getPaymentChannelsByLedgerMock = {
  jsonrpc: "2.0",
  id: 1684516515966,
  result: [
    {
      ID: "0x128c577ea4da25d7c91df9efa88ce8df4d41a262c969f3dc21558180ec7af044",
      Status: "Open",
      Balance: {
        AssetAddress: "0x0000000000000000000000000000000000000000",
        Payee: "0xbbb676f9cff8d242e9eac39d063848807d3d1d94",
        Payer: "0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce",
        PaidSoFar: 900,
        RemainingFunds: 100,
      },
    },
    {
      ID: "0xa3c1dd747ebe7e0886574b4405451b14abd63d107989821aef64b28ab908d215",
      Status: "Open",
      Balance: {
        AssetAddress: "0x0000000000000000000000000000000000000000",
        Payee: "0xbbb676f9cff8d242e9eac39d063848807d3d1d94",
        Payer: "0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce",
        PaidSoFar: 100,
        RemainingFunds: 900,
      },
    },
    {
      ID: "0xddc70aa382e7bdeddf982683c4a2e99ce0e73f9b4f5d1b84d77203f0c7971e7f",
      Status: "Open",
      Balance: {
        AssetAddress: "0x0000000000000000000000000000000000000000",
        Payee: "0xbbb676f9cff8d242e9eac39d063848807d3d1d94",
        Payer: "0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce",
        PaidSoFar: 500,
        RemainingFunds: 500,
      },
    },
  ],
  error: null,
};
