// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


export default class Provider {
    constructor() {
        this.latestPrefix = '';
        this.latestComplete = true;
    }

    handlePretextChanged(suggestionId, pretext) { // eslint-disable-line no-unused-vars
        // NO-OP for inherited classes to override
    }

    startNewRequest(prefix) {
        this.latestPrefix = prefix;
        this.latestComplete = false;
    }

    shouldCancelDispatch(prefix) {
        if (prefix === this.latestPrefix) {
            this.latestComplete = true;
        } else if (this.latestComplete) {
            return true;
        }

        return false;
    }
}
