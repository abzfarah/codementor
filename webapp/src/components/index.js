// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.
import Grommet from '../index-commonjs';

const BaseIcons = {};

Object.keys(Grommet.Icons.Base).forEach((icon) => {
  BaseIcons[icon.replace('Icon', '')] = Grommet.Icons.Base[icon];
});

export const Icons = Object.assign({}, Grommet.Icons, { Base: BaseIcons });

export { default as Accordion } from './Accordion';
export { default as AccordionPanel } from './AccordionPanel';
export { default as Anchor } from './Anchor';
export { default as Animate } from './Animate';
export { default as App } from './App';
export { default as Article } from './Article';
export { default as Box } from './Box';
export { default as Button } from './Button';
export { default as Card } from './Card';
export { default as Carousel } from './Carousel';
export * from './chart/index';
export { default as CheckBox } from './CheckBox';
export { default as Columns } from './Columns';
export { default as DateTime } from './DateTime';
export { default as Distribution } from './Distribution';
export { default as Footer } from './Footer';
export { default as Form } from './Form';
export { default as FormattedMessage } from './FormattedMessage';
export { default as FormField } from './FormField';
export { default as FormFields } from './FormFields';
export { default as Grommet } from './Grommet';
export { default as Header } from './Header';
export { default as Heading } from './Heading';
export { default as Headline } from './Headline';
export { default as Hero } from './Hero';
export { default as Image } from './Image';
export { default as Label } from './Label';
export { default as Layer } from './Layer';
export { default as Legend } from './Legend';
export { default as List } from './List';
export { default as ListItem } from './ListItem';
export { default as LoginForm } from './LoginForm';
export { default as Map } from './Map';
export { default as Markdown } from './Markdown';
export { default as Menu } from './Menu';
export { default as Meter } from './Meter';
export { default as Notification } from './Notification';
export { default as NumberInput } from './NumberInput';
export { default as Object } from './Object';
export { default as Paragraph } from './Paragraph';
export { default as Quote } from './Quote';
export { default as RadioButton } from './RadioButton';
export { default as Search } from './Search';
export { default as SearchInput } from './SearchInput';
export { default as Section } from './Section';
export { default as Select } from './Select';
export { default as Sidebar } from './Sidebar';
export { default as SkipLinkAnchor } from './SkipLinkAnchor';
export { default as SkipLinks } from './SkipLinks';
export { default as SocialShare } from './SocialShare';
export { default as Split } from './Split';
export { default as SunBurst } from './SunBurst';
export { default as SVGIcon } from './SVGIcon';
export { default as Tab } from './Tab';
export { default as Table } from './Table';
export { default as TableHeader } from './TableHeader';
export { default as TableRow } from './TableRow';
export { default as Tabs } from './Tabs';
export { default as TBD } from './TBD';
export { default as TextInput } from './TextInput';
export { default as Tile } from './Tile';
export { default as Tiles } from './Tiles';
export { default as Timestamp } from './Timestamp';
export { default as Tip } from './Tip';
export { default as Title } from './Title';
export { default as Toast } from './Toast';
export { default as Topology } from './Topology';
export { default as Value } from './Value';
export { default as Video } from './Video';
export { default as WorldMap } from './WorldMap';
export * from './icons/index';
export { default as Cookies } from '../utils/Cookies';
export { default as DOM } from '../utils/DOM';
export { default as KeyboardAccelerators } from '../utils/KeyboardAccelerators';
export { default as Locale } from '../utils/Locale';
export { default as Responsive } from '../utils/Responsive';
export { default as Rest } from '../utils/Rest';
export { default as Validator } from '../utils/Validator';

export default { ...Grommet };
