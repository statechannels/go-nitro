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
      {"."}
    </Typography>
  );
}
