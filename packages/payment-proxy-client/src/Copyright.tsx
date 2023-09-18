import { Link, Typography } from "@mui/material";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function Copyright(props: any) {
  return (
    <Typography
      variant="body2"
      color="text.secondary"
      align="center"
      {...props}
    >
      {"Copyright © "}
      <Link color="inherit" href="https://statechannels.org/">
        statechannels.org
      </Link>{" "}
      {new Date().getFullYear()}
      {" | "}
      <Link
        href="https://statechannels.notion.site/Filecoin-Paid-Retrieval-Demo-bf6ad9ec92a74e139331ce77900305fc?pvs=4"
        variant="body2"
      >
        How does this work?
      </Link>
    </Typography>
  );
}
