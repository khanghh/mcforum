// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  // ssr: false,
  srcDir: 'src/',
  compatibilityDate: '2024-09-15',
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
    // '@nuxtjs/color-mode',
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
    }
  },

  css: [
    '~/assets/css/index.scss',
  ],

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
        proxy: `http://localhost:3001/api/**`,
      },
      '/admin/**': {
        proxy: `${import.meta.env.SERVER_URL}/admin/**`,
      },
    },
  },

  runtimeConfig: {
    public: {
      // BASE_URL: 'http://localhost:3001',
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
        'circle-notch',
        'clock',
        'cloud-upload',
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
        'image',
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
        'sign-out-alt',
        'spinner',
        'star',
        'th-list',
        'thumbs-up',
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
