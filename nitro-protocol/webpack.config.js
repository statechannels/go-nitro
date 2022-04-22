/* eslint-disable @typescript-eslint/no-var-requires */
import {resolve} from 'path';

const config = {
  entry: './lib/src/index.js',
  module: {},
  target: 'web',
  mode: 'production',
  output: {
    filename: 'nitro-protocol.min.js',
    libraryTarget: 'commonjs',
    // eslint-disable-next-line no-undef
    path: resolve(__dirname, 'dist'),
  },
  node: {
    fs: 'empty',
    child_process: 'empty',
  },
};

export default [config];
