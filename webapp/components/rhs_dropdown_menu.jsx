// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


import {Dropdown} from 'react-bootstrap';
import React from 'react';

export default class RhsDropdownMenu extends Dropdown.Menu {
    constructor(props) { //eslint-disable-line no-useless-constructor
        super(props);
    }

    render() {
        return (
            <div
                className='dropdown-menu__content'
                onClick={this.props.onClose}
            >
                {super.render()}
            </div>
        );
    }
}
