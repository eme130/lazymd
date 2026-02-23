import React from 'react';
import ComponentCreator from '@docusaurus/ComponentCreator';

export default [
  {
    path: '/__docusaurus/debug',
    component: ComponentCreator('/__docusaurus/debug', '5ff'),
    exact: true
  },
  {
    path: '/__docusaurus/debug/config',
    component: ComponentCreator('/__docusaurus/debug/config', '5ba'),
    exact: true
  },
  {
    path: '/__docusaurus/debug/content',
    component: ComponentCreator('/__docusaurus/debug/content', 'a2b'),
    exact: true
  },
  {
    path: '/__docusaurus/debug/globalData',
    component: ComponentCreator('/__docusaurus/debug/globalData', 'c3c'),
    exact: true
  },
  {
    path: '/__docusaurus/debug/metadata',
    component: ComponentCreator('/__docusaurus/debug/metadata', '156'),
    exact: true
  },
  {
    path: '/__docusaurus/debug/registry',
    component: ComponentCreator('/__docusaurus/debug/registry', '88c'),
    exact: true
  },
  {
    path: '/__docusaurus/debug/routes',
    component: ComponentCreator('/__docusaurus/debug/routes', '000'),
    exact: true
  },
  {
    path: '/blog',
    component: ComponentCreator('/blog', 'df4'),
    exact: true
  },
  {
    path: '/blog/2025/01/01/welcome',
    component: ComponentCreator('/blog/2025/01/01/welcome', '413'),
    exact: true
  },
  {
    path: '/blog/archive',
    component: ComponentCreator('/blog/archive', '182'),
    exact: true
  },
  {
    path: '/blog/authors',
    component: ComponentCreator('/blog/authors', '0b7'),
    exact: true
  },
  {
    path: '/blog/tags',
    component: ComponentCreator('/blog/tags', '287'),
    exact: true
  },
  {
    path: '/blog/tags/announcement',
    component: ComponentCreator('/blog/tags/announcement', '4fa'),
    exact: true
  },
  {
    path: '/blog/tags/markdown',
    component: ComponentCreator('/blog/tags/markdown', '4b4'),
    exact: true
  },
  {
    path: '/blog/tags/terminal-editor',
    component: ComponentCreator('/blog/tags/terminal-editor', '36a'),
    exact: true
  },
  {
    path: '/blog/tags/vim',
    component: ComponentCreator('/blog/tags/vim', '95d'),
    exact: true
  },
  {
    path: '/blog/tags/zig',
    component: ComponentCreator('/blog/tags/zig', '773'),
    exact: true
  },
  {
    path: '/docs',
    component: ComponentCreator('/docs', '104'),
    routes: [
      {
        path: '/docs',
        component: ComponentCreator('/docs', '7e6'),
        routes: [
          {
            path: '/docs',
            component: ComponentCreator('/docs', 'bd2'),
            routes: [
              {
                path: '/docs/architecture/module-reference',
                component: ComponentCreator('/docs/architecture/module-reference', '980'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/architecture/system-design',
                component: ComponentCreator('/docs/architecture/system-design', 'e0b'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/configuration/file-formats',
                component: ComponentCreator('/docs/configuration/file-formats', '8be'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/configuration/themes',
                component: ComponentCreator('/docs/configuration/themes', '8f6'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/contributing/building',
                component: ComponentCreator('/docs/contributing/building', '548'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/contributing/how-to-contribute',
                component: ComponentCreator('/docs/contributing/how-to-contribute', '5cf'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/contributing/testing',
                component: ComponentCreator('/docs/contributing/testing', '128'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/getting-started/first-file',
                component: ComponentCreator('/docs/getting-started/first-file', '83e'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/getting-started/installation',
                component: ComponentCreator('/docs/getting-started/installation', 'f1f'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/getting-started/quick-start',
                component: ComponentCreator('/docs/getting-started/quick-start', '835'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/mcp-server/overview',
                component: ComponentCreator('/docs/mcp-server/overview', 'b38'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/plugins/api-reference',
                component: ComponentCreator('/docs/plugins/api-reference', '278'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/plugins/catalog',
                component: ComponentCreator('/docs/plugins/catalog', 'bb4'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/plugins/development',
                component: ComponentCreator('/docs/plugins/development', '24c'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/plugins/overview',
                component: ComponentCreator('/docs/plugins/overview', '1ca'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/usage/brain-graph',
                component: ComponentCreator('/docs/usage/brain-graph', 'a67'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/usage/commands',
                component: ComponentCreator('/docs/usage/commands', 'aab'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/usage/editing',
                component: ComponentCreator('/docs/usage/editing', '295'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/usage/editing-modes',
                component: ComponentCreator('/docs/usage/editing-modes', 'fb4'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/usage/mouse-support',
                component: ComponentCreator('/docs/usage/mouse-support', '2ba'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/usage/navigation',
                component: ComponentCreator('/docs/usage/navigation', 'bdb'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/usage/panels-layout',
                component: ComponentCreator('/docs/usage/panels-layout', 'ff9'),
                exact: true,
                sidebar: "docsSidebar"
              }
            ]
          }
        ]
      }
    ]
  },
  {
    path: '/',
    component: ComponentCreator('/', 'e5f'),
    exact: true
  },
  {
    path: '*',
    component: ComponentCreator('*'),
  },
];
