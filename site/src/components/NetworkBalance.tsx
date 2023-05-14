import bigDecimal from "js-big-decimal-esm";
import { PieChart } from "react-minimal-pie-chart";
import "./NetworkBalance.scss";
import Typography from "@mui/material/Typography";
import LinearProgress, {
  LinearProgressProps,
} from "@mui/material/LinearProgress";
import Box from "@mui/material/Box";

import colors from "../styles/colors.module.css";
import { prettyPrintWei } from "../utils";

// This prevents a runtime error in storybook
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore

BigInt.prototype.toJSON = function (): string {
  return this.toString();
};

export type NetworkBalanceProps = {
  asTable?: boolean;
  status: "running" | "unresponsive-peer" | "under-challenge";
  myBalanceFree: bigint;
  theirBalanceFree: bigint;
  lockedBalances: VirtualChannelBalanceProps[];
};

export type VirtualChannelBalanceProps = {
  budget: bigint;
  myPercentage: number;
};

function percentageOfTotal(quantity: bigint, total: bigint): number {
  return parseFloat(bigDecimal.divide(quantity, total, 8)) * 100;
}

/**
 * returns intermediate color between colorA and colorB, based on percentage
 *
 * @param colorA
 * @param colorB
 * @param percentage
 * @returns
 */
function interpolateColor(
  colorA: string,
  colorB: string,
  percentage: number
): string {
  const hex = (x: string) => parseInt(x, 16);

  const r = Math.round(
    hex(colorA.substr(1, 2)) * (1 - percentage) +
      hex(colorB.substr(1, 2)) * percentage
  ).toString(16);

  const g = Math.round(
    hex(colorA.substr(3, 2)) * (1 - percentage) +
      hex(colorB.substr(3, 2)) * percentage
  ).toString(16);

  const b = Math.round(
    hex(colorA.substr(5, 2)) * (1 - percentage) +
      hex(colorB.substr(5, 2)) * percentage
  ).toString(16);

  return `#${r}${g}${b}`;
}

function sortToExtremes(
  channels: VirtualChannelBalanceProps[]
): VirtualChannelBalanceProps[] {
  channels.sort((a, b) => a.myPercentage - b.myPercentage);

  const toExtremes: VirtualChannelBalanceProps[] = [];

  let next: "l" | "r" = "l";

  // from smallest to largest, alternate adding to the left or right of the array
  // End result is smallest in the middle, largest to the outsides
  for (let i = 0; i < channels.length; i++) {
    if (next === "l") {
      toExtremes.unshift(channels[i]);
      next = "r";
    } else {
      toExtremes.push(channels[i]);
      next = "l";
    }
  }

  return toExtremes;
}

