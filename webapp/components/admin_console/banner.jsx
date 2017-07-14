// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


import React from 'react';
import {FormattedMessage} from 'react-intl';

export default class Banner extends React.Component {
    render() {
        let title = (
            <FormattedMessage
                id='admin.banner.heading'
                defaultMessage='Note:'
            />
        );

        if (this.props.title) {
            title = this.props.title;
        }

        return (
            <div className='banner'>
                <div className='banner__content'>
                    <h4 className='banner__heading'>
                        {title}
                    </h4>
                    <p>
                        {this.props.description}
                    </p>
                </div>
            </div>
        );
    }
}

Banner.defaultProps = {
};
Banner.propTypes = {
    title: React.PropTypes.node,
    description: React.PropTypes.node.isRequired
};
