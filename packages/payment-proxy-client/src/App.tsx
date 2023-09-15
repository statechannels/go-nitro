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
const provider = import.meta.env.VITE_PROVIDER;
const hub = import.meta.env.VITE_HUB;
const defaultNitroRPCUrl = import.meta.env.VITE_NITRO_RPC_URL;
const defaultFileUrl = import.meta.env.VITE_FILE_URL;
const CHANNEL_ID_KEY = "channelId";
const initialChannelBalance = parseInt(
  import.meta.env.VITE_INITIAL_CHANNEL_BALANCE,
  10
);

const costPerByte = 1;
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

  const [dataSize, setDataSize] = useState<number>(12);

  const [errorText, setErrorText] = useState<string>("");
  const [chunkSize, setChunkSize] = useState<number>(100);
  const [useMicroPayments, setUseMicroPayments] = useState<boolean>();
  const [microPaymentProgress, setMicroPaymentProgress] = useState<number>(0);
  useEffect(() => {
    NitroRpcClient.CreateHttpNitroClient(url)
      .then((c) => setNitroClient(c))
      .catch((e) => {
        console.error(e);
        setErrorText(e.message);
      });
  }, [url]);

  useEffect(() => {
    if (!nitroClient) return;
    const fetchedId = localStorage.getItem(CHANNEL_ID_KEY);
    if (fetchedId && fetchedId != "") {
      setPaymentChannelId(fetchedId);

      nitroClient
        .GetPaymentChannel(fetchedId)
        .then((paymentChannel) => {
          console.log(paymentChannel);
          setPaymentChannelInfo(paymentChannel);
        })
        .catch(() => {
          console.warn(
            `Channel from storage ${fetchedId} not found. Removing from storage...`
          );
          localStorage.removeItem(CHANNEL_ID_KEY);
          setPaymentChannelId("");
        });
    }
  }, [nitroClient]);

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
      initialChannelBalance
    );

    localStorage.setItem(CHANNEL_ID_KEY, result.ChannelId);
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
      console.error(e);
      setErrorText((e as Error).message);
    }
  };

  return (
    <Box>
      <Box p={10}>
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
      {useMicroPayments && (
        <Box>
          <TextField
            label="Chunk size(bytes)"
            onChange={(e: ChangeEvent<HTMLInputElement>) => {
              setChunkSize(parseInt(e.target.value));
            }}
            value={chunkSize}
            type="number"
          ></TextField>
        </Box>
      )}
      {!useMicroPayments && (
        <Box>
          <TextField
            label="Payment amount (bytes)"
            onChange={(e: ChangeEvent<HTMLInputElement>) => {
              setDataSize(parseInt(e.target.value));
            }}
            value={dataSize}
            type="number"
          ></TextField>
        </Box>
      )}
      <Button
        id="createChannel"
        onClick={() => {
          createPaymentChannel();
        }}
        disabled={paymentChannelId != "" || !nitroClient}
      >
        Create Channel
      </Button>
      <Button onClick={fetchAndDownloadFile} disabled={paymentChannelId == ""}>
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
