/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/components/**/*.{js,vue,ts}',
    './src/layouts/**/*.vue',
    './src/pages/**/*.vue',
    './src/plugins/**/*.{js,ts}',
    './src/app.vue',
    './src/error.vue',
  ],
  theme: {
    extend: {
      fontFamily: {
        // sans: ["ui-sans-serif", "system-ui", "sans-serif", "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji"],
        sans: ['Inter', 'sans-serif'],
        display: ['Rajdhani', 'sans-serif'],
      },
      colors: {
        gaming: {
          bg: '#0f0f23',
          card: 'rgba(30, 30, 60, 0.8)',
        },
      },
    },
  },
  plugins: [],
}
