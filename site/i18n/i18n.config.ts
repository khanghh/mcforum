import en from './en-US'
import vi from './vi-VN'

export default defineI18nConfig(() => ({
  legacy: false,
  locale: 'vi',
  messages: { en, vi },
}))
