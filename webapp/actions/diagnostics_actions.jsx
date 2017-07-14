// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


import Client from 'client/web_client.jsx';

export function trackEvent(category, event, properties) {
    Client.trackEvent(category, event, properties);
}
