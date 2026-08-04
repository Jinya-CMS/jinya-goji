import type { SidebarsConfig } from '@docusaurus/plugin-content-docs';
import apiSidebar from './docs/api/sidebar';

const sidebars: SidebarsConfig = {
  developers: [
    {
      type: 'category',
      label: 'API Documentation',
      items: apiSidebar,
      collapsed: true,
    },
  ],
  users: [
    {
      type: 'category',
      label: 'Getting Started',
      link: {
        type: 'doc',
        id: 'getting-started',
      },
      items: ['getting-started/docker'],
    },
  ],
};

export default sidebars;
