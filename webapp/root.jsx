import React from 'react'
import ReactDOM from 'react-dom'
import createStore from './src/store/createStore'
import App from './src/components/App'
import './sass/index.scss'
import $ from 'jquery';
// Store Initialization
// ------------------------------------
const store = createStore(window.__INITIAL_STATE__)

// Render Setup
// ------------------------------------

const MOUNT_NODE = document.getElementsByTagName('body')

console.log(MOUNT_NODE)
ReactDOM.render(
  <App />,
  MOUNT_NODE
)



