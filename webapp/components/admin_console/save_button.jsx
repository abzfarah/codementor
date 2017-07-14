// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


import React from 'react';

import {FormattedMessage} from 'react-intl';

export default class SaveButton extends React.Component {
    static get propTypes() {
        return {
            saving: React.PropTypes.bool.isRequired,
            disabled: React.PropTypes.bool
        };
    }

    static get defaultProps() {
        return {
            disabled: false
        };
    }

    render() {
        const {saving, disabled, ...props} = this.props; // eslint-disable-line no-use-before-define

        let contents;
        if (saving) {
            contents = (
                <span>
                    <span className='icon fa fa-refresh icon--rotate'/>
                    <FormattedMessage
                        id='admin.saving'
                        defaultMessage='Saving Config...'
                    />
                </span>
            );
        } else {
            contents = (
                <FormattedMessage
                    id='admin.save'
                    defaultMessage='Save'
                />
            );
        }

        let className = 'save-button btn';
        if (!disabled) {
            className += ' btn-primary';
        }

        return (
            <button
                type='submit'
                className={className}
                disabled={disabled}
                {...props}
            >
                {contents}
            </button>
        );
    }
}
