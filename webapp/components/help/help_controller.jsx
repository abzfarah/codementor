// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


import React from 'react';
import ReactDOM from 'react-dom';

export default class HelpController extends React.Component {
    static get propTypes() {
        return {
            children: React.PropTypes.node.isRequired
        };
    }

    componentWillUpdate() {
        ReactDOM.findDOMNode(this).scrollIntoView();
    }

    render() {
        return (
            <div className='help'>
                <div className='container col-sm-10 col-sm-offset-1'>
                    {this.props.children}
                </div>
            </div>
        );
    }
}
