import en from './en-US'
import vi from './vi-VN'

export default defineI18nConfig(() => ({
  legacy: false,
  locale: 'en',
  messages: { en, vi },
}));
