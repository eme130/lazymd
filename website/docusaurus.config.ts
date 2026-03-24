import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'LazyMD — The Compiler for Mortals',
  tagline: 'An LLVM-inspired markdown compiler for mortals. Plain text in, anything out. Runs everywhere — terminal, web, native apps, any device. Written in Go.',
  favicon: 'img/favicon.ico',

  url: 'https://lazymd.com',
  baseUrl: '/',
  trailingSlash: false,

  organizationName: 'user',
  projectName: 'LazyMD',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  headTags: [
    {
      tagName: 'meta',
      attributes: {name: 'keywords', content: 'markdown compiler, compiler for mortals, LLVM architecture, plain text compiler, markdown editor, terminal editor, vim editor, go editor, MCP server, LazyMD, lazymd, multi-platform compiler, collaboration'},
    },
    {
      tagName: 'meta',
      attributes: {name: 'author', content: 'LazyMD'},
    },
    {
      tagName: 'link',
      attributes: {rel: 'canonical', href: 'https://lazymd.com'},
    },
    {
      tagName: 'script',
      attributes: {type: 'application/ld+json'},
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'SoftwareApplication',
        name: 'LazyMD',
        url: 'https://lazymd.com',
        description: 'An LLVM-inspired markdown compiler for mortals. Plain text in, anything out. Runs everywhere — terminal, web, native apps, any device. Written in Go.',
        applicationCategory: 'BusinessApplication',
        operatingSystem: 'Linux, macOS',
        offers: {
          '@type': 'Offer',
          price: '0',
          priceCurrency: 'USD',
        },
        license: 'https://opensource.org/licenses/MIT',
        programmingLanguage: 'Go',
        codeRepository: 'https://github.com/EME130/lazymd',
      }),
    },
    {
      tagName: 'script',
      attributes: {type: 'application/ld+json'},
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'WebSite',
        name: 'LazyMD',
        url: 'https://lazymd.com',
        potentialAction: {
          '@type': 'SearchAction',
          target: 'https://lazymd.com/search?q={search_term_string}',
          'query-input': 'required name=search_term_string',
        },
      }),
    },
  ],

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/EME130/lazymd/tree/main/website/',
        },
        blog: {
          showReadingTime: true,
          feedOptions: {
            type: ['rss', 'atom'],
            xslt: true,
          },
          editUrl: 'https://github.com/EME130/lazymd/tree/main/website/',
        },
        theme: {
          customCss: './src/css/custom.css',
        },
        sitemap: {
          lastmod: 'date',
          changefreq: 'weekly',
          priority: 0.5,
          filename: 'sitemap.xml',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    metadata: [
      {name: 'description', content: 'LazyMD is an LLVM-inspired markdown compiler for mortals. Plain text in, anything out. Runs everywhere — terminal, web, native apps, any device. Written in Go.'},
      {property: 'og:type', content: 'website'},
      {property: 'og:site_name', content: 'LazyMD'},
      {property: 'og:title', content: 'LazyMD — The Compiler for Mortals'},
      {property: 'og:description', content: 'An LLVM-inspired markdown compiler for mortals. Plain text in, anything out. Runs everywhere — terminal, web, native apps, any device.'},
      {property: 'og:url', content: 'https://lazymd.com'},
      {name: 'twitter:card', content: 'summary_large_image'},
      {name: 'twitter:title', content: 'LazyMD — The Compiler for Mortals'},
      {name: 'twitter:description', content: 'An LLVM-inspired markdown compiler for mortals. Plain text in, anything out. Runs everywhere. Written in Go.'},
      {name: 'robots', content: 'index, follow'},
      {name: 'theme-color', content: '#000000'},
    ],
    colorMode: {
      defaultMode: 'dark',
      disableSwitch: false,
      respectPrefersColorScheme: false,
    },
    navbar: {
      title: 'LazyMD',
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Docs',
        },
        {to: '/blog', label: 'Blog', position: 'left'},
        {
          href: 'https://github.com/EME130/lazymd',
          position: 'right',
          className: 'header-github-link',
          'aria-label': 'GitHub repository',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {label: 'Getting Started', to: '/docs/getting-started/installation'},
            {label: 'Usage Guide', to: '/docs/usage/editing-modes'},
            {label: 'Plugin System', to: '/docs/plugins/overview'},
            {label: 'MCP Server', to: '/docs/mcp-server/overview'},
          ],
        },
        {
          title: 'Community',
          items: [
            {label: 'GitHub Discussions', href: 'https://github.com/EME130/lazymd/discussions'},
            {label: 'Issues', href: 'https://github.com/EME130/lazymd/issues'},
          ],
        },
        {
          title: 'More',
          items: [
            {label: 'Blog', to: '/blog'},
            {label: 'GitHub', href: 'https://github.com/EME130/lazymd'},
            {label: 'Releases', href: 'https://github.com/EME130/lazymd/releases'},
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} LazyMD. Open source under the MIT License.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['bash', 'go', 'json', 'yaml'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