export const NetworkBalance: React.FC<NetworkBalanceProps> = (props) => {
  const { myBalanceFree, theirBalanceFree, lockedBalances } = props;
  const sortedVirtualChannels: VirtualChannelBalanceProps[] =
    sortToExtremes(lockedBalances);
  const lockedTotal = sortedVirtualChannels.reduce(
    (acc, x) => acc + x.budget,
    BigInt(0)
  );

  const color = ((status: NetworkBalanceProps["status"]): string => {
    switch (status) {
      case "running":
        return colors.cGreen;
      case "unresponsive-peer":
        return colors.cAmber;
      case "under-challenge":
        return colors.cRed;
      default:
        return colors.cOrange;
    }
  })(props.status);

  const total = myBalanceFree + theirBalanceFree + lockedTotal;
  const myTotal =
    myBalanceFree +
    sortedVirtualChannels.reduce(
      (acc, x) =>
        acc +
        (x.budget * BigInt(Math.round(10_000 * x.myPercentage))) / 10_000n,
      BigInt(0)
    );
  let data = [];
  let myBalanceFreePercentage: number;
  let theirBalanceFreePercentage: number;

  const virtualChannelData = sortedVirtualChannels.map((x) => ({
    title: `${prettyPrintWei(x.budget)} Locked, ${prettyPrintWei(
      (x.budget * BigInt(Math.round(10_000 * x.myPercentage))) / 10_000n
    )} Mine`,
    value: percentageOfTotal(x.budget, total),
    color: interpolateColor(colors.cGrey, color, x.myPercentage),
  }));

  if (total > 0) {
    [myBalanceFreePercentage, theirBalanceFreePercentage] = [
      myBalanceFree,
      theirBalanceFree,
    ].map((x) => percentageOfTotal(x, total));

    data = [];

    if (myBalanceFreePercentage > 0) {
      data.push({
        title: `${prettyPrintWei(myBalanceFree)}`,
        value: myBalanceFreePercentage,
        color,
      });
    }

    const firstHalfCutoff = virtualChannelData.length / 2;

    // starting from the "my-balance" on the bottom and moving clockwise:

    // stack the locked channels sorter "high-me" to "low-me" balances to the left
    data.push(...virtualChannelData.slice(0, firstHalfCutoff));

    // then the "capacity" balance on top (their balance)
    if (theirBalanceFreePercentage > 0) {
      data.push({
        title: `Receive capacity: ${prettyPrintWei(theirBalanceFree)}`,
        value: theirBalanceFreePercentage,
        color: colors.cGrey,
      });
    }

    // and then the locked balances sorted "low-me" to "high-me" to the right of that
    data.push(...virtualChannelData.slice(firstHalfCutoff));
  } else {
    // failure case: no received balance information
    data = [{ title: "0", value: 100, color: "red" }];
    myBalanceFreePercentage = 0;
    theirBalanceFreePercentage = 0;
  }

  if (props.asTable) {
    return (
      <table>
        <tbody>
          <tr>
            <td></td>
            <td className="budget-progress-bars">
              <span>Available spend capacity</span>
              <LinearProgressWithLabel
                variant="determinate"
                value={myBalanceFreePercentage ?? 0}
                label={prettyPrintWei(myBalanceFree)}
                className={"bar me"}
              />
              <span>Available receive capacity</span>
              <LinearProgressWithLabel
                variant="determinate"
                value={theirBalanceFreePercentage ?? 0}
                label={prettyPrintWei(theirBalanceFree)}
                className={"bar their"}
              />

              <span>Locked Capacity</span>
              <LinearProgressWithLabel
                variant="determinate"
                value={
                  100 - myBalanceFreePercentage - theirBalanceFreePercentage
                }
                label={prettyPrintWei(lockedTotal)}
                className={"bar locked-me"}
              />
            </td>
          </tr>
        </tbody>
      </table>
    );
  }

  // The pie-chart renders the first data point starting 0 degrees (3 o'clock)
  // and then rotates clockwise. We want the first data point - the user balance,
  // to sit at the bottom (centered at 12 o'clock).
  //
  // we:
  //   - add 90 degrees to the angle (to start at 6 o'clock instead of 3)
  //   - subtract half of the angle spanned by the user balance (to center it)
  const angleOfMyBalance = (myBalanceFreePercentage * 360) / 100;
  const angleOffset = 90 - angleOfMyBalance / 2;

  return (
    <PieChart
      className="budget-pie-chart"
      animate
      lineWidth={18}
      labelStyle={() => ({
        fill: color,
        fontSize: "10px",
        fontFamily: "sans-serif",
      })}
      radius={42}
      data={data}
      label={() => prettyPrintWei(myTotal)}
      labelPosition={0}
      segmentsStyle={() => ({ color: "red" })}
      paddingAngle={data.length > 1 ? 0.5 : 0}
      startAngle={angleOffset}
    />
  );
};

function LinearProgressWithLabel(
  props: LinearProgressProps & { value: number; label: string }
) {
  return (
    <Box display="flex" alignItems="center">
      <Box width="100%" mr={1}>
        <LinearProgress variant="determinate" {...props} />
      </Box>
      <Box minWidth={100} maxWidth={100}>
        <Typography variant="caption" color="textSecondary">
          {props.label}
        </Typography>
      </Box>
    </Box>
  );
}
