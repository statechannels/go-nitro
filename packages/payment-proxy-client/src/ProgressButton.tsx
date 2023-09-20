import React from "react";
import Button, { ButtonProps } from "@mui/material/Button";
import { makeStyles } from "@mui/styles";

const useStyles = makeStyles({
  customButton: (props: { fillPercentage: number }) => ({
    background: `linear-gradient(92deg, white 0%, white ${props.fillPercentage}%, transparent ${props.fillPercentage}%, transparent 100%)`,
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
    <Button className={classes.customButton} {...props}>
      {props.label}
    </Button>
  );
};

export default CustomButton;
