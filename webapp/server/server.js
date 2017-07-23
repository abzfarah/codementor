const express = require('express');
const httpProxy = require('http-proxy');
const webpack = require('webpack');
const webpackDevMiddleware = require('webpack-dev-middleware');
const webpackHotMiddleware = require('webpack-hot-middleware');
const config = require('../webpack.config');
const path = require('path')
const project = require('../project.config')


const app = express();
const apiProxy = httpProxy.createProxyServer();
const compiler = webpack(config);

// Start a webpack-dev-server
app.use(webpackDevMiddleware(compiler, {
  // server and middleware options
  publicPath: config.output.publicPath,
  contentBase: path.resolve(project.basePath, project.outDir),
  stats: {
    colors: true,
  }
}));

// Enables HMR
app.use(webpackHotMiddleware(compiler));

// Proxy api requests
app.use('*', (req, res) => {
  req.url = req.baseUrl; // Janky hack...
  apiProxy.web(req, res, {
    target: {
      port: 8065,
      host: 'localhost'
    },
    ws: true,
    secure: false
  });
});

app.listen(8080, 'localhost');
