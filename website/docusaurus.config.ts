import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Linear CLI',
  tagline: 'A JSON-first CLI for Linear automation',
  favicon: 'img/logo.svg',
  url: 'https://linear-cli.enolalab.com',
  baseUrl: '/',
  organizationName: 'enolalab',
  projectName: 'linear-cli',
  onBrokenLinks: 'throw',
  markdown: {hooks: {onBrokenMarkdownLinks: 'throw'}},
  i18n: {defaultLocale: 'en', locales: ['en']},
  presets: [
    [
      'classic',
      {
        docs: {
          routeBasePath: '/',
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/enolalab/linear-cli/edit/main/website/',
        },
        blog: false,
        theme: {customCss: './src/css/custom.css'},
      } satisfies Preset.Options,
    ],
  ],
  themeConfig: {
    image: 'img/logo.svg',
    navbar: {
      title: 'Linear CLI',
      logo: {alt: 'Linear CLI', src: 'img/logo.svg'},
      items: [
        {type: 'docSidebar', sidebarId: 'docs', position: 'left', label: 'Documentation'},
        {href: 'https://enolalab.com', label: 'Enolalab', position: 'right'},
        {href: 'https://github.com/enolalab/linear-cli', label: 'GitHub', position: 'right'},
        {href: 'https://github.com/enolalab/linear-cli/releases', label: 'Releases', position: 'right'},
        {href: 'https://github.com/enolalab/linear-cli', label: 'Catalogue', position: 'right'},
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {title: 'Documentation', items: [{label: 'Getting started', to: '/'}, {label: 'Command reference', to: '/commands'}]},
        {title: 'Project', items: [{label: 'GitHub', href: 'https://github.com/enolalab/linear-cli'}, {label: 'Releases', href: 'https://github.com/enolalab/linear-cli/releases'}]},
        {title: 'Enolalab', items: [{label: 'Enolalab', href: 'https://enolalab.com'}, {label: 'Catalogue', href: 'https://github.com/enolalab/linear-cli'}]},
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Enolalab.`,
    },
    prism: {additionalLanguages: ['bash', 'json', 'yaml']},
  } satisfies Preset.ThemeConfig,
};

export default config;
