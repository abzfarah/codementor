// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.

import Graph from './Graph';

export default class Area extends Graph {};

Area.defaultProps = {
  ...Graph.defaultProps,
  type: 'area'
};

Area.displayName = 'Area';
