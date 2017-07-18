// (C) Copyright 2014-2015 Hewlett Packard Enterprise Development LP

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import CSSClassnames from '../../../utils/CSSClassnames';
import Intl from '../../../utils/Intl';
import Props from '../../../utils/Props';

const CLASS_ROOT = CSSClassnames.CONTROL_ICON;
const COLOR_INDEX = CSSClassnames.COLOR_INDEX;

export default class Icon extends Component {
  render () {
    const { className, colorIndex } = this.props;
    let { a11yTitle, size, responsive } = this.props;
    let { intl } = this.context;

    const classes = classnames(
      CLASS_ROOT,
      `${CLASS_ROOT}-social-google-plus`,
      className,
      {
        [`${CLASS_ROOT}--${size}`]: size,
        [`${CLASS_ROOT}--responsive`]: responsive,
        [`${COLOR_INDEX}-${colorIndex}`]: colorIndex
      }
    );

    a11yTitle = a11yTitle || Intl.getMessage(intl, 'social-google-plus');

    const restProps = Props.omit(this.props, Object.keys(Icon.propTypes));
    return <svg {...restProps} version="1.1" viewBox="0 0 24 24" width="24px" height="24px" role="img" className={classes} aria-label={a11yTitle}><path fill="#4285F4" fillRule="evenodd" d="M0.0202537528,11.2032618 C0.0858884305,7.29693926 3.69605513,3.87603603 7.62298049,4.00635861 C9.5045944,3.91939087 11.2733582,4.73384249 12.7144674,5.87784249 C12.0993693,6.57306829 11.4624794,7.24274571 10.7820061,7.86906829 C9.04956189,6.67810055 6.58579693,6.33771345 4.85335275,7.71345539 C2.37505996,9.41848765 2.26220982,13.4442941 4.64607169,15.2796489 C6.96429889,17.3722941 11.3459973,16.3333264 11.9865191,13.1291973 C10.534514,13.1075199 9.07861756,13.1291973 7.62661245,13.0822296 C7.62298049,12.2205522 7.61934853,11.3588747 7.62298049,10.4974554 C10.0504259,10.4902296 12.4778712,10.4863586 14.9089486,10.5046812 C15.0544863,12.532036 14.7852026,14.6897135 13.5332146,16.3658425 C11.6370729,19.0195199 7.83052098,19.7942296 4.86061667,18.6574554 C1.88007591,17.527907 -0.230870232,14.3962941 0.0202537528,11.2032618 M19.6437253,8.32145539 L21.8091508,8.32145539 C21.8127827,9.04197152 21.8164147,9.76610055 21.8236786,10.4863586 C22.5479949,10.4938425 23.2759431,10.4938425 24,10.5010683 L24,12.6551328 C23.2759431,12.6623586 22.5516269,12.6659715 21.8236786,12.6731973 C21.8164147,13.3970683 21.8127827,14.1175844 21.8091508,14.8381005 L19.6400933,14.8381005 C19.6328294,14.1175844 19.6328294,13.3970683 19.6255655,12.6768102 C18.9012492,12.6695844 18.1735604,12.6623586 17.4492441,12.6551328 L17.4492441,10.5010683 C18.1735604,10.4938425 18.8976172,10.4902296 19.6255655,10.4863586 C19.6291974,9.76248765 19.6364613,9.04197152 19.6437253,8.32145539" stroke="none"/></svg>;
  }
};

Icon.contextTypes = {
  intl: PropTypes.object
};

Icon.defaultProps = {
  responsive: true
};

Icon.displayName = 'SocialGooglePlus';

Icon.icon = true;

Icon.propTypes = {
  a11yTitle: PropTypes.string,
  colorIndex: PropTypes.string,
  size: PropTypes.oneOf(['xsmall', 'small', 'medium', 'large', 'xlarge', 'huge']),
  responsive: PropTypes.bool
};

