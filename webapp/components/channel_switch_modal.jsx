// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


import SuggestionList from './suggestion/suggestion_list.jsx';
import SuggestionBox from './suggestion/suggestion_box.jsx';
import SwitchChannelProvider from './suggestion/switch_channel_provider.jsx';

import {FormattedMessage} from 'react-intl';
import {Modal} from 'react-bootstrap';

import {goToChannel, openDirectChannelToUser} from 'actions/channel_actions.jsx';

import ChannelStore from 'stores/channel_store.jsx';
import UserStore from 'stores/user_store.jsx';

import Constants from 'utils/constants.jsx';
import * as Utils from 'utils/utils.jsx';

import React from 'react';
import $ from 'jquery';

export default class SwitchChannelModal extends React.Component {
    constructor() {
        super();

        this.onChange = this.onChange.bind(this);
        this.onItemSelected = this.onItemSelected.bind(this);
        this.onShow = this.onShow.bind(this);
        this.onHide = this.onHide.bind(this);
        this.onExited = this.onExited.bind(this);
        this.handleKeyDown = this.handleKeyDown.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.switchToChannel = this.switchToChannel.bind(this);

        this.suggestionProviders = [new SwitchChannelProvider()];

        this.state = {
            text: '',
            error: ''
        };
    }

    componentDidUpdate(prevProps) {
        if (this.props.show && !prevProps.show) {
            const textbox = this.refs.search.getTextbox();
            textbox.focus();
            Utils.placeCaretAtEnd(textbox);
        }
    }

    onShow() {
        this.setState({
            text: '',
            error: ''
        });
    }

    onHide() {
        this.setState({
            text: '',
            error: ''
        });
        this.props.onHide();
    }

    onExited() {
        this.selected = null;
        setTimeout(() => {
            $('#post_textbox').get(0).focus();
        });
    }

    onChange(e) {
        this.setState({text: e.target.value});
        this.selected = null;
    }

    onItemSelected(item) {
        this.selected = item;
    }

    handleKeyDown(e) {
        this.setState({
            error: ''
        });
        if (e.keyCode === Constants.KeyCodes.ENTER) {
            this.handleSubmit();
        }
    }

    handleSubmit() {
        let channel = null;

        if (!this.selected) {
            if (this.state.text !== '') {
                this.setState({
                    error: Utils.localizeMessage('channel_switch_modal.not_found', 'No matches found.')
                });
            }
            return;
        }

        if (this.selected.type === Constants.DM_CHANNEL) {
            const user = UserStore.getProfileByUsername(this.selected.name);

            if (user) {
                openDirectChannelToUser(
                    user.id,
                    (ch) => {
                        channel = ch;
                        this.switchToChannel(channel);
                    },
                    () => {
                        channel = null;
                        this.switchToChannel(channel);
                    }
                );
            }
        } else {
            channel = ChannelStore.get(this.selected.id);
            this.switchToChannel(channel);
        }
    }

    switchToChannel(channel) {
        if (channel !== null) {
            goToChannel(channel);
            this.onHide();
        } else if (this.state.text !== '') {
            this.setState({
                error: Utils.localizeMessage('channel_switch_modal.failed_to_open', 'Failed to open channel.')
            });
        }
    }

    render() {
        const message = this.state.error;
        return (
            <Modal
                dialogClassName='channel-switch-modal'
                ref='modal'
                show={this.props.show}
                onHide={this.onHide}
                onExited={this.onExited}
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title>
                        <span>
                            <FormattedMessage
                                id='channel_switch_modal.title'
                                defaultMessage='Switch Channels'
                            />
                        </span>
                    </Modal.Title>
                </Modal.Header>

                <Modal.Body>
                    <div className='modal__hint'>
                        <FormattedMessage
                            id='channel_switch_modal.help'
                            defaultMessage='Type channel name. Use ↑↓ to browse, TAB to select, ↵ to confirm, ESC to dismiss'
                        />
                    </div>
                    <SuggestionBox
                        ref='search'
                        className='form-control focused'
                        type='input'
                        onChange={this.onChange}
                        value={this.state.text}
                        onKeyDown={this.handleKeyDown}
                        onItemSelected={this.onItemSelected}
                        listComponent={SuggestionList}
                        maxLength='64'
                        providers={this.suggestionProviders}
                        listStyle='bottom'
                    />
                </Modal.Body>
                <Modal.Footer>
                    <div className='modal__error'>
                        {message}
                    </div>
                    <button
                        type='button'
                        className='btn btn-default'
                        onClick={this.onHide}
                    >
                        <FormattedMessage
                            id='edit_channel_header_modal.cancel'
                            defaultMessage='Cancel'
                        />
                    </button>
                    <button
                        type='button'
                        className='btn btn-primary'
                        onClick={this.handleSubmit}
                    >
                        <FormattedMessage
                            id='channel_switch_modal.submit'
                            defaultMessage='Switch'
                        />
                    </button>
                </Modal.Footer>
            </Modal>
        );
    }
}

SwitchChannelModal.propTypes = {
    show: React.PropTypes.bool.isRequired,
    onHide: React.PropTypes.func.isRequired
};

