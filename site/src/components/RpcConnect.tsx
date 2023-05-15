import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import React, { useState } from "react";
import Typography from "@mui/material/Typography";

import { QUERY_KEY } from "../constants";

export type RPCConnectProps = {
  url: string;
};

export default function RpcConnect({ url }: RPCConnectProps) {
  const [urlToEdit, setUrlToEdit] = useState(url);
  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    setUrlToEdit(e.target.value);
  }

  return (
    <form style={{ display: "flex", alignItems: "center" }}>
      <Typography display="inline">Nitro RPC Connect:</Typography>
      <TextField
        sx={{ ml: 2 }}
        name={QUERY_KEY}
        value={urlToEdit}
        onChange={handleChange}
      />
      <Button type="submit" sx={{ ml: 2 }} variant="contained">
        Connect
      </Button>
    </form>
  );
}
