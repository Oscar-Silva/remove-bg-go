/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: '#0A0A0F',
        surface: '#111118',
        surfaceElevated: '#1C1C27',
        border: '#2A2A3D',
        accent: {
          DEFAULT: '#7C6EFA',
          hover: '#9D8FFF',
        },
        success: '#34D399',
        error: '#F87171',
        textPrimary: '#F0F0FF',
        textMuted: '#8B8BA7',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      borderRadius: {
        DEFAULT: '12px',
        lg: '16px',
        full: '999px',
      },
      boxShadow: {
        'glow': '0 0 20px rgba(124, 110, 250, 0.25)',
        'glow-lg': '0 0 40px rgba(124, 110, 250, 0.35)',
      },
      animation: {
        'float': 'float 3s ease-in-out infinite',
        'shimmer': 'shimmer 2s linear infinite',
        'progress': 'progress 1.5s ease-in-out infinite',
      },
      keyframes: {
        float: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-8px)' },
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' },
        },
        progress: {
          '0%': { backgroundPosition: '0% 50%' },
          '100%': { backgroundPosition: '100% 50%' },
        },
      },
    },
  },
  plugins: [],
}
