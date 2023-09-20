import React from "react";
import { Button, ButtonProps } from "@mui/material";
import { styled } from "@mui/system";

interface CustomButtonProps extends ButtonProps {
  fillPercentage: number;
}

const CustomButton = styled(Button)({
  background: `linear-gradient(90deg, var(--primary-color) 0%, var(--primary-color) var(--fill-percentage), transparent var(--fill-percentage), transparent 100%)`,
  transition: "background 0.3s ease-in-out",
}) as React.FC<CustomButtonProps>;

export default CustomButton;
