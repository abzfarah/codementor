var webpack = require('webpack')
const WebpackDevServer = require('webpack-dev-server')
const project = require('../project.config')
const path = require('path')
const config = require('../webpack.config')
const compiler = webpack(config)

const server = new WebpackDevServer(compiler, {
  publicPath: config.output.publicPath,
  contentBase : path.resolve(project.basePath, project.outDir),
  hot: true,
  proxy: {
    '/': {
      target: 'http://localhost:8065/',
      secure: false
    },
  },
  historyApiFallback:  true,
  // It suppress error shown in console, so it has to be set to false.
  quiet: false,
  // It suppress everything except error, so it has to be set to false as well
  // to see success build.
  noInfo: false,
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
})

server.listen(8080, 'localhost', function () {})

/*
new WebpackDevServer(webpack(config), {
  publicPath: config.output.publicPath,
  contentBase : path.resolve(project.basePath, project.srcDir),
  hot: false,
  historyApiFallback: false,
  // It suppress error shown in console, so it has to be set to false.
  quiet: false,
  // It suppress everything except error, so it has to be set to false as well
  // to see success build.
  noInfo: false,
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
}).listen(3000, 'localhost', function (err) {
  if (err) {
    console.log(err)
  }

  console.log('Listening at localhost:3000')
}) */
