import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  docs: [
    'intro',
    'installation',
    'authentication',
    'configuration',
    'output-and-errors',
    'pagination',
    {
      type: 'category',
      label: 'Command reference',
      link: {type: 'doc', id: 'commands/index'},
      items: ['commands/auth', 'commands/config', 'commands/teams', 'commands/users', 'commands/issues', 'commands/comments', 'commands/labels', 'commands/workflow-states', 'commands/projects', 'commands/cycles', 'commands/documents', 'commands/attachments'],
    },
    {
      type: 'category',
      label: 'Automation',
      items: ['automation/ai-agents', 'automation/shell-scripting'],
    },
    {
      type: 'category',
      label: 'Development',
      items: ['development/contributing', 'development/releases'],
    },
    'reference/limitations',
  ],
};

export default sidebars;
