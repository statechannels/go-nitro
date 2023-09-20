import React from "react";
import Button, { ButtonProps } from "@mui/material/Button";
import { makeStyles } from "@mui/styles";

const useStyles = makeStyles({
  customButton: (props: { fillPercentage: number }) => ({
    background: `linear-gradient(92deg, white 0%, white ${props.fillPercentage}%, transparent ${props.fillPercentage}%, transparent 100%)`,
    transition: "background 0.3s ease-in-out",
    color: props.fillPercentage > 0 ? "transparent" : undefined,
  }),
});

interface CustomButtonProps extends ButtonProps {
  fillPercentage: number;
  label: string;
}

const CustomButton: React.FC<CustomButtonProps> = (
  props: CustomButtonProps
) => {
  const classes = useStyles({ fillPercentage: props.fillPercentage });
  return (
    <Button
      className={classes.customButton}
      {...props}
      disabled={props.fillPercentage > 0}
    >
      {props.label}
    </Button>
  );
};

export default CustomButton;
