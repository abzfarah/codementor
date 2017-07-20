// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.

import Graph from './Graph';

export default class Bar extends Graph {};

Bar.defaultProps = {
  ...Graph.defaultProps,
  type: 'bar'
};

Bar.displayName = 'Bar';
