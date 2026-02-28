/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: [
    './internal/renderer/templates/**/*.html',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['"Google Sans"', 'ui-sans-serif', 'system-ui'],
      },
      colors: {
        blue: { DEFAULT: '#1a6eb5', light: '#3b9eff', bg: '#ddeeff' },
        navy: '#0d2a4a',
      },
    },
  },
  plugins: [],
}
