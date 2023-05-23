import type { Preview } from "@storybook/react";
import { initialize, mswDecorator } from "msw-storybook-addon";
import { WebSocket } from "mock-socket";

// Override the global WebSocket object with the mock one
global.WebSocket = WebSocket;

// Initialize MSW
initialize({
  onUnhandledRequest: "bypass",
});

// Provide the MSW addon decorator globally
export const decorators = [mswDecorator];

const preview: Preview = {
  parameters: {
    actions: { argTypesRegex: "^on[A-Z].*" },
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/,
      },
    },
  },
};

export default preview;
