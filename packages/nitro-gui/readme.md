# React + Typescript Front-end

## Available Scripts

In the project directory, you can run:

### `yarn storybook`

Runs the storybook instance for the project, for\
previewing components in a variety of configured states.
Open [http://localhost:6006](http://localhost:6006) to view it in the browser.

### `yarn dev`

Runs the app in the development mode with hot reload at [http://localhost:5173/](http://localhost:5173/). The UI will "point at" (i.e. launch RPC requests against) the `VITE_RPC_HOST` env var. This can be set in `.env.development` and should correspond to a Nitro RPC-Server-Enabled node.

### `yarn build`

Builds the app for production to the `dist` folder. By leaving the `VITE_RPC_HOST` env var unset, the UI will "point at" the same host and port that serves the UI itself. This works well when the output is to be embedded into the `go-nitro` binary. This can be done by using the [appropriate build tags](../../readme.md).

## Wireframe

This wireframe gives a rough idea of the layout and components that will be used in the app:

![Wireframe](./wireframe.png)

## Tooling References

This project was bootstrapped with [Vite](https://vitejs.dev/).
