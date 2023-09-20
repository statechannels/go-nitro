import React from "react";
import Button, { ButtonProps } from "@mui/material/Button";
import Box from "@mui/material/Box";
import { makeStyles } from "@mui/styles";

interface ProgressButtonProps extends Omit<ButtonProps, "onClick"> {
  value: number;
  onClick: () => void;
}

const useStyles = makeStyles({
  progressStyle: (props: ProgressButtonProps) => ({
    position: "absolute",
    top: 0,
    left: 0,
    bottom: 0,
    backgroundColor: "rgba(0, 0, 255, 0.5)", // Adjust color as needed
    width: `${props.value}%`,
    // transition: "width 0.3s ease",
  }),
  customDisabled: {
    cursor: "not-allowed",
    pointerEvents: "none",
    opacity: 0.5,
  },
});

const ProgressButton: React.FC<ProgressButtonProps> = ({
  value,
  onClick,
  children,
  ...props
}) => {
  const isDisabled = value > 0;
  const classes = useStyles({
    value,
    onClick,
    children,
    ...props,
  });

  return (
    <Box position="relative">
      <Button
        variant="contained"
        fullWidth
        onClick={onClick}
        className={isDisabled ? classes.customDisabled : ""}
        {...props}
      >
        {children}
      </Button>
      <Box className={classes.progressStyle}></Box>
    </Box>
  );
};

export default ProgressButton;
