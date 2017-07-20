// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.

import Graph from './Graph';

export default class Line extends Graph {};

Line.defaultProps = {
  ...Graph.defaultProps,
  type: 'line'
};

Line.displayName = 'Line';
