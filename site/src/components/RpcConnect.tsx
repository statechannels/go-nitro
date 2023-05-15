import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import React, { useState } from "react";
import Typography from "@mui/material/Typography";

export type RPCConnectProps = {
  url: string;
  setUrl: (url: string) => void;
};

export default function RpcConnect({ url, setUrl }: RPCConnectProps) {
  const [urlToEdit, setUrlToEdit] = useState(url);
  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    setUrlToEdit(e.target.value);
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setUrl(urlToEdit);
  };

  return (
    <form
      style={{ display: "flex", alignItems: "center" }}
      onSubmit={handleSubmit}
    >
      <Typography display="inline">Nitro RPC Connect:</Typography>
      <TextField sx={{ ml: 2 }} value={urlToEdit} onChange={handleChange} />
      <Button type="submit" sx={{ ml: 2 }} variant="contained">
        Connect
      </Button>
    </form>
  );
}
