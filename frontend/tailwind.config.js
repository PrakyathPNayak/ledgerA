/** @type {import('tailwindcss').Config} */
export default {
  darkMode: ['class'],
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        app: 'var(--color-bg-app)',
        surface: 'var(--color-bg-surface)',
        elevated: 'var(--color-bg-elevated)',
        'surface-hover': 'var(--color-bg-hover)',
        border: 'var(--color-border)',
        'border-subtle': 'var(--color-border-subtle)',
        foreground: 'var(--color-text)',
        secondary: 'var(--color-text-secondary)',
        muted: 'var(--color-text-muted)',
        accent: 'var(--color-accent)',
        positive: 'var(--color-positive)',
        'positive-muted': 'var(--color-positive-muted)',
        negative: 'var(--color-negative)',
        'negative-muted': 'var(--color-negative-muted)',
      },
    },
  },
  plugins: [],
}

