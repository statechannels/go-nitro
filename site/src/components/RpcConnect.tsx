import { useState } from "react";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Box from "@mui/material/Box";
import * as React from "react";

export default function RpcConnect() {
  const [text, setText] = useState("localhost:4005");

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    setText(e.target.value);
  }

  return (
    <Box>
      Nitro RPC Connect:
      <TextField value={text} onChange={handleChange} />
      <Button
        variant="contained"
        onClick={() => {
          console.log("Button click");
        }}
      >
        Connect
      </Button>
    </Box>
  );
}
