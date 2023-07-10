import { getAndValidateResult } from "./serde";

describe("get_address", () => {
  it("success: validate response string", () => {
    const getAddressResponse = {
      jsonrpc: "2.0",
      id: 168513765,
      result: "0x111A00868581f73AB42FEEF67D235Ca09ca1E8db",
    };
    const validatedResponse = getAndValidateResult(
      getAddressResponse,
      "get_address"
    );
    expect(validatedResponse).toEqual(getAddressResponse.result);
  });
});

describe("get_ledger_channel", () => {
  it("success: validate response object", () => {
    const getLedgerChannelResponse = {
      jsonrpc: "2.0",
      id: 168513765,
      result: {
        ID: "0x586d127530f69177d790bb940eae132922e7648c29264648af5375de2c19e270",
        Status: "Open",
        Balance: {
          AssetAddress: "0x0000000000000000000000000000000000000000",
          Hub: "0x111a00868581f73ab42feef67d235ca09ca1e8db",
          Client: "0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce",
          HubBalance: "0xf368a",
          ClientBalance: "0xf3686",
        },
      },
    };

    const validatedGetLedgerChannelResponse = {
      ID: "0x586d127530f69177d790bb940eae132922e7648c29264648af5375de2c19e270",
      Status: "Open",
      Balance: {
        AssetAddress: "0x0000000000000000000000000000000000000000",
        Hub: "0x111a00868581f73ab42feef67d235ca09ca1e8db",
        Client: "0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce",
        HubBalance: 997002n,
        ClientBalance: 996998n,
      },
    };

    const validatedResponse = getAndValidateResult(
      getLedgerChannelResponse,
      "get_ledger_channel"
    );
    expect(validatedResponse).toEqual(validatedGetLedgerChannelResponse);
  });
});

describe("create_ledger_channel", () => {
  it("error", () => {
    const failedCreateLedgerResponse = {
      jsonrpc: "2.0",
      id: 168513765,
      error: {
        code: -32603,
        message: "Internal Server Error",
      },
    };

    try {
      getAndValidateResult(failedCreateLedgerResponse, "create_ledger_channel");
    } catch (err) {
      if (err instanceof Error) {
        expect(err.message).toEqual("jsonrpc response: Internal Server Error");
      } else {
        expect(false);
      }
    }
  });

  it("success: validate response object", () => {
    const successCreateLedgerResponse = {
      jsonrpc: "2.0",
      id: 995772692,
      result: {
        Id: "123",
        ChannelId: "456",
      },
    };

    const validatedResponse = getAndValidateResult(
      successCreateLedgerResponse,
      "create_ledger_channel"
    );
    expect(validatedResponse).toEqual(successCreateLedgerResponse.result);
  });
});
