import en from './locales/en-US'
import vi from './locales/vi-VN'

export default defineI18nConfig(() => ({
  legacy: false,
  locale: 'vi',
  defaultLocale: 'vi',
  messages: { en, vi },
}))
