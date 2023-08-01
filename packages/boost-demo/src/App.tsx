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

const QUERY_KEY = "rpcUrl";

import "./App.css";
import { fetchFile } from "./utils";

const retrievalProvider = "0xbbb676f9cff8d242e9eac39d063848807d3d1d94";
const hub = "0x111a00868581f73ab42feef67d235ca09ca1e8db";
const defaultNitroRPCUrl = "localhost:4005/api/v1";
const defaultFileUrl = "http://localhost:5511/test.txt";
const costPerByte = 1;

function App() {
  const url =
    new URLSearchParams(window.location.search).get(QUERY_KEY) ??
    defaultNitroRPCUrl;

  const [nitroClient, setNitroClient] = useState<NitroRpcClient | null>(null);
  const [paymentChannels, setPaymentChannels] = useState<PaymentChannelInfo[]>(
    []
  );
  const [selectedChannel, setSelectedChannel] = useState<string>("");
  const [selectedChannelInfo, setSelectedChannelInfo] = useState<
    PaymentChannelInfo | undefined
  >();

  const [fileUrl, setFileUrl] = useState<string>(defaultFileUrl);

  const [paymentAmount, setPaymentAmount] = useState<number>(5);

  const [errorText, setErrorText] = useState<string>("");

  useEffect(() => {
    NitroRpcClient.CreateHttpNitroClient(url)
      .then((c) => setNitroClient(c))
      .catch((e) => {
        setErrorText(e.message);
      });
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

  const updateProxyUrl = (e: ChangeEvent<HTMLInputElement>) => {
    setFileUrl(e.target.value);
  };

  const updatePaymentAmount = (e: ChangeEvent<HTMLInputElement>) => {
    setPaymentAmount(parseInt(e.target.value));
  };

  const triggerFileDownload = (file: File) => {
    // This will prompt the browser to download the file
    const blob = new Blob([file], { type: file.type });

    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = file.name;
    link.click();
    URL.revokeObjectURL(url);
  };

  const fetchAndDownloadFile = async () => {
    setErrorText("");

    if (!nitroClient) {
      setErrorText("Nitro client not initialized");
      return;
    }
    if (!selectedChannel) {
      setErrorText("Please select a channel");
      return;
    }

    try {
      const file = await fetchFile(
        fileUrl,
        costPerByte,
        paymentAmount,
        selectedChannel,
        nitroClient
      );

      console.log(file);

      triggerFileDownload(file);

      // TODO: Slightly hacky but we wait a beat before querying so we see the updated balance
      setTimeout(() => {
        updateChannelInfo(selectedChannel);
      }, 50);
    } catch (e: unknown) {
      setErrorText((e as Error).message);
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
          label="Proxy URL"
          onChange={updateProxyUrl}
          value={fileUrl}
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
        <Button onClick={fetchAndDownloadFile}>Fetch</Button>
        <Box>{errorText}</Box>
      </Box>
    </Box>
  );
}

export default App;
