module.exports = {
  purge: ["./client/**/*.tsx", "./client/**/*.html"],
  darkMode: 'media', // or 'media' or 'class'
  theme: {
    extend: {
      "colors": theme => ({
        "nord-1": "#242932",
        "nord-2": "#20242c",
        "nord-3": "#1a1d23",
        "nord-4": "#16181d",
        "nord-5": "#0b0c0f",
      })
    },
    fontFamily: {
      'quicksand': ['Quicksand', 'sans-serif']
    }
  },
  variants: {
    extend: {},
  },
  plugins: [
    require('tailwind-nord'),
  ],
}
