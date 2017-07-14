// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


import * as AsyncClient from 'utils/async_client.jsx';
import Client from 'client/web_client.jsx';

export function uploadFile(file, name, channelId, clientId, success, error) {
    Client.uploadFile(
        file,
        name,
        channelId,
        clientId,
        (data) => {
            if (success) {
                success(data);
            }
        },
        (err) => {
            AsyncClient.dispatchError(err, 'uploadFile');

            if (error) {
                error(err);
            }
        }
    );
}
