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
  InputLabel,
  Checkbox,
  FormControlLabel,
  LinearProgress,
} from "@mui/material";

const QUERY_KEY = "rpcUrl";

import "./App.css";
import { fetchFile, fetchFileInChunks } from "./file";
const provider = "0xbbb676f9cff8d242e9eac39d063848807d3d1d94";
const hub = "0x111a00868581f73ab42feef67d235ca09ca1e8db";
const defaultNitroRPCUrl = "localhost:4005/api/v1";
const defaultFileUrl = "http://localhost:5511/test.txt";

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

  const [costPerByte, setCostPerByte] = useState<number>(1);
  const [dataSize, setDataSize] = useState<number>(12);
  const [totalCost, setTotalCost] = useState<number>(costPerByte * dataSize);
  const [errorText, setErrorText] = useState<string>("");
  const [chunkSize, setChunkSize] = useState<number>(100);
  const [useMicroPayments, setUseMicroPayments] = useState<boolean>(false);
  const [microPaymentProgress, setMicroPaymentProgress] = useState<number>(0);
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
              (p) => p.Balance.Payee == provider
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

  const proxyUrlChanged = (e: ChangeEvent<HTMLInputElement>) => {
    setFileUrl(e.target.value);
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
      setMicroPaymentProgress(0);
      const file = useMicroPayments
        ? await fetchFileInChunks(
            chunkSize,
            fileUrl,
            costPerByte,
            selectedChannel,
            nitroClient,
            setMicroPaymentProgress
          )
        : await fetchFile(
            fileUrl,
            costPerByte * dataSize,
            selectedChannel,
            nitroClient
          );

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
        <InputLabel id="select-channel">Select a payment channel</InputLabel>
        <Select
          onChange={handleSelectedChannelChanged}
          value={selectedChannel}
          inputProps={{
            id: "select-channel",
          }}
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
          onChange={proxyUrlChanged}
          value={fileUrl}
        ></TextField>
      </Box>
      <FormControlLabel
        label="Use micropayments"
        control={
          <Checkbox
            onChange={(e: ChangeEvent<HTMLInputElement>) =>
              setUseMicroPayments(e.target.checked)
            }
            value={useMicroPayments}
          ></Checkbox>
        }
      />
      <Box visibility={useMicroPayments ? "visible" : "hidden"}>
        <TextField
          label="Cost Per Byte(wei)"
          onChange={(e: ChangeEvent<HTMLInputElement>) => {
            setCostPerByte(parseInt(e.target.value));
          }}
          value={costPerByte}
          type="number"
        ></TextField>
        <TextField
          label="Chunk size(bytes)"
          onChange={(e: ChangeEvent<HTMLInputElement>) => {
            setChunkSize(parseInt(e.target.value));
          }}
          value={chunkSize}
          type="number"
        ></TextField>
      </Box>
      <Box visibility={useMicroPayments ? "hidden" : "visible"}>
        <TextField
          label="Cost Per Byte(wei)"
          onChange={(e: ChangeEvent<HTMLInputElement>) => {
            setCostPerByte(parseInt(e.target.value));
            setTotalCost(dataSize * parseInt(e.target.value));
          }}
          value={costPerByte}
          type="number"
        ></TextField>
        <TextField
          label="Data Size(bytes)"
          onChange={(e: ChangeEvent<HTMLInputElement>) => {
            setDataSize(parseInt(e.target.value));
            setTotalCost(costPerByte * parseInt(e.target.value));
          }}
          value={dataSize}
          type="number"
        ></TextField>
        <TextField
          inputProps={{ readOnly: true }}
          label="Total cost(wei)"
          value={totalCost}
          type="number"
        ></TextField>
      </Box>

      <Button onClick={fetchAndDownloadFile}>
        {useMicroPayments ? "Fetch with micropayments" : "Fetch"}
      </Button>
      <Box visibility={useMicroPayments ? "visible" : "hidden"}>
        <LinearProgress value={microPaymentProgress} variant="determinate" />
      </Box>

      <Box>{errorText}</Box>
    </Box>
  );
}

export default App;
