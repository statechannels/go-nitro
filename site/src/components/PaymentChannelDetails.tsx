import {
  Box,
  LinearProgress,
  Typography,
  Stack,
  Button,
  SvgIcon,
} from "@mui/material";
import { FC, useEffect, useMemo, useState } from "react";
import { makeStyles } from "tss-react/mui";
import { PhoneArrowUpRightIcon, UserIcon } from "@heroicons/react/24/outline";

interface PaymentChannelDetails {
  channelID: string;
  counterparty: string; // address for now - can abstract to something richer later
  capacity: number; // total value locked in channel
  myBalance: number; // balance of the viewing participant
  status: "prefund" | "running" | "unresponsive-peer" | "under-challenge";
}

const useStyles = makeStyles()(() => ({
  typography: {
    marginTop: "0 !important",
  },
  leftPeer: {},
  rightPeer: {
    marginTop: "-6rem !important",
  },
  icons: {
    lineHeight: "0 !important",
    fontSize: "3rem",
  },
  iconLeft: {
    fontSize: "6rem",
  },
  iconRight: {
    fontSize: "2.5rem",
  },
}));

const PaymentChannelDetails: FC<PaymentChannelDetails> = ({
  channelID,
  counterparty,
  capacity,
  myBalance,
  status,
}: PaymentChannelDetails) => {
  const [progress, setProgress] = useState<number>(myBalance);
  const { classes, cx } = useStyles();

  useEffect(() => {
    setProgress((myBalance * 100) / capacity);
  }, [myBalance, capacity]);

  const capitalizedStatus = useMemo(() => {
    return status.charAt(0).toUpperCase() + status.slice(1);
  }, [status]);

  return (
    <Stack
      direction="column"
      alignItems="center"
      justifyContent="space-between"
      spacing={20}
    >
      <Stack direction="column" alignItems="center" width="100%" spacing={2}>
        <SvgIcon fontSize="large">
          <PhoneArrowUpRightIcon />
        </SvgIcon>
        <Typography variant="h6" component="h6">
          Outbound Payment Channel
        </Typography>
        <Typography variant="h6" component="h6">
          {channelID}
        </Typography>
      </Stack>
      <Stack
        direction="row"
        alignItems="center"
        justifyContent="center"
        width="100%"
        spacing={2}
      >
        <Stack
          minWidth="fit-content"
          direction="column"
          alignItems="center"
          spacing={2}
        >
          <Typography
            variant="h2"
            component="h2"
            className={cx(classes.icons, classes.iconLeft)}
          >
            <SvgIcon fontSize="inherit">
              <UserIcon />
            </SvgIcon>
          </Typography>
          <Typography
            variant="body1"
            component="span"
            className={classes.typography}
          >
            {myBalance} wei
          </Typography>
        </Stack>
        <Stack
          direction="row"
          alignItems="center"
          justifyContent="center"
          spacing={2}
          width="100%"
        >
          <Box sx={{ display: "flex", alignItems: "center", width: "100%" }}>
            <Box sx={{ width: "100%" }}>
              <LinearProgress
                variant="determinate"
                value={progress}
                color={"primary"}
              />
            </Box>
          </Box>
        </Stack>
        <Stack
          minWidth="fit-content"
          direction="column"
          alignItems="center"
          spacing={2}
          className={classes.rightPeer}
        >
          <Typography
            variant="h2"
            component="h2"
            className={cx(classes.icons, classes.iconRight)}
          >
            <SvgIcon fontSize="inherit">
              <UserIcon />
            </SvgIcon>
          </Typography>
          <Typography
            variant="body1"
            component="span"
            className={classes.typography}
          >
            {capacity - myBalance} wei
          </Typography>
          <Typography
            variant="body1"
            component="span"
            className={classes.typography}
          >
            {counterparty}
          </Typography>
        </Stack>
      </Stack>
      <Stack direction="column" alignItems="center" spacing={2}>
        <Button variant="contained">1 wei</Button>
      </Stack>
      <Stack direction="column" alignItems="center" spacing={2}>
        <Typography
          variant="body1"
          component="span"
          className={classes.typography}
        >
          {capitalizedStatus}
        </Typography>
      </Stack>
    </Stack>
  );
};

export default PaymentChannelDetails;
