// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Sedge documentation',
  tagline: 'Sedge was created for the setup and deployment of an Ethereum node with ease',
  url: 'https://docs.sedge.nethermind.io/',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',

  // GitHub pages deployment config.
  organizationName: 'NethermindEth',
  projectName: 'sedge',
  trailingSlash: false,

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl:
              'https://github.com/NethermindEth/sedge/tree/main/docs/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
  /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
      ({
        navbar: {
          title: '',
          logo: {
            alt: 'Sedge Logo',
            src: 'img/Sedge_Horizontal_Dark.svg', // Dark theme logo
            srcDark: 'img/Sedge_Horizontal_Light.svg', // Light theme logo
          },
          items: [
            {
              type: 'doc',
              docId: 'intro',
              position: 'left',
              label: 'Documentation',
            },
            {
              type: 'dropdown',
              label: 'Networks',
              position: 'left',
              items: [
                {
                  label: 'Ethereum',
                  to: '/docs/networks/mainnet',
                },
                {
                  label: 'Holesky',
                  to: '/docs/networks/holesky',
                },
                {
                  label: 'Sepolia',
                  to: '/docs/networks/sepolia',
                },
                {
                  label: 'Gnosis',
                  to: '/docs/networks/gnosis',
                },
                {
                  label: 'Chiado',
                  to: '/docs/networks/chiado',
                },
                {
                  label: 'Optimism',
                  to: '/docs/quickstart/running-optimism-node',
                },
                {
                  label: 'Base',
                  to: '/docs/quickstart/running-optimism-node#base-support',
                },
                {
                  label: 'Lido CSM',
                  to: '/docs/quickstart/staking-with-lido',
                },
              ]
            },
            {
              href: 'https://github.com/NethermindEth/sedge',
              label: 'GitHub',
              position: 'right',
            },
          ],
        },
        footer: {
          // style: 'dark',
          links: [
            {
              title: 'Networks',
              items: [
                {
                  label: 'Mainnet',
                  to: '/docs/networks/mainnet',
                },
                {
                  label: 'Sepolia',
                  to: '/docs/networks/sepolia',
                },
                {
                  label: 'Holesky',
                  to: '/docs/networks/holesky',
                },
                {
                  label: 'Gnosis',
                  to: '/docs/networks/gnosis',
                },
                {
                  label: 'Chiado',
                  to: '/docs/networks/chiado',
                },
                {
                  label: 'Optimism',
                  to: '/docs/quickstart/running-optimism-node',
                },
                {
                  label: 'Base',
                  to: '/docs/quickstart/running-optimism-node#base-support',
                },
                {
                  label: 'Lido CSM',
                  to: '/docs/quickstart/staking-with-lido',
                },
              ],
            },
            {
              title: 'Community',
              items: [
                {
                  label: 'Discord',
                  href: 'https://discord.com/invite/PaCMRFdvWT',
                },
                {
                  label: 'Twitter',
                  href: 'https://twitter.com/nethermindeth',
                },
                {
                  label: 'GitHub',
                  href: 'https://github.com/NethermindEth/sedge',
                },
              ],
            },
            {
              title: 'GitHub Repositories',
              items: [
                {
                  label: 'Sedge',
                  href: 'https://github.com/NethermindEth/sedge',
                },
                {
                  label: 'Nethermind',
                  href: 'https://github.com/NethermindEth/nethermind',
                },
              ],
            },
            {
              title: 'More',
              items: [
                {
                  label: 'Nethermind',
                  href: 'https://nethermind.io/',
                },
                {
                  label: 'Documentation',
                  to: '/docs/intro',
                },
              ],
            },
          ],
          copyright: `Copyright © ${new Date().getFullYear()} Nethermind. Built with Docusaurus.`,
        },
        prism: {
          theme: lightCodeTheme,
          darkTheme: darkCodeTheme,
        },
        algolia: {
          appId: 'HR7BYPH22J',
          apiKey: '9fd1f2f1dfdc41503fea4f38c11fe89f',
          indexName: 'sedge-nethermind',
          contextualSearch: true,
          externalUrlRegex: 'external\\.com|domain\\.com',
          searchParameters: {},
          searchPagePath: 'false',
        },
      }),
};

module.exports = config;
