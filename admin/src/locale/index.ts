import { createI18n } from 'vue-i18n';
import en from './en-US';

export const LOCALE_OPTIONS = [
  { label: 'English', value: 'en-US' },
];
const defaultLocale = localStorage.getItem('arco-locale') || 'en-US';

const i18n = createI18n({
  locale: defaultLocale,
  fallbackLocale: 'en-US',
  legacy: false,
  allowComposition: true,
  messages: {
    'en-US': en,
  },
});

export default i18n;
