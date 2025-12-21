// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  // ssr: false,
  modules: [
    '@vueuse/nuxt', // https://color-mode.nuxtjs.org/#configuration
    '@pinia/nuxt',
    '@nuxtjs/tailwindcss',
    '@vesp/nuxt-fontawesome',
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
    '@nuxtjs/color-mode',
    '@nuxtjs/i18n',
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
      script: [],
      link: [
        { rel: 'stylesheet', href: 'https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css' },
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
        proxy: `http://localhost:3001/api/**`,
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
  fontawesome: {
    icons: {
      solid: [
        'arrow-up',
        'ban',
        'bookmark',
        'calendar-plus',
        'chart-line',
        'check-circle',
        'chevron-down',
        'clock',
        'cog',
        'comment',
        'comment-dots',
        'comments',
        'crown',
        'cube',
        'envelope',
        'exclamation-triangle',
        'eye',
        'filter',
        'fire',
        'flag',
        'gamepad',
        'heart',
        'home',
        'info-circle',
        'map-marker-alt',
        'medal',
        'palette',
        'play',
        'plus',
        'puzzle-piece',
        'question-circle',
        'quote-left',
        'server',
        'share',
        'shield-alt',
        'shopping-cart',
        'spinner',
        'star',
        'th-list',
        'thumbtack',
        'times',
        'trophy',
        'user',
        'user-circle',
        'user-plus',
        'users',
      ],
      regular: [
        'bookmark',
        'comment',
        'comments',
        'eye',
      ],
      brands: [
        'discord',
        'twitch',
        'twitter',
        'youtube',
      ],
    },
  },

  i18n: {
    vueI18n: './i18n.config.ts',
  },

  tailwindcss: {
    exposeConfig: true,
    viewer: true,
  },
})
