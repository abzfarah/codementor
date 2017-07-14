// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


import AppDispatcher from 'dispatcher/app_dispatcher.jsx';
import {WebrtcActionTypes} from 'utils/constants.jsx';

import Client from 'client/web_client.jsx';

export function initWebrtc(userId, isCalling) {
    AppDispatcher.handleServerAction({
        type: WebrtcActionTypes.INITIALIZE,
        user_id: userId,
        is_calling: isCalling
    });
}

export function handle(message) {
    AppDispatcher.handleServerAction({
        type: message.action,
        message
    });
}

export function webrtcToken(success, error) {
    Client.webrtcToken(
        (data) => {
            if (success) {
                success(data);
            }
        },
        () => {
            if (error) {
                error();
            }
        });
}
