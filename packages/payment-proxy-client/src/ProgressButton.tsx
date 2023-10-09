import React from "react";
import { Button, ButtonProps } from "@mui/material";
import { styled } from "@mui/system";

const CustomButton = styled(Button)({
  background: `linear-gradient(90deg, var(--primary-color) 0%, var(--primary-color) var(--fill-percentage), transparent var(--fill-percentage), transparent 100%)`,
}) as React.FC<ButtonProps>;

export default CustomButton;
