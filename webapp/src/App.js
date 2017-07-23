// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.

import React, { Component } from 'react';
import { App, Box, Header, Button } from './components';
import  GrommetIcon  from './components/icons/base/BrandGrommetOutline';

export default class SampleApp extends Component {


  render () {
    return (
      <App centered={false} >
        <Header justify="center" colorIndex="neutral-4">
          <Box
            size={{ width: { max: 'xxlarge' } }}
            direction="row"
            responsive={false}
            justify="between"
            align="center"
            pad={{ horizontal: 'medium' }}
            flex="grow"
          >fghfghddsssdddddddss
            <GrommetIcon colorIndex="brand" size="small" />
            <Button
              label="Sign in"
              primary
              href="#"
            />
          </Box>
        </Header>
      </App>
    );
  }
}
