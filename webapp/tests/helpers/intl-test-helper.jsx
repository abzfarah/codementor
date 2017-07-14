// Copyright (c) 2017 Mattermost, Inc. All Rights Reserved.


import {mount, shallow} from 'enzyme';
import React from 'react';
import {IntlProvider, intlShape} from 'react-intl';

const intlProvider = new IntlProvider({locale: 'en'}, {});
const {intl} = intlProvider.getChildContext();

export function shallowWithIntl(node, {context} = {}) {
    return shallow(React.cloneElement(node, {intl}), {
        context: Object.assign({}, context, {intl})
    });
}

export function mountWithIntl(node, {context, childContextTypes} = {}) {
    return mount(React.cloneElement(node, {intl}), {
        context: Object.assign({}, context, {intl}),
        childContextTypes: Object.assign({}, {intl: intlShape}, childContextTypes)
    });
}
