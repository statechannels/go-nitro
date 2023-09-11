import { Table, TableRow, TableCell, TableBody } from "@mui/material";
import { PaymentChannelInfo } from "@statechannels/nitro-rpc-client/src/types";

function ChannelDetails({
  info,
}: {
  info: PaymentChannelInfo | undefined;
}): JSX.Element {
  return (
    <Table>
      <TableBody>
        <TableRow>
          <TableCell>Channel Id</TableCell>
          <TableCell>{info && info.ID}</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>Paid so far</TableCell>
          <TableCell>
            {info &&
              // TODO: We shouldn't have to cast to a BigInt here, the client should return a BigInt
              BigInt(info?.Balance.PaidSoFar).toString(10)}
          </TableCell>
        </TableRow>
        <TableRow>
          <TableCell>Remaining funds</TableCell>
          <TableCell>
            {info &&
              // TODO: We shouldn't have to cast to a BigInt here, the client should return a BigInt
              BigInt(info?.Balance.RemainingFunds).toString(10)}
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
  );
}

export default ChannelDetails;
