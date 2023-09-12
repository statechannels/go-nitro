import { ChangeEvent, useEffect, useState } from "react";
import { NitroRpcClient } from "@statechannels/nitro-rpc-client";
import { PaymentChannelInfo } from "@statechannels/nitro-rpc-client/src/types";
import {
  Button,
  TextField,
  Box,
  Checkbox,
  FormControlLabel,
  LinearProgress,
} from "@mui/material";

const QUERY_KEY = "rpcUrl";

import "./App.css";
import { fetchFile, fetchFileInChunks } from "./file";
import ChannelDetails from "./ChannelDetails";
const provider = "0xbbb676f9cff8d242e9eac39d063848807d3d1d94";
const hub = "0x111a00868581f73ab42feef67d235ca09ca1e8db";
const defaultNitroRPCUrl = "localhost:4005/api/v1";
const defaultFileUrl = "http://localhost:5511/test.txt";
const defaultPaymentChannelAmount = 100_000;

function App() {
  const url =
    new URLSearchParams(window.location.search).get(QUERY_KEY) ??
    defaultNitroRPCUrl;

  const [nitroClient, setNitroClient] = useState<NitroRpcClient | null>(null);

  const [paymentChannelId, setPaymentChannelId] = useState<string>("");
  const [paymentChannelInfo, setPaymentChannelInfo] = useState<
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

  const updateChannelInfo = async (channelId: string) => {
    if (channelId == "") {
      throw new Error("Empty channel id provided");
    }
    const paymentChannel = await nitroClient?.GetPaymentChannel(channelId);
    setPaymentChannelInfo(paymentChannel);
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

  const createPaymentChannel = async () => {
    if (!nitroClient) {
      setErrorText("Nitro client not initialized");
      return;
    }
    const result = await nitroClient.CreatePaymentChannel(
      provider,
      [hub],
      defaultPaymentChannelAmount
    );
    console.log(result);
    setPaymentChannelId(result.ChannelId);
    updateChannelInfo(result.ChannelId);

    // TODO: Slightly hacky but we wait a beat before querying so we see the updated balance
    setTimeout(() => {
      updateChannelInfo(result.ChannelId);
    }, 1000);
  };

  const fetchAndDownloadFile = async () => {
    setErrorText("");

    if (!nitroClient) {
      setErrorText("Nitro client not initialized");
      return;
    }
    if (!paymentChannelInfo) {
      setErrorText("No payment channel to use");
      return;
    }

    try {
      setMicroPaymentProgress(0);
      const file = useMicroPayments
        ? await fetchFileInChunks(
            chunkSize,
            fileUrl,
            costPerByte,
            paymentChannelInfo.ID,
            nitroClient,
            setMicroPaymentProgress
          )
        : await fetchFile(
            fileUrl,
            costPerByte * dataSize,
            paymentChannelInfo.ID,
            nitroClient
          );

      triggerFileDownload(file);

      // TODO: Slightly hacky but we wait a beat before querying so we see the updated balance
      setTimeout(() => {
        updateChannelInfo(paymentChannelInfo.ID);
      }, 50);
    } catch (e: unknown) {
      setErrorText((e as Error).message);
    }
  };

  return (
    <Box>
      <Box p={10} minHeight={200}>
        <Button
          id="createChannel"
          onClick={() => {
            createPaymentChannel();
          }}
          disabled={paymentChannelId != ""}
        >
          Create Channel
        </Button>
        <ChannelDetails info={paymentChannelInfo} />
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
