import React from "react";
import { Button, ButtonProps } from "@mui/material";
import { styled } from "@mui/system";

interface CustomButtonProps extends ButtonProps {
  fillPercentage: number;
}

const CustomButton = styled(Button)({
  background: `linear-gradient(90deg, transparent 0%, transparent var(--fill-percentage), var(--primary-color) var(--fill-percentage), var(--primary-color) 100%)`,
  transition: "background 0.3s ease-in-out",
}) as React.FC<CustomButtonProps>;

export default CustomButton;
