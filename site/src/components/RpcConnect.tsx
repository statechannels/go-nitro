import { useState } from "react";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import * as React from "react";
import "../styles/styles.css";
import { Container, Grid, Typography } from "@material-ui/core";

export default function RpcConnect() {
  const [text, setText] = useState("localhost:4005");

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    setText(e.target.value);
  }

  return (
    <Grid>
      <Typography display="inline">Nitro RPC Connect:</Typography>
      <TextField value={text} onChange={handleChange} />
      <Button
        variant="contained"
        onClick={() => {
          console.log("Button click");
        }}
      >
        Connect
      </Button>
    </Grid>
  );
}
