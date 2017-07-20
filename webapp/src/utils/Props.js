// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.

export default {
  pick (props, fields) {
    const has = (p) => props.hasOwnProperty(p);
    const obj = {};
    (fields || []).forEach((field) => {
      if (has(field))
        obj[field] = props[field];
    });
    return obj;
  },
  omit (props, fields) {
    const obj = {};
    Object.keys(props).forEach((p) => {
      if ((fields || []).indexOf(p) === -1) {
        obj[p] = props[p];
      }
    });
    return obj;
  }
};
