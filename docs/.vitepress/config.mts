import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Proxctl',
  description: 'Lightweight proxy deployment and operations CLI',
  lang: 'en-US',
  cleanUrls: true,
  lastUpdated: true,
  sitemap: {
    hostname: 'https://proxctl.kaixuan.ai'
  },
  head: [
    ['link', { rel: 'icon', href: '/favicon.svg' }],
    ['meta', { name: 'theme-color', content: '#0f766e' }]
  ],
  themeConfig: {
    logo: '/favicon.svg',
    siteTitle: 'Proxctl',
    search: {
      provider: 'local'
    },
    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Operations', link: '/operations/backup-restore' },
      { text: 'Reference', link: '/reference/configuration' },
      { text: 'Development', link: '/development/architecture' }
    ],
    sidebar: [
      {
        text: 'Guide',
        items: [
          { text: 'Getting Started', link: '/guide/getting-started' },
          { text: 'Installation', link: '/guide/installation' },
          { text: 'Xray Provider', link: '/guide/xray' },
          { text: 'Client Export', link: '/guide/client-export' }
        ]
      },
      {
        text: 'Operations',
        items: [
          { text: 'Backup and Restore', link: '/operations/backup-restore' },
          { text: 'Credential Rotation', link: '/operations/credential-rotation' },
          { text: 'SSH Hardening', link: '/operations/ssh-hardening' },
          { text: 'Firewall', link: '/operations/firewall' },
          { text: 'Boot Check', link: '/operations/boot-check' }
        ]
      },
      {
        text: 'Reference',
        items: [
          { text: 'Configuration', link: '/reference/configuration' },
          { text: 'Artifacts', link: '/reference/artifacts' },
          { text: 'Release Process', link: '/reference/release' },
          { text: 'Command Reference', link: '/commands/proxctl' }
        ]
      },
      {
        text: 'Development',
        items: [
          { text: 'Architecture', link: '/development/architecture' },
          { text: 'Development Guide', link: '/development/development' }
        ]
      }
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/kehr/Proxctl' }
    ],
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2026 Proxctl contributors'
    }
  }
})
