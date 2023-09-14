import * as React from "react";
import Button from "@mui/material/Button";
import CssBaseline from "@mui/material/CssBaseline";
import FormControlLabel from "@mui/material/FormControlLabel";
import Link from "@mui/material/Link";
import Grid from "@mui/material/Grid";
import { createTheme, ThemeProvider } from "@mui/material/styles";
import { Slider, Stack, Switch } from "@mui/material";
import Box from "@mui/material/Box";
import Stepper from "@mui/material/Stepper";
import Step from "@mui/material/Step";
import StepLabel from "@mui/material/StepLabel";
import StepContent from "@mui/material/StepContent";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import PersonIcon from "@mui/icons-material/Person";
import StorageIcon from "@mui/icons-material/Storage";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function Copyright(props: any) {
  return (
    <Typography
      variant="body2"
      color="text.secondary"
      align="center"
      {...props}
    >
      {"Copyright © "}
      <Link color="inherit" href="https://statechannels.org/">
        statechannels.org
      </Link>{" "}
      {new Date().getFullYear()}
      {"."}
    </Typography>
  );
}

// TODO remove, this demo shouldn't need to reset the theme.
const defaultTheme = createTheme();

export default function SignInSide() {
  return (
    <ThemeProvider theme={defaultTheme}>
      <Grid container component="main" sx={{ height: "100vh" }}>
        <CssBaseline />
        <Grid
          item
          xs={false}
          sm={4}
          md={7}
          sx={{
            backgroundImage:
              "url(https://source.unsplash.com/random?wallpapers)",
            backgroundRepeat: "no-repeat",
            backgroundColor: (t) =>
              t.palette.mode === "light"
                ? t.palette.grey[50]
                : t.palette.grey[900],
            backgroundSize: "cover",
            backgroundPosition: "center",
          }}
        />
        <Grid item xs={12} sm={8} md={5} component={Paper} elevation={6} square>
          <Box
            sx={{
              my: 8,
              mx: 4,
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
            }}
          >
            <Stack spacing={3}>
              <Typography component="h1" variant="h5">
                Filecoin Paid Retrieval Demo
              </Typography>
              <VerticalLinearStepper />
              <Link href="#" variant="body2">
                How does this work?
              </Link>
              <Copyright sx={{ mt: 5 }} />
            </Stack>
          </Box>
        </Grid>
      </Grid>
    </ThemeProvider>
  );
}

function VerticalLinearStepper() {
  const [activeStep, setActiveStep] = React.useState(0);

  const handleNext = () => {
    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };

  return (
    <Box sx={{ maxWidth: 400 }}>
      <Stepper activeStep={activeStep} orientation="vertical">
        <Step key={"Join the Nitro Payment Network"}>
          <StepLabel>{"Join the Nitro Payment Network"}</StepLabel>
          <StepContent>
            <Typography>{`In this demonstration, you will be sharing in a prefunded network account on Calibration Tesnet with all other users.`}</Typography>
            <Box sx={{ mb: 2 }}>
              <div>
                <Button
                  variant="contained"
                  onClick={handleNext}
                  sx={{ mt: 1, mr: 1 }}
                >
                  OK
                </Button>
              </div>
            </Box>
          </StepContent>
        </Step>

        <Step key={"Connect to a Retrieval Provider"}>
          <StepLabel>{"Connect to a Retrieval Provider"}</StepLabel>
          <StepContent>
            <Typography>
              {
                "Create a virtual payment with enough capacity to pay for 10 retrievals."
              }
            </Typography>
            <Box sx={{ mb: 2 }}>
              <div>
                <Button
                  variant="contained"
                  onClick={handleNext}
                  sx={{ mt: 1, mr: 1 }}
                >
                  Create Channel
                </Button>
              </div>
            </Box>
          </StepContent>
        </Step>

        <Step key={"Execute a Paid Retrieval"}>
          <StepLabel>{"Execute a Paid Retrieval"}</StepLabel>
          <StepContent>
            <Stack spacing={5} direction="column">
              <Stack>
                <Typography>
                  {
                    "Create a payment voucher, and attach it to a request for the provider."
                  }
                </Typography>
                <Box sx={{ mb: 2 }}>
                  <div>
                    <Box
                      component="form"
                      noValidate
                      onSubmit={() => {
                        /* TODO */
                      }}
                      sx={{ mt: 1 }}
                    >
                      <Stack direction="row" spacing={2}></Stack>
                      <FormControlLabel
                        control={<Switch value="skippayment" color="primary" />}
                        label="Skip payment"
                      />
                      <FormControlLabel
                        control={
                          <Switch value="usemicropayments" color="primary" />
                        }
                        label="Use micropayments"
                      />
                      <Button
                        variant="contained"
                        onClick={() => {
                          /* TODO */
                        }}
                        sx={{ mt: 1, mr: 1 }}
                      >
                        Pay & Download
                      </Button>
                    </Box>
                  </div>
                </Box>
              </Stack>

              <Stack
                spacing={2}
                direction="row"
                sx={{ mb: 1 }}
                alignItems="center"
              >
                <Typography variant="body2" color="text.secondary">
                  150 FIL
                </Typography>
                <PersonIcon />
                <Slider aria-label="Volume" value={75} valueLabelDisplay="on" />
                <StorageIcon />{" "}
                <Typography variant="body2" color="text.secondary">
                  250 FIL
                </Typography>
              </Stack>
            </Stack>
          </StepContent>
        </Step>
      </Stepper>
    </Box>
  );
}
