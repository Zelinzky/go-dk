/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./**/*.go'],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('@tailwindcss/forms'),
  ],
}

