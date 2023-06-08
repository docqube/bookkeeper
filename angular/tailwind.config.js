/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{html,scss,ts}",
  ],
  theme: {
    borderRadius: {
      'md': '0.375rem',
      DEFAULT: '0.75rem',
    },
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}

