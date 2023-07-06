import { ChangeEvent, useEffect, useState } from "react";
import { NitroRpcClient } from "@statechannels/nitro-rpc-client";
import { PaymentChannelInfo } from "@statechannels/nitro-rpc-client/src/types";
import {
  Select,
  MenuItem,
  SelectChangeEvent,
  Button,
  TextField,
  Box,
  Table,
  TableRow,
  TableCell,
  TableBody,
} from "@mui/material";
import axios, { isAxiosError } from "axios";

const QUERY_KEY = "rpcUrl";

import "./App.css";

function App() {
  const retrievalProvider = "0xbbb676f9cff8d242e9eac39d063848807d3d1d94";
  const hub = "0x111a00868581f73ab42feef67d235ca09ca1e8db";
  const defaultUrl = "localhost:4005/api/v1";

  const url =
    new URLSearchParams(window.location.search).get(QUERY_KEY) ?? defaultUrl;

  const [nitroClient, setNitroClient] = useState<NitroRpcClient | null>(null);
  const [paymentChannels, setPaymentChannels] = useState<PaymentChannelInfo[]>(
    []
  );
  const [selectedChannel, setSelectedChannel] = useState<string>("");
  const [selectedChannelInfo, setSelectedChannelInfo] = useState<
    PaymentChannelInfo | undefined
  >();

  // TODO: For now the default is a hardcoded value based on a local file
  // If you're running this locally you'll need to override this value
  // Ideally we should just query boost/lotus for the list of available payloads>
  const [payloadId, setPayloadId] = useState<string>(
    "bafk2bzacec3jst4tkh424chatp273o6rxvipfg54kphd56gaxobpcdtr2sgco"
  );

  const [paymentAmount, setPaymentAmount] = useState<number>(5);

  const [errorText, setErrorText] = useState<string>("");

  useEffect(() => {
    NitroRpcClient.CreateHttpNitroClient(url).then((c) => setNitroClient(c));
  }, [url]);

  // Fetch all the payment channels for the retrieval provider
  useEffect(() => {
    if (nitroClient) {
      // TODO: We should consider adding a API function so this ins't as painful
      nitroClient.GetAllLedgerChannels().then((ledgers) => {
        for (const l of ledgers) {
          if (l.Balance.Them != hub) continue;

          nitroClient.GetPaymentChannelsByLedger(l.ID).then((payChs) => {
            const withProvider = payChs.filter(
              (p) => p.Balance.Payee == retrievalProvider
            );
            setPaymentChannels(withProvider);
          });
        }
      });
    }
  }, [nitroClient]);

  const updateChannelInfo = async (channelId: string) => {
    const paymentChannel = await nitroClient?.GetPaymentChannel(channelId);
    setSelectedChannelInfo(paymentChannel);
  };

  const handleSelectedChannelChanged = async (event: SelectChangeEvent) => {
    setSelectedChannel(event.target.value);
    updateChannelInfo(event.target.value);
  };

  const updatePayloadId = (e: ChangeEvent<HTMLInputElement>) => {
    setPayloadId(e.target.value);
  };

  const updatePaymentAmount = (e: ChangeEvent<HTMLInputElement>) => {
    setPaymentAmount(parseInt(e.target.value));
  };
  const fetchFile = async () => {
    setErrorText("");

    if (!nitroClient) {
      setErrorText("Nitro client not initialized");
      return;
    }
    if (!selectedChannel) {
      setErrorText("Please select a channel");
      return;
    }

    const voucher = await nitroClient.CreateVoucher(
      selectedChannel,
      paymentAmount
    );
    // TODO: Slightly hacky but we wait a beat before querying so we see the updated balance
    setTimeout(() => {
      updateChannelInfo(selectedChannel);
    }, 50);

    try {
      const result = await axios.get(
        `http://localhost:7777/ipfs/${payloadId}?channelId=${voucher.ChannelId}&amount=${voucher.Amount}&signature=${voucher.Signature}`,
        {
          responseType: "blob", // This lets us download the file
          headers: {
            Accept: "*/*", // TODO: Do we need to specify this?
          },
        }
      );

      // This will prompt the browser to download the file
      const blob = result.data;
      const blobUrl = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = blobUrl;
      a.download = "fetched-file-from-ipfs";
      a.click();
      a.remove();
      window.URL.revokeObjectURL(blobUrl);
    } catch (e) {
      if (isAxiosError(e)) {
        const { message } = e;
        e.response?.data.text().then((text: string) => {
          setErrorText(`${message}: ${text}`);
        });
      } else {
        setErrorText(JSON.stringify(e));
      }
    }
  };

  return (
    <Box>
      <Box p={10} minHeight={200}>
        <Select
          label="virtual channels"
          onChange={handleSelectedChannelChanged}
          value={selectedChannel}
        >
          {...paymentChannels.map((p) => (
            <MenuItem value={p.ID}>{p.ID}</MenuItem>
          ))}
        </Select>

        <Table>
          <TableBody>
            <TableRow>
              <TableCell>Paid so far</TableCell>
              <TableCell>
                {selectedChannelInfo &&
                  // TODO: We shouldn't have to cast to a BigInt here, the client should return a BigInt
                  BigInt(selectedChannelInfo?.Balance.PaidSoFar).toString(10)}
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>Remaining funds</TableCell>
              <TableCell>
                {selectedChannelInfo &&
                  // TODO: We shouldn't have to cast to a BigInt here, the client should return a BigInt
                  BigInt(selectedChannelInfo?.Balance.RemainingFunds).toString(
                    10
                  )}
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>Payee</TableCell>
              <TableCell>
                {selectedChannelInfo && selectedChannelInfo.Balance.Payee}
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>Payer</TableCell>
              <TableCell>
                {selectedChannelInfo && selectedChannelInfo.Balance.Payer}
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </Box>
      <Box>
        <TextField
          fullWidth={true}
          label="Payload Id"
          onChange={updatePayloadId}
          value={payloadId}
        ></TextField>
      </Box>
      <br></br>
      <Box>
        <TextField
          label="Payment Amount"
          onChange={updatePaymentAmount}
          value={paymentAmount}
          type="number"
        ></TextField>
        <Button onClick={fetchFile}>Fetch</Button>
        <Box>{errorText}</Box>
      </Box>
    </Box>
  );
}

export default App;
