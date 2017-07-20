import React, { Component } from 'react'

class Home extends Component {
  login () {
    this.props.auth.login()
  }
  render () {
    const { isAuthenticated } = this.props.auth
    return (
      <div className='container'>
        {
          isAuthenticated() && (
            <h4>
              You are loggdded in!
            </h4>
            )
        }
        {
          !isAuthenticated() && (
            <h4>
              You addre not loggdrddded in! Pleadsse{' '}
              <a
                style={{ cursor: 'pointer' }}
                onClick={this.login.bind(this)}
              >
                Log In
              </a>
              {' '}to continue.
            </h4>
            )
        }
      </div>
    )
  }
}

export default Home

