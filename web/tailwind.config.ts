import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        surface: '#1a1a1a',
        border: '#2a2a2a',
        accent: '#FF9900',
        muted: '#666666',
      },
    },
  },
  plugins: [],
}

export default config
