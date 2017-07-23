import React from 'react';
import ReactDOM from 'react-dom';
import { makeMainRoutes } from './routes';
import './styles/index.scss';

const routes = makeMainRoutes()

// Render Setup
// ------------------------------------
var MOUNT_NODE = document.getElementById('root')

let render = () => {
  ReactDOM.render(
    routes,
    MOUNT_NODE
  );
}

// Development Tools
// ------------------------------------

if (module.hot) {
  const renderApp = render;

  render = () => {
    renderApp();
  }

  // Setup hot module replacement

  module.hot.accept('./components/App', () => {
    console.log('fark you HMR');
    render();
  });
}

// Let's Go!
// ------------------------------------
render()

