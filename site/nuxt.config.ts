// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  // ssr: false,

  modules: [
    '@nuxtjs/tailwindcss',
    '@vueuse/nuxt',
    '@pinia/nuxt',
    '@nuxt/icon',
    '@nuxtjs/i18n',
    '@nuxt/eslint',
    '@nuxtjs/seo',
    '@nuxtjs/sitemap',
  ],

  plugins: [
  ],

  imports: {
    autoImport: true,
    dirs: [
      'apis',
      'stores',
    ],
  },
  devtools: { enabled: true },

  app: {
    head: {
      title: 'Minecraft Việt Nam - Diễn đàn MineViet Network',
      link: [
        { rel: 'icon', href: '/favicon.ico' },
        { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png' },
        { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon-16x16.png' },
        { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' },
      ],
      htmlAttrs: {
        lang: 'vi',
      },
      script: [],
    },
  },

  css: [
    '~/assets/css/index.scss',
  ],

  runtimeConfig: {
    public: {
      loginUrl: `https://auth.mineviet.com/login?service=forum-test`,
    },
  },
  devServer: {
    port: 3000, // default: 3000
    host: '0.0.0.0', // default: localhost
  },
  compatibilityDate: '2024-09-15',

  eslint: {
    config: {
      stylistic: true,
    },
  },
  unhead: {
    legacy: true,
  },
  i18n: {
    locales: [
      { code: 'vi', file: 'vi-VN.ts', iso: 'vi-VN', name: 'Tiếng Việt' },
    ],
    defaultLocale: 'vi',
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'i18n_redirected',
      redirectOn: 'root',
    },
  },

  icon: {
    serverBundle: {
      collections: [],
    },
  },

  tailwindcss: {
    exposeConfig: true,
    viewer: true,
  },
  site: {
    name: 'MineViet Network',
    url: 'https://mineviet.com',
  }
})