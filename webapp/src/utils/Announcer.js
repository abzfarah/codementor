// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.

import CSSClassnames from './CSSClassnames';

const CLASS_ROOT = CSSClassnames.APP;

function clearAnnouncer() {
  const announcer = document.querySelector(`.${CLASS_ROOT}__announcer`);
  if(announcer) {
    announcer.innerHTML = '';
  }
};

export function announcePageLoaded (title) {
  announce(`${title} page was loaded`);
}

export function announce (message, mode = 'assertive') {
  const announcer = document.querySelector(`.${CLASS_ROOT}__announcer`);
  if(announcer) {
    announcer.setAttribute('aria-live', 'off');
    announcer.innerHTML = message;
    setTimeout(clearAnnouncer, 500);
    announcer.setAttribute('aria-live', mode);
  }
}

export default { announce, announcePageLoaded };
