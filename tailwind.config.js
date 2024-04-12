/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: ["./**/*.{html,js,templ}", "./node_modules/flowbite/**/*.js"],
  theme: {
    extend: {
    },
  },
  plugins: [
    require('flowbite/plugin')
  ],
}

