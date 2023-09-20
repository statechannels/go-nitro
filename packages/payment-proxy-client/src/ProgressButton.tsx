import React from "react";
import Button, { ButtonProps } from "@mui/material/Button";
import { styled } from "@mui/material/styles";

interface ProgressButtonProps extends ButtonProps {
  value: number;
  children: React.ReactNode;
}

const StyledButton = styled(Button)<{ progress: number }>(
  ({ theme, progress }) => ({
    position: "relative",
    "&:before": {
      content: '""',
      display: progress < 100 ? "block" : "none",
      position: "absolute",
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
      backgroundColor: theme.palette.primary.main,
      opacity: 0.4,
      zIndex: -1,
      width: `${progress}%`,
      transition: "width 0.3s ease",
    },
  })
);

const ProgressButton: React.FC<ProgressButtonProps> = ({
  value,
  onClick,
  children,
  ...props
}) => {
  const isDisabled = value != 0;

  return (
    <StyledButton
      variant="contained"
      disabled={isDisabled}
      onClick={onClick}
      progress={value}
      {...props}
    >
      {children}
    </StyledButton>
  );
};

export default ProgressButton;
