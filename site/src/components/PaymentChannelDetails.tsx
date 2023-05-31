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
import { ChannelStatus } from "@statechannels/nitro-rpc-client/src/types";

interface PaymentChannelDetails {
  channelID: string;
  payee: string;
  payer: string;
  paidSoFar: bigint;
  remainingFunds: bigint;
  status: ChannelStatus;
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

const shortString = (value: string, count: number) => {
  return `${value.slice(0, count)}...`;
};

const PaymentChannelDetails: FC<PaymentChannelDetails> = ({
  channelID,
  payee,
  payer,
  paidSoFar,
  remainingFunds,
  status,
}: PaymentChannelDetails) => {
  const [progress, setProgress] = useState<number>(0);
  const { classes, cx } = useStyles();

  useEffect(() => {
    setProgress(
      (Number(paidSoFar) / (Number(remainingFunds) + Number(paidSoFar))) * 100
    );
  }, [paidSoFar, remainingFunds]);

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
          {shortString(channelID, 5)}
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
            {paidSoFar.toString()} wei
          </Typography>
          <Typography
            variant="body1"
            component="span"
            className={classes.typography}
          >
            {shortString(payee, 8)}
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
            {remainingFunds.toString()} wei
          </Typography>
          <Typography
            variant="body1"
            component="span"
            className={classes.typography}
          >
            {shortString(payer, 8)}
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
