// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({

  // ssr: false,
  modules: [
    '@pinia/nuxt',
    '@vueuse/nuxt',
    // https://color-mode.nuxtjs.org/#configuration
    '@nuxtjs/color-mode',
    '@nuxtjs/i18n',
    '@element-plus/nuxt',
    ['nuxt-lazy-load', {
      images: true,
      videos: true,
      audios: true,
      iframes: true,
      native: true,
      directiveOnly: false,

      // Default image must be in the public folder
      // defaultImage: '/images/default-image.jpg',

      // To remove class set value to false
      loadingClass: 'isLoading',
      loadedClass: 'isLoaded',
      appendClass: 'lazyLoad',

      observerConfig: {
        // See IntersectionObserver documentation
      },
    }],
    '@nuxt/eslint',
  ],

  plugins: [
  ],

  imports: {
    dirs: [
      'apis',
      'stores',
    ],
  },
  devtools: { enabled: true },

  app: {
    head: {
      title: 'BBS-GO',
      htmlAttrs: { class: 'theme-light has-navbar-fixed-top' },
      script: [
        {
          src: 'https://hm.baidu.com/hm.js?f14b836e09b72aedce29a86e809936de',
          type: 'text/javascript',
          async: true,
        },
      ],
    },
  },

  css: [
    '~/assets/css/index.scss',
  ],

  colorMode: {
    preference: 'system', // default value of $colorMode.preference
    fallback: 'light', // fallback value if not system preference found
    storageKey: 'bbsgo-color-mode',
    classPrefix: 'theme-',
    classSuffix: '',
  },
  srcDir: 'src/',

  compatibilityDate: '2024-09-15',

  nitro: {
    routeRules: {
      '/api/**': {
        proxy: `${import.meta.env.SERVER_URL}/api/**`,
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
})
