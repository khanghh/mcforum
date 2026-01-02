// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  // ssr: false,
  modules: [// https://color-mode.nuxtjs.org/#configuration
    '@vueuse/nuxt',
    '@pinia/nuxt',
    '@nuxtjs/tailwindcss',
    '@nuxt/icon',
    '@nuxtjs/i18n',
    '@nuxt/eslint',
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
      link: [
        { rel: 'icon', href: '/favicon.ico' },
        { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png' },
        { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon-16x16.png' },
        { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' },
      ],
      htmlAttrs: { class: 'theme-light has-navbar-fixed-top' },
      script: [],
    },
  },

  css: [
    '~/assets/css/index.scss',
  ],

  runtimeConfig: {
    public: {
      loginUrl: 'https://auth.mineviet.com/login?service=test',
      // BASE_URL: 'http://localhost:3001',
    },
  },
  srcDir: 'src/',
  devServer: {
    port: 3000, // default: 3000
    host: '0.0.0.0', // default: localhost
  },
  compatibilityDate: '2024-09-15',

  // colorMode: {
  //   preference: 'system', // default value of $colorMode.preference
  //   fallback: 'light', // fallback value if not system preference found
  //   storageKey: 'bbsgo-color-mode',
  //   classPrefix: 'theme-',
  //   classSuffix: '',
  // },

  nitro: {
    routeRules: {
      '/api/**': {
        proxy: `http://172.17.0.1:3001/api/**`,
      },
      '/upload/**': {
        proxy: `http://172.17.0.1:3001/upload/**`,
      },
      '/admin/**': {
        proxy: `${import.meta.env.SERVER_URL}/admin/**`,
      },
    },
  },

  eslint: {
    config: {
      stylistic: true,
    },
  },

  i18n: {
    vueI18n: './i18n.config.ts',
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

})
