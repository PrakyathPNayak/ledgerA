import type { Config } from 'tailwindcss'

const config: Config = {
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
            keyframes: {
                'accordion-down': {
                    from: { height: '0' },
                    to: { height: 'var(--radix-accordion-content-height)' },
                },
                'accordion-up': {
                    from: { height: 'var(--radix-accordion-content-height)' },
                    to: { height: '0' },
                },
            },
            animation: {
                'accordion-down': 'accordion-down 0.2s ease-out',
                'accordion-up': 'accordion-up 0.2s ease-out',
            },
        },
    },
    plugins: [],
}

export default config
