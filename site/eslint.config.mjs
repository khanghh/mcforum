// @ts-check
// import stylistic from '@stylistic/eslint-plugin'
import withNuxt from './.nuxt/eslint.config.mjs'

export default withNuxt().overrideRules({
  '@typescript-eslint/no-unused-vars': 'off',
  'vue/html-closing-bracket-newline': 'off',
  'vue/max-attributes-per-line': 'off',
  'vue/first-attribute-linebreak': 'off',
  'vue/multi-word-component-names': 'off',
  'vue/html-self-closing': 'off',
  'vue/v-on-event-hyphenation': 'off',
  'vue/attribute-hyphenation': 'off',
  'vue/singleline-html-element-content-newline': 'off',
  'vue/html-indent': ['error', 2, { alignAttributesVertically: false }],
  'vue/multiline-html-element-content-newline': ['error', {
    ignoreWhenEmpty: false,
    ignores: ['pre', 'textarea'],
    allowEmptyLines: false,
  }],
  '@stylistic/brace-style': ['error', '1tbs', { allowSingleLine: true }],
})
