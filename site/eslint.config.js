import { createConfigForNuxt } from '@nuxt/eslint-config/flat'

export default createConfigForNuxt({
  features: {
    stylistic: {
      semi: false,
      indent: 2,
      quotes: 'single',
    },
  },
}).override('nuxt/vue/rules', {
  rules: {
    '@typescript-eslint/no-unused-vars': 'off',
    'vue/html-closing-bracket-newline': 'off',
    'vue/max-attributes-per-line': 'off',
    'vue/first-attribute-linebreak': 'off',
    'vue/multi-word-component-names': 'off',
    'vue/html-self-closing': 'off',
    'vue/v-on-event-hyphenation': 'off',
    'vue/html-indent': ['error', 2, { alignAttributesVertically: false }],
    'vue/multiline-html-element-content-newline': ['error', {
      ignoreWhenEmpty: false,
      ignores: ['pre', 'textarea'],
      allowEmptyLines: false,
    }],
  },
})
