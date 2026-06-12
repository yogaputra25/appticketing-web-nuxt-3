import type { Config } from 'tailwindcss'

export default <Partial<Config>>{
  content: [
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './composables/**/*.{js,ts}',
    './plugins/**/*.{js,ts}',
    './app.vue',
    './error.vue',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f0f4ff',
          100: '#dbe5ff',
          200: '#bccfff',
          300: '#8eaeff',
          400: '#5a82ff',
          500: '#3358ff',
          600: '#1f3ce0',
          700: '#1a2fae',
          800: '#1a2a87',
          900: '#1c296c',
        },
        accent: {
          DEFAULT: '#ff3d6e',
          dark: '#cc264f',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      animation: {
        'pulse-slow': 'pulse 2.5s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
    },
  },
  plugins: [],
}
