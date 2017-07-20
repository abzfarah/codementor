// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.

import React, { Component } from 'react'
import { Navbar, Button } from 'react-bootstrap'

import App from './components/App'

export default class SampleApp extends Component {
  goTo (route) {
    this.props.history.replace(`/${route}`)
  }

  login () {
    this.props.auth.login()
  }

  logout () {
    this.props.auth.logout()
  }
  render () {
    const { isAuthenticated } = this.props.auth
    return (
      <App centered={false}>
        <Navbar fluid>
          <Navbar.Header>
            <Navbar.Brand>
              <a href='#'>Auth0 - React</a>
            </Navbar.Brand>
            <Button
              bsStyle='primary'
              className='btn-margin'
              onClick={this.goTo.bind(this, 'home')}
            >
              Homesddd
            </Button>
            {
              !isAuthenticated() && (
                <Button
                  bsStyle='primary'
                  className='btn-margin'
                  onClick={this.login.bind(this)}
                >
                  Log ssddIn
                </Button>
              )
            }
            {
              isAuthenticated() && (
                <Button
                  bsStyle='primary'
                  className='btn-margin'
                  onClick={this.logout.bind(this)}
                >
                  Log Out
                </Button>
              )
            }
          </Navbar.Header>
        </Navbar>
      </App>
    )
  }
}
