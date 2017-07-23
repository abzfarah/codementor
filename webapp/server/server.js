const express = require('express');
const httpProxy = require('http-proxy');
const webpack = require('webpack');
const webpackDevMiddleware = require('webpack-dev-middleware');
const webpackHotMiddleware = require('webpack-hot-middleware');
const config = require('../webpack.config');


const app = express();
const apiProxy = httpProxy.createProxyServer();
const compiler = webpack(config);

// Start a webpack-dev-server
app.use(webpackDevMiddleware(compiler, {
  // server and middleware options
  publicPath: config.output.publicPath,
}));

// Enables HMR
app.use(webpackHotMiddleware(compiler));

// Proxy api requests
app.use('/', (req, res) => {
  req.url = req.baseUrl; // Janky hack...
  apiProxy.web(req, res, {
    target: {
      port: 8065,
      host: 'localhost'
    },
    ws: true,
    secure: false,
    historyApiFallback: true,
    stats: {
      // Config for minimal console.log mess.
      assets: true,
      colors: true,
      version: false,
      hash: true,
      timings: false,
      chunks: true,
      chunkModules: false
    }
  });
});

app.listen(8080, 'localhost');
