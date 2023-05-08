import { Outcome, RPCMethod, RPCRequestAndResponses } from "./types";

/**
 * createDirectFundOutcome creates a basic outcome for a directly funded channel
 *
 * @param asset - The asset to fund the channel with
 * @param alpha - The address of the first participant
 * @param beta - The address of the second participant
 * @returns An outcome for a directly funded channel with 100 wei allocated to each participant
 */
export function createDirectFundOutcome(
  asset: string,
  alpha: string,
  beta: string
): Outcome {
  return [
    {
      Asset: asset,
      AssetMetadata: {
        AssetType: 0,
        Metadata: null,
      },

      Allocations: [
        {
          Destination: convertAddressToBytes32(alpha),
          Amount: 100,
          AllocationType: 0,
          Metadata: null,
        },
        {
          Destination: convertAddressToBytes32(beta),
          Amount: 100,
          AllocationType: 0,
          Metadata: null,
        },
      ],
    },
  ];
}

/**
 * Left pads a 20 byte address hex string with zeros until it is a 32 byte hex string
 * e.g.,
 * 0x9546E319878D2ca7a21b481F873681DF344E0Df8 becomes
 * 0x0000000000000000000000009546E319878D2ca7a21b481F873681DF344E0Df8
 *
 * @param address - 20 byte hex string
 * @returns 32 byte padded hex string
 */
export function convertAddressToBytes32(address: string): string {
  const digits = address.startsWith("0x") ? address.substring(2) : address;
  return `0x${digits.padStart(24, "0")}`;
}

/**
 * generateRequest is a helper function that generates a request object for the given method and params
 *
 * @param method - The RPC method to generate a request for
 * @param params - The params to include in the request
 * @returns A request object of the correct type
 */
export function generateRequest<
  K extends RPCMethod,
  T extends RPCRequestAndResponses[K][0]
>(method: K, params: T["params"]): T {
  return {
    jsonrpc: "2.0",
    method,
    params,
    id: Date.now(),
  } as T; // TODO: We shouldn't have to cast here
}
